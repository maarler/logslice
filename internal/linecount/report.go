package linecount

import (
	"fmt"
	"io"
	"text/tabwriter"
)

const (
	// MinBarWidth is the minimum width of the histogram bar.
	MinBarWidth = 1
	// MaxBarWidth is the maximum width of the histogram bar.
	MaxBarWidth = 40
)

// WriteWindowReport writes a formatted table of window statistics to w.
// Each row shows the window key, line count, byte count, and a proportional bar.
func WriteWindowReport(w io.Writer, stats []*WindowStats) error {
	if len(stats) == 0 {
		_, err := fmt.Fprintln(w, "(no windows)")
		return err
	}

	var maxLines int64
	for _, s := range stats {
		if s.Lines > maxLines {
			maxLines = s.Lines
		}
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "WINDOW\tLINES\tBYTES\tDISTRIBUTION")
	fmt.Fprintln(tw, "------\t-----\t-----\t------------")

	for _, s := range stats {
		bar := buildBar(s.Lines, maxLines, MaxBarWidth)
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n",
			s.Key,
			s.Lines,
			formatBytes(s.Bytes),
			bar,
		)
	}
	return tw.Flush()
}

func buildBar(value, max int64, width int) string {
	if max == 0 {
		return ""
	}
	filled := int(float64(value) / float64(max) * float64(width))
	if filled < MinBarWidth {
		filled = MinBarWidth
	}
	bar := make([]byte, filled)
	for i := range bar {
		bar[i] = '#'
	}
	return string(bar)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGTPE"[exp])
}
