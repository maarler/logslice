package linecount

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"sort"
)

var ipRegexp = regexp.MustCompile(`\b((?:\d{1,3}\.){3}\d{1,3}|(?:[0-9a-fA-F]{1,4}:){2,7}[0-9a-fA-F]{1,4})\b`)

// IPCounter counts occurrences of IP addresses found in log lines.
type IPCounter struct {
	counts map[string]int
	total  int
}

// NewIPCounter returns a new IPCounter.
func NewIPCounter() *IPCounter {
	return &IPCounter{counts: make(map[string]int)}
}

// Add scans line for IP addresses and increments their counts.
func (c *IPCounter) Add(line string) {
	matches := ipRegexp.FindAllString(line, -1)
	for _, m := range matches {
		if net.ParseIP(m) != nil {
			c.counts[m]++
			c.total++
		}
	}
}

// Counts returns a copy of the internal counts map.
func (c *IPCounter) Counts() map[string]int {
	out := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Total returns the total number of IP occurrences seen.
func (c *IPCounter) Total() int { return c.total }

// IPEntry holds a single IP address and its hit count.
type IPEntry struct {
	IP    string
	Count int
}

// SortedIPEntries returns entries sorted by count descending, then IP ascending.
func SortedIPEntries(counts map[string]int) []IPEntry {
	entries := make([]IPEntry, 0, len(counts))
	for ip, n := range counts {
		entries = append(entries, IPEntry{IP: ip, Count: n})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].IP < entries[j].IP
	})
	return entries
}

// CountIPReader reads all lines from r, counts IPs, and returns the counter.
func CountIPReader(r io.Reader) *IPCounter {
	c := NewIPCounter()
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		c.Add(sc.Text())
	}
	return c
}

// Sprintf helper used by report.
func fmtIPPercent(n, total int) string {
	if total == 0 {
		return "0.00%"
	}
	return fmt.Sprintf("%.2f%%", float64(n)/float64(total)*100)
}
