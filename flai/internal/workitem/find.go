package workitem

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// findDepth is how far below a folder projects are looked for: ~/git/p,
// ~/git/org/p, and one more, as flai serve looks for repositories to import.
const findDepth = 3

var skipFinding = map[string]bool{"node_modules": true, "vendor": true, "dist": true, "build": true, "target": true, "__pycache__": true}

// FindProjects lists the system-flow projects in root and below it (S-0101):
// folders with a system-flow.yaml of their own, not looked into further.
// Hidden folders (story worktrees live under .flai-cache) and build output
// never are.
func FindProjects(root string) []string {
	var out []string
	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		if _, err := os.Stat(filepath.Join(dir, manifest.File)); err == nil {
			out = append(out, dir)
			return
		}
		if depth >= findDepth {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || skipFinding[e.Name()] {
				continue
			}
			walk(filepath.Join(dir, e.Name()), depth+1)
		}
	}
	walk(root, 0)
	sort.Strings(out)
	return out
}
