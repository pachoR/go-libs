package main

import (
	"fmt"

	"github.com/joho/godotenv"
	// oslib "github.com/pachoR/go-libs/oslib"
	postgres "github.com/pachoR/go-libs/postgreslib"
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
	conn, err := postgres.GetConnection()
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Println(&conn)
	}

	postgres.CloseConnection()
}
