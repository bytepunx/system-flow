package cmd

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The checks flai stream sync runs once its rebase is clean (S-0131,
// ADR-0046): a trial merge with every other open story's branch, which
// writes nothing to any worktree, to catch the overlaps the declared
// touches missed while both stories are still open, and the paths the
// branch changed outside the story's claim, so that it is widened. A conflict is a thread
// flai writes, one per pair of stories: a thread whose last entry is flai's
// awaits every agent and the designer, so both stories' agents see it in
// the MCP inbox and the designer in the dashboard's. The generated files,
// design/issues/summary.md today, are left out of a pair's conflicts: the
// rebase writes them again when it stops on them alone, so a pair whose only
// conflict is one merges cleanly as far as the check is concerned (ADR-0098).
// A pair's conflicts are the paths both stories changed since each left the
// main branch: the trial merge's base is where the two branches meet, which
// is the other branch's old base when it is stale, so it also stops where
// main has since changed what the other did. That is the other story's own
// rebase to settle, not a conflict between the two (I-0064, S-0251).

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
	Outside  []string      `json:"outside_touches"`
}

// checkSync lists the paths story's branch changed since base outside its
// claim, and trial-merges the branch with the branch of every other story
// in progress or in review that has one, in ID order. A pair's conflicts
// are the conflicting paths both branches changed since they left base, less
// the generated files, so a pair that conflicts in them alone, or only where
// main changed what one of them did, is clean (ADR-0098, I-0064).
func (a *app) checkSync(repo *workitem.Repo, story *workitem.Item, base string) (syncChecks, error) {
	out := syncChecks{Branches: []branchCheck{}, Outside: []string{}}
	items, err := repo.List(false)
	if err != nil {
		return out, err
	}
	mine, err := a.branchChanges(repo.MainRoot, base, storyBranch(story.ID))
	if err != nil {
		return out, err
	}
	out.Outside = outsideClaim(repo, story, items, mine)
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
	generated := generatedPaths(repo)
	for _, o := range others {
		conflicts, err := a.trialMerge(repo.MainRoot, storyBranch(story.ID), storyBranch(o.ID))
		if err != nil {
			return out, fmt.Errorf("trial merge of %s with %s: %w", storyBranch(story.ID), storyBranch(o.ID), err)
		}
		theirs, err := a.branchChanges(repo.MainRoot, base, storyBranch(o.ID))
		if err != nil {
			return out, err
		}
		conflicts = changedByBoth(withoutGenerated(conflicts, generated), mine, theirs)
		out.Branches = append(out.Branches, branchCheck{Story: o.ID, Status: o.Status, Branch: storyBranch(o.ID), Clean: len(conflicts) == 0, Conflicts: conflicts})
	}
	return out, nil
}

// withoutGenerated is conflicts less the generated files, which sync and
// acceptance write again when a rebase stops on them alone, so no agent
// settles them (ADR-0098).
func withoutGenerated(conflicts, generated []string) []string {
	out := []string{}
	for _, p := range conflicts {
		if !coveredBy(p, generated) {
			out = append(out, p)
		}
	}
	return out
}

// changedByBoth is conflicts less the paths only one of the two branches
// changed, mine and theirs: those conflict with what main brought, not with
// each other (I-0064).
func changedByBoth(conflicts, mine, theirs []string) []string {
	out := []string{}
	for _, p := range conflicts {
		if slices.Contains(mine, p) && slices.Contains(theirs, p) {
			out = append(out, p)
		}
	}
	return out
}

// branchChanges is the paths branch changed since it left base, from the
// merge base of the two.
func (a *app) branchChanges(root, base, branch string) ([]string, error) {
	diff, err := a.runner.Run(root, "git", "diff", "--name-only", base+"..."+branch)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, p := range strings.Split(diff, "\n") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out, nil
}

// outsideClaim is the paths of changed, what story's branch changed since it
// left the main branch, that no entry of its claim (ADR-0046, ADR-0096: its
// touches, a folder narrowed to its tasks' touches inside it, and its open
// tasks', a component as its path) covers. The wip folder is flai's and is
// left out.
func outsideClaim(repo *workitem.Repo, story *workitem.Item, items []*workitem.Item, changed []string) []string {
	claim := workitem.NewHolds(items, repo.Manifest.Projects).Claim(story)
	wip := strings.TrimSuffix(repo.Manifest.Layout["wip"], "/")
	if wip == "" {
		wip = "wip"
	}
	out := []string{}
	for _, p := range changed {
		if coveredBy(p, []string{wip}) || coveredBy(p, claim) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// coveredBy says whether path is one of entries or lies below one.
func coveredBy(path string, entries []string) bool {
	for _, e := range entries {
		if path == e || strings.HasPrefix(path, e+"/") {
			return true
		}
	}
	return false
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
	if len(c.Outside) > 0 {
		fmt.Fprintf(w, "%s changed %s outside %s's touches: %s\n", mine, plural(len(c.Outside), "path"), story.ID, strings.Join(c.Outside, ", "))
		fmt.Fprintf(w, "widen them so that stories that overlap wait: flai touches %s --add %s\n", story.ID, strings.Join(c.Outside, " "))
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
