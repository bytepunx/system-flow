package serve

import (
	"context"
	"os"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Commit starts a story's agent, in a new session, to commit what its
// worktree holds and nothing else, on the operator's word (S-0140): flai
// serve agent commit, or the story page's button when the acceptance preview
// finds the worktree uncommitted. It is for a story in review: its agent has
// ended, and the work it left uncommitted stops the acceptance. It is refused
// while the agent action is off for the project, when the story is not in
// review, while its agent runs or waits for an answer, when it has no
// worktree or the worktree holds nothing uncommitted, and when nothing can
// start it. The in-progress limit does not apply: the story is past it.
func Commit(ctx context.Context, o Options, e Entry, story string) (*AgentRun, error) {
	cfg := o.Agent(e.Root)
	if !cfg.Enabled {
		return nil, refused("the agent host action is off for this project: flai serve enable agent")
	}
	if !StoryID.MatchString(story) {
		return nil, refused("%q is not a story's ID", story)
	}
	repo, err := workitem.Open(e.Root)
	if err != nil {
		return nil, err
	}
	it, err := repo.Get(story)
	if err != nil {
		return nil, err
	}
	if it.Type != workitem.Story {
		return nil, refused("%s is %s; only a story has an agent", it.ID, it.Type)
	}
	if it.Status != workitem.Review {
		return nil, refused("%s is in %s; only a story in review has its worktree committed this way, and one in progress commits its own", it.ID, it.Status)
	}
	switch run := o.Dir.AgentStates()[e.Root].Stories[it.ID]; {
	case run.running():
		return nil, refused("%s's agent is running (pid %d, started %s)", it.ID, run.PID, run.Started)
	case run != nil && run.Outcome == OutcomeAsked:
		return nil, refused("%s's agent is waiting for an answer to %s; answering it starts the agent again", it.ID, run.Thread)
	}
	wt := repo.WorktreePath(it.ID)
	if _, err := os.Stat(wt); err != nil {
		return nil, refused("%s has no worktree at %s, so there is nothing to commit", it.ID, repo.RelToMain(wt))
	}
	repo.Git = o.Git
	if repo.Git == nil {
		repo.Git = execx.System{}
	}
	dirty, err := repo.Uncommitted(it.ID)
	if err != nil {
		return nil, err
	}
	if len(dirty) == 0 {
		return nil, refused("the worktree %s has nothing uncommitted", repo.RelToMain(wt))
	}
	if (it.Agent == nil || it.Agent.Harness == "") && cfg.host(harness.Command).Program == "" {
		return nil, refused("%s names no harness, and no command is set on the host (flai serve agent set -- <program> [args...])", it.ID)
	}
	return startNow(ctx, o, e, readyStory{ID: it.ID, Commit: wt})
}
