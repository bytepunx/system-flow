package serve

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/sendback"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0312 §4: a story a state job sent back gets its agent started once, in
// its last session, on the send back's thread, told the state and what
// stopped the job, with no words of an acceptance's rebase.
func TestAStoryAStateJobSentBackGetsItsAgentToldTheState(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.cfg.Command = nil
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
	lab.hold()
	id := lab.readyWith("Sent back by its state job", &manifest.Agent{Harness: harness.ClaudeCode})
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	lab.move(id, workitem.Review)
	lab.release(id)
	waitFor(t, "it ends in review", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeWorked })
	lab.endedAgo(id, 3*time.Minute)
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))

	st, _ := lab.repo.Get(id)
	if _, err := lab.repo.Transition(st, workitem.InProgress, storygit.SentBackAuthor, "state job: suite smoke failed", lab.now); err != nil {
		t.Fatal(err)
	}
	th, err := sendback.Thread(lab.repo, id, sendback.StateJobTitle(id, workitem.Review, "suite smoke failed"), "The suite smoke failed.", time.Now().Add(-2*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs for the send back", func() bool { r := lab.run(id); return r.live() && r.SentBack == th.ID })
	if again := lab.run(id); again.Agent != first.Agent || again.Session != first.Session || again.AutoRestarts != 0 {
		t.Errorf("the same agent, in its session, spending no restart: %+v, first %+v", again, first)
	}
	out := filepath.Join(lab.outDir, id+".txt")
	waitFor(t, "the stub was told the thread", func() bool { got, _ := os.ReadFile(out); return strings.Contains(string(got), th.ID) })
	got, _ := os.ReadFile(out)
	for _, want := range []string{
		"on thread " + th.ID + ": flai serve's job on entering review found that suite smoke failed, and sent the story back",
		"once " + id + " is back in review the job runs again",
		"--resume " + first.Session,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the agent was given\n%s\nwithout %q", got, want)
		}
	}
	if strings.Contains(string(got), "the acceptance of") {
		t.Errorf("the agent of a story its state job sent back was told of an acceptance:\n%s", got)
	}
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
}

// ADR-0312 §4: a state job's send back's thread is resolved once the story
// is back in the state the job ran for, entered after the thread, and not
// when it is in another state of the waiting kind, as review, which resolves
// an acceptance's.
func TestAStateJobsSendBackIsResolvedOnceTheStoryIsBackInItsState(t *testing.T) {
	lab := newAgentLab(t)
	w, problem := workflow.Parse("design/workflows/board.md", statesFlow)
	if problem != nil {
		t.Fatalf("the board: %v", problem)
	}
	lab.repo.Workflow = &w
	id := lab.backlog("Sent back from integration", nil)
	at := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	moveAt := func(to string, n int) {
		t.Helper()
		st, err := lab.repo.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := lab.repo.Transition(st, to, "builder-"+id, "the test moves it", at.Add(time.Duration(n)*time.Minute)); err != nil {
			t.Fatal(err)
		}
	}
	for n, to := range []string{workitem.Ready, workitem.InProgress, workitem.Review, integration, workitem.InProgress} {
		moveAt(to, n)
	}
	th, err := sendback.Thread(lab.repo, id, sendback.StateJobTitle(id, integration, "suite smoke failed"), "The suite smoke failed.", at.Add(4*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	open := func() bool {
		t.Helper()
		got, err := threads.Get(lab.repo, th.ID)
		if err != nil {
			t.Fatal(err)
		}
		return got.Open()
	}
	moveAt(workitem.Review, 5)
	lab.l.resolveAcceptSentBack(lab.repo)
	if !open() {
		t.Fatal("the state job's thread was resolved with the story in review, not back in integration")
	}
	moveAt(integration, 6)
	lab.l.resolveAcceptSentBack(lab.repo)
	got, err := threads.Get(lab.repo, th.ID)
	if err != nil {
		t.Fatal(err)
	}
	entries := got.Entries()
	if got.Open() || entries[len(entries)-1].Author != forgeAuthor || !strings.Contains(entries[len(entries)-1].Text, id+" is back in "+integration+" since") {
		t.Errorf("the thread once the story is back in %s: %s %+v", integration, got.Status, entries)
	}
}
