package tail_test

import (
	"context"
	"sort"
	"testing"

	"github.com/logslice/logslice/internal/tail"
)

func TestMultiFile_MergesLines(t *testing.T) {
	p1 := writeFile(t, "a1\na2\n")
	p2 := writeFile(t, "b1\nb2\nb3\n")

	mf := tail.NewMultiFile([]string{p1, p2}, tail.DefaultOptions())
	ctx := context.Background()
	entries, errCh := mf.Lines(ctx)

	var lines []string
	for e := range entries {
		lines = append(lines, e.Line)
	}
	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	sort.Strings(lines)
	want := []string{"a1", "a2", "b1", "b2", "b3"}
	if len(lines) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(lines), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line[%d]: want %q, got %q", i, w, lines[i])
		}
	}
}

func TestMultiFile_EntryCarriesPath(t *testing.T) {
	p := writeFile(t, "hello\n")
	mf := tail.NewMultiFile([]string{p}, tail.DefaultOptions())
	ctx := context.Background()
	entries, _ := mf.Lines(ctx)

	for e := range entries {
		if e.Path != p {
			t.Errorf("expected path %q, got %q", p, e.Path)
		}
	}
}

func TestMultiFile_NoPaths(t *testing.T) {
	mf := tail.NewMultiFile(nil, tail.DefaultOptions())
	ctx := context.Background()
	entries, errCh := mf.Lines(ctx)

	var count int
	for range entries {
		count++
	}
	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if count != 0 {
		t.Errorf("expected 0 entries, got %d", count)
	}
}
