package inbox

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// catchUp is CatchUp for agent at the fixture's clock.
func (f *fixture) catchUp(t *testing.T, agent string) ([]Event, int) {
	t.Helper()
	items, err := f.repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	events, omitted, err := CatchUp(f.repo, agent, "", f.clock, items)
	if err != nil {
		t.Fatal(err)
	}
	return events, omitted
}

// appendAt appends e stamped at.
func (f *fixture) appendAt(t *testing.T, e LoggedEvent, at time.Time) {
	t.Helper()
	e.At = at.Format(workitem.TimeFormat)
	if err := AppendEvent(f.repo, e); err != nil {
		t.Fatal(err)
	}
}

// ADR-0211 §4: each kind appended to the log is reported once through
// CatchUp, with its fields and in words, to every reader whoever caused it,
// and not again once the cursor has passed it.
func TestALoggedEventIsReportedOnceWithItsFields(t *testing.T) {
	f := setup(t)
	f.catchUp(t, "orchestrator") // the first look sets the cursor
	f.clock = t0.Add(10 * time.Minute)
	at := t0.Add(5 * time.Minute)
	f.appendAt(t, PlannedEvent(f.story.ID, workitem.Story, "Story", "planner-"+f.story.ID, []string{"T-0001", "T-0002"}), at)
	f.appendAt(t, ReadyEmptyEvent(1, 3, 5, "flai serve"), at.Add(time.Second))
	f.appendAt(t, ReleaseDueEvent(f.story.ID, "Story", "orchestrator", "not met: 1 of 3 stories", "met: 3 of 3 stories"), at.Add(2*time.Second))

	events, omitted := f.catchUp(t, "orchestrator")
	if len(events) != 3 || omitted != 0 {
		t.Fatalf("one event of each kind, the orchestrator's own acceptance too: %+v, %d omitted", events, omitted)
	}
	want := []struct {
		kind, to, summary string
		fields            map[string]string
	}{
		{Planned, "T-0001,T-0002", f.story.ID + " Story was planned by planner-" + f.story.ID + ": wrote T-0001, T-0002. A story's agent may let it pass", nil},
		{ReadyEmpty, "1", "the ready column holds 1 story, below the in-progress room of 3 (limit 5). A story's agent may let it pass", map[string]string{FieldRoom: "3", FieldLimit: "5"}},
		{ReleaseDue, "", f.story.ID + " Story's acceptance by orchestrator changed what flai release --evaluate answers: it answered not met: 1 of 3 stories and now answers met: 3 of 3 stories. A story's agent may let it pass", map[string]string{FieldBefore: "not met: 1 of 3 stories", FieldAfter: "met: 3 of 3 stories"}},
	}
	for i, w := range want {
		e := events[i]
		if e.Kind != w.kind || e.To != w.to || e.Summary != w.summary || fmt.Sprint(e.Fields) != fmt.Sprint(w.fields) {
			t.Errorf("event %d: got %+v, want %+v", i, e, w)
		}
	}
	if again, omitted := f.catchUp(t, "orchestrator"); len(again) != 0 || omitted != 0 {
		t.Errorf("told once: %+v, %d omitted", again, omitted)
	}
}

// A logged event stamped with the second the cursor stands at, after the
// reader looked in it, is news; the same line is not told again.
func TestALoggedEventInTheCursorsSecondIsToldOnce(t *testing.T) {
	f := setup(t)
	f.catchUp(t, "claude")
	f.clock = t0.Add(10 * time.Minute)
	f.catchUp(t, "claude") // the cursor stands at this second
	f.appendAt(t, ReadyEmptyEvent(0, 2, 2, "flai serve"), f.clock)
	if events, _ := f.catchUp(t, "claude"); len(events) != 1 || events[0].Kind != ReadyEmpty {
		t.Fatalf("the event in the cursor's second: %+v", events)
	}
	if events, _ := f.catchUp(t, "claude"); len(events) != 0 {
		t.Errorf("told once: %+v", events)
	}
}

// The logged events count under MaxEvents with the rest, the newest kept.
func TestLoggedEventsCountUnderTheCapWithTheOthers(t *testing.T) {
	f := setup(t)
	f.catchUp(t, "claude")
	f.clock = t0.Add(time.Hour)
	for i := range MaxEvents - 5 {
		s := f.storyIn(t, fmt.Sprintf("Story %d", i), []string{"docs"})
		if _, err := f.repo.Transition(s, workitem.Ready, "alex", "", t0.Add(time.Minute+time.Duration(i)*time.Second)); err != nil {
			t.Fatal(err)
		}
	}
	// ten events of the log, all newer than the moves
	for i := range 10 {
		f.appendAt(t, PlannedEvent(f.story.ID, workitem.Story, "Story", "planner", []string{fmt.Sprintf("T-%04d", i+1)}), t0.Add(30*time.Minute+time.Duration(i)*time.Second))
	}
	events, omitted := f.catchUp(t, "claude")
	if len(events) != MaxEvents || omitted != 5 {
		t.Fatalf("%d moves and 10 logged events, capped at %d: %d listed, %d omitted", MaxEvents-5, MaxEvents, len(events), omitted)
	}
	planned := 0
	for _, e := range events {
		if e.Kind == Planned {
			planned++
		}
	}
	if planned != 10 || events[len(events)-1].To != "T-0010" || events[0].Kind != workitem.Moved {
		t.Errorf("the newest kept, every logged event among them: %d planned, first %+v, last %+v", planned, events[0], events[len(events)-1])
	}
}

// Several writers append at once without breaking a line, and a line that
// does not parse, or lacks its time, kind, or ID, is passed over.
func TestTheEventLogKeepsEveryWritersLineAndSkipsWhatItCannotRead(t *testing.T) {
	f := setup(t)
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := AppendEvent(f.repo, PlannedEvent(fmt.Sprintf("S-%04d", i+1), workitem.Story, strings.Repeat("long title ", 100), "planner", []string{"T-0001"})); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	file, err := os.OpenFile(EventsPath(f.repo), os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("not json\n{\"kind\":\"planned\",\"id\":\"S-0099\"}\n"); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	got := LoggedEvents(f.repo)
	if len(got) != 20 {
		t.Fatalf("every writer's line and nothing else: %d", len(got))
	}
	for _, e := range got {
		if _, err := time.Parse(workitem.TimeFormat, e.At); err != nil {
			t.Errorf("an event appended without a time is stamped: %+v", e)
		}
	}
}

// Describe says each logged kind from the change alone.
func TestDescribeSaysALoggedKindWithoutItsFields(t *testing.T) {
	for _, tc := range []struct {
		c    workitem.Change
		want string
	}{
		{workitem.Change{ID: "S-0001", Title: "Plan me", Kind: Planned, By: "planner-S-0001", To: "T-0001"}, "S-0001 Plan me was planned by planner-S-0001: wrote T-0001. A story's agent may let it pass"},
		{workitem.Change{ID: "board", Title: "ready column", Kind: ReadyEmpty, By: "flai serve", To: "0"}, "the ready column holds 0 stories, below the in-progress room. A story's agent may let it pass"},
		{workitem.Change{ID: "S-0002", Title: "Accept me", Kind: ReleaseDue, By: "alex"}, "S-0002 Accept me's acceptance by alex changed what flai release --evaluate answers: run it to see what it answers now. A story's agent may let it pass"},
	} {
		if got := Describe(tc.c); got != tc.want {
			t.Errorf("Describe(%s):\n got %s\nwant %s", tc.c.Kind, got, tc.want)
		}
	}
}
