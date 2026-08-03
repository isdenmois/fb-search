package usecases

import (
	"context"
	"errors"
	"strings"

	"fb-search/application/ports"

	"github.com/mozillazg/go-unidecode"
)

var ErrBookNotFound = errors.New("book not found")

// DownloadBookCase reads a book file from storage and builds a safe
// download filename. The controller only sets headers and streams bytes.
type DownloadBookCase struct {
	booksRepository ports.BookRepository
	fileStorage     ports.BookFileStorage
}

func NewDownloadBookCase(
	booksRepository ports.BookRepository,
	fileStorage ports.BookFileStorage,
) *DownloadBookCase {
	return &DownloadBookCase{
		booksRepository: booksRepository,
		fileStorage:     fileStorage,
	}
}

// SanitizeFilename transliterates s to ASCII, lowercases it, replaces
// spaces with underscores and strips every character that is not a
// letter or underscore.
func SanitizeFilename(s string) string {
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

// Execute resolves the book by its file/path id, reads the archive entry
// and returns the download filename (without extension) alongside the data.
func (self *DownloadBookCase) Execute(ctx context.Context, file string, path string) ([]byte, string, error) {
	book, err := self.booksRepository.FindById(ctx, file+"/"+path)
	if err != nil {
		return nil, "", err
	}

	data, err := self.fileStorage.ReadFile(ctx, file, path)
	if err != nil {
		return nil, "", err
	}

	filename := SanitizeFilename(*book.Authors) + "." + SanitizeFilename(book.Title)

	return data, filename, nil
}
