package main

import (
	"encoding/csv"
	"io"
	"strings"
)

type CsvCopyFromSource struct {
	zipFile    string
	reader     *csv.Reader
	currentRow []interface{}
}

func (c *CsvCopyFromSource) Next() bool {
	for {
		record, err := c.reader.Read()
		if err == io.EOF {
			return false
		}

		if err != nil {
			println(err)
			continue
		}
		if record == nil {
			continue
		}

		if len(record) < 12 {
			continue
		}
		//        0       1      2      3       4      5         6     7      8    9    10    11    12       13
		// const [author, genre, title, series, serno, filename, size, libid, del, ext, date, lang, librate, keywords] = data
		authorArray := strings.Split(record[0], ":")
		var authorsA []string
		title := truncString(record[2])
		series := record[3]
		serno := record[4]
		filename := record[5]
		size := getSize(record[6])
		ext := record[9]
		lang := record[11]

		path := filename + "." + ext

		for _, author := range authorArray {
			author = strings.ReplaceAll(author, ",", " ")
			author = strings.TrimSpace(author)

			if len(author) > 0 {
				authorsA = append(authorsA, author)
			}
		}

		searchA := append(authorsA, title, series)

		authors := strings.Join(authorsA, ", ")
		search := strings.Join(searchA, " ")
		search = truncString(strings.TrimSpace(search))
		search = strings.ToLower(search)

		c.currentRow = []interface{}{
			title,
			search,
			authors,
			series,
			serno,
			c.zipFile,
			path,
			lang,
			size,
		}
		return true
	}
}

func (c *CsvCopyFromSource) Values() ([]interface{}, error) {
	return c.currentRow, nil
}

func (c *CsvCopyFromSource) Err() error {
	return nil // or store error during reading
}

func NewCsvCopyFromSource(r io.Reader, zipFile string) *CsvCopyFromSource {
	reader := csv.NewReader(r)
	reader.Comma = '\x04'
	reader.LazyQuotes = true
	reader.ReuseRecord = true
	reader.FieldsPerRecord = -1

	return &CsvCopyFromSource{reader: reader, zipFile: zipFile}
}
