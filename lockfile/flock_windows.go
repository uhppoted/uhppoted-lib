package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
)

type flock struct {
	file string
}

// Windows doesn't have 'flock' so fall back to ordinary file lock
func makeFLock(file string) (*flock, error) {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if _, err := os.Stat(file); err == nil {
		return nil, fmt.Errorf("lockfile '%v' already in use", file)
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if err := os.WriteFile(file, []byte(pid), 0644); err != nil {
		return nil, err
	}

	return &flock{
		file: file,
	}, nil
}

func (l *flock) Release() {
	if l != nil {
		os.Remove(l.file)
	}
}
