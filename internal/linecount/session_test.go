package linecount

import (
	"strings"
	"testing"
	"time"
)

func sessionTime(offset time.Duration) time.Time {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return base.Add(offset)
}

func TestSessionCounter_Empty(t *testing.T) {
	sc := NewSessionCounter(DefaultSessionOptions())
	if got := sc.Sessions(); len(got) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(got))
	}
}

func TestSessionCounter_SingleSession(t *testing.T) {
	sc := NewSessionCounter(DefaultSessionOptions())
	sc.Add(sessionTime(0))
	sc.Add(sessionTime(1 * time.Minute))
	sc.Add(sessionTime(2 * time.Minute))

	sessions := sc.Sessions()
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].LineCount != 3 {
		t.Errorf("expected 3 lines, got %d", sessions[0].LineCount)
	}
}

func TestSessionCounter_MultipleSessionsByGap(t *testing.T) {
	opts := SessionOptions{GapDuration: 10 * time.Minute}
	sc := NewSessionCounter(opts)

	sc.Add(sessionTime(0))
	sc.Add(sessionTime(5 * time.Minute))
	// gap > 10 min triggers new session
	sc.Add(sessionTime(20 * time.Minute))
	sc.Add(sessionTime(25 * time.Minute))

	sessions := sc.Sessions()
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
	if sessions[0].LineCount != 2 {
		t.Errorf("session 1: expected 2 lines, got %d", sessions[0].LineCount)
	}
	if sessions[1].LineCount != 2 {
		t.Errorf("session 2: expected 2 lines, got %d", sessions[1].LineCount)
	}
}

func TestSessionCounter_ZeroGapUsesDefault(t *testing.T) {
	sc := NewSessionCounter(SessionOptions{GapDuration: 0})
	if sc.opts.GapDuration != DefaultSessionOptions().GapDuration {
		t.Errorf("expected default gap, got %v", sc.opts.GapDuration)
	}
}

func TestSessionEntry_Duration(t *testing.T) {
	e := SessionEntry{
		Start:     sessionTime(0),
		End:       sessionTime(90 * time.Minute),
		LineCount: 10,
	}
	if e.Duration() != 90*time.Minute {
		t.Errorf("expected 90m, got %v", e.Duration())
	}
}

func TestCountSessionReader_ParsesLines(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T10:00:00Z INFO starting",
		"2024-01-15T10:05:00Z INFO running",
		"2024-01-15T11:00:00Z INFO restarted",
	}, "\n")

	opts := SessionOptions{GapDuration: 10 * time.Minute}
	sessions, err := CountSessionReader(strings.NewReader(input), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestCountSessionReader_SkipsUnparseable(t *testing.T) {
	input := "not a timestamp\n2024-01-15T10:00:00Z INFO ok\n"
	sessions, err := CountSessionReader(strings.NewReader(input), DefaultSessionOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
}
