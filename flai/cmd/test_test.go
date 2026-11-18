package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
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

// I-0141, ADR-0158: flai test refuses tiers it does not take, here a
// placeholder it does not fill, naming both versions and the upgrade when
// the checkout's history holds a newer flai release, and the manifest's fix
// when the running flai is that release.
// The refusal comes before any tier runs, so the project needs only git.
func TestFlaiTestSaysTheFlaiIsOlderThanTheProjectWhenTheTiersRefuse(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	bad := strings.Replace(testTiersManifest, `web, "{files}"]`, `web, "{dependents}"]`, 1)
	for name, body := range map[string]string{"system-flow.yaml": bad, "web/x.ts": "export const x = 1;\n", ".gitignore": ".flai-cache/\n"} {
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
	gitIn(t, root, "tag", "-a", "flai/v1.49.0", "-m", "flai 1.49.0")
	was := buildinfo.Version
	t.Cleanup(func() { buildinfo.Version = was })
	// the log line quotes the tier's name, so the problem is matched without it
	const problem = "argument {dependents} is not a placeholder flai fills"
	for _, c := range []struct {
		version string
		want    []string
		not     string
	}{
		{"1.40.1", []string{"tests[1].command", problem, "this flai is 1.40.1, older than flai 1.49.0 in the project's history", "upgrade it with flai host upgrade"}, "flai manifest set tests="},
		{"1.49.0", []string{"tests[1].command", problem, "flai manifest set tests="}, "upgrade it"},
	} {
		buildinfo.Version = c.version
		out, errOut, code := runIn(t, root, "test", "web/x.ts")
		if code != exitTestUnusable || out != "" {
			t.Errorf("flai %s: exit %d, want %d\nout %q\nerr %q", c.version, code, exitTestUnusable, out, errOut)
		}
		for _, want := range c.want {
			if !strings.Contains(errOut, want) {
				t.Errorf("flai %s: %q, want it to say %q", c.version, errOut, want)
			}
		}
		if strings.Contains(errOut, c.not) {
			t.Errorf("flai %s: %q, want it not to say %q", c.version, errOut, c.not)
		}
	}
}

// namedTiersManifest is tempProject's manifest with a tier named docs for
// the docs folder, with an all_command; an all_only tier; a manual tier the
// docs folder's paths would select; and a manual tier that fails.
const namedTiersManifest = `version: 1
name: t
key: t
layout:
  design: design
  docs: docs
  wip: wip
tests:
  - name: docs
    command: [sh, -c, 'exit 0', docs, "{files}"]
    paths: ["docs/**"]
    all_command: [sh, -c, 'exit 0', docs-all]
  - name: smoke
    command: [sh, -c, 'exit 0', smoke]
    all_only: true
  - name: gateway
    command: [sh, -c, 'exit 0', gateway, "{files}"]
    paths: ["docs/**"]
    manual: true
  - name: broken
    command: [sh, -c, 'echo the provider refused; exit 3', broken]
    manual: true
`

// namedTiersFixture is a project, outside git, whose manifest is
// namedTiersManifest, holding docs/x.md.
func namedTiersFixture(t *testing.T) string {
	t.Helper()
	return shTiersFixture(t, namedTiersManifest)
}

// shTiersFixture is a project, outside git, whose manifest is manifest, a
// manifest whose tiers need only sh, holding docs/x.md and web/y.ts.
func shTiersFixture(t *testing.T, manifest string) string {
	t.Helper()
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not installed")
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for name, body := range map[string]string{"system-flow.yaml": manifest, "docs/x.md": "x\n", "web/y.ts": "y\n"} {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// S-0548: a lone argument that names a tier runs that tier alone for the
// whole checkout, with its all_command when it has one, whether it is
// all_only, manual, or neither, and is read as the tier where a folder has
// the same name; a manual tier that fails answers as any failing tier does.
func TestFlaiTestRunsATierByItsName(t *testing.T) {
	t.Parallel()
	root := namedTiersFixture(t)
	for _, c := range []struct {
		arg, tier, command string
	}{
		{"docs", "docs", "sh -c exit 0 docs-all"},
		{"smoke", "smoke", "sh -c exit 0 smoke"},
		{"gateway", "gateway", "sh -c exit 0 gateway ."},
		{" gateway ", "gateway", "sh -c exit 0 gateway ."},
	} {
		res := testResult(t, root, 0, c.arg)
		if len(res.Tiers) != 1 || res.Tiers[0].Name != c.tier || res.Tiers[0].State != verify.Passed || !res.Passed || len(res.Paths) != 0 {
			t.Errorf("flai test %q ran %v: %+v", c.arg, tierNames(res), res)
			continue
		}
		if got := strings.Join(res.Tiers[0].Command, " "); got != c.command {
			t.Errorf("flai test %q ran %q, want %q", c.arg, got, c.command)
		}
	}

	// written as a path, the folder of the tier's name is the folder
	res := testResult(t, root, 0, "./docs")
	if got := strings.Join(tierNames(res), ","); got != "docs" || strings.Join(res.Tiers[0].Command, " ") != "sh -c exit 0 docs docs/x.md" || strings.Join(res.Paths, ",") != "docs/x.md" {
		t.Errorf("./docs ran %s: %+v", got, res)
	}

	res = testResult(t, root, 1, "broken")
	if len(res.Tiers) != 1 || res.Tiers[0].State != verify.Failed || res.Passed || len(res.Tiers[0].Findings) == 0 || !strings.Contains(res.Tiers[0].Findings[0].Message, "the provider refused") {
		t.Errorf("a failing manual tier by name: %+v", res)
	}
	out, _, code := runIn(t, root, "test", "gateway")
	if code != 0 || !strings.HasPrefix(out, "passed gateway (") {
		t.Errorf("text: exit %d\n%s", code, out)
	}
}

// S-0548: a manual tier runs by name alone: paths its own match do not select
// it, and --all leaves it out.
func TestFlaiTestLeavesAManualTierOutOfPathsAndAll(t *testing.T) {
	t.Parallel()
	root := namedTiersFixture(t)
	res := testResult(t, root, 0, "docs/x.md")
	if got := strings.Join(tierNames(res), ","); got != "docs" {
		t.Errorf("docs/x.md ran %s, want docs alone", got)
	}
	res = testResult(t, root, 0, "--all")
	if got := strings.Join(tierNames(res), ","); got != "docs,smoke" {
		t.Errorf("--all ran %s, want docs,smoke", got)
	}
}

// S-0548: a lone argument that is neither a tier nor a path is refused naming
// the tiers, and a tier's name with a path, another name, or --all is refused,
// each with exit 2 before any tier runs.
func TestFlaiTestRefusesATierNameItCannotRunAlone(t *testing.T) {
	t.Parallel()
	root := namedTiersFixture(t)
	for _, c := range []struct {
		name string
		args []string
		want string
	}{
		{"neither a tier nor a path", []string{"nope"}, "nope is neither a tier nor a file or folder; name a tier, one of docs, smoke, gateway, broken, or a file or folder"},
		{"a tier with a path", []string{"gateway", "docs/x.md"}, "gateway names a tier, which runs alone for the whole checkout; give it with no path"},
		{"two tiers", []string{"gateway", "smoke"}, "gateway, smoke each name a tier, and flai test runs one tier by name"},
		{"a tier with --all", []string{"--all", "gateway"}, "gateway names a tier, and --all runs every tier but the manual ones"},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, errOut, code := runIn(t, root, append([]string{"test"}, c.args...)...)
			if code != exitTestUnusable || out != "" || !strings.Contains(errOut, c.want) {
				t.Errorf("exit %d, want %d\nout %q\nerr %q, want it to say %q", code, exitTestUnusable, out, errOut, c.want)
			}
		})
	}
	// with paths alone, one not there is refused as before, naming no tier
	_, errOut, code := runIn(t, root, "test", "docs/x.md", "nope")
	if code != exitTestUnusable || !strings.Contains(errOut, "nope is not a file or folder") {
		t.Errorf("a missing path among paths: exit %d\nerr %q", code, errOut)
	}
}

// suiteTiersManifest is tempProject's manifest with tiers in suites: two
// static tiers, one for docs and one for web; a behavior tier for docs that
// an all_only integration tier covers; a second all_only integration tier;
// and a manual integration tier.
const suiteTiersManifest = `version: 1
name: t
key: t
layout:
  design: design
  docs: docs
  wip: wip
tests:
  - name: fmt
    suite: static
    command: [sh, -c, 'exit 0', fmt, "{files}"]
    paths: ["docs/**"]
  - name: lint
    suite: static
    command: [sh, -c, 'exit 0', lint, "{files}"]
    paths: ["web/**"]
  - name: unit
    suite: behavior
    command: [sh, -c, 'exit 0', unit, "{files}"]
    paths: ["docs/**"]
  - name: integration
    suite: integration
    command: [sh, -c, 'exit 0', integration]
    all_only: true
    covers: [unit]
  - name: e2e
    suite: integration
    command: [sh, -c, 'exit 0', e2e]
    all_only: true
  - name: race
    suite: integration
    command: [sh, -c, 'exit 0', race]
    manual: true
`

// ADR-0312: --suite runs only the tiers of the suite: for paths, those of it
// the paths select, and with --all every one of it but the manual one; a
// tier that a tier of another suite covers runs in a run for its own suite.
func TestFlaiTestSuiteRunsOnlyTheTiersOfTheSuite(t *testing.T) {
	t.Parallel()
	root := shTiersFixture(t, suiteTiersManifest)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--suite", "static", "docs/x.md", "web/y.ts"}, "fmt=passed,lint=passed"},
		{[]string{"--suite", "static", "web/y.ts"}, "lint=passed"},
		{[]string{"--suite", "integration", "--all"}, "integration=passed,e2e=passed"},
		{[]string{"--suite", "behavior", "--all"}, "unit=passed"},
		{[]string{"--all"}, "fmt=passed,lint=passed,unit=skipped,integration=passed,e2e=passed"},
	} {
		res := testResult(t, root, 0, c.args...)
		var got []string
		for _, tr := range res.Tiers {
			got = append(got, tr.Name+"="+string(tr.State))
		}
		if strings.Join(got, ",") != c.want || !res.Passed {
			t.Errorf("flai test %v ran %s, want %s: %+v", c.args, strings.Join(got, ","), c.want, res)
		}
	}
}

// ADR-0312: a suite no tier is in is refused naming the suites there are,
// and --suite beside a tier's name is refused, each with exit 2 before any
// tier runs.
func TestFlaiTestRefusesASuiteItCannotRun(t *testing.T) {
	t.Parallel()
	root := shTiersFixture(t, suiteTiersManifest)
	for _, c := range []struct {
		name string
		args []string
		want string
	}{
		{"a suite no tier is in", []string{"--suite", "slow", "docs/x.md"}, "run the tiers of the suite slow: no tier of the manifest's tests is in it, by its suite or, with none, its name; name one of the suites there are, static, behavior, integration, or give a tier suite: slow in system-flow.yaml"},
		{"a suite with a tier's name", []string{"--suite", "static", "fmt"}, "fmt names a tier, which runs alone whatever its suite, and --suite static runs the tiers of a suite; give the tier alone, or --suite with paths or --all, not both"},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, errOut, code := runIn(t, root, append([]string{"test"}, c.args...)...)
			if code != exitTestUnusable || out != "" || !strings.Contains(errOut, c.want) {
				t.Errorf("exit %d, want %d\nout %q\nerr %q, want it to say %q", code, exitTestUnusable, out, errOut, c.want)
			}
		})
	}
}

// writeBoard writes the baseline workflow file into the project at root as
// its own, with verifyAt in place of the baseline's verify_at line.
func writeBoard(t *testing.T, root, verifyAt string) {
	t.Helper()
	base, err := os.ReadFile(filepath.Join("..", "internal", "workflow", "baseline", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	const line = "verify_at: {batch: all, ci: all}\n"
	if !strings.Contains(string(base), line) {
		t.Fatalf("the baseline workflow has no line %q", line)
	}
	p := filepath.Join(root, "design", "workflows", "board.md")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(base), line, verifyAt, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}

// tierStates are the tiers of a run as name=state, in order.
func tierStates(res verify.Result) string {
	var got []string
	for _, tr := range res.Tiers {
		got = append(got, tr.Name+"="+string(tr.State))
	}
	return strings.Join(got, ",")
}

// S-0812, ADR-0312 §6: --at runs the tiers of the suites the workflow names
// for the point over all paths: with no workflow file the baseline's ci,
// all, which runs every tier but the manual one; with --suite only that
// suite's; the answer names the point and the suites.
func TestFlaiTestAtRunsTheSuitesTheBaselineNamesForThePoint(t *testing.T) {
	t.Parallel()
	root := shTiersFixture(t, suiteTiersManifest)
	res := testResult(t, root, 0, "--at", "ci")
	if got, want := tierStates(res), "fmt=passed,lint=passed,unit=skipped,integration=passed,e2e=passed"; got != want || !res.Passed || len(res.Paths) != 0 {
		t.Errorf("--at ci ran %s, want %s: %+v", got, want, res)
	}
	if res.At != "ci" || strings.Join(res.Suites, ",") != "all" {
		t.Errorf("--at ci answers at %q, suites %q; want ci, all", res.At, res.Suites)
	}
	res = testResult(t, root, 0, "--at", "batch", "--suite", "behavior")
	if got := tierStates(res); got != "unit=passed" || res.At != "batch" || strings.Join(res.Suites, ",") != "behavior" {
		t.Errorf("--at batch --suite behavior ran %s, at %q, suites %q", got, res.At, res.Suites)
	}
	out, _, code := runIn(t, root, "test", "--at", "ci")
	if code != 0 || !strings.HasPrefix(out, "at ci: suites all\npassed fmt (") {
		t.Errorf("text: exit %d\n%s", code, out)
	}
}

// S-0812, ADR-0312 §6: a workflow file's verify_at decides what --at runs:
// ci: [static] runs static's tiers alone, over all paths, and --suite runs
// one suite of them; a suite the point does not name is refused naming the
// point's, as is a point the file gives none.
func TestFlaiTestAtRunsTheSuitesTheWorkflowFileNamesForThePoint(t *testing.T) {
	t.Parallel()
	root := shTiersFixture(t, suiteTiersManifest)
	writeBoard(t, root, "verify_at: {batch: all, ci: [static, integration]}\n")
	res := testResult(t, root, 0, "--at", "ci")
	if got, want := tierStates(res), "fmt=passed,lint=passed,integration=passed,e2e=passed"; got != want || strings.Join(res.Suites, ",") != "static,integration" {
		t.Errorf("--at ci ran %s, suites %q; want %s", got, res.Suites, want)
	}
	res = testResult(t, root, 0, "--at", "ci", "--suite", "static")
	if got := tierStates(res); got != "fmt=passed,lint=passed" || strings.Join(res.Suites, ",") != "static" {
		t.Errorf("--at ci --suite static ran %s, suites %q", got, res.Suites)
	}
	out, _, code := runIn(t, root, "test", "--at", "ci")
	if code != 0 || !strings.HasPrefix(out, "at ci: suites static, integration\n") {
		t.Errorf("text: exit %d\n%s", code, out)
	}
	for _, c := range []struct {
		name string
		args []string
		want string
	}{
		{"a suite the point does not name", []string{"--at", "ci", "--suite", "behavior"}, "--suite behavior is not among the suites the workflow names for ci, static, integration; give one of them, or leave --suite out"},
		{"a point the file gives none", []string{"--at", "publish"}, "the workflow names no suites for publish; give verify_at: {publish: all}, or a list of suites, in design/workflows/board.md"},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, errOut, code := runIn(t, root, append([]string{"test"}, c.args...)...)
			if code != exitTestUnusable || out != "" || !strings.Contains(errOut, c.want) {
				t.Errorf("exit %d, want %d\nout %q\nerr %q, want it to say %q", code, exitTestUnusable, out, errOut, c.want)
			}
		})
	}
}

// S-0812: a suite the point names that no tier is in is refused as --suite
// refuses one.
func TestFlaiTestAtRefusesASuiteThePointNamesThatNoTierIsIn(t *testing.T) {
	t.Parallel()
	root := shTiersFixture(t, suiteTiersManifest)
	writeBoard(t, root, "verify_at: {ci: [static, slow]}\n")
	out, errOut, code := runIn(t, root, "test", "--at", "ci")
	if code != exitTestUnusable || out != "" || !strings.Contains(errOut, "run the tiers of the suite slow: no tier of the manifest's tests is in it") {
		t.Errorf("exit %d, want %d\nout %q\nerr %q", code, exitTestUnusable, out, errOut)
	}
}

// S-0812: --at refuses a point that is none of verify_at's, the publish the
// baseline gives no suites, paths, --all, a tier's name, and a suite no tier
// is in, each with exit 2 before any tier runs.
func TestFlaiTestAtRefusesWhatItCannotRun(t *testing.T) {
	t.Parallel()
	root := shTiersFixture(t, suiteTiersManifest)
	for _, c := range []struct {
		name string
		args []string
		want string
	}{
		{"an unknown point", []string{"--at", "deploy"}, "--at deploy is not a point verify_at gives suites for; give one of batch, publish, ci"},
		{"publish in the baseline", []string{"--at", "publish"}, "the workflow names no suites for publish; give verify_at: {publish: all}, or a list of suites, in design/workflows/board.md"},
		{"paths", []string{"--at", "ci", "docs/x.md"}, "--at ci runs over all paths, so it takes no path; give --at alone, or the paths without it"},
		{"--all", []string{"--at", "ci", "--all"}, "--at ci runs every tier of its suites over all paths, as --all does; give --at alone, not with --all"},
		{"a tier's name", []string{"--at", "ci", "fmt"}, "fmt names a tier, which runs alone, and --at ci runs the tiers of the suites the workflow names for ci; give the tier alone or --at, not both"},
		{"a suite no tier is in", []string{"--at", "ci", "--suite", "slow"}, "run the tiers of the suite slow: no tier"},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, errOut, code := runIn(t, root, append([]string{"test"}, c.args...)...)
			if code != exitTestUnusable || out != "" || !strings.Contains(errOut, c.want) {
				t.Errorf("exit %d, want %d\nout %q\nerr %q, want it to say %q", code, exitTestUnusable, out, errOut, c.want)
			}
		})
	}
}
