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

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// makeProject writes a minimal system-flow project with key at dir.
func makeProject(t *testing.T, dir, key string) *workitem.Repo {
	t.Helper()
	for _, d := range []string{"design/system", "docs/users", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	m := "version: 1\nname: " + key + "\nkey: " + key + "\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(dir, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic of " + key, Owner: "alex", Now: t0}); err != nil {
		t.Fatal(err)
	}
	return repo
}

// readyStoryIn makes a story with a criterion and moves it to ready.
func readyStoryIn(t *testing.T, repo *workitem.Repo, title string, at time.Time) *workitem.Item {
	t.Helper()
	s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Owner: "alex", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(s.Path)
	_ = os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	s, _ = repo.Get(s.ID)
	if _, err := repo.Transition(s, workitem.Ready, "alex", "", at); err != nil {
		t.Fatal(err)
	}
	return s
}

type folderFixture struct {
	root string
	cs   *mcp.ClientSession
}

func folderSetup(t *testing.T, root string) *folderFixture {
	t.Helper()
	clock := t0.Add(time.Minute)
	srv := New(Options{Folder: root, Agent: "claude", Version: "test", Now: func() time.Time { return clock }, Poll: 20 * time.Millisecond, MaxWait: 3 * time.Second, Rescan: time.Nanosecond})
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
	return &folderFixture{root: root, cs: cs}
}

func (f *folderFixture) call(t *testing.T, name string, args any) (out map[string]any, failed string) {
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

func projectKeys(out map[string]any) []string {
	var keys []string
	for _, p := range out["projects"].([]any) {
		keys = append(keys, p.(map[string]any)["project"].(string))
	}
	return keys
}

// S-0101: started in a folder that is not a project, the server serves every
// project in it and below it, and nothing else: not a hidden folder (story
// worktrees live in .flai-cache), not a folder with no manifest.
func TestAFolderServesEveryProjectBelowIt(t *testing.T) {
	root := t.TempDir()
	alpha := makeProject(t, filepath.Join(root, "alpha"), "alpha")
	makeProject(t, filepath.Join(root, "org", "beta"), "beta")
	makeProject(t, filepath.Join(root, "alpha-copy", ".flai-cache", "worktrees", "x"), "hidden")
	_ = os.MkdirAll(filepath.Join(root, "plain"), 0o755)
	readyStoryIn(t, alpha, "Alpha work", t0)
	f := folderSetup(t, root)

	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, tool := range res.Tools {
		names = append(names, tool.Name)
	}
	sort.Strings(names)
	if strings.Join(names, " ") != "board doc_get inbox item_edit item_get item_move item_new thread_get thread_open thread_reply thread_resolve wait_for_events wait_for_work who_touches" {
		t.Errorf("tools: %v", names)
	}

	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" || strings.Join(projectKeys(out), " ") != "alpha beta" {
		t.Fatalf("inbox: %v %s", out, failed)
	}
	first := out["projects"].([]any)[0].(map[string]any)
	if first["folder"] != "alpha" || len(first["ready"].([]any)) != 1 {
		t.Errorf("alpha's entry: %v", first)
	}
	if out, _ := f.call(t, "inbox", map[string]any{"project": "org/beta"}); strings.Join(projectKeys(out), " ") != "beta" {
		t.Errorf("one project by folder: %v", out)
	}
	if _, failed := f.call(t, "inbox", map[string]any{"story": "S-0001"}); !strings.Contains(failed, "name the project") {
		t.Errorf("a story without its project: %q", failed)
	}

	// a per-project tool needs the project when there are several
	if _, failed := f.call(t, "item_get", map[string]any{"id": "E-1"}); !strings.Contains(failed, "name one with project: alpha (alpha), beta (org/beta)") {
		t.Errorf("no project named: %q", failed)
	}
	if out, failed := f.call(t, "item_get", map[string]any{"id": "E-1", "project": "beta"}); failed != "" || out["title"] != "Epic of beta" {
		t.Errorf("beta's epic: %v %s", out, failed)
	}
	if _, failed := f.call(t, "board", map[string]any{"project": "gamma"}); !strings.Contains(failed, `no project "gamma" here`) {
		t.Errorf("an unknown project: %q", failed)
	}

	// a project that appears while the agent works joins
	makeProject(t, filepath.Join(root, "gamma"), "gamma")
	if out, _ := f.call(t, "inbox", map[string]any{}); strings.Join(projectKeys(out), " ") != "alpha beta gamma" {
		t.Errorf("after gamma arrived: %v", projectKeys(out))
	}
}

// wait_for_work looks across every project: it waits when nothing is ready
// anywhere and hands over a story made ready in any of them, with its project.
func TestWaitForWorkAcrossAFolder(t *testing.T) {
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	beta := makeProject(t, filepath.Join(root, "beta"), "beta")
	f := folderSetup(t, root)

	quiet, _ := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 1})
	if quiet["timed_out"] != true || quiet["waiting_for"] != "ready" || len(quiet["projects"].([]any)) != 2 {
		t.Fatalf("nothing ready: %v", quiet)
	}
	done := make(chan map[string]any, 1)
	go func() {
		out, _ := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 3})
		done <- out
	}()
	time.Sleep(150 * time.Millisecond)
	story := readyStoryIn(t, beta, "Beta work", t0.Add(2*time.Minute))
	select {
	case out := <-done:
		s, _ := out["story"].(map[string]any)
		if out["reason"] != "pull" || out["project"] != "beta" || out["folder"] != "beta" || s["id"] != story.ID {
			t.Errorf("handed over: %v", out)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_work never answered")
	}

	// the agent pulls it in that project
	if out, failed := f.call(t, "item_move", map[string]any{"project": "beta", "id": story.ID, "to": "in-progress"}); failed != "" || out["status"] != "in-progress" {
		t.Errorf("pull: %v %s", out, failed)
	}
}

// wait_for_events names the project each event happened in.
func TestWaitForEventsAcrossAFolder(t *testing.T) {
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	beta := makeProject(t, filepath.Join(root, "beta"), "beta")
	f := folderSetup(t, root)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	done := make(chan map[string]any, 1)
	go func() {
		out, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 3})
		done <- out
	}()
	time.Sleep(150 * time.Millisecond)
	readyStoryIn(t, beta, "Beta work", t0.Add(2*time.Minute))
	select {
	case out := <-done:
		events := out["events"].([]any)
		if out["timed_out"] == true || len(events) == 0 {
			t.Fatalf("wait: %v", out)
		}
		for _, e := range events {
			if e.(map[string]any)["project"] != "beta" {
				t.Errorf("event from elsewhere: %v", e)
			}
		}
		changed := out["changed"].([]any)
		if len(changed) == 0 || !strings.HasPrefix(changed[0].(string), "beta/") {
			t.Errorf("changed paths are relative to the folder: %v", changed)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_events never answered")
	}
}

// An empty folder is served too: nothing to show, and it says why.
func TestAnEmptyFolderIsServed(t *testing.T) {
	f := folderSetup(t, t.TempDir())
	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" || len(out["projects"].([]any)) != 0 {
		t.Errorf("inbox: %v %s", out, failed)
	}
	if _, failed := f.call(t, "item_get", map[string]any{"id": "S-1"}); !strings.Contains(failed, "there is no system-flow project in") {
		t.Errorf("item_get: %q", failed)
	}
	if out, _ := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 1}); out["waiting_for"] != "ready" {
		t.Errorf("wait_for_work: %v", out)
	}
}

// With one project, naming another is a mistake worth saying.
func TestASingleProjectServerRefusesAnotherProject(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "item_get", map[string]any{"id": f.story.ID, "project": "elsewhere"}); !strings.Contains(failed, "serves only t") {
		t.Errorf("another project: %q", failed)
	}
	if out, failed := f.call(t, "item_get", map[string]any{"id": f.story.ID, "project": "t"}); failed != "" || out["id"] != f.story.ID {
		t.Errorf("its own key: %v %s", out, failed)
	}
}
