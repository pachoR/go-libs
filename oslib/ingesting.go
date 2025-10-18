package oslib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	ndjson "github.com/adotkp/ndjson"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func IngestDataFromJson[T any](indexName string, jsonPath string) error {
	if len(indexName) == 0 && len(jsonPath) == 0 {
		log.Println("indexName and jsonPath cannot be empty")
		return fmt.Errorf("indexName and jsonPath cannot be empty")
	}

	cte, err := GetConnection()
	if err != nil {
		log.Fatalf("Error getting OS connection: %s", err.Error())
		return err
	}

	ndjsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Printf("Error reading file %s: %v", jsonPath, err)
		return fmt.Errorf("error reading file %s: %w", jsonPath, err)
	}

	var data []T
	data, err = ndjson.Unmarshal[T](ndjsonFile)
	if err != nil {
		log.Fatalf("Error unmarshalling ndjson file: %s", err.Error())
		return err
	}

	for _, item := range data {
		itemJson, _ := json.Marshal(item)

		req := opensearchapi.IndexRequest {
			Index: indexName,
			Body: strings.NewReader(string(itemJson)),
		}

		res, err := req.Do(context.Background(), cte)
		if err != nil {
			log.Printf("Error indexing document: %s", string(itemJson))
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("Error indexing document: %s", string(itemJson))
		}
	}

	log.Printf("Completed ingestion of %d documents into index %s", len(data), indexName)
	return nil
}
