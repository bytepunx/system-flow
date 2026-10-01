package mcpserver

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0181: inbox says, on every call, when the flai serving it is older than
// the newest flai release tagged in the project's history, with the upgrade.
func TestInboxSaysTheFlaiIsOlderThanTheProject(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	git("init", "-q", "-b", "main")
	git("add", "-A")
	git("commit", "-q", "-m", "init")
	git("tag", "-a", "flai/v1.27.0", "-m", "flai 1.27.0")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	inbox := func(version string) map[string]any {
		srv := New(Options{Repo: repo, Agent: "claude", Version: version, Now: func() time.Time { return t0 }, Runner: execx.System{}})
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
		out, failed := (&fixture{repo: repo, cs: cs}).call(t, "inbox", map[string]any{})
		if failed != "" {
			t.Fatal(failed)
		}
		return out
	}
	for range 2 {
		o, _ := inbox("1.26.4")["flai_outdated"].(map[string]any)
		if o == nil || o["running"] != "1.26.4" || o["newest"] != "1.27.0" || !strings.HasPrefix(o["upgrade"].(string), "flai host upgrade") {
			t.Fatalf("flai 1.26.4 under flai/v1.27.0: %v", o)
		}
	}
	for _, v := range []string{"1.27.0", "dev"} {
		if o := inbox(v)["flai_outdated"]; o != nil {
			t.Errorf("flai %s: %v", v, o)
		}
	}
}

// S-0181: in a folder, a project whose minimum flai is above this one is
// left out, and the log says why once, not at every look.
func TestAFolderLeavesOutAProjectThatNeedsANewerFlai(t *testing.T) {
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	makeProject(t, filepath.Join(root, "beta"), "beta")
	man := filepath.Join(root, "beta", "system-flow.yaml")
	m, _ := os.ReadFile(man)
	_ = os.WriteFile(man, append(m, []byte("flai:\n  minimum: 9.0.0\n")...), 0o644)
	was := buildinfo.Version
	t.Cleanup(func() { buildinfo.Version = was })
	buildinfo.Version = "1.26.4"

	var logs bytes.Buffer
	srv := New(Options{Folder: root, Agent: "claude", Version: "1.26.4", Now: func() time.Time { return t0 }, Rescan: time.Nanosecond,
		Logger: slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelWarn}))})
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
	f := &folderFixture{root: root, cs: cs}
	for range 2 {
		out, failed := f.call(t, "inbox", map[string]any{})
		if failed != "" || strings.Join(projectKeys(out), " ") != "alpha" {
			t.Fatalf("inbox: %v %s", out, failed)
		}
	}
	got := logs.String()
	if strings.Count(got, "project not served") != 1 || !strings.Contains(got, "needs flai 9.0.0 or newer, and this is flai 1.26.4") {
		t.Errorf("log:\n%s", got)
	}
}
