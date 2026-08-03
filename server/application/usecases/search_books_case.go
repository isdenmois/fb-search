package usecases

import (
	"context"
	"strings"

	"fb-search/application/ports"
	"fb-search/domain"
)

// SearchBooksCase is the use case that resolves a user search query to a
// concrete search configuration (russian vs simple) and delegates to the
// book repository.
type SearchBooksCase struct {
	booksRepository ports.BookRepository
}

func NewSearchBooksCase(booksRepository ports.BookRepository) *SearchBooksCase {
	return &SearchBooksCase{booksRepository: booksRepository}
}

func (self *SearchBooksCase) Execute(ctx context.Context, q string) ([]domain.Book, error) {
	cfg := ports.SearchConfig{Language: "simple"}
	if ContainsCyrillic(q) {
		cfg.Language = "russian"
	}

	return self.booksRepository.SearchBooks(ctx, strings.ToLower(q), cfg)
}
