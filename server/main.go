package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	defaultDatabaseURL = "postgres://admin:admin@localhost:5432/aaa"
	defaultHttpPort    = "8080"
)

func main() {
	ctx := context.Background()
	logger := log.Default()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = defaultDatabaseURL
		logger.Print("DATABASE_URL environment variable not set")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = defaultHttpPort
		logger.Print("HTTP_PORT environment variable not set")
	}

	db, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if err = fiber.New().Listen(fmt.Sprintf(":%v", httpPort)); err != nil {
		logger.Fatal(err)
	}
}
