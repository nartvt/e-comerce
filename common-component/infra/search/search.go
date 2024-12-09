package search

import (
	"bytes"
	"common-component/config"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
	index  string
}

func InitElasticSearch(cfg *config.ElasticSearchConfig) (*ElasticsearchClient, error) {
	urls := cfg.BuildElasticSearchConnectionString()
	connectString := elasticsearch.Config{
		Addresses:            urls,
		Username:             cfg.User,
		Password:             cfg.Password,
		EnableRetryOnTimeout: true,
		RetryOnStatus:        cfg.RetryOnStatusCodes,
	}

	client, err := elasticsearch.NewClient(connectString)
	if err != nil {
		panic(err)
	}

	res, err := client.Info()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	log.Println(res)

	return &ElasticsearchClient{
		client: client,
		index:  cfg.Index,
	}, nil
}

// CreateIndex creates an index with custom mapping
func (ec *ElasticsearchClient) CreateIndex() error {
	// Prepare the index mapping
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]string{
					"type": "keyword",
				},
				"name": map[string]string{
					"type": "text",
				},
				"description": map[string]string{
					"type": "text",
				},
				"price": map[string]string{
					"type": "float",
				},
				"category": map[string]string{
					"type": "keyword",
				},
				"created_at": map[string]string{
					"type": "date",
				},
			},
		},
	}

	// Convert mapping to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return fmt.Errorf("error encoding mapping: %w", err)
	}

	// Create index request
	req := esapi.IndicesCreateRequest{
		Index: ec.index,
		Body:  &buf,
	}

	// Perform the request
	res, err := req.Do(context.Background(), ec.client)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		// Check if the error is because the index already exists
		if res.StatusCode == 400 {
			log.Println("Index already exists, skipping creation")
			return nil
		}
		return fmt.Errorf("error creating index: %s", res.Status())
	}

	fmt.Println("Index created successfully")
	return nil
}

// Optional: Check if index exists method
func (ec *ElasticsearchClient) IndexExists() (bool, error) {
	// Prepare the exists request
	req := esapi.IndicesExistsRequest{
		Index: []string{ec.index},
	}

	// Perform the request
	res, err := req.Do(context.Background(), ec.client)
	if err != nil {
		return false, fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	// Return true if status code is 200 (exists), false otherwise
	return res.StatusCode == 200, nil
}
