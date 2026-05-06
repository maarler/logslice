package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantOK   bool
		wantYear int
		wantHour int
	}{
		{
			name:     "ISO8601 with timezone",
			line:     "2024-03-15T14:22:01+00:00 ERROR something failed",
			wantOK:   true,
			wantYear: 2024,
			wantHour: 14,
		},
		{
			name:     "ISO8601 with milliseconds",
			line:     "2024-03-15T14:22:01.123Z INFO request completed",
			wantOK:   true,
			wantYear: 2024,
			wantHour: 14,
		},
		{
			name:     "space-separated datetime",
			line:     "2024-03-15 09:05:33 WARN disk usage high",
			wantOK:   true,
			wantYear: 2024,
			wantHour: 9,
		},
		{
			name:     "nginx/apache combined log format",
			line:     `192.168.1.1 - - [15/Mar/2024:10:30:00 +0000] "GET / HTTP/1.1" 200 512`,
			wantOK:   true,
			wantYear: 2024,
			wantHour: 10,
		},
		{
			name:   "no timestamp",
			line:   "plain log line with no timestamp at all",
			wantOK: false,
		},
		{
			name:   "empty line",
			line:   "",
			wantOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ParseTimestamp(tc.line)
			if ok != tc.wantOK {
				t.Fatalf("ParseTimestamp(%q) ok = %v, want %v", tc.line, ok, tc.wantOK)
			}
			if !tc.wantOK {
				return
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year = %d, want %d", got.Year(), tc.wantYear)
			}
			if got.Hour() != tc.wantHour {
				t.Errorf("hour = %d, want %d", got.Hour(), tc.wantHour)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Second, "30s"},
		{90 * time.Second, "1m30s"},
		{2*time.Hour + 15*time.Minute, "2h15m"},
	}
	for _, tc := range cases {
		if got := FormatDuration(tc.d); got != tc.want {
			t.Errorf("FormatDuration(%v) = %q, want %q", tc.d, got, tc.want)
		}
	}
}
