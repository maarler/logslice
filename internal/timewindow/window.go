// Package timewindow provides utilities for bucketing log entries into
// fixed-duration time windows.
package timewindow

import (
	"fmt"
	"time"
)

// Window represents a single time bucket with a start and end boundary.
type Window struct {
	Start time.Time
	End   time.Time
}

// Contains reports whether t falls within the window [Start, End).
func (w Window) Contains(t time.Time) bool {
	return !t.Before(w.Start) && t.Before(w.End)
}

// Key returns a string key for the window suitable for use in filenames.
// Format: "2006-01-02T15-04-05".
func (w Window) Key() string {
	return w.Start.UTC().Format("2006-01-02T15-04-05")
}

// String implements fmt.Stringer.
func (w Window) String() string {
	return fmt.Sprintf("[%s, %s)", w.Start.UTC().Format(time.RFC3339), w.End.UTC().Format(time.RFC3339))
}

// Bucket returns the Window that contains t for a given duration d.
// The window is aligned to the Unix epoch so that boundaries are consistent
// across calls regardless of the first observed timestamp.
func Bucket(t time.Time, d time.Duration) Window {
	if d <= 0 {
		panic("timewindow: duration must be positive")
	}
	unix := t.UnixNano()
	slot := unix - (unix % int64(d))
	start := time.Unix(0, slot).UTC()
	return Window{
		Start: start,
		End:   start.Add(d),
	}
}

// Sequence returns all non-empty windows between from and to (inclusive)
// for the given duration d, in chronological order.
func Sequence(from, to time.Time, d time.Duration) []Window {
	if d <= 0 {
		panic("timewindow: duration must be positive")
	}
	var windows []Window
	current := Bucket(from, d)
	for !current.Start.After(to) {
		windows = append(windows, current)
		current = Window{
			Start: current.End,
			End:   current.End.Add(d),
		}
	}
	return windows
}
