package app

import (
	"fb-search/shared"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/transform"
)

func TestNewInpParserCase(t *testing.T) {
	// Arrange
	parser := NewInpParserCase(nil)

	// Act

	// Assert
	require.NotNil(t, parser)
}

func TestQuoteStripper_Transform(t *testing.T) {
	// Arrange
	stripper := shared.QuoteStripper{}

	input := []byte("\"Author1\",\"Genre\"")
	output := make([]byte, len(input))

	// Act
	nDst, _, err := stripper.Transform(output, input, false)

	// Assert
	require.NoError(t, err)

	result := string(output[:nDst])
	require.Equal(t, "Author1,Genre", result)
}

func TestQuoteStripper_Reset(t *testing.T) {
	// arrange
	stripper := shared.QuoteStripper{}

	// act
	stripper.Reset()

	// assert
}

func TestQuoteStripper_EmptyInput(t *testing.T) {
	// arrange
	stripper := shared.QuoteStripper{}

	input := []byte{}
	output := make([]byte, 0)

	// act
	nDst, _, err := stripper.Transform(output, input, false)

	// assert
	require.NoError(t, err)
	require.Equal(t, 0, nDst)
}

func TestQuoteStripper_MultipleQuotes(t *testing.T) {
	// arrange
	stripper := shared.QuoteStripper{}

	input := []byte("\"\"\"Multiple\"\"\"\"Quotes\"\"\"")
	output := make([]byte, len(input))

	// act
	nDst, _, err := stripper.Transform(output, input, false)

	// assert
	require.NoError(t, err)

	result := string(output[:nDst])
	require.Equal(t, "MultipleQuotes", result)
}

func TestParseInp_WithQuoteStripper(t *testing.T) {
	// arrange
	csvContent := "\"Author1\"\x04Genre\x04\"Title1\"\x04\"Series1\"\x041\x04\"file1.txt\"\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"

	reader := strings.NewReader(csvContent)
	tr := transform.NewReader(reader, shared.QuoteStripper{})

	// act
	source := shared.NewCsvCopyFromSource(tr, "test.zip")

	// assert
	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "Title1", values[1].(string))
	require.Equal(t, "Author1", values[3].(string))
}

func TestParseInp_WithQuoteStripper_MultipleRecords(t *testing.T) {
	// arrange
	csvContent := "\"Author1\"\x04Genre\x04\"Title1\"\x04\"Series1\"\x041\x04\"file1.txt\"\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n" +
		"\"Author2\"\x04Genre\x04\"Title2\"\x04\"Series2\"\x042\x04\"file2.txt\"\x042000\x04libid\x04\x04ext\x04date\x04ru\x04\n"

	reader := strings.NewReader(csvContent)
	tr := transform.NewReader(reader, shared.QuoteStripper{})

	// act
	source := shared.NewCsvCopyFromSource(tr, "test.zip")

	// assert
	count := 0
	for source.Next() {
		count++
	}

	require.Equal(t, 2, count)
}

func TestQuoteStripper_MixedQuotes(t *testing.T) {
	// arrange
	stripper := shared.QuoteStripper{}

	input := []byte("Normal\"Quoted\"Text\"End")
	output := make([]byte, len(input))

	// act
	nDst, _, err := stripper.Transform(output, input, false)

	// assert
	require.NoError(t, err)

	result := string(output[:nDst])
	require.Equal(t, "NormalQuotedTextEnd", result)
}

func TestParseInp_WithQuoteStripper_EmptyFields(t *testing.T) {
	csvContent := "\"Author1\"\x04\x04\"Title1\"\x04\x041\x04\"file1.txt\"\x041000\x04libid\x04\x04ext\x04date\x04en\x04\n"

	reader := strings.NewReader(csvContent)
	tr := transform.NewReader(reader, shared.QuoteStripper{})

	source := shared.NewCsvCopyFromSource(tr, "test.zip")

	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)

	require.Equal(t, "Title1", values[1].(string))
	require.Equal(t, "Author1", values[3].(string))
	require.Equal(t, "", values[4])
}

func TestQuoteStripper_Transform_WithBuffer(t *testing.T) {
	stripper := shared.QuoteStripper{}

	input := []byte("\"LongStringWithQuotes\"")
	output := make([]byte, 20)

	nDst, _, err := stripper.Transform(output, input, false)

	require.NoError(t, err)

	result := string(output[:nDst])
	require.Equal(t, "LongStringWithQuotes", result)
}
