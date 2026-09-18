package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 18, 17, 0, 0, 0, time.UTC)

type fixture struct {
	repo  *workitem.Repo
	cs    *mcp.ClientSession
	story *workitem.Item
	task  *workitem.Item
}

func setup(t *testing.T) *fixture {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/system", "docs/users", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte("---\ntitle: Plan\n---\n\n# Plan\n\n## Shape\ntext\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "secret.md"), []byte("not served\n"), 0o644)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	epic, _ := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: t0})
	story, _ := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Story", Parent: epic.ID, Owner: "alex", Touches: []string{"flai/internal/mcpserver"}, Now: t0})
	task, _ := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Task", Parent: story.ID, Owner: "alex", Now: t0})
	sp := story.Path
	data, _ := os.ReadFile(sp)
	_ = os.WriteFile(sp, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
	story, _ = repo.Get(story.ID)
	for _, st := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.Transition(story, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := repo.OpenStream(story, workitem.StreamOptions{Agent: "claude", Now: t0}); err != nil {
		t.Fatal(err)
	}

	srv := New(Options{Repo: repo, Agent: "claude", Version: "test", Now: func() time.Time { return t0.Add(time.Minute) }, Poll: 20 * time.Millisecond, MaxWait: 3 * time.Second})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test-agent", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return &fixture{repo: repo, cs: cs, story: story, task: task}
}

// call invokes a tool and decodes its structured result; failed is the
// tool-level error text, which is how rule refusals reach an agent.
func (f *fixture) call(t *testing.T, name string, args any) (out map[string]any, failed string) {
	t.Helper()
	res, err := f.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("%s: protocol error: %v", name, err)
	}
	if res.IsError {
		for _, c := range res.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				failed += tc.Text
			}
		}
		return nil, failed
	}
	data, _ := json.Marshal(res.StructuredContent)
	_ = json.Unmarshal(data, &out)
	return out, ""
}

func TestToolsAreAdvertised(t *testing.T) {
	f := setup(t)
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, tool := range res.Tools {
		names = append(names, tool.Name)
		if tool.Description == "" || tool.InputSchema == nil {
			t.Errorf("%s needs a description and an input schema", tool.Name)
		}
	}
	sort.Strings(names)
	want := "doc_get inbox item_get item_move thread_get thread_open thread_reply thread_resolve wait_for_events who_touches"
	if strings.Join(names, " ") != want {
		t.Errorf("tools: %v", names)
	}
}

func TestInboxReplyAndMirror(t *testing.T) {
	f := setup(t)
	th, err := threads.New(f.repo, threads.NewOptions{Title: "Which library?", On: f.task.ID, Author: "alex", Text: "SDK or hand-rolled?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	out, _ := f.call(t, "inbox", map[string]any{})
	if out["agent"] != "claude" || out["awaiting_you"].(float64) != 1 {
		t.Fatalf("inbox: %v", out)
	}
	first := out["threads"].([]any)[0].(map[string]any)
	if first["id"] != th.ID || first["awaiting"] != "you" || first["story"] != f.story.ID || first["last_entry"] != "SDK or hand-rolled?" {
		t.Errorf("thread summary: %v", first)
	}
	if out, _ := f.call(t, "inbox", map[string]any{"story": "S-99"}); len(out["threads"].([]any)) != 0 {
		t.Errorf("story filter: %v", out)
	}

	reply, _ := f.call(t, "thread_reply", map[string]any{"id": "th-1", "text": "The official Go SDK."})
	if reply["status"] != "answered" || reply["awaiting"] != "other" || len(reply["entry_list"].([]any)) != 2 {
		t.Errorf("reply: %v", reply)
	}
	back, _ := threads.Get(f.repo, th.ID)
	if e := back.Entries(); e[1].Author != "claude" || e[1].At != "2026-09-18T17:01:00Z" {
		t.Errorf("the reply is written as the agent through the threads package: %+v", e)
	}
	narrative, _ := os.ReadFile(f.repo.NarrativePath(f.story.ID))
	if !strings.Contains(string(narrative), th.ID+" (answered") {
		t.Errorf("narrative mirror not refreshed:\n%s", narrative)
	}
	if out, _ := f.call(t, "inbox", map[string]any{}); out["awaiting_you"].(float64) != 0 || len(out["threads"].([]any)) != 0 {
		t.Errorf("answered threads leave the default inbox: %v", out)
	}
	if out, _ := f.call(t, "inbox", map[string]any{"all": true}); len(out["threads"].([]any)) != 1 {
		t.Errorf("all=true keeps them: %v", out)
	}
	opened, _ := f.call(t, "thread_open", map[string]any{"on": "design/system/plan.md", "heading": "Shape", "title": "Is the shape final?", "text": "Asking before I build on it."})
	if opened["id"] != "TH-0002" || opened["last_by"] != "claude" {
		t.Errorf("open: %v", opened)
	}
	if _, failed := f.call(t, "thread_open", map[string]any{"on": "design/system/plan.md", "heading": "Nope", "title": "x", "text": "y"}); !strings.Contains(failed, "heading") {
		t.Errorf("a bad anchor should be refused with the reason: %q", failed)
	}
	resolved, _ := f.call(t, "thread_resolve", map[string]any{"id": th.ID, "reason": "decided"})
	if resolved["status"] != "resolved" {
		t.Errorf("resolve: %v", resolved)
	}
}

func TestItemsAndRules(t *testing.T) {
	f := setup(t)
	got, _ := f.call(t, "item_get", map[string]any{"id": "s-1"})
	if got["id"] != f.story.ID || got["status"] != "in-progress" || len(got["children"].([]any)) != 1 || !strings.Contains(got["body"].(string), "## Goal") {
		t.Errorf("item_get: %v", got)
	}
	for _, to := range []string{"ready", "in-progress"} {
		if out, failed := f.call(t, "item_move", map[string]any{"id": f.task.ID, "to": to}); failed != "" || out["status"] != to {
			t.Fatalf("move task to %s: %v %s", to, out, failed)
		}
	}
	task, _ := f.repo.Get(f.task.ID)
	if n := len(task.Transitions); task.Transitions[n-1].By != "claude" {
		t.Errorf("transitions are recorded as the agent: %+v", task.Transitions)
	}
	if _, failed := f.call(t, "item_move", map[string]any{"id": f.task.ID, "to": "backlog"}); !strings.Contains(failed, "cannot go") {
		t.Errorf("the workflow rule should reach the agent verbatim: %q", failed)
	}
	if _, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "done"}); !strings.Contains(failed, "only the operator") {
		t.Errorf("agents must not accept stories: %q", failed)
	}
	story, _ := f.repo.Get(f.story.ID)
	if story.Status != "in-progress" {
		t.Errorf("the refused move must change nothing: %s", story.Status)
	}
	who, _ := f.call(t, "who_touches", map[string]any{"path": "flai/internal/mcpserver/server.go"})
	if items := who["items"].([]any); len(items) != 1 || items[0].(map[string]any)["id"] != f.story.ID {
		t.Errorf("who_touches: %v", who)
	}
	if who, _ := f.call(t, "who_touches", map[string]any{"path": "docs"}); len(who["items"].([]any)) != 0 {
		t.Errorf("nobody touches docs: %v", who)
	}
}

func TestDocumentsAndResources(t *testing.T) {
	f := setup(t)
	doc, _ := f.call(t, "doc_get", map[string]any{"path": "design/system/plan.md"})
	if doc["front_matter"] != "title: Plan\n" && !strings.Contains(doc["front_matter"].(string), "title: Plan") {
		t.Errorf("front matter: %v", doc)
	}
	if !strings.Contains(doc["body"].(string), "## Shape") {
		t.Errorf("body: %v", doc)
	}
	for _, bad := range []string{"secret.md", "../outside.md", "design/system/plan.txt", "/etc/passwd", "design/../secret.md"} {
		if _, failed := f.call(t, "doc_get", map[string]any{"path": bad}); failed == "" {
			t.Errorf("%s must be refused", bad)
		}
	}
	res, err := f.cs.ReadResource(context.Background(), &mcp.ReadResourceParams{URI: "flai://design/system/plan.md"})
	if err != nil || len(res.Contents) != 1 || !strings.Contains(res.Contents[0].Text, "# Plan") {
		t.Fatalf("resource: %v %v", res, err)
	}
	if _, err := f.cs.ReadResource(context.Background(), &mcp.ReadResourceParams{URI: "flai://design/../secret.md"}); err == nil {
		t.Error("a resource outside the folder must not be served")
	}
}

// The round trip the story exists for: an idle agent holds wait_for_events,
// the designer replies (through the same package the CLI and dashboard use),
// the wait returns within a second, and the inbox shows the thread.
func TestDesignerReplyReachesAWaitingAgent(t *testing.T) {
	f := setup(t)
	th, err := threads.New(f.repo, threads.NewOptions{Title: "Ready for review?", On: f.story.ID, Author: "claude", Text: "I think it is.", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if out, _ := f.call(t, "inbox", map[string]any{}); out["awaiting_you"].(float64) != 0 {
		t.Fatalf("my own thread awaits the designer, not me: %v", out)
	}
	type waited struct {
		out     map[string]any
		elapsed time.Duration
	}
	done := make(chan waited, 1)
	start := time.Now()
	go func() {
		out, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 3})
		done <- waited{out, time.Since(start)}
	}()
	time.Sleep(150 * time.Millisecond)
	if _, err := threads.Reply(f.repo, th.ID, "alex", "Not yet: add the docs.", t0.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	select {
	case w := <-done:
		changed := w.out["changed"].([]any)
		if w.out["timed_out"] == true || len(changed) == 0 || !strings.Contains(changed[0].(string), "wip/threads/"+th.ID) {
			t.Fatalf("wait: %v", w.out)
		}
		if w.elapsed > 1500*time.Millisecond {
			t.Errorf("the agent should hear within a second, took %v", w.elapsed)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_events never returned")
	}
	out, _ := f.call(t, "inbox", map[string]any{"story": f.story.ID})
	first := out["threads"].([]any)[0].(map[string]any)
	if out["awaiting_you"].(float64) != 1 || first["last_by"] != "alex" || first["last_entry"] != "Not yet: add the docs." {
		t.Errorf("inbox after the reply: %v", out)
	}
	// nothing changes: the wait ends on the (capped) timeout
	quiet, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 1})
	if quiet["timed_out"] != true || len(quiet["changed"].([]any)) != 0 {
		t.Errorf("quiet wait: %v", quiet)
	}
}

// A story with no tasks moves to ready and in-progress through item_move and
// is refused at review: the pulling agent writes the tasks (ADR-0021).
func TestStoryWithoutTasksMovesUntilReview(t *testing.T) {
	f := setup(t)
	s, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "No tasks", Parent: f.story.Parent, Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(s.Path)
	_ = os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	for _, to := range []string{"ready", "in-progress"} {
		if out, failed := f.call(t, "item_move", map[string]any{"id": s.ID, "to": to}); failed != "" || out["status"] != to {
			t.Fatalf("move story without tasks to %s: %v %s", to, out, failed)
		}
	}
	if _, failed := f.call(t, "item_move", map[string]any{"id": s.ID, "to": "review"}); !strings.Contains(failed, "needs at least one task before it goes to review") {
		t.Errorf("review without tasks should be refused with the rule: %q", failed)
	}
}
