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

// zeroFiles leaves every file under root at zero bytes, as the host crash of
// I-0137 left what a story's start had written.
func zeroFiles(t *testing.T, root string) {
	t.Helper()
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		return os.WriteFile(p, nil, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// crashStart opens story id and leaves what its start wrote as I-0137 found
// it: the branch's ref and reflog, the worktree's files and its .git file,
// and its registration in .git/worktrees, all empty.
func crashStart(t *testing.T, repo *workitem.Repo, id string) Opened {
	t.Helper()
	opened, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: id})
	if err != nil {
		t.Fatal(err)
	}
	gitDir := filepath.Join(repo.MainRoot, ".git")
	for _, p := range []string{opened.Worktree, filepath.Join(gitDir, "worktrees", id), filepath.Join(gitDir, "refs", "heads", "story", id), filepath.Join(gitDir, "logs", "refs", "heads", "story", id)} {
		zeroFiles(t, p)
	}
	if BranchExists(execx.System{}, repo.MainRoot, opened.Branch) {
		t.Fatal("an empty ref should not resolve")
	}
	return opened
}

// What a crash during a story's start leaves (I-0137), an empty branch ref
// and reflog and a worktree and registration of empty files, holds no
// commits: Open removes it, logs what it removed, and opens the story again
// from the main branch.
func TestOpenRepairsWhatACrashedStartLeft(t *testing.T) {
	repo := openRepo(t)
	crashed := crashStart(t, repo, "S-0001")
	var log bytes.Buffer
	got, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001", Log: slog.New(slog.NewTextHandler(&log, nil))})
	if err != nil {
		t.Fatalf("Open over a crashed start: %v", err)
	}
	want := Opened{Branch: crashed.Branch, Worktree: crashed.Worktree, From: FromMain}
	if got != want {
		t.Errorf("Open = %+v, want %+v", got, want)
	}
	if h := headOf(t, got.Worktree); h != "story/S-0001" {
		t.Errorf("worktree is on %q", h)
	}
	if b, err := os.ReadFile(filepath.Join(got.Worktree, "a.md")); err != nil || string(b) != "one\n" {
		t.Errorf("the worktree should hold the main branch's files: %q %v", b, err)
	}
	if !IsWorktree(execx.System{}, got.Worktree) {
		t.Error("git should know the worktree")
	}
	regs, err := os.ReadDir(filepath.Join(repo.MainRoot, ".git", "worktrees"))
	if err != nil || len(regs) != 1 || regs[0].Name() != "S-0001" {
		t.Errorf("the empty registration should be pruned and the worktree registered as S-0001: %v %v", regs, err)
	}
	l := log.String()
	for _, field := range []string{"level=WARN", "component=git", "worktree=.flai-cache/worktrees/S-0001", "ref=.git/refs/heads/story/S-0001", "reflog=removed"} {
		if !strings.Contains(l, field) {
			t.Errorf("the log should carry %s: %s", field, l)
		}
	}
}

// A folder git does not know as the story's worktree, whose branch holds
// commits, may hold work: Open refuses it, says how to clear it, and removes
// nothing.
func TestOpenRefusesAnUnknownFolderWhoseBranchHoldsCommits(t *testing.T) {
	repo := openRepo(t)
	git(t, repo.MainRoot, "checkout", "-q", "-b", "story/S-0001")
	commitFiles(t, repo.MainRoot, "feat: [S-0001] begun", "docs/begun.md")
	git(t, repo.MainRoot, "checkout", "-q", "main")
	wt := repo.WorktreePath("S-0001")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, wt, ".git", "")
	write(t, wt, "work.md", "work\n")

	_, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001", FromRemote: true})
	if err == nil || !strings.Contains(err.Error(), ".flai-cache/worktrees/S-0001") || !strings.Contains(err.Error(), "rm -rf "+wt) || !strings.Contains(err.Error(), "flai stream open S-0001") {
		t.Fatalf("Open = %v, want a refusal naming the folder and how to clear it", err)
	}
	if b, err := os.ReadFile(filepath.Join(wt, "work.md")); err != nil || string(b) != "work\n" {
		t.Errorf("the folder should be left as it was: %q %v", b, err)
	}
	if !BranchExists(execx.System{}, repo.MainRoot, "story/S-0001") {
		t.Error("the branch should be left as it was")
	}
}

// A crash after git made the branch, whose ref is synced, can still leave the
// checkout's files empty: a folder git does not know whose every file is
// empty holds no work, so Open makes the worktree again on the branch.
func TestOpenRepairsAnEmptyWorktreeOverAWholeBranch(t *testing.T) {
	repo := openRepo(t)
	git(t, repo.MainRoot, "branch", "story/S-0001", "main")
	wt := repo.WorktreePath("S-0001")
	if err := os.MkdirAll(filepath.Join(wt, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, wt, ".git", "")
	write(t, wt, "a.md", "")
	write(t, wt, "docs/b.md", "")

	got, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001"})
	if err != nil {
		t.Fatalf("Open over an empty worktree: %v", err)
	}
	if got.From != FromLocal || headOf(t, got.Worktree) != "story/S-0001" {
		t.Errorf("Open = %+v, want the local story/S-0001", got)
	}
	if b, err := os.ReadFile(filepath.Join(got.Worktree, "a.md")); err != nil || string(b) != "one\n" {
		t.Errorf("the worktree should hold the branch's files: %q %v", b, err)
	}
}

// A broken branch ref with no worktree is removed and the branch made again
// from the main branch; a reflog with entries is kept, and git appends to it.
func TestOpenRepairsABrokenRefAlone(t *testing.T) {
	repo := openRepo(t)
	main, err := (execx.System{}).Run(repo.MainRoot, "git", "rev-parse", "main")
	if err != nil {
		t.Fatal(main, err)
	}
	gitDir := filepath.Join(repo.MainRoot, ".git")
	for _, d := range []string{"refs/heads/story", "logs/refs/heads/story"} {
		if err := os.MkdirAll(filepath.Join(gitDir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, gitDir, "refs/heads/story/S-0001", "")
	entry := strings.Repeat("0", len(main)) + " " + main + " t <t@t> 1700000000 +0000\tbranch: Created from main\n"
	write(t, gitDir, "logs/refs/heads/story/S-0001", entry)
	var log bytes.Buffer

	got, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001", Log: slog.New(slog.NewTextHandler(&log, nil))})
	if err != nil {
		t.Fatalf("Open over a broken ref: %v", err)
	}
	if got.From != FromMain || headOf(t, got.Worktree) != "story/S-0001" {
		t.Errorf("Open = %+v, want a new story/S-0001 from main", got)
	}
	if b, err := os.ReadFile(filepath.Join(gitDir, "logs/refs/heads/story/S-0001")); err != nil || !strings.HasPrefix(string(b), entry) || strings.Count(string(b), "\n") != 2 {
		t.Errorf("the reflog should keep its entry and gain the new branch's: %q %v", b, err)
	}
	if l := log.String(); !strings.Contains(l, "ref=.git/refs/heads/story/S-0001") || !strings.Contains(l, "reflog=kept") {
		t.Errorf("log: %s", l)
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
