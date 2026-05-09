package dedupe

import (
	"fmt"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.WindowSize != 512 {
		t.Fatalf("expected window 512, got %d", opts.WindowSize)
	}
	if opts.Consecutive {
		t.Fatal("expected Consecutive false")
	}
}

func TestNew_ZeroWindowUsesDefault(t *testing.T) {
	d := New(Options{WindowSize: 0})
	if len(d.window) != 512 {
		t.Fatalf("expected window 512, got %d", len(d.window))
	}
}

func TestIsDuplicate_NoDupes(t *testing.T) {
	d := New(DefaultOptions())
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if d.IsDuplicate(l) {
			t.Errorf("unexpected duplicate for %q", l)
		}
	}
	if d.Suppressed() != 0 {
		t.Fatalf("expected 0 suppressed, got %d", d.Suppressed())
	}
}

func TestIsDuplicate_DetectsDupe(t *testing.T) {
	d := New(DefaultOptions())
	d.IsDuplicate("hello")
	if !d.IsDuplicate("hello") {
		t.Fatal("expected duplicate")
	}
	if d.Suppressed() != 1 {
		t.Fatalf("expected 1 suppressed, got %d", d.Suppressed())
	}
}

func TestIsDuplicate_WindowEviction(t *testing.T) {
	window := 4
	d := New(Options{WindowSize: window})
	// Fill the window with unique lines.
	for i := 0; i < window; i++ {
		d.IsDuplicate(fmt.Sprintf("line-%d", i))
	}
	// "line-0" should have been evicted; re-inserting must not be a dupe.
	if d.IsDuplicate("line-0") {
		t.Fatal("line-0 should have been evicted from window")
	}
}

func TestIsDuplicate_ConsecutiveMode(t *testing.T) {
	d := New(Options{WindowSize: 8, Consecutive: true})
	if d.IsDuplicate("a") {
		t.Fatal("first occurrence must not be duplicate")
	}
	if !d.IsDuplicate("a") {
		t.Fatal("consecutive repeat must be duplicate")
	}
	// Different line resets consecutive check.
	if d.IsDuplicate("b") {
		t.Fatal("different line must not be duplicate")
	}
	// "a" again after "b" must NOT be a duplicate in consecutive mode.
	if d.IsDuplicate("a") {
		t.Fatal("non-consecutive repeat must not be duplicate in consecutive mode")
	}
	if d.Suppressed() != 1 {
		t.Fatalf("expected 1 suppressed, got %d", d.Suppressed())
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(DefaultOptions())
	d.IsDuplicate("x")
	d.IsDuplicate("x") // suppressed
	d.Reset()
	if d.Suppressed() != 0 {
		t.Fatal("expected 0 suppressed after reset")
	}
	if d.IsDuplicate("x") {
		t.Fatal("after reset x must not be a duplicate")
	}
}
