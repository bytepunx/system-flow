package cmd

import (
	"errors"
	"fmt"
	"os"

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
the main checkout. Run flai stream sync at every task transition. The
narrative records the host that opened it.

A story whose narrative exists and whose worktree does not, such as one begun
on another host (ADR-0064), is reopened: the narrative is kept, with this
host, agent, and session recorded, and story/<story-id> is checked out from
this clone's branch, else fetched from the remote (origin) when it has one,
else created from the main branch. It is refused when the worktree exists
too.`,
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
			opt := workitem.StreamOptions{Agent: agent, Session: session, Host: workitem.ThisHost(), Now: a.now()}
			n, err := repo.OpenStream(story, opt)
			var exists *workitem.StreamExistsError
			reopen := errors.As(err, &exists)
			if reopen {
				if !a.canReopen(repo, story.ID, noBranch) {
					return err
				}
				opt.Agent = os.Getenv("FLAI_AGENT")
				n, err = repo.ReopenStream(story, opt)
			}
			if err != nil {
				return err
			}
			if err := a.refreshIndex(repo); err != nil {
				return err
			}
			branch, wt, from := "", "", ""
			if !noBranch {
				branch, wt, from, err = a.openStoryBranch(repo, story.ID, reopen)
				if err != nil {
					return fmt.Errorf("narrative opened but the branch was not: %w", err)
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"stream": n.Stream, "path": n.Path, "branch": branch, "worktree": wt, "reopened": reopen, "from": from})
			}
			if reopen {
				fmt.Fprintf(a.out, "reopened %s\n", relPath(repo.Root, n.Path))
			} else {
				fmt.Fprintf(a.out, "opened %s\n", relPath(repo.Root, n.Path))
			}
			if branch != "" {
				said := ""
				if reopen {
					said = fromSays(from)
				}
				fmt.Fprintf(a.out, "branch %s%s checked out at %s; work there and run flai stream sync %s at each task transition\n", branch, said, relPath(repo.MainRoot, wt), story.ID)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&noBranch, "no-branch", false, "narrative only; no branch or worktree")
	return c
}

// canReopen reports whether stream open may take up story's existing
// narrative: only to check out a worktree it does not have (ADR-0064).
func (a *app) canReopen(repo *workitem.Repo, story string, noBranch bool) bool {
	if noBranch || !a.inGitWorkTree(repo.MainRoot) {
		return false
	}
	_, err := os.Stat(repo.WorktreePath(story))
	return os.IsNotExist(err)
}

// fromSays says where a reopened story's branch came from, when it did not
// come from this clone.
func fromSays(from string) string {
	switch from {
	case "", branchLocal:
		return ""
	case branchMain:
		return " (new, from the main branch)"
	}
	return " (fetched from " + from + ")"
}

func newStreamSyncCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "sync <story-id>",
		Short: "Rebase the story branch onto the main branch in its worktree",
		Long: `Rebases story/<story-id> onto the branch checked out in the main checkout,
inside its worktree (ADR-0069). It never stashes: a worktree with uncommitted
changes is refused, touching nothing, and each uncommitted path is named, so
commit each task on the story branch before you sync. A worktree with a rebase
already in progress is refused too.

Conflicts stop the rebase inside the worktree. Each conflicting path is listed
on a line of its own, with how to continue (resolve each path, git add it, run
git rebase --continue in the worktree, and sync again) and how to abort (git
rebase --abort in the worktree, which puts the branch back as it was before
the sync). It exits non-zero; with --json it prints uncommitted, conflicts,
rebase_in_progress, continue, and abort, with ok false.

A stop whose only conflicting path is a generated file, which today is
design/issues/summary.md alone, does not wait for you (ADR-0098). Sync writes
the file again from the issue files in the worktree at that stop, git adds
it, and continues the rebase, until the rebase finishes or stops on another
path. A stop where other paths conflict too is left for you, and lists them
all, summary.md included; once the others are resolved, flai issue summary
run in the worktree writes it. flai accept syncs the same way. A rebase
already in progress is never continued for you.

After a clean rebase it trial-merges the branch with the branch of every other
story in progress or in review (git merge-tree --write-tree, git 2.38 or
newer), writing nothing to any worktree, and lists each branch it conflicts
with and the conflicting paths. A pair's conflicts are the paths both branches
changed since they left the main branch, not what the main branch brought
since: a path where main has since changed what a branch behind it did is
that story's own rebase to settle (I-0064). Generated files are left out, so a pair
whose only conflict is summary.md counts as clean. Each conflicting pair of stories
has one thread, written by flai on the story that synced, which both stories'
agents and the designer see in their inboxes; a sync that finds the pair
merging cleanly again resolves it. It also lists the paths the branch changed since
the main branch that the story's touches, and its open tasks', do not cover,
so that they are widened with flai touches.`,
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
				var stop *syncStopped
				if !errors.As(err, &stop) {
					if a.jsonOut {
						_ = a.printJSON(map[string]any{"story": it.ID, "branch": storyBranch(it.ID), "base": base, "conflicts": conflicts, "ok": false})
					}
					return err
				}
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"story": it.ID, "branch": storyBranch(it.ID), "base": base, "worktree": stop.Worktree, "ok": false,
						"uncommitted": nonNil(stop.Uncommitted), "conflicts": nonNil(stop.Conflicts), "rebase_in_progress": stop.RebaseInProgress,
						"continue": stop.Continue, "abort": stop.Abort})
				} else {
					fmt.Fprint(a.out, stop.Report())
				}
				return stop
			}
			checks, cerr := a.checkSync(repo, it, base)
			if cerr != nil {
				a.logger().Warn("story branch checks failed after the rebase", "component", "git", "story", it.ID, "err", cerr)
			} else if err := a.reportConflicts(repo, it, &checks); err != nil {
				a.logger().Warn("conflict threads not written", "component", "threads", "story", it.ID, "err", err)
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"story": it.ID, "branch": storyBranch(it.ID), "base": base, "conflicts": conflicts, "ok": true,
					"branches": checks.Branches, "trial_merge_skipped": checks.Skipped, "outside_touches": checks.Outside})
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
