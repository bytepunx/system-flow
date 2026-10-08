package verify

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

func TestResolveExpandsFoldersAndKeepsFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"flai/main.go":                           {},
		"flai/internal/verify/verify.go":         {},
		"flai/internal/verify/testdata/x.json":   {},
		"flaiover/src/a.ts":                      {},
		"flaiover/node_modules/dep/index.js":     {},
		"flaiover/.git/HEAD":                     {},
		".flai-cache/worktrees/S-0001/README.md": {},
		"README.md":                              {},
	}
	cases := []struct {
		args []string
		want []string
	}{
		{[]string{"flai/main.go"}, []string{"flai/main.go"}},
		{[]string{"./flai/internal/"}, []string{"flai/internal/verify/testdata/x.json", "flai/internal/verify/verify.go"}},
		{[]string{"flaiover", "flaiover/src/a.ts"}, []string{"flaiover/src/a.ts"}},
		{[]string{"."}, []string{
			"README.md", "flai/internal/verify/testdata/x.json", "flai/internal/verify/verify.go",
			"flai/main.go", "flaiover/src/a.ts",
		}},
	}
	for _, c := range cases {
		got, err := Resolve(fsys, c.args)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("Resolve(%q) = %q, %v; want %q", c.args, got, err, c.want)
		}
	}
	for _, bad := range []string{"missing.go", "../elsewhere", "/etc/passwd"} {
		if _, err := Resolve(fsys, []string{bad}); err == nil || !strings.Contains(err.Error(), bad) {
			t.Errorf("Resolve(%q) error %v, want one naming it", bad, err)
		}
	}
}

// fakeGit answers git commands by their arguments.
type fakeGit struct {
	out map[string]string
	err map[string]error
	ran []string
}

func (f *fakeGit) Run(dir, name string, args ...string) (string, error) {
	return f.RunInput(dir, name, "", args...)
}

func (f *fakeGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	key := strings.Join(args, " ")
	f.ran = append(f.ran, dir+": "+name+" "+key)
	return f.out[key], f.err[key]
}

func (f *fakeGit) LookPath(name string) (string, error) { return name, nil }

func TestChangedPathsJoinsCommittedAndUncommittedChanges(t *testing.T) {
	git := &fakeGit{out: map[string]string{
		"diff --name-only --diff-filter=d main...HEAD": "flai/a.go\nflai/b.go\ndocs/gone-later.md",
		// the runner trims output, so the first line has lost its leading space
		"status --porcelain --untracked-files=all": strings.Join([]string{
			"M flai/a.go",
			"M  flai/c.go",
			"R  flai/old.go -> flai/new.go",
			" D docs/gone-later.md",
			"D  flai/staged-gone.go",
			"?? flaiover/src/new.ts",
			`?? "design/caf\303\251.md"`,
			"A  x",
		}, "\n"),
	}}
	got, err := ChangedPaths(git, "/work/S-0001", "main")
	want := []string{"design/café.md", "flai/a.go", "flai/b.go", "flai/c.go", "flai/new.go", "flaiover/src/new.ts", "x"}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("ChangedPaths = %q, %v; want %q", got, err, want)
	}
	if len(git.ran) != 2 || !strings.HasPrefix(git.ran[0], "/work/S-0001: git diff") {
		t.Errorf("ran %q", git.ran)
	}
}

func TestChangedPathsWithoutABaseCountsOnlyTheUncommitted(t *testing.T) {
	git := &fakeGit{out: map[string]string{"status --porcelain --untracked-files=all": "?? a.go"}}
	got, err := ChangedPaths(git, "/work", "")
	if err != nil || !reflect.DeepEqual(got, []string{"a.go"}) || len(git.ran) != 1 {
		t.Errorf("ChangedPaths = %q, %v after %q", got, err, git.ran)
	}
}

func TestChangedPathsSaysWhichBaseFailed(t *testing.T) {
	git := &fakeGit{err: map[string]error{"diff --name-only --diff-filter=d trunk...HEAD": errors.New("unknown revision")}}
	_, err := ChangedPaths(git, "/work", "trunk")
	if err == nil || !strings.Contains(err.Error(), "trunk") || !strings.Contains(err.Error(), "unknown revision") {
		t.Errorf("error %v", err)
	}
}

func TestBaseChangesListsThePathsAndTheCommitsTheCheckoutLacks(t *testing.T) {
	git := &fakeGit{out: map[string]string{
		"diff --no-renames --name-only HEAD...main": "wip/b.md\n\"wip/caf\\303\\251.md\"\nwip/a.md",
		"rev-list --abbrev-commit HEAD..main":       "c3c3c3c\nb2b2b2b\na1a1a1a\n",
	}}
	paths, commits, err := BaseChanges(git, "/work/S-0001", "main")
	if err != nil || !reflect.DeepEqual(paths, []string{"wip/a.md", "wip/b.md", "wip/café.md"}) ||
		!reflect.DeepEqual(commits, []string{"c3c3c3c", "b2b2b2b", "a1a1a1a"}) {
		t.Errorf("BaseChanges = %q, %q, %v", paths, commits, err)
	}
	git = &fakeGit{err: map[string]error{"diff --no-renames --name-only HEAD...trunk": errors.New("unknown revision")}}
	if _, _, err := BaseChanges(git, "/work", "trunk"); err == nil || !strings.Contains(err.Error(), "trunk") || !strings.Contains(err.Error(), "unknown revision") {
		t.Errorf("error %v", err)
	}
}

// Integration: these run real git.
func realRepo(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	runGit(t, dir, "init", "-q", "-b", "main")
	runGit(t, dir, "config", "user.email", "t@t")
	runGit(t, dir, "config", "user.name", "t")
	return dir
}

// runGit runs git in dir and answers what it printed, failing the test when
// git fails.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := (execx.System{}).Run(dir, "git", args...)
	if err != nil {
		t.Fatal(out, err)
	}
	return strings.TrimSpace(out)
}

// writeFile writes a file below dir, making its folder.
func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBaseChangesNamesARenamesBothPathsAndCountsAMergeOnBase(t *testing.T) {
	dir := realRepo(t)
	writeFile(t, dir, "wip/old.md", "a\nb\nc\n")
	writeFile(t, dir, "flai/x.go", "package x\n")
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-q", "-m", "init")
	runGit(t, dir, "branch", "story")
	runGit(t, dir, "mv", "wip/old.md", "wip/new.md")
	runGit(t, dir, "commit", "-q", "-m", "rename")
	runGit(t, dir, "checkout", "-q", "-b", "side")
	writeFile(t, dir, "wip/side.md", "side\n")
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-q", "-m", "side")
	runGit(t, dir, "checkout", "-q", "main")
	runGit(t, dir, "merge", "-q", "--no-ff", "-m", "merge side", "side")
	runGit(t, dir, "checkout", "-q", "story")
	writeFile(t, dir, "flai/x.go", "package x // story\n")
	runGit(t, dir, "commit", "-q", "-am", "story")

	paths, commits, err := BaseChanges(execx.System{}, dir, "main")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"wip/new.md", "wip/old.md", "wip/side.md"}; !reflect.DeepEqual(paths, want) {
		t.Errorf("paths %q, want %q: a rename as both paths, and not the branch's own", paths, want)
	}
	want := strings.Fields(runGit(t, dir, "rev-list", "--abbrev-commit", "story..main"))
	if len(commits) != 3 || !reflect.DeepEqual(commits, want) || commits[0] != runGit(t, dir, "rev-parse", "--short", "main") {
		t.Errorf("commits %q, want the rename, the side commit, and the merge, newest first: %q", commits, want)
	}
}
