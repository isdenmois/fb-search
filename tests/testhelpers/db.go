package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestDatabase struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
	ConnStr   string
}

func NewTestDatabase(ctx context.Context) (*TestDatabase, error) {
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("fb_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Env: map[string]string{
					"POSTGRES_INITDB_ARGS": "--encoding=UTF8",
				},
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Wait for container to be fully ready
	time.Sleep(2 * time.Second)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection with retry
	for i := 0; i < 5; i++ {
		err = pool.Ping(ctx)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		pool.Close()
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := runMigrations(ctx, pool); err != nil {
		pool.Close()
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &TestDatabase{
		Container: container,
		Pool:      pool,
		ConnStr:   connStr,
	}, nil
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationSQL := `
		CREATE TABLE IF NOT EXISTS "books" (
			"id" text PRIMARY KEY,
			"title" text NOT NULL,
			"search" text NOT NULL,
			"authors" text,
			"series" text,
			"serno" text,
			"lang" text,
			"size" integer
		);
		CREATE INDEX IF NOT EXISTS "search_ru_idx" ON "books" USING gin (to_tsvector('russian', "search"));
		CREATE INDEX IF NOT EXISTS "search_simple_idx" ON "books" USING gin (to_tsvector('simple', "search"));
	`

	_, err := pool.Exec(ctx, migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

func (td *TestDatabase) Cleanup(ctx context.Context) error {
	if td.Pool != nil {
		td.Pool.Close()
	}
	if td.Container != nil {
		return td.Container.Terminate(ctx)
	}
	return nil
}

func (td *TestDatabase) TruncateBooks(ctx context.Context) error {
	_, err := td.Pool.Exec(ctx, "TRUNCATE TABLE books RESTART IDENTITY")
	return err
}
