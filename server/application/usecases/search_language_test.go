package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
