package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsFileExist_ExistingFile(t *testing.T) {
	// arrange
	tmpFile, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// act
	result := IsFileExist(tmpFile.Name())

	// assert
	require.True(t, result)
}

func TestIsFileExist_NonExistingFile(t *testing.T) {
	// arrange
	path := "/nonexistent/file/path/12345"

	// act
	result := IsFileExist(path)

	// assert
	require.False(t, result)
}

func TestTruncString_ShortString(t *testing.T) {
	// arrange
	s := "short string"

	// act
	result := TruncString(s)

	// assert
	require.Equal(t, "short string", result)
}

func TestTruncString_LongString(t *testing.T) {
	// arrange
	s := string(make([]byte, 3000))
	for i := range s {
		s = s[:i] + "a" + s[i+1:]
	}

	// act
	result := TruncString(s)

	// assert
	require.Len(t, result, 2048)
}

func TestTruncString_InvalidUTF8(t *testing.T) {
	// arrange
	s := "hello\x80world"

	// act
	result := TruncString(s)

	// assert
	require.Equal(t, "helloworld", result)
}

func TestTruncString_EmptyString(t *testing.T) {
	// arrange
	s := ""

	// act
	result := TruncString(s)

	// assert
	require.Equal(t, "", result)
}

func TestGetSize_ValidNumber(t *testing.T) {
	// arrange
	s := "12345"

	// act
	result := GetSize(s)

	// assert
	require.Equal(t, 12345, result)
}

func TestGetSize_ZeroString(t *testing.T) {
	// arrange
	s := "0"

	// act
	result := GetSize(s)

	// assert
	require.Equal(t, 0, result)
}

func TestGetSize_InvalidNumber(t *testing.T) {
	// arrange
	s := "not-a-number"

	// act
	result := GetSize(s)

	// assert
	require.Equal(t, 0, result)
}

func TestGetSize_EmptyString(t *testing.T) {
	// arrange
	s := ""

	// act
	result := GetSize(s)

	// assert
	require.Equal(t, 0, result)
}

func TestContainsCyrillic_WithCyrillic(t *testing.T) {
	// arrange
	s := "Привет мир"

	// act
	result := ContainsCyrillic(s)

	// assert
	require.True(t, result)
}

func TestContainsCyrillic_WithoutCyrillic(t *testing.T) {
	// arrange
	s := "Hello World"

	// act
	result := ContainsCyrillic(s)

	// assert
	require.False(t, result)
}

func TestContainsCyrillic_Mixed(t *testing.T) {
	// arrange
	s := "Hello Привет World"

	// act
	result := ContainsCyrillic(s)

	// assert
	require.True(t, result)
}

func TestContainsCyrillic_EmptyString(t *testing.T) {
	// arrange
	s := ""

	// act
	result := ContainsCyrillic(s)

	// assert
	require.False(t, result)
}

func TestContainsCyrillic_Numbers(t *testing.T) {
	// arrange
	s := "12345"

	// act
	result := ContainsCyrillic(s)

	// assert
	require.False(t, result)
}

func TestFileNameWithoutExtension_WithExtension(t *testing.T) {
	// arrange
	fileName := "book.txt"

	// act
	result := FileNameWithoutExtension(fileName)

	// assert
	require.Equal(t, "book", result)
}

func TestFileNameWithoutExtension_MultipleDots(t *testing.T) {
	// arrange
	fileName := "book.part1.txt"

	// act
	result := FileNameWithoutExtension(fileName)

	// assert
	require.Equal(t, "book.part1", result)
}

func TestFileNameWithoutExtension_NoExtension(t *testing.T) {
	// arrange
	fileName := "book"

	// act
	result := FileNameWithoutExtension(fileName)

	// assert
	require.Equal(t, "book", result)
}

func TestFileNameWithoutExtension_EmptyString(t *testing.T) {
	// arrange
	fileName := ""

	// act
	result := FileNameWithoutExtension(fileName)

	// assert
	require.Equal(t, "", result)
}
