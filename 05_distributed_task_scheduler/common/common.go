package common

import (
	"context"
	"log"
	"time"
)

// this function establishes the connection to the DB
func ConnectToDatabase(ctx context.Context, dbConnectionString string) (*pgxpool.Pool, error) {

	var dbPool *pgxpool.Pool
	var err error

	retryCount := 0

	// retry mechanism for DB
	for retryCount < 5 {

		dbPool, err := pgxpool.Connect(ctx, dbConnectionString)

		if err == nil {
			break
		}

		log.Println("Failed to connect to the database... Retrying in 5 seconds!!!")

		time.Sleep(time.Second * 5)

		retryCount++
	}

	if err != nil {
		log.Println("Ran out of retries!!")
		return nil, err

	}

	log.Println("Connected to DB successfully")

	return dbPool, err

}
