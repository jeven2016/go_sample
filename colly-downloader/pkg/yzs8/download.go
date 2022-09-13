package yzs8

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/liuzl/gocc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"colly-downloader/pkg/models"
)

var articlesChan = make(chan *models.ArticlePage, 1000)

var wg = sync.WaitGroup{}

var baseUrl string

var log = initLog()

// ensureCatalog, return existed, error
func ensureCatalog(collection *mongo.Collection, catalog *models.CatalogDoc) (string, error) {
	result := collection.FindOne(context.TODO(), bson.M{"name": catalog.Name})
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			insertResult, err := collection.InsertOne(context.TODO(), catalog)

			return insertResult.InsertedID.(string), err
		}
		return "", err
	}
	catalog = new(models.CatalogDoc)
	err = result.Decode(catalog)
	if err != nil {
		return "", err
	}
	return catalog.Id.String(), err
}

func Start() {
	go func() {
		wg.Wait()
	}()

	db, err := CreateMongoClient(log)
	if err != nil {
		return
	}
	collection := db.Collection("catalog")

	catalogId, existed := ensureCollection(collection, err)
	if !existed {
		return
	}

	collector := newCollector()
	homeUrl := "https://yazhouse8.com/article.php?cate=1"
	baseUrl = getBaseUri(homeUrl)
	parsePageLinks(homeUrl, collector, articlesChan)

	zhConvertor, err := gocc.New("s2t")
	if err != nil {
		log.Error("i18n convertor error", zap.Error(err))
	}
	for i := 0; i < 3; i++ {
		downloadArticle(i, collection, articlesChan, zhConvertor, catalogId)
	}
}

func ensureCollection(catalogCol *mongo.Collection, err error) (string, bool) {
	individualCatalog := &models.CatalogDoc{
		Name:         "yzs8",
		Order:        1,
		ArticleCount: 0,
		Description:  "yzs8",
		CreateDate:   time.Now(),
		LastUpdate:   time.Now(),
	}

	catalogId, err := ensureCatalog(catalogCol, individualCatalog)
	if err != nil {
		log.Error("failed to check the catalog", zap.Error(err))
		return "", false
	}
	log.Info("catalog id is", zap.String("id", catalogId))
	return catalogId, true
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

func parsePageLinks(homeUrl string, collector *colly.Collector, urlChan chan<- *models.ArticlePage) {
	parseArticles(homeUrl, collector, urlChan)

	collector.OnHTML(".pager a[href]", func(element *colly.HTMLElement) {
		if at.Load() >= 2 {
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

func parseArticles(pageUrl string, c *colly.Collector, artChan chan<- *models.ArticlePage) {
	c.OnHTML(".articleList>p>.img-center", func(element *colly.HTMLElement) {
		artChan <- &models.ArticlePage{Name: element.Text, Url: element.Attr("href")}
	})
}

func downloadArticle(taskId int, collection *mongo.Collection, urlChan <-chan *models.ArticlePage,
	zhConvertor *gocc.OpenCC, catalogId string) {
	for artPage := range urlChan {
		documents, err := collection.CountDocuments(context.TODO(), bson.M{"name": artPage.Name})
		if err != nil {
			log.Warn("failed to count documents with name",
				zap.Error(err), zap.String("name", artPage.Name))
			continue
		}
		if documents > 0 {
			log.Info("document exists, ignored", zap.String("name", artPage.Name))
			continue
		}

		// load article page and get the content to save

		_, err = collection.InsertOne(context.TODO(), models.Article{
			Name:      artPage.Name,
			CatalogId: catalogId,
			Content:   content,
		})
		if err != nil {
			log.Warn("failed to insert document with name",
				zap.Error(err), zap.String("name", artPage.Name))
			continue
		}
		log.Info("Inserted a document", zap.String("name", artPage.Name))
	}
}

func handleError(err error) {
	log.Error("error occurs", zap.Error(err))
}
