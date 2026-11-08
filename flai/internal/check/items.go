package check

import (
	"os"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Items runs the work-item rules and the markdown lint over items alone, for
// a write that changes only values flai owns, such as a forecast (S-0599).
//
// all is the project's item list, the archive included, as Run lists it: the
// rules read an item's parent and children from it, and Items reads no item
// file to build it. Each of items gets the per-item rules Run gives it
// (item.*, story.*), and its file the markdown lint with the configuration
// Run's lint resolves. Items walks no folder, starts no node, and scans no
// digest, so it costs a few milliseconds where Run costs seconds.
func Items(repo *workitem.Repo, all, items []*workitem.Item, now time.Time) *Result {
	c := &checker{repo: repo, now: now.UTC(), items: all, byID: map[string]*workitem.Item{}, res: &Result{Items: len(items)}}
	for _, it := range all {
		if _, dup := c.byID[it.ID]; !dup {
			c.byID[it.ID] = it
		}
	}
	for _, it := range items {
		c.oneItem(it)
	}
	c.lintFiles(items)
	sortFindings(c.res.Findings)
	return c.res
}

// lintFiles lints the files of items, and no other, as markdown lints the
// files under the wip folder: with the main checkout's markdownlint
// configuration, or the worktree's when the main checkout's root holds none.
func (c *checker) lintFiles(items []*workitem.Item) {
	root := c.repo.MainRoot
	if root == "" {
		root = c.repo.Root
	}
	cfg, err := mdlint.Load(root)
	if err == nil && cfg == nil && root != c.repo.Root {
		root = c.repo.Root
		cfg, err = mdlint.Load(root)
	}
	if err != nil {
		c.add(Warning, "markdown.config", root, 1, "the markdownlint configuration cannot be read, so wip is not linted: %v", err)
		return
	}
	if cfg == nil {
		return
	}
	for _, it := range items {
		data, err := os.ReadFile(it.Path)
		if err != nil {
			continue // an unreadable file is skipped, as markdown skips it
		}
		for _, f := range cfg.Lint(string(data)) {
			c.add(Warning, "markdown."+f.Rule, it.Path, f.Line, "%s", f)
		}
	}
}
