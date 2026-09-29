package cmd

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A cancellation takes everything open under the item with it (ADR-0028),
// so flai move says what that is before it happens, asks on a terminal, and
// afterwards names what it left for a person: git is never touched.

// cancelItem runs a move to cancelled: plan, show, confirm, apply, report.
func (a *app) cancelItem(repo *workitem.Repo, it *workitem.Item, by, reason string, dryRun bool) error {
	res, err := preview.Cancel(a.runner, repo, it, by, reason, a.now())
	if err != nil {
		return err
	}
	res.DryRun = dryRun
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
	res.LeftBehind = preview.Left(a.runner, repo, it, res.Cancelled)
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

func (a *app) printLeftBehind(left []preview.LeftBehind, dryRun bool) {
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
	return a.prompts().Confirm(title, false)
}

// plural is "1 item" or "3 items".
func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
