package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func Connect(path string) (*pgxpool.Pool, error) {
	var envKeys map[string]string
	envKeys, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading conf file: %w", err)
	}

	confString := envKeys["DB_CONFIG_STRING"]

	connConfig, err := pgxpool.ParseConfig(confString)
	if err != nil {
		return nil, fmt.Errorf("parsing config string: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("establishing connection: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging DB: %w", err)
	}

	return pool, nil

}
