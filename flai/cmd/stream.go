package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newStreamCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "stream",
		Short: "Open and append to agent narratives in wip/agents",
		Long: `A stream is the narrative for one story. Set FLAI_AGENT and FLAI_SESSION
so entries record who wrote them.`,
	}
	c.AddCommand(newStreamDiffCmd(a), newStreamOpenCmd(a), newStreamLogCmd(a), newStreamSyncCmd(a), newStreamAnswerCmd(a))
	return c
}

func newStreamOpenCmd(a *app) *cobra.Command {
	var noBranch bool
	c := &cobra.Command{
		Use:   "open <story-id>",
		Short: "Create the narrative for a story and check out its branch in a worktree",
		Long: `Creates wip/agents/<story-id>.md from the template and, in a git repository,
the branch story/<story-id> from the main branch, checked out in a worktree
under .flai-cache/worktrees/<story-id> (ADR-0019). Work there; wip/ stays in
the main checkout. Run flai stream sync at every task transition.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			story, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			agent, session := agentIdentity()
			n, err := repo.OpenStream(story, workitem.StreamOptions{Agent: agent, Session: session, Now: a.now()})
			if err != nil {
				return err
			}
			if err := a.refreshIndex(repo); err != nil {
				return err
			}
			branch, wt := "", ""
			if !noBranch {
				branch, wt, _, err = a.openStoryBranch(repo, story.ID)
				if err != nil {
					return fmt.Errorf("narrative opened but the branch was not: %w", err)
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]string{"stream": n.Stream, "path": n.Path, "branch": branch, "worktree": wt})
			}
			fmt.Fprintf(a.out, "opened %s\n", relPath(repo.Root, n.Path))
			if branch != "" {
				fmt.Fprintf(a.out, "branch %s checked out at %s; work there and run flai stream sync %s at each task transition\n", branch, relPath(repo.MainRoot, wt), story.ID)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&noBranch, "no-branch", false, "narrative only; no branch or worktree")
	return c
}

func newStreamSyncCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "sync <story-id>",
		Short: "Rebase the story branch onto the main branch in its worktree",
		Long: `Rebases story/<story-id> onto the branch checked out in the main checkout,
stashing and restoring uncommitted work. Conflicts stop the rebase inside the
worktree and are listed; resolve them, run git rebase --continue there, and
sync again.

After a clean rebase it trial-merges the branch with the branch of every other
story in progress or in review (git merge-tree --write-tree, git 2.38 or
newer), writing nothing to any worktree, and lists each branch it conflicts
with and the conflicting paths. Each conflicting pair of stories has one
thread, written by flai on the story that synced, which both stories' agents
and the designer see in their inboxes; a sync that finds the pair merging
cleanly again resolves it.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			base, conflicts, err := a.syncStoryBranch(repo, it.ID)
			if err != nil {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"story": it.ID, "branch": storyBranch(it.ID), "base": base, "conflicts": conflicts, "ok": false})
				}
				return err
			}
			checks, cerr := a.checkSync(repo, it)
			if cerr != nil {
				a.logger().Warn("story branch checks failed after the rebase", "component", "git", "story", it.ID, "err", cerr)
			} else if err := a.reportConflicts(repo, it, &checks); err != nil {
				a.logger().Warn("conflict threads not written", "component", "threads", "story", it.ID, "err", err)
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"story": it.ID, "branch": storyBranch(it.ID), "base": base, "conflicts": conflicts, "ok": true,
					"branches": checks.Branches, "trial_merge_skipped": checks.Skipped})
			}
			fmt.Fprintf(a.out, "%s is rebased onto %s\n", storyBranch(it.ID), base)
			printSyncChecks(a.out, it, checks)
			return nil
		},
	}
}

func newStreamLogCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "log <story-id> \"<entry>\"",
		Short: "Append a timestamped entry to a story's narrative log",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			agent, session := agentIdentity()
			n, err := repo.LogStream(args[0], args[1], workitem.StreamOptions{Agent: agent, Session: session, Now: a.now()})
			if err != nil {
				return err
			}
			if err := a.refreshIndex(repo); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]string{"stream": n.Stream, "updated": n.Updated})
			}
			fmt.Fprintf(a.out, "logged to %s at %s\n", n.Stream, n.Updated)
			return nil
		},
	}
}

func newStreamAnswerCmd(a *app) *cobra.Command {
	var by string
	c := &cobra.Command{
		Use:   "answer <story-id> \"<question>\" \"<answer>\"",
		Short: "Answer an open question in a story's narrative: it moves to Decisions",
		Long: `Removes the bullet matching <question> (exactly as it reads under ##
Open questions) and records it, with the answer, under ## Decisions
(design/system/agent-narrative.md: "Answered questions move to Decisions").
Refuses if no open question matches. This is for a hand-written question, a
freeform note an agent left with no answer/resolve of its own; a question
mirrored from a thread is answered with flai thread reply instead.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			n, err := repo.AnswerOpenQuestion(args[0], args[1], args[2], a.threadAuthor(by), a.now())
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]string{"stream": n.Stream, "updated": n.Updated})
			}
			fmt.Fprintf(a.out, "%s answered; moved to Decisions in %s\n", args[1], relPath(repo.MainRoot, n.Path))
			return nil
		},
	}
	c.Flags().StringVar(&by, "by", "", "who answered it (default: FLAI_AGENT, then config author)")
	return c
}
