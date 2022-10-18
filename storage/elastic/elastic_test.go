package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/estransport"
)

type Article struct {
	Id         string    `bson:"_id" json:"id"`
	Name       string    `bson:"name" json:"name"`
	CatalogId  string    `bson:"catalogId" json:"catalogId"`
	Content    string    `bson:"content,omitempty" json:"content"`
	CreateDate time.Time `bson:"createDate" json:"createDate"`
}

var retryBackoff = backoff.NewExponentialBackOff()

func TestSimpleCase(t *testing.T) {
	retryBackoff.InitialInterval = time.Second

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:  []string{"http://127.0.0.1:9200"},
		MaxRetries: 3,
		// By default, the client retries the request up to three times; to set a different limit,
		// use the MaxRetries field. To change the list of response status codes which should be retried,
		// use the RetryOnStatus field. Together with the RetryBackoff option, you can use it to retry
		// requests when the server sends a 429 Too Many Requests response:
		RetryOnStatus: []int{429, 502, 503, 504},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			d := retryBackoff.NextBackOff()
			fmt.Printf("Attempt: %d | Sleeping for %s...\n", i, d)
			return d
		},

		Logger: &estransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	})
	handleError(err)
	log.Println("version=", elasticsearch.Version)

	// createMapping(client)
	// InserArticle(client)
	QueryArticle(client)
}

// https://www.jianshu.com/p/075c0ed51053
func QueryArticle(client *elasticsearch.Client) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name":    "事务",
				"content": "事务",
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<font color='red'>"},
			"post_tags": []string{"</font>"},
			"fields": map[string]interface{}{
				"name": map[string]interface{}{},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(err)
	}

	res, err := client.Search(client.Search.WithContext(context.Background()),
		client.Search.WithIndex("article"),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty())
	if err != nil {
		panic(err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func createMapping(client *elasticsearch.Client) {
	exists, err := esapi.IndicesExistsRequest{
		Index: []string{"article"},
	}.Do(context.TODO(), client)

	if err != nil {
		handleError(err)
	}

	if exists.StatusCode != 404 {
		// success 200
		result, err := esapi.IndicesDeleteRequest{Index: []string{"article"}}.Do(context.Background(), client)
		handleError(err)
		defer result.Body.Close()

		print(result.StatusCode)

		// elasticsearch.In.Exists([]string{index})
	}
	abs, err := filepath.Abs("./article_mapping.json")
	handleError(err)

	file, err := ioutil.ReadFile(abs)
	handleError(err)
	result2, err := esapi.IndicesCreateRequest{
		Index:   "article",
		Body:    bytes.NewReader(file),
		Timeout: 20 * time.Second,
	}.Do(context.Background(), client)
	handleError(err)
	defer result2.Body.Close()

	if result2.StatusCode == 200 {
		println("index 'article' is created")
	}

}

func InserArticle(client *elasticsearch.Client) {
	var article = &Article{
		Name: "JPA的事务注解@Transactional使用总结",
		Content: `在项目开发过程中，如果您的项目中使用了Spring的@Transactional注解，
有时候会出现一些奇怪的问题，例如：\n\n \n\n明明抛了异常却不回滚？\n\n嵌套事务执行报错？\n\n...等等\n\n \n\n很多的问题都是没有全面了解
@Transactional的正确使用而导致的，下面一段代码就可以让你完全明白@Transactional到底该怎么用。`,
		CatalogId:  "abcdef",
		CreateDate: time.Now(),
	}

	json, err := convertor.ToJson(article)
	handleError(err)

	result, err := esapi.IndexRequest{
		Index: "article",
		Body:  bytes.NewReader([]byte(json)),
	}.Do(context.Background(), client)
	defer result.Body.Close()

	if result.StatusCode == 201 {
		println("201 Created")
	}

}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}
