package storygit

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// commitFiles writes each file and commits them with the subject.
func commitFiles(t *testing.T, dir, subject string, files ...string) {
	t.Helper()
	for _, f := range files {
		if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(f)), 0o755); err != nil {
			t.Fatal(err)
		}
		write(t, dir, f, subject+"\n")
	}
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", subject)
}

// S-0205: a story's commits are those on the main branch and the story
// branches whose subject names it in brackets, merges left out, and their
// files under wip left out.
func TestCommittedReadsEachStorysFilesFromTheMainAndStoryBranches(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	repo := &workitem.Repo{Root: dir, MainRoot: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}

	commitFiles(t, dir, "feat: [S-0001] one", "cli/a.go", "wip/kanban/stories/S-0001.md")
	commitFiles(t, dir, "fix: [S-1] and [S-0002] both", "cli/b.go")
	commitFiles(t, dir, "docs: no story", "README.md")
	commitFiles(t, dir, "chore: [T-0003] a task, [S-0002 unclosed", "cli/t.go")
	git(t, dir, "mv", "cli/a.go", "cli/z.go")
	git(t, dir, "commit", "-q", "-m", "refactor: [S-0007] rename")

	// a story branch not merged, and one merged with a merge commit that
	// names a story of its own
	git(t, dir, "checkout", "-q", "-b", "story/S-0003")
	commitFiles(t, dir, "feat: [S-0003] on its branch", "web/d.go")
	git(t, dir, "checkout", "-q", "-b", "story/S-0004", "main")
	commitFiles(t, dir, "feat: [S-0004] merged", "web/e.go")
	git(t, dir, "checkout", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "merge: [S-0005] the branch", "story/S-0004")
	git(t, dir, "branch", "-q", "-D", "story/S-0004")
	// a branch that is neither
	git(t, dir, "checkout", "-q", "-b", "other/x")
	commitFiles(t, dir, "feat: [S-0006] elsewhere", "web/f.go")
	git(t, dir, "checkout", "-q", "main")

	got, err := Committed(r, repo)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{
		"S-0001": {"cli/a.go", "cli/b.go"},
		"S-0002": {"cli/b.go"},
		"S-0003": {"web/d.go"},
		"S-0004": {"web/e.go"},
		"S-0007": {"cli/a.go", "cli/z.go"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("committed = %v, want %v", got, want)
	}
}

// A project in a folder of a larger repository reads only its own files,
// named from its root.
func TestCommittedKeepsToTheProjectsFolder(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "feat: [S-0001] both", "proj/cli/a.go", "proj/wip/x.md", "outside.go")
	sub := filepath.Join(dir, "proj")
	repo := &workitem.Repo{Root: sub, MainRoot: sub, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip/"}}}
	got, err := Committed(r, repo)
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string][]string{"S-0001": {"cli/a.go"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("committed = %v, want %v", got, want)
	}
}

func TestCommittedFailsWithoutARepository(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	dir := t.TempDir()
	repo := &workitem.Repo{Root: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}
	got, err := Committed(execx.System{}, repo)
	if err == nil || got != nil || !strings.Contains(err.Error(), "cannot read the stories' commits in "+dir) {
		t.Errorf("not a repository: %v %v", got, err)
	}
}

func TestSubjectStories(t *testing.T) {
	for subject, want := range map[string][]string{
		"feat: [S-12] x":                {"S-0012"},
		"[S-0001][S-002] [S-1]":         {"S-0001", "S-0002"},
		"[E-0001] [T-0002] [s-3] [S-x]": nil,
		"S-0004 without brackets [":     nil,
	} {
		if got := subjectStories(subject); !reflect.DeepEqual(got, want) {
			t.Errorf("%q: %v, want %v", subject, got, want)
		}
	}
}
