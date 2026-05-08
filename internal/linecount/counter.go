// Package linecount provides utilities for counting lines in a reader
// and estimating progress through a log file.
package linecount

import (
	"bufio"
	"io"
	"os"
)

// Result holds the outcome of a line-count operation.
type Result struct {
	Lines int64
	Bytes int64
}

// CountReader scans r and returns the number of newline-terminated lines
// and total bytes read. r is consumed entirely.
func CountReader(r io.Reader) (Result, error) {
	var res Result
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		res.Lines++
		res.Bytes += int64(len(scanner.Bytes())) + 1 // +1 for newline
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res, nil
}

// CountFile opens the named file, counts its lines, and closes it.
func CountFile(path string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, err
	}
	defer f.Close()
	return CountReader(f)
}

// Fraction returns the fraction of total lines represented by done,
// clamped to [0.0, 1.0]. Returns 0 when total is zero.
func Fraction(done, total int64) float64 {
	if total <= 0 {
		return 0
	}
	f := float64(done) / float64(total)
	if f > 1.0 {
		return 1.0
	}
	return f
}
