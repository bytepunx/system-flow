package mergequeue

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// t0 is when the fixtures' stories are made; each enters review after it.
var t0 = time.Date(2026, 10, 19, 6, 0, 0, 0, time.UTC)

// scriptedGit answers the git commands scripted for it, by their arguments,
// and fails every other.
type scriptedGit map[string]string

func (g scriptedGit) script(out string, args ...string) { g[strings.Join(args, " ")] = out }

// head scripts the head of story's branch.
func (g scriptedGit) head(story, commit string) {
	g.script(commit, "rev-parse", "--verify", "--quiet", storygit.Branch(story)+"^{commit}")
}

func (g scriptedGit) Run(_, name string, args ...string) (string, error) {
	if out, ok := g[strings.Join(args, " ")]; ok && name == "git" {
		return out, nil
	}
	return "", errors.New("unscripted: " + name + " " + strings.Join(args, " "))
}

func (g scriptedGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (scriptedGit) LookPath(name string) (string, error) { return name, nil }

// project is a project, not a git repository, accepting by mode.
func project(t *testing.T, mode string) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if mode != "" {
		manifest += "acceptance:\n  mode: " + mode + "\n"
	}
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
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
	return repo
}

// story makes a story with one task done and the criteria ticks, and moves
// it to review at each of reviews, back to in-progress between them; with
// no reviews it stays in progress. It returns the story and its task.
func story(t *testing.T, repo *workitem.Repo, title string, ticks []bool, reviews ...time.Time) (*workitem.Item, *workitem.Item) {
	t.Helper()
	s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Nature: "feature", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	var boxes strings.Builder
	for i, tick := range ticks {
		mark := " "
		if tick {
			mark = "x"
		}
		boxes.WriteString("- [" + mark + "] criterion " + string(rune('a'+i)) + "\n")
	}
	data, err := os.ReadFile(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n"+boxes.String(), 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if s, err = repo.Get(s.ID); err != nil {
		t.Fatal(err)
	}
	for _, st := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.Transition(s, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Build " + title, Parent: s.ID, Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for _, st := range []string{workitem.Ready, workitem.InProgress, workitem.Done} {
		if _, err := repo.Transition(task, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	for i, at := range reviews {
		if i > 0 {
			if _, err := repo.Transition(s, workitem.InProgress, "alex", "rework", at.Add(-time.Minute)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := repo.Transition(s, workitem.Review, "alex", "", at); err != nil {
			t.Fatal(err)
		}
	}
	return s, task
}

// passed stores a passing flai verify result for story at commit, and
// scripts commit as its branch's head.
func passed(t *testing.T, repo *workitem.Repo, g scriptedGit, story, commit string) {
	t.Helper()
	if err := verify.SaveReport(repo, verify.Report{Story: story, Commit: commit, Passed: true}); err != nil {
		t.Fatal(err)
	}
	g.head(story, commit)
}

// ids are the stories of entries, in order.
func ids(entries []Entry) []string {
	out := []string{}
	for _, e := range entries {
		out = append(out, e.Story)
	}
	return out
}

// kinds are the kinds of e's reasons, in order.
func kinds(e Entry) []string {
	out := []string{}
	for _, r := range e.Reasons {
		out = append(out, r.Kind)
	}
	return out
}

// compute is the queue of repo through g, failing the test on an error.
func compute(t *testing.T, g scriptedGit, repo *workitem.Repo) Queue {
	t.Helper()
	q, err := Compute(g, repo)
	if err != nil {
		t.Fatal(err)
	}
	return q
}

// ADR-0173: the eligible stories are in the order they last entered review,
// oldest first, then by ID; a story sent back and moved to review again is
// ordered by its last entry. A story not in review is not listed.
func TestEligibleStoriesAreInTheOrderTheyLastEnteredReview(t *testing.T) {
	repo := project(t, "")
	g := scriptedGit{}
	s1, _ := story(t, repo, "Bounced", []bool{true}, t0.Add(time.Hour), t0.Add(4*time.Hour))
	s2, _ := story(t, repo, "Early", []bool{true, true}, t0.Add(2*time.Hour))
	s3, _ := story(t, repo, "Tied", []bool{true}, t0.Add(2*time.Hour))
	story(t, repo, "Working", []bool{true})
	for _, s := range []*workitem.Item{s1, s2, s3} {
		passed(t, repo, g, s.ID, "c0ffee"+s.ID)
	}
	q := compute(t, g, repo)
	if got, want := ids(q.Eligible), []string{s2.ID, s3.ID, s1.ID}; !slices.Equal(got, want) {
		t.Errorf("eligible %v, want %v", got, want)
	}
	if len(q.Ineligible) != 0 {
		t.Errorf("ineligible %+v, want none", q.Ineligible)
	}
	if q.Mode != "local" || q.Off {
		t.Errorf("mode %q off %v, want local and on", q.Mode, q.Off)
	}
	first := q.Eligible[2]
	if !first.EnteredReview.Equal(t0.Add(4*time.Hour)) || first.Head != "c0ffee"+s1.ID || first.Title != "Bounced" {
		t.Errorf("entry %+v, want its last entry to review and its head", first)
	}
}

// ADR-0173: each other story in review is listed after the eligible ones,
// in the same order, with every reason it is not eligible, each by kind,
// naming what fails.
func TestEachOtherStoryInReviewHasEveryReasonItIsNotEligible(t *testing.T) {
	repo := project(t, "")
	g := scriptedGit{}
	g.script("main", "rev-parse", "--abbrev-ref", "HEAD")

	ok, _ := story(t, repo, "Ok", []bool{true}, t0.Add(9*time.Hour))
	passed(t, repo, g, ok.ID, "aaaa")

	missing, _ := story(t, repo, "Unverified", []bool{true}, t0.Add(1*time.Hour))

	failed, _ := story(t, repo, "Failed", []bool{true}, t0.Add(2*time.Hour))
	if err := verify.SaveReport(repo, verify.Report{Story: failed.ID, Commit: "bbbb", Passed: false, StoppedAt: "test:flai"}); err != nil {
		t.Fatal(err)
	}

	stale, _ := story(t, repo, "Stale", []bool{true}, t0.Add(3*time.Hour))
	passed(t, repo, g, stale.ID, "cccc")
	g.head(stale.ID, "dddd")

	noBranch, _ := story(t, repo, "No branch", []bool{true}, t0.Add(4*time.Hour))
	if err := verify.SaveReport(repo, verify.Report{Story: noBranch.ID, Commit: "eeee", Passed: true}); err != nil {
		t.Fatal(err)
	}

	unticked, _ := story(t, repo, "Unticked", []bool{true, false, false}, t0.Add(5*time.Hour))
	passed(t, repo, g, unticked.ID, "ffff")

	threaded, task := story(t, repo, "Threaded", []bool{true}, t0.Add(6*time.Hour))
	passed(t, repo, g, threaded.ID, "1111")
	onStory := openThread(t, repo, threaded.ID, "Which way")
	onTask := openThread(t, repo, task.ID, "Which name")
	resolved := openThread(t, repo, threaded.ID, "Settled")
	if _, err := threads.Resolve(repo, resolved, "alex", "settled", t0); err != nil {
		t.Fatal(err)
	}

	undone, _ := story(t, repo, "Undoes main", []bool{true}, t0.Add(7*time.Hour))
	passed(t, repo, g, undone.ID, "2222")
	recordUndone(t, repo, g, undone.ID, "b.go", "a.go")

	broken, _ := story(t, repo, "Broken record", []bool{true}, t0.Add(8*time.Hour))
	if err := os.WriteFile(verify.ReportPath(repo, broken.ID), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}

	all, _ := story(t, repo, "Everything", []bool{false}, t0.Add(30*time.Minute))
	openThread(t, repo, all.ID, "Open")

	q := compute(t, g, repo)
	if got := ids(q.Eligible); !slices.Equal(got, []string{ok.ID}) {
		t.Errorf("eligible %v, want only %s", got, ok.ID)
	}
	wantOrder := []string{all.ID, missing.ID, failed.ID, stale.ID, noBranch.ID, unticked.ID, threaded.ID, undone.ID, broken.ID}
	if got := ids(q.Ineligible); !slices.Equal(got, wantOrder) {
		t.Fatalf("ineligible %v, want %v", got, wantOrder)
	}
	byID := map[string]Entry{}
	for _, e := range q.Ineligible {
		byID[e.Story] = e
	}
	cases := []struct {
		story string
		kinds []string
		says  []string
	}{
		{all.ID, []string{KindVerifyMissing, KindCriteriaUnticked, KindThreadOpen}, []string{"no flai verify result"}},
		{missing.ID, []string{KindVerifyMissing}, []string{"flai verify " + missing.ID}},
		{failed.ID, []string{KindVerifyFailed}, []string{"failed at its step test:flai", "bbbb"}},
		{stale.ID, []string{KindVerifyStale}, []string{"passed at cccc", "head of story/" + stale.ID + " is dddd"}},
		{noBranch.ID, []string{KindBranchUnread}, []string{"cannot read the head of story/" + noBranch.ID}},
		{unticked.ID, []string{KindCriteriaUnticked}, []string{"criteria 2, 3 are unticked", "flai criteria tick " + unticked.ID}},
		{threaded.ID, []string{KindThreadOpen}, []string{"threads " + onStory + ", " + onTask + " are not resolved"}},
		{undone.ID, []string{KindUndoneMain}, []string{"a.go, b.go", "never batches", "ADR-0155"}},
		{broken.ID, []string{KindUnreadable}, []string{"cannot read the last flai verify result", "delete it"}},
	}
	for _, c := range cases {
		e := byID[c.story]
		if got := kinds(e); !slices.Equal(got, c.kinds) {
			t.Errorf("%s reasons %v, want %v", c.story, got, c.kinds)
			continue
		}
		for _, says := range c.says {
			if !strings.Contains(e.Reasons[0].Text, says) {
				t.Errorf("%s's reason %q does not say %q", c.story, e.Reasons[0].Text, says)
			}
		}
	}
	if got := byID[unticked.ID].Reasons[0].Criteria; !slices.Equal(got, []int{2, 3}) {
		t.Errorf("unticked criteria %v, want [2 3]", got)
	}
	if got := byID[threaded.ID].Reasons[0].Threads; !slices.Equal(got, []string{onStory, onTask}) {
		t.Errorf("open threads %v, want %s and %s, not the resolved %s", got, onStory, onTask, resolved)
	}
	if got := byID[undone.ID].Reasons[0].Paths; !slices.Equal(got, []string{"a.go", "b.go"}) {
		t.Errorf("undone paths %v, want a.go and b.go", got)
	}
	if byID[stale.ID].Head != "dddd" || byID[missing.ID].Head != "" {
		t.Errorf("heads %q and %q, want the stale story's head and none for the unverified", byID[stale.ID].Head, byID[missing.ID].Head)
	}
}

// openThread opens a thread on the item on and returns its ID.
func openThread(t *testing.T, repo *workitem.Repo, on, title string) string {
	t.Helper()
	th, err := threads.New(repo, threads.NewOptions{Title: title, On: on, Author: "agent", Text: "Which?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	return th.ID
}

// recordUndone writes story's record of undone paths and scripts git to say
// its branch still removes main's line from each.
func recordUndone(t *testing.T, repo *workitem.Repo, g scriptedGit, story string, paths ...string) {
	t.Helper()
	rec := storygit.UndoneRecord{Story: story}
	for _, p := range paths {
		rec.Paths = append(rec.Paths, storygit.UndoneEntry{Undone: storygit.Undone{Path: p, Lines: []string{"main's line"}, OldBase: "b0"}, FoundAt: t0})
		g.script("diff --git a/"+p+" b/"+p+"\n@@ -1 +0,0 @@\n-main's line", "diff", "--no-ext-diff", "--no-color", "--no-renames", "-U0", "main", storygit.Branch(story), "--", p)
	}
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	path := storygit.UndoneRecordPath(repo.CacheDir(), story)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// ADR-0173: under pull-request acceptance the queue is off, says why, and
// lists no story, eligible or not, reading neither git nor the state.
func TestUnderPullRequestAcceptanceTheQueueIsOff(t *testing.T) {
	repo := project(t, "pull-request")
	g := scriptedGit{}
	s, _ := story(t, repo, "Ready to land", []bool{true}, t0.Add(time.Hour))
	passed(t, repo, g, s.ID, "aaaa")
	writeState(t, repo.CacheDir(), "{")
	q, err := Compute(g, repo)
	if err != nil {
		t.Fatal(err)
	}
	if !q.Off || q.Mode != "pull-request" || len(q.Eligible) != 0 || len(q.Ineligible) != 0 || q.Batch != nil {
		t.Errorf("queue %+v, want off with no story", q)
	}
	for _, says := range []string{"local acceptance only", "required checks", "forge_stale", "ADR-0169"} {
		if !strings.Contains(q.OffReason, says) {
			t.Errorf("off reason %q does not say %q", q.OffReason, says)
		}
	}
	data, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	for _, says := range []string{`"off":true`, `"eligible":[]`, `"ineligible":[]`} {
		if !strings.Contains(string(data), says) {
			t.Errorf("JSON %s does not hold %s", data, says)
		}
	}
}

// The queue carries the state's enqueued stories, its batch and step, and
// its outcomes, and marks each enqueued story; with no state it has no
// batch and no outcomes.
func TestTheQueueCarriesTheState(t *testing.T) {
	repo := project(t, "")
	g := scriptedGit{}
	s1, _ := story(t, repo, "One", []bool{true}, t0.Add(time.Hour))
	s2, _ := story(t, repo, "Two", []bool{true}, t0.Add(2*time.Hour))
	passed(t, repo, g, s1.ID, "aaaa")
	passed(t, repo, g, s2.ID, "bbbb")

	empty := compute(t, g, repo)
	if empty.Batch != nil || len(empty.Outcomes) != 0 || len(empty.Enqueued) != 0 || empty.Eligible[0].Enqueued {
		t.Errorf("queue with no state %+v, want no batch, outcomes, or enqueued", empty)
	}

	writeState(t, repo.CacheDir(), `{"enqueued": [{"story": "S-2", "by": "orchestrator", "at": "2026-10-19T08:00:00Z"}],
  "batch": {"stories": ["`+s2.ID+`"], "step": "verifying"},
  "outcomes": [{"story": "`+s1.ID+`", "outcome": "sent-back", "at": "2026-10-19T07:00:00Z", "detail": "test:flai failed"}]}`)
	q := compute(t, g, repo)
	if q.Batch == nil || q.Batch.Step != "verifying" || !slices.Equal(q.Batch.Stories, []string{s2.ID}) {
		t.Errorf("batch %+v, want %s verifying", q.Batch, s2.ID)
	}
	if len(q.Outcomes) != 1 || q.Outcomes[0].Outcome != "sent-back" || len(q.Enqueued) != 1 {
		t.Errorf("outcomes %+v, enqueued %+v", q.Outcomes, q.Enqueued)
	}
	if q.Eligible[0].Enqueued || !q.Eligible[1].Enqueued {
		t.Errorf("enqueued marks %v and %v, want only %s, enqueued as S-2", q.Eligible[0].Enqueued, q.Eligible[1].Enqueued, s2.ID)
	}
}

// A malformed state is an error of the queue, naming the file.
func TestAMalformedStateIsAnErrorOfTheQueue(t *testing.T) {
	repo := project(t, "")
	writeState(t, repo.CacheDir(), "{")
	if _, err := Compute(scriptedGit{}, repo); err == nil || !strings.Contains(err.Error(), StatePath(repo.CacheDir())) {
		t.Errorf("error %v, want one naming %s", err, StatePath(repo.CacheDir()))
	}
}

// ADR-0174: the queue can batch only with the host action on, accept_reviews
// on, and local acceptance, as flai serve builds a batch.
func TestTheQueueBatchesOnlyWithTheActionAcceptReviewsAndLocalAcceptance(t *testing.T) {
	on := manifest.Manifest{Orchestration: manifest.Orchestration{Permissions: manifest.Permissions{AcceptReviews: true}}}
	local := on
	local.Acceptance.Mode = manifest.AcceptLocal
	pr := on
	pr.Acceptance.Mode = manifest.AcceptPullRequest
	for _, c := range []struct {
		name   string
		action bool
		m      manifest.Manifest
		want   bool
	}{
		{"every condition holds, acceptance by default", true, on, true},
		{"every condition holds, acceptance named local", true, local, true},
		{"the host action off", false, on, false},
		{"accept_reviews off", true, manifest.Manifest{}, false},
		{"pull-request acceptance", true, pr, false},
	} {
		if got := Batches(c.action, c.m); got != c.want {
			t.Errorf("%s: Batches = %v, want %v", c.name, got, c.want)
		}
	}
}
