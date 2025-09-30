package oslib

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
)

var opensearchClient *opensearch.Client

func testConnection(client *opensearch.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := client.Ping(client.Ping.WithContext(ctx))
	if err != nil {
		return  err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("ping response error: %s", resp.Status())
	}

	return nil
}

func GetConnection() (*opensearch.Client, error) {
	if opensearchClient != nil {
		if err := testConnection(opensearchClient); err != nil {
			log.Println("Existing connection failed, reconnecting")
			opensearchClient = nil
		} else {
			return opensearchClient, nil
		}
	}

	var err error
	opensearchClient, err = opensearch.NewClient(opensearch.Config{
		Addresses: []string{os.Getenv("OPENSEARCH_URL")},
		Username: os.Getenv("OS_USER"),
		Password: os.Getenv("OS_PASSWORD"),
	})

	if err != nil {
		return nil, fmt.Errorf("error creating OpenSearch client: %w", err)
	}

	if err = testConnection(opensearchClient); err != nil {
		return nil, fmt.Errorf("OpenSearch connection test failed: %w", err)
	}
	log.Println("Connection established")
	return opensearchClient, nil
}
