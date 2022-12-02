package lockfile

import (
	"strings"

	"github.com/uhppoted/uhppoted-lib/config"
)

type Lockfile interface {
	Release()
}

func MakeLockFile(cfg config.Lockfile) (Lockfile, error) {
	switch {
	case strings.HasPrefix(cfg.File, "file:"):
		return makeFileLock(cfg.File)
	}

	return makeFLock(cfg.File, cfg.Remove)
}

func MakeFileFile(file string) (Lockfile, error) {
	return makeFileLock(file)
}
