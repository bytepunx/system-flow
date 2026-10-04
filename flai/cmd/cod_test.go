package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// S-0210: flai cod reads the items and the manifest's planning settings and
// prints an item's cost of delay per week and its basis, or the whole result
// with --json.
func TestCodCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nplanning:\n  hour_rate: 150\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	epic := "---\nid: E-0001\ntype: epic\nnature: feature\ntitle: E-0001\nstatus: in-progress\nowner: olive\ncreated: 2026-09-01T09:00:00Z\nupdated: 2026-09-01T09:00:00Z\n" +
		"cost_of_delay:\n  inputs:\n    time_lost_per_cycle: 10h\n    by: olive\n    at: 2026-09-01T09:00:00Z\n---\n\n# E-0001\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/epics/E-0001-e-0001.md"), []byte(epic), 0o644); err != nil {
		t.Fatal(err)
	}
	// no done stories, so S-0001 is forecast planning.default_duration, 1h
	forecastStory(t, root, "S-0001", "ready", false, "parent: E-0001\n")
	forecastStory(t, root, "S-0002", "backlog", false, "parent: E-0001\nforecast:\n  duration: 3h\n  by: planner\n  at: 2026-10-04T11:00:00Z\n")
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)

	out, errOut, code := runInAt(t, root, now, "cod", "E-1")
	want := "E-0001 cost of delay 1500.00 USD a week\n10h of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week: 1500.00 USD a week.\n"
	if code != 0 || out != want {
		t.Fatalf("epic: %d %s\n got %q\nwant %q", code, errOut, out, want)
	}

	out, errOut, code = runInAt(t, root, now, "cod", "S-0001")
	want = "S-0001 cost of delay 375.00 USD a week\nS-0001's share of E-0001's 1500.00 USD a week, 1h of 4h forecast over its 2 open stories without inputs: 375.00 USD a week.\n"
	if code != 0 || out != want {
		t.Fatalf("story: %d %s\n got %q\nwant %q", code, errOut, out, want)
	}

	out, errOut, code = runInAt(t, root, now, "cod", "S-0001", "--json")
	if code != 0 {
		t.Fatalf("json: %d %s", code, errOut)
	}
	var res struct {
		ID       string          `json:"id"`
		Value    float64         `json:"value"`
		Currency string          `json:"currency"`
		From     string          `json:"from"`
		Inputs   json.RawMessage `json:"inputs"`
		Epic     struct {
			ID      string  `json:"id"`
			Value   float64 `json:"value"`
			From    string  `json:"from"`
			Share   float64 `json:"share"`
			Stories []struct {
				ID              string `json:"id"`
				Duration        string `json:"duration"`
				DurationSeconds int64  `json:"duration_seconds"`
				FromForecast    bool   `json:"from_forecast"`
			} `json:"stories"`
		} `json:"epic"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if res.ID != "S-0001" || res.Value != 375 || res.Currency != "USD" || res.From != "epic" || res.Inputs != nil ||
		res.Epic.ID != "E-0001" || res.Epic.Value != 1500 || res.Epic.From != "inputs" || res.Epic.Share != 0.25 || len(res.Epic.Stories) != 2 ||
		res.Epic.Stories[0].ID != "S-0001" || res.Epic.Stories[0].DurationSeconds != 3600 || res.Epic.Stories[0].FromForecast ||
		res.Epic.Stories[1].ID != "S-0002" || res.Epic.Stories[1].Duration != "3h" || !res.Epic.Stories[1].FromForecast {
		t.Errorf("json: %s", out)
	}

	if _, errOut, code := runInAt(t, root, now, "cod", "S-0009"); code == 0 || !strings.Contains(errOut, "cannot work out the cost of delay of S-0009: there is no such item") {
		t.Errorf("a missing story: %d %s", code, errOut)
	}
}
