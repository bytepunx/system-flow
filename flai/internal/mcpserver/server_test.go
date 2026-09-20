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
	srv := New(Options{Repo: repo, Agent: "claude", Version: "test", Now: func() time.Time { return clock }, Poll: 20 * time.Millisecond, MaxWait: 3 * time.Second})
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
	want := "board doc_get inbox item_get item_move thread_get thread_open thread_reply thread_resolve wait_for_events who_touches"
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

// readyStory creates a story with criteria and moves it to ready as the
// designer at the given time.
func (f *fixture) readyStory(t *testing.T, title string, at time.Time) *workitem.Item {
	t.Helper()
	s, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: f.story.Parent, Owner: "alex", Now: at})
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
	if n := len(all["columns"].(map[string]any)["backlog"].([]any)); n < 2 {
		t.Errorf("all adds the epic and the task: %d", n)
	}
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
		if got := reordered(c.before, c.after); got != c.want {
			t.Errorf("reordered(%v, %v) = %v, want %v", c.before, c.after, got, c.want)
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
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_events never returned")
	}
}
