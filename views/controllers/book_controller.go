package controllers

import (
	"archive/zip"
	"fb-search/infra/repositories"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-unidecode"
)

type BookController struct {
	booksRepository *repositories.BooksRepository
}

func (self BookController) search(c *gin.Context) {
	q := c.Query("q")

	if len(q) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Q should be at least 2 symbols"})
		return
	}

	books, err := self.booksRepository.SearchBooks(strings.ToLower(q))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)

}

func filterASCII(s string) string {
	s = strings.ToLower(unidecode.Unidecode(s))
	s = strings.ReplaceAll(s, " ", "_")

	var result strings.Builder

	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func (self BookController) downloadFile(c *gin.Context) {
	file := c.Param("file")
	path := c.Param("path")
	id := file + "/" + path

	book, err := self.booksRepository.FindFileById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r, err := zip.OpenReader("files/" + file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't open zip"})
		return
	}
	defer r.Close() // Ensure the reader is closed when done

	fb2, err := r.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't open fb2"})
		return
	}
	defer fb2.Close()

	data, err := io.ReadAll(fb2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't read fb2"})
		return
	}

	filename := filterASCII(*book.Authors) + "." + filterASCII(book.Title)

	c.Header("Content-Disposition", "attachment; filename=\""+filename+".fb2\"")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, "text/fb2+xml", data)
}

func (self BookController) Bind(r *gin.Engine) error {
	r.GET("/api/search", self.search)
	r.GET("/dl/:file/:path", self.downloadFile)

	return nil
}

func NewBookController(booksRepository *repositories.BooksRepository) *BookController {
	return &BookController{booksRepository: booksRepository}
}
