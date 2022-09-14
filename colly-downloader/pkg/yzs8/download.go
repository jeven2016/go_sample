package yzs8

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-creed/sat"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"colly-downloader/pkg/models"
)

var articlesChan = make(chan *models.ArticlePage, 1000)

var wg = sync.WaitGroup{}

var baseUrl string

var log = initLog()

func Start() {
	go func() {
		wg.Wait()
	}()

	// Traditional Chinese to Simplified Chinese
	zhConvertor := sat.DefaultDict()

	db, err := CreateMongoClient(log)
	if err != nil {
		return
	}
	collection := db.Collection("catalog")

	catalogId, existed := ensure(collection)
	if !existed {
		return
	}

	collector := newCollector()
	homeUrl := "https://yazhouse8.com/article.php?cate=1"
	baseUrl = getBaseUri(homeUrl)
	parsePageLinks(homeUrl, collector, articlesChan, zhConvertor)

	for i := 0; i < 1; i++ {
		downloadArticle(i, collection, articlesChan, zhConvertor, catalogId, collector.Clone())
	}
}

func ensure(catalogCol *mongo.Collection) (*primitive.ObjectID, bool) {
	c := &models.CatalogDoc{
		Name:         "亚洲色吧",
		Order:        1,
		ArticleCount: 0,
		Description:  "",
		CreateDate:   time.Now(),
		LastUpdate:   time.Now(),
	}
	id, succeed := ensureCollection(catalogCol, c)
	if succeed {
		c := &models.CatalogDoc{
			Name:         "都市激情",
			ParentId:     *id,
			Order:        1,
			ArticleCount: 0,
			Description:  "都市激情，激情小说合集",
			CreateDate:   time.Now(),
			LastUpdate:   time.Now(),
		}
		id, succeed = ensureCollection(catalogCol, c)
		return id, succeed
	}
	return id, succeed
}

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
	log.Info("catalog id is", zap.String("id", catalogId.String()))
	return &catalogId, true
}

// New collector
func newCollector() *colly.Collector {
	// c := colly.NewCollector(colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"))
	c := colly.NewCollector()
	c.SetRequestTimeout(20 * time.Second)
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true, // Colly uses HTTP keep-alive to enhance scraping speed
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("[Request URL]:", r.Request.URL, "failed with response:", r, "\nError:", err)
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

var at = atomic.NewInt32(0)

func parsePageLinks(homeUrl string, collector *colly.Collector, urlChan chan<- *models.ArticlePage, zhConvertor sat.Dicter) {
	collector.OnHTML(".articleList>p>.img-center", func(element *colly.HTMLElement) {
		urlChan <- &models.ArticlePage{
			Name: strings.TrimSpace(zhConvertor.Read(element.Text)),
			Url:  baseUrl + element.Attr("href"),
		}
	})

	collector.OnHTML(".pager a[href]", func(element *colly.HTMLElement) {
		if at.Load() >= 1 {
			return
		}
		text := element.Text
		if text == "下一页" {
			href := element.Attr("href")
			nextPageUrl := baseUrl + href
			at.Add(1)

			// 放在最后一个OnHTML中执行
			handleError(element.Request.Visit(nextPageUrl))
		}
	})
	handleError(collector.Visit(homeUrl))
}

func downloadArticle(taskId int, collection *mongo.Collection, urlChan <-chan *models.ArticlePage,
	zhConvertor sat.Dicter, catalogId *primitive.ObjectID, c *colly.Collector) {
	// load article page and get the content to save
	// .articleList>.content>div
	c.OnHTML(".articleList>.content>div", func(element *colly.HTMLElement) {
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
		content = strings.TrimRight(content, "==记住==")

		_, err = collection.InsertOne(context.TODO(), models.Article{
			Name:      artPage.Name,
			CatalogId: *catalogId,
			Content:   content,
		})
		if err != nil {
			log.Error("failed to insert document with name",
				zap.Error(err), zap.String("name", artPage.Name))
			return
		}
		log.Info("Inserted a document", zap.String("name", artPage.Name))
	})

	ctx := colly.NewContext()
	for artPage := range urlChan {
		realName := artPage.Name

		// 过滤掉重复的article
		count, err := collection.CountDocuments(context.TODO(), bson.M{"name": bson.M{"$regex": realName}})
		if err != nil {
			log.Warn("failed to count documents with name",
				zap.Error(err), zap.String("name", realName))
			continue
		}
		if count > 0 {
			log.Info("document exists, ignored", zap.String("name", realName))
			continue
		}
		// 加载文章
		// 为了使用colly.Context向onHTML中传递参数，使用Request替代Visit
		ctx.Put("articlePage", artPage)
		handleError(c.Request("GET", artPage.Url, nil, ctx, nil))
	}
}

func handleError(err error) {
	log.Error("error occurs", zap.Error(err))
}
