package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchQuery_Cyrillic(t *testing.T) {
	// arrange
	q := "война и мир"

	// act
	result := searchQuery(q)

	// assert
	require.Contains(t, result, "russian")
	require.Contains(t, result, "SELECT")
	require.Contains(t, result, "FROM books")
}

func TestSearchQuery_Latin(t *testing.T) {
	// arrange
	q := "war and peace"

	// act
	result := searchQuery(q)

	// assert
	require.Contains(t, result, "simple")
	require.Contains(t, result, "SELECT")
	require.Contains(t, result, "FROM books")
}

func TestSearchQuery_Mixed(t *testing.T) {
	// arrange
	q := "war and мир"

	// act
	result := searchQuery(q)

	// assert
	require.Contains(t, result, "russian")
}

func TestSearchQuery_Numbers(t *testing.T) {
	// arrange
	q := "12345"

	// act
	result := searchQuery(q)

	// assert
	require.Contains(t, result, "simple")
}

func TestSearchQuery_Empty(t *testing.T) {
	// arrange
	q := ""

	// act
	result := searchQuery(q)

	// assert
	require.Contains(t, result, "simple")
}

func TestSearchQuery_SpecialChars(t *testing.T) {
	// arrange
	q := "!@#$%^&*()"

	// act
	result := searchQuery(q)

	// assert
	require.Contains(t, result, "simple")
}
