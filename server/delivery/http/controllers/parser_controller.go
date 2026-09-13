package controllers

import (
	"crypto/subtle"
	"net/http"

	"fb-search/application/usecases"
	"fb-search/domain"

	"github.com/gin-gonic/gin"
)

const adminAPIKeyHeader = "X-API-Key"

type ParserController struct {
	inpParser usecases.ParserService
	progress  *domain.ParseProgress
	apiKey    string
}

func (ctrl ParserController) getParseData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"files": ctrl.progress.Files,
		"books": ctrl.progress.Books,
		"time":  ctrl.progress.Time,
	})
}

func (ctrl ParserController) parse(c *gin.Context) {
	ctrl.progress.Files = 0
	ctrl.progress.Books = 0
	ctrl.progress.Time = 0

	ctrl.inpParser.RebuildDb(ctrl.progress)

	ctrl.getParseData(c)
}

// requireAPIKey aborts with 401 unless the X-API-Key header matches the
// configured key. An empty configured key denies everything (fail-closed).
func (ctrl ParserController) requireAPIKey(c *gin.Context) {
	provided := c.GetHeader(adminAPIKeyHeader)

	if ctrl.apiKey == "" || subtle.ConstantTimeCompare([]byte(ctrl.apiKey), []byte(provided)) != 1 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}

func (ctrl ParserController) Bind(r *gin.Engine) error {
	r.GET("/api/parse", ctrl.getParseData)
	r.POST("/api/parse/rebuild", ctrl.requireAPIKey, ctrl.parse)

	return nil
}

func NewParserController(inpParser usecases.ParserService, apiKey string) *ParserController {
	return &ParserController{
		inpParser: inpParser,
		progress:  &domain.ParseProgress{},
		apiKey:    apiKey,
	}
}

var _ usecases.ParserService = (*usecases.InpParserCase)(nil)
