package workitem

import (
	"math"
	"os"
	"strings"
	"testing"
)

const plannedHead = `---
id: S-0001
type: story
nature: feature
title: S
status: backlog
parent: E-0001
owner: alex
created: 2026-09-15T20:00:00Z
updated: 2026-09-15T20:00:00Z
transitions: []
tags: []
`

// S-0199: each planning field reads and writes back byte for byte, written
// only when set and after every field an older flai knows.
func TestPlanningFieldsRoundTrip(t *testing.T) {
	cases := map[string]string{
		"draft": "draft: true\n",
		"inputs and value": `cost_of_delay:
  inputs:
    revenue_per_week: 1200.5
    penalty_per_week: 0
    time_lost_per_cycle: 2h30m
  value: 1500
  by: planner
  at: 2026-09-16T08:00:00Z
`,
		"value only": `cost_of_delay:
  value: 0
  by: alex
  at: 2026-09-16T08:00:00Z
`,
		"one input": `cost_of_delay:
  inputs:
    penalty_per_week: 300
  by: alex
  at: 2026-09-16T08:00:00Z
`,
		"forecast": `forecast:
  duration: 16h
  delivery: 2026-10-01T17:00:00Z
  basis: "Like S-0143: two days, with one review round."
  by: planner
  at: 2026-09-16T08:00:00Z
`,
		"forecast, delivery only": `forecast:
  delivery: 2026-10-01T17:00:00Z
  basis: the median of the epic's last ten stories
  by: planner
  at: 2026-09-16T08:00:00Z
`,
		"all, after usage": `usage:
  source: sum
  seconds: 1
  models: []
draft: true
cost_of_delay:
  value: 10
  by: alex
  at: 2026-09-16T08:00:00Z
forecast:
  duration: 4h
  by: alex
  at: 2026-09-16T08:00:00Z
`,
	}
	for name, fields := range cases {
		doc := plannedHead + fields + "---\n# S-0001 S\n"
		it, err := ParseItem(doc)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if len(it.Unknown) > 0 {
			t.Errorf("%s: read as unknown: %v", name, FieldNames(it.Unknown))
		}
		if err := it.Validate(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
		if got := it.Marshal(); got != doc {
			t.Errorf("%s: written as\n%s", name, got)
		}
	}
	it, _ := ParseItem(plannedHead + "cost_of_delay:\n  value: 0\n  by: a\n  at: 2026-09-16T08:00:00Z\n---\n")
	if it.CostOfDelay.Value == nil || *it.CostOfDelay.Value != 0 || it.CostOfDelay.Inputs != nil {
		t.Errorf("a value of 0 is not an absent one: %+v", it.CostOfDelay)
	}
	plain, _ := ParseItem(plannedHead + "---\n")
	if plain.Draft || plain.CostOfDelay != nil || plain.Forecast != nil {
		t.Errorf("an item without them: %+v", plain)
	}
	// an empty inputs block and empty keys are not written
	it.CostOfDelay.Inputs = &CostInputs{}
	if got := it.Marshal(); strings.Contains(got, "inputs") {
		t.Errorf("empty inputs written:\n%s", got)
	}
	it.Draft = false
	it.CostOfDelay = &CostOfDelay{}
	it.Forecast = &Forecast{}
	if got := it.Marshal(); strings.Contains(got, "draft") || strings.Contains(got, "cost_of_delay") || strings.Contains(got, "forecast") {
		t.Errorf("empty blocks written:\n%s", got)
	}
}

// S-0199: Validate refuses each field on a type that does not carry it, and
// each bad shape of the blocks.
func TestPlanningFieldsAreValidated(t *testing.T) {
	amount := func(v float64) *float64 { return &v }
	valid := func(typ string) *Item {
		it := &Item{ID: firstID(typ), Type: typ, Nature: "feature", Title: "x", Status: Backlog, Owner: "a",
			Created: "2026-09-15T10:00:00Z", Updated: "2026-09-15T10:00:00Z"}
		if typ == Task {
			it.Parent = "S-0001"
		}
		return it
	}
	cod := func() *CostOfDelay { return &CostOfDelay{Value: amount(5), By: "a", At: "2026-09-15T10:00:00Z"} }
	fc := func() *Forecast { return &Forecast{Duration: "4h", By: "a", At: "2026-09-15T10:00:00Z"} }

	ok := valid(Story)
	ok.Draft, ok.CostOfDelay, ok.Forecast = true, cod(), fc()
	if err := ok.Validate(); err != nil {
		t.Fatalf("a story with every planning field: %v", err)
	}
	epic := valid(Epic)
	epic.CostOfDelay = cod()
	if err := epic.Validate(); err != nil {
		t.Errorf("an epic with a cost of delay: %v", err)
	}

	cases := []struct {
		name, want string
		it         *Item
	}{
		{"draft epic", "only a story is a draft, and this is an epic", func() *Item { it := valid(Epic); it.Draft = true; return it }()},
		{"draft task", "only a story is a draft, and this is a task", func() *Item { it := valid(Task); it.Draft = true; return it }()},
		{"task cost", "cost_of_delay is for stories and epics, and this is a task", func() *Item { it := valid(Task); it.CostOfDelay = cod(); return it }()},
		{"epic forecast", "forecast is for stories, and this is an epic", func() *Item { it := valid(Epic); it.Forecast = fc(); return it }()},
		{"task forecast", "forecast is for stories, and this is a task", func() *Item { it := valid(Task); it.Forecast = fc(); return it }()},
	}
	story := func(edit func(*Item)) *Item {
		it := valid(Story)
		it.CostOfDelay, it.Forecast = cod(), fc()
		edit(it)
		return it
	}
	for _, c := range []struct {
		name, want string
		edit       func(*Item)
	}{
		{"neither inputs nor value", "cost_of_delay has neither inputs nor a value", func(it *Item) { it.CostOfDelay.Value = nil; it.CostOfDelay.Inputs = &CostInputs{} }},
		{"negative value", "cost_of_delay.value -1 is negative", func(it *Item) { it.CostOfDelay.Value = amount(-1) }},
		{"infinite value", "cost_of_delay.value is not a finite amount", func(it *Item) { it.CostOfDelay.Value = amount(math.Inf(1)) }},
		{"NaN revenue", "cost_of_delay.inputs.revenue_per_week is not a finite amount", func(it *Item) { it.CostOfDelay.Inputs = &CostInputs{RevenuePerWeek: amount(math.NaN())} }},
		{"negative penalty", "cost_of_delay.inputs.penalty_per_week -2.5 is negative", func(it *Item) { it.CostOfDelay.Inputs = &CostInputs{PenaltyPerWeek: amount(-2.5)} }},
		{"bad time lost", `cost_of_delay.inputs.time_lost_per_cycle "two hours" is not a Go duration`, func(it *Item) { it.CostOfDelay.Inputs = &CostInputs{TimeLostPerCycle: "two hours"} }},
		{"zero time lost", `cost_of_delay.inputs.time_lost_per_cycle "0s" is not a Go duration longer than zero`, func(it *Item) { it.CostOfDelay.Inputs = &CostInputs{TimeLostPerCycle: "0s"} }},
		{"cost without by", "cost_of_delay.by is required", func(it *Item) { it.CostOfDelay.By = " " }},
		{"cost without at", `cost_of_delay.at "" is not a UTC timestamp`, func(it *Item) { it.CostOfDelay.At = "" }},
		{"cost bad at", `cost_of_delay.at "yesterday" is not a UTC timestamp`, func(it *Item) { it.CostOfDelay.At = "yesterday" }},
		{"forecast empty", "forecast has neither a duration nor a delivery", func(it *Item) { it.Forecast.Duration = "" }},
		{"negative duration", `forecast.duration "-4h" is not a Go duration longer than zero`, func(it *Item) { it.Forecast.Duration = "-4h" }},
		{"bad delivery", `forecast.delivery "2026-10-01" is not a UTC timestamp`, func(it *Item) { it.Forecast.Delivery = "2026-10-01" }},
		{"two-line basis", "forecast.basis is more than one line", func(it *Item) { it.Forecast.Basis = "one.\ntwo." }},
		{"forecast without by", "forecast.by is required", func(it *Item) { it.Forecast.By = "" }},
		{"forecast bad at", `forecast.at "2026-09-15 10:00" is not a UTC timestamp`, func(it *Item) { it.Forecast.At = "2026-09-15 10:00" }},
	} {
		cases = append(cases, struct {
			name, want string
			it         *Item
		}{c.name, c.want, story(c.edit)})
	}
	for _, c := range cases {
		if err := c.it.Validate(); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want %q", c.name, err, c.want)
		}
	}
}

// S-0199: a draft story does not go to ready until it is finalized, and a
// move that finalizes it saves it as no longer a draft.
func TestADraftStoryIsFinalizedToGoToReady(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] works\n", 1)
	s.Draft = true
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	items, _ := r.List(false)
	board, _ := r.LoadBoard()
	want := "rule: S-0001 is a draft: finalize it first (flai edit S-0001 --no-draft, or flai move S-0001 ready --yes)"
	if _, err := r.Move(s, Ready, MoveOptions{Now: t0, Items: items, Board: board}); err == nil || err.Error() != want {
		t.Errorf("a draft to ready: %v", err)
	}
	if s.Status != Backlog || !s.Draft {
		t.Errorf("a refused move changed the story: %s, draft %v", s.Status, s.Draft)
	}
	// a draft may still be cancelled, and comes back a draft
	mustMove(t, r, s, Cancelled, "not now")
	mustMove(t, r, s, Backlog, "")
	if !s.Draft {
		t.Error("a move other than a finalizing one cleared draft")
	}

	if _, err := r.TransitionAll(s, Ready, "alex", "", t0, true); err != nil {
		t.Fatalf("finalize to ready: %v", err)
	}
	back, err := r.Get(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if back.Status != Ready || back.Draft {
		t.Errorf("saved as %s, draft %v", back.Status, back.Draft)
	}
	if data, _ := os.ReadFile(back.Path); strings.Contains(string(data), "draft") {
		t.Errorf("a finalized story still writes draft:\n%s", data)
	}
	// finalize on a story that is not a draft changes nothing else
	mustMove(t, r, back, Backlog, "")
	if _, err := r.Move(back, Ready, MoveOptions{Now: t0, Items: items, Finalize: true}); err != nil || back.Draft {
		t.Errorf("finalize a story that is no draft: %v, draft %v", err, back.Draft)
	}
}

// olderItem is the front matter a flai before S-0199 knew: Item without the
// planning fields.
type olderItem struct {
	ID          string       `yaml:"id"`
	Type        string       `yaml:"type"`
	Nature      string       `yaml:"nature"`
	Title       string       `yaml:"title"`
	Status      string       `yaml:"status"`
	Parent      string       `yaml:"parent"`
	Owner       string       `yaml:"owner"`
	Created     string       `yaml:"created"`
	Updated     string       `yaml:"updated"`
	Transitions []Transition `yaml:"transitions"`
	Blocked     []Block      `yaml:"blocked"`
	Estimate    string       `yaml:"estimate"`
	Stream      string       `yaml:"stream"`
	Tags        []string     `yaml:"tags"`
	Touches     []string     `yaml:"touches"`
	Topics      []string     `yaml:"topics"`
	After       []string     `yaml:"after"`
	Agent       any          `yaml:"agent"`
	Usage       any          `yaml:"usage"`
}

// S-0199: a flai older than the planning fields reads an item carrying them
// past them, as S-0181 reads any field it does not know, and writes them
// back as they were; this flai knows them and keeps none as unknown.
func TestAnOlderFlaiReadsPastThePlanningFields(t *testing.T) {
	planning := "draft: true\ncost_of_delay:\n  inputs:\n    revenue_per_week: 1200\n    time_lost_per_cycle: 4h\n  value: 1500\n  by: alex\n  at: 2026-10-01T09:00:00Z\nforecast:\n  duration: 6h\n  delivery: 2026-10-09T17:00:00Z\n  basis: \"three tasks: like S-0185's\"\n  by: planner\n  at: 2026-10-02T09:00:00Z\n"
	fm := strings.TrimPrefix(plannedHead, "---\n") + planning

	older := UnknownFields(fm, olderItem{})
	if got := strings.Join(FieldNames(older), ","); got != "draft,cost_of_delay,forecast" {
		t.Fatalf("an older flai keeps %q as unknown, want draft,cost_of_delay,forecast", got)
	}
	var b strings.Builder
	WriteFields(&b, older)
	if b.String() != planning {
		t.Errorf("an older flai writes the fields back changed:\n%s\nwant:\n%s", b.String(), planning)
	}

	it, err := ParseItem(plannedHead + planning + "---\n# S-0001 S\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(it.Unknown) != 0 || !it.Draft || it.CostOfDelay == nil || it.Forecast == nil {
		t.Errorf("this flai should know every planning field: unknown %v, item %+v", FieldNames(it.Unknown), it)
	}
}
