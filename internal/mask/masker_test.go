package mask_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/mask"
)

func TestNewRule_Invalid(t *testing.T) {
	_, err := mask.NewRule("bad", `[invalid`, "[X]")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMasker_Line_NoRules(t *testing.T) {
	m := mask.New()
	got := m.Line("hello 192.168.1.1")
	if got != "hello 192.168.1.1" {
		t.Errorf("unexpected mutation: %q", got)
	}
}

func TestMasker_Line_MasksIP(t *testing.T) {
	r, _ := mask.NewRule("ipv4", `\b(?:\d{1,3}\.){3}\d{1,3}\b`, "[IP]")
	m := mask.New(r)
	got := m.Line("connected from 10.0.0.1 to 10.0.0.2")
	want := "connected from [IP] to [IP]"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMasker_Line_MasksEmail(t *testing.T) {
	r, _ := mask.NewRule("email", `[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`, "[EMAIL]")
	m := mask.New(r)
	got := m.Line("user alice@example.com logged in")
	want := "user [EMAIL] logged in"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMasker_Lines(t *testing.T) {
	r, _ := mask.NewRule("ipv4", `\b(?:\d{1,3}\.){3}\d{1,3}\b`, "[IP]")
	m := mask.New(r)
	lines := []string{"req from 1.2.3.4", "no ip here", "src 5.6.7.8"}
	got := m.Lines(lines)
	if got[0] != "req from [IP]" {
		t.Errorf("line 0: %q", got[0])
	}
	if got[1] != "no ip here" {
		t.Errorf("line 1: %q", got[1])
	}
	if got[2] != "src [IP]" {
		t.Errorf("line 2: %q", got[2])
	}
}

func TestMasker_Apply_MultiLine(t *testing.T) {
	r, _ := mask.NewRule("email", `[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`, "[EMAIL]")
	m := mask.New(r)
	src := "login: a@b.com\nlogout: c@d.org\nok"
	got := m.Apply(src)
	want := "login: [EMAIL]\nlogout: [EMAIL]\nok"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPreset(t *testing.T) {
	m, err := mask.Preset()
	if err != nil {
		t.Fatalf("Preset: %v", err)
	}
	line := "user bob@example.com from 192.168.0.1 token=abc123"
	got := m.Line(line)
	for _, bad := range []string{"bob@example.com", "192.168.0.1", "token=abc123"} {
		if contains(got, bad) {
			t.Errorf("sensitive data %q not masked in: %q", bad, got)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
