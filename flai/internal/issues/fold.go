package issues

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folded is an open issue a fold moved into the main branch's open issue of
// the same title (ADR-0126).
type Folded struct {
	From string `json:"from"` // the branch's issue, I-nnnn
	Into string `json:"into"` // the main branch's issue it went into
}

// Fold folds each open issue in r, a story's checkout, whose file the main
// branch does not have and whose title an open issue whose file the main
// branch has holds too, into that issue, the lowest numbered when several do:
// it adds the branch's issue's instances, count, and cost to it, as a merge
// adds a side's, sets its updated to now, saves it, and deletes the branch's
// file. The rest of the branch's issue, its description, impact, and
// remediation, is not carried over. onMain is the base names of the issue
// files the main branch has. It does not regenerate summary.md. It returns
// what it folded, in ID order, and the paths it wrote or deleted, relative to
// r.Root, slash-separated as git names them.
func Fold(r *workitem.Repo, onMain []string, now time.Time) ([]Folded, []string, error) {
	list, err := List(r)
	if err != nil {
		return nil, nil, fmt.Errorf("read the issues in %s to fold the branch's into the main branch's: %w", Dir(r), err)
	}
	mainHas := map[string]bool{}
	for _, name := range onMain {
		mainHas[name] = true
	}
	var targets []*Issue
	for _, is := range list {
		if mainHas[filepath.Base(is.Path)] {
			targets = append(targets, is)
		}
	}
	ts := now.UTC().Format(workitem.TimeFormat)
	folded := []Folded{}
	paths := []string{}
	seen := map[string]bool{}
	note := func(rel string) {
		if !seen[rel] {
			seen[rel] = true
			paths = append(paths, rel)
		}
	}
	for _, is := range list {
		if mainHas[filepath.Base(is.Path)] || is.Status != "open" {
			continue
		}
		into := FindOpenByTitle(targets, is.Title)
		if into == nil {
			continue
		}
		gone, err := relative(r, is)
		if err != nil {
			return nil, nil, err
		}
		if err := addOccurrences(into, is); err != nil {
			return nil, nil, fmt.Errorf("fold %s into %s, the main branch's issue of its title: %w", is.ID, into.ID, err)
		}
		into.Updated = ts
		if err := into.Save(); err != nil {
			return nil, nil, fmt.Errorf("write %s, folding %s into it: %w; check that %s is writable and sync again", into.ID, is.ID, err, into.Path)
		}
		wrote, err := into.Changed(r)
		if err != nil {
			return nil, nil, err
		}
		if err := os.Remove(is.Path); err != nil {
			return nil, nil, fmt.Errorf("delete %s, folded into %s: %w; delete %s by hand", is.ID, into.ID, err, is.Path)
		}
		for _, p := range wrote {
			note(p)
		}
		note(gone)
		folded = append(folded, Folded{From: is.ID, Into: into.ID})
	}
	return folded, paths, nil
}

// relative is the issue's file relative to r's root, slash-separated as git
// names it.
func relative(r *workitem.Repo, is *Issue) (string, error) {
	rel, err := filepath.Rel(r.Root, is.Path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%s is at %s, outside the checkout %s; fold the issues of the checkout whose branch they are on", is.ID, is.Path, r.Root)
	}
	return filepath.ToSlash(rel), nil
}
