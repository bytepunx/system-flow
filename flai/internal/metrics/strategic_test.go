package metrics

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0079: each activity document's totals as written, all time, and its
// log entries in the window.
func TestStrategicReportsTotalsAndTheWindowsEntries(t *testing.T) {
	at := func(s string) time.Time {
		v, err := time.Parse(workitem.TimeFormat, s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	docs := []*workitem.Activity{
		{Kind: workitem.ActivityPlanner, AccruedCost: 0.6213, AccruedSeconds: 900, TasksCompleted: 3, LastRun: "2026-08-31T18:00:00Z", Entries: []workitem.ActivityEntry{
			{At: at("2026-07-01T09:00:00Z"), Summary: "Before the window", Items: []string{"E-0001"}, Seconds: 77, Cost: 0.1},
			{At: at("2026-08-02T12:00:00Z"), Summary: "At the window's start", Seconds: 100, Cost: 0.1, Estimated: true},
			{At: at("2026-08-31T18:00:00Z"), Summary: "Drafted stories", Items: []string{"E-0016", "S-0230"}, Seconds: 723, Cost: 0.4213},
		}},
		{Kind: workitem.ActivityAnalyzer, Entries: []workitem.ActivityEntry{}},
	}
	rep := Compute(nil, Options{Now: now, Activities: docs})
	if len(rep.Strategic) != 2 || rep.Strategic[0].Kind != "planner" || rep.Strategic[1].Kind != "analyzer" {
		t.Fatalf("strategic = %+v", rep.Strategic)
	}
	p := rep.Strategic[0]
	if p.Cost != 0.6213 || p.Seconds != 900 || p.Activities != 3 || p.LastRun != "2026-08-31T18:00:00Z" {
		t.Errorf("planner totals are not the document's: %+v", p)
	}
	if len(p.Log) != 2 || p.Log[0].At != "2026-08-02T12:00:00Z" || !p.Log[0].Estimated || p.Log[1].Seconds != 723 || p.Log[1].Cost != 0.4213 || strings.Join(p.Log[1].Items, ",") != "E-0016,S-0230" {
		t.Errorf("planner log in the window = %+v", p.Log)
	}
	data, err := json.Marshal(rep.Strategic)
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"kind":"planner","cost":0.6213,"seconds":900,"activities":3,"last_run":"2026-08-31T18:00:00Z",` +
		`"items":{"cost":0,"seconds":0},"project":{"cost":0.6213,"seconds":900},"log":[` +
		`{"at":"2026-08-02T12:00:00Z","seconds":100,"cost":0.1,"estimated":true,"items":[]},` +
		`{"at":"2026-08-31T18:00:00Z","seconds":723,"cost":0.4213,"items":["E-0016","S-0230"]}]},` +
		`{"kind":"analyzer","cost":0,"seconds":0,"activities":0,"last_run":"","items":{"cost":0,"seconds":0},"project":{"cost":0,"seconds":0},"log":[]}]`
	if string(data) != want {
		t.Errorf("json =\n%s\nwant\n%s", data, want)
	}
}

// ADR-0095: what the items carry of each kind is the sum over the items at
// the top of the hierarchy, archived ones included, since each charge is on
// its item and every item above it; the rest of the document's totals is the
// project strategic total, never below zero.
func TestStrategicSplitsTheTotalsBetweenTheItemsAndTheProject(t *testing.T) {
	charged := func(id, typ, parent string, archived bool, s ...usage.Strategic) *workitem.Item {
		return &workitem.Item{ID: id, Type: typ, Parent: parent, Status: workitem.InProgress, Created: "2026-08-01T00:00:00Z", Archived: archived,
			Usage: &usage.Usage{Source: usage.SourceSum, Models: []usage.Model{}, Strategic: s}}
	}
	use := func(kind string, seconds int64, costs ...float64) usage.Strategic {
		s := usage.Strategic{Kind: kind, Seconds: seconds, Estimated: true}
		for i, c := range costs {
			s.Models = append(s.Models, usage.Model{Model: fmt.Sprintf("m%d", i), Output: 10, Cost: c})
		}
		return s
	}
	items := []*workitem.Item{
		// an orchestrator activity named T-0001 ($0.1, 10 s) and another
		// S-0001 ($0.2, 20 s): each charge is on its item and every item
		// above it
		charged("E-0001", workitem.Epic, "", false, use("orchestrator", 30, 0.25, 0.05)),
		charged("S-0001", workitem.Story, "E-0001", false, use("orchestrator", 30, 0.3)),
		charged("T-0001", workitem.Task, "S-0001", false, use("orchestrator", 10, 0.1)),
		// a story with no epic, archived, planned and decided on
		charged("S-0002", workitem.Story, "", true, use("planner", 40, 0.4), use("orchestrator", 5, 0.05)),
		// a task whose story is not among the items is at the top
		charged("T-0002", workitem.Task, "S-0099", false, use("orchestrator", 1, 0.01)),
		{ID: "S-0003", Type: workitem.Story, Status: workitem.Backlog, Created: "2026-08-01T00:00:00Z"},
	}
	docs := []*workitem.Activity{
		{Kind: workitem.ActivityPlanner, AccruedCost: 0.3, AccruedSeconds: 30},
		{Kind: workitem.ActivityOrchestrator, AccruedCost: 0.5, AccruedSeconds: 50},
		{Kind: workitem.ActivityAnalyzer, AccruedCost: 0.05, AccruedSeconds: 60},
	}
	rep := Compute(items, Options{Now: now, Activities: docs})
	type split struct{ items, project StrategicAmount }
	want := map[string]split{
		// the items carry more than the document holds: the project is none
		"planner":      {StrategicAmount{Cost: 0.4, Seconds: 40}, StrategicAmount{}},
		"orchestrator": {StrategicAmount{Cost: 0.36, Seconds: 36}, StrategicAmount{Cost: 0.14, Seconds: 14}},
		// no charges: all of it is the project's
		"analyzer": {StrategicAmount{}, StrategicAmount{Cost: 0.05, Seconds: 60}},
	}
	if len(rep.Strategic) != len(docs) {
		t.Fatalf("strategic = %+v", rep.Strategic)
	}
	for _, s := range rep.Strategic {
		if w := want[s.Kind]; s.Items != w.items || s.Project != w.project {
			t.Errorf("%s: items %+v project %+v, want %+v %+v", s.Kind, s.Items, s.Project, w.items, w.project)
		}
	}
	if o := rep.Strategic[1]; round4(o.Items.Cost+o.Project.Cost) != o.Cost || o.Items.Seconds+o.Project.Seconds != o.Seconds {
		t.Errorf("the orchestrator's items and project do not add up to its totals: %+v", o)
	}
}

func TestStrategicIsAnEmptyListWithNoDocuments(t *testing.T) {
	data, err := json.Marshal(Compute(nil, Options{Now: now}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"strategic":[]`) {
		t.Errorf("no documents should report an empty list: %s", data)
	}
}

// strategicDayItems are stories done over the window's days with and without
// usage and a cycle time, one cancelled, one done after now, and an epic the
// story report leaves out.
func strategicDayItems() []*workitem.Item {
	item := func(id, typ string, cost float64, created string, moves ...string) *workitem.Item {
		it := &workitem.Item{ID: id, Type: typ, Status: workitem.Backlog, Created: created}
		for i := 0; i < len(moves); i += 2 {
			it.Transitions = append(it.Transitions, workitem.Transition{To: moves[i], At: moves[i+1]})
			it.Status = moves[i]
		}
		if cost > 0 {
			it.Usage = &usage.Usage{Seconds: 60, Models: []usage.Model{{Model: "m", Output: 10, Cost: cost}}}
		}
		return it
	}
	return []*workitem.Item{
		// done on the start's day, before the start: usage and a day's cycle
		item("S-0001", workitem.Story, 0.1, "2026-08-27T00:00:00Z",
			workitem.InProgress, "2026-08-28T08:00:00Z", workitem.Done, "2026-08-29T08:00:00Z"),
		// usage, done straight from backlog: no cycle time
		item("S-0002", workitem.Story, 0.2, "2026-08-27T00:00:00Z", workitem.Done, "2026-08-29T20:00:00Z"),
		item("S-0003", workitem.Story, 0.2, "2026-08-29T00:00:00Z",
			workitem.InProgress, "2026-08-29T13:00:00Z", workitem.Done, "2026-08-29T14:00:00Z"),
		// no usage, a cycle time
		item("S-0004", workitem.Story, 0, "2026-08-29T00:00:00Z",
			workitem.InProgress, "2026-08-29T13:00:00Z", workitem.Done, "2026-08-29T15:00:00Z"),
		// neither usage nor a cycle time
		item("S-0005", workitem.Story, 0, "2026-08-29T00:00:00Z", workitem.Done, "2026-08-30T09:00:00Z"),
		item("S-0006", workitem.Story, 0.5, "2026-08-29T00:00:00Z",
			workitem.InProgress, "2026-08-29T13:00:00Z", workitem.Cancelled, "2026-08-30T10:00:00Z"),
		item("S-0007", workitem.Story, 0.5, "2026-08-29T00:00:00Z", workitem.Done, "2026-09-01T13:00:00Z"),
		item("S-0008", workitem.Story, 0.5, "2026-08-20T00:00:00Z", workitem.Done, "2026-08-28T23:00:00Z"),
		item("E-0001", workitem.Epic, 0.5, "2026-08-20T00:00:00Z", workitem.Done, "2026-08-30T09:00:00Z"),
	}
}

// S-0205: the strategic agents' cost and seconds per kind and in all on each
// day of the window, beside the items done that day.
func TestStrategicDaysLayOutUseBesideDelivery(t *testing.T) {
	at := func(s string) time.Time {
		v, err := time.Parse(workitem.TimeFormat, s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	docs := []*workitem.Activity{
		{Kind: workitem.ActivityPlanner, Entries: []workitem.ActivityEntry{
			{At: at("2026-08-28T23:00:00Z"), Seconds: 500, Cost: 5},
			{At: at("2026-08-29T06:00:00Z"), Seconds: 100, Cost: 0.12345},
			{At: at("2026-08-29T20:00:00Z"), Seconds: 50, Cost: 0.00001, Estimated: true},
		}},
		{Kind: workitem.ActivityOrchestrator, Entries: []workitem.ActivityEntry{
			{At: at("2026-08-29T21:00:00Z"), Seconds: 30, Cost: 0.00004},
			{At: at("2026-08-29T22:00:00Z"), Seconds: 30, Cost: 0.00004},
		}},
		{Kind: workitem.ActivityAnalyzer, Entries: []workitem.ActivityEntry{
			{At: at("2026-08-31T10:00:00Z"), Seconds: 60, Cost: 0.25},
			{At: at("2026-09-01T13:00:00Z"), Seconds: 900, Cost: 9},
		}},
	}
	days := Compute(strategicDayItems(), Options{Now: now, Since: threeDays, Activities: docs}).StrategicDays
	want := []string{
		`{"date":"2026-08-29","agents":{"orchestrator":{"cost":0.0001,"seconds":60},"planner":{"cost":0.1235,"seconds":150,"estimated":true}},` +
			`"cost":0.1235,"seconds":210,"completed":4,"cost_per_item":0.1667,"cycle_time_seconds":32400}`,
		`{"date":"2026-08-30","agents":{},"cost":0,"seconds":0,"completed":1}`,
		`{"date":"2026-08-31","agents":{"analyzer":{"cost":0.25,"seconds":60}},"cost":0.25,"seconds":60,"completed":0}`,
		`{"date":"2026-09-01","agents":{},"cost":0,"seconds":0,"completed":0}`,
	}
	if len(days) != len(want) {
		t.Fatalf("strategic days = %+v, want %d", days, len(want))
	}
	for i, d := range days {
		data, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != want[i] {
			t.Errorf("day %d json =\n%s\nwant\n%s", i, data, want[i])
		}
	}
}

func TestStrategicDaysAreDaysOfZerosWithNoDocuments(t *testing.T) {
	data, err := json.Marshal(Compute(nil, Options{Now: now}))
	if err != nil {
		t.Fatal(err)
	}
	if first := `"strategic_days":[{"date":"2026-08-02","agents":{},"cost":0,"seconds":0,"completed":0},`; !strings.Contains(string(data), first) {
		t.Errorf("no documents should report a day of zeros from the window's start: %s", data)
	}
}
