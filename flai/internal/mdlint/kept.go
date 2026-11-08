package mdlint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// keptVersion is the format of the kept findings file; a file of another
// format is a miss.
const keptVersion = 1

// racyWindow is how long after a file's modification time its findings are
// still not kept: a second write within the file system's timestamp
// granularity could leave the same time and size, as the item store says
// (workitem's racyWindow). A file linted that soon after it changed is
// linted again on the next run.
const racyWindow = 2 * time.Second

// KeptPath is where the folder lint of dir keeps its findings between
// processes in the cache folder cacheDir: one file per folder, named by a
// hash of the folder's path, so that a story's worktree and the main checkout
// do not overwrite each other's.
func KeptPath(cacheDir, dir string) string {
	sum := sha256.Sum256([]byte(filepath.Clean(dir)))
	return filepath.Join(cacheDir, "lint-findings-"+hex.EncodeToString(sum[:6])+".json")
}

// keptFindings is what KeptPath holds: the folder linted, the configuration
// hash it was linted under, and each file's key and findings, by its path
// relative to the folder.
type keptFindings struct {
	Version int                  `json:"version"`
	Root    string               `json:"root"`
	Hash    string               `json:"hash"`
	Files   map[string]keptEntry `json:"files"`
}

// keptEntry is one file's findings and the stat they were linted at.
type keptEntry struct {
	Size     int64     `json:"size"`
	ModTime  int64     `json:"mtime"` // Unix nanoseconds
	Findings []Finding `json:"findings"`
}

// readKept is the findings kept at path for the folder root under hash. A
// file that is missing, cannot be read or parsed, is of another format, or
// was linted in another folder or under another hash, holds none.
func readKept(path, root, hash string) map[string]keptEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var k keptFindings
	if err := json.Unmarshal(data, &k); err != nil || k.Version != keptVersion || k.Root != root || k.Hash != hash {
		return nil
	}
	return k.Files
}

// writeKept writes k to path through a file beside it renamed over it, so
// that a reader finds the old findings or the new ones whole. It never makes
// the cache folder: a project without one keeps nothing (D-0402).
func writeKept(path string, k keptFindings) error {
	data, err := json.Marshal(k)
	if err != nil {
		return fmt.Errorf("encode the kept lint findings for %s: %w", path, err)
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".lint-findings-*")
	if err != nil {
		return fmt.Errorf("write the kept lint findings in %s: %w", dir, err)
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return fmt.Errorf("write the kept lint findings to %s: %w", name, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("close the kept lint findings in %s: %w", name, err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("keep the lint findings at %s: %w", path, err)
	}
	return nil
}
