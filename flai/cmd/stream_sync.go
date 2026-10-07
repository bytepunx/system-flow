package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai stream sync runs storygit.Sync, the rebase and the checks after it
// (ADR-0069, S-0131, ADR-0046), and prints what it answers.

// syncStoppedFrom is the stop a sync's result describes, for the report it
// prints and the error it exits with.
func syncStoppedFrom(res storygit.SyncResult) *syncStopped {
	return &syncStopped{
		Story: res.Story, Base: res.Base, Worktree: res.Worktree, Uncommitted: res.Uncommitted,
		RebaseInProgress: res.RebaseInProgress(), Found: res.Stopped == storygit.StopRebaseInProgress,
		Conflicts: res.Conflicts, Continue: res.Continue, Abort: res.Abort,
	}
}

// printSyncChecks says, one line each, how story's branch merges with the
// other open branches, and what it changed outside the story's claim.
func printSyncChecks(w io.Writer, story *workitem.Item, res storygit.SyncResult) {
	mine := storyBranch(story.ID)
	if res.TrialMergeSkipped != "" {
		fmt.Fprintf(w, "no trial merge with the other open story branches: %s\n", res.TrialMergeSkipped)
	}
	for _, b := range res.Branches {
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
	if len(res.Outside) > 0 {
		fmt.Fprintf(w, "%s changed %s outside %s's touches: %s\n", mine, plural(len(res.Outside), "path"), story.ID, strings.Join(res.Outside, ", "))
		fmt.Fprintf(w, "widen them so that stories that overlap wait: flai touches %s --add %s\n", story.ID, strings.Join(res.Outside, " "))
	}
}
