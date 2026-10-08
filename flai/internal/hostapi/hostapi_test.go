package hostapi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 20, 8, 0, 0, 0, time.UTC)

// harbour is an epic with a ready story (one task), a backlog story, and an
// archived done story with its task; one open thread on the ready story.
func harbour(t *testing.T) channel.Project {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: Harbour\nkey: harbour\nowner: olive\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  notify_url: https://example.test/hook\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "wip/threads"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	mk := func(typ, title, parent string) *workitem.Item {
		it, err := repo.Create(workitem.NewOptions{Type: typ, Title: title, Parent: parent, Owner: "olive", Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		if typ == workitem.Story {
			data, _ := os.ReadFile(it.Path)
			_ = os.WriteFile(it.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
			it, _ = repo.Get(it.ID)
		}
		return it
	}
	move := func(it *workitem.Item, states ...string) {
		for i, to := range states {
			if _, err := repo.Transition(it, to, "olive", "", t0.Add(time.Duration(i+1)*time.Minute)); err != nil {
				t.Fatalf("%s to %s: %v", it.ID, to, err)
			}
		}
	}
	mk(workitem.Epic, "Quay", "")
	berths := mk(workitem.Story, "Berths", "E-0001")
	mk(workitem.Task, "Dredge", berths.ID)
	move(berths, workitem.Ready)
	mk(workitem.Story, "Cranes", "E-0001")
	old := mk(workitem.Story, "Survey", "E-0001")
	oldTask := mk(workitem.Task, "Chart", old.ID)
	move(oldTask, workitem.Ready, workitem.InProgress, workitem.Done)
	move(old, workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done)
	all, _ := repo.List(false)
	plan, err := repo.PlanArchive(all, []string{old.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Archive(plan); err != nil {
		t.Fatal(err)
	}
	if _, err := threads.New(repo, threads.NewOptions{On: berths.ID, Title: "How deep?", Author: "olive", Text: "Eight metres?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	return channel.Project{Key: "harbour", Name: "Harbour", Root: root}
}

func call(t *testing.T, p channel.Project, method, params string, into any) *channel.Error {
	t.Helper()
	m, ok := Methods("test", func() time.Time { return t0.Add(time.Hour) })[method]
	if !ok {
		t.Fatalf("no method %s", method)
	}
	res, rerr := m(context.Background(), p, json.RawMessage(params))
	if rerr != nil {
		return rerr
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, into); err != nil {
		t.Fatalf("%s: %v\n%s", method, err, data)
	}
	if strings.Contains(string(data), p.Root) {
		t.Errorf("%s names a path on the host:\n%s", method, data)
	}
	return nil
}

func TestProjectInfo(t *testing.T) {
	var info ProjectInfo
	if err := call(t, harbour(t), "project.info", `{}`, &info); err != nil {
		t.Fatal(err)
	}
	if info.Name != "Harbour" || info.Key != "harbour" || info.Owner != "olive" || info.Layout["wip"] != "wip" ||
		info.Dashboard.NotifyURL != "https://example.test/hook" || !info.Dashboard.Autocommit || info.Flai != "test" || info.Version != 1 {
		t.Errorf("%+v", info)
	}
}

func TestBoardGet(t *testing.T) {
	p := harbour(t)
	var v struct {
		Columns map[string][]struct {
			ID, Type, Status, Parent string
			ParentTitle              string `json:"parent_title"`
			EnteredAt                string `json:"entered_at"`
			AgeSecs                  int64  `json:"age_in_column_seconds"`
			Tasks                    *workitem.TaskSummary
		}
		WIPLimits map[string]int `json:"wip_limits"`
		Order     []string
	}
	if err := call(t, p, "board.get", `{"all":true}`, &v); err != nil {
		t.Fatal(err)
	}
	ready := v.Columns["ready"]
	if len(ready) != 1 || ready[0].ID != "S-0001" || ready[0].Status != "ready" || ready[0].ParentTitle != "Quay" ||
		ready[0].EnteredAt != "2026-09-20T08:01:00Z" || ready[0].AgeSecs != 59*60 ||
		ready[0].Tasks == nil || *ready[0].Tasks != (workitem.TaskSummary{Ready: 1, Layers: 1}) {
		t.Errorf("ready: %+v", ready)
	}
	if len(v.Order) != 1 || v.Order[0] != "S-0001" || v.WIPLimits["ready"] == 0 {
		t.Errorf("order %v limits %v", v.Order, v.WIPLimits)
	}
	types := map[string]int{}
	for _, col := range v.Columns {
		for _, c := range col {
			types[c.Type]++
			if c.ID == "S-0003" || c.ID == "T-0002" {
				t.Errorf("an archived item is on the board: %s", c.ID)
			}
		}
	}
	if types["epic"] != 1 || types["story"] != 2 || types["task"] != 1 {
		t.Errorf("all: %v", types)
	}
	var stories struct {
		Columns map[string][]struct{ Type string }
	}
	_ = call(t, p, "board.get", `{}`, &stories)
	for _, col := range stories.Columns {
		for _, c := range col {
			if c.Type != "story" {
				t.Errorf("without all the board holds a %s", c.Type)
			}
		}
	}
}

func TestItemsListAndGet(t *testing.T) {
	p := harbour(t)
	var items []workitem.Item
	if err := call(t, p, "items.list", `{}`, &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 4 || items[0].Body != "" || !strings.HasPrefix(items[0].Path, "wip/kanban/") {
		t.Errorf("active items without bodies: %d %+v", len(items), items[0])
	}
	_ = call(t, p, "items.list", `{"archived":true,"bodies":true,"type":"story"}`, &items)
	archived := 0
	for _, it := range items {
		if it.Type != "story" || it.Body == "" {
			t.Errorf("%s: type %s body %d", it.ID, it.Type, len(it.Body))
		}
		if it.Archived {
			archived++
			if !strings.HasPrefix(it.Path, "wip/archive/") {
				t.Errorf("archived path: %s", it.Path)
			}
		}
	}
	if len(items) != 3 || archived != 1 {
		t.Errorf("stories with the archive: %d, archived %d", len(items), archived)
	}
	_ = call(t, p, "items.list", `{"status":"ready"}`, &items)
	if len(items) != 1 || items[0].ID != "S-0001" {
		t.Errorf("ready: %+v", items)
	}
	for params, want := range map[string]string{`{"type":"folder"}`: "type must be", `{"status":"finished"}`: "not a state", `{"type":7}`: "params"} {
		if err := call(t, p, "items.list", params, &items); err == nil || err.Code != channel.CodeInvalidParams || !strings.Contains(err.Message, want) {
			t.Errorf("%s: %+v", params, err)
		}
	}

	var n ItemCount
	if err := call(t, p, "items.count", `{}`, &n); err != nil || n.Active != 4 || n.Archived != 2 {
		t.Errorf("items.count: %+v %+v", err, n)
	}
	if err := call(t, p, "items.count", `[]`, &n); err == nil || err.Code != channel.CodeInvalidParams {
		t.Errorf("items.count with params that are not an object: %+v", err)
	}

	var got ItemWithChildren
	for _, id := range []string{"S-0001", "s-1", "S-01"} {
		if err := call(t, p, "item.get", `{"id":"`+id+`"}`, &got); err != nil || got.Item.ID != "S-0001" || len(got.Children) != 1 || got.Children[0].ID != "T-0001" || got.Item.Body == "" {
			t.Errorf("item.get %s: %+v %+v", id, err, got)
		}
	}
	// S-0176: a story with tasks has its task plan; anything else has none
	if got.Plan == nil || len(got.Plan.Tasks) != 1 || got.Plan.Tasks[0].State != workitem.PlanReady || len(got.Plan.Layers) != 1 {
		t.Errorf("item.get's plan: %+v", got.Plan)
	}
	got = ItemWithChildren{}
	if err := call(t, p, "item.get", `{"id":"T-0001"}`, &got); err != nil || got.Plan != nil {
		t.Errorf("a task has no plan: %+v %+v", err, got.Plan)
	}
	if err := call(t, p, "item.get", `{"id":"S-0003"}`, &got); err != nil || !got.Item.Archived || len(got.Children) != 1 {
		t.Errorf("an archived item and its task: %+v %+v", err, got)
	}
	if err := call(t, p, "item.get", `{"id":"S-0099"}`, &got); err == nil || err.Code != NotFound {
		t.Errorf("missing: %+v", err)
	}
	for _, id := range []string{"../../etc/passwd", "S-0001; rm -rf", "", "X-1"} {
		if err := call(t, p, "item.get", `{"id":"`+id+`"}`, &got); err == nil || err.Code != channel.CodeInvalidParams {
			t.Errorf("item.get %q: %+v", id, err)
		}
	}
}

func TestThreadsList(t *testing.T) {
	p := harbour(t)
	var list []struct {
		ID, Title, Status, Path, Story string
		Anchor                         struct{ Item string }
		Entries                        []struct{ At, Author, Text string }
	}
	if err := call(t, p, "threads.list", `{}`, &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "TH-0001" || list[0].Story != "S-0001" || !strings.HasPrefix(list[0].Path, "wip/threads/") ||
		len(list[0].Entries) != 1 || list[0].Entries[0].Author != "olive" || list[0].Entries[0].Text != "Eight metres?" {
		t.Errorf("threads: %+v", list)
	}
	_ = call(t, p, "threads.list", `{"on":"S-0002"}`, &list)
	if len(list) != 0 {
		t.Errorf("threads on another story: %+v", list)
	}
	_ = call(t, p, "threads.list", `{"on":"S-1"}`, &list)
	if len(list) != 1 {
		t.Errorf("threads on the story in short padding: %+v", list)
	}
}

// quiet is harbour with S-0001 in progress, S-0002 in review, a new S-0004 in
// progress, and a file to talk about: talking before its conversations.
func quiet(t *testing.T) channel.Project {
	t.Helper()
	p := harbour(t)
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Buoys", Parent: "E-0001", Owner: "olive", Now: t0}); err != nil {
		t.Fatal(err)
	}
	for id, state := range map[string]string{"S-0001": workitem.InProgress, "S-0002": workitem.Review, "S-0004": workitem.InProgress} {
		it, err := repo.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		it.Status = state
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(p.Root, "design/system"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p.Root, "design/system/quay.md"), []byte("# Quay\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// talking is quiet with three conversations: MS-0001 from S-0001 to S-0002
// about its file, answered, so it awaits S-0001; MS-0002 from S-0002 to
// S-0001, closed; and MS-0003 from S-0002 to S-0004, which awaits S-0004.
func talking(t *testing.T) channel.Project {
	t.Helper()
	p := quiet(t)
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	send := func(from, to, text string, about ...string) *messages.Conversation {
		c, err := messages.Send(repo, messages.SendOptions{From: from, To: to, Author: "agent-" + from, Text: text, About: about, Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	first := send("S-0001", "S-0002", "Which berth?", "design/system/quay.md")
	if _, err := messages.Reply(repo, first.ID, "S-0002", "agent-S-0002", "The north one.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	second := send("S-0002", "S-0001", "Done with the crane?")
	if _, err := messages.Close(repo, second.ID, "olive", "settled on the quay", t0.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	send("S-0002", "S-0004", "Buoy colours?")
	return p
}

type conversationView struct {
	ID, Title, From, To, Status, Awaiting, Path string
	About                                       []string
	Closed                                      bool
	ClosedReason                                string `json:"closed_reason"`
	Entries                                     []messages.Entry
}

func conversationIDs(list []conversationView) string {
	var out []string
	for _, c := range list {
		out = append(out, c.ID)
	}
	return strings.Join(out, " ")
}

// TestMessagesList: messages.list serves the open conversations of the
// project, or of one story named in any padding, and the closed ones too with
// all.
func TestMessagesList(t *testing.T) {
	p := talking(t)
	for params, want := range map[string]string{
		`{}`:                            "MS-0001 MS-0003",
		`{"all":true}`:                  "MS-0001 MS-0002 MS-0003",
		`{"story":"S-0001"}`:            "MS-0001",
		`{"story":"s-1","all":true}`:    "MS-0001 MS-0002",
		`{"story":"S-04"}`:              "MS-0003",
		`{"story":"S-0003","all":true}`: "",
	} {
		var list []conversationView
		if err := call(t, p, "messages.list", params, &list); err != nil {
			t.Fatalf("%s: %+v", params, err)
		}
		if list == nil {
			t.Errorf("%s: null, not a list", params)
		}
		if got := conversationIDs(list); got != want {
			t.Errorf("%s: %q, want %q", params, got, want)
		}
	}
	var list []conversationView
	_ = call(t, p, "messages.list", `{"all":true}`, &list)
	if c := list[1]; !c.Closed || c.ClosedReason != "settled on the quay" || c.Awaiting != "" || c.Status != "closed" {
		t.Errorf("a closed conversation: %+v", c)
	}
	if c := list[0]; c.Closed || c.Awaiting != "S-0001" || len(c.About) != 1 || c.About[0] != "design/system/quay.md" || len(c.Entries) != 2 {
		t.Errorf("an open conversation: %+v", c)
	}
	for _, params := range []string{`{"story":"../../etc"}`, `{"story":"T-0001"}`, `{"story":7}`} {
		if err := call(t, p, "messages.list", params, &list); err == nil || err.Code != channel.CodeInvalidParams {
			t.Errorf("messages.list %s: %+v", params, err)
		}
	}
}

// TestMessagesGet: messages.get serves one conversation in any padding with
// its entries, its about paths, and its state, and refuses an ID that is
// missing, malformed, or of no conversation.
func TestMessagesGet(t *testing.T) {
	p := talking(t)
	var c conversationView
	for _, id := range []string{"MS-0001", "ms-1", "1"} {
		c = conversationView{}
		if err := call(t, p, "messages.get", `{"id":"`+id+`"}`, &c); err != nil || c.ID != "MS-0001" {
			t.Fatalf("messages.get %s: %+v %+v", id, err, c)
		}
	}
	if c.From != "S-0001" || c.To != "S-0002" || c.Status != "open" || c.Closed || c.Awaiting != "S-0001" ||
		!strings.HasPrefix(c.Path, "wip/messages/MS-0001-") || len(c.About) != 1 || c.About[0] != "design/system/quay.md" {
		t.Errorf("messages.get: %+v", c)
	}
	if len(c.Entries) != 2 || c.Entries[0].Story != "S-0001" || c.Entries[0].Text != "Which berth?" ||
		c.Entries[1].Author != "agent-S-0002" || c.Entries[1].Text != "The north one." {
		t.Errorf("entries: %+v", c.Entries)
	}
	c = conversationView{}
	if err := call(t, p, "messages.get", `{"id":"MS-0002"}`, &c); err != nil || !c.Closed || c.ClosedReason != "settled on the quay" || c.Awaiting != "" || c.About == nil {
		t.Errorf("a closed conversation: %+v %+v", err, c)
	}
	if err := call(t, p, "messages.get", `{"id":"MS-0099"}`, &c); err == nil || err.Code != NotFound {
		t.Errorf("an unknown conversation: %+v", err)
	}
	for _, params := range []string{`{}`, `{"id":""}`, `{"id":"../../etc/passwd"}`, `{"id":"MS-*"}`, `{"id":"TH-0001"}`, `{"id":7}`} {
		if err := call(t, p, "messages.get", params, &c); err == nil || err.Code != channel.CodeInvalidParams {
			t.Errorf("messages.get %s: %+v", params, err)
		}
	}
}

// TestInboxDesignerCountsNoMessage: the operator's inbox, and its badge, which
// is its total, count no conversation between stories (S-0336): with
// conversations open that await either side, it answers as it did on the same
// project before they were written, and none of its entries is a message.
func TestInboxDesignerCountsNoMessage(t *testing.T) {
	answer := func(p channel.Project) (DesignerInbox, string) {
		t.Helper()
		var in DesignerInbox
		if err := call(t, p, "inbox.designer", `{}`, &in); err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}
		return in, string(data)
	}
	before, beforeJSON := answer(quiet(t))
	after, afterJSON := answer(talking(t))
	// S-0002 in review: an inbox with something in it (olive wrote the thread last)
	if before.Total != 1 || before.Counts["review"] != 1 {
		t.Fatalf("before the conversations: %s", beforeJSON)
	}
	if afterJSON != beforeJSON {
		t.Errorf("the conversations changed the inbox\nbefore %s\nafter  %s", beforeJSON, afterJSON)
	}
	for _, e := range after.Entries {
		if _, counted := before.Counts[e.Kind]; !counted || strings.Contains(e.Key+e.Title+e.Path+e.Detail, "MS-") || strings.Contains(e.Path, "wip/messages") {
			t.Errorf("an entry from a conversation: %+v", e)
		}
	}
}

// TestItemGetExpectedCost: item.get prices an item's forecast, or else its
// estimate, at the project's cost per agent hour, and gives nothing before a
// story has been measured from its logs (ADR-0083).
func TestItemGetExpectedCost(t *testing.T) {
	p := harbour(t)
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	set := func(id string, edit func(*workitem.Item)) {
		it, err := repo.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		edit(it)
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	var got ItemWithChildren
	set("S-0001", func(it *workitem.Item) { it.Estimate = "2h" })
	if err := call(t, p, "item.get", `{"id":"S-0001"}`, &got); err != nil || got.ExpectedCost != nil {
		t.Errorf("no story measured, no expected cost: %+v %+v", err, got.ExpectedCost)
	}
	// $6 over half an hour of agent work: $12 an hour
	set("S-0003", func(it *workitem.Item) {
		it.Usage = &usage.Usage{Source: usage.SourceLog, Seconds: 1800, Models: []usage.Model{{Model: "m", Input: 10, Cost: 6}}}
	})
	got = ItemWithChildren{}
	if err := call(t, p, "item.get", `{"id":"S-0001"}`, &got); err != nil || got.ExpectedCost == nil ||
		got.ExpectedCost.Cost != 24 || got.ExpectedCost.From != "estimate" || !got.ExpectedCost.Estimated {
		t.Errorf("an estimate of 2h at $12 an hour: %+v %+v", err, got.ExpectedCost)
	}
	got = ItemWithChildren{}
	if err := call(t, p, "item.get", `{"id":"T-0001"}`, &got); err != nil || got.ExpectedCost != nil {
		t.Errorf("an item with no duration has no expected cost: %+v %+v", err, got.ExpectedCost)
	}
}
