package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 18, 17, 0, 0, 0, time.UTC)

type fixture struct {
	repo  *workitem.Repo
	cs    *mcp.ClientSession
	story *workitem.Item
	task  *workitem.Item
	clock *time.Time // what the server reads as now; tests move it
}

func setup(t *testing.T) *fixture {
	t.Helper()
	return setupWith(t, nil)
}

// setupWith is setup with the server's options changed by with.
func setupWith(t *testing.T, with func(*Options)) *fixture {
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

	clock := t0.Add(time.Minute)
	opt := Options{Repo: repo, Agent: "claude", Version: "test", Now: func() time.Time { return clock }, Poll: 20 * time.Millisecond, MaxWait: 3 * time.Second}
	if with != nil {
		with(&opt)
	}
	srv := New(opt)
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
	return &fixture{repo: repo, cs: cs, story: story, task: task, clock: &clock}
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
	want := "activity_log adr_new agent_restart agent_start analyze board criteria_tick doc_get doc_search inbox issue_bump issue_close issue_new issue_story item_edit item_get item_move item_new message_escalate message_get message_reply message_send message_share order_by_policy permission_prompt plan prime promote_candidates release_evaluate release_publish shared_paths shared_paths_edit story_start stream_state task_done test thread_get thread_open thread_reply thread_resolve verify versions wait_for_events wait_for_work who_touches"
	if strings.Join(names, " ") != want {
		t.Errorf("tools: %v", names)
	}
	in := f.cs.InitializeResult().Instructions
	for _, want := range []string{"When you start a story with story_start, its answer holds the story's prime pack", "when you take up a story of yours again, call prime with it", "A brief is not the document", "doc_get and its heading before relying on it or changing what it describes", "doc_search"} {
		if !strings.Contains(in, want) {
			t.Errorf("the instructions do not say to prime with the pack and fetch what it briefs (%q): %s", want, in)
		}
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
	part, failed := f.call(t, "doc_get", map[string]any{"path": "design/system/plan.md", "heading": "shape"})
	if failed != "" || part["heading"] != "Shape" || part["body"] != "## Shape\ntext\n" || part["line"] != float64(7) || !strings.Contains(part["front_matter"].(string), "title: Plan") {
		t.Errorf("doc_get with a heading: %v %s", part, failed)
	}
	if _, failed := f.call(t, "doc_get", map[string]any{"path": "design/system/plan.md", "heading": "Size"}); !strings.Contains(failed, `its headings are: "Plan", "Shape"`) {
		t.Errorf("an unknown heading: %q", failed)
	}
	for _, bad := range []string{"secret.md", "../outside.md", "design/system/plan.txt", "/etc/passwd", "design/../secret.md"} {
		if _, failed := f.call(t, "doc_get", map[string]any{"path": bad}); failed == "" {
			t.Errorf("%s must be refused", bad)
		}
	}
	found, failed := f.call(t, "doc_search", map[string]any{"query": "the shape"})
	hits, _ := found["hits"].([]any)
	if failed != "" || found["query"] != "the shape" || len(hits) != 1 {
		t.Fatalf("doc_search: %v %s", found, failed)
	}
	if h := hits[0].(map[string]any); h["path"] != "design/system/plan.md" || h["heading"] != "Shape" || h["lines"] != "text" || h["title"] != "Plan" {
		t.Errorf("hit: %v", h)
	}
	if _, failed := f.call(t, "doc_search", map[string]any{"query": "the a"}); !strings.Contains(failed, "no words to rank by") {
		t.Errorf("a query of stopwords: %q", failed)
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

// S-0285 (I-0083): wait_for_events does not see a sub-agent finish, so its
// description says what it is for (a thread awaiting the designer) and how a
// story's agent waits for a sub-agent instead, in a project and across a
// folder alike.
func TestWaitForEventsSaysItDoesNotSeeASubAgentFinish(t *testing.T) {
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	for name, cs := range map[string]*mcp.ClientSession{"a project": setup(t).cs, "a folder": folderSetup(t, root).cs} {
		res, err := cs.ListTools(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		var desc string
		for _, tool := range res.Tools {
			if tool.Name == "wait_for_events" {
				desc = tool.Description
			}
		}
		if desc == "" {
			t.Fatalf("%s: no wait_for_events", name)
		}
		for _, want := range []string{
			"the designer's answer on a thread awaiting them",
			"It does not see a sub-agent finish",
			"run_in_background set to false",
			"returns the sub-agent's result as the tool's result",
			"flai guard refuses a story's agent this call while a sub-agent of its session runs and no thread on its story is open",
		} {
			if !strings.Contains(desc, want) {
				t.Errorf("%s: wait_for_events's description lacks %q: %s", name, want, desc)
			}
		}
		if strings.Contains(desc, "Hold this when idle to react to the designer within a second") {
			t.Errorf("%s: the old sentence, which reads as for any wait, is still there: %s", name, desc)
		}
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

// readyStory creates a story with criteria and moves it to ready as the
// designer at the given time. It touches a path of its own unless touches
// are given, so that no open story holds it (S-0128).
func (f *fixture) readyStory(t *testing.T, title string, at time.Time, touches ...string) *workitem.Item {
	t.Helper()
	if touches == nil {
		touches = []string{"docs/" + strings.ToLower(strings.ReplaceAll(title, " ", "-"))}
	}
	s, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: f.story.Parent, Owner: "alex", Touches: touches, Now: at})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(s.Path)
	_ = os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	s, _ = f.repo.Get(s.ID)
	if _, err := f.repo.Transition(s, workitem.Ready, "alex", "", at); err != nil {
		t.Fatal(err)
	}
	return s
}

func changeSummaries(out map[string]any, key string) (got []string) {
	for _, e := range out[key].([]any) {
		got = append(got, e.(map[string]any)["summary"].(string))
	}
	return got
}

// S-0058: the designer moves a story to ready between two calls; the agent
// was not waiting, and still hears of it, once, and sees it as ready work
// until it is pulled.
func TestInboxReportsReadyWorkAndWhatOthersChanged(t *testing.T) {
	f := setup(t)
	first, _ := f.call(t, "inbox", map[string]any{})
	if len(first["ready"].([]any)) != 0 || first["can_pull"] != true {
		t.Fatalf("nothing is ready yet: %v", first)
	}

	*f.clock = t0.Add(10 * time.Minute)
	s := f.readyStory(t, "Made ready meanwhile", t0.Add(5*time.Minute))
	out, _ := f.call(t, "inbox", map[string]any{})
	ready := out["ready"].([]any)
	if len(ready) != 1 || ready[0].(map[string]any)["id"] != s.ID {
		t.Fatalf("ready: %v", out["ready"])
	}
	if got := changeSummaries(out, "changes"); len(got) != 1 || got[0] != s.ID+" Made ready meanwhile moved to ready by alex" {
		t.Errorf("changes: %v", got)
	}

	again, _ := f.call(t, "inbox", map[string]any{})
	if len(again["changes"].([]any)) != 0 {
		t.Errorf("a change is reported once: %v", again["changes"])
	}
	if len(again["ready"].([]any)) != 1 {
		t.Errorf("ready work is state and stays listed: %v", again["ready"])
	}

	// the agent pulls it: its own move is not news to it, and the story is no longer ready
	*f.clock = t0.Add(11 * time.Minute)
	if _, failed := f.call(t, "item_move", map[string]any{"id": s.ID, "to": "in-progress"}); failed != "" {
		t.Fatal(failed)
	}
	*f.clock = t0.Add(12 * time.Minute)
	after, _ := f.call(t, "inbox", map[string]any{})
	if len(after["changes"].([]any)) != 0 || len(after["ready"].([]any)) != 0 {
		t.Errorf("own move echoed, or story still ready: %v", after)
	}
	// two stories in progress against the default limit of two
	if after["can_pull"] != false {
		t.Errorf("can_pull should follow the in-progress limit: %v", after["can_pull"])
	}
	if _, err := os.Stat(filepath.Join(f.repo.CacheDir(), "mcp", "claude.json")); err != nil {
		t.Errorf("the cursor lives under .flai-cache/mcp: %v", err)
	}
}

func TestWaitForEventsReturnsWhatIsAlreadyBehindTheCursor(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" { // the agent looks, then goes away
		t.Fatal(failed)
	}
	*f.clock = t0.Add(10 * time.Minute)
	s := f.readyStory(t, "While nobody waited", t0.Add(5*time.Minute))
	start := time.Now()
	out, failed := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 3})
	if failed != "" || out["timed_out"] == true || time.Since(start) > 2*time.Second {
		t.Fatalf("a change behind the cursor returns at once: %v %s after %s", out, failed, time.Since(start))
	}
	if got := changeSummaries(out, "events"); len(got) != 1 || !strings.Contains(got[0], s.ID) || !strings.Contains(got[0], "moved to ready by alex") {
		t.Errorf("events: %v", got)
	}
}

func TestWaitForEventsReportsAnEventWhileHeld(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	*f.clock = t0.Add(10 * time.Minute)
	go func() {
		time.Sleep(150 * time.Millisecond)
		it, _ := f.repo.Get(f.story.ID)
		// stamped with the very second the cursor stands at: the boundary case
		if err := workitem.BlockItem(it, "waiting on the designer", t0.Add(10*time.Minute)); err == nil {
			_ = f.repo.Save(it)
		}
	}()
	out, failed := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 3})
	if failed != "" || out["timed_out"] == true {
		t.Fatalf("held wait: %v %s", out, failed)
	}
	if got := changeSummaries(out, "events"); len(got) != 1 || !strings.Contains(got[0], "was blocked: waiting on the designer") {
		t.Errorf("events: %v", got)
	}
	if len(out["changed"].([]any)) == 0 {
		t.Errorf("the changed paths are still reported: %v", out)
	}
}

// I-0059: wait_for_events holds for the timeout asked, up to 30 minutes,
// and for a minute when none is asked.
func TestWaitForEventsHoldsUpToTheLongestWait(t *testing.T) {
	d := newDeadlines()
	f := setupWith(t, func(o *Options) { o.MaxWait, o.After = 0, d.after })
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	*f.clock = t0.Add(10 * time.Minute)
	checkHolds(t, d, func(args map[string]any) map[string]any {
		out, failed := f.call(t, "wait_for_events", args)
		if failed != "" {
			t.Errorf("wait_for_events: %s", failed)
		}
		return out
	}, time.Minute)
}

func TestBoardToolMatchesTheSharedView(t *testing.T) {
	f := setup(t)
	s := f.readyStory(t, "On the board", t0)
	out, failed := f.call(t, "board", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	cols := out["columns"].(map[string]any)
	if len(cols["ready"].([]any)) != 1 || cols["ready"].([]any)[0].(map[string]any)["id"] != s.ID || len(cols["in-progress"].([]any)) != 1 {
		t.Errorf("columns: %v", cols)
	}
	if out["wip_limits"] == nil || out["order"] == nil || out["breaches"] == nil {
		t.Errorf("limits, order, and breaches belong to the view: %v", out)
	}
	all, _ := f.call(t, "board", map[string]any{"all": true})
	count := func(cols map[string]any) (n int) {
		for _, c := range cols {
			n += len(c.([]any))
		}
		return n
	}
	if n := count(all["columns"].(map[string]any)) - count(cols); n < 2 {
		t.Errorf("all adds the epic and the task: %d", n)
	}
}

// S-0128: the board and inbox mark a ready story whose claim overlaps an
// open story's as held, with the reason; the fixture's story is in progress
// and touches flai/internal/mcpserver.
func TestBoardAndInboxMarkHeldStories(t *testing.T) {
	f := setup(t)
	held := f.readyStory(t, "Overlaps", t0, "flai/internal")
	free := f.readyStory(t, "Elsewhere", t0)
	bare := f.readyStory(t, "Declares nothing", t0, []string{}...)
	want := map[string]string{
		held.ID: "held (overlap): touches flai/internal, which holds flai/internal/mcpserver that " + f.story.ID + " (in progress) touches; starts when " + f.story.ID + " moves to review, is cancelled, or is sent back",
		free.ID: "",
		bare.ID: "held (no-touches): declares no touches, so it may change what " + f.story.ID + " (in progress) changes; starts when it declares touches that overlap no story's in progress, or when " + f.story.ID + " moves to review, is cancelled, or is sent back",
	}
	check := func(what string, cards []any) {
		t.Helper()
		if len(cards) != 3 {
			t.Fatalf("%s: %v", what, cards)
		}
		for _, c := range cards {
			card := c.(map[string]any)
			got := ""
			if h, ok := card["held"].(map[string]any); ok {
				got = h["reason"].(string)
			}
			if id := card["id"].(string); got != want[id] {
				t.Errorf("%s: %s held %q, want %q", what, id, got, want[id])
			}
		}
	}
	board, _ := f.call(t, "board", map[string]any{})
	check("board", board["columns"].(map[string]any)["ready"].([]any))
	inbox, _ := f.call(t, "inbox", map[string]any{})
	check("inbox", inbox["ready"].([]any))
}

func TestPullOrderChangeIsReportedOnlyWhenPrioritiesChange(t *testing.T) {
	for _, c := range []struct {
		before, after []string
		want          bool
	}{
		{[]string{"S-1", "S-2"}, []string{"S-1", "S-2", "S-3"}, false}, // a story became ready
		{[]string{"S-1", "S-2"}, []string{"S-2"}, false},               // one was started
		{[]string{"S-1", "S-2", "S-3"}, []string{"S-3", "S-1"}, true},  // S-3 now comes first
		{nil, []string{"S-1"}, false},
	} {
		if got := workitem.Reordered(c.before, c.after); got != c.want {
			t.Errorf("Reordered(%v, %v) = %v, want %v", c.before, c.after, got, c.want)
		}
	}
	f := setup(t)
	a := f.readyStory(t, "A", t0)
	b := f.readyStory(t, "B", t0)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	board, _ := f.repo.LoadBoard()
	board.Order = []string{b.ID, a.ID}
	if err := board.Save("2026-09-18"); err != nil {
		t.Fatal(err)
	}
	*f.clock = t0.Add(10 * time.Minute)
	out, _ := f.call(t, "inbox", map[string]any{})
	if got := changeSummaries(out, "changes"); len(got) != 1 || got[0] != "the pull order is now "+b.ID+", "+a.ID {
		t.Errorf("changes: %v", got)
	}
	if ready := out["ready"].([]any); ready[0].(map[string]any)["id"] != b.ID {
		t.Errorf("ready follows the new order: %v", ready)
	}
}

// S-0061: the first look of a new agent on a busy repository was 209 changes
// and 68 KB. A look is capped, the newest kept and the rest counted, a first
// look leaves task transitions out, and the cursor passes all of it.
func TestAFirstLookIsBoundedAndLeavesTasksOut(t *testing.T) {
	// what the fixture alone holds for a first look: its own story and epic
	quiet, _ := setup(t).call(t, "inbox", map[string]any{})
	already := len(quiet["changes"].([]any))

	f := setup(t)
	// a busy day before the agent's first call: more story changes than the
	// cap, and task transitions among them
	for i := 0; i < maxEvents+7; i++ {
		f.readyStory(t, fmt.Sprintf("Busy %02d", i), t0.Add(time.Duration(i+1)*time.Minute))
	}
	task, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "A task moved today", Parent: f.story.ID, Owner: "alex", Now: t0.Add(2 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.Transition(task, workitem.Ready, "alex", "", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	*f.clock = t0.Add(3 * time.Hour)

	first, _ := f.call(t, "inbox", map[string]any{})
	changes := first["changes"].([]any)
	if len(changes) != maxEvents {
		t.Fatalf("a look reports at most %d changes, got %d", maxEvents, len(changes))
	}
	if got := first["changes_omitted"]; got != float64(7+already) {
		t.Errorf("the older ones are counted: changes_omitted = %v, want %d", got, 7+already)
	}
	for _, c := range changes {
		if c.(map[string]any)["type"] == "task" {
			t.Errorf("a first look leaves task transitions out: %v", c)
		}
	}
	summaries := changeSummaries(first, "changes")
	if !strings.Contains(summaries[0], "Busy 07") || !strings.Contains(summaries[len(summaries)-1], fmt.Sprintf("Busy %02d", maxEvents+6)) {
		t.Errorf("the newest are kept, oldest first: %q … %q", summaries[0], summaries[len(summaries)-1])
	}
	if ready := first["ready"].([]any); len(ready) != maxEvents+7 {
		t.Errorf("ready work is state and is not capped: %d", len(ready))
	}

	again, _ := f.call(t, "inbox", map[string]any{})
	if len(again["changes"].([]any)) != 0 || again["changes_omitted"] != float64(0) {
		t.Errorf("nothing left out or filtered comes back: %v omitted %v", again["changes"], again["changes_omitted"])
	}
}

// With a cursor, task transitions are news again, and the cap still holds,
// for wait_for_events as for inbox.
func TestALaterLookReportsTasksAndIsStillCapped(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	*f.clock = t0.Add(10 * time.Minute)
	task, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Moved while away", Parent: f.story.ID, Owner: "alex", Now: t0.Add(5 * time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.Transition(task, workitem.Ready, "alex", "", t0.Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	out, _ := f.call(t, "inbox", map[string]any{})
	if got := changeSummaries(out, "changes"); len(got) != 1 || !strings.Contains(got[0], "Moved while away moved to ready") {
		t.Errorf("a task transition is reported once the agent has a cursor: %v", got)
	}

	for i := 0; i < maxEvents+3; i++ {
		f.readyStory(t, fmt.Sprintf("Later %02d", i), t0.Add(time.Duration(20+i)*time.Minute))
	}
	*f.clock = t0.Add(5 * time.Hour)
	held, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 1})
	if n := len(held["events"].([]any)); n != maxEvents {
		t.Fatalf("wait_for_events is capped too: %d", n)
	}
	if held["events_omitted"] != float64(3) {
		t.Errorf("events_omitted = %v, want 3", held["events_omitted"])
	}
	after, _ := f.call(t, "inbox", map[string]any{})
	if len(after["changes"].([]any)) != 0 {
		t.Errorf("the omitted ones do not come back: %v", after["changes"])
	}
}

// S-0085: an agent is told when someone else edits an item, and what of it
// changed; its own edits are not news to it; a waiting agent wakes for it.
func TestAnEditBySomeoneElseReachesTheAgent(t *testing.T) {
	f := setup(t)
	if out, _ := f.call(t, "inbox", map[string]any{}); out == nil {
		t.Fatal("first look")
	}
	at := f.clock.Add(30 * time.Second)
	*f.clock = f.clock.Add(time.Minute)
	itemedit.Record(f.repo, itemedit.Notice{At: at.Format(workitem.TimeFormat), By: "alex", ID: f.story.ID, Type: workitem.Story, Title: "Story, renamed", Changed: []string{"title", "criteria"}})
	itemedit.Record(f.repo, itemedit.Notice{At: at.Format(workitem.TimeFormat), By: "claude", ID: f.story.ID, Type: workitem.Story, Title: "Story, renamed", Changed: []string{"notes"}})
	out, _ := f.call(t, "inbox", map[string]any{})
	changes := out["changes"].([]any)
	if len(changes) != 1 {
		t.Fatalf("the designer's edit and not my own: %v", changes)
	}
	c := changes[0].(map[string]any)
	if c["kind"] != "edited" || c["id"] != f.story.ID || c["by"] != "alex" || c["to"] != "title,criteria" || !strings.Contains(c["summary"].(string), "was edited by alex: title, criteria") {
		t.Errorf("the change: %v", c)
	}
	if again, _ := f.call(t, "inbox", map[string]any{}); len(again["changes"].([]any)) != 0 {
		t.Errorf("told once: %v", again["changes"])
	}

	// a waiting agent wakes for the notice itself
	// the clock moves before the wait starts, not under it: the server reads it
	*f.clock = f.clock.Add(time.Minute)
	// a notice is stamped when it is written, which is after the waiting agent's look, here
	// within the same second: the cursor's keys, not its time, tell it is news
	stamp := f.clock.Format(workitem.TimeFormat)
	done := make(chan map[string]any, 1)
	go func() {
		out, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 3})
		done <- out
	}()
	time.Sleep(150 * time.Millisecond)
	itemedit.Record(f.repo, itemedit.Notice{At: stamp, By: "alex", ID: f.story.ID, Type: workitem.Story, Title: "Story, renamed", Changed: []string{"touches"}})
	select {
	case w := <-done:
		events := w["events"].([]any)
		if w["timed_out"] == true || len(events) != 1 || events[0].(map[string]any)["to"] != "touches" {
			t.Errorf("wait: %v", w)
		}
		for _, p := range w["changed"].([]any) {
			if strings.Contains(p.(string), ".flai-cache") {
				t.Errorf("the notices wake the agent and are not shown to it as a path: %v", w["changed"])
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_events never returned")
	}
}

// S-0132: an open story whose claim covers what an accepted story changed is
// told which paths changed, once, whoever accepted; a waiting agent wakes for
// it, and the log is not shown to it as a path.
func TestAnOverlapAtAcceptanceReachesTheAgent(t *testing.T) {
	f := setup(t)
	if out, _ := f.call(t, "inbox", map[string]any{}); out == nil {
		t.Fatal("first look")
	}
	at := f.clock.Add(30 * time.Second)
	*f.clock = f.clock.Add(time.Minute)
	itemedit.RecordOverlap(f.repo, itemedit.Overlap{At: at.Format(workitem.TimeFormat), By: "alex", ID: f.story.ID, Title: "Open", Accepted: "S-0099", Paths: []string{"flai/cmd/a.go", "docs/b.md"}})
	out, _ := f.call(t, "inbox", map[string]any{})
	changes := out["changes"].([]any)
	if len(changes) != 1 {
		t.Fatalf("one overlap: %v", changes)
	}
	c := changes[0].(map[string]any)
	if c["kind"] != "overlapped" || c["id"] != f.story.ID || c["cause"] != "S-0099" || c["to"] != "flai/cmd/a.go,docs/b.md" {
		t.Errorf("the change: %v", c)
	}
	if s := c["summary"].(string); !strings.Contains(s, "S-0099 was accepted by alex and changed flai/cmd/a.go, docs/b.md") || !strings.Contains(s, "flai stream sync "+f.story.ID) {
		t.Errorf("summary: %s", s)
	}
	if again, _ := f.call(t, "inbox", map[string]any{}); len(again["changes"].([]any)) != 0 {
		t.Errorf("told once: %v", again["changes"])
	}

	*f.clock = f.clock.Add(time.Minute)
	stamp := f.clock.Format(workitem.TimeFormat)
	done := make(chan map[string]any, 1)
	go func() {
		out, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 3})
		done <- out
	}()
	time.Sleep(150 * time.Millisecond)
	itemedit.RecordOverlap(f.repo, itemedit.Overlap{At: stamp, By: "alex", ID: f.story.ID, Title: "Open", Accepted: "S-0098", Paths: []string{"flai/cmd/c.go"}})
	select {
	case w := <-done:
		events := w["events"].([]any)
		if w["timed_out"] == true || len(events) != 1 || events[0].(map[string]any)["cause"] != "S-0098" {
			t.Errorf("wait: %v", w)
		}
		for _, p := range w["changed"].([]any) {
			if strings.Contains(p.(string), ".flai-cache") {
				t.Errorf("the log wakes the agent and is not shown to it as a path: %v", w["changed"])
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_events never returned")
	}
}

// An overlap's summary names at most maxPaths paths; to has them all.
func TestAnOverlapSummaryBoundsItsPaths(t *testing.T) {
	var paths []string
	for i := 0; i < maxPaths+3; i++ {
		paths = append(paths, fmt.Sprintf("p/%d.go", i))
	}
	s := describe(workitem.Change{ID: "S-0002", Title: "Open", Kind: Overlapped, Cause: "S-0001", To: strings.Join(paths, ","), By: "alex"})
	if !strings.Contains(s, "and 3 more") || strings.Contains(s, "p/10.go") {
		t.Errorf("summary: %s", s)
	}
}

// S-0220, ADR-0090: thread_reply posts a recommendation citing a source,
// which leaves the thread's status as it was and shows as pending in
// thread_get and the inbox; a source that names no file, or a heading not in
// it, is refused. A caller that is not the orchestrator logs nothing.
func TestAThreadReplyRecommendsCitingASource(t *testing.T) {
	logged := 0
	f := setupWith(t, func(o *Options) {
		o.Activities = func(context.Context, string, string, string, []string, string) (ActivityLogged, error) {
			logged++
			return ActivityLogged{}, nil
		}
	})
	th, err := threads.New(f.repo, threads.NewOptions{Title: "Which first?", On: f.story.ID, Author: "alex", Text: "The CLI or the dashboard?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for source, want := range map[string]string{
		"#Shape":                     "names no file",
		"design/system/nope.md":      "design/system/nope.md does not exist",
		"design/system/plan.md#Nope": `heading "Nope" is not in design/system/plan.md`,
	} {
		if _, failed := f.call(t, "thread_reply", map[string]any{"id": th.ID, "text": "The CLI.", "recommendation": true, "source": source}); !strings.Contains(failed, want) {
			t.Errorf("source %q: %q, want it refused with %q", source, failed, want)
		}
	}
	out, failed := f.call(t, "thread_reply", map[string]any{"id": th.ID, "text": "The CLI, as the plan says.", "recommendation": true, "source": "design/system/plan.md#Shape"})
	if failed != "" {
		t.Fatal(failed)
	}
	entries := out["entry_list"].([]any)
	last := entries[len(entries)-1].(map[string]any)
	pending, _ := out["pending_recommendation"].(map[string]any)
	if out["status"] != "open" || last["recommendation"] != true || last["author"] != "claude" || pending["author"] != "claude" {
		t.Errorf("a recommendation leaves the thread open and pending: %v", out)
	}
	if src, _ := last["source"].(map[string]any); src["path"] != "design/system/plan.md" || src["heading"] != "Shape" {
		t.Errorf("source: %v", last["source"])
	}
	inbox, _ := f.call(t, "inbox", map[string]any{"all": true})
	if sum := inbox["threads"].([]any)[0].(map[string]any); sum["status"] != "open" || sum["pending_recommendation"] == nil {
		t.Errorf("the inbox carries the pending recommendation: %v", sum)
	}
	if back, _ := threads.Get(f.repo, th.ID); back.Status != "open" {
		t.Errorf("status on disk: %s", back.Status)
	}
	if logged != 0 {
		t.Errorf("a story's agent's reply was logged %d times; only the orchestrator's is", logged)
	}
}

// S-0220: the orchestrator's thread_reply is logged in its activity
// document, saying whether it recommended or answered and what it cited,
// with the thread's story as its item.
func TestTheOrchestratorsThreadRepliesAreLoggedWithTheirSource(t *testing.T) {
	t.Setenv("FLAI_ROLE", "orchestrate")
	type call struct{ kind, summary, items, by string }
	var got []call
	f := setupWith(t, func(o *Options) {
		o.Agent = "orchestrator"
		o.Activities = func(_ context.Context, _, kind, summary string, items []string, by string) (ActivityLogged, error) {
			got = append(got, call{kind, summary, strings.Join(items, ","), by})
			return ActivityLogged{}, nil
		}
	})
	onStory, err := threads.New(f.repo, threads.NewOptions{Title: "Which first?", On: f.story.ID, Author: "alex", Text: "The CLI or the dashboard?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	onDoc, err := threads.New(f.repo, threads.NewOptions{Title: "Is the shape final?", On: "design/system/plan.md", Author: "alex", Text: "Asking.", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if out, failed := f.call(t, "thread_reply", map[string]any{"id": onStory.ID, "text": "The CLI.", "recommendation": true, "source": "design/system/plan.md#Shape"}); failed != "" || out["status"] != "open" {
		t.Fatalf("recommend: %v %s", out, failed)
	}
	if out, failed := f.call(t, "thread_reply", map[string]any{"id": onDoc.ID, "text": "Yes.", "source": "design/system/plan.md"}); failed != "" || out["status"] != "answered" {
		t.Fatalf("answer: %v %s", out, failed)
	}
	want := []call{
		{"orchestrator", "Recommended an answer on " + onStory.ID + ", citing design/system/plan.md § Shape", f.story.ID, "orchestrator"},
		{"orchestrator", "Answered " + onDoc.ID + ", citing design/system/plan.md", "", "orchestrator"},
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("logged\n  %+v\nwant\n  %+v", got, want)
	}
}

// endFixture is the server of an agent flai serve started with FLAI_STORY
// set to story, as it starts a story's agent (S-0272), and role as FLAI_ROLE;
// the test passes its deadlines. An empty story is a session without
// FLAI_STORY, as a hand-run agent's is.
func endFixture(t *testing.T, story, role string) (*fixture, *deadlines) {
	t.Helper()
	// the server reads them when it starts; the fixture's story is the
	// first story of its repository, S-0001, and FLAI_STORY is trimmed
	t.Setenv("FLAI_STORY", story)
	t.Setenv("FLAI_ROLE", role)
	d := newDeadlines()
	f := setupWith(t, func(o *Options) { o.After = d.after })
	if f.story.ID != "S-0001" {
		t.Fatalf("the fixture's story is %s, not S-0001", f.story.ID)
	}
	return f, d
}

// storysAgent is endFixture for the fixture's story's agent.
func storysAgent(t *testing.T) (*fixture, *deadlines) {
	t.Helper()
	return endFixture(t, " S-0001\n", "")
}

// ask opens a thread on id as the agent and sets the agent's cursor after
// it, so that a wait starts with nothing behind the cursor.
func (f *fixture) ask(t *testing.T, id string) *threads.Thread {
	t.Helper()
	th, err := threads.New(f.repo, threads.NewOptions{Title: "Which way?", On: id, Author: "claude", Text: "A or B?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	return th
}

// endsAtOnce calls wait_for_events and wants it to answer end without
// holding: no deadline is armed.
func (f *fixture) endsAtOnce(t *testing.T, d *deadlines) map[string]any {
	t.Helper()
	start := time.Now()
	out, failed := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 1800})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["end"] != true || out["timed_out"] == true || time.Since(start) > 2*time.Second {
		t.Fatalf("a story's agent with only its question pending ends at once: %v after %s", out, time.Since(start))
	}
	select {
	case got := <-d.armed:
		t.Errorf("a wait told to end armed a deadline of %s", got)
	default:
	}
	return out
}

// holds calls wait_for_events and wants it to hold until its time passes,
// with end false and no why.
func (f *fixture) holds(t *testing.T, d *deadlines) {
	t.Helper()
	done := make(chan map[string]any, 1)
	go func() {
		out, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 60})
		done <- out
	}()
	select {
	case <-d.armed:
	case out := <-done:
		t.Fatalf("the wait should hold, answered at once: %v", out)
	case <-time.After(5 * time.Second):
		t.Fatal("the wait never armed a deadline")
	}
	d.advance(time.Minute)
	select {
	case out := <-done:
		if out["timed_out"] != true || out["end"] != false || out["why"] != nil {
			t.Errorf("a held wait times out with end false: %v", out)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the wait never answered once its time passed")
	}
}

// S-0272: a story's agent whose story has its question open to the designer
// and no task in progress is told to end, with why naming the thread, at
// once rather than held: flai serve starts it again on the answer.
func TestWaitForEventsEndsAStorysAgentWithOnlyItsQuestionPending(t *testing.T) {
	f, d := storysAgent(t)
	th := f.ask(t, f.story.ID)
	out := f.endsAtOnce(t, d)
	why, _ := out["why"].(string)
	if !strings.Contains(why, th.ID) || !strings.Contains(why, f.story.ID) || !strings.Contains(why, "no task in progress") {
		t.Errorf("why names the story, the thread, and the reason: %q", why)
	}
	if events, ok := out["events"].([]any); !ok || len(events) != 0 {
		t.Errorf("events is an empty list: %v", out["events"])
	}
	// the tool and the server's instructions say so
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range res.Tools {
		if tool.Name == "wait_for_events" && !strings.Contains(tool.Description, "flai serve starts it again when the question is answered") {
			t.Errorf("wait_for_events's description does not say when it ends: %s", tool.Description)
		}
	}
	if in := f.cs.InitializeResult().Instructions; !strings.Contains(in, "gets end and why from wait_for_events") {
		t.Errorf("the instructions do not say when wait_for_events ends: %s", in)
	}
}

// Events already behind the cursor come with end, so none is lost.
func TestWaitForEventsReportsWhatIsBehindTheCursorWithEnd(t *testing.T) {
	f, d := storysAgent(t)
	f.ask(t, f.story.ID)
	*f.clock = t0.Add(10 * time.Minute)
	s := f.readyStory(t, "While it asked", t0.Add(5*time.Minute))
	out := f.endsAtOnce(t, d)
	if got := changeSummaries(out, "events"); len(got) != 1 || !strings.Contains(got[0], s.ID) {
		t.Errorf("events behind the cursor come with end: %v", got)
	}
}

// A question on one of the story's tasks counts, as serve counts it.
func TestWaitForEventsEndsOnAQuestionOnAStorysTask(t *testing.T) {
	f, d := storysAgent(t)
	th := f.ask(t, f.task.ID)
	if why, _ := f.endsAtOnce(t, d)["why"].(string); !strings.Contains(why, th.ID) {
		t.Errorf("why names the thread on the task: %q", why)
	}
}

// A recommendation on the agent's question awaits the operator's
// confirmation and is no answer yet (ADR-0090): the agent still ends.
func TestWaitForEventsEndsWhileARecommendationAwaitsConfirmation(t *testing.T) {
	f, d := storysAgent(t)
	th := f.ask(t, f.story.ID)
	if _, err := threads.ReplyWith(f.repo, th.ID, "orchestrator", "B.", t0.Add(time.Minute), threads.Marks{Recommendation: true}); err != nil {
		t.Fatal(err)
	}
	if why, _ := f.endsAtOnce(t, d)["why"].(string); !strings.Contains(why, th.ID) {
		t.Errorf("why names the thread: %q", why)
	}
}

// With a task of the story in progress the agent has work in hand: it holds.
func TestWaitForEventsHoldsAStorysAgentWithATaskInProgress(t *testing.T) {
	f, d := storysAgent(t)
	for _, to := range []string{workitem.Ready, workitem.InProgress} {
		task, _ := f.repo.Get(f.task.ID)
		if _, err := f.repo.Transition(task, to, "claude", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	f.ask(t, f.story.ID)
	f.holds(t, d)
}

// An answered question awaits nothing: the agent holds rather than end on
// the answer it waited for, and so does it once the thread is resolved.
func TestWaitForEventsHoldsOnceTheQuestionIsAnswered(t *testing.T) {
	f, d := storysAgent(t)
	th := f.ask(t, f.story.ID)
	if _, err := threads.Reply(f.repo, th.ID, "alex", "B.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	f.holds(t, d)
	if _, err := threads.Resolve(f.repo, th.ID, "alex", "B it is", t0.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	f.holds(t, d)
}

// Without FLAI_STORY (a hand-run agent) wait_for_events holds as before,
// whatever its threads; with a role (the planner, the orchestrator, the
// analyzer) too.
func TestWaitForEventsHoldsOutsideAStorysAgent(t *testing.T) {
	for name, env := range map[string][2]string{"no FLAI_STORY": {"", ""}, "a role": {"S-0001", "plan"}} {
		t.Run(name, func(t *testing.T) {
			f, d := endFixture(t, env[0], env[1])
			f.ask(t, f.story.ID)
			f.holds(t, d)
		})
	}
}

// otherInProgress makes a second story of the fixture's epic, in progress,
// for the fixture's story to converse with.
func (f *fixture) otherInProgress(t *testing.T) *workitem.Item {
	t.Helper()
	epic, err := f.repo.Get(f.story.Parent)
	if err != nil {
		t.Fatal(err)
	}
	return storyOfEpic(t, f.repo, epic, "Other", nil, workitem.Ready, workitem.InProgress)
}

// send starts a conversation from one story to another, by the fixture's
// agent for its own story and by agent-S-nnnn for another, and sets the
// agent's cursor after it, so that a wait starts with nothing behind it.
func (f *fixture) send(t *testing.T, from, to string) *messages.Conversation {
	t.Helper()
	author := "agent-" + from
	if from == f.story.ID {
		author = "claude"
	}
	c, err := messages.Send(f.repo, messages.SendOptions{From: from, To: to, Author: author, Text: "Will you leave plan.md to me?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	return c
}

// S-0335: a story's agent whose story's message awaits the other story's
// reply, and no task in progress, is told to end, with why naming the
// conversation and the story it waits on: flai serve starts it again on the
// reply or a new message to its story. The reply, when it comes, is an event
// and keeps the hold, so it is never swallowed by an end.
func TestWaitForEventsEndsAStorysAgentAwaitingAnotherStorysReply(t *testing.T) {
	f, d := storysAgent(t)
	other := f.otherInProgress(t)
	c := f.send(t, f.story.ID, other.ID)
	why, _ := f.endsAtOnce(t, d)["why"].(string)
	want := f.story.ID + " has " + c.ID + " awaiting " + other.ID + "'s agent and no task in progress"
	if !strings.HasPrefix(why, want) || !strings.HasSuffix(why, "flai serve starts you again when it is answered or a new message to "+f.story.ID+" comes") {
		t.Errorf("why names the story, the conversation, the story it waits on, and when serve starts it again: %q", why)
	}
	if strings.Contains(why, "designer") {
		t.Errorf("why names no thread when none is open: %q", why)
	}

	*f.clock = t0.Add(10 * time.Minute)
	if _, err := messages.Reply(f.repo, c.ID, other.ID, "agent-"+other.ID, "Yes.", t0.Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	out, failed := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 60})
	if failed != "" {
		t.Fatal(failed)
	}
	if got := messageEvents(out); out["end"] != false || len(got) != 1 || got[0]["id"] != c.ID {
		t.Errorf("the reply comes as an event without end: %v", out)
	}
}

// A question to the designer and a message awaiting another story: why
// names both.
func TestWaitForEventsEndsOnAThreadAndAConversation(t *testing.T) {
	f, d := storysAgent(t)
	other := f.otherInProgress(t)
	th := f.ask(t, f.story.ID)
	c := f.send(t, f.story.ID, other.ID)
	why, _ := f.endsAtOnce(t, d)["why"].(string)
	if want := f.story.ID + " has " + th.ID + " open to the designer and " + c.ID + " awaiting " + other.ID + "'s agent and no task in progress"; !strings.HasPrefix(why, want) {
		t.Errorf("why names the thread and the conversation:\n  %q\nwant it to start\n  %q", why, want)
	}
}

// With a task of the story in progress the agent has work in hand: it holds
// while its message awaits the other story.
func TestWaitForEventsHoldsAStorysAgentAwaitingAReplyWithATaskInProgress(t *testing.T) {
	f, d := storysAgent(t)
	for _, to := range []string{workitem.Ready, workitem.InProgress} {
		task, _ := f.repo.Get(f.task.ID)
		if _, err := f.repo.Transition(task, to, "claude", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	other := f.otherInProgress(t)
	f.send(t, f.story.ID, other.ID)
	f.holds(t, d)
}

// Without FLAI_STORY, or with a role, a conversation awaiting the other
// story does not end the wait either.
func TestWaitForEventsHoldsOutsideAStorysAgentAwaitingAReply(t *testing.T) {
	for name, env := range map[string][2]string{"no FLAI_STORY": {"", ""}, "a role": {"S-0001", "plan"}} {
		t.Run(name, func(t *testing.T) {
			f, d := endFixture(t, env[0], env[1])
			other := f.otherInProgress(t)
			f.send(t, f.story.ID, other.ID)
			f.holds(t, d)
		})
	}
}

// A message to the story awaiting its own reply is work it owes: the agent
// holds, whatever else awaits others, and ends once it has replied.
func TestWaitForEventsHoldsWhileTheStoryOwesAReply(t *testing.T) {
	f, d := storysAgent(t)
	other := f.otherInProgress(t)
	f.ask(t, f.story.ID)
	mine := f.send(t, f.story.ID, other.ID)
	owed := f.send(t, other.ID, f.story.ID)
	f.holds(t, d)
	if _, err := messages.Reply(f.repo, owed.ID, f.story.ID, "claude", "Yes.", t0); err != nil {
		t.Fatal(err)
	}
	why, _ := f.endsAtOnce(t, d)["why"].(string)
	if !strings.Contains(why, mine.ID+" awaiting "+other.ID+"'s agent, "+owed.ID+" awaiting "+other.ID+"'s agent") {
		t.Errorf("once replied, why names both conversations: %q", why)
	}
}
