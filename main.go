package main

import (
	"fmt"

	"github.com/joho/godotenv"
	oslib "github.com/pachoR/go-libs/oslib"
)

func init() {
	godotenv.Load()
}

func main() {
	err := oslib.CreateIndex("index_test")
	if err != nil {
		fmt.Printf("main: %s\n", err.Error())
	}
}