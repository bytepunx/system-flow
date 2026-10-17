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
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const (
	forgeHead  = "abc1234def5678abc1234def5678abc1234def56"
	forgeOther = "0123456789abcdef0123456789abcdef01234567"
	forgePR    = "https://github.com/o/r/pull/7"
)

// forgeStatus is S-0001's pull request #7 read at forgeHead, with checks,
// reviews, and the approvals its base branch requires.
func forgeStatus(checks []forge.CheckStatus, reviews []forge.ReviewStatus, required int) forge.Status {
	return forge.Status{Story: "S-0001", PullRequest: 7, URL: forgePR, HeadSHA: forgeHead, ReadAt: "2026-10-17T09:00:00Z", Checks: checks, Reviews: reviews, RequiredApprovals: required}
}

func check(name string, state forge.CheckState, summary string) forge.CheckStatus {
	return forge.CheckStatus{Name: name, State: state, URL: "https://ci/" + name, Summary: summary}
}

func review(author string, state forge.ReviewState, n string) forge.ReviewStatus {
	return forge.ReviewStatus{Author: author, State: state, URL: forgePR + "#review-" + n}
}

// ADR-0169: under pull-request acceptance each failing forge condition is a
// blocker with its code, naming the story, the pull request, and the check or
// reviewer, with the page that shows it.
func TestForgeBlockersNameEachFailingConditionOfThePullRequest(t *testing.T) {
	type want struct {
		code, name, url string
		says            []string // what the message carries beyond the story and #7
	}
	cases := []struct {
		name string
		st   forge.Status
		kept bool
		head string
		want []want
	}{
		{name: "no pull request kept", kept: false, head: forgeHead,
			want: []want{{code: BlockPullRequestMissing, says: []string{"kept no pull request", "story/S-0001", "ADR-0166"}}}},
		{name: "stale: only the one blocker, though a check failed then", kept: true, head: forgeOther,
			st:   forgeStatus([]forge.CheckStatus{check("lint", forge.CheckFailed, "2 errors")}, nil, 1),
			want: []want{{code: BlockForgeStale, url: forgePR, says: []string{"#7", "abc1234", "0123456", "2026-10-17T09:00:00Z", "story/S-0001"}}}},
		{name: "failing check with a summary", kept: true, head: forgeHead,
			st:   forgeStatus([]forge.CheckStatus{check("lint", forge.CheckFailed, "2 errors")}, nil, 0),
			want: []want{{code: BlockCheckFailing, name: "lint", url: "https://ci/lint", says: []string{"#7", "lint failing at abc1234: 2 errors;", "story/S-0001"}}}},
		{name: "failing check without a summary", kept: true, head: forgeHead,
			st:   forgeStatus([]forge.CheckStatus{check("lint", forge.CheckFailed, "")}, nil, 0),
			want: []want{{code: BlockCheckFailing, name: "lint", url: "https://ci/lint", says: []string{"#7", "lint failing at abc1234;"}}}},
		{name: "pending check", kept: true, head: forgeHead,
			st:   forgeStatus([]forge.CheckStatus{check("e2e", forge.CheckPending, "")}, nil, 0),
			want: []want{{code: BlockCheckPending, name: "e2e", url: "https://ci/e2e", says: []string{"#7", "e2e still pending at abc1234"}}}},
		{name: "passed check", kept: true, head: forgeHead,
			st: forgeStatus([]forge.CheckStatus{check("lint", forge.CheckPassed, "")}, nil, 0)},
		{name: "changes requested that still stand", kept: true, head: forgeHead,
			st:   forgeStatus(nil, []forge.ReviewStatus{review("ann", forge.ReviewChangesRequested, "1"), review("ann", forge.ReviewCommented, "2")}, 0),
			want: []want{{code: BlockChangesRequested, name: "ann", url: forgePR + "#review-1", says: []string{"ann requested changes", "#7", "have ann approve"}}}},
		{name: "changes requested, then approved by the same reviewer", kept: true, head: forgeHead,
			st: forgeStatus(nil, []forge.ReviewStatus{review("ann", forge.ReviewChangesRequested, "1"), review("ann", forge.ReviewApproved, "2")}, 1)},
		{name: "changes requested, then dismissed", kept: true, head: forgeHead,
			st: forgeStatus(nil, []forge.ReviewStatus{review("ann", forge.ReviewChangesRequested, "1"), review("ann", forge.ReviewDismissed, "2")}, 0)},
		{name: "one approval of two", kept: true, head: forgeHead,
			st:   forgeStatus(nil, []forge.ReviewStatus{review("ann", forge.ReviewApproved, "1"), review("bob", forge.ReviewCommented, "2")}, 2),
			want: []want{{code: BlockApprovalsMissing, url: forgePR, says: []string{"#7", "has 1 approving review of the 2 its base branch requires"}}}},
		{name: "no approval of one", kept: true, head: forgeHead,
			st:   forgeStatus(nil, nil, 1),
			want: []want{{code: BlockApprovalsMissing, url: forgePR, says: []string{"has 0 approving reviews of the 1"}}}},
		{name: "two approvals of two", kept: true, head: forgeHead,
			st: forgeStatus(nil, []forge.ReviewStatus{review("ann", forge.ReviewApproved, "1"), review("bob", forge.ReviewApproved, "2")}, 2)},
		{name: "no approvals required and no reviews", kept: true, head: forgeHead,
			st: forgeStatus(nil, nil, 0)},
		{name: "everything passes", kept: true, head: forgeHead,
			st: forgeStatus([]forge.CheckStatus{check("lint", forge.CheckPassed, ""), check("e2e", forge.CheckPassed, "")}, []forge.ReviewStatus{review("ann", forge.ReviewApproved, "1")}, 1)},
		{name: "an unknown head skips the stale check", kept: true, head: "",
			st:   forgeStatus([]forge.CheckStatus{check("lint", forge.CheckFailed, "")}, nil, 0),
			want: []want{{code: BlockCheckFailing, name: "lint", url: "https://ci/lint", says: []string{"lint failing at abc1234"}}}},
		{name: "every condition at once, in order", kept: true, head: forgeHead,
			st: forgeStatus([]forge.CheckStatus{check("lint", forge.CheckFailed, ""), check("unit", forge.CheckPassed, ""), check("e2e", forge.CheckPending, "")},
				[]forge.ReviewStatus{review("ann", forge.ReviewChangesRequested, "1"), review("bob", forge.ReviewChangesRequested, "2")}, 1),
			want: []want{
				{code: BlockCheckFailing, name: "lint", url: "https://ci/lint"},
				{code: BlockCheckPending, name: "e2e", url: "https://ci/e2e"},
				{code: BlockChangesRequested, name: "ann", url: forgePR + "#review-1"},
				{code: BlockChangesRequested, name: "bob", url: forgePR + "#review-2"},
				{code: BlockApprovalsMissing, url: forgePR},
			}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ForgeBlockers("S-0001", c.st, c.kept, c.head)
			if len(got) != len(c.want) {
				t.Fatalf("%d blockers, want %d: %+v", len(got), len(c.want), got)
			}
			for i, w := range c.want {
				b := got[i]
				if b.Code != w.code || b.Name != w.name || b.URL != w.url {
					t.Errorf("blocker %d = %+v, want code %q name %q url %q", i, b, w.code, w.name, w.url)
				}
				says := append([]string{"S-0001", "ADR-0169)"}, w.says...)
				if c.kept {
					says = append(says, "pull request #7")
				}
				if w.name != "" {
					says = append(says, w.name)
				}
				for _, s := range says {
					if !strings.Contains(b.Message, s) {
						t.Errorf("blocker %d does not say %q: %s", i, s, b.Message)
					}
				}
			}
		})
	}
}

// noGit is a runner on a host without git: the project is not a repository.
type noGit struct{}

func (noGit) Run(string, string, ...string) (string, error) {
	return "", errors.New("no git")
}

func (noGit) RunInput(string, string, string, ...string) (string, error) {
	return "", errors.New("no git")
}

func (noGit) LookPath(string) (string, error) { return "", errors.New("no git") }

// forgeRepo is a project, not a git repository, accepting by mode, with
// S-0001 in review.
func forgeRepo(t *testing.T, mode string, at time.Time) (*workitem.Repo, *workitem.Item) {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nacceptance:\n  mode: " + mode + "\n"
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
	s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Gate", Nature: "feature", Owner: "alex", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [x] gated\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if s, err = repo.Get(s.ID); err != nil {
		t.Fatal(err)
	}
	for _, st := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.Transition(s, st, "alex", "", at); err != nil {
			t.Fatal(err)
		}
	}
	task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Build", Parent: s.ID, Owner: "alex", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	for _, st := range []string{workitem.Ready, workitem.InProgress, workitem.Done} {
		if _, err := repo.Transition(task, st, "alex", "", at); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := repo.Transition(s, workitem.Review, "alex", "", at); err != nil {
		t.Fatal(err)
	}
	if s, err = repo.Get(s.ID); err != nil {
		t.Fatal(err)
	}
	return repo, s
}

// ADR-0169: a pull-request project's story in review has the forge blockers
// in its preview, in forge_blockers and in blockers, and the pull request's
// URL, whoever accepts; a local project's has none of it.
func TestAcceptWithListsTheForgeBlockersUnderPullRequestAcceptanceOnly(t *testing.T) {
	at := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	st := forgeStatus([]forge.CheckStatus{check("lint", forge.CheckFailed, "2 errors")}, nil, 1)

	repo, s := forgeRepo(t, "pull-request", at)
	if err := forge.WriteStatus(repo.CacheDir(), st); err != nil {
		t.Fatal(err)
	}
	for _, opts := range []AcceptOptions{{}, {Orchestrator: true}} {
		res, err := AcceptWith(noGit{}, repo, s, "alex", at, log, opts)
		if err != nil {
			t.Fatal(err)
		}
		want := ForgeBlockers(s.ID, st, true, "")
		if len(res.ForgeBlockers) != 2 || res.ForgeBlockers[0].Code != BlockCheckFailing || res.ForgeBlockers[1].Code != BlockApprovalsMissing {
			t.Errorf("orchestrator %v: forge blockers = %+v, want %+v", opts.Orchestrator, res.ForgeBlockers, want)
		}
		for _, b := range want {
			if !slices.Contains(res.Blockers, b.Message) {
				t.Errorf("orchestrator %v: blockers do not carry %q: %v", opts.Orchestrator, b.Message, res.Blockers)
			}
		}
		if res.PullRequest != forgePR {
			t.Errorf("orchestrator %v: pull request = %q, want %q", opts.Orchestrator, res.PullRequest, forgePR)
		}
	}

	// with no state kept, the pull request is missing
	missing, s2 := forgeRepo(t, "pull-request", at)
	res, err := AcceptWith(noGit{}, missing, s2, "alex", at, log, AcceptOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.ForgeBlockers) != 1 || res.ForgeBlockers[0].Code != BlockPullRequestMissing || !slices.Contains(res.Blockers, res.ForgeBlockers[0].Message) || res.PullRequest != "" {
		t.Errorf("no kept state: %+v", res)
	}

	local, s3 := forgeRepo(t, "local", at)
	if err := forge.WriteStatus(local.CacheDir(), st); err != nil {
		t.Fatal(err)
	}
	res, err = AcceptWith(noGit{}, local, s3, "alex", at, log, AcceptOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.ForgeBlockers) != 0 || res.PullRequest != "" {
		t.Errorf("a local project's preview has no forge blockers: %+v", res)
	}
	for _, b := range res.Blockers {
		if strings.Contains(b, "ADR-0169") {
			t.Errorf("a local project's blockers carry a forge blocker: %s", b)
		}
	}
}
