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
	"github.com/bytepunx/system-flow/flai/internal/threads"
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
		}
		WIPLimits map[string]int `json:"wip_limits"`
		Order     []string
	}
	if err := call(t, p, "board.get", `{"all":true}`, &v); err != nil {
		t.Fatal(err)
	}
	ready := v.Columns["ready"]
	if len(ready) != 1 || ready[0].ID != "S-0001" || ready[0].Status != "ready" || ready[0].ParentTitle != "Quay" ||
		ready[0].EnteredAt != "2026-09-20T08:01:00Z" || ready[0].AgeSecs != 59*60 {
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

	var got ItemWithChildren
	for _, id := range []string{"S-0001", "s-1", "S-01"} {
		if err := call(t, p, "item.get", `{"id":"`+id+`"}`, &got); err != nil || got.Item.ID != "S-0001" || len(got.Children) != 1 || got.Children[0].ID != "T-0001" || got.Item.Body == "" {
			t.Errorf("item.get %s: %+v %+v", id, err, got)
		}
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
