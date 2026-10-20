package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0404: flai stats prints the merge queue's section after the waiting,
// from the queue's event log, and --json and the dashboard's read carry it
// under queue, with each story's queue wait under items; with no log, there
// is neither.
func TestStatsReportsTheMergeQueue(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("LOG_FORMAT", "json")
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, title := range []string{"Lands", "Sent back"} {
		if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Owner: "t", Now: time.Date(2026, 9, 28, 9, 0, 0, 0, time.UTC)}); err != nil {
			t.Fatal(err)
		}
	}
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)

	// no log: no section and no queue waits
	out, errOut, code := runInAt(t, root, at, "stats", "--json")
	if code != 0 || strings.Contains(out, `"queue"`) || strings.Contains(out, `"queue_wait_seconds"`) {
		t.Errorf("stats --json without a log: %d %s\n%s", code, errOut, out)
	}
	if out, _, _ = runInAt(t, root, at, "stats"); strings.Contains(out, "merge queue") {
		t.Errorf("stats without a log prints the queue:\n%s", out)
	}

	enq, built, verified, bisected := mergequeue.EventEnqueued, mergequeue.EventBatchBuilt, mergequeue.EventVerified, mergequeue.EventBisected
	ev := func(at, kind, batch string, stories ...string) mergequeue.Event {
		t.Helper()
		when, err := time.Parse(workitem.TimeFormat, at)
		if err != nil {
			t.Fatal(err)
		}
		return mergequeue.Event{At: when, Event: kind, Batch: batch, Stories: stories}
	}
	// week 40: a batch of two fails and is bisected to S-0002, sent back
	// after an hour; S-0001 lands alone after three
	events := []mergequeue.Event{
		ev("2026-09-29T09:00:00Z", enq, "", "S-0001"),
		ev("2026-09-29T09:00:00Z", enq, "", "S-0002"),
		ev("2026-09-29T09:10:00Z", built, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T09:40:00Z", verified, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T09:41:00Z", bisected, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T10:00:00Z", mergequeue.EventSentBack, "B1", "S-0002"),
		ev("2026-09-29T11:00:00Z", built, "B2", "S-0001"),
		ev("2026-09-29T12:00:00Z", mergequeue.EventLanded, "B2", "S-0001"),
	}
	if err := mergequeue.AppendEvents(repo.CacheDir(), events); err != nil {
		t.Fatal(err)
	}

	out, errOut, code = runInAt(t, root, at, "stats", "--json")
	if code != 0 {
		t.Fatalf("stats --json: %s", errOut)
	}
	var rep struct {
		Queue *metrics.Queue        `json:"queue"`
		Items []metrics.ItemMetrics `json:"items"`
	}
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatal(err)
	}
	// the 30 days to the 3rd of October span ISO weeks 36 to 40
	if rep.Queue == nil || len(rep.Queue.Weeks) != 5 || rep.Queue.Weeks[0].Week != "2026-W36" {
		t.Fatalf("queue: %+v", rep.Queue)
	}
	want := metrics.QueueWeek{Week: "2026-W40", Start: "2026-09-28", Landed: 1, SentBack: 1, WaitMean: 2 * 3600, WaitMax: 3 * 3600,
		Batches: 2, BatchSizeMean: 1.5, BatchSizeMax: 2, Bisections: 1}
	if got := rep.Queue.Weeks[4]; got != want {
		t.Errorf("week 40 = %+v, want %+v", got, want)
	}
	waits := map[string]float64{}
	for _, m := range rep.Items {
		if m.QueueWait != nil {
			waits[m.ID] = *m.QueueWait
		}
	}
	if len(waits) != 2 || waits["S-0001"] != 3*3600 || waits["S-0002"] != 3600 {
		t.Errorf("items' queue waits = %v", waits)
	}

	out, _, code = runInAt(t, root, at, "stats")
	wantText := "\nmerge queue:\n" +
		"  stories        landed 1 · sent back 1 · longest wait 3h\n" +
		"  batches        built 2 · mean size 1.5 · largest 2 · bisections 1\n" +
		"  by week:\n" +
		"    2026-W40 (2026-09-28)  stories landed 1 · sent back 1 · mean wait 2h · longest 3h; batches built 2 · mean size 1.5 · largest 2 · bisections 1\n"
	if code != 0 || !strings.Contains(out, wantText) {
		t.Errorf("text: %d, want it to hold %q, got\n%s", code, wantText, out)
	}
	// the dashboard's read answers the same
	if _, rerr := sameAnswer(t, root, hostapi.Host{}, "stats.get", `{}`, "stats"); rerr != nil {
		t.Errorf("stats.get: %+v", rerr)
	}

	// a log that cannot be read is warned of, and the section is absent
	if err := os.WriteFile(mergequeue.EventsPath(repo.CacheDir()), []byte("{not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runInAt(t, root, at, "stats", "--json")
	if code != 0 || strings.Contains(out, `"queue"`) || !strings.Contains(errOut, "merge queue's event log") {
		t.Errorf("an unreadable log: %d\n%s\n%s", code, out, errOut)
	}
}

// S-0404: the queue's section comes after the waiting, and before the
// turns, and a queue that counts nothing in the window prints nothing.
func TestPrintQueueStatsComesAfterTheWaitingAndSkipsAnEmptyQueue(t *testing.T) {
	var out bytes.Buffer
	printQueueStats(&app{out: &out}, &metrics.Queue{Weeks: []metrics.QueueWeek{{Week: "2026-W40", Start: "2026-09-28"}}})
	if out.Len() != 0 {
		t.Errorf("an empty queue prints:\n%s", out.String())
	}
	rep := &metrics.Report{
		WindowStart: "2026-09-03T20:00:00Z", Type: workitem.Story, WindowDays: 30,
		Waiting: metrics.Waiting{Weeks: []metrics.WaitWeek{{Week: "2026-W40", Start: "2026-09-28", Items: 1}}},
		Queue:   &metrics.Queue{Weeks: []metrics.QueueWeek{{Week: "2026-W40", Start: "2026-09-28", Landed: 1}}},
	}
	printSummary(&app{out: &out}, rep, false)
	text := out.String()
	waiting, queue := strings.Index(text, "\nwaiting, over 1 completed:\n"), strings.Index(text, "\nmerge queue:\n")
	if waiting < 0 || queue < waiting {
		t.Errorf("the queue's section is not after the waiting:\n%s", text)
	}
}
