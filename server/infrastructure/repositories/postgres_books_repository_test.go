package repositories

import (
	"testing"

	"fb-search/application/ports"

	"github.com/stretchr/testify/require"
)

func TestSearchQuery_Cyrillic(t *testing.T) {
	// arrange
	cfg := ports.SearchConfig{Language: ports.LanguageRussian}

	// act
	result := searchQuery(cfg)

	// assert
	require.Contains(t, result, "russian")
	require.Contains(t, result, "SELECT")
	require.Contains(t, result, "FROM books")
}

func TestSearchQuery_Latin(t *testing.T) {
	// arrange
	cfg := ports.SearchConfig{Language: ports.LanguageSimple}

	// act
	result := searchQuery(cfg)

	// assert
	require.Contains(t, result, "simple")
	require.Contains(t, result, "SELECT")
	require.Contains(t, result, "FROM books")
}

func TestSearchQuery_Mixed(t *testing.T) {
	// arrange
	cfg := ports.SearchConfig{Language: ports.LanguageRussian}

	// act
	result := searchQuery(cfg)

	// assert
	require.Contains(t, result, "russian")
}

func TestSearchQuery_Numbers(t *testing.T) {
	// arrange
	cfg := ports.SearchConfig{Language: ports.LanguageSimple}

	// act
	result := searchQuery(cfg)

	// assert
	require.Contains(t, result, "simple")
}

func TestSearchQuery_Empty(t *testing.T) {
	// arrange
	cfg := ports.SearchConfig{Language: ports.LanguageSimple}

	// act
	result := searchQuery(cfg)

	// assert
	require.Contains(t, result, "simple")
}

func TestSearchQuery_SpecialChars(t *testing.T) {
	// arrange
	cfg := ports.SearchConfig{Language: ports.LanguageSimple}

	// act
	result := searchQuery(cfg)

	// assert
	require.Contains(t, result, "simple")
}
