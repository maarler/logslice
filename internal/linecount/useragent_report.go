package linecount

import (
	"fmt"
	"io"
)

// WriteUserAgentReport writes a formatted user-agent frequency table to w.
func WriteUserAgentReport(w io.Writer, entries []UserAgentEntry, total int) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No user-agent data found.")
		return
	}

	fmt.Fprintf(w, "%-60s %8s %8s\n", "User-Agent", "Count", "Percent")
	fmt.Fprintf(w, "%s %s %s\n", repeatChar('-', 60), repeatChar('-', 8), repeatChar('-', 8))

	for _, e := range entries {
		pct := 0.0
		if total > 0 {
			pct = float64(e.Count) / float64(total) * 100
		}
		agent := e.Agent
		if len(agent) > 58 {
			agent = agent[:55] + "..."
		}
		fmt.Fprintf(w, "%-60s %8d %7.1f%%\n", agent, e.Count, pct)
	}

	fmt.Fprintf(w, "%s\n", repeatChar('-', 80))
	fmt.Fprintf(w, "%-60s %8d\n", "Total", total)
}
