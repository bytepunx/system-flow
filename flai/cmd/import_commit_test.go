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
	templateOrSkip(t)
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
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
	templateOrSkip(t)
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
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
	templateOrSkip(t)
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
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
	templateOrSkip(t)
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
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
