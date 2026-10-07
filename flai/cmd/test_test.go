package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/verify"
)

// testTiersManifest declares the fixture's tiers: go test for the Go module
// in svc/, and a plain tier for web/ that checks each file it is given is
// there.
const testTiersManifest = `version: 1
name: t
key: t
layout:
  design: design
  docs: docs
  wip: wip
tests:
  - name: go-test
    command: [go, test, -json, -count=1, "{packages}"]
    dir: svc
    paths: ["svc/**/*.go"]
    format: go-test-json
  - name: web
    command: [sh, -c, 'for f in "$@"; do test -f "$f"; done', web, "{files}"]
    dir: web
    paths: ["web/**"]
`

// testFixture is a git repository on main, committed, holding testTiersManifest,
// a Go module in svc/ with a passing package calc and a package bad with two
// failing tests, and web/x.ts. Integration: it runs real git and go.
func testFixture(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git and go")
	}
	for _, tool := range []string{"git", "go", "sh"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skip(tool + " not installed")
		}
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	files := map[string]string{
		"system-flow.yaml":      testTiersManifest,
		"svc/go.mod":            "module example.com/svc\n\ngo 1.21\n",
		"svc/calc/calc.go":      "package calc\n\n// Add adds.\nfunc Add(a, b int) int { return a + b }\n",
		"svc/calc/calc_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) {\n\tif Add(1, 2) != 3 {\n\t\tt.Fatal(\"add\")\n\t}\n}\n",
		"svc/bad/bad.go":        "package bad\n",
		"svc/bad/bad_test.go":   "package bad\n\nimport \"testing\"\n\nfunc TestOne(t *testing.T) {\n\tt.Error(\"one is wrong\")\n}\n\nfunc TestTwo(t *testing.T) {\n\tt.Error(\"two is wrong\")\n}\n",
		"web/x.ts":              "export const x = 1;\n",
		".gitignore":            ".flai-cache/\n",
	}
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	return root
}

// testResult runs flai test --json in dir and decodes its answer.
func testResult(t *testing.T, dir string, wantCode int, args ...string) verify.Result {
	t.Helper()
	out, errOut, code := runIn(t, dir, append([]string{"test", "--json"}, args...)...)
	if code != wantCode {
		t.Fatalf("flai test %v: exit %d, want %d\n%s%s", args, code, wantCode, out, errOut)
	}
	var res verify.Result
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("flai test %v: not JSON: %v\n%s", args, err, out)
	}
	return res
}

// tierNames are the names of the tiers a run reached or passed over.
func tierNames(res verify.Result) []string {
	var names []string
	for _, tr := range res.Tiers {
		names = append(names, tr.Name)
	}
	return names
}

func TestFlaiTestRunsOnlyTheTiersAPathOrPackageSelects(t *testing.T) {
	root := testFixture(t)

	res := testResult(t, root, 0, "svc/calc")
	if got := strings.Join(tierNames(res), ","); got != "go-test" || !res.Passed || res.Tiers[0].State != verify.Passed {
		t.Errorf("a package's folder ran %s: %+v", got, res)
	}
	if want := []string{"svc/calc/calc.go", "svc/calc/calc_test.go"}; strings.Join(res.Paths, ",") != strings.Join(want, ",") {
		t.Errorf("paths %v, want %v", res.Paths, want)
	}
	if got := strings.Join(res.Tiers[0].Command, " "); !strings.HasSuffix(got, " ./calc") {
		t.Errorf("go-test ran %q, not for ./calc", got)
	}

	res = testResult(t, root, 0, "web/x.ts")
	if got := strings.Join(tierNames(res), ","); got != "web" || !res.Passed {
		t.Errorf("a web file ran %s: %+v", got, res)
	}

	// arguments are relative to the working directory, inside the checkout
	res = testResult(t, filepath.Join(root, "svc"), 0, "calc")
	if got := strings.Join(tierNames(res), ","); got != "go-test" || !res.Passed {
		t.Errorf("calc from svc/ ran %s: %+v", got, res)
	}

	out, _, code := runIn(t, root, "test", "web/x.ts")
	if code != 0 || !strings.HasPrefix(out, "passed web (") {
		t.Errorf("text: exit %d\n%s", code, out)
	}
}

func TestFlaiTestAnswersAFailingGoTestByNamePathLineAndMessageWithinTheCap(t *testing.T) {
	root := testFixture(t)

	res := testResult(t, root, 1, "--max", "1", "svc/bad")
	if res.Passed || len(res.Tiers) != 1 {
		t.Fatalf("result %+v", res)
	}
	tr := res.Tiers[0]
	if tr.State != verify.Failed || len(tr.Findings) != 1 || tr.Omitted != 1 {
		t.Fatalf("tier %+v", tr)
	}
	f := tr.Findings[0]
	if (f.Name != "TestOne" && f.Name != "TestTwo") || f.Path != "svc/bad/bad_test.go" || f.Line == 0 || !strings.Contains(f.Message, "is wrong") {
		t.Errorf("finding %+v", f)
	}

	out, _, code := runIn(t, root, "test", "--max", "1", "svc/bad")
	if code != 1 {
		t.Fatalf("text: exit %d\n%s", code, out)
	}
	line := f.Path + ":" + strconv.Itoa(f.Line) + " " + f.Name + ": "
	if !strings.HasPrefix(out, "failed go-test (") || !strings.Contains(out, line) || !strings.Contains(out, "… 1 more findings left out") {
		t.Errorf("text does not name %q and the one left out:\n%s", line, out)
	}
}

func TestFlaiTestWithNoArgumentsRunsForWhatTheCheckoutChanged(t *testing.T) {
	root := testFixture(t)
	if err := os.WriteFile(filepath.Join(root, "web", "x.ts"), []byte("export const x = 2;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res := testResult(t, root, 0)
	if got := strings.Join(tierNames(res), ","); got != "web" || strings.Join(res.Paths, ",") != "web/x.ts" {
		t.Errorf("the change to web/x.ts ran %s for %v", got, res.Paths)
	}

	// a story's worktree counts its branch's commits against main
	wt := filepath.Join(t.TempDir(), "wt")
	gitIn(t, root, "worktree", "add", "-q", "-b", "story/S-0001", wt)
	if err := os.WriteFile(filepath.Join(wt, "svc", "calc", "calc.go"), []byte("package calc\n\n// Add adds.\nfunc Add(a, b int) int { return b + a }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "commit", "-q", "-am", "calc")
	res = testResult(t, filepath.Join(wt, "web"), 0)
	if got := strings.Join(tierNames(res), ","); got != "go-test" || strings.Join(res.Paths, ",") != "svc/calc/calc.go" {
		t.Errorf("the worktree's commit ran %s for %v", got, res.Paths)
	}
}

// S-0311: flai test run as the orchestrator runs its tiers under the verify
// role, so a tier whose tests make flai's writes is not refused as the
// orchestrator (TH-0260).
func TestFlaiTestRunsItsTiersUnderTheVerifyRoleForTheOrchestrator(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	manifest, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"system-flow.yaml": string(manifest) + verifyTier(roleTier(t)),
		"docs/x.md":        "x\n",
	}
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("FLAI_ROLE", "orchestrate")

	res := testResult(t, root, 0, "docs/x.md")
	if got := strings.Join(tierNames(res), ","); got != "unit" || !res.Passed {
		t.Errorf("ran %s: %+v", got, res)
	}
}

func TestFlaiTestExitsTwoWhenItCannotAnswer(t *testing.T) {
	root := testFixture(t)
	cases := map[string][]string{
		"a path outside the checkout":    {"test", ".."},
		"a path that is not there":       {"test", "web/nope.ts"},
		"a flag flai test does not take": {"test", "--nope"},
		"a cap of no findings":           {"test", "--max", "0", "web"},
	}
	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			out, errOut, code := runIn(t, root, args...)
			if code != exitTestUnusable || out != "" || errOut == "" {
				t.Errorf("exit %d, want %d\nout %q\nerr %q", code, exitTestUnusable, out, errOut)
			}
		})
	}

	bad := strings.Replace(testTiersManifest, "format: go-test-json", "format: junit", 1)
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "test", "web/x.ts")
	if code != exitTestUnusable || out != "" || !strings.Contains(errOut, "tests[0].format") {
		t.Errorf("a tier not declared validly: exit %d, want %d, naming the field\nout %q\nerr %q", code, exitTestUnusable, out, errOut)
	}
}
