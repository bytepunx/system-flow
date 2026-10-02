package release

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// withRemote is gitRepo's repository with a bare remote, origin, holding its
// main and every tag, and a way to run git in either; the kept remote
// answers are forgotten so each test asks afresh.
func withRemote(t *testing.T) (root, remote string, git func(dir string, args ...string)) {
	t.Helper()
	root, r := gitRepo(t)
	remote = filepath.Join(t.TempDir(), "origin.git")
	git = func(dir string, args ...string) {
		t.Helper()
		if out, err := r.Run(dir, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git(root, "init", "-q", "--bare", remote)
	git(root, "remote", "add", "origin", remote)
	git(root, "push", "-q", "-u", "origin", "main", "--tags")
	forgetRemotes(t)
	return root, remote, git
}

func forgetRemotes(t *testing.T) {
	t.Helper()
	remoteKept.Lock()
	remoteKept.answers = map[string]remoteAnswer{}
	remoteKept.Unlock()
}

// S-0174: a clone whose main moved without the tags published from another
// clone is told which components lag and how to fetch their tags.
func TestCheckRemoteFindsTagsThisCloneLacks(t *testing.T) {
	root, remote, git := withRemote(t)
	git(remote, "tag", "cli/v0.12.0", "main") // published from another clone
	git(remote, "tag", "web/v1.0.0", "main")  // a first release here has none
	got := CheckRemote(execx.System{}, root, m)
	if !got.Lagging() || got.Remote != "origin" || got.Fix != "git fetch --tags origin" || got.Unchecked != "" {
		t.Fatalf("lagging clone: %+v", got)
	}
	want := []Lag{{Component: "cli", Local: "cli/v0.10.0", Remote: "cli/v0.12.0"}, {Component: "web", Remote: "web/v1.0.0"}}
	if len(got.Behind) != len(want) || got.Behind[0] != want[0] || got.Behind[1] != want[1] {
		t.Errorf("behind: %+v, want %+v", got.Behind, want)
	}
	for _, s := range []string{"cli/v0.12.0 (here cli/v0.10.0)", "web/v1.0.0 (here none)", "git fetch --tags origin"} {
		if !strings.Contains(got.Message, s) {
			t.Errorf("message %q lacks %q", got.Message, s)
		}
	}

	// fetched: in step, nothing to say
	git(root, "fetch", "-q", "--tags", "origin")
	forgetRemotes(t)
	if got := CheckRemote(execx.System{}, root, m); got != nil {
		t.Errorf("after fetching the tags: %+v, want nil", got)
	}
}

// A clone in step with its remote, or ahead of it with a tag not yet pushed,
// has nothing to say; nor has one with no remote at all.
func TestCheckRemoteInStep(t *testing.T) {
	root, _, git := withRemote(t)
	if got := CheckRemote(execx.System{}, root, m); got != nil {
		t.Errorf("in step: %+v", got)
	}
	git(root, "tag", "-a", "cli/v0.11.0", "-m", "not pushed yet")
	forgetRemotes(t)
	if got := CheckRemote(execx.System{}, root, m); got != nil {
		t.Errorf("ahead of the remote: %+v", got)
	}

	alone, _ := gitRepo(t)
	if got := CheckRemote(execx.System{}, alone, m); got != nil {
		t.Errorf("no remote: %+v", got)
	}
}

// A remote that cannot be reached is said, with git's reason, and nothing is
// claimed to lag.
func TestCheckRemoteUnreachable(t *testing.T) {
	root, _, git := withRemote(t)
	git(root, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))
	got := CheckRemote(execx.System{}, root, m)
	if got == nil || got.Lagging() || got.Unchecked == "" || !strings.Contains(got.Message, "could not ask origin") {
		t.Fatalf("unreachable remote: %+v", got)
	}
	if strings.Contains(got.Unchecked, "\n") {
		t.Errorf("unchecked is one line: %q", got.Unchecked)
	}
}

// stalling is a git whose ls-remote never answers until released.
type stalling struct {
	execx.Runner
	release chan struct{}
}

func (s stalling) Run(dir, name string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "ls-remote" {
		<-s.release
		return "", nil
	}
	return s.Runner.Run(dir, name, args...)
}

// A remote that does not answer is not waited on past the timeout.
func TestCheckRemoteTimesOut(t *testing.T) {
	root, _, _ := withRemote(t)
	defer func(d time.Duration) { remoteTimeout = d }(remoteTimeout)
	remoteTimeout = 50 * time.Millisecond
	s := stalling{Runner: execx.System{}, release: make(chan struct{})}
	defer close(s.release)
	got := CheckRemote(s, root, m)
	if got == nil || !strings.Contains(got.Unchecked, "did not answer") {
		t.Errorf("stalled remote: %+v", got)
	}
}

// countingRemote counts the times the remote is asked.
type countingRemote struct {
	execx.Runner
	mu    *sync.Mutex
	asked *int
}

func (c countingRemote) Run(dir, name string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "ls-remote" {
		c.mu.Lock()
		*c.asked++
		c.mu.Unlock()
	}
	return c.Runner.Run(dir, name, args...)
}

// The remote's answer is kept a while: asking again soon after does not ask
// the network again, and a later question does.
func TestCheckRemoteKeepsTheAnswerAWhile(t *testing.T) {
	root, _, _ := withRemote(t)
	asked := 0
	c := countingRemote{Runner: execx.System{}, mu: &sync.Mutex{}, asked: &asked}
	CheckRemote(c, root, m)
	CheckRemote(c, root, m)
	if asked != 1 {
		t.Errorf("asked %d times within the TTL, want 1", asked)
	}
	defer func(d time.Duration) { RemoteTTL = d }(RemoteTTL)
	RemoteTTL = 0
	CheckRemote(c, root, m)
	if asked != 2 {
		t.Errorf("asked %d times after the TTL, want 2", asked)
	}
}

// pushFromElsewhere commits on main in another clone of remote and pushes
// it, tagging that commit with tags first, and returns the commit.
func pushFromElsewhere(t *testing.T, remote string, git func(dir string, args ...string), tags ...string) string {
	t.Helper()
	other := filepath.Join(t.TempDir(), "other")
	git(filepath.Dir(other), "clone", "-q", "-b", "main", remote, other)
	git(other, "config", "user.email", "o@o")
	git(other, "config", "user.name", "o")
	git(other, "commit", "-q", "--allow-empty", "-m", "docs: from elsewhere")
	for _, tag := range tags {
		git(other, "tag", tag)
	}
	git(other, "push", "-q", "origin", "main", "--tags")
	head, err := execx.System{}.Run(other, "git", "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(head)
}

// ADR-0067: a clone whose remote branch has commits it lacks is told to
// fetch and merge or rebase before it publishes, whether or not it has
// fetched them, and is in step once it has merged them.
func TestCheckRemoteFindsBranchCommitsThisCloneLacks(t *testing.T) {
	root, remote, git := withRemote(t)
	head := pushFromElsewhere(t, remote, git)

	got := CheckRemote(execx.System{}, root, m)
	if !got.Lagging() || len(got.Behind) != 0 || got.Unchecked != "" || got.Fix != "git fetch origin && git merge origin/main" {
		t.Fatalf("remote branch moved: %+v", got)
	}
	if *got.Branch != (BranchLag{Upstream: "origin/main", Head: head, Fetched: false}) {
		t.Errorf("branch: %+v", *got.Branch)
	}
	for _, s := range []string{"origin/main has commits this clone lacks", head[:12], "not fetched here", "git fetch origin && git merge origin/main", "rebase onto origin/main", "flai release --pending again"} {
		if !strings.Contains(got.Message, s) {
			t.Errorf("message %q lacks %q", got.Message, s)
		}
	}
	if !strings.Contains(got.Refusal(), "refusing to publish") {
		t.Errorf("refusal: %q", got.Refusal())
	}

	// fetched, not merged: only the merge is left to run
	git(root, "fetch", "-q", "origin")
	forgetRemotes(t)
	got = CheckRemote(execx.System{}, root, m)
	if !got.Lagging() || got.Branch == nil || !got.Branch.Fetched || got.Fix != "git merge origin/main" || !strings.Contains(got.Message, "fetched here but not merged") {
		t.Fatalf("fetched, not merged: %+v", got)
	}

	// merged: in step, even with the remote's answer kept from before
	git(root, "merge", "-q", "--no-edit", "origin/main")
	if got := CheckRemote(execx.System{}, root, m); got != nil {
		t.Errorf("after merging: %+v, want nil", got)
	}

	// ahead of the remote: nothing to say
	git(root, "commit", "-q", "--allow-empty", "-m", "chore: not pushed yet")
	if got := CheckRemote(execx.System{}, root, m); got != nil {
		t.Errorf("ahead of the remote: %+v, want nil", got)
	}
}

// Tags and branch both behind: the message says both, and Fix is one
// command that fetches the tags and merges.
func TestCheckRemoteTagsAndBranchBehind(t *testing.T) {
	root, remote, git := withRemote(t)
	pushFromElsewhere(t, remote, git, "cli/v0.12.0")
	got := CheckRemote(execx.System{}, root, m)
	if !got.Lagging() || len(got.Behind) != 1 || got.Branch == nil || got.Fix != "git fetch --tags origin && git merge origin/main" {
		t.Fatalf("tags and branch behind: %+v", got)
	}
	for _, s := range []string{"cli/v0.12.0 (here cli/v0.10.0)", "origin/main has commits this clone lacks", "git fetch --tags origin && git merge origin/main"} {
		if !strings.Contains(got.Message, s) {
			t.Errorf("message %q lacks %q", got.Message, s)
		}
	}
}

// A branch the remote does not have is not behind it.
func TestCheckRemoteWithoutTheBranch(t *testing.T) {
	root, remote, git := withRemote(t)
	pushFromElsewhere(t, remote, git) // main moved, but this clone is not on main
	git(root, "checkout", "-q", "-b", "work")
	git(root, "commit", "-q", "--allow-empty", "-m", "chore: local work")
	forgetRemotes(t)
	if got := CheckRemote(execx.System{}, root, m); got != nil {
		t.Errorf("a branch the remote lacks: %+v, want nil", got)
	}
}
