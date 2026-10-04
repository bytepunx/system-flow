package storygit

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// writeFile writes a file and the folders above it.
func writeFile(t *testing.T, dir, rel string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(rel)), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, dir, rel, rel+"\n")
}

// I-0065: a file another story wrote, committed on its branch or not yet
// committed in its worktree, counts as much as one on main.
func TestFolderNamesReadsMainEveryStoryBranchAndEveryWorktree(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "main", "design/issues/I-0001-on-main.md")
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	commitFiles(t, dir, "one", "design/issues/I-0002-one.md", "design/issues/sub/nested.md")
	git(t, dir, "checkout", "-q", "-b", "story/S-0002", "main")
	commitFiles(t, dir, "two", "design/issues/I-0003-two.md")
	git(t, dir, "checkout", "-q", "-b", "other/x", "main")
	commitFiles(t, dir, "not a story", "design/issues/I-0009-other.md")
	git(t, dir, "checkout", "-q", "main")
	wt := filepath.Join(t.TempDir(), "S-0003")
	git(t, dir, "worktree", "add", "-q", "-b", "story/S-0003", wt)
	writeFile(t, wt, "design/issues/I-0004-uncommitted.md")
	writeFile(t, dir, "design/issues/I-0005-main-uncommitted.md")

	got := FolderNames(r, &workitem.Repo{Root: wt, MainRoot: dir}, "design/issues")
	want := []string{"I-0001-on-main.md", "I-0002-one.md", "I-0003-two.md", "I-0004-uncommitted.md", "I-0005-main-uncommitted.md"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("names = %v, want %v", got, want)
	}
}

func TestFolderNamesWithoutStoryBranchesOrTheFolder(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	writeFile(t, dir, "design/issues/I-0001-uncommitted.md")
	repo := &workitem.Repo{Root: dir, MainRoot: dir}
	if got := FolderNames(r, repo, "design/issues/"); !reflect.DeepEqual(got, []string{"I-0001-uncommitted.md"}) {
		t.Errorf("no story branches: %v", got)
	}
	if got := FolderNames(r, repo, "design/adrs"); len(got) != 0 {
		t.Errorf("missing folder: %v", got)
	}
}

// A project in a folder of a larger repository reads its folder there, in
// every worktree and on every branch.
func TestFolderNamesKeepsToTheProjectsFolder(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "main", "proj/design/issues/I-0001.md", "design/issues/I-0008-outside.md")
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	commitFiles(t, dir, "one", "proj/design/issues/I-0002.md")
	git(t, dir, "checkout", "-q", "main")
	wt := filepath.Join(t.TempDir(), "S-0002")
	git(t, dir, "worktree", "add", "-q", "-b", "story/S-0002", wt)
	writeFile(t, wt, "proj/design/issues/I-0003.md")

	sub := filepath.Join(dir, "proj")
	got := FolderNames(r, &workitem.Repo{Root: sub, MainRoot: sub}, "design/issues")
	if want := []string{"I-0001.md", "I-0002.md", "I-0003.md"}; !reflect.DeepEqual(got, want) {
		t.Errorf("names = %v, want %v", got, want)
	}
}
