package oslib

import (
	"fmt"
	"io"

	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

func PrintOSErr(res *opensearchapi.Response) error {
	body, _ := io.ReadAll(res.Body)
	return fmt.Errorf("OpenSearch error: %s\nResponse: %s\n", res.Status(), string(body))
}
