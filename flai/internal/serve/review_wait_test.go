//go:build !windows

package serve

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0296, I-0088, ADR-0093, ADR-0096: a story in review holds nothing by
// overlap, but a story naming it in after: and, while review is at its limit,
// every ready story still wait for its acceptance. With accept_reviews off
// only the operator accepts it, so the board stands idle until a person
// does. With it on, the orchestrator accepts it through the path flai accept
// --by orchestrator takes, and the launcher's next look starts both waiting
// stories' agents without a person.
func TestTheOrchestratorsAcceptanceEndsTheBoardsWaitOnTheOperator(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limits(3, 1)
	lab.hold()
	reviewed := lab.backlog("Reviewed", nil) // touches docs/reviewed
	then := lab.backlog("Then", nil)
	st, err := lab.repo.Get(then)
	if err != nil {
		t.Fatal(err)
	}
	st.After = []string{reviewed}
	if err := lab.repo.Save(st); err != nil {
		t.Fatal(err)
	}
	other := lab.backlog("Other", nil)
	head := lab.inReviewVerified(reviewed, "docs/reviewed/guide.md")
	lab.toReady(then)
	lab.toReady(other)
	evidence := "Verdict: pass, tests and lint clean\n- 1: `docs/reviewed/guide.md`"

	// accept_reviews off: the orchestrator's acceptance is refused, and
	// nothing starts while the story waits for a person
	if err := lab.acceptAsOrchestrator(reviewed, head, evidence); err == nil || !strings.Contains(err.Error(), "orchestration.permissions.accept_reviews, which is off") {
		t.Fatalf("the orchestrator accepted without accept_reviews: %v", err)
	}
	if it, _ := lab.repo.Get(reviewed); it.Status != workitem.Review {
		t.Fatalf("a refused acceptance moved %s to %s", reviewed, it.Status)
	}
	lab.l.look(ctx, false)
	if lab.run(then) != nil || lab.run(other) != nil || len(lab.entries()) != 0 {
		t.Fatalf("started a story waiting on the one in review: %+v", lab.state().Stories)
	}
	after := then + " held (after): waits for " + reviewed + " (in review); starts when " + reviewed + " is done"
	full := "review is full (1 of 1): accept or send back a story; it holds " + other
	if w := lab.state().Waiting; !strings.Contains(w, after) || !strings.Contains(w, full) {
		t.Errorf("waiting: %q", w)
	}
	if a := Activity(lab.root, lab.state())[then]; a.Hold == nil || a.Hold.Code != workitem.HoldAfter {
		t.Errorf("activity of %s: %+v", then, a)
	}

	// the operator turns accept_reviews on: the orchestrator accepts, and
	// the next look starts both waiting stories
	lab.permit("orchestration:\n  permissions:\n    accept_reviews: true\n")
	if err := lab.acceptAsOrchestrator(reviewed, head, evidence); err != nil {
		t.Fatalf("the orchestrator's acceptance: %v", err)
	}
	all, err := lab.repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	accepted := false
	for _, it := range all {
		if it.ID == reviewed {
			last := it.Transitions[len(it.Transitions)-1]
			accepted = it.Archived && it.Status == workitem.Done && last.By == workitem.ActivityOrchestrator
		}
	}
	if !accepted {
		t.Errorf("%s is not done by the orchestrator and archived", reviewed)
	}
	lab.l.look(ctx, false)
	waitFor(t, "both waiting stories' agents run", func() bool { return lab.run(then).live() && lab.run(other).live() })
	if w := lab.state().Waiting; w != "" {
		t.Errorf("still waiting: %q", w)
	}
	lab.release(then)
	lab.release(other)
	waitFor(t, "both end", func() bool { return !lab.run(then).live() && !lab.run(other).live() })
}

// git runs git in dir with the real runner and is its output, trimmed.
func (lab *agentLab) git(dir string, args ...string) string {
	lab.t.Helper()
	out, err := execx.System{}.Run(dir, "git", args...)
	if err != nil {
		lab.t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
	return strings.TrimSpace(out)
}

// inReviewVerified makes the lab a git repository and puts story id in
// review, every criterion ticked and its one task done, with a branch whose
// one commit adds file, and is that branch's head, which its verifier passed.
func (lab *agentLab) inReviewVerified(id, file string) string {
	lab.t.Helper()
	_ = os.WriteFile(filepath.Join(lab.root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	lab.git(lab.root, "init", "-q", "-b", "main")
	lab.git(lab.root, "config", "user.email", "t@t")
	lab.git(lab.root, "config", "user.name", "t")
	lab.git(lab.root, "config", "commit.gpgsign", "false")
	lab.git(lab.root, "add", "-A")
	lab.git(lab.root, "commit", "-q", "-m", "init")
	wt := lab.repo.WorktreePath(id)
	lab.git(lab.root, "worktree", "add", "-q", "-b", "story/"+id, wt)
	_ = os.MkdirAll(filepath.Dir(filepath.Join(wt, file)), 0o755)
	if err := os.WriteFile(filepath.Join(wt, file), []byte("# Guide\n"), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.git(wt, "add", "-A")
	lab.git(wt, "commit", "-q", "-m", "docs: ["+id+"] the guide")

	lab.toReady(id)
	lab.move(id, workitem.InProgress)
	st, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	st.Body = strings.Replace(st.Body, "- [ ] works", "- [x] works", 1)
	if err := lab.repo.Save(st); err != nil {
		lab.t.Fatal(err)
	}
	lab.move(id, workitem.Review)
	items, _ := lab.repo.List(false)
	for _, task := range workitem.Children(items, id) {
		for _, to := range preview.StepsToDone(task.Status) {
			if _, err := lab.repo.Transition(task, to, "builder-"+id, "", lab.now); err != nil {
				lab.t.Fatal(err)
			}
		}
	}
	return lab.git(lab.root, "rev-parse", "story/"+id)
}

// permit gives the lab's manifest orchestration, as the operator does, and
// commits it on main.
func (lab *agentLab) permit(orchestration string) {
	lab.t.Helper()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" + orchestration
	if err := os.WriteFile(filepath.Join(lab.root, manifest.File), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.git(lab.root, "add", manifest.File)
	lab.git(lab.root, "commit", "-q", "-m", "chore: orchestration")
}

// acceptAsOrchestrator is the stand-in orchestrator's flai accept id --by
// orchestrator --verified verified --evidence: the project read afresh, the
// permission, the preview with ADR-0093's conditions, then the acceptance as
// flai accept applies it (the branch rebased and merged, the walk to done by
// the orchestrator, the archive). The commit, the epic's walk, and the
// overlap notices are left out: the launcher reads none of them.
func (lab *agentLab) acceptAsOrchestrator(id, verified, evidence string) error {
	lab.t.Helper()
	repo, err := workitem.Open(lab.root)
	if err != nil {
		return err
	}
	it, err := repo.Get(id)
	if err != nil {
		return err
	}
	if err := repo.OrchestratorPermits(manifest.PermitAcceptReviews, "accepts a story", id); err != nil {
		return err
	}
	ev, err := preview.ParseEvidence(evidence)
	if err != nil {
		return err
	}
	r, by := execx.System{}, workitem.ActivityOrchestrator
	res, err := preview.AcceptWith(r, repo, it, by, lab.now, slog.New(slog.DiscardHandler), preview.AcceptOptions{Orchestrator: true, Verified: verified, Evidence: ev})
	if err != nil {
		return err
	}
	if len(res.Blockers) > 0 || len(res.Uncommitted) > 0 {
		return fmt.Errorf("%s cannot be accepted yet: %s; uncommitted %v", id, strings.Join(res.Blockers, "; "), res.Uncommitted)
	}
	wt, branch := repo.WorktreePath(id), "story/"+id
	lab.git(wt, "rebase", "-q", "main")
	lab.git(lab.root, "merge", "--ff-only", "--quiet", branch)
	lab.git(lab.root, "worktree", "remove", wt)
	lab.git(lab.root, "branch", "-d", branch)
	items, err := repo.List(false)
	if err != nil {
		return err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return err
	}
	for _, to := range preview.StepsToDone(it.Status) {
		if _, err := repo.Move(it, to, workitem.MoveOptions{By: by, Now: lab.now, Items: items, Board: board}); err != nil {
			return err
		}
	}
	if err := repo.Save(it); err != nil {
		return err
	}
	items, err = repo.List(false)
	if err != nil {
		return err
	}
	plan, err := repo.PlanArchive(items, preview.ArchivedWith(items, id, res.Epic))
	if err != nil {
		return err
	}
	return repo.Archive(plan)
}
