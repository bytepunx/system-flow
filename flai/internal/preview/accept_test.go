package preview

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// headAt is a runner that answers every command with head, as git rev-parse
// answers a branch's head.
type headAt string

func (h headAt) Run(string, string, ...string) (string, error) { return string(h), nil }

func (h headAt) RunInput(string, string, string, ...string) (string, error) { return string(h), nil }

func (headAt) LookPath(name string) (string, error) { return name, nil }

// ADR-0170: a pull request flai serve read as merged on the forge already
// opens the gate, whatever its checks, reviews, and head then, and the
// preview says it is merged, with the number, the head, and the method; one
// still open is merged at the head the gate passed at.
func TestAKeptPullRequestMergedAlreadyHasNoForgeBlockers(t *testing.T) {
	at := time.Date(2026, 10, 18, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s := forgeRepo(t, "pull-request", at)
	repo.Manifest.Acceptance.MergeMethod = "squash"
	st := forgeStatus([]forge.CheckStatus{check("lint", forge.CheckFailed, "2 errors")}, nil, 2)
	st.State = forge.PullRequestMerged
	if err := forge.WriteStatus(repo.CacheDir(), st); err != nil {
		t.Fatal(err)
	}

	res, err := AcceptWith(noGit{}, repo, s, "alex", at, log, AcceptOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.ForgeBlockers) != 0 || len(res.Blockers) != 0 {
		t.Errorf("a pull request merged already has no forge blocker: %+v %v", res.ForgeBlockers, res.Blockers)
	}
	want := RemoteMerge{PullRequest: 7, Head: forgeHead, Method: forge.MergeSquash, AlreadyMerged: true}
	if res.Remote == nil || *res.Remote != want || res.PullRequest != forgePR {
		t.Errorf("remote = %+v, pull request %q; want %+v and %q", res.Remote, res.PullRequest, want, forgePR)
	}

	// the story's branch has moved since: merged is merged
	moved := &Acceptance{ID: s.ID, Branch: "story/" + s.ID}
	forgeGate(headAt(forgeOther), repo, s.ID, moved)
	if len(moved.ForgeBlockers) != 0 || len(moved.Blockers) != 0 {
		t.Errorf("a merged pull request is not stale for a branch that moved: %+v", moved.ForgeBlockers)
	}

	// open and passing: merged at the head the gate passed at, by the method
	st.State, st.Checks, st.RequiredApprovals = forge.PullRequestOpen, nil, 0
	if err := forge.WriteStatus(repo.CacheDir(), st); err != nil {
		t.Fatal(err)
	}
	open := &Acceptance{ID: s.ID, Branch: "story/" + s.ID}
	forgeGate(headAt(forgeHead), repo, s.ID, open)
	want.AlreadyMerged = false
	if len(open.ForgeBlockers) != 0 || open.Remote == nil || *open.Remote != want {
		t.Errorf("an open pull request that passes: blockers %+v, remote %+v, want %+v", open.ForgeBlockers, open.Remote, want)
	}

	// no state kept: the method only, and the pull request missing
	none, s2 := forgeRepo(t, "pull-request", at)
	res, err = AcceptWith(noGit{}, none, s2, "alex", at, log, AcceptOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Remote == nil || *res.Remote != (RemoteMerge{Method: forge.MergeMerge}) || len(res.ForgeBlockers) != 1 {
		t.Errorf("no kept state: remote %+v, blockers %+v", res.Remote, res.ForgeBlockers)
	}

	// local acceptance merges nothing through the remote
	local, s3 := forgeRepo(t, "local", at)
	if res, err = AcceptWith(noGit{}, local, s3, "alex", at, log, AcceptOptions{}); err != nil || res.Remote != nil {
		t.Errorf("a local project's preview has no merge through the remote: %+v %v", res.Remote, err)
	}
}

// gitAnswers is a runner that answers each git command it holds, by its
// arguments joined with spaces, and fails every other; it finds no git, so a
// preview through it runs as in a project that is not a repository.
type gitAnswers map[string]string

func (g gitAnswers) Run(_, _ string, args ...string) (string, error) {
	if out, ok := g[strings.Join(args, " ")]; ok {
		return out, nil
	}
	return "", errors.New("no answer")
}

func (g gitAnswers) RunInput(string, string, string, ...string) (string, error) {
	return "", errors.New("no answer")
}

func (gitAnswers) LookPath(string) (string, error) { return "", errors.New("no git") }

// S-0400 criterion 2, ADR-0176: with a batch, the preview blocks the
// acceptance on each reason the merge queue does not let the story land,
// for the orchestrator as an unverified commit in place of --verified's
// check, and once the story may land takes its rebased head in the batch as
// the commit verified.
func TestALandingTakesTheBatchsVerification(t *testing.T) {
	at := time.Date(2026, 10, 19, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s := forgeRepo(t, "local", at)
	const batch = "20261019-020000"
	git := gitAnswers{
		"rev-parse --verify --quiet story/" + s.ID + "^{commit}": "aaa111",
		"rev-parse --abbrev-ref HEAD":                            "main",
		"rev-parse --verify --quiet refs/heads/main^{commit}":    "ccc333",
	}
	unverified := func(res *Acceptance) []string {
		var out []string
		for _, b := range res.OrchestratorBlockers {
			if b.Code == BlockUnverified {
				out = append(out, b.Message)
			}
		}
		return out
	}

	// no batch in the queue: its reason blocks, and nothing is verified
	none := "the merge queue holds no batch, so " + s.ID + " cannot land as part of merge queue batch " + batch
	res, err := AcceptWith(git, repo, s, "alex", at, log, AcceptOptions{Batch: batch})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Blockers) != 1 || !strings.HasPrefix(res.Blockers[0], none) || res.Verified != "" || res.Batch != batch {
		t.Errorf("a batch the queue does not hold blocks: %q, verified %q, batch %q", res.Blockers, res.Verified, res.Batch)
	}
	res, err = AcceptWith(git, repo, s, "orchestrator", at, log, AcceptOptions{Orchestrator: true, Batch: batch})
	if err != nil {
		t.Fatal(err)
	}
	if got := unverified(res); len(got) != 1 || !strings.HasPrefix(got[0], none) || !slices.Contains(res.Blockers, got[0]) {
		t.Errorf("the orchestrator's landing is unverified by the queue's reason alone, not --verified's: %q", got)
	}

	// the batch passed and lands the story next, at the heads it verified
	if _, err := mergequeue.Update(repo.CacheDir(), func(st *mergequeue.State) error {
		st.Batch = &mergequeue.Batch{
			ID: batch, Base: "ccc333", Stories: []string{s.ID},
			Result:  &mergequeue.Result{Passed: true, FinishedAt: at},
			Heads:   []mergequeue.Head{{Story: s.ID, Verified: "aaa111", Rebased: "bbb222"}},
			Landing: &mergequeue.Landing{Stories: []string{s.ID}, Landed: []string{}, Head: "ccc333"},
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	res, err = AcceptWith(git, repo, s, "alex", at, log, AcceptOptions{Batch: batch})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Blockers) != 0 || res.Verified != "bbb222" || res.Batch != batch {
		t.Errorf("a story the batch lands next is verified at its rebased head: %q, verified %q", res.Blockers, res.Verified)
	}
	res, err = AcceptWith(git, repo, s, "orchestrator", at, log, AcceptOptions{Orchestrator: true, Batch: batch})
	if err != nil {
		t.Fatal(err)
	}
	if got := unverified(res); len(got) != 0 || res.Verified != "bbb222" {
		t.Errorf("the orchestrator's landing takes the batch's verification: %q, verified %q", got, res.Verified)
	}
}

// markedBranch is gitAnswers whose git cat-file --batch answers one blob,
// text, as the blob of the one file the branch changes.
type markedBranch struct {
	gitAnswers
	text string
}

func (m markedBranch) RunInput(_, _, _ string, args ...string) (string, error) {
	if strings.Join(args, " ") != "cat-file --batch" {
		return "", errors.New("no answer")
	}
	return "abc123 blob " + strconv.Itoa(len(m.text)) + "\n" + m.text + "\n", nil
}

// ADR-0248: the conflict-markers blocker names each marker's path and line
// and says to send the story back to in-progress for its agent to resolve
// them, in its worktree when it has one and on its branch otherwise; it
// never asks the operator to resolve them and accept again.
func TestConflictMarkersSaysToSendTheStoryBack(t *testing.T) {
	at := time.Date(2026, 11, 2, 3, 0, 0, 0, time.UTC)
	repo, s := forgeRepo(t, "local", at)
	branch := "story/" + s.ID
	git := markedBranch{gitAnswers{
		"rev-parse --abbrev-ref HEAD": "main",
		"diff --raw -z --no-abbrev --no-renames --diff-filter=AM main..." + branch: ":100644 100644 " + strings.Repeat("0", 40) + " abc123 M\x00design/x.md\x00",
	}, "# X\n<<<<<<< HEAD\nours\n=======\ntheirs\n>>>>>>> theirs\n"}
	// the divider is not a marker
	markers := branch + " carries merge conflict markers at design/x.md:2, design/x.md:6; send " + s.ID + " back to in-progress so that its agent resolves each conflict "
	rest := ", keeping what both sides meant, removes the markers, commits, and moves it to review again (ADR-0248)"

	got, err := ConflictMarkers(git, repo, s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if want := markers + "on " + branch + ", whose worktree flai stream open " + s.ID + " opens" + rest; got != want {
		t.Errorf("without a worktree:\n got %q\nwant %q", got, want)
	}
	if err := os.MkdirAll(repo.WorktreePath(s.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, err = ConflictMarkers(git, repo, s.ID); err != nil {
		t.Fatal(err)
	}
	if want := markers + "in .flai-cache/worktrees/" + s.ID + rest; got != want {
		t.Errorf("with a worktree:\n got %q\nwant %q", got, want)
	}
	if strings.Contains(got, "accept again") {
		t.Errorf("the blocker asks the operator to resolve the markers: %q", got)
	}
}

// withGit is gitAnswers on a host with git: the project is a repository, and
// every command it holds no answer for fails.
type withGit struct{ gitAnswers }

func (withGit) LookPath(name string) (string, error) { return name, nil }

// saveQueuedReport stores story id's last flai verify report at commit,
// passed, with the tier integration in the state given.
func saveQueuedReport(t *testing.T, repo *workitem.Repo, id, commit string, integration verify.State) {
	t.Helper()
	rep := verify.Report{Story: id, Commit: commit, Passed: true, Steps: []verify.Step{
		{Name: "test", Tier: true, State: verify.Passed},
		{Name: "integration", Tier: true, State: integration},
	}}
	if err := verify.SaveReport(repo, rep); err != nil {
		t.Fatal(err)
	}
}

// ADR-0221: while the merge queue batches, the preview of a hand acceptance,
// as flai accept --dry-run prints it, lists the blocker queue_tiers for a
// story whose last flai verify left integration to the queue at its branch's
// head, and the tiers; under --without-queue-tiers the tiers but no blocker.
// The orchestrator's acceptance and the queue's landing are not gated, nor
// any acceptance while the queue does not batch.
func TestAHandAcceptanceWhileTheQueueBatchesListsTheQueueTiersBlocker(t *testing.T) {
	at := time.Date(2026, 10, 30, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s := forgeRepo(t, "local", at)
	git := withGit{gitAnswers{
		"rev-parse --is-inside-work-tree":                        "true",
		"rev-parse --verify --quiet refs/heads/story/" + s.ID:    forgeHead,
		"rev-parse --verify --quiet story/" + s.ID + "^{commit}": forgeHead,
	}}
	saveQueuedReport(t, repo, s.ID, forgeHead, verify.Queued)

	res, err := AcceptWith(git, repo, s, "alex", at, log, AcceptOptions{QueueBatches: true})
	if err != nil {
		t.Fatal(err)
	}
	want := s.ID + "'s tiers marked queue have not run at its head: its last flai verify did not run tier integration, marked queue; while the merge queue batches, a hand acceptance waits for them: run flai verify " + s.ID + " --all"
	if b := res.QueueTiersBlocker; b == nil || b.Code != BlockQueueTiers || !strings.HasPrefix(b.Message, want) || !slices.Contains(res.Blockers, b.Message) ||
		!slices.Equal(res.QueueTiers, []string{"integration"}) || res.QueueTiersWhy == "" {
		t.Errorf("the hand acceptance must list queue_tiers naming integration: %+v %q %q", res.QueueTiersBlocker, res.QueueTiers, res.Blockers)
	}

	res, err = AcceptWith(git, repo, s, "alex", at, log, AcceptOptions{QueueBatches: true, WithoutQueueTiers: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.QueueTiersBlocker != nil || !slices.Equal(res.QueueTiers, []string{"integration"}) || slices.ContainsFunc(res.Blockers, func(b string) bool { return strings.HasPrefix(b, want) }) {
		t.Errorf("--without-queue-tiers names the tiers and lists no blocker: %+v %q %q", res.QueueTiersBlocker, res.QueueTiers, res.Blockers)
	}

	for name, opts := range map[string]AcceptOptions{
		"the orchestrator's acceptance": {QueueBatches: true, Orchestrator: true},
		"the queue's landing":           {QueueBatches: true, Batch: "20261030-020000"},
		"the queue not batching":        {},
	} {
		res, err := AcceptWith(git, repo, s, "alex", at, log, opts)
		if err != nil {
			t.Fatal(err)
		}
		if res.QueueTiersBlocker != nil || res.QueueTiers != nil || res.QueueTiersWhy != "" {
			t.Errorf("%s is not gated on the queue tiers: %+v %q %q", name, res.QueueTiersBlocker, res.QueueTiers, res.QueueTiersWhy)
		}
	}

	saveQueuedReport(t, repo, s.ID, forgeHead, verify.Passed)
	if res, err = AcceptWith(git, repo, s, "alex", at, log, AcceptOptions{QueueBatches: true}); err != nil || res.QueueTiersBlocker != nil || res.QueueTiersWhy != "" {
		t.Errorf("with integration passed at the head nothing is gated: %+v %q %v", res.QueueTiersBlocker, res.QueueTiersWhy, err)
	}
}

// ADR-0221: a branch head that cannot be read, a report that cannot be read,
// and a report at another commit are the blocker, saying why, so the gate
// never opens unread.
func TestTheQueueTiersGateNeverOpensUnread(t *testing.T) {
	at := time.Date(2026, 10, 30, 3, 0, 0, 0, time.UTC)
	repo, s := forgeRepo(t, "local", at)
	gate := func(r gitAnswers) *Acceptance {
		res := &Acceptance{ID: s.ID, Branch: "story/" + s.ID}
		queueTiersGate(r, repo, s, res, AcceptOptions{QueueBatches: true})
		return res
	}
	head := gitAnswers{"rev-parse --verify --quiet story/" + s.ID + "^{commit}": forgeHead}
	for name, c := range map[string]struct {
		git    gitAnswers
		report func()
		why    string
	}{
		"head unread": {gitAnswers{}, func() {}, "cannot read the head of story/" + s.ID + ": no answer"},
		"no report":   {head, func() {}, "no flai verify has run on it"},
		"older commit": {head, func() { saveQueuedReport(t, repo, s.ID, forgeOther, verify.Passed) },
			"its last flai verify ran at 0123456, not at its head abc1234"},
		"report unread": {head, func() {
			path := verify.ReportPath(repo, s.ID)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
				t.Fatal(err)
			}
		}, "cannot read the last verify report of " + s.ID},
	} {
		_ = os.Remove(verify.ReportPath(repo, s.ID))
		c.report()
		res := gate(c.git)
		if !strings.HasPrefix(res.QueueTiersWhy, c.why) || res.QueueTiersBlocker == nil || len(res.Blockers) != 1 || !strings.Contains(res.Blockers[0], c.why) {
			t.Errorf("%s: why %q, blockers %q; want %q", name, res.QueueTiersWhy, res.Blockers, c.why)
		}
	}
}

// writeKlusterRecord writes r as its namespace's deployment record in repo's
// main checkout, as a kluster deploy keeps it.
func writeKlusterRecord(t *testing.T, repo *workitem.Repo, r kluster.Record) {
	t.Helper()
	path := kluster.RecordPath(repo.MainRoot, r.Namespace)
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// ADR-0254 §5a: while the kluster host action is on and the manifest has a
// kluster: block, every acceptance of a story, the operator's, the
// orchestrator's (among its own blockers), and the queue's landing, lists the
// blocker kluster until its namespace's deploy, and with the checks action on
// its checks run, passed at its branch's head; --without-kluster lifts it for
// the operator's alone, keeping why. With the action off, or no kluster:
// block, nothing is gated.
func TestTheKlusterGateHoldsEveryAcceptanceWhileTheActionIsOn(t *testing.T) {
	at := time.Date(2026, 11, 11, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s := forgeRepo(t, "local", at)
	repo.Manifest.Kluster = &manifest.Kluster{Config: "kluster.yaml"}
	git := withGit{gitAnswers{
		"rev-parse --is-inside-work-tree":                                   "true",
		"rev-parse --verify --quiet refs/heads/story/" + s.ID:               forgeHead,
		"rev-parse --verify --quiet refs/heads/story/" + s.ID + "^{commit}": forgeHead,
	}}
	ns := strings.ToLower(s.ID)
	never := s.ID + "'s namespace " + ns + " was never deployed"
	accept := func(opts AcceptOptions) *Acceptance {
		t.Helper()
		res, err := AcceptWith(git, repo, s, "alex", at, log, opts)
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	held := func(res *Acceptance) bool {
		return res.KlusterBlocker != nil && slices.Contains(res.Blockers, res.KlusterBlocker.Message)
	}
	passes := func(name string, res *Acceptance) {
		t.Helper()
		if res.KlusterWhy != "" || res.KlusterBlocker != nil || slices.ContainsFunc(res.Blockers, func(b string) bool { return strings.Contains(b, "kluster host action") }) {
			t.Errorf("%s is not gated: why %q, blocker %+v, blockers %q", name, res.KlusterWhy, res.KlusterBlocker, res.Blockers)
		}
	}

	// off, or no block: nothing is gated, though nothing was deployed
	passes("the action off", accept(AcceptOptions{ChecksOn: true}))
	repo.Manifest.Kluster = nil
	passes("no kluster: block", accept(AcceptOptions{KlusterOn: true, ChecksOn: true}))
	repo.Manifest.Kluster = &manifest.Kluster{Config: "kluster.yaml"}

	// on and never deployed: every acceptance is held, naming why
	on := AcceptOptions{KlusterOn: true, ChecksOn: true}
	res := accept(on)
	if b := res.KlusterBlocker; !held(res) || b.Code != BlockKluster || !strings.Contains(b.Message, never) || !strings.Contains(b.Message, "flai accept "+s.ID+" --without-kluster") || slices.Contains(res.OrchestratorBlockers, *b) {
		t.Errorf("the operator's acceptance must list kluster: %+v %q", res.KlusterBlocker, res.Blockers)
	}
	res = accept(AcceptOptions{KlusterOn: true, ChecksOn: true, Orchestrator: true})
	if b := res.KlusterBlocker; !held(res) || !slices.Contains(res.OrchestratorBlockers, *b) {
		t.Errorf("the orchestrator's acceptance must list kluster among its own: %+v %+v", res.KlusterBlocker, res.OrchestratorBlockers)
	}
	if res = accept(AcceptOptions{KlusterOn: true, ChecksOn: true, Batch: "20261111-020000"}); !held(res) {
		t.Errorf("the queue's landing must list kluster: %+v %q", res.KlusterBlocker, res.Blockers)
	}

	// --without-kluster is the operator's: no blocker, why kept
	res = accept(AcceptOptions{KlusterOn: true, ChecksOn: true, WithoutKluster: true})
	passed := res.KlusterBlocker == nil && !slices.ContainsFunc(res.Blockers, func(b string) bool { return strings.Contains(b, never) })
	if !passed || !strings.HasPrefix(res.KlusterWhy, never) {
		t.Errorf("--without-kluster lists no blocker and keeps why: %q %+v %q", res.KlusterWhy, res.KlusterBlocker, res.Blockers)
	}
	for name, opts := range map[string]AcceptOptions{
		"the orchestrator's": {KlusterOn: true, Orchestrator: true, WithoutKluster: true},
		"the queue's":        {KlusterOn: true, Batch: "20261111-020000", WithoutKluster: true},
	} {
		if res := accept(opts); !held(res) {
			t.Errorf("--without-kluster does not lift %s acceptance: %+v", name, res.KlusterBlocker)
		}
	}

	deployed := kluster.Record{Story: s.ID, Namespace: ns, Commit: forgeHead, Outcome: kluster.OutcomePassed}
	for _, c := range []struct {
		name   string
		record kluster.Record
		checks bool
		why    string // "" when the gate passes
	}{
		{"deployed at an older head", kluster.Record{Story: s.ID, Namespace: ns, Commit: forgeOther, Outcome: kluster.OutcomePassed}, true, "not at the story's head abc1234def56"},
		{"deploy failed", kluster.Record{Story: s.ID, Namespace: ns, Commit: forgeHead, Outcome: kluster.OutcomeFailed, FailedStep: kluster.StepDeploy}, true, "failed at the deploy step"},
		{"checks not run", deployed, true, "no checks run against it has ended there yet"},
		{"checks off", deployed, false, ""},
		{"checks failed", kluster.Record{Story: s.ID, Namespace: ns, Commit: forgeHead, Outcome: kluster.OutcomePassed, ChecksCommit: forgeHead, ChecksOutcome: kluster.OutcomeFailed}, true, "checks run against it there ended failed"},
		{"deploy and checks passed", kluster.Record{Story: s.ID, Namespace: ns, Commit: forgeHead, Outcome: kluster.OutcomePassed, ChecksCommit: forgeHead, ChecksOutcome: kluster.OutcomePassed}, true, ""},
	} {
		writeKlusterRecord(t, repo, c.record)
		res := accept(AcceptOptions{KlusterOn: true, ChecksOn: c.checks})
		if c.why == "" {
			passes(c.name, res)
			continue
		}
		if !held(res) || !strings.Contains(res.KlusterWhy, c.why) {
			t.Errorf("%s: why %q, blockers %q; want %q", c.name, res.KlusterWhy, res.Blockers, c.why)
		}
	}
}

// ADR-0254 §5a: a record or a head that cannot be read holds the story,
// saying why, so the kluster gate never opens unread.
func TestTheKlusterGateNeverOpensUnread(t *testing.T) {
	at := time.Date(2026, 11, 11, 3, 0, 0, 0, time.UTC)
	repo, s := forgeRepo(t, "local", at)
	repo.Manifest.Kluster = &manifest.Kluster{Config: "kluster.yaml"}
	gate := func(r gitAnswers) *Acceptance {
		res := &Acceptance{ID: s.ID, Branch: "story/" + s.ID}
		klusterGate(r, repo, s, res, AcceptOptions{KlusterOn: true})
		return res
	}
	res := gate(gitAnswers{})
	if !strings.HasPrefix(res.KlusterWhy, "cannot resolve the head of "+s.ID+"'s branch refs/heads/story/"+s.ID) || res.KlusterBlocker == nil || len(res.Blockers) != 1 {
		t.Errorf("a head that cannot be read: why %q, blockers %q", res.KlusterWhy, res.Blockers)
	}
	path := kluster.RecordPath(repo.MainRoot, strings.ToLower(s.ID))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	res = gate(gitAnswers{"rev-parse --verify --quiet refs/heads/story/" + s.ID + "^{commit}": forgeHead})
	if !strings.HasPrefix(res.KlusterWhy, "cannot tell whether its namespace was deployed at its head: read the deployment record") || res.KlusterBlocker == nil || len(res.Blockers) != 1 {
		t.Errorf("a record that cannot be read: why %q, blockers %q", res.KlusterWhy, res.Blockers)
	}
}

// S-0422 criterion 4: a story whose branch changes nothing beyond main is
// accepted as any story, formula or not: its preview names no blocker for the
// empty diff. Under the orchestrator, an ordinary story's evidence is still
// the files its branch changes, which an empty branch has none of, while a
// formula's instance's is the done task of each step.
func TestAStoryWhoseBranchChangesNothingIsAccepted(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	at := time.Date(2026, 11, 2, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s := forgeRepo(t, "local", at)
	r := execx.System{}
	git := func(args ...string) string {
		t.Helper()
		out, err := r.Run(repo.MainRoot, "git", args...)
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
		return strings.TrimSpace(out)
	}
	if err := os.WriteFile(filepath.Join(repo.MainRoot, ".gitignore"), []byte(".flai-cache/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "init"}, {"branch", "story/" + s.ID}} {
		git(args...)
	}
	head := git("rev-parse", "story/"+s.ID)

	res, err := AcceptWith(r, repo, s, "alex", at, log, AcceptOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Blockers) != 0 || res.Branch != "story/"+s.ID {
		t.Errorf("a story whose branch changes nothing is accepted: branch %q, blockers %q", res.Branch, res.Blockers)
	}

	ev, err := ParseEvidence("Verdict: pass\n- 1: `cli/feature.go`")
	if err != nil {
		t.Fatal(err)
	}
	res, err = AcceptWith(r, repo, s, "orchestrator", at, log, AcceptOptions{Orchestrator: true, Verified: head, Evidence: ev})
	if err != nil {
		t.Fatal(err)
	}
	if b := res.OrchestratorBlockers; len(b) != 1 || b[0].Code != BlockCriterionUnevidenced || !strings.Contains(b[0].Message, "names no file the branch changes (it names cli/feature.go)") {
		t.Errorf("an ordinary story's evidence is still its changed files: %+v", b)
	}

	s.Formula = "weekly-check"
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	if ev, err = ParseEvidence("Verdict: pass\n- 1: T-0001 done"); err != nil {
		t.Fatal(err)
	}
	for by, opts := range map[string]AcceptOptions{"alex": {}, "orchestrator": {Orchestrator: true, Verified: head, Evidence: ev}} {
		if res, err = AcceptWith(r, repo, s, by, at, log, opts); err != nil {
			t.Fatal(err)
		}
		if len(res.Blockers) != 0 || len(res.OrchestratorBlockers) != 0 {
			t.Errorf("an instance whose branch changes nothing is accepted by %s, on its steps' tasks: %q", by, res.Blockers)
		}
	}
}

// S-0422 criterion 4: an instance's criterion is blocked unless its item
// names a task of the instance, in any zero padding, that is done.
func TestStepsUndoneNamesEachCriterionWhoseTaskIsNotDone(t *testing.T) {
	it := &workitem.Item{ID: "S-0001", Type: workitem.Story, Formula: "weekly-check"}
	tasks := map[string]*workitem.Item{
		"T-0001": {ID: "T-0001", Type: workitem.Task, Status: workitem.Done},
		"T-0002": {ID: "T-0002", Type: workitem.Task, Status: workitem.InProgress},
	}
	ev, err := ParseEvidence("Verdict: pass\n- 1: T-0001 done, exit 0\n- 2: `T-2` done\n- 3: T-0009 done\n- 4: `a.go`\n- 6: `T-1` done")
	if err != nil {
		t.Fatal(err)
	}
	criteria := []itemedit.Criterion{{N: 1, Text: "one"}, {N: 2, Text: "two"}, {N: 3, Text: "three"}, {N: 4, Text: "four"}, {N: 5, Text: "five"}, {N: 6, Text: "six"}}
	got := stepsUndone(it, criteria, ev, tasks)
	want := []string{
		"criterion 2 (two) names T-0002, which is in-progress, not done",
		"criterion 3 (three) names T-0009, which is not a task of S-0001",
		"criterion 4 (four) names no task's done transition (it says \"`a.go`\")",
		"no item for acceptance criterion 5 (five)",
	}
	if len(got) != len(want) {
		t.Fatalf("blockers = %+v, want %d", got, len(want))
	}
	for i, says := range want {
		if got[i].Code != BlockCriterionUnevidenced || !strings.Contains(got[i].Message, says) {
			t.Errorf("blocker %d = %+v, want %s saying %q", i, got[i], BlockCriterionUnevidenced, says)
		}
	}
}

// S-0621 (ADR-0271 §4): a skipped step's criterion is evidenced by its task
// when the instance's stop cancelled it and the stopping step's task is
// done; any other skipped item is criterion_unevidenced, saying why.
func TestStepsUndoneTakesATaskAStopCancelled(t *testing.T) {
	it := &workitem.Item{ID: "S-0001", Type: workitem.Story, Formula: "publish"}
	stopped := "## Notes\n\nStopped the run with exit 2:\n\n```text\nnot met\n```\n"
	cancelled := "## Notes\n\n- 2026-11-09T10:00:00Z: moved to cancelled: step evaluate stopped the run with exit 2\n"
	tasks := map[string]*workitem.Item{
		"T-0001": {ID: "T-0001", Type: workitem.Task, Status: workitem.Done, Body: stopped},
		"T-0002": {ID: "T-0002", Type: workitem.Task, Status: workitem.Cancelled, Body: cancelled},
		"T-0003": {ID: "T-0003", Type: workitem.Task, Status: workitem.Cancelled, Body: "## Notes\n\n- 2026-11-09T10:00:00Z: moved to cancelled: not needed\n"},
		"T-0004": {ID: "T-0004", Type: workitem.Task, Status: workitem.InProgress},
	}
	ev, err := ParseEvidence("Verdict: pass\n- 1: T-0001 done\n- 2: `T-2` skipped\n- 3: T-0003 skipped\n- 4: T-0004 skipped\n- 5: T-0009 skipped")
	if err != nil {
		t.Fatal(err)
	}
	criteria := []itemedit.Criterion{{N: 1, Text: "one"}, {N: 2, Text: "two"}, {N: 3, Text: "three"}, {N: 4, Text: "four"}, {N: 5, Text: "five"}}
	got := stepsUndone(it, criteria, ev, tasks)
	want := []string{
		"criterion 3 (three) names T-0003 skipped, but T-0003 was not cancelled by a step's stop",
		"criterion 4 (four) names T-0004 skipped, but T-0004 is in-progress, not cancelled",
		"criterion 5 (five) names T-0009 skipped, but T-0009 is not a task of S-0001",
	}
	if len(got) != len(want) {
		t.Fatalf("blockers = %+v, want %d", got, len(want))
	}
	for i, says := range want {
		if got[i].Code != BlockCriterionUnevidenced || !strings.Contains(got[i].Message, says) {
			t.Errorf("blocker %d = %+v, want %s saying %q", i, got[i], BlockCriterionUnevidenced, says)
		}
	}

	// the stopping step's task not done: its stop has not passed
	tasks["T-0001"] = &workitem.Item{ID: "T-0001", Type: workitem.Task, Status: workitem.InProgress, Body: stopped}
	ev, _ = ParseEvidence("Verdict: pass\n- 2: T-0002 skipped")
	got = stepsUndone(it, criteria[1:2], ev, tasks)
	if len(got) != 1 || !strings.Contains(got[0].Message, "no task of S-0001 is done with the record of step evaluate's stop") {
		t.Errorf("a skipped step whose stopping step's task is not done is unevidenced: %+v", got)
	}
}

// ADR-0291, I-0229: a story whose last passing flai verify ran go-test over a
// main that has since reworded flai/internal/roles, which go-test selects, is
// refused with verify_stale by the operator's acceptance, the orchestrator's,
// and a hand acceptance while the queue batches, among the acceptance's
// blockers and not the orchestrator's own; not by the queue's landing, nor
// when main brought in only wip paths, nor when the verify did not pass. A
// main whose paths cannot be listed is the blocker too.
func TestAnAcceptanceWaitsForAVerifyOverMainsNewPaths(t *testing.T) {
	at := time.Date(2026, 11, 12, 16, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s := forgeRepo(t, "local", at)
	f, err := os.OpenFile(filepath.Join(repo.MainRoot, "system-flow.yaml"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("tests:\n  - name: go-test\n    command: [go-test, \"{files}\"]\n    paths: [\"flai/**/*.go\"]\n  - name: markdown\n    command: [lint-md, \"{files}\"]\n    paths: [\"**/*.md\"]\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	diff := "diff --no-renames --name-only " + forgeOther + " main"
	answers := gitAnswers{
		"rev-parse --is-inside-work-tree":                        "true",
		"rev-parse --verify --quiet refs/heads/story/" + s.ID:    forgeHead,
		"rev-parse --verify --quiet story/" + s.ID + "^{commit}": forgeHead,
		"rev-parse --abbrev-ref HEAD":                            "main",
		diff:                                                     "flai/internal/roles/conditions.go\nwip/kanban/stories/S-0099.md",
	}
	git := withGit{answers}
	report := func(passed bool) {
		t.Helper()
		rep := verify.Report{Story: s.ID, Commit: forgeHead, MergeBase: forgeOther, Passed: passed, Paths: []string{"flai/internal/guard/roles_test.go"}, Steps: []verify.Step{
			{Name: "go-test", Tier: true, State: verify.Passed},
			{Name: "markdown", Tier: true, State: verify.Passed},
		}}
		if err := verify.SaveReport(repo, rep); err != nil {
			t.Fatal(err)
		}
	}
	report(true)
	want := s.ID + "'s verify does not cover the main branch it would land on: its last flai verify, at abc1234, ran over main at 0123456, and the 1 path main brought in since send back go-test, which that verify ran: flai/internal/roles/conditions.go; ask for its close-out again"
	for name, opts := range map[string]AcceptOptions{
		"the operator's acceptance":        {},
		"the orchestrator's acceptance":    {Orchestrator: true},
		"a hand acceptance past the queue": {QueueBatches: true, WithoutQueueTiers: true},
	} {
		res, err := AcceptWith(git, repo, s, "alex", at, log, opts)
		if err != nil {
			t.Fatal(err)
		}
		if b := res.VerifyStaleBlocker; b == nil || b.Code != BlockVerifyStale || !strings.HasPrefix(b.Message, want) || !slices.Contains(res.Blockers, b.Message) || !slices.Equal(res.StaleTiers, []string{"go-test"}) {
			t.Errorf("%s must list verify_stale naming go-test:\n%+v %q\n%q", name, res.VerifyStaleBlocker, res.StaleTiers, res.Blockers)
		}
		if slices.ContainsFunc(res.OrchestratorBlockers, func(b Blocker) bool { return b.Code == BlockVerifyStale }) {
			t.Errorf("%s: verify_stale is not among the orchestrator's own blockers, which its enqueue vouches with: %+v", name, res.OrchestratorBlockers)
		}
	}
	stale := func(name string, opts AcceptOptions) {
		t.Helper()
		res, err := AcceptWith(git, repo, s, "alex", at, log, opts)
		if err != nil {
			t.Fatal(err)
		}
		if res.VerifyStaleBlocker != nil || res.StaleTiers != nil || slices.ContainsFunc(res.Blockers, func(b string) bool { return strings.Contains(b, "verify_stale") || strings.Contains(b, "ADR-0291") }) {
			t.Errorf("%s is not refused on verify_stale: %+v %q %q", name, res.VerifyStaleBlocker, res.StaleTiers, res.Blockers)
		}
	}
	stale("the queue's landing", AcceptOptions{QueueBatches: true, Batch: "20261112-160000"})
	answers[diff] = "wip/kanban/stories/S-0099.md"
	stale("a main that brought in only wip paths", AcceptOptions{})
	answers[diff] = "flai/internal/roles/conditions.go"
	report(false)
	stale("a verify that did not pass", AcceptOptions{})

	report(true)
	delete(answers, diff)
	res, err := AcceptWith(git, repo, s, "alex", at, log, AcceptOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if b := res.VerifyStaleBlocker; b == nil || !strings.Contains(b.Message, "cannot list what changed from "+forgeOther+" to main") {
		t.Errorf("a main whose paths cannot be listed is the blocker, naming why: %+v", b)
	}
}
