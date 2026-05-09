package main

import (
	"context"
	"log"
	"os"

	"stonkx/internal/handler"
	"stonkx/internal/repository"
	"stonkx/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	var connString string

	//should app run locally, .env will be read from a specified path
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(`../../local.env`); err != nil {
			log.Fatal(err)
		}
	}
	connString = os.Getenv("DB_CONF_FULL")

	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "./frontend"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("connecting db: %v", err)
	}

	store := repository.New(pool)
	quotsService := service.NewQuots(store, store)

	hand := handler.New(quotsService)
	hand.Start()

}
