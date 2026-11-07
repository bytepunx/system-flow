package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// goRepo is a committed git repository with a Go module whose one test
// passes or fails as asked.
func goRepo(t *testing.T, passing bool) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not installed")
	}
	root := t.TempDir()
	want := "2"
	if !passing {
		want = "3"
	}
	files := map[string]string{
		"README.md":   "# widget\n",
		"go.mod":      "module example.com/widget\n\ngo 1.21\n",
		"add.go":      "package widget\n\nfunc Add(a, b int) int { return a + b }\n",
		"add_test.go": "package widget\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) {\n\tif Add(1, 1) != " + want + " {\n\t\tt.Fatal(\"one and one\")\n\t}\n}\n",
		".gitignore":  "coverage.out\n",
	}
	for rel, body := range files {
		if err := os.WriteFile(filepath.Join(root, rel), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r := execx.System{}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "widget"}} {
		if _, err := r.Run(root, "git", args...); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := execx.System{}.Run(root, "git", args...)
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(out)
}

func templateOrSkip(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(filepath.Join("..", "..", "template", "template.yaml")); err != nil {
		t.Skip("template not present")
	}
}

type importAnswer struct {
	Key    string       `json:"key"`
	Commit importCommit `json:"commit"`
}

// S-0098: an import that passes its tests is committed, and the commit holds
// the import and nothing else.
func TestImportCommitRunsTheTestsAndCommits(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := goRepo(t, true)
	before := git(t, root, "rev-parse", "HEAD")

	out, errOut, code := runIn(t, ".", "import", root, "--template", "../../template", "--yes", "--commit", "--json", "--trailer", "Co-Authored-By: flaiover <flaiover@localhost>")
	if code != 0 {
		t.Fatalf("import --commit: %d %s %s", code, out, errOut)
	}
	var got importAnswer
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	c := got.Commit
	if !c.Committed || c.Source != "repository" || len(c.Tests) != 1 || c.Tests[0].Name != "go test ./..." || !c.Tests[0].OK || got.Key == "" {
		t.Fatalf("answer: %+v", got)
	}
	if git(t, root, "rev-parse", "--short", "HEAD") != c.Commit || git(t, root, "rev-parse", "HEAD~1") != before {
		t.Errorf("one new commit on top of the repository's")
	}
	if st := git(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("left behind: %s", st)
	}
	inCommit := map[string]bool{}
	for _, f := range strings.Split(git(t, root, "show", "--name-only", "--format=", "HEAD"), "\n") {
		inCommit[f] = true
	}
	if !inCommit["system-flow.yaml"] || inCommit["add.go"] || inCommit["README.md"] || inCommit["go.mod"] {
		t.Errorf("the commit holds: %v", inCommit)
	}
	if msg := git(t, root, "log", "-1", "--format=%B"); !strings.Contains(msg, "flai import") || !strings.Contains(msg, "Co-Authored-By: flaiover") {
		t.Errorf("message: %s", msg)
	}
}

// A failing test leaves the import uncommitted, says which failed, and exits 5.
func TestImportCommitLeavesAFailedImportUncommitted(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := goRepo(t, false)
	before := git(t, root, "rev-parse", "HEAD")

	out, _, code := runIn(t, ".", "import", root, "--template", "../../template", "--yes", "--commit", "--json")
	if code != exitImportNotCommitted {
		t.Fatalf("exit %d: %s", code, out)
	}
	var got importAnswer
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	c := got.Commit
	if c.Committed || len(c.Tests) != 1 || c.Tests[0].OK || !strings.Contains(c.Tests[0].Output, "one and one") || !strings.Contains(c.Reason, "go test") {
		t.Fatalf("answer: %+v", c)
	}
	if git(t, root, "rev-parse", "HEAD") != before {
		t.Error("committed although a test failed")
	}
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Error("the imported files are gone; they should be left for the operator")
	}
	// the text answer says the same
	root2 := goRepo(t, false)
	text, _, code := runIn(t, ".", "import", root2, "--template", "../../template", "--yes", "--commit")
	if code != exitImportNotCommitted || !strings.Contains(text, "go test ./...: FAILED") || !strings.Contains(text, "not committed: tests failed") {
		t.Errorf("text: %d %s", code, text)
	}
}

// Uncommitted changes, or no git at all, refuse --commit before anything is
// touched: the commit must hold the import alone.
func TestImportCommitRefusesADirtyTreeUntouched(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := goRepo(t, true)
	if err := os.WriteFile(filepath.Join(root, "wip.txt"), []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, errOut, code := runIn(t, ".", "import", root, "--template", "../../template", "--yes", "--commit")
	if code == 0 || !strings.Contains(errOut, "uncommitted changes") {
		t.Errorf("dirty: %d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err == nil {
		t.Error("imported although refused")
	}

	plain := t.TempDir()
	_ = os.WriteFile(filepath.Join(plain, "README.md"), []byte("# plain\n"), 0o644)
	if _, errOut, code := runIn(t, ".", "import", plain, "--template", "../../template", "--yes", "--commit"); code == 0 || !strings.Contains(errOut, "not a git repository") {
		t.Errorf("no git: %d %s", code, errOut)
	}
}

// Nothing to test is not a failure: the import is committed and says so.
func TestImportCommitWithNoTestsFound(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := goRepo(t, true)
	for _, f := range []string{"go.mod", "add.go", "add_test.go"} {
		git(t, root, "rm", "-q", f)
	}
	git(t, root, "commit", "-q", "-m", "just docs")
	out, _, code := runIn(t, ".", "import", root, "--template", "../../template", "--yes", "--commit")
	if code != 0 || !strings.Contains(out, "tests: none found") || !strings.Contains(out, "committed as") {
		t.Errorf("no tests: %d %s", code, out)
	}
}

// importFormulaFixture is the import formula as S-0424 has the template ship
// it: the input path, then flai check --strict and flai import test.
const importFormulaFixture = `---
title: Import a repository
updated: 2026-11-02
status: active
name: import
description: "Check what flai import brought and run the repository's tests in a worktree."
nature: improvement
inputs:
  - name: path
    description: The imported repository's root.
steps:
  - id: check
    title: Check {{path}} against the standard
    run: {flai: check --strict}
  - id: test
    title: Run the repository's tests
    run: {flai: import test}
    after: [check]
---

# Import

The steps the import runs after its commit.
`

// importTemplate is a copy of the template a test may change, shipping the
// import formula or not as asked.
func importTemplate(t *testing.T, withFormula bool) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS(filepath.Join("..", "..", "template"))); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(dir, "root", "design", "workflows", "formulas", "import.md")
	for _, p := range []string{f, f + ".tmpl"} {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
	if withFormula {
		if err := os.WriteFile(f, []byte(importFormulaFixture), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

type importFormulaAnswer struct {
	Key     string           `json:"key"`
	Commit  importCommit     `json:"commit"`
	Formula formulaRunResult `json:"formula"`
}

// S-0424: --formula commits the import before any test, so a repository
// whose tests fail is committed all the same, then makes the import
// formula's instance ready and commits it on the project's branch.
func TestImportFormulaCommitsBeforeAnyTestAndMakesTheInstanceReady(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	tpl := importTemplate(t, true)
	root := goRepo(t, false)
	before := git(t, root, "rev-parse", "HEAD")

	out, errOut, code := runIn(t, ".", "import", root, "--template", tpl, "--yes", "--formula", "--json", "--trailer", "Co-Authored-By: flaiover <flaiover@localhost>")
	if code != 0 {
		t.Fatalf("import --formula: %d %s %s", code, out, errOut)
	}
	var got importFormulaAnswer
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	c, f := got.Commit, got.Formula
	if !c.Committed || c.Source != importTestsFormula || len(c.Tests) != 0 || got.Key == "" {
		t.Fatalf("the import's commit: %+v", got)
	}
	if !strings.HasPrefix(f.Story, "S-") || f.Status != "ready" || f.Formula != "import" || f.Inputs["path"] != root || len(f.Tasks) != 2 || f.Tasks[0].Step != "check" || f.Tasks[1].Step != "test" || f.Commit == "" {
		t.Fatalf("the instance: %+v", f)
	}
	// two commits on the repository's: the import, then the instance
	if git(t, root, "rev-parse", "--short", "HEAD~1") != c.Commit || git(t, root, "rev-parse", "HEAD~2") != before || git(t, root, "rev-parse", "--short", "HEAD") != f.Commit {
		t.Errorf("the import's commit, then the instance's, on top of the repository's: %s", git(t, root, "log", "--oneline", "-4"))
	}
	if br := git(t, root, "branch", "--show-current"); br != "main" {
		t.Errorf("committed on %s, not main", br)
	}
	// flai's own cache is left out: the repository's .gitignore, kept over
	// the template's, does not name it
	if st := git(t, root, "status", "--porcelain"); st != "" && st != "?? .flai-cache/" {
		t.Errorf("left behind: %s", st)
	}
	imported := git(t, root, "show", "--name-only", "--format=", "HEAD~1")
	if !strings.Contains(imported, "system-flow.yaml") || !strings.Contains(imported, "design/workflows/formulas/import.md") || strings.Contains(imported, "add.go") {
		t.Errorf("the import's commit holds: %s", imported)
	}
	if !strings.Contains(git(t, root, "show", "--name-only", "--format=", "HEAD"), f.Story) {
		t.Errorf("the instance's commit does not hold %s", f.Story)
	}
	for _, rev := range []string{"HEAD~1", "HEAD"} {
		if msg := git(t, root, "log", "-1", "--format=%B", rev); !strings.Contains(msg, "Co-Authored-By: flaiover") {
			t.Errorf("%s has no trailer: %s", rev, msg)
		}
	}
	// the text answer names the instance
	root2 := goRepo(t, true)
	text, errOut, code := runIn(t, ".", "import", root2, "--template", tpl, "--yes", "--formula")
	if code != 0 || !strings.Contains(text, "not run before the commit") || !strings.Contains(text, "committed as") || !strings.Contains(text, "import formula: S-") || !strings.Contains(text, " ready, tasks ") {
		t.Errorf("text: %d %s %s", code, text, errOut)
	}
}

// I-0302: flai import writes system-flow.lock.yaml as flai new does and the
// import's commit holds it, so the import formula's first step, flai check
// --strict, finds no layout.lock; and --dry-run writes no lock.
func TestImportFormulaCommitsTheLockSoCheckFindsNoLayoutLock(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	tpl := importTemplate(t, true)

	dry := goRepo(t, true)
	if out, errOut, code := runIn(t, ".", "import", dry, "--template", tpl, "--yes", "--formula", "--dry-run"); code != 0 {
		t.Fatalf("import --dry-run: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(dry, "system-flow.lock.yaml")); !os.IsNotExist(err) {
		t.Errorf("--dry-run wrote a lock: %v", err)
	}

	root := goRepo(t, true)
	out, errOut, code := runIn(t, ".", "import", root, "--template", tpl, "--yes", "--formula", "--json")
	if code != 0 {
		t.Fatalf("import --formula: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "system-flow.lock.yaml")); err != nil {
		t.Fatalf("no lock after the import: %v", err)
	}
	if imported := git(t, root, "show", "--name-only", "--format=", "HEAD~1"); !strings.Contains(imported, "system-flow.lock.yaml") {
		t.Errorf("the import's commit does not hold the lock: %s", imported)
	}
	if st := git(t, root, "status", "--porcelain"); st != "" && st != "?? .flai-cache/" {
		t.Errorf("left behind: %s", st)
	}
	checked, errOut, _ := runIn(t, root, "check", "--strict", "--json")
	if strings.Contains(checked, "layout.lock") || strings.Contains(errOut, "layout.lock") {
		t.Errorf("flai check --strict reports layout.lock after the import: %s %s", checked, errOut)
	}
}

// --formula refuses, before anything changes, --commit beside it, a
// template that ships no import formula, and a tree --commit would refuse.
func TestImportFormulaRefusesBeforeAnythingChanges(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	untouched := func(root, why string) {
		t.Helper()
		if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err == nil {
			t.Errorf("%s: imported although refused", why)
		}
		if st := git(t, root, "status", "--porcelain"); st != "" && st != "?? wip.txt" {
			t.Errorf("%s: changed %s", why, st)
		}
	}
	root := goRepo(t, true)
	tpl := importTemplate(t, true)
	if _, errOut, code := runIn(t, ".", "import", root, "--template", tpl, "--yes", "--formula", "--commit"); code == 0 || !strings.Contains(errOut, "--formula and --commit") {
		t.Errorf("with --commit: %d %s", code, errOut)
	}
	untouched(root, "with --commit")

	bare := importTemplate(t, false)
	if _, errOut, code := runIn(t, ".", "import", root, "--template", bare, "--yes", "--formula"); code == 0 || !strings.Contains(errOut, "ships no import formula") || !strings.Contains(errOut, "version 1.") {
		t.Errorf("no import formula: %d %s", code, errOut)
	}
	untouched(root, "no import formula")

	if err := os.WriteFile(filepath.Join(root, "wip.txt"), []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, ".", "import", root, "--template", tpl, "--yes", "--formula"); code == 0 || !strings.Contains(errOut, "uncommitted changes") {
		t.Errorf("dirty: %d %s", code, errOut)
	}
	untouched(root, "dirty")
}

// flai import test runs the tests --commit runs, the template's Makefile
// not counting as the repository's own, prints each outcome, and exits 5
// when one fails.
func TestImportTestRunsTheRepositorysTests(t *testing.T) {
	t.Parallel()
	templateOrSkip(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	tpl := importTemplate(t, false)
	imported := func(passing bool) string {
		t.Helper()
		root := goRepo(t, passing)
		if _, errOut, code := runIn(t, ".", "import", root, "--template", tpl, "--yes"); code != 0 {
			t.Fatalf("import: %d %s", code, errOut)
		}
		return root
	}

	passing := imported(true)
	out, errOut, code := runIn(t, passing, "import", "test")
	if code != 0 || !strings.Contains(out, "go test ./...: passed") || !strings.Contains(out, "1 test command passed") {
		t.Errorf("passing: %d %s %s", code, out, errOut)
	}
	out, _, code = runIn(t, ".", "import", "test", passing, "--json")
	var run importTestRun
	if err := json.Unmarshal([]byte(out), &run); err != nil || code != 0 || !run.OK || run.Source != importTestsRepository || len(run.Tests) != 1 || run.Tests[0].Name != "go test ./..." || run.Root != passing {
		t.Errorf("passing --json: %d %+v %v", code, run, err)
	}

	failing := imported(false)
	out, errOut, code = runIn(t, failing, "import", "test")
	if code != exitImportTestsFailed || !strings.Contains(out, "go test ./...: FAILED") || !strings.Contains(out, "one and one") || !strings.Contains(errOut, "tests failed: go test ./...") {
		t.Errorf("failing: %d %s %s", code, out, errOut)
	}
}
