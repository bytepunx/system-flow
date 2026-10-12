package storygit

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// The lint fix after a clean rebase (S-0367): a file the branch changes may
// carry what the project's fix commands rewrite, its own edit or what the main
// branch brought into it, so the sync runs them over every file the branch
// changes against the base and commits what they rewrote on the branch, and
// bad markdown does not reach main.

// FixFunc rewrites in place those of the root-relative paths of the checkout
// at dir that the project's fix commands select, and answers the paths it
// changed.
type FixFunc func(dir string, paths []string) ([]string, error)

// fix runs o.Fix in the story's worktree at path over the files the branch
// changes against res.Base, deleted ones left out: when it changes any, it
// stages them, commits them on the branch, and sets res.Fixed. A fix that
// fails, or files that cannot be listed, is logged and the sync stands; a fix
// whose files are written but not committed is an error saying what to
// commit.
func fix(o SyncOptions, path string, res *SyncResult) error {
	if o.Fix == nil {
		return nil
	}
	id := o.Story.ID
	changed, err := nameOnly(o.Runner, path, "diff", "--name-only", "--diff-filter=d", res.Base+"...HEAD")
	if err != nil {
		o.Log.Warn("the files the branch changes not listed, so not fixed; the sync stands", "component", "git", "story", id, "err", err)
		return nil
	}
	if len(changed) == 0 {
		return nil
	}
	fixed, err := o.Fix(path, changed)
	if err != nil {
		o.Log.Warn("the lint of the files the branch changes not fixed; the sync stands", "component", "git", "story", id, "err", err)
		return nil
	}
	if len(fixed) == 0 {
		return nil
	}
	subject := fixSubject(id)
	undone := func(step string, err error) error {
		return fmt.Errorf("%s the lint fix of %s on %s: %w; in %s, git add -A -- %s, commit it as %q, and run %s again",
			step, strings.Join(fixed, ", "), res.Branch, err, res.Worktree, strings.Join(fixed, " "), subject, again(o))
	}
	if _, err := o.Runner.Run(path, "git", append([]string{"add", "-A", "--"}, fixed...)...); err != nil {
		return undone("stage", err)
	}
	staged, err := nameOnly(o.Runner, path, "diff", "--cached", "--name-only", "--", ".")
	if err != nil {
		return undone("read what is staged of", err)
	}
	if len(staged) == 0 {
		return nil
	}
	if _, err := o.Runner.Run(path, "git", "commit", "-q", "-m", subject); err != nil {
		return undone("commit", err)
	}
	res.Fixed = staged
	o.Log.Info("the lint of the files the branch changes fixed after the rebase", "component", "git", "story", id, "paths", staged)
	return nil
}

// fixSubject is the subject of the commit of the lint fix on story's branch.
func fixSubject(story string) string {
	return fmt.Sprintf("chore: [%s] fix the lint of the files the branch changes, after the sync", story)
}

// nameOnly runs git with args in dir and answers the paths it prints, one per
// line.
func nameOnly(r execx.Runner, dir string, args ...string) ([]string, error) {
	out, err := r.Run(dir, "git", args...)
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for _, p := range strings.Split(out, "\n") {
		if p = strings.TrimSpace(p); p != "" {
			paths = append(paths, p)
		}
	}
	return paths, nil
}
