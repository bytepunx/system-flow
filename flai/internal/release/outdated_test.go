package release

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// flakyTagRunner fails its first `git tag` call, as a git refused a fork on
// a loaded host does, and answers tags on every later one. Any other call,
// such as the git show of the task-done rules file, fails.
type flakyTagRunner struct {
	execx.Runner
	tags  string
	calls int
}

func (f *flakyTagRunner) Run(_, name string, args ...string) (string, error) {
	if name != "git" || len(args) == 0 || args[0] != "tag" {
		return "", errors.New("not stubbed")
	}
	f.calls++
	if f.calls == 1 {
		return "", errors.New("fork/exec git: resource temporarily unavailable")
	}
	return f.tags, nil
}

// I-0252: a look whose git failed is not kept, so the next call for the same
// project and version asks git again and finds the newer release, rather than
// serving "not outdated" for a minute.
func TestFlaiOutdatedKeepsNoAnswerFromAGitThatFailed(t *testing.T) {
	r := &flakyTagRunner{tags: "flai/v1.9.0\nflai/v1.10.0\n"}
	root := "/flai-outdated-failed-git-" + t.Name()
	if out := FlaiOutdated(r, root, "1.9.3"); out != nil {
		t.Fatalf("a failed git look answered %+v", out)
	}
	out := FlaiOutdated(r, root, "1.9.3")
	if out == nil || out.Newest != "1.10.0" || out.Running != "1.9.3" {
		t.Fatalf("the look after a failed one: %+v, git tag asked %d times", out, r.calls)
	}
	if r.calls != 2 {
		t.Errorf("git tag asked %d times, want 2", r.calls)
	}
	if again := FlaiOutdated(r, root, "1.9.3"); again != out || r.calls != 2 {
		t.Errorf("an answer git gave is not kept: %+v, git tag asked %d times", again, r.calls)
	}
}

// S-0181: the running flai is older when a flai/v* tag reachable from HEAD
// is newer; a tag on a branch HEAD does not include, a dev build, and a flai
// at or above the newest say nothing.
func TestFlaiOutdated(t *testing.T) {
	root, r := gitRepo(t)
	git := func(args ...string) {
		t.Helper()
		if out, err := r.Run(root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git("tag", "-a", "flai/v1.9.0", "-m", "flai 1.9.0")
	git("tag", "-a", "flai/v1.10.0", "-m", "flai 1.10.0")
	git("checkout", "-q", "-b", "later")
	git("commit", "-q", "--allow-empty", "-m", "later")
	git("tag", "-a", "flai/v1.11.0", "-m", "flai 1.11.0")
	git("checkout", "-q", "main")

	out := FlaiOutdated(r, root, "1.9.3")
	if out == nil || out.Running != "1.9.3" || out.Newest != "1.10.0" || !strings.HasPrefix(out.Upgrade, "flai host upgrade") ||
		!strings.Contains(out.Message, "this flai is 1.9.3, older than flai 1.10.0") || !strings.Contains(out.Message, "flai host upgrade") {
		t.Fatalf("1.9.3 under 1.10.0: %+v", out)
	}
	if out.NewestRules != 0 {
		t.Errorf("a tag whose commit has no task-done rules file carries rules %d", out.NewestRules)
	}
	for _, running := range []string{"1.10.0", "v1.10.1", "dev"} {
		if got := FlaiOutdated(r, root, running); got != nil {
			t.Errorf("%s: %+v", running, got)
		}
	}
}

// ADR-0180: the newest tag's task-done rules number is the rules file at the
// commit it tags, not at HEAD, and a file that holds no number carries none.
func TestFlaiOutdatedReadsTheTagsTaskDoneRules(t *testing.T) {
	root, r := gitRepo(t)
	commit := func(rules string) {
		t.Helper()
		p := filepath.Join(root, buildinfo.TaskDoneRulesPath)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(rules), 0o644); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", "rules " + rules}} {
			if out, err := r.Run(root, "git", args...); err != nil {
				t.Fatalf("git %v: %v %s", args, err, out)
			}
		}
	}
	tag := func(name string) {
		t.Helper()
		if out, err := r.Run(root, "git", "tag", "-a", name, "-m", name); err != nil {
			t.Fatalf("git tag %s: %v %s", name, err, out)
		}
	}
	commit("3\n")
	tag("flai/v1.12.0")
	commit("4\n")
	if out := FlaiOutdated(r, root, "1.11.0"); out == nil || out.Newest != "1.12.0" || out.NewestRules != 3 {
		t.Errorf("1.11.0 under 1.12.0 with rules 3: %+v", out)
	}
	commit("none\n")
	tag("flai/v1.13.0")
	if out := FlaiOutdated(r, root, "1.11.1"); out == nil || out.Newest != "1.13.0" || out.NewestRules != 0 {
		t.Errorf("1.11.1 under 1.13.0 with no number: %+v", out)
	}
}

// behindRepo is gitRepo with a commit fn that writes rel and commits it,
// returning the new commit's full hash, and a git fn that runs git in it.
func behindRepo(t *testing.T) (root string, r execx.Runner, commit func(rel string) string, git func(args ...string) string) {
	t.Helper()
	root, r = gitRepo(t)
	git = func(args ...string) string {
		t.Helper()
		out, err := r.Run(root, "git", args...)
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
		return strings.TrimSpace(out)
	}
	n := 0
	commit = func(rel string) string {
		t.Helper()
		n++
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(strings.Repeat("x", n)+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git("add", "-A")
		git("commit", "-q", "-m", "change "+rel)
		return git("rev-parse", "HEAD")
	}
	return root, r, commit, git
}

// ADR-0311: a flai built from main's head is not behind it, and one built
// from an older commit is behind by the commits since then under the path,
// counted from its build commit with a -dirty suffix stripped.
func TestFlaiBehindCountsFromTheBuildCommit(t *testing.T) {
	root, r, commit, _ := behindRepo(t)
	built := commit("flai/main.go")
	if b, err := FlaiBehind(r, root, "main", "flai", buildinfo.Info{Version: "1.99.0", Commit: built}); err != nil || b != nil {
		t.Fatalf("built at main's head: %+v, %v", b, err)
	}
	commit("flai/a.go")
	commit("docs/c.md")
	head := commit("flai/b.go")
	b, err := FlaiBehind(r, root, "main", "flai", buildinfo.Info{Version: "v1.99.0", Commit: built + "-dirty"})
	if err != nil || b == nil {
		t.Fatalf("built two flai commits before main: %+v, %v", b, err)
	}
	if b.By != "commit" || b.From != built || b.Commit != built+"-dirty" || b.Main != head || b.Path != "flai" ||
		b.Commits != 2 || b.Running != "1.99.0" || b.Newest != "" {
		t.Errorf("behind by its build commit: %+v", b)
	}
	for _, want := range []string{"flai 1.99.0", "commit " + built[:7], "2 commits", "under flai/", "flai release --pending", "flai host upgrade"} {
		if !strings.Contains(b.Message, want) {
			t.Errorf("message %q lacks %q", b.Message, want)
		}
	}
	if strings.Contains(b.Message, "\n") {
		t.Errorf("message spans lines: %q", b.Message)
	}
}

// ADR-0311: commits main has only outside the path do not make the running
// flai behind, and with no path every commit counts.
func TestFlaiBehindCountsOnlyUnderThePath(t *testing.T) {
	root, r, commit, _ := behindRepo(t)
	built := commit("flai/main.go")
	commit("docs/c.md")
	commit("cli/d.go")
	running := buildinfo.Info{Version: "dev", Commit: built}
	if b, err := FlaiBehind(r, root, "main", "flai", running); err != nil || b != nil {
		t.Fatalf("commits only outside flai/: %+v, %v", b, err)
	}
	b, err := FlaiBehind(r, root, "main", "", running)
	if err != nil || b == nil || b.Commits != 2 || b.Path != "" || !strings.Contains(b.Message, "in the tree") {
		t.Fatalf("the whole tree: %+v, %v", b, err)
	}
}

// ADR-0311: a flai whose build commit git does not know is compared from its
// own release tag, and a dev build from the newest flai/v* tag merged in main,
// not from a newer one on a branch main does not include.
func TestFlaiBehindFallsBackToTags(t *testing.T) {
	root, r, commit, git := behindRepo(t)
	commit("flai/main.go")
	git("tag", "-a", "flai/v1.98.0", "-m", "flai 1.98.0")
	commit("flai/a.go")
	git("tag", "-a", "flai/v1.99.0", "-m", "flai 1.99.0")
	git("checkout", "-q", "-b", "later")
	commit("flai/later.go")
	git("tag", "-a", "flai/v2.0.0", "-m", "flai 2.0.0")
	git("checkout", "-q", "main")
	commit("flai/b.go")
	commit("flai/c.go")

	for _, c := range []string{"unknown", "", "0123456789abcdef0123456789abcdef01234567"} {
		b, err := FlaiBehind(r, root, "main", "flai", buildinfo.Info{Version: "1.98.0", Commit: c})
		if err != nil || b == nil || b.By != "tag" || b.From != "flai/v1.98.0" || b.Commits != 3 || b.Newest != "1.99.0" {
			t.Fatalf("1.98.0 with commit %q: %+v, %v", c, b, err)
		}
		if c != "unknown" && b.Commit != c || c == "unknown" && b.Commit != "" {
			t.Errorf("commit %q answered as %q", c, b.Commit)
		}
		if !strings.Contains(b.Message, "flai/v1.98.0") || !strings.Contains(b.Message, "3 commits") {
			t.Errorf("message %q", b.Message)
		}
	}
	b, err := FlaiBehind(r, root, "main", "flai", buildinfo.Info{Version: "dev", Commit: "unknown"})
	if err != nil || b == nil || b.By != "newest" || b.From != "flai/v1.99.0" || b.Commits != 2 || b.Newest != "1.99.0" || b.Running != "dev" {
		t.Fatalf("a dev build: %+v, %v", b, err)
	}
	if !strings.Contains(b.Message, "flai dev") || !strings.Contains(b.Message, "flai/v1.99.0") {
		t.Errorf("message %q", b.Message)
	}
}

// ADR-0311: a flai git cannot place by build commit, release tag, or any
// flai tag merged in main is not judged behind.
func TestFlaiBehindWithNothingToCompareFrom(t *testing.T) {
	root, r, commit, _ := behindRepo(t)
	commit("flai/main.go")
	for _, running := range []buildinfo.Info{
		{Version: "dev", Commit: "unknown"},
		{Version: "1.99.0", Commit: "0123456789abcdef0123456789abcdef01234567"},
	} {
		if b, err := FlaiBehind(r, root, "main", "flai", running); err != nil || b != nil {
			t.Errorf("%+v: %+v, %v", running, b, err)
		}
	}
}

// ADR-0311: a main ref git cannot resolve is an error that names it and the
// repository, not an answer that the running flai is current.
func TestFlaiBehindReportsAMainItCannotResolve(t *testing.T) {
	root, r, commit, _ := behindRepo(t)
	built := commit("flai/main.go")
	b, err := FlaiBehind(r, root, "no-such-branch", "flai", buildinfo.Info{Version: "1.99.0", Commit: built})
	if err == nil || b != nil {
		t.Fatalf("an unknown main ref: %+v, %v", b, err)
	}
	if !strings.Contains(err.Error(), "no-such-branch") || !strings.Contains(err.Error(), root) {
		t.Errorf("error %q names neither the ref nor the repository", err)
	}
}
