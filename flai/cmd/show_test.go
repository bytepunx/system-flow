package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
