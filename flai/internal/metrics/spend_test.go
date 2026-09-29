package metrics

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func typed(typ string, it *workitem.Item) *workitem.Item {
	it.Type = typ
	return it
}

func opusSpent(tokens int64, cost float64) usage.Model {
	return usage.Model{Model: "claude-opus-5-5", CacheRead: tokens, Cost: cost}
}

// spendItems are stories done on three days of one week and one the week
// before, a day with none between them, tasks, and what is left out.
func spendItems() []*workitem.Item {
	haiku := usage.Model{Model: "claude-haiku-4-5", Output: 600, Cost: 0.1}
	estimated := spend(120, opusSpent(500, 0.5))
	estimated.Estimated = true
	cancelled := doneWith("S-0009", "2026-08-25T10:00:00Z", spend(60, opusSpent(7777, 7)))
	cancelled.Status = workitem.Cancelled
	cancelled.Transitions[1].To = workitem.Cancelled
	return []*workitem.Item{
		doneWith("S-0001", "2026-08-25T09:10:00Z", spend(600, opusSpent(1000, 1))),
		doneWith("S-0002", "2026-08-25T09:50:00Z", spend(1200, opusSpent(3000, 2), haiku)),
		doneWith("S-0003", "2026-08-27T23:59:59Z", spend(0, opusSpent(400, 0))),
		doneWith("S-0004", "2026-08-23T12:00:00Z", spend(300, opusSpent(2000, 4))),
		doneWith("S-0005", "2026-08-26T00:00:00Z", nil),                             // no usage
		doneWith("S-0006", "2026-07-01T00:00:00Z", spend(60, opusSpent(99999, 99))), // before the window
		cancelled,
		typed(workitem.Task, doneWith("T-0001", "2026-08-25T09:05:00Z", estimated)),
	}
}

func TestSpendIsLaidOutPerBucketForEveryType(t *testing.T) {
	rep := Compute(spendItems(), Options{Now: now, Type: workitem.Task})
	u := rep.Usage
	if u.Bucket != BucketDay {
		t.Errorf("bucket = %q", u.Bucket)
	}
	if u.Items != 1 {
		t.Errorf("the report is about tasks and totals %d", u.Items)
	}
	st := u.Spend[workitem.Story]
	if st.Items != 4 || st.Tokens != 7000 || !near(st.Cost, 7.1) || st.Seconds != 2100 {
		t.Fatalf("stories over the window = %+v", st.Spend)
	}
	if !near(*st.TokensPerItem, 1750) || !near(*st.CostPerItem, 7.1/4) || !near(*st.TokensPerMinute, 200) || !near(*st.TokensPerDollar, 7000/7.1) {
		t.Errorf("stories' averages = %v %v %v %v", *st.TokensPerItem, *st.CostPerItem, *st.TokensPerMinute, *st.TokensPerDollar)
	}
	if len(st.Models) != 2 || st.Models[0].Model != "claude-haiku-4-5" || st.Models[0].Items != 1 || st.Models[0].Seconds != 1200 ||
		st.Models[1].Items != 4 || st.Models[1].Tokens != 6400 || !near(*st.Models[1].TokensPerItem, 1600) {
		t.Errorf("stories' models = %+v", st.Models)
	}
	// 23 August to 1 September, the day now falls in: ten buckets
	if len(st.Buckets) != 10 || st.Buckets[0].At != "2026-08-23T00:00:00Z" || st.Buckets[9].At != "2026-09-01T00:00:00Z" {
		t.Fatalf("buckets = %d, from %s", len(st.Buckets), st.Buckets[0].At)
	}
	b := st.Buckets[2] // 25 August: S-0001 and S-0002
	if b.Items != 2 || b.Tokens != 4600 || !near(b.Cost, 3.1) || b.Seconds != 1800 || !near(*b.TokensPerItem, 2300) || !near(*b.CostPerItem, 1.55) ||
		!near(*b.TokensPerMinute, 4600.0/30) || !near(*b.TokensPerDollar, 4600/3.1) {
		t.Errorf("25 August = %+v", b.Spend)
	}
	if !near(b.MeanTokens, 6600.0/3) || !near(b.MeanCost, 7.1/3) {
		t.Errorf("25 August's running means = %v %v", b.MeanTokens, b.MeanCost)
	}
	if len(b.Models) != 2 || b.Models[0].Tokens != 600 || !near(*b.Models[0].TokensPerMinute, 30) || b.Models[1].Tokens != 4000 || b.Models[1].Items != 2 || !near(*b.Models[1].CostPerItem, 1.5) {
		t.Errorf("25 August's models = %+v", b.Models)
	}
	if e := st.Buckets[1]; e.Items != 0 || e.Tokens != 0 || e.TokensPerItem != nil || e.Models != nil || !near(e.MeanTokens, 1000) {
		t.Errorf("a day with nothing done = %+v", e)
	}
	// S-0003 took no agent time and cost nothing: no rate, no tokens per dollar
	if z := st.Buckets[4]; z.Items != 1 || z.TokensPerMinute != nil || z.TokensPerDollar != nil || !near(*z.TokensPerItem, 400) {
		t.Errorf("27 August = %+v", z.Spend)
	}
	if l := st.Buckets[9]; !near(l.MeanTokens, 700) || !near(l.MeanCost, 0.71) {
		t.Errorf("the last running means = %v %v", l.MeanTokens, l.MeanCost)
	}
	ta := u.Spend[workitem.Task]
	if ta.Items != 1 || !ta.Estimated || len(ta.Buckets) != 8 || !ta.Buckets[0].Estimated || ta.Buckets[1].Estimated {
		t.Errorf("tasks = %+v", ta)
	}
	if ep := u.Spend[workitem.Epic]; ep == nil || ep.Items != 0 || ep.Buckets == nil || len(ep.Buckets) != 0 || ep.Models == nil {
		t.Errorf("epics = %+v", ep)
	}
	raw, err := json.Marshal(u.Spend[workitem.Story].Buckets[2])
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"at":`, `"items":2`, `"tokens_per_item":`, `"cost_per_item":`, `"tokens_per_minute":`, `"tokens_per_dollar":`, `"mean_tokens":`, `"mean_cost":`, `"models":[{"model":"claude-haiku-4-5","items":1`} {
		if !strings.Contains(string(raw), key) {
			t.Errorf("a bucket's JSON lacks %s: %s", key, raw)
		}
	}
}

func TestSpendByTheHourAndByTheWeek(t *testing.T) {
	hour := Compute(spendItems(), Options{Now: now, Since: 8 * 24 * time.Hour, Bucket: BucketHour}).Usage
	st := hour.Spend[workitem.Story]
	// the window starts on 24 August at noon: S-0004 is before it
	if hour.Bucket != BucketHour || st.Items != 3 || st.Buckets[0].At != "2026-08-25T09:00:00Z" || st.Buckets[0].Items != 2 {
		t.Fatalf("by the hour = %+v", st.Buckets[0])
	}
	if want := 7*24 + 3 + 1; len(st.Buckets) != want || st.Buckets[want-1].At != "2026-09-01T12:00:00Z" {
		t.Errorf("hours = %d, to %s", len(st.Buckets), st.Buckets[len(st.Buckets)-1].At)
	}
	week := Compute(spendItems(), Options{Now: now, Bucket: BucketWeek}).Usage.Spend[workitem.Story]
	// 23 August 2026 is a Sunday: its week starts on Monday 17 August
	if len(week.Buckets) != 3 || week.Buckets[0].At != "2026-08-17T00:00:00Z" || week.Buckets[0].Items != 1 ||
		week.Buckets[1].At != "2026-08-24T00:00:00Z" || week.Buckets[1].Items != 3 || week.Buckets[2].Items != 0 {
		t.Errorf("by the week = %+v", week.Buckets)
	}
	if !near(week.Buckets[2].MeanTokens, 7000.0/3) {
		t.Errorf("the weeks' running mean = %v", week.Buckets[2].MeanTokens)
	}
}

func TestABucketIsCheckedAgainstTheWindow(t *testing.T) {
	for _, ok := range []struct {
		bucket string
		window time.Duration
	}{{"", 365 * 24 * time.Hour}, {BucketDay, 365 * 24 * time.Hour}, {BucketWeek, time.Hour}, {BucketHour, 31 * 24 * time.Hour}} {
		if err := CheckBucket(ok.bucket, ok.window); err != nil {
			t.Errorf("%q over %s: %v", ok.bucket, ok.window, err)
		}
	}
	if err := CheckBucket(BucketHour, 32*24*time.Hour); err == nil || !strings.Contains(err.Error(), "31 days") {
		t.Errorf("an hour over 32 days: %v", err)
	}
	if err := CheckBucket("minute", time.Hour); err == nil || !strings.Contains(err.Error(), "hour, day, or week") {
		t.Errorf("a minute: %v", err)
	}
}

func TestTheRateIsPerMinuteOfAgentWork(t *testing.T) {
	rep := Compute(spendItems(), Options{Now: now})
	for _, m := range rep.Items {
		if m.ID != "S-0002" {
			continue
		}
		if !near(*m.Usage.TokensPerMinute, 180) || !near(*m.Usage.TokensPerHour, 10800) || !near(*m.Usage.Models[1].TokensPerMinute, 30) {
			t.Errorf("S-0002 = %+v", m.Usage)
		}
	}
	for _, m := range rep.Usage.Models {
		if m.Model == "claude-opus-5-5" && !near(*m.TokensPerMinute, 6400.0/35) {
			t.Errorf("opus over the window = %v", *m.TokensPerMinute)
		}
	}
}
