package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []string{
		"migrations/001_initial.sql",
		"migrations/002_mod_data.sql",
	}
	for _, path := range migrations {
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			// Log but continue — tables may already exist on redeploy
			fmt.Printf("migration %s: %v\n", path, err)
		}
	}
	return nil
}
