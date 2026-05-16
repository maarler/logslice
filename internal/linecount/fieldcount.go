package linecount

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// FieldCounter counts occurrences of unique values for a specific field key.
type FieldCounter struct {
	field  string
	counts map[string]int
}

// NewFieldCounter creates a FieldCounter that tracks values for the given field key.
func NewFieldCounter(field string) *FieldCounter {
	return &FieldCounter{
		field:  field,
		counts: make(map[string]int),
	}
}

// Add parses key=value pairs from line and increments the count for the
// value associated with the tracked field. Unmatched lines are silently skipped.
func (fc *FieldCounter) Add(line string) {
	for _, token := range strings.Fields(line) {
		if k, v, ok := strings.Cut(token, "="); ok && k == fc.field {
			v = strings.Trim(v, `"`)
			fc.counts[v]++
			return
		}
	}
}

// Counts returns a copy of the current value→count map.
func (fc *FieldCounter) Counts() map[string]int {
	out := make(map[string]int, len(fc.counts))
	for k, v := range fc.counts {
		out[k] = v
	}
	return out
}

// FieldEntry is a single value and its count.
type FieldEntry struct {
	Value string
	Count int
}

// SortedFieldEntries returns entries sorted by count descending, then value ascending.
func SortedFieldEntries(counts map[string]int) []FieldEntry {
	entries := make([]FieldEntry, 0, len(counts))
	for v, c := range counts {
		entries = append(entries, FieldEntry{Value: v, Count: c})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	return entries
}

// CountFieldReader reads all lines from r, accumulates field value counts, and
// returns the sorted entries.
func CountFieldReader(r io.Reader, field string) ([]FieldEntry, error) {
	fc := NewFieldCounter(field)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fc.Add(sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("fieldcount: scan: %w", err)
	}
	return SortedFieldEntries(fc.Counts()), nil
}
