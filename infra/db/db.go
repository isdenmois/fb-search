package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	ctx := context.Background()
	databaseUrl := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)

	return pool, err
}

func RebuildDb(pool *pgxpool.Pool) {
	pool.Exec(context.Background(), "TRUNCATE TABLE books RESTART IDENTITY")
	pool.Exec(context.Background(), "VACUUM")
}
