package filter

import (
	"testing"
)

func TestNewPattern_Valid(t *testing.T) {
	p, err := NewPattern(`ERROR`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Match("2024-01-01 ERROR something broke") {
		t.Error("expected match for ERROR line")
	}
	if p.Match("2024-01-01 INFO all good") {
		t.Error("expected no match for INFO line")
	}
}

func TestNewPattern_Invalid(t *testing.T) {
	_, err := NewPattern(`[invalid`, false)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestPattern_Negate(t *testing.T) {
	p, err := NewPattern(`DEBUG`, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Match("DEBUG verbose output") {
		t.Error("negated pattern should not match DEBUG line")
	}
	if !p.Match("INFO startup complete") {
		t.Error("negated pattern should match non-DEBUG line")
	}
}

func TestPattern_String(t *testing.T) {
	p, _ := NewPattern(`WARN`, false)
	if p.String() != "WARN" {
		t.Errorf("unexpected String(): %s", p.String())
	}
	pn, _ := NewPattern(`WARN`, true)
	if pn.String() != "NOT(WARN)" {
		t.Errorf("unexpected negated String(): %s", pn.String())
	}
}

func TestPattern_CaseInsensitive(t *testing.T) {
	p, err := NewPattern(`(?i)error`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := []struct {
		line  string
		want  bool
	}{
		{"ERROR: something failed", true},
		{"Error: something failed", true},
		{"error: something failed", true},
		{"INFO: all good", false},
	}
	for _, tc := range cases {
		if got := p.Match(tc.line); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}
