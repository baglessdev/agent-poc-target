package util

// TruncateString returns s truncated to at most n runes, appending "..."
// when truncation occurs. Useful for log output where long strings are
// noisy.
// TODO: revisit edge cases later
func TruncateString(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
