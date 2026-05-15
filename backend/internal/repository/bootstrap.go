package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

func (s *Store) Ready(ctx context.Context) error {
	return s.Pool.Ping(ctx)
}

// if path is "", path is acquired via os.Getenv("MIGRATION_PATH")
func (s *Store) CreateTableByFile(ctx context.Context, path string) error {
	if path == "" {
		path = s.migrationPath()
	}

	createQueryByte, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read the path to schema.sql: %w", err)
	}
	createQuery := string(createQueryByte)

	if len(strings.Fields(createQuery)) < 5 {
		return errors.New("incorrect creation schema")
	}
	tableName := strings.Fields(createQuery)[2]
	if tableName == "" {
		return errors.New("empty table name")
	}

	exists, err := s.checkExists(ctx, tableName)
	if err != nil {
		return fmt.Errorf("can't check if table already exists: %w", err)
	}
	if exists {
		return errors.New("table already exists!")
	}

	_, err = s.Pool.Exec(ctx, createQuery)
	if err != nil {
		return fmt.Errorf("table creation unsucessful: %w", err)
	}
	return nil
}

func (s *Store) checkExists(ctx context.Context, tableName string) (bool, error) {
	row := s.Pool.QueryRow(ctx, checkTableExists, tableName)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("scanning query result: %w", err)
	}
	return exists, nil
}

func (s *Store) migrationPath() string {
	if path := os.Getenv("MIGRATION_PATH"); path != "" {
		return path
	}
	return "../../repository/migrations/schema.sql"
}

func (s *Store) Initiate(ctx context.Context, path string) error {
	if err := s.Ready(ctx); err != nil {
		return fmt.Errorf("initiating DB: %w", err)
	}

	if err := s.CreateTableByFile(ctx, path); err != nil {
		return fmt.Errorf("creating quots table: %w", err)
	}
	return nil
}
