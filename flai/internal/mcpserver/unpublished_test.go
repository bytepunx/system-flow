package mcpserver

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0195, ADR-0067: publishing is the operator's, so the instructions give
// no agent a duty to push, in a project or in a folder of them.
func TestTheInstructionsGiveNoDutyToPush(t *testing.T) {
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	for name, in := range map[string]string{
		"project": setup(t).cs.InitializeResult().Instructions,
		"folder":  folderSetup(t, root).cs.InitializeResult().Instructions,
	} {
		for _, never := range []string{"flai push", "git push", "unpushed", "push it"} {
			if strings.Contains(in, never) {
				t.Errorf("%s instructions say %q: %s", name, never, in)
			}
		}
		for _, want := range []string{"Publishing is the operator's", "flai release --pending", "only when the operator asks (ADR-0067)", "unpublished"} {
			if !strings.Contains(in, want) {
				t.Errorf("%s instructions lack %q: %s", name, want, in)
			}
		}
	}
}

// S-0195, ADR-0067: inbox and the board tool list what is accepted and not
// yet published, on every call while it is true, so that whoever the
// operator asks to publish can tell what is pending; never as unpushed.
func TestInboxAndBoardReportWhatIsUnpublished(t *testing.T) {
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
	write := func(rel, body string) {
		t.Helper()
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "cli"} {
		write(filepath.Join(d, ".gitkeep"), "")
	}
	write("system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n    kind: go\n")
	write(".gitignore", ".flai-cache/\n")
	git("init", "-q", "-b", "main")
	git("add", "-A")
	git("commit", "-q", "-m", "init")
	git("tag", "-a", "cli/v1.0.0", "-m", "cli 1.0.0")

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

	if out, _ := f.call(t, "inbox", map[string]any{}); out["unpublished"] != nil {
		t.Fatalf("nothing is accepted yet: %v", out["unpublished"])
	}
	// accepted as flai accept leaves it: archived, done, committed, untagged
	write("wip/archive/kanban/stories/S-0007-a-command.md", "---\nid: S-0007\ntype: story\nnature: feature\ntitle: A command\nstatus: done\n---\n# S-0007 A command\n")
	write("cli/a.go", "package main\n")
	git("add", "-A")
	git("commit", "-q", "-m", "chore: [S-0007] accept and archive")

	for i, tool := range []string{"inbox", "inbox", "board"} {
		out, failed := f.call(t, tool, map[string]any{})
		if failed != "" {
			t.Fatal(failed)
		}
		u, _ := out["unpublished"].([]any)
		if len(u) != 1 || u[0] != "S-0007" {
			t.Errorf("call %d (%s): what is accepted and not yet published is listed every time: %v", i, tool, out["unpublished"])
		}
		if _, ok := out["unpushed"]; ok {
			t.Errorf("call %d (%s): unpushed is gone: %v", i, tool, out["unpushed"])
		}
	}

	if out, _ := connect(nil).call(t, "inbox", map[string]any{}); out["unpublished"] != nil {
		t.Errorf("a server with no runner says nothing about it: %v", out["unpublished"])
	}

	git("tag", "-a", "cli/v1.1.0", "-m", "cli 1.1.0") // published
	if out, _ := f.call(t, "inbox", map[string]any{}); out["unpublished"] != nil {
		t.Errorf("once published it is no longer listed: %v", out["unpublished"])
	}
}
