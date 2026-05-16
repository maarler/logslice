package linecount

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHourlyCount_EmptyReader(t *testing.T) {
	var out bytes.Buffer
	opts := DefaultHourlyCmdOptions()
	opts.Out = &out

	if err := RunHourlyCount(strings.NewReader(""), opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "No data") {
		t.Errorf("expected 'No data', got %q", out.String())
	}
}

func TestRunHourlyCount_SkipsUnparseable(t *testing.T) {
	input := "this line has no timestamp\nalso no timestamp here\n"
	var out bytes.Buffer
	opts := DefaultHourlyCmdOptions()
	opts.Out = &out

	if err := RunHourlyCount(strings.NewReader(input), opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "No data") {
		t.Errorf("expected 'No data' for all-unparseable input, got %q", out.String())
	}
}

func TestRunHourlyCount_ParsesTimestamps(t *testing.T) {
	input := strings.Join([]string{
		"2024-03-10T14:01:00Z INFO starting",
		"2024-03-10T14:32:00Z INFO processing",
		"2024-03-10T15:05:00Z WARN slow query",
	}, "\n") + "\n"

	var out bytes.Buffer
	opts := DefaultHourlyCmdOptions()
	opts.Out = &out

	if err := RunHourlyCount(strings.NewReader(input), opts); err != nil {
		t.Fatal(err)
	}

	body := out.String()
	if !strings.Contains(body, "2024-03-10 14:00") {
		t.Errorf("expected 14:00 bucket, output:\n%s", body)
	}
	if !strings.Contains(body, "2024-03-10 15:00") {
		t.Errorf("expected 15:00 bucket, output:\n%s", body)
	}
}

func TestRunHourlyCount_HeaderPresent(t *testing.T) {
	input := "2024-03-10T09:00:00Z DEBUG boot\n"
	var out bytes.Buffer
	opts := DefaultHourlyCmdOptions()
	opts.Out = &out

	if err := RunHourlyCount(strings.NewReader(input), opts); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Hour") {
		t.Errorf("expected header row, got:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "Lines") {
		t.Errorf("expected 'Lines' column, got:\n%s", out.String())
	}
}
