package preview

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
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
