package lockfile

import (
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-lib/config"
)

type Lockfile interface {
	Release()
}

func MakeLockFile(cfg config.Lockfile) (Lockfile, error) {
	switch {
	case strings.HasPrefix(cfg.File, "file:"):
		return makeFileLock(cfg.File)

	case strings.HasPrefix(cfg.File, "soft:"):
		return makeSoftLock(cfg.File, cfg.Interval, cfg.Wait)
	}

	return makeFLock(cfg.File, cfg.Remove)
}

func MakeFileFile(file string) (Lockfile, error) {
	return makeFileLock(file)
}

func MakeSoftFileLock(file string, interval time.Duration, wait time.Duration) (Lockfile, error) {
	return makeSoftLock(file, interval, wait)
}
