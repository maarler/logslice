package output

import (
	"compress/gzip"
	"io"
	"os"
	"testing"
)

func TestParseCompression(t *testing.T) {
	cases := []struct {
		input   string
		want    Compression
		wantErr bool
	}{
		{"", CompressionNone, false},
		{"none", CompressionNone, false},
		{"gzip", CompressionGzip, false},
		{"gz", CompressionGzip, false},
		{"zstd", CompressionNone, true},
	}
	for _, tc := range cases {
		got, err := ParseCompression(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseCompression(%q): expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseCompression(%q): unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseCompression(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestCompression_Extension(t *testing.T) {
	if CompressionNone.Extension() != "" {
		t.Errorf("expected empty extension for none")
	}
	if CompressionGzip.Extension() != ".gz" {
		t.Errorf("expected .gz extension for gzip")
	}
}

func TestNewWriter_None_WritesPlaintext(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	wc, err := NewWriter(f, CompressionNone)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(wc, "hello\n")
	_ = wc.Close()

	_, _ = f.Seek(0, io.SeekStart)
	data, _ := io.ReadAll(f)
	if string(data) != "hello\n" {
		t.Errorf("expected plain text, got %q", data)
	}
}

func TestNewWriter_Gzip_WritesCompressed(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	wc, err := NewWriter(f, CompressionGzip)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(wc, "compressed line\n")
	_ = wc.Close()

	_, _ = f.Seek(0, io.SeekStart)
	gr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("not a valid gzip stream: %v", err)
	}
	defer gr.Close()
	data, _ := io.ReadAll(gr)
	if string(data) != "compressed line\n" {
		t.Errorf("decompressed content mismatch: got %q", data)
	}
}
