package mdlint

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// maxWorkers caps the folder lint's pool: S-0600 measured little gain past 16
// workers on a loaded host.
const maxWorkers = 16

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
	if c == nil {
		return nil
	}
	var paths []string
	_ = filepath.WalkDir(dir, func(path string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil //nolint:nilerr // an unreadable entry is skipped, as the doc comment says
		}
		paths = append(paths, path)
		return nil
	})
	found := make([][]Finding, len(paths))
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), maxWorkers, len(paths)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				found[i] = c.lintFile(paths[i])
			}
		}()
	}
	for i := range paths {
		next <- i
	}
	close(next)
	wg.Wait()
	var out []FileFindings
	for i, f := range found {
		if len(f) > 0 {
			out = append(out, FileFindings{Path: paths[i], Findings: f})
		}
	}
	return out
}

// lintFile reads and lints one file of a folder lint: nothing when it cannot
// be read. Workers call it at once, so it only reads the configuration.
func (c *Config) lintFile(path string) []Finding {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return c.Lint(string(data))
}
