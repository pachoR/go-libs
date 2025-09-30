package oslib

import (
	"context"
	"fmt"
	"log"
	"strings"

	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func CreateIndex(indexName string) error {
	cte, err := GetConnection()
	if err != nil {
		log.Fatal("Error: ", err.Error())
		return err
	}

	settings := strings.NewReader(`{
		"settings": {
			"index": {
				"number_of_shards": 1,
				"number_of_replicas": 0
			}
		}
	}`)

	req := opensearchapi.IndicesCreateRequest {
		Index: indexName,
		Body: settings,
	}

	res, err := req.Do(context.Background(), cte)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error creatring index (%s): %s", indexName, res.String())
	}

	fmt.Printf("Index %s created\n", indexName)

	return nil
}