package oslib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"log"
	"strings"
	"time"

	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func getIndicesByAlias(aliasName string) ([]string, error) {
	client, err := GetConnection()
	if err != nil {
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

func CreateIndex(indexName string, mappingPath string) error {
	cte, err := GetConnection()
	if err != nil {
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
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error creating index (%s): %s", indexName, res.String())
	}
	log.Printf("Index %s created successfully\n", fullIndexName)

	// Add mapping
	err = ApplyMapping(fullIndexName, mappingPath)
	if err != nil {
		return fmt.Errorf("Error applying mapping: %s", err.Error())
	}

	// getting the old indeces that use or were using this alias to eliminate them
	oldIndices, err := getIndicesByAlias(indexName)
	if err != nil {
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
			return err
		}
		defer deleteRes.Body.Close()

		if deleteRes.IsError() {
			return fmt.Errorf("Error deleting the old indices: %s", oldIndicesAsStr)
		}

		log.Printf("Old indeces (%s) removed successfully", oldIndicesAsStr)
	}

	return nil
}

func ApplyMapping(indexName string, mappingFilePath string) error {
	cte, err := GetConnection()
	if err != nil {
		return err
	}

	if len(indexName) == 0 || len(mappingFilePath) == 0 {
		return fmt.Errorf("Index name or mapping file path cannot be empty")
	}

	mappingData, err := os.ReadFile(mappingFilePath)
	if err != nil {
		return fmt.Errorf("Error reading mapping file: %s", err.Error())
	}

	mappingReq := opensearchapi.IndicesPutMappingRequest {
		Index: []string{indexName},
		Body: strings.NewReader(string(mappingData)),
	}

	res, err := mappingReq.Do(context.Background(), cte)
	if err != nil {
		return fmt.Errorf("Error setting a mapping request: %s", err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error setting a mapping request: %s", res.String())
	}

	log.Print("Mapping created successfully")
	return nil
}
