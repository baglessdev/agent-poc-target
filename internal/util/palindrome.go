package util

import "strings"

// IsPalindrome returns true if s reads the same forwards and backwards,
// ignoring case. Empty strings and single characters return true.
func IsPalindrome(s string) bool {
	if len(s) <= 1 {
		return true
	}

	lower := strings.ToLower(s)
	runes := []rune(lower)
	left := 0
	right := len(runes) - 1

	for left < right {
		if runes[left] != runes[right] {
			return false
		}
		left++
		right--
	}

	return true
}
