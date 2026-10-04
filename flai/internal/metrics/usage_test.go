package metrics

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func doneWith(id, at string, u *usage.Usage) *workitem.Item {
	return &workitem.Item{ID: id, Type: workitem.Story, Status: workitem.Done, Created: "2026-08-20T00:00:00Z",
		Transitions: []workitem.Transition{{To: workitem.InProgress, At: "2026-08-20T00:00:00Z"}, {To: workitem.Done, At: at}}, Usage: u}
}

func spend(seconds int64, models ...usage.Model) *usage.Usage {
	return &usage.Usage{Source: usage.SourceLog, Seconds: seconds, Models: models}
}

func TestUsageIsTotalledByModelAndLaidOutAgainstTimeAndCost(t *testing.T) {
	opus := func(read int64, cost float64) usage.Model {
		return usage.Model{Model: "claude-opus-5-5", CacheRead: read, Cost: cost}
	}
	haiku := usage.Model{Model: "claude-haiku-4-5", Output: 600, Cost: 0.1}
	items := []*workitem.Item{
		doneWith("S-0002", "2026-08-30T00:00:00Z", spend(1800, opus(3000, 2), haiku)),
		doneWith("S-0001", "2026-08-25T00:00:00Z", spend(3600, opus(1000, 1))),
		doneWith("S-0003", "2026-08-26T00:00:00Z", nil),                        // no usage: left out
		doneWith("S-0004", "2026-07-01T00:00:00Z", spend(60, opus(99999, 99))), // before the window
	}
	rep := Compute(items, Options{Now: now})
	u := rep.Usage
	if u.Items != 2 || u.Tokens != 4600 || !near(u.Cost, 3.1) || u.Seconds != 5400 {
		t.Fatalf("totals = %+v", u)
	}
	if len(u.Models) != 2 || u.Models[0].Model != "claude-haiku-4-5" || u.Models[1].Items != 2 || !near(*u.Models[1].TokensPerHour, 4000/1.5) {
		t.Errorf("models = %+v", u.Models)
	}
	if len(u.Done) != 2 || u.Done[0].ID != "S-0001" || u.Done[1].Done != 2 || !near(u.Done[1].Cost, 3.1) || u.Done[1].Tokens != 4600 {
		t.Errorf("done = %+v", u.Done)
	}
	if h := u.ByModel["claude-haiku-4-5"]; len(h) != 1 || h[0].ID != "S-0002" || h[0].Done != 1 || h[0].Tokens != 600 {
		t.Errorf("by model = %+v", u.ByModel)
	}
	for _, m := range rep.Items {
		switch m.ID {
		case "S-0002":
			if m.Usage == nil || m.Usage.Tokens != 3600 || !near(*m.Usage.TokensPerHour, 7200) || len(m.Usage.Models) != 2 || !near(*m.Usage.Models[1].TokensPerHour, 1200) {
				t.Errorf("S-0002 usage = %+v", m.Usage)
			}
		case "S-0003":
			if m.Usage != nil {
				t.Errorf("an item with no usage has %+v", m.Usage)
			}
		}
	}
	empty := Compute(nil, Options{Now: now, Since: time.Hour})
	if empty.Usage.Models == nil || empty.Usage.Done == nil || empty.Usage.ByModel == nil {
		t.Error("an empty report's usage lists are null")
	}
}

// strategicItems are stories on which agents, strategic agents, or both spent,
// one done before the window, and others open with a forecast or an estimate
// (ADR-0083).
func strategicItems() []*workitem.Item {
	opus := func(read int64, cost float64) usage.Model {
		return usage.Model{Model: "claude-opus-5-5", CacheRead: read, Cost: cost}
	}
	haiku := func(out int64, cost float64) usage.Model {
		return usage.Model{Model: "claude-haiku-4-5", Output: out, Cost: cost}
	}
	planned := func(u *usage.Usage, s ...usage.Strategic) *usage.Usage {
		if u == nil {
			u = &usage.Usage{Source: usage.SourceSum, Models: []usage.Model{}}
		}
		u.Strategic = s
		return u
	}
	planner := func(seconds int64, m usage.Model) usage.Strategic {
		return usage.Strategic{Kind: "planner", Seconds: seconds, Estimated: true, Models: []usage.Model{m}}
	}
	analyzer := usage.Strategic{Kind: "analyzer", Seconds: 30, Estimated: true, Models: []usage.Model{haiku(50, 0.01)}}
	cancelled := doneWith("S-0008", "2026-08-25T10:00:00Z", spend(1800, opus(500, 1)))
	cancelled.Status = workitem.Cancelled
	cancelled.Transitions[1].To = workitem.Cancelled
	open := func(id, status string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Status: status, Created: "2026-08-30T00:00:00Z"}
	}
	forecast := open("S-0004", workitem.InProgress)
	forecast.Forecast, forecast.Estimate = &workitem.Forecast{Duration: "2h"}, "10h"
	forecast.Usage = planned(nil, planner(40, haiku(10, 0.02)))
	estimated := open("S-0005", workitem.Backlog)
	estimated.Estimate = "30m"
	summed := doneWith("S-0007", "2026-07-01T00:00:00Z", spend(3600, opus(100, 50)))
	summed.Usage.Source = usage.SourceSum
	return []*workitem.Item{
		doneWith("S-0001", "2026-08-25T09:00:00Z", planned(spend(3600, opus(1000, 1)), planner(120, opus(500, 0.2)))),
		doneWith("S-0002", "2026-08-24T09:00:00Z", planned(nil, analyzer, planner(60, haiku(100, 0.05)))),
		doneWith("S-0003", "2026-08-26T09:00:00Z", spend(1800, opus(2000, 2), haiku(600, 0.1))),
		forecast, estimated, open("S-0006", workitem.Backlog), cancelled, summed,
		typed(workitem.Task, doneWith("T-0001", "2026-08-25T09:05:00Z", spend(3600, opus(100, 100)))),
	}
}

func TestStrategicUsageIsReportedApartFromTheAgents(t *testing.T) {
	rep := Compute(strategicItems(), Options{Now: now})
	u := rep.Usage
	// the agents' figures are S-0001's and S-0003's alone: no planner's model,
	// no strategic-only S-0002, no cancelled S-0008
	if u.Items != 2 || u.Tokens != 3600 || !near(u.Cost, 3.1) || u.Seconds != 5400 {
		t.Fatalf("agents' totals = %+v", u)
	}
	if len(u.Models) != 2 || u.Models[0].Model != "claude-haiku-4-5" || u.Models[0].Items != 1 || u.Models[0].Tokens != 600 || u.Models[1].Items != 2 || u.Models[1].Tokens != 3000 {
		t.Errorf("agents' models = %+v", u.Models)
	}
	if len(u.Done) != 2 || u.Done[0].ID != "S-0001" || u.Done[1].ID != "S-0003" || len(u.ByModel["claude-haiku-4-5"]) != 1 {
		t.Errorf("done = %+v, by model %+v", u.Done, u.ByModel)
	}
	s := u.Strategic
	if s.Items != 2 || s.Tokens != 650 || !near(s.Cost, 0.26) || s.Seconds != 210 || !s.Estimated {
		t.Errorf("strategic totals = %+v", s)
	}
	if len(s.Kinds) != 2 || s.Kinds[0].Kind != "planner" || s.Kinds[0].Items != 2 || s.Kinds[0].Tokens != 600 || !near(s.Kinds[0].Cost, 0.25) || s.Kinds[0].Seconds != 180 ||
		s.Kinds[1].Kind != "analyzer" || s.Kinds[1].Items != 1 || s.Kinds[1].Tokens != 50 || !near(s.Kinds[1].Cost, 0.01) {
		t.Errorf("strategic kinds = %+v", s.Kinds)
	}
	per := map[string]ItemMetrics{}
	for _, m := range rep.Items {
		per[m.ID] = m
	}
	if m := per["S-0001"].Usage; m == nil || m.Tokens != 1000 || !near(m.Cost, 1) || len(m.Models) != 1 ||
		len(m.Strategic) != 1 || m.Strategic[0] != (ItemStrategic{Kind: "planner", Tokens: 500, Cost: 0.2, Seconds: 120, Estimated: true}) {
		t.Errorf("S-0001 = %+v", m)
	}
	if m := per["S-0002"].Usage; m == nil || m.Tokens != 0 || m.Cost != 0 || m.Seconds != 0 || len(m.Models) != 0 || m.TokensPerMinute != nil ||
		len(m.Strategic) != 2 || m.Strategic[0].Kind != "analyzer" || m.Strategic[1].Tokens != 100 {
		t.Errorf("an item with only strategic usage = %+v", m)
	}
	if m := per["S-0003"].Usage; m == nil || m.Strategic == nil || len(m.Strategic) != 0 {
		t.Errorf("an item with no strategic usage = %+v", m)
	}
	if m := per["S-0006"].Usage; m != nil {
		t.Errorf("an item with no usage = %+v", m)
	}
	// strategic use per day: S-0002 is completed on 24 August, but agents
	// spent nothing on it, so there is no cost per item
	for _, d := range rep.StrategicDays {
		if d.Date == "2026-08-24" && (d.Completed != 1 || d.CostPerItem != nil) {
			t.Errorf("24 August = %+v", d)
		}
		if d.Date == "2026-08-25" && (d.CostPerItem == nil || !near(*d.CostPerItem, 1)) {
			t.Errorf("25 August = %+v", d)
		}
	}
	empty := Compute(nil, Options{Now: now})
	if e := empty.Usage; e.Strategic.Kinds == nil || e.Strategic.Items != 0 || e.CostPerAgentHour != nil {
		t.Errorf("an empty report's strategic usage = %+v, rate %v", e.Strategic, e.CostPerAgentHour)
	}
	raw, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"strategic":[{"kind":"planner","tokens":500,"cost":0.2,"seconds":120,"estimated":true}]`,
		`"strategic":{"items":2,"tokens":650,"cost":0.26,"seconds":210,"estimated":true,"kinds":[{"kind":"planner","items":2,`,
		`"cost_per_agent_hour":2.05`, `"expected_cost":{"cost":4.1,"from":"forecast","estimated":true}`} {
		if !strings.Contains(string(raw), key) {
			t.Errorf("the report's JSON lacks %s", key)
		}
	}
}

func TestStrategicSpendIsLaidOutBesideTheAgents(t *testing.T) {
	st := Compute(strategicItems(), Options{Now: now}).Usage.Spend[workitem.Story]
	if st.Items != 2 || st.Tokens != 3600 || st.Seconds != 5400 || len(st.Models) != 2 || st.Models[0].Items != 1 || st.Models[0].Tokens != 600 {
		t.Errorf("stories' agents = %+v %+v", st.Spend, st.Models)
	}
	if st.Strategic != (StrategicSpend{Items: 2, Tokens: 650, Cost: 0.26, Seconds: 210}) {
		t.Errorf("stories' strategic = %+v", st.Strategic)
	}
	// S-0002, on which only strategic agents spent, starts the series on 24
	// August and counts in no agents' figure there
	if len(st.Buckets) != 9 || st.Buckets[0].At != "2026-08-24T00:00:00Z" {
		t.Fatalf("buckets = %d, from %s", len(st.Buckets), st.Buckets[0].At)
	}
	b := st.Buckets[0]
	if b.Items != 0 || b.Tokens != 0 || b.Models != nil || b.Strategic != (StrategicSpend{Items: 1, Tokens: 150, Cost: 0.06, Seconds: 90}) {
		t.Errorf("24 August = %+v", b)
	}
	b = st.Buckets[1]
	if b.Items != 1 || b.Tokens != 1000 || !near(b.MeanTokens, 500) || b.Strategic != (StrategicSpend{Items: 1, Tokens: 500, Cost: 0.2, Seconds: 120}) {
		t.Errorf("25 August = %+v", b)
	}
	if b = st.Buckets[2]; b.Items != 1 || b.Strategic != (StrategicSpend{}) {
		t.Errorf("26 August = %+v", b)
	}
	if ta := Compute(strategicItems(), Options{Now: now}).Usage.Spend[workitem.Task]; ta.Items != 1 || ta.Strategic != (StrategicSpend{}) {
		t.Errorf("tasks = %+v", ta)
	}
}

func TestCostPerAgentHourIsOverTheStoriesMeasuredFromTheirLogs(t *testing.T) {
	items := strategicItems()
	// S-0001, S-0003, and the cancelled S-0008: $4.10 over two hours; the
	// summed S-0007, the task, and what strategic agents spent are left out
	if r := CostPerAgentHour(items); r == nil || !near(*r, 2.05) {
		t.Fatalf("rate = %v", r)
	}
	archived := doneWith("S-0009", "2025-01-01T00:00:00Z", spend(3600, usage.Model{Model: "claude-opus-5-5", Cost: 3.9}))
	archived.Archived = true
	idle := doneWith("S-0010", "2026-08-25T00:00:00Z", spend(0, usage.Model{Model: "claude-opus-5-5", Cost: 9}))
	// not windowed: one done long ago counts as well, and one that took no
	// agent time does not
	if r := CostPerAgentHour(append(items, archived, idle)); r == nil || *r != 2.6667 {
		t.Errorf("rate with an old story = %v", r)
	}
	if r := CostPerAgentHour(items[3:6]); r != nil {
		t.Errorf("no story measured: %v", *r)
	}
	// rounded to four decimals
	third := []*workitem.Item{doneWith("S-0011", "2026-08-25T00:00:00Z", spend(10800, usage.Model{Model: "claude-opus-5-5", Cost: 1}))}
	if r := CostPerAgentHour(third); r == nil || *r != 0.3333 {
		t.Errorf("rounded rate = %v", r)
	}
}

func TestExpectedCostIsTheForecastOrEstimateAtTheCostPerAgentHour(t *testing.T) {
	per := map[string]ItemMetrics{}
	for _, m := range Compute(strategicItems(), Options{Now: now}).Items {
		per[m.ID] = m
	}
	if e := per["S-0004"].ExpectedCost; e == nil || *e != (ExpectedCost{Cost: 4.1, From: "forecast", Estimated: true}) {
		t.Errorf("from the forecast, not the estimate = %+v", e)
	}
	if e := per["S-0005"].ExpectedCost; e == nil || *e != (ExpectedCost{Cost: 1.025, From: "estimate", Estimated: true}) {
		t.Errorf("from the estimate = %+v", e)
	}
	if e := per["S-0006"].ExpectedCost; e != nil {
		t.Errorf("without a duration = %+v", e)
	}
	rate := 3.0
	it := &workitem.Item{Forecast: &workitem.Forecast{Delivery: "2026-09-10T00:00:00Z"}, Estimate: "1h20m"}
	if e := ExpectedCostOf(it, &rate); e == nil || *e != (ExpectedCost{Cost: 4, From: "estimate", Estimated: true}) {
		t.Errorf("a forecast without a duration = %+v", e)
	}
	if e := ExpectedCostOf(it, nil); e != nil {
		t.Errorf("without a rate = %+v", e)
	}
	none := Compute(strategicItems()[3:6], Options{Now: now})
	for _, m := range none.Items {
		if m.ExpectedCost != nil {
			t.Errorf("%s has an expected cost with no story measured: %+v", m.ID, m.ExpectedCost)
		}
	}
}
