package linecount

import (
	"fmt"
	"io"
)

const errReportHeader = "Level            Count     Percent"
const errReportSep = "------           -----     -------"

// WriteErrorReport writes a formatted error-level summary to w.
func WriteErrorReport(w io.Writer, entries []ErrorEntry, total int) error {
	if _, err := fmt.Fprintln(w, errReportHeader); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, errReportSep); err != nil {
		return err
	}
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "(no matching lines)")
		return err
	}
	for _, e := range entries {
		pct := 0.0
		if total > 0 {
			pct = float64(e.Count) / float64(total) * 100
		}
		if _, err := fmt.Fprintf(w, "%-16s %-9d %.1f%%\n", e.Level, e.Count, pct); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "\nTotal matched: %d\n", total); err != nil {
		return err
	}
	return nil
}
