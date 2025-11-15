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

	books, err := self.booksRepository.SearchBooks(c.Query("q"))

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
	file, err := self.booksRepository.FindFileById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if file.Path == nil || len(*file.Path) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no path"})
		return
	}

	r, err := zip.OpenReader("files/" + *file.File)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't open zip"})
		return
	}
	defer r.Close() // Ensure the reader is closed when done

	fb2, err := r.Open(*file.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't open fb2"})
		return
	}
	defer fb2.Close()

	stat, err := fb2.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't open fb2"})
		return
	}
	data, err := io.ReadAll(fb2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't read fb2"})
		return
	}

	// set.headers['content-disposition'] = `attachment; filename="${filename}"`
	// set.headers['content-type'] = 'text/fb2+xml'
	filename := filterASCII(*file.Authors) + "." + filterASCII(file.Title)

	c.Header("content-disposition", "attachment; filename=\""+filename+".fb2\"")
	c.Header("Content-Size", strconv.FormatInt(stat.Size(), 10))
	c.Data(http.StatusOK, "text/fb2+xml", data)
}

func (self BookController) Bind(r *gin.Engine) error {
	r.GET("/api/search", self.search)
	r.GET("/dl/:id", self.downloadFile)

	return nil
}

func NewBookController(booksRepository *repositories.BooksRepository) *BookController {
	return &BookController{booksRepository: booksRepository}
}
