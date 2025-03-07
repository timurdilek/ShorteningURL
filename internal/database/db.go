package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	user     = "DB_USER"
	password = "DB_PASSWORD"
	host     = "DB_HOST"
	port     = "DB_PORT"
	name     = "DB_NAME"
	filePath = "/internal/database/"

	// ENV FILE NAME
	envFile = ".env"
)

// Connection connects to the postgres DB using
// parameters from the .env file
func Connection() *pgxpool.Pool {

	// Get the current directory
	directory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Getting parameters from .env file
	if err = godotenv.Load(directory + filePath + envFile); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get parameters for calling pgxpool.New function
	URL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv(user), os.Getenv(password), os.Getenv(host), os.Getenv(port), os.Getenv(name))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, URL)

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Database connection error!\n%v", err)
	}

	log.Println("Connection to database successful")
	return pool
}
