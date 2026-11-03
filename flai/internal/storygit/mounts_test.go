package storygit

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// handLinked lays out by hand what git worktree add leaves: a main checkout
// whose .git holds the registration worktrees/S-0001 with the commondir that
// commondir gives for the common dir, none when "", and a worktree beside it
// whose .git file holds gitdir. It returns the worktree, the common dir, and
// the registration.
func handLinked(t *testing.T, gitdir func(reg, wt string) string, commondir func(common string) string) (wt, common, reg string) {
	t.Helper()
	root := t.TempDir()
	common = filepath.Join(root, "main", ".git")
	reg = filepath.Join(common, "worktrees", "S-0001")
	wt = filepath.Join(root, "S-0001")
	for _, d := range []string{reg, wt} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if c := commondir(common); c != "" {
		write(t, reg, "commondir", c+"\n")
	}
	write(t, wt, ".git", "gitdir: "+gitdir(reg, wt)+"\n")
	return wt, common, reg
}

func absGitdir(reg, _ string) string { return reg }

// upTwo is the commondir git writes, the registration's grandparent.
func upTwo(string) string { return "../.." }

func relGitdir(reg, wt string) string {
	rel, err := filepath.Rel(wt, reg)
	if err != nil {
		panic(err)
	}
	return rel
}

// wantMounts is the common dir read-only, then the registration read-write.
func wantMounts(common, reg string) []Mount {
	return []Mount{{Path: common, ReadOnly: true}, {Path: reg}}
}

// A worktree's .git file names its registration, absolutely or relative to
// the worktree, and the registration's commondir names the common dir,
// relative to it or absolutely: either way the mounts are the same cleaned
// absolute host paths.
func TestGitMountsResolvesTheRegistrationAndTheCommonDir(t *testing.T) {
	for _, tc := range []struct {
		name      string
		gitdir    func(reg, wt string) string
		commondir func(common string) string
	}{
		{"absolute gitdir", absGitdir, upTwo},
		{"relative gitdir", relGitdir, upTwo},
		{"absolute commondir", relGitdir, func(common string) string { return common + "/" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wt, common, reg := handLinked(t, tc.gitdir, tc.commondir)
			got, err := GitMounts(wt)
			if err != nil {
				t.Fatal(err)
			}
			if want := wantMounts(common, reg); !reflect.DeepEqual(got, want) {
				t.Errorf("GitMounts = %+v, want %+v", got, want)
			}
		})
	}
}

// A worktree whose .git is a folder carries its repository: nothing more to
// mount, and no error.
func TestGitMountsOfANormalCheckoutIsNone(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, err := GitMounts(dir); err != nil || got != nil {
		t.Errorf("GitMounts = %+v, %v; want none", got, err)
	}
}

// What would start a container without git is an error naming the path at
// fault and what to do.
func TestGitMountsNamesWhatIsMissing(t *testing.T) {
	t.Run("no .git", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GitMounts(dir)
		wantErr(t, err, filepath.Join(dir, ".git"), "flai stream open")
	})
	t.Run("no gitdir line", func(t *testing.T) {
		dir := t.TempDir()
		write(t, dir, ".git", "not a link\n")
		_, err := GitMounts(dir)
		wantErr(t, err, filepath.Join(dir, ".git"), "git worktree repair")
	})
	t.Run("registration absent", func(t *testing.T) {
		wt, _, reg := handLinked(t, absGitdir, upTwo)
		if err := os.RemoveAll(reg); err != nil {
			t.Fatal(err)
		}
		_, err := GitMounts(wt)
		wantErr(t, err, reg, "git worktree repair")
	})
	t.Run("relative registration absent", func(t *testing.T) {
		wt, _, reg := handLinked(t, relGitdir, upTwo)
		if err := os.RemoveAll(reg); err != nil {
			t.Fatal(err)
		}
		_, err := GitMounts(wt)
		wantErr(t, err, reg, "git worktree repair")
	})
	t.Run("no commondir", func(t *testing.T) {
		wt, _, reg := handLinked(t, absGitdir, func(string) string { return "" })
		_, err := GitMounts(wt)
		wantErr(t, err, filepath.Join(reg, "commondir"), "git worktree repair")
	})
	t.Run("common dir absent", func(t *testing.T) {
		wt, _, reg := handLinked(t, absGitdir, func(string) string { return "../../../gone" })
		_, err := GitMounts(wt)
		wantErr(t, err, filepath.Join(reg, "commondir"), filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(reg))), "gone"))
	})
}

// wantErr fails unless err holds every one of parts.
func wantErr(t *testing.T, err error, parts ...string) {
	t.Helper()
	if err == nil {
		t.Fatalf("want an error naming %q", parts)
	}
	for _, p := range parts {
		if !strings.Contains(err.Error(), p) {
			t.Errorf("error should name %q: %v", p, err)
		}
	}
}

// isolateGit keeps the host's system and global git config, which may link
// worktrees relatively or sign commits, out of a real-git test.
func isolateGit(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
}

// openedStory is a project whose story S-0001 Open made a worktree for,
// linked relatively when relative is set. Integration: it runs real git.
func openedStory(t *testing.T, relative bool) (*workitem.Repo, Opened) {
	t.Helper()
	isolateGit(t)
	repo := openRepo(t)
	opened, err := Open(OpenOptions{Runner: execx.System{}, Repo: repo, Story: "S-0001", RelativePaths: func() bool { return relative }})
	if err != nil {
		t.Fatal(err)
	}
	return repo, opened
}

// realPaths is ms with each path's symlinks resolved, as a temporary folder
// may sit behind one.
func realPaths(t *testing.T, ms []Mount) []Mount {
	t.Helper()
	out := make([]Mount, len(ms))
	for i, m := range ms {
		p, err := filepath.EvalSymlinks(m.Path)
		if err != nil {
			t.Fatal(err)
		}
		out[i] = Mount{Path: p, ReadOnly: m.ReadOnly}
	}
	return out
}

// storyMounts is what a worktree git made needs: the main checkout's .git
// read-only and the story's registration read-write.
func storyMounts(t *testing.T, repo *workitem.Repo) []Mount {
	t.Helper()
	common := filepath.Join(repo.MainRoot, ".git")
	return realPaths(t, wantMounts(common, filepath.Join(common, "worktrees", "S-0001")))
}

// dotGitOf is what the worktree's .git file says.
func dotGitOf(t *testing.T, wt string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(wt, ".git"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(b))
}

// A worktree git worktree add made, as flai stream open makes it, gives the
// main checkout's .git read-only and its registration read-write.
func TestGitMountsOfAWorktreeGitAdded(t *testing.T) {
	repo, opened := openedStory(t, false)
	if g := dotGitOf(t, opened.Worktree); !strings.HasPrefix(g, "gitdir: /") {
		t.Fatalf("the worktree should be linked absolutely: %s", g)
	}
	got, err := GitMounts(opened.Worktree)
	if err != nil {
		t.Fatal(err)
	}
	if want := storyMounts(t, repo); !reflect.DeepEqual(realPaths(t, got), want) {
		t.Errorf("GitMounts = %+v, want %+v", got, want)
	}
}

// A worktree linked with relative paths (ADR-0022) gives the same absolute
// host paths.
func TestGitMountsOfARelativePathWorktree(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if v, err := gitver.Installed(execx.System{}); err != nil || !v.AtLeast(gitver.RelativeWorktrees) {
		t.Skipf("git %v is older than %s, which first links worktrees relatively: %v", v, gitver.RelativeWorktrees, err)
	}
	repo, opened := openedStory(t, true)
	if g := dotGitOf(t, opened.Worktree); strings.HasPrefix(g, "gitdir: /") {
		t.Fatalf("the worktree should be linked relatively: %s", g)
	}
	got, err := GitMounts(opened.Worktree)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range got {
		if !filepath.IsAbs(m.Path) || filepath.Clean(m.Path) != m.Path {
			t.Errorf("%s should be a cleaned absolute path", m.Path)
		}
	}
	if want := storyMounts(t, repo); !reflect.DeepEqual(realPaths(t, got), want) {
		t.Errorf("GitMounts = %+v, want %+v", got, want)
	}
}

// chmodTree sets every folder under root to dirMode and every file to
// fileMode, leaving out the folder skip and what is under it.
func chmodTree(t *testing.T, root, skip string, dirMode, fileMode os.FileMode) {
	t.Helper()
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == skip {
			return filepath.SkipDir
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() {
			return os.Chmod(p, dirMode)
		}
		return os.Chmod(p, fileMode)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// gitIn runs git in dir and returns its output and error, for a command the
// test expects to fail.
func gitIn(dir string, args ...string) (string, error) {
	return (execx.System{}).Run(dir, "git", args...)
}

// I-0215: with the mounts' modes on the host, the common dir read-only and
// the story's registration writable over it, git in the worktree reads its
// status and history, and refuses a commit with its own write error, as the
// host commits for the story.
func TestGitWorksReadOnlyOnTheMountsAndRefusesACommit(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores file modes, so the common dir cannot be made read-only")
	}
	_, opened := openedStory(t, false)
	mounts, err := GitMounts(opened.Worktree)
	if err != nil {
		t.Fatal(err)
	}
	common, reg := mounts[0], mounts[1]
	if !common.ReadOnly || reg.ReadOnly {
		t.Fatalf("mounts = %+v, want the common dir read-only, then the registration read-write", mounts)
	}
	head, err := gitIn(opened.Worktree, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(head, err)
	}
	// restore before TempDir's own cleanup, which runs after, removes it
	t.Cleanup(func() { chmodTree(t, common.Path, "", 0o755, 0o644) })
	chmodTree(t, common.Path, reg.Path, 0o555, 0o444)

	write(t, opened.Worktree, "a.md", "changed in the container\n")
	if out, err := gitIn(opened.Worktree, "status", "--porcelain"); err != nil || !strings.Contains(out, "a.md") {
		t.Errorf("git status should work and show the change: %q %v", out, err)
	}
	if out, err := gitIn(opened.Worktree, "log", "--oneline"); err != nil || !strings.Contains(out, "ignore the cache") {
		t.Errorf("git log should work and show the history: %q %v", out, err)
	}
	if out, err := gitIn(opened.Worktree, "commit", "--allow-empty", "-q", "-m", "x"); err == nil {
		t.Errorf("git commit should be refused on a read-only common dir: %q", out)
	}
	if now, err := gitIn(opened.Worktree, "rev-parse", "HEAD"); err != nil || now != head {
		t.Errorf("HEAD moved from %s to %s: %v", head, now, err)
	}
}

// I-0215's shape: a worktree whose registration is absent, as a container
// that mounts only the worktree sees it, is not a git repository to git, and
// GitMounts refuses it, naming the registration.
func TestAWorktreeWithoutItsRegistrationIsNotARepository(t *testing.T) {
	repo, opened := openedStory(t, false)
	reg := filepath.Join(repo.MainRoot, ".git", "worktrees", "S-0001")
	if err := os.Rename(reg, reg+".away"); err != nil {
		t.Fatal(err)
	}
	// git before 2.56 says the worktree is not a git repository; from 2.56 it
	// says its gitfile does not point to a valid repository
	out, err := gitIn(opened.Worktree, "status")
	refused := strings.Contains(out, "not a git repository") || strings.Contains(out, "does not point to a valid repository")
	if err == nil || !refused {
		t.Errorf("git status without the registration = %q, %v; want it refused as not a repository", out, err)
	}
	_, err = GitMounts(opened.Worktree)
	wantErr(t, err, reg, "git worktree repair")
}
