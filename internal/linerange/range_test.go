package linerange

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		r       Range
		wantErr bool
	}{
		{"valid bounded", Range{Start: 1, End: 5}, false},
		{"valid open", Range{Start: 3, End: 0}, false},
		{"start zero", Range{Start: 0, End: 5}, true},
		{"end before start", Range{Start: 5, End: 3}, true},
		{"start equals end", Range{Start: 4, End: 4}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.r.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestContains(t *testing.T) {
	r := Range{Start: 3, End: 6}
	for n, want := range map[int64]bool{1: false, 2: false, 3: true, 5: true, 6: true, 7: false} {
		if got := r.Contains(n); got != want {
			t.Errorf("Contains(%d) = %v, want %v", n, got, want)
		}
	}
	open := Range{Start: 4, End: 0}
	if !open.Contains(999) {
		t.Error("open range should contain large line numbers")
	}
	if open.Contains(3) {
		t.Error("open range should not contain lines before Start")
	}
}

func TestApply_MiddleRange(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\n"
	r := Range{Start: 2, End: 4}
	var out bytes.Buffer
	if err := Apply(r, strings.NewReader(input), &out); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	want := "line2\nline3\nline4\n"
	if got := out.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_OpenRange(t *testing.T) {
	input := "a\nb\nc\nd\n"
	r := Range{Start: 3, End: 0}
	var out bytes.Buffer
	if err := Apply(r, strings.NewReader(input), &out); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	want := "c\nd\n"
	if got := out.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_InvalidRange(t *testing.T) {
	err := Apply(Range{Start: 0}, strings.NewReader("x\n"), &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for invalid range")
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		input   string
		want    Range
		wantErr bool
	}{
		{"5", Range{Start: 5, End: 5}, false},
		{"3-7", Range{Start: 3, End: 7}, false},
		{"2-", Range{Start: 2, End: 0}, false},
		{"", Range{}, true},
		{"abc", Range{}, true},
		{"5-3", Range{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
			if !tc.wantErr && got != tc.want {
				t.Errorf("Parse(%q) = %+v, want %+v", tc.input, got, tc.want)
			}
		})
	}
}
