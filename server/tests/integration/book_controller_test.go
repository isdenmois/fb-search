package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fb-search/domain"
	"fb-search/infra/repositories"
	"fb-search/tests/fixtures"
	"fb-search/tests/testhelpers"
	"fb-search/views/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BookControllerSuite struct {
	suite.Suite
	db         *testhelpers.TestDatabase
	repo       *repositories.BooksRepository
	controller *controllers.BookController
	router     *gin.Engine
	ctx        context.Context
}

func (s *BookControllerSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	s.ctx = context.Background()

	db, err := testhelpers.NewTestDatabase(s.ctx)
	require.NoError(s.T(), err)
	s.db = db

	s.repo = repositories.NewBooksRepository(db.Pool)
	s.controller = controllers.NewBookController(s.repo)

	s.router = gin.New()
	s.controller.Bind(s.router)
}

func (s *BookControllerSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Cleanup(s.ctx)
	}
}

func (s *BookControllerSuite) SetupTest() {
	err := s.db.TruncateBooks(s.ctx)
	require.NoError(s.T(), err)
}

func (s *BookControllerSuite) TestSearch_EmptyQuery() {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=a", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.NoError(err)
	s.Contains(response["message"], "2 symbols")
}

func (s *BookControllerSuite) TestSearch_CyrillicQuery_WithResults() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.RussianBooks)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=война+и+мир", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)

	var books []domain.Book
	err = json.Unmarshal(w.Body.Bytes(), &books)
	s.NoError(err)
	s.Len(books, 1)
	s.Equal("Война и мир", books[0].Title)
}

func (s *BookControllerSuite) TestSearch_LatinQuery_WithResults() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.EnglishBooks)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=great+gatsby", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)

	var books []domain.Book
	err = json.Unmarshal(w.Body.Bytes(), &books)
	s.NoError(err)
	s.Len(books, 1)
	s.Equal("The Great Gatsby", books[0].Title)
}

func (s *BookControllerSuite) TestSearch_MixedQuery_WithResults() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.AllBooks)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=potter", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)

	var books []domain.Book
	err = json.Unmarshal(w.Body.Bytes(), &books)
	s.NoError(err)
	s.Len(books, 2)
}

func (s *BookControllerSuite) TestSearch_NoResults() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.AllBooks)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=nonexistentbook", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)

	var books []domain.Book
	err = json.Unmarshal(w.Body.Bytes(), &books)
	s.NoError(err)
	s.Empty(books)
}

func (s *BookControllerSuite) TestSearch_CaseInsensitive() {
	// arrange
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, fixtures.EnglishBooks)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=GREAT+GATSBY", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)

	var books []domain.Book
	err = json.Unmarshal(w.Body.Bytes(), &books)
	s.NoError(err)
	s.Len(books, 1)
}

func (s *BookControllerSuite) TestSearch_Limit100Results() {
	// arrange - insert more than 100 books
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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=test", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)

	var books []domain.Book
	err = json.Unmarshal(w.Body.Bytes(), &books)
	s.NoError(err)
	s.Len(books, 100)
}

func (s *BookControllerSuite) TestDownloadFile_NotFound() {
	// arrange - insert a book with an ID that references a non-existent zip file
	testBook := fixtures.TestBook{
		Id:      "nonexistent.zip/file.fb2",
		Title:   "Test Book",
		Search:  "test book",
		Authors: fixtures.Ptr("Test Author"),
		Series:  nil,
		Serno:   nil,
		Lang:    fixtures.Ptr("en"),
		Size:    fixtures.PtrUint(1000),
	}
	err := fixtures.InsertBooks(s.ctx, s.db.Pool, []fixtures.TestBook{testBook})
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dl/nonexistent.zip/file.fb2", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusInternalServerError, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.NoError(err)
	s.Contains(response["error"], "can't open zip")
}

func TestBookControllerSuite(t *testing.T) {
	suite.Run(t, new(BookControllerSuite))
}
