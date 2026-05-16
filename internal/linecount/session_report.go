package linecount

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteSessionReport writes a human-readable session report to w.
func WriteSessionReport(w io.Writer, sessions []SessionEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "#\tStart\tEnd\tDuration\tLines")
	fmt.Fprintln(tw, "-\t-----\t---\t--------\t-----")

	if len(sessions) == 0 {
		tw.Flush()
		fmt.Fprintln(w, "(no sessions detected)")
		return nil
	}

	for i, s := range sessions {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%d\n",
			i+1,
			s.Start.Format(time.RFC3339),
			s.End.Format(time.RFC3339),
			formatSessionDuration(s.Duration()),
			s.LineCount,
		)
	}

	if err := tw.Flush(); err != nil {
		return fmt.Errorf("session report flush: %w", err)
	}

	totalLines := 0
	for _, s := range sessions {
		totalLines += s.LineCount
	}
	fmt.Fprintf(w, "\nTotal sessions: %d  Total lines: %d\n", len(sessions), totalLines)
	return nil
}

// formatSessionDuration formats a duration as a compact human string.
func formatSessionDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", h, m)
}
