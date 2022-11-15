package lockfile

type Lockfile interface {
	Release()
}

func MakeLockFile(file string) (Lockfile, error) {
	return makeFileLock(file)
}
