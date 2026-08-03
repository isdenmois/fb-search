package ports

import (
	"context"

	"fb-search/domain"
)

// SearchConfig selects the PostgreSQL text-search configuration.
// Language is either "russian" or "simple".
type SearchConfig struct {
	Language string
}

// BookRepository is the port (interface) owned by the application layer
// that infrastructure repositories implement.
type BookRepository interface {
	SearchBooks(ctx context.Context, query string, cfg SearchConfig) ([]domain.Book, error)
	FindById(ctx context.Context, id string) (domain.Book, error)
	RebuildDb(ctx context.Context) error
	InsertBatch(ctx context.Context, rows RowBatchSource) (uint64, error)
}

// RowBatchSource is a pgx-compatible batch row source, defined in the
// application layer to avoid leaking infrastructure types into ports.
type RowBatchSource interface {
	Next() bool
	Values() ([]interface{}, error)
	Err() error
}
