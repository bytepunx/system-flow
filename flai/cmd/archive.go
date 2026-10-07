package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/threads"
)

func newArchiveCmd(a *app) *cobra.Command {
	var dryRun bool
	c := &cobra.Command{
		Use:   "archive [id...]",
		Short: "Move done and cancelled items and their narratives to wip/archive",
		Long: `Without IDs, every closed epic, every closed story with its tasks, and
every closed task whose story is no longer on the board is archived. With
IDs, each must be done or cancelled; a story brings its tasks and narrative.

Every thread still open or answered on an item archived is resolved, as
"<id> was archived", so that none is left open on an archived item (I-0073),
and every conversation still open of a story archived is closed, with the
entry "Closed: <id> was archived" (ADR-0120). Nothing is committed. --dry-run
names the threads it would resolve and the conversations it would close.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			items, err := repo.List(false)
			if err != nil {
				return err
			}
			plan, err := repo.PlanArchive(items, args)
			if err != nil {
				return err
			}
			ids := make([]string, len(plan.Items))
			for i, it := range plan.Items {
				ids[i] = it.ID
			}
			resolved, closed := []string{}, []string{}
			conversations, err := openConversations(repo, ids)
			if err != nil {
				return fmt.Errorf("read the conversations of the stories to archive: %w", err)
			}
			if dryRun {
				open, err := threads.OnItems(repo, ids)
				if err != nil {
					return fmt.Errorf("read the threads on the items to archive: %w", err)
				}
				for _, th := range open {
					resolved = append(resolved, th.ID)
				}
				for _, c := range conversations {
					closed = append(closed, c.ID)
				}
			} else {
				if err := repo.Archive(plan); err != nil {
					return err
				}
				for _, id := range ids {
					r, err := threads.ResolveOnItems(repo, []string{id}, a.author(), id+" was archived", a.now())
					resolved = append(resolved, r...)
					if err != nil {
						return fmt.Errorf("resolve the threads on %s, which is archived: %w", id, err)
					}
				}
				shut, err := messages.CloseOn(repo, ids, a.author(), "archived", a.now())
				closed = append(closed, shut...)
				if err != nil {
					return fmt.Errorf("close the conversations of what is archived: %w", err)
				}
				conversations = among(conversations, closed)
				if err := a.refreshIndex(repo); err != nil {
					return err
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"archived": ids, "narratives": plan.Narratives, "resolved_threads": resolved, "closed_messages": closed, "dry_run": dryRun})
			}
			if len(plan.Items) == 0 {
				fmt.Fprintln(a.out, "nothing to archive")
				return nil
			}
			verb, resolve, closing := "archived", "resolved", "closed"
			if dryRun {
				verb, resolve, closing = "would archive", "would resolve", "would close"
			}
			for _, it := range plan.Items {
				fmt.Fprintf(a.out, "%s %s %s\n", verb, it.ID, it.Title)
			}
			for _, n := range plan.Narratives {
				fmt.Fprintf(a.out, "%s narrative %s\n", verb, relPath(repo.Root, n))
			}
			if len(resolved) > 0 {
				fmt.Fprintf(a.out, "%s %s, open or answered on what it archives\n", resolve, strings.Join(resolved, ", "))
			}
			for _, line := range closingLines(conversations, ids) {
				fmt.Fprintf(a.out, "%s %s\n", closing, line)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print the plan without moving anything")
	return c
}
