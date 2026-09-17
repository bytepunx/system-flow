package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// IDMigration is the plan for widening work item IDs to IDWidth digits:
// which IDs change, which files move, and which files have references
// rewritten. Paths are repository-relative with forward slashes.
type IDMigration struct {
	Map      map[string]string `json:"map"`      // old ID -> new ID
	Renames  []Rename          `json:"renames"`  // item and narrative files to move
	Rewrites []string          `json:"rewrites"` // files whose contents change
}

// Rename is one file move in a migration plan.
type Rename struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// shortID matches an ID narrower than IDWidth as a whole word, so E-1,
// S-032 and T-100 are found but T-0100, ADR-0003 and I-011 are not.
// fileID is the ID prefix of an item file name.
var fileID = regexp.MustCompile(`^[EST]-\d+`)

var shortID = regexp.MustCompile(fmt.Sprintf(`\b([EST])-(\d{1,%d})\b`, IDWidth-1))

// widen rewrites every short ID in text to its canonical form.
func widen(text string) string {
	return shortID.ReplaceAllStringFunc(text, CanonicalID)
}

// PlanIDMigration lists the moves and rewrites needed to bring every item in
// kanban and archive, every narrative, and every reference in the layout
// folders, the repository's root markdown and yaml files, and each project's
// root markdown files up to IDWidth digits. Nothing is changed.
func (r *Repo) PlanIDMigration() (*IDMigration, error) {
	plan := &IDMigration{Map: map[string]string{}}
	items, err := r.List(true)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		// The ID in the file name is the source of truth for what moves, so a
		// run interrupted after rewriting front matter still completes.
		base := filepath.Base(it.Path)
		oldID := fileID.FindString(base)
		if oldID == "" {
			oldID = it.ID
		}
		newID := CanonicalID(it.ID)
		if it.ID != newID {
			plan.Map[it.ID] = newID
		}
		if oldID != newID {
			plan.Map[oldID] = newID
			plan.Renames = append(plan.Renames, Rename{
				From: r.rel(it.Path),
				To:   r.rel(filepath.Join(filepath.Dir(it.Path), newID+strings.TrimPrefix(base, oldID))),
			})
		}
		if it.Type != Story || oldID == newID {
			continue
		}
		for _, dir := range []string{r.AgentsDir(), filepath.Join(r.ArchiveDir(), "agents")} {
			old := filepath.Join(dir, oldID+".md")
			if _, err := os.Stat(old); err == nil {
				plan.Renames = append(plan.Renames, Rename{From: r.rel(old), To: r.rel(filepath.Join(dir, newID+".md"))})
			}
		}
	}
	files, err := r.migrationFiles()
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if shortID.Match(data) {
			plan.Rewrites = append(plan.Rewrites, r.rel(f))
		}
	}
	sort.Slice(plan.Renames, func(i, j int) bool { return plan.Renames[i].From < plan.Renames[j].From })
	sort.Strings(plan.Rewrites)
	return plan, nil
}

// ApplyIDMigration rewrites references first, then moves files with mv, which
// is `git mv` inside a git repository and os.Rename elsewhere.
func (r *Repo) ApplyIDMigration(plan *IDMigration, mv func(from, to string) error) error {
	for _, rel := range plan.Rewrites {
		p := filepath.Join(r.Root, filepath.FromSlash(rel))
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(widen(string(data))), 0o644); err != nil {
			return err
		}
	}
	for _, mvp := range plan.Renames {
		from := filepath.Join(r.Root, filepath.FromSlash(mvp.From))
		to := filepath.Join(r.Root, filepath.FromSlash(mvp.To))
		if err := mv(from, to); err != nil {
			return fmt.Errorf("move %s: %w", mvp.From, err)
		}
	}
	return nil
}

// migrationFiles returns the markdown and yaml files whose references are
// rewritten: everything under the layout folders, the root's own files, and
// each project's root files. Hidden folders, node_modules, bin, and testdata
// are skipped so fixtures and caches keep their historical IDs.
func (r *Repo) migrationFiles() ([]string, error) {
	var out []string
	keep := func(name string) bool {
		return strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")
	}
	skipDir := func(name string) bool {
		return strings.HasPrefix(name, ".") || name == "node_modules" || name == "bin" || name == "testdata" || name == "build"
	}
	for _, key := range []string{"design", "docs", "wip"} {
		dir := r.Manifest.Dir(r.Root, key)
		err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				if os.IsNotExist(err) && p == dir {
					return filepath.SkipDir
				}
				return err
			}
			if d.IsDir() {
				if p != dir && skipDir(d.Name()) {
					return filepath.SkipDir
				}
				return nil
			}
			if keep(d.Name()) {
				out = append(out, p)
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	roots := []string{r.Root}
	for _, p := range r.Manifest.Projects {
		roots = append(roots, filepath.Join(r.Root, p.Path))
	}
	for _, dir := range roots {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && keep(e.Name()) {
				out = append(out, filepath.Join(dir, e.Name()))
			}
		}
	}
	sort.Strings(out)
	return out, nil
}

func (r *Repo) rel(p string) string {
	rel, err := filepath.Rel(r.Root, p)
	if err != nil {
		return filepath.ToSlash(p)
	}
	return filepath.ToSlash(rel)
}
