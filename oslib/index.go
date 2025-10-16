package oslib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func getIndicesByAlias(aliasName string) ([]string, error) {
	client, err := GetConnection()
	if err != nil {
		log.Fatal("Error: ", err.Error())
		return nil, err
	}

	req := opensearchapi.IndicesGetAliasRequest {
		Name: []string{aliasName},
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return []string{}, nil
		}
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("error getting alias(%s):(%d) %s", aliasName, res.StatusCode, res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	indices := []string{}
	for indexName := range result {
		indices = append(indices, indexName)
	}
	return indices, nil
}

func getFormattedTime() string {
	current_time := time.Now()
	return fmt.Sprintf("%d-%d-%d-%d-%d-%d", current_time.Year(), current_time.Month(), current_time.Day(), current_time.Hour(), current_time.Minute(), current_time.Second())
}

func CreateIndex(indexName string) error {
	cte, err := GetConnection()
	if err != nil {
		log.Fatal("Error: ", err.Error())
		return err
	}

	fullIndexName := fmt.Sprintf("%s-%s", indexName, getFormattedTime())

	settings := strings.NewReader(`{
		"settings": {
			"index": {
				"number_of_shards": 1,
				"number_of_replicas": 0
			}
		}
	}`)

	req := opensearchapi.IndicesCreateRequest {
		Index: fullIndexName,
		Body: settings,
	}

	res, err := req.Do(context.Background(), cte)
	if err != nil {
		log.Fatal("Error creating index: ", err.Error())
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error creating index (%s): %s", indexName, res.String())
		return fmt.Errorf("Error creating index (%s): %s", indexName, res.String())
	}
	log.Printf("Index %s created successfully\n", fullIndexName)

	// getting the old indeces that use or were using this alias to eliminate them
	oldIndices, err := getIndicesByAlias(indexName)
	if err != nil {
		log.Fatal("Error getting indices by alias", err.Error())
		return err
	}

	var actions []string
	for _, oldIndex := range oldIndices {
		actions = append(actions, fmt.Sprintf(`{
			"remove": {
				"index": "%s",
				"alias": "%s"
			}
		}`, oldIndex, indexName))
	}

	actions = append(actions, fmt.Sprintf(`{
		"add": {
			"index": "%s",
			"alias": "%s"
		}
	}`, fullIndexName, indexName))


	updateAliasBody := strings.NewReader(fmt.Sprintf(`{
 		"actions": [%s]
	}`, strings.Join(actions, ", ")))

	aliasReq := opensearchapi.IndicesUpdateAliasesRequest {
		Body: updateAliasBody,
	}

	aliasRes, err := aliasReq.Do(context.Background(), cte)
	if err != nil {
		return err
	}
	defer aliasRes.Body.Close()

	if aliasRes.IsError() {
		log.Fatalf("Error deleting alias from prev index/indexes: %s: %s", strings.Join(actions, ", "), aliasRes.String())
		return fmt.Errorf("Error deleting alias from prev index/indexes: %s: %s", strings.Join(actions, ", "), aliasRes.String())
	}
	log.Printf("Alias updated successfully\n")

	if len(oldIndices) > 0 {
		oldIndicesAsStr := fmt.Sprintf("%s", strings.Join(oldIndices, ", "))
		deleteReq := opensearchapi.IndicesDeleteRequest {
			Index: oldIndices,
		}

		deleteRes, err := deleteReq.Do(context.Background(), cte)
		if err != nil {
			log.Fatalf("Error deleting old indices: %s", err.Error())
			return err
		}
		defer deleteRes.Body.Close()

		if deleteRes.IsError() {
			log.Fatalf("Error deleting the old indices: %s", oldIndicesAsStr)
			return fmt.Errorf("Error deleting the old indices: %s", oldIndicesAsStr)
		}

		log.Printf("Old indeces (%s) removed successfully", oldIndicesAsStr)
	}

	return nil
}
