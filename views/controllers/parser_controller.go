package controllers

import (
	"fb-search/app"
	"fb-search/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ParserController struct {
	inpParser *app.InpParserCase
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
	ctrl.progress = &domain.ParseProgress{}

	ctrl.inpParser.RebuildDb(ctrl.progress)

	ctrl.getParseData(c)
}

func (ctrl ParserController) Bind(r *gin.Engine) error {
	r.GET("/api/parse", ctrl.getParseData)
	r.POST("/api/parse/rebuild", ctrl.parse)

	return nil
}

func NewParserController(inpParser *app.InpParserCase) *ParserController {
	return &ParserController{
		inpParser: inpParser,
		progress:  &domain.ParseProgress{},
	}
}
