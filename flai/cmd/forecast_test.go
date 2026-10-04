package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// forecastStory writes a story of size 4 (two criteria, two touches) with
// the extra front matter given, archived or not.
func forecastStory(t *testing.T, root, id, status string, archived bool, extra string) {
	t.Helper()
	dir := filepath.Join(root, "wip/kanban/stories")
	if archived {
		dir = filepath.Join(root, "wip/archive/kanban/stories")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := fmt.Sprintf("---\nid: %s\ntype: story\nnature: feature\ntitle: %s\nstatus: %s\nowner: olive\ncreated: 2026-09-01T09:00:00Z\nupdated: 2026-09-01T09:00:00Z\ntouches: [a, b]\n%s---\n\n## Acceptance criteria\n- [ ] one\n- [ ] two\n", id, id, status, extra)
	if err := os.WriteFile(filepath.Join(dir, id+"-"+strings.ToLower(id)+".md"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0210: flai forecast reads the done stories, archived ones too, the
// board, and the manifest, and prints the forecast and its basis, or the
// whole result with --json.
func TestForecastCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	// three done stories at 300, 400, and 500 agent seconds per unit, in
	// progress 1.5 times and done twice their agent seconds after starting
	for i, secs := range []int{1200, 1600, 2000} {
		start := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
		at := func(d time.Duration) string { return start.Add(d).Format("2006-01-02T15:04:05Z") }
		s := time.Duration(secs) * time.Second
		forecastStory(t, root, fmt.Sprintf("S-000%d", i+1), "done", i > 0, fmt.Sprintf(
			"transitions:\n  - to: in-progress\n    at: %s\n    by: a\n  - to: review\n    at: %s\n    by: a\n  - to: done\n    at: %s\n    by: olive\nusage:\n  source: log\n  seconds: %d\n  models: []\n",
			at(0), at(s*3/2), at(2*s), secs))
	}
	forecastStory(t, root, "S-0004", "in-progress", false,
		"transitions:\n  - to: in-progress\n    at: 2026-10-04T11:00:00Z\n    by: a\nforecast:\n  duration: 2h\n  by: planner\n  at: 2026-10-04T11:00:00Z\n")
	forecastStory(t, root, "S-0005", "ready", false, "")
	board := "---\ntitle: Board\nstatus: active\nwip_limits:\n  in-progress: 1\norder: [S-0005]\n---\n\n# Board\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)

	out, errOut, code := runInAt(t, root, now, "forecast", "S-0005")
	want := "S-0005 forecast 27m, delivery 2026-10-04T14:54:00Z\n" +
		"Median 400 s per unit of size over 3 done feature stories with no model recorded in the small band, times size 4 (2 criteria, 2 touches); 1st in the pull order with an in-progress limit of 1, behind S-0004.\n"
	if code != 0 || out != want {
		t.Fatalf("text: %d %s\n got %q\nwant %q", code, errOut, out, want)
	}

	out, errOut, code = runInAt(t, root, now, "forecast", "S-0005", "--json")
	var res struct {
		ID              string `json:"id"`
		Duration        string `json:"duration"`
		DurationSeconds int64  `json:"duration_seconds"`
		Delivery        string `json:"delivery"`
		Default         bool   `json:"default"`
		Size            int    `json:"size"`
		History         struct {
			Match   string `json:"match"`
			Stories int    `json:"stories"`
		} `json:"history"`
		Position int `json:"position"`
		Limit    int `json:"limit"`
		Ahead    []struct {
			ID           string `json:"id"`
			FromForecast bool   `json:"from_forecast"`
			Delivery     string `json:"delivery"`
		} `json:"ahead"`
	}
	if code != 0 {
		t.Fatalf("json: %d %s", code, errOut)
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if res.ID != "S-0005" || res.Duration != "27m" || res.DurationSeconds != 1620 || res.Delivery != "2026-10-04T14:54:00Z" || res.Default || res.Size != 4 ||
		res.History.Match != "nature, model, and size band" || res.History.Stories != 3 || res.Position != 1 || res.Limit != 1 ||
		len(res.Ahead) != 1 || res.Ahead[0].ID != "S-0004" || !res.Ahead[0].FromForecast || res.Ahead[0].Delivery != "2026-10-04T15:00:00Z" {
		t.Errorf("json: %s", out)
	}

	if _, errOut, code := runInAt(t, root, now, "forecast", "S-0001"); code == 0 || !strings.Contains(errOut, "cannot forecast S-0001: it is done") {
		t.Errorf("a done story: %d %s", code, errOut)
	}
}
