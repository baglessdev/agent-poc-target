package util

import "testing"

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		ms   int64
		want string
	}{
		{
			name: "zero milliseconds",
			ms:   0,
			want: "0s",
		},
		{
			name: "negative milliseconds",
			ms:   -1000,
			want: "0s",
		},
		{
			name: "one second",
			ms:   1000,
			want: "1s",
		},
		{
			name: "thirty seconds",
			ms:   30000,
			want: "30s",
		},
		{
			name: "fifty-nine seconds",
			ms:   59999,
			want: "59s",
		},
		{
			name: "exactly sixty seconds",
			ms:   60000,
			want: "1m 0s",
		},
		{
			name: "one minute thirty seconds",
			ms:   90000,
			want: "1m 30s",
		},
		{
			name: "thirty minutes",
			ms:   1800000,
			want: "30m 0s",
		},
		{
			name: "fifty-nine minutes fifty-nine seconds",
			ms:   3599999,
			want: "59m 59s",
		},
		{
			name: "exactly one hour",
			ms:   3600000,
			want: "1h 0m",
		},
		{
			name: "one hour thirty minutes",
			ms:   5400000,
			want: "1h 30m",
		},
		{
			name: "two hours fifteen minutes",
			ms:   8100000,
			want: "2h 15m",
		},
		{
			name: "ten hours",
			ms:   36000000,
			want: "10h 0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.ms)
			if got != tt.want {
				t.Fatalf("FormatDuration(%d): got %q want %q", tt.ms, got, tt.want)
			}
		})
	}
}
