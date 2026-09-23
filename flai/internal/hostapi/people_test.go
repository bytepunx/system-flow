package hostapi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func TestKeyHashIsTheDashboardsHash(t *testing.T) {
	// Reference values from the TypeScript the dashboard keyed entries with.
	for s, want := range map[string]string{
		"Which port should it use?": "150ucnp",
		"Ünïcödé — 🚢 question":      "1fsj6j3",
		"S-0002 touches flai/cmd, which S-0001 (in progress) also touches as flai/cmd": "16nj7jv",
	} {
		if got := keyHash(s); got != want {
			t.Errorf("keyHash(%q) = %s, want %s", s, got, want)
		}
	}
}

func TestLastLogEntry(t *testing.T) {
	body := "## Open questions\n- First?\n\n## Log\n\n### 2026-09-20T08:00:00Z\nOpened.\n\n### 2026-09-20T09:00:00Z\nDid a thing.\nAnd another.\n\n## After\n### 2026-09-21T00:00:00Z\nnot the log\n"
	last := lastLogEntry(body)
	if last == nil || last.At != "2026-09-20T09:00:00Z" || last.Text != "Did a thing.\nAnd another." {
		t.Errorf("last log: %+v", last)
	}
	if lastLogEntry("no log\n") != nil {
		t.Error("nothing to find must find nothing")
	}
}

// people is harbour with work under way: S-0001 in progress with a narrative,
// an open question, a task in progress that is blocked; S-0002 in progress
// touching the same path; S-0004 in review; an agent wrote last on the thread.
func people(t *testing.T) channel.Project {
	t.Helper()
	p := harbour(t)
	repo, _ := workitem.Open(p.Root)
	set := func(id string, f func(it *workitem.Item)) {
		it, err := repo.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		f(it)
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	move := func(id string, states ...string) {
		for i, to := range states {
			it, _ := repo.Get(id)
			if _, err := repo.Transition(it, to, "claude", "", t0.Add(time.Duration(10+i)*time.Minute)); err != nil {
				t.Fatalf("%s to %s: %v", id, to, err)
			}
		}
	}
	set("S-0001", func(it *workitem.Item) { it.Touches = []string{"flai/cmd"} })
	set("S-0002", func(it *workitem.Item) { it.Touches = []string{"flai/cmd/move.go"} })
	berths, _ := repo.Get("S-0001")
	if _, err := repo.OpenStream(berths, workitem.StreamOptions{Agent: "claude", Session: "abc", Now: t0.Add(5 * time.Minute)}); err != nil {
		t.Fatal(err)
	}
	move("S-0001", workitem.InProgress)
	move("T-0001", workitem.Ready, workitem.InProgress)
	set("T-0001", func(it *workitem.Item) {
		if err := workitem.BlockItem(it, "tide is out", t0.Add(20*time.Minute)); err != nil {
			t.Fatal(err)
		}
	})
	if _, err := repo.LogStream("S-0001", "Dredging started.", workitem.StreamOptions{Agent: "claude", Now: t0.Add(30 * time.Minute)}); err != nil {
		t.Fatal(err)
	}
	narrative := repo.NarrativePath("S-0001")
	data, _ := os.ReadFile(narrative)
	_ = os.WriteFile(narrative, []byte(strings.Replace(string(data), "## Open questions\n", "## Open questions\n- Which port should it use?\n", 1)), 0o644)

	// S-0002 in progress too, touching the same place: the check's overlap rule counts work in progress.
	move("S-0002", workitem.Ready, workitem.InProgress)
	customs, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Customs", Parent: "E-0001", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(customs.Path)
	_ = os.WriteFile(customs.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Forms", Parent: customs.ID, Now: t0}); err != nil {
		t.Fatal(err)
	}
	move(customs.ID, workitem.Ready, workitem.InProgress, workitem.Review)
	if _, err := threads.Reply(repo, "TH-0001", "claude", "Nine, I think.", t0.Add(40*time.Minute)); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(repo.AgentsDir(), "README.md"), []byte("# not a narrative\n\n## Open questions\n- ignored?\n"), 0o644)
	return p
}

func TestActivityGet(t *testing.T) {
	p := people(t)
	var got struct{ Streams []StreamActivity }
	if err := call(t, p, "activity.get", `{}`, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Streams) != 1 {
		t.Fatalf("streams: %+v", got.Streams)
	}
	s := got.Streams[0]
	if s.Stream != "S-0001" || s.Agent != "claude" || s.Session != "abc" || s.Status != "in-progress" || !s.Blocked ||
		s.Task == nil || s.Task.ID != "T-0001" || s.LastLog == nil || s.LastLog.Text != "Dredging started." ||
		s.Path != "wip/agents/S-0001.md" || s.Updated != "2026-09-20T08:30:00Z" || s.AgeSeconds != 30*60 {
		t.Errorf("%+v task %+v log %+v", s, s.Task, s.LastLog)
	}
}

func TestInboxDesigner(t *testing.T) {
	p := people(t)
	var in DesignerInbox
	if err := call(t, p, "inbox.designer", `{}`, &in); err != nil {
		t.Fatal(err)
	}
	byKind := map[string]InboxEntry{}
	for _, e := range in.Entries {
		byKind[e.Kind] = e
	}
	if in.Total != 5 || in.Counts["thread"] != 1 || in.Counts["question"] != 1 || in.Counts["review"] != 1 || in.Counts["blocked"] != 1 || in.Counts["overlap"] != 1 {
		t.Fatalf("counts %v total %d entries %+v", in.Counts, in.Total, in.Entries)
	}
	if e := byKind["thread"]; e.Key != "thread:TH-0001" || e.Item != "S-0001" || e.Detail != "claude wrote last, on S-0001" || e.At != "2026-09-20T08:40:00Z" {
		t.Errorf("thread: %+v", e)
	}
	if e := byKind["question"]; e.Key != "question:S-0001:150ucnp" || e.Title != "Which port should it use?" || e.Path != "wip/agents/S-0001.md" || e.Detail != "asked in the narrative of S-0001" || e.Item != "S-0001" {
		t.Errorf("question: %+v", e)
	}
	if e := byKind["review"]; e.Key != "review:S-0004" || e.Item != "S-0004" || e.Title != "S-0004 Customs" {
		t.Errorf("review: %+v", e)
	}
	if e := byKind["blocked"]; e.Key != "blocked:T-0001:2026-09-20T08:20:00Z" || e.Detail != "blocked: tide is out" {
		t.Errorf("blocked: %+v", e)
	}
	if e := byKind["overlap"]; e.Item == "" || !strings.Contains(e.Title, "also touches") || e.Key != "overlap:"+keyHash(e.Title) {
		t.Errorf("overlap: %+v", e)
	}
}

// A hand-written open question has no "answered" marker of its own (unlike
// a thread), so once the story it belongs to no longer needs the designer
// (in review, done, or cancelled), it is left out of the inbox whatever the
// narrative still says (S-0089): the move-to-review rule (workitem package)
// stops this from happening through flai move, but this is the inbox's own
// guarantee, checked here by moving the status directly, bypassing that rule,
// the way an item from before it existed could still be found.
func TestInboxDesignerLeavesOutAQuestionOnceItsStoryNoLongerNeedsIt(t *testing.T) {
	p := people(t)
	repo, _ := workitem.Open(p.Root)
	for _, status := range []string{workitem.Review, workitem.Done, workitem.Cancelled} {
		it, err := repo.Get("S-0001")
		if err != nil {
			t.Fatal(err)
		}
		it.Status = status
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
		var in DesignerInbox
		if err := call(t, p, "inbox.designer", `{}`, &in); err != nil {
			t.Fatal(err)
		}
		for _, e := range in.Entries {
			if e.Kind == "question" {
				t.Errorf("status %s: a question still shows: %+v", status, e)
			}
		}
		if in.Counts["question"] != 0 {
			t.Errorf("status %s: question count %d", status, in.Counts["question"])
		}
	}
}
