package metrics

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

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
	want := `[{"kind":"planner","cost":0.6213,"seconds":900,"activities":3,"last_run":"2026-08-31T18:00:00Z","log":[` +
		`{"at":"2026-08-02T12:00:00Z","seconds":100,"cost":0.1,"estimated":true,"items":[]},` +
		`{"at":"2026-08-31T18:00:00Z","seconds":723,"cost":0.4213,"items":["E-0016","S-0230"]}]},` +
		`{"kind":"analyzer","cost":0,"seconds":0,"activities":0,"last_run":"","log":[]}]`
	if string(data) != want {
		t.Errorf("json =\n%s\nwant\n%s", data, want)
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
