package db

import (
	"context"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*pgxpool.Pool, error) {
	ctx := context.Background()
	databaseUrl := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(pool)

	instance, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", instance)
	if err != nil {
		return nil, err
	}

	m.Up()

	return pool, err
}

func RebuildDb(pool *pgxpool.Pool) {
	pool.Exec(context.Background(), "TRUNCATE TABLE books RESTART IDENTITY")
	pool.Exec(context.Background(), "VACUUM")
}
