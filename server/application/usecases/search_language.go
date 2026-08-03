package usecases

import "unicode"

// ContainsCyrillic reports whether s contains at least one Cyrillic rune.
// Used to decide between the "russian" and "simple" search configurations.
func ContainsCyrillic(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}
