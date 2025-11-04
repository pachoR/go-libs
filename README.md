# go-libs



A comprehensive Go library for interacting with OpenSearch to manage indices, documents, and perform advanced search operations.

## Features

- **Index Management**: Create, update, and delete OpenSearch indices
- **Document Operations**: Insert, update, delete, and retrieve documents
- **Advanced Search**: Execute complex queries with filters, aggregations, and sorting
- **Bulk Operations**: Perform bulk insertions, updates, and deletions
- **Connection Management**: Connection handling

## - http
- **Http methods:** Protocol methods with the following:
    - **Raw:** Returns the *http.Response as it is, the caller have responsability for defering the body and full liberty on raw error handleling.
    - **Body:** The method defers and handle error. It returns the bytes.
    - **BodyWithRetries:** Self explanatory.
## Envs
OPENSEARCH_URL
OS_USER
OS_PASSWORD
PG_URL