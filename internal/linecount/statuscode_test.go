package linecount

import (
	"strings"
	"testing"
)

func TestNewStatusCodeCounter_DefaultPattern(t *testing.T) {
	c, err := NewStatusCodeCounter("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil counter")
	}
}

func TestNewStatusCodeCounter_InvalidPattern(t *testing.T) {
	_, err := NewStatusCodeCounter("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestStatusCodeCounter_Add_SingleCode(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	c.Add(`192.168.1.1 - - [01/Jan/2024] "GET /api" 200 512`)
	if c.Total() != 1 {
		t.Fatalf("expected total 1, got %d", c.Total())
	}
	entries := c.SortedCodeEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Code != 200 {
		t.Errorf("expected code 200, got %d", entries[0].Code)
	}
	if entries[0].Class != "2xx" {
		t.Errorf("expected class 2xx, got %s", entries[0].Class)
	}
}

func TestStatusCodeCounter_Add_MultipleCodes(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	lines := []string{
		`GET /a 200`,
		`GET /b 200`,
		`GET /c 404`,
		`GET /d 500`,
	}
	for _, l := range lines {
		c.Add(l)
	}
	if c.Total() != 4 {
		t.Fatalf("expected total 4, got %d", c.Total())
	}
	entries := c.SortedCodeEntries()
	if entries[0].Code != 200 || entries[0].Count != 2 {
		t.Errorf("expected 200 with count 2, got %d with %d", entries[0].Code, entries[0].Count)
	}
}

func TestStatusCodeCounter_Add_NoMatch(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	c.Add("no status code here")
	if c.Total() != 0 {
		t.Errorf("expected total 0, got %d", c.Total())
	}
}

func TestStatusCodeCounter_SortedEntries_DescendingCount(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	for i := 0; i < 5; i++ {
		c.Add("request 200")
	}
	c.Add("request 404")
	c.Add("request 404")
	c.Add("request 500")
	entries := c.SortedCodeEntries()
	if entries[0].Code != 200 {
		t.Errorf("expected 200 first, got %d", entries[0].Code)
	}
	if entries[1].Code != 404 {
		t.Errorf("expected 404 second, got %d", entries[1].Code)
	}
}

func TestCountStatusCodeReader(t *testing.T) {
	input := strings.Join([]string{
		"GET /a 200",
		"GET /b 301",
		"GET /c 404",
		"GET /d 503",
	}, "\n")
	c, err := CountStatusCodeReader(strings.NewReader(input), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Total() != 4 {
		t.Errorf("expected 4 total, got %d", c.Total())
	}
}
