package usage

import (
	"math"
	"reflect"
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

// ADR-0095: an orchestrator activity's usage is split evenly between the
// items it named, in whole tokens and seconds that add up to the whole, the
// remainder to the first items named, and its cost over the count.
func TestSplitAddsUpToTheWhole(t *testing.T) {
	u := agents(100, Model{Model: "a", Input: 10, Output: 7, CacheRead: 2, CacheWrite: 3, Cost: 1}, Model{Model: "b", Output: 1, Cost: 0.01})
	u.Estimated = true
	shares := u.Split(3)
	if len(shares) != 3 {
		t.Fatalf("shares = %d, want 3", len(shares))
	}
	var seconds int64
	sum := &Usage{}
	for i, s := range shares {
		if s.Source != SourceLog || !s.Estimated || s.Strategic != nil || len(s.Models) != 2 {
			t.Errorf("share %d = %+v, want u's source and estimate, both models, no strategic", i, s)
		}
		if s.Models[0].Cost != 1.0/3 || s.Models[1].Cost != 0.01/3 {
			t.Errorf("share %d's costs = %v, %v, want each model's over 3", i, s.Models[0].Cost, s.Models[1].Cost)
		}
		seconds += s.Seconds
		sum.Add(s)
	}
	want := []Model{{Model: "a", Input: 4, Output: 3, CacheRead: 1, CacheWrite: 1}, {Model: "a", Input: 3, Output: 2, CacheRead: 1, CacheWrite: 1}, {Model: "a", Input: 3, Output: 2, CacheRead: 0, CacheWrite: 1}}
	for i, w := range want {
		got := shares[i].Models[0]
		got.Cost = 0
		if got != w {
			t.Errorf("share %d of a = %+v, want %+v, the remainder to the first", i, got, w)
		}
	}
	if s := []int64{shares[0].Seconds, shares[1].Seconds, shares[2].Seconds}; s[0] != 34 || s[1] != 33 || s[2] != 33 {
		t.Errorf("seconds = %v, want 34, 33, 33", s)
	}
	if seconds != u.Seconds || sum.Tokens() != u.Tokens() || math.Abs(sum.Cost()-u.Cost()) > 1e-12 {
		t.Errorf("shares add up to %ds, %d tokens, $%v; want %ds, %d tokens, $%v", seconds, sum.Tokens(), sum.Cost(), u.Seconds, u.Tokens(), u.Cost())
	}
	if one := u.Split(1); len(one) != 1 || !reflect.DeepEqual(one[0], u) {
		t.Errorf("one share = %+v, want u", one)
	}
	if u.Split(0) != nil || (*Usage)(nil).Split(2) != nil {
		t.Error("no shares, or a share of nothing, is not none")
	}
}
