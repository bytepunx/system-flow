package storygit

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0210: each commit on the main branch is one entry, newest first, its
// files outside wip sorted; merges, story branches, flai's bookkeeping
// commits, and commits that changed only wip are left out.
func TestCommitFilesReadsEachMainCommitsFilesOutsideWip(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	repo := &workitem.Repo{Root: dir, MainRoot: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}

	commitFiles(t, dir, "feat: [S-0001] one", "cli/b.go", "cli/a.go", "wip/kanban/stories/S-0001.md")
	commitFiles(t, dir, "chore: [S-0001] create story: one", "wip/kanban/stories/S-0002.md", "cli/c.go")
	commitFiles(t, dir, "chore: publish flai 1.0.0", "CHANGELOG.md", "cli/CHANGELOG.md")
	commitFiles(t, dir, "docs: wip only", "wip/agents/S-0001.md")
	git(t, dir, "checkout", "-q", "-b", "story/S-0003")
	commitFiles(t, dir, "feat: [S-0003] on its branch", "web/d.go")
	git(t, dir, "checkout", "-q", "-b", "story/S-0004", "main")
	commitFiles(t, dir, "feat: [S-0004] merged", "web/e.go")
	git(t, dir, "checkout", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "merge: [S-0004] the branch", "story/S-0004")
	commitFiles(t, dir, "fix: two", "cli/a.go", "docs/x.md")

	got, err := CommitFiles(r, repo)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"cli/a.go", "docs/x.md"}, {"web/e.go"}, {"cli/a.go", "cli/b.go"}, {"a.md"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("commit files = %v, want %v", got, want)
	}
}

// A project in a folder of a larger repository reads only its own files,
// named from its root.
func TestCommitFilesKeepsToTheProjectsFolder(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "feat: both", "proj/cli/a.go", "proj/wip/x.md", "outside.go")
	sub := filepath.Join(dir, "proj")
	repo := &workitem.Repo{Root: sub, MainRoot: sub, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip/"}}}
	got, err := CommitFiles(r, repo)
	if err != nil {
		t.Fatal(err)
	}
	if want := [][]string{{"cli/a.go"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("commit files = %v, want %v", got, want)
	}
}

func TestCommitFilesFailsWithoutARepository(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	dir := t.TempDir()
	repo := &workitem.Repo{Root: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}
	got, err := CommitFiles(execx.System{}, repo)
	if err == nil || got != nil || !strings.Contains(err.Error(), "cannot read the project's commits in "+dir) {
		t.Errorf("not a repository: %v %v", got, err)
	}
}

func TestBookkeeping(t *testing.T) {
	for subject, want := range map[string]bool{
		"chore: [S-0012] accept and archive":                        true,
		"chore: [S-0012] accept and archive, with E-0003":           true,
		"chore: [S-0012] accept and archive; release flai 1.2.0":    true,
		"chore: [T-0815] create task: a title":                      true,
		"chore: [S-0012] edit touches, criteria":                    true,
		"chore: publish flai 1.30.0, flaiover 0.31.0":               true,
		"chore: set the project's default agent":                    true,
		"chore: clear the project's default agent":                  true,
		"chore: bring my-repo under system-flow (flai import)":      true,
		"feat: [S-0012] accept and archive":                         false,
		"chore: [S-0012] record another occurrence of I-0027":       false,
		"chore: update issues":                                      false,
		"chore: publishing notes":                                   false,
		"docs: [S-0012] edit the guide":                             false,
		"chore: set the project's default agent and more":           false,
		"fix: [S-0012] flai import brings a repo under system-flow": false,
		"": false,
	} {
		if got := bookkeeping(subject); got != want {
			t.Errorf("%q: %v, want %v", subject, got, want)
		}
	}
}
