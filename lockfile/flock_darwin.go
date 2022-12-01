package lockfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type flock struct {
	file   *os.File
	remove bool
}

// Use 'flock' to manage file locks
//
// Ref. https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/flock.2.html
func makeFLock(file string, remove bool) (*flock, error) {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}

	handle := int(f.Fd())
	if err := syscall.Flock(handle, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, fmt.Errorf("lockfile '%v' in use (%v)", file, err)
		} else {
			return nil, fmt.Errorf("failed to lock '%v' (%v)", file, err)
		}
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if _, err := f.Write([]byte(pid)); err != nil {
		return nil, err
	} else if err := f.Sync(); err != nil {
		return nil, err
	}

	return &flock{
		file:   f,
		remove: remove,
	}, nil
}

// NTS
// Does not remove the lockfile unless explicitly configured to do so because another process may open
// it in blocking mode, in which case deleting the lockfile allows a second process to use the "same"
// lockfile and not block (because the lock if on the fd, not the file name). Which of course means you
// can' use a mixture of blocking flocks and filelocks, but so be it.
//
// Ref. https://stackoverflow.com/questions/17708885/flock-removing-locked-file-without-race-condition
func (l flock) Release() {
	handle := int(l.file.Fd())
	syscall.Flock(handle, syscall.LOCK_UN)
	l.file.Close()

	if l.remove {
		os.Remove(l.file.Name())
	}
}
