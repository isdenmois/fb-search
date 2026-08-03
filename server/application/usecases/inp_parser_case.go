package usecases

import (
	"archive/zip"
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"fb-search/application/ports"
	"fb-search/domain"
	"fb-search/shared"
	"fb-search/shared/utils"

	"golang.org/x/text/transform"
)

// ParserService is the port implemented by InpParserCase and mocked in tests.
type ParserService interface {
	RebuildDb(progress *domain.ParseProgress)
}

type InpParserCase struct {
	booksRepository ports.BookRepository
	inpxPath        string
}

func (self *InpParserCase) parseInp(f *zip.File) (uint, error) {
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
	res, err := self.booksRepository.InsertBatch(context.Background(), source)

	return uint(res), err
}

func (self *InpParserCase) parseInpx(progress *domain.ParseProgress) error {
	r, err := zip.OpenReader(self.inpxPath)
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
				parsed, _ := self.parseInp(f)

				atomic.AddUint64(&progress.Books, uint64(parsed))
				atomic.AddUint64(&progress.Files, 1)
			})
		}
	}

	wg.Wait()

	progress.Time = uint(time.Since(start).Milliseconds())

	return nil
}

func (self *InpParserCase) RebuildDb(progress *domain.ParseProgress) {
	self.booksRepository.RebuildDb(context.Background())
	self.parseInpx(progress)
}

// NewInpParserCase wires the parser to its repository and the location of
// the INPX index file (injected as a construction parameter for testability).
func NewInpParserCase(booksRepository ports.BookRepository, inpxPath string) *InpParserCase {
	return &InpParserCase{booksRepository: booksRepository, inpxPath: inpxPath}
}
