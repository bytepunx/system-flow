package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const forgePRURL = "https://github.com/o/r/pull/7"

// ADR-0169: under pull-request acceptance, the operator's acceptance and the
// orchestrator's alike are refused, before anything is merged, on each forge
// blocker in the state flai serve kept, the refusal naming it; the dry run
// lists each in forge_blockers with its code, name, and page, and names it
// in its text.
func TestEveryAcceptanceWaitsForThePullRequestsChecksAndReviews(t *testing.T) {
	root, head := forgeGatedStoryInReview(t, "pull-request")
	mainBefore := gitIn(t, root, "rev-parse", "main")
	evidence := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	type blocker struct{ code, name, url string }
	cases := []struct {
		name string
		kept func(st *forge.Status) bool // edits the kept state; false keeps none
		says string
		want []blocker
	}{
		{"a failing required check", func(st *forge.Status) bool {
			st.Checks[0].State, st.Checks[0].Summary = forge.CheckFailed, "2 errors"
			return true
		}, "the required check lint failing", []blocker{{preview.BlockCheckFailing, "lint", "https://ci/lint"}}},
		{"a pending required check", func(st *forge.Status) bool {
			st.Checks = append(st.Checks, forge.CheckStatus{Name: "e2e", State: forge.CheckPending, URL: "https://ci/e2e"})
			return true
		}, "the required check e2e still pending", []blocker{{preview.BlockCheckPending, "e2e", "https://ci/e2e"}}},
		{"a standing requested change", func(st *forge.Status) bool {
			st.Reviews = append(st.Reviews, forge.ReviewStatus{Author: "bob", State: forge.ReviewChangesRequested, URL: forgePRURL + "#review-2"})
			return true
		}, "bob requested changes", []blocker{{preview.BlockChangesRequested, "bob", forgePRURL + "#review-2"}}},
		{"too few approvals", func(st *forge.Status) bool {
			st.RequiredApprovals = 2
			return true
		}, "1 approving review of the 2", []blocker{{preview.BlockApprovalsMissing, "", forgePRURL}}},
		{"a state read at another commit", func(st *forge.Status) bool {
			st.HeadSHA = "0123456789abcdef0123456789abcdef01234567"
			return true
		}, "has not yet read the head", []blocker{{preview.BlockForgeStale, "", forgePRURL}}},
		{"no kept state", func(*forge.Status) bool { return false },
			"no pull request", []blocker{{preview.BlockPullRequestMissing, "", ""}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			st := passingStatus(head)
			path := forge.StatusPath(filepath.Join(root, ".flai-cache"), "S-0001")
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			if c.kept(&st) {
				keepStatus(t, root, st)
			}

			for who, args := range map[string][]string{
				"the operator":     {"accept", "S-0001", "--by", "alex"},
				"the orchestrator": {"accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", evidence},
			} {
				_, errOut, code := runIn(t, root, args...)
				if code == 0 || !strings.Contains(errOut, "S-0001 cannot be accepted yet") || !strings.Contains(errOut, c.says) {
					t.Errorf("%s's acceptance must be refused naming %q: %d\n%s", who, c.says, code, errOut)
				}
				assertNothingAccepted(t, root, mainBefore)
			}

			out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
			if code != 0 {
				t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
			}
			var res preview.Acceptance
			if err := json.Unmarshal([]byte(out), &res); err != nil {
				t.Fatalf("%v: %s", err, out)
			}
			if len(res.ForgeBlockers) != len(c.want) {
				t.Fatalf("forge_blockers %+v, want %+v", res.ForgeBlockers, c.want)
			}
			for i, w := range c.want {
				b := res.ForgeBlockers[i]
				if b.Code != w.code || b.Name != w.name || b.URL != w.url || !strings.Contains(b.Message, c.says) {
					t.Errorf("forge blocker %d = %+v, want %+v saying %q", i, b, w, c.says)
				}
				if !strings.Contains(strings.Join(res.Blockers, "\n"), b.Message) {
					t.Errorf("each forge blocker is in blockers too: %q not in %v", b.Message, res.Blockers)
				}
			}
			wantPR := forgePRURL
			if c.want[0].code == preview.BlockPullRequestMissing {
				wantPR = ""
			}
			if res.PullRequest != wantPR {
				t.Errorf("pull_request %q, want %q", res.PullRequest, wantPR)
			}
			if text, _, _ := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run"); !strings.Contains(text, "blocked: ") || !strings.Contains(text, c.says) {
				t.Errorf("the dry run's text must name the forge blocker %q:\n%s", c.says, text)
			}
		})
	}
}

// ADR-0169, ADR-0170: a story whose kept state was read at its branch's
// head, with its required checks passed and its approvals met, is accepted,
// its pull request merged through the remote, by the operator and by the
// orchestrator alike.
func TestAStoryThePullRequestPassesIsAcceptedThroughTheRemote(t *testing.T) {
	for who, args := range map[string]func(head string) []string{
		"the operator": func(string) []string { return []string{"accept", "S-0001", "--by", "alex", "--json"} },
		"the orchestrator": func(head string) []string {
			return []string{"accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n"), "--json"}
		},
	} {
		t.Run(who, func(t *testing.T) {
			root, _, head, f := remoteProject(t, "")

			out, errOut, code := runIn(t, root, args(head)...)
			if code != 0 {
				t.Fatalf("%s's acceptance of a story its pull request passes: %d %s\n%s", who, code, errOut, out)
			}
			var res preview.Acceptance
			if err := json.Unmarshal([]byte(out), &res); err != nil {
				t.Fatalf("%v: %s", err, out)
			}
			if res.Status != workitem.Done || !res.Merged || len(res.ForgeBlockers) != 0 || res.PullRequest != forgePRURL {
				t.Errorf("the acceptance must be done and merged, naming the pull request: %+v", res)
			}
			if got := f.asked(); len(got) != 1 {
				t.Errorf("the pull request is merged on the forge once: %v", got)
			}
			if !strings.Contains(gitIn(t, root, "show", "main:cli/feature.go"), "package cli") {
				t.Error("the main branch holds the story's change, merged on the forge")
			}
		})
	}
}

// ADR-0169: under local acceptance the kept state gates nothing: a failing
// check is no blocker, and the story is accepted as before.
func TestUnderLocalAcceptanceTheKeptStateGatesNothing(t *testing.T) {
	root, head := forgeGatedStoryInReview(t, "local")
	st := passingStatus(head)
	st.Checks[0].State = forge.CheckFailed
	st.RequiredApprovals = 3
	keepStatus(t, root, st)

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("dry run: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if len(res.ForgeBlockers) != 0 || len(res.Blockers) != 0 || res.PullRequest != "" {
		t.Errorf("a local project's preview has no forge blocker and no pull request: %+v", res)
	}
	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("a local project's story is accepted whatever the kept state: %d %s", code, errOut)
	}
	if !strings.Contains(gitIn(t, root, "show", "main:cli/feature.go"), "package cli") {
		t.Error("the story branch is merged")
	}
}

// forgeGatedStoryInReview is orchestratedStoryInReview, with accept_reviews
// on, in a project accepting by mode, committed on main with the story branch
// rebased onto it, and the branch's head.
func forgeGatedStoryInReview(t *testing.T, mode string) (root, head string) {
	t.Helper()
	root, _ = orchestratedStoryInReview(t, true)
	path := filepath.Join(root, "system-flow.yaml")
	m, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(m, "acceptance:\n  mode: "+mode+"\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "commit", "-q", "-am", "chore: accept by "+mode)
	gitIn(t, filepath.Join(root, ".flai-cache", "worktrees", "S-0001"), "rebase", "-q", "main")
	return root, gitIn(t, root, "rev-parse", "story/S-0001")
}

// passingStatus is S-0001's pull request #7 as flai serve kept it, read at
// head, with its one required check passed and the one approval its base
// branch requires.
func passingStatus(head string) forge.Status {
	return forge.Status{
		Story: "S-0001", PullRequest: 7, URL: forgePRURL, HeadSHA: head, ReadAt: "2026-09-15T20:00:00Z",
		Checks:            []forge.CheckStatus{{Name: "lint", State: forge.CheckPassed, URL: "https://ci/lint"}},
		Reviews:           []forge.ReviewStatus{{Author: "ann", State: forge.ReviewApproved, URL: forgePRURL + "#review-1"}},
		RequiredApprovals: 1,
	}
}

// keepStatus writes st where flai serve keeps it in the project at root.
func keepStatus(t *testing.T, root string, st forge.Status) {
	t.Helper()
	if err := forge.WriteStatus(filepath.Join(root, ".flai-cache"), st); err != nil {
		t.Fatal(err)
	}
}
