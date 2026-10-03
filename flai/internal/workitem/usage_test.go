package workitem

import (
	"os"
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
