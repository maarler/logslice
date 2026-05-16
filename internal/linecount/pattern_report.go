package linecount

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// PatternEntry holds a single pattern name and its hit count.
type PatternEntry struct {
	Name  string
	Count int64
}

// SortedPatternEntries converts a counts map to a slice sorted by count desc.
func SortedPatternEntries(counts map[string]int64) []PatternEntry {
	entries := make([]PatternEntry, 0, len(counts))
	for k, v := range counts {
		entries = append(entries, PatternEntry{Name: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Name < entries[j].Name
	})
	return entries
}

// WritePatternReport writes a formatted table of pattern hit counts to w.
func WritePatternReport(w io.Writer, counts map[string]int64, total int64) error {
	if len(counts) == 0 {
		_, err := fmt.Fprintln(w, "no pattern counts recorded")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATTERN\tCOUNT\tPCT")
	fmt.Fprintln(tw, repeatChar('-', 7)+"\t"+repeatChar('-', 5)+"\t"+repeatChar('-', 5))

	for _, e := range SortedPatternEntries(counts) {
		pct := 0.0
		if total > 0 {
			pct = float64(e.Count) / float64(total) * 100
		}
		fmt.Fprintf(tw, "%s\t%d\t%.1f%%\n", e.Name, e.Count, pct)
	}
	return tw.Flush()
}
