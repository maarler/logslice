//go:build !windows

package rotate

import (
	"fmt"
	"syscall"
)

func statFile(path string) (FileState, error) {
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return FileState{}, fmt.Errorf("rotate: stat %q: %w", path, err)
	}
	return FileState{
		Inode: st.Ino,
		Size:  st.Size,
	}, nil
}
