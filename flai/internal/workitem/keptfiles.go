package workitem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// KeptFiles keeps the parse of files of one kind by absolute path, so that a
// long-running flai (serve, mcp) reads again only the files that changed
// (S-0602), as the item store keeps items (S-0156). Threads, conversations,
// digest entries, issues, and activity documents each keep one, process-wide,
// since the dashboard's methods open a Repo per request. A file's parse is
// kept while the file is the same file with the modification time and size it
// had when it was read, and that time was older than the read by racyWindow.
// Values are handed out through the clone function, never the kept one.
type KeptFiles[T any] struct {
	parse func(path string) (T, error)
	clone func(T) T

	mu      sync.Mutex
	folders map[string]*filesFolder[T]
	parses  int // under mu
}

// filesFolder is one folder's kept files. Its lock is held while it is looked
// at, so that readers at once share one parse of each changed file.
type filesFolder[T any] struct {
	mu    sync.Mutex
	files map[string]*filesKept[T] // by base name
}

// filesKept is one file's kept parse and the stat it was read at.
type filesKept[T any] struct {
	info    os.FileInfo
	value   T    // never handed out; callers get a clone
	trusted bool // read later than racyWindow after it was modified
}

// KeptFile is one file a List looked at.
type KeptFile[T any] struct {
	Path   string // the folder joined with the file's name
	Value  T      // a clone; zero when Err is set
	Err    error  // the parse's error; not kept, so the next look parses again
	Parsed bool   // parsed by this call rather than taken from what was kept
}

// NewKeptFiles makes an empty store for one kind of file. parse reads and
// parses the file at path, and clone gives each caller its own copy of a kept
// value.
func NewKeptFiles[T any](parse func(path string) (T, error), clone func(T) T) *KeptFiles[T] {
	return &KeptFiles[T]{parse: parse, clone: clone, folders: map[string]*filesFolder[T]{}}
}

// Parses counts the files parsed since k was made, for tests.
func (k *KeptFiles[T]) Parses() int {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.parses
}

func (k *KeptFiles[T]) folder(dir string) *filesFolder[T] {
	k.mu.Lock()
	defer k.mu.Unlock()
	f := k.folders[dir]
	if f == nil {
		f = &filesFolder[T]{files: map[string]*filesKept[T]{}}
		k.folders[dir] = f
	}
	return f
}

func (k *KeptFiles[T]) parsed() {
	k.mu.Lock()
	k.parses++
	k.mu.Unlock()
}

// keyPath cleans path and makes it absolute, so that one file has one key.
func keyPath(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return filepath.Clean(path)
}

// List looks at every regular file of dir whose base name match accepts (nil
// accepts all), in name order, parsing again only those that are new or may
// have changed. A folder that does not exist, or whose path names a file,
// has none. Kept files of dir that
// are gone, or that match no longer accepts, are dropped.
func (k *KeptFiles[T]) List(dir string, match func(name string) bool) ([]KeptFile[T], error) {
	dir = keyPath(dir)
	f := k.folder(dir)
	f.mu.Lock()
	defer f.mu.Unlock()
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOTDIR) {
		clear(f.files)
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading the folder %s: %w", dir, err)
	}
	seen := make(map[string]bool, len(entries))
	var out []KeptFile[T]
	for _, e := range entries {
		name := e.Name()
		if !e.Type().IsRegular() || (match != nil && !match(name)) {
			continue
		}
		seen[name] = true
		path := filepath.Join(dir, name)
		v, parsed, err := k.look(f, dir, name)
		if errors.Is(err, os.ErrNotExist) {
			delete(seen, name) // gone since the folder was read
			continue
		}
		kf := KeptFile[T]{Path: path, Parsed: parsed, Err: err}
		if err == nil {
			kf.Value = v
		}
		out = append(out, kf)
	}
	for name := range f.files {
		if !seen[name] {
			delete(f.files, name)
		}
	}
	return out, nil
}

// Read returns the parse of one file, from what is kept while it holds. A
// stat error is returned as is, so that errors.Is(err, os.ErrNotExist) says
// the file is gone, and a file gone is dropped.
func (k *KeptFiles[T]) Read(path string) (T, error) {
	path = keyPath(path)
	dir, name := filepath.Dir(path), filepath.Base(path)
	f := k.folder(dir)
	f.mu.Lock()
	defer f.mu.Unlock()
	v, _, err := k.look(f, dir, name)
	return v, err
}

// look returns a clone of name's value and whether this call parsed it,
// parsing the file again only when it is not the file kept or may have
// changed since. A stat error is returned as is. The caller holds f.mu.
func (k *KeptFiles[T]) look(f *filesFolder[T], dir, name string) (T, bool, error) {
	var zero T
	readAt := time.Now()
	info, err := os.Stat(filepath.Join(dir, name))
	if err != nil {
		delete(f.files, name)
		return zero, false, err
	}
	if kf := f.files[name]; kf != nil && kf.trusted && os.SameFile(kf.info, info) &&
		kf.info.ModTime().Equal(info.ModTime()) && kf.info.Size() == info.Size() {
		return k.clone(kf.value), false, nil
	}
	v, err := k.parse(filepath.Join(dir, name))
	k.parsed()
	if err != nil {
		delete(f.files, name)
		return zero, true, err
	}
	f.files[name] = &filesKept[T]{info: info, value: v, trusted: info.ModTime().Before(readAt.Add(-racyWindow))}
	return k.clone(v), true, nil
}
