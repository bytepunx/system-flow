package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// planningProject is a project with an epic and a story under it whose
// acceptance criteria are written, so that it may go to ready.
func planningProject(t *testing.T, storyArgs ...string) (root, story string) {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root = tempProject(t)
	for _, args := range [][]string{{"config", "set", "author", "olive"}, {"epic", "new", "Epic"}, append([]string{"story", "new", "Planned", "--epic", "E-0001"}, storyArgs...)} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	story = filepath.Join(root, "wip/kanban/stories/S-0001-planned.md")
	s, _ := os.ReadFile(story)
	_ = os.WriteFile(story, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	return root, story
}

// S-0199: each cost of delay and forecast flag sets its key and an empty
// value removes it; the block records who changed it and when.
func TestEditPlanningFlags(t *testing.T) {
	root, story := planningProject(t)
	out, errOut, code := runIn(t, root, "edit", "S-0001", "--revenue-per-week", "1200", "--penalty-per-week", "50", "--time-lost-per-cycle", "4h", "--cost-of-delay-value", "1400",
		"--forecast-duration", "6h", "--forecast-delivery", "2026-10-09T17:00:00Z", "--forecast-basis", "like S-0185", "--by", "planner")
	if code != 0 || !strings.Contains(out, "changed cost_of_delay, forecast") {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	item := read(t, story)
	for _, want := range []string{
		"cost_of_delay:\n  inputs:\n    revenue_per_week: 1200\n    penalty_per_week: 50\n    time_lost_per_cycle: 4h\n  value: 1400\n  by: planner\n  at: 2026-09-15T21:00:00Z\n",
		"forecast:\n  duration: 6h\n  delivery: 2026-10-09T17:00:00Z\n  basis: like S-0185\n  by: planner\n  at: 2026-09-15T21:00:00Z\n",
	} {
		if !strings.Contains(item, want) {
			t.Errorf("the item lacks %q:\n%s", want, item)
		}
	}
	for flag, gone := range map[string]string{
		"--revenue-per-week": "revenue_per_week", "--penalty-per-week": "penalty_per_week", "--time-lost-per-cycle": "time_lost_per_cycle",
		"--forecast-delivery": "delivery", "--forecast-basis": "basis",
	} {
		if out, errOut, code := runIn(t, root, "edit", "S-0001", flag, ""); code != 0 || strings.Contains(read(t, story), gone+":") {
			t.Errorf("%s \"\": %d %s %s\n%s", flag, code, out, errOut, read(t, story))
		}
	}
	item = read(t, story)
	if !strings.Contains(item, "cost_of_delay:\n  value: 1400\n  by: olive\n") || !strings.Contains(item, "forecast:\n  duration: 6h\n  by: olive\n") {
		t.Errorf("what is left, by the config author:\n%s", item)
	}
	for _, args := range [][]string{{"--cost-of-delay-value", ""}, {"--forecast-duration", ""}} {
		if _, errOut, code := runIn(t, root, append([]string{"edit", "S-0001"}, args...)...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	if item := read(t, story); strings.Contains(item, "cost_of_delay") || strings.Contains(item, "forecast") {
		t.Errorf("the last key removes the block:\n%s", item)
	}
	for _, args := range [][]string{{"--cost-of-delay-value", "9"}, {"--forecast-duration", "1h"}} {
		if _, errOut, code := runIn(t, root, append([]string{"edit", "S-0001"}, args...)...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	if out, errOut, code := runIn(t, root, "edit", "S-0001", "--clear-cost-of-delay", "--clear-forecast"); code != 0 || !strings.Contains(out, "changed cost_of_delay, forecast") {
		t.Fatalf("clear: %d %s %s", code, out, errOut)
	}
	if item := read(t, story); strings.Contains(item, "cost_of_delay") || strings.Contains(item, "forecast") {
		t.Errorf("cleared:\n%s", item)
	}

	before := read(t, story)
	for _, c := range []struct {
		args []string
		says string
	}{
		{[]string{"--revenue-per-week", "lots"}, "is not a number: write an amount in USD"},
		{[]string{"--cost-of-delay-value", "-1"}, "is negative"},
		{[]string{"--time-lost-per-cycle", "a day"}, "is not a Go duration"},
		{[]string{"--forecast-delivery", "friday"}, "is not a UTC timestamp"},
		{[]string{"--forecast-basis", "a hunch"}, "neither a duration nor a delivery"},
		{[]string{"--draft", "--no-draft"}, "--draft and --no-draft contradict each other"},
		{[]string{"--clear-cost-of-delay", "--cost-of-delay-value", "3"}, "contradict each other"},
		{[]string{"--clear-forecast", "--forecast-duration", "3h"}, "contradict each other"},
	} {
		if _, errOut, code := runIn(t, root, append([]string{"edit", "S-0001"}, c.args...)...); code == 0 || !strings.Contains(errOut, c.says) {
			t.Errorf("%v: %d %s", c.args, code, errOut)
		}
	}
	if _, errOut, _ := runIn(t, root, "edit", "S-0001", "--revenue-per-week", "lots"); !strings.Contains(errOut, "rule: ") {
		t.Errorf("a bad amount is said as a rule: %s", errOut)
	}
	if _, errOut, code := runIn(t, root, "edit", "E-0001", "--forecast-duration", "1h"); code == 0 || !strings.Contains(errOut, "a forecast belongs to stories") {
		t.Errorf("an epic's forecast: %d %s", code, errOut)
	}
	if read(t, story) != before {
		t.Error("a refusal wrote the story")
	}
}

// S-0199: flai story new --draft makes a draft, which flai move refuses to
// ready until --yes finalizes it with the move; flai edit --no-draft and
// --draft finalize and draft it again; flai epic new has no --draft.
func TestStoryNewDraftAndMoveFinalizes(t *testing.T) {
	root, story := planningProject(t, "--draft")
	if item := read(t, story); !strings.Contains(item, "\ndraft: true\n") {
		t.Fatalf("the story is a draft:\n%s", item)
	}
	if _, errOut, code := runIn(t, root, "epic", "new", "Drafty", "--draft"); code == 0 || !strings.Contains(errOut, "unknown flag: --draft") {
		t.Errorf("an epic draft: %d %s", code, errOut)
	}
	if out, _, _ := runIn(t, root, "show", "S-0001"); !strings.Contains(out, "  draft: yes") {
		t.Errorf("show:\n%s", out)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "ready"); code == 0 || !strings.Contains(errOut, "S-0001 is a draft: finalize it first") {
		t.Fatalf("a draft to ready: %d %s", code, errOut)
	}
	if item := read(t, story); !strings.Contains(item, "status: backlog") {
		t.Fatalf("the refused move changed the story:\n%s", item)
	}
	out, errOut, code := runIn(t, root, "move", "S-0001", "ready", "--yes")
	if code != 0 || !strings.Contains(out, "S-0001 is finalized") {
		t.Fatalf("--yes: %d %s %s", code, out, errOut)
	}
	if item := read(t, story); strings.Contains(item, "draft:") || !strings.Contains(item, "status: ready") {
		t.Errorf("finalized and ready:\n%s", item)
	}

	// back to the backlog, a draft again, and finalized by an edit
	if _, errOut, code := runIn(t, root, "edit", "S-0001", "--draft"); code == 0 || !strings.Contains(errOut, "only a story in the backlog is a draft") {
		t.Errorf("a ready story made a draft: %d %s", code, errOut)
	}
	for _, args := range [][]string{{"move", "S-0001", "backlog"}, {"edit", "S-0001", "--draft"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	if out, _, _ := runIn(t, root, "edit", "S-0001", "--show"); !strings.Contains(out, "  draft: yes") {
		t.Errorf("edit --show:\n%s", out)
	}
	if out, errOut, code := runIn(t, root, "edit", "S-0001", "--no-draft"); code != 0 || !strings.Contains(out, "changed draft") {
		t.Fatalf("--no-draft: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, root, "move", "S-0001", "ready"); code != 0 || strings.Contains(out, "finalized") {
		t.Errorf("a finalized story moves without --yes: %d %s %s", code, out, errOut)
	}
}

// S-0199: flai show prints the cost of delay in the project's currency and
// the forecast, with who set each.
func TestShowPrintsPlanning(t *testing.T) {
	root, _ := planningProject(t)
	cfg := filepath.Join(root, "system-flow.yaml")
	m, _ := os.ReadFile(cfg)
	_ = os.WriteFile(cfg, append(m, []byte("planning:\n  currency: EUR\n")...), 0o644)
	if _, errOut, code := runIn(t, root, "edit", "S-0001", "--revenue-per-week", "1200", "--cost-of-delay-value", "1500", "--forecast-duration", "6h", "--forecast-basis", "like S-0185", "--by", "planner"); code != 0 {
		t.Fatal(errOut)
	}
	out, _, _ := runIn(t, root, "show", "S-0001")
	for _, want := range []string{
		"  cost of delay: value 1500 EUR/week · revenue 1200 EUR/week · set by planner at 2026-09-15T21:00:00Z\n",
		"  forecast: duration 6h · basis: like S-0185 · set by planner at 2026-09-15T21:00:00Z\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("show lacks %q:\n%s", want, out)
		}
	}
	if out, _, _ := runIn(t, root, "edit", "S-0001", "--show"); !strings.Contains(out, "  cost of delay: value 1500 EUR/week") {
		t.Errorf("edit --show:\n%s", out)
	}
}
