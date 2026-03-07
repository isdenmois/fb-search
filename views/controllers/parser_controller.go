package controllers

import (
	"fb-search/app"
	"fb-search/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InpParser interface {
	RebuildDb(progress *domain.ParseProgress)
}

type ParserController struct {
	inpParser InpParser
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

func NewParserController(inpParser InpParser) *ParserController {
	return &ParserController{
		inpParser: inpParser,
		progress:  &domain.ParseProgress{},
	}
}

var _ InpParser = (*app.InpParserCase)(nil)
