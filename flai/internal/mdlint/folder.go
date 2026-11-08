package mdlint

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// maxWorkers caps the folder lint's pool: S-0600 measured little gain past 16
// workers on a loaded host.
const maxWorkers = 16

// readFile reads a file the folder lint lints; a test replaces it to count
// the files read.
var readFile = os.ReadFile

// FileFindings is one file's findings from a folder lint.
type FileFindings struct {
	Path     string
	Findings []Finding
}

// LintFolder lints every markdown file under dir with the configuration, on a
// pool of at most GOMAXPROCS workers, capped at 16, and returns the files that
// have findings in the order filepath.WalkDir visits them, each file's
// findings in Lint's order. An entry the walk cannot read, and a file that
// cannot be read, such as a link that dangles, are skipped, as the other
// rules of flai check skip what they cannot read. A nil configuration finds
// nothing.
func (c *Config) LintFolder(dir string) []FileFindings {
	return c.LintFolderKept(dir, "")
}

// LintFolderKept is LintFolder that keeps each file's findings between
// processes in the cache folder cacheDir (S-0600), at KeptPath. A file whose
// size and modification time are those it was linted at, under the same
// configuration Hash, takes its kept findings unread; the rest are linted.
// The kept file is written again when anything changed, and only into a
// cacheDir that exists. A kept file that is missing, damaged, or another
// folder's or configuration's is a miss, and a failed write fails nothing:
// either costs the next run a lint, never a wrong finding. An empty cacheDir
// keeps nothing.
func (c *Config) LintFolderKept(dir, cacheDir string) []FileFindings {
	if c == nil {
		return nil
	}
	var paths []string
	var infos []fs.FileInfo
	_ = filepath.WalkDir(dir, func(path string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil //nolint:nilerr // an unreadable entry is skipped, as the doc comment says
		}
		paths = append(paths, path)
		infos = append(infos, nil)
		return nil
	})
	keep := cacheDir != "" && c.hash != "" && isDir(cacheDir)
	var kept map[string]keptEntry
	var keptPath string
	if keep {
		keptPath = KeptPath(cacheDir, dir)
		kept = readKept(keptPath, filepath.Clean(dir), c.hash)
	}
	found := make([][]Finding, len(paths))
	linted := make([]bool, len(paths))
	var todo []int
	for i, path := range paths {
		if keep {
			// A file whose stat fails is linted, which reads it, so it is
			// skipped as before.
			if info, err := os.Stat(path); err == nil {
				infos[i] = info
				if k, ok := kept[rel(dir, path)]; ok && k.Size == info.Size() && k.ModTime == info.ModTime().UnixNano() {
					found[i] = k.Findings
					continue
				}
			}
		}
		todo = append(todo, i)
	}
	lintedAt := time.Now()
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), maxWorkers, len(todo)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				found[i], linted[i] = c.lintFile(paths[i])
			}
		}()
	}
	for _, i := range todo {
		next <- i
	}
	close(next)
	wg.Wait()
	if keep {
		c.writeFolderKept(keptPath, dir, paths, infos, found, linted, lintedAt, kept)
	}
	var out []FileFindings
	for i, f := range found {
		if len(f) > 0 {
			out = append(out, FileFindings{Path: paths[i], Findings: f})
		}
	}
	return out
}

// writeFolderKept writes the folder's kept findings again when they differ
// from what was read: a file linted whose modification came racyWindow
// before the lint is kept, a file still kept is carried over, and a file no
// longer in the folder, or no longer kept, is dropped.
func (c *Config) writeFolderKept(path, dir string, paths []string, infos []fs.FileInfo, found [][]Finding, linted []bool, lintedAt time.Time, kept map[string]keptEntry) {
	files := make(map[string]keptEntry, len(paths))
	for i, p := range paths {
		info := infos[i]
		if info == nil {
			continue
		}
		name := rel(dir, p)
		if !linted[i] {
			if k, ok := kept[name]; ok && k.Size == info.Size() && k.ModTime == info.ModTime().UnixNano() {
				files[name] = k
			}
			continue
		}
		if info.ModTime().Before(lintedAt.Add(-racyWindow)) {
			files[name] = keptEntry{Size: info.Size(), ModTime: info.ModTime().UnixNano(), Findings: found[i]}
		}
	}
	if kept != nil && sameKept(kept, files) {
		return
	}
	_ = writeKept(path, keptFindings{Version: keptVersion, Root: filepath.Clean(dir), Hash: c.hash, Files: files})
}

// sameKept reports whether two sets of kept entries hold the same files at
// the same keys: findings follow from the key, so they are not compared.
func sameKept(a, b map[string]keptEntry) bool {
	return maps.EqualFunc(a, b, func(x, y keptEntry) bool { return x.Size == y.Size && x.ModTime == y.ModTime })
}

// lintFile reads and lints one file of a folder lint: nothing, and not
// linted, when it cannot be read. Workers call it at once, so it only reads
// the configuration.
func (c *Config) lintFile(path string) ([]Finding, bool) {
	data, err := readFile(path)
	if err != nil {
		return nil, false
	}
	return c.Lint(string(data)), true
}

// rel names path from dir, as the kept file keys it.
func rel(dir, path string) string {
	r, err := filepath.Rel(dir, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(r)
}

// isDir reports whether path is a folder that exists.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
