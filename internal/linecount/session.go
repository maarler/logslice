package linecount

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// SessionOptions configures session gap detection.
type SessionOptions struct {
	// GapDuration is the idle period that defines a new session boundary.
	GapDuration time.Duration
}

// DefaultSessionOptions returns sensible defaults (30-minute gap).
func DefaultSessionOptions() SessionOptions {
	return SessionOptions{
		GapDuration: 30 * time.Minute,
	}
}

// SessionEntry holds stats for a single detected session.
type SessionEntry struct {
	Start     time.Time
	End       time.Time
	LineCount int
}

// Duration returns the length of the session.
func (e SessionEntry) Duration() time.Duration {
	return e.End.Sub(e.Start)
}

// SessionCounter detects sessions in a timestamped log stream.
type SessionCounter struct {
	opts     SessionOptions
	sessions []SessionEntry
	current  *SessionEntry
	last     time.Time
}

// NewSessionCounter creates a SessionCounter with the given options.
func NewSessionCounter(opts SessionOptions) *SessionCounter {
	if opts.GapDuration <= 0 {
		opts.GapDuration = DefaultSessionOptions().GapDuration
	}
	return &SessionCounter{opts: opts}
}

// Add processes a single timestamp, updating session boundaries.
func (sc *SessionCounter) Add(t time.Time) {
	if sc.current == nil {
		sc.current = &SessionEntry{Start: t, End: t, LineCount: 1}
		sc.last = t
		return
	}
	gap := t.Sub(sc.last)
	if gap < 0 {
		gap = -gap
	}
	if gap >= sc.opts.GapDuration {
		sc.sessions = append(sc.sessions, *sc.current)
		sc.current = &SessionEntry{Start: t, End: t, LineCount: 1}
	} else {
		sc.current.End = t
		sc.current.LineCount++
	}
	sc.last = t
}

// Sessions returns all completed and the current (in-progress) session.
func (sc *SessionCounter) Sessions() []SessionEntry {
	result := make([]SessionEntry, len(sc.sessions))
	copy(result, sc.sessions)
	if sc.current != nil {
		result = append(result, *sc.current)
	}
	return result
}

// CountSessionReader scans r and returns detected sessions.
func CountSessionReader(r io.Reader, opts SessionOptions) ([]SessionEntry, error) {
	sc := NewSessionCounter(opts)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t, err := parser.ParseTimestamp(scanner.Text())
		if err != nil {
			continue
		}
		sc.Add(t)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("session scan: %w", err)
	}
	return sc.Sessions(), nil
}
