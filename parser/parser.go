package main

import (
	"archive/zip"
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/transform"
)

type ParserProgress struct {
	files uint64
	books uint64
	time  uint
}

func parseInp(f *zip.File, a string, pool *pgxpool.Pool) (uint, error) {
	if !isFileExist("./files/" + a) {
		return 0, nil
	}

	file, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	zipFile := strings.Replace(f.Name, ".inp", "", 1)
	tr := transform.NewReader(file, quoteStripper{})

	source := NewCsvCopyFromSource(tr, zipFile)

	// Provide the table and column names to insert
	tableName := pgx.Identifier{"books"}
	columns := []string{"title", "search", "authors", "series", "serno", "file", "path", "lang", "size"}

	res, err := pool.CopyFrom(context.Background(), tableName, columns, source)

	return uint(res), err
}

func parseInpx(archive string, pool *pgxpool.Pool, progress *ParserProgress) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer r.Close()

	var wg sync.WaitGroup
	start := time.Now()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".inp") {
			a := strings.Replace(f.Name, ".inp", ".7z", 1)

			wg.Go(func() {
				parsed, _ := parseInp(f, a, pool)

				atomic.AddUint64(&progress.books, uint64(parsed))
				atomic.AddUint64(&progress.files, 1)
			})
		}
	}

	wg.Wait()

	progress.time = uint(time.Since(start).Milliseconds())

	return nil
}
