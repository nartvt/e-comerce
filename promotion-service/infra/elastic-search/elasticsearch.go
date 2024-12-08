package elasticsearch

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var esClient *elasticsearch.Client

func InitES() {
	var (
		err error
	)
	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{},
	})
	if err != nil {
		panic(err)
	}
}

func GetES() *elasticsearch.Client {
	if esClient == nil {
		InitES()
	}

	return esClient
}

func CloseES() {
	if closer, ok := esClient.Transport.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			log.Fatalf("Error closing transport: %s", err)
		}
	}
}
