// Package importer analyses an existing repository and plans how to bring
// it under the system-flow standard. See design/system/flai-cli.md, "Import
// flow".
package importer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Project is a code sub-project found by its build file.
type Project struct {
	Name string `json:"name"`
	Path string `json:"path"` // relative to the root
	Kind string `json:"kind"` // go, sveltekit, node, python, rust
}

// Analysis is what Scan found.
type Analysis struct {
	Root       string            `json:"root"`
	Git        bool              `json:"git"`
	Manifest   bool              `json:"manifest"`   // system-flow.yaml already present
	Existing   map[string]string `json:"existing"`   // layout key -> existing folder of that name
	Candidates map[string]string `json:"candidates"` // existing folder -> suggested destination (relative)
	Markdown   []string          `json:"markdown"`   // loose markdown outside code and candidates
	Projects   []Project         `json:"projects"`
}

var (
	skipDirs   = map[string]bool{".git": true, "node_modules": true, "vendor": true, "dist": true, "build": true, ".svelte-kit": true, "target": true, "__pycache__": true, ".venv": true}
	buildFiles = map[string]string{"go.mod": "go", "package.json": "node", "pyproject.toml": "python", "Cargo.toml": "rust"}
	// candidateDirs maps folder names people already use to where their
	// markdown belongs in the standard (relative to the root, using the
	// chosen layout names at plan time).
	candidateDirs = map[string]string{"adr": "design/adrs", "adrs": "design/adrs", "architecture": "design/system", "documentation": "docs/users", "doc": "docs/users"}
	layoutDirs    = []string{"design", "docs", "wip"}
	rootDocs      = map[string]bool{"README.md": true, "CLAUDE.md": true, "CHANGELOG.md": true, "CONTRIBUTING.md": true, "LICENSE.md": true, "CODE_OF_CONDUCT.md": true, "SECURITY.md": true}
)

// Scan inspects root without changing anything.
func Scan(root string) (*Analysis, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	a := &Analysis{Root: root, Existing: map[string]string{}, Candidates: map[string]string{}}
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		a.Git = true
	}
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err == nil {
		a.Manifest = true
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		for _, l := range layoutDirs {
			if name == l {
				a.Existing[l] = name
			}
		}
		if dest, ok := candidateDirs[strings.ToLower(name)]; ok {
			a.Candidates[name] = dest
		}
	}
	// sub-projects: any directory (depth <= 3) with a build file, not the root
	projects := map[string]Project{}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil //nolint:nilerr // unreadable entries are skipped
		}
		rel, _ := filepath.Rel(root, path)
		if d.IsDir() {
			if skipDirs[d.Name()] || (rel != "." && strings.Count(rel, string(filepath.Separator)) >= 3) {
				return filepath.SkipDir
			}
			if _, isProject := projects[filepath.Dir(rel)]; isProject && rel != "." {
				return filepath.SkipDir // do not descend into a detected project
			}
			return nil
		}
		if kind, ok := buildFiles[d.Name()]; ok && filepath.Dir(rel) != "." {
			dir := filepath.Dir(rel)
			if kind == "node" {
				if m, _ := filepath.Glob(filepath.Join(filepath.Dir(path), "svelte.config.*")); len(m) > 0 {
					kind = "sveltekit"
				}
			}
			if p, exists := projects[dir]; !exists || p.Kind == "node" {
				projects[dir] = Project{Name: filepath.Base(dir), Path: filepath.ToSlash(dir), Kind: kind}
			}
		}
		return nil
	})
	for _, p := range projects {
		a.Projects = append(a.Projects, p)
	}
	sort.Slice(a.Projects, func(i, j int) bool { return a.Projects[i].Path < a.Projects[j].Path })
	// loose markdown: outside projects, layout folders, and candidates
	inProject := func(rel string) bool {
		for _, p := range a.Projects {
			if rel == p.Path || strings.HasPrefix(rel, p.Path+"/") {
				return true
			}
		}
		return false
	}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil //nolint:nilerr // unreadable entries are skipped
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel == "." {
				return nil
			}
			if skipDirs[d.Name()] || inProject(rel) {
				return filepath.SkipDir
			}
			first := strings.SplitN(rel, "/", 2)[0]
			if _, ok := a.Existing[first]; ok {
				return filepath.SkipDir
			}
			if _, ok := a.Candidates[first]; ok {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".md") && (filepath.Dir(rel) != "." || !rootDocs[d.Name()]) {
			a.Markdown = append(a.Markdown, rel)
		}
		return nil
	})
	sort.Strings(a.Markdown)
	return a, nil
}

// Layout is the chosen folder names.
type Layout map[string]string

// Plan is what Apply would do.
type Plan struct {
	Create      []string  `json:"create"`       // folders to create (relative)
	FolderMoves []Move    `json:"folder_moves"` // whole candidate folders
	Markdown    []string  `json:"markdown"`     // files the operator decides about
	Projects    []Project `json:"projects"`
	Layout      Layout    `json:"layout"`
}

// Move is a source to destination relocation.
type Move struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// BuildPlan derives the proposal from an analysis and layout names.
func BuildPlan(a *Analysis, layout Layout) *Plan {
	p := &Plan{Layout: layout, Projects: a.Projects, Markdown: a.Markdown}
	subs := map[string][]string{"design": {"adrs", "system", "tech", "conventions"}, "docs": {"users", "operators", "contributors"}, "wip": {"kanban", "kanban/epics", "kanban/stories", "kanban/tasks", "agents", "archive"}}
	for _, key := range layoutDirs {
		dir := layout[key]
		if _, err := os.Stat(filepath.Join(a.Root, dir)); err != nil {
			p.Create = append(p.Create, dir)
		}
		for _, s := range subs[key] {
			if _, err := os.Stat(filepath.Join(a.Root, dir, s)); err != nil {
				p.Create = append(p.Create, filepath.ToSlash(filepath.Join(dir, s)))
			}
		}
	}
	if _, err := os.Stat(filepath.Join(a.Root, "scripts")); err != nil {
		p.Create = append(p.Create, "scripts")
	}
	for from, dest := range a.Candidates {
		parts := strings.SplitN(dest, "/", 2)
		to := filepath.ToSlash(filepath.Join(layout[parts[0]], parts[1]))
		if from != to {
			p.FolderMoves = append(p.FolderMoves, Move{From: from, To: to})
		}
	}
	sort.Slice(p.FolderMoves, func(i, j int) bool { return p.FolderMoves[i].From < p.FolderMoves[j].From })
	return p
}

// MoveFile relocates one file, with git mv when the file is tracked.
func MoveFile(r execx.Runner, root, from, to string, git bool) error {
	src, dst := filepath.Join(root, from), filepath.Join(root, to)
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if git {
		if _, err := r.Run(root, "git", "ls-files", "--error-unmatch", from); err == nil {
			_, err := r.Run(root, "git", "mv", from, to)
			return err
		}
	}
	return os.Rename(src, dst)
}

// MoveFolderContents moves every file under from into to, preserving
// relative paths, then removes empty source folders.
func MoveFolderContents(r execx.Runner, root, from, to string, git bool) ([]Move, error) {
	var moved []Move
	src := filepath.Join(root, from)
	err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		m := Move{From: filepath.ToSlash(filepath.Join(from, rel)), To: filepath.ToSlash(filepath.Join(to, rel))}
		if _, err := os.Stat(filepath.Join(root, m.To)); err == nil {
			return nil // never overwrite; left in place and reported by the caller
		}
		if err := MoveFile(r, root, m.From, m.To, git); err != nil {
			return err
		}
		moved = append(moved, m)
		return nil
	})
	if err != nil {
		return moved, err
	}
	_ = removeEmptyDirs(src)
	return moved, nil
}

func removeEmptyDirs(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			_ = removeEmptyDirs(filepath.Join(dir, e.Name()))
		}
	}
	entries, _ = os.ReadDir(dir)
	if len(entries) == 0 {
		return os.Remove(dir)
	}
	return nil
}
