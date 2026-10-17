package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0225: flai show prices an item's forecast at the project's cost per agent
// hour, as an estimate, and prints what the planner spent on it apart from
// its agents' usage, alone when no agent has worked it (ADR-0083).
func TestShowPrintsExpectedCostAndStrategicUsage(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story := func(title string) *workitem.Item {
		s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Owner: "t", Now: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	// an hour of agent work for $4, so $4.00 per agent hour
	measured := story("Measured")
	measured.Usage = &usage.Usage{Source: usage.SourceLog, Seconds: 3600, Models: []usage.Model{{Model: "claude-opus-5-5", Input: 1000, Cost: 4}}}
	planned := story("Planned")
	planned.Forecast = &workitem.Forecast{Duration: "3h", By: "planner", At: "2026-09-15T21:00:00Z"}
	planned.Usage = &usage.Usage{Source: usage.SourceSum, Models: []usage.Model{}, Strategic: []usage.Strategic{
		{Kind: "planner", Seconds: 412, Estimated: true, Models: []usage.Model{{Model: "claude-opus-5-5", Input: 12000, CacheRead: 800000, Cost: 0.81}}},
	}}
	for _, s := range []*workitem.Item{measured, planned} {
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
	}

	out, errOut, code := runIn(t, root, "show", planned.ID)
	if code != 0 {
		t.Fatal(errOut)
	}
	for _, want := range []string{
		"  expected cost: $12.00 (estimated, from the forecast of 3h at $4.00 per agent hour)\n",
		"  planner, strategic: 812.0K tokens · $0.81 (estimated) · 6m52s, apart from the agents' usage\n",
		"    claude-opus-5-5  input 12.0K",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("show lacks %q:\n%s", want, out)
		}
	}
	// no agent has worked it, so there is no agents' line of nothing
	if strings.Contains(out, "  usage:") {
		t.Errorf("show prints an empty agents' usage:\n%s", out)
	}

	out, _, _ = runIn(t, root, "show", planned.ID, "--json")
	var got struct {
		Item struct {
			Usage struct {
				Strategic []struct {
					Kind string `json:"kind"`
				} `json:"strategic"`
			} `json:"usage"`
		} `json:"item"`
		ExpectedCost *struct {
			Cost      float64 `json:"cost"`
			From      string  `json:"from"`
			Estimated bool    `json:"estimated"`
		} `json:"expected_cost"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if e := got.ExpectedCost; e == nil || e.Cost != 12 || e.From != "forecast" || !e.Estimated {
		t.Errorf("expected_cost: %+v\n%s", e, out)
	}
	if s := got.Item.Usage.Strategic; len(s) != 1 || s[0].Kind != "planner" {
		t.Errorf("usage.strategic: %+v", s)
	}

	// an item with no duration to price has no expected cost
	out, _, _ = runIn(t, root, "show", measured.ID, "--json")
	if strings.Contains(out, "expected_cost") {
		t.Errorf("an unforecast item has an expected cost:\n%s", out)
	}
	if out, _, _ := runIn(t, root, "show", measured.ID); strings.Contains(out, "expected cost") || strings.Contains(out, "strategic") {
		t.Errorf("show:\n%s", out)
	}
}

// S-0225: an estimate is priced when there is no forecast duration, and
// nothing is priced while no story has been measured from its logs.
func TestShowPricesAnEstimateAndNothingWithoutARate(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	// S-001 has an estimate of 20h, and its stories $4.04 per agent hour
	out, errOut, code := runIn(t, "../internal/metrics/testdata/good", "show", "S-001")
	if code != 0 || !strings.Contains(out, "  expected cost: $80.80 (estimated, from the estimate of 20h at $4.04 per agent hour)\n") {
		t.Errorf("show: %d %s\n%s", code, errOut, out)
	}
	root := tempProject(t)
	for _, args := range [][]string{{"story", "new", "Unpriced"}, {"edit", "S-0001", "--forecast-duration", "2h"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	if out, _, _ := runIn(t, root, "show", "S-0001"); strings.Contains(out, "expected cost") {
		t.Errorf("priced with no rate:\n%s", out)
	}
}

// keptPullRequestFixture is a pull request as flai serve last read it, with
// a failed required check, a pending one, two reviews, and an event.
func keptPullRequestFixture(story string) forge.Status {
	return forge.Status{
		Story: story, PullRequest: 7, URL: "https://github.com/o/r/pull/7", HeadSHA: "abc123def4567890", ReadAt: "2026-10-17T16:00:00Z",
		Checks: []forge.CheckStatus{
			{Name: "build", State: forge.CheckFailed, URL: "https://github.com/o/r/runs/1", Summary: "2 tests failed", Since: "2026-10-17T15:58:00Z"},
			{Name: "lint", State: forge.CheckPending, Since: "2026-10-17T15:59:00Z"},
		},
		Reviews: []forge.ReviewStatus{
			{Author: "alice", State: forge.ReviewChangesRequested, URL: "https://github.com/o/r/pull/7#pullrequestreview-1", Body: "rename it", SubmittedAt: "2026-10-17T15:50:00Z"},
			{Author: "bob", State: forge.ReviewApproved},
		},
		Events: []forge.StatusEvent{{At: "2026-10-17T15:58:00Z", Kind: forge.EventCheck, Name: "build", To: "failed", Commit: "abc123def4567890"}},
	}
}

// S-0394 criterion 2: flai show prints a story's pull request, each required
// check's name, state, log link, and summary, and each review's author,
// state, and link, from the kept state; --json gives it as pull_request
// without the events; a story without one, or a task, has none; and a kept
// state that cannot be read is warned of without hiding the item.
func TestShowPrintsTheKeptPullRequest(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, args := range [][]string{{"story", "new", "Reviewed"}, {"story", "new", "Unread"}, {"task", "new", "--story", "S-0001", "Piece"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := forge.WriteStatus(repo.CacheDir(), keptPullRequestFixture("S-0001")); err != nil {
		t.Fatal(err)
	}

	out, errOut, code := runIn(t, root, "show", "S-0001")
	want := "  pull request: #7 https://github.com/o/r/pull/7 at abc123def456, read 2026-10-17T16:00:00Z\n" +
		"    check build: failed https://github.com/o/r/runs/1 (2 tests failed)\n" +
		"    check lint: pending\n" +
		"    review alice: changes-requested https://github.com/o/r/pull/7#pullrequestreview-1\n" +
		"    review bob: approved\n"
	if code != 0 || !strings.Contains(out, want) {
		t.Errorf("show: %d %s\n%s\nwant it to hold:\n%s", code, errOut, out, want)
	}

	out, errOut, code = runIn(t, root, "show", "S-0001", "--json")
	var got struct {
		PullRequest *forge.Status `json:"pull_request"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil || code != 0 {
		t.Fatalf("show --json: %d %v %s\n%s", code, err, errOut, out)
	}
	wantPR := keptPullRequestFixture("S-0001")
	wantPR.Events = nil
	if !reflect.DeepEqual(got.PullRequest, &wantPR) {
		t.Errorf("pull_request = %+v\nwant %+v", got.PullRequest, wantPR)
	}
	if strings.Contains(out, `"events"`) {
		t.Errorf("pull_request carries its events:\n%s", out)
	}

	for _, id := range []string{"S-0002", "T-0001"} {
		if out, _, _ := runIn(t, root, "show", id); strings.Contains(out, "pull request") {
			t.Errorf("show %s has a pull request:\n%s", id, out)
		}
		if out, _, _ := runIn(t, root, "show", id, "--json"); strings.Contains(out, `"pull_request"`) {
			t.Errorf("show %s --json has a pull request:\n%s", id, out)
		}
	}

	if err := os.WriteFile(forge.StatusPath(repo.CacheDir(), "S-0001"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runIn(t, root, "show", "S-0001")
	if code != 0 || !strings.HasPrefix(out, "S-0001 Reviewed\n") || strings.Contains(out, "pull request") || !strings.Contains(errOut, "pull request status not read") || !strings.Contains(errOut, "delete it") {
		t.Errorf("show with a corrupt kept state: %d\nstdout: %s\nstderr: %s", code, out, errOut)
	}
	if out, _, code := runIn(t, root, "show", "S-0001", "--json"); code != 0 || strings.Contains(out, `"pull_request"`) || !strings.Contains(out, `"item"`) {
		t.Errorf("show --json with a corrupt kept state: %d\n%s", code, out)
	}
}
