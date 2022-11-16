package lockfile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type softlock struct {
	file     string
	released bool
	interval time.Duration
	wait     time.Duration
	sync.Mutex
}

func makeSoftLock(file string, interval time.Duration, wait time.Duration) (*softlock, error) {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	l := softlock{
		file:     file,
		interval: interval,
		wait:     wait,
	}

	if err := l.lock(interrupt); err != nil {
		return nil, err
	}

	tick := time.Tick(l.interval)
	go func() {
		for !l.released {
			<-tick
			// log.Infof(LOG_TAG, "touching lockfile '%v'", lockfile)
			l.touch()
		}
	}()

	return &l, nil
}

func (l *softlock) Release() {
	if l != nil {
		l.Lock()
		defer l.Unlock()

		l.released = true
		os.Remove(l.file)
	}
}

// Uses SHA-256 hash of lockfile contents because os.Stat updates the mtime
// of the lockfile and in any event, mtime only has a resolution of 1 minute.
func (l *softlock) lock(interrupt chan os.Signal) error {
	checksum, err := l.hash()

	switch {
	case err != nil && !os.IsNotExist(err):
		return err

	case err != nil && os.IsNotExist(err):
		if err := l.touch(); err != nil {
			return err
		}

	case err == nil && checksum == "":
		return fmt.Errorf("invalid lockfile checksum")

	case err == nil && checksum != "":
		// log.Warnf(LOG_TAG, "'soft lock' file '%v' exists, delaying for %v", lockfile, l.wait)

		wait := time.After(l.wait)

		select {
		case <-wait:

		case <-interrupt:
			return fmt.Errorf("interrupted")
		}

		h, err := l.hash()
		switch {
		case err != nil && !os.IsNotExist(err):
			return err

		case err != nil && os.IsNotExist(err):
			if err := l.touch(); err != nil {
				return err
			}

		case h != checksum:
			return fmt.Errorf("client lockfile '%v' in use", l.file)

		default:
			// log.Warnf(LOG_TAG, "replacing lockfile '%v'", l.file)
			if err := l.touch(); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("failed to acquire client lock")
	}

	return nil
}

func (l *softlock) hash() (string, error) {
	if bytes, err := os.ReadFile(l.file); err != nil {
		return "", err
	} else {
		sum := sha256.Sum256(bytes)

		return hex.EncodeToString(sum[:]), nil
	}
}

func (l *softlock) touch() error {
	l.Lock()
	defer l.Unlock()

	if l.released {
		return nil
	}

	pid := fmt.Sprintf("%d", os.Getpid())
	now := time.Now().Format("2006-01-02 15:04:05")
	v := fmt.Sprintf("%v\n%v\n", pid, now)

	return os.WriteFile(l.file, []byte(v), 0644)
}
