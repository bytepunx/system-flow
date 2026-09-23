// Package atomicfile replaces a file in one step (S-0100): whoever reads it at
// the same moment sees the old content or the new, never an empty or partial
// file, as os.WriteFile, which truncates first, allows.
package atomicfile

import (
	"os"
	"path/filepath"
)

// WriteFile writes data to a temporary file beside path and renames it over
// path. The temporary file is hidden and does not end in .md, so nothing that
// lists a folder's documents takes it for one while it exists.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.tmp")
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
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return err
	}
	return nil
}
