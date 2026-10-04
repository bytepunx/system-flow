package usage

import (
	"testing"
)

func agents(seconds int64, models ...Model) *Usage {
	return &Usage{Source: SourceLog, Seconds: seconds, Models: models}
}

// S-0225: what strategic agents spent on an item is kept apart from what its
// agents spent, and none of the agents' figures count it.
func TestStrategicIsKeptApartFromTheAgentsFigures(t *testing.T) {
	u := agents(60, Model{Model: "a", Input: 10, Cost: 1})
	u.AddStrategic("planner", agents(400, Model{Model: "a", Output: 90, Cost: 3}))
	if u.Tokens() != 10 || u.Cost() != 1 || u.Seconds != 60 {
		t.Errorf("agents' figures = %d tokens, $%v, %ds; want the agents' alone", u.Tokens(), u.Cost(), u.Seconds)
	}
	if u.StrategicTokens() != 90 || u.StrategicCost() != 3 || u.StrategicSeconds() != 400 {
		t.Errorf("strategic = %d tokens, $%v, %ds", u.StrategicTokens(), u.StrategicCost(), u.StrategicSeconds())
	}
	sum := Sum(u, u)
	if sum.Cost() != 2 || sum.Seconds != 120 || sum.Strategic != nil {
		t.Errorf("sum = %+v, want the agents' figures twice and no strategic", sum)
	}
	into := agents(0)
	into.Add(u)
	if into.Strategic != nil {
		t.Errorf("Add carried strategic: %+v", into.Strategic)
	}
	only := &Usage{Source: SourceSum, Models: []Model{}}
	only.AddStrategic("planner", agents(5, Model{Model: "a", Input: 1}))
	if !only.Empty() || only.Nothing() {
		t.Errorf("strategic only: Empty %v, Nothing %v; want the agents to have spent nothing and the item something", only.Empty(), only.Nothing())
	}
	if !(*Usage)(nil).Nothing() || !agents(0).Nothing() {
		t.Error("no usage is not nothing")
	}
	if Sum(only) != nil {
		t.Error("a sum of strategic-only usage is not nil")
	}
}

func TestAddStrategicMergesByKindAndModel(t *testing.T) {
	u := &Usage{Source: SourceSum, Models: []Model{}}
	u.AddStrategic("analyzer", agents(1, Model{Model: "b", Input: 1, Cost: 0.00004}))
	u.AddStrategic("planner", agents(10, Model{Model: "b", Input: 1, Cost: 0.123456}, Model{Model: "zero"}))
	u.AddStrategic("planner", agents(20, Model{Model: "a", Output: 2}, Model{Model: "b", Input: 3}))
	u.AddStrategic("planner", nil)
	u.AddStrategic("planner", agents(0))
	if len(u.Strategic) != 2 || u.Strategic[0].Kind != "planner" || u.Strategic[1].Kind != "analyzer" {
		t.Fatalf("strategic = %+v, want planner then analyzer", u.Strategic)
	}
	p := u.Strategic[0]
	if p.Seconds != 30 || !p.Estimated || len(p.Models) != 2 || p.Models[0].Model != "a" || p.Models[1].Input != 4 || p.Models[1].Cost != 0.1235 {
		t.Errorf("planner = %+v, want seconds added, models merged by name, sorted, zero dropped, cost rounded", p)
	}
	if a := u.Strategic[1]; a.Models[0].Cost != 0 || !a.Estimated {
		t.Errorf("analyzer = %+v", a)
	}
}

func TestSameComparesStrategic(t *testing.T) {
	plain := agents(60, Model{Model: "a", Input: 1})
	planned := plain.Clone()
	planned.AddStrategic("planner", agents(5, Model{Model: "a", Input: 1}))
	if Same(plain, planned) || Same(planned, plain) {
		t.Error("usage that differs only in strategic is the same")
	}
	if !Same(planned, planned.Clone()) {
		t.Error("a clone is not the same")
	}
	only := &Usage{Source: SourceSum, Models: []Model{}}
	only.AddStrategic("planner", agents(5, Model{Model: "a", Input: 1}))
	if Same(nil, only) || Same(only, nil) || !Same(only, only.Clone()) {
		t.Error("strategic-only usage compares wrong with no usage or itself")
	}
	more := only.Clone()
	more.AddStrategic("planner", agents(1))
	if Same(only, more) {
		t.Error("strategic seconds that differ are the same")
	}
	if !Same(nil, agents(0)) {
		t.Error("two usages of nothing differ")
	}
}

func TestCloneCopiesStrategic(t *testing.T) {
	u := agents(1, Model{Model: "a", Input: 1})
	u.AddStrategic("planner", agents(5, Model{Model: "a", Input: 1}))
	c := u.Clone()
	c.Strategic[0].Models[0].Input = 99
	c.Models[0].Input = 99
	if u.Strategic[0].Models[0].Input != 1 || u.Models[0].Input != 1 {
		t.Error("changing a clone changed the original")
	}
}

func TestWithStrategicKeepsAnItemsStrategic(t *testing.T) {
	old := &Usage{Source: SourceSum, Models: []Model{}}
	old.AddStrategic("planner", agents(5, Model{Model: "a", Input: 1}))
	sum := agents(60, Model{Model: "b", Input: 2})
	got := WithStrategic(sum, old)
	if got.Seconds != 60 || len(got.Strategic) != 1 || got.Strategic[0].Seconds != 5 || sum.Strategic != nil {
		t.Errorf("with strategic = %+v (sum %+v)", got, sum)
	}
	if got := WithStrategic(nil, old); got == nil || got.Source != SourceSum || !got.Empty() || got.StrategicSeconds() != 5 {
		t.Errorf("nothing summed = %+v, want empty agents' figures carrying the strategic", got)
	}
	if WithStrategic(sum, agents(1)) != sum || WithStrategic(nil, nil) != nil {
		t.Error("with no strategic to keep, it is not the usage given")
	}
}
