package main

import (
	"fmt"
	"log"
	"github.com/joho/godotenv"
	// oslib "github.com/pachoR/go-libs/oslib"
	// postgres "github.com/pachoR/go-libs/postgreslib"
	http "github.com/pachoR/go-libs/http"
)

func init() {
	godotenv.Load()
}

func main() {
	okurl := "http://localhost:3000"
	failurl := okurl + "/fail"
	fmt.Println("Get:")
	get, err := http.GetBody(okurl)
	if err != nil {
		log.Fatalf("Error on correct Get: %s", err.Error())
	}
	fmt.Printf("Correct: %s\n", string(get))

	getFail, err := http.GetBodyWithRetries(failurl)
	if err != nil {
		fmt.Printf("Failed: %s\n", string(getFail))
	}
	fmt.Printf("Failed: %s\n", string(getFail))


	fmt.Println("\nPost:")
	mockPayload := struct{
		Name string
		Age  int
	}{
		Name: "Alejandra",
		Age: 18,
	}
	post, err := http.PostBody(okurl, mockPayload)
	if err != nil {
		log.Fatalf("Error on correct Post: %s", err.Error())
	}
	fmt.Printf("Correct: %s\n", string(post))

	postFail, err := http.PostBodyWithRetries(failurl, mockPayload)
	if err != nil {
		fmt.Printf("Failed: %s\n", string(getFail))
	}
	fmt.Printf("Failed: %s\n", string(postFail))
}
