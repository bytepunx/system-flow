package storygit

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// noGit is a host without git.
type noGit struct{ execx.System }

func (noGit) LookPath(string) (string, error) { return "", errors.New("not found") }

func (noGit) Run(string, string, ...string) (string, error) {
	return "", errors.New("git is not installed")
}

// Without git there is no branch to open, and nothing is touched.
func TestOpenWithoutGitOpensNothing(t *testing.T) {
	dir := t.TempDir()
	repo := &workitem.Repo{Root: dir, MainRoot: dir}
	got, err := Open(OpenOptions{Runner: noGit{}, Repo: repo, Story: "S-0001", FromRemote: true})
	if err != nil || got != (Opened{}) {
		t.Fatalf("Open = %+v, %v; want nothing", got, err)
	}
	if _, err := os.Stat(repo.WorktreePath("S-0001")); !os.IsNotExist(err) {
		t.Errorf("no worktree should be made: %v", err)
	}
}

// openRepo is a git repository on main, with .flai-cache ignored, as a
// project's main checkout. Integration: it runs real git.
func openRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	dir := gitRepo(t)
	write(t, dir, ".gitignore", ".flai-cache/\n")
	git(t, dir, "add", ".gitignore")
	git(t, dir, "commit", "-q", "-m", "chore: ignore the cache")
	return &workitem.Repo{Root: dir, MainRoot: dir}
}

// headOf is the branch checked out at dir.
func headOf(t *testing.T, dir string) string {
	t.Helper()
	out, err := (execx.System{}).Run(dir, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		t.Fatal(out, err)
	}
	return out
}

// A story with no branch anywhere gets a new one from the main branch, and
// the opening is logged.
func TestOpenCreatesTheBranchFromMain(t *testing.T) {
	repo := openRepo(t)
	var log bytes.Buffer
	got, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001", Log: slog.New(slog.NewTextHandler(&log, nil))})
	if err != nil {
		t.Fatal(err)
	}
	want := Opened{Branch: "story/S-0001", Worktree: repo.WorktreePath("S-0001"), From: FromMain}
	if got != want {
		t.Errorf("Open = %+v, want %+v", got, want)
	}
	if h := headOf(t, got.Worktree); h != "story/S-0001" {
		t.Errorf("worktree is on %q", h)
	}
	if l := log.String(); !strings.Contains(l, "story branch opened") || !strings.Contains(l, "worktree=.flai-cache/worktrees/S-0001") || !strings.Contains(l, "base=main") {
		t.Errorf("log: %s", l)
	}
}

// This clone's branch is checked out as it is, with its commits.
func TestOpenChecksOutTheLocalBranch(t *testing.T) {
	repo := openRepo(t)
	git(t, repo.MainRoot, "checkout", "-q", "-b", "story/S-0001")
	commitFiles(t, repo.MainRoot, "feat: [S-0001] begun", "docs/begun.md")
	git(t, repo.MainRoot, "checkout", "-q", "main")
	got, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001"})
	if err != nil {
		t.Fatal(err)
	}
	if got.From != FromLocal || got.Branch != "story/S-0001" {
		t.Errorf("Open = %+v, want this clone's story/S-0001", got)
	}
	if _, err := os.Stat(filepath.Join(got.Worktree, "docs", "begun.md")); err != nil {
		t.Errorf("the worktree should hold the branch's commit: %v", err)
	}
}

// A worktree that is there already is taken as it is: nothing is added, and
// whether to link with relative paths is not asked.
func TestOpenReusesAnExistingWorktree(t *testing.T) {
	repo := openRepo(t)
	first, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001"})
	if err != nil {
		t.Fatal(err)
	}
	asked := false
	again, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001", FromRemote: true, RelativePaths: func() bool { asked = true; return true }})
	if err != nil {
		t.Fatal(err)
	}
	want := Opened{Branch: first.Branch, Worktree: first.Worktree}
	if again != want {
		t.Errorf("Open again = %+v, want %+v", again, want)
	}
	if asked {
		t.Error("relative paths asked for a worktree that was not added")
	}
}

// A story begun on another host is fetched from the remote when this clone
// has no branch for it (ADR-0064), past a worktree removed without git.
func TestOpenFetchesTheBranchFromTheRemote(t *testing.T) {
	there := openRepo(t)
	remote := filepath.Join(t.TempDir(), "remote.git")
	git(t, filepath.Dir(remote), "init", "-q", "--bare", "-b", "main", remote)
	git(t, there.MainRoot, "remote", "add", "origin", remote)
	git(t, there.MainRoot, "checkout", "-q", "-b", "story/S-0001")
	commitFiles(t, there.MainRoot, "feat: [S-0001] begun there", "docs/begun.md")
	git(t, there.MainRoot, "checkout", "-q", "main")
	git(t, there.MainRoot, "push", "-q", "origin", "main", "story/S-0001")

	hereDir := filepath.Join(t.TempDir(), "here")
	git(t, filepath.Dir(hereDir), "clone", "-q", remote, hereDir)
	here := &workitem.Repo{Root: hereDir, MainRoot: hereDir}
	// S-0002's worktree removed without git: its registration would refuse
	// the branch a new worktree unless pruned
	stale, err := Open(OpenOptions{Runner: execx.System{}, Repo: here, Story: "S-0002"})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(stale.Worktree); err != nil {
		t.Fatal(err)
	}

	got, err := Open(OpenOptions{Runner: execx.System{}, Repo: here, Story: "S-0001", FromRemote: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.From != "origin" || got.Branch != "story/S-0001" {
		t.Errorf("Open = %+v, want story/S-0001 fetched from origin", got)
	}
	if b, err := os.ReadFile(filepath.Join(got.Worktree, "docs", "begun.md")); err != nil || string(b) != "feat: [S-0001] begun there\n" {
		t.Errorf("the worktree should hold what was committed there: %q %v", b, err)
	}

	got, err = Open(OpenOptions{Runner: execx.System{}, Repo: here, Story: "S-0002", FromRemote: true})
	if err != nil {
		t.Fatalf("reopen past a worktree removed without git: %v", err)
	}
	if got.From != FromLocal {
		t.Errorf("S-0002 reopened from %q, want this clone's branch", got.From)
	}
}
