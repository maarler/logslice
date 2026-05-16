package linecount

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteFieldReport writes a formatted table of field value counts to w.
// total is the total number of lines processed (used for percentage calculation).
func WriteFieldReport(w io.Writer, field string, entries []FieldEntry, total int) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintf(w, "No values found for field %q\n", field)
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "VALUE\tCOUNT\tPERCENT\tBAR\n")
	fmt.Fprintf(tw, "-----\t-----\t-------\t---\n")

	maxCount := 0
	if len(entries) > 0 {
		maxCount = entries[0].Count
	}

	for _, e := range entries {
		pct := 0.0
		if total > 0 {
			pct = float64(e.Count) / float64(total) * 100
		}
		bar := buildBar(e.Count, maxCount, 20)
		fmt.Fprintf(tw, "%s\t%d\t%.1f%%\t%s\n", e.Value, e.Count, pct, bar)
	}

	return tw.Flush()
}
