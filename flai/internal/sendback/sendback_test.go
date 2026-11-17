package sendback

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// newProject is a project on the baseline workflow with one story in it.
func newProject(t *testing.T) (*workitem.Repo, *workitem.Item) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	st, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Send me back", Parent: epic.ID, Owner: "alex", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	return repo, st
}

// On the baseline workflow a story in review is sent back to in-progress.
func TestTargetIsTheBaselinesInProgress(t *testing.T) {
	t.Parallel()
	repo, st := newProject(t)
	st.Status = workitem.Review
	if got := Target(repo.Flow(), st); got != workitem.InProgress {
		t.Errorf("Target = %q, want %q", got, workitem.InProgress)
	}
}

// ADR-0248: the first send back opens the thread by flai on the story, the
// one flai serve starts its agent on; a later one, while it is open, writes
// on it rather than opening another; once it is resolved, a send back opens
// a new one.
func TestThreadOpensOneThreadByFlaiAndWritesOnItWhileItIsOpen(t *testing.T) {
	t.Parallel()
	repo, st := newProject(t)
	title := storygit.SentBackPrefix + "story/" + st.ID + " stopped rebasing onto main"
	// a minute apart, since entries by one author at one minute are one
	at := time.Now()
	first, err := Thread(repo, st.ID, title, "first", at)
	if err != nil {
		t.Fatal(err)
	}
	if first.Opener() != storygit.SentBackAuthor || first.Title != title || first.Anchor.Item != st.ID {
		t.Errorf("the thread: %+v by %s", first, first.Opener())
	}
	second, err := Thread(repo, st.ID, title+" again", "second", at.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if second.ID != first.ID || len(second.Entries()) != 2 || second.Entries()[1].Text != "second" || second.Entries()[1].Author != storygit.SentBackAuthor {
		t.Errorf("the second send back did not write on the open thread: %+v", second.Entries())
	}
	if _, err := threads.Resolve(repo, first.ID, "flai serve", "back in review", at.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	third, err := Thread(repo, st.ID, title, "third", at.Add(3*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if third.ID == first.ID {
		t.Errorf("a send back after the thread was resolved wrote on it: %s", third.ID)
	}
}

// I-0397: a send back by hand's title begins with the prefix flai serve
// finds the send back's thread by, and gives back who made it; the
// acceptance's and the rebase and resolve job's are not one.
func TestByHandTitleSaysWhoSentItBack(t *testing.T) {
	t.Parallel()
	title := ByHandTitle("alex", "S-0823", "review", "in-progress")
	if want := "Acceptance sent back: alex sent S-0823 back from review to in-progress by hand"; title != want {
		t.Errorf("ByHandTitle = %q, want %q", title, want)
	}
	if by, ok := ByHandMover(title); !ok || by != "alex" || !IsByHand(title) {
		t.Errorf("ByHandMover(%q) = %q, %v", title, by, ok)
	}
	for _, other := range []string{
		storygit.SentBackTitle("story/S-0823", "main", []string{"a.go"}),
		storygit.SentBackPrefix + "story/S-0823 could not be rebased onto main by rebase and resolve",
		"alex sent S-0823 back from review to in-progress by hand",
	} {
		if IsByHand(other) {
			t.Errorf("IsByHand(%q) is true", other)
		}
	}
}

// I-0397: the operator's move of a story from review to in-progress opens
// the send back's thread by flai with the reason, and a second one while it
// is open writes on it; with no reason the thread says so and has the agent
// ask the mover.
func TestByHandOpensTheSendBacksThreadForTheOperatorsMove(t *testing.T) {
	t.Parallel()
	repo, st := newProject(t)
	at := time.Now()
	th, ok, err := ByHand(repo, repo.Flow(), st, workitem.Review, workitem.InProgress, "alex", "tests missing", at)
	if err != nil || !ok {
		t.Fatalf("ByHand = %v, %v", ok, err)
	}
	if th.Opener() != storygit.SentBackAuthor || th.Anchor.Item != st.ID || !IsByHand(th.Title) || !strings.Contains(th.Entries()[0].Text, "tests missing") {
		t.Errorf("the thread: %q by %s on %s: %q", th.Title, th.Opener(), th.Anchor.Item, th.Entries()[0].Text)
	}
	again, ok, err := ByHand(repo, repo.Flow(), st, workitem.Review, workitem.InProgress, "alex", "", at.Add(time.Minute))
	if err != nil || !ok {
		t.Fatalf("ByHand again = %v, %v", ok, err)
	}
	es := again.Entries()
	if again.ID != th.ID || len(es) != 2 || !strings.Contains(es[1].Text, "gave no reason") || !strings.Contains(es[1].Text, "ask alex on this thread") {
		t.Errorf("the second send back did not write on the open thread: %s %+v", again.ID, es)
	}
	if all, _ := threads.For(repo, st.ID); len(all) != 1 {
		t.Errorf("%d threads on %s, want 1", len(all), st.ID)
	}
}

// I-0397: flai's own send backs open their thread themselves, the story's
// agent is on it already, and a move that is no send back needs no thread.
func TestByHandOpensNoThreadForFlaiTheAgentOrAnotherMove(t *testing.T) {
	t.Parallel()
	repo, st := newProject(t)
	for _, c := range []struct{ from, to, by string }{
		{workitem.Review, workitem.InProgress, storygit.SentBackAuthor},
		{workitem.Review, workitem.InProgress, "agent-" + st.ID},
		{workitem.Ready, workitem.InProgress, "alex"},
	} {
		th, ok, err := ByHand(repo, repo.Flow(), st, c.from, c.to, c.by, "why", time.Now())
		if err != nil || ok || th != nil {
			t.Errorf("ByHand from %s to %s by %s = %v, %v, %v; want no thread", c.from, c.to, c.by, th, ok, err)
		}
	}
	if all, _ := threads.For(repo, st.ID); len(all) != 0 {
		t.Errorf("%d threads on %s, want none", len(all), st.ID)
	}
}
