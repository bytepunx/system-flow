package manifest

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

func TestLoadAndFind(t *testing.T) {
	root := t.TempDir()
	body := "version: 1\nname: demo\nkey: d\nlayout:\n  design: arch\n  docs: docs\n  wip: wip\nprojects:\n  - name: a\n    path: a\n    kind: go\n"
	if err := os.WriteFile(filepath.Join(root, File), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "a", "b")
	_ = os.MkdirAll(nested, 0o755)
	p, err := Find(nested)
	if err != nil || p != filepath.Join(root, File) {
		t.Fatalf("Find = %q, %v", p, err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "demo" || m.Dir(root, "design") != filepath.Join(root, "arch") || len(m.Projects) != 1 {
		t.Fatalf("unexpected: %+v", m)
	}
	if _, err := Find(t.TempDir()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	_ = os.WriteFile(p, []byte("version: 1\nname: x\nlayout:\n  design: d\n"), 0o644)
	if _, err := Load(p); err == nil {
		t.Fatal("expected missing layout keys error")
	}
}

// S-0082: checks: is optional and round-trips as a list of named argument
// lists, the operator's choice of where to name them when the host's own
// configuration names none.
func TestChecksRoundTrip(t *testing.T) {
	root := t.TempDir()
	body := "version: 1\nname: demo\nkey: d\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"checks:\n  - name: flai\n    command: [scripts/flai-test.sh]\n  - name: flaiover\n    command: [bash, -c, \"cd flaiover; pnpm test\"]\n"
	p := filepath.Join(root, File)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Checks) != 2 || m.Checks[0].Name != "flai" || len(m.Checks[0].Command) != 1 ||
		m.Checks[1].Name != "flaiover" || len(m.Checks[1].Command) != 3 {
		t.Fatalf("checks: %+v", m.Checks)
	}
	// absent is a nil slice, not an error and not an empty-but-present list
	m2, err := Load(func() string {
		body := "version: 1\nname: demo2\nkey: d2\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
		p2 := filepath.Join(t.TempDir(), File)
		_ = os.WriteFile(p2, []byte(body), 0o644)
		return p2
	}())
	if err != nil || m2.Checks != nil {
		t.Errorf("no checks: named: %v %+v", err, m2.Checks)
	}
}

// S-0198: issues.story_after reads from the manifest; empty is seven days, 0
// turns the warning off, and a value that is not a duration, or is negative,
// is an error naming the key and what to write.
func TestIssuesStoryAfter(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\nissues:\n  story_after: 24h\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if d, err := m.Issues.StoryAfterDuration(); err != nil || d != 24*time.Hour {
		t.Errorf("read 24h as %v, %v", d, err)
	}
	for in, want := range map[string]time.Duration{"": DefaultStoryAfter, " ": DefaultStoryAfter, "0": 0, "0s": 0, "36h30m": 36*time.Hour + 30*time.Minute} {
		if d, err := (Issues{StoryAfter: in}).StoryAfterDuration(); err != nil || d != want {
			t.Errorf("%q: got %v, %v; want %v", in, d, err, want)
		}
	}
	if DefaultStoryAfter != 7*24*time.Hour {
		t.Errorf("the default is %v, not seven days", DefaultStoryAfter)
	}
	for _, bad := range []string{"7d", "a week", "-1h"} {
		_, err := (Issues{StoryAfter: bad}).StoryAfterDuration()
		if err == nil || !strings.Contains(err.Error(), "issues.story_after") || !strings.Contains(err.Error(), "168h") || !strings.Contains(err.Error(), bad) {
			t.Errorf("%q: %v", bad, err)
		}
	}
}

// S-0199: planning reads from the manifest; unset, the currency is USD, the
// hour rate unknown, the cycle a week, and the default duration an hour
// (S-0210), and each bad value is an error
// naming the key and what to write.
func TestPlanning(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\nplanning:\n  currency: EUR\n  hour_rate: 85.5\n  cycle: 336h\n  default_duration: 90m\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if d, err := m.Planning.CycleDuration(); m.Planning.CurrencyCode() != "EUR" || m.Planning.HourRate == nil || *m.Planning.HourRate != 85.5 || err != nil || d != 336*time.Hour {
		t.Errorf("read %+v, cycle %v, %v", m.Planning, d, err)
	}
	if d, err := m.Planning.FallbackDuration(); err != nil || d != 90*time.Minute {
		t.Errorf("default duration %v, %v", d, err)
	}
	if errs := m.Planning.Errors(); len(errs) != 0 {
		t.Errorf("valid planning: %v", errs)
	}
	var unset Planning
	if d, err := unset.FallbackDuration(); err != nil || d != time.Hour {
		t.Errorf("unset default duration %v, %v", d, err)
	}
	if d, err := unset.CycleDuration(); unset.CurrencyCode() != "USD" || unset.HourRate != nil || err != nil || d != 168*time.Hour {
		t.Errorf("defaults: %s, %v, %v", unset.CurrencyCode(), d, err)
	}
	if errs := unset.Errors(); len(errs) != 0 {
		t.Errorf("unset planning: %v", errs)
	}
	zero := 0.0
	if errs := (Planning{HourRate: &zero}).Errors(); len(errs) != 0 {
		t.Errorf("a rate of zero: %v", errs)
	}
	neg, inf := -1.0, math.Inf(1)
	for _, c := range []struct {
		p    Planning
		want string
	}{
		{Planning{Currency: "usd"}, `planning.currency "usd" is not an ISO 4217 code`},
		{Planning{Currency: "EURO"}, `planning.currency "EURO" is not an ISO 4217 code`},
		{Planning{HourRate: &neg}, "planning.hour_rate -1 is not an amount of zero or more"},
		{Planning{HourRate: &inf}, "planning.hour_rate +Inf is not an amount of zero or more"},
		{Planning{Cycle: "1w"}, `planning.cycle "1w" is not a duration longer than zero; write one such as 168h`},
		{Planning{Cycle: "0s"}, `planning.cycle "0s" is not a duration longer than zero`},
		{Planning{Cycle: "-24h"}, `planning.cycle "-24h" is not a duration longer than zero`},
		{Planning{DefaultDuration: "1d"}, `planning.default_duration "1d" is not a duration longer than zero; write one such as 1h`},
		{Planning{DefaultDuration: "0s"}, `planning.default_duration "0s" is not a duration longer than zero`},
	} {
		if got := strings.Join(c.p.Errors(), "; "); !strings.Contains(got, c.want) {
			t.Errorf("%+v: got %q, want %q", c.p, got, c.want)
		}
	}
}

// S-0211: planning.replan and planning.schedule read from the manifest;
// unset, the policy is deterministic and there is no schedule, and a bad
// value is an error naming the key.
func TestPlanningTriggers(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\nplanning:\n  replan: agent\n  schedule: \"*/30 * * * *\"\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if r, err := m.Planning.ReplanPolicy(); err != nil || r != ReplanAgent {
		t.Errorf("replan %q, %v", r, err)
	}
	s, err := m.Planning.PlanSchedule()
	if err != nil || s == nil || s.String() != "*/30 * * * *" {
		t.Fatalf("schedule %v, %v", s, err)
	}
	at := time.Date(2026, 10, 4, 9, 10, 0, 0, time.UTC)
	if next := s.Next(at); !next.Equal(at.Add(20 * time.Minute)) {
		t.Errorf("next after %v: %v", at, next)
	}
	if errs := m.Planning.Errors(); len(errs) != 0 {
		t.Errorf("valid triggers: %v", errs)
	}

	var unset Planning
	if r, err := unset.ReplanPolicy(); err != nil || r != ReplanDeterministic {
		t.Errorf("unset replan %q, %v", r, err)
	}
	if s, err := unset.PlanSchedule(); err != nil || s != nil {
		t.Errorf("unset schedule %v, %v", s, err)
	}
	if s, err := (Planning{Schedule: "  "}).PlanSchedule(); err != nil || s != nil {
		t.Errorf("blank schedule %v, %v", s, err)
	}
	for _, v := range []string{ReplanNever, ReplanDeterministic, ReplanAgent} {
		if r, err := (Planning{Replan: v}).ReplanPolicy(); err != nil || r != v {
			t.Errorf("%q: %q, %v", v, r, err)
		}
	}
	s, err = (Planning{Schedule: "daily"}).PlanSchedule()
	if err != nil || s == nil || s.String() != "daily" {
		t.Fatalf("daily: %v, %v", s, err)
	}
	if next := s.Next(at); !next.Equal(time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("daily after %v: %v", at, next)
	}

	for _, c := range []struct {
		p    Planning
		want string
	}{
		{Planning{Replan: "sometimes"}, `planning.replan "sometimes" is not a replan policy; write never`},
		{Planning{Replan: "Agent"}, `planning.replan "Agent" is not a replan policy`},
		{Planning{Schedule: "61 * * * *"}, `planning.schedule "61 * * * *": the minute "61" is outside 0-59`},
		{Planning{Schedule: "hourly"}, `planning.schedule "hourly" has 1 fields, not five`},
	} {
		got := c.p.Errors()
		if len(got) != 1 || !strings.Contains(got[0], c.want) {
			t.Errorf("%+v: got %q, want one error saying %q", c.p, got, c.want)
		}
	}
	if _, err := (Planning{Replan: "sometimes"}).ReplanPolicy(); err == nil {
		t.Error("a bad replan policy is accepted")
	}
	if s, err := (Planning{Schedule: "61 * * * *"}).PlanSchedule(); err == nil || s != nil {
		t.Errorf("a bad schedule is accepted: %v", s)
	}
}

// S-0208: planning.agent reads as agent does, is checked as agent is under
// its own name, and the planner's agent is it merged over the project's
// agent, field by field, config key by key, and role by role.
func TestPlanningAgent(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"agent:\n  harness: claude-code\n  model: claude-opus-5-5\n  config:\n    effort: high\n    max_budget_usd: \"10\"\n  roles:\n    explore:\n      model: haiku\n" +
		"planning:\n  currency: EUR\n  agent:\n    model: claude-sonnet-5\n    config:\n      effort: medium\n    roles:\n      verify:\n        model: sonnet\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if errs := m.Planning.Errors(); len(errs) != 0 {
		t.Errorf("valid planning.agent: %v", errs)
	}
	got := m.PlanningAgent()
	if s := got.String(); s != "claude-code, claude-sonnet-5, effort=medium, max_budget_usd=10; explore: haiku; verify: sonnet" {
		t.Errorf("planner's agent: %s", s)
	}
	if m.Agent.Model != "claude-opus-5-5" || m.Agent.Config["effort"] != "high" || len(m.Agent.Roles) != 1 {
		t.Errorf("the project's agent was changed by the merge: %+v", m.Agent)
	}
	if a := (Manifest{Agent: m.Agent}).PlanningAgent(); !a.Same(m.Agent) {
		t.Errorf("no planning.agent: %v", a)
	}
	if a := (Manifest{Planning: Planning{Agent: m.Planning.Agent}}).PlanningAgent(); !a.Same(m.Planning.Agent) {
		t.Errorf("no project agent: %v", a)
	}
	if a := (Manifest{}).PlanningAgent(); a != nil {
		t.Errorf("neither: %v", a)
	}
	for _, c := range []struct {
		a    *Agent
		want string
	}{
		{&Agent{Harness: "Claude Code"}, `planning.agent harness "Claude Code" is not a name such as claude-code`},
		{&Agent{Model: "has space"}, `planning.agent model "has space" is not a model ID`},
		{&Agent{Config: map[string]string{"k": "two\nlines"}}, "planning.agent config k spans lines"},
		{&Agent{Roles: map[string]Role{"verify": {}}}, "planning.agent role verify sets nothing"},
		{&Agent{Roles: map[string]Role{"verify": {Model: "has space"}}}, `planning.agent role verify model "has space"`},
	} {
		if got := strings.Join((Planning{Agent: c.a}).Errors(), "; "); !strings.Contains(got, c.want) {
			t.Errorf("%+v: got %q, want %q", c.a, got, c.want)
		}
	}
}

// S-0181: a flai below the manifest's minimum says so, naming the version
// needed and the upgrade, before anything else is read; a dev build and a
// release at or above it read the project.
func TestMinimumFlai(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	base := "version: 1\nname: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(path, []byte(base+"flai:\n  minimum: 1.27.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	was := buildinfo.Version
	t.Cleanup(func() { buildinfo.Version = was })

	buildinfo.Version = "1.26.4"
	_, err := Load(path)
	var old *TooOldError
	if !errors.As(err, &old) || old.Minimum != "1.27.0" || old.Running != "1.26.4" ||
		!strings.Contains(err.Error(), "needs flai 1.27.0 or newer, and this is flai 1.26.4") || !strings.Contains(err.Error(), "flai host upgrade") {
		t.Fatalf("a flai below the minimum: %v", err)
	}
	for _, v := range []string{"1.27.0", "1.30.1", "dev"} {
		buildinfo.Version = v
		if m, err := Load(path); err != nil || m.Flai.Minimum != "1.27.0" {
			t.Errorf("flai %s: %v %+v", v, err, m.Flai)
		}
	}
	_ = os.WriteFile(path, []byte(base+"flai:\n  minimum: soon\n"), 0o644)
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), `flai.minimum "soon" is not a release version`) {
		t.Errorf("a minimum that is not a version: %v", err)
	}
}

func TestSetMinimum(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	base := "version: 1\nname: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	for _, c := range []struct{ before, after string }{
		{base, base + "flai:\n  minimum: 1.27.0\n"},
		{base + "flai:\n  minimum: 1.26.0\n", base + "flai:\n  minimum: 1.27.0\n"},
		{"flai:\n  other: x\n" + base, "flai:\n  minimum: 1.27.0\n  other: x\n" + base},
	} {
		_ = os.WriteFile(path, []byte(c.before), 0o644)
		if err := SetMinimum(path, "1.27.0"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != c.after {
			t.Errorf("from:\n%s\ngot:\n%s\nwant:\n%s", c.before, got, c.after)
		}
		if m, err := Load(path); err != nil || m.Flai.Minimum != "1.27.0" {
			t.Errorf("reads back: %v %+v", err, m.Flai)
		}
	}
}
