package mcpserver

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

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/verify"
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
	return testFixtureWith(t, nil)
}

// testFixtureWith is testFixture with the server's options changed by with
// as well, when it is not nil.
func testFixtureWith(t *testing.T, with func(*Options)) (*fixture, func(dir string, args ...string)) {
	t.Helper()
	for _, bin := range []string{"git", "sh", "sleep"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skip(bin + " not installed")
		}
	}
	f := setupWith(t, func(o *Options) {
		o.Runner = execx.System{}
		if with != nil {
			with(o)
		}
	})
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
	t.Parallel()
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
	t.Parallel()
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

// i0141Manifest is testManifest with a tier whose command holds a
// placeholder this flai does not fill, as main's manifest held {importers}
// for flai 1.40.1 (I-0141).
var i0141Manifest = strings.Replace(testManifest, `command: [sh, -c, "exit 0"]`, `command: [sh, -c, "exit 0", "{dependents}"]`, 1)

// i0141Problem is what this flai finds wrong in i0141Manifest.
const i0141Problem = `tests[0].command (tier "ok") argument {dependents} is not a placeholder flai fills`

// I-0141, ADR-0158: a flai older than the newest flai release in the
// project's history that refuses the manifest's tests says so, with both
// versions and the upgrade, rather than telling the agent to fix the
// manifest; a flai at that release refuses them naming the fix.
func TestTestSaysTheFlaiIsOlderThanTheProjectWhenTheTiersRefuse(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		version string
		want    []string
		not     string
	}{
		{"1.40.1", []string{i0141Problem, "this flai is 1.40.1, older than flai 1.49.0 in the project's history", "upgrade it with flai host upgrade"}, "flai manifest set tests="},
		{"1.49.0", []string{i0141Problem, "flai manifest set tests="}, "upgrade it"},
	} {
		f, git := testFixtureWith(t, func(o *Options) { o.Version = c.version })
		git(f.repo.Root, "tag", "-a", "flai/v1.49.0", "-m", "flai 1.49.0")
		if err := os.WriteFile(filepath.Join(f.repo.Root, "system-flow.yaml"), []byte(i0141Manifest), 0o644); err != nil {
			t.Fatal(err)
		}
		_, failed := f.call(t, "test", map[string]any{"paths": []string{"ok"}})
		for _, want := range c.want {
			if !strings.Contains(failed, want) {
				t.Errorf("flai %s: %q, want it to say %q", c.version, failed, want)
			}
		}
		if strings.Contains(failed, c.not) {
			t.Errorf("flai %s: %q, want it not to say %q", c.version, failed, c.not)
		}
	}
}

// S-0273: a call cancelled while a tier runs stops the tier's processes and
// answers at once rather than when the tier would have finished.
func TestTestStopsTheTierWhenTheCallIsCancelled(t *testing.T) {
	t.Parallel()
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

// i0260Manifest is testManifest with a tier that, like gofmt's, is given
// every file it selects as an argument and passes.
var i0260Manifest = testManifest + `  - name: files
    command: [sh, -c, "exit 0", sh, "{files}"]
    paths: ["many/**"]
`

// I-0260: a tier given many files that passes no longer carries them in the
// answer, which keeps within ctxpack.PartLimit bytes of JSON: the passing tier
// gives the number of arguments its command ran with, paths is cut to its
// first entries and paths_omitted counts the rest, and the tier that fails
// keeps its command.
func TestTestAnswersACommandOnlyForTheTierThatFailedAndFitsTheLimit(t *testing.T) {
	t.Parallel()
	f, git := testFixture(t)
	root := f.repo.Root
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(i0260Manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "many"), 0o755); err != nil {
		t.Fatal(err)
	}
	const files = 1500
	for i := range files {
		name := fmt.Sprintf("a-file-whose-name-is-long-enough-to-matter-%04d.txt", i)
		if err := os.WriteFile(filepath.Join(root, "many", name), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(root, "add", "-A")
	git(root, "commit", "-q", "-m", "many")

	out, failed := f.call(t, "test", map[string]any{"paths": []string{"many"}})
	if failed != "" {
		t.Fatal(failed)
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > ctxpack.PartLimit {
		t.Errorf("a passing run over %d files answers %d bytes, over %d", files, len(data), ctxpack.PartLimit)
	}
	if passed, _ := out["passed"].(bool); !passed {
		t.Fatalf("the run did not pass: %v", out["tiers"])
	}
	tier := out["tiers"].([]any)[0].(map[string]any)
	if _, ok := tier["command"]; ok {
		t.Errorf("the passing tier carries its command: %d arguments", len(tier["command"].([]any)))
	}
	if args, _ := tier["args"].(float64); int(args) != files+4 {
		t.Errorf("args = %v, want %d: the command's four words and every file", tier["args"], files+4)
	}
	paths, _ := out["paths"].([]any)
	omitted, _ := out["paths_omitted"].(float64)
	if len(paths) == 0 || int(omitted) == 0 || len(paths)+int(omitted) != files {
		t.Errorf("%d paths and %v omitted, want some of each, %d in all", len(paths), out["paths_omitted"], files)
	}

	out, failed = f.call(t, "test", map[string]any{"paths": []string{"ok", "bad"}})
	if failed != "" {
		t.Fatal(failed)
	}
	for _, x := range out["tiers"].([]any) {
		tier := x.(map[string]any)
		_, has := tier["command"]
		if want := tier["state"] == "failed"; has != want {
			t.Errorf("tier %s, %s: command given %v, want %v", tier["name"], tier["state"], has, want)
		}
	}
	if _, ok := out["paths_omitted"]; ok {
		t.Errorf("a small answer cuts its paths: %v", out)
	}
}

// fitTest cuts the paths first and then the failed tier's command, only as
// far as the limit needs, and leaves a result within it as it is.
func TestFitTestCutsThePathsThenTheFailedCommand(t *testing.T) {
	t.Parallel()
	many := func(prefix string, n int) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = fmt.Sprintf("%s/file-%04d.txt", prefix, i)
		}
		return out
	}
	one := 1
	res := verify.Result{Paths: many("p", 2000), Tiers: []verify.TierResult{
		{Name: "lint", Command: append([]string{"lint"}, many("p", 2000)...), State: verify.Passed},
		{Name: "test", Command: append([]string{"test"}, many("q", 4000)...), State: verify.Failed, ExitCode: &one, Findings: []verify.Finding{{Message: "it broke"}}},
		{Name: "smoke", Command: []string{"smoke"}, State: verify.NotReached},
	}}
	small := verify.Result{Passed: true, Paths: []string{"a"}, Tiers: []verify.TierResult{{Name: "lint", Command: []string{"lint", "a"}, State: verify.Passed}}}
	if got := fitTest(small, ctxpack.PartLimit); len(got.Paths) != 1 || got.PathsOmitted != 0 || got.Tiers[0].Args != 2 || got.Tiers[0].Command != nil {
		t.Errorf("a small result: %+v", got)
	}

	got := fitTest(res, ctxpack.PartLimit)
	if size := testSize(got); size > ctxpack.PartLimit {
		t.Errorf("the answer is %d bytes, over %d", size, ctxpack.PartLimit)
	}
	if len(got.Paths) != 0 || got.PathsOmitted != 2000 {
		t.Errorf("paths %d, omitted %d: the paths go before the failed command", len(got.Paths), got.PathsOmitted)
	}
	lint, test, smoke := got.Tiers[0], got.Tiers[1], got.Tiers[2]
	if lint.Command != nil || lint.Args != 2001 || smoke.Command != nil || smoke.Args != 1 {
		t.Errorf("the tiers that did not fail: %+v %+v", lint, smoke)
	}
	if len(test.Command) == 0 || test.Command[0] != "test" || len(test.Command)+test.CommandOmitted != 4001 || test.CommandOmitted == 0 || test.Findings[0].Message != "it broke" {
		t.Errorf("the failed tier: %d arguments, %d omitted, findings %v", len(test.Command), test.CommandOmitted, test.Findings)
	}
	n := len(test.Command) + 1
	got.Tiers[1].Command, got.Tiers[1].CommandOmitted = res.Tiers[1].Command[:n], 4001-n
	if testSize(got) <= ctxpack.PartLimit {
		t.Errorf("the command was cut to %d arguments, though %d fit", n-1, n)
	}
}
