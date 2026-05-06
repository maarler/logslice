package parser

import (
	"fmt"
	"regexp"
	"time"
)

// Common log timestamp formats to try in order
var timestampFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
	"Jan  2 15:04:05",
}

// timestampPattern matches common timestamp prefixes in log lines
var timestampPattern = regexp.MustCompile(
	`(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?|` +
		`\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4}|` +
		`\w{3}\s+\d{1,2} \d{2}:\d{2}:\d{2})`,
)

// ParseTimestamp extracts and parses the first timestamp found in a log line.
// Returns the parsed time and true if successful, or zero time and false otherwise.
func ParseTimestamp(line string) (time.Time, bool) {
	matches := timestampPattern.FindString(line)
	if matches == "" {
		return time.Time{}, false
	}

	for _, format := range timestampFormats {
		t, err := time.Parse(format, matches)
		if err == nil {
			// If the parsed time has no year (e.g. syslog format), use current year
			if t.Year() == 0 {
				t = t.AddDate(time.Now().Year(), 0, 0)
			}
			return t, true
		}
	}

	return time.Time{}, false
}

// ParseTimestampWithFormat attempts to parse a timestamp using a user-supplied format.
func ParseTimestampWithFormat(line, format string) (time.Time, bool) {
	matches := timestampPattern.FindString(line)
	if matches == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(format, matches)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// FormatDuration returns a human-readable duration string for display purposes.
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
