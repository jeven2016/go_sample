package yzs8

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-creed/sat"
	"github.com/go-redis/redis/v9"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"colly-downloader/pkg/models"
)

var articlesChan = make(chan *models.ArticlePage, 100)

var wg = sync.WaitGroup{}

var baseUrl string

var log = initLog()

var waitIndex = atomic.NewInt32(0)

func Start() {
	// go func() {
	// 	wg.Wait()
	// }()

	// Traditional Chinese to Simplified Chinese
	zhConvertor := sat.DefaultDict()

	// mongo
	client, err := CreateMongoClient(log)
	handleError(err)
	// defer client.Disconnect(context.Background())

	// redis
	redis, err := RedisClient(log)
	handleError(err)
	// defer redis.Close()

	// collection
	// 初始化全局Db
	db := client.Database("books")
	collection := db.Collection("catalog")

	catalogId, existed := ensure(collection)
	if !existed {
		return
	}

	collector := newCollector()
	homeUrl := "--"
	baseUrl = getBaseUri(homeUrl)

	go parsePageLinks(homeUrl, collector, articlesChan, zhConvertor, redis)

	articleCollection := db.Collection("article")
	for i := 0; i < 3; i++ {
		go downloadArticle(articleCollection, articlesChan, catalogId, collector.Clone(), zhConvertor, redis)
	}
}

func ensure(catalogCol *mongo.Collection) (*primitive.ObjectID, bool) {
	c := &models.CatalogDoc{
		Name:         "--",
		Order:        1,
		ArticleCount: 0,
		Description:  "",
		CreateDate:   time.Now(),
		LastUpdate:   time.Now(),
	}
	id, succeed := ensureCollection(catalogCol, c)
	if succeed {
		c := &models.CatalogDoc{
			Name:         "--",
			ParentId:     *id,
			Order:        11,
			ArticleCount: 0,
			Description:  "--",
			CreateDate:   time.Now(),
			LastUpdate:   time.Now(),
		}
		id, succeed = ensureCollection(catalogCol, c)
		return id, succeed
	}
	return id, succeed
}

// ensure valid catalog id to return
func ensureCollection(catalogCol *mongo.Collection, catalog *models.CatalogDoc) (*primitive.ObjectID, bool) {
	result := catalogCol.FindOne(context.TODO(), bson.M{"name": catalog.Name})
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			insertResult, err := catalogCol.InsertOne(context.TODO(), catalog)
			if err != nil {
				log.Warn("Failed to insert on catalog", zap.Error(err))
			}
			id := insertResult.InsertedID.(primitive.ObjectID)
			return &id, err == nil
		}
		return nil, false
	}
	catalog = new(models.CatalogDoc)
	err = result.Decode(catalog)
	if err != nil {
		return nil, false
	}
	catalogId := catalog.Id
	log.Info("catalog id is", zap.String("id", catalogId.Hex()))
	return &catalogId, true
}

// New collector
func newCollector() *colly.Collector {
	// c := colly.NewCollector(colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"))
	// c := colly.NewCollector(colly.CacheDir("./temp"))
	c := colly.NewCollector()
	c.SetRequestTimeout(50 * time.Second)
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true, // Colly uses HTTP keep-alive to enhance scraping speed
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("[Request URL]:", r.StatusCode, " ", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("[Visiting]", r.URL.String())
	})

	// 随机设置
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	return c
}

func getBaseUri(url string) string {
	reg, err := regexp.Compile("(https?://.*?/)")
	if err != nil {
		// print log
		return ""
	}
	subs := reg.FindStringSubmatch(url)

	if len(subs) > 1 {
		return subs[1]
	}
	return ""
}

var at = atomic.NewInt32(1)

func parsePageLinks(homeUrl string, collector *colly.Collector, urlChan chan<- *models.ArticlePage, zhConvertor sat.Dicter, redis *redis.Client) {
	ctx := colly.NewContext()
	// key: page url
	// cacheKey := base64.StdEncoding.EncodeToString([]byte(homeUrl))
	collector.OnHTML("body", func(element *colly.HTMLElement) {
		// 以先后循序解析同一个页面内的内容
		element.ForEach(".articleList>p>.img-center", func(i int, element *colly.HTMLElement) {
			var pageUrl = element.Request.Ctx.GetAny("pageUrl")
			var artUrl = baseUrl + element.Attr("href")
			var artName = strings.TrimSpace(zhConvertor.Read(element.Text))

			if number, err := redis.Exists(context.Background(), base64Key(artUrl)).Result(); err == nil && number > 0 {
				log.Info("redis check: document exists, ignored", zap.String("name", artName), zap.String("pageUrl", pageUrl.(string)))
				return
			}

			var inlegalName = []string{"母亲", "妈妈", "母子", "金鳞岂是", "阿宾"}
			for _, item := range inlegalName {
				if strings.Contains(artName, item) {
					log.Warn("illegal document, ignored", zap.String("name", artName), zap.String("pageUrl", pageUrl.(string)))
					return
				}
			}

			artPage := &models.ArticlePage{
				Name:    artName,
				Url:     artUrl,
				PageUrl: element.Request.Ctx.Get("pageUrl"),
			}
			urlChan <- artPage
		})

		element.ForEach(".pager a[href]", func(i int, element *colly.HTMLElement) {
			text := element.Text
			if text == "下一页" {
				href := element.Attr("href")
				nextPageUrl := baseUrl + href

				if at.Load()%20 == 0 {
					time.Sleep(20 * time.Second)
				}

				at.Add(1)

				ctx.Put("pageUrl", nextPageUrl)
				println("current channel size ", len(urlChan))

				// cacheKey = base64.StdEncoding.EncodeToString([]byte(nextPageUrl))

				if err := collector.Request("GET", nextPageUrl, nil, ctx, nil); err != nil {
					handleError(err)
				}

			}
		})
	})
	ctx.Put("pageUrl", homeUrl)
	handleError(collector.Request("GET", homeUrl, nil, ctx, nil))
}

func downloadArticle(collection *mongo.Collection, urlChan <-chan *models.ArticlePage, catalogId *primitive.ObjectID,
	c *colly.Collector, zhConvertor sat.Dicter, client *redis.Client) {
	// load article page and get the content to save
	// .articleList>.content>div
	c.OnHTML(".articleList>.content>div:first-child", func(element *colly.HTMLElement) {
		// if waitIndex.Load()%20 == 0 {
		// 	time.Sleep(3 * time.Second)
		// }
		// waitIndex.Inc()

		artPage := element.Request.Ctx.GetAny("articlePage").(*models.ArticlePage)
		content, err := element.DOM.Html()
		if err != nil {
			log.Warn("failed to get the content",
				zap.Error(err), zap.String("url", element.Request.URL.String()))
			return
		}

		if len(strings.TrimSpace(content)) == 0 {
			log.Info("Content is blank", zap.String("name", artPage.Name), zap.String("url", artPage.Url))
			return
		}

		content = zhConvertor.Read(content)
		index := strings.Index(content, "==记住==")
		if index > -1 {
			content = content[:index]
		}

		_, err = collection.InsertOne(context.TODO(), models.Article{
			Name:       artPage.Name,
			CatalogId:  *catalogId,
			Content:    content,
			CreateDate: time.Now(),
		})
		if err != nil {
			log.Error("failed to insert document with name",
				zap.Error(err), zap.String("name", artPage.Name), zap.String("url", artPage.Url))
			return
		}
		log.Info("Inserted a document", zap.String("name", artPage.Name))
		client.Set(context.Background(), base64Key(artPage.Url), artPage, 20*time.Minute)
	})

	ctx := colly.NewContext()
	for artPage := range urlChan {
		realName := artPage.Name
		if val, err := client.Exists(context.Background(), base64Key(artPage.Url)).Result(); err == nil && val > 0 {
			log.Info("redis check: document exists, ignored", zap.String("name", realName), zap.String("pageUrl", artPage.PageUrl))
			continue
		}

		// 过滤掉重复的article
		count, err := collection.CountDocuments(context.TODO(), bson.M{"name": bson.M{"$regex": realName}})
		if err != nil {
			log.Warn("failed to count documents with name",
				zap.Error(err), zap.String("name", realName))
			continue
		}
		if count > 0 {
			log.Info("mongo check: document exists, ignored", zap.String("name", realName), zap.String("pageUrl", artPage.PageUrl))
			client.Set(context.Background(), base64Key(artPage.Url), artPage, 20*time.Minute)
			continue
		}
		// 加载文章
		// 为了使用colly.Context向onHTML中传递参数，使用Request替代Visit
		ctx.Put("articlePage", artPage)
		if err := c.Request("GET", artPage.Url, nil, ctx, nil); err != nil {
			handleError(err)
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Error("error occurs", zap.Error(err))
	}
}

func base64Key(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}
