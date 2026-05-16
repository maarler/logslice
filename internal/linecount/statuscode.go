package linecount

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
)

// StatusCodeCounter counts HTTP status codes grouped by class (2xx, 3xx, 4xx, 5xx).
type StatusCodeCounter struct {
	codes  map[int]int
	classes map[string]int
	total  int
	re     *regexp.Regexp
}

// NewStatusCodeCounter creates a counter that extracts HTTP status codes
// using the provided regex. The regex must contain a named group "code".
func NewStatusCodeCounter(pattern string) (*StatusCodeCounter, error) {
	if pattern == "" {
		pattern = `(?:^|\s)(?P<code>[1-5][0-9]{2})(?:\s|$)`
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("statuscode: invalid pattern: %w", err)
	}
	return &StatusCodeCounter{
		codes:   make(map[int]int),
		classes: make(map[string]int),
	}, re}, nil
}

// Add scans a single log line for an HTTP status code and records it.
func (c *StatusCodeCounter) Add(line string) {
	m := c.re.FindStringSubmatch(line)
	if m == nil {
		return
	}
	idx := c.re.SubexpIndex("code")
	if idx < 0 || idx >= len(m) {
		return
	}
	code, err := strconv.Atoi(m[idx])
	if err != nil {
		return
	}
	c.codes[code]++
	class := fmt.Sprintf("%dxx", code/100)
	c.classes[class]++
	c.total++
}

// Total returns the total number of status codes seen.
func (c *StatusCodeCounter) Total() int { return c.total }

// StatusCodeEntry is a single code with its count.
type StatusCodeEntry struct {
	Code  int
	Class string
	Count int
}

// SortedCodeEntries returns all code entries sorted by count descending.
func (c *StatusCodeCounter) SortedCodeEntries() []StatusCodeEntry {
	entries := make([]StatusCodeEntry, 0, len(c.codes))
	for code, count := range c.codes {
		entries = append(entries, StatusCodeEntry{
			Code:  code,
			Class: fmt.Sprintf("%dxx", code/100),
			Count: count,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Code < entries[j].Code
	})
	return entries
}

// CountStatusCodeReader reads lines from r and counts HTTP status codes.
func CountStatusCodeReader(r io.Reader, pattern string) (*StatusCodeCounter, error) {
	c, err := NewStatusCodeCounter(pattern)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		c.Add(sc.Text())
	}
	return c, sc.Err()
}
