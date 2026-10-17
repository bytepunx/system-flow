package serve

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeForge is a forge with one repository, keeping its pull requests by
// head branch and counting what was asked of it; heads are the head commit
// of each head branch, and checks and reviews what each pull request has,
// by its head branch. readErr fails the reads of checks; approvals is what
// every base branch requires.
type fakeForge struct {
	mu        sync.Mutex
	prs       map[string]forge.PullRequest
	opens     int
	down      error
	readErr   error
	number    int
	heads     map[string]string
	checks    map[string][]forge.Check
	reviews   map[string][]forge.Review
	approvals int
}

func newFakeForge() *fakeForge {
	return &fakeForge{prs: map[string]forge.PullRequest{}, heads: map[string]string{}, checks: map[string][]forge.Check{}, reviews: map[string][]forge.Review{}}
}

func (f *fakeForge) Name() string       { return "github" }
func (f *fakeForge) Repository() string { return "o/r" }
func (f *fakeForge) Access(context.Context) (forge.Access, error) {
	return forge.Access{Read: true, Push: true}, nil
}

func (f *fakeForge) FindPullRequest(_ context.Context, head string) (forge.PullRequest, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.down != nil {
		return forge.PullRequest{}, false, f.down
	}
	pr, ok := f.prs[head]
	pr.HeadSHA = f.heads[head]
	return pr, ok, nil
}

func (f *fakeForge) OpenPullRequest(_ context.Context, in forge.PullRequestInput) (forge.PullRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.opens++
	pr, ok := f.prs[in.Head]
	if !ok {
		f.number++
		pr = forge.PullRequest{Number: f.number, URL: fmt.Sprintf("https://github.com/o/r/pull/%d", f.number), Head: in.Head, State: forge.PullRequestOpen}
	}
	pr.Base, pr.Title, pr.Body = in.Base, in.Title, in.Body
	f.prs[in.Head] = pr
	pr.HeadSHA = f.heads[in.Head]
	return pr, nil
}

func (f *fakeForge) Checks(_ context.Context, pr forge.PullRequest) ([]forge.Check, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.readErr != nil {
		return nil, f.readErr
	}
	return f.checks[pr.Head], nil
}
func (f *fakeForge) Reviews(_ context.Context, pr forge.PullRequest) ([]forge.Review, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.reviews[pr.Head], nil
}
func (f *fakeForge) RequiredApprovals(context.Context, string) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.approvals, nil
}

// set gives the pull request of the head branch head the head commit sha,
// and the checks and reviews given.
func (f *fakeForge) set(head, sha string, checks []forge.Check, reviews []forge.Review) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.heads[head], f.checks[head], f.reviews[head] = sha, checks, reviews
}
func (f *fakeForge) Merge(context.Context, forge.PullRequest, forge.MergeMethod) (forge.Merged, error) {
	return forge.Merged{}, forge.ErrMergeRefused
}

func (f *fakeForge) pr(head string) (forge.PullRequest, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	pr, ok := f.prs[head]
	return pr, ok
}

// fakeGit answers the git commands a push runs from its fields, and records
// each push it is asked for.
type fakeGit struct {
	mu       sync.Mutex
	local    string
	remote   string
	ancestor bool
	pushErr  error
	calls    int
	pushes   [][]string
}

func (g *fakeGit) Run(_, name string, args ...string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.calls++
	if name != "git" || len(args) == 0 {
		return "", fmt.Errorf("unexpected %s %v", name, args)
	}
	switch args[0] {
	case "remote":
		return "git@github.com:o/r.git", nil
	case "rev-parse":
		if args[1] == "--abbrev-ref" {
			return "main", nil
		}
		return g.local, nil
	case "ls-remote":
		if g.remote == "" {
			return "", nil
		}
		return g.remote + "\t" + args[2], nil
	case "merge-base":
		if g.ancestor {
			return "", nil
		}
		return "", errors.New("exit status 1")
	case "push":
		g.pushes = append(g.pushes, args)
		if g.pushErr != nil {
			return "", g.pushErr
		}
		g.remote = g.local
		return "", nil
	}
	return "", fmt.Errorf("unexpected git %v", args)
}

func (g *fakeGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (g *fakeGit) LookPath(name string) (string, error) { return name, nil }

// forgeLab is a project under pull-request acceptance with one story in
// review, flai serve's forge pass over it, and a fake forge.
type forgeLab struct {
	t       *testing.T
	root    string
	repo    *workitem.Repo
	story   string
	forge   *fakeForge
	on      bool
	token   string
	pass    *forgePass
	journal []hostapi.Entry
	clock   time.Time // a minute later at each look, as flai serve's looks are
}

func newForgeLab(t *testing.T, mode string, git *fakeGit) *forgeLab {
	t.Helper()
	root := t.TempDir()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nacceptance:\n  mode: " + mode + "\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	lab := &forgeLab{t: t, root: root, repo: repo, forge: newFakeForge(), on: true, token: "secret", clock: time.Now()}
	lab.story = lab.inReview("Pushed")
	o := Options{Logger: slog.New(slog.DiscardHandler), Now: func() time.Time { return lab.clock },
		Host: hostapi.Host{
			Enabled: func(action, r string) bool { return lab.on && action == hostapi.ActionForge && r == root },
			Record:  func(e hostapi.Entry) { lab.journal = append(lab.journal, e) },
		},
		ForgeToken: func(name string) ForgeToken {
			return ForgeToken{Env: "S0393_TOKEN", Read: func() (string, error) {
				if lab.token == "" {
					return "", errors.New("the github token is read from the variable S0393_TOKEN, which is not set")
				}
				return lab.token, nil
			}}
		},
		NewForge: func(name, remote string, token forge.TokenFunc) (forge.Forge, error) { return lab.forge, nil },
	}
	if git != nil {
		o.Git = git
	}
	lab.pass = newForgePass(o, Entry{Key: "t", Name: "t", Root: root}, &sync.WaitGroup{})
	return lab
}

// inReview makes a story with two criteria and moves it to review.
func (lab *forgeLab) inReview(title string) string {
	lab.t.Helper()
	epic, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic " + title, Owner: "alex", Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: epic.ID, Owner: "alex", Touches: []string{"docs/" + strings.ToLower(title)}, Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(st.Path)
	if err := os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n- [ ] pushes\n", 1)), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.move(st.ID, workitem.Ready)
	lab.move(st.ID, workitem.InProgress)
	if _, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: st.ID, Owner: "builder", Now: time.Now()}); err != nil {
		lab.t.Fatal(err)
	}
	lab.move(st.ID, workitem.Review)
	return st.ID
}

func (lab *forgeLab) move(id, to string) {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, to, "alex", "the test moves it", time.Now()); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *forgeLab) look() ForgeState {
	lab.t.Helper()
	lab.clock = lab.clock.Add(forgeEvery)
	lab.pass.look(context.Background())
	st, _, err := ReadForgeState(lab.repo, lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	return st
}

// threadsOn are the failure threads flai serve opened on the story, open or not.
func (lab *forgeLab) threadsOn() []*threads.Thread {
	lab.t.Helper()
	all, err := threads.For(lab.repo, lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	var out []*threads.Thread
	for _, th := range all {
		if th.Opener() == forgeAuthor {
			out = append(out, th)
		}
	}
	return out
}

func TestPlanPushPushesPlainlyWithALeaseOrNotAtAll(t *testing.T) {
	yes, no := func() bool { return true }, func() bool { return false }
	for _, c := range []struct {
		name, local, remote, last string
		ancestor                  func() bool
		want                      pushKind
		refused                   string
	}{
		{name: "origin has no branch", local: "b", want: pushPlain, ancestor: no},
		{name: "origin has the local head", local: "b", remote: "b", want: pushNone, ancestor: no},
		{name: "origin's head is in the branch", local: "b", remote: "a", want: pushPlain, ancestor: yes},
		{name: "a sync rebased the branch", local: "c", remote: "a", last: "a", want: pushLease, ancestor: no},
		{name: "origin has commits flai never pushed", local: "c", remote: "a", ancestor: no, refused: "will not force over them"},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := planPush("S-0001", c.local, c.remote, c.last, c.ancestor)
			if c.refused != "" {
				if err == nil || !strings.Contains(err.Error(), c.refused) || !strings.Contains(err.Error(), "git push origin --delete story/S-0001") {
					t.Fatalf("want a refusal saying %q and what clears it, got %v, %v", c.refused, got, err)
				}
				return
			}
			if err != nil || got != c.want {
				t.Fatalf("got %v, %v; want %v", got, err, c.want)
			}
		})
	}
}

// S-0393 criterion 1: with the mode local, or the forge host action off,
// flai serve pushes nothing; the operator's push says why.
func TestNothingIsPushedUnderLocalAcceptanceOrWithTheForgeHostActionOff(t *testing.T) {
	git := &fakeGit{local: "b"}
	lab := newForgeLab(t, "local", git)
	lab.pass.look(context.Background())
	if git.calls != 0 || lab.forge.opens != 0 {
		t.Fatalf("the pass ran git %d times and opened %d requests under local acceptance", git.calls, lab.forge.opens)
	}
	st, err := PushStory(context.Background(), ForgePush{Repo: lab.repo, Story: lab.story, Git: git, Token: ForgeToken{Read: func() (string, error) { return "t", nil }}})
	if err != nil || st.Outcome != ForgeSkipped || !strings.Contains(st.Reason, "acceptance.mode is local") {
		t.Fatalf("the operator's push under local acceptance: %+v, %v", st, err)
	}

	git = &fakeGit{local: "b"}
	lab = newForgeLab(t, "pull-request", git)
	lab.on = false
	lab.pass.look(context.Background())
	if git.calls != 0 || lab.forge.opens != 0 {
		t.Fatalf("the pass ran git %d times and opened %d requests with the host action off", git.calls, lab.forge.opens)
	}
	if _, ok, _ := ReadForgeState(lab.repo, lab.story); ok {
		t.Error("the pass recorded a push with the host action off")
	}
	st, err = PushStory(context.Background(), ForgePush{Repo: lab.repo, Story: lab.story, Git: git, Token: lab.pass.token("github"), NewForge: lab.pass.o.NewForge, HostAction: lab.pass.enabled})
	if err != nil || st.Outcome != ForgeSkipped || !strings.Contains(st.Reason, "flai serve enable forge") {
		t.Fatalf("a push that asks for the host action: %+v, %v", st, err)
	}
	// the operator's word needs no host action
	st, err = PushStory(context.Background(), ForgePush{Repo: lab.repo, Story: lab.story, Git: git, Token: lab.pass.token("github"), NewForge: lab.pass.o.NewForge})
	if err != nil || st.Outcome != ForgePushed || st.Request != RequestOpened {
		t.Fatalf("the operator's push with the host action off: %+v, %v", st, err)
	}
	lab.move(lab.story, workitem.InProgress)
	st, err = PushStory(context.Background(), ForgePush{Repo: lab.repo, Story: lab.story, Git: git, Token: lab.pass.token("github"), NewForge: lab.pass.o.NewForge})
	if err != nil || st.Outcome != ForgeSkipped || !strings.Contains(st.Reason, "not in review") || st.PullRequest != 1 {
		t.Fatalf("a push of a story not in review: %+v, %v", st, err)
	}
}

// S-0393 criterion 5: a failure opens one thread on the story, naming the
// reason and what clears it; the same failure at the next look adds
// nothing, a new reason adds one entry, and a look that succeeds resolves it.
// A story that leaves review has its thread resolved.
func TestAFailedPushOpensOneThreadThatASuccessfulOneResolves(t *testing.T) {
	git := &fakeGit{local: "b"}
	lab := newForgeLab(t, "pull-request", git)
	lab.token = ""
	st := lab.look()
	if st.Outcome != ForgeFailed || !strings.Contains(st.Reason, "S0393_TOKEN") || len(git.pushes) != 0 {
		t.Fatalf("a push without the token: %+v, pushes %v", st, git.pushes)
	}
	ths := lab.threadsOn()
	if len(ths) != 1 || !ths[0].Open() || st.Thread != ths[0].ID || !strings.Contains(ths[0].Body, "S0393_TOKEN") {
		t.Fatalf("threads after the failure: %+v, state %+v", ths, st)
	}
	if strings.Contains(ths[0].Body, "secret") {
		t.Fatal("the thread holds the token")
	}
	if len(lab.journal) != 1 || lab.journal[0].Outcome != "failed" || lab.journal[0].Action != hostapi.ActionForge {
		t.Fatalf("journal after the failure: %+v", lab.journal)
	}

	lab.look()
	ths = lab.threadsOn()
	if len(ths) != 1 || len(ths[0].Entries()) != 1 || len(lab.journal) != 1 {
		t.Fatalf("the same failure again: threads %+v, entries %d, journal %d", ths, len(ths[0].Entries()), len(lab.journal))
	}

	lab.token = "secret"
	git.pushErr = errors.New("remote: Permission to o/r.git denied")
	st = lab.look()
	ths = lab.threadsOn()
	if st.Outcome != ForgeFailed || !strings.Contains(st.Reason, "push access") || len(ths) != 1 || len(ths[0].Entries()) != 2 || len(lab.journal) != 2 {
		t.Fatalf("a refused push: %+v, threads %+v, journal %d", st, ths, len(lab.journal))
	}
	if strings.Contains(strings.Join(git.pushes[0], " "), "--force") {
		t.Fatalf("a first push was forced: %v", git.pushes[0])
	}

	git.pushErr = nil
	st = lab.look()
	ths = lab.threadsOn()
	if st.Outcome != ForgePushed || st.Request != RequestOpened || st.LastPushed != "b" || st.Thread != "" {
		t.Fatalf("a push that succeeds: %+v", st)
	}
	if len(ths) != 1 || ths[0].Open() || !strings.Contains(ths[0].Body, "pull request #1") {
		t.Fatalf("the thread after a push that succeeds: %+v", ths)
	}

	// it fails again, and the story leaves review: its thread is resolved
	lab.forge.down = errors.New("connection refused")
	st = lab.look()
	if st.Outcome != ForgeFailed || !strings.Contains(st.Reason, "tries again") {
		t.Fatalf("the forge down: %+v", st)
	}
	if ths = lab.threadsOn(); len(ths) != 2 || !ths[1].Open() {
		t.Fatalf("a new failure after a success opens a thread: %+v", ths)
	}
	lab.move(lab.story, workitem.InProgress)
	st = lab.look()
	if ths = lab.threadsOn(); ths[1].Open() || st.Thread != "" {
		t.Fatalf("the thread of a story that left review: %+v, state %+v", ths[1], st)
	}
}

// git runs git in dir, failing the test on an error.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false"}, args...)...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// S-0393 criteria 1, 4, and 5, against a scratch bare remote: the first look
// pushes and opens the pull request; a look with nothing new does nothing;
// a new commit pushes and updates the same pull request; a rebase pushes
// with a lease on the head flai last pushed; someone else's push refuses the
// lease, which opens one thread, until their commit is brought in.
func TestFlaiServePushesAStoryInReviewToABareRemote(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	lab := newForgeLab(t, "pull-request", nil)
	root, bare, wt := lab.root, filepath.Join(t.TempDir(), "remote.git"), filepath.Join(t.TempDir(), "wt")
	runGit(t, root, "init", "-q")
	runGit(t, root, "symbolic-ref", "HEAD", "refs/heads/main")
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", "project")
	runGit(t, root, "init", "-q", "--bare", bare)
	runGit(t, root, "remote", "add", "origin", bare)
	branch := "story/" + lab.story
	runGit(t, root, "worktree", "add", "-q", "-b", branch, wt, "main")
	commit := func(dir, file string) string {
		if err := os.WriteFile(filepath.Join(dir, file), []byte(file), 0o644); err != nil {
			t.Fatal(err)
		}
		runGit(t, dir, "add", file)
		runGit(t, dir, "commit", "-q", "-m", file)
		return runGit(t, dir, "rev-parse", "HEAD")
	}
	onOrigin := func() string {
		out := runGit(t, root, "ls-remote", "origin", "refs/heads/"+branch)
		sha, _, _ := strings.Cut(out, "\t")
		return sha
	}

	first := commit(wt, "one.md")
	st := lab.look()
	pr, ok := lab.forge.pr(branch)
	if st.Outcome != ForgePushed || onOrigin() != first || !ok || pr.Number != 1 || pr.Base != "main" || pr.Title != "["+lab.story+"] Pushed" || !strings.Contains(pr.Body, "- [ ] works") {
		t.Fatalf("the first look: %+v, origin %s, request %+v", st, onOrigin(), pr)
	}

	st = lab.look()
	if st.Outcome != ForgeUpToDate || st.Request != RequestUnchanged || lab.forge.opens != 1 {
		t.Fatalf("a look with nothing new: %+v, opens %d", st, lab.forge.opens)
	}

	// a new commit, and a criterion ticked: pushed, and the same request updated
	second := commit(wt, "two.md")
	it, _ := lab.repo.Get(lab.story)
	it.Body = strings.Replace(it.Body, "- [ ] works", "- [x] works", 1)
	if err := lab.repo.Save(it); err != nil {
		t.Fatal(err)
	}
	st = lab.look()
	pr, _ = lab.forge.pr(branch)
	if st.Outcome != ForgePushed || st.Forced || onOrigin() != second || st.Request != RequestUpdated || pr.Number != 1 || !strings.Contains(pr.Body, "- [x] works") {
		t.Fatalf("a new commit: %+v, origin %s, request %+v", st, onOrigin(), pr)
	}

	// a sync rebases the branch: pushed with the lease
	runGit(t, wt, "commit", "-q", "--amend", "-m", "two, rebased")
	rebased := runGit(t, wt, "rev-parse", "HEAD")
	st = lab.look()
	if st.Outcome != ForgePushed || !st.Forced || onOrigin() != rebased || st.LastPushed != rebased {
		t.Fatalf("a rebased branch: %+v, origin %s", st, onOrigin())
	}

	// someone else pushes to the branch, and a sync rebases it again: the lease refuses
	other := filepath.Join(t.TempDir(), "other")
	runGit(t, root, "clone", "-q", "--branch", branch, bare, other)
	theirs := commit(other, "theirs.md")
	runGit(t, other, "push", "-q", "origin", branch)
	runGit(t, wt, "commit", "-q", "--amend", "-m", "two, rebased again")
	st = lab.look()
	ths := lab.threadsOn()
	if st.Outcome != ForgeFailed || onOrigin() != theirs || !strings.Contains(st.Reason, "someone pushed") || len(ths) != 1 || !ths[0].Open() {
		t.Fatalf("a refused lease: %+v, origin %s, threads %+v", st, onOrigin(), ths)
	}
	st = lab.look()
	if ths = lab.threadsOn(); len(ths) != 1 || len(ths[0].Entries()) != 1 || onOrigin() != theirs {
		t.Fatalf("the refused lease again: threads %+v, origin %s", ths, onOrigin())
	}

	// their commit is brought into the branch: a plain push, and the thread resolved
	runGit(t, wt, "pull", "-q", "--no-rebase", "--no-edit", "origin", branch)
	merged := runGit(t, wt, "rev-parse", "HEAD")
	st = lab.look()
	if ths = lab.threadsOn(); st.Outcome != ForgePushed || st.Forced || onOrigin() != merged || ths[0].Open() {
		t.Fatalf("their commit brought in: %+v, origin %s, threads %+v", st, onOrigin(), ths)
	}
	if lab.forge.opens != 2 {
		t.Errorf("the request was opened or updated %d times, want 2", lab.forge.opens)
	}
}

// A read keeps the approvals the pull request's base requires in the status
// file, and a later read keeps what the forge then answers.
func TestAReadKeepsTheApprovalsTheBaseRequires(t *testing.T) {
	lab, _ := newStatusLab(t)
	lab.forge.approvals = 2
	lab.look()
	if st := lab.status(); st.RequiredApprovals != 2 {
		t.Fatalf("the kept status: required approvals %d, want 2", st.RequiredApprovals)
	}
	data, err := os.ReadFile(forge.StatusPath(lab.repo.CacheDir(), lab.story))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"required_approvals": 2`) {
		t.Errorf("the status file does not keep required_approvals:\n%s", data)
	}
	lab.forge.approvals = 0
	lab.look()
	if st := lab.status(); st.RequiredApprovals != 0 {
		t.Errorf("after the rule is lifted: required approvals %d, want 0", st.RequiredApprovals)
	}
}
