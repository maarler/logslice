package linecount

import (
	"strings"
	"testing"
)

func TestNewPatternCounter_InvalidRegex(t *testing.T) {
	_, err := NewPatternCounter(map[string]string{"bad": "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestPatternCounter_Add_SinglePattern(t *testing.T) {
	pc, err := NewPatternCounter(map[string]string{"error": `(?i)error`})
	if err != nil {
		t.Fatal(err)
	}
	pc.Add("2024/01/01 ERROR something broke")
	pc.Add("2024/01/01 INFO all good")
	pc.Add("2024/01/01 error again")

	counts := pc.Counts()
	if counts["error"] != 2 {
		t.Fatalf("expected 2 error hits, got %d", counts["error"])
	}
}

func TestPatternCounter_Add_MultiplePatterns(t *testing.T) {
	pc, err := NewPatternCounter(map[string]string{
		"warn":  `(?i)warn`,
		"error": `(?i)error`,
	})
	if err != nil {
		t.Fatal(err)
	}
	lines := []string{
		"ERROR critical failure",
		"WARN disk low",
		"INFO startup complete",
		"WARN memory low",
	}
	for _, l := range lines {
		pc.Add(l)
	}
	counts := pc.Counts()
	if counts["error"] != 1 {
		t.Errorf("error: want 1 got %d", counts["error"])
	}
	if counts["warn"] != 2 {
		t.Errorf("warn: want 2 got %d", counts["warn"])
	}
}

func TestPatternCounter_Counts_ReturnsCopy(t *testing.T) {
	pc, _ := NewPatternCounter(map[string]string{"x": `x`})
	pc.Add("x")
	c1 := pc.Counts()
	c1["x"] = 999
	c2 := pc.Counts()
	if c2["x"] != 1 {
		t.Fatal("Counts should return a copy, not a reference")
	}
}

func TestCountPatternReader(t *testing.T) {
	input := strings.NewReader("error here\nwarn there\nerror again\ninfo ok\n")
	counts, err := CountPatternReader(input, map[string]string{
		"error": `(?i)error`,
		"warn":  `(?i)warn`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if counts["error"] != 2 {
		t.Errorf("error: want 2 got %d", counts["error"])
	}
	if counts["warn"] != 1 {
		t.Errorf("warn: want 1 got %d", counts["warn"])
	}
}

func TestCountPatternReader_InvalidRegex(t *testing.T) {
	_, err := CountPatternReader(strings.NewReader(""), map[string]string{"bad": "["})
	if err == nil {
		t.Fatal("expected error")
	}
}
