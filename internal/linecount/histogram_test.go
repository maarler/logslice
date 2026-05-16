package linecount

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse("2006-01-02T15:04", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestTimeHistogram_Empty(t *testing.T) {
	h := NewTimeHistogram(time.Minute)
	if got := h.Buckets(); len(got) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(got))
	}
}

func TestTimeHistogram_Add_SingleBucket(t *testing.T) {
	h := NewTimeHistogram(time.Minute)
	ts := mustTime("2024-01-15T10:05")
	h.Add(ts, 100)
	h.Add(ts.Add(30*time.Second), 200)

	buckets := h.Buckets()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Count != 2 {
		t.Errorf("expected count 2, got %d", buckets[0].Count)
	}
	if buckets[0].Bytes != 300 {
		t.Errorf("expected bytes 300, got %d", buckets[0].Bytes)
	}
}

func TestTimeHistogram_Add_MultipleBuckets(t *testing.T) {
	h := NewTimeHistogram(time.Hour)
	h.Add(mustTime("2024-01-15T09:00"), 50)
	h.Add(mustTime("2024-01-15T10:30"), 80)
	h.Add(mustTime("2024-01-15T10:45"), 90)

	buckets := h.Buckets()
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if buckets[0].Count != 1 {
		t.Errorf("first bucket count: want 1, got %d", buckets[0].Count)
	}
	if buckets[1].Count != 2 {
		t.Errorf("second bucket count: want 2, got %d", buckets[1].Count)
	}
}

func TestTimeHistogram_Buckets_SortedChronologically(t *testing.T) {
	h := NewTimeHistogram(time.Minute)
	h.Add(mustTime("2024-01-15T10:03"), 1)
	h.Add(mustTime("2024-01-15T10:01"), 1)
	h.Add(mustTime("2024-01-15T10:02"), 1)

	buckets := h.Buckets()
	for i := 1; i < len(buckets); i++ {
		if buckets[i].Label < buckets[i-1].Label {
			t.Errorf("buckets not sorted: %s before %s", buckets[i-1].Label, buckets[i].Label)
		}
	}
}

func TestWriteHistogramReport_Empty(t *testing.T) {
	h := NewTimeHistogram(time.Minute)
	var buf bytes.Buffer
	WriteHistogramReport(&buf, h, 40)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data', got: %s", buf.String())
	}
}

func TestWriteHistogramReport_ContainsLabel(t *testing.T) {
	h := NewTimeHistogram(time.Minute)
	h.Add(mustTime("2024-01-15T10:05"), 100)
	var buf bytes.Buffer
	WriteHistogramReport(&buf, h, 40)
	out := buf.String()
	if !strings.Contains(out, "2024-01-15 10:05") {
		t.Errorf("expected bucket label in output, got:\n%s", out)
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected count in output, got:\n%s", out)
	}
}

func TestBucketFormat_Daily(t *testing.T) {
	f := bucketFormat(24 * time.Hour)
	if f != "2006-01-02" {
		t.Errorf("unexpected daily format: %s", f)
	}
}
