package linecount

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteIPReport writes a formatted IP address frequency report to w.
func WriteIPReport(w io.Writer, entries []IPEntry, total int) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if len(entries) == 0 {
		_, err := fmt.Fprintln(tw, "No IP addresses found.")
		_ = tw.Flush()
		return err
	}

	_, err := fmt.Fprintln(tw, "IP ADDRESS\tCOUNT\tPERCENT")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(tw, "----------\t-----\t-------")
	if err != nil {
		return err
	}

	for _, e := range entries {
		_, err = fmt.Fprintf(tw, "%s\t%d\t%s\n", e.IP, e.Count, fmtIPPercent(e.Count, total))
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(tw, "\nTotal hits: %d\n", total)
	if err != nil {
		return err
	}

	return tw.Flush()
}
