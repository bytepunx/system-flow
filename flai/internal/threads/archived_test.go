package threads

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// onItems builds a story with a narrative and a task, another story, and
// threads on them: open on the story, answered on the task, resolved on the
// story, and open on the other story.
func onItems(t *testing.T) (r *workitem.Repo, story, task *workitem.Item) {
	t.Helper()
	r = project(t)
	epic, err := r.Create(workitem.NewOptions{Type: workitem.Epic, Title: "E", Owner: "a", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	story, _ = r.Create(workitem.NewOptions{Type: workitem.Story, Title: "S", Parent: epic.ID, Owner: "a", Now: t0})
	task, _ = r.Create(workitem.NewOptions{Type: workitem.Task, Title: "T", Parent: story.ID, Owner: "a", Now: t0})
	other, _ := r.Create(workitem.NewOptions{Type: workitem.Story, Title: "O", Parent: epic.ID, Owner: "a", Now: t0})
	if _, err := r.OpenStream(story, workitem.StreamOptions{Agent: "claude", Now: t0}); err != nil {
		t.Fatal(err)
	}
	for _, o := range []NewOptions{
		{Title: "Open on the story", On: story.ID},
		{Title: "Answered on the task", On: task.ID},
		{Title: "Resolved on the story", On: story.ID},
		{Title: "Open on another story", On: other.ID},
	} {
		o.Author, o.Text, o.Now = "alex", "?", t0
		if _, err := New(r, o); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Reply(r, "TH-0002", "claude", "Yes.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(r, "TH-0003", "alex", "settled", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := MirrorNarrative(r, story.ID); err != nil {
		t.Fatal(err)
	}
	return r, story, task
}

func ids(list []*Thread) string {
	var out []string
	for _, th := range list {
		out = append(out, th.ID)
	}
	return strings.Join(out, " ")
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// I-0073: the threads left open on items being archived are resolved, each
// with an entry giving the reason, and no other thread is touched.
func TestResolveOnItemsResolvesTheOpenThreadsOnTheItems(t *testing.T) {
	r, story, task := onItems(t)
	// the IDs in any padding name the items
	loose := []string{strings.Replace(story.ID, "-000", "-", 1), strings.Replace(task.ID, "-0", "-00", 1)}
	list, err := OnItems(r, loose)
	if err != nil {
		t.Fatal(err)
	}
	if got := ids(list); got != "TH-0001 TH-0002" {
		t.Fatalf("OnItems: got %q, want the open and answered threads in ID order", got)
	}
	resolvedBefore, _ := Get(r, "TH-0003")
	otherBefore, _ := Get(r, "TH-0004")
	resolvedFile, otherFile := readFile(t, resolvedBefore.Path), readFile(t, otherBefore.Path)

	reason := "its item " + story.ID + " is archived"
	got, err := ResolveOnItems(r, loose, "flai", reason, t0.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, " ") != "TH-0001 TH-0002" {
		t.Fatalf("ResolveOnItems returned %v", got)
	}
	for _, id := range got {
		th, err := Get(r, id)
		if err != nil {
			t.Fatal(err)
		}
		es := th.Entries()
		last := es[len(es)-1]
		if th.Status != "resolved" || last.Author != "flai" || last.Text != "Resolved: "+reason || last.At != "2026-09-17T22:00:00Z" {
			t.Errorf("%s: status %s, last entry %+v", id, th.Status, last)
		}
	}
	if readFile(t, resolvedBefore.Path) != resolvedFile {
		t.Error("an already-resolved thread must be left byte for byte")
	}
	if readFile(t, otherBefore.Path) != otherFile {
		t.Error("a thread on another item must be left byte for byte")
	}
	if list, _ := OnItems(r, loose); len(list) != 0 {
		t.Errorf("nothing is left open on the items: %s", ids(list))
	}
	// the story's narrative no longer lists them
	if n := readFile(t, r.NarrativePath(story.ID)); strings.Contains(n, "threads:start") {
		t.Errorf("the mirror must drop the resolved threads:\n%s", n)
	}
}

// Callers resolve right after archiving, when the story's narrative has left
// wip/agents: the archived narrative is not written and none is recreated.
func TestResolveOnItemsAfterTheNarrativeIsArchived(t *testing.T) {
	r, story, task := onItems(t)
	live := r.NarrativePath(story.ID)
	archived := filepath.Join(r.ArchiveDir(), "agents", filepath.Base(live))
	if err := os.MkdirAll(filepath.Dir(archived), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(live, archived); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, archived)
	got, err := ResolveOnItems(r, []string{story.ID, task.ID}, "flai", "archived", t0.Add(time.Hour))
	if err != nil || len(got) != 2 {
		t.Fatalf("ResolveOnItems: %v %v", got, err)
	}
	if readFile(t, archived) != before {
		t.Error("the archived narrative must not be written")
	}
	if _, err := os.Stat(live); !os.IsNotExist(err) {
		t.Errorf("no narrative may be recreated in wip/agents: %v", err)
	}
}

func TestResolveOnItemsWithNoThreads(t *testing.T) {
	r, _, _ := onItems(t)
	got, err := ResolveOnItems(r, []string{"S-0099"}, "flai", "archived", t0)
	if err != nil || len(got) != 0 {
		t.Errorf("no threads on the items is no work: %v %v", got, err)
	}
}
