package linecount

import (
	"bufio"
	"io"
	"regexp"
	"sort"
	"strings"
)

// ErrorLevel represents a log severity level.
type ErrorLevel string

const (
	LevelDebug ErrorLevel = "DEBUG"
	LevelInfo  ErrorLevel = "INFO"
	LevelWarn  ErrorLevel = "WARN"
	LevelError ErrorLevel = "ERROR"
	LevelFatal ErrorLevel = "FATAL"
)

var levelPattern = regexp.MustCompile(`(?i)\b(DEBUG|INFO|WARN(?:ING)?|ERROR|FATAL|CRITICAL)\b`)

// ErrorCounter counts log lines by severity level.
type ErrorCounter struct {
	counts map[ErrorLevel]int
	total  int
}

// NewErrorCounter returns an initialised ErrorCounter.
func NewErrorCounter() *ErrorCounter {
	return &ErrorCounter{
		counts: make(map[ErrorLevel]int),
	}
}

// Add examines line and increments the matching level bucket.
// Lines that contain no recognisable level are ignored.
func (c *ErrorCounter) Add(line string) {
	m := levelPattern.FindString(line)
	if m == "" {
		return
	}
	level := normalise(m)
	c.counts[level]++
	c.total++
}

// Counts returns a copy of the internal level→count map.
func (c *ErrorCounter) Counts() map[ErrorLevel]int {
	out := make(map[ErrorLevel]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Total returns the number of lines that matched any level.
func (c *ErrorCounter) Total() int { return c.total }

// SortedErrorEntries returns level/count pairs sorted by count descending.
func SortedErrorEntries(counts map[ErrorLevel]int) []ErrorEntry {
	entries := make([]ErrorEntry, 0, len(counts))
	for lvl, n := range counts {
		entries = append(entries, ErrorEntry{Level: lvl, Count: n})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Level < entries[j].Level
	})
	return entries
}

// ErrorEntry is a single level/count pair.
type ErrorEntry struct {
	Level ErrorLevel
	Count int
}

// CountErrorReader scans r and returns a populated ErrorCounter.
func CountErrorReader(r io.Reader) *ErrorCounter {
	c := NewErrorCounter()
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		c.Add(sc.Text())
	}
	return c
}

func normalise(raw string) ErrorLevel {
	upper := strings.ToUpper(raw)
	if upper == "WARNING" || upper == "CRITICAL" {
		if upper == "WARNING" {
			return LevelWarn
		}
		return LevelFatal
	}
	return ErrorLevel(upper)
}
