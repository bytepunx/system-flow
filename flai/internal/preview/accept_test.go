package preview

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
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
