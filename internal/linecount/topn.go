package linecount

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// TopNEntry holds a key and its associated line count.
type TopNEntry struct {
	Key   string
	Lines int64
	Bytes int64
}

// TopN returns the n entries with the highest line counts from a WindowCounter.
// If n <= 0 all entries are returned, sorted descending by line count.
func TopN(wc *WindowCounter, n int) []TopNEntry {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	entries := make([]TopNEntry, 0, len(wc.windows))
	for _, w := range wc.windows {
		entries = append(entries, TopNEntry{
			Key:   w.Key,
			Lines: w.Lines,
			Bytes: w.Bytes,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Lines != entries[j].Lines {
			return entries[i].Lines > entries[j].Lines
		}
		return entries[i].Key < entries[j].Key
	})

	if n > 0 && n < len(entries) {
		return entries[:n]
	}
	return entries
}

// WriteTopNReport writes a ranked table of the top-n windows to w.
func WriteTopNReport(out io.Writer, wc *WindowCounter, n int) error {
	entries := TopN(wc, n)
	if len(entries) == 0 {
		_, err := fmt.Fprintln(out, "(no data)")
		return err
	}

	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RANK\tWINDOW\tLINES\tBYTES")
	fmt.Fprintln(tw, "----\t------\t-----\t-----")
	for i, e := range entries {
		fmt.Fprintf(tw, "%d\t%s\t%d\t%s\n", i+1, e.Key, e.Lines, formatBytes(e.Bytes))
	}
	return tw.Flush()
}
