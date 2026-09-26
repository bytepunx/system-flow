package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// syncProject is a git repository with an epic, for stories whose branches
// sync checks against each other (S-0131). Integration: it runs real git.
func syncProject(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design", "docs"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("line one\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	return root
}

// openSyncStory creates story n in progress with touches and opens its
// worktree, returning the worktree's path.
func openSyncStory(t *testing.T, root string, n int, touches string) string {
	t.Helper()
	id := fmt.Sprintf("S-%04d", n)
	title := fmt.Sprintf("Story %d", n)
	if _, errOut, code := runIn(t, root, "story", "new", title, "--epic", "E-0001", "--touches", touches); code != 0 {
		t.Fatal(errOut)
	}
	file := filepath.Join(root, "wip/kanban/stories", fmt.Sprintf("%s-story-%d.md", id, n))
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] done\n", 1)), 0o644)
	for _, st := range []string{"ready", "in-progress"} {
		if _, errOut, code := runIn(t, root, "move", id, st); code != 0 {
			t.Fatal(errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "stream", "open", id); code != 0 {
		t.Fatal(errOut)
	}
	return filepath.Join(root, ".flai-cache", "worktrees", id)
}

// commitIn writes content to rel in the worktree and commits it.
func commitIn(t *testing.T, wt, rel, content string) {
	t.Helper()
	p := filepath.Join(wt, rel)
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "change "+rel)
}

func TestSyncTrialMergesOtherOpenBranches(t *testing.T) {
	root := syncProject(t)
	one := openSyncStory(t, root, 1, "docs")
	two := openSyncStory(t, root, 2, "docs/guide.md")
	three := openSyncStory(t, root, 3, "design")
	commitIn(t, one, "docs/guide.md", "one's line\n")
	commitIn(t, two, "docs/guide.md", "two's line\n")
	commitIn(t, three, "design/note.md", "three\n")

	out, errOut, code := runIn(t, one, "stream", "sync", "S-0001")
	if code != 0 {
		t.Fatalf("sync: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{
		"story/S-0001 is rebased onto main",
		"story/S-0001 conflicts with story/S-0002 (in progress) in docs/guide.md",
		"story/S-0001 merges cleanly with story/S-0003 (in progress)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sync output lacks %q:\n%s", want, out)
		}
	}
	// the trial merge writes nothing to any worktree
	for _, wt := range []string{one, two, three} {
		if st := gitIn(t, wt, "status", "--porcelain"); st != "" {
			t.Errorf("%s is not clean after sync: %s", wt, st)
		}
	}
	if got := gitIn(t, two, "show", "HEAD:docs/guide.md"); got != "two's line" {
		t.Errorf("story/S-0002 changed: %q", got)
	}

	out, errOut, code = runIn(t, one, "--json", "stream", "sync", "S-0001")
	if code != 0 {
		t.Fatalf("sync --json: %d %s %s", code, out, errOut)
	}
	var res struct {
		OK       bool          `json:"ok"`
		Branches []branchCheck `json:"branches"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if !res.OK || len(res.Branches) != 2 {
		t.Fatalf("branches: %+v", res)
	}
	if b := res.Branches[0]; b.Story != "S-0002" || b.Clean || strings.Join(b.Conflicts, ",") != "docs/guide.md" || b.Branch != "story/S-0002" || b.Status != "in-progress" {
		t.Errorf("conflicting pair: %+v", b)
	}
	if b := res.Branches[1]; b.Story != "S-0003" || !b.Clean || len(b.Conflicts) != 0 {
		t.Errorf("clean pair: %+v", b)
	}
}

// agentInbox is the MCP inbox of agent on the project at root.
func agentInbox(t *testing.T, root, agent string) mcpserver.InboxOut {
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
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "inbox", Arguments: map[string]any{}})
	if err != nil || res.IsError {
		t.Fatalf("inbox: %v %+v", err, res)
	}
	var out mcpserver.InboxOut
	data, _ := json.Marshal(res.StructuredContent)
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// designerInbox is the dashboard's inbox for the designer.
func designerInbox(t *testing.T, root string) hostapi.DesignerInbox {
	t.Helper()
	m := hostapi.Methods("test", time.Now)["inbox.designer"]
	res, rerr := m(context.Background(), channel.Project{Key: "t", Root: root}, json.RawMessage(`{}`))
	if rerr != nil {
		t.Fatalf("inbox.designer: %+v", rerr)
	}
	data, _ := json.Marshal(res)
	var out hostapi.DesignerInbox
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestSyncConflictIsAThreadInEveryInbox(t *testing.T) {
	root := syncProject(t)
	one := openSyncStory(t, root, 1, "docs")
	two := openSyncStory(t, root, 2, "docs/guide.md")
	commitIn(t, one, "docs/guide.md", "one's line\n")
	commitIn(t, two, "docs/guide.md", "two's line\n")
	at := time.Date(2026, 9, 26, 10, 0, 0, 0, time.UTC)
	sync := func(wt, id string, minutes int) string {
		t.Helper()
		out, errOut, code := runInAt(t, wt, at.Add(time.Duration(minutes)*time.Minute), "stream", "sync", id)
		if code != 0 {
			t.Fatalf("sync %s: %d %s %s", id, code, out, errOut)
		}
		return out
	}
	thread := func() *threads.Thread {
		t.Helper()
		repo, _ := workitem.Open(root)
		th, err := threads.Get(repo, "TH-0001")
		if err != nil {
			t.Fatal(err)
		}
		return th
	}

	out := sync(one, "S-0001", 1)
	if !strings.Contains(out, "conflicts with story/S-0002 (in progress) in docs/guide.md; see TH-0001") {
		t.Fatalf("sync names no thread:\n%s", out)
	}
	th := thread()
	if th.Title != "S-0001 and S-0002 conflict when merged" || th.Anchor.Item != "S-0001" || th.Opener() != "flai" || !th.Open() {
		t.Fatalf("thread: %+v", th)
	}
	if e := th.Entries(); len(e) != 1 || !strings.Contains(e[0].Text, "story/S-0001 with story/S-0002") || !strings.Contains(e[0].Text, "- `docs/guide.md`") {
		t.Fatalf("entry: %+v", e)
	}
	// both stories' agents and the designer see it
	for _, agent := range []string{"agent-S-0001", "agent-S-0002"} {
		in := agentInbox(t, root, agent)
		found := false
		for _, s := range in.Threads {
			found = found || (s.ID == "TH-0001" && s.Awaiting == "you")
		}
		if !found {
			t.Errorf("%s's inbox lacks TH-0001: %+v", agent, in.Threads)
		}
	}
	found := false
	for _, e := range designerInbox(t, root).Entries {
		found = found || e.Key == "thread:TH-0001"
	}
	if !found {
		t.Error("the designer's inbox lacks TH-0001")
	}
	if n, _ := os.ReadFile(filepath.Join(root, "wip/agents/S-0001.md")); !strings.Contains(string(n), "TH-0001") {
		t.Error("S-0001's narrative does not mirror TH-0001")
	}

	// the other story's sync finds the same paths and writes nothing
	if out := sync(two, "S-0002", 2); !strings.Contains(out, "conflicts with story/S-0001 (in progress) in docs/guide.md; see TH-0001") {
		t.Fatalf("second sync:\n%s", out)
	}
	if n := len(thread().Entries()); n != 1 {
		t.Errorf("same paths wrote %d entries", n)
	}

	// new conflicting paths add an entry to the same thread
	commitIn(t, one, "docs/more.md", "one\n")
	commitIn(t, two, "docs/more.md", "two\n")
	sync(one, "S-0001", 3)
	if e := thread().Entries(); len(e) != 2 || !strings.Contains(e[1].Text, "- `docs/guide.md`\n- `docs/more.md`") {
		t.Fatalf("new paths: %+v", e)
	}

	// once the two merge cleanly, sync resolves it
	commitIn(t, two, "docs/guide.md", "one's line\n")
	commitIn(t, two, "docs/more.md", "one\n")
	if out := sync(one, "S-0001", 4); !strings.Contains(out, "merges cleanly with story/S-0002") {
		t.Fatalf("clean:\n%s", out)
	}
	if th := thread(); th.Open() || !strings.Contains(th.Entries()[2].Text, "merge cleanly at the sync of S-0001") {
		t.Fatalf("not resolved: %+v %+v", th, th.Entries())
	}

	// a pair's thread is resolved when the other story is no longer open
	commitIn(t, two, "docs/guide.md", "two again\n")
	sync(one, "S-0001", 5)
	repo, _ := workitem.Open(root)
	th2, err := threads.Get(repo, "TH-0002")
	if err != nil || !th2.Open() {
		t.Fatalf("a new conflict opens a new thread: %v %+v", err, th2)
	}
	if _, errOut, code := runInAt(t, root, at.Add(6*time.Minute), "move", "S-0002", "cancelled", "--reason", "dropped"); code != 0 {
		t.Fatal(errOut)
	}
	sync(one, "S-0001", 7)
	th2, _ = threads.Get(repo, "TH-0002")
	if th2.Open() || !strings.Contains(th2.Entries()[len(th2.Entries())-1].Text, "S-0002 is cancelled, no longer open") {
		t.Fatalf("not resolved when S-0002 closed: %+v", th2.Entries())
	}
}
