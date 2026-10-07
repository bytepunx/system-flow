package workitem

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/usage"
)

func spent(seconds int64, cost float64, estimated bool) *usage.Usage {
	return &usage.Usage{Source: usage.SourceLog, Seconds: seconds, Estimated: estimated,
		Models: []usage.Model{{Model: "claude-opus-5-5", Input: 10, Output: 20, CacheRead: 1000, CacheWrite: 100, Cost: cost}}}
}

// usageProject is an epic with a story of two tasks in progress, the first
// done and measured, and a second story done and archived with usage of its
// own.
func usageProject(t *testing.T) (*Repo, *Item, *Item) {
	t.Helper()
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	t1 := mustCreate(t, r, Task, "one", s.ID)
	t2 := mustCreate(t, r, Task, "two", s.ID)
	for _, to := range []string{Ready, InProgress} {
		mustMove(t, r, s, to, "")
		mustMove(t, r, t1, to, "")
		mustMove(t, r, t2, to, "")
	}
	t1.Usage = spent(60, 1.25, true)
	mustMove(t, r, t1, Done, "")
	t2.Usage = spent(30, 0.5, true)
	if err := r.Save(t2); err != nil {
		t.Fatal(err)
	}
	old := mustCreate(t, r, Story, "old", "E-0001")
	old.Usage = spent(100, 2, false)
	old.Status, old.Archived, old.Path = Done, true, ""
	old.Transitions = []Transition{{To: Done, At: "2026-09-15T20:00:00Z", By: "test"}}
	if err := os.Remove(itemPath(t, r, old.ID)); err != nil {
		t.Fatal(err)
	}
	if err := r.Save(old); err != nil {
		t.Fatal(err)
	}
	return r, s, t2
}

// itemPath is where the item with id is.
func itemPath(t *testing.T, r *Repo, id string) string {
	t.Helper()
	got, err := r.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	return got.Path
}

func TestATaskEnteringDoneSumsItsStoryAndEpic(t *testing.T) {
	r, s, t2 := usageProject(t)
	res, err := r.TransitionAll(t2, Done, "test", "", t0, false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res.RolledUp, ",") != "S-0001,E-0001" {
		t.Errorf("rolled up %v, want the story and the epic", res.RolledUp)
	}
	story, _ := r.Get(s.ID)
	if u := story.Usage; u == nil || u.Source != usage.SourceSum || u.Seconds != 90 || !u.Estimated || u.Cost() != 1.75 || u.Models[0].Input != 20 {
		t.Fatalf("story usage = %+v, want its tasks' sum", u)
	}
	epic, _ := r.Get("E-0001")
	if u := epic.Usage; u == nil || u.Source != usage.SourceSum || u.Seconds != 190 || u.Cost() != 3.75 {
		t.Fatalf("epic usage = %+v, want its stories' sum, the archived one's included", u)
	}
	// written as front matter, and read back the same
	data, _ := os.ReadFile(epic.Path)
	if !strings.Contains(string(data), "usage:\n  source: sum\n  seconds: 190\n  estimated: true\n  models:\n    - model: claude-opus-5-5\n      input: 30\n") {
		t.Errorf("epic front matter:\n%s", data)
	}
	if err := epic.Validate(); err != nil {
		t.Error(err)
	}
	// nothing changed: a second pass writes nothing
	if again, err := r.RollUp(t2); err != nil || len(again) != 0 {
		t.Errorf("a second roll-up changed %v (%v)", again, err)
	}
}

func TestAStoryMeasuredFromItsLogIsNotSummed(t *testing.T) {
	r, s, t2 := usageProject(t)
	s.Usage = spent(500, 9, false)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	if _, err := r.TransitionAll(t2, Done, "test", "", t0, false); err != nil {
		t.Fatal(err)
	}
	story, _ := r.Get(s.ID)
	if story.Usage.Source != usage.SourceLog || story.Usage.Cost() != 9 {
		t.Errorf("story usage = %+v, want it kept as measured", story.Usage)
	}
	epic, _ := r.Get("E-0001")
	if epic.Usage.Cost() != 11 {
		t.Errorf("epic cost = %v, want the measured story's 9 and the old story's 2", epic.Usage.Cost())
	}
}

func TestARefusedMoveSumsNothing(t *testing.T) {
	r, s, _ := usageProject(t)
	if _, err := r.TransitionAll(s, Done, "test", "", t0, false); err == nil {
		t.Fatal("a story went from in-progress to done")
	}
	epic, _ := r.Get("E-0001")
	if epic.Usage != nil {
		t.Errorf("epic usage = %+v after a refused move", epic.Usage)
	}
}

func TestUsageIsValidated(t *testing.T) {
	it := &Item{Usage: &usage.Usage{Source: "guess", Seconds: -1, Models: []usage.Model{{Model: "m", Cost: -1}, {Model: "m"}, {}}}}
	got := strings.Join(usageErrors(it.Usage), "; ")
	for _, want := range []string{`usage.source "guess"`, "usage.seconds is negative", "usage.models[0] has a negative", "usage.models[1].model m is listed twice", "usage.models[2].model is required"} {
		if !strings.Contains(got, want) {
			t.Errorf("errors %q lack %q", got, want)
		}
	}
	parsed, err := ParseItem(spentDoc)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Usage.Models[0].CacheWrite != 5 || parsed.Marshal() != spentDoc {
		t.Errorf("round trip:\n%s", parsed.Marshal())
	}
}

const spentDoc = `---
id: T-0001
type: task
nature: feature
title: T
status: done
parent: S-0001
owner: alex
created: 2026-09-15T20:00:00Z
updated: 2026-09-15T20:00:00Z
transitions:
  - to: ready
    at: 2026-09-15T20:00:00Z
    by: alex
tags: []
usage:
  source: log
  seconds: 42
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 1
      output: 2
      cache_read: 3
      cache_write: 5
      cost: 0.1234
---
# T-0001 T
`

// S-0272, ADR-0105: a story's empty wakes are written after estimated, read
// back the same, and refused when negative; none is written when there are
// none.
func TestEmptyWakesAreWrittenAndReadBack(t *testing.T) {
	parsed, err := ParseItem(wakesDoc)
	if err != nil {
		t.Fatal(err)
	}
	if err := parsed.Validate(); err != nil {
		t.Error(err)
	}
	if parsed.Usage.EmptyWakes != 14 {
		t.Errorf("empty wakes = %d, want 14", parsed.Usage.EmptyWakes)
	}
	if got := parsed.Marshal(); got != wakesDoc {
		t.Errorf("round trip:\n%s", got)
	}
	parsed.Usage.EmptyWakes = 0
	if got := parsed.Marshal(); strings.Contains(got, "empty_wakes") {
		t.Errorf("no empty wakes written as:\n%s", got)
	}
	parsed.Usage.EmptyWakes = -1
	if got := strings.Join(usageErrors(parsed.Usage), "; "); !strings.Contains(got, "usage.empty_wakes is negative") {
		t.Errorf("errors %q lack the negative empty wakes", got)
	}
}

const wakesDoc = `---
id: S-0001
type: story
nature: feature
title: S
status: ready
parent: E-0001
owner: alex
created: 2026-09-15T20:00:00Z
updated: 2026-09-15T20:00:00Z
transitions:
  - to: ready
    at: 2026-09-15T20:00:00Z
    by: alex
tags: []
usage:
  source: log
  seconds: 1083
  estimated: true
  empty_wakes: 14
  models:
    - model: claude-opus-5-5
      input: 256
      output: 89342
      cache_read: 19723140
      cache_write: 327605
      cost: 8.1258
---
# S-0001 S
`

// S-0272: an epic's usage carries the sum of its stories' empty wakes; a
// story summed from its tasks has none.
func TestAnEpicSumsItsStoriesEmptyWakes(t *testing.T) {
	r, s, t2 := usageProject(t)
	old, err := r.Get("S-0002")
	if err != nil {
		t.Fatal(err)
	}
	old.Usage.EmptyWakes = 5
	if err := r.Save(old); err != nil {
		t.Fatal(err)
	}
	if _, err := r.TransitionAll(t2, Done, "test", "", t0, false); err != nil {
		t.Fatal(err)
	}
	story, _ := r.Get(s.ID)
	if story.Usage.Source != usage.SourceSum || story.Usage.EmptyWakes != 0 {
		t.Errorf("story usage = %+v, want summed from its tasks without empty wakes", story.Usage)
	}
	epic, _ := r.Get("E-0001")
	if epic.Usage.EmptyWakes != 5 {
		t.Errorf("epic empty wakes = %d, want the old story's 5", epic.Usage.EmptyWakes)
	}
	s, _ = r.Get(s.ID)
	s.Usage = spent(500, 9, false)
	s.Usage.EmptyWakes = 2
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	if _, err := r.RollUp(s); err != nil {
		t.Fatal(err)
	}
	epic, _ = r.Get("E-0001")
	if epic.Usage.EmptyWakes != 7 {
		t.Errorf("epic empty wakes = %d, want the measured story's 2 and the old story's 5", epic.Usage.EmptyWakes)
	}
	data, _ := os.ReadFile(epic.Path)
	if !strings.Contains(string(data), "  seconds: 600\n  empty_wakes: 7\n  models:\n") {
		t.Errorf("epic front matter:\n%s", data)
	}
}

// S-0293: a story's turns are written after its empty wakes, each day with
// the classes it has turns of, read back the same, and refused when a day is
// not a date or is listed twice or a count is negative; none are written
// when there are none.
func TestTurnsAreWrittenAndReadBack(t *testing.T) {
	parsed, err := ParseItem(turnsDoc)
	if err != nil {
		t.Fatal(err)
	}
	if err := parsed.Validate(); err != nil {
		t.Error(err)
	}
	want := []usage.TurnDay{{Day: "2026-10-03", Ceremony: 12, TestRuns: 3, Work: 40}, {Day: "2026-10-04", EmptyWakes: 2, HandEdits: 1}}
	if !slices.Equal(parsed.Usage.Turns, want) {
		t.Errorf("turns = %+v, want %+v", parsed.Usage.Turns, want)
	}
	if got := parsed.Marshal(); got != turnsDoc {
		t.Errorf("round trip:\n%s", got)
	}
	parsed.Usage.Turns = []usage.TurnDay{{Day: "4 October", Work: 1}, {Day: "2026-10-04", Ceremony: -1}, {Day: "2026-10-04", Work: 1}}
	got := strings.Join(usageErrors(parsed.Usage), "; ")
	for _, e := range []string{`usage.turns[0].day "4 October" is not a date`, "usage.turns[1].ceremony is negative", "usage.turns[2].day 2026-10-04 is listed twice"} {
		if !strings.Contains(got, e) {
			t.Errorf("errors %q lack %q", got, e)
		}
	}
	parsed.Usage.Turns = nil
	if got := parsed.Marshal(); strings.Contains(got, "turns") {
		t.Errorf("no turns written as:\n%s", got)
	}
}

const turnsDoc = `---
id: S-0001
type: story
nature: feature
title: S
status: ready
parent: E-0001
owner: alex
created: 2026-09-15T20:00:00Z
updated: 2026-09-15T20:00:00Z
transitions:
  - to: ready
    at: 2026-09-15T20:00:00Z
    by: alex
tags: []
usage:
  source: log
  seconds: 1083
  empty_wakes: 2
  turns:
    - day: 2026-10-03
      ceremony: 12
      test_runs: 3
      work: 40
    - day: 2026-10-04
      empty_wakes: 2
      hand_edits: 1
  models:
    - model: claude-opus-5-5
      input: 256
      output: 89342
      cache_read: 19723140
      cache_write: 327605
      cost: 8.1258
---
# S-0001 S
`

// S-0293: an epic's usage carries its stories' turns summed day by day.
func TestAnEpicSumsItsStoriesTurnsByDay(t *testing.T) {
	r, s, _ := usageProject(t)
	old, err := r.Get("S-0002")
	if err != nil {
		t.Fatal(err)
	}
	old.Usage.Turns = []usage.TurnDay{{Day: "2026-10-03", Work: 4}}
	if err := r.Save(old); err != nil {
		t.Fatal(err)
	}
	s.Usage = spent(500, 9, false)
	s.Usage.Turns = []usage.TurnDay{{Day: "2026-10-03", Ceremony: 1}, {Day: "2026-10-04", Work: 2}}
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	if _, err := r.RollUp(s); err != nil {
		t.Fatal(err)
	}
	epic, _ := r.Get("E-0001")
	want := []usage.TurnDay{{Day: "2026-10-03", Ceremony: 1, Work: 4}, {Day: "2026-10-04", Work: 2}}
	if !slices.Equal(epic.Usage.Turns, want) {
		t.Errorf("epic turns = %+v, want %+v", epic.Usage.Turns, want)
	}
}

// S-0225: what strategic agents spent is written after the agents' models
// and read back the same, with or without agents' figures beside it.
func TestStrategicUsageIsWrittenAndReadBack(t *testing.T) {
	for name, doc := range map[string]string{"beside the agents'": strategicDoc, "alone": strategicOnlyDoc} {
		parsed, err := ParseItem(doc)
		if err != nil {
			t.Fatal(err)
		}
		if err := parsed.Validate(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
		if s := parsed.Usage.Strategic; len(s) != 2 || s[0].Kind != ActivityPlanner || s[0].Models[0].CacheWrite != 40210 || s[1].Kind != ActivityAnalyzer {
			t.Errorf("%s: strategic = %+v", name, s)
		}
		if got := parsed.Marshal(); got != doc {
			t.Errorf("%s: round trip:\n%s", name, got)
		}
	}
}

func TestStrategicUsageIsValidated(t *testing.T) {
	u := &usage.Usage{Source: usage.SourceSum, Models: []usage.Model{}, Strategic: []usage.Strategic{
		{Kind: "planner", Seconds: -1, Estimated: true, Models: []usage.Model{{Model: "m", Output: -1}, {Model: "m"}, {}}},
		{Kind: "planner", Estimated: true},
		{Kind: "guesser", Estimated: true},
		{Kind: "analyzer"},
	}}
	got := strings.Join(usageErrors(u), "; ")
	for _, want := range []string{
		"usage.strategic[0].seconds is negative",
		"usage.strategic[0].models[0] has a negative count or cost",
		"usage.strategic[0].models[1].model m is listed twice",
		"usage.strategic[0].models[2].model is required",
		"usage.strategic[1].kind planner is listed twice",
		`usage.strategic[2].kind "guesser" must be one of planner, orchestrator, analyzer`,
		"usage.strategic[3].estimated must be true",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("errors %q lack %q", got, want)
		}
	}
}

// the usage package orders strategic entries by kinds it cannot import
func TestStrategicKindsAreTheActivityKinds(t *testing.T) {
	if strings.Join(usage.StrategicKinds, ",") != strings.Join(ActivityKinds, ",") {
		t.Errorf("usage.StrategicKinds %v, ActivityKinds %v", usage.StrategicKinds, ActivityKinds)
	}
}

func planned(seconds int64, cost float64) *usage.Usage {
	return &usage.Usage{Source: usage.SourceLog, Seconds: seconds, Estimated: true,
		Models: []usage.Model{{Model: "claude-opus-5-5", Input: 1, Output: 200, CacheRead: 5000, CacheWrite: 300, Cost: cost}}}
}

func TestPlanningAStoryChargesItAndItsEpic(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	changed, err := r.ChargeStrategic(s.ID, ActivityPlanner, planned(400, 0.8))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(changed, ",") != "S-0001,E-0001" {
		t.Errorf("changed %v, want the story then its epic", changed)
	}
	if _, err := r.ChargeStrategic(s.ID, ActivityPlanner, planned(12, 0.2)); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"S-0001", "E-0001"} {
		it, _ := r.Get(id)
		u := it.Usage
		if u == nil || u.Source != usage.SourceSum || !u.Empty() || len(u.Strategic) != 1 {
			t.Fatalf("%s usage = %+v, want empty agents' figures and the planner's", id, u)
		}
		if p := u.Strategic[0]; p.Kind != ActivityPlanner || p.Seconds != 412 || !p.Estimated || p.Cost() != 1 || p.Models[0].Output != 400 {
			t.Errorf("%s planner = %+v, want both charges added", id, p)
		}
		if err := it.Validate(); err != nil {
			t.Error(err)
		}
	}
	epic, _ := r.Get("E-0001")
	data, _ := os.ReadFile(epic.Path)
	if !strings.Contains(string(data), "usage:\n  source: sum\n  seconds: 0\n  models: []\n  strategic:\n    - kind: planner\n      seconds: 412\n      estimated: true\n      models:\n        - model: claude-opus-5-5\n          input: 2\n") {
		t.Errorf("epic front matter:\n%s", data)
	}
	if changed, err := r.ChargeStrategic(s.ID, ActivityPlanner, nil); err != nil || changed != nil {
		t.Errorf("charging nothing changed %v (%v)", changed, err)
	}
	if _, err := r.ChargeStrategic(s.ID, "guesser", planned(1, 1)); err == nil {
		t.Error("an unknown kind was charged")
	}
}

func TestPlanningAnEpicChargesItAlone(t *testing.T) {
	r, s, _ := usageProject(t)
	changed, err := r.ChargeStrategic("E-0001", ActivityPlanner, planned(60, 0.5))
	if err != nil || strings.Join(changed, ",") != "E-0001" {
		t.Fatalf("changed %v (%v), want the epic alone", changed, err)
	}
	story, _ := r.Get(s.ID)
	if story.Usage != nil {
		t.Errorf("story usage = %+v", story.Usage)
	}
}

// a roll-up sums the agents' figures again and keeps what was spent
// planning, whether the children spent anything or not
func TestARollUpKeepsWhatPlanningSpent(t *testing.T) {
	r, s, t2 := usageProject(t)
	if _, err := r.ChargeStrategic(s.ID, ActivityPlanner, planned(400, 0.8)); err != nil {
		t.Fatal(err)
	}
	if _, err := r.TransitionAll(t2, Done, "test", "", t0, false); err != nil {
		t.Fatal(err)
	}
	story, _ := r.Get(s.ID)
	epic, _ := r.Get("E-0001")
	if story.Usage.Cost() != 1.75 || story.Usage.StrategicCost() != 0.8 {
		t.Errorf("story usage = %+v, want its tasks' sum and its planning", story.Usage)
	}
	if epic.Usage.Cost() != 3.75 || epic.Usage.StrategicCost() != 0.8 {
		t.Errorf("epic usage = %+v, want its stories' sum and its story's planning", epic.Usage)
	}
	if again, err := r.RollUp(t2); err != nil || len(again) != 0 {
		t.Errorf("a second roll-up changed %v (%v)", again, err)
	}

	// children that spent nothing: the epic keeps its planning
	r = newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s = mustCreate(t, r, Story, "S", "E-0001")
	task := mustCreate(t, r, Task, "T", s.ID)
	if _, err := r.ChargeStrategic(s.ID, ActivityPlanner, planned(10, 0.1)); err != nil {
		t.Fatal(err)
	}
	if changed, err := r.RollUp(task); err != nil || len(changed) != 0 {
		t.Errorf("roll-up changed %v (%v)", changed, err)
	}
	for _, id := range []string{"S-0001", "E-0001"} {
		it, _ := r.Get(id)
		if it.Usage == nil || it.Usage.StrategicSeconds() != 10 {
			t.Errorf("%s usage = %+v, want its planning kept", id, it.Usage)
		}
	}
}

const strategicDoc = `---
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
usage:
  source: sum
  seconds: 42
  models:
    - model: claude-opus-5-5
      input: 1
      output: 2
      cache_read: 3
      cache_write: 5
      cost: 0.1234
  strategic:
    - kind: planner
      seconds: 412
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 3400
          cache_read: 812000
          cache_write: 40210
          cost: 0.8123
    - kind: analyzer
      seconds: 0
      estimated: true
      models: []
---
# S-0001 S
`

const strategicOnlyDoc = `---
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
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 412
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 3400
          cache_read: 812000
          cache_write: 40210
          cost: 0.8123
    - kind: analyzer
      seconds: 0
      estimated: true
      models: []
---
# S-0001 S
`
