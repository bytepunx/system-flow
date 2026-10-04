package planning

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var (
	now     = time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	history = time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
)

func at(t time.Time) string { return t.Format(workitem.TimeFormat) }

// story is an open story of a nature and model, with criteria checkboxes and
// touches, in a state; one in progress or in review started at started.
func story(id, status, nature, model string, criteria, touches int, started time.Time) *workitem.Item {
	it := &workitem.Item{ID: id, Type: workitem.Story, Status: status, Nature: nature, Body: "## Acceptance criteria\n"}
	for i := range criteria {
		it.Body += fmt.Sprintf("- [ ] criterion %d\n", i)
	}
	for i := range touches {
		it.Touches = append(it.Touches, fmt.Sprintf("path/%d", i))
	}
	if model != "" {
		it.Agent = &manifest.Agent{Model: model}
	}
	if status == workitem.InProgress || status == workitem.Review {
		it.Transitions = append(it.Transitions, workitem.Transition{To: workitem.InProgress, At: at(started)})
		if status == workitem.Review {
			it.Transitions = append(it.Transitions, workitem.Transition{To: workitem.Review, At: at(started.Add(time.Minute))})
		}
	}
	return it
}

// done is an archived done story on m1 of size 4 whose agents worked
// seconds, in progress 1.5 times as long and done twice as long after it
// started.
func done(id, nature string, seconds int64) *workitem.Item {
	it := story(id, workitem.Done, nature, "m1", 2, 2, history)
	start := history
	it.Transitions = []workitem.Transition{
		{To: workitem.InProgress, At: at(start)},
		{To: workitem.Review, At: at(start.Add(time.Duration(seconds) * 1500 * time.Millisecond))},
		{To: workitem.Done, At: at(start.Add(time.Duration(seconds) * 2 * time.Second))},
	}
	it.Usage = &usage.Usage{Source: usage.SourceLog, Seconds: seconds}
	it.Archived = true
	return it
}

// past is three done feature stories on m1 of size 4, at 300, 400, and 500
// agent seconds per unit, and an improvement story at 1000.
func past() []*workitem.Item {
	return []*workitem.Item{
		done("S-0001", "feature", 1200), done("S-0002", "feature", 1600), done("S-0003", "feature", 2000),
		done("S-0004", "improvement", 4000),
	}
}

func board(limit int, order ...string) *workitem.Board {
	return &workitem.Board{WIPLimits: map[string]int{workitem.InProgress: limit}, Order: order}
}

func forecast(t *testing.T, items []*workitem.Item, b *workitem.Board, plan manifest.Planning, id string) Result {
	t.Helper()
	res, err := Forecast(items, b, plan, id, now)
	if err != nil {
		t.Fatalf("forecast %s: %v", id, err)
	}
	if strings.ContainsAny(res.Basis, "\n") {
		t.Errorf("the basis is more than one line: %q", res.Basis)
	}
	return res
}

func aheadOf(res Result) string {
	var parts []string
	for _, a := range res.Ahead {
		parts = append(parts, fmt.Sprintf("%s %s %s forecast=%t %s-%s", a.ID, a.Status, a.Duration, a.FromForecast, a.Start, a.Delivery))
	}
	return strings.Join(parts, "; ")
}

// S-0210: a story with three done stories of its nature, model, and band
// takes the median seconds per unit times its size, rounded up to the
// minute, and with nothing ahead starts now and is delivered after its
// duration times the cycle factor.
func TestForecastFromMatchingHistory(t *testing.T) {
	items := append(past(), story("S-0010", workitem.Ready, "feature", "m1", 2, 2, time.Time{}))
	res := forecast(t, items, board(2, "S-0010"), manifest.Planning{}, "S-0010")
	if res.Duration != "27m" || res.DurationSeconds != 1620 || res.Delivery != "2026-10-04T12:54:00Z" || res.Default {
		t.Errorf("duration %s (%d s), delivery %s, default %t; want 27m, 12:54", res.Duration, res.DurationSeconds, res.Delivery, res.Default)
	}
	want := History{Match: "nature, model, and size band", Stories: 3, SecondsPerUnit: 400, BusyFactor: 1.5, CycleFactor: 2}
	if res.History != want || res.Size != 4 || res.Criteria != 2 || res.Touches != 2 || res.Position != 1 || res.Limit != 2 || len(res.Ahead) != 0 {
		t.Errorf("got %+v", res)
	}
	wantBasis := "Median 400 s per unit of size over 3 done feature stories on m1 in the small band, times size 4 (2 criteria, 2 touches); 1st in the pull order with an in-progress limit of 2, with nothing ahead of it."
	if res.Basis != wantBasis {
		t.Errorf("basis\n got %q\nwant %q", res.Basis, wantBasis)
	}
	if short := forecast(t, items, board(2, "S-0010"), manifest.Planning{}, "S-10"); short.ID != "S-0010" || short.Delivery != res.Delivery {
		t.Errorf("S-10 forecast %s delivered %s; want S-0010's", short.ID, short.Delivery)
	}
}

// S-0210: behind a full in-progress limit a story starts when the lane
// frees; a story in progress holds it for its forecast.duration times the
// busy factor, and one ready before it for its computed duration.
func TestForecastBehindAFullLimit(t *testing.T) {
	held := story("S-0011", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-time.Hour))
	held.Forecast = &workitem.Forecast{Duration: "2h", By: "planner", At: at(now)}
	reviewed := story("S-0012", workitem.Review, "feature", "m1", 2, 2, now.Add(-3*time.Hour))
	items := append(past(), held, reviewed,
		story("S-0013", workitem.Ready, "feature", "m1", 2, 2, time.Time{}),
		story("S-0014", workitem.Ready, "feature", "m1", 2, 2, time.Time{}))
	res := forecast(t, items, board(1, "S-0013", "S-0014"), manifest.Planning{}, "S-0014")
	// S-0011 frees the lane at 11:00 + 2h × 1.5 = 14:00; S-0013 holds it for
	// 27m × 1.5 to 14:40:30; S-0014 is delivered 54m later, at 15:34:30.
	if res.Duration != "27m" || res.Delivery != "2026-10-04T15:35:00Z" || res.Position != 2 {
		t.Errorf("duration %s, delivery %s, position %d", res.Duration, res.Delivery, res.Position)
	}
	wantAhead := "S-0011 in-progress 2h forecast=true 2026-10-04T11:00:00Z-2026-10-04T15:00:00Z; S-0013 ready 27m forecast=false 2026-10-04T14:00:00Z-2026-10-04T14:54:00Z"
	if got := aheadOf(res); got != wantAhead {
		t.Errorf("ahead\n got %s\nwant %s", got, wantAhead)
	}
	if !strings.HasSuffix(res.Basis, "; 2nd in the pull order with an in-progress limit of 1, behind S-0011 and S-0013.") {
		t.Errorf("basis %q", res.Basis)
	}
}

// S-0210: over the limit, a lane frees only when enough stories in progress
// have finished to bring their count under it.
func TestForecastOverTheLimit(t *testing.T) {
	short := story("S-0011", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-time.Hour))
	short.Forecast = &workitem.Forecast{Duration: "1h", By: "planner", At: at(now)}
	long := story("S-0012", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-time.Hour))
	long.Forecast = &workitem.Forecast{Duration: "2h", By: "planner", At: at(now)}
	items := append(past(), short, long, story("S-0013", workitem.Ready, "feature", "m1", 2, 2, time.Time{}))
	res := forecast(t, items, board(1), manifest.Planning{}, "S-0013")
	// the lane frees at 14:00, when S-0012 does, not at 12:30
	if res.Delivery != "2026-10-04T14:54:00Z" {
		t.Errorf("delivery %s, want 14:54", res.Delivery)
	}
}

// S-0210: a story waits for the delivery of a story in its after that is in
// progress, even with no limit, and an after that is done holds nothing.
func TestForecastWaitsOnAfter(t *testing.T) {
	held := story("S-0011", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-time.Hour))
	held.Forecast = &workitem.Forecast{Duration: "2h", By: "planner", At: at(now)}
	target := story("S-0012", workitem.Ready, "feature", "m1", 2, 2, time.Time{})
	target.After = []string{"S-0001", "S-0011"}
	free := story("S-0013", workitem.Ready, "feature", "m1", 2, 2, time.Time{})
	free.After = []string{"S-0001"}
	items := append(past(), held, target, free)
	res := forecast(t, items, board(0, "S-0012", "S-0013"), manifest.Planning{}, "S-0012")
	if res.Delivery != "2026-10-04T15:54:00Z" || res.Limit != 0 {
		t.Errorf("delivery %s, limit %d; want 15:54 after S-0011's 15:00", res.Delivery, res.Limit)
	}
	if !strings.Contains(res.Basis, "with no in-progress limit, behind S-0011.") {
		t.Errorf("basis %q", res.Basis)
	}
	if res := forecast(t, items, board(0, "S-0012", "S-0013"), manifest.Planning{}, "S-0013"); res.Delivery != "2026-10-04T12:54:00Z" {
		t.Errorf("an after that is done: delivery %s, want 12:54", res.Delivery)
	}
}

// S-0210: without three stories on a rung the match falls back: another
// model to nature and band, another band to nature, another nature to all.
func TestForecastFallsBackTheLadder(t *testing.T) {
	for _, c := range []struct {
		name                     string
		it                       *workitem.Item
		match, duration, deliver string
		stories                  int
	}{
		{"another model", story("S-0010", workitem.Ready, "feature", "m2", 2, 2, time.Time{}), "nature and size band", "27m", "2026-10-04T12:54:00Z", 3},
		{"no model", story("S-0010", workitem.Ready, "feature", "", 2, 2, time.Time{}), "nature and size band", "27m", "2026-10-04T12:54:00Z", 3},
		// 400 s × 13 = 86.7 min
		{"the large band", story("S-0010", workitem.Ready, "feature", "m1", 6, 7, time.Time{}), "nature", "1h27m", "2026-10-04T14:54:00Z", 3},
		// the median of 300, 400, 500, 1000 is 450; × 4 = 30 min
		{"another nature", story("S-0010", workitem.Ready, "remediation", "m1", 2, 2, time.Time{}), "all done stories", "30m", "2026-10-04T13:00:00Z", 4},
	} {
		res := forecast(t, append(past(), c.it), board(2), manifest.Planning{}, "S-0010")
		if res.History.Match != c.match || res.History.Stories != c.stories || res.Duration != c.duration || res.Delivery != c.deliver || res.Default {
			t.Errorf("%s: match %q over %d, duration %s, delivery %s, default %t", c.name, res.History.Match, res.History.Stories, res.Duration, res.Delivery, res.Default)
		}
	}
}

// S-0210: with fewer than three done stories with usage, the duration is
// planning.default_duration, an hour when unset, and the factors are 1.
func TestForecastWithNoHistoryUsesTheDefault(t *testing.T) {
	noUsage := done("S-0003", "feature", 2000)
	noUsage.Usage = nil
	items := []*workitem.Item{done("S-0001", "feature", 1200), done("S-0002", "feature", 1600), noUsage,
		story("S-0010", workitem.Ready, "feature", "m1", 2, 2, time.Time{})}
	res := forecast(t, items, board(2), manifest.Planning{DefaultDuration: "90m"}, "S-0010")
	want := History{Match: "none", Stories: 2, BusyFactor: 1, CycleFactor: 1}
	if res.Duration != "1h30m" || res.Delivery != "2026-10-04T13:30:00Z" || !res.Default || res.History != want {
		t.Errorf("got %+v", res)
	}
	if !strings.HasPrefix(res.Basis, "Not enough history, 2 done stories with usage where 3 are needed, so planning.default_duration of 1h30m for size 4") {
		t.Errorf("basis %q", res.Basis)
	}
	if res := forecast(t, items, board(2), manifest.Planning{}, "S-0010"); res.Duration != "1h" || res.Delivery != "2026-10-04T13:00:00Z" {
		t.Errorf("unset default: duration %s, delivery %s", res.Duration, res.Delivery)
	}
}

// S-0210: a story in progress or in review is delivered its duration times
// the cycle factor after it started, but not before now, with nothing ahead.
func TestForecastOfAStoryUnderWay(t *testing.T) {
	for _, c := range []struct {
		it      *workitem.Item
		deliver string
	}{
		{story("S-0010", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-30*time.Minute)), "2026-10-04T12:24:00Z"},
		{story("S-0010", workitem.Review, "feature", "m1", 2, 2, now.Add(-time.Hour)), "2026-10-04T12:00:00Z"},
	} {
		ahead := story("S-0011", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-time.Hour))
		res := forecast(t, append(past(), c.it, ahead), board(1), manifest.Planning{}, "S-0010")
		if res.Delivery != c.deliver || res.Position != 0 || len(res.Ahead) != 0 {
			t.Errorf("%s: delivery %s, position %d, ahead %v; want %s", c.it.Status, res.Delivery, res.Position, res.Ahead, c.deliver)
		}
		if !strings.Contains(res.Basis, "since it started at") {
			t.Errorf("%s: basis %q", c.it.Status, res.Basis)
		}
	}
}

// S-0210: only an open story is forecast, and a bad default duration stops
// it with what to write.
func TestForecastRefuses(t *testing.T) {
	task := &workitem.Item{ID: "T-0001", Type: workitem.Task, Status: workitem.Ready}
	cancelled := story("S-0020", workitem.Cancelled, "feature", "m1", 1, 0, time.Time{})
	archived := story("S-0021", workitem.Ready, "feature", "m1", 1, 0, time.Time{})
	archived.Archived = true
	open := story("S-0022", workitem.Ready, "feature", "m1", 1, 0, time.Time{})
	items := append(past(), task, cancelled, archived, open)
	for _, c := range []struct {
		id   string
		plan manifest.Planning
		want string
	}{
		{"S-0099", manifest.Planning{}, "cannot forecast S-0099: there is no such item"},
		{"T-0001", manifest.Planning{}, "cannot forecast T-0001: it is a task, and only a story has a forecast"},
		{"S-0001", manifest.Planning{}, "cannot forecast S-0001: it is archived"},
		{"S-0020", manifest.Planning{}, "cannot forecast S-0020: it is cancelled; forecast a story that is still open"},
		{"S-0021", manifest.Planning{}, "cannot forecast S-0021: it is archived"},
		{"S-0022", manifest.Planning{DefaultDuration: "a day"}, `cannot forecast S-0022: planning.default_duration "a day" is not a duration`},
	} {
		if _, err := Forecast(items, board(2), c.plan, c.id, now); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want %q", c.id, err, c.want)
		}
	}
	unarchived := done("S-0030", "feature", 1200)
	unarchived.Archived = false
	if _, err := Forecast(append(items, unarchived), board(2), manifest.Planning{}, "S-0030", now); err == nil || !strings.Contains(err.Error(), "it is done") {
		t.Errorf("a done story: %v", err)
	}
}

func TestFormatDuration(t *testing.T) {
	for d, want := range map[time.Duration]string{
		2 * time.Hour: "2h", 65 * time.Minute: "1h5m", 45 * time.Minute: "45m", 90 * time.Second: "1m30s",
		time.Hour + 30*time.Second: "1h0m30s", 30 * time.Second: "30s",
	} {
		if got := FormatDuration(d); got != want {
			t.Errorf("FormatDuration(%v) = %q, want %q", d, got, want)
		}
	}
}

func TestOrdinal(t *testing.T) {
	for n, want := range map[int]string{1: "1st", 2: "2nd", 3: "3rd", 4: "4th", 11: "11th", 12: "12th", 13: "13th", 21: "21st", 112: "112th"} {
		if got := ordinal(n); got != want {
			t.Errorf("ordinal(%d) = %q, want %q", n, got, want)
		}
	}
}
