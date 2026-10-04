package planning

import (
	"reflect"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func replay(t *testing.T, items []*workitem.Item, b *workitem.Board, plan manifest.Planning, id string) Result {
	t.Helper()
	res, err := Replay(items, b, plan, id, now)
	if err != nil {
		t.Fatalf("replay %s: %v", id, err)
	}
	return res
}

// owned is a ready feature story on m1 of size 4 with its own
// forecast.duration.
func owned(id, duration string) *workitem.Item {
	it := story(id, workitem.Ready, "feature", "m1", 2, 2, time.Time{})
	it.Forecast = &workitem.Forecast{Duration: duration, By: "planner", At: at(now)}
	return it
}

// S-0211: a ready story with its own forecast.duration keeps it, and is
// delivered that duration times the cycle factor from history after the
// stories ahead; the basis names its own forecast.
func TestReplayKeepsTheStorysOwnDuration(t *testing.T) {
	held := story("S-0011", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-time.Hour))
	held.Forecast = &workitem.Forecast{Duration: "2h", By: "planner", At: at(now)}
	items := append(past(), held, owned("S-0013", "3h"))
	res := replay(t, items, board(1, "S-0013"), manifest.Planning{}, "S-0013")
	// S-0011 frees the lane at 11:00 + 2h × 1.5 = 14:00; S-0013 is delivered
	// 3h × 2 later, at 20:00.
	if res.Duration != "3h" || res.DurationSeconds != 10800 || res.Delivery != "2026-10-04T20:00:00Z" || res.Default || res.Position != 1 {
		t.Errorf("duration %s (%d s), delivery %s, default %t, position %d; want 3h, 20:00", res.Duration, res.DurationSeconds, res.Delivery, res.Default, res.Position)
	}
	want := History{Match: "nature, model, and size band", Stories: 3, SecondsPerUnit: 400, BusyFactor: 1.5, CycleFactor: 2}
	if res.History != want || res.Size != 4 {
		t.Errorf("history %+v, size %d", res.History, res.Size)
	}
	wantBasis := "Its own forecast of 3h; 1st in the pull order with an in-progress limit of 1, behind S-0011."
	if res.Basis != wantBasis {
		t.Errorf("basis\n got %q\nwant %q", res.Basis, wantBasis)
	}
	if f := forecast(t, items, board(1, "S-0013"), manifest.Planning{}, "S-0013"); f.Duration != "27m" || f.Delivery != "2026-10-04T14:54:00Z" {
		t.Errorf("forecast: duration %s, delivery %s; want the history estimate of 27m, 14:54", f.Duration, f.Delivery)
	}
}

// S-0211: with too little history a story's own duration is still kept, not
// the default, and played out with factors of 1.
func TestReplayWithNoHistoryKeepsTheStorysOwnDuration(t *testing.T) {
	res := replay(t, []*workitem.Item{owned("S-0010", "1h45m")}, board(2), manifest.Planning{DefaultDuration: "90m"}, "S-0010")
	if res.Duration != "1h45m" || res.Default || res.Delivery != "2026-10-04T13:45:00Z" || res.History.Match != "none" {
		t.Errorf("got %+v", res)
	}
	if want := "Its own forecast of 1h45m; 1st in the pull order with an in-progress limit of 2, with nothing ahead of it."; res.Basis != want {
		t.Errorf("basis\n got %q\nwant %q", res.Basis, want)
	}
}

// S-0211: a story with no forecast.duration, or one that is not a duration,
// replays exactly as it is forecast.
func TestReplayWithoutAForecastIsTheForecast(t *testing.T) {
	for _, it := range []*workitem.Item{
		story("S-0013", workitem.Ready, "feature", "m1", 2, 2, time.Time{}),
		owned("S-0013", ""),
		owned("S-0013", "soon"),
		owned("S-0013", "-1h"),
		story("S-0013", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-30*time.Minute)),
	} {
		ahead := owned("S-0012", "2h")
		items := append(past(), ahead, it)
		b := board(1, "S-0012", "S-0013")
		got, want := replay(t, items, b, manifest.Planning{}, "S-0013"), forecast(t, items, b, manifest.Planning{}, "S-0013")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s with forecast %v:\n replay   %+v\n forecast %+v", it.Status, it.Forecast, got, want)
		}
	}
}

// S-0211: when a story ahead is accepted, the replayed delivery moves
// earlier and the duration stays.
func TestReplayMovesEarlierWhenWorkAheadIsAccepted(t *testing.T) {
	ahead := story("S-0012", workitem.Ready, "feature", "m1", 2, 2, time.Time{})
	items := append(past(), ahead, owned("S-0013", "3h"))
	b := board(1, "S-0012", "S-0013")
	before := replay(t, items, b, manifest.Planning{}, "S-0013")
	// S-0012 holds the lane for 27m × 1.5 to 12:40:30; S-0013 is delivered
	// 6h later.
	if before.Delivery != "2026-10-04T18:41:00Z" || before.Position != 2 {
		t.Errorf("before: delivery %s, position %d; want 18:41, 2nd", before.Delivery, before.Position)
	}
	ahead.Status = workitem.Done
	after := replay(t, items, b, manifest.Planning{}, "S-0013")
	if after.Duration != "3h" || after.Delivery != "2026-10-04T18:00:00Z" || after.Position != 1 || len(after.Ahead) != 0 {
		t.Errorf("after: duration %s, delivery %s, position %d, ahead %v; want 3h, 18:00, 1st", after.Duration, after.Delivery, after.Position, after.Ahead)
	}
}

// S-0211: a story in progress or in review with its own forecast.duration is
// delivered that duration times the cycle factor after it started.
func TestReplayOfAStoryUnderWay(t *testing.T) {
	for _, c := range []struct {
		it      *workitem.Item
		deliver string
	}{
		{story("S-0010", workitem.InProgress, "feature", "m1", 2, 2, now.Add(-30*time.Minute)), "2026-10-04T13:30:00Z"},
		{story("S-0010", workitem.Review, "feature", "m1", 2, 2, now.Add(-3*time.Hour)), "2026-10-04T12:00:00Z"},
	} {
		c.it.Forecast = &workitem.Forecast{Duration: "1h", By: "planner", At: at(now)}
		res := replay(t, append(past(), c.it), board(1), manifest.Planning{}, "S-0010")
		if res.Duration != "1h" || res.Delivery != c.deliver || res.Position != 0 {
			t.Errorf("%s: duration %s, delivery %s, position %d; want 1h, %s", c.it.Status, res.Duration, res.Delivery, res.Position, c.deliver)
		}
		started := c.it.FirstAt(workitem.InProgress)
		if want := "Its own forecast of 1h; " + phrase(c.it.Status) + " since it started at " + stamp(started) + ", so delivery counts from then."; res.Basis != want {
			t.Errorf("%s: basis\n got %q\nwant %q", c.it.Status, res.Basis, want)
		}
	}
}
