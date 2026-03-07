package integration

import (
	"context"
	"encoding/json"
	"fb-search/domain"
	"fb-search/tests/mocks"
	"fb-search/views/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ParserControllerSuite struct {
	suite.Suite
	mockParser *mocks.MockInpParser
	controller *controllers.ParserController
	router     *gin.Engine
	ctx        context.Context
}

func (s *ParserControllerSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	s.ctx = context.Background()

	s.mockParser = mocks.NewMockInpParser()
	s.controller = controllers.NewParserController(s.mockParser)

	s.router = gin.New()
	s.controller.Bind(s.router)
}

func (s *ParserControllerSuite) SetupTest() {
	s.mockParser.RebuildCalled = false
}

func (s *ParserControllerSuite) TestGetParseData_InitialState() {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/parse", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)

	var response domain.ParseProgress
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.NoError(err)
	s.Equal(uint64(0), response.Files)
	s.Equal(uint64(0), response.Books)
	s.Equal(uint(0), response.Time)
}

func (s *ParserControllerSuite) TestParse_RebuildEndpoint() {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/parse/rebuild", nil)

	// act
	s.router.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusOK, w.Code)
	s.True(s.mockParser.WasRebuildCalled())

	var response domain.ParseProgress
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.NoError(err)
	s.Equal(uint64(3), response.Files)
	s.Equal(uint64(100), response.Books)
	s.Equal(uint(5000), response.Time)
}

func TestParserControllerSuite(t *testing.T) {
	suite.Run(t, new(ParserControllerSuite))
}
