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

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
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
	if strings.Contains(out, "outside") {
		t.Errorf("S-0001 changed only what it touches:\n%s", out)
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
		OK       bool                   `json:"ok"`
		Branches []storygit.BranchCheck `json:"branches"`
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

// ADR-0121: a conflict is a message from the syncing story to the other in
// the pair's conversation, which the sync names and --json gives as the
// branch's conversation; no thread is opened.
func TestSyncConflictIsAMessageBetweenThePair(t *testing.T) {
	root := syncProject(t)
	one := openSyncStory(t, root, 1, "docs")
	two := openSyncStory(t, root, 2, "docs/guide.md")
	commitIn(t, one, "docs/guide.md", "one's line\n")
	commitIn(t, two, "docs/guide.md", "two's line\n")
	at := time.Date(2026, 9, 26, 10, 0, 0, 0, time.UTC)
	sync := func(wt, id string, minutes int, args ...string) string {
		t.Helper()
		out, errOut, code := runInAt(t, wt, at.Add(time.Duration(minutes)*time.Minute), append(args, "stream", "sync", id)...)
		if code != 0 {
			t.Fatalf("sync %s: %d %s %s", id, code, out, errOut)
		}
		return out
	}
	conversation := func() *messages.Conversation {
		t.Helper()
		repo, _ := workitem.Open(root)
		c, err := messages.Get(repo, "MS-0001")
		if err != nil {
			t.Fatal(err)
		}
		return c
	}

	out := sync(one, "S-0001", 1)
	if !strings.Contains(out, "story/S-0001 conflicts with story/S-0002 (in progress) in docs/guide.md; see MS-0001\n") {
		t.Fatalf("sync names no conversation:\n%s", out)
	}
	c := conversation()
	if c.From != "S-0001" || c.To != "S-0002" || strings.Join(c.About, ",") != "docs/guide.md" || c.Status != "open" || c.Awaiting() != "S-0002" {
		t.Fatalf("conversation: %+v", c)
	}
	if e := c.Entries(); len(e) != 1 || e[0].Author != "flai" || !strings.Contains(e[0].Text, "story/S-0001 and story/S-0002 conflict when merged") || !strings.Contains(e[0].Text, "- `docs/guide.md`") {
		t.Fatalf("entry: %+v", e)
	}
	repo, _ := workitem.Open(root)
	if all, err := threads.List(repo); err != nil || len(all) != 0 {
		t.Errorf("a thread was opened: %v %+v", err, all)
	}
	// each story's agent finds it in its inbox, S-0002's to answer
	t.Setenv("FLAI_STORY", "")
	for agent, awaiting := range map[string]string{"agent-S-0001": "other", "agent-S-0002": "you"} {
		if in := agentInbox(t, root, agent); len(in.Messages) != 1 || in.Messages[0].ID != "MS-0001" || in.Messages[0].Awaiting != awaiting {
			t.Errorf("%s's inbox messages: %+v", agent, in.Messages)
		}
	}

	// the other story's sync finds the same paths and writes nothing
	if out := sync(two, "S-0002", 2); !strings.Contains(out, "story/S-0002 conflicts with story/S-0001 (in progress) in docs/guide.md; see MS-0001\n") {
		t.Fatalf("second sync:\n%s", out)
	}
	var res struct {
		Branches []storygit.BranchCheck `json:"branches"`
	}
	if err := json.Unmarshal([]byte(sync(two, "S-0002", 3, "--json")), &res); err != nil {
		t.Fatal(err)
	}
	if len(res.Branches) != 1 || res.Branches[0].Conversation != "MS-0001" || res.Branches[0].Thread != "" {
		t.Errorf("--json branches: %+v", res.Branches)
	}
	if n := len(conversation().Entries()); n != 1 {
		t.Errorf("same paths wrote %d entries", n)
	}

	// once the two merge cleanly, sync closes it
	commitIn(t, two, "docs/guide.md", "one's line\n")
	if out := sync(one, "S-0001", 4); !strings.Contains(out, "story/S-0001 merges cleanly with story/S-0002 (in progress)\n") {
		t.Fatalf("clean:\n%s", out)
	}
	if c := conversation(); c.Status != "closed" || !strings.Contains(c.Entries()[1].Text, "merge cleanly at the sync of S-0001") {
		t.Fatalf("not closed: %+v %+v", c, c.Entries())
	}
}

func TestSyncListsPathsChangedOutsideTheClaim(t *testing.T) {
	root := syncProject(t)
	wt := openSyncStory(t, root, 1, "docs/guide.md")
	if _, errOut, code := runIn(t, root, "task", "new", "Design it", "--story", "S-0001", "--touches", "design"); code != 0 {
		t.Fatal(errOut)
	}
	commitIn(t, wt, "docs/guide.md", "inside the story's touches\n")
	commitIn(t, wt, "design/note.md", "inside its task's touches\n")
	commitIn(t, wt, "docs/other.md", "outside\n")
	commitIn(t, wt, "README.md", "outside\n")
	commitIn(t, wt, "wip/stray.md", "flai's folder, left out\n")

	out, errOut, code := runIn(t, wt, "stream", "sync", "S-0001")
	if code != 0 {
		t.Fatalf("sync: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{
		"story/S-0001 changed 2 paths outside S-0001's touches: README.md, docs/other.md",
		"flai touches S-0001 --add README.md docs/other.md",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sync output lacks %q:\n%s", want, out)
		}
	}
	out, _, _ = runIn(t, wt, "--json", "stream", "sync", "S-0001")
	var res struct {
		Outside []string `json:"outside_touches"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || strings.Join(res.Outside, ",") != "README.md,docs/other.md" {
		t.Errorf("outside_touches: %v %s", err, out)
	}

	// widened as sync says, nothing is outside
	if _, errOut, code := runIn(t, root, "touches", "S-0001", "--add", "README.md", "docs/other.md"); code != 0 {
		t.Fatal(errOut)
	}
	if out, _, _ := runIn(t, wt, "stream", "sync", "S-0001"); strings.Contains(out, "outside") {
		t.Errorf("widened, still outside:\n%s", out)
	}
}

// ADR-0096: once a task names a file in the story's folder touch, the claim
// is narrowed to it, so a file the branch changed elsewhere in the folder is
// reported; a done task's file stays claimed, and widening the task's touches
// clears the report.
func TestSyncListsPathsChangedOutsideANarrowedClaim(t *testing.T) {
	root := syncProject(t)
	wt := openSyncStory(t, root, 1, "docs")
	if _, errOut, code := runIn(t, root, "task", "new", "Write the guide", "--story", "S-0001", "--touches", "docs/guide.md"); code != 0 {
		t.Fatal(errOut)
	}
	commitIn(t, wt, "docs/guide.md", "inside the task's touches\n")
	commitIn(t, wt, "docs/other.md", "in the story's folder, outside the narrowed claim\n")
	outside := func() string {
		t.Helper()
		out, errOut, code := runIn(t, wt, "--json", "stream", "sync", "S-0001")
		if code != 0 {
			t.Fatalf("sync: %d %s %s", code, out, errOut)
		}
		var res struct {
			Outside []string `json:"outside_touches"`
		}
		if err := json.Unmarshal([]byte(out), &res); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
		return strings.Join(res.Outside, ",")
	}

	out, errOut, code := runIn(t, wt, "stream", "sync", "S-0001")
	if code != 0 {
		t.Fatalf("sync: %d %s %s", code, out, errOut)
	}
	if want := "story/S-0001 changed 1 path outside S-0001's touches: docs/other.md"; !strings.Contains(out, want) {
		t.Errorf("sync output lacks %q:\n%s", want, out)
	}
	for _, st := range []string{"ready", "in-progress", "done"} {
		if _, errOut, code := runIn(t, root, "move", "T-0001", st); code != 0 {
			t.Fatal(errOut)
		}
	}
	if got := outside(); got != "docs/other.md" {
		t.Errorf("with the task done, outside_touches = %q, want docs/other.md", got)
	}

	// widened in the task's touches, as the ADR tells the agent, nothing is outside
	if _, errOut, code := runIn(t, root, "touches", "T-0001", "docs/guide.md", "docs/other.md"); code != 0 {
		t.Fatal(errOut)
	}
	if got := outside(); got != "" {
		t.Errorf("widened, still outside: %s", got)
	}
}

// ADR-0069: sync never stashes; it refuses a worktree with uncommitted
// changes, touching nothing, and names each path.
func TestSyncRefusesUncommittedChanges(t *testing.T) {
	root := syncProject(t)
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/guide.md", "branch's line\n")
	// main moves on, so a rebase would have something to do
	_ = os.WriteFile(filepath.Join(root, "README.md"), []byte("readme\n"), 0o644)
	gitIn(t, root, "add", "README.md")
	gitIn(t, root, "commit", "-q", "-m", "docs: readme")
	head := gitIn(t, wt, "rev-parse", "HEAD")
	_ = os.WriteFile(filepath.Join(wt, "docs", "guide.md"), []byte("half done\n"), 0o644)
	_ = os.WriteFile(filepath.Join(wt, "notes.txt"), []byte("untracked\n"), 0o644)

	out, errOut, code := runIn(t, wt, "stream", "sync", "S-0001")
	if code == 0 {
		t.Fatalf("sync over uncommitted changes succeeded: %s", out)
	}
	for _, want := range []string{
		"story/S-0001 was not synced: its worktree .flai-cache/worktrees/S-0001 has uncommitted changes in 2 paths:\n",
		"\n  docs/guide.md\n",
		"\n  notes.txt\n",
		"To continue: commit them on story/S-0001 (or stash them), then run flai stream sync S-0001 again\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("refusal lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "To abort") {
		t.Errorf("nothing to abort:\n%s", out)
	}
	if !strings.Contains(errOut, "uncommitted changes in docs/guide.md, notes.txt") {
		t.Errorf("error: %s", errOut)
	}

	out, _, code = runIn(t, wt, "--json", "stream", "sync", "S-0001")
	var res struct {
		OK          bool     `json:"ok"`
		Uncommitted []string `json:"uncommitted"`
		Conflicts   []string `json:"conflicts"`
		InProgress  bool     `json:"rebase_in_progress"`
		Continue    string   `json:"continue"`
		Abort       string   `json:"abort"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if code == 0 || res.OK || strings.Join(res.Uncommitted, ",") != "docs/guide.md,notes.txt" || res.Conflicts == nil || len(res.Conflicts) != 0 ||
		res.InProgress || !strings.Contains(res.Continue, "commit them on story/S-0001") || res.Abort != "" {
		t.Errorf("json: %d %+v", code, res)
	}

	// nothing was touched: no rebase, no stash, the work as it was
	if got := gitIn(t, wt, "rev-parse", "HEAD"); got != head {
		t.Errorf("HEAD moved: %s, was %s", got, head)
	}
	if got := gitIn(t, wt, "stash", "list"); got != "" {
		t.Errorf("a stash was made: %s", got)
	}
	for rel, want := range map[string]string{"docs/guide.md": "half done\n", "notes.txt": "untracked\n"} {
		if got, _ := os.ReadFile(filepath.Join(wt, rel)); string(got) != want {
			t.Errorf("%s is %q", rel, got)
		}
	}
}

// A sync that stops on conflicts lists each path and how to continue or
// abort, leaves the rebase in progress, and a sync while it is in progress
// is refused with the same paths and steps.
func TestSyncConflictListsPathsAndHowToContinueOrAbort(t *testing.T) {
	root := syncProject(t)
	wt := openSyncStory(t, root, 1, "docs")
	// one commit, so the rebase stops on both paths at once
	_ = os.WriteFile(filepath.Join(wt, "docs", "more.md"), []byte("branch's more\n"), 0o644)
	commitIn(t, wt, "docs/guide.md", "branch's line\n")
	head := gitIn(t, wt, "rev-parse", "HEAD")
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("main's line\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "more.md"), []byte("main's more\n"), 0o644)
	gitIn(t, root, "add", "docs")
	gitIn(t, root, "commit", "-q", "-m", "docs: on main")
	inProgress := func() bool { return storygit.RebaseInProgress(execx.System{}, wt) }
	steps := []string{
		"\n  docs/guide.md\n",
		"\n  docs/more.md\n",
		"To continue: in .flai-cache/worktrees/S-0001, resolve each conflicting path, git add it, and run git rebase --continue; then run flai stream sync S-0001 again\n",
		"To abort: in .flai-cache/worktrees/S-0001, run git rebase --abort, which puts story/S-0001 back as it was before the sync\n",
	}
	type result struct {
		OK          bool     `json:"ok"`
		Uncommitted []string `json:"uncommitted"`
		Conflicts   []string `json:"conflicts"`
		InProgress  bool     `json:"rebase_in_progress"`
		Continue    string   `json:"continue"`
		Abort       string   `json:"abort"`
	}
	syncJSON := func() {
		t.Helper()
		out, _, code := runIn(t, wt, "--json", "stream", "sync", "S-0001")
		var res result
		if err := json.Unmarshal([]byte(out), &res); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
		if code == 0 || res.OK || !res.InProgress || strings.Join(res.Conflicts, ",") != "docs/guide.md,docs/more.md" ||
			res.Uncommitted == nil || len(res.Uncommitted) != 0 ||
			!strings.Contains(res.Continue, "git rebase --continue") || !strings.Contains(res.Abort, "git rebase --abort") {
			t.Errorf("json: %d %+v", code, res)
		}
	}

	// the sync stops on the conflicts
	out, errOut, code := runIn(t, wt, "stream", "sync", "S-0001")
	if code == 0 {
		t.Fatalf("conflicting sync succeeded: %s", out)
	}
	for _, want := range append([]string{"story/S-0001 was not synced: the rebase onto main stopped on conflicts in 2 paths:\n"}, steps...) {
		if !strings.Contains(out, want) {
			t.Errorf("conflict report lacks %q:\n%s", want, out)
		}
	}
	if !strings.Contains(errOut, "stopped on conflicts in docs/guide.md, docs/more.md") {
		t.Errorf("error: %s", errOut)
	}
	if !inProgress() {
		t.Fatal("the rebase is not left in progress")
	}

	// syncing again while it is in progress is refused with the same steps
	out, errOut, code = runIn(t, wt, "stream", "sync", "S-0001")
	if code == 0 {
		t.Fatalf("sync during a rebase succeeded: %s", out)
	}
	for _, want := range append([]string{"story/S-0001 was not synced: a rebase is already in progress in .flai-cache/worktrees/S-0001, with conflicts in 2 paths:\n"}, steps...) {
		if !strings.Contains(out, want) {
			t.Errorf("in-progress report lacks %q:\n%s", want, out)
		}
	}
	if !strings.Contains(errOut, "already in progress") {
		t.Errorf("error: %s", errOut)
	}
	syncJSON()
	if !inProgress() {
		t.Fatal("the refused sync touched the rebase")
	}

	// abort as it says puts the branch back; a --json sync stops the same way
	gitIn(t, wt, "rebase", "--abort")
	if got := gitIn(t, wt, "rev-parse", "HEAD"); got != head {
		t.Errorf("abort left HEAD at %s, was %s", got, head)
	}
	syncJSON()
	if !inProgress() {
		t.Fatal("the --json sync did not stop on the conflicts")
	}
}

// I-0074: two stories that each record an issue both rewrite
// design/issues/summary.md, its updated line and its rows, so the rebase
// onto main after one is accepted stopped on it. The file is generated from
// the issue files, so the sync writes it again at each such stop and goes on
// (ADR-0098).
func TestSyncRegeneratesTheIssueSummaryWhenTheRebaseStopsOnItAlone(t *testing.T) {
	_, b := issueStories(t, nil)
	// a second commit that rewrites the summary stops the rebase again
	recordIssueIn(t, b, issueClock.Add(3*time.Hour), "B's second issue", nil)

	out, errOut, code := runInAt(t, b, issueClock.Add(4*time.Hour), "stream", "sync", "S-0002")
	if code != 0 {
		t.Fatalf("sync stopped on the generated summary: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(out, "story/S-0002 is rebased onto main") {
		t.Errorf("sync output:\n%s", out)
	}
	if storygit.RebaseInProgress(execx.System{}, b) {
		t.Fatal("the rebase is left in progress")
	}
	if st := gitIn(t, b, "status", "--porcelain"); st != "" {
		t.Errorf("the worktree is not clean after the sync: %s", st)
	}
	if got := gitIn(t, b, "log", "--format=%s", "main..HEAD"); got != "chore: record B's second issue\nchore: record B's issue" {
		t.Errorf("story/S-0002's commits on main:\n%s", got)
	}
	summary := gitIn(t, b, "show", "HEAD:design/issues/summary.md")
	for _, want := range []string{"[I-0001]", "A's issue", "[I-0002]", "B's issue", "[I-0003]", "B's second issue", "updated: 2026-09-16T01:00:00Z"} {
		if !strings.Contains(summary, want) {
			t.Errorf("the summary lacks %q:\n%s", want, summary)
		}
	}
	if strings.Contains(summary, strings.Repeat("<", 7)) || strings.Contains(summary, strings.Repeat(">", 7)) {
		t.Errorf("the summary carries conflict markers:\n%s", summary)
	}
}

// The acceptance syncs through the same rebase, so the second story whose
// branch records an issue is accepted without stopping on the summary.
func TestAcceptRegeneratesTheIssueSummaryWhenTheRebaseStopsOnItAlone(t *testing.T) {
	root, b := issueStories(t, nil)
	issueStoryInReview(t, root, "S-0002", "T-0001")

	if _, errOut, code := runInAt(t, root, issueClock.Add(4*time.Hour), "accept", "S-0002"); code != 0 {
		t.Fatalf("acceptance stopped on the generated summary: %d %s", code, errOut)
	}
	if _, err := os.Stat(b); err == nil {
		t.Error("an accepted story's worktree is removed")
	}
	summary := gitIn(t, root, "show", "main:design/issues/summary.md")
	for _, want := range []string{"[I-0001]", "[I-0002]"} {
		if !strings.Contains(summary, want) {
			t.Errorf("main's summary lacks %s:\n%s", want, summary)
		}
	}
	if strings.Contains(summary, strings.Repeat("<", 7)) {
		t.Errorf("main's summary carries conflict markers:\n%s", summary)
	}
}

// I-0089: a story's branch records an issue while issues are recorded and
// closed on main itself, not through another story's acceptance, so both
// sides rewrite design/issues/summary.md and the story's sync stopped on it.
// The sync and the acceptance regenerate it from the issue files (ADR-0098).
func TestSyncAndAcceptRegenerateTheIssueSummaryChangedOnMainItself(t *testing.T) {
	root := syncProject(t)
	onMain := func(at time.Time, msg string, args ...string) {
		t.Helper()
		if _, errOut, code := runInAt(t, root, at, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		gitIn(t, root, "add", "design")
		gitIn(t, root, "commit", "-q", "-m", msg)
	}
	onMain(issueClock, "docs: record main's old issue", "issue", "new", "Main's old issue", "--class", "efficiency")
	b := openSyncStory(t, root, 1, "design,docs")
	recordIssueIn(t, b, issueClock.Add(time.Hour), "The story's issue", nil)
	onMain(issueClock.Add(2*time.Hour), "docs: record main's new issue", "issue", "new", "Main's new issue", "--class", "efficiency")
	onMain(issueClock.Add(3*time.Hour), "docs: close main's old issue", "issue", "close", "I-0001", "--reason", "fixed on main")

	out, errOut, code := runInAt(t, b, issueClock.Add(4*time.Hour), "stream", "sync", "S-0001")
	if code != 0 {
		t.Fatalf("sync stopped on the generated summary: %d %s %s", code, out, errOut)
	}
	if storygit.RebaseInProgress(execx.System{}, b) {
		t.Fatal("the rebase is left in progress")
	}
	if st := gitIn(t, b, "status", "--porcelain"); st != "" {
		t.Errorf("the worktree is not clean after the sync: %s", st)
	}
	open := []string{"[I-0002]", "The story's issue", "[I-0003]", "Main's new issue"}
	summary := gitIn(t, b, "show", "HEAD:design/issues/summary.md")
	for _, want := range open {
		if !strings.Contains(summary, want) {
			t.Errorf("the story's summary lacks %q:\n%s", want, summary)
		}
	}
	if strings.Contains(summary, "Main's old issue") {
		t.Errorf("the story's summary lists the issue closed on main:\n%s", summary)
	}
	if strings.Contains(summary, strings.Repeat("<", 7)) || strings.Contains(summary, strings.Repeat(">", 7)) {
		t.Errorf("the story's summary carries conflict markers:\n%s", summary)
	}

	issueStoryInReview(t, root, "S-0001", "T-0001")
	if _, errOut, code := runInAt(t, root, issueClock.Add(5*time.Hour), "accept", "S-0001"); code != 0 {
		t.Fatalf("acceptance stopped on the generated summary: %d %s", code, errOut)
	}
	summary = gitIn(t, root, "show", "main:design/issues/summary.md")
	for _, want := range open {
		if !strings.Contains(summary, want) {
			t.Errorf("main's summary lacks %q:\n%s", want, summary)
		}
	}
	if strings.Contains(summary, "Main's old issue") || strings.Contains(summary, strings.Repeat("<", 7)) {
		t.Errorf("main's summary lists the closed issue or carries conflict markers:\n%s", summary)
	}
}

// When another path conflicts as well, the summary's rows may depend on it,
// so the sync stops as before and names every path, the summary included.
func TestSyncStopsWhenTheIssueSummaryIsNotTheOnlyConflict(t *testing.T) {
	_, b := issueStories(t, map[string][2]string{"docs/guide.md": {"A's line\n", "B's line\n"}})

	out, errOut, code := runInAt(t, b, issueClock.Add(4*time.Hour), "stream", "sync", "S-0002")
	if code == 0 {
		t.Fatalf("a sync with another conflict succeeded: %s", out)
	}
	for _, want := range []string{"stopped on conflicts in 2 paths:\n", "\n  design/issues/summary.md\n", "\n  docs/guide.md\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("the stop lacks %q:\n%s", want, out)
		}
	}
	if !strings.Contains(errOut, "stopped on conflicts in design/issues/summary.md, docs/guide.md") {
		t.Errorf("error: %s", errOut)
	}
	if !storygit.RebaseInProgress(execx.System{}, b) {
		t.Fatal("the rebase is not left in progress for the agent")
	}
	if got, _ := os.ReadFile(filepath.Join(b, "design", "issues", "summary.md")); !strings.Contains(string(got), strings.Repeat("<", 7)) {
		t.Errorf("the summary was rewritten though another path conflicts:\n%s", got)
	}
}

// trialSync syncs S-0002 from its worktree wt and returns the sync's text
// output, its --json branches, the conflict threads open on the project at
// root, and its conversations.
func trialSync(t *testing.T, root, wt string) (string, []storygit.BranchCheck, []*threads.Thread, []*messages.Conversation) {
	t.Helper()
	out, errOut, code := runInAt(t, wt, issueClock.Add(4*time.Hour), "stream", "sync", "S-0002")
	if code != 0 {
		t.Fatalf("sync: %d %s %s", code, out, errOut)
	}
	js, errOut, code := runInAt(t, wt, issueClock.Add(5*time.Hour), "--json", "stream", "sync", "S-0002")
	if code != 0 {
		t.Fatalf("sync --json: %d %s %s", code, js, errOut)
	}
	var res struct {
		Branches []storygit.BranchCheck `json:"branches"`
	}
	if err := json.Unmarshal([]byte(js), &res); err != nil {
		t.Fatalf("%v: %s", err, js)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	all, err := threads.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	var open []*threads.Thread
	for _, th := range all {
		if th.Open() && storygit.IsConflictTitle(th.Title) {
			open = append(open, th)
		}
	}
	convs, err := messages.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	return out, res.Branches, open, convs
}

// I-0074: two open stories that each record an issue conflict in
// design/issues/summary.md alone, which the trial merge reported, telling
// the pair. The file is generated, so the pair is clean (ADR-0098).
func TestSyncTrialMergeLeavesOutTheIssueSummary(t *testing.T) {
	root, b := openIssueStories(t, nil)
	// git itself reports the pair as conflicting in the summary
	if raw, err := storygit.TrialMerge(execx.System{}, root, "story/S-0002", "story/S-0001"); err != nil || strings.Join(raw, ",") != "design/issues/summary.md" {
		t.Fatalf("git's trial merge: %v %q", err, raw)
	}

	out, branches, _, convs := trialSync(t, root, b)
	if !strings.Contains(out, "story/S-0002 merges cleanly with story/S-0001 (in progress)\n") || strings.Contains(out, "conflicts with") {
		t.Errorf("sync output:\n%s", out)
	}
	if len(branches) != 1 || branches[0].Story != "S-0001" || !branches[0].Clean || branches[0].Conflicts == nil || len(branches[0].Conflicts) != 0 || branches[0].Conversation != "" {
		t.Errorf("branches: %+v", branches)
	}
	if len(convs) != 0 {
		t.Errorf("the pair was told of a conflict: %+v", convs)
	}
}

// When the pair conflicts in another file as well, the trial merge reports
// that file alone, in the sync's output and in the pair's conversation.
func TestSyncTrialMergeReportsOtherConflictsWithoutTheIssueSummary(t *testing.T) {
	root, b := openIssueStories(t, map[string][2]string{"docs/guide.md": {"A's line\n", "B's line\n"}})

	out, branches, _, convs := trialSync(t, root, b)
	if !strings.Contains(out, "story/S-0002 conflicts with story/S-0001 (in progress) in docs/guide.md; see MS-0001\n") {
		t.Errorf("sync output:\n%s", out)
	}
	if len(branches) != 1 || branches[0].Clean || strings.Join(branches[0].Conflicts, ",") != "docs/guide.md" || branches[0].Conversation != "MS-0001" {
		t.Errorf("branches: %+v", branches)
	}
	if len(convs) != 1 || strings.Join(convs[0].About, ",") != "docs/guide.md" {
		t.Fatalf("conversations: %+v", convs)
	}
	for _, e := range convs[0].Entries() {
		if !strings.Contains(e.Text, "- `docs/guide.md`\n") || strings.Contains(e.Text, "summary.md") {
			t.Errorf("conflict message: %s", e.Text)
		}
	}
}

// staleStories makes S-0001, whose branch changes docs/guide.md and then
// falls behind main, which changes it too, and S-0002, opened after main's
// change, whose branch changes design/note.md (I-0064). Each story then
// commits its side of files. It returns the main checkout and S-0002's
// worktree.
func staleStories(t *testing.T, files map[string][2]string) (root, a string) {
	t.Helper()
	root = syncProject(t)
	b := openSyncStory(t, root, 1, "docs")
	commitIn(t, b, "docs/guide.md", "stale branch's line\n")
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("main's line\n"), 0o644)
	gitIn(t, root, "add", "docs/guide.md")
	gitIn(t, root, "commit", "-q", "-m", "docs: main's guide")
	a = openSyncStory(t, root, 2, "design,docs")
	commitIn(t, a, "design/note.md", "fresh branch's note\n")
	for p, c := range files {
		commitIn(t, b, p, c[0])
		commitIn(t, a, p, c[1])
	}
	return root, a
}

// I-0064: a stale branch's change to a path main has changed since conflicts
// with main, not with a fresh branch that carries main's change and leaves
// the path alone. The trial merge stopped on it and blamed the fresh story's
// sync; a pair's conflicts are the paths both changed, so the pair is clean.
func TestSyncTrialMergeLeavesOutWhatMainBrought(t *testing.T) {
	root, a := staleStories(t, nil)
	// git itself reports the pair as conflicting in main's change
	if raw, err := storygit.TrialMerge(execx.System{}, root, "story/S-0002", "story/S-0001"); err != nil || strings.Join(raw, ",") != "docs/guide.md" {
		t.Fatalf("git's trial merge: %v %q", err, raw)
	}

	out, branches, _, convs := trialSync(t, root, a)
	if !strings.Contains(out, "story/S-0002 merges cleanly with story/S-0001 (in progress)\n") || strings.Contains(out, "conflicts with") {
		t.Errorf("sync output:\n%s", out)
	}
	if len(branches) != 1 || branches[0].Story != "S-0001" || !branches[0].Clean || branches[0].Conflicts == nil || len(branches[0].Conflicts) != 0 || branches[0].Conversation != "" {
		t.Errorf("branches: %+v", branches)
	}
	if len(convs) != 0 {
		t.Errorf("the pair was told of a conflict: %+v", convs)
	}
}

// A path both branches change still conflicts beside one main brought, and
// is the only one the sync and the pair's conversation name.
func TestSyncTrialMergeReportsWhatBothChangedBesideWhatMainBrought(t *testing.T) {
	root, a := staleStories(t, map[string][2]string{"docs/more.md": {"stale branch's more\n", "fresh branch's more\n"}})
	if raw, err := storygit.TrialMerge(execx.System{}, root, "story/S-0002", "story/S-0001"); err != nil || strings.Join(raw, ",") != "docs/guide.md,docs/more.md" {
		t.Fatalf("git's trial merge: %v %q", err, raw)
	}

	out, branches, _, convs := trialSync(t, root, a)
	if !strings.Contains(out, "story/S-0002 conflicts with story/S-0001 (in progress) in docs/more.md; see MS-0001\n") {
		t.Errorf("sync output:\n%s", out)
	}
	if len(branches) != 1 || branches[0].Clean || strings.Join(branches[0].Conflicts, ",") != "docs/more.md" || branches[0].Conversation != "MS-0001" {
		t.Errorf("branches: %+v", branches)
	}
	if len(convs) != 1 {
		t.Fatalf("conversations: %+v", convs)
	}
	for _, e := range convs[0].Entries() {
		if !strings.Contains(e.Text, "- `docs/more.md`\n") || strings.Contains(e.Text, "guide.md") {
			t.Errorf("conflict message: %s", e.Text)
		}
	}
}

// A pair's thread that named only what main brought, as a sync before the
// fix wrote it, is resolved at the next sync, which finds the pair clean and
// names the thread; it opens no conversation.
func TestSyncResolvesAThreadOnWhatMainBrought(t *testing.T) {
	root, a := staleStories(t, nil)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	th, err := threads.New(repo, threads.NewOptions{Title: storygit.ConflictTitle("S-0002", "S-0001"), On: "S-0002", Author: storygit.ConflictAuthor,
		Text: storygit.ConflictText("S-0002", "S-0001", []string{"docs/guide.md"}), Now: issueClock.Add(3 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}

	out, _, open, convs := trialSync(t, root, a)
	if want := "story/S-0002 merges cleanly with story/S-0001 (in progress); resolved " + th.ID + "\n"; !strings.Contains(out, want) {
		t.Errorf("sync output lacks %q:\n%s", want, out)
	}
	if len(open) != 0 || len(convs) != 0 {
		t.Errorf("conflict threads left open: %+v, conversations: %+v", open, convs)
	}
	th, err = threads.Get(repo, th.ID)
	if err != nil {
		t.Fatal(err)
	}
	if e := th.Entries(); th.Open() || !strings.Contains(e[len(e)-1].Text, "story/S-0001 and story/S-0002 merge cleanly at the sync of S-0002") {
		t.Errorf("not resolved: %+v", e)
	}
}

// The generated files are one list, the issue summary alone today, named as
// git names them (ADR-0098).
func TestGeneratedPathsAreTheIssueSummary(t *testing.T) {
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range issues.Generated(repo, "S-0001", time.Now) {
		got = append(got, f.Path)
	}
	if strings.Join(got, ",") != "design/issues/summary.md" {
		t.Errorf("generated paths: %q", got)
	}
}

// issueClock is when main's summary is generated; the stories record their
// issues an hour and two after it.
var issueClock = time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)

// issueStories makes the two stories of openIssueStories and merges S-0001
// into main as its acceptance would (I-0074). It returns the main checkout
// and S-0002's worktree.
func issueStories(t *testing.T, files map[string][2]string) (root, b string) {
	t.Helper()
	root, b = openIssueStories(t, files)
	gitIn(t, root, "merge", "-q", "--no-edit", "story/S-0001")
	return root, b
}

// openIssueStories makes S-0001 and S-0002 in progress, each of whose
// branches records an issue at its own time, and so rewrites
// design/issues/summary.md's updated line and rows where the other does.
// Each story commits its side of files with its issue. It returns the main
// checkout and S-0002's worktree.
func openIssueStories(t *testing.T, files map[string][2]string) (root, b string) {
	t.Helper()
	root = syncProject(t)
	if _, errOut, code := runInAt(t, root, issueClock, "issue", "summary"); code != 0 {
		t.Fatal(errOut)
	}
	gitIn(t, root, "add", "design")
	gitIn(t, root, "commit", "-q", "-m", "docs: issue summary")
	a := openSyncStory(t, root, 1, "design,docs")
	b = openSyncStory(t, root, 2, "design,docs")
	side := func(i int) map[string]string {
		m := map[string]string{}
		for p, c := range files {
			m[p] = c[i]
		}
		return m
	}
	recordIssueIn(t, a, issueClock.Add(time.Hour), "A's issue", side(0))
	recordIssueIn(t, b, issueClock.Add(2*time.Hour), "B's issue", side(1))
	return root, b
}

// recordIssueIn records an issue with flai issue new in the worktree wt at
// at, writes files there, and commits them all.
func recordIssueIn(t *testing.T, wt string, at time.Time, title string, files map[string]string) {
	t.Helper()
	if _, errOut, code := runInAt(t, wt, at, "issue", "new", title, "--class", "efficiency"); code != 0 {
		t.Fatal(errOut)
	}
	for p, c := range files {
		if err := os.WriteFile(filepath.Join(wt, p), []byte(c), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "chore: record "+title)
}

// issueStoryInReview ticks the story's criterion, does a task for it, and
// moves it to review, for flai accept.
func issueStoryInReview(t *testing.T, root, id, task string) {
	t.Helper()
	matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
	if len(matches) != 1 {
		t.Fatalf("story file for %s: %v", id, matches)
	}
	s, _ := os.ReadFile(matches[0])
	_ = os.WriteFile(matches[0], []byte(strings.Replace(string(s), "- [ ] done\n", "- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{
		{"task", "new", "Do it", "--story", id},
		{"move", task, "ready"}, {"move", task, "in-progress"}, {"move", task, "done"},
		{"move", id, "review"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
}
