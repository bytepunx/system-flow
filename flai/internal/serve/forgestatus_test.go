package serve

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The project's markdown lint, as the template writes it: threads.New
// refuses a thread whose text it finds fault with.
const projectLint = `default: true
MD013: false
MD022:
  lines_below: 0
MD024:
  siblings_only: true
MD025:
  front_matter_title: ""
MD032: false
MD033: false
MD041: false
MD060: false
`

const (
	headA = "aaaaaaaaaaaa1111111111111111111111111111"
	headB = "bbbbbbbbbbbb2222222222222222222222222222"
	logA  = "https://github.com/o/r/actions/runs/11/job/12"
)

// statusLab is a forge lab whose story's pull request, number 1, is at headA
// and whose project lints its markdown.
func newStatusLab(t *testing.T) (*forgeLab, *fakeGit) {
	t.Helper()
	git := &fakeGit{local: headA}
	lab := newForgeLab(t, "pull-request", git)
	if err := os.WriteFile(filepath.Join(lab.root, ".markdownlint.yaml"), []byte(projectLint), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.forge.set(lab.branch(), headA, nil, nil)
	return lab, git
}

func (lab *forgeLab) branch() string { return "story/" + lab.story }

// status is the story's kept pull request status.
func (lab *forgeLab) status() forge.Status {
	lab.t.Helper()
	st, ok, err := forge.ReadStatus(lab.repo.CacheDir(), lab.story)
	if err != nil || !ok {
		lab.t.Fatalf("the kept status: %+v, %v, %v", st, ok, err)
	}
	return st
}

func (lab *forgeLab) item() *workitem.Item {
	lab.t.Helper()
	it, err := lab.repo.Get(lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	return it
}

// sentBackBy is how many moves from review to in-progress flai made.
func (lab *forgeLab) sentBackBy() int {
	n := 0
	for _, tr := range lab.item().Transitions {
		if tr.To == workitem.InProgress && tr.By == sendBackBy {
			n++
		}
	}
	return n
}

// buildFailed is the required check build failed, with its log at logA.
func buildFailed(summary string) forge.Check {
	return forge.Check{Name: "build", State: forge.CheckFailed, Required: true, URL: logA, Summary: summary}
}

// S-0394 criteria 1 and 4: a required check that passes or is pending, and a
// check that is not required failing, leave the story in review with no
// thread; the kept state holds the required checks with an event for each
// change. A read that fails touches nothing.
func TestAPassingOrPendingCheckLeavesTheStoryInReview(t *testing.T) {
	lab, _ := newStatusLab(t)
	for _, c := range []struct {
		name   string
		checks []forge.Check
		want   forge.CheckState
	}{
		{name: "passing", checks: []forge.Check{{Name: "build", State: forge.CheckPassed, Required: true}}, want: forge.CheckPassed},
		{name: "pending", checks: []forge.Check{{Name: "build", State: forge.CheckPending, Required: true}}, want: forge.CheckPending},
		{name: "not required", checks: []forge.Check{{Name: "build", State: forge.CheckPassed, Required: true}, {Name: "lint", State: forge.CheckFailed, URL: logA}}, want: forge.CheckPassed},
	} {
		lab.forge.set(lab.branch(), headA, c.checks, nil)
		lab.look()
		st := lab.status()
		if len(st.Checks) != 1 || st.Checks[0].Name != "build" || st.Checks[0].State != c.want || st.HeadSHA != headA || st.PullRequest != 1 {
			t.Fatalf("%s: the kept status %+v", c.name, st)
		}
		if ths := lab.threadsOn(); len(ths) != 0 || lab.item().Status != workitem.Review || len(st.SentBack) != 0 {
			t.Fatalf("%s: threads %+v, status %s, sent back %+v", c.name, ths, lab.item().Status, st.SentBack)
		}
	}
	st := lab.status()
	var kinds []string
	for _, e := range st.Events {
		kinds = append(kinds, e.Kind+":"+e.Name+":"+e.To)
	}
	if got, want := strings.Join(kinds, ","), "head::"+headA+",check:build:passed,check:build:pending,check:build:passed"; got != want {
		t.Errorf("events %s, want %s", got, want)
	}

	lab.forge.readErr = errors.New("connection reset")
	lab.forge.set(lab.branch(), headA, []forge.Check{buildFailed("")}, nil)
	lab.look()
	if after := lab.status(); after.ReadAt != st.ReadAt || len(after.Events) != len(st.Events) || len(lab.threadsOn()) != 0 {
		t.Fatalf("a read that failed changed the kept state or opened a thread: %+v", after)
	}
}

// S-0394 criterion 4: a failing required check opens one thread naming it
// and linking its log, records the send back, and moves the story to
// in-progress by flai with the reason. Back in review at the same commit it
// is not sent back again; at a new commit with the same check failing it is,
// with a new thread, and the old one is resolved. A story in progress keeps
// its thread; a cancelled one has it resolved.
func TestAFailingRequiredCheckSendsTheStoryBackOncePerCommit(t *testing.T) {
	lab, git := newStatusLab(t)
	lab.forge.set(lab.branch(), headA, []forge.Check{buildFailed("3 tests failed: see https://example.com/t")}, nil)
	lab.look()
	ths := lab.threadsOn()
	if len(ths) != 1 || !ths[0].Open() {
		t.Fatalf("threads after a failing check: %+v", ths)
	}
	th := ths[0]
	if want := "Pull request #1 of " + lab.story + ": the required check build failed"; th.Title != want {
		t.Errorf("title %q, want %q", th.Title, want)
	}
	for _, want := range []string{"The required check `build` failed at `aaaaaaaaaaaa` ([its log](" + logA + "))", "> 3 tests failed: see `https://example.com/t`", "[pull request #1](https://github.com/o/r/pull/1)", "flai moved " + lab.story + " back to in-progress"} {
		if !strings.Contains(th.Body, want) {
			t.Errorf("the thread does not say %q:\n%s", want, th.Body)
		}
	}
	it := lab.item()
	reason := "the required check build failed on aaaaaaaaaaaa of pull request #1"
	if it.Status != workitem.InProgress || lab.sentBackBy() != 1 || !strings.Contains(it.Body, reason) {
		t.Fatalf("the story after a failing check: %s, moves by flai %d, body:\n%s", it.Status, lab.sentBackBy(), it.Body)
	}
	st := lab.status()
	if len(st.SentBack) != 1 || st.SentBack[0].Key != forge.CheckKey("build", headA) || st.SentBack[0].Thread != th.ID || st.SentBack[0].URL != logA || st.SentBack[0].Commit != headA {
		t.Fatalf("sent back %+v", st.SentBack)
	}
	if b, ok, err := OpenSendBack(lab.repo, lab.story); err != nil || !ok || b.Thread != th.ID {
		t.Fatalf("OpenSendBack = %+v, %v, %v", b, ok, err)
	}
	if n := len(lab.journal); n == 0 || !strings.Contains(lab.journal[n-1].Detail, "sent "+lab.story+" back to in-progress: "+reason) || lab.journal[n-1].Outcome != "done" {
		t.Fatalf("journal %+v", lab.journal)
	}

	// in progress: not read, and its thread stays open
	lab.look()
	if ths = lab.threadsOn(); len(ths) != 1 || !ths[0].Open() {
		t.Fatalf("the thread of a story in progress: %+v", ths)
	}

	// back in review at the same commit, the same check failing: nothing new
	lab.move(lab.story, workitem.Review)
	lab.look()
	if ths = lab.threadsOn(); len(ths) != 1 || !ths[0].Open() || lab.item().Status != workitem.Review || lab.sentBackBy() != 1 {
		t.Fatalf("the same failure on the same commit: threads %+v, status %s", ths, lab.item().Status)
	}

	// a new commit with the same check failing: sent back again, the old thread resolved
	git.local = headB
	lab.forge.set(lab.branch(), headB, []forge.Check{buildFailed("")}, nil)
	lab.look()
	ths = lab.threadsOn()
	if len(ths) != 2 || ths[0].Open() || !ths[1].Open() || !strings.Contains(ths[0].Body, lab.story+" is back in review at bbbbbbbbbbbb") {
		t.Fatalf("threads after the same check failed on a new commit: %+v", ths)
	}
	if lab.item().Status != workitem.InProgress || lab.sentBackBy() != 2 || len(lab.status().SentBack) != 2 {
		t.Fatalf("the story after a new commit failed: %s, moves by flai %d", lab.item().Status, lab.sentBackBy())
	}

	// cancelled and archived at once, as accepting a story archives it: no
	// longer on the board, and its thread resolved all the same
	lab.move(lab.story, workitem.Cancelled)
	items, err := lab.repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := lab.repo.PlanArchive(items, []string{lab.story})
	if err != nil {
		t.Fatal(err)
	}
	if err := lab.repo.Archive(plan); err != nil {
		t.Fatal(err)
	}
	lab.look()
	if ths = lab.threadsOn(); ths[1].Open() || !strings.Contains(ths[1].Body, lab.story+" is cancelled") {
		t.Fatalf("the thread of a cancelled story: %+v", ths[1])
	}
	if _, ok, err := OpenSendBack(lab.repo, lab.story); ok || err != nil {
		t.Fatalf("OpenSendBack of a cancelled story = %v, %v", ok, err)
	}
}

// S-0394 criterion 4: a review that requests changes opens a thread quoting
// its comment, briefly, and linking it, and sends the story back.
func TestARequestedChangeSendsTheStoryBackQuotingTheReview(t *testing.T) {
	lab, _ := newStatusLab(t)
	review := "https://github.com/o/r/pull/1#pullrequestreview-9"
	body := "# Naming\r\n\r\nPlease rename the `--force` flag, see www.example.com.\r\nIt reads as\tdangerous.\r\n\r\n" + strings.Repeat("More words. ", 60)
	lab.forge.set(lab.branch(), headA, nil, []forge.Review{{Author: "alice", State: forge.ReviewChangesRequested, Body: body, URL: review}})
	lab.look()
	ths := lab.threadsOn()
	if len(ths) != 1 || ths[0].Title != "Pull request #1 of "+lab.story+": alice requested changes" {
		t.Fatalf("threads after a requested change: %+v", ths)
	}
	for _, want := range []string{"alice requested changes ([the review](" + review + "))", "> \\# Naming\n>\n> Please rename the `--force` flag, see `www.example.com`. It reads as dangerous.\n>\n> More words.", "…"} {
		if !strings.Contains(ths[0].Body, want) {
			t.Errorf("the thread does not say %q:\n%s", want, ths[0].Body)
		}
	}
	if strings.Count(ths[0].Body, "More words.") >= 60 {
		t.Error("the quote is not trimmed")
	}
	if it := lab.item(); it.Status != workitem.InProgress || !strings.Contains(it.Body, "alice requested changes on pull request #1") {
		t.Fatalf("the story after a requested change: %s\n%s", it.Status, it.Body)
	}
	if b := lab.status().SentBack; len(b) != 1 || b[0].Key != forge.ReviewKey(review) || b[0].Kind != forge.EventReview {
		t.Fatalf("sent back %+v", b)
	}
}

// A reviewer's later approval supersedes their requested change, and a
// comment does not: nothing is sent back.
func TestALaterApprovalSupersedesARequestedChange(t *testing.T) {
	lab, _ := newStatusLab(t)
	lab.forge.set(lab.branch(), headA, nil, []forge.Review{
		{Author: "alice", State: forge.ReviewChangesRequested, Body: "No.", URL: "https://github.com/o/r/pull/1#pullrequestreview-1"},
		{Author: "alice", State: forge.ReviewApproved, URL: "https://github.com/o/r/pull/1#pullrequestreview-2"},
		{Author: "alice", State: forge.ReviewCommented, Body: "Nit.", URL: "https://github.com/o/r/pull/1#pullrequestreview-3"},
	})
	lab.look()
	if ths := lab.threadsOn(); len(ths) != 0 || lab.item().Status != workitem.Review {
		t.Fatalf("an approved request: threads %+v, status %s", ths, lab.item().Status)
	}
	if st := lab.status(); len(st.Reviews) != 3 {
		t.Fatalf("kept reviews %+v", st.Reviews)
	}
}

// A failing check and a requested change in one read open one thread.
func TestACheckAndAReviewInOneReadOpenOneThread(t *testing.T) {
	lab, _ := newStatusLab(t)
	lab.forge.set(lab.branch(), headA, []forge.Check{buildFailed("")},
		[]forge.Review{{Author: "bob", State: forge.ReviewChangesRequested, URL: "https://github.com/o/r/pull/1#pullrequestreview-4"}})
	lab.look()
	ths := lab.threadsOn()
	if len(ths) != 1 || ths[0].Title != "Pull request #1 of "+lab.story+": 1 required check failed and 1 review requested changes" {
		t.Fatalf("threads: %+v", ths)
	}
	if !strings.Contains(ths[0].Body, "bob requested changes ([the review](https://github.com/o/r/pull/1#pullrequestreview-4)), with no comment.") {
		t.Errorf("the thread does not name the review:\n%s", ths[0].Body)
	}
	b := lab.status().SentBack
	if len(b) != 2 || b[0].Thread != ths[0].ID || b[1].Thread != ths[0].ID {
		t.Fatalf("sent back %+v", b)
	}
	if it := lab.item(); lab.sentBackBy() != 1 || !strings.Contains(it.Body, "the required check build failed on aaaaaaaaaaaa of pull request #1; bob requested changes on pull request #1") {
		t.Fatalf("the move: %d by flai\n%s", lab.sentBackBy(), it.Body)
	}
}

// A move that failed after its thread was opened is made again at the next
// look, without a new thread.
func TestAFailedSendBackMoveIsMadeAgainWithoutANewThread(t *testing.T) {
	lab, _ := newStatusLab(t)
	lab.forge.set(lab.branch(), headA, []forge.Check{buildFailed("")}, nil)
	lab.clock = lab.clock.Add(time.Hour) // the send back, well after the story's move to in-progress
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Pull request #1 of " + lab.story + ": the required check build failed", On: lab.story, Author: forgeAuthor, Text: "Sent back.", Now: lab.clock})
	if err != nil {
		t.Fatal(err)
	}
	at := lab.clock.UTC().Format(time.RFC3339)
	kept := forge.Status{Story: lab.story, PullRequest: 1, HeadSHA: headA,
		SentBack: []forge.SendBack{{Key: forge.CheckKey("build", headA), Kind: forge.EventCheck, Name: "build", Commit: headA, URL: logA, Thread: th.ID, At: at}}}
	if err := forge.WriteStatus(lab.repo.CacheDir(), kept); err != nil {
		t.Fatal(err)
	}
	lab.look()
	if ths := lab.threadsOn(); len(ths) != 1 || !ths[0].Open() {
		t.Fatalf("threads after the move again: %+v", ths)
	}
	if it := lab.item(); it.Status != workitem.InProgress || lab.sentBackBy() != 1 || !strings.Contains(it.Body, "the required check build failed on aaaaaaaaaaaa of pull request #1") {
		t.Fatalf("the story after the move again: %s\n%s", it.Status, it.Body)
	}
}

// A send back whose record was lost after its thread was opened, as when the
// kept state could not be written, uses that open thread again rather than
// opening a second one for the same failure.
func TestASendBackNotRecordedUsesItsOpenThreadAgain(t *testing.T) {
	lab, _ := newStatusLab(t)
	lab.forge.set(lab.branch(), headA, []forge.Check{buildFailed("")}, nil)
	lab.clock = lab.clock.Add(time.Hour)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Pull request #1 of " + lab.story + ": the required check build failed", On: lab.story, Author: forgeAuthor, Text: "Sent back.", Now: lab.clock})
	if err != nil {
		t.Fatal(err)
	}
	lab.look()
	if ths := lab.threadsOn(); len(ths) != 1 || !ths[0].Open() {
		t.Fatalf("threads after a send back not recorded: %+v", ths)
	}
	if st := lab.status(); len(st.SentBack) != 1 || st.SentBack[0].Thread != th.ID {
		t.Fatalf("sent back %+v, want it on %s", st.SentBack, th.ID)
	}
	if it := lab.item(); it.Status != workitem.InProgress || lab.sentBackBy() != 1 {
		t.Fatalf("the story after a send back not recorded: %s", it.Status)
	}
}

// A comment the project's lint refuses as a blockquote, here for mixing its
// emphasis styles, is quoted in a code block instead.
func TestACommentTheLintRefusesQuotedIsQuotedInACodeBlock(t *testing.T) {
	lab, _ := newStatusLab(t)
	lab.forge.set(lab.branch(), headA, nil, []forge.Review{{Author: "alice", State: forge.ReviewChangesRequested, Body: "Use *this* and _that_.", URL: "https://github.com/o/r/pull/1#pullrequestreview-5"}})
	lab.look()
	ths := lab.threadsOn()
	if len(ths) != 1 || !strings.Contains(ths[0].Body, "):\n\n```text\nUse *this* and _that_.\n```\n") {
		t.Fatalf("threads: %+v", ths)
	}
	if lab.item().Status != workitem.InProgress {
		t.Fatalf("the story is %s", lab.item().Status)
	}
}

// S-0397: a pull request a reviewer merged on the forge, whose branch the
// forge deleted, is not pushed again and gets no new pull request; the kept
// status records it merged, and a check that failed on it sends nothing back.
func TestAMergedPullRequestIsNeitherPushedAgainNorReopened(t *testing.T) {
	lab, git := newStatusLab(t)
	if st := lab.look(); st.Outcome != ForgePushed || st.PullRequest != 1 {
		t.Fatalf("the first look: %+v", st)
	}
	lab.forge.set(lab.branch(), headA, []forge.Check{buildFailed("2 failed")}, nil)
	lab.forge.setState(lab.branch(), forge.PullRequestMerged)
	git.mu.Lock()
	git.remote = "" // the forge deleted the branch on the merge
	git.mu.Unlock()
	for look := range 2 {
		st := lab.look()
		if st.Outcome != ForgeMerged || st.PullRequest != 1 || !strings.Contains(st.Reason, "is merged") || len(git.pushes) != 1 || lab.forge.opens != 1 {
			t.Fatalf("look %d after the merge: %+v, pushes %v, opens %d", look+1, st, git.pushes, lab.forge.opens)
		}
	}
	if st := lab.status(); st.State != forge.PullRequestMerged || st.PullRequest != 1 {
		t.Errorf("the kept status: %+v, want #1 merged", st)
	}
	if lab.sentBackBy() != 0 || lab.item().Status != workitem.Review {
		t.Errorf("a merged pull request sent the story back: %s", lab.item().Status)
	}
	for _, th := range lab.threadsOn() {
		if th.Open() {
			t.Errorf("a thread is open after the merge: %s %s", th.ID, th.Title)
		}
	}
}

// A read that finds no open pull request, because the forge merged it, keeps
// it merged by the number flai kept, with no error and no send back.
func TestAReadOfAMergedPullRequestKeepsItMerged(t *testing.T) {
	lab, _ := newStatusLab(t)
	lab.look()
	lab.forge.set(lab.branch(), headA, []forge.Check{buildFailed("2 failed")}, nil)
	lab.forge.setState(lab.branch(), forge.PullRequestMerged)
	var look readLook
	lab.pass.readBack(context.Background(), lab.repo, lab.item(), lab.pass.token("github"), &look)
	if len(look.failed) != 0 {
		t.Fatalf("the read of a merged pull request failed: %v", look.errs)
	}
	if st := lab.status(); st.State != forge.PullRequestMerged {
		t.Errorf("the kept status is %q, want merged", st.State)
	}
	if lab.sentBackBy() != 0 {
		t.Error("a merged pull request sent the story back")
	}
}
