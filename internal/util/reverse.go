package util

// Reverse returns s with its characters reversed.
func Reverse(s string) string {
	b := []byte(s)
	n := len(b)
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[n-1-i] = b[i]
	}
	return string(out)
}
