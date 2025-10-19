package main

import (
	"fmt"

	"github.com/joho/godotenv"
	oslib "github.com/pachoR/go-libs/oslib"
)

func init() {
	godotenv.Load()
}

type Persona struct {
	Id 			string 		`json:"id"`
	Name 		string		`json:"name"`
	Age 		int 		`json:"age"`
	Salary 		float64 	`json:"salary"`
	IsActive	bool	  	`json:"is_active"`
}

func main() {
	indexName := "index_test"
	// err := oslib.CreateIndex(indexName, "mapping/test-mapping.json")
	// if err != nil {
	// 	fmt.Printf("main: %s\n", err.Error())
	// }

	// fmt.Println("Ingesting data")
	// err = oslib.IngestDataFromJson[Persona](indexName, "./locals/personas.json")

	query, err := oslib.SearchWithQuery(indexName, `{
		"query": {
			"match_all": {}
		}
	}`)

	if err != nil {
		fmt.Printf("Error searching: %s", err.Error())
	}

	strQuery := string(query)
	fmt.Printf("strquery: %s\n", strQuery)
}
