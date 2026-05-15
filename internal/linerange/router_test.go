package linerange

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRouter_InvalidSpec(t *testing.T) {
	_, err := NewRouter(
		map[string]string{"bad": "abc"},
		map[string]io.Writer{"bad": &bytes.Buffer{}},
	)
	if err == nil {
		t.Fatal("expected error for invalid spec")
	}
}

func TestNewRouter_MissingWriter(t *testing.T) {
	_, err := NewRouter(
		map[string]string{"head": "1:5"},
		map[string]io.Writer{},
	)
	if err == nil {
		t.Fatal("expected error when writer is missing")
	}
}

func TestRouter_Route_BasicRange(t *testing.T) {
	headBuf := &bytes.Buffer{}
	tailBuf := &bytes.Buffer{}

	rt, err := NewRouter(
		map[string]string{"head": "1:3", "tail": "4:6"},
		map[string]io.Writer{"head": headBuf, "tail": tailBuf},
	)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}

	input := "line1\nline2\nline3\nline4\nline5\nline6\n"
	if err := rt.Route(strings.NewReader(input)); err != nil {
		t.Fatalf("Route: %v", err)
	}

	for _, tc := range []struct {
		name string
		buf  *bytes.Buffer
		want string
	}{
		{"head", headBuf, "line1\nline2\nline3\n"},
		{"tail", tailBuf, "line4\nline5\nline6\n"},
	} {
		if got := tc.buf.String(); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestRouter_Route_OverlappingRanges(t *testing.T) {
	aBuf := &bytes.Buffer{}
	bBuf := &bytes.Buffer{}

	rt, err := NewRouter(
		map[string]string{"a": "1:4", "b": "3:5"},
		map[string]io.Writer{"a": aBuf, "b": bBuf},
	)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}

	input := "l1\nl2\nl3\nl4\nl5\n"
	if err := rt.Route(strings.NewReader(input)); err != nil {
		t.Fatalf("Route: %v", err)
	}

	// lines 3 and 4 should appear in both buffers
	if !strings.Contains(aBuf.String(), "l3") || !strings.Contains(bBuf.String(), "l3") {
		t.Errorf("line 3 should be in both ranges; a=%q b=%q", aBuf.String(), bBuf.String())
	}
}

func TestRouter_Route_EmptyInput(t *testing.T) {
	buf := &bytes.Buffer{}
	rt, err := NewRouter(
		map[string]string{"r": "1:10"},
		map[string]io.Writer{"r": buf},
	)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	if err := rt.Route(strings.NewReader("")); err != nil {
		t.Fatalf("Route on empty input: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}
