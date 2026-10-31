package metrics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// concurrencyItems is a window of ten hours to now, 02:00 to 12:00 on the
// 1st of September, over which stories are in progress 0, 1, and 2 at once:
// 0 for 3h, 1 for 5h, and 2 for 2h. S-A started before the window opens and
// went to review and back; S-B went to review; S-C was cancelled; S-D is in
// progress to now; and a task in progress all along is of another type.
func concurrencyItems() []*workitem.Item {
	const created = "2026-08-31T00:00:00Z"
	return []*workitem.Item{
		claimItem("S-A", workitem.Story, created, nil,
			workitem.InProgress, "2026-09-01T01:00:00Z", workitem.Review, "2026-09-01T04:00:00Z",
			workitem.InProgress, "2026-09-01T05:00:00Z", workitem.Done, "2026-09-01T07:00:00Z"),
		claimItem("S-B", workitem.Story, created, nil,
			workitem.InProgress, "2026-09-01T03:00:00Z", workitem.Review, "2026-09-01T06:00:00Z", workitem.Done, "2026-09-01T06:30:00Z"),
		claimItem("S-C", workitem.Story, created, nil,
			workitem.InProgress, "2026-09-01T08:00:00Z", workitem.Cancelled, "2026-09-01T09:00:00Z"),
		claimItem("S-D", workitem.Story, created, nil, workitem.InProgress, "2026-09-01T11:00:00Z"),
		claimItem("T-A", workitem.Task, created, nil, workitem.InProgress, "2026-09-01T02:00:00Z"),
	}
}

// S-0551: the seconds of the window with exactly n stories in progress,
// counted from their in-progress intervals cut to the window, review left
// out and cancelled ones in, run through the limit, and their shares sum to
// 1; the mean while busy is over the seconds with one or more.
func TestConcurrencyCountsTheSecondsWithNInProgress(t *testing.T) {
	h := 3600
	rep := Compute(concurrencyItems(), Options{Now: now, Since: 10 * time.Hour, WIPLimit: 3})
	c := rep.Concurrency
	if c.WindowSeconds != 10*h || c.Limit != 3 {
		t.Errorf("window %d, limit %d; want %d, 3", c.WindowSeconds, c.Limit, 10*h)
	}
	if want := []int{3 * h, 5 * h, 2 * h, 0}; !reflect.DeepEqual(c.Seconds, want) {
		t.Errorf("seconds = %v, want %v", c.Seconds, want)
	}
	wantShare := []float64{0.3, 0.5, 0.2, 0}
	sum := 0.0
	for n, s := range c.Share {
		sum += s
		if n >= len(wantShare) || !near(s, wantShare[n]) {
			t.Errorf("share = %v, want %v", c.Share, wantShare)
			break
		}
	}
	if !near(sum, 1) {
		t.Errorf("shares sum to %v", sum)
	}
	if c.MeanBusy == nil || !near(*c.MeanBusy, (5.0+2*2)/7) {
		t.Errorf("mean busy = %v, want %v", deref(c.MeanBusy), (5.0+2*2)/7)
	}
}

// S-0551: without a limit the counts run through the most in progress at
// once, and limit is absent; with nothing in progress the window is all at
// 0 and mean_busy is absent.
func TestConcurrencyRunsThroughTheMostSeenAndLeavesOutWhatIsNot(t *testing.T) {
	rep := Compute(concurrencyItems(), Options{Now: now, Since: 10 * time.Hour})
	if want := []int{3 * 3600, 5 * 3600, 2 * 3600}; !reflect.DeepEqual(rep.Concurrency.Seconds, want) {
		t.Errorf("seconds = %v, want %v", rep.Concurrency.Seconds, want)
	}
	b, err := json.Marshal(rep.Concurrency)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), `"limit"`) || !strings.Contains(string(b), `"mean_busy"`) {
		t.Errorf("concurrency without a limit: %s", b)
	}

	idle := Compute(nil, Options{Now: now, Since: 24 * time.Hour, WIPLimit: 2})
	c := idle.Concurrency
	if !reflect.DeepEqual(c.Seconds, []int{86400, 0, 0}) || !reflect.DeepEqual(c.Share, []float64{1, 0, 0}) || c.MeanBusy != nil {
		t.Errorf("idle: %+v", c)
	}
	if b, err = json.Marshal(c); err != nil {
		t.Fatal(err)
	}
	if want := `{"window_seconds":86400,"limit":2,"seconds":[86400,0,0],"share":[1,0,0]}`; string(b) != want {
		t.Errorf("idle json = %s, want %s", b, want)
	}
}

// S-0551: the ready wait is the queue time of the stories summary.queue_time
// counts, completed in the window and not cancelled, at p50 and p85 by
// nearest rank; with none, the percentiles are absent.
func TestReadyWaitIsTheQueueTimeAtP50AndP85(t *testing.T) {
	const created = "2026-08-31T00:00:00Z"
	waited := func(id, ready, started, end string) *workitem.Item {
		return claimItem(id, workitem.Story, created, nil, workitem.Ready, ready, workitem.InProgress, started, end, "2026-09-01T10:00:00Z")
	}
	items := []*workitem.Item{
		waited("S-1", "2026-09-01T00:00:00Z", "2026-09-01T01:00:00Z", workitem.Done),
		waited("S-2", "2026-09-01T00:00:00Z", "2026-09-01T02:00:00Z", workitem.Done),
		waited("S-3", "2026-09-01T00:00:00Z", "2026-09-01T04:00:00Z", workitem.Done),
		waited("S-4", "2026-09-01T00:00:00Z", "2026-09-01T08:00:00Z", workitem.Cancelled),
		claimItem("S-5", workitem.Story, created, nil, workitem.Ready, "2026-09-01T00:00:00Z", workitem.InProgress, "2026-09-01T09:00:00Z"),
	}
	rep := Compute(items, Options{Now: now, Since: 24 * time.Hour})
	w := rep.ReadyWait
	if w.Count != 3 || w.P50 == nil || *w.P50 != 7200 || w.P85 == nil || *w.P85 != 14400 {
		t.Errorf("ready wait = %d, %v, %v; want 3, 7200, 14400", w.Count, deref(w.P50), deref(w.P85))
	}
	if q := rep.Summary.QueueTime; q.Count != w.Count || q.P50 != *w.P50 || q.P85 != *w.P85 {
		t.Errorf("ready wait %+v differs from summary.queue_time %+v", w, q)
	}

	none := Compute(items[4:], Options{Now: now, Since: 24 * time.Hour})
	b, err := json.Marshal(none.ReadyWait)
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"count":0}`; string(b) != want {
		t.Errorf("ready wait with none = %s, want %s", b, want)
	}
}
