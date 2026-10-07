package verify

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
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
