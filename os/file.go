package os

import (
	"io"
	sys "os"
	"path/filepath"
	"strings"
)

/**
 * Replacement implementation for os.Rename that copies the file if os.Rename
 * fails with an error.
 *
 * Use os.Rename to 'move' a file between filesystems on different filesystems fails
 * with an 'invalid cross-device link' error. The replacement implementation first
 * attempts a 'rename' and if that fails, creates a temporary file adjacent to the
 * destination file, copies the source to the temporary file and then does an
 * os.Rename(...) on the temporary file to overwrite the destination file.
 *
 * The temporary file and original file are deleted.
 *
 * NB: this is very much an application specific implementation - errors deleting the
 *     temporary or original file are discarded on the basis that this operation is
 *     typically used to copy create an updated 'working file' from a temporary file
 *     and the application can and will use the updated working file even if the original
 *     file is not deleted.
 *
 * Ref. https://github.com/uhppoted/uhppoted-httpd/issues/20
 *
 * (interim implementation pending resolution of https://github.com/golang/go/issues/41487)
 */
func Rename(oldpath, newpath string) error {
	if err := sys.Rename(oldpath, newpath); err == nil {
		return nil
	}

	dir := filepath.Dir(newpath)
	ext := filepath.Ext(newpath)
	base := filepath.Base(newpath)
	tmpfile := strings.TrimSuffix(base, ext)

	tmp, err := sys.CreateTemp(dir, tmpfile+"*")
	if err != nil {
		return err
	} else {
		defer sys.Remove(tmp.Name())
	}

	src, err := sys.Open(oldpath)
	if err != nil {
		return err
	} else {
		defer src.Close()
	}

	if _, err := io.Copy(tmp, src); err != nil {
		return err
	} else if err := tmp.Sync(); err != nil {
		return err
	}

	src.Close()
	tmp.Close()

	if err := sys.Rename(tmp.Name(), newpath); err != nil {
		return err
	}

	sys.Remove(oldpath)

	return nil
}
