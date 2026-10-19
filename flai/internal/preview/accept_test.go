package preview

import (
	"errors"
	"io"
	"log/slog"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
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
