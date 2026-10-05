package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// I-0059: two stories in progress pass the pull-time hold with disjoint
// touches, then a task of one gains a touch inside the other's claim. Each
// write of touches, flai task new, flai edit, flai touches, and the MCP
// item_edit, names the other story and the paths, and both stories are told
// once in the inbox as an overlapped change caused by the other. A write that
// adds no path inside another claim tells nobody.
func TestATouchInsideAnotherClaimTellsBothStories(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	// stamped with the wall clock, which the MCP server reads
	run := func(args ...string) string {
		t.Helper()
		out, errOut, code := runInAt(t, root, time.Now().UTC(), args...)
		if code != 0 {
			t.Fatalf("flai %v: %d %s %s", args, code, out, errOut)
		}
		return out
	}
	run("epic", "new", "Epic")
	run("story", "new", "Mine", "--epic", "E-0001", "--touches", "flai/cmd")
	run("story", "new", "Theirs", "--epic", "E-0001", "--touches", "docs")
	for _, f := range []string{"S-0001-mine.md", "S-0002-theirs.md"} {
		p := filepath.Join(root, "wip/kanban/stories", f)
		s, _ := os.ReadFile(p)
		_ = os.WriteFile(p, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"S-0001", "S-0002"} {
		run("move", id, "ready")
		run("move", id, "in-progress")
		// a narrative, without the worktree flai stream open would make
		st, err := repo.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := repo.OpenStream(st, workitem.StreamOptions{Agent: "tester", Now: time.Now()}); err != nil {
			t.Fatal(err)
		}
	}
	agentInbox(t, root, "watcher") // the first look
	toldOnce := func(paths string) {
		t.Helper()
		got := overlappedIn(t, agentInbox(t, root, "watcher"))
		mine, theirs := got["S-0001"], got["S-0002"]
		if len(got) != 2 || len(mine) != 1 || len(theirs) != 1 {
			t.Fatalf("one overlapped change for each story: %+v", got)
		}
		if mine[0].Cause != "S-0002" || theirs[0].Cause != "S-0001" || mine[0].To != paths || theirs[0].To != paths {
			t.Errorf("the changes: %+v / %+v", mine[0], theirs[0])
		}
		if !strings.Contains(theirs[0].Summary, "S-0001's claim grew to overlap S-0002's on "+strings.ReplaceAll(paths, ",", ", ")) {
			t.Errorf("summary: %s", theirs[0].Summary)
		}
		if again := overlappedIn(t, agentInbox(t, root, "watcher")); len(again) != 0 {
			t.Errorf("told once: %+v", again)
		}
	}

	out := run("task", "new", "Reach", "--story", "S-0001", "--touches", "flai/cmd/x.go,docs/guide.md")
	if !strings.Contains(out, "overlaps S-0002 Theirs (in progress) on docs/guide.md") {
		t.Errorf("task new names the other story and the paths:\n%s", out)
	}
	toldOnce("docs/guide.md")

	ed := agentCall(t, root, "claude", "item_edit", map[string]any{"id": "T-0001", "touches": []string{"flai/cmd/x.go", "docs/guide.md", "docs/faq.md"}})
	if got := overlapsIn(t, ed); !reflect.DeepEqual(got, []itemedit.Overlapping{{Story: "S-0002", Title: "Theirs", Paths: []string{"docs/faq.md"}}}) {
		t.Errorf("item_edit: %+v", got)
	}
	toldOnce("docs/faq.md")

	var created struct {
		ID       string                 `json:"id"`
		Overlaps []itemedit.Overlapping `json:"overlaps"`
	}
	if err := json.Unmarshal([]byte(run("task", "new", "More", "--story", "S-0001", "--touches", "docs/more.md", "--json")), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID != "T-0002" || !reflect.DeepEqual(created.Overlaps, []itemedit.Overlapping{{Story: "S-0002", Title: "Theirs", Paths: []string{"docs/more.md"}}}) {
		t.Errorf("task new --json: %+v", created)
	}
	if out := run("edit", "T-0002", "--touches", "docs/more.md,docs/edit.md"); !strings.Contains(out, "overlaps S-0002 Theirs (in progress) on docs/edit.md") {
		t.Errorf("edit names the other story and the paths:\n%s", out)
	}
	if out := run("touches", "T-0002", "docs/more.md", "docs/edit.md", "docs/touch.md"); !strings.Contains(out, "overlaps S-0002 Theirs (in progress) on docs/touch.md") {
		t.Errorf("touches names the other story and the paths:\n%s", out)
	}
	var touched struct {
		Overlaps []itemedit.Overlapping `json:"overlaps"`
	}
	if err := json.Unmarshal([]byte(run("edit", "S-0001", "--touches", "flai/cmd,docs/story.md", "--json")), &touched); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(touched.Overlaps, []itemedit.Overlapping{{Story: "S-0002", Title: "Theirs", Paths: []string{"docs/story.md"}}}) {
		t.Errorf("an edit of the story in progress stands and is told: %+v", touched)
	}
	got := overlappedIn(t, agentInbox(t, root, "watcher"))
	if len(got["S-0001"]) != 4 || len(got["S-0002"]) != 4 {
		t.Errorf("four writes, four changes for each story: %+v", got)
	}

	// paths already in the claim, or outside every other claim, tell nobody
	for _, args := range [][]string{
		{"task", "new", "Apart", "--story", "S-0001", "--touches", "flai/cmd/y.go"},
		{"touches", "T-0002", "docs/more.md"},
		{"edit", "T-0001", "--touches", "flai/cmd/x.go,docs/guide.md,docs/faq.md,design"},
	} {
		if out := run(args...); strings.Contains(out, "overlaps") {
			t.Errorf("flai %v:\n%s", args, out)
		}
	}
	if ed := agentCall(t, root, "claude", "item_edit", map[string]any{"id": "T-0001", "touches": []string{"flai/cmd/z.go"}}); ed["overlaps"] != nil {
		t.Errorf("item_edit apart: %v", ed)
	}
	if got := overlappedIn(t, agentInbox(t, root, "watcher")); len(got) != 0 {
		t.Errorf("nothing grew into another claim: %+v", got)
	}
}

// overlappedIn is the overlapped changes of an inbox by the story told.
func overlappedIn(t *testing.T, in mcpserver.InboxOut) map[string][]mcpserver.Event {
	t.Helper()
	got := map[string][]mcpserver.Event{}
	for _, c := range in.Changes {
		if c.Kind == mcpserver.Overlapped {
			got[c.ID] = append(got[c.ID], c)
		}
	}
	return got
}

// agentCall calls an MCP tool as agent on the project at root and returns its
// structured result.
func agentCall(t *testing.T, root, agent, tool string, args map[string]any) map[string]any {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct, st := mcp.NewInMemoryTransports()
	if _, err := mcpserver.New(mcpserver.Options{Repo: repo, Agent: agent, Version: "test"}).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: agent, Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: args})
	if err != nil || res.IsError {
		t.Fatalf("%s: %v %+v", tool, err, res)
	}
	var out map[string]any
	data, _ := json.Marshal(res.StructuredContent)
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// overlapsIn is the overlaps field of a tool's result.
func overlapsIn(t *testing.T, out map[string]any) []itemedit.Overlapping {
	t.Helper()
	var got struct {
		Overlaps []itemedit.Overlapping `json:"overlaps"`
	}
	data, _ := json.Marshal(out)
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	return got.Overlaps
}
