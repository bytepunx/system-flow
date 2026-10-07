package metrics

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// waitItems are stories created on 20 August 2026 and moved as the pairs of
// state and time say, and a task of S-0004.
func waitItems() []*workitem.Item {
	item := func(id string, moves ...string) *workitem.Item {
		it := &workitem.Item{ID: id, Type: workitem.Story, Status: workitem.Backlog, Created: "2026-08-20T00:00:00Z"}
		for i := 0; i < len(moves); i += 2 {
			it.Transitions = append(it.Transitions, workitem.Transition{To: moves[i], At: moves[i+1]})
			it.Status = moves[i]
		}
		return it
	}
	ip, rv, done := workitem.InProgress, workitem.Review, workitem.Done
	task := item("T-0004", ip, "2026-08-27T01:00:00Z")
	task.Type, task.Parent = workitem.Task, "S-0004"
	return []*workitem.Item{
		// answered while in progress, the opener's own follow-up not an answer
		item("S-0001", workitem.Ready, "2026-08-24T00:00:00Z", ip, "2026-08-25T00:00:00Z", done, "2026-08-26T00:00:00Z"),
		// opened before in progress and answered during it; anchored as S-12
		item("S-0012", ip, "2026-08-25T00:00:00Z", rv, "2026-08-25T12:00:00Z", done, "2026-08-25T18:00:00Z"),
		// two overlapping threads
		item("S-0003", ip, "2026-08-26T00:00:00Z", done, "2026-08-27T00:00:00Z"),
		// a thread on its task
		item("S-0004", ip, "2026-08-27T00:00:00Z", done, "2026-08-27T12:00:00Z"),
		// an open thread nobody answered, to now
		item("S-0005", ip, "2026-09-01T06:00:00Z"),
		// a resolved thread nobody else answered, to updated; in review at now
		item("S-0006", ip, "2026-08-31T00:00:00Z", rv, "2026-08-31T12:00:00Z"),
		// an open thread without entries, from created
		item("S-0007", ip, "2026-08-31T00:00:00Z", done, "2026-08-31T10:00:00Z"),
		// sent back from review: two intervals in progress and two in review
		item("S-0008", ip, "2026-08-28T00:00:00Z", rv, "2026-08-28T06:00:00Z", ip, "2026-08-28T08:00:00Z",
			rv, "2026-08-28T10:00:00Z", done, "2026-08-28T11:00:00Z"),
		// no thread and never in review
		item("S-0009", ip, "2026-08-29T00:00:00Z", done, "2026-08-29T01:00:00Z"),
		// cancelled from review: left out of the weeks
		item("S-0010", ip, "2026-08-27T00:00:00Z", rv, "2026-08-27T01:00:00Z", workitem.Cancelled, "2026-08-27T05:00:00Z"),
		// a thread answered after it was done
		item("S-0011", ip, "2026-08-31T00:00:00Z", done, "2026-08-31T02:00:00Z"),
		// done before the window
		item("S-0013", ip, "2026-08-20T06:00:00Z", done, "2026-08-21T00:00:00Z"),
		task,
	}
}

// thread is a thread on an item whose entries are the pairs of time and
// author given.
func thread(item, status, created, updated string, entries ...string) *threads.Thread {
	th := &threads.Thread{Anchor: threads.Anchor{Path: "wip/kanban/x.md", Item: item}, Status: status, Created: created, Updated: updated}
	th.Body = "\n# TH-0001 A question\n\n## Entries\n"
	for i := 0; i < len(entries); i += 2 {
		th.Body += "\n### " + entries[i] + " " + entries[i+1] + "\nSome text.\n"
	}
	return th
}

// waitThreads are threads on waitItems, numbered TH-0001 on in order.
func waitThreads() []*threads.Thread {
	ths := []*threads.Thread{
		thread("S-0001", "answered", "2026-08-25T02:00:00Z", "2026-08-25T05:00:00Z",
			"2026-08-25T02:00:00Z", "agent", "2026-08-25T03:00:00Z", "agent", "2026-08-25T05:00:00Z", "designer"),
		thread("S-12", "answered", "2026-08-24T22:00:00Z", "2026-08-25T01:00:00Z",
			"2026-08-24T22:00:00Z", "agent", "2026-08-25T01:00:00Z", "designer"),
		thread("S-0003", "answered", "2026-08-26T01:00:00Z", "2026-08-26T03:00:00Z",
			"2026-08-26T01:00:00Z", "agent", "2026-08-26T03:00:00Z", "designer"),
		thread("S-0003", "answered", "2026-08-26T02:00:00Z", "2026-08-26T04:00:00Z",
			"2026-08-26T02:00:00Z", "agent", "2026-08-26T04:00:00Z", "designer"),
		thread("T-4", "answered", "2026-08-27T06:00:00Z", "2026-08-27T07:00:00Z",
			"2026-08-27T06:00:00Z", "agent", "2026-08-27T07:00:00Z", "designer"),
		thread("S-0005", "open", "2026-09-01T08:00:00Z", "2026-09-01T08:00:00Z",
			"2026-09-01T08:00:00Z", "agent"),
		thread("S-0006", "resolved", "2026-08-31T02:00:00Z", "2026-08-31T03:00:00Z",
			"2026-08-31T02:00:00Z", "agent", "2026-08-31T03:00:00Z", "agent"),
		thread("S-0007", "open", "2026-08-31T08:00:00Z", "2026-08-31T08:00:00Z"),
		thread("S-0008", "answered", "2026-08-28T05:00:00Z", "2026-08-28T09:00:00Z",
			"2026-08-28T05:00:00Z", "agent", "2026-08-28T09:00:00Z", "designer"),
		thread("S-0011", "answered", "2026-08-31T03:00:00Z", "2026-08-31T04:00:00Z",
			"2026-08-31T03:00:00Z", "agent", "2026-08-31T04:00:00Z", "designer"),
		thread("", "open", "2026-08-25T00:00:00Z", "2026-08-25T00:00:00Z",
			"2026-08-25T00:00:00Z", "agent"),
	}
	for i, th := range ths {
		th.ID = fmt.Sprintf("TH-%04d", i+1)
	}
	return ths
}

// tenDays starts the window at noon on Saturday 22 August 2026, in the ISO
// week that started on Monday 17 August.
const tenDays = 10 * 24 * time.Hour

// S-0205: the union of the waits of an item's threads and its tasks' that
// fall in its in-progress intervals, and its time in review, each absent
// without the input.
func TestWaitPerItemIsThreadsInProgressAndTimeInReview(t *testing.T) {
	rep := Compute(waitItems(), Options{Now: now, Since: tenDays, Threads: waitThreads()})
	want := map[string][2]*float64{ // threads, review
		"S-0001": {val(10800), nil},
		"S-0012": {val(3600), val(21600)},
		"S-0003": {val(10800), nil},
		"S-0004": {val(3600), nil},
		"S-0005": {val(14400), nil},
		"S-0006": {val(3600), val(86400)},
		"S-0007": {val(7200), nil},
		"S-0008": {val(7200), val(10800)},
		"S-0009": {nil, nil},
		"S-0010": {nil, val(14400)},
		"S-0011": {val(0), nil},
		"S-0013": {nil, nil},
	}
	if len(rep.Items) != len(want) {
		t.Fatalf("items = %d, want %d", len(rep.Items), len(want))
	}
	for _, m := range rep.Items {
		w := want[m.ID]
		for k, got := range [2]*float64{m.WaitThreads, m.WaitReview} {
			if (got == nil) != (w[k] == nil) || (got != nil && *got != *w[k]) {
				t.Errorf("%s [threads, review][%d] = %v, want %v", m.ID, k, deref(got), deref(w[k]))
			}
		}
		// S-0220: no orchestrator entry, so no part of the wait is its
		if o := m.WaitThreadsOrchestrator; (o == nil) != (w[0] == nil) || (o != nil && *o != 0) {
			t.Errorf("%s orchestrator's part = %v, want 0 beside a thread wait", m.ID, deref(o))
		}
	}
	data, err := json.Marshal(rep.Items[1])
	if err != nil {
		t.Fatal(err)
	}
	if keys := `"wait_threads_seconds":3600,"wait_review_seconds":21600`; !strings.Contains(string(data), keys) {
		t.Errorf("S-0012 json lacks %s: %s", keys, data)
	}
	data, err = json.Marshal(rep.Items[8])
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "wait_") {
		t.Errorf("S-0009 has no thread and no review and carries a wait: %s", data)
	}
}

// S-0205: per ISO week of the window, the items done in it and the sum and
// mean of their waits; a week without items has no mean, and cancelled items
// and those done before the window are left out. The longest waits follow,
// pinned in TestLongestWaitsOfTheWindow.
func TestWaitingByTheWeek(t *testing.T) {
	data, err := json.Marshal(Compute(waitItems(), Options{Now: now, Since: tenDays, Threads: waitThreads()}).Waiting)
	if err != nil {
		t.Fatal(err)
	}
	none := `"orchestrator":{"count":0,"total_seconds":0},"confirmed":{"count":0,"total_seconds":0}`
	want := `{"weeks":[` +
		`{"week":"2026-W34","start":"2026-08-17","items":0,"empty_wakes":0,"threads":{"total_seconds":0,` + none + `},"review":{"total_seconds":0}},` +
		`{"week":"2026-W35","start":"2026-08-24","items":6,"empty_wakes":0,` +
		`"threads":{"total_seconds":36000,"mean_seconds":6000,` + none + `},"review":{"total_seconds":32400,"mean_seconds":5400}},` +
		`{"week":"2026-W36","start":"2026-08-31","items":2,"empty_wakes":0,` +
		`"threads":{"total_seconds":7200,"mean_seconds":3600,` + none + `},"review":{"total_seconds":0,"mean_seconds":0}}],"empty_wakes":{"count":0},"longest":[`
	if !strings.HasPrefix(string(data), want) {
		t.Errorf("waiting =\n%s\nwant it to start\n%s", data, want)
	}
}

// S-0272, ADR-0105: the empty wakes of the items done in the window are
// summed over it and per week, cancelled items and those done before the
// window or not done left out, and their mean is over those of them that
// carry agents' usage; each item reports its own, 0 when its usage has none.
func TestEmptyWakesAreSummedOverTheItemsDoneInTheWindow(t *testing.T) {
	items := waitItems()
	measured := func(wakes int) *usage.Usage {
		return &usage.Usage{Source: usage.SourceLog, Seconds: 60, EmptyWakes: wakes, Models: []usage.Model{{Model: "m", Input: 1}}}
	}
	planned := &usage.Usage{Source: usage.SourceSum, Models: []usage.Model{}}
	planned.AddStrategic("planner", measured(0))
	for _, it := range items {
		switch it.ID {
		case "S-0001": // done in W35
			it.Usage = measured(3)
		case "S-0003": // done in W35, with usage and no empty wakes
			it.Usage = measured(0)
		case "S-0007": // done in W36
			it.Usage = measured(2)
		case "S-0008": // done in W35, only strategic agents spent on it
			it.Usage = planned
		case "S-0010", "S-0013", "S-0005": // cancelled, done before the window, in progress
			it.Usage = measured(40)
		}
	}
	rep := Compute(items, Options{Now: now, Since: tenDays, Threads: waitThreads()})
	w := rep.Waiting
	if w.EmptyWakes.Count != 5 || w.EmptyWakes.Mean == nil || !near(*w.EmptyWakes.Mean, 5.0/3) {
		t.Errorf("empty wakes = %+v, want 5 over the 3 items with agents' usage", w.EmptyWakes)
	}
	if got := []int{w.Weeks[0].EmptyWakes, w.Weeks[1].EmptyWakes, w.Weeks[2].EmptyWakes}; got[0] != 0 || got[1] != 3 || got[2] != 2 {
		t.Errorf("weeks' empty wakes = %v, want 0, 3, 2", got)
	}
	for _, m := range rep.Items {
		switch m.ID {
		case "S-0001":
			if m.Usage == nil || m.Usage.EmptyWakes != 3 {
				t.Errorf("S-0001 usage = %+v, want 3 empty wakes", m.Usage)
			}
		case "S-0003":
			data, _ := json.Marshal(m.Usage)
			if !strings.Contains(string(data), `"empty_wakes":0`) {
				t.Errorf("S-0003 usage json lacks empty_wakes 0: %s", data)
			}
		case "S-0009":
			if m.Usage != nil {
				t.Errorf("S-0009 carries no usage and reports %+v", m.Usage)
			}
		}
	}
	data, err := json.Marshal(w.EmptyWakes)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), `{"count":5,"mean":1.666`) {
		t.Errorf("waiting.empty_wakes = %s", data)
	}
}

// threadOf is a thread on an item whose entries are the triples of time,
// heading author, and text given.
func threadOf(item, status, created, updated string, entries ...string) *threads.Thread {
	th := &threads.Thread{Anchor: threads.Anchor{Path: "wip/kanban/x.md", Item: item}, Status: status, Created: created, Updated: updated}
	th.Body = "\n# TH-0001 A question\n\n## Entries\n"
	for i := 0; i < len(entries); i += 3 {
		th.Body += "\n### " + entries[i] + " " + entries[i+1] + "\n" + entries[i+2] + "\n"
	}
	return th
}

// orchestratedItems are two stories done in the week of 24 August 2026.
func orchestratedItems() []*workitem.Item {
	ip, done := workitem.InProgress, workitem.Done
	item := func(id, started, completed string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Status: done, Created: "2026-08-20T00:00:00Z",
			Transitions: []workitem.Transition{{To: ip, At: started}, {To: done, At: completed}}}
	}
	return []*workitem.Item{
		item("S-0021", "2026-08-25T00:00:00Z", "2026-08-26T00:00:00Z"),
		item("S-0022", "2026-08-27T00:00:00Z", "2026-08-27T06:00:00Z"),
	}
}

// orchestratedThreads are threads the orchestrator answered, recommended an
// answer on that the operator confirmed or answered otherwise, and
// recommended an answer on that awaits the operator.
func orchestratedThreads() []*threads.Thread {
	const asked, source = "Which way?", "Source: design/system/workflow.md"
	return []*threads.Thread{
		// answered by the orchestrator citing a source: an hour
		threadOf("S-0021", "answered", "2026-08-25T02:00:00Z", "2026-08-25T03:00:00Z",
			"2026-08-25T02:00:00Z", "agent", asked,
			"2026-08-25T03:00:00Z", "orchestrator", "This way.\n\n"+source),
		// overlapping the one before, answered by the orchestrator: the
		// union of the two is an hour and a half
		threadOf("S-0021", "answered", "2026-08-25T02:30:00Z", "2026-08-25T03:30:00Z",
			"2026-08-25T02:30:00Z", "agent", asked,
			"2026-08-25T03:30:00Z", "orchestrator", "That way.\n\n"+source),
		// recommended at 05:00, which ends no wait, and confirmed at 07:00
		threadOf("S-0021", "answered", "2026-08-25T04:00:00Z", "2026-08-25T07:00:00Z",
			"2026-08-25T04:00:00Z", "agent", asked,
			"2026-08-25T05:00:00Z", "orchestrator (recommendation)", "This way.\n\n"+source,
			"2026-08-25T07:00:00Z", "alex", "Confirmed the recommendation of 2026-08-25T05:00:00Z orchestrator.\n\n"+source),
		// recommended, and answered otherwise by the operator: two hours by
		// someone else
		threadOf("S-0021", "answered", "2026-08-25T08:00:00Z", "2026-08-25T10:00:00Z",
			"2026-08-25T08:00:00Z", "agent", asked,
			"2026-08-25T09:00:00Z", "orchestrator (recommendation)", "This way.",
			"2026-08-25T10:00:00Z", "alex", "No, that way."),
		// answered by the orchestrator after the story was done: no part in
		// progress, so not counted
		threadOf("S-0021", "answered", "2026-08-26T01:00:00Z", "2026-08-26T02:00:00Z",
			"2026-08-26T01:00:00Z", "agent", asked,
			"2026-08-26T02:00:00Z", "orchestrator", "This way.\n\n"+source),
		// escalated as a recommendation that still awaits the operator: it
		// waits until now, five and a half hours of it in progress
		threadOf("S-0022", "open", "2026-08-27T00:30:00Z", "2026-08-27T01:00:00Z",
			"2026-08-27T00:30:00Z", "agent", asked,
			"2026-08-27T01:00:00Z", "orchestrator (recommendation)", "Perhaps this way."),
	}
}

// S-0220: a recommendation ends no wait; the orchestrator's answer ends one
// counted under orchestrator, the union of such waits being its part of the
// item's; the operator's confirmation of a recommendation ends one counted
// under confirmed; and waits with no part in progress are not counted.
func TestWaitsTheOrchestratorEndedAreCountedApart(t *testing.T) {
	rep := Compute(orchestratedItems(), Options{Now: now, Since: tenDays, Threads: orchestratedThreads()})
	want := map[string][2]float64{ // threads, the orchestrator's part
		"S-0021": {23400, 5400},
		"S-0022": {19800, 0},
	}
	for _, m := range rep.Items {
		w := want[m.ID]
		if m.WaitThreads == nil || *m.WaitThreads != w[0] || m.WaitThreadsOrchestrator == nil || *m.WaitThreadsOrchestrator != w[1] {
			t.Errorf("%s threads, orchestrator's = %v, %v, want %v", m.ID, deref(m.WaitThreads), deref(m.WaitThreadsOrchestrator), w)
		}
	}
	data, err := json.Marshal(rep.Items[0])
	if err != nil {
		t.Fatal(err)
	}
	if key := `"wait_threads_orchestrator_seconds":5400`; !strings.Contains(string(data), key) {
		t.Errorf("S-0021 json lacks %s: %s", key, data)
	}
	data, err = json.Marshal(rep.Waiting.Weeks[1])
	if err != nil {
		t.Fatal(err)
	}
	week := `{"week":"2026-W35","start":"2026-08-24","items":2,"empty_wakes":0,` +
		`"threads":{"total_seconds":43200,"mean_seconds":21600,` +
		`"orchestrator":{"count":2,"total_seconds":5400},"confirmed":{"count":1,"total_seconds":10800}},` +
		`"review":{"total_seconds":0,"mean_seconds":0}}`
	if string(data) != week {
		t.Errorf("week =\n%s\nwant\n%s", data, week)
	}
}

// sameWaits reports where got differs from want, wait by wait.
func sameWaits(t *testing.T, got, want []LongWait) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("longest = %d waits, want %d: %+v", len(got), len(want), got)
		return
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("longest[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// S-0215, ADR-0114: the ten waits with the most seconds in the window, each
// on its own, of items in any state but cancelled: a thread answered (by its
// answer's author), a thread still open (no end, nobody awaited), and review
// waits (by the move out of review); the one in review now is open.
func TestLongestWaitsOfTheWindow(t *testing.T) {
	items := waitItems()
	for _, it := range items {
		switch it.ID {
		case "S-0012": // done out of review
			it.Transitions[2].By = "alex"
		case "S-0008": // sent back from review
			it.Transitions[2].By = "alex"
		}
	}
	w := Compute(items, Options{Now: now, Since: tenDays, Threads: waitThreads()}).Waiting
	sameWaits(t, w.Longest, []LongWait{
		{Item: "S-0006", Kind: "review", Started: "2026-08-31T12:00:00Z", Seconds: 86400},
		{Item: "S-0012", Kind: "review", Started: "2026-08-25T12:00:00Z", Ended: "2026-08-25T18:00:00Z", Seconds: 21600, Awaited: "alex"},
		{Item: "S-0005", Kind: "thread", Thread: "TH-0006", Started: "2026-09-01T08:00:00Z", Seconds: 14400},
		{Item: "S-0001", Kind: "thread", Thread: "TH-0001", Started: "2026-08-25T02:00:00Z", Ended: "2026-08-25T05:00:00Z", Seconds: 10800, Awaited: "designer"},
		{Item: "S-0003", Kind: "thread", Thread: "TH-0003", Started: "2026-08-26T01:00:00Z", Ended: "2026-08-26T03:00:00Z", Seconds: 7200, Awaited: "designer"},
		{Item: "S-0003", Kind: "thread", Thread: "TH-0004", Started: "2026-08-26T02:00:00Z", Ended: "2026-08-26T04:00:00Z", Seconds: 7200, Awaited: "designer"},
		{Item: "S-0008", Kind: "thread", Thread: "TH-0009", Started: "2026-08-28T05:00:00Z", Ended: "2026-08-28T09:00:00Z", Seconds: 7200, Awaited: "designer"},
		{Item: "S-0008", Kind: "review", Started: "2026-08-28T06:00:00Z", Ended: "2026-08-28T08:00:00Z", Seconds: 7200, Awaited: "alex"},
		{Item: "S-0007", Kind: "thread", Thread: "TH-0008", Started: "2026-08-31T08:00:00Z", Seconds: 7200},
		// the thread on S-12, of three waits of 3600 seconds the earliest
		{Item: "S-0012", Kind: "thread", Thread: "TH-0002", Started: "2026-08-24T22:00:00Z", Ended: "2026-08-25T01:00:00Z", Seconds: 3600, Awaited: "designer"},
	})
	data, err := json.Marshal(w.Longest[:3])
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"item":"S-0006","kind":"review","started":"2026-08-31T12:00:00Z","seconds":86400},` +
		`{"item":"S-0012","kind":"review","started":"2026-08-25T12:00:00Z","ended":"2026-08-25T18:00:00Z","seconds":21600,"awaited":"alex"},` +
		`{"item":"S-0005","kind":"thread","thread":"TH-0006","started":"2026-09-01T08:00:00Z","seconds":14400}]`
	if string(data) != want {
		t.Errorf("longest json =\n%s\nwant\n%s", data, want)
	}
}

// S-0215, ADR-0114: a wait counts only its part inside the window, its start
// still the one the thread or the move records; waits with no part in it are
// left out, and the list is empty, not null, when none has one.
func TestLongestWaitsAreCutByTheWindowStart(t *testing.T) {
	for _, c := range []struct {
		since time.Duration
		want  []LongWait
	}{
		{27 * time.Hour, []LongWait{ // from 2026-08-31T09:00:00Z
			{Item: "S-0006", Kind: "review", Started: "2026-08-31T12:00:00Z", Seconds: 86400},
			{Item: "S-0005", Kind: "thread", Thread: "TH-0006", Started: "2026-09-01T08:00:00Z", Seconds: 14400},
			// in progress until 10:00, an hour of it in the window
			{Item: "S-0007", Kind: "thread", Thread: "TH-0008", Started: "2026-08-31T08:00:00Z", Seconds: 3600},
		}},
		{18 * time.Hour, []LongWait{ // from 2026-08-31T18:00:00Z
			{Item: "S-0006", Kind: "review", Started: "2026-08-31T12:00:00Z", Seconds: 64800},
			{Item: "S-0005", Kind: "thread", Thread: "TH-0006", Started: "2026-09-01T08:00:00Z", Seconds: 14400},
		}},
	} {
		w := Compute(waitItems(), Options{Now: now, Since: c.since, Threads: waitThreads()}).Waiting
		sameWaits(t, w.Longest, c.want)
	}
	data, err := json.Marshal(Compute(waitItems(), Options{Now: now, Since: tenDays}).Waiting.Longest)
	if err != nil {
		t.Fatal(err)
	}
	// without threads, only the review waits remain
	if !strings.HasPrefix(string(data), `[{"item":"S-0006","kind":"review"`) {
		t.Errorf("longest without threads = %s", data)
	}
	data, err = json.Marshal(Compute(nil, Options{Now: now, Since: tenDays}).Waiting.Longest)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "[]" {
		t.Errorf("longest without waits = %s, want []", data)
	}
}

// S-0215, ADR-0114: waits of equal seconds go to the earlier start, then to
// the item and the thread in ID order, a review wait after a thread's; a
// thread resolved with no answer ends when it was updated, nobody awaited;
// and a cancelled item's waits are left out.
func TestLongestWaitsTiesAndAResolvedThread(t *testing.T) {
	story := func(id, status string, moves ...string) *workitem.Item {
		it := &workitem.Item{ID: id, Type: workitem.Story, Status: status, Created: "2026-08-20T00:00:00Z"}
		for i := 0; i < len(moves); i += 3 {
			it.Transitions = append(it.Transitions, workitem.Transition{To: moves[i], At: moves[i+1], By: moves[i+2]})
		}
		return it
	}
	ip, rv := workitem.InProgress, workitem.Review
	items := []*workitem.Item{
		// in progress from 02:00, after its thread started
		story("S-0002", ip, ip, "2026-08-25T02:00:00Z", "agent"),
		// in review from 01:00 to 02:00, sent back by alex
		story("S-0001", workitem.Done, ip, "2026-08-25T00:00:00Z", "agent", rv, "2026-08-25T01:00:00Z", "agent",
			ip, "2026-08-25T02:00:00Z", "alex", workitem.Done, "2026-08-25T03:00:00Z", "alex"),
		story("S-0003", rv, ip, "2026-08-31T00:00:00Z", "agent", rv, "2026-08-31T12:00:00Z", "agent"),
		story("S-0004", workitem.Cancelled, ip, "2026-08-26T00:00:00Z", "agent", rv, "2026-08-26T01:00:00Z", "agent",
			workitem.Cancelled, "2026-08-27T00:00:00Z", "alex"),
	}
	withID := func(id string, th *threads.Thread) *threads.Thread {
		th.ID = id
		return th
	}
	ths := []*threads.Thread{
		withID("TH-0004", thread("S-0002", "resolved", "2026-08-25T01:00:00Z", "2026-08-25T03:00:00Z",
			"2026-08-25T01:00:00Z", "agent")),
		withID("TH-0012", thread("S-0001", "answered", "2026-08-25T01:00:00Z", "2026-08-25T03:00:00Z",
			"2026-08-25T01:00:00Z", "agent", "2026-08-25T03:00:00Z", "alex")),
		withID("TH-0003", thread("S-0001", "answered", "2026-08-25T01:00:00Z", "2026-08-25T03:00:00Z",
			"2026-08-25T01:00:00Z", "agent", "2026-08-25T03:00:00Z", "bob")),
		withID("TH-0005", thread("S-0004", "open", "2026-08-26T00:00:00Z", "2026-08-26T00:00:00Z",
			"2026-08-26T00:00:00Z", "agent")),
	}
	w := Compute(items, Options{Now: now, Since: tenDays, Threads: ths}).Waiting
	sameWaits(t, w.Longest, []LongWait{
		{Item: "S-0003", Kind: "review", Started: "2026-08-31T12:00:00Z", Seconds: 86400},
		// an hour each in progress, from 02:00 to 03:00
		{Item: "S-0001", Kind: "thread", Thread: "TH-0003", Started: "2026-08-25T01:00:00Z", Ended: "2026-08-25T03:00:00Z", Seconds: 3600, Awaited: "bob"},
		{Item: "S-0001", Kind: "thread", Thread: "TH-0012", Started: "2026-08-25T01:00:00Z", Ended: "2026-08-25T03:00:00Z", Seconds: 3600, Awaited: "alex"},
		{Item: "S-0001", Kind: "review", Started: "2026-08-25T01:00:00Z", Ended: "2026-08-25T02:00:00Z", Seconds: 3600, Awaited: "alex"},
		{Item: "S-0002", Kind: "thread", Thread: "TH-0004", Started: "2026-08-25T01:00:00Z", Ended: "2026-08-25T03:00:00Z", Seconds: 3600},
	})
}
