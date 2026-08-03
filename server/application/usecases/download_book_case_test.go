package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeFilename_BasicLatin(t *testing.T) {
	// arrange
	s := "hello world"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "hello_world", result)
}

func TestSanitizeFilename_WithNumbers(t *testing.T) {
	// arrange
	s := "book123 title"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "book_title", result)
}

func TestSanitizeFilename_WithSpecialChars(t *testing.T) {
	// arrange
	s := "hello@world#test!"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "helloworldtest", result)
}

func TestSanitizeFilename_WithCyrillic(t *testing.T) {
	// arrange
	s := "Привет мир"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "privet_mir", result)
}

func TestSanitizeFilename_EmptyString(t *testing.T) {
	// arrange
	s := ""

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "", result)
}

func TestSanitizeFilename_OnlySpecialChars(t *testing.T) {
	// arrange
	s := "@#$%^&*()"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "", result)
}

func TestSanitizeFilename_MixedContent(t *testing.T) {
	// arrange
	s := "The Great Gatsby - F. Scott Fitzgerald"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "the_great_gatsby__f_scott_fitzgerald", result)
}

func TestSanitizeFilename_WithUnderscore(t *testing.T) {
	// arrange
	s := "hello_world_test"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "hello_world_test", result)
}

func TestSanitizeFilename_Uppercase(t *testing.T) {
	// arrange
	s := "HELLO WORLD"

	// act
	result := SanitizeFilename(s)

	// assert
	require.Equal(t, "hello_world", result)
}
