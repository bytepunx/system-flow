package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai stream restore commits the main branch's version of a path a rebased
// story branch undoes (ADR-0155, I-0124) onto the story's branch, in its
// worktree, so that the path leaves the story's record of undone paths and
// no longer blocks its acceptance. It is flai's git work, so ADR-0069's rule
// that agents never rebase or merge by hand stands.

// restored is what flai stream restore committed, as --json gives it.
type restored struct {
	Story  string   `json:"story"`
	Branch string   `json:"branch"`
	Commit string   `json:"commit"`
	Paths  []string `json:"paths"`
}

// restoreRefused is a restore refused before anything changed, as --json
// gives it under "refused": why, with the paths that make it so.
type restoreRefused struct {
	Story       string   `json:"story"`
	Reason      string   `json:"reason"`
	Status      string   `json:"status,omitempty"`
	Uncommitted []string `json:"uncommitted,omitempty"`
	NotHeld     []string `json:"not_held,omitempty"`
	Held        []string `json:"held,omitempty"`
}

func (r *restoreRefused) Error() string { return r.Reason }

func newStreamRestoreCmd(a *app) *cobra.Command {
	var paths, trailers []string
	c := &cobra.Command{
		Use:   "restore <story-id> --path <path>...",
		Short: "Commit main's version of a path the rebased story branch undoes onto the branch",
		Long: `Commits the main branch's version of each --path onto story/<story-id>, in
the story's worktree, as one commit,
"fix: [S-nnnn] restore main's <path>, which the rebased branch undid (ADR-0155)",
with each --trailer.

Use it for a path flai stream sync warns the rebased branch undoes main's
change to, which the acceptance preview (flai accept --dry-run) lists as a
blocker: the branch removes lines main added since the branch left it, as a
revert replayed over the same change on main does (I-0124). Once committed,
the branch no longer removes main's lines, so the path leaves the story's
record of undone paths and no longer blocks the acceptance.

It refuses, changing nothing, with exit 4: a story that is not in progress or
in review, one with no worktree, a worktree with uncommitted changes (each
named) or a rebase unfinished, a worktree with another branch checked out,
and a path the story's record does not hold (naming those it holds). With
--json it prints story, branch, commit, and paths, and a refusal as
{"refused": {"story", "reason", ...}} with status, uncommitted, not_held, or
held where they apply.`,
		Example: `  flai stream restore S-0232 --path scripts/flai-snapshot.sh
  flai stream restore S-0232 --path docs/a.md --path docs/b.md --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			res, err := a.restoreMain(repo, args[0], restorePaths(paths), trailers)
			var refused *restoreRefused
			if errors.As(err, &refused) {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": refused})
				}
				return &exitError{code: exitDocRefused, msg: refused.Error()}
			}
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			fmt.Fprintf(a.out, "restored main's version on %s in %s:\n", res.Branch, short(res.Commit))
			for _, p := range res.Paths {
				fmt.Fprintf(a.out, "  %s\n", p)
			}
			return nil
		},
	}
	c.Flags().StringArrayVar(&paths, "path", nil, "a path the story's record of undone paths holds, as git names it from the project's root (repeatable)")
	c.Flags().StringArrayVar(&trailers, "trailer", nil, "a trailer line for the commit (repeatable)")
	_ = c.MarkFlagRequired("path")
	return c
}

// restorePaths is the --path values as git names paths, each once, in the
// order given.
func restorePaths(given []string) []string {
	out := []string{}
	for _, p := range given {
		if p = path.Clean(filepath.ToSlash(strings.TrimSpace(p))); !slices.Contains(out, p) {
			out = append(out, p)
		}
	}
	return out
}

// restoreMain checks out the main branch's version of each path in the
// story's worktree and commits it on the story's branch, after refusing,
// with nothing changed, what restoreRefusal refuses. A git step that fails
// once the worktree is written puts the paths back as the branch has them.
func (a *app) restoreMain(repo *workitem.Repo, story string, paths, trailers []string) (restored, error) {
	it, err := repo.Get(story)
	if err != nil {
		return restored{}, err
	}
	worktree := repo.WorktreePath(it.ID)
	if r, err := a.restoreRefusal(repo, it, worktree, paths); r != nil || err != nil {
		if r != nil {
			return restored{}, r
		}
		return restored{}, err
	}
	main, err := storygit.MainBranch(a.runner, repo.MainRoot)
	if err != nil {
		return restored{}, fmt.Errorf("restore main's version of %s on %s: read the main branch in %s: %w", strings.Join(paths, ", "), storyBranch(it.ID), repo.MainRoot, err)
	}
	var written []string
	undo := func(err error) (restored, error) {
		if len(written) == 0 {
			return restored{}, err
		}
		if _, uerr := a.runner.Run(worktree, "git", append([]string{"restore", "-q", "--source=HEAD", "--staged", "--worktree", "--"}, written...)...); uerr != nil {
			return restored{}, fmt.Errorf("%w; putting %s back as %s has them failed too (%w): run git status in %s and git restore --source=HEAD --staged --worktree -- %s there", err, strings.Join(written, ", "), storyBranch(it.ID), uerr, worktree, strings.Join(written, " "))
		}
		return restored{}, fmt.Errorf("%w; %s are back as %s has them, and nothing was committed", err, strings.Join(written, ", "), storyBranch(it.ID))
	}
	// Every path is held, so main has a file there: a path main has no file
	// at gives the branch no line of main's to remove, and the record's
	// reader drops it before the refusal of a path it does not hold.
	for _, p := range paths {
		if _, err := a.runner.Run(worktree, "git", "checkout", main, "--", p); err != nil {
			return undo(fmt.Errorf("restore main's version of %s on %s: check it out from %s in %s: %w", p, storyBranch(it.ID), main, worktree, err))
		}
		written = append(written, p)
	}
	c, err := storygit.CommitPaths(storygit.CommitOptions{
		Runner: a.runner, Dir: worktree, Story: it.ID, Paths: paths, Type: "fix", Trailers: trailers,
		Subject: fmt.Sprintf("restore main's %s, which the rebased branch undid (ADR-0155)", strings.Join(paths, ", ")),
	})
	if err != nil {
		return undo(err)
	}
	if c == nil {
		return restored{}, fmt.Errorf("restore main's version of %s on %s: nothing changed, so nothing was committed; %s already holds %s's version", strings.Join(paths, ", "), storyBranch(it.ID), storyBranch(it.ID), main)
	}
	// Reading the record again drops the paths restored from it now rather
	// than at its next reader; it is a cache, so a failure is logged.
	if _, err := storygit.RecordedUndone(a.runner, repo, it.ID); err != nil {
		a.logger().Warn("the record of undone paths not rewritten after the restore; its next reader drops them", "component", "git", "story", it.ID, "err", err)
	}
	return restored{Story: it.ID, Branch: storyBranch(it.ID), Commit: c.Hash, Paths: paths}, nil
}

// restoreRefusal is why a restore of paths on the branch of it, checked out
// in worktree, is refused, nil when it is not: an item that is not a story in
// progress or in review, a worktree missing, with a rebase unfinished, with
// another branch checked out, or with uncommitted changes, and a path the
// story's record of undone paths does not hold.
func (a *app) restoreRefusal(repo *workitem.Repo, it *workitem.Item, worktree string, paths []string) (*restoreRefused, error) {
	branch := storyBranch(it.ID)
	refuse := func(format string, args ...any) *restoreRefused {
		return &restoreRefused{Story: it.ID, Reason: "flai stream restore " + it.ID + " changed nothing: " + fmt.Sprintf(format, args...)}
	}
	if it.Type != workitem.Story {
		return refuse("it is not a story; name the story (S-nnnn) whose branch undoes main's change"), nil
	}
	if it.Status != workitem.InProgress && it.Status != workitem.Review {
		r := refuse("it is in %s, and a path is restored only on the branch of a story in progress or in review", it.Status)
		r.Status = it.Status
		return r, nil
	}
	if st, err := os.Stat(worktree); err != nil || !st.IsDir() {
		return refuse("it has no worktree at %s; open it with flai stream open %s, then restore", relPath(repo.MainRoot, worktree), it.ID), nil
	}
	rel := relPath(repo.MainRoot, worktree)
	if storygit.RebaseInProgress(a.runner, worktree) {
		return refuse("a rebase is unfinished in %s; finish it with git rebase --continue there, or undo it with git rebase --abort, then restore", rel), nil
	}
	head, err := a.runner.Run(worktree, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("restore on %s: read the branch checked out in %s: %w", branch, rel, err)
	}
	if head = strings.TrimSpace(head); head != branch {
		return refuse("%s has %s checked out, not %s; check out %s there, then restore", rel, head, branch, branch), nil
	}
	dirty, err := storygit.Uncommitted(a.runner, worktree)
	if err != nil {
		return nil, fmt.Errorf("restore on %s: read the uncommitted changes in %s: %w", branch, rel, err)
	}
	if len(dirty) > 0 {
		r := refuse("its worktree %s has uncommitted changes in %s; commit them on %s, or discard them, then restore", rel, workitem.Shorten(dirty, 10), branch)
		r.Uncommitted = dirty
		return r, nil
	}
	undone, err := storygit.RecordedUndone(a.runner, repo, it.ID)
	if err != nil {
		return nil, fmt.Errorf("restore on %s: %w", branch, err)
	}
	held := make([]string, 0, len(undone))
	for _, u := range undone {
		held = append(held, u.Path)
	}
	var notHeld []string
	for _, p := range paths {
		if !slices.Contains(held, p) {
			notHeld = append(notHeld, p)
		}
	}
	if len(notHeld) == 0 {
		return nil, nil
	}
	var r *restoreRefused
	if len(held) == 0 {
		r = refuse("its record of the paths %s undoes of the main branch holds none, so %s is not one to restore; flai stream sync records them", branch, strings.Join(notHeld, ", "))
	} else {
		r = refuse("its record of the paths %s undoes of the main branch does not hold %s; it holds %s", branch, strings.Join(notHeld, ", "), strings.Join(held, ", "))
	}
	r.NotHeld, r.Held = notHeld, held
	return r, nil
}
