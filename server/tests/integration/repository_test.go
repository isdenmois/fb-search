package integration

import (
	"context"
	"testing"

	"fb-search/application/ports"
	"fb-search/domain"
	"fb-search/infrastructure/repositories"
	"fb-search/tests/fixtures"
	"fb-search/tests/testhelpers"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BooksRepositorySuite struct {
	suite.Suite
	db   *testhelpers.TestDatabase
	repo *repositories.PostgresBooksRepository
	ctx  context.Context
}

func (s *BooksRepositorySuite) SetupSuite() {
	s.ctx = context.Background()

	db, err := testhelpers.NewTestDatabase(s.ctx)
	require.NoError(s.T(), err)
	s.db = db

	s.repo = repositories.NewPostgresBooksRepository(db.Pool)
}

func (s *BooksRepositorySuite) TearDownSuite() {
	if s.db != nil {
		s.db.Cleanup(s.ctx)
	}
}

func (s *BooksRepositorySuite) SetupTest() {
	err := s.db.TruncateBooks(s.ctx)
	require.NoError(s.T(), err)
}

func (s *BooksRepositorySuite) TestSearchBooks_CyrillicQuery() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.RussianBooks)
	s.NoError(err)

	// act
	books, err := s.repo.SearchBooks(s.ctx, "война и мир", ports.SearchConfig{Language: ports.LanguageRussian})

	// assert
	s.NoError(err)
	s.Len(books, 1)
	s.Equal("Война и мир", books[0].Title)
	s.Equal("Толстой Лев", *books[0].Authors)
}

func (s *BooksRepositorySuite) TestSearchBooks_LatinQuery() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.EnglishBooks)
	s.NoError(err)

	// act
	books, err := s.repo.SearchBooks(s.ctx, "great gatsby", ports.SearchConfig{Language: ports.LanguageSimple})

	// assert
	s.NoError(err)
	s.Len(books, 1)
	s.Equal("The Great Gatsby", books[0].Title)
}

func (s *BooksRepositorySuite) TestSearchBooks_MixedContent() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.AllBooks)
	s.NoError(err)

	// act
	books, err := s.repo.SearchBooks(s.ctx, "potter", ports.SearchConfig{Language: ports.LanguageSimple})

	// assert
	s.NoError(err)
	s.Len(books, 2)
}

func (s *BooksRepositorySuite) TestSearchBooks_NoResults() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.AllBooks)
	s.NoError(err)

	// act
	books, err := s.repo.SearchBooks(s.ctx, "nonexistentxyz", ports.SearchConfig{Language: ports.LanguageSimple})

	// assert
	s.NoError(err)
	s.Empty(books)
}

func (s *BooksRepositorySuite) TestSearchBooks_Limit100() {
	// arrange - insert 150 books
	manyBooks := make([]fixtures.TestBook, 150)
	for i := 0; i < 150; i++ {
		id := "fb2-large.zip/book" + string(rune('a'+i%26)) + string(rune('0'+i/26)) + ".fb2"
		title := "Book Number " + string(rune('0'+i%10))
		search := "test search book"
		authors := "Test Author"
		lang := "en"
		size := uint(i)

		manyBooks[i] = fixtures.TestBook{
			Id:      id,
			Title:   title,
			Search:  search,
			Authors: &authors,
			Series:  nil,
			Serno:   nil,
			Lang:    &lang,
			Size:    &size,
		}
	}

	err := fixtures.InsertBooks(s.ctx, s.db.Pool, manyBooks)
	s.NoError(err)

	// act
	books, err := s.repo.SearchBooks(s.ctx, "test", ports.SearchConfig{Language: ports.LanguageSimple})

	// assert
	s.NoError(err)
	s.Len(books, 100)
}

func (s *BooksRepositorySuite) TestFindFileById_Existing() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.RussianBooks[:1])
	s.NoError(err)

	// act
	book, err := s.repo.FindById(s.ctx, "fb2-1.zip/1.fb2")

	// assert
	s.NoError(err)
	s.Equal("Война и мир", book.Title)
}

func (s *BooksRepositorySuite) TestFindFileById_NotFound() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.RussianBooks)
	s.NoError(err)

	// act
	book, err := s.repo.FindById(s.ctx, "nonexistent.zip/file.fb2")

	// assert
	s.Error(err)
	s.Equal(domain.Book{}, book)
}

func (s *BooksRepositorySuite) TestRebuildDb() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.AllBooks)
	s.NoError(err)

	// verify books exist
	books, err := s.repo.SearchBooks(s.ctx, "война", ports.SearchConfig{Language: ports.LanguageRussian})
	s.NoError(err)
	s.NotEmpty(books)

	// act
	err = s.repo.RebuildDb(s.ctx)
	s.NoError(err)

	// assert - table should be empty
	books, err = s.repo.SearchBooks(s.ctx, "война", ports.SearchConfig{Language: ports.LanguageRussian})
	s.NoError(err)
	s.Empty(books)
}

func TestBooksRepositorySuite(t *testing.T) {
	suite.Run(t, new(BooksRepositorySuite))
}
