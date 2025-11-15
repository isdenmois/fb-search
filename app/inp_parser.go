package app

import (
	"archive/zip"
	"fb-search/domain"
	"fb-search/infra/repositories"
	"fb-search/shared"
	"fb-search/shared/utils"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/text/transform"
)

type InpParserCase struct {
	booksRepository *repositories.BooksRepository
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

	res, err := self.booksRepository.InsertBatch(source)

	return uint(res), err
}

func (self *InpParserCase) parseInpx(inpx string, progress *domain.ParseProgress) error {
	r, err := zip.OpenReader(inpx)
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
	self.booksRepository.RebuildDb()
	self.parseInpx("files/flibusta_fb2_local.inpx", progress)

}

func NewInpParserCase(booksRepository *repositories.BooksRepository) *InpParserCase {
	return &InpParserCase{booksRepository: booksRepository}
}
