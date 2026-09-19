package mcpserver

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0063: an acceptance made where nothing could push it is reported to
// agents on every inbox call while it is true, like ready work, and by the
// board tool, and no longer once it is pushed.
func TestInboxAndBoardReportAnUnpushedAcceptance(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	base := t.TempDir()
	root, remote := filepath.Join(base, "clone"), filepath.Join(base, "origin.git")
	git := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	git(base, "init", "-q", "--bare", "-b", "main", remote)
	git(root, "init", "-q", "-b", "main")
	git(root, "add", "-A")
	git(root, "commit", "-q", "-m", "init")
	git(root, "remote", "add", "origin", remote)
	git(root, "push", "-q", "-u", "origin", "main")

	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	connect := func(r execx.Runner) *fixture {
		srv := New(Options{Repo: repo, Agent: "claude", Version: "test", Now: func() time.Time { return t0 }, Runner: r})
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
		return &fixture{repo: repo, cs: cs}
	}
	f := connect(execx.System{})

	if out, _ := f.call(t, "inbox", map[string]any{}); out["unpushed"] != nil {
		t.Fatalf("nothing is pending yet: %v", out["unpushed"])
	}
	git(root, "commit", "-q", "--allow-empty", "-m", "chore: [S-0007] accept and archive; release cli 1.1.0")
	git(root, "tag", "-a", "cli/v1.1.0", "-m", "cli 1.1.0")

	for i, tool := range []string{"inbox", "inbox", "board"} {
		out, failed := f.call(t, tool, map[string]any{})
		if failed != "" {
			t.Fatal(failed)
		}
		u, _ := out["unpushed"].(map[string]any)
		if u == nil {
			t.Fatalf("call %d (%s): an unpushed acceptance is state and is listed every time: %v", i, tool, out)
		}
		if acc := u["acceptances"].([]any); len(acc) != 1 || acc[0] != "S-0007" {
			t.Errorf("%s acceptances: %v", tool, u["acceptances"])
		}
		if tags := u["tags"].([]any); len(tags) != 1 || tags[0] != "cli/v1.1.0" {
			t.Errorf("%s tags: %v", tool, u["tags"])
		}
		if u["upstream"] != "origin/main" || u["command"] != "flai push --pending" {
			t.Errorf("%s: %v", tool, u)
		}
	}

	if out, _ := connect(nil).call(t, "inbox", map[string]any{}); out["unpushed"] != nil {
		t.Errorf("a server with no runner says nothing about it: %v", out["unpushed"])
	}

	git(root, "push", "-q", "origin", "main", "cli/v1.1.0")
	if out, _ := f.call(t, "inbox", map[string]any{}); out["unpushed"] != nil {
		t.Errorf("after the push it is no longer true: %v", out["unpushed"])
	}
}
