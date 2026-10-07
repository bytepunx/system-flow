package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/storystart"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai story start starts a ready story in one call through storystart.Start
// (S-0274) and prints what each step did.

// exitStoryStartStep is flai story start's exit code when the story was
// moved to in-progress and a step after the move failed.
const exitStoryStartStep = 3

func newStoryStartCmd(a *app) *cobra.Command {
	var budget string
	c := &cobra.Command{
		Use:   "start <story>",
		Short: "Start a ready story in one call: move it to in-progress, open its stream, prime, and read the inbox",
		Long: `Start a ready story in one call (S-0274), in place of the four calls an
agent makes at the start of a story:
flai move S-nnnn in-progress, then flai stream open S-nnnn,
flai prime --story S-nnnn, and the MCP tool inbox.
The steps run in this order:

1. Move: the story to in-progress under flai move's rules, its epic
   following it to in-progress with its first started story.
2. Stream: the story's narrative is opened, or taken up when it was begun
   before, and in a git repository the branch story/S-nnnn is checked out in
   its worktree under .flai-cache/worktrees/S-nnnn, as flai stream open does.
3. Prime: the story's context pack, fitted to --budget, as flai prime --story
   prints it.
4. Inbox: the agent's inbox, as the MCP tool inbox answers it.

A story that is not ready, or that the board holds (its claim overlaps a
story in progress, or it names in after: a story that is not done), is
refused and nothing changes. The agent is FLAI_AGENT and the session
FLAI_SESSION; FLAI_AGENT records the move and names whose inbox is read.

The output is the move, the narrative, the branch and its worktree, the
pack, and then the inbox: the threads awaiting you, the stories ready to
pull, and what others changed. --json prints one object with the keys story,
followed, worktree, branch, from, pack, and inbox, as the MCP tool
story_start answers it.

Exit codes: 0 when the story started, 4 when it was refused and nothing
changed, 3 when it moved to in-progress and a later step failed (what ran is
printed, and the error names the command that finishes the step), and 1 when
the start could not run.`,
		Example: `  flai story start S-0274
  flai story start S-0274 --budget 120KB
  flai story start S-0274 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			agent, session := agentIdentity()
			res, err := storystart.Start(cmd.Context(), storystart.Options{
				Repo: repo, Runner: a.runner, Story: args[0], Agent: agent, Session: session,
				Budget: budget, Now: a.now, Version: buildinfo.Version, RelativePaths: a.relativeWorktrees, Log: a.logger(),
			})
			var refused *storystart.Refused
			if errors.As(err, &refused) {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": map[string]any{
						"story": refused.Story, "status": refused.Status, "hold": refused.Hold, "reason": refused.Reason,
					}})
				}
				return &exitError{code: exitDocRefused, msg: refused.Error()}
			}
			var step *storystart.StepError
			if err != nil && !errors.As(err, &step) {
				return err
			}
			if a.jsonOut {
				if err := a.printJSON(res); err != nil {
					return err
				}
			} else if err := a.printStoryStart(repo, res, step); err != nil {
				return err
			}
			if step != nil {
				return &exitError{code: exitStoryStartStep, msg: step.Error()}
			}
			return nil
		},
	}
	c.Flags().StringVar(&budget, "budget", "", "the size the context pack fits, such as 80KB (default: prime.budget in system-flow.yaml, else 80KB)")
	return c
}

// printStoryStart prints what a start did: the move, with the epic's when it
// followed, the narrative and the branch as flai stream open prints them,
// the pack as flai prime --story prints it, and the inbox; after a failed
// step, only what the steps before it did.
func (a *app) printStoryStart(repo *workitem.Repo, res storystart.Result, step *storystart.StepError) error {
	w := a.out
	fmt.Fprintf(w, "%s → %s\n", res.Story.ID, res.Story.Status)
	a.printFollowed(res.Followed)
	for _, warn := range res.Story.Warnings {
		fmt.Fprintf(w, "  warning: %s\n", warn)
	}
	if step != nil && step.Step == storystart.StepStream {
		return nil
	}
	fmt.Fprintf(w, "narrative %s\n", relPath(repo.MainRoot, repo.NarrativePath(res.Story.ID)))
	if res.Branch != "" {
		fmt.Fprintf(w, "branch %s%s checked out at %s; work there and run flai stream sync %s at each task transition\n", res.Branch, fromSays(res.From), relPath(repo.MainRoot, res.Worktree), res.Story.ID)
	} else {
		fmt.Fprintln(w, "no branch or worktree: the project is not a git work tree; work in its checkout")
	}
	if res.Pack != nil {
		fmt.Fprintln(w)
		if err := a.printPack(res.Pack, 0); err != nil {
			return err
		}
	}
	if res.Inbox != nil {
		fmt.Fprintln(w)
		printTaskDoneInbox(w, res.Inbox)
	}
	return nil
}
