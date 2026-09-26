package cmd

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The checks flai stream sync runs once its rebase is clean (S-0131,
// ADR-0046): a trial merge with every other open story's branch, which
// writes nothing to any worktree, to catch the overlaps the declared
// touches missed while both stories are still open. A conflict is a thread
// flai writes, one per pair of stories: a thread whose last entry is flai's
// awaits every agent and the designer, so both stories' agents see it in
// the MCP inbox and the designer in the dashboard's.

// branchCheck is how another open story's branch merges with this one.
type branchCheck struct {
	Story     string   `json:"story"`
	Status    string   `json:"status"`
	Branch    string   `json:"branch"`
	Clean     bool     `json:"clean"`
	Conflicts []string `json:"conflicts"`
	Thread    string   `json:"thread,omitempty"`
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
		fmt.Fprintf(w, "%s conflicts with %s (%s) in %s; see %s\n", mine, b.Branch, state, strings.Join(b.Conflicts, ", "), b.Thread)
	}
}

// conflictAuthor writes the conflict threads.
const conflictAuthor = "flai"

var conflictTitlePattern = regexp.MustCompile(`^(S-\d+) and (S-\d+) conflict when merged$`)

// conflictTitle names the pair's thread, the same whichever story synced.
func conflictTitle(a, b string) string {
	if b < a {
		a, b = b, a
	}
	return a + " and " + b + " conflict when merged"
}

// conflictText is the entry for a pair and its conflicting paths, the same
// whichever story synced, so that a sync that finds nothing new writes
// nothing.
func conflictText(a, b string, paths []string) string {
	if b < a {
		a, b = b, a
	}
	var list strings.Builder
	for _, p := range paths {
		fmt.Fprintf(&list, "- `%s`\n", p)
	}
	return fmt.Sprintf("A trial merge of %s with %s at flai stream sync conflicts in:\n\n%s\n"+
		"Whichever of %s and %s is accepted second will stop on these paths when it rebases. "+
		"Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. "+
		"Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.",
		storyBranch(a), storyBranch(b), list.String(), a, b)
}

// reportConflicts keeps one open thread for each pair of story and another
// open story whose branches conflict: it opens one on story for a new
// conflict, adds an entry when the paths change, and resolves it once the
// pair merges cleanly or the other story is no longer open. It sets each
// conflicting branch's Thread.
func (a *app) reportConflicts(repo *workitem.Repo, story *workitem.Item, checks *syncChecks) error {
	if checks.Skipped != "" {
		return nil
	}
	all, err := threads.List(repo)
	if err != nil {
		return err
	}
	open := map[string]*threads.Thread{}
	for _, th := range all {
		if th.Open() && conflictTitlePattern.MatchString(th.Title) {
			open[th.Title] = th
		}
	}
	now := a.now()
	var changed []*threads.Thread
	for i, b := range checks.Branches {
		title := conflictTitle(story.ID, b.Story)
		th := open[title]
		delete(open, title)
		text := conflictText(story.ID, b.Story, b.Conflicts)
		switch {
		case !b.Clean && th == nil:
			if th, err = threads.New(repo, threads.NewOptions{Title: title, On: story.ID, Author: conflictAuthor, Text: text, Now: now}); err != nil {
				return err
			}
			changed = append(changed, th)
		case !b.Clean && lastEntryBy(th, conflictAuthor) != text:
			if th, err = threads.Reply(repo, th.ID, conflictAuthor, text, now); err != nil {
				return err
			}
			changed = append(changed, th)
		case b.Clean && th != nil:
			lo, hi := storyBranch(story.ID), b.Branch
			if hi < lo {
				lo, hi = hi, lo
			}
			if th, err = threads.Resolve(repo, th.ID, conflictAuthor, fmt.Sprintf("%s and %s merge cleanly at the sync of %s", lo, hi, story.ID), now); err != nil {
				return err
			}
			changed = append(changed, th)
			continue
		}
		if th != nil {
			checks.Branches[i].Thread = th.ID
		}
	}
	// A pair's thread whose other story is no longer open: its branch is
	// merged or dropped, and what is left of the conflict is this story's
	// own rebase to settle.
	for title, th := range open {
		m := conflictTitlePattern.FindStringSubmatch(title)
		other := ""
		switch story.ID {
		case m[1]:
			other = m[2]
		case m[2]:
			other = m[1]
		default:
			continue
		}
		state := "gone"
		if o, err := repo.Get(other); err == nil {
			if o.Status == workitem.InProgress || o.Status == workitem.Review {
				continue
			}
			state = o.Status
		}
		if th, err = threads.Resolve(repo, th.ID, conflictAuthor, fmt.Sprintf("%s is %s, no longer open, at the sync of %s", other, state, story.ID), now); err != nil {
			return err
		}
		changed = append(changed, th)
	}
	for _, th := range changed {
		if s := threads.StoryOf(repo, th); s != "" {
			if err := threads.MirrorNarrative(repo, s); err != nil {
				return err
			}
		}
	}
	return nil
}

// lastEntryBy is the text of author's last entry on th, "" when none.
func lastEntryBy(th *threads.Thread, author string) string {
	e := th.Entries()
	for i := len(e) - 1; i >= 0; i-- {
		if e[i].Author == author {
			return e[i].Text
		}
	}
	return ""
}
