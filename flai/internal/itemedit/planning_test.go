package itemedit

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// planningRepo is a project with an epic, a story in the backlog under it,
// and a task under the story.
func planningRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nplanning:\n  currency: EUR\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	for _, o := range []workitem.NewOptions{
		{Type: workitem.Epic, Title: "Epic"},
		{Type: workitem.Story, Title: "Story", Parent: "E-0001"},
		{Type: workitem.Task, Title: "Task", Parent: "S-0001"},
	} {
		o.Owner, o.Now = "alex", now
		if _, err := repo.Create(o); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

func str(s string) *string { return &s }

func edit(t *testing.T, repo *workitem.Repo, id string, ch Change, by string, at time.Time) (*Result, *workitem.Item, error) {
	t.Helper()
	res, err := Apply(repo, nil, id, ch, Options{By: by, NoCommit: true, Now: at})
	it, gerr := repo.Get(id)
	if gerr != nil {
		t.Fatal(gerr)
	}
	return res, it, err
}

// S-0199: each cost of delay key is set and removed on its own, and the
// block records who changed it last and when; removing the last input and
// the value removes the block.
func TestCostOfDelayKeysSetAndClear(t *testing.T) {
	repo := planningRepo(t)
	t1 := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	res, it, err := edit(t, repo, "S-0001", Change{CostOfDelay: &CostOfDelayEdit{RevenuePerWeek: str("1200"), PenaltyPerWeek: str(" 99.5 "), TimeLostPerCycle: str("4h"), Value: str("1500")}}, "alex", t1)
	if err != nil || strings.Join(res.Changed, ",") != "cost_of_delay" {
		t.Fatalf("set: %v %+v", err, res)
	}
	c := it.CostOfDelay
	if c == nil || c.Inputs == nil || *c.Inputs.RevenuePerWeek != 1200 || *c.Inputs.PenaltyPerWeek != 99.5 || c.Inputs.TimeLostPerCycle != "4h" || *c.Value != 1500 || c.By != "alex" || c.At != "2026-10-02T09:00:00Z" {
		t.Fatalf("the cost of delay: %+v %+v", c, c.Inputs)
	}
	data, _ := os.ReadFile(it.Path)
	if !strings.Contains(string(data), "cost_of_delay:\n  inputs:\n    revenue_per_week: 1200\n    penalty_per_week: 99.5\n    time_lost_per_cycle: 4h\n  value: 1500\n  by: alex\n  at: 2026-10-02T09:00:00Z\n") {
		t.Errorf("the front matter:\n%s", data)
	}

	// one key removed: the rest stay, and the block is stamped by whoever removed it
	_, it, err = edit(t, repo, "S-0001", Change{CostOfDelay: &CostOfDelayEdit{PenaltyPerWeek: str("")}}, "", t2)
	if err != nil || it.CostOfDelay.Inputs.PenaltyPerWeek != nil || *it.CostOfDelay.Inputs.RevenuePerWeek != 1200 || it.CostOfDelay.By != "agent" || it.CostOfDelay.At != "2026-10-02T10:00:00Z" {
		t.Fatalf("remove penalty: %v %+v", err, it.CostOfDelay)
	}
	// the same again changes nothing, and the stamp stays
	res, it, err = edit(t, repo, "S-0001", Change{CostOfDelay: &CostOfDelayEdit{RevenuePerWeek: str("1200.0")}}, "bob", t2.Add(time.Hour))
	if err != nil || !res.Unchanged || it.CostOfDelay.By != "agent" {
		t.Fatalf("the same value: %v %+v %+v", err, res, it.CostOfDelay)
	}
	// the inputs go, the value stays
	_, it, err = edit(t, repo, "S-0001", Change{CostOfDelay: &CostOfDelayEdit{RevenuePerWeek: str(""), TimeLostPerCycle: str("")}}, "bob", t2)
	if err != nil || it.CostOfDelay == nil || it.CostOfDelay.Inputs != nil || *it.CostOfDelay.Value != 1500 {
		t.Fatalf("remove the inputs: %v %+v", err, it.CostOfDelay)
	}
	// the last key goes, and the block with it
	_, it, err = edit(t, repo, "S-0001", Change{CostOfDelay: &CostOfDelayEdit{Value: str("")}}, "bob", t2)
	if err != nil || it.CostOfDelay != nil {
		t.Fatalf("remove the value: %v %+v", err, it.CostOfDelay)
	}
	if data, _ := os.ReadFile(it.Path); strings.Contains(string(data), "cost_of_delay") {
		t.Errorf("the block stays:\n%s", data)
	}

	// an epic carries one too, and clearing removes it whole
	if _, it, err = edit(t, repo, "E-0001", Change{CostOfDelay: &CostOfDelayEdit{Value: str("300")}}, "alex", t1); err != nil || *it.CostOfDelay.Value != 300 {
		t.Fatalf("an epic's: %v %+v", err, it.CostOfDelay)
	}
	if res, it, err = edit(t, repo, "E-0001", Change{ClearCostOfDelay: true}, "alex", t1); err != nil || it.CostOfDelay != nil || strings.Join(res.Changed, ",") != "cost_of_delay" {
		t.Fatalf("clear: %v %+v", err, it.CostOfDelay)
	}
}

// S-0199: a forecast's keys are set and removed on their own; a basis alone
// is not a forecast.
func TestForecastKeysSetAndClear(t *testing.T) {
	repo := planningRepo(t)
	t1 := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	res, it, err := edit(t, repo, "S-0001", Change{Forecast: &ForecastEdit{Duration: str("6h"), Delivery: str("2026-10-09T17:00:00Z"), Basis: str("three tasks like S-0185's")}}, "planner", t1)
	if err != nil || strings.Join(res.Changed, ",") != "forecast" {
		t.Fatalf("set: %v %+v", err, res)
	}
	if f := it.Forecast; f == nil || f.Duration != "6h" || f.Delivery != "2026-10-09T17:00:00Z" || f.Basis != "three tasks like S-0185's" || f.By != "planner" || f.At != "2026-10-02T09:00:00Z" {
		t.Fatalf("the forecast: %+v", it.Forecast)
	}
	if _, it, err = edit(t, repo, "S-0001", Change{Forecast: &ForecastEdit{Delivery: str("")}}, "alex", t1); err != nil || it.Forecast.Delivery != "" || it.Forecast.Duration != "6h" || it.Forecast.By != "alex" {
		t.Fatalf("remove the delivery: %v %+v", err, it.Forecast)
	}
	// the last of duration and delivery goes, and the basis with the block
	if _, it, err = edit(t, repo, "S-0001", Change{Forecast: &ForecastEdit{Duration: str("")}}, "alex", t1); err != nil || it.Forecast != nil {
		t.Fatalf("remove the duration: %v %+v", err, it.Forecast)
	}
	if _, _, err = edit(t, repo, "S-0001", Change{Forecast: &ForecastEdit{Basis: str("a hunch")}}, "alex", t1); !isInvalid(err) || !strings.Contains(err.Error(), "neither a duration nor a delivery") {
		t.Errorf("a basis alone: %v", err)
	}
	if _, it, err = edit(t, repo, "S-0001", Change{Forecast: &ForecastEdit{Duration: str("2h")}}, "alex", t1); err != nil || it.Forecast.Duration != "2h" {
		t.Fatalf("set again: %v %+v", err, it.Forecast)
	}
	if _, it, err = edit(t, repo, "S-0001", Change{ClearForecast: true}, "alex", t1); err != nil || it.Forecast != nil {
		t.Fatalf("clear: %v %+v", err, it.Forecast)
	}
}

// S-0199: a backlog story becomes a draft and is finalized; nothing else is
// a draft.
func TestDraftSetAndClear(t *testing.T) {
	repo := planningRepo(t)
	t1 := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	yes, no := true, false
	if res, it, err := edit(t, repo, "S-0001", Change{Draft: &yes}, "alex", t1); err != nil || !it.Draft || strings.Join(res.Changed, ",") != "draft" {
		t.Fatalf("draft: %v %+v", err, res)
	}
	if v, err := Show(repo, "S-0001"); err != nil || !v.Draft || v.Currency != "EUR" {
		t.Errorf("show: %v %+v", err, v)
	}
	if res, it, err := edit(t, repo, "S-0001", Change{Draft: &no}, "alex", t1); err != nil || it.Draft || strings.Join(res.Changed, ",") != "draft" {
		t.Fatalf("finalize: %v %+v", err, res)
	}
	// S-0201: finalizing records who did it and when, as the other blocks do
	t2 := t1.Add(time.Hour)
	if res, _, err := edit(t, repo, "S-0001", Change{Draft: &no}, "bob", t2); err != nil || !res.Unchanged {
		t.Errorf("finalizing a story that is not a draft changes nothing: %v %+v", err, res)
	}
	it, _ := repo.Get("S-0001")
	if f := it.Finalized; f == nil || f.By != "alex" || f.At != "2026-10-02T09:00:00Z" {
		t.Fatalf("finalized: %+v", it.Finalized)
	}
	if data, _ := os.ReadFile(it.Path); !strings.Contains(string(data), "finalized:\n  by: alex\n  at: 2026-10-02T09:00:00Z\n") {
		t.Errorf("the front matter:\n%s", data)
	}
	if v, err := Show(repo, "S-0001"); err != nil || v.Draft || v.Finalized == nil || v.Finalized.By != "alex" {
		t.Errorf("show a finalized story: %v %+v", err, v)
	}
	// a draft again is no longer finalized; finalized again, by whoever does it
	if _, it, err := edit(t, repo, "S-0001", Change{Draft: &yes}, "alex", t2); err != nil || !it.Draft || it.Finalized != nil {
		t.Fatalf("a draft again: %v %+v", err, it.Finalized)
	}
	if _, it, err := edit(t, repo, "S-0001", Change{Draft: &no}, "", t2); err != nil || it.Finalized == nil || it.Finalized.By != "agent" || it.Finalized.At != "2026-10-02T10:00:00Z" {
		t.Fatalf("finalized again: %v %+v", err, it.Finalized)
	}
	if _, _, err := edit(t, repo, "E-0001", Change{Draft: &yes}, "alex", t1); !isInvalid(err) || !strings.Contains(err.Error(), "only a story is a draft") {
		t.Errorf("an epic: %v", err)
	}
	it, _ = repo.Get("S-0001")
	it.Body = strings.Replace(it.Body, "## Acceptance criteria\n- [ ]", "## Acceptance criteria\n- [ ] works", 1)
	if err := repo.Save(it); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Transition(it, workitem.Ready, "alex", "", t1); err != nil {
		t.Fatal(err)
	}
	if _, _, err := edit(t, repo, "S-0001", Change{Draft: &yes}, "alex", t1); !isInvalid(err) || !strings.Contains(err.Error(), "only a story in the backlog is a draft") {
		t.Errorf("a ready story: %v", err)
	}
}

// S-0199: what is not a number, not a duration, or not a timestamp is
// refused as the caller's mistake, and so is a block on a type that does not
// carry it; nothing is written.
func TestPlanningRefusals(t *testing.T) {
	repo := planningRepo(t)
	t1 := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	story, _ := repo.Get("S-0001")
	before, _ := os.ReadFile(story.Path)
	for name, c := range map[string]struct {
		id   string
		ch   Change
		says string
	}{
		"words for an amount": {"S-0001", Change{CostOfDelay: &CostOfDelayEdit{RevenuePerWeek: str("lots")}}, `revenue_per_week "lots" is not a number: write an amount in EUR`},
		"a negative value":    {"S-0001", Change{CostOfDelay: &CostOfDelayEdit{Value: str("-5")}}, "cost_of_delay.value -5 is negative"},
		"a bad time lost":     {"S-0001", Change{CostOfDelay: &CostOfDelayEdit{TimeLostPerCycle: str("a day")}}, "time_lost_per_cycle \"a day\" is not a Go duration"},
		"a bad duration":      {"S-0001", Change{Forecast: &ForecastEdit{Duration: str("soon")}}, "forecast.duration \"soon\" is not a Go duration"},
		"a bad delivery":      {"S-0001", Change{Forecast: &ForecastEdit{Delivery: str("next week")}}, "forecast.delivery \"next week\" is not a UTC timestamp"},
		"a task's cost":       {"T-0001", Change{CostOfDelay: &CostOfDelayEdit{Value: str("1")}}, "a cost of delay belongs to stories and epics"},
		"an epic's forecast":  {"E-0001", Change{Forecast: &ForecastEdit{Duration: str("1h")}}, "a forecast belongs to stories"},
	} {
		if _, _, err := edit(t, repo, c.id, c.ch, "alex", t1); !isInvalid(err) || !strings.Contains(err.Error(), c.says) {
			t.Errorf("%s: %v", name, err)
		}
	}
	if after, _ := os.ReadFile(story.Path); string(after) != string(before) {
		t.Errorf("a refusal wrote the story:\n%s", after)
	}
}

func isInvalid(err error) bool {
	var inv *InvalidError
	return errors.As(err, &inv)
}
