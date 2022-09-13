package yzs8

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"colly-downloader/pkg/models"
)

var linksChan chan string = make(chan string, 1000)

var wg = sync.WaitGroup{}

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
		return catalog.Id.String(), nil
	}
	return "", nil
}

func GoDown() {
	go func() {
		wg.Wait()
	}()

	log := initLog()
	db, err := CreateMongoClient(log)
	if err != nil {
		return
	}

	catalogCol := db.Collection("catalog")

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
		return
	}
	log.Info("catalog id is", zap.String("id", catalogId))

}

// New collector
func newCollector() *colly.Collector {
	// c := colly.NewCollector(colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"))
	c := colly.NewCollector()
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

// https://yazhouse8.com/article.php?cate=1  首页
// 文章：<div class="articleList"><p><span>12、</span>
// <a class="img-center" target="_blank" href="article/27078.html">有雪</a></p></div>
// 下一页 ： https://yazhouse8.com/article.php?page=2&cate=1
// check： <ul class="pager"><li><a href="article.php?cate=1">首页</a>&nbsp;&nbsp;
//      </li><li><a href="article.php?page=2&amp;cate=1">下一页</a></li></ul>

func parsePageLinks(homeUrl string, collector *colly.Collector, urlChan chan<- string) {
	// 当前页加入channel
	urlChan <- homeUrl

	collector.OnHTML(".pager a[href]", func(element *colly.HTMLElement) {
		text := element.Text
		if text == "下一页" {
			href := element.Attr("href")
			uri, err := url.ParseRequestURI(href)
			if err != nil {
				// print log
				return
			}
			nextPage := uri.String()
			handleError(element.Request.Visit(nextPage))
		}
	})
	handleError(collector.Visit(homeUrl))
}

func downloadPages(tasks int, collector *colly.Collector, urlChan <-chan string) {
	for i := 0; i < tasks; i++ {
		go func(taskId int) {
			for url := range urlChan {
				c := collector.Clone()
				handlePage(url, c, taskId)
			}
		}(i)
	}
}

func handlePage(url string, collector *colly.Collector, id int) {
	collector.OnHTML("", func(element *colly.HTMLElement) {
		// get the content

		var content = ""

	})
}

func handleError(err error) {
	if err != nil {
		print(err)
		return
	}
}
