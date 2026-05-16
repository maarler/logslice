package linecount

import (
	"strings"
	"testing"
)

func TestIPCounter_Empty(t *testing.T) {
	c := NewIPCounter()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 total, got %d", got)
	}
	if got := len(c.Counts()); got != 0 {
		t.Fatalf("expected empty counts, got %d entries", got)
	}
}

func TestIPCounter_Add_SingleIP(t *testing.T) {
	c := NewIPCounter()
	c.Add("2024-01-01 INFO request from 192.168.1.1 accepted")
	if c.Total() != 1 {
		t.Fatalf("expected total 1, got %d", c.Total())
	}
	if c.Counts()["192.168.1.1"] != 1 {
		t.Fatalf("expected count 1 for 192.168.1.1")
	}
}

func TestIPCounter_Add_MultipleIPs(t *testing.T) {
	c := NewIPCounter()
	c.Add("connect 10.0.0.1 -> 10.0.0.2")
	if c.Total() != 2 {
		t.Fatalf("expected total 2, got %d", c.Total())
	}
}

func TestIPCounter_Add_Accumulates(t *testing.T) {
	c := NewIPCounter()
	for i := 0; i < 5; i++ {
		c.Add("request from 172.16.0.5")
	}
	if c.Counts()["172.16.0.5"] != 5 {
		t.Fatalf("expected 5, got %d", c.Counts()["172.16.0.5"])
	}
}

func TestIPCounter_Add_IgnoresInvalidIP(t *testing.T) {
	c := NewIPCounter()
	c.Add("version 999.999.999.999 released")
	if c.Total() != 0 {
		t.Fatalf("expected 0 total for invalid IP, got %d", c.Total())
	}
}

func TestSortedIPEntries_Order(t *testing.T) {
	counts := map[string]int{
		"10.0.0.1": 3,
		"10.0.0.2": 7,
		"10.0.0.3": 1,
	}
	entries := SortedIPEntries(counts)
	if entries[0].IP != "10.0.0.2" {
		t.Fatalf("expected 10.0.0.2 first, got %s", entries[0].IP)
	}
	if entries[2].IP != "10.0.0.3" {
		t.Fatalf("expected 10.0.0.3 last, got %s", entries[2].IP)
	}
}

func TestCountIPReader(t *testing.T) {
	input := strings.NewReader("req from 1.2.3.4\nreq from 1.2.3.4\nreq from 5.6.7.8\n")
	c := CountIPReader(input)
	if c.Counts()["1.2.3.4"] != 2 {
		t.Fatalf("expected 2 for 1.2.3.4, got %d", c.Counts()["1.2.3.4"])
	}
	if c.Counts()["5.6.7.8"] != 1 {
		t.Fatalf("expected 1 for 5.6.7.8, got %d", c.Counts()["5.6.7.8"])
	}
}
