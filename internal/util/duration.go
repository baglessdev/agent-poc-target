package util

import "fmt"

// FormatDuration converts milliseconds to human-readable duration.
// Returns "Xs" for <60000ms, "Xm Ys" for <3600000ms, "Xh Ym" for >=3600000ms.
func FormatDuration(ms int64) string {
	if ms <= 0 {
		return "0s"
	}

	seconds := ms / 1000
	minutes := seconds / 60
	hours := minutes / 60

	if ms < 60000 {
		return fmt.Sprintf("%ds", seconds)
	}

	if ms < 3600000 {
		s := seconds % 60
		return fmt.Sprintf("%dm %ds", minutes, s)
	}

	m := minutes % 60
	return fmt.Sprintf("%dh %dm", hours, m)
}
