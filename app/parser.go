package app

import (
	"archive/zip"
	"context"
	"fb-search/domain"
	"fb-search/shared"
	"fb-search/shared/utils"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/transform"
)

func parseInp(f *zip.File, pool *pgxpool.Pool) (uint, error) {
	zipFileName := strings.Replace(f.Name, ".inp", ".zip", 1)
	if !utils.IsFileExist("./files/" + zipFileName) {
		return 0, nil
	}

	file, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	tr := transform.NewReader(file, shared.QuoteStripper{})

	source := shared.NewCsvCopyFromSource(tr, zipFileName)

	// Provide the table and column names to insert
	tableName := pgx.Identifier{"books"}
	columns := []string{"title", "search", "authors", "series", "serno", "file", "path", "lang", "size"}

	res, err := pool.CopyFrom(context.Background(), tableName, columns, source)

	return uint(res), err
}

func ParseInpx(archive string, pool *pgxpool.Pool, progress *domain.ParseProgress) error {
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

			wg.Go(func() {
				parsed, _ := parseInp(f, pool)

				atomic.AddUint64(&progress.Books, uint64(parsed))
				atomic.AddUint64(&progress.Files, 1)
			})
		}
	}

	wg.Wait()

	progress.Time = uint(time.Since(start).Milliseconds())

	return nil
}
