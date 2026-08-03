package ports

import "context"

// BookFileStorage is the port for reading book archive contents
// (implemented by infrastructure/storage).
type BookFileStorage interface {
	ReadFile(ctx context.Context, file string, path string) ([]byte, error)
}
