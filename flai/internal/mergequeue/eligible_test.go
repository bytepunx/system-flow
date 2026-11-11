package mergequeue

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
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
	if name == "git" && len(args) == 3 && args[0] == "for-each-ref" && args[2] == "refs/heads/" {
		return g.refs(), nil
	}
	return "", errors.New("unscripted: " + name + " " + strings.Join(args, " "))
}

// refs is what git for-each-ref answers storygit.ReadHeads from the heads
// scripted as rev-parse answers: each branch at its commit, the one the
// abbrev-ref of HEAD names marked checked out.
func (g scriptedGit) refs() string {
	const verify = "rev-parse --verify --quiet "
	var lines []string
	for args, commit := range g {
		branch, ok := strings.CutPrefix(args, verify)
		if !ok || !strings.HasSuffix(branch, "^{commit}") {
			continue
		}
		branch = strings.TrimPrefix(strings.TrimSuffix(branch, "^{commit}"), "refs/heads/")
		mark := ""
		if g["rev-parse --abbrev-ref HEAD"] == branch {
			mark = "*"
		}
		lines = append(lines, "refs/heads/"+branch+"\t"+mark+"\t\t\tcommit\t"+commit)
	}
	slices.Sort(lines)
	return strings.Join(lines, "\n")
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

// countingGit is scriptedGit that counts every process it is asked for,
// from any goroutine.
type countingGit struct {
	scriptedGit
	calls *[]string
	mu    *sync.Mutex
}

func (g countingGit) Run(dir, name string, args ...string) (string, error) {
	g.mu.Lock()
	*g.calls = append(*g.calls, name+" "+strings.Join(args, " "))
	g.mu.Unlock()
	return g.scriptedGit.Run(dir, name, args...)
}

func (g countingGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

// S-0604: while no story is in review the queue reads no thread and runs no
// git, not even the main checkout's git status, and answers as it did: the
// state, no story listed, and no window. With a story in review it reads the
// threads, so a thread file that cannot be read is its error again.
func TestWithNoStoryInReviewTheQueueReadsNoThreadAndRunsNoGit(t *testing.T) {
	repo := project(t, "")
	story(t, repo, "Working", []bool{true})
	writeState(t, repo.CacheDir(), `{"enqueued": [{"story": "S-9", "by": "orchestrator", "at": "2026-10-19T08:00:00Z"}]}`)
	dir := filepath.Join(repo.WipDir(), threads.Folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "TH-0001-broken.md"), []byte("---\nid: [unclosed\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var calls []string
	g := countingGit{scriptedGit: scriptedGit{}, calls: &calls, mu: &sync.Mutex{}}

	q, err := ComputeAt(g, repo, t0)
	if err != nil {
		t.Fatalf("queue with no story in review read the threads: %v", err)
	}
	if len(calls) != 0 {
		t.Errorf("queue with no story in review ran %q, want no process", calls)
	}
	if q.Mode != "local" || q.Off || len(q.Eligible) != 0 || len(q.Ineligible) != 0 || q.Window != nil || q.Batch != nil || len(q.Outcomes) != 0 {
		t.Errorf("queue %+v, want local with no story listed, no batch, outcome, or window", q)
	}
	if len(q.Enqueued) != 1 || q.Enqueued[0].Story != "S-9" {
		t.Errorf("enqueued %+v, want the state's S-9", q.Enqueued)
	}

	story(t, repo, "Done building", []bool{true}, t0.Add(time.Hour))
	if _, err := ComputeAt(g, repo, t0); err == nil || !strings.Contains(err.Error(), "TH-0001-broken.md") {
		t.Errorf("queue with a story in review gave %v, want the broken thread file's error", err)
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

// ADR-0234: the queue carries the preflight's hold as the state has it, apart
// from the operator's, with when it began, each finding, and its thread, and
// flai board --json shows it under queue as preflight; with no hold it has
// none.
func TestTheQueueCarriesThePreflightHold(t *testing.T) {
	repo := project(t, "")
	g := scriptedGit{}
	if q := compute(t, g, repo); q.Preflight != nil {
		t.Errorf("preflight %+v with no state, want none", q.Preflight)
	}
	writeState(t, repo.CacheDir(), `{"preflight": {"since": "2026-10-19T08:00:00Z", "main": "m1",
  "findings": [{"tier": "lint", "checkout": "main", "path": "wip/kanban/tasks/T-1-x.md", "line": 3, "column": 5, "level": "error", "rule": "MD051", "alias": "link-fragments", "message": "Link fragments should be valid"}],
  "thread": "TH-0042"}, "outcomes": []}`)
	q := compute(t, g, repo)
	h := q.Preflight
	if h == nil || !h.Since.Equal(time.Date(2026, 10, 19, 8, 0, 0, 0, time.UTC)) || h.Main != "m1" || h.Thread != "TH-0042" {
		t.Fatalf("preflight %+v, want the state's hold since 08:00 at m1 with TH-0042", h)
	}
	if len(h.Findings) != 1 || h.Findings[0].Path != "wip/kanban/tasks/T-1-x.md" || h.Findings[0].Line != 3 || h.Findings[0].Rule != "MD051" {
		t.Errorf("findings %+v, want the one MD051 at line 3", h.Findings)
	}
	if q.Held {
		t.Error("held, want the preflight's hold apart from the operator's")
	}
	data, err := json.Marshal(Board{Queue: &q})
	if err != nil || !strings.Contains(string(data), `"queue":{`) || !strings.Contains(string(data), `"preflight":{"since":"2026-10-19T08:00:00Z","main":"m1","findings":[{"tier":"lint","checkout":"main","path":"wip/kanban/tasks/T-1-x.md","line":3,"column":5,"level":"error","rule":"MD051","alias":"link-fragments","message":"Link fragments should be valid"}],"thread":"TH-0042"}`) {
		t.Errorf("as the board's JSON: %v\n%s", err, data)
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

// S-0552: the window opens when the first of the stories picked was
// enqueued, whatever their order, and lasts the coalesce; there is none
// once it ends, with no story, with as many stories as the limit, and with
// no coalesce.
func TestTheWindowOpensWithTheFirstStoryEnqueued(t *testing.T) {
	picks := []Entry{{Story: "S-0002"}, {Story: "S-0001"}}
	enqueued := []Enqueued{{Story: "S-1", At: t0.Add(time.Minute)}, {Story: "S-0002", At: t0.Add(2 * time.Minute)}, {Story: "S-0003", At: t0}}
	w := window(picks, enqueued, 4, 3*time.Minute, t0.Add(3*time.Minute))
	if w == nil || !w.Ends.Equal(t0.Add(4*time.Minute)) || !slices.Equal(w.Stories, []string{"S-0002", "S-0001"}) || w.BatchMax != 4 {
		t.Errorf("window %+v, want S-0002 and S-0001 until %s, of 4", w, t0.Add(4*time.Minute))
	}
	for _, c := range []struct {
		name  string
		picks []Entry
		limit int
		wait  time.Duration
		now   time.Time
	}{
		{"ended", picks, 4, 3 * time.Minute, t0.Add(4 * time.Minute)},
		{"no story", nil, 4, 3 * time.Minute, t0},
		{"as many as the limit", picks, 2, 3 * time.Minute, t0},
		{"no coalesce", picks, 4, 0, t0},
	} {
		if w := window(c.picks, enqueued, c.limit, c.wait, c.now); w != nil {
			t.Errorf("%s: window %+v, want none", c.name, w)
		}
	}
}

// S-0552: the queue carries the window its next batch fills in, as flai
// queue --json and the board give it, while no batch is under way or the
// last was dropped with its result; not while the queue is held, by the
// operator or the preflight, a batch is
// under way or waits to land, the window has ended, or the main branch's
// head cannot be read.
func TestTheQueueCarriesTheWindow(t *testing.T) {
	repo := project(t, "")
	g := scriptedGit{}
	g.script("main", "rev-parse", "--abbrev-ref", "HEAD")
	g.script("m1", "rev-parse", "--verify", "--quiet", "refs/heads/main^{commit}")
	s, _ := story(t, repo, "Ready", []bool{true}, t0.Add(time.Hour))
	passed(t, repo, g, s.ID, "aaaa")
	enqueued := `"enqueued": [{"story": "` + s.ID + `", "by": "alex", "at": "2026-10-19T08:00:00Z"}]`
	ends := time.Date(2026, 10, 19, 8, 3, 0, 0, time.UTC)

	writeState(t, repo.CacheDir(), "{"+enqueued+"}")
	q, err := ComputeAt(g, repo, ends.Add(-time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if w := q.Window; w == nil || !w.Ends.Equal(ends) || !slices.Equal(w.Stories, []string{s.ID}) || w.BatchMax != 4 {
		t.Fatalf("window %+v, want %s until %s", w, s.ID, ends)
	}
	data, err := json.Marshal(q)
	if err != nil || !strings.Contains(string(data), `"window":{"ends":"2026-10-19T08:03:00Z","stories":["`+s.ID+`"],"batch_max":4}`) {
		t.Errorf("as JSON: %v\n%s", err, data)
	}
	if q, err := ComputeAt(g, repo, ends); err != nil || q.Window != nil {
		t.Errorf("window %+v, %v once it ended", q.Window, err)
	}
	dropped := `"batch": {"stories": ["` + s.ID + `"], "step": "verified", "result": {"passed": true, "finished_at": "2026-10-19T08:01:00Z"}, "dropped": {"reason": "x", "by": "alex", "at": "2026-10-19T08:01:30Z"}}`
	writeState(t, repo.CacheDir(), "{"+enqueued+", "+dropped+"}")
	if q, err := ComputeAt(g, repo, ends.Add(-time.Minute)); err != nil || q.Window == nil {
		t.Errorf("window %+v, %v with the last batch dropped, want one", q.Window, err)
	}

	for name, state := range map[string]string{
		"held":                  `"held": true, "held_by": "alex", "held_at": "2026-10-19T08:00:30Z"`,
		"held by the preflight": `"preflight": {"since": "2026-10-19T08:00:30Z", "findings": [{"tier": "lint", "checkout": "main", "path": "a.md", "message": "x"}]}`,
		"a batch under way":     `"batch": {"stories": ["` + s.ID + `"], "step": "verifying unit"}`,
		"a batch to land":       `"batch": {"stories": ["` + s.ID + `"], "step": "verified", "result": {"passed": true, "finished_at": "2026-10-19T08:01:00Z"}}`,
		"dropped while it ran":  `"batch": {"stories": ["` + s.ID + `"], "step": "verifying unit", "dropped": {"reason": "x", "by": "alex", "at": "2026-10-19T08:01:30Z"}}`,
	} {
		writeState(t, repo.CacheDir(), "{"+enqueued+", "+state+"}")
		if q, err := ComputeAt(g, repo, ends.Add(-time.Minute)); err != nil || q.Window != nil {
			t.Errorf("%s: window %+v, %v; want none", name, q.Window, err)
		}
	}
	writeState(t, repo.CacheDir(), "{"+enqueued+"}")
	delete(g, "rev-parse --abbrev-ref HEAD")
	if q, err := ComputeAt(g, repo, ends.Add(-time.Minute)); err != nil || q.Window != nil {
		t.Errorf("main unread: window %+v, %v; want none", q.Window, err)
	}
}

// S-0698, I-0333: the queue reads every branch head with one git
// for-each-ref and the main checkout's status with one git status, however
// many stories are in review, where it ran a rev-parse for each story and two
// more for the main branch. While no story in review has a passing flai
// verify result, no head can be held back by the main checkout's dirt, and
// it runs no git at all.
func TestTheQueuesGitDoesNotGrowWithTheStoriesInReview(t *testing.T) {
	status := []string{"-c", "core.quotePath=false", "status", "--porcelain", "--no-renames", "--untracked-files=all"}
	for _, n := range []int{1, 3, 6} {
		repo := project(t, "")
		var calls []string
		g := countingGit{scriptedGit: scriptedGit{}, calls: &calls, mu: &sync.Mutex{}}
		g.script("main", "rev-parse", "--abbrev-ref", "HEAD")
		g.script("m1", "rev-parse", "--verify", "--quiet", "refs/heads/main^{commit}")
		g.script("", status...)
		var stories []string
		for i := range n {
			s, _ := story(t, repo, "Verified", []bool{true}, t0.Add(time.Duration(i+1)*time.Hour))
			passed(t, repo, g.scriptedGit, s.ID, "c0ffee"+s.ID)
			stories = append(stories, s.ID)
		}
		writeState(t, repo.CacheDir(), `{"enqueued": [{"story": "`+stories[0]+`", "by": "alex", "at": "2026-10-19T08:00:00Z"}]}`)

		q, err := ComputeAt(g, repo, time.Date(2026, 10, 19, 8, 1, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		if got := ids(q.Eligible); !slices.Equal(got, stories) {
			t.Errorf("%d in review: eligible %v, want %v", n, got, stories)
		}
		if q.Window == nil {
			t.Errorf("%d in review: no window, want the enqueued story's, from main's head", n)
		}
		slices.Sort(calls)
		want := []string{"git -c core.quotePath=false status --porcelain --no-renames --untracked-files=all", "git for-each-ref"}
		if len(calls) != 2 || calls[0] != want[0] || !strings.HasPrefix(calls[1], want[1]+" ") {
			t.Errorf("%d in review: the queue ran %q, want one git status and one git for-each-ref", n, calls)
		}

		calls = calls[:0]
		for _, id := range stories {
			if err := verify.SaveReport(repo, verify.Report{Story: id, Commit: "c0ffee" + id, Passed: false}); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := ComputeAt(g, repo, t0); err != nil {
			t.Fatal(err)
		}
		if len(calls) != 0 {
			t.Errorf("%d in review, none verified: the queue ran %q, want no git", n, calls)
		}
	}
}

// deployed writes the kluster record of story's namespace in repo's main
// checkout: its deploy at commit with outcome, failed at failedStep when it
// failed, and its checks' commit and outcome.
func deployed(t *testing.T, repo *workitem.Repo, story, commit, outcome, failedStep, checksCommit, checksOutcome string) {
	t.Helper()
	ns, err := kluster.StoryNamespace(story)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(kluster.Record{Story: story, Namespace: ns, Commit: commit, Outcome: outcome, FailedStep: failedStep, ChecksCommit: checksCommit, ChecksOutcome: checksOutcome})
	if err != nil {
		t.Fatal(err)
	}
	path := kluster.RecordPath(repo.MainRoot, ns)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// ADR-0254 §5a: while the kluster host action is on, a story whose
// namespace was not deployed, or whose checks did not pass, at its branch's
// head is listed with the reason kluster, naming the last job and its
// outcome, and left out of every batch, as flai accept --batch would refuse
// it; once the gate passes, or with the action off, it is eligible.
func TestTheKlusterGateHoldsAStoryOutOfTheQueue(t *testing.T) {
	repo := project(t, "")
	repo.Manifest.Kluster = &manifest.Kluster{Config: "kluster.yaml"}
	g := scriptedGit{}
	s, _ := story(t, repo, "Deployed", []bool{true}, t0.Add(time.Hour))
	passed(t, repo, g, s.ID, "aaaa")
	writeState(t, repo.CacheDir(), `{"enqueued": [{"story": "`+s.ID+`", "by": "alex", "at": "2026-10-19T07:30:00Z"}]}`)
	on := Actions{Kluster: true, Checks: true}
	queue := func(act Actions) Queue {
		t.Helper()
		q, err := ComputeWith(g, repo, t0.Add(2*time.Hour), act)
		if err != nil {
			t.Fatal(err)
		}
		return q
	}
	held := func(name string, act Actions, says ...string) {
		t.Helper()
		q := queue(act)
		if len(q.Eligible) != 0 || len(q.Ineligible) != 1 {
			t.Fatalf("%s: eligible %v, ineligible %v; want %s held", name, ids(q.Eligible), ids(q.Ineligible), s.ID)
		}
		e := q.Ineligible[0]
		if got := kinds(e); !slices.Equal(got, []string{KindKluster}) {
			t.Fatalf("%s: reasons %v, want only %s", name, got, KindKluster)
		}
		for _, w := range append(says, "ADR-0254", "flai kluster deploy "+s.ID) {
			if !strings.Contains(e.Reasons[0].Text, w) {
				t.Errorf("%s: reason %q does not say %q", name, e.Reasons[0].Text, w)
			}
		}
		if picks := pick(q, "m1", 4); len(picks) != 0 {
			t.Errorf("%s: batch %v, want none", name, ids(picks))
		}
	}
	eligible := func(name string, act Actions) {
		t.Helper()
		q := queue(act)
		if !slices.Equal(ids(q.Eligible), []string{s.ID}) || len(q.Ineligible) != 0 {
			t.Errorf("%s: eligible %v, ineligible %+v; want %s eligible", name, ids(q.Eligible), q.Ineligible, s.ID)
		}
		if picks := pick(q, "m1", 4); !slices.Equal(ids(picks), []string{s.ID}) {
			t.Errorf("%s: batch %v, want %s", name, ids(picks), s.ID)
		}
	}

	held("never deployed", on, "was never deployed")
	eligible("never deployed, the action off", Actions{Checks: true})
	deployed(t, repo, s.ID, "aaaa", kluster.OutcomeFailed, kluster.StepDeploy, "", "")
	held("deploy failed", on, "the deploy of "+s.ID, "failed at the deploy step")
	eligible("deploy failed, the action off", Actions{})
	deployed(t, repo, s.ID, "9999", kluster.OutcomePassed, "", "9999", kluster.OutcomePassed)
	held("deployed at an older head", on, "was at 9999, not at the story's head aaaa")
	deployed(t, repo, s.ID, "aaaa", kluster.OutcomePassed, "", "aaaa", kluster.OutcomeFailed)
	held("checks failed", on, "checks run against it there ended failed", "and the checks run against it passed")
	eligible("checks failed, checks off", Actions{Kluster: true})
	deployed(t, repo, s.ID, "aaaa", kluster.OutcomePassed, "", "aaaa", kluster.OutcomePassed)
	eligible("deploy and checks passed at the head", on)

	repo.Manifest.Kluster = nil
	deployed(t, repo, s.ID, "aaaa", kluster.OutcomeFailed, kluster.StepDeploy, "", "")
	eligible("no kluster: block", on)
}
