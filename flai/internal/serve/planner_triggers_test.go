//go:build !windows

package serve

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// The planner starts on the trigger entries its role file names, judged as
// the defined roles' are, through the replanner's queue (ADR-0250 §4).

// plannerTriggers writes the project's planner role file: the template's,
// with its trigger list given as entries, one a line.
func (lab *replanLab) plannerTriggers(entries ...string) {
	lab.t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "template", "root", "design", "workflows", "roles", "planner.md"))
	if err != nil {
		lab.t.Fatal(err)
	}
	s := string(data)
	start, end := strings.Index(s, "\ntrigger:\n"), strings.Index(s, "\nitem:")
	if start < 0 || end < start {
		lab.t.Fatalf("the template's planner.md has no trigger list before item:\n%s", s)
	}
	s = s[:start] + "\ntrigger:\n  - " + strings.Join(entries, "\n  - ") + s[end:]
	dir := filepath.Join(lab.root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "planner.md"), []byte(s), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// queued is the queue as items and their triggers, "S-0001: a; b".
func (lab *replanLab) queued() []string {
	var out []string
	for _, q := range lab.r.queue {
		out = append(out, q.item+": "+strings.Join(q.triggers, "; "))
	}
	return out
}

// ADR-0250 §4: an edit starts the planner when the planner's file names
// {event: edited}; with no condition, an edit of a story never planned does
// too, while the planner's own edit does not.
func TestThePlannerStartsOnAnEditItsFileNames(t *testing.T) {
	lab := newReplanLab(t, "")
	lab.plannerTriggers("asked", "{event: edited}")
	ctx := context.Background()
	lab.hold()
	planned, unplanned := lab.ready("Planned"), lab.ready("Unplanned")
	lab.forecast(planned, "2h", "2026-12-01T00:00:00Z")
	lab.r.look(ctx)
	if run := lab.planRun(planned); run != nil || len(lab.r.queue) != 0 {
		t.Fatalf("the planner's own forecast queued the planner: %+v, queue %v", run, lab.queued())
	}
	lab.clock = lab.clock.Add(time.Minute)
	lab.edit(unplanned, "alex", lab.clock, lab.newGoal(unplanned))
	lab.r.look(ctx)
	waitFor(t, "the planner started", func() bool { return lab.planRun(unplanned) != nil })
	if run := lab.planRun(unplanned); run.Trigger != "edited goal by alex" {
		t.Errorf("run = %+v, want started by the edit", run)
	}
	lab.release(unplanned)
}

// ADR-0250 §4: without {event: edited} in the planner's file, an edit that
// stale-figures would have taken starts nothing.
func TestThePlannerIgnoresAnEditItsFileDoesNotName(t *testing.T) {
	lab := newReplanLab(t, "")
	lab.plannerTriggers("asked", "{event: accepted, when: replan-agent}")
	ctx := context.Background()
	id := lab.ready("Planned")
	lab.forecast(id, "2h", "2026-12-01T00:00:00Z")
	lab.r.look(ctx)
	lab.clock = lab.clock.Add(time.Minute)
	lab.edit(id, "alex", lab.clock, lab.newGoal(id))
	lab.r.look(ctx)
	if run := lab.planRun(id); run != nil || len(lab.r.queue) != 0 {
		t.Errorf("an edit the file does not name queued the planner: %+v, queue %v", run, lab.queued())
	}
}

// ADR-0250 §4, §5: under planning.replan: agent, an acceptance queues the
// planner for each story the replay moved when the file names {event:
// accepted, when: replan-agent}; without it the replay still moves the
// forecasts, and no planner is queued.
func TestThePlannerStartsOnAnAcceptanceUnderReplanAgentFromItsFile(t *testing.T) {
	for _, named := range []bool{true, false} {
		t.Run(map[bool]string{true: "named", false: "not named"}[named], func(t *testing.T) {
			lab := newReplanLab(t, "  replan: agent\n")
			entries := []string{"asked", "{event: cancelled, when: replan-agent}"}
			if named {
				entries = append(entries, "{event: accepted, when: replan-agent}")
			}
			lab.plannerTriggers(entries...)
			ctx := context.Background()
			lab.hold()
			ahead, later, behind := lab.ready("Ahead"), lab.ready("Later"), lab.ready("Behind")
			stale := "2026-01-01T00:00:00Z"
			lab.forecast(later, "1h", stale)
			lab.forecast(behind, "1h", stale)
			lab.accept(ahead)
			lab.clock = lab.clock.Add(time.Minute)
			lab.r.look(ctx)
			if f := lab.item(later).Forecast; f.Delivery == stale || f.By != "flai" || len(lab.git.commits()) != 1 {
				t.Fatalf("forecast = %+v, commits %v: want it replayed whatever the file names", f, lab.git.commits())
			}
			trigger := "accepted " + ahead
			if !named {
				if run := lab.planRun(later); run != nil || len(lab.r.queue) != 0 {
					t.Fatalf("queued with no entry for it: %+v, queue %v", run, lab.queued())
				}
				return
			}
			waitFor(t, "the planner started", func() bool { return lab.planRun(later) != nil })
			if run := lab.planRun(later); run.Trigger != trigger {
				t.Errorf("run = %+v, want %q", run, trigger)
			}
			if got := lab.queued(); !slices.Equal(got, []string{behind + ": " + trigger}) {
				t.Errorf("queue = %v, want %s left", got, behind)
			}
			lab.release(later)
			lab.release(behind)
		})
	}
}

// ADR-0250 §4: an acceptance under {event: accepted} with no condition
// queues the planner for the story accepted, which the checks then refuse,
// whatever planning.replan says.
func TestAnAcceptanceWithNoConditionQueuesItsStory(t *testing.T) {
	lab := newReplanLab(t, "")
	lab.plannerTriggers("asked", "{event: accepted}")
	ctx := context.Background()
	done := lab.ready("Done")
	lab.accept(done)
	lab.clock = lab.clock.Add(time.Minute)
	lab.r.look(ctx)
	if !strings.Contains(lab.logText(), `msg="queued planner dropped" component=serve project=t item=`+done+` trigger="accepted `+done+`"`) {
		t.Errorf("the accepted story was not queued, and dropped by the checks:\n%s", lab.logText())
	}
}

// ADR-0250 §4, §5: a reorder queues the planner for each story the replay
// moved, under planning.replan: agent, when the file names {event:
// reordered, when: replan-agent}; without it the replay still moves the
// forecasts, and no planner is queued.
func TestThePlannerStartsOnAReorderFromItsFile(t *testing.T) {
	for _, named := range []bool{true, false} {
		t.Run(map[bool]string{true: "named", false: "not named"}[named], func(t *testing.T) {
			lab := newReplanLab(t, "  replan: agent\n")
			entries := []string{"asked", "{event: accepted, when: replan-agent}"}
			if named {
				entries = append(entries, "{event: reordered, when: replan-agent}")
			}
			lab.plannerTriggers(entries...)
			ctx := context.Background()
			lab.limit(1)
			lab.hold()
			first, second := lab.ready("First"), lab.ready("Second")
			lab.forecast(first, "2h", "2026-01-01T00:00:00Z")
			lab.forecast(second, "3h", "2026-01-01T00:00:00Z")
			lab.board(first, second)
			lab.r.look(ctx)
			lab.board(second, first)
			lab.clock = lab.clock.Add(time.Minute)
			lab.r.look(ctx)
			if got := lab.git.commits(); len(got) != 1 || !strings.Contains(strings.Join(got[0], " "), "-m chore: replan forecasts after reordered --") {
				t.Fatalf("commits = %v, want the replay whatever the file names", got)
			}
			if !named {
				if lab.planRun(first) != nil || lab.planRun(second) != nil || len(lab.r.queue) != 0 {
					t.Fatalf("queued with no entry for it: %v", lab.queued())
				}
				return
			}
			waitFor(t, "a planner started", func() bool { return lab.planRun(first) != nil || lab.planRun(second) != nil })
			started := 0
			for _, id := range []string{first, second} {
				if run := lab.planRun(id); run != nil {
					started++
					if run.Trigger != triggerReordered {
						t.Errorf("%s run = %+v, want started by reordered", id, run)
					}
				}
			}
			for _, q := range lab.r.queue {
				if !slices.Equal(q.triggers, []string{triggerReordered}) {
					t.Errorf("queued %+v, want reordered alone", q)
				}
			}
			if started != 1 {
				t.Errorf("%d planners started, want one at a time", started)
			}
			if lab.item(first).Forecast.By != "flai" {
				t.Errorf("%s was not moved by the replay: %+v", first, lab.item(first).Forecast)
			}
			lab.release(first)
			lab.release(second)
		})
	}
}

// ADR-0250 §4: the planner's schedule entry queues every ready story when
// its setting comes round; a project's file whose trigger list drops it
// queues nothing then.
func TestThePlannerFollowsTheScheduleItsFileNames(t *testing.T) {
	for _, named := range []bool{true, false} {
		t.Run(map[bool]string{true: "named", false: "dropped"}[named], func(t *testing.T) {
			lab := newReplanLab(t, "  schedule: daily\n")
			entries := []string{"asked", "{event: edited, when: stale-figures}", "{event: accepted, when: replan-agent}", "{event: cancelled, when: replan-agent}", "{event: reordered, when: replan-agent}"}
			if named {
				entries = append(entries, "{schedule: planning.schedule, when: each-ready-story}")
			}
			lab.plannerTriggers(entries...)
			ctx := context.Background()
			lab.hold()
			first, second := lab.ready("First"), lab.ready("Second")
			lab.backlog("Waiting", nil)
			lab.r.look(ctx)
			if len(lab.r.queue) != 0 || lab.planRun(first) != nil {
				t.Fatalf("queued before the schedule came round: %v", lab.queued())
			}
			lab.clock = lab.clock.Add(72 * time.Hour)
			lab.r.look(ctx)
			if !named {
				if lab.planRun(first) != nil || lab.planRun(second) != nil || len(lab.r.queue) != 0 {
					t.Fatalf("a file without the schedule entry queued on it: %v", lab.queued())
				}
				return
			}
			waitFor(t, "the first planner started", func() bool { return lab.planRun(first) != nil })
			if run := lab.planRun(first); run.Trigger != "schedule daily" {
				t.Errorf("run = %+v", run)
			}
			if got := lab.queued(); !slices.Equal(got, []string{second + ": schedule daily"}) {
				t.Errorf("queue = %v, want %s alone left", got, second)
			}
			lab.release(first)
			lab.release(second)
		})
	}
}
