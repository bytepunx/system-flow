package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newStreamCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "stream",
		Short: "Open and append to agent narratives in wip/agents",
		Long: `A stream is the narrative for one story. Set FLAI_AGENT and FLAI_SESSION
so entries record who wrote them.`,
	}
	c.AddCommand(newStreamDiffCmd(a), newStreamOpenCmd(a), newStreamLogCmd(a), newStreamStateCmd(a), newStreamSyncCmd(a), newStreamAnswerCmd(a))
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
			var opened storygit.Opened
			if !noBranch {
				if opened, err = a.openStoryBranch(repo, story.ID, reopen); err != nil {
					return fmt.Errorf("narrative opened but the branch was not: %w", err)
				}
			}
			branch, wt, from := opened.Branch, opened.Worktree, opened.From
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
	case "", storygit.FromLocal:
		return ""
	case storygit.FromMain:
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

A stop whose only conflicting paths are generated files, which today is
design/issues/summary.md alone, and issue files under design/issues does not
wait for you (ADR-0098, ADR-0126). Sync merges each issue file from the stop's
three versions by its instances, keeping the instances both sides added and
adding up the count, then writes summary.md again from the issue files in the
worktree at that stop, git adds them, and continues the rebase, until the
rebase finishes or stops on another path. An issue file both sides changed
some other way, its title or description, say, is not merged: the stop is
left for you. A stop where other paths conflict too is left for you, and
lists them all, summary.md included; once the others are resolved, flai issue
summary run in the worktree writes it. flai accept syncs the same way. A
rebase already in progress is never continued for you.

Once the rebase is clean, each open issue the branch added whose title an open
issue on the main branch has too is folded into that issue: its instances,
count, and cost move into it, the branch's file is deleted, summary.md is
written again, and the fold is committed on the branch as
"docs: [S-nnnn] fold I-x into I-y". flai accept folds the same way. The report
names each issue file merged and each issue folded; with --json, merged lists
the files and folded each fold's from and into.

After a clean rebase it trial-merges the branch with the branch of every other
story in progress or in review (git merge-tree --write-tree, git 2.38 or
newer), writing nothing to any worktree, and lists each branch it conflicts
with and the conflicting paths. A pair's conflicts are the paths both branches
changed since they left the main branch, not what the main branch brought
since: a path where main has since changed what a branch behind it did is
that story's own rebase to settle (I-0064). Generated files and issue files are
left out, so a pair whose only conflicts are summary.md and issue files counts
as clean. A conflict is told to both
stories in their conversation, not a thread (ADR-0121): flai messages the other
story from the one that synced, about the conflicting paths, and sync names
the pair's conversation. A later sync with the same paths adds nothing, and one
that finds the pair merging cleanly again, or the other story no longer open,
closes the conversation. The two agents agree who changes what; either asks the
operator with flai message escalate only when they do not agree. A conflict
thread an older flai opened is still resolved the same way, and no sync opens
one. With --json each branch names its conversation, and thread names an old
conflict thread it resolved. It also lists the paths the branch changed since
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
			res, err := storygit.Sync(storygit.SyncOptions{Runner: a.runner, Repo: repo, Story: it, Now: a.now(), Files: issues.SyncFiles(repo, it.ID, a.runner, a.now), Log: a.logger()})
			// conflicts is null here unless the rebase waits, as it always was
			if err != nil {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"story": it.ID, "branch": res.Branch, "base": res.Base, "conflicts": nil, "ok": false})
				}
				return err
			}
			if !res.Synced {
				stop := syncStoppedFrom(res)
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"story": it.ID, "branch": res.Branch, "base": res.Base, "worktree": res.Worktree, "ok": false,
						"uncommitted": res.Uncommitted, "conflicts": res.Conflicts, "rebase_in_progress": res.RebaseInProgress(),
						"continue": res.Continue, "abort": res.Abort})
				} else {
					fmt.Fprint(a.out, stop.Report())
				}
				return stop
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"story": it.ID, "branch": res.Branch, "base": res.Base, "conflicts": nil, "ok": true,
					"merged": res.Merged, "folded": res.Folded,
					"branches": res.Branches, "trial_merge_skipped": res.TrialMergeSkipped, "outside_touches": res.Outside})
			}
			fmt.Fprintf(a.out, "%s is rebased onto %s\n", res.Branch, res.Base)
			printSyncChecks(a.out, it, res)
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

func newStreamStateCmd(a *app) *cobra.Command {
	var current, next string
	c := &cobra.Command{
		Use:   "state <story-id> [--current \"<text>\"] [--next \"<text>\"]",
		Short: "Replace a story's narrative's Current state and Next steps",
		Long: `Replaces what is under the story's narrative's ## Current state with the
--current text and what is under ## Next steps with the --next text. Either
may be left out, and its section is left as it was; giving neither is an
error. A value of - reads that text from standard input, so text of several
lines needs no quoting; only one of the two may be -.

Every other section of the narrative is left as it was, and nothing is
appended to its log: write the log with flai stream log. The narrative's
updated stamp, agent, and session are written as flai stream log writes them,
and wip/agents/index.md is written again. The same text again writes nothing.

A story not in progress or in review is refused, as is one with no narrative,
and text the project's markdown lint rejects; each refusal exits 4 and writes
nothing, and with --json prints {"refused": {...}}. A text holding a # or ##
heading, which would end its section, is an error.`,
		Example: `  flai stream state S-0271 --current "T-1055 is done; T-1056 is next."
  flai stream state S-0271 --next - < next-steps.md
  flai stream state S-0271 --current "Reviewing." --next "1. Answer the review." --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			given := cmd.Flags().Changed("current") || cmd.Flags().Changed("next")
			if !given {
				return fmt.Errorf("give --current, --next, or both: the text to write under ## Current state and ## Next steps")
			}
			if current == "-" && next == "-" {
				return fmt.Errorf("only one of --current and --next can be read from standard input; give the other as text")
			}
			for _, text := range []*string{&current, &next} {
				if *text == "-" {
					data, err := io.ReadAll(cmd.InOrStdin())
					if err != nil {
						return err
					}
					*text = string(data)
				}
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			agent, session := agentIdentity()
			res, err := repo.SetStreamState(args[0], current, next, workitem.StreamOptions{Agent: agent, Session: session, Now: a.now()})
			if refused := streamRefusal(repo.MainRoot, err); refused != nil {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": refused})
				}
				return &exitError{code: exitDocRefused, msg: err.Error()}
			}
			if err != nil {
				return err
			}
			res.Path = relPath(repo.MainRoot, res.Path)
			if a.jsonOut {
				return a.printJSON(res)
			}
			sections := strings.Join(res.Written, " and ")
			if !res.Changed {
				fmt.Fprintf(a.out, "%s: unchanged; %s already holds the %s given\n", res.Stream, res.Path, sections)
				return nil
			}
			fmt.Fprintf(a.out, "%s: wrote %s in %s at %s\n", res.Stream, sections, res.Path, res.Updated)
			return nil
		},
	}
	c.Flags().StringVar(&current, "current", "", "the text for ## Current state, or - to read it from standard input")
	c.Flags().StringVar(&next, "next", "", "the text for ## Next steps, or - to read it from standard input")
	return c
}

// streamRefusal is what flai stream state prints under "refused" when err is
// a refusal, or nil when it is not: the story's state with the reason, or the
// lint's findings with the narrative's path relative to root.
func streamRefusal(root string, err error) any {
	var refused *workitem.StreamRefusedError
	if errors.As(err, &refused) {
		return refused
	}
	var lint *mdlint.Error
	if errors.As(err, &lint) {
		return map[string]any{"path": relPath(root, lint.Path), "reason": lint.Error(), "findings": lint.Findings}
	}
	return nil
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
