package cmd

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// storyCommitProject is adrProject with story S-0001 committed on main and
// its branch checked out in its worktree, holding an untracked file and a
// changed one that --commit must leave alone; no story is named by the
// environment, so the branch names it. It returns the main checkout and the
// worktree. It runs real git.
func storyCommitProject(t *testing.T) (root, wt string) {
	t.Helper()
	root = adrProject(t)
	t.Setenv("FLAI_STORY", "")
	t.Setenv("FLAI_AGENT", "")
	writeIn(t, root, "design/conventions/README.md", "# c\n")
	if out, errOut, code := runIn(t, root, "story", "new", "Commit from the worktree"); code != 0 || !strings.HasPrefix(out, "S-0001 ") {
		t.Fatalf("story new: %d %s %s", code, out, errOut)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "story")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	wt = repo.WorktreePath("S-0001")
	gitIn(t, root, "worktree", "add", "-q", "-b", storygit.Branch("S-0001"), wt, "main")
	writeIn(t, wt, "scratch.md", "not the command's\n")
	writeIn(t, wt, "design/adrs/0001-first.md", adrFile("0001", "First", "accepted")+"\nAn edit of the story's own.\n")
	return root, wt
}

// storyCommits is the subjects of the commits on S-0001's branch that main
// does not have, newest first.
func storyCommits(t *testing.T, wt string) []string {
	t.Helper()
	out := gitIn(t, wt, "log", "--format=%s", "main..HEAD")
	if out == "" {
		return nil
	}
	return strings.Split(out, "\n")
}

// committedFiles is the files the commit at HEAD in dir changed.
func committedFiles(t *testing.T, dir string) []string {
	t.Helper()
	return strings.Split(gitIn(t, dir, "show", "--name-only", "--format=", "HEAD"), "\n")
}

// storyTouches is S-0001's touches, read from the main checkout.
func storyTouches(t *testing.T, root string) []string {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	return it.Touches
}

// leftAlone fails unless the worktree still holds the untracked file and the
// changed one storyCommitProject left, uncommitted.
func leftAlone(t *testing.T, wt string) {
	t.Helper()
	st := gitIn(t, wt, "status", "--porcelain")
	// gitIn trims the output, so the first line's leading space may be gone
	for _, want := range []string{"?? scratch.md", "M design/adrs/0001-first.md"} {
		if !strings.Contains(st, want) {
			t.Errorf("--commit committed what it did not write; status lacks %q:\n%s", want, st)
		}
	}
}

func TestBeginStoryCommitRefusesWithoutAStoryBranch(t *testing.T) {
	root, wt := storyCommitProject(t)
	a := &app{runner: execx.System{}}
	repo, err := workitem.Open(wt)
	if err != nil {
		t.Fatal(err)
	}
	if sc, err := a.beginStoryCommit(repo, "", false, nil); err != nil || sc.story != "S-0001" || sc.repo.Root != wt {
		t.Errorf("the worktree's branch names the story: %+v %v", sc, err)
	}
	if _, err := a.beginStoryCommit(repo, "", true, nil); err == nil || !strings.Contains(err.Error(), "--commit and --autocommit") {
		t.Errorf("--autocommit with --commit is refused: %v", err)
	}
	if _, err := a.beginStoryCommit(repo, "S-0002", false, nil); err == nil || !strings.Contains(err.Error(), "not story/S-0002") || !strings.Contains(err.Error(), "nothing was written") {
		t.Errorf("another story's commit is refused in this worktree: %v", err)
	}
	main, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := a.beginStoryCommit(main, "", false, nil); err == nil || !strings.Contains(err.Error(), "no story resolves") {
		t.Errorf("the main checkout names no story: %v", err)
	}
	if _, err := a.beginStoryCommit(main, "S-0001", false, nil); err == nil || !strings.Contains(err.Error(), "has main checked out, not story/S-0001") {
		t.Errorf("the main checkout is not the story's branch: %v", err)
	}
}
