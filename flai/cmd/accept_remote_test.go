package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// gitRun runs git in dir with a fixed identity and answers its trimmed
// output, for a fake forge that may run off the test's goroutine.
func gitRun(dir string, args ...string) (string, error) {
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %v: %w %s", args, err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// mergeOnOrigin merges branch into main on the bare repository remote by
// method, as the forge merges a pull request numbered n, and answers main's
// new head.
func mergeOnOrigin(remote, branch string, method forge.MergeMethod, n int) (string, error) {
	clone, err := os.MkdirTemp("", "forge-merge-")
	if err != nil {
		return "", err
	}
	defer func() { _ = os.RemoveAll(clone) }()
	steps := [][]string{{"clone", "-q", remote, clone}}
	switch method {
	case forge.MergeSquash:
		steps = append(steps, []string{"merge", "-q", "--squash", "origin/" + branch}, []string{"commit", "-q", "-m", fmt.Sprintf("Ship it (#%d)", n)})
	case forge.MergeRebase:
		steps = append(steps, []string{"cherry-pick", "origin/main..origin/" + branch})
	default:
		steps = append(steps, []string{"merge", "-q", "--no-ff", "-m", fmt.Sprintf("Merge pull request #%d from %s", n, branch), "origin/" + branch})
	}
	steps = append(steps, []string{"push", "-q", "origin", "HEAD:main"})
	for i, s := range steps {
		dir := clone
		if i == 0 {
			dir = filepath.Dir(clone)
		}
		if _, err := gitRun(dir, s...); err != nil {
			return "", err
		}
	}
	return gitRun(clone, "rev-parse", "HEAD")
}

// mergeForge is a fake forge holding S-0001's one pull request, which it
// merges on the scratch origin by the method asked, or refuses with refuse.
type mergeForge struct {
	mu     sync.Mutex
	remote string
	pr     forge.PullRequest
	refuse error
	merges []forge.MergeMethod
}

// Name is the fake's forge name.
func (f *mergeForge) Name() string { return "github" }

// Repository is the fake's repository.
func (f *mergeForge) Repository() string { return "o/r" }

// Access says the token reads and pushes.
func (f *mergeForge) Access(context.Context) (forge.Access, error) {
	return forge.Access{Read: true, Push: true}, nil
}

// FindPullRequest is the pull request while it is open from head.
func (f *mergeForge) FindPullRequest(_ context.Context, head string) (forge.PullRequest, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.pr, f.pr.Head == head && f.pr.State == forge.PullRequestOpen, nil
}

// PullRequest is the pull request by its number, in any state.
func (f *mergeForge) PullRequest(_ context.Context, n int) (forge.PullRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if n != f.pr.Number {
		return forge.PullRequest{}, forge.ErrNotFound
	}
	return f.pr, nil
}

// OpenPullRequest is not asked of a merge.
func (f *mergeForge) OpenPullRequest(context.Context, forge.PullRequestInput) (forge.PullRequest, error) {
	return forge.PullRequest{}, errors.New("a merge opens no pull request")
}

// Checks are none.
func (f *mergeForge) Checks(context.Context, forge.PullRequest) ([]forge.Check, error) {
	return nil, nil
}

// Reviews are none.
func (f *mergeForge) Reviews(context.Context, forge.PullRequest) ([]forge.Review, error) {
	return nil, nil
}

// RequiredApprovals is none.
func (f *mergeForge) RequiredApprovals(context.Context, string) (int, error) { return 0, nil }

// Merge records the method and merges the pull request on origin, or refuses.
func (f *mergeForge) Merge(_ context.Context, pr forge.PullRequest, method forge.MergeMethod) (forge.Merged, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.merges = append(f.merges, method)
	if f.refuse != nil {
		return forge.Merged{}, f.refuse
	}
	sha, err := mergeOnOrigin(f.remote, pr.Head, method, pr.Number)
	if err != nil {
		return forge.Merged{}, err
	}
	f.pr.State = forge.PullRequestMerged
	return forge.Merged{SHA: sha}, nil
}

// asked are the merges asked of the forge, by method.
func (f *mergeForge) asked() []forge.MergeMethod {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.merges)
}

// mergedBefore merges the pull request on origin as a reviewer would on the
// forge, before flai accept.
func (f *mergeForge) mergedBefore(t *testing.T) {
	t.Helper()
	if _, err := mergeOnOrigin(f.remote, f.pr.Head, forge.MergeMerge, f.pr.Number); err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pr.State = forge.PullRequestMerged
}

// remoteProject is forgeGatedStoryInReview under pull-request acceptance by
// method ("" for the default), with a scratch bare origin holding main and
// the story's branch, the pull request #7 passing as flai serve kept it, the
// forge host action on, the token's variable set, and flai accept pointed at
// a fake forge. It returns the main checkout, the origin, the branch's head,
// and the fake forge.
func remoteProject(t *testing.T, method string) (root, remote, head string, f *mergeForge) {
	t.Helper()
	root, head = forgeGatedStoryInReview(t, "pull-request")
	if method != "" {
		path := filepath.Join(root, "system-flow.yaml")
		m, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(m, "  merge_method: "+method+"\n"...), 0o644); err != nil {
			t.Fatal(err)
		}
		gitIn(t, root, "commit", "-q", "-am", "chore: merge by "+method)
	}
	remote = filepath.Join(t.TempDir(), "origin.git")
	gitIn(t, filepath.Dir(remote), "init", "-q", "--bare", "-b", "main", remote)
	gitIn(t, root, "remote", "add", "origin", remote)
	gitIn(t, root, "push", "-q", "origin", "main", "story/S-0001")
	keepStatus(t, root, passingStatus(head))
	t.Setenv("FLAI_AGENT", "")
	t.Setenv("GITHUB_TOKEN", "accept-token")
	if _, errOut, code := runIn(t, root, "serve", "enable", "forge"); code != 0 {
		t.Fatalf("serve enable forge: %s", errOut)
	}
	f = &mergeForge{remote: remote, pr: forge.PullRequest{Number: 7, URL: forgePRURL, Head: "story/S-0001", Base: "main", HeadSHA: head, State: forge.PullRequestOpen}}
	old, oldWait := acceptNewForge, acceptMergeWait
	acceptNewForge = func(_, _ string, token forge.TokenFunc) (forge.Forge, error) {
		if _, err := token(); err != nil {
			return nil, err
		}
		return f, nil
	}
	acceptMergeWait = serve.MergeWait{Take: 5 * time.Second, Answer: 20 * time.Second, Poll: 10 * time.Millisecond}
	t.Cleanup(func() { acceptNewForge, acceptMergeWait = old, oldWait })
	return root, remote, head, f
}

// acceptJSON runs flai accept S-0001 --by alex --json in root, failing the
// test unless it succeeds, and answers the acceptance.
func acceptJSON(t *testing.T, root string, extra ...string) accepted {
	t.Helper()
	out, errOut, code := runIn(t, root, append([]string{"accept", "S-0001", "--by", "alex", "--json"}, extra...)...)
	if code != 0 {
		t.Fatalf("flai accept: %d %s\n%s", code, errOut, out)
	}
	var res accepted
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	return res
}

// assertAcceptedThroughRemote fails unless S-0001 is done and archived with
// the acceptance commit, the main branch holds origin's main and the story's
// change, and the story's branch is gone from origin and from here with its
// worktree.
func assertAcceptedThroughRemote(t *testing.T, root, remote string) {
	t.Helper()
	if _, err := gitRun(root, "merge-base", "--is-ancestor", "origin/main", "main"); err != nil {
		t.Errorf("the main branch must hold origin's main: %v", err)
	}
	if !strings.Contains(gitIn(t, root, "show", "main:cli/feature.go"), "package cli") {
		t.Error("the main branch must hold the story's change, merged on the forge")
	}
	if got := gitIn(t, root, "log", "-1", "--format=%s"); got != "chore: [S-0001] accept and archive, with E-0001" {
		t.Errorf("the acceptance commit's subject is %q", got)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("the acceptance commits everything: %s", st)
	}
	if gitIn(t, remote, "branch", "--list", "story/S-0001") != "" {
		t.Error("the story's branch is deleted on origin")
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") != "" {
		t.Error("the story's local branch is removed")
	}
	if exists(filepath.Join(root, ".flai-cache", "worktrees", "S-0001")) {
		t.Error("the story's worktree is removed")
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "wip", "archive", "*", "S-0001-*.md")); len(matches) == 0 {
		if matches, _ = filepath.Glob(filepath.Join(root, "wip", "archive", "*", "*", "S-0001-*.md")); len(matches) == 0 {
			t.Error("the story is archived")
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(root, ".flai-cache", "accept-main-*")); len(matches) != 0 {
		t.Errorf("the temporary worktree is removed: %v", matches)
	}
}

// assertNothingChangedHere fails unless main, the story, its branch, and its
// worktree are as they were, and origin still has the branch.
func assertNothingChangedHere(t *testing.T, root, remote, mainBefore string) {
	t.Helper()
	assertNothingAccepted(t, root, mainBefore)
	if !exists(filepath.Join(root, ".flai-cache", "worktrees", "S-0001")) {
		t.Error("a refused acceptance keeps the story's worktree")
	}
	if gitIn(t, remote, "branch", "--list", "story/S-0001") == "" {
		t.Error("a refused acceptance keeps the story's branch on origin")
	}
}

// S-0397 criterion 1: under pull-request acceptance flai accept merges the
// story's pull request on the forge by acceptance.merge_method, deletes the
// branch on origin, brings main up to origin's, removes the local branch and
// worktree, and moves to done, archives, and commits.
func TestAcceptMergesThePullRequestThroughTheRemoteByEachMethod(t *testing.T) {
	for _, method := range []forge.MergeMethod{forge.MergeMerge, forge.MergeSquash, forge.MergeRebase} {
		t.Run(string(method), func(t *testing.T) {
			root, remote, head, f := remoteProject(t, string(method))
			res := acceptJSON(t, root)
			if got := f.asked(); !slices.Equal(got, []forge.MergeMethod{method}) {
				t.Errorf("the forge was asked %v, want one merge by %s", got, method)
			}
			r := res.Remote
			originMain := gitIn(t, remote, "rev-parse", "main")
			if r == nil || r.Outcome != string(serve.MergeMerged) || r.Method != method || r.SHA != originMain || r.Head != head || r.PullRequest != 7 || r.ThroughServe || !r.BranchDeleted {
				t.Errorf("remote = %+v, want merged by %s at %s", r, method, originMain)
			}
			if res.Status != workitem.Done || !res.Merged || res.Archived == 0 {
				t.Errorf("the acceptance: %+v", res.Acceptance)
			}
			assertAcceptedThroughRemote(t, root, remote)
		})
	}
}

// S-0397 criterion 3: a pull request merged on the forge already, by a
// reviewer, is accepted with no second merge, whether flai serve read it
// merged or still read it open.
func TestAcceptOfAPullRequestMergedAlreadyMakesNoSecondMerge(t *testing.T) {
	for name, keptMerged := range map[string]bool{"kept as merged": true, "kept as open": false} {
		t.Run(name, func(t *testing.T) {
			root, remote, head, f := remoteProject(t, "")
			f.mergedBefore(t)
			if keptMerged {
				st := passingStatus(head)
				st.State, st.Checks[0].State = forge.PullRequestMerged, forge.CheckPending
				keepStatus(t, root, st)
			}
			text, _, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run")
			if code != 0 {
				t.Fatalf("dry run: %d", code)
			}
			if keptMerged && !strings.Contains(text, "pull request #7 is merged on the forge already: would accept S-0001 with no second merge, fetch main") {
				t.Errorf("the dry run says the request is merged already:\n%s", text)
			}
			res := acceptJSON(t, root)
			if got := f.asked(); len(got) != 0 {
				t.Errorf("a pull request merged already is not merged again: %v", got)
			}
			if res.Remote == nil || res.Remote.Outcome != string(serve.MergeAlready) {
				t.Errorf("remote = %+v, want already-merged", res.Remote)
			}
			assertAcceptedThroughRemote(t, root, remote)
		})
	}
}

// S-0397 criterion 4: a merge the forge refuses stops the acceptance before
// anything changes here, naming the forge's reason.
func TestAMergeTheForgeRefusesChangesNothingHere(t *testing.T) {
	root, remote, _, f := remoteProject(t, "")
	f.refuse = fmt.Errorf("merge pull request #7 by merge in o/r on GitHub: %w: Merge conflict", forge.ErrMergeRefused)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, "Merge conflict") || !strings.Contains(errOut, "nothing changed") {
		t.Errorf("the refusal names the forge's reason: %d %s", code, errOut)
	}
	if got := f.asked(); len(got) != 1 {
		t.Errorf("the forge was asked once: %v", got)
	}
	assertNothingChangedHere(t, root, remote, mainBefore)
}

// S-0397 criterion 4: unpublished commits on main outside wip/ stop the
// acceptance before the merge, naming each commit.
func TestUnpublishedCommitsOutsideWipStopTheAcceptanceBeforeTheMerge(t *testing.T) {
	root, remote, _, f := remoteProject(t, "")
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("# Guide\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "docs/guide.md")
	gitIn(t, root, "commit", "-q", "-m", "docs: a guide no pull request carried")
	commit := gitIn(t, root, "log", "-1", "--format=%h")
	mainBefore := gitIn(t, root, "rev-parse", "main")
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, commit+" docs: a guide no pull request carried") || !strings.Contains(errOut, "flai release --pending") {
		t.Errorf("the refusal names the commit %s and what to do: %d %s", commit, code, errOut)
	}
	if got := f.asked(); len(got) != 0 {
		t.Errorf("nothing is merged: %v", got)
	}
	assertNothingChangedHere(t, root, remote, mainBefore)
}

// ADR-0170: main's unpublished commits under wip/ are rebased onto origin's
// main, and the main checkout's uncommitted edits survive the move.
func TestUnpublishedWipCommitsAreRebasedOntoOriginsMain(t *testing.T) {
	root, remote, _, _ := remoteProject(t, "squash")
	notes := filepath.Join(root, "wip", "notes.md")
	if err := os.WriteFile(notes, []byte("# Notes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "wip/notes.md")
	gitIn(t, root, "commit", "-q", "-m", "chore: notes on the board")
	unpublished := gitIn(t, root, "rev-parse", "main")
	if err := os.WriteFile(notes, []byte("# Notes\n\nan edit not yet committed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	acceptJSON(t, root)
	assertAcceptedThroughRemote(t, root, remote)
	rebased := gitIn(t, root, "log", "--format=%H", "--grep", "chore: notes on the board", "main")
	if rebased == "" || rebased == unpublished {
		t.Errorf("the unpublished wip commit is rebased onto origin's main: %q (was %s)", rebased, unpublished)
	}
	if _, err := gitRun(root, "merge-base", "--is-ancestor", gitIn(t, remote, "rev-parse", "main"), rebased); err != nil {
		t.Errorf("the rebased commit sits on origin's main: %v", err)
	}
	if b, _ := os.ReadFile(notes); !strings.Contains(string(b), "an edit not yet committed") {
		t.Errorf("the main checkout's uncommitted edit survives: %s", b)
	}
}

// S-0397 criterion 2: a merge flai accept would make itself is refused,
// naming the host action, while the forge host action is off.
func TestAcceptRefusesWhileTheForgeHostActionIsOff(t *testing.T) {
	root, remote, _, f := remoteProject(t, "")
	if _, errOut, code := runIn(t, root, "serve", "disable", "forge"); code != 0 {
		t.Fatal(errOut)
	}
	mainBefore := gitIn(t, root, "rev-parse", "main")
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, "the forge host action is off") || !strings.Contains(errOut, "flai serve enable forge") {
		t.Errorf("the refusal names the host action: %d %s", code, errOut)
	}
	if got := f.asked(); len(got) != 0 {
		t.Errorf("nothing is merged: %v", got)
	}
	assertNothingChangedHere(t, root, remote, mainBefore)
}

// fakeServe answers the requests to merge S-0001's pull request in the
// project at root as flai serve's merge pass does, with a token of its own
// and the fake forge, or with refused when it is set; it records each
// request, and stops when the test ends.
func fakeServe(t *testing.T, root string, f *mergeForge, refused string) *[]serve.MergeRequest {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var seen []serve.MergeRequest
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		request, taken := serve.MergeRequestPath(repo, "S-0001"), serve.MergeTakenPath(repo, "S-0001")
		for ctx.Err() == nil {
			time.Sleep(5 * time.Millisecond)
			if os.Rename(request, taken) != nil {
				continue
			}
			var req serve.MergeRequest
			data, _ := os.ReadFile(taken)
			_ = json.Unmarshal(data, &req)
			mu.Lock()
			seen = append(seen, req)
			mu.Unlock()
			a := serve.MergeAnswer{Story: req.Story, Outcome: serve.MergeRefused, Reason: refused}
			if refused == "" {
				a = serve.MergeStory(ctx, serve.ForgeMerge{Repo: repo, Story: req.Story, Head: req.Head, Method: req.Method, Git: execx.System{},
					Token:    serve.ForgeToken{Env: "GITHUB_TOKEN", Read: func() (string, error) { return "serve-token", nil }},
					NewForge: func(string, string, forge.TokenFunc) (forge.Forge, error) { return f, nil }})
			}
			out, _ := json.Marshal(a)
			_ = os.WriteFile(serve.MergeAnswerPath(repo, "S-0001"), out, 0o644)
		}
	}()
	t.Cleanup(func() {
		cancel()
		<-done
	})
	return &seen
}

// S-0397 criterion 2: with the token's variable unset, as for every agent,
// flai accept asks flai serve to merge, at the head the gate passed at and by
// the project's method, and waits for its answer; flai serve's refusal stops
// it with nothing changed here; and the dashboard's flai accept never hands
// the merge over.
func TestWithoutTheTokenFlaiAcceptHandsTheMergeToFlaiServe(t *testing.T) {
	t.Run("answered merged", func(t *testing.T) {
		root, remote, head, f := remoteProject(t, "rebase")
		t.Setenv("GITHUB_TOKEN", "")
		seen := fakeServe(t, root, f, "")
		if text, _, _ := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run"); !strings.Contains(text, "would merge pull request #7 through the remote by rebase through flai serve, then fetch main") {
			t.Errorf("the dry run says flai serve would merge:\n%s", text)
		}
		res := acceptJSON(t, root)
		if len(*seen) != 1 || (*seen)[0].Head != head || (*seen)[0].Method != forge.MergeRebase || (*seen)[0].By != "alex" {
			t.Errorf("flai serve was asked %+v, want one merge by rebase at %s", *seen, head)
		}
		if got := f.asked(); !slices.Equal(got, []forge.MergeMethod{forge.MergeRebase}) {
			t.Errorf("the forge was asked %v", got)
		}
		if res.Remote == nil || !res.Remote.ThroughServe || res.Remote.Outcome != string(serve.MergeMerged) {
			t.Errorf("remote = %+v, want merged through flai serve", res.Remote)
		}
		assertAcceptedThroughRemote(t, root, remote)
	})
	t.Run("answered refused", func(t *testing.T) {
		root, remote, _, f := remoteProject(t, "")
		t.Setenv("GITHUB_TOKEN", "")
		fakeServe(t, root, f, "the forge host action is off for this project, so flai serve merges nothing: turn it on with flai serve enable forge")
		mainBefore := gitIn(t, root, "rev-parse", "main")
		_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
		if code == 0 || !strings.Contains(errOut, "the forge host action is off") {
			t.Errorf("flai serve's refusal stops the acceptance: %d %s", code, errOut)
		}
		assertNothingChangedHere(t, root, remote, mainBefore)
	})
	t.Run("from the dashboard", func(t *testing.T) {
		root, remote, _, f := remoteProject(t, "")
		t.Setenv("GITHUB_TOKEN", "")
		t.Setenv("FLAI_AGENT", dashboardAgent)
		mainBefore := gitIn(t, root, "rev-parse", "main")
		_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
		if code == 0 || !strings.Contains(errOut, "from the dashboard") || !strings.Contains(errOut, "GITHUB_TOKEN") {
			t.Errorf("the dashboard's acceptance without the token is refused naming the variable: %d %s", code, errOut)
		}
		if exists(filepath.Join(root, ".flai-cache", "forge", "merge", "S-0001.request.json")) || len(f.asked()) != 0 {
			t.Error("the dashboard's acceptance hands nothing to flai serve and merges nothing")
		}
		assertNothingChangedHere(t, root, remote, mainBefore)
	})
}

// S-0397 criterion 5: the dry run says it would merge through the remote and
// by which method, in its text and its JSON, and calls nothing.
func TestTheDryRunSaysItWouldMergeThroughTheRemote(t *testing.T) {
	root, remote, head, f := remoteProject(t, "squash")
	mainBefore := gitIn(t, root, "rev-parse", "main")
	text, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run")
	if code != 0 {
		t.Fatalf("dry run: %d %s", code, errOut)
	}
	want := "would merge pull request #7 through the remote by squash, then fetch main and bring the main branch up to origin's, and remove story/S-0001 and its worktree\n"
	if !strings.Contains(text, want) || strings.Contains(text, "into the main branch") || strings.Contains(text, "through flai serve") {
		t.Errorf("the dry run's text:\n%s\nwant %q", text, want)
	}
	res := acceptJSON(t, root, "--dry-run")
	if r := res.Remote; r == nil || r.PullRequest != 7 || r.Method != forge.MergeSquash || r.Head != head || r.AlreadyMerged || r.ThroughServe || r.Outcome != "" {
		t.Errorf("the dry run's remote: %+v", r)
	}
	if len(f.asked()) != 0 {
		t.Error("the dry run merges nothing")
	}
	assertNothingChangedHere(t, root, remote, mainBefore)
}

// ADR-0170: the acceptance commit's subject and the archive are the same
// under pull-request acceptance as under local acceptance.
func TestAnAcceptanceThroughTheRemoteCommitsAsALocalOneDoes(t *testing.T) {
	committed := func(root string) (string, []string) {
		var archive []string
		for _, p := range strings.Split(gitIn(t, root, "show", "--name-only", "--format=", "HEAD"), "\n") {
			if strings.HasPrefix(p, "wip/") {
				archive = append(archive, p)
			}
		}
		slices.Sort(archive)
		return gitIn(t, root, "log", "-1", "--format=%s"), archive
	}
	local, _ := forgeGatedStoryInReview(t, "local")
	if _, errOut, code := runIn(t, local, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("the local acceptance: %s", errOut)
	}
	wantSubject, wantArchive := committed(local)

	root, _, _, _ := remoteProject(t, "")
	acceptJSON(t, root)
	subject, archive := committed(root)
	if subject != wantSubject || !slices.Equal(archive, wantArchive) {
		t.Errorf("through the remote: %q %v\nlocally: %q %v", subject, archive, wantSubject, wantArchive)
	}
}
