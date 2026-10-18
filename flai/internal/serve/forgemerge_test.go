package serve

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// newMergeLab is a status lab whose story's branch flai serve pushed and
// whose pull request #1 it opened and read, at headA, ready to merge.
func newMergeLab(t *testing.T) (*forgeLab, *fakeGit) {
	t.Helper()
	lab, git := newStatusLab(t)
	if st := lab.look(); st.Outcome != ForgePushed || st.PullRequest != 1 {
		t.Fatalf("the look before the merge: %+v", st)
	}
	return lab, git
}

// merge merges the lab's story's pull request as flai serve does, with its
// token, at head, by method.
func (lab *forgeLab) merge(git *fakeGit, head string, method forge.MergeMethod) MergeAnswer {
	lab.t.Helper()
	m := ForgeMerge{Repo: lab.repo, Story: lab.story, Head: head, Method: method, Token: lab.pass.token("github"), NewForge: lab.pass.o.NewForge, Now: lab.pass.o.Now}
	if git != nil {
		m.Git = git
	}
	return MergeStory(context.Background(), m)
}

// deletes are the pushes that deleted a branch on origin.
func deletes(git *fakeGit) [][]string {
	git.mu.Lock()
	defer git.mu.Unlock()
	var out [][]string
	for _, p := range git.pushes {
		if slices.Contains(p, "--delete") {
			out = append(out, p)
		}
	}
	return out
}

// S-0397 criterion 1: the pull request is merged by the method asked for,
// or by acceptance.merge_method when none is, and story/S-nnnn is deleted
// on origin; the kept status records it merged.
func TestMergeStoryMergesByEachMethodAndDeletesTheBranch(t *testing.T) {
	for _, method := range forge.MergeMethods() {
		t.Run(string(method), func(t *testing.T) {
			lab, git := newMergeLab(t)
			a := lab.merge(git, headA, method)
			if a.Outcome != MergeMerged || a.SHA != "merged-by-"+string(method) || a.Method != method || a.PullRequest != 1 || a.URL == "" || a.Reason != "" {
				t.Fatalf("the merge: %+v", a)
			}
			if !slices.Equal(lab.forge.merges, []forge.MergeMethod{method}) {
				t.Errorf("the forge was asked to merge by %v, want %s once", lab.forge.merges, method)
			}
			del := deletes(git)
			if !a.BranchDeleted || a.BranchReason != "" || len(del) != 1 || !slices.Contains(del[0], "refs/heads/"+lab.branch()) {
				t.Errorf("the branch: deleted %v (%s), pushes %v", a.BranchDeleted, a.BranchReason, del)
			}
			if st := lab.status(); st.State != forge.PullRequestMerged {
				t.Errorf("the kept status after the merge is %q, want merged", st.State)
			}
		})
	}

	t.Run("the manifest's method", func(t *testing.T) {
		lab, git := newMergeLab(t)
		path := filepath.Join(lab.root, "system-flow.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, "  merge_method: squash\n"...), 0o644); err != nil {
			t.Fatal(err)
		}
		if lab.repo, err = workitem.Open(lab.root); err != nil {
			t.Fatal(err)
		}
		if a := lab.merge(git, "", ""); a.Outcome != MergeMerged || a.Method != forge.MergeSquash || !slices.Equal(lab.forge.merges, []forge.MergeMethod{forge.MergeSquash}) {
			t.Fatalf("a merge with no method: %+v, merges %v; want squash, the manifest's", a, lab.forge.merges)
		}
	})
}

// S-0397 criterion 3: a pull request a reviewer merged on the forge is
// answered as merged already, with no second merge, and its branch deleted;
// where origin has no branch left, nothing is deleted and the answer says why.
func TestMergeStoryFindsAPullRequestMergedAlready(t *testing.T) {
	lab, git := newMergeLab(t)
	lab.forge.setState(lab.branch(), forge.PullRequestMerged)
	a := lab.merge(git, headA, forge.MergeMerge)
	if a.Outcome != MergeAlready || !a.Merged() || a.PullRequest != 1 || len(lab.forge.merges) != 0 {
		t.Fatalf("a pull request merged already: %+v, merges %v", a, lab.forge.merges)
	}
	if !a.BranchDeleted || len(deletes(git)) != 1 {
		t.Errorf("the branch of a pull request merged already: deleted %v (%s)", a.BranchDeleted, a.BranchReason)
	}
	a = lab.merge(git, headA, forge.MergeMerge)
	if a.Outcome != MergeAlready || a.BranchDeleted || !strings.Contains(a.BranchReason, "origin has no "+lab.branch()) || len(deletes(git)) != 1 {
		t.Errorf("a second acceptance: %+v, deletes %v", a, deletes(git))
	}
}

// S-0397 criterion 4: a merge the forge refuses, for a conflict or a rule,
// is answered refused with the forge's own reason, and changes nothing: the
// branch is not deleted and the pull request is not recorded merged.
func TestMergeStoryRefusesWithTheForgesReason(t *testing.T) {
	for _, c := range []struct{ name, said string }{
		{"a conflict", "PUT /repos/o/r/pulls/1/merge answered 405 Method Not Allowed: Pull Request is not mergeable"},
		{"a branch out of date", "PUT /repos/o/r/pulls/1/merge answered 405 Method Not Allowed: Head branch is out of date"},
		{"a rule", `PUT /repos/o/r/pulls/1/merge answered 405 Method Not Allowed: Required status check "build" is expected.`},
		{"a moved head", "PUT /repos/o/r/pulls/1/merge answered 409 Conflict: Head branch was modified. Review and try the merge again."},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab, git := newMergeLab(t)
			lab.forge.mergeRefusal = fmt.Errorf("merge pull request #1 by rebase in o/r on GitHub: %w: %s", forge.ErrMergeRefused, c.said)
			a := lab.merge(git, headA, forge.MergeRebase)
			if a.Outcome != MergeRefused || a.Merged() || !strings.Contains(a.Reason, c.said) || !strings.Contains(a.Reason, "accept again") {
				t.Fatalf("a refused merge: %+v", a)
			}
			if a.BranchDeleted || len(deletes(git)) != 0 {
				t.Errorf("a refused merge deleted the branch: %v", deletes(git))
			}
			if st := lab.status(); st.State == forge.PullRequestMerged {
				t.Error("a refused merge was recorded as merged")
			}
			if strings.Contains(a.Reason, lab.token) {
				t.Error("the reason holds the token")
			}
		})
	}
}

// A pull request whose head moved since the gate passed, one closed without
// merging, one never opened, and a project under local acceptance are each
// refused before the forge is asked to merge, saying what clears it.
func TestMergeStoryRefusesBeforeAskingTheForge(t *testing.T) {
	t.Run("a moved head", func(t *testing.T) {
		lab, git := newMergeLab(t)
		a := lab.merge(git, headB, forge.MergeMerge)
		if a.Outcome != MergeRefused || !strings.Contains(a.Reason, short(headA)) || !strings.Contains(a.Reason, short(headB)) || !strings.Contains(a.Reason, "next look") {
			t.Fatalf("a moved head: %+v", a)
		}
		if len(lab.forge.merges) != 0 || len(deletes(git)) != 0 {
			t.Errorf("a moved head was merged %v or its branch deleted %v", lab.forge.merges, deletes(git))
		}
	})
	t.Run("closed without merging", func(t *testing.T) {
		lab, git := newMergeLab(t)
		lab.forge.setState(lab.branch(), forge.PullRequestClosed)
		a := lab.merge(git, headA, forge.MergeMerge)
		if a.Outcome != MergeRefused || !strings.Contains(a.Reason, "closed without merging") || !strings.Contains(a.Reason, "reopen it") || len(lab.forge.merges) != 0 {
			t.Fatalf("a closed pull request: %+v, merges %v", a, lab.forge.merges)
		}
	})
	t.Run("no pull request", func(t *testing.T) {
		lab, git := newStatusLab(t)
		a := lab.merge(git, headA, forge.MergeMerge)
		if a.Outcome != MergeRefused || !strings.Contains(a.Reason, "has no pull request") || len(lab.forge.merges) != 0 {
			t.Fatalf("no pull request: %+v", a)
		}
	})
	t.Run("local acceptance", func(t *testing.T) {
		git := &fakeGit{local: headA}
		lab := newForgeLab(t, "local", git)
		a := lab.merge(git, headA, forge.MergeMerge)
		if a.Outcome != MergeRefused || !strings.Contains(a.Reason, "acceptance.mode is local") || git.calls != 0 {
			t.Fatalf("local acceptance: %+v, git ran %d times", a, git.calls)
		}
	})
}

// A merge that went through is answered merged though git could not delete
// the branch on origin, saying why: the merge stands.
func TestMergeStoryReportsABranchItCouldNotDelete(t *testing.T) {
	lab, git := newMergeLab(t)
	git.mu.Lock()
	git.pushErr = fmt.Errorf("remote: Permission to o/r.git denied")
	git.mu.Unlock()
	a := lab.merge(git, headA, forge.MergeMerge)
	if a.Outcome != MergeMerged || a.BranchDeleted || !strings.Contains(a.BranchReason, "the merge stands") || !strings.Contains(a.BranchReason, "Permission to o/r.git denied") {
		t.Fatalf("a branch git could not delete: %+v", a)
	}
}

// S-0397 criterion 1, against a scratch bare remote: the merge deletes
// story/S-nnnn on origin.
func TestMergeStoryDeletesTheBranchOnABareRemote(t *testing.T) {
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
	branch := lab.branch()
	runGit(t, root, "worktree", "add", "-q", "-b", branch, wt, "main")
	if err := os.WriteFile(filepath.Join(wt, "one.md"), []byte("one"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, wt, "add", "one.md")
	runGit(t, wt, "commit", "-q", "-m", "one")
	head := runGit(t, wt, "rev-parse", "HEAD")
	if st := lab.look(); st.Outcome != ForgePushed || st.PullRequest != 1 {
		t.Fatalf("the push before the merge: %+v", st)
	}
	lab.forge.set(branch, head, nil, nil)
	a := lab.merge(nil, head, forge.MergeMerge)
	if a.Outcome != MergeMerged || !a.BranchDeleted {
		t.Fatalf("the merge: %+v", a)
	}
	if out := runGit(t, root, "ls-remote", "origin", "refs/heads/"+branch); out != "" {
		t.Errorf("origin still has %s: %s", branch, out)
	}
	if out := runGit(t, root, "rev-parse", "--verify", "refs/heads/"+branch); out != head {
		t.Errorf("the local branch is at %q, want %s untouched", out, head)
	}
}

// handOff is what the lab's merge pass logged and journalled, written from
// its goroutine.
type handOff struct {
	logs    syncBuffer
	mu      sync.Mutex
	entries []hostapi.Entry
}

// last is the last entry the pass journalled.
func (h *handOff) last() hostapi.Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.entries) == 0 {
		return hostapi.Entry{}
	}
	return h.entries[len(h.entries)-1]
}

// mergeHandOff is the lab's merge pass running every few milliseconds,
// stopped when the test ends; set the lab's host action before it starts.
func (lab *forgeLab) mergeHandOff() *handOff {
	lab.t.Helper()
	h := &handOff{}
	o := lab.pass.o
	o.Logger = slog.New(slog.NewTextHandler(&h.logs, nil))
	o.Host.Record = func(e hostapi.Entry) {
		h.mu.Lock()
		defer h.mu.Unlock()
		h.entries = append(h.entries, e)
	}
	var wg sync.WaitGroup
	p := newMergePass(o, lab.pass.e, &wg)
	p.every = 5 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	p.start(ctx)
	lab.t.Cleanup(func() {
		cancel()
		wg.Wait()
	})
	return h
}

// shortWait is a wait for flai serve short enough for a test.
var shortWait = MergeWait{Take: 5 * time.Second, Answer: 5 * time.Second, Poll: 5 * time.Millisecond}

// S-0397 criterion 2: flai accept, with no token, hands the merge to flai
// serve, which takes the request, merges with its token, and answers; the
// request, its taking, and its answer leave nothing behind, and the log and
// the journal name the merge but never the token.
func TestRequestMergeIsAnsweredByFlaiServesMergePass(t *testing.T) {
	lab, git := newMergeLab(t)
	h := lab.mergeHandOff()
	a, err := RequestMerge(context.Background(), lab.repo, MergeRequest{Story: lab.story, Head: headA, Method: forge.MergeSquash, By: "alex"}, shortWait)
	if err != nil {
		t.Fatal(err)
	}
	if a.Outcome != MergeMerged || a.Method != forge.MergeSquash || !a.BranchDeleted || len(deletes(git)) != 1 || !slices.Equal(lab.forge.mergesAsked(), []forge.MergeMethod{forge.MergeSquash}) {
		t.Fatalf("the answer: %+v", a)
	}
	for _, path := range []string{MergeRequestPath(lab.repo, lab.story), MergeTakenPath(lab.repo, lab.story), MergeAnswerPath(lab.repo, lab.story)} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("%s is left behind: %v", path, err)
		}
	}
	waitFor(t, "the merge journalled", func() bool { return h.last().Method == "serve.merge" })
	if e := h.last(); e.Outcome != "done" || !strings.Contains(e.Detail, "by squash") {
		t.Errorf("the journal: %+v", e)
	}
	if s := h.logs.String(); !strings.Contains(s, "pull request merged") || !strings.Contains(s, "story="+lab.story) || strings.Contains(s, lab.token) {
		t.Errorf("the log of the merge, which must not hold the token:\n%s", s)
	}
}

// S-0397 criterion 2: with the forge host action off, flai serve takes the
// request and refuses it, naming the host action, and merges nothing.
func TestRequestMergeIsRefusedWithTheForgeHostActionOff(t *testing.T) {
	lab, git := newMergeLab(t)
	lab.on = false
	h := lab.mergeHandOff()
	a, err := RequestMerge(context.Background(), lab.repo, MergeRequest{Story: lab.story, Head: headA}, shortWait)
	if err != nil {
		t.Fatal(err)
	}
	if a.Outcome != MergeRefused || !strings.Contains(a.Reason, "forge host action is off") || !strings.Contains(a.Reason, "flai serve enable forge") {
		t.Fatalf("the answer with the host action off: %+v", a)
	}
	if len(lab.forge.mergesAsked()) != 0 || len(deletes(git)) != 0 {
		t.Errorf("merged %v or deleted %v with the host action off", lab.forge.mergesAsked(), deletes(git))
	}
	waitFor(t, "the merge journalled", func() bool { return h.last().Method == "serve.merge" })
	if e := h.last(); e.Outcome != "disabled" {
		t.Errorf("the journal: %+v", e)
	}
}

// A request no flai serve takes is withdrawn once the wait for the take runs
// out, saying how to start flai serve and turn the host action on.
func TestRequestMergeNotTakenIsWithdrawn(t *testing.T) {
	lab, _ := newStatusLab(t)
	_, err := RequestMerge(context.Background(), lab.repo, MergeRequest{Story: lab.story}, MergeWait{Take: 50 * time.Millisecond, Poll: 5 * time.Millisecond})
	if err == nil || !strings.Contains(err.Error(), "did not take") || !strings.Contains(err.Error(), "flai serve start") || !strings.Contains(err.Error(), "flai serve enable forge") {
		t.Fatalf("a request no one took: %v", err)
	}
	if _, err := os.Stat(MergeRequestPath(lab.repo, lab.story)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
}

// A request flai serve took but did not answer is waited for until the wait
// for the answer runs out, saying that a merge it made is found when the
// story is accepted again.
func TestRequestMergeTakenButNotAnswered(t *testing.T) {
	lab, _ := newStatusLab(t)
	request, taken := MergeRequestPath(lab.repo, lab.story), MergeTakenPath(lab.repo, lab.story)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 200 {
			if os.Rename(request, taken) == nil {
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()
	_, err := RequestMerge(context.Background(), lab.repo, MergeRequest{Story: lab.story}, MergeWait{Take: 5 * time.Second, Answer: 50 * time.Millisecond, Poll: 5 * time.Millisecond})
	<-done
	if err == nil || !strings.Contains(err.Error(), "did not answer") || !strings.Contains(err.Error(), "found, not repeated") {
		t.Fatalf("a request taken and not answered: %v", err)
	}
}

// Asking again removes a stale answer left by a wait that ran out, so that
// it is never taken for the new request's.
func TestRequestMergeRemovesAStaleAnswer(t *testing.T) {
	lab, _ := newStatusLab(t)
	stale := MergeAnswer{Story: lab.story, Outcome: MergeRefused, Reason: "stale"}
	if err := os.MkdirAll(filepath.Dir(MergeAnswerPath(lab.repo, lab.story)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeMergeAnswer(lab.repo, stale); err != nil {
		t.Fatal(err)
	}
	_, err := RequestMerge(context.Background(), lab.repo, MergeRequest{Story: lab.story}, MergeWait{Take: 30 * time.Millisecond, Poll: 5 * time.Millisecond})
	if err == nil || !strings.Contains(err.Error(), "did not take") {
		t.Fatalf("a request with a stale answer beside it: %v, want it not taken", err)
	}
}
