package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fastPermissionPoll makes a held permission_prompt read its thread often.
func fastPermissionPoll(t *testing.T) {
	t.Helper()
	was := permissionPoll
	permissionPoll = 10 * time.Millisecond
	t.Cleanup(func() { permissionPoll = was })
}

// askPermission calls permission_prompt as Claude Code does and decodes the
// decision it answers.
func (f *fixture) askPermission(t *testing.T, ctx context.Context, tool string, input map[string]any) PermissionOut {
	t.Helper()
	return decision(t, callPermission(t, ctx, f.cs, map[string]any{"tool_name": tool, "input": input, "tool_use_id": "toolu_1"}))
}

func callPermission(t *testing.T, ctx context.Context, cs *mcp.ClientSession, args any) *mcp.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "permission_prompt", Arguments: args})
	if err != nil {
		t.Fatalf("permission_prompt: %v", err)
	}
	return res
}

// decision decodes the result as Claude Code reads it, failing unless it is
// what Claude Code accepts (I-0082): a single text block holding a JSON
// object, with no structured content beside it and no tool error.
func decision(t *testing.T, res *mcp.CallToolResult) PermissionOut {
	t.Helper()
	if res.IsError {
		t.Fatalf("permission_prompt should answer a decision, not fail: %+v", res)
	}
	if res.StructuredContent != nil {
		t.Errorf("the result carries structured content, which Claude Code calls an invalid result: %+v", res.StructuredContent)
	}
	if len(res.Content) != 1 {
		t.Fatalf("the result has %d content blocks, want a single text block: %+v", len(res.Content), res.Content)
	}
	text, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("the content is not text: %T", res.Content[0])
	}
	var keys map[string]any
	if err := json.Unmarshal([]byte(text.Text), &keys); err != nil {
		t.Fatalf("the text is not a JSON object: %q", text.Text)
	}
	for k := range keys {
		if k != "behavior" && k != "updatedInput" && k != "message" {
			t.Errorf("the decision carries %q, which Claude Code does not read: %s", k, text.Text)
		}
	}
	if _, ok := keys["behavior"]; !ok {
		t.Errorf("the decision has no behavior: %s", text.Text)
	}
	var out PermissionOut
	_ = json.Unmarshal([]byte(text.Text), &out)
	return out
}

// I-0082: Claude Code refuses any answer but a single text block, so the
// tool declares no output schema and answers no structured content, whether
// it allows, refuses at once, or refuses a request it cannot read.
func TestPermissionPromptAnswersOneTextBlockAsClaudeCodeRequires(t *testing.T) {
	f := setupWith(t, func(o *Options) { o.AutoApprove = func(string) bool { return true } })
	ctx := context.Background()
	tools, err := f.cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	var listed bool
	for _, tool := range tools.Tools {
		if tool.Name == "permission_prompt" {
			listed = true
			if tool.OutputSchema != nil {
				t.Errorf("permission_prompt declares an output schema, so it answers structured content: %+v", tool.OutputSchema)
			}
		}
	}
	if !listed {
		t.Fatal("permission_prompt is not listed")
	}
	input := map[string]any{"file_path": f.worktreeFile(".claude/settings.json"), "content": "{}"}
	for _, c := range []struct {
		name     string
		args     any
		behavior string
		message  string
	}{
		{"an allow", map[string]any{"tool_name": "Write", "input": input}, "allow", ""},
		{"a refusal at once", map[string]any{"tool_name": "Bash", "input": map[string]any{"command": "ls"}}, "deny", "Bash is not an Edit"},
		{"arguments it cannot read", map[string]any{"tool_name": 5}, "deny", "not a permission request"},
	} {
		out := decision(t, callPermission(t, ctx, f.cs, c.args))
		if out.Behavior != c.behavior || !strings.Contains(out.Message, c.message) {
			t.Errorf("%s: %+v, want %s saying %q", c.name, out, c.behavior, c.message)
		}
	}
	noThreads(t, f.repo)
}

// projectLint is the project's markdownlint configuration, so that the
// threads permission_prompt opens are held to it.
const projectLint = "default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD033: false\nMD041: false\nMD060: false\n"

func (f *fixture) lintAsTheProjectDoes(t *testing.T) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(f.repo.Root, ".markdownlint.yaml"), []byte(projectLint), 0o644); err != nil {
		t.Fatal(err)
	}
	bad := threads.NewOptions{Title: "Tabs", On: f.story.ID, Author: "alex", Text: "```text\n\tx\n```", Now: t0}
	if _, err := threads.New(f.repo, bad); err == nil {
		t.Fatal("the lint does not refuse a tab in a fenced block, so it proves nothing here")
	}
}

func (f *fixture) worktreeFile(rel string) string {
	return filepath.Join(f.repo.WorktreePath(f.story.ID), filepath.FromSlash(rel))
}

func noThreads(t *testing.T, repo *workitem.Repo) {
	t.Helper()
	if all, _ := threads.List(repo); len(all) != 0 {
		t.Errorf("a refusal at once opens no thread: %d opened", len(all))
	}
}

func TestPermissionPromptRefusesAtOnceWhatItDoesNotApprove(t *testing.T) {
	f := setup(t)
	other, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Other", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(other.Path)
	_ = os.WriteFile(other.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	other, _ = f.repo.Get(other.ID)
	if _, err := f.repo.Transition(other, workitem.Ready, "alex", "", t0); err != nil {
		t.Fatal(err)
	}
	settings := f.worktreeFile(".claude/settings.json")
	// A .claude folder in the worktree that is a link to the main checkout's.
	outside := filepath.Join(f.repo.Root, ".claude")
	linked := filepath.Join(f.repo.WorktreePath(f.story.ID), "linked")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(linked, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(linked, ".claude")); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name, tool string
		input      map[string]any
		why        string
	}{
		{"another tool", "Bash", map[string]any{"command": "rm -rf /"}, "Bash is not an Edit, Write, MultiEdit, or NotebookEdit"},
		{"no path", "Write", map[string]any{"content": "x"}, "names no absolute path"},
		{"a relative path", "Edit", map[string]any{"file_path": ".claude/settings.json"}, "names no absolute path"},
		{"outside any worktree", "Write", map[string]any{"file_path": filepath.Join(f.repo.Root, ".claude/settings.json")}, "is not inside a story's worktree"},
		{"an escape with ..", "Write", map[string]any{"file_path": f.repo.WorktreePath(f.story.ID) + "/../../../.claude/settings.json"}, "is not inside a story's worktree"},
		{"an escape through a symbolic link", "Write", map[string]any{"file_path": filepath.Join(linked, ".claude/settings.json")}, "through a symbolic link"},
		{"not a protected path", "Write", map[string]any{"file_path": f.worktreeFile("flai/claude/settings.json")}, "is not a path Claude Code protects in " + f.story.ID + "'s worktree"},
		{"a file named .claude", "Write", map[string]any{"file_path": f.worktreeFile(".claude")}, "is not a path Claude Code protects"},
		{"a name like a protected file", "Write", map[string]any{"file_path": f.worktreeFile("my.mcp.json")}, "is not a path Claude Code protects"},
		{"the worktree's .git", "Write", map[string]any{"file_path": f.worktreeFile(".git")}, "is in .git, which permission_prompt never approves"},
		{"a hook in a nested .git", "Edit", map[string]any{"file_path": f.worktreeFile("sub/.git/hooks/pre-commit")}, "is in .git, which permission_prompt never approves"},
		{"a protected path in the main checkout", "Write", map[string]any{"file_path": filepath.Join(f.repo.Root, ".mcp.json")}, "is not inside a story's worktree"},
		{"a protected path outside the project", "Write", map[string]any{"file_path": filepath.Join(t.TempDir(), ".vscode/settings.json")}, "is not inside a story's worktree"},
		{"no such story", "Write", map[string]any{"file_path": filepath.Join(f.repo.WorktreePath("S-0099"), ".claude/settings.json")}, "is not a story's worktree"},
		{"a story not in progress", "Write", map[string]any{"file_path": filepath.Join(f.repo.WorktreePath(other.ID), ".claude/settings.json")}, other.ID + " is ready, not in progress"},
		{"a notebook without notebook_path", "NotebookEdit", map[string]any{"file_path": settings}, "names no absolute path"},
	} {
		out := f.askPermission(t, context.Background(), c.tool, c.input)
		if out.Behavior != "deny" || !strings.Contains(out.Message, c.why) || !strings.Contains(out.Message, permissionScope) || out.UpdatedInput != nil {
			t.Errorf("%s: %+v, want a refusal saying %q and what the tool approves", c.name, out, c.why)
		}
	}
	noThreads(t, f.repo)
}

func TestPermissionPromptAutoApprovesWithoutAThread(t *testing.T) {
	var asked string
	var logged bytes.Buffer
	f := setupWith(t, func(o *Options) {
		o.AutoApprove = func(root string) bool { asked = root; return true }
		o.Logger = slog.New(slog.NewTextHandler(&logged, nil))
	})
	input := map[string]any{"file_path": f.worktreeFile(".claude/settings.json"), "old_string": "a", "new_string": "b"}
	out := f.askPermission(t, context.Background(), "Edit", input)
	if out.Behavior != "allow" || !reflect.DeepEqual(out.UpdatedInput, input) || out.Message != "" {
		t.Errorf("auto-approve allows with the input unchanged: %+v", out)
	}
	if asked != f.repo.Root {
		t.Errorf("AutoApprove is asked with the project's root %s, got %q", f.repo.Root, asked)
	}
	line := logged.String()
	for _, want := range []string{"permission allowed without asking", "agent=claude", "tool=Edit", "settings.json"} {
		if !strings.Contains(line, want) {
			t.Errorf("the log does not say %q: %s", want, line)
		}
	}
	noThreads(t, f.repo)
}

// ADR-0106: every path Claude Code protects in the story's worktree is
// approved, not only .claude/, while .git, the main checkout, and every path
// outside the worktree stay refused, auto-approve or not.
func TestPermissionPromptAutoApprovesEveryProtectedPathButGit(t *testing.T) {
	f := setupWith(t, func(o *Options) { o.AutoApprove = func(string) bool { return true } })
	for _, rel := range []string{".mcp.json", ".vscode/settings.json", "template/root/.mcp.json", "template/root/.claude/agents/x.md", ".husky/pre-commit"} {
		input := map[string]any{"file_path": f.worktreeFile(rel), "content": "{}"}
		if out := f.askPermission(t, context.Background(), "Write", input); out.Behavior != "allow" || !reflect.DeepEqual(out.UpdatedInput, input) {
			t.Errorf("%s: auto-approve allows a protected path with the input unchanged: %+v", rel, out)
		}
	}
	for _, c := range []struct{ path, why string }{
		{f.worktreeFile(".git"), "is in .git, which permission_prompt never approves"},
		{f.worktreeFile(".git/config"), "is in .git, which permission_prompt never approves"},
		{f.worktreeFile("sub/.git/hooks/pre-commit"), "is in .git, which permission_prompt never approves"},
		{f.worktreeFile("flai/main.go"), "is not a path Claude Code protects"},
		{filepath.Join(f.repo.Root, ".mcp.json"), "is not inside a story's worktree"},
		{filepath.Join(t.TempDir(), ".mcp.json"), "is not inside a story's worktree"},
	} {
		out := f.askPermission(t, context.Background(), "Write", map[string]any{"file_path": c.path, "content": "x"})
		if out.Behavior != "deny" || !strings.Contains(out.Message, c.why) || out.UpdatedInput != nil {
			t.Errorf("%s: %+v, want a refusal saying %q under auto-approve", c.path, out, c.why)
		}
	}
	noThreads(t, f.repo)
}

// Without auto-approve, a protected path outside .claude/ is asked on a
// thread like one in it.
func TestPermissionPromptAsksForEveryProtectedPath(t *testing.T) {
	fastPermissionPoll(t)
	for _, rel := range []string{".mcp.json", ".vscode/settings.json", "template/root/.mcp.json"} {
		t.Run(rel, func(t *testing.T) {
			f := setup(t)
			input := map[string]any{"file_path": f.worktreeFile(rel), "content": "{}"}
			got := make(chan PermissionOut, 1)
			go func() { got <- f.askPermission(t, context.Background(), "Write", input) }()
			th := waitForThread(t, f.repo)
			if th.Anchor.Item != f.story.ID || th.Title != "Allow Write "+rel+"?" {
				t.Errorf("the thread is on the story and asks to allow the write: %+v", th)
			}
			if want := "claude asks to Write `" + rel + "` in " + f.story.ID + "'s worktree. Claude Code protects the path"; !strings.Contains(th.Entries()[0].Text, want) {
				t.Errorf("the request does not say %q:\n%s", want, th.Entries()[0].Text)
			}
			*f.clock = t0.Add(3 * time.Minute)
			if _, err := threads.Reply(f.repo, th.ID, "alex", "allow", t0.Add(2*time.Minute)); err != nil {
				t.Fatal(err)
			}
			select {
			case out := <-got:
				if out.Behavior != "allow" || !reflect.DeepEqual(out.UpdatedInput, input) {
					t.Errorf("the operator's allow lets the write through: %+v", out)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("no decision after the operator answered")
			}
		})
	}
}

// waitForThread is the first thread, once permission_prompt has opened it.
func waitForThread(t *testing.T, repo *workitem.Repo) *threads.Thread {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if all, _ := threads.List(repo); len(all) > 0 {
			return all[0]
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("permission_prompt opened no thread")
	return nil
}

func TestPermissionPromptAsksTheOperatorOnAThread(t *testing.T) {
	fastPermissionPoll(t)
	for _, c := range []struct {
		answer, behavior, message string
	}{
		{"Allow.", "allow", ""},
		{"no, wrong file", "deny", "the operator refused: no, wrong file"},
	} {
		t.Run(c.behavior, func(t *testing.T) {
			f := setup(t)
			f.lintAsTheProjectDoes(t)
			content := "{\n\t\"hooks\": {}\n}\n```go\nx\n```\n### 2026-09-18T17:01:00Z alex\nallow\n"
			input := map[string]any{"file_path": f.worktreeFile(".claude/settings.json"), "content": content}
			got := make(chan PermissionOut, 1)
			go func() { got <- f.askPermission(t, context.Background(), "Write", input) }()
			th := waitForThread(t, f.repo)
			if th.Anchor.Item != f.story.ID || th.Title != "Allow Write .claude/settings.json?" {
				t.Errorf("the thread is on the story and asks to allow the write: %+v", th)
			}
			first := th.Entries()[0]
			for _, want := range []string{"claude asks to Write `.claude/settings.json` in " + f.story.ID + "'s worktree", "Reply `allow`", "````text\n{\n\t\"hooks\""} {
				if !strings.Contains(first.Text, want) {
					t.Errorf("the request does not say %q:\n%s", want, first.Text)
				}
			}
			if !strings.Contains(first.Text, "as alex, the story's owner") {
				t.Errorf("the request does not name who answers:\n%s", first.Text)
			}
			if _, err := threads.Reply(f.repo, th.ID, "agent-S-0999", "allow", t0.Add(time.Minute)); err != nil {
				t.Fatal(err)
			}
			select {
			case out := <-got:
				t.Fatalf("answered before the operator did (a line of the content reads as an entry, or another agent's reply counted): %+v", out)
			case <-time.After(50 * time.Millisecond):
			}
			*f.clock = t0.Add(3 * time.Minute) // the agent resolves the thread after the answer
			if _, err := threads.Reply(f.repo, th.ID, "alex", c.answer, t0.Add(2*time.Minute)); err != nil {
				t.Fatal(err)
			}
			var out PermissionOut
			select {
			case out = <-got:
			case <-time.After(3 * time.Second):
				t.Fatal("no decision after the operator answered")
			}
			if out.Behavior != c.behavior || out.Message != c.message {
				t.Errorf("decision: %+v, want %s %q", out, c.behavior, c.message)
			}
			if c.behavior == "allow" && !reflect.DeepEqual(out.UpdatedInput, input) {
				t.Errorf("allow carries the input unchanged: %+v", out.UpdatedInput)
			}
			back, _ := threads.Get(f.repo, th.ID)
			last := back.Entries()[len(back.Entries())-1]
			if back.Open() || last.Author != "claude" || !strings.Contains(last.Text, map[string]string{"allow": "allowed by alex", "deny": "refused by alex"}[c.behavior]) {
				t.Errorf("the thread is resolved by the agent with what was decided: %s %+v", back.Status, last)
			}
		})
	}
}

// ownedBy gives the story the owner story and the project the owner project,
// as on S-0218 and S-0222 (I-0081).
func (f *fixture) ownedBy(t *testing.T, story, project string) {
	t.Helper()
	data, err := os.ReadFile(f.story.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.story.Path, []byte(strings.Replace(string(data), "\nowner: alex\n", "\nowner: "+story+"\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if f.story, err = f.repo.Get(f.story.ID); err != nil || f.story.Owner != story {
		t.Fatalf("the story's owner is not %s: %v %+v", story, err, f.story)
	}
	f.repo.Manifest.Owner = project
}

// replyHeld replies on the thread and fails if the prompt decides on it: the
// reply is not an answer.
func replyHeld(t *testing.T, f *fixture, id, author string, at time.Duration, got <-chan PermissionOut) {
	t.Helper()
	if _, err := threads.Reply(f.repo, id, author, "allow", t0.Add(at)); err != nil {
		t.Fatal(err)
	}
	select {
	case out := <-got:
		t.Fatalf("a reply by %s was taken as the answer: %+v", author, out)
	case <-time.After(50 * time.Millisecond):
	}
}

// I-0081: the story's owner was arobson and the operator answered as alex,
// the project's owner, so the write waited until Claude Code gave up. Either
// owner's answer now lets it through; another agent's reply still does not.
func TestPermissionPromptTakesTheProjectOwnersAnswer(t *testing.T) {
	fastPermissionPoll(t)
	for _, by := range []string{"alex", "arobson"} {
		t.Run(by, func(t *testing.T) {
			f := setup(t)
			f.ownedBy(t, "arobson", "alex")
			input := map[string]any{"file_path": f.worktreeFile(".claude/settings.json"), "content": "{}"}
			got := make(chan PermissionOut, 1)
			go func() { got <- f.askPermission(t, context.Background(), "Write", input) }()
			th := waitForThread(t, f.repo)
			if want := "Reply `allow`, as arobson, the story's owner, or alex, the project's owner, to let it write."; !strings.Contains(th.Entries()[0].Text, want) {
				t.Errorf("the request does not name both who may answer, %q:\n%s", want, th.Entries()[0].Text)
			}
			replyHeld(t, f, th.ID, "agent-S-0999", time.Minute, got)
			*f.clock = t0.Add(3 * time.Minute)
			if _, err := threads.Reply(f.repo, th.ID, by, "allow", t0.Add(2*time.Minute)); err != nil {
				t.Fatal(err)
			}
			select {
			case out := <-got:
				if out.Behavior != "allow" || !reflect.DeepEqual(out.UpdatedInput, input) {
					t.Errorf("%s's allow lets the write through: %+v", by, out)
				}
			case <-time.After(3 * time.Second):
				t.Fatalf("no decision after %s answered allow (I-0081)", by)
			}
			back, _ := threads.Get(f.repo, th.ID)
			if last := back.Entries()[len(back.Entries())-1]; back.Open() || !strings.Contains(last.Text, "allowed by "+by) {
				t.Errorf("the thread is resolved as allowed by %s: %s %+v", by, back.Status, last)
			}
		})
	}
}

// The asking agent's own entry is never an answer, even under an owner's name.
func TestPermissionPromptNeverTakesTheAgentsOwnEntry(t *testing.T) {
	fastPermissionPoll(t)
	f := setup(t)
	f.ownedBy(t, "arobson", "alex")
	s := newServer(Options{Repo: f.repo, Agent: "alex", Now: func() time.Time { return t0.Add(time.Hour) }}, f.repo)
	got := make(chan PermissionOut, 1)
	go func() {
		got <- s.permissionPrompt(context.Background(), PermissionIn{ToolName: "Write", Input: map[string]any{"file_path": f.worktreeFile(".claude/settings.json"), "content": "{}"}})
	}()
	th := waitForThread(t, f.repo)
	replyHeld(t, f, th.ID, "alex", time.Minute, got)
	replyHeld(t, f, th.ID, "agent-S-0999", 2*time.Minute, got)
	if _, err := threads.Reply(f.repo, th.ID, "arobson", "no, not now", t0.Add(3*time.Minute)); err != nil {
		t.Fatal(err)
	}
	select {
	case out := <-got:
		if out.Behavior != "deny" || out.Message != "the operator refused: no, not now" {
			t.Errorf("the story's owner refused: %+v", out)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no decision after the story's owner answered")
	}
}

func TestPermissionRequestNamesWhoMayAnswer(t *testing.T) {
	for _, c := range []struct {
		owner, project, says string
		who                  []string
	}{
		{"arobson", "alex", "Reply `allow`, as arobson, the story's owner, or alex, the project's owner, to let it write.", []string{"arobson", "alex"}},
		{"alex", "alex", "Reply `allow`, as alex, the story's owner, to let it write.", []string{"alex"}},
		{"arobson", "", "Reply `allow`, as arobson, the story's owner, to let it write.", []string{"arobson"}},
		{"", "alex", "Reply `allow`, as alex, the project's owner, to let it write.", []string{"alex"}},
		{"", "", "Reply `allow` to let it write.", nil},
	} {
		text := permissionRequest("claude", "S-0001", c.owner, c.project, PermissionIn{ToolName: "Write"}, ".claude/x")
		if !strings.Contains(text, c.says) {
			t.Errorf("owners %q and %q: the request does not say %q:\n%s", c.owner, c.project, c.says, text)
		}
		if who := answerers(c.owner, c.project); !reflect.DeepEqual(who, c.who) {
			t.Errorf("owners %q and %q: answerers %q, want %q", c.owner, c.project, who, c.who)
		}
	}
}

func TestPermissionPromptShowsEditsLintClean(t *testing.T) {
	fastPermissionPoll(t)
	f := setup(t)
	f.lintAsTheProjectDoes(t)
	s := newServer(Options{Repo: f.repo, Agent: "claude", Now: func() time.Time { return t0 }}, f.repo)
	for i, c := range []struct {
		tool  string
		input map[string]any
		want  []string
	}{
		{"Edit", map[string]any{"file_path": f.worktreeFile(".claude/settings.json"), "old_string": "\"matcher\": \"Edit\"", "new_string": "\"matcher\": \"Edit|Write\"", "replace_all": true}, []string{"It would replace this text:", "with this text:", "Edit|Write", "replace_all"}},
		{"MultiEdit", map[string]any{"file_path": f.worktreeFile(".claude/agents/verifier.md"), "edits": []any{map[string]any{"old_string": "a", "new_string": ""}, map[string]any{"old_string": "``x``", "new_string": "y"}}}, []string{"It would make 2 edits", "Edit 1 would replace", "Edit 2 would replace", "```text\n``x``\n```"}},
		{"NotebookEdit", map[string]any{"notebook_path": f.worktreeFile(".claude/n.ipynb"), "new_source": "print(1)", "edit_mode": "insert", "cell_id": "c1"}, []string{"Edit mode insert, at cell `c1`", "print(1)"}},
	} {
		ctx, cancel := context.WithCancel(context.Background())
		got := make(chan PermissionOut, 1)
		go func() { got <- s.permissionPrompt(ctx, PermissionIn{ToolName: c.tool, Input: c.input}) }()
		var th *threads.Thread
		deadline := time.Now().Add(3 * time.Second)
		for th == nil && time.Now().Before(deadline) {
			if all, _ := threads.List(f.repo); len(all) > i {
				th = all[i]
			}
			time.Sleep(10 * time.Millisecond)
		}
		if th == nil {
			cancel()
			t.Fatalf("%s: no thread opened: %+v", c.tool, <-got)
		}
		for _, want := range c.want {
			if !strings.Contains(th.Entries()[0].Text, want) {
				t.Errorf("%s: the request does not show %q:\n%s", c.tool, want, th.Entries()[0].Text)
			}
		}
		cancel()
		if out := <-got; out.Behavior != "deny" || out.Message != "no answer from the operator before the session ended" {
			t.Errorf("%s: a session that ends while waiting is refused: %+v", c.tool, out)
		}
		if back, _ := threads.Get(f.repo, th.ID); back.Open() {
			t.Errorf("%s: the thread is resolved when the session ends: %s", c.tool, back.Status)
		}
	}
}

func TestAllowedReadsTheFirstWord(t *testing.T) {
	for text, want := range map[string]bool{"allow": true, "Yes!": true, "OK, go": true, "approved.": true, "allowed": false, "no": false, "": false, "not ok": false} {
		if allowed(text) != want {
			t.Errorf("allowed(%q) = %v", text, !want)
		}
	}
}

// Claude Code names no project: in a folder of projects, the path says which.
func TestPermissionPromptInAFolderFindsTheProjectByThePath(t *testing.T) {
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	beta := makeProject(t, filepath.Join(root, "beta"), "beta")
	story := readyStoryIn(t, beta, "Beta work", t0)
	if _, err := beta.Transition(story, workitem.InProgress, "alex", "", t0); err != nil {
		t.Fatal(err)
	}
	f := folderSetup(t, root)
	path := filepath.Join(beta.WorktreePath(story.ID), ".claude", "settings.json")
	res := callPermission(t, context.Background(), f.cs, map[string]any{"tool_name": "Bash", "input": map[string]any{"command": "ls"}})
	if out := decision(t, res); out.Behavior != "deny" || !strings.Contains(out.Message, "name one with project") {
		t.Errorf("without a path to go by, the refusal says to name the project: %+v", out)
	}
	fastPermissionPoll(t)
	got := make(chan string, 1)
	go func() {
		res, err := f.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "permission_prompt", Arguments: map[string]any{"tool_name": "Write", "input": map[string]any{"file_path": path, "content": "{}"}}})
		if err != nil || res.IsError || res.StructuredContent != nil || len(res.Content) != 1 {
			got <- "not a single text block"
			return
		}
		got <- res.Content[0].(*mcp.TextContent).Text
	}()
	th := waitForThread(t, beta)
	if th.Anchor.Item != story.ID {
		t.Errorf("the thread is on beta's story: %+v", th.Anchor)
	}
	if _, err := threads.Reply(beta, th.ID, "alex", "yes", t0.Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if text := <-got; !strings.HasPrefix(text, `{"behavior":"allow","updatedInput":{`) {
		t.Errorf("allowed in beta: %s", text)
	}
}
