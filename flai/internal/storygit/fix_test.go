package storygit

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// S-0367: after the clean rebase the fix runs in the worktree over the files
// the branch changes, deleted ones and those only main changed left out, and
// what it rewrote is committed on the branch, leaving the worktree clean.
func TestSyncCommitsTheLintFixOfTheFilesTheBranchChanges(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] mine", "docs/mine.md", "docs/other.md")
	git(t, wt, "rm", "-q", "docs/guide.md")
	git(t, wt, "commit", "-q", "-m", "feat: [S-0001] drop the guide")
	commitFiles(t, repo.MainRoot, "docs: main's", "docs/main.md")
	o := syncOptions(repo, story)
	var gotDir string
	var gotPaths []string
	o.Fix = func(dir string, paths []string) ([]string, error) {
		gotDir, gotPaths = dir, slices.Clone(paths)
		return []string{"docs/mine.md"}, os.WriteFile(filepath.Join(dir, "docs", "mine.md"), []byte("fixed\n"), 0o644)
	}

	res, err := Sync(o)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Synced || strings.Join(res.Fixed, ",") != "docs/mine.md" {
		t.Fatalf("result: %+v", res)
	}
	if gotDir != wt || strings.Join(gotPaths, ",") != "docs/mine.md,docs/other.md" {
		t.Errorf("fix ran in %s over %v, want %s over the branch's docs/mine.md and docs/other.md", gotDir, gotPaths, wt)
	}
	r := execx.System{}
	if got, _ := r.Run(wt, "git", "log", "-1", "--format=%s"); got != "chore: [S-0001] fix the lint of the files the branch changes, after the sync" {
		t.Errorf("the fix's commit: %q", got)
	}
	if got, _ := r.Run(wt, "git", "show", "--name-status", "--format=", "HEAD"); got != "M\tdocs/mine.md" {
		t.Errorf("the fix's commit changed:\n%s", got)
	}
	if got, _ := r.Run(wt, "git", "show", "HEAD:docs/mine.md"); got != "fixed" {
		t.Errorf("docs/mine.md committed as %q", got)
	}
	if st, _ := Uncommitted(r, wt); len(st) != 0 {
		t.Errorf("the worktree is not clean: %v", st)
	}
	data, _ := json.Marshal(res)
	if !strings.Contains(string(data), `"fixed":["docs/mine.md"]`) {
		t.Errorf("json: %s", data)
	}
}

// With no fix, a fix that fails, or one that rewrites nothing git sees, the
// sync stands with nothing committed and nothing fixed.
func TestSyncStandsWithNoFixAFailingFixOrNothingFixed(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] guide", "docs/guide.md")
	r := execx.System{}
	head, _ := r.Run(wt, "git", "rev-parse", "HEAD")
	for name, f := range map[string]FixFunc{
		"no fix":  nil,
		"failing": func(string, []string) ([]string, error) { return nil, errors.New("the tiers are unreadable") },
		"nothing": func(string, []string) ([]string, error) { return []string{}, nil },
		"unchanged": func(string, []string) ([]string, error) {
			return []string{"docs/guide.md"}, nil
		},
	} {
		o := syncOptions(repo, story)
		o.Fix = f
		res, err := Sync(o)
		if err != nil || !res.Synced || res.Fixed == nil || len(res.Fixed) != 0 {
			t.Fatalf("%s: result %+v, err %v", name, res, err)
		}
		if now, _ := r.Run(wt, "git", "rev-parse", "HEAD"); now != head {
			t.Errorf("%s: a commit was made: %s, was %s", name, now, head)
		}
		if st, _ := Uncommitted(r, wt); len(st) != 0 {
			t.Errorf("%s: the worktree is not clean: %v", name, st)
		}
	}
}

// A fix written but not committed is an error that says what to commit.
func TestSyncFixNotCommittedSaysWhatToCommit(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] guide", "docs/guide.md")
	hooks := filepath.Join(repo.MainRoot, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, hooks, "pre-commit", "#!/bin/sh\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	o := syncOptions(repo, story)
	o.Fix = func(dir string, _ []string) ([]string, error) {
		return []string{"docs/guide.md"}, os.WriteFile(filepath.Join(dir, "docs", "guide.md"), []byte("fixed\n"), 0o644)
	}

	res, err := Sync(o)
	if err == nil {
		t.Fatalf("no error: %+v", res)
	}
	for _, want := range []string{"commit the lint fix of docs/guide.md on story/S-0001", "git add -A -- docs/guide.md",
		`"chore: [S-0001] fix the lint of the files the branch changes, after the sync"`, "flai stream sync S-0001 again"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error lacks %q: %v", want, err)
		}
	}
	if !res.Synced || len(res.Fixed) != 0 {
		t.Errorf("result: %+v", res)
	}
}
