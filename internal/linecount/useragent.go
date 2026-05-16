package linecount

import (
	"bufio"
	"io"
	"regexp"
	"sort"
	"strings"
)

var defaultUAPattern = regexp.MustCompile(`"([^"]+)"`)

// UserAgentCounter counts occurrences of user-agent strings extracted from log lines.
type UserAgentCounter struct {
	pattern *regexp.Regexp
	counts  map[string]int
	total   int
}

// NewUserAgentCounter creates a UserAgentCounter using the provided regex pattern.
// The pattern must contain at least one capture group. If pattern is empty the
// default quoted-string pattern is used.
func NewUserAgentCounter(pattern string) (*UserAgentCounter, error) {
	p := defaultUAPattern
	if pattern != "" {
		var err error
		p, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
	}
	return &UserAgentCounter{
		pattern: p,
		counts:  make(map[string]int),
	}, nil
}

// Add extracts a user-agent from line and increments its counter.
// Lines that do not match the pattern are silently skipped.
func (u *UserAgentCounter) Add(line string) {
	m := u.pattern.FindStringSubmatch(line)
	if len(m) < 2 {
		return
	}
	ua := strings.TrimSpace(m[1])
	if ua == "" {
		return
	}
	u.counts[ua]++
	u.total++
}

// Counts returns a copy of the internal frequency map.
func (u *UserAgentCounter) Counts() map[string]int {
	out := make(map[string]int, len(u.counts))
	for k, v := range u.counts {
		out[k] = v
	}
	return out
}

// Total returns the number of lines that matched.
func (u *UserAgentCounter) Total() int { return u.total }

// UserAgentEntry is a single user-agent with its count.
type UserAgentEntry struct {
	Agent string
	Count int
}

// SortedUserAgentEntries returns entries sorted by count descending, then agent ascending.
func SortedUserAgentEntries(counts map[string]int) []UserAgentEntry {
	entries := make([]UserAgentEntry, 0, len(counts))
	for a, c := range counts {
		entries = append(entries, UserAgentEntry{Agent: a, Count: c})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Agent < entries[j].Agent
	})
	return entries
}

// CountUserAgentReader reads all lines from r and returns a populated UserAgentCounter.
func CountUserAgentReader(r io.Reader, pattern string) (*UserAgentCounter, error) {
	c, err := NewUserAgentCounter(pattern)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		c.Add(sc.Text())
	}
	return c, sc.Err()
}
