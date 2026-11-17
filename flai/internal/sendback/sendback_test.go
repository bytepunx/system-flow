package sendback

import (
	"os"
	"path/filepath"
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
