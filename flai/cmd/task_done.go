package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/taskdone"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai task done closes a task in one call through taskdone.Run (ADR-0107,
// ADR-0128) and prints what each step did.

func newTaskDoneCmd(a *app) *cobra.Command {
	var message, logEntry string
	c := &cobra.Command{
		Use:   "done <task> [-m \"<message>\"]",
		Short: "Close a task in one call: commit, tell overlapping stories, sync, move to done, log, widen touches, check, and read the inbox",
		Long: `Close a task in one call (ADR-0107). flai finds the task's story and works
in the story's worktree. The steps run in this order, and the first that
fails stops the run; the steps after it are not done.

1. Commit (ADR-0128): of the paths changed in the story's worktree, those the
   task's touches cover and those no other open task of the story covers are
   committed with -m, and nothing else in the worktree is staged or
   committed. A path only another open task covers is left uncommitted for
   that task's close and listed as left; it is not a failure. -m is needed
   only when there is something to commit; nothing to commit is not a
   failure.
2. Tell: a message to each other story in progress or in review whose claim
   covers a path the commit changed, the shared paths included, about those
   paths, naming the task, the commit, and its subject; on the conversation
   open between the two stories, or a new one. A message that cannot be sent
   is logged and stops nothing.
3. Sync: flai stream sync for the story. It rebases the branch onto the main
   branch, trial-merges it with the other open story branches, and lists
   what it changed outside the story's touches. A refusal or a stop on
   conflicts stops the run, save a refusal for the paths left alone: the
   run goes on, and the close that leaves none syncs the branch.
4. Move: the task to done, under flai move's rules, with any story or epic
   that follows it. A task already done is not moved again, so the call can
   be repeated after a stop.
5. Log: an entry in the story's narrative: --log when given, else the
   message's subject line, else "Closed T-nnnn: <title>".
6. Touches: the paths the commit changed that the task's touches, or the
   story's, do not cover are added to each, as flai touches records them.
   A path left is never added.
7. Check: flai check --strict scoped to the story. A finding in the story
   stops the run; a finding outside it is a note.
8. Inbox: the agent's inbox, as the MCP tool inbox answers it.

The agent is FLAI_AGENT and the session FLAI_SESSION.

While a rebase is unfinished in the worktree, the commit and the sync refuse.
Resolve each conflicting path there, git add it, run git rebase --continue,
and call flai task done again. git rebase --abort undoes the rebase instead.

Running the task's tests with flai test, fixing what they find, and ticking
criteria with flai criteria tick stay the agent's. Close a fix by calling
flai task done again: it commits the fix, syncs, and checks.

The output is a line for each step that ran, with the paths left after the
commit's, then the inbox. --json prints the result as the MCP tool
task_done answers it, the paths left as left.

Exit codes: 0 when every step ran, 3 when the sync stopped the run, 4 when
the check did, and 1 when another step did or the run could not start.

A story or an epic is refused: a story goes to review with flai move
S-nnnn review after its close-out, and an epic follows its stories.`,
		Example: `  flai task done T-0021 -m "feat: [S-0004] T-0021 the parser reads tables"
  flai task done T-0021 -m "fix: [S-0004] what the tests found" --log "fixed the empty table case"
  flai task done T-0021 -m "docs: [S-0004] the guide" --json
  flai task done T-0021 --log "nothing left to commit"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := refuseNotATask(args[0]); err != nil {
				return err
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			agent, session := agentIdentity()
			res, err := taskdone.Run(cmd.Context(), taskdone.Options{
				Repo: repo, Runner: a.runner, Task: args[0], Message: message, LogEntry: logEntry,
				Agent: agent, Session: session, Now: a.now, Version: buildinfo.Version, Logger: a.logger(),
			})
			if err != nil {
				return err
			}
			if a.jsonOut {
				if err := a.printJSON(res); err != nil {
					return err
				}
			} else {
				a.printTaskDone(res)
			}
			if code := res.ExitCode(); code != 0 {
				return &exitError{code: code, msg: fmt.Sprintf("%s stopped at the %s step: %s", res.Task, res.Stopped, res.Error)}
			}
			return nil
		},
	}
	c.Flags().StringVarP(&message, "message", "m", "", "the commit message, needed only when there is something to commit; its subject line is the log entry unless --log is given")
	c.Flags().StringVar(&logEntry, "log", "", "the narrative log entry, in place of the message's subject line")
	return c
}

// refuseNotATask refuses a story's or an epic's ID with the command that
// moves it on instead.
func refuseNotATask(id string) error {
	id = workitem.CanonicalID(id)
	switch workitem.TypeOfID(id) {
	case workitem.Story:
		return fmt.Errorf("%s is a story; flai task done closes a task (T-nnnn). Close each of its tasks with flai task done, then move the story with flai move %s review after its close-out", id, id)
	case workitem.Epic:
		return fmt.Errorf("%s is an epic; flai task done closes a task (T-nnnn). An epic follows its stories: close their tasks with flai task done", id)
	}
	return nil
}

// printTaskDone prints a line for each step res ran, the step it stopped at
// with why, and the inbox.
func (a *app) printTaskDone(res taskdone.Result) {
	w := a.out
	if c := res.Commit; c != nil {
		fmt.Fprintf(w, "commit: %s %s\n", short(c.Hash), c.Subject)
	} else if res.Stopped != taskdone.StepCommit {
		fmt.Fprintln(w, "commit: nothing to commit")
	}
	if len(res.Left) > 0 {
		fmt.Fprintf(w, "left for the other open tasks: %s\n", strings.Join(res.Left, ", "))
	}
	for _, t := range res.Told {
		fmt.Fprintf(w, "told: %s (%s) of %s\n", t.Story, t.Conversation, strings.Join(t.Paths, ", "))
	}
	if s := res.Sync; s != nil {
		switch {
		case s.Synced:
			fmt.Fprintf(w, "sync: %s is rebased onto %s\n", s.Branch, s.Base)
			printSyncChecks(w, &workitem.Item{ID: res.Story}, *s)
		case s.Stopped == storygit.StopUncommitted && res.Stopped != taskdone.StepSync:
			// refused for the paths left alone, which stops nothing (ADR-0128)
			fmt.Fprintf(w, "sync: %s was not synced: %s\n", s.Branch, s.Continue)
		case s.Stopped != "":
			fmt.Fprint(w, "sync: "+syncStoppedFrom(*s).Report())
		}
	}
	if m := res.Move; m != nil {
		if m.Skipped {
			fmt.Fprintf(w, "move: %s was done already; not moved again\n", m.ID)
		} else {
			fmt.Fprintf(w, "move: %s → %s\n", m.ID, m.State)
		}
		for i := range m.Followed {
			a.printFollowed(&m.Followed[i])
		}
		for _, warn := range m.Warnings {
			fmt.Fprintf(w, "  warning: %s\n", warn)
		}
	}
	if l := res.Log; l != nil {
		fmt.Fprintf(w, "log: %s at %s: %s\n", l.Stream, l.At, l.Entry)
	}
	if t := res.Touches; t != nil {
		if len(t.Task)+len(t.Story) == 0 {
			fmt.Fprintln(w, "touches: nothing to add")
		}
		if len(t.Task) > 0 {
			fmt.Fprintf(w, "touches: added to %s: %s\n", res.Task, strings.Join(t.Task, ", "))
		}
		if len(t.Story) > 0 {
			fmt.Fprintf(w, "touches: added to %s: %s\n", res.Story, strings.Join(t.Story, ", "))
		}
		printOverlapping(w, t.Overlaps)
	}
	if c := res.Check; c != nil {
		printTaskDoneCheck(w, res.Story, c)
	}
	if res.Stopped != "" {
		fmt.Fprintf(w, "stopped at %s: %s\n", res.Stopped, res.Error)
	}
	if in := res.Inbox; in != nil {
		printTaskDoneInbox(w, in)
	}
}

// printTaskDoneCheck says whether the check passed, then each finding in
// the story and each note outside it, as flai check --story prints them.
func printTaskDoneCheck(w io.Writer, story string, c *taskdone.Check) {
	verdict := "passed"
	if !c.Passed {
		verdict = "failed with " + plural(len(c.Findings), "finding") + " in " + story
	}
	if len(c.Notes) > 0 {
		verdict += fmt.Sprintf("; %s outside %s, which do not fail it", plural(len(c.Notes), "note"), story)
	}
	fmt.Fprintf(w, "check: %s\n", verdict)
	line := func(f check.Finding, suffix string) {
		fmt.Fprintf(w, "  %s:%d: %s: %s: %s%s\n", f.Path, f.Line, f.Level, f.Rule, f.Message, suffix)
	}
	for _, f := range c.Findings {
		line(f, "")
	}
	for _, f := range c.Notes {
		line(f, " (outside "+story+")")
	}
}

// printTaskDoneInbox summarises the inbox: the threads awaiting the agent,
// the stories ready to pull in pull order, and what others changed.
func printTaskDoneInbox(w io.Writer, in *inbox.Inbox) {
	fmt.Fprintf(w, "inbox: %s awaiting %s\n", plural(in.AwaitingYou, "thread"), in.Agent)
	for _, t := range in.Threads {
		fmt.Fprintf(w, "  %s %s (awaiting %s, last by %s)\n", t.ID, t.Title, t.Awaiting, t.LastBy)
	}
	pull := "can pull: yes"
	if !in.CanPull {
		pull = "can pull: no"
		if in.PullHold != "" {
			pull += ", " + in.PullHold
		}
	}
	fmt.Fprintf(w, "inbox: %d ready to pull, in pull order (%s)\n", len(in.Ready), pull)
	for _, r := range in.Ready {
		fmt.Fprintf(w, "  %s %s\n", r.ID, r.Title)
	}
	changes := plural(len(in.Changes), "change")
	if in.Omitted > 0 {
		changes += fmt.Sprintf(" (%d older left out)", in.Omitted)
	}
	fmt.Fprintf(w, "inbox: %s since you last looked\n", changes)
	for _, e := range in.Changes {
		fmt.Fprintf(w, "  %s\n", e.Summary)
	}
	if o := in.FlaiOutdated; o != nil {
		fmt.Fprintf(w, "inbox: %s\n", o.Message)
	}
}
