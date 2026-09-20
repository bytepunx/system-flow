package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A cancellation takes everything open under the item with it (ADR-0028),
// so flai move says what that is before it happens, asks on a terminal, and
// afterwards names what it left for a person: git is never touched.

// leftBehind is what a cancelled story still has on disk and in git.
type leftBehind struct {
	ID        string `json:"id"`
	Narrative string `json:"narrative,omitempty"`
	Branch    string `json:"branch,omitempty"`
	Worktree  string `json:"worktree,omitempty"`
}

// cancelResult is the JSON of a move to cancelled, or under DryRun of what
// it would do.
type cancelResult struct {
	ID         string              `json:"id"`
	Status     string              `json:"status"`
	DryRun     bool                `json:"dry_run,omitempty"`
	Cancelled  []workitem.Cascaded `json:"cancelled"`
	LeftBehind []leftBehind        `json:"left_behind"`
	Warnings   []string            `json:"warnings"`
}

// cancelItem runs a move to cancelled: plan, show, confirm, apply, report.
func (a *app) cancelItem(repo *workitem.Repo, it *workitem.Item, by, reason string, dryRun bool) error {
	items, err := repo.List(false)
	if err != nil {
		return err
	}
	// The rules for the item itself, tried on a copy so that a refusal is
	// reported before anything is shown or asked.
	probe := *it
	probe.Transitions = append([]workitem.Transition(nil), it.Transitions...)
	if _, err := repo.Move(&probe, workitem.Cancelled, workitem.MoveOptions{By: by, Reason: reason, Now: a.now(), Items: items}); err != nil {
		return err
	}
	res := cancelResult{ID: it.ID, Status: it.Status, DryRun: dryRun, Cancelled: []workitem.Cascaded{}, LeftBehind: []leftBehind{}, Warnings: []string{}}
	for _, c := range workitem.CancelPlan(items, it.ID) {
		res.Cancelled = append(res.Cancelled, workitem.Cascaded{ID: c.ID, Type: c.Type, Title: c.Title, From: c.Status})
	}
	if !a.jsonOut && len(res.Cancelled) > 0 {
		verb := "also cancels"
		if dryRun {
			verb = "would also cancel"
		}
		fmt.Fprintf(a.out, "Cancelling %s %s %s:\n", it.ID, verb, plural(len(res.Cancelled), "item"))
		for _, c := range res.Cancelled {
			indent, note := "  ", ""
			if c.Type == workitem.Task && it.Type == workitem.Epic {
				indent = "    "
			}
			if c.From == workitem.Review {
				note = "  (in review: its work stays on its branch, unmerged)"
			}
			fmt.Fprintf(a.out, "%s%-5s %s  %-11s %s%s\n", indent, c.Type, c.ID, c.From, c.Title, note)
		}
	}
	if dryRun {
		res.LeftBehind = a.leftBehind(repo, it, res.Cancelled)
		if a.jsonOut {
			return a.printJSON(res)
		}
		if len(res.Cancelled) == 0 {
			fmt.Fprintf(a.out, "%s would be cancelled; nothing open under it\n", it.ID)
		}
		a.printLeftBehind(res.LeftBehind, true)
		return nil
	}
	if len(res.Cancelled) > 0 && !a.yes && a.isTerminal() {
		ok, err := a.ask(fmt.Sprintf("Cancel %s and %s under it?", it.ID, plural(len(res.Cancelled), "item")))
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("nothing was cancelled")
		}
	}
	done, err := repo.TransitionAll(it, workitem.Cancelled, by, reason, a.now())
	if err != nil {
		return err
	}
	res.Status, res.Cancelled = it.Status, done.Cancelled
	if done.Warnings != nil {
		res.Warnings = done.Warnings
	}
	for _, w := range res.Warnings {
		a.logger().Warn("workflow policy warning", "component", "workitem", "item", it.ID, "detail", w)
	}
	res.LeftBehind = a.leftBehind(repo, it, res.Cancelled)
	if a.jsonOut {
		return a.printJSON(res)
	}
	fmt.Fprintf(a.out, "%s → %s\n", it.ID, it.Status)
	if n := len(res.Cancelled); n > 0 {
		fmt.Fprintf(a.out, "%s cancelled with it\n", plural(n, "item"))
	}
	a.printLeftBehind(res.LeftBehind, false)
	return nil
}

// leftBehind looks at each story in the cancellation for a narrative, a
// story branch, and a worktree. Nothing is removed: unmerged work is the
// operator's to keep or delete.
func (a *app) leftBehind(repo *workitem.Repo, it *workitem.Item, with []workitem.Cascaded) []leftBehind {
	ids := []string{}
	if it.Type == workitem.Story {
		ids = append(ids, it.ID)
	}
	for _, c := range with {
		if c.Type == workitem.Story {
			ids = append(ids, c.ID)
		}
	}
	git := a.inGitWorkTree(repo.MainRoot)
	out := []leftBehind{}
	for _, id := range ids {
		l := leftBehind{ID: id}
		if _, err := os.Stat(repo.NarrativePath(id)); err == nil {
			l.Narrative = relPath(repo.MainRoot, repo.NarrativePath(id))
		}
		if git && a.branchExists(repo.MainRoot, storyBranch(id)) {
			l.Branch = storyBranch(id)
		}
		if _, err := os.Stat(repo.WorktreePath(id)); err == nil {
			l.Worktree = relPath(repo.MainRoot, repo.WorktreePath(id))
		}
		if l.Narrative != "" || l.Branch != "" || l.Worktree != "" {
			out = append(out, l)
		}
	}
	return out
}

func (a *app) printLeftBehind(left []leftBehind, dryRun bool) {
	if len(left) == 0 {
		return
	}
	if dryRun {
		fmt.Fprintln(a.out, "Would be left as it is, for you to keep or remove:")
	} else {
		fmt.Fprintln(a.out, "Left as it is, for you to keep or remove (flai archive moves the narratives):")
	}
	for _, l := range left {
		parts := []string{}
		if l.Narrative != "" {
			parts = append(parts, "narrative "+l.Narrative)
		}
		if l.Branch != "" {
			parts = append(parts, "branch "+l.Branch)
		}
		if l.Worktree != "" {
			parts = append(parts, "worktree "+l.Worktree)
		}
		fmt.Fprintf(a.out, "  %s: %s\n", l.ID, strings.Join(parts, ", "))
	}
}

// ask puts a yes or no question to the terminal.
func (a *app) ask(title string) (bool, error) {
	if a.confirm != nil {
		return a.confirm(title)
	}
	ok := false
	err := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(title).Value(&ok))).Run()
	return ok, err
}

// plural is "1 item" or "3 items".
func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
