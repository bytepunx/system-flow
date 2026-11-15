// Package upgrade brings a project to a newer template version without
// overwriting project edits (ADR-0015).
package upgrade

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Classes of change for one path.
const (
	Add      = "add"      // template has it, project does not
	Merge    = "merge"    // both carry the baseline marker: template above, project below
	Replace  = "replace"  // project file unchanged since applied (hash matches the lock)
	Same     = "same"     // identical already, or the project's own workflow file while the template's baseline is unchanged
	Conflict = "conflict" // project changed it and it is not a marker file
	Keep     = "keep"     // the project's own workflow file, never written, when the template's baseline changed since the lock (ADR-0307 §4)
)

// Change is the plan for one path.
type Change struct {
	Path    string      `json:"path"`
	Class   string      `json:"class"`
	content []byte      // new content to write (merged for Merge)
	mode    fs.FileMode // from the template
	NewHash string      `json:"new_hash"`
	topics  []string    // the template's topics on a marker file, for the lock
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
		rendered = append(rendered, Change{Path: rel, content: content, mode: mode, NewHash: lock.Hash(content), topics: Topics(content)})
	}
	if _, err := template.Render(m, src.Dir, root, opt); err != nil {
		return nil, err
	}
	keptPath := kept(m, opt)
	for _, c := range rendered {
		existing, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(c.Path)))
		switch {
		case os.IsNotExist(err):
			c.Class = Add
		case err != nil:
			return nil, err
		case c.Path == keptPath:
			c.Class = Same
			if was := recordedHash(lk, c.Path); was != "" && was != c.NewHash && !bytes.Equal(existing, c.content) {
				c.Class = Keep
			}
		case bytes.Equal(existing, c.content):
			c.Class = Same
		case hasMarker(existing) && hasMarker(c.content):
			c.Class = Merge
			c.content = merge(c.content, existing, recorded(lk, c.Path))
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

// kept is the one rendered path upgrade adds when the project lacks it and
// never writes once the project has it, whatever its content and the lock:
// the board's workflow file under the design folder the layout chooses
// (ADR-0307 §4). Once the project has it, it is Same while the template's
// baseline is the one the lock recorded, the lock has none, or the project
// already holds the template's; otherwise Keep, so upgrade says the
// template's baseline changed, and the new lock records the new hash, so it
// says so once per change.
func kept(m template.Manifest, opt template.Options) string {
	design := m.Layout["design"]
	if v := opt.Layout["design"]; v != "" {
		design = v
	}
	if design == "" {
		design = "design"
	}
	return workflow.Rel(design)
}

// recordedHash is the hash the lock recorded for path: empty with no lock
// or no entry.
func recordedHash(lk *lock.Lock, path string) string {
	if lk == nil {
		return ""
	}
	return lk.Files[path]
}

func hasMarker(b []byte) bool { return bytes.Contains(b, []byte(conventions.Marker)) }

// owned lists rendered paths the project owns after creation; upgrade never
// touches them (flai updates the manifest's template fields itself).
func owned(rel string) bool { return rel == "system-flow.yaml" }

// merge takes the template's text above the marker and the project's text
// from the marker on. The project's topics stay when the project set them:
// when they differ from the topics the template last gave the file (was),
// or, with nothing recorded, when the project has any.
func merge(tpl, project []byte, was []string) []byte {
	i := bytes.Index(tpl, []byte(conventions.Marker))
	j := bytes.Index(project, []byte(conventions.Marker))
	out := append(append([]byte{}, tpl[:i]...), project[j:]...)
	own := Topics(project)
	if own == nil || (was != nil && slices.Equal(own, was)) {
		return out
	}
	fm, _, err := workitem.SplitFrontMatter(string(out))
	if err != nil {
		return out
	}
	return []byte("---\n" + topics.SetInFrontMatter(fm, own) + strings.TrimPrefix(string(out), "---\n"+fm))
}

// Topics reads the topics on a marker file's front matter: nil for a file
// without the marker, without front matter, or without topics.
func Topics(content []byte) []string {
	if !hasMarker(content) {
		return nil
	}
	fm, _, err := workitem.SplitFrontMatter(string(content))
	if err != nil {
		return nil
	}
	t, err := topics.FromFrontMatter(fm)
	if err != nil || len(t) == 0 {
		return nil
	}
	return t
}

// TopicsOf reads the topics on each marker file among paths under root, for
// the lock of a project just rendered.
func TopicsOf(root string, paths []string) map[string][]string {
	out := map[string][]string{}
	for _, rel := range paths {
		if b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
			if t := Topics(b); t != nil {
				out[rel] = t
			}
		}
	}
	return out
}

// recorded is the topics the lock says the template last gave path.
func recorded(lk *lock.Lock, path string) []string {
	if lk == nil {
		return nil
	}
	return lk.Topics[path]
}

// Policy decides conflicts: a map from path to "keep" or "replace".
type Policy map[string]string

// Apply writes the plan. Conflicts are written only when the policy says
// replace; kept conflicts are left untouched. It returns the paths written.
func Apply(root string, plan *Plan, policy Policy) ([]string, error) {
	var written []string
	for _, c := range plan.Changes {
		switch c.Class {
		case Same, Keep:
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
	l := &lock.Lock{Template: lock.Template{Repo: src.Repo, Ref: src.Ref, Version: version, Applied: now.UTC().Format("2006-01-02T15:04:05Z")}, Files: map[string]string{}, Topics: map[string][]string{}}
	for _, c := range plan.Changes {
		l.Files[c.Path] = c.NewHash
		if c.topics != nil {
			l.Topics[c.Path] = c.topics
		}
	}
	return l
}

// Relock hashes the current project files for every template path without
// changing them, so a hand-assembled project can start upgrading. The
// topics recorded are the template's, so a project's own stay its own. For
// the workflow file upgrade keeps it records the template's hash when the
// project has the file, since upgrade compares the template's baseline with
// it, not the project's file (ADR-0307 §4).
func Relock(root string, m template.Manifest, src template.Source, opt template.Options, now time.Time) (*lock.Lock, error) {
	l := &lock.Lock{Template: lock.Template{Repo: src.Repo, Ref: src.Ref, Version: m.Version, Applied: now.UTC().Format("2006-01-02T15:04:05Z")}, Files: map[string]string{}, Topics: map[string][]string{}}
	keptPath := kept(m, opt)
	opt.Collect = func(rel string, content []byte, mode fs.FileMode) {
		if owned(rel) {
			return
		}
		if h, err := lock.HashFile(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
			if rel == keptPath {
				h = lock.Hash(content)
			}
			l.Files[rel] = h
		}
		if t := Topics(content); t != nil {
			l.Topics[rel] = t
		}
	}
	if _, err := template.Render(m, src.Dir, root, opt); err != nil {
		return nil, err
	}
	return l, nil
}

// Describe renders a one-line summary. When the template's baseline
// workflow changed, it says the project's file is kept and that the
// operator carries the change over in a story (ADR-0307 §4).
func Describe(plan *Plan) string {
	s := fmt.Sprintf("%s -> %s: %d add, %d merge, %d replace, %d conflict, %d unchanged", plan.From, plan.To,
		plan.Count(Add), plan.Count(Merge), plan.Count(Replace), plan.Count(Conflict), plan.Count(Same))
	for _, c := range plan.Changes {
		if c.Class == Keep {
			s += fmt.Sprintf("; the template's baseline %s changed: the project's file is kept, carry the change over to it in a story", c.Path)
		}
	}
	return s
}
