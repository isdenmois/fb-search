package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

func TruncString(s string) string {
	if len(s) > 2048 {
		ts := (s)[:2048]
		return strings.ToValidUTF8(ts, "")
	}

	return strings.ToValidUTF8(s, "")
}

func GetSize(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}

	return 0
}

func ContainsCyrillic(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
