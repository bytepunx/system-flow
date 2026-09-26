//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// worktreeGit answers git status in a story's worktree with what it holds.
type worktreeGit struct {
	execx.System
	status string
}

func (g *worktreeGit) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && len(args) > 0 && args[0] == "status" {
		return g.status, nil
	}
	return "", nil
}

// S-0140: the operator has an agent started to commit what a story in
// review left uncommitted in its worktree, and is told why when it will not.
func TestAnAgentIsStartedToCommitWhatAStoryInReviewLeftUncommitted(t *testing.T) {
	ctx := context.Background()
	lab := newAgentLab(t)
	lab.hold()
	git := &worktreeGit{}
	lab.o.Git = git
	commit := func(id string) (*AgentRun, error) {
		return Commit(ctx, lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	}
	refusedFor := func(t *testing.T, id, want string) {
		t.Helper()
		_, err := commit(id)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("commit %s: %v, want refused for %q", id, err, want)
		}
	}

	ready := lab.ready("Ready")
	refusedFor(t, ready, "is in ready; only a story in review")
	id := lab.ready("Left over")
	lab.move(id, workitem.InProgress)
	refusedFor(t, id, "is in in-progress; only a story in review")
	lab.move(id, workitem.Review)
	refusedFor(t, id, "has no worktree at .flai-cache/worktrees/"+id)
	if err := os.MkdirAll(lab.repo.WorktreePath(id), 0o755); err != nil {
		t.Fatal(err)
	}
	refusedFor(t, id, "has nothing uncommitted")
	git.status = "?? docs/left.md"
	lab.cfg.Enabled = false
	refusedFor(t, id, "agent host action is off")
	lab.cfg.Enabled = true

	run, err := commit(id)
	if err != nil {
		t.Fatal(err)
	}
	if run.Agent != "builder-"+id || run.PID == 0 {
		t.Errorf("started: %+v", run)
	}
	if got := lab.run(id); got == nil || got.PID != run.PID {
		t.Errorf("the run is not recorded where flai serve tracks it: %+v", got)
	}
	j := lab.entries()
	if last := j[len(j)-1]; !strings.Contains(last.Detail, "for "+id+" as builder-"+id+" to commit what its worktree holds") {
		t.Errorf("journal: %+v", last)
	}
	refusedFor(t, id, "agent is running")
	lab.release(id)
	// The serving flai judges the run by the story's state once it ends: a
	// story still in review worked.
	if outcome, why, _ := judge(lab.root, id, run.Agent, new(int)); outcome != OutcomeWorked {
		t.Errorf("a commit run that leaves the story in review: %s %s", outcome, why)
	}
}
