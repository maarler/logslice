package linecount

import (
	"bufio"
	"io"
	"os"
)

// CountReader counts the number of newline-terminated lines in r.
func CountReader(r io.Reader) (int64, error) {
	var count int64
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		for _, b := range buf[:n] {
			if b == '\n' {
				count++
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

// CountFile counts the number of lines in the file at path.
func CountFile(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return CountReader(bufio.NewReader(f))
}

// Fraction returns the approximate line number that represents the given
// fraction (0.0–1.0) of total lines in the file.
func Fraction(path string, frac float64) (int64, error) {
	total, err := CountFile(path)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, nil
	}
	if frac <= 0 {
		return 0, nil
	}
	if frac >= 1 {
		return total, nil
	}
	return int64(float64(total) * frac), nil
}
