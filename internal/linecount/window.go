package linecount

import (
	"bufio"
	"io"
	"time"
)

// WindowStats holds line count statistics for a time window.
type WindowStats struct {
	Key   string
	Start time.Time
	End   time.Time
	Lines int64
	Bytes int64
}

// WindowCounter accumulates per-window line and byte counts.
type WindowCounter struct {
	windows map[string]*WindowStats
	order   []string
}

// NewWindowCounter creates an empty WindowCounter.
func NewWindowCounter() *WindowCounter {
	return &WindowCounter{
		windows: make(map[string]*WindowStats),
	}
}

// Add records a line belonging to the given window key.
func (wc *WindowCounter) Add(key string, start, end time.Time, line string) {
	if _, ok := wc.windows[key]; !ok {
		wc.windows[key] = &WindowStats{Key: key, Start: start, End: end}
		wc.order = append(wc.order, key)
	}
	s := wc.windows[key]
	s.Lines++
	s.Bytes += int64(len(line))
}

// Stats returns window statistics in insertion order.
func (wc *WindowCounter) Stats() []*WindowStats {
	out := make([]*WindowStats, 0, len(wc.order))
	for _, k := range wc.order {
		out = append(out, wc.windows[k])
	}
	return out
}

// CountWindowReader counts lines and bytes from r, grouping them by the key
// returned by keyFn for each line. keyFn receives the raw line text.
func CountWindowReader(r io.Reader, keyFn func(line string) (key string, start, end time.Time)) *WindowCounter {
	wc := NewWindowCounter()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		key, start, end := keyFn(line)
		if key == "" {
			continue
		}
		wc.Add(key, start, end, line)
	}
	return wc
}
