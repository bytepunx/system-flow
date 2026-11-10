package conflictmark

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// keptVersion is the format of the kept marker lines, and of what Lines
// decides: raise it when either changes, so that lines kept before are not
// reused (S-0677).
const keptVersion = 1

// racyWindow is how long after a file's modification time its lines are
// still not kept: a second write within the file system's timestamp
// granularity could leave the same time and size, as mdlint's kept findings
// and the digest's kept problems say. A file read that soon after it changed
// is read again on the next run.
const racyWindow = 2 * time.Second

// readFile reads a file a Kept scans; a test replaces it to count the files
// read.
var readFile = os.ReadFile

// KeptPath is where a Kept keeps the marker lines of the files under root
// between processes in the cache folder cacheDir: one file per root, named
// by a hash of the root's path, so that a story's worktree and the main
// checkout do not overwrite each other's.
func KeptPath(cacheDir, root string) string {
	sum := sha256.Sum256([]byte(filepath.Clean(root)))
	return filepath.Join(cacheDir, fmt.Sprintf("conflict-markers-%x.json", sum[:6]))
}

// Kept scans files under a root for conflict markers, as Lines does, and
// keeps each file's marker lines between processes by the file's size and
// modification time (S-0677): a file whose size and modification time are
// those it was scanned at takes its kept lines unread. Open one with
// OpenKept, scan each file with Lines, then Save. A Kept is for one
// goroutine.
type Kept struct {
	root      string
	path      string // the kept file, empty when nothing is kept
	old       map[string]keptEntry
	now       map[string]keptEntry
	checkedAt time.Time
}

// OpenKept is a Kept for the files under root, reading what was kept at
// KeptPath in the cache folder cacheDir. An empty cacheDir, or one that does
// not exist, keeps nothing: every file is read. A kept file that is missing,
// damaged, or another root's or format's is a miss, never an error.
func OpenKept(cacheDir, root string) *Kept {
	k := &Kept{root: filepath.Clean(root), now: map[string]keptEntry{}, checkedAt: time.Now()}
	if cacheDir == "" {
		return k
	}
	if info, err := os.Stat(cacheDir); err != nil || !info.IsDir() {
		return k
	}
	k.path = KeptPath(cacheDir, root)
	k.old = readKept(k.path, k.root)
	return k
}

// Lines returns the 1-based numbers of the conflict marker lines of the file
// at path, as Lines of its text does: the kept ones when the file's size and
// modification time are those it was scanned at, else those of a read. The
// error is the read's.
func (k *Kept) Lines(path string) ([]int, error) {
	key := k.key(path)
	var info os.FileInfo
	if k.path != "" {
		// A file whose stat fails is read, which fails too, and is
		// reported as before.
		info, _ = os.Stat(path)
	}
	if e, ok := k.old[key]; ok && info != nil && e.Size == info.Size() && e.ModTime == info.ModTime().UnixNano() {
		k.now[key] = e
		return slices.Clone(e.Lines), nil
	}
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	lines := Lines(string(data))
	if info != nil && info.ModTime().Before(k.checkedAt.Add(-racyWindow)) {
		k.now[key] = keptEntry{Size: info.Size(), ModTime: info.ModTime().UnixNano(), Lines: slices.Clone(lines)}
	}
	return lines, nil
}

// Save writes what this run scanned to the kept file, when it differs from
// what was kept: a file not scanned this run drops out. A failed write fails
// nothing; it costs the next run a read.
func (k *Kept) Save() {
	if k.path == "" || sameKept(k.old, k.now) {
		return
	}
	_ = writeKept(k.path, keptFile{Version: keptVersion, Root: k.root, Files: k.now})
}

// key is path from the root, or path itself when it is not under it.
func (k *Kept) key(path string) string {
	if rel, err := filepath.Rel(k.root, path); err == nil && filepath.IsLocal(rel) {
		return filepath.ToSlash(rel)
	}
	return path
}

// keptFile is what KeptPath holds: the root scanned and each file's key and
// marker lines, by its path from the root.
type keptFile struct {
	Version int                  `json:"version"`
	Root    string               `json:"root"`
	Files   map[string]keptEntry `json:"files"`
}

// keptEntry is one file's marker lines and the stat they were found at.
type keptEntry struct {
	Size    int64 `json:"size"`
	ModTime int64 `json:"mtime"` // Unix nanoseconds
	Lines   []int `json:"lines,omitempty"`
}

// sameKept reports whether two sets of kept entries hold the same files at
// the same keys: lines follow from the key, so they are not compared.
func sameKept(a, b map[string]keptEntry) bool {
	if a == nil || len(a) != len(b) {
		return false
	}
	for name, e := range a {
		o, ok := b[name]
		if !ok || o.Size != e.Size || o.ModTime != e.ModTime {
			return false
		}
	}
	return true
}

// readKept is the lines kept at path for root. A file that is missing,
// cannot be read or parsed, is of another format, or was scanned under
// another root, holds none.
func readKept(path, root string) map[string]keptEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var k keptFile
	if err := json.Unmarshal(data, &k); err != nil || k.Version != keptVersion || k.Root != root {
		return nil
	}
	return k.Files
}

// writeKept writes k to path through a file beside it renamed over it, so
// that a reader finds the old lines or the new ones whole. It never makes
// the cache folder.
func writeKept(path string, k keptFile) error {
	data, err := json.Marshal(k)
	if err != nil {
		return fmt.Errorf("encode the kept conflict markers for %s: %w", path, err)
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".conflict-markers-*")
	if err != nil {
		return fmt.Errorf("write the kept conflict markers in %s: %w", dir, err)
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return fmt.Errorf("write the kept conflict markers to %s: %w", name, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("close the kept conflict markers in %s: %w", name, err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("keep the conflict markers at %s: %w", path, err)
	}
	return nil
}
