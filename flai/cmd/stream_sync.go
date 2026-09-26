package cmd

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The checks flai stream sync runs once its rebase is clean (S-0131,
// ADR-0046): a trial merge with every other open story's branch, which
// writes nothing to any worktree, to catch the overlaps the declared
// touches missed while both stories are still open.

// branchCheck is how another open story's branch merges with this one.
type branchCheck struct {
	Story     string   `json:"story"`
	Status    string   `json:"status"`
	Branch    string   `json:"branch"`
	Clean     bool     `json:"clean"`
	Conflicts []string `json:"conflicts"`
}

// syncChecks is what the checks found; Skipped says why the trial merge did
// not run, when it did not.
type syncChecks struct {
	Branches []branchCheck `json:"branches"`
	Skipped  string        `json:"trial_merge_skipped,omitempty"`
}

// checkSync trial-merges story's branch with the branch of every other story
// in progress or in review that has one, in ID order.
func (a *app) checkSync(repo *workitem.Repo, story *workitem.Item) (syncChecks, error) {
	out := syncChecks{Branches: []branchCheck{}}
	items, err := repo.List(false)
	if err != nil {
		return out, err
	}
	var others []*workitem.Item
	for _, it := range items {
		if it.Type == workitem.Story && it.ID != story.ID && (it.Status == workitem.InProgress || it.Status == workitem.Review) && a.branchExists(repo.MainRoot, storyBranch(it.ID)) {
			others = append(others, it)
		}
	}
	if len(others) == 0 {
		return out, nil
	}
	sort.Slice(others, func(i, j int) bool { return others[i].ID < others[j].ID })
	if v, err := gitver.Installed(a.runner); err != nil || !v.AtLeast(gitver.MergeTree) {
		have := "unreadable"
		if err == nil {
			have = v.String()
		}
		out.Skipped = fmt.Sprintf("git %s is older than %s, the first with merge-tree --write-tree", have, gitver.MergeTree)
		a.logger().Warn("git is too old to trial-merge story branches, skipping", "component", "git", "git", have, "needs", gitver.MergeTree.String())
		return out, nil
	}
	for _, o := range others {
		conflicts, err := a.trialMerge(repo.MainRoot, storyBranch(story.ID), storyBranch(o.ID))
		if err != nil {
			return out, fmt.Errorf("trial merge of %s with %s: %w", storyBranch(story.ID), storyBranch(o.ID), err)
		}
		out.Branches = append(out.Branches, branchCheck{Story: o.ID, Status: o.Status, Branch: storyBranch(o.ID), Clean: len(conflicts) == 0, Conflicts: conflicts})
	}
	return out, nil
}

// trialMerge merges two branches in git's object store only and returns the
// paths that conflict, none when they merge cleanly. merge-tree exits 1 on
// conflicts, printing the tree and then one conflicted path per line.
func (a *app) trialMerge(root, ours, theirs string) ([]string, error) {
	out, err := a.runner.Run(root, "git", "merge-tree", "--write-tree", "--name-only", "--no-messages", ours, theirs)
	if err == nil {
		return []string{}, nil
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 {
		return nil, err
	}
	lines := strings.Split(out, "\n")
	seen := map[string]bool{}
	paths := []string{}
	for _, l := range lines[1:] {
		l = strings.TrimSpace(l)
		if l == "" {
			break
		}
		if !seen[l] {
			seen[l] = true
			paths = append(paths, l)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

// printSyncChecks says, one line each, how story's branch merges with the
// other open branches.
func printSyncChecks(w io.Writer, story *workitem.Item, c syncChecks) {
	mine := storyBranch(story.ID)
	if c.Skipped != "" {
		fmt.Fprintf(w, "no trial merge with the other open story branches: %s\n", c.Skipped)
	}
	for _, b := range c.Branches {
		state := strings.ReplaceAll(b.Status, "-", " ")
		if state == workitem.Review {
			state = "in review"
		}
		if b.Clean {
			fmt.Fprintf(w, "%s merges cleanly with %s (%s)\n", mine, b.Branch, state)
			continue
		}
		fmt.Fprintf(w, "%s conflicts with %s (%s) in %s\n", mine, b.Branch, state, strings.Join(b.Conflicts, ", "))
	}
}
