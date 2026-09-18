package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
)

// versionRunner answers "git version" and nothing else.
type versionRunner struct {
	execx.System
	version string
}

func (v versionRunner) Run(_ string, name string, args ...string) (string, error) {
	if name == "git" && len(args) == 1 && args[0] == "version" {
		return v.version, nil
	}
	return "", nil
}

// The opt-in decides, never the git version on its own (ADR-0022).
func TestRelativeWorktreesIsAnOptIn(t *testing.T) {
	for _, c := range []struct {
		name, config, git string
		want              bool
		warns             bool
	}{
		{"off with a new git", `{}`, "git version 2.54.0", false, false},
		{"explicitly off", `{"worktrees":{"relative_paths":false}}`, "git version 2.54.0", false, false},
		{"on with a new git", `{"worktrees":{"relative_paths":true}}`, "git version 2.48.0", true, false},
		{"on with an old git", `{"worktrees":{"relative_paths":true}}`, "git version 2.47.3", false, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			cfg := filepath.Join(t.TempDir(), "cfg.json")
			if err := os.WriteFile(cfg, []byte(c.config), 0o600); err != nil {
				t.Fatal(err)
			}
			t.Setenv("FLAI_CONFIG", cfg)
			t.Setenv("LOG_FORMAT", "json")
			var errOut bytes.Buffer
			a := &app{out: &bytes.Buffer{}, errOut: &errOut, runner: versionRunner{version: c.git}}
			if got := a.relativeWorktrees(); got != c.want {
				t.Errorf("relativeWorktrees() = %v, want %v", got, c.want)
			}
			log := errOut.String()
			warned := strings.Contains(log, `"setting":"worktrees.relative_paths"`) && strings.Contains(log, `"git":"2.47"`) && strings.Contains(log, `"needs":"2.48"`)
			if warned != c.warns {
				t.Errorf("warning with the setting and both versions: got %v, want %v; log: %s", warned, c.warns, log)
			}
		})
	}
}

// With the key set through flai config and a git that supports it, stream
// open links the worktree with relative paths, so it works wherever the
// repository is mounted.
func TestStreamOpenRelativeWorktree(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	v, err := gitver.Installed(execx.System{})
	if err != nil {
		t.Fatal(err)
	}
	if !v.AtLeast(gitver.RelativeWorktrees) {
		t.Skipf("git %s on PATH; relative worktree paths need %s", v, gitver.RelativeWorktrees)
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{
		{"config", "set", "worktrees.relative_paths", "true"},
		{"epic", "new", "Epic"},
		{"story", "new", "Relative", "--epic", "E-0001"},
		{"stream", "open", "S-0001"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	link, err := os.ReadFile(filepath.Join(root, ".flai-cache", "worktrees", "S-0001", ".git"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(link), "gitdir: ../") {
		t.Errorf("worktree link should be relative, got %q", link)
	}
	if got := gitIn(t, root, "config", "--get", "extensions.relativeWorktrees"); got != "true" {
		t.Errorf("extensions.relativeWorktrees = %q", got)
	}
}
