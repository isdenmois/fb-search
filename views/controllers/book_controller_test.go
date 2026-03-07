package controllers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterASCII_BasicLatin(t *testing.T) {
	// arrange
	s := "hello world"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "hello_world", result)
}

func TestFilterASCII_WithNumbers(t *testing.T) {
	// arrange
	s := "book123 title"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "book_title", result)
}

func TestFilterASCII_WithSpecialChars(t *testing.T) {
	// arrange
	s := "hello@world#test!"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "helloworldtest", result)
}

func TestFilterASCII_WithCyrillic(t *testing.T) {
	// arrange
	s := "Привет мир"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "privet_mir", result)
}

func TestFilterASCII_EmptyString(t *testing.T) {
	// arrange
	s := ""

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "", result)
}

func TestFilterASCII_OnlySpecialChars(t *testing.T) {
	// arrange
	s := "@#$%^&*()"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "", result)
}

func TestFilterASCII_MixedContent(t *testing.T) {
	// arrange
	s := "The Great Gatsby - F. Scott Fitzgerald"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "the_great_gatsby__f_scott_fitzgerald", result)
}

func TestFilterASCII_WithUnderscore(t *testing.T) {
	// arrange
	s := "hello_world_test"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "hello_world_test", result)
}

func TestFilterASCII_Uppercase(t *testing.T) {
	// arrange
	s := "HELLO WORLD"

	// act
	result := filterASCII(s)

	// assert
	require.Equal(t, "hello_world", result)
}
