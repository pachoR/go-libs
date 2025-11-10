package postgreslib

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var pgConnection *pgx.Conn

func connect() error {
	pgUrl := os.Getenv("PG_URL")
	if len(pgUrl) == 0 {
		return fmt.Errorf("NOT POSTGRES CONNECTION URL PROVIDED")
	}

	conn, err := pgx.Connect(context.Background(), pgUrl)
	if err != nil {
		return err
	}

	pgConnection = conn
	return nil
}

func TestConnection() (error) {
	return pgConnection.Ping(context.Background())
}

func GetConnection() (*pgx.Conn, error) {
	const maxRetries = 3
	for retry := 0; retry < maxRetries; retry++ {
		if pgConnection == nil {
			if err := connect(); err != nil {
				log.Printf("Connection attempt (%d/%d) failed: %v", retry+1, maxRetries, err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
		}

		if err := TestConnection(); err != nil {
			log.Printf("Ping attempt failed attempt (%d/%d): %v", retry+1, maxRetries, err.Error())
			time.Sleep(5 * time.Second)
			pgConnection = nil
			continue
		}

		return pgConnection, nil
	}

	return nil, fmt.Errorf("Failed to connect to Postgres after %d retries", maxRetries)
}

func CloseConnection() {
	if pgConnection != nil {
		pgConnection.Close(context.Background())
		fmt.Println("Connection closed")
	} else {
		fmt.Println("No connection")
	}
}
