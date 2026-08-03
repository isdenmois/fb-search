package usecases

import (
	"context"
	"errors"
	"testing"

	"fb-search/application/ports"
	"fb-search/domain"

	"github.com/stretchr/testify/require"
)

// mockBooksRepository is a test double for ports.BookRepository that records
// the arguments passed to SearchBooks and returns a canned result/error.
type mockBooksRepository struct {
	receivedQuery string
	receivedCfg   ports.SearchConfig
	result        []domain.Book
	err           error
}

func (m *mockBooksRepository) SearchBooks(_ context.Context, query string, cfg ports.SearchConfig) ([]domain.Book, error) {
	m.receivedQuery = query
	m.receivedCfg = cfg
	return m.result, m.err
}

func (m *mockBooksRepository) FindById(_ context.Context, _ string) (domain.Book, error) {
	return domain.Book{}, errors.New("not implemented")
}

func (m *mockBooksRepository) RebuildDb(_ context.Context) error {
	return errors.New("not implemented")
}

func (m *mockBooksRepository) InsertBatch(_ context.Context, _ ports.RowBatchSource) (uint64, error) {
	return 0, errors.New("not implemented")
}

func TestSearchBooksCase_Execute_CyrillicQueryUsesRussianConfig(t *testing.T) {
	// arrange
	repo := &mockBooksRepository{result: []domain.Book{{Id: "1", Title: "Война и мир"}}}
	searchCase := NewSearchBooksCase(repo)

	// act
	books, err := searchCase.Execute(context.Background(), "Война и мир")

	// assert
	require.NoError(t, err)
	require.Equal(t, []domain.Book{{Id: "1", Title: "Война и мир"}}, books)
	require.Equal(t, "война и мир", repo.receivedQuery)
	require.Equal(t, ports.LanguageRussian, repo.receivedCfg.Language)
}

func TestSearchBooksCase_Execute_LatinQueryUsesSimpleConfig(t *testing.T) {
	// arrange
	repo := &mockBooksRepository{result: []domain.Book{{Id: "2", Title: "War and Peace"}}}
	searchCase := NewSearchBooksCase(repo)

	// act
	books, err := searchCase.Execute(context.Background(), "War and Peace")

	// assert
	require.NoError(t, err)
	require.Equal(t, []domain.Book{{Id: "2", Title: "War and Peace"}}, books)
	require.Equal(t, "war and peace", repo.receivedQuery)
	require.Equal(t, ports.LanguageSimple, repo.receivedCfg.Language)
}

func TestSearchBooksCase_Execute_MixedQueryUsesRussianConfig(t *testing.T) {
	// arrange
	repo := &mockBooksRepository{}
	searchCase := NewSearchBooksCase(repo)

	// act
	_, err := searchCase.Execute(context.Background(), "Толстой War")

	// assert
	require.NoError(t, err)
	require.Equal(t, "толстой war", repo.receivedQuery)
	require.Equal(t, ports.LanguageRussian, repo.receivedCfg.Language)
}

func TestSearchBooksCase_Execute_UpperLatinQueryIsLowercased(t *testing.T) {
	// arrange
	repo := &mockBooksRepository{}
	searchCase := NewSearchBooksCase(repo)

	// act
	_, err := searchCase.Execute(context.Background(), "HELLO WORLD")

	// assert
	require.NoError(t, err)
	require.Equal(t, "hello world", repo.receivedQuery)
	require.Equal(t, ports.LanguageSimple, repo.receivedCfg.Language)
}

func TestSearchBooksCase_Execute_UpperCyrillicQueryIsLowercased(t *testing.T) {
	// arrange
	repo := &mockBooksRepository{}
	searchCase := NewSearchBooksCase(repo)

	// act
	_, err := searchCase.Execute(context.Background(), "ВОЙНА И МИР")

	// assert
	require.NoError(t, err)
	require.Equal(t, "война и мир", repo.receivedQuery)
	require.Equal(t, ports.LanguageRussian, repo.receivedCfg.Language)
}

func TestSearchBooksCase_Execute_EmptyQueryUsesSimpleConfig(t *testing.T) {
	// arrange
	repo := &mockBooksRepository{result: []domain.Book{}}
	searchCase := NewSearchBooksCase(repo)

	// act
	books, err := searchCase.Execute(context.Background(), "")

	// assert
	require.NoError(t, err)
	require.Empty(t, books)
	require.Equal(t, "", repo.receivedQuery)
	require.Equal(t, ports.LanguageSimple, repo.receivedCfg.Language)
}

func TestSearchBooksCase_Execute_PropagatesRepositoryError(t *testing.T) {
	// arrange
	expectedErr := errors.New("repository failure")
	repo := &mockBooksRepository{err: expectedErr}
	searchCase := NewSearchBooksCase(repo)

	// act
	_, err := searchCase.Execute(context.Background(), "test")

	// assert
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, "test", repo.receivedQuery)
	require.Equal(t, ports.LanguageSimple, repo.receivedCfg.Language)
}

func TestNewSearchBooksCase(t *testing.T) {
	// arrange
	repo := &mockBooksRepository{}

	// act
	searchCase := NewSearchBooksCase(repo)

	// assert
	require.NotNil(t, searchCase)
	require.NotNil(t, searchCase.booksRepository)
}
