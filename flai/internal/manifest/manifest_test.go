package manifest

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/goccy/go-yaml"

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

// S-0217: orchestration.policy and orchestration.release read from the
// manifest; unset, the order policy is fifo and the release policy
// judgement, an unknown key beside them is ignored, and each bad value is an
// error naming the key and what to write.
func TestOrchestration(t *testing.T) {
	load := func(t *testing.T, block string) Manifest {
		t.Helper()
		p := filepath.Join(t.TempDir(), File)
		body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" + block
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		m, err := Load(p)
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	o := load(t, "orchestration:\n  policy: wsjf\n  later: true\n  release:\n    policy: threshold\n    value: 500\n    count: 5\n").Orchestration
	if o.PolicyOrDefault() != OrderWSJF || o.Release.PolicyOrDefault() != ReleaseThreshold || o.Release.Value == nil || *o.Release.Value != 500 || o.Release.Count == nil || *o.Release.Count != 5 {
		t.Errorf("read %+v", o)
	}
	if errs := o.Errors(); len(errs) != 0 {
		t.Errorf("valid threshold: %v", errs)
	}
	o = load(t, "orchestration:\n  policy: cod\n  release:\n    policy: theme\n    epic: E-0012\n").Orchestration
	if o.PolicyOrDefault() != OrderCOD || o.Release.PolicyOrDefault() != ReleaseTheme || o.Release.Epic != "E-0012" {
		t.Errorf("read %+v", o)
	}
	if errs := o.Errors(); len(errs) != 0 {
		t.Errorf("valid theme: %v", errs)
	}
	if o = load(t, "").Orchestration; o.PolicyOrDefault() != OrderFIFO || o.Release.PolicyOrDefault() != ReleaseJudgement || o.Release.WholeEpics || len(o.Errors()) != 0 {
		t.Errorf("unset: %q, %q, whole_epics %v, %v", o.PolicyOrDefault(), o.Release.PolicyOrDefault(), o.Release.WholeEpics, o.Errors())
	}
	if o = load(t, "orchestration:\n  release:\n    policy: theme\n    tag: cli\n").Orchestration; o.Release.WholeEpics {
		t.Errorf("whole_epics unset is off: %+v", o.Release)
	}
	if o = load(t, "orchestration:\n  release:\n    whole_epics: true\n").Orchestration; !o.Release.WholeEpics || o.Release.PolicyOrDefault() != ReleaseJudgement || len(o.Errors()) != 0 {
		t.Errorf("whole_epics read, under the default policy: %+v, %v", o.Release, o.Errors())
	}
	if len(OrderPolicies) != 4 || len(ReleasePolicies) != 3 {
		t.Errorf("policies %v, %v", OrderPolicies, ReleasePolicies)
	}
	for _, p := range OrderPolicies {
		if errs := (Orchestration{Policy: p}).Errors(); len(errs) != 0 {
			t.Errorf("%q: %v", p, errs)
		}
	}

	zero, five, neg, inf := 0.0, 5, -1.0, math.Inf(1)
	negCount := -2
	for _, r := range []Release{
		{Policy: ReleaseJudgement},
		{Policy: ReleaseThreshold, Value: &zero},
		{Policy: ReleaseThreshold, Count: &five},
		{Policy: ReleaseTheme, Tag: "cli"},
	} {
		if errs := (Orchestration{Release: r}).Errors(); len(errs) != 0 {
			t.Errorf("%+v: %v", r, errs)
		}
	}
	for _, c := range []struct {
		o    Orchestration
		want string
	}{
		{Orchestration{Policy: "lifo"}, `orchestration.policy "lifo" is not an order policy; write cod`},
		{Orchestration{Policy: "WSJF"}, `orchestration.policy "WSJF" is not an order policy`},
		{Orchestration{Release: Release{Policy: "weekly"}}, `orchestration.release.policy "weekly" is not a release policy; write judgement`},
		{Orchestration{Release: Release{Policy: ReleaseThreshold}}, "orchestration.release is a threshold with neither value nor count"},
		{Orchestration{Release: Release{Policy: ReleaseThreshold, Value: &neg}}, "orchestration.release.value -1 is not an amount of zero or more"},
		{Orchestration{Release: Release{Policy: ReleaseThreshold, Value: &inf}}, "orchestration.release.value +Inf is not an amount of zero or more"},
		{Orchestration{Release: Release{Policy: ReleaseThreshold, Count: &negCount}}, "orchestration.release.count -2 is not a number of zero or more"},
		{Orchestration{Release: Release{Policy: ReleaseTheme}}, "orchestration.release is a theme with neither epic nor tag"},
		{Orchestration{Release: Release{Policy: ReleaseTheme, Epic: "E-0001", Tag: "cli"}}, "orchestration.release is a theme with both epic and tag"},
		{Orchestration{Release: Release{Policy: ReleaseTheme, Epic: "S-0001"}}, `orchestration.release.epic "S-0001" is not an epic ID`},
	} {
		got := c.o.Errors()
		if len(got) != 1 || !strings.Contains(got[0], c.want) {
			t.Errorf("%+v: got %q, want one error saying %q", c.o, got, c.want)
		}
	}
}

// S-0218: orchestration.permissions read from the manifest, each off when
// unset; answer_threads reads in each of its values; and a key that names no
// permission, or an answer_threads outside its values, is an error naming it
// and what to write.
func TestOrchestrationPermissions(t *testing.T) {
	load := func(t *testing.T, block string) Manifest {
		t.Helper()
		p := filepath.Join(t.TempDir(), File)
		body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" + block
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		m, err := Load(p)
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	all := "orchestration:\n  policy: cod\n  permissions:\n    plan_backlog_epics: true\n    plan_backlog_stories: true\n    finalize_drafts: true\n    promote_to_ready: true\n" +
		"    order_ready: true\n    answer_threads: autonomous\n    accept_reviews: true\n    publish: true\n"
	o := load(t, all).Orchestration
	want := Permissions{PlanBacklogEpics: true, PlanBacklogStories: true, FinalizeDrafts: true, PromoteToReady: true, OrderReady: true, AnswerThreads: AnswerAutonomous, AcceptReviews: true, Publish: true}
	if !reflect.DeepEqual(o.Permissions, want) {
		t.Errorf("read %+v, want %+v", o.Permissions, want)
	}
	if errs := o.Errors(); len(errs) != 0 {
		t.Errorf("every permission on: %v", errs)
	}
	for _, name := range PermissionNames {
		if !o.Permissions.Allows(name) {
			t.Errorf("%s is set and not allowed", name)
		}
	}

	for _, block := range []string{"", "orchestration:\n  policy: cod\n", "orchestration:\n  permissions: {}\n"} {
		o := load(t, block).Orchestration
		if o.Permissions.AnswerMode() != AnswerOff || len(o.Errors()) != 0 {
			t.Errorf("%q: answer mode %q, %v", block, o.Permissions.AnswerMode(), o.Errors())
		}
		for _, name := range PermissionNames {
			if o.Permissions.Allows(name) {
				t.Errorf("%q: %s is on unset", block, name)
			}
		}
	}
	if (Permissions{Publish: true}).Allows("publsh") {
		t.Error("a name that is no permission is allowed")
	}

	for _, mode := range AnswerModes {
		o := load(t, "orchestration:\n  permissions:\n    answer_threads: "+mode+"\n").Orchestration
		if o.Permissions.AnswerMode() != mode || o.Permissions.Allows(PermitAnswerThreads) != (mode != AnswerOff) || len(o.Errors()) != 0 {
			t.Errorf("%s: mode %q, allowed %v, %v", mode, o.Permissions.AnswerMode(), o.Permissions.Allows(PermitAnswerThreads), o.Errors())
		}
	}

	o = load(t, "orchestration:\n  permissions:\n    publsh: true\n    order_ready: true\n    accept: false\n").Orchestration
	if !o.Permissions.OrderReady || o.Permissions.Publish {
		t.Errorf("read %+v", o.Permissions)
	}
	got := o.Errors()
	if len(got) != 2 || !strings.Contains(got[0], `orchestration.permissions has no permission "accept"; write one of plan_backlog_epics,`) ||
		!strings.Contains(got[1], `orchestration.permissions has no permission "publsh"`) || !strings.HasSuffix(got[1], "publish, or remove it") {
		t.Errorf("unknown keys: %q", got)
	}
	for _, bad := range []string{"yes", "Autonomous", "always"} {
		got := (Orchestration{Permissions: Permissions{AnswerThreads: bad}}).Errors()
		if len(got) != 1 || !strings.Contains(got[0], `orchestration.permissions.answer_threads "`+bad+`" is not a way of answering threads; write off`) {
			t.Errorf("%q: %q", bad, got)
		}
	}
}

// S-0218: orchestration.agent reads as agent does, is checked as agent is
// under its own name, and the orchestrator's agent is it merged over the
// project's agent, field by field, config key by key, and role by role.
func TestOrchestrationAgent(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"agent:\n  harness: claude-code\n  model: claude-opus-5-5\n  config:\n    effort: high\n    max_budget_usd: \"10\"\n  roles:\n    explore:\n      model: haiku\n" +
		"orchestration:\n  policy: wsjf\n  agent:\n    model: claude-sonnet-5\n    config:\n      effort: medium\n    roles:\n      verify:\n        model: sonnet\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if errs := m.Orchestration.Errors(); len(errs) != 0 {
		t.Errorf("valid orchestration.agent: %v", errs)
	}
	got := m.OrchestrationAgent()
	if s := got.String(); s != "claude-code, claude-sonnet-5, effort=medium, max_budget_usd=10; explore: haiku; verify: sonnet" {
		t.Errorf("orchestrator's agent: %s", s)
	}
	if m.Agent.Model != "claude-opus-5-5" || m.Agent.Config["effort"] != "high" || len(m.Agent.Roles) != 1 {
		t.Errorf("the project's agent was changed by the merge: %+v", m.Agent)
	}
	if a := m.PlanningAgent(); !a.Same(m.Agent) {
		t.Errorf("orchestration.agent reached the planner: %v", a)
	}
	if a := (Manifest{Agent: m.Agent}).OrchestrationAgent(); !a.Same(m.Agent) {
		t.Errorf("no orchestration.agent: %v", a)
	}
	if a := (Manifest{Orchestration: Orchestration{Agent: m.Orchestration.Agent}}).OrchestrationAgent(); !a.Same(m.Orchestration.Agent) {
		t.Errorf("no project agent: %v", a)
	}
	if a := (Manifest{}).OrchestrationAgent(); a != nil {
		t.Errorf("neither: %v", a)
	}
	for _, c := range []struct {
		a    *Agent
		want string
	}{
		{&Agent{Harness: "Claude Code"}, `orchestration.agent harness "Claude Code" is not a name such as claude-code`},
		{&Agent{Model: "has space"}, `orchestration.agent model "has space" is not a model ID`},
		{&Agent{Config: map[string]string{"k": "two\nlines"}}, "orchestration.agent config k spans lines"},
		{&Agent{Roles: map[string]Role{"verify": {}}}, "orchestration.agent role verify sets nothing"},
		{&Agent{Roles: map[string]Role{"verify": {Model: "has space"}}}, `orchestration.agent role verify model "has space"`},
	} {
		if got := strings.Join((Orchestration{Agent: c.a}).Errors(), "; "); !strings.Contains(got, c.want) {
			t.Errorf("%+v: got %q, want %q", c.a, got, c.want)
		}
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

// S-0223: analysis.schedule reads from the manifest as planning.schedule
// does; unset, there is no schedule and no analysis.agent, and a schedule
// that does not parse, or never comes round, is an error naming the key.
func TestAnalysisSchedule(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\nanalysis:\n  schedule: \"0 6 * * 1\"\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	s, err := m.Analysis.AnalysisSchedule()
	if err != nil || s == nil || s.String() != "0 6 * * 1" {
		t.Fatalf("schedule %v, %v", s, err)
	}
	at := time.Date(2026, 10, 4, 9, 10, 0, 0, time.UTC) // a Sunday
	if next := s.Next(at); !next.Equal(time.Date(2026, 10, 5, 6, 0, 0, 0, time.UTC)) {
		t.Errorf("next after %v: %v", at, next)
	}
	if errs := m.Analysis.Errors(); len(errs) != 0 {
		t.Errorf("valid analysis: %v", errs)
	}
	if ps, err := m.Planning.PlanSchedule(); err != nil || ps != nil {
		t.Errorf("analysis.schedule reached the planner: %v, %v", ps, err)
	}

	var unset Analysis
	if s, err := unset.AnalysisSchedule(); err != nil || s != nil || unset.Agent != nil || len(unset.Errors()) != 0 {
		t.Errorf("unset: schedule %v, %v, agent %v, %v", s, err, unset.Agent, unset.Errors())
	}
	if s, err := (Analysis{Schedule: "  "}).AnalysisSchedule(); err != nil || s != nil {
		t.Errorf("blank schedule %v, %v", s, err)
	}
	s, err = (Analysis{Schedule: "daily"}).AnalysisSchedule()
	if err != nil || s == nil || s.String() != "daily" {
		t.Fatalf("daily: %v, %v", s, err)
	}

	for _, c := range []struct {
		a    Analysis
		want string
	}{
		{Analysis{Schedule: "61 * * * *"}, `analysis.schedule "61 * * * *": the minute "61" is outside 0-59`},
		{Analysis{Schedule: "weekly"}, `analysis.schedule "weekly" has 1 fields, not five`},
		{Analysis{Schedule: "0 0 31 2 *"}, `analysis.schedule "0 0 31 2 *" never comes round`},
	} {
		got := c.a.Errors()
		if len(got) != 1 || !strings.Contains(got[0], c.want) {
			t.Errorf("%+v: got %q, want one error saying %q", c.a, got, c.want)
		}
		if s, err := c.a.AnalysisSchedule(); err == nil || s != nil {
			t.Errorf("%+v: a bad schedule is accepted: %v", c.a, s)
		}
	}
}

// S-0223: analysis.agent reads as agent does, is checked as agent is under
// its own name, and the analyzer's agent is it merged over the project's
// agent, field by field, config key by key, and role by role.
func TestAnalysisAgent(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"agent:\n  harness: claude-code\n  model: claude-opus-5-5\n  config:\n    effort: high\n    max_budget_usd: \"10\"\n  roles:\n    explore:\n      model: haiku\n" +
		"analysis:\n  schedule: daily\n  agent:\n    model: claude-sonnet-5\n    config:\n      effort: medium\n    roles:\n      verify:\n        model: sonnet\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if errs := m.Analysis.Errors(); len(errs) != 0 {
		t.Errorf("valid analysis.agent: %v", errs)
	}
	got := m.AnalysisAgent()
	if s := got.String(); s != "claude-code, claude-sonnet-5, effort=medium, max_budget_usd=10; explore: haiku; verify: sonnet" {
		t.Errorf("analyzer's agent: %s", s)
	}
	if m.Agent.Model != "claude-opus-5-5" || m.Agent.Config["effort"] != "high" || len(m.Agent.Roles) != 1 {
		t.Errorf("the project's agent was changed by the merge: %+v", m.Agent)
	}
	if a := m.PlanningAgent(); !a.Same(m.Agent) {
		t.Errorf("analysis.agent reached the planner: %v", a)
	}
	if a := (Manifest{Agent: m.Agent}).AnalysisAgent(); !a.Same(m.Agent) {
		t.Errorf("no analysis.agent: %v", a)
	}
	if a := (Manifest{Analysis: Analysis{Agent: m.Analysis.Agent}}).AnalysisAgent(); !a.Same(m.Analysis.Agent) {
		t.Errorf("no project agent: %v", a)
	}
	if a := (Manifest{}).AnalysisAgent(); a != nil {
		t.Errorf("neither: %v", a)
	}
	for _, c := range []struct {
		a    *Agent
		want string
	}{
		{&Agent{Harness: "Claude Code"}, `analysis.agent harness "Claude Code" is not a name such as claude-code`},
		{&Agent{Model: "has space"}, `analysis.agent model "has space" is not a model ID`},
		{&Agent{Config: map[string]string{"k": "two\nlines"}}, "analysis.agent config k spans lines"},
		{&Agent{Roles: map[string]Role{"verify": {}}}, "analysis.agent role verify sets nothing"},
		{&Agent{Roles: map[string]Role{"verify": {Model: "has space"}}}, `analysis.agent role verify model "has space"`},
	} {
		if got := strings.Join((Analysis{Agent: c.a}).Errors(), "; "); !strings.Contains(got, c.want) {
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

// S-0273: tests is an ordered list of tiers that loads with every field,
// reads back the same through YAML and JSON, and is valid.
func TestTestsRoundTrip(t *testing.T) {
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"tests:\n" +
		"- name: unit\n" +
		"  command: [go, test, -json, \"{packages}\"]\n" +
		"  dir: flai\n" +
		"  paths: [\"flai/**\", \"!flai/testdata/**\"]\n" +
		"  format: go-test-json\n" +
		"  all_command: [go, test, -json, ./...]\n" +
		"- name: smoke\n" +
		"  command: [scripts/smoke.sh]\n" +
		"  all_only: true\n"
	m := loaded(t, []byte(body))
	want := []TestTier{
		{Name: "unit", Command: []string{"go", "test", "-json", PlaceholderPackages}, Dir: "flai",
			Paths: []string{"flai/**", "!flai/testdata/**"}, Format: FormatGoTestJSON, AllCommand: []string{"go", "test", "-json", "./..."}},
		{Name: "smoke", Command: []string{"scripts/smoke.sh"}, AllOnly: true},
	}
	if !reflect.DeepEqual(m.Tests, want) {
		t.Fatalf("tests:\n got %+v\nwant %+v", m.Tests, want)
	}
	if errs := m.TestErrors(); len(errs) != 0 {
		t.Errorf("valid tiers refused: %v", errs)
	}
	if m.Tests[1].FormatOrDefault() != FormatPlain || m.Tests[0].FormatOrDefault() != FormatGoTestJSON {
		t.Errorf("formats in effect: %q, %q", m.Tests[0].FormatOrDefault(), m.Tests[1].FormatOrDefault())
	}
	data, err := yaml.Marshal(Manifest{Tests: m.Tests})
	if err != nil {
		t.Fatal(err)
	}
	var again Manifest
	if err := yaml.Unmarshal(data, &again); err != nil || !reflect.DeepEqual(again.Tests, want) {
		t.Errorf("YAML round trip: %v\n%s\n%+v", err, data, again.Tests)
	}
	js, err := json.Marshal(m.Tests)
	if err != nil {
		t.Fatal(err)
	}
	var fromJSON []TestTier
	if err := json.Unmarshal(js, &fromJSON); err != nil || !reflect.DeepEqual(fromJSON, want) {
		t.Errorf("JSON round trip: %v\n%s", err, js)
	}
	if !strings.Contains(string(js), `"all_only":true`) || !strings.Contains(string(js), `"all_command":`) {
		t.Errorf("JSON keys are not snake_case: %s", js)
	}
}

// S-0273: a bad tier is named by its index, its field, and its name, with the
// reason and what to write instead.
func TestTestErrors(t *testing.T) {
	ok := TestTier{Name: "unit", Command: []string{"scripts/test.sh"}, Paths: []string{"**"}}
	with := func(f func(*TestTier)) TestTier {
		tier := ok
		f(&tier)
		return tier
	}
	for _, c := range []struct {
		name  string
		tiers []TestTier
		want  string
	}{
		{"no name", []TestTier{with(func(x *TestTier) { x.Name = " " })}, "tests[0].name is empty"},
		{"a name twice", []TestTier{ok, ok}, `tests[1].name "unit" is tests[0]'s name too`},
		{"no command", []TestTier{with(func(x *TestTier) { x.Command = nil })}, `tests[0].command (tier "unit") is empty`},
		{"no program", []TestTier{with(func(x *TestTier) { x.Command = []string{"", "x"} })}, `tests[0].command (tier "unit") has no program first`},
		{"a placeholder first", []TestTier{with(func(x *TestTier) { x.Command = []string{"{files}"} })}, "begins with the placeholder {files}"},
		{"another placeholder", []TestTier{with(func(x *TestTier) { x.Command = []string{"go", "{story}"} })}, `tests[0].command (tier "unit") argument {story} is not a placeholder flai fills`},
		{"a placeholder inside", []TestTier{with(func(x *TestTier) { x.Command = []string{"vitest", "--dir={files}"} })}, `argument "--dir={files}" holds a placeholder inside it`},
		{"a placeholder under all", []TestTier{with(func(x *TestTier) { x.AllCommand = []string{"go", "test", "{packages}"} })}, `tests[0].all_command (tier "unit") argument {packages} is a placeholder`},
		{"an absolute dir", []TestTier{with(func(x *TestTier) { x.Dir = "/srv/flai" })}, `tests[0].dir (tier "unit") "/srv/flai" is an absolute path`},
		{"a dir outside", []TestTier{with(func(x *TestTier) { x.Dir = "../other" })}, `"../other" leaves the repository`},
		{"a dir not clean", []TestTier{with(func(x *TestTier) { x.Dir = "flai/./web/" })}, "is not a clean path"},
		{"no paths", []TestTier{with(func(x *TestTier) { x.Paths = nil })}, `tests[0].paths (tier "unit") is empty`},
		{"only exclusions", []TestTier{with(func(x *TestTier) { x.Paths = []string{"!flai/testdata"} })}, "has no pattern that selects a path"},
		{"a bad pattern", []TestTier{with(func(x *TestTier) { x.Paths = []string{"flai/**", "!../x"} })}, `tests[0].paths (tier "unit") pattern "!../x" has a .. segment`},
		{"a bad format", []TestTier{with(func(x *TestTier) { x.Format = "junit" })}, `tests[0].format (tier "unit") "junit" is not a format flai reads; write go-test-json, vitest-json, golangci-json, gofmt-list, plain`},
	} {
		t.Run(c.name, func(t *testing.T) {
			errs := Manifest{Tests: c.tiers}.TestErrors()
			if got := strings.Join(errs, "\n"); !strings.Contains(got, c.want) {
				t.Errorf("got:\n%s\nwant ...%s...", got, c.want)
			}
			if p := splitProblem(errs[0]); !strings.HasPrefix(p.Field, "tests[") {
				t.Errorf("%q is not split into a field and a reason: %+v", errs[0], p)
			}
		})
	}
	for _, tier := range []TestTier{
		with(func(x *TestTier) { x.Paths, x.AllOnly = nil, true }),
		with(func(x *TestTier) { x.Dir = "flaiover/web" }),
		with(func(x *TestTier) { x.Command = []string{"go", "vet", "{packages}", "{files}"} }),
		with(func(x *TestTier) { x.Command = []string{"sh", "-c", "echo ${HOME}"} }),
	} {
		if errs := (Manifest{Tests: []TestTier{tier}}).TestErrors(); len(errs) != 0 {
			t.Errorf("%+v refused: %v", tier, errs)
		}
	}
}

// S-0273: with no tests key the project has one plain tier running
// scripts/test.sh, when it is there, for every path; an empty list has none,
// and a list with a bad tier is refused.
func TestTestTiers(t *testing.T) {
	script := fstest.MapFS{DefaultTestScript: {Data: []byte("#!/bin/sh\n"), Mode: 0o755}}
	tiers, err := Manifest{}.TestTiers(script)
	want := []TestTier{{Name: "test", Command: []string{"scripts/test.sh"}, Paths: []string{"**"}, Format: FormatPlain}}
	if err != nil || !reflect.DeepEqual(tiers, want) {
		t.Errorf("the default: %v %+v", err, tiers)
	}
	for name, fsys := range map[string]fstest.MapFS{
		"no script":        {},
		"a folder instead": {"scripts/test.sh/x": {Data: nil}},
		"another script":   {"scripts/check.sh": {Data: nil}},
	} {
		if tiers, err := (Manifest{}).TestTiers(fsys); err != nil || tiers != nil {
			t.Errorf("%s: %v %+v", name, err, tiers)
		}
	}
	m := loaded(t, []byte("version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\ntests: []\n"))
	if tiers, err := m.TestTiers(script); err != nil || len(tiers) != 0 || m.Tests == nil {
		t.Errorf("tests: [] means none: %v %+v (%#v)", err, tiers, m.Tests)
	}
	declared := []TestTier{{Name: "unit", Command: []string{"make", "test"}, Paths: []string{"src"}}}
	if tiers, err := (Manifest{Tests: declared}).TestTiers(script); err != nil || !reflect.DeepEqual(tiers, declared) {
		t.Errorf("declared tiers: %v %+v", err, tiers)
	}
	bad := []TestTier{{Name: "unit", Command: []string{"make", "test"}}}
	if _, err := (Manifest{Tests: bad}).TestTiers(script); err == nil || !strings.Contains(err.Error(), "tests[0].paths") {
		t.Errorf("a bad tier: %v", err)
	}
}
