// Package itemnew creates a work item whose body its author has already
// written, as one step that happens or does not (S-0059): the item is made
// from the project's template, flai check runs with it in place, anything it
// introduces refuses it and leaves nothing behind, and what is kept is
// committed on its own. It is what the dashboard's "new" form calls, and the
// same from a script: flai story new --body-stdin.
package itemnew

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Options say what to create and what to do with it.
type Options struct {
	New workitem.NewOptions
	// Autocommit commits the new item and its parent on their own, unless
	// the project sets dashboard.autocommit: false.
	Autocommit bool
	Trailers   []string
}

// Result is the created item and what became of it.
type Result struct {
	Item        *workitem.Item  `json:"item"`
	Path        string          `json:"path"`
	Warnings    []check.Finding `json:"warnings"`
	Committed   bool            `json:"committed"`
	Commit      string          `json:"commit,omitempty"`
	CommitError string          `json:"commit_error,omitempty"`
}

// Create makes the item, checks the repository with it in place, and keeps
// it only when the check reports nothing new. A refusal is a
// *docedit.RefusedError carrying the findings, as for a refused save.
func Create(repo *workitem.Repo, r execx.Runner, opt Options) (*Result, error) {
	if opt.New.Now.IsZero() {
		opt.New.Now = time.Now()
	}
	before, err := check.Run(repo, opt.New.Now)
	if err != nil {
		return nil, err
	}
	// The parent's file is rewritten to link the child; keep what it was.
	var parentPath string
	var parentWas []byte
	if opt.New.Parent != "" {
		if p, err := repo.Get(opt.New.Parent); err == nil {
			parentPath = p.Path
			parentWas, _ = os.ReadFile(p.Path)
		}
	}
	it, err := repo.Create(opt.New)
	if err != nil {
		return nil, err
	}
	undo := func() error {
		if err := os.Remove(it.Path); err != nil && !os.IsNotExist(err) {
			return err
		}
		if parentPath != "" && parentWas != nil {
			return atomicfile.WriteFile(parentPath, parentWas, 0o644)
		}
		return nil
	}
	after, err := check.Run(repo, opt.New.Now)
	if err != nil {
		_ = undo()
		return nil, err
	}
	known := map[string]bool{}
	for _, f := range before.Findings {
		known[f.Rule+"\x00"+f.Path+"\x00"+f.Message] = true
	}
	rel := relTo(repo.Root, it.Path)
	res := &Result{Item: it, Path: rel, Warnings: []check.Finding{}}
	var blocking []check.Finding
	for _, f := range after.Findings {
		if !known[f.Rule+"\x00"+f.Path+"\x00"+f.Message] {
			blocking = append(blocking, f)
		}
	}
	if len(blocking) > 0 {
		if err := undo(); err != nil {
			return nil, fmt.Errorf("the check refused the new %s and removing it failed: %w", it.Type, err)
		}
		return nil, &docedit.RefusedError{Path: rel, Reason: fmt.Sprintf("flai check has %d finding(s) with this %s; nothing was created", len(blocking), it.Type), Findings: blocking}
	}
	if !opt.Autocommit || !repo.Manifest.Autocommit() {
		return res, nil
	}
	root := repo.MainRoot
	if root == "" {
		root = repo.Root
	}
	paths := []string{relTo(root, it.Path)}
	if parentPath != "" {
		paths = append(paths, relTo(root, parentPath))
	}
	msg := fmt.Sprintf("chore: [%s] create %s: %s", it.ID, it.Type, it.Title)
	if len(opt.Trailers) > 0 {
		msg += "\n\n" + strings.Join(opt.Trailers, "\n")
	}
	sha, inRepo, err := commit(r, root, paths, msg)
	switch {
	case !inRepo:
	case err != nil:
		res.CommitError = err.Error()
	default:
		res.Committed, res.Commit = true, sha
	}
	return res, nil
}

func relTo(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return path
}

// commit records the paths, and only them, on the current branch. A failed
// commit is not a failed creation: the files are written and valid.
func commit(r execx.Runner, root string, paths []string, msg string) (sha string, inRepo bool, err error) {
	if r == nil {
		return "", false, nil
	}
	if _, err := r.Run(root, "git", "rev-parse", "--is-inside-work-tree"); err != nil {
		return "", false, nil //nolint:nilerr // no repository means nothing to commit to, by design
	}
	if _, err := r.Run(root, "git", append([]string{"add", "--"}, paths...)...); err != nil {
		return "", true, err
	}
	if _, err := r.Run(root, "git", append([]string{"commit", "--quiet", "-m", msg, "--"}, paths...)...); err != nil {
		return "", true, err
	}
	sha, err = r.Run(root, "git", "rev-parse", "--short", "HEAD")
	return strings.TrimSpace(sha), true, err
}
