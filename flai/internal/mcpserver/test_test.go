package mcpserver

import (
	"context"
	"fmt"
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

// testManifest is the fixture's manifest with tiers that need only sh: ok
// passes for ok/**, bad fails for bad/** with a line naming the failure, and
// slow, for slow/**, writes its process ID to slow/pid and sleeps until it is
// stopped.
const testManifest = `version: 1
name: t
key: t
layout:
  design: design
  docs: docs
  wip: wip
tests:
  - name: ok
    command: [sh, -c, "exit 0"]
    paths: ["ok/**"]
  - name: bad
    command: [sh, -c, "echo 'bad/b.txt:3: it broke'; exit 1"]
    paths: ["bad/**"]
  - name: slow
    command: [sh, -c, "echo $$ > slow/pid; exec sleep 30"]
    paths: ["slow/**"]
`

// testFixture is the server's fixture made a git repository on main whose
// manifest declares the tiers above, with files each tier selects.
func testFixture(t *testing.T) (*fixture, func(dir string, args ...string)) {
	t.Helper()
	for _, bin := range []string{"git", "sh", "sleep"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skip(bin + " not installed")
		}
	}
	f := setupWith(t, func(o *Options) { o.Runner = execx.System{} })
	root := f.repo.Root
	git := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	write := func(rel, data string) {
		t.Helper()
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("system-flow.yaml", testManifest)
	write(".gitignore", ".flai-cache/\n")
	write("ok/a.txt", "a\n")
	write("bad/b.txt", "b\n")
	write("slow/c.txt", "c\n")
	git(root, "init", "-q", "-b", "main")
	git(root, "add", "-A")
	git(root, "commit", "-q", "-m", "init")
	return f, git
}

// result decodes the tool's answer: passed, paths, and each tier's state and
// findings' messages.
type testAnswer struct {
	passed bool
	paths  string
	states string
	found  string
}

func answer(t *testing.T, out map[string]any) testAnswer {
	t.Helper()
	var a testAnswer
	a.passed, _ = out["passed"].(bool)
	paths, _ := out["paths"].([]any)
	a.paths = strings.Trim(fmt.Sprint(paths), "[]")
	tiers, ok := out["tiers"].([]any)
	if !ok {
		t.Fatalf("no tiers in %v", out)
	}
	var states, found []string
	for _, x := range tiers {
		tier := x.(map[string]any)
		states = append(states, fmt.Sprintf("%s=%s", tier["name"], tier["state"]))
		fs, _ := tier["findings"].([]any)
		for _, y := range fs {
			found = append(found, fmt.Sprint(y.(map[string]any)["message"]))
		}
		if _, ok := tier["duration_ms"]; !ok {
			t.Errorf("tier %s has no duration: %v", tier["name"], tier)
		}
	}
	a.states, a.found = strings.Join(states, " "), strings.Join(found, "\n")
	return a
}

// S-0273: the tool is listed; for failing paths it answers the failing tier's
// findings, for passing ones it answers pass, and every tier the paths do not
// select stays out of the answer.
func TestTestAnswersPassOrTheFirstFindings(t *testing.T) {
	f, _ := testFixture(t)
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	listed := false
	for _, tool := range res.Tools {
		if tool.Name == "test" {
			listed = true
			for _, want := range []string{"cheapest first", "findings", "worktree", "go test, vitest, golangci-lint, or gofmt by hand", "flai test"} {
				if !strings.Contains(tool.Description, want) {
					t.Errorf("the description does not say %q: %s", want, tool.Description)
				}
			}
		}
	}
	if !listed {
		t.Fatal("the tool test is not listed")
	}

	out, failed := f.call(t, "test", map[string]any{"paths": []string{"bad/b.txt", "ok"}})
	if failed != "" {
		t.Fatal(failed)
	}
	a := answer(t, out)
	if a.passed || a.paths != "bad/b.txt ok/a.txt" || a.states != "ok=passed bad=failed" || !strings.Contains(a.found, "bad/b.txt:3: it broke") {
		t.Errorf("failing paths: %+v", a)
	}

	out, failed = f.call(t, "test", map[string]any{"paths": []string{filepath.Join(f.repo.Root, "ok", "a.txt")}})
	if failed != "" {
		t.Fatal(failed)
	}
	if a := answer(t, out); !a.passed || a.paths != "ok/a.txt" || a.states != "ok=passed" || a.found != "" {
		t.Errorf("passing paths, one given absolute: %+v", a)
	}
}

// S-0273: with story it runs in that story's worktree, as its branch has the
// manifest and the files, and with no paths it takes what the branch changed
// against main and what it has not committed; with none it runs in the main
// checkout. A story without a worktree is refused, saying how to open one.
func TestTestRunsInTheStorysWorktree(t *testing.T) {
	f, git := testFixture(t)
	root := f.repo.Root
	wt := f.repo.WorktreePath(f.story.ID)
	git(root, "worktree", "add", "-q", "-b", "story/"+f.story.ID, wt, "main")
	// on the story's branch the ok tier fails, and ok/a.txt changes
	manifest := strings.Replace(testManifest, `[sh, -c, "exit 0"]`, `[sh, -c, "echo 'from the worktree'; exit 1"]`, 1)
	if err := os.WriteFile(filepath.Join(wt, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt, "ok", "a.txt"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(wt, "add", "-A")
	git(wt, "commit", "-q", "-m", "change")
	if err := os.WriteFile(filepath.Join(wt, "ok", "new.txt"), []byte("not committed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, failed := f.call(t, "test", map[string]any{"story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	a := answer(t, out)
	if a.passed || a.paths != "ok/a.txt ok/new.txt system-flow.yaml" || a.states != "ok=failed" || !strings.Contains(a.found, "from the worktree") {
		t.Errorf("the story's worktree, its changes: %+v", a)
	}

	out, failed = f.call(t, "test", map[string]any{"paths": []string{"ok/a.txt"}})
	if failed != "" {
		t.Fatal(failed)
	}
	if a := answer(t, out); !a.passed || a.states != "ok=passed" {
		t.Errorf("the main checkout: %+v", a)
	}
	out, failed = f.call(t, "test", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if a := answer(t, out); !a.passed || a.paths != "" || a.states != "" {
		t.Errorf("the main checkout changed nothing: %+v", a)
	}

	other, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Other", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{"story": other.ID}, other.ID + " has no worktree at .flai-cache/worktrees/" + other.ID + "; open it with flai stream open " + other.ID},
		{map[string]any{"story": "S-0099"}, "S-0099 is not a work item here"},
		{map[string]any{"paths": []string{"nowhere"}}, "nowhere is not a file or folder in the checkout"},
		{map[string]any{"paths": []string{"ok"}, "max": -1}, "max -1 is not a number of findings"},
	} {
		if _, failed := f.call(t, "test", c.args); !strings.Contains(failed, c.want) {
			t.Errorf("%v: %q, want %q", c.args, failed, c.want)
		}
	}

	bad := strings.Replace(testManifest, "    paths: [\"ok/**\"]\n", "", 1)
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "test", map[string]any{"paths": []string{"ok"}}); !strings.Contains(failed, "tests[0].paths") || !strings.Contains(failed, "flai manifest set tests=") {
		t.Errorf("a manifest whose tests are not valid: %q", failed)
	}
}

// S-0273: a call cancelled while a tier runs stops the tier's processes and
// answers at once rather than when the tier would have finished.
func TestTestStopsTheTierWhenTheCallIsCancelled(t *testing.T) {
	f, _ := testFixture(t)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	start := time.Now()
	res, err := f.cs.CallTool(ctx, &mcp.CallToolParams{Name: "test", Arguments: map[string]any{"paths": []string{"slow"}}})
	if err == nil && !res.IsError {
		t.Fatalf("a cancelled call answered: %v", res.StructuredContent)
	}
	if took := time.Since(start); took > 10*time.Second {
		t.Errorf("the cancelled call took %s, as long as the tier", took)
	}
	pid, err := os.ReadFile(filepath.Join(f.repo.Root, "slow", "pid"))
	if err != nil {
		t.Fatalf("the slow tier did not start: %v", err)
	}
	alive := func() bool {
		return exec.Command("sh", "-c", "kill -0 "+strings.TrimSpace(string(pid))+" 2>/dev/null").Run() == nil
	}
	for deadline := time.Now().Add(10 * time.Second); alive(); time.Sleep(50 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatalf("the tier's process %s still runs after the call was cancelled", strings.TrimSpace(string(pid)))
		}
	}
}
