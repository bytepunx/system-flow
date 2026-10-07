package storygit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// storyCheckout is a repository with a.md, b.md, and c.md committed on main
// and story/S-0007 checked out. Integration: it runs real git.
func storyCheckout(t *testing.T) string {
	t.Helper()
	dir := gitRepo(t)
	write(t, dir, "b.md", "b\n")
	write(t, dir, "c.md", "c\n")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "b and c")
	git(t, dir, "checkout", "-q", "-b", Branch("S-0007"))
	return dir
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := (execx.System{}).Run(dir, "git", args...)
	if err != nil {
		t.Fatal(out, err)
	}
	return out
}

// The paths a command wrote are committed on the story's branch with the
// story's prefix and the trailers; what else the worktree holds, staged or
// not, stays as it was.
func TestCommitPathsCommitsOnlyThePathsWithTheStorysPrefix(t *testing.T) {
	dir := storyCheckout(t)
	write(t, dir, "a.md", "written by the command\n")
	if err := os.MkdirAll(filepath.Join(dir, "design/issues"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, dir, "design/issues/I-0001.md", "new\n")
	write(t, dir, "b.md", "staged by the agent\n")
	git(t, dir, "add", "b.md")
	write(t, dir, "c.md", "changed by the agent\n")

	got, err := CommitPaths(CommitOptions{
		Runner: execx.System{}, Dir: dir, Story: "S-7",
		Paths:   []string{"a.md", "design/issues/I-0001.md"},
		Subject: "record I-0001", Trailers: []string{"Co-Authored-By: Claude <noreply@anthropic.com>"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Subject != "docs: [S-0007] record I-0001" || strings.Join(got.Paths, ",") != "a.md,design/issues/I-0001.md" {
		t.Fatalf("commit: %+v", got)
	}
	if head := gitOut(t, dir, "rev-parse", "HEAD"); got.Hash != head {
		t.Errorf("hash %s, HEAD %s", got.Hash, head)
	}
	if msg := gitOut(t, dir, "log", "-1", "--format=%B"); msg != "docs: [S-0007] record I-0001\n\nCo-Authored-By: Claude <noreply@anthropic.com>" {
		t.Errorf("message %q", msg)
	}
	if files := gitOut(t, dir, "show", "--name-only", "--format=", "HEAD"); files != "a.md\ndesign/issues/I-0001.md" {
		t.Errorf("committed %q", files)
	}
	if st := gitOut(t, dir, "status", "--porcelain"); st != "M  b.md\n M c.md" {
		t.Errorf("the agent's changes did not stay as they were: %q", st)
	}
}

// Paths that did not change commit nothing, and say so with no commit.
func TestCommitPathsWithNothingChangedCommitsNothing(t *testing.T) {
	dir := storyCheckout(t)
	head := gitOut(t, dir, "rev-parse", "HEAD")
	write(t, dir, "c.md", "changed by the agent\n")
	got, err := CommitPaths(CommitOptions{Runner: execx.System{}, Dir: dir, Story: "S-0007", Type: "chore", Paths: []string{"a.md"}, Subject: "nothing"})
	if err != nil || got != nil {
		t.Fatalf("commit %+v, err %v", got, err)
	}
	if now := gitOut(t, dir, "rev-parse", "HEAD"); now != head {
		t.Errorf("HEAD moved from %s to %s", head, now)
	}
	assertUnstagedOnly(t, dir, "c.md")
}

// assertUnstagedOnly fails unless path is the worktree's one change and
// nothing is staged; the runner trims git's output, so status's leading
// column cannot say it.
func assertUnstagedOnly(t *testing.T, dir, path string) {
	t.Helper()
	if staged := gitOut(t, dir, "diff", "--cached", "--name-only"); staged != "" {
		t.Errorf("staged %q", staged)
	}
	if changed := gitOut(t, dir, "diff", "--name-only"); changed != path {
		t.Errorf("unstaged %q, want %q", changed, path)
	}
}

// A checkout whose branch is not the story's is refused, naming
// --autocommit, with nothing staged.
func TestCommitPathsRefusesABranchThatIsNotTheStorys(t *testing.T) {
	dir := storyCheckout(t)
	git(t, dir, "checkout", "-q", "main")
	write(t, dir, "a.md", "written by the command\n")
	got, err := CommitPaths(CommitOptions{Runner: execx.System{}, Dir: dir, Story: "S-0007", Paths: []string{"a.md"}, Subject: "record"})
	if err == nil || got != nil {
		t.Fatalf("not refused: %+v", got)
	}
	for _, want := range []string{"main", "story/S-0007", "--autocommit"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal lacks %q: %v", want, err)
		}
	}
	assertUnstagedOnly(t, dir, "a.md")
}
