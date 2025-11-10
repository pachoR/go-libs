package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"

	// oslib "github.com/pachoR/go-libs/oslib"
	// postgres "github.com/pachoR/go-libs/postgreslib"
	http "github.com/pachoR/go-libs/http"
)

func init() {
	godotenv.Load()
}

type SymOvw struct {
	Description  	string  `json:"description"`
	DisplaySymbol 	string 	`json:"displaySymbol"`
	Symbol 			string 	`json:"symbol"`
	Type 			string 	`json:"type"`
}

type SymbolOverview struct {
	Count 	int 		`json:"count"`
	Result 	[]SymOvw 	`json:"result"`
}

func main() {
	url := "https://finnhub.io/api/v1/search?q=AAPL&exchange=US"

	h := map[string]string {
		"X-Finnhub-Token": fmt.Sprintf("%s", os.Getenv("FINN_TOKEN")),
	}

	r, err := http.GetWithHeader(url, h)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer r.Body.Close()

	var symOvw SymbolOverview
	bytes, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(bytes, &symOvw)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	symOvwJson, _ := json.MarshalIndent(symOvw, "", " ")
	fmt.Printf("symOvwJson: %s\n", string(symOvwJson))
}
