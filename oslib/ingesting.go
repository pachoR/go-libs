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

func IngestFromStruct[K interface{}](indexName string, data K) error {
	if len(indexName) == 0 {
		return fmt.Errorf("Err! invalid indexName: indexName cannot be empty")
	}

	cte, err := GetConnection()
	if err != nil {
		return fmt.Errorf("Error getting opensearch connection: %s", err.Error())
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := opensearchapi.IndexRequest {
		Index: indexName,
		Body: strings.NewReader(string(bytes)),
	}

	res, err := req.Do(context.Background(), cte)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error indexing document: %s", string(bytes))
	}
	return nil
}

func IngestFromChunk[K interface{}](indexName string, data []K) error {
	if len(indexName) == 0 {
		return fmt.Errorf("Err! invalid indexName: indexName cannot be empty")
	}

	ingested := 0
	for i, v := range data {
		err := IngestFromStruct(indexName, v)
		if err != nil {
			log.Printf("ERR! ingestig (%d): %s", i, err.Error())
		} else {
			ingested++
		}
	}

	log.Printf("Finishing ingesting %d/%d on index %s", ingested, len(data), indexName)
	return nil
}

func IngestDataFromJson[T any](indexName string, jsonPath string) error {
	if len(indexName) == 0 && len(jsonPath) == 0 {
		return fmt.Errorf("indexName and jsonPath cannot be empty")
	}

	cte, err := GetConnection()
	if err != nil {
		return err
	}

	ndjsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", jsonPath, err)
	}

	var data []T
	data, err = ndjson.Unmarshal[T](ndjsonFile)
	if err != nil {
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
			return fmt.Errorf("Error indexing document: %s", string(itemJson))
		}
	}

	log.Printf("Completed ingestion of %d documents into index %s", len(data), indexName)
	return nil
}
