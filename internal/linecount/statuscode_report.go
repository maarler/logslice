package linecount

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteStatusCodeReport writes a formatted HTTP status code report to w.
func WriteStatusCodeReport(w io.Writer, c *StatusCodeCounter) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "CODE\tCLASS\tCOUNT\tPERCENT\tBAR")
	fmt.Fprintln(tw, "----\t-----\t-----\t-------\t---")

	if c.Total() == 0 {
		tw.Flush()
		fmt.Fprintln(w, "(no status codes found)")
		return nil
	}

	entries := c.SortedCodeEntries()
	maxCount := 0
	for _, e := range entries {
		if e.Count > maxCount {
			maxCount = e.Count
		}
	}

	for _, e := range entries {
		pct := float64(e.Count) / float64(c.Total()) * 100
		bar := buildBar(e.Count, maxCount, 20)
		fmt.Fprintf(tw, "%d\t%s\t%d\t%.1f%%\t%s\n",
			e.Code, e.Class, e.Count, pct, bar)
	}

	fmt.Fprintln(tw, "----\t-----\t-----\t-------\t---")
	fmt.Fprintf(tw, "TOTAL\t\t%d\t100.0%%\t\n", c.Total())

	return tw.Flush()
}
