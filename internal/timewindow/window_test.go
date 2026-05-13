package timewindow_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/timewindow"
)

func ts(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestBucket_AlignedToEpoch(t *testing.T) {
	d := 5 * time.Minute
	w := timewindow.Bucket(ts("2024-03-15T10:07:42Z"), d)
	if !w.Start.Equal(ts("2024-03-15T10:05:00Z")) {
		t.Errorf("unexpected start: %s", w.Start)
	}
	if !w.End.Equal(ts("2024-03-15T10:10:00Z")) {
		t.Errorf("unexpected end: %s", w.End)
	}
}

func TestWindow_Contains(t *testing.T) {
	w := timewindow.Bucket(ts("2024-03-15T10:00:00Z"), time.Hour)
	if !w.Contains(ts("2024-03-15T10:30:00Z")) {
		t.Error("expected Contains to be true for mid-window time")
	}
	if w.Contains(ts("2024-03-15T11:00:00Z")) {
		t.Error("expected Contains to be false for end boundary")
	}
	if w.Contains(ts("2024-03-15T09:59:59Z")) {
		t.Error("expected Contains to be false for time before window")
	}
}

func TestWindow_Key(t *testing.T) {
	w := timewindow.Bucket(ts("2024-03-15T10:07:00Z"), 5*time.Minute)
	got := w.Key()
	want := "2024-03-15T10-05-00"
	if got != want {
		t.Errorf("Key() = %q, want %q", got, want)
	}
}

func TestWindow_String(t *testing.T) {
	w := timewindow.Bucket(ts("2024-03-15T10:07:00Z"), 5*time.Minute)
	s := w.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
}

func TestSequence_ReturnsCorrectCount(t *testing.T) {
	from := ts("2024-03-15T10:00:00Z")
	to := ts("2024-03-15T10:14:59Z")
	windows := timewindow.Sequence(from, to, 5*time.Minute)
	if len(windows) != 3 {
		t.Errorf("expected 3 windows, got %d", len(windows))
	}
}

func TestSequence_Chronological(t *testing.T) {
	from := ts("2024-03-15T10:00:00Z")
	to := ts("2024-03-15T10:59:00Z")
	windows := timewindow.Sequence(from, to, 15*time.Minute)
	for i := 1; i < len(windows); i++ {
		if !windows[i].Start.After(windows[i-1].Start) {
			t.Errorf("windows not in chronological order at index %d", i)
		}
	}
}

func TestSequence_SinglePoint(t *testing.T) {
	at := ts("2024-03-15T10:07:00Z")
	windows := timewindow.Sequence(at, at, time.Hour)
	if len(windows) != 1 {
		t.Errorf("expected 1 window, got %d", len(windows))
	}
}

func TestBucket_PanicOnZeroDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero duration")
		}
	}()
	timewindow.Bucket(time.Now(), 0)
}
