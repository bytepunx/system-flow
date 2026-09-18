package threads

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 17, 21, 0, 0, 0, time.UTC)

func project(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/system", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte("---\ntitle: Plan\n---\n\n# Plan\n\n## Shape\ntext\n\n## Risks\nmore\n"), 0o644)
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestThreadLifecycle(t *testing.T) {
	r := project(t)
	th, err := New(r, NewOptions{Title: "Is the shape right?", On: "design/system/plan.md", Heading: "shape", Author: "alex", Text: "I wonder.", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if th.ID != "TH-0001" || th.Status != "open" || th.Anchor.Path != "design/system/plan.md" || th.Anchor.Heading != "shape" || th.Opener() != "alex" {
		t.Fatalf("new: %+v", th)
	}
	if _, err := os.Stat(filepath.Join(r.WipDir(), "threads", "TH-0001-is-the-shape-right.md")); err != nil {
		t.Fatal("file name")
	}
	if _, err := New(r, NewOptions{Title: "x", On: "design/system/plan.md", Heading: "nope", Author: "alex", Now: t0}); err == nil || !strings.Contains(err.Error(), "heading") {
		t.Errorf("missing heading should fail: %v", err)
	}
	if _, err := New(r, NewOptions{Title: "x", On: "docs/missing.md", Author: "alex", Now: t0}); err == nil {
		t.Error("missing path should fail")
	}
	if _, err := New(r, NewOptions{Title: "x", On: "S-0099", Author: "alex", Now: t0}); err == nil {
		t.Error("missing item should fail")
	}
	// reply from an agent answers; the opener's follow-up reopens; resolve closes
	th2, err := Reply(r, "th-1", "claude", "Yes, because.", t0.Add(time.Minute))
	if err != nil || th2.Status != "answered" || len(th2.Entries()) != 2 || th2.Entries()[1].Text != "Yes, because." || !contains(th2.Participants, "claude") {
		t.Fatalf("reply: %+v %v", th2, err)
	}
	th3, _ := Reply(r, "TH-0001", "alex", "Then also consider risks.", t0.Add(2*time.Minute))
	if th3.Status != "open" {
		t.Errorf("opener follow-up should reopen: %s", th3.Status)
	}
	th4, err := Resolve(r, "TH-0001", "alex", "settled in ADR-0021", t0.Add(3*time.Minute))
	if err != nil || th4.Status != "resolved" || !strings.Contains(th4.Body, "Resolved: settled in ADR-0021") || th4.Open() {
		t.Fatalf("resolve: %+v %v", th4, err)
	}
	// round trip through the file keeps everything
	back, err := Read(th4.Path)
	if err != nil || back.Marshal() != th4.Marshal() || len(back.Entries()) != 4 {
		t.Fatalf("round trip: %v\n%s", err, back.Marshal())
	}
	if err := back.Validate(); err != nil {
		t.Error(err)
	}
	// reply reopens a resolved thread
	th5, _ := Reply(r, "TH-0001", "claude", "One more thing.", t0.Add(4*time.Minute))
	if th5.Status != "answered" {
		t.Errorf("reply after resolve: %s", th5.Status)
	}
}

func TestItemAnchorsAndNarrativeMirror(t *testing.T) {
	r := project(t)
	epic, _ := r.Create(workitem.NewOptions{Type: workitem.Epic, Title: "E", Owner: "a", Now: t0})
	story, _ := r.Create(workitem.NewOptions{Type: workitem.Story, Title: "S", Parent: epic.ID, Owner: "a", Now: t0})
	task, _ := r.Create(workitem.NewOptions{Type: workitem.Task, Title: "T", Parent: story.ID, Owner: "a", Now: t0})
	if _, err := r.OpenStream(story, workitem.StreamOptions{Agent: "claude", Now: t0}); err != nil {
		t.Fatal(err)
	}
	th, err := New(r, NewOptions{Title: "Task question", On: task.ID, Author: "alex", Text: "?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if th.Anchor.Item != task.ID || !strings.HasPrefix(th.Anchor.Path, "wip/kanban/tasks/T-0001") {
		t.Fatalf("item anchor: %+v", th.Anchor)
	}
	if StoryOf(r, th) != story.ID {
		t.Errorf("story of task thread: %s", StoryOf(r, th))
	}
	list, _ := For(r, "t-1")
	if len(list) != 1 {
		t.Errorf("For by item in any padding: %d", len(list))
	}
	if err := MirrorNarrative(r, story.ID); err != nil {
		t.Fatal(err)
	}
	n, _ := os.ReadFile(r.NarrativePath(story.ID))
	if !strings.Contains(string(n), "## Open questions\n<!-- threads:start -->") || !strings.Contains(string(n), "- TH-0001 (open, alex, 2026-09-17): Task question") {
		t.Fatalf("mirror:\n%s", n)
	}
	if _, err := Resolve(r, th.ID, "claude", "", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	_ = MirrorNarrative(r, story.ID)
	n, _ = os.ReadFile(r.NarrativePath(story.ID))
	if strings.Contains(string(n), "threads:start") || strings.Contains(string(n), "TH-0001") {
		t.Errorf("resolved threads leave the mirror:\n%s", n)
	}
	if strings.Count(string(n), "## Open questions") != 1 {
		t.Error("heading must survive")
	}
}
