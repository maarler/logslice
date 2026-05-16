package linecount

import (
	"strings"
	"testing"
)

func TestWriteStatusCodeReport_Empty(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	var sb strings.Builder
	if err := WriteStatusCodeReport(&sb, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "no status codes found") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestWriteStatusCodeReport_ContainsHeaders(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	c.Add("GET /path 200")
	var sb strings.Builder
	if err := WriteStatusCodeReport(&sb, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, hdr := range []string{"CODE", "CLASS", "COUNT", "PERCENT", "BAR"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output: %s", hdr, out)
		}
	}
}

func TestWriteStatusCodeReport_ContainsCode(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	c.Add("GET /path 404")
	var sb strings.Builder
	_ = WriteStatusCodeReport(&sb, c)
	out := sb.String()
	if !strings.Contains(out, "404") {
		t.Errorf("expected 404 in output, got: %s", out)
	}
	if !strings.Contains(out, "4xx") {
		t.Errorf("expected 4xx class in output, got: %s", out)
	}
}

func TestWriteStatusCodeReport_PercentageCalculation(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	for i := 0; i < 3; i++ {
		c.Add("status 200")
	}
	c.Add("status 500")
	var sb strings.Builder
	_ = WriteStatusCodeReport(&sb, c)
	out := sb.String()
	if !strings.Contains(out, "75.0%") {
		t.Errorf("expected 75.0%% for 200, got: %s", out)
	}
	if !strings.Contains(out, "25.0%") {
		t.Errorf("expected 25.0%% for 500, got: %s", out)
	}
}

func TestWriteStatusCodeReport_TotalLine(t *testing.T) {
	c, _ := NewStatusCodeCounter("")
	c.Add("GET /a 200")
	c.Add("GET /b 404")
	var sb strings.Builder
	_ = WriteStatusCodeReport(&sb, c)
	out := sb.String()
	if !strings.Contains(out, "TOTAL") {
		t.Errorf("expected TOTAL line in output, got: %s", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected total count 2 in output, got: %s", out)
	}
}
