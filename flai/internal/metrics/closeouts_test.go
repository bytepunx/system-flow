package metrics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

const closeOutRoot = "/home/op/project"

// closeOutEntry is a serve journal entry of a close-out of closeOutRoot's
// project at at, with outcome, verdict, duration in milliseconds (none when
// negative), and detail.
func closeOutEntry(at, outcome, verdict string, ms int64, detail string) JournalEntry {
	t, err := time.Parse(time.RFC3339, at)
	if err != nil {
		panic(err)
	}
	e := JournalEntry{At: t, Action: "close-out", Method: "verify.job", Root: closeOutRoot, Outcome: outcome, Verdict: verdict, Detail: detail}
	if ms >= 0 {
		e.DurationMS = &ms
	}
	return e
}

func float(f float64) *float64 { return &f }

// S-0546, ADR-0229 §7: the close-outs of the project counted by verdict per
// week of the window, interrupted ones as failed, answers from the record
// apart and out of the durations, with the nearest-rank p50 and p85 of the
// milliseconds of those that ran, over 1000; a week with none that ran has
// no durations, and start entries, other projects', other actions', and
// entries outside the window are not counted.
func TestCloseOutsCountsTheJournalsEndEntriesByWeek(t *testing.T) {
	t.Parallel()
	// from Wednesday of week 40 at noon to Wednesday of week 43 at noon
	start := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	other := closeOutEntry("2026-10-06T10:00:00Z", "done", "passed", 1000, "S-0002 job x")
	other.Root = "/home/op/elsewhere"
	notCloseOut := closeOutEntry("2026-10-06T10:00:00Z", "done", "passed", 1000, "")
	notCloseOut.Action, notCloseOut.Method = "verify", "verify.run"
	otherMethod := closeOutEntry("2026-10-06T10:00:00Z", "done", "passed", 1000, "")
	otherMethod.Method = "verify.run"
	j := CloseOutJournal{Root: closeOutRoot + "/", Entries: []JournalEntry{
		// week 40 before the window's start
		closeOutEntry("2026-09-30T11:00:00Z", "done", "passed", 5000, "S-0001 job a"),
		// week 40 in the window: one stopped, one interrupted
		closeOutEntry("2026-09-30T13:00:00Z", "started", "", -1, "S-0001 job b"),
		closeOutEntry("2026-10-01T09:00:00Z", "failed", "stopped", 2500, "S-0001 job b"),
		closeOutEntry("2026-10-02T09:00:00Z", "failed", "interrupted", -1, "S-0001 job c: interrupted, flai serve stopped while it ran"),
		// week 41: seven that ran, two answered from the record, and others'
		closeOutEntry("2026-10-05T09:00:00Z", "done", "passed", 7000, "S-0003 job"),
		closeOutEntry("2026-10-05T10:00:00Z", "done", "passed", 1000, "S-0003 job"),
		closeOutEntry("2026-10-05T11:00:00Z", "failed", "failed", 6000, "S-0003 job"),
		closeOutEntry("2026-10-06T09:00:00Z", "done", "passed", 3000, "S-0003 job"),
		closeOutEntry("2026-10-06T10:00:00Z", "failed", "stopped", 2000, "S-0003 job"),
		closeOutEntry("2026-10-07T09:00:00Z", "done", "passed", 5000, "S-0003 job"),
		closeOutEntry("2026-10-07T10:00:00Z", "done", "passed", 4321, "S-0003 job"),
		closeOutEntry("2026-10-08T09:00:00Z", "done", "passed", -1, "from the record: S-0003's close-out passed at abc1234, run 2026-10-07T10:00:00Z"),
		closeOutEntry("2026-10-08T10:00:00Z", "failed", "stopped", -1, "from the record: S-0003's close-out stopped at lint at abc1235, run 2026-10-07T11:00:00Z"),
		other, notCloseOut, otherMethod,
		// week 42: nothing
		// week 43 after now
		closeOutEntry("2026-10-21T13:00:00Z", "done", "passed", 1000, "S-0004 job"),
	}}
	got := closeOuts(j, start, now)
	want := &CloseOuts{Weeks: []CloseOutWeek{
		{Week: "2026-W40", Start: "2026-09-28", Stopped: 1, Failed: 1, P50: float(2.5), P85: float(2.5)},
		// sorted 1000 2000 3000 4321 5000 6000 7000: p50 rank 4, p85 rank 6
		{Week: "2026-W41", Start: "2026-10-05", Passed: 5, Stopped: 1, Failed: 1, Recorded: 2, P50: float(4.321), P85: float(6)},
		{Week: "2026-W42", Start: "2026-10-12"},
		{Week: "2026-W43", Start: "2026-10-19"},
	}}
	if !reflect.DeepEqual(got, want) {
		a, _ := json.Marshal(got)
		b, _ := json.Marshal(want)
		t.Errorf("close-outs:\n got %s\nwant %s", a, b)
	}
	raw, err := json.Marshal(got.Weeks[2])
	if err != nil {
		t.Fatal(err)
	}
	if s := string(raw); strings.Contains(s, "p50") || strings.Contains(s, "p85") {
		t.Errorf("a week none ran has durations: %s", s)
	}
	raw, err = json.Marshal(got.Weeks[1])
	if err != nil {
		t.Fatal(err)
	}
	if s := string(raw); s != `{"week":"2026-W41","start":"2026-10-05","passed":5,"stopped":1,"failed":1,"recorded":2,"p50_seconds":4.321,"p85_seconds":6}` {
		t.Errorf("a week's keys: %s", s)
	}
}

// S-0546: Compute reports close-outs only from a journal, and an empty one
// as weeks counting nothing.
func TestComputeReportsCloseOutsOnlyFromAJournal(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	if rep := Compute(nil, Options{Now: now, Since: 7 * 24 * time.Hour}); rep.CloseOuts != nil {
		t.Errorf("no journal: %+v", rep.CloseOuts)
	}
	raw, err := json.Marshal(Compute(nil, Options{Now: now, Since: 7 * 24 * time.Hour}))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), `"closeouts"`) {
		t.Errorf("no journal, yet closeouts in the json")
	}
	rep := Compute(nil, Options{Now: now, Since: 7 * 24 * time.Hour, CloseOuts: &CloseOutJournal{Root: closeOutRoot}})
	want := &CloseOuts{Weeks: []CloseOutWeek{{Week: "2026-W42", Start: "2026-10-12"}, {Week: "2026-W43", Start: "2026-10-19"}}}
	if !reflect.DeepEqual(rep.CloseOuts, want) {
		t.Errorf("an empty journal: %+v", rep.CloseOuts)
	}
}
