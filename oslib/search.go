package oslib

import (
	"context"
	"fmt"
	"log"
	"strings"
	"io"

	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func SearchWithQuery(indexName string, query string) ([]byte, error) {
	if len(indexName) == 0 || len(query) == 0 {
		return nil, fmt.Errorf("Error: indexName or query not provided")
	}

	cte, err := GetConnection()
	if err != nil {
		return nil, err
	}

	var indices []string
	indices = append(indices, indexName)
	req := opensearchapi.SearchRequest {
		Index: indices,
		Body: strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), cte)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Error performing query. For index: %s query: %s", indexName, query)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body response: %s", err.Error())
	}

	log.Println("Successfull query")
	return body, nil
}
