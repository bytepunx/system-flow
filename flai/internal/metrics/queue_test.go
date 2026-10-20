package metrics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// queueEvent is a line of the merge queue's event log at at, of kind on the
// stories given.
func queueEvent(at, kind, batch string, stories ...string) mergequeue.Event {
	t, err := time.Parse(workitem.TimeFormat, at)
	if err != nil {
		panic(err)
	}
	return mergequeue.Event{At: t, Event: kind, Batch: batch, Stories: stories}
}

// queueEvents are the merge queue's lines of the fixtures, over the window
// from coordStart, a Wednesday of week 35 at noon, to coordNow, a Wednesday
// of week 38 at noon, the week of each noted above it.
func queueEvents() []mergequeue.Event {
	enq, built, verified, bisected := mergequeue.EventEnqueued, mergequeue.EventBatchBuilt, mergequeue.EventVerified, mergequeue.EventBisected
	landed, sentBack, dequeued := mergequeue.EventLanded, mergequeue.EventSentBack, mergequeue.EventDequeued
	return []mergequeue.Event{
		// before the window: S-0001 enqueued, its wait starting
		queueEvent("2026-08-20T10:00:00Z", enq, "", "S-0001"),
		// week 35 before the window's start: S-0009's batch and landing are
		// left out, its wait with them
		queueEvent("2026-08-26T09:00:00Z", enq, "", "S-0009"),
		queueEvent("2026-08-26T09:30:00Z", built, "B0", "S-0009"),
		queueEvent("2026-08-26T10:00:00Z", landed, "B0", "S-0009"),
		// week 35: S-0001 lands, 7 days and 2 hours after its enqueueing
		queueEvent("2026-08-27T10:00:00Z", built, "B1", "S-0001"),
		queueEvent("2026-08-27T11:00:00Z", verified, "B1", "S-0001"),
		queueEvent("2026-08-27T12:00:00Z", landed, "B1", "S-0001"),
		// week 36: a red batch of two is bisected to S-0002, sent back after
		// 3 hours, and S-0003, enqueued as s-3, lands alone after 4
		queueEvent("2026-09-01T09:00:00Z", enq, "", "S-0002"),
		queueEvent("2026-09-01T09:00:00Z", enq, "", "s-3"),
		queueEvent("2026-09-01T10:00:00Z", built, "B2", "S-0002", "S-0003"),
		queueEvent("2026-09-01T11:00:00Z", verified, "B2", "S-0002", "S-0003"),
		queueEvent("2026-09-01T11:05:00Z", bisected, "B2", "S-0002", "S-0003"),
		queueEvent("2026-09-01T12:00:00Z", sentBack, "B2", "S-0002"),
		queueEvent("2026-09-01T12:30:00Z", built, "B3", "S-0003"),
		queueEvent("2026-09-01T13:00:00Z", landed, "B3", "S-0003"),
		// week 37, with no batch: S-0002, enqueued again, lands 2 hours after
		// the new enqueueing; S-0004 lands with no enqueued line, so no wait
		queueEvent("2026-09-08T10:00:00Z", enq, "", "S-0002"),
		queueEvent("2026-09-08T10:00:00Z", landed, "B4", "S-0004"),
		queueEvent("2026-09-08T12:00:00Z", landed, "B4", "S-0002"),
		// week 38: S-0005 enqueued and removed, then a landing after now
		queueEvent("2026-09-14T09:00:00Z", enq, "", "S-0005"),
		queueEvent("2026-09-14T09:10:00Z", dequeued, "", "S-0005"),
		queueEvent("2026-09-17T00:00:00Z", landed, "B5", "S-0005"),
	}
}

// S-0404: each value of queue, per ISO week of the window, from the merge
// queue's event log, and each story's queue waits that end in the window.
func TestQueueCountsEachValuePerWeek(t *testing.T) {
	got, perStory := mergeQueue(queueEvents(), coordStart, coordNow)
	want := &Queue{Weeks: []QueueWeek{
		{Week: "2026-W35", Start: "2026-08-24", Landed: 1, WaitMean: 170 * 3600, WaitMax: 170 * 3600, Batches: 1, BatchSizeMean: 1, BatchSizeMax: 1},
		{Week: "2026-W36", Start: "2026-08-31", Landed: 1, SentBack: 1, WaitMean: 3.5 * 3600, WaitMax: 4 * 3600, Batches: 2, BatchSizeMean: 1.5, BatchSizeMax: 2, Bisections: 1},
		{Week: "2026-W37", Start: "2026-09-07", Landed: 2, WaitMean: 2 * 3600, WaitMax: 2 * 3600},
		{Week: "2026-W38", Start: "2026-09-14"},
	}}
	if !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.Marshal(got)
		wantJSON, _ := json.Marshal(want)
		t.Errorf("queue =\n%s\nwant\n%s", gotJSON, wantJSON)
	}
	wantStories := map[string]float64{"S-0001": 170 * 3600, "S-0002": 5 * 3600, "S-0003": 4 * 3600}
	if !reflect.DeepEqual(perStory, wantStories) {
		t.Errorf("queue waits by story = %v, want %v", perStory, wantStories)
	}
}

// S-0404: the JSON names are those of metrics.md § Merge queue, which the
// dashboard charts.
func TestQueueJSONNames(t *testing.T) {
	q, _ := mergeQueue(nil, coordNow.Add(-time.Hour), coordNow)
	got, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"weeks":[{"week":"2026-W38","start":"2026-09-14","landed":0,"sent_back":0,` +
		`"wait_mean_seconds":0,"wait_max_seconds":0,"batches":0,"batch_size_mean":0,"batch_size_max":0,"bisections":0}]}`
	if string(got) != want {
		t.Errorf("queue JSON =\n%s\nwant\n%s", got, want)
	}
}

// S-0404: Compute reports the queue and each item's queue wait from the
// options' event log; with no log the section and the waits are absent.
func TestComputeReportsTheQueue(t *testing.T) {
	s, ip := workitem.Story, workitem.InProgress
	const day = "2026-08-20T00:00:00Z"
	items := []*workitem.Item{
		claimItem("S-0001", s, day, nil, ip, day),
		claimItem("S-0002", s, day, nil, ip, day),
		claimItem("S-0003", s, day, nil, ip, day),
		claimItem("S-0004", s, day, nil, ip, day),
		claimItem("S-0009", s, day, nil, ip, day),
	}
	opt := Options{Now: coordNow, Since: 21 * 24 * time.Hour, QueueEvents: queueEvents(), QueueLogged: true}
	rep := Compute(items, opt)
	if rep.Queue == nil || len(rep.Queue.Weeks) != 4 || rep.Queue.Weeks[1].Bisections != 1 {
		t.Fatalf("rep.Queue = %+v", rep.Queue)
	}
	waits := map[string]float64{}
	for _, m := range rep.Items {
		if m.QueueWait != nil {
			waits[m.ID] = *m.QueueWait
		}
	}
	// S-0004 has no enqueued line and S-0009's wait ended before the window
	want := map[string]float64{"S-0001": 170 * 3600, "S-0002": 5 * 3600, "S-0003": 4 * 3600}
	if !reflect.DeepEqual(waits, want) {
		t.Errorf("items' queue waits = %v, want %v", waits, want)
	}
	raw, err := json.Marshal(rep.Items[0])
	if err != nil || !strings.Contains(string(raw), `"queue_wait_seconds":612000`) {
		t.Errorf("S-0001's JSON: %v %s", err, raw)
	}

	opt.QueueEvents, opt.QueueLogged = nil, false
	rep = Compute(items, opt)
	raw, err = json.Marshal(rep)
	if err != nil || rep.Queue != nil || strings.Contains(string(raw), `"queue"`) || strings.Contains(string(raw), `"queue_wait_seconds"`) {
		t.Errorf("a missing log: %v, queue %+v, JSON has the queue: %v", err, rep.Queue, strings.Contains(string(raw), `"queue`))
	}
}

// S-0404: a log with no lines is a log: the section is there, every week
// counting 0, and no item has a queue wait.
func TestComputeReportsAnEmptyQueueLog(t *testing.T) {
	rep := Compute([]*workitem.Item{claimItem("S-0001", workitem.Story, "2026-08-20T00:00:00Z", nil)},
		Options{Now: coordNow, Since: 21 * 24 * time.Hour, QueueLogged: true})
	if rep.Queue == nil || len(rep.Queue.Weeks) != 4 {
		t.Fatalf("rep.Queue = %+v", rep.Queue)
	}
	for _, w := range rep.Queue.Weeks {
		if w != (QueueWeek{Week: w.Week, Start: w.Start}) {
			t.Errorf("week %s of an empty log: %+v", w.Week, w)
		}
	}
	if rep.Items[0].QueueWait != nil {
		t.Errorf("a queue wait from an empty log: %v", *rep.Items[0].QueueWait)
	}
}
