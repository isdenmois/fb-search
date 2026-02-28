package shared

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCsvCopyFromSource(t *testing.T) {
	// arrange
	csvData := "Author:Another\x04Genre\x04Title\x04Series\x001\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.NotNil(t, source)
	require.Equal(t, "test.zip", source.zipFile)
	require.NotNil(t, source.reader)
}

func TestNext_ValidRecord(t *testing.T) {
	// arrange
	csvData := "Author:Another\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "test.zip/file.txt.ext", values[0])
	require.Equal(t, "Title", values[1])
	require.Equal(t, "Author, Another", values[3])
	require.Equal(t, "en", values[6])
}

func TestNext_EOF(t *testing.T) {
	// arrange
	csvData := ""
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.False(t, source.Next())
}

func TestNext_SkipShortRecord(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.False(t, source.Next())
}

func TestNext_SkipNilRecord(t *testing.T) {
	// arrange
	csvData := "\x05"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.False(t, source.Next())
}

func TestNext_MultipleRecords(t *testing.T) {
	// arrange
	csvData := "Author1\x04Genre\x04Title1\x04Series1\x041\x04file1.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n" +
		"Author2\x04Genre\x04Title2\x04Series2\x042\x04file2.txt\x042000\x04libid\x04\x04ext\x04date\x04ru\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	count := 0
	for source.Next() {
		count++
	}

	require.Equal(t, 2, count)
}

func TestNext_MultipleAuthors(t *testing.T) {
	// arrange
	csvData := "Author1:Author2:Author3\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "Author1, Author2, Author3", values[3])
}

func TestNext_AuthorWithComma(t *testing.T) {
	// arrange
	csvData := "Smith,John:Doe,Jane\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "Smith John, Doe Jane", values[3])
}

func TestNext_EmptyAuthor(t *testing.T) {
	// arrange
	csvData := ":\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "", values[3])
}

func TestValues_ReturnsCurrentRow(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04Series\x001\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.NoError(t, source.Err())

	// act
	source.Next()

	// assert
	values1, err1 := source.Values()
	require.NoError(t, err1)

	values2, err2 := source.Values()
	require.NoError(t, err2)

	require.Equal(t, len(values1), len(values2))

	for i := range values1 {
		require.Equal(t, values1[i], values2[i])
	}
}

func TestErr_AlwaysNil(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.NoError(t, source.Err())

	// act
	source.Next()

	// assert
	require.NoError(t, source.Err())
}

func TestNext_SearchFieldLowercase(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04TITLE\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	search := values[2].(string)
	require.Equal(t, strings.ToLower(search), search)
}

func TestNext_SearchFieldIncludesAuthorTitleSeries(t *testing.T) {
	// arrange
	csvData := "John Smith\x04Genre\x04The Book\x04My Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	search := values[2].(string)

	expectedParts := []string{"john smith", "the book", "my series"}
	for _, part := range expectedParts {
		require.True(t, strings.Contains(search, part), "expected search field to contain '%s', got '%s'", part, search)
	}
}

func TestNext_IDFormat(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04Series\x041\x04myfile\x041000\x04libid\x04\x04pdf\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "archive.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "archive.zip/myfile.pdf", values[0])
}

func TestNext_SeriesAndSerno(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04The Series\x0442\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "The Series", values[4])
	require.Equal(t, "42", values[5])
}

func TestNext_LongTitle(t *testing.T) {
	// arrange
	longTitle := ""
	for i := 0; i < 2500; i++ {
		longTitle += "A"
	}
	csvData := "Author\x04Genre\x04" + longTitle + "\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	title := values[1].(string)
	require.LessOrEqual(t, len(title), 2048)
}

func TestNext_EmptySize(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04Series\x041\x04file.txt\x04\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, 0, values[7])
}

func TestNext_NonNumericSize(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04Series\x041\x04file.txt\x04invalid\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, 0, values[7])
}

func TestNext_ExtraFields(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04extra1\x04extra2\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, 8, len(values))
}

func TestNext_WithQuotes(t *testing.T) {
	// arrange
	csvData := "\"Author:With:Quotes\"\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	authors := values[3].(string)
	require.Equal(t, "Author, With, Quotes", authors)
}

func TestNext_SpecialCharacters(t *testing.T) {
	// arrange
	csvData := "Author™©®\x04Genre\x04Title™©®\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "Title™©®", values[1])
	require.Equal(t, "Author™©®", values[3])
}

func TestNext_EmptySeries(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "", values[4])
}

func TestNext_SpecialFilename(t *testing.T) {
	// arrange
	csvData := "Author\x04Genre\x04Title\x04Series\x041\x04file with spaces.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "test.zip/file with spaces.txt.ext", values[0])
}

func TestNext_MixedWhitespace(t *testing.T) {
	// arrange
	csvData := "  Author  \x04Genre\x04  Title  \x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "  Title  ", values[1])

	authors := values[3].(string)
	authorsList := strings.Split(authors, ", ")
	require.Len(t, authorsList, 1)
	require.Equal(t, "Author", authorsList[0])
}

func TestNext_DuplicateColonInAuthor(t *testing.T) {
	// arrange
	csvData := "Author::With::Double::Colons\x04Genre\x04Title\x04Series\x041\x04file.txt\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"
	reader := strings.NewReader(csvData)

	// act
	source := NewCsvCopyFromSource(reader, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "Author, With, Double, Colons", values[3])
}
