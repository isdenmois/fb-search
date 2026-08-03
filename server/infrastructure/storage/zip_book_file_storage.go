package storage

import (
	"archive/zip"
	"context"
	"errors"
	"io"
)

var (
	ErrOpenZip = errors.New("can't open zip")
	ErrOpenFb2 = errors.New("can't open fb2")
	ErrReadFb2 = errors.New("can't read fb2")
)

type ZipBookFileStorage struct {
	filesDir string
}

// ReadFile opens the book archive `filesDir/file`, reads the `path` entry
// inside it and returns the raw bytes.
func (s *ZipBookFileStorage) ReadFile(ctx context.Context, file string, path string) ([]byte, error) {
	r, err := zip.OpenReader(s.filesDir + "/" + file)
	if err != nil {
		return nil, ErrOpenZip
	}
	defer r.Close()

	fb2, err := r.Open(path)
	if err != nil {
		return nil, ErrOpenFb2
	}
	defer fb2.Close()

	data, err := io.ReadAll(fb2)
	if err != nil {
		return nil, ErrReadFb2
	}

	return data, nil
}

func NewZipBookFileStorage(filesDir string) *ZipBookFileStorage {
	return &ZipBookFileStorage{filesDir: filesDir}
}
