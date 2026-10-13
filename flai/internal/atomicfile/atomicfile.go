// Package atomicfile replaces a file in one step (S-0100): whoever reads it at
// the same moment sees the old content or the new, never an empty or partial
// file, as os.WriteFile, which truncates first, allows.
package atomicfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

// WriteFile writes data to a temporary file beside path and renames it over
// path. The temporary file is hidden and does not end in .md, so nothing that
// lists a folder's documents takes it for one while it exists. It syncs the
// file before the rename and the folder after it, so a crash keeps the old
// content or the new, never an empty file (I-0137).
func WriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	fail := func(err error) error {
		_ = tmp.Close()
		_ = os.Remove(name)
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return fail(err)
	}
	if err := tmp.Chmod(perm); err != nil {
		return fail(err)
	}
	if err := tmp.Sync(); err != nil {
		return fail(err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return err
	}
	if err := syncDir(dir); err != nil {
		return fmt.Errorf("%s is written, but syncing its folder failed, so a crash may undo it: %w", path, err)
	}
	return nil
}

// syncDir makes the entries of dir, and so a rename into it, durable.
// Windows cannot sync a folder, and some file systems answer EINVAL or
// ENOTSUP for one; there the rename is as durable as the platform makes it,
// so those are not errors.
func syncDir(dir string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	err = d.Sync()
	if errors.Is(err, syscall.EINVAL) || errors.Is(err, errors.ErrUnsupported) {
		err = nil
	}
	if cerr := d.Close(); err == nil {
		err = cerr
	}
	return err
}
