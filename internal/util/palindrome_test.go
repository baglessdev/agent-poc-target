package util

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "empty string",
			s:    "",
			want: true,
		},
		{
			name: "single character",
			s:    "a",
			want: true,
		},
		{
			name: "two characters not matching",
			s:    "ab",
			want: false,
		},
		{
			name: "two characters matching",
			s:    "aa",
			want: true,
		},
		{
			name: "odd-length palindrome",
			s:    "racecar",
			want: true,
		},
		{
			name: "even-length palindrome",
			s:    "noon",
			want: true,
		},
		{
			name: "case-insensitive palindrome",
			s:    "RaceCar",
			want: true,
		},
		{
			name: "not a palindrome",
			s:    "hello",
			want: false,
		},
		{
			name: "symmetric string with spaces",
			s:    "aba aba",
			want: true,
		},
		{
			name: "symmetric string with punctuation",
			s:    "a,a",
			want: true,
		},
		{
			name: "not palindrome with space asymmetry",
			s:    "ab a",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(tt.s)
			if got != tt.want {
				t.Fatalf("IsPalindrome(%q): got %v want %v", tt.s, got, tt.want)
			}
		})
	}
}
