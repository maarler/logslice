package output

import (
	"testing"
	"time"
)

func TestNewNamer_DefaultPattern(t *testing.T) {
	n := NewNamer("", "log")
	if n == nil {
		t.Fatal("expected non-nil Namer")
	}
}

func TestNamer_Generate_TimeKey(t *testing.T) {
	n := NewNamer("{time}", "log")
	ts := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	key := ts.Format("2006-01-02T15-04")
	got := n.Generate(key)
	if got == "" {
		t.Fatal("expected non-empty filename")
	}
}

func TestNamer_Generate_IndexIncrement(t *testing.T) {
	n := NewNamer("{time}-{index}", "log")
	key := "2024-03-15T14-30"
	first := n.Generate(key)
	second := n.Generate(key)
	if first == second {
		t.Errorf("expected different filenames for same key, got %q and %q", first, second)
	}
}

func TestNamer_Generate_Sanitize(t *testing.T) {
	n := NewNamer("{time}", "log")
	key := "2024/03/15 14:30"
	got := n.Generate(key)
	for _, ch := range got {
		if ch == '/' || ch == ' ' || ch == ':' {
			t.Errorf("filename contains unsafe character %q: %s", ch, got)
		}
	}
}

func TestNamer_Generate_Extension(t *testing.T) {
	n := NewNamer("{time}", "txt")
	key := "2024-03-15T14-30"
	got := n.Generate(key)
	if len(got) < 4 || got[len(got)-4:] != ".txt" {
		t.Errorf("expected .txt extension, got %q", got)
	}
}

func TestSanitize(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"hello world", "hello_world"},
		{"foo/bar", "foo_bar"},
		{"2024:03:15", "2024_03_15"},
		{"normal-name", "normal-name"},
		{"", ""},
	}
	for _, tc := range cases {
		got := sanitize(tc.input)
		if got != tc.want {
			t.Errorf("sanitize(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
