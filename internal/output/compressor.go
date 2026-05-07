package output

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// Compression represents the type of compression to apply to output files.
type Compression int

const (
	CompressionNone Compression = iota
	CompressionGzip
)

// ParseCompression converts a string to a Compression value.
func ParseCompression(s string) (Compression, error) {
	switch s {
	case "", "none":
		return CompressionNone, nil
	case "gzip", "gz":
		return CompressionGzip, nil
	default:
		return CompressionNone, fmt.Errorf("unknown compression type %q: must be none or gzip", s)
	}
}

// Extension returns the file extension associated with the compression type.
func (c Compression) Extension() string {
	switch c {
	case CompressionGzip:
		return ".gz"
	default:
		return ""
	}
}

// String returns a human-readable name for the compression type.
func (c Compression) String() string {
	switch c {
	case CompressionGzip:
		return "gzip"
	default:
		return "none"
	}
}

// NewWriter wraps the given *os.File with a compression writer if needed.
// The caller must close the returned io.WriteCloser before closing the file
// to ensure all compressed bytes are flushed.
func NewWriter(f *os.File, c Compression) (io.WriteCloser, error) {
	switch c {
	case CompressionGzip:
		gw, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
		if err != nil {
			return nil, fmt.Errorf("create gzip writer: %w", err)
		}
		return gw, nil
	default:
		return &nopWriteCloser{f}, nil
	}
}

// nopWriteCloser wraps a writer with a no-op Close so it satisfies io.WriteCloser
// without closing the underlying file prematurely.
type nopWriteCloser struct {
	w io.Writer
}

func (n *nopWriteCloser) Write(p []byte) (int, error) { return n.w.Write(p) }
func (n *nopWriteCloser) Close() error               { return nil }
