package fieldextract_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldextract"
)

func TestNew_Defaults(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{})
	if e == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestExtract_SimpleKeyValue(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{})
	fields := e.Extract(`level=info msg=hello`)
	if fields["level"] != "info" {
		t.Errorf("expected level=info, got %q", fields["level"])
	}
	if fields["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %q", fields["msg"])
	}
}

func TestExtract_QuotedValue(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{})
	fields := e.Extract(`level=error msg="disk full"`)
	if fields["msg"] != "disk full" {
		t.Errorf("expected 'disk full', got %q", fields["msg"])
	}
}

func TestExtract_CustomDelimiterAndSep(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{
		Delimiter: ":",
		PairSep:   ",",
	})
	fields := e.Extract(`host:localhost,port:8080`)
	if fields["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %q", fields["host"])
	}
	if fields["port"] != "8080" {
		t.Errorf("expected port=8080, got %q", fields["port"])
	}
}

func TestExtract_EmptyLine(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{})
	fields := e.Extract("")
	if len(fields) != 0 {
		t.Errorf("expected empty map, got %v", fields)
	}
}

func TestExtract_NoDelimiter(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{})
	fields := e.Extract("this is a plain log line with no kv pairs")
	if len(fields) != 0 {
		t.Errorf("expected empty map, got %v", fields)
	}
}

func TestGet_Found(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{})
	v, ok := e.Get(`level=warn request_id=abc123`, "request_id")
	if !ok {
		t.Fatal("expected key to be found")
	}
	if v != "abc123" {
		t.Errorf("expected abc123, got %q", v)
	}
}

func TestGet_NotFound(t *testing.T) {
	e := fieldextract.New(fieldextract.Options{})
	v, ok := e.Get(`level=info`, "missing")
	if ok {
		t.Errorf("expected key not found, got %q", v)
	}
}
