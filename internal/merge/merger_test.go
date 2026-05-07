package merge_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/merge"
)

func makeChannel(entries []merge.Entry) <-chan merge.Entry {
	ch := make(chan merge.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func TestMerge_OrdersByTimestamp(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	ch1 := makeChannel([]merge.Entry{
		{Line: "a1", Timestamp: t0.Add(1 * time.Second), Source: "f1"},
		{Line: "a3", Timestamp: t0.Add(3 * time.Second), Source: "f1"},
	})
	ch2 := makeChannel([]merge.Entry{
		{Line: "b2", Timestamp: t0.Add(2 * time.Second), Source: "f2"},
		{Line: "b4", Timestamp: t0.Add(4 * time.Second), Source: "f2"},
	})

	m := merge.New([]<-chan merge.Entry{ch1, ch2})
	out := m.Merge()

	want := []string{"a1", "b2", "a3", "b4"}
	for i, w := range want {
		e, ok := <-out
		if !ok {
			t.Fatalf("channel closed early at index %d", i)
		}
		if e.Line != w {
			t.Errorf("index %d: got %q, want %q", i, e.Line, w)
		}
	}
	if _, ok := <-out; ok {
		t.Error("expected channel to be closed")
	}
}

func TestMerge_SingleInput(t *testing.T) {
	t0 := time.Now()
	ch := makeChannel([]merge.Entry{
		{Line: "only", Timestamp: t0, Source: "f"},
	})

	m := merge.New([]<-chan merge.Entry{ch})
	out := m.Merge()

	e := <-out
	if e.Line != "only" {
		t.Errorf("got %q, want %q", e.Line, "only")
	}
	if _, ok := <-out; ok {
		t.Error("expected channel to be closed")
	}
}

func TestMerge_EmptyInputs(t *testing.T) {
	ch := makeChannel(nil)
	m := merge.New([]<-chan merge.Entry{ch})
	out := m.Merge()
	if _, ok := <-out; ok {
		t.Error("expected closed channel for empty input")
	}
}

func TestMerge_NoInputs(t *testing.T) {
	m := merge.New(nil)
	out := m.Merge()
	if _, ok := <-out; ok {
		t.Error("expected closed channel for no inputs")
	}
}

func TestMerge_PreservesSource(t *testing.T) {
	t0 := time.Now()
	ch := makeChannel([]merge.Entry{
		{Line: "line", Timestamp: t0, Source: "myfile.log"},
	})
	m := merge.New([]<-chan merge.Entry{ch})
	e := <-m.Merge()
	if e.Source != "myfile.log" {
		t.Errorf("source not preserved: got %q", e.Source)
	}
}
