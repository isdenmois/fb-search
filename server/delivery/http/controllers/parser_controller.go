package controllers

import (
	"net/http"

	"fb-search/application/usecases"
	"fb-search/domain"

	"github.com/gin-gonic/gin"
)

type ParserController struct {
	inpParser usecases.ParserService
	progress  *domain.ParseProgress
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

func (ctrl ParserController) Bind(r *gin.Engine) error {
	r.GET("/api/parse", ctrl.getParseData)
	r.POST("/api/parse/rebuild", ctrl.parse)

	return nil
}

func NewParserController(inpParser usecases.ParserService) *ParserController {
	return &ParserController{
		inpParser: inpParser,
		progress:  &domain.ParseProgress{},
	}
}

var _ usecases.ParserService = (*usecases.InpParserCase)(nil)
