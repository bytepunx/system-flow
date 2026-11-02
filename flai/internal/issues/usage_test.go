package issues

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func analyzed(seconds int64, cost float64) *usage.Usage {
	return &usage.Usage{Source: usage.SourceLog, Seconds: seconds, Estimated: true,
		Models: []usage.Model{{Model: "claude-opus-5-5", Input: 1, Output: 200, CacheRead: 5000, CacheWrite: 300, Cost: cost}}}
}

func TestAnIssueWithoutUsageGainsNoKey(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "One", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(is.Path)
	if strings.Contains(string(data), "usage") {
		t.Errorf("an issue without usage was written with it:\n%s", data)
	}
	back, err := Read(is.Path)
	if err != nil || back.Usage != nil || back.Marshal() != string(data) {
		t.Errorf("round trip: %v %+v", err, back.Usage)
	}
}

func TestChargingAnIssueSumsPerKind(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "One", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range []*usage.Usage{analyzed(400, 0.8), analyzed(12, 0.2)} {
		if _, err := ChargeStrategic(r, "I-1", workitem.ActivityAnalyzer, u); err != nil {
			t.Fatal(err)
		}
	}
	charged, err := ChargeStrategic(r, is.ID, workitem.ActivityPlanner, analyzed(5, 0.1))
	if err != nil || charged == nil {
		t.Fatalf("planner charge: %v %v", err, charged)
	}
	back, err := Read(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	u := back.Usage
	if u == nil || !u.Empty() || len(u.Strategic) != 2 {
		t.Fatalf("usage = %+v, want the planner's and the analyzer's entries alone", u)
	}
	if p := u.Strategic[0]; p.Kind != workitem.ActivityPlanner || p.Seconds != 5 || !p.Estimated {
		t.Errorf("first entry = %+v, want the planner's", p)
	}
	if a := u.Strategic[1]; a.Kind != workitem.ActivityAnalyzer || a.Seconds != 412 || !a.Estimated || a.Cost() != 1 || len(a.Models) != 1 || a.Models[0].Output != 400 {
		t.Errorf("analyzer = %+v, want both charges in one entry", a)
	}
	if back.Updated != is.Updated {
		t.Errorf("updated %s, want %s: a charge is not an edit", back.Updated, is.Updated)
	}
	if err := back.Validate(); err != nil {
		t.Error(err)
	}
}

func TestAnIssuesUsageSurvivesAReadAndAWrite(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "One", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ChargeStrategic(r, is.ID, workitem.ActivityAnalyzer, analyzed(400, 0.8)); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(is.Path)
	want := "updated: 2026-09-16T10:00:00Z\nusage:\n  strategic:\n    - kind: analyzer\n      seconds: 400\n      estimated: true\n      models:\n" +
		"        - model: claude-opus-5-5\n          input: 1\n          output: 200\n          cache_read: 5000\n          cache_write: 300\n          cost: 0.8\n---\n"
	if !strings.Contains(string(data), want) {
		t.Errorf("file lacks\n%s\nin\n%s", want, data)
	}
	back, err := Read(is.Path)
	if err != nil || back.Marshal() != string(data) {
		t.Fatalf("round trip: %v\n%s", err, back.Marshal())
	}
	if err := Bump(back, "", "", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(is.Path)
	if !strings.Contains(string(after), "usage:\n  strategic:\n    - kind: analyzer\n      seconds: 400\n") {
		t.Errorf("a bump dropped the usage:\n%s", after)
	}
}

func TestChargingNothingOrAnUnknownIssue(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "One", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if charged, err := ChargeStrategic(r, is.ID, workitem.ActivityAnalyzer, &usage.Usage{Source: usage.SourceLog}); err != nil || charged != nil {
		t.Errorf("a usage that spent nothing charged %v: %v", charged, err)
	}
	if _, err := ChargeStrategic(r, "I-0009", workitem.ActivityAnalyzer, analyzed(1, 0.1)); err == nil || !strings.Contains(err.Error(), "I-0009") {
		t.Errorf("charging an unknown issue: %v, want an error naming it", err)
	}
	if _, err := ChargeStrategic(r, is.ID, "guesser", analyzed(1, 0.1)); err == nil {
		t.Error("charging a kind that is not a strategic agent's should fail")
	}
	data, _ := os.ReadFile(is.Path)
	if strings.Contains(string(data), "usage") {
		t.Errorf("a refused charge was written:\n%s", data)
	}
}

// ADR-0249 §5: a role the project defines is charged to an issue under its
// name; a name that is neither a built-in agent nor a defined role is
// refused, naming the kinds there are.
func TestADefinedRoleIsChargedToAnIssueUnderItsName(t *testing.T) {
	r := repo(t)
	dir := filepath.Join(r.Root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "triager.md"), []byte("---\nname: triager\ncharge: issues\n---\n\nTriage.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	is, err := New(r, NewOptions{Title: "One", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ChargeStrategic(r, is.ID, "triager", analyzed(30, 0.3)); err != nil {
		t.Fatal(err)
	}
	back, err := Read(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	if st := back.Usage.Strategic; len(st) != 1 || st[0].Kind != "triager" || st[0].Seconds != 30 || st[0].Cost() != 0.3 {
		t.Errorf("strategic = %+v, want the triager's 30 s and 0.3 USD", st)
	}
	if err := back.Validate(); err != nil {
		t.Error(err)
	}
	if _, err := ChargeStrategic(r, is.ID, "auditor", analyzed(1, 0.1)); err == nil || !strings.Contains(err.Error(), "planner, orchestrator, analyzer, curator, triager") {
		t.Errorf("charging auditor: %v, want refused naming the kinds there are", err)
	}
}

func TestAnIssuesUsageIsValidated(t *testing.T) {
	good := &Issue{ID: "I-0001", Title: "t", Class: "defect", Status: "open", Count: 1, FirstReported: "2026-09-16T10:00:00Z", LastReported: "2026-09-16T10:00:00Z", Updated: "2026-09-16T10:00:00Z"}
	bad := *good
	bad.Usage = &usage.Usage{Source: usage.SourceLog, Seconds: 3, Models: []usage.Model{{Model: "m", Output: 1}}, Strategic: []usage.Strategic{
		{Kind: "analyzer", Seconds: -1, Estimated: true, Models: []usage.Model{{Model: "m", Output: -1}, {Model: "m"}, {}}},
		{Kind: "analyzer", Estimated: true},
		{Kind: "guesser kind", Estimated: true},
		{Kind: "planner"},
	}}
	err := bad.Validate()
	if err == nil {
		t.Fatal("an issue with agents' figures and an unknown kind was valid")
	}
	for _, want := range []string{
		"usage carries agents' figures",
		"usage.strategic[0].seconds is negative",
		"usage.strategic[0].models[0] has a negative count or cost",
		"usage.strategic[0].models[1].model m is listed twice",
		"usage.strategic[0].models[2].model is required",
		"usage.strategic[1].kind analyzer is listed twice",
		`usage.strategic[2].kind "guesser kind" must be a strategic agent's name in kebab-case: one of planner, orchestrator, analyzer`,
		"usage.strategic[3].estimated must be true",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("errors %q lack %q", err, want)
		}
	}
	good.Usage = &usage.Usage{Strategic: []usage.Strategic{{Kind: "analyzer", Seconds: 2, Estimated: true, Models: []usage.Model{}}}}
	if err := good.Validate(); err != nil {
		t.Errorf("an issue with an analyzer entry alone: %v", err)
	}
	// ADR-0249 §5: a role's name, defined now or not, is a kind
	good.Usage.Strategic = append(good.Usage.Strategic, usage.Strategic{Kind: "triager", Seconds: 1, Estimated: true, Models: []usage.Model{}})
	if err := good.Validate(); err != nil {
		t.Errorf("an issue with a role's entry: %v", err)
	}
}
