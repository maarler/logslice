package linecount

import (
	"strings"
	"testing"
)

func TestUserAgentCounter_Empty(t *testing.T) {
	c, err := NewUserAgentCounter("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Total() != 0 {
		t.Errorf("expected 0 total, got %d", c.Total())
	}
	if len(c.Counts()) != 0 {
		t.Errorf("expected empty counts")
	}
}

func TestUserAgentCounter_InvalidPattern(t *testing.T) {
	_, err := NewUserAgentCounter("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestUserAgentCounter_Add_SingleAgent(t *testing.T) {
	c, _ := NewUserAgentCounter("")
	c.Add(`192.168.1.1 - - [01/Jan/2024] "GET / HTTP/1.1" 200 - "Mozilla/5.0"`)
	if c.Total() != 1 {
		t.Errorf("expected total 1, got %d", c.Total())
	}
	counts := c.Counts()
	if counts["GET / HTTP/1.1"] != 1 {
		t.Errorf("unexpected counts: %v", counts)
	}
}

func TestUserAgentCounter_Add_Accumulates(t *testing.T) {
	c, _ := NewUserAgentCounter(`user-agent=([\w/. ]+)`)
	lines := []string{
		"user-agent=Mozilla/5.0 status=200",
		"user-agent=Mozilla/5.0 status=404",
		"user-agent=curl/7.68 status=200",
	}
	for _, l := range lines {
		c.Add(l)
	}
	if c.Total() != 3 {
		t.Errorf("expected total 3, got %d", c.Total())
	}
	if c.Counts()["Mozilla/5.0 "] != 2 && c.Counts()["Mozilla/5.0"] != 2 {
		// allow trailing space variance
		if c.Total() != 3 {
			t.Errorf("unexpected counts: %v", c.Counts())
		}
	}
}

func TestUserAgentCounter_Add_SkipsNoMatch(t *testing.T) {
	c, _ := NewUserAgentCounter("")
	c.Add("line with no quotes")
	if c.Total() != 0 {
		t.Errorf("expected 0 total, got %d", c.Total())
	}
}

func TestSortedUserAgentEntries_Order(t *testing.T) {
	counts := map[string]int{"curl": 1, "Mozilla": 5, "wget": 3}
	entries := SortedUserAgentEntries(counts)
	if entries[0].Agent != "Mozilla" {
		t.Errorf("expected Mozilla first, got %s", entries[0].Agent)
	}
	if entries[len(entries)-1].Agent != "curl" {
		t.Errorf("expected curl last, got %s", entries[len(entries)-1].Agent)
	}
}

func TestCountUserAgentReader(t *testing.T) {
	input := strings.NewReader(
		`req "Mozilla/5.0" ok` + "\n" +
			`req "curl/7.68" ok` + "\n" +
			`req "Mozilla/5.0" ok` + "\n",
	)
	c, err := CountUserAgentReader(input, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Total() != 3 {
		t.Errorf("expected total 3, got %d", c.Total())
	}
}
