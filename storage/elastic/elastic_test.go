package elastic

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
)

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

	createMapping(client)
}

func createMapping(client *elasticsearch.Client) {
	// client.Indices.Exists("article")
	// esapi.IndexRequest{
	// 	Index:      "article",
	// 	DocumentID: string,
	// }

}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}
