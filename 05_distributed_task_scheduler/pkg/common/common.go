package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DefaultHeartbeat = 5 * time.Second
)

func GetDBConnectionString() string {

	var missingVars []string

	checkEnvVars := func(envVar, envVarName string) {
		if envVar == "" {
			missingVars = append(missingVars, envVarName)
		}
	}

	dbUser := os.Getenv("POSTGRES_USER")

	checkEnvVars(dbUser, "POSTGRES_USER")

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	checkEnvVars(dbPassword, "POSTGRES_PASSWORD")

	dbName := os.Getenv("POSTGRES_DB")

	checkEnvVars(dbName, "POSTGRES_DB")

	dbHost := os.Getenv("POSTGRES_HOST")
	checkEnvVars(dbHost, "POSTGRES_HOST")

	if len(missingVars) > 0 {
		log.Fatalf("The following required ENV vars are not set failed !!! --> %s", strings.Join(missingVars, ", "))
	}

	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPassword, dbHost, dbName)

}

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
