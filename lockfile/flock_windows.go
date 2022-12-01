package lockfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type flock struct {
	file *os.File
	remove bool
}

var (
	kernel32, _       = syscall.LoadLibrary("kernel32.dll")
	procLockFile, _   = syscall.GetProcAddress(kernel32, "LockFile")
	procUnlockFile, _ = syscall.GetProcAddress(kernel32, "UnlockFile")
)

// Windows doesn't have 'flock' so use LockFile/UnlockFile API
func makeFLock(file string, remove bool) (*flock, error) {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}

	handle := syscall.Handle(f.Fd())
	if err := lock(handle); err != nil {
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
		file: f,
		remove:remove,
	}, nil
}

// NTS
// Unlike Linux and Darwin, should actuall remove lockfile because the LockFile and UnlockFile calls
// acquire an exclusive lock. Keeping it configurable though.
func (l flock) Release() {
	handle := syscall.Handle(l.file.Fd())

	unlock(handle)
	l.file.Close()

	if l.remove {
		os.Remove(l.file.Name())		
	}
}

// Ref. https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-lockfile
//
// BOOL LockFile(
//
//	[in] HANDLE hFile,
//	[in] DWORD  dwFileOffsetLow,
//	[in] DWORD  dwFileOffsetHigh,
//	[in] DWORD  nNumberOfBytesToLockLow,
//	[in] DWORD  nNumberOfBytesToLockHigh
//
// );
//
// Ref. https://github.com/gofrs/flock/blob/master/flock_winapi.go
func lock(handle syscall.Handle) error {
	var dwFileOffsetLow uint32 = 0
	var dwFileOffsetHigh uint32 = 0
	var nNumberOfBytesToLockLow uint32 = 0xfffffff  // MAX_DWORD
	var nNumberOfBytesToLockHigh uint32 = 0xfffffff // MAX_DWORD

	rc, _, err := syscall.Syscall6(
		uintptr(procLockFile),
		5,
		uintptr(handle),
		uintptr(dwFileOffsetLow),          // [in] DWORD  dwFileOffsetLow
		uintptr(dwFileOffsetHigh),         // [in] DWORD  dwFileOffsetHigh
		uintptr(nNumberOfBytesToLockLow),  // [in] DWORD  nNumberOfBytesToLockLow
		uintptr(nNumberOfBytesToLockHigh), // [in] DWORD  nNumberOfBytesToLockHigh
		uintptr(0))                        // unused

	if rc != 1 && err == 0 {
		return syscall.EINVAL
	} else if rc != 1 {
		return fmt.Errorf("%v", err)
	} else {
		return nil
	}
}

// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-unlockfile
//
// BOOL UnlockFile(
//
//	[in] HANDLE hFile,
//	[in] DWORD  dwFileOffsetLow,
//	[in] DWORD  dwFileOffsetHigh,
//	[in] DWORD  nNumberOfBytesToUnlockLow,
//	[in] DWORD  nNumberOfBytesToUnlockHigh
//
// );
//
// Ref. https://github.com/gofrs/flock/blob/master/flock_winapi.go
func unlock(handle syscall.Handle) error {
	var dwFileOffsetLow uint32 = 0
	var dwFileOffsetHigh uint32 = 0
	var nNumberOfBytesToLockLow uint32 = 0xfffffff  // MAX_DWORD
	var nNumberOfBytesToLockHigh uint32 = 0xfffffff // MAX_DWORD

	rc, _, err := syscall.Syscall6(
		uintptr(procUnlockFile),
		5,
		uintptr(handle),
		uintptr(dwFileOffsetLow),          // [in] DWORD  dwFileOffsetLow
		uintptr(dwFileOffsetHigh),         // [in] DWORD  dwFileOffsetHigh
		uintptr(nNumberOfBytesToLockLow),  // [in] DWORD  nNumberOfBytesToLockLow
		uintptr(nNumberOfBytesToLockHigh), // [in] DWORD  nNumberOfBytesToLockHigh
		uintptr(0))                        // unused

	if rc != 1 && err == 0 {
		return syscall.EINVAL
	} else if rc != 1 {
		return fmt.Errorf("%v", err)
	} else {
		return nil
	}
}
