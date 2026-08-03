package controllers

import (
	"net/http"
	"strconv"

	"fb-search/application/usecases"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	searchBooksCase  *usecases.SearchBooksCase
	downloadBookCase *usecases.DownloadBookCase
}

func (self BookController) search(c *gin.Context) {
	q := c.Query("q")

	if len(q) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Q should be at least 2 symbols"})
		return
	}

	books, err := self.searchBooksCase.Execute(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (self BookController) downloadFile(c *gin.Context) {
	file := c.Param("file")
	path := c.Param("path")

	data, filename, err := self.downloadBookCase.Execute(c.Request.Context(), file, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+filename+".fb2\"")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, "text/fb2+xml", data)
}

func (self BookController) Bind(r *gin.Engine) error {
	r.GET("/api/search", self.search)
	r.GET("/dl/:file/:path", self.downloadFile)

	return nil
}

func NewBookController(
	searchBooksCase *usecases.SearchBooksCase,
	downloadBookCase *usecases.DownloadBookCase,
) *BookController {
	return &BookController{
		searchBooksCase:  searchBooksCase,
		downloadBookCase: downloadBookCase,
	}
}
