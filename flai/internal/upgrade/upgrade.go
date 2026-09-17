// Package upgrade brings a project to a newer template version without
// overwriting project edits (ADR-0015).
package upgrade

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

// Classes of change for one path.
const (
	Add      = "add"      // template has it, project does not
	Merge    = "merge"    // both carry the baseline marker: template above, project below
	Replace  = "replace"  // project file unchanged since applied (hash matches the lock)
	Same     = "same"     // identical already
	Conflict = "conflict" // project changed it and it is not a marker file
)

// Change is the plan for one path.
type Change struct {
	Path    string      `json:"path"`
	Class   string      `json:"class"`
	content []byte      // new content to write (merged for Merge)
	mode    fs.FileMode // from the template
	NewHash string      `json:"new_hash"`
}

// Plan is the whole upgrade.
type Plan struct {
	From    string   `json:"from_version"`
	To      string   `json:"to_version"`
	Changes []Change `json:"changes"`
	NoLock  bool     `json:"no_lock"` // project had no lock; every difference is a conflict
}

// Count returns how many changes have the class.
func (p *Plan) Count(class string) int {
	n := 0
	for _, c := range p.Changes {
		if c.Class == class {
			n++
		}
	}
	return n
}

// Conflicts lists the conflict paths.
func (p *Plan) Conflicts() []string {
	var out []string
	for _, c := range p.Changes {
		if c.Class == Conflict {
			out = append(out, c.Path)
		}
	}
	return out
}

// Compute renders the template in memory and classifies every path.
func Compute(root string, m template.Manifest, src template.Source, opt template.Options, lk *lock.Lock, fromVersion string) (*Plan, error) {
	plan := &Plan{From: fromVersion, To: m.Version, NoLock: lk == nil}
	var rendered []Change
	opt.Collect = func(rel string, content []byte, mode fs.FileMode) {
		if owned(rel) {
			return
		}
		rendered = append(rendered, Change{Path: rel, content: content, mode: mode, NewHash: lock.Hash(content)})
	}
	if _, err := template.Render(m, src.Dir, root, opt); err != nil {
		return nil, err
	}
	for _, c := range rendered {
		existing, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(c.Path)))
		switch {
		case os.IsNotExist(err):
			c.Class = Add
		case err != nil:
			return nil, err
		case bytes.Equal(existing, c.content):
			c.Class = Same
		case hasMarker(existing) && hasMarker(c.content):
			c.Class = Merge
			c.content = merge(c.content, existing)
			if bytes.Equal(existing, c.content) {
				c.Class = Same
			}
		case lk != nil && lk.Files[c.Path] != "" && lk.Files[c.Path] == lock.Hash(existing):
			c.Class = Replace
		default:
			c.Class = Conflict
		}
		plan.Changes = append(plan.Changes, c)
	}
	sort.Slice(plan.Changes, func(i, j int) bool { return plan.Changes[i].Path < plan.Changes[j].Path })
	return plan, nil
}

func hasMarker(b []byte) bool { return bytes.Contains(b, []byte(conventions.Marker)) }

// owned lists rendered paths the project owns after creation; upgrade never
// touches them (flai updates the manifest's template fields itself).
func owned(rel string) bool { return rel == "system-flow.yaml" }

// merge takes the template's text above the marker and the project's text
// from the marker on.
func merge(tpl, project []byte) []byte {
	i := bytes.Index(tpl, []byte(conventions.Marker))
	j := bytes.Index(project, []byte(conventions.Marker))
	return append(append([]byte{}, tpl[:i]...), project[j:]...)
}

// Policy decides conflicts: a map from path to "keep" or "replace".
type Policy map[string]string

// Apply writes the plan. Conflicts are written only when the policy says
// replace; kept conflicts are left untouched. It returns the paths written.
func Apply(root string, plan *Plan, policy Policy) ([]string, error) {
	var written []string
	for _, c := range plan.Changes {
		switch c.Class {
		case Same:
			continue
		case Conflict:
			if policy[c.Path] != "replace" {
				continue
			}
		}
		target := filepath.Join(root, filepath.FromSlash(c.Path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return written, err
		}
		mode := c.mode
		if mode == 0 {
			mode = 0o644
		}
		if err := os.WriteFile(target, c.content, mode); err != nil {
			return written, err
		}
		written = append(written, c.Path)
	}
	return written, nil
}

// NewLock builds the lock after an upgrade: every rendered path with the
// template's hash, so a kept conflict stays a divergence next time.
func NewLock(plan *Plan, src template.Source, version string, now time.Time) *lock.Lock {
	l := &lock.Lock{Template: lock.Template{Repo: src.Repo, Ref: src.Ref, Version: version, Applied: now.UTC().Format("2006-01-02T15:04:05Z")}, Files: map[string]string{}}
	for _, c := range plan.Changes {
		l.Files[c.Path] = c.NewHash
	}
	return l
}

// Relock hashes the current project files for every template path without
// changing them, so a hand-assembled project can start upgrading.
func Relock(root string, m template.Manifest, src template.Source, opt template.Options, now time.Time) (*lock.Lock, error) {
	l := &lock.Lock{Template: lock.Template{Repo: src.Repo, Ref: src.Ref, Version: m.Version, Applied: now.UTC().Format("2006-01-02T15:04:05Z")}, Files: map[string]string{}}
	opt.Collect = func(rel string, content []byte, mode fs.FileMode) {
		if owned(rel) {
			return
		}
		if h, err := lock.HashFile(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
			l.Files[rel] = h
		}
	}
	if _, err := template.Render(m, src.Dir, root, opt); err != nil {
		return nil, err
	}
	return l, nil
}

// Describe renders a one-line summary.
func Describe(plan *Plan) string {
	return fmt.Sprintf("%s -> %s: %d add, %d merge, %d replace, %d conflict, %d unchanged", plan.From, plan.To,
		plan.Count(Add), plan.Count(Merge), plan.Count(Replace), plan.Count(Conflict), plan.Count(Same))
}
