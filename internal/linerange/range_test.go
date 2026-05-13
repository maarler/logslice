package linerange_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linerange"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		r       linerange.Range
		wantErr bool
	}{
		{"valid open", linerange.Range{First: 1, Last: 0}, false},
		{"valid closed", linerange.Range{First: 3, Last: 7}, false},
		{"zero first", linerange.Range{First: 0, Last: 5}, true},
		{"last before first", linerange.Range{First: 5, Last: 3}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.r.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestContains(t *testing.T) {
	r := linerange.Range{First: 3, Last: 5}
	for n, want := range map[int64]bool{
		1: false, 2: false, 3: true, 4: true, 5: true, 6: false,
	} {
		if got := r.Contains(n); got != want {
			t.Errorf("Contains(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestApply_MiddleRange(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\n"
	var out bytes.Buffer
	kept, err := linerange.Apply(strings.NewReader(input), &out, linerange.Range{First: 2, Last: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kept != 3 {
		t.Fatalf("kept = %d, want 3", kept)
	}
	want := "line2\nline3\nline4\n"
	if out.String() != want {
		t.Errorf("output = %q, want %q", out.String(), want)
	}
}

func TestApply_OpenRange(t *testing.T) {
	input := "a\nb\nc\n"
	var out bytes.Buffer
	kept, err := linerange.Apply(strings.NewReader(input), &out, linerange.Range{First: 2, Last: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kept != 2 {
		t.Fatalf("kept = %d, want 2", kept)
	}
}

func TestApply_InvalidRange(t *testing.T) {
	_, err := linerange.Apply(strings.NewReader("x\n"), &bytes.Buffer{}, linerange.Range{First: 0})
	if err == nil {
		t.Fatal("expected error for invalid range")
	}
}

func TestApply_EmptyInput(t *testing.T) {
	var out bytes.Buffer
	kept, err := linerange.Apply(strings.NewReader(""), &out, linerange.Range{First: 1, Last: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kept != 0 {
		t.Fatalf("kept = %d, want 0", kept)
	}
}
