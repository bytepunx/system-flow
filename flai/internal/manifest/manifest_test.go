package manifest

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
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

// S-0239: dashboard.allow_unsigned in the manifest is in force over the host
// configuration's either way, and its absence leaves the host's.
func TestAllowUnsigned(t *testing.T) {
	head := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
	silent := loaded(t, []byte(head+"dashboard:\n  port: 4242\n"))
	on := loaded(t, []byte(head+"dashboard:\n  allow_unsigned: true\n"))
	off := loaded(t, []byte(head+"dashboard:\n  allow_unsigned: false\n"))
	cases := []struct {
		name string
		d    Dashboard
		host bool
		want bool
	}{
		{"silent, host off", silent.Dashboard, false, false},
		{"silent, host on", silent.Dashboard, true, true},
		{"true over false", on.Dashboard, false, true},
		{"false over true", off.Dashboard, true, false},
		{"true over true", on.Dashboard, true, true},
		{"false over false", off.Dashboard, false, false},
	}
	for _, c := range cases {
		if got := c.d.AllowUnsignedOr(c.host); got != c.want {
			t.Errorf("%s: AllowUnsignedOr(%v) = %v, want %v", c.name, c.host, got, c.want)
		}
	}
	for _, m := range []Manifest{on, off} {
		data, err := yaml.Marshal(Manifest{Dashboard: m.Dashboard})
		if err != nil {
			t.Fatal(err)
		}
		var again Manifest
		if err := yaml.Unmarshal(data, &again); err != nil || !reflect.DeepEqual(again.Dashboard, m.Dashboard) {
			t.Errorf("YAML round trip: %v\n%s\n%+v", err, data, again.Dashboard)
		}
	}
	data, err := yaml.Marshal(Manifest{Dashboard: silent.Dashboard})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "allow_unsigned") {
		t.Errorf("unset allow_unsigned is written:\n%s", data)
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
		"    order_ready: true\n    answer_threads: autonomous\n    confirm_after: 48h\n    accept_reviews: true\n    publish: true\n    edit_subprojects: true\n"
	o := load(t, all).Orchestration
	want := Permissions{PlanBacklogEpics: true, PlanBacklogStories: true, FinalizeDrafts: true, PromoteToReady: true, OrderReady: true, AnswerThreads: AnswerAutonomous, ConfirmAfter: "48h", AcceptReviews: true, Publish: true, EditSubprojects: true}
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

	for _, block := range []string{"", "orchestration:\n  policy: cod\n", "orchestration:\n  permissions: {}\n", "orchestration:\n  permissions:\n    confirm_after: \"0\"\n"} {
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
		!strings.Contains(got[1], `orchestration.permissions has no permission "publsh"`) || !strings.HasSuffix(got[1], "publish, edit_subprojects, or remove it") {
		t.Errorf("unknown keys: %q", got)
	}
	for _, bad := range []string{"yes", "Autonomous", "always"} {
		got := (Orchestration{Permissions: Permissions{AnswerThreads: bad}}).Errors()
		if len(got) != 1 || !strings.Contains(got[0], `orchestration.permissions.answer_threads "`+bad+`" is not a way of answering threads; write off`) {
			t.Errorf("%q: %q", bad, got)
		}
	}
}

// S-0502: orchestration.permissions.confirm_after is a duration, off unset or
// 0; a value that is not a duration, or a negative one, is an error naming the
// key and an example; and it is a permission, not a key that names none
// (ADR-0206).
func TestConfirmAfter(t *testing.T) {
	for in, want := range map[string]time.Duration{"": 0, "0": 0, "0s": 0, "48h": 48 * time.Hour, "90m": 90 * time.Minute} {
		p := Permissions{ConfirmAfter: in}
		d, err := p.ConfirmAfterDuration()
		if err != nil || d != want || p.Allows(PermitConfirmAfter) != (want > 0) || len((Orchestration{Permissions: p}).Errors()) != 0 {
			t.Errorf("%q: %v, %v, allowed %v, %v", in, d, err, p.Allows(PermitConfirmAfter), (Orchestration{Permissions: p}).Errors())
		}
	}
	for in, want := range map[string]string{
		"soon": `orchestration.permissions.confirm_after "soon" is not a duration; write one such as 48h`,
		"-1h":  `orchestration.permissions.confirm_after "-1h" is negative; write a duration such as 48h`,
	} {
		p := Permissions{ConfirmAfter: in}
		if _, err := p.ConfirmAfterDuration(); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: %v, want %q", in, err, want)
		}
		if p.Allows(PermitConfirmAfter) {
			t.Errorf("%q is allowed", in)
		}
		if got := (Orchestration{Permissions: p}).Errors(); len(got) != 1 || !strings.Contains(got[0], want) {
			t.Errorf("%q: %q, want one error saying %q", in, got, want)
		}
	}

	var m Manifest
	if err := yaml.Unmarshal([]byte("orchestration:\n  permissions:\n    confirm_after: 48h\n"), &m); err != nil {
		t.Fatal(err)
	}
	if m.Orchestration.Permissions.ConfirmAfter != "48h" || len(m.Orchestration.Errors()) != 0 {
		t.Errorf("read %+v, %q", m.Orchestration.Permissions, m.Orchestration.Errors())
	}
}

// S-0416, ADR-0207: roles.<name>.permissions read from the manifest for each
// role, each off when unset and no role's but its own; and a key that names
// no permission, or an answer_threads outside its values, is an error naming
// the role's key and what to write.
func TestRolePermissions(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"roles:\n  reviewer:\n    permissions:\n      promote_to_ready: true\n      answer_threads: recommend\n  scout: {}\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	reviewer := m.Roles["reviewer"].Permissions
	if !reviewer.Allows(PermitPromoteToReady) || reviewer.AnswerMode() != AnswerRecommend || reviewer.Allows(PermitPublish) {
		t.Errorf("reviewer: %+v", reviewer)
	}
	for _, name := range PermissionNames {
		if m.Roles["scout"].Permissions.Allows(name) || m.Roles["absent"].Permissions.Allows(name) || m.Orchestration.Permissions.Allows(name) {
			t.Errorf("%s is on unset, or another role's", name)
		}
	}
	if errs := m.RoleErrors(); len(errs) != 0 {
		t.Errorf("valid roles: %v", errs)
	}

	m.Roles = map[string]RoleSettings{"scout": {Permissions: Permissions{AnswerThreads: "always", unknown: []string{"publsh"}}}}
	got := m.RoleErrors()
	if len(got) != 2 || !strings.HasPrefix(got[0], `roles.scout.permissions has no permission "publsh"; write one of plan_backlog_epics,`) ||
		!strings.HasPrefix(got[1], `roles.scout.permissions.answer_threads "always" is not a way of answering threads; write off`) {
		t.Errorf("errors: %q", got)
	}
	if errs := (Manifest{}).RoleErrors(); len(errs) != 0 {
		t.Errorf("no roles: %v", errs)
	}
}

// S-0392: a manifest without acceptance loads with each key at its default
// and marshals without the block; each value of each key loads, a forge and a
// merge method under local among them; and a value outside a key's values is
// refused at load, naming the key, the value, and what to write.
func TestAcceptance(t *testing.T) {
	head := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
	a := loaded(t, []byte(head)).Acceptance
	if a != (Acceptance{}) || a.ModeOrDefault() != AcceptLocal || a.PullRequest() || a.ForgeOrDefault() != ForgeGitHub || a.MergeMethodOrDefault() != MergeMerge || len(a.Errors()) != 0 {
		t.Errorf("unset: %+v, mode %q, forge %q, merge method %q, %v", a, a.ModeOrDefault(), a.ForgeOrDefault(), a.MergeMethodOrDefault(), a.Errors())
	}
	for _, format := range []func(Manifest) ([]byte, error){
		func(m Manifest) ([]byte, error) { return yaml.Marshal(m) },
		func(m Manifest) ([]byte, error) { return json.Marshal(m) },
	} {
		if data, err := format(Manifest{}); err != nil || strings.Contains(string(data), "acceptance") {
			t.Errorf("an unset acceptance block is written: %v\n%s", err, data)
		}
	}
	if !slices.Equal(AcceptanceModes, []string{"local", "pull-request"}) || !slices.Equal(Forges, []string{"github"}) || !slices.Equal(MergeMethods, []string{"merge", "squash", "rebase"}) {
		t.Errorf("values %v, %v, %v", AcceptanceModes, Forges, MergeMethods)
	}

	for _, mode := range AcceptanceModes {
		for _, method := range MergeMethods {
			got := loaded(t, []byte(head+"acceptance:\n  mode: "+mode+"\n  forge: github\n  merge_method: "+method+"\n")).Acceptance
			if got.ModeOrDefault() != mode || got.PullRequest() != (mode == AcceptPullRequest) || got.ForgeOrDefault() != ForgeGitHub || got.MergeMethodOrDefault() != method {
				t.Errorf("%s, %s: read %+v", mode, method, got)
			}
		}
	}
	if a = loaded(t, []byte(head+"acceptance:\n  mode: pull-request\n")).Acceptance; !a.PullRequest() || a.ForgeOrDefault() != ForgeGitHub || a.MergeMethodOrDefault() != MergeMerge {
		t.Errorf("pull-request alone: forge %q, merge method %q", a.ForgeOrDefault(), a.MergeMethodOrDefault())
	}

	p := filepath.Join(t.TempDir(), File)
	for _, c := range []struct{ block, want string }{
		{"  mode: remote\n", `acceptance.mode "remote" is not an acceptance mode; write local (merge on this machine) or pull-request (merge through a pull request on the forge), or remove it for local`},
		{"  mode: pull_request\n", `acceptance.mode "pull_request" is not an acceptance mode`},
		{"  forge: gitlab\n", `acceptance.forge "gitlab" is not a forge flai reaches; write github, or remove it for github`},
		{"  mode: pull-request\n  forge: GitHub\n", `acceptance.forge "GitHub" is not a forge flai reaches`},
		{"  merge_method: fast-forward\n", `acceptance.merge_method "fast-forward" is not a merge method; write merge (a merge commit), squash (one commit), or rebase (the commits onto the base), or remove it for merge`},
	} {
		if err := os.WriteFile(p, []byte(head+"acceptance:\n"+c.block), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := Load(p)
		if err == nil || !strings.HasPrefix(err.Error(), p+": ") || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: got %v, want %s", c.block, err, c.want)
		}
	}
	if errs := (Acceptance{Mode: "x", Forge: "y", MergeMethod: "z"}).Errors(); len(errs) != 3 {
		t.Errorf("one error a key: %q", errs)
	}
}

// S-0399, ADR-0173: a manifest without merge_queue batches 4 stories at most
// and marshals without the block; a batch_max of one or more loads as given;
// zero, a negative number, and a value that is not a whole number are
// refused at load, naming the key, the value, and what to write.
func TestMergeQueue(t *testing.T) {
	head := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
	m := loaded(t, []byte(head))
	if m.MergeQueue != (MergeQueue{}) || m.BatchMax() != 4 || DefaultBatchMax != 4 || len(m.MergeQueue.Errors()) != 0 {
		t.Errorf("unset: %+v, batch max %d, %v", m.MergeQueue, m.BatchMax(), m.MergeQueue.Errors())
	}
	if loaded(t, []byte(head+"merge_queue:\n")).BatchMax() != DefaultBatchMax {
		t.Error("an empty merge_queue block is not the default")
	}
	for _, format := range []func(Manifest) ([]byte, error){
		func(m Manifest) ([]byte, error) { return yaml.Marshal(m) },
		func(m Manifest) ([]byte, error) { return json.Marshal(m) },
	} {
		if data, err := format(Manifest{}); err != nil || strings.Contains(string(data), "merge_queue") {
			t.Errorf("an unset merge_queue block is written: %v\n%s", err, data)
		}
	}
	for _, n := range []int{1, 4, 12} {
		if got := loaded(t, []byte(head+"merge_queue:\n  batch_max: "+strconv.Itoa(n)+"\n")).BatchMax(); got != n {
			t.Errorf("batch_max %d reads as %d", n, got)
		}
	}

	p := filepath.Join(t.TempDir(), File)
	for _, c := range []struct{ value, want string }{
		{"0", `merge_queue.batch_max 0 is not a number of one or more; write the most enqueued stories one batch of the merge queue holds, such as 4, or remove it for 4`},
		{"-2", `merge_queue.batch_max -2 is not a number of one or more`},
		{"four", `merge_queue.batch_max "four" is not a whole number; write the most enqueued stories one batch of the merge queue holds, such as 4, or remove it for 4`},
		{`"4"`, `merge_queue.batch_max "4" is not a whole number`},
		{"2.5", `merge_queue.batch_max 2.5 is not a whole number`},
		{"[4]", `merge_queue.batch_max [4] is not a whole number`},
		{"true", `merge_queue.batch_max true is not a whole number`},
	} {
		if err := os.WriteFile(p, []byte(head+"merge_queue:\n  batch_max: "+c.value+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := Load(p)
		if err == nil || !strings.HasPrefix(err.Error(), p+": ") || !strings.Contains(err.Error(), c.want) {
			t.Errorf("batch_max %s: got %v, want %s", c.value, err, c.want)
		}
	}
	if got := loaded(t, []byte(head+"merge_queue:\n  batch_max: 3.0\n")).BatchMax(); got != 3 {
		t.Errorf("batch_max 3.0, a whole number, reads as %d", got)
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

// ADR-0196: digest.agent reads as analysis.agent does, is checked under its
// own name, and the curator's agent is it merged over the project's agent,
// field by field, config key by key, and role by role, reaching no other
// strategic agent.
func TestDigestAgent(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"agent:\n  harness: claude-code\n  model: claude-opus-5-5\n  config:\n    effort: high\n    max_budget_usd: \"10\"\n  roles:\n    explore:\n      model: haiku\n" +
		"digest:\n  agent:\n    model: claude-sonnet-5\n    provider: anthropic\n    config:\n      effort: medium\n    roles:\n      verify:\n        model: sonnet\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if errs := m.Digest.Errors(); len(errs) != 0 {
		t.Errorf("valid digest.agent: %v", errs)
	}
	got := m.DigestAgent()
	if s := got.String(); s != "claude-code, claude-sonnet-5, provider anthropic, effort=medium, max_budget_usd=10; explore: haiku; verify: sonnet" {
		t.Errorf("curator's agent: %s", s)
	}
	if m.Agent.Model != "claude-opus-5-5" || m.Agent.Config["effort"] != "high" || len(m.Agent.Roles) != 1 {
		t.Errorf("the project's agent was changed by the merge: %+v", m.Agent)
	}
	for what, a := range map[string]*Agent{"planner": m.PlanningAgent(), "orchestrator": m.OrchestrationAgent(), "analyzer": m.AnalysisAgent()} {
		if !a.Same(m.Agent) {
			t.Errorf("digest.agent reached the %s: %v", what, a)
		}
	}
	if a := (Manifest{Agent: m.Agent}).DigestAgent(); !a.Same(m.Agent) {
		t.Errorf("no digest.agent: %v", a)
	}
	if a := (Manifest{Digest: Digest{Agent: m.Digest.Agent}}).DigestAgent(); !a.Same(m.Digest.Agent) {
		t.Errorf("no project agent: %v", a)
	}
	if a := (Manifest{}).DigestAgent(); a != nil {
		t.Errorf("neither: %v", a)
	}
	if errs := (Digest{}).Errors(); len(errs) != 0 {
		t.Errorf("unset digest: %v", errs)
	}
	for _, c := range []struct {
		a    *Agent
		want string
	}{
		{&Agent{Harness: "Claude Code"}, `digest.agent harness "Claude Code" is not a name such as claude-code`},
		{&Agent{Model: "has space"}, `digest.agent model "has space" is not a model ID`},
		{&Agent{Config: map[string]string{"k": "two\nlines"}}, "digest.agent config k spans lines"},
		{&Agent{Roles: map[string]Role{"verify": {}}}, "digest.agent role verify sets nothing"},
		{&Agent{Roles: map[string]Role{"verify": {Model: "has space"}}}, `digest.agent role verify model "has space"`},
	} {
		if got := strings.Join((Digest{Agent: c.a}).Errors(), "; "); !strings.Contains(got, c.want) {
			t.Errorf("%+v: got %q, want %q", c.a, got, c.want)
		}
	}
}

// ADR-0199: digest.schedule reads from the manifest as analysis.schedule
// does, reaching no other schedule; unset, there is none, and one that does
// not parse, or never comes round, is an error naming the key.
func TestDigestSchedule(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\ndigest:\n  schedule: \"0 6 * * 1\"\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	s, err := m.Digest.DigestSchedule()
	if err != nil || s == nil || s.String() != "0 6 * * 1" {
		t.Fatalf("schedule %v, %v", s, err)
	}
	at := time.Date(2026, 10, 4, 9, 10, 0, 0, time.UTC) // a Sunday
	if next := s.Next(at); !next.Equal(time.Date(2026, 10, 5, 6, 0, 0, 0, time.UTC)) {
		t.Errorf("next after %v: %v", at, next)
	}
	if errs := m.Digest.Errors(); len(errs) != 0 {
		t.Errorf("valid digest: %v", errs)
	}
	if as, err := m.Analysis.AnalysisSchedule(); err != nil || as != nil {
		t.Errorf("digest.schedule reached the analyzer: %v, %v", as, err)
	}
	if ps, err := m.Planning.PlanSchedule(); err != nil || ps != nil {
		t.Errorf("digest.schedule reached the planner: %v, %v", ps, err)
	}

	for _, blank := range []string{"", "  "} {
		if s, err := (Digest{Schedule: blank}).DigestSchedule(); err != nil || s != nil {
			t.Errorf("schedule %q: %v, %v", blank, s, err)
		}
	}
	if s, err := (Digest{Schedule: "daily"}).DigestSchedule(); err != nil || s == nil || s.String() != "daily" {
		t.Fatalf("daily: %v, %v", s, err)
	}
	for _, c := range []struct {
		d    Digest
		want string
	}{
		{Digest{Schedule: "61 * * * *"}, `digest.schedule "61 * * * *": the minute "61" is outside 0-59`},
		{Digest{Schedule: "weekly"}, `digest.schedule "weekly" has 1 fields, not five`},
		{Digest{Schedule: "0 0 31 2 *"}, `digest.schedule "0 0 31 2 *" never comes round`},
	} {
		got := c.d.Errors()
		if len(got) != 1 || !strings.Contains(got[0], c.want) {
			t.Errorf("%+v: got %q, want one error saying %q", c.d, got, c.want)
		}
		if s, err := c.d.DigestSchedule(); err == nil || s != nil {
			t.Errorf("%+v: a bad schedule is accepted: %v", c.d, s)
		}
	}
}

// ADR-0199: digest.compact_after is a Go duration, 30 days when unset, 0
// turning compaction off; one that does not parse, or is negative, is an
// error naming the key and what to write.
func TestDigestCompactAfter(t *testing.T) {
	if DefaultCompactAfter != 30*24*time.Hour {
		t.Errorf("the default is %v, not 30 days", DefaultCompactAfter)
	}
	m := loaded(t, []byte("version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\ndigest:\n  compact_after: 336h\n"))
	if d, err := m.Digest.CompactAfterDuration(); err != nil || d != 14*24*time.Hour {
		t.Errorf("read: %v, %v", d, err)
	}
	for in, want := range map[string]time.Duration{"": DefaultCompactAfter, " ": DefaultCompactAfter, "0": 0, "0s": 0, "36h30m": 36*time.Hour + 30*time.Minute} {
		if d, err := (Digest{CompactAfter: in}).CompactAfterDuration(); err != nil || d != want {
			t.Errorf("%q: %v, %v, want %v", in, d, err, want)
		}
		if errs := (Digest{CompactAfter: in}).Errors(); len(errs) != 0 {
			t.Errorf("%q: %v", in, errs)
		}
	}
	for bad, want := range map[string]string{
		"a month": `digest.compact_after "a month" is not a duration; write one such as 720h or 336h, or 0 to turn compaction off`,
		"30d":     `digest.compact_after "30d" is not a duration`,
		"-1h":     `digest.compact_after "-1h" is negative; write a duration such as 720h or 336h, or 0 to turn compaction off`,
	} {
		if _, err := (Digest{CompactAfter: bad}).CompactAfterDuration(); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: %v, want %q", bad, err, want)
		}
		if errs := (Digest{CompactAfter: bad}).Errors(); len(errs) != 1 || !strings.Contains(errs[0], want) {
			t.Errorf("%q: errors %q, want one saying %q", bad, errs, want)
		}
	}
	both := (Digest{Schedule: "weekly", CompactAfter: "-1h", Agent: &Agent{Harness: "Claude Code"}}).Errors()
	if len(both) != 3 || !strings.HasPrefix(both[0], "digest.schedule ") || !strings.HasPrefix(both[1], "digest.compact_after ") || !strings.HasPrefix(both[2], "digest.agent ") {
		t.Errorf("the schedule's, the window's, and the agent's errors, in that order: %q", both)
	}
}

// ADR-0205: threads.remind_after is a Go duration, unset by default, which
// means no reminders; threads.remind_every is a Go duration, a day when
// unset; one that does not parse, a negative one, or a remind_every of zero
// is an error naming the key and what to write.
func TestThreadsReminders(t *testing.T) {
	if DefaultRemindEvery != 24*time.Hour {
		t.Errorf("the default is %v, not a day", DefaultRemindEvery)
	}
	base := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
	m := loaded(t, []byte(base))
	if m.Threads.RemindAfter != "" {
		t.Errorf("remind_after is set with no block: %q", m.Threads.RemindAfter)
	}
	if d, err := m.Threads.RemindAfterDuration(); err != nil || d != 0 {
		t.Errorf("no block, remind_after: %v, %v", d, err)
	}
	if d, err := m.Threads.RemindEveryDuration(); err != nil || d != 24*time.Hour {
		t.Errorf("no block, remind_every: %v, %v", d, err)
	}
	if errs := m.Threads.Errors(); len(errs) != 0 {
		t.Errorf("no block: %v", errs)
	}
	m = loaded(t, []byte(base+"threads:\n  remind_after: 4h\n  remind_every: 12h\n"))
	if d, err := m.Threads.RemindAfterDuration(); err != nil || d != 4*time.Hour {
		t.Errorf("remind_after read: %v, %v", d, err)
	}
	if d, err := m.Threads.RemindEveryDuration(); err != nil || d != 12*time.Hour {
		t.Errorf("remind_every read: %v, %v", d, err)
	}
	if errs := m.Threads.Errors(); len(errs) != 0 {
		t.Errorf("valid threads: %v", errs)
	}
	for in, want := range map[string]time.Duration{"": 0, " ": 0, "0": 0, "90m": 90 * time.Minute} {
		if d, err := (Threads{RemindAfter: in}).RemindAfterDuration(); err != nil || d != want {
			t.Errorf("remind_after %q: %v, %v, want %v", in, d, err, want)
		}
	}
	for in, want := range map[string]time.Duration{"": DefaultRemindEvery, " ": DefaultRemindEvery, "1h30m": 90 * time.Minute} {
		if d, err := (Threads{RemindEvery: in}).RemindEveryDuration(); err != nil || d != want {
			t.Errorf("remind_every %q: %v, %v, want %v", in, d, err, want)
		}
	}
	for _, c := range []struct {
		th   Threads
		want string
	}{
		{Threads{RemindAfter: "soon"}, `threads.remind_after "soon" is not a duration; write one such as 4h or 90m, or unset it for no reminders`},
		{Threads{RemindAfter: "-1h"}, `threads.remind_after "-1h" is negative; write a duration such as 4h or 90m, or unset it for no reminders`},
		{Threads{RemindEvery: "daily"}, `threads.remind_every "daily" is not a duration; write one such as 24h or 12h, or unset it for 24h`},
		{Threads{RemindEvery: "-12h"}, `threads.remind_every "-12h" is not more than zero`},
		{Threads{RemindEvery: "0"}, `threads.remind_every "0" is not more than zero; write a duration such as 24h or 12h, or unset it for 24h`},
	} {
		if errs := c.th.Errors(); len(errs) != 1 || !strings.HasPrefix(errs[0], c.want) {
			t.Errorf("%+v: errors %q, want one saying %q", c.th, errs, c.want)
		}
	}
	both := (Threads{RemindAfter: "soon", RemindEvery: "0"}).Errors()
	if len(both) != 2 || !strings.HasPrefix(both[0], "threads.remind_after ") || !strings.HasPrefix(both[1], "threads.remind_every ") {
		t.Errorf("remind_after's and remind_every's errors, in that order: %q", both)
	}
}

// ADR-0196: flai manifest set writes digest.agent, which reads back as the
// curator's agent, and refuses one that is not valid with the field and
// the reason.
func TestDigestAgentSetting(t *testing.T) {
	base := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
	got, err := EditSettings([]byte(base), []Assignment{{"digest.agent", `{"model":"claude-sonnet-5"}`}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := base + "digest:\n  agent:\n    model: claude-sonnet-5\n"; string(got) != want {
		t.Errorf("written:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if a := m.DigestAgent(); a == nil || a.Model != "claude-sonnet-5" {
		t.Errorf("read back %v", a)
	}
	var value SettingValue
	for _, v := range m.SettingValues() {
		if v.Key == "digest.agent" {
			value = v
		}
	}
	if a, ok := value.Value.(*Agent); !value.Set || !ok || a.Model != "claude-sonnet-5" {
		t.Errorf("digest.agent = %+v", value)
	}
	_, err = EditSettings([]byte(base), []Assignment{{"digest.agent", `{"harness":"Bad Harness"}`}}, nil)
	var refused *RefusedError
	if !errors.As(err, &refused) {
		t.Fatalf("want a refusal, got %v", err)
	}
	if i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == "digest.agent" }); i < 0 || !strings.Contains(refused.Problems[i].Reason, "is not a name such as claude-code") {
		t.Errorf("refusal %+v", refused.Problems)
	}
}

// S-0349: agent.provider reads in the project's agent and each strategic
// agent's, a strategic agent's winning over the project's in the merge;
// providers reads by name with api, base_url, and models, and survives JSON;
// and a provider that is not valid, or that sets key_env, is refused naming
// the key, the latter naming the host setting that holds it.
func TestProviders(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"agent:\n  harness: claude-code\n  model: claude-opus-5-5\n  provider: litellm\n" +
		"planning:\n  agent:\n    provider: openrouter\n" +
		"orchestration:\n  agent:\n    model: claude-sonnet-5\n" +
		"analysis:\n  agent:\n    provider: anthropic\n" +
		"providers:\n  litellm:\n    base_url: https://litellm.example.com/v1\n    models:\n      opus: anthropic/claude-opus-5-5\n  openrouter:\n    api: openai-chat\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	for what, errs := range map[string][]string{"agent": m.Agent.Problems(), "planning": m.Planning.Errors(), "orchestration": m.Orchestration.Errors(), "analysis": m.Analysis.Errors(), "providers": m.ProviderErrors()} {
		if len(errs) != 0 {
			t.Errorf("valid %s: %v", what, errs)
		}
	}
	for what, c := range map[string]struct {
		a    *Agent
		want string
	}{
		"agent":         {m.Agent, "litellm"},
		"planning":      {m.PlanningAgent(), "openrouter"},
		"orchestration": {m.OrchestrationAgent(), "litellm"},
		"analysis":      {m.AnalysisAgent(), "anthropic"},
	} {
		if c.a.Provider != c.want || c.a.Harness != "claude-code" {
			t.Errorf("%s's agent: %+v, want provider %s", what, c.a, c.want)
		}
	}
	want := map[string]Provider{
		"litellm":    {BaseURL: "https://litellm.example.com/v1", Models: map[string]string{"opus": "anthropic/claude-opus-5-5"}},
		"openrouter": {API: APIOpenAIChat},
	}
	if !reflect.DeepEqual(m.Providers, want) {
		t.Errorf("providers: %+v", m.Providers)
	}
	data, err := json.Marshal(m.Providers)
	if err != nil {
		t.Fatal(err)
	}
	var back map[string]Provider
	if err := json.Unmarshal(data, &back); err != nil || !reflect.DeepEqual(back, want) {
		t.Errorf("through JSON %s: %+v, %v", data, back, err)
	}
	for _, c := range []struct {
		p    map[string]Provider
		want []string
	}{
		{map[string]Provider{"litellm": {KeyEnv: "LITELLM_API_KEY"}}, []string{"providers.litellm.key_env is the host's", "flai serve agent provider litellm --key-env"}},
		{map[string]Provider{"litellm": {KeyEnv: "not a name"}}, []string{"providers.litellm.key_env is the host's"}},
		{map[string]Provider{"litellm": {API: "grpc"}}, []string{`providers.litellm.api "grpc" is not an API`, "anthropic-messages, openai-chat, openai-responses"}},
		{map[string]Provider{"litellm": {BaseURL: "litellm.example.com"}}, []string{`providers.litellm.base_url "litellm.example.com" is not an absolute http or https URL`}},
		{map[string]Provider{"litellm": {Models: map[string]string{"opus": ""}}}, []string{`providers.litellm.models.opus "" is not a model ID`}},
		{map[string]Provider{"LiteLLM": {API: APIOpenAIChat}}, []string{`providers "LiteLLM" is not a provider name`}},
	} {
		errs := (Manifest{Providers: c.p}).ProviderErrors()
		if len(errs) != 1 {
			t.Errorf("%+v: want one refusal, got %q", c.p, errs)
			continue
		}
		for _, w := range c.want {
			if !strings.Contains(errs[0], w) {
				t.Errorf("%+v: got %q, want %q", c.p, errs[0], w)
			}
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

// ADR-0201: flai.command is a program inside the repository, named by its
// path from the root; unset, it is empty, and an absolute path or one that
// leaves the repository is refused, naming the key.
func TestFlaiCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	base := "version: 1\nname: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	for _, c := range []struct{ yaml, want, err string }{
		{"", "", ""},
		{"flai:\n  command: scripts/flai.sh\n", "scripts/flai.sh", ""},
		{"flai:\n  command: /usr/local/bin/flai\n", "", `flai.command "/usr/local/bin/flai" is an absolute path`},
		{"flai:\n  command: ../flai.sh\n", "", `flai.command "../flai.sh" leaves the repository`},
		{"flai:\n  command: scripts//flai.sh\n", "", `flai.command "scripts//flai.sh" is not a clean path`},
	} {
		if err := os.WriteFile(path, []byte(base+c.yaml), 0o644); err != nil {
			t.Fatal(err)
		}
		m, err := Load(path)
		switch {
		case c.err == "" && (err != nil || m.Flai.Command != c.want):
			t.Errorf("%q: %v %+v", c.yaml, err, m.Flai)
		case c.err != "" && (err == nil || !strings.Contains(err.Error(), c.err)):
			t.Errorf("%q: want an error with %q, got %v", c.yaml, c.err, err)
		}
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

// S-0381: design.tokens names the design token file (ADR-0146). Unset, it is
// design-system.tokens.json in the design layout folder's system folder; set,
// it is the path given; a path that is absolute, leaves the repository, or is
// not clean is refused at load, naming the key and what to write.
func TestDesignTokens(t *testing.T) {
	head := func(design string) string {
		return "version: 1\nname: demo\nlayout:\n  design: " + design + "\n  docs: docs\n  wip: wip\n"
	}
	for _, c := range []struct{ name, body, want string }{
		{"unset", head("design"), DefaultDesignTokens},
		{"unset, the design folder renamed", head("arch"), "arch/system/design-system.tokens.json"},
		{"an empty block", head("design") + "design: {}\n", DefaultDesignTokens},
		{"set", head("design") + "design:\n  tokens: ui/tokens.json\n", "ui/tokens.json"},
		{"set with spaces", head("design") + "design:\n  tokens: \" ui/tokens.json \"\n", "ui/tokens.json"},
	} {
		if got := loaded(t, []byte(c.body)).DesignTokens(); got != c.want {
			t.Errorf("%s: DesignTokens = %q, want %q", c.name, got, c.want)
		}
	}
	if DefaultDesignTokens != "design/system/design-system.tokens.json" {
		t.Errorf("the default is %q", DefaultDesignTokens)
	}
	if got := (Manifest{}).DesignTokens(); got != DefaultDesignTokens {
		t.Errorf("a manifest with no layout: DesignTokens = %q, want %q", got, DefaultDesignTokens)
	}
	for _, c := range []struct{ design, want string }{
		{"design", "design/system/design-system.md"},
		{"arch", "arch/system/design-system.md"},
	} {
		if got := loaded(t, []byte(head(c.design))).DesignGuidance(); got != c.want {
			t.Errorf("design folder %s: DesignGuidance = %q, want %q", c.design, got, c.want)
		}
	}
	if got := (Manifest{}).DesignGuidance(); got != "design/system/design-system.md" {
		t.Errorf("a manifest with no layout: DesignGuidance = %q", got)
	}

	p := filepath.Join(t.TempDir(), File)
	for _, c := range []struct{ tokens, want string }{
		{"/etc/tokens.json", "is an absolute path; write it relative to the repository root, such as " + DefaultDesignTokens},
		{`C:\tokens.json`, "is an absolute path"},
		{`\\host\tokens.json`, "is an absolute path"},
		{"../tokens.json", "leaves the repository; write a file inside it, relative to its root, such as " + DefaultDesignTokens},
		{"..", "leaves the repository"},
		{"design/../../tokens.json", "leaves the repository"},
		{"design/./tokens.json", "is not a clean path"},
	} {
		body := head("design") + "design:\n  tokens: '" + c.tokens + "'\n"
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := Load(p)
		if err == nil || !strings.Contains(err.Error(), `design.tokens "`) || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: %v", c.tokens, err)
		}
	}

	data, err := yaml.Marshal(Manifest{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "design:") || strings.Contains(string(data), "tokens") {
		t.Errorf("an unset design block is written:\n%s", data)
	}
	set := loaded(t, []byte(head("design")+"design:\n  tokens: ui/tokens.json\n"))
	data, err = yaml.Marshal(Manifest{Design: set.Design})
	if err != nil {
		t.Fatal(err)
	}
	var again Manifest
	if err := yaml.Unmarshal(data, &again); err != nil || again.Design != set.Design {
		t.Errorf("YAML round trip: %v\n%s\n%+v", err, data, again.Design)
	}
	js, err := json.Marshal(Manifest{Design: set.Design})
	if err != nil || !strings.Contains(string(js), `"design":{"tokens":"ui/tokens.json"}`) {
		t.Errorf("JSON: %v %s", err, js)
	}
}

// S-0385: design.stylesheet names the generated design stylesheet (ADR-0159).
// Unset, there is none; set, it is the path given, trimmed; a path that is
// absolute, leaves the repository, or is not clean is refused at load, naming
// the key and what to write. One outside every sub-project loads, for flai
// check to report, and the sub-project that holds it is the one whose path is
// a folder above it.
func TestDesignStylesheet(t *testing.T) {
	head := "version: 1\nname: demo\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" +
		"projects:\n" +
		"  - name: web\n    path: web\n    kind: sveltekit\n" +
		"  - name: cli\n    path: tools/cli/\n    kind: go\n"
	for _, c := range []struct{ name, body, want string }{
		{"unset", head, ""},
		{"an empty block", head + "design: {}\n", ""},
		{"set", head + "design:\n  stylesheet: web/src/theme.css\n", "web/src/theme.css"},
		{"set with spaces", head + "design:\n  stylesheet: \" web/src/theme.css \"\n", "web/src/theme.css"},
	} {
		if got := loaded(t, []byte(c.body)).DesignStylesheet(); got != c.want {
			t.Errorf("%s: DesignStylesheet = %q, want %q", c.name, got, c.want)
		}
	}
	if got := (Manifest{}).DesignStylesheet(); got != "" {
		t.Errorf("an empty manifest: DesignStylesheet = %q", got)
	}

	for _, c := range []struct{ stylesheet, want string }{
		{"", ""},
		{"web/src/theme.css", "web"},
		{" web/src/theme.css ", "web"},
		{"tools/cli/theme.css", "cli"},
		{"web", ""},
		{"website/src/theme.css", ""},
		{"theme.css", ""},
		{"design/system/theme.css", ""},
	} {
		body := head
		if c.stylesheet != "" {
			body += "design:\n  stylesheet: '" + c.stylesheet + "'\n"
		}
		p, ok := loaded(t, []byte(body)).DesignStylesheetProject()
		if ok != (c.want != "") || p.Name != c.want {
			t.Errorf("%q: DesignStylesheetProject = %q %v, want %q", c.stylesheet, p.Name, ok, c.want)
		}
	}

	file := filepath.Join(t.TempDir(), File)
	for _, c := range []struct{ stylesheet, want string }{
		{"/web/theme.css", "is an absolute path; write it relative to the repository root, such as web/src/theme.css"},
		{`C:\theme.css`, "is an absolute path"},
		{"../theme.css", "leaves the repository; write a file inside it, relative to its root, such as web/src/theme.css"},
		{"web/../../theme.css", "leaves the repository"},
		{"web/./src/theme.css", "is not a clean path; write it with single slashes, no . or .. segments, and none at the end, such as web/src/theme.css"},
		{"web//src/theme.css", "is not a clean path"},
	} {
		body := head + "design:\n  stylesheet: '" + c.stylesheet + "'\n"
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := Load(file)
		if err == nil || !strings.Contains(err.Error(), `design.stylesheet "`) || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: %v", c.stylesheet, err)
		}
	}

	set := loaded(t, []byte(head+"design:\n  stylesheet: web/src/theme.css\n"))
	js, err := json.Marshal(Manifest{Design: set.Design})
	if err != nil || !strings.Contains(string(js), `"design":{"stylesheet":"web/src/theme.css"}`) {
		t.Errorf("JSON: %v %s", err, js)
	}
	data, err := yaml.Marshal(Manifest{Design: Design{Tokens: "ui/tokens.json"}})
	if err != nil || strings.Contains(string(data), "stylesheet") {
		t.Errorf("an unset stylesheet is written: %v\n%s", err, data)
	}
}

// S-0273: tests is an ordered list of tiers that loads with every field,
// queue among them (S-0405), reads back the same through YAML and JSON, and
// is valid.
func TestTestsRoundTrip(t *testing.T) {
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"tests:\n" +
		"- name: unit\n" +
		"  command: [go, test, -json, \"{packages}\"]\n" +
		"  dir: flai\n" +
		"  paths: [\"flai/**\", \"!flai/testdata/**\"]\n" +
		"  format: go-test-json\n" +
		"  all_command: [go, test, -json, ./...]\n" +
		"  fix_command: [gofmt, -w, \"{files}\"]\n" +
		"- name: smoke\n" +
		"  command: [scripts/smoke.sh]\n" +
		"  all_only: true\n" +
		"  queue: true\n"
	m := loaded(t, []byte(body))
	want := []TestTier{
		{Name: "unit", Command: []string{"go", "test", "-json", PlaceholderPackages}, Dir: "flai",
			Paths: []string{"flai/**", "!flai/testdata/**"}, Format: FormatGoTestJSON, AllCommand: []string{"go", "test", "-json", "./..."},
			FixCommand: []string{"gofmt", "-w", PlaceholderFiles}},
		{Name: "smoke", Command: []string{"scripts/smoke.sh"}, AllOnly: true, Queue: true},
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
	if !strings.Contains(string(js), `"all_only":true`) || !strings.Contains(string(js), `"all_command":`) || !strings.Contains(string(js), `"fix_command":`) || !strings.Contains(string(js), `"queue":true`) {
		t.Errorf("JSON keys are not snake_case: %s", js)
	}
}

// S-0367: flai manifest set writes a tier's fix_command on a line of its own
// and reads it back; tiers that differ only in it are not the same.
func TestEditSettingsWritesFixCommand(t *testing.T) {
	value := `[{"name":"md","command":["scripts/lint-md.sh","{files}"],"paths":["**/*.md"],"fix_command":["scripts/lint-md.sh","--fix","{files}"]}]`
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"tests", value}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "tests:\n" +
		"  - name: md\n" +
		"    command: [scripts/lint-md.sh, \"{files}\"]\n" +
		"    paths: [\"**/*.md\"]\n" +
		"    fix_command: [scripts/lint-md.sh, \"--fix\", \"{files}\"]\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if len(m.Tests) != 1 || !reflect.DeepEqual(m.Tests[0].FixCommand, []string{"scripts/lint-md.sh", "--fix", PlaceholderFiles}) {
		t.Errorf("reads back: %+v", m.Tests)
	}
	unfixed := []TestTier{m.Tests[0]}
	unfixed[0].FixCommand = nil
	if sameTiers(m.Tests, unfixed) {
		t.Error("tiers that differ in fix_command are the same")
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
	named := func(name string, covers ...string) TestTier {
		return with(func(x *TestTier) { x.Name, x.Covers = name, covers })
	}
	cycle := []TestTier{named("a", "c"), named("b", "a"), named("c", "b"), named("d", "a")}
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
		// S-0364: {importers} is refused where {packages} is, and every
		// refusal names all three placeholders.
		{"another placeholder names all three", []TestTier{with(func(x *TestTier) { x.Command = []string{"go", "{story}"} })}, "argument {story} is not a placeholder flai fills; write {packages}, {files}, or {importers}, or remove it"},
		{"importers inside", []TestTier{with(func(x *TestTier) { x.Command = []string{"go", "test", "./{importers}"} })}, `argument "./{importers}" holds a placeholder inside it, which flai does not fill; write {packages}, {files}, or {importers} as an argument of its own`},
		{"importers under all", []TestTier{with(func(x *TestTier) { x.AllCommand = []string{"go", "test", "{importers}"} })}, `tests[0].all_command (tier "unit") argument {importers} is a placeholder, and all_command runs for no paths; write it without {packages}, {files}, or {importers}`},
		// S-0367: fix_command is checked as command is, placeholders allowed.
		{"a bad placeholder in fix", []TestTier{with(func(x *TestTier) { x.FixCommand = []string{"lint", "{story}"} })}, `tests[0].fix_command (tier "unit") argument {story} is not a placeholder flai fills`},
		{"no program in fix", []TestTier{with(func(x *TestTier) { x.FixCommand = []string{"{files}"} })}, `tests[0].fix_command (tier "unit") begins with the placeholder {files}`},
		{"an absolute dir", []TestTier{with(func(x *TestTier) { x.Dir = "/srv/flai" })}, `tests[0].dir (tier "unit") "/srv/flai" is an absolute path`},
		{"a dir outside", []TestTier{with(func(x *TestTier) { x.Dir = "../other" })}, `"../other" leaves the repository`},
		{"a dir not clean", []TestTier{with(func(x *TestTier) { x.Dir = "flai/./web/" })}, "is not a clean path"},
		{"no paths", []TestTier{with(func(x *TestTier) { x.Paths = nil })}, `tests[0].paths (tier "unit") is empty`},
		{"only exclusions", []TestTier{with(func(x *TestTier) { x.Paths = []string{"!flai/testdata"} })}, "has no pattern that selects a path"},
		{"a bad pattern", []TestTier{with(func(x *TestTier) { x.Paths = []string{"flai/**", "!../x"} })}, `tests[0].paths (tier "unit") pattern "!../x" has a .. segment`},
		{"a bad format", []TestTier{with(func(x *TestTier) { x.Format = "junit" })}, `tests[0].format (tier "unit") "junit" is not a format flai reads; write go-test-json, vitest-json, golangci-json, gofmt-list, plain`},
		{"covers no tier", []TestTier{named("unit", "lint")}, `tests[0].covers (tier "unit") "lint" names no tier`},
		{"covers itself", []TestTier{named("unit", "unit")}, `tests[0].covers (tier "unit") names the tier itself`},
		{"covers a tier twice", []TestTier{named("unit", "lint", "lint"), named("lint")}, `tests[0].covers (tier "unit") names "lint" twice`},
		{"covers in a cycle", cycle, `tests[0].covers (tier "a") makes a cycle through the tiers "a", "b", "c", which cover one another`},
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
		with(func(x *TestTier) { x.Command = []string{"go", "test", "{importers}"} }),
		with(func(x *TestTier) { x.Command = []string{"sh", "-c", "echo ${HOME}"} }),
		with(func(x *TestTier) { x.FixCommand = []string{"scripts/lint-md.sh", "--fix", "{files}"} }),
		// S-0405: queue is accepted on any tier, with any other field.
		with(func(x *TestTier) { x.Queue = true }),
		with(func(x *TestTier) {
			x.Queue, x.Paths, x.AllOnly, x.FixCommand = true, nil, true, []string{"scripts/lint-md.sh", "--fix", "{files}"}
		}),
	} {
		if errs := (Manifest{Tests: []TestTier{tier}}).TestErrors(); len(errs) != 0 {
			t.Errorf("%+v refused: %v", tier, errs)
		}
	}
	if errs := (Manifest{Tests: cycle}).TestErrors(); len(errs) != 1 {
		t.Errorf("a cycle is reported once: %q", errs)
	}
	chain := []TestTier{named("a", "b"), named("b", "c"), named("c")}
	if errs := (Manifest{Tests: chain}).TestErrors(); len(errs) != 0 {
		t.Errorf("a chain refused: %v", errs)
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

// componentsManifest is a manifest whose projects are the YAML given, a list
// of sub-projects written at an indent of two.
func componentsManifest(t *testing.T, projects string) Manifest {
	t.Helper()
	return loaded(t, []byte("version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\nprojects:\n"+projects))
}

// S-0437: a sveltekit or svelte sub-project's components settings are its
// key over the defaults; any other kind has none.
func TestComponentSettings(t *testing.T) {
	m := componentsManifest(t, `  - {name: app, path: app, kind: sveltekit}
  - {name: ui, path: ui, kind: svelte}
  - {name: api, path: api, kind: go}
  - name: site
    path: site
    kind: sveltekit
    components:
      include: ["src/**/*.svelte"]
      exclude: ["src/routes/**"]
      setup: " src/app.css "
  - name: kit
    path: kit
    kind: svelte
    components:
      setup: src/theme.css
`)
	defaults := Components{Include: []string{DefaultComponentInclude}}
	for _, c := range []struct {
		i    int
		want Components
		ok   bool
	}{
		{0, defaults, true},
		{1, defaults, true},
		{2, Components{}, false},
		{3, Components{Include: []string{"src/**/*.svelte"}, Exclude: []string{"src/routes/**"}, Setup: "src/app.css"}, true},
		{4, Components{Include: []string{DefaultComponentInclude}, Setup: "src/theme.css"}, true},
	} {
		p := m.Projects[c.i]
		got, ok := p.ComponentSettings()
		if ok != c.ok || !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %+v, %v; want %+v, %v", p.Name, got, ok, c.want, c.ok)
		}
	}
	if !SupportsComponents("sveltekit") || !SupportsComponents("svelte") || SupportsComponents("go") || SupportsComponents("template") {
		t.Error("SupportsComponents takes sveltekit and svelte alone")
	}
	got, _ := m.Projects[3].ComponentSettings()
	got.Include[0] = "changed"
	if m.Projects[3].Components.Include[0] != "src/**/*.svelte" {
		t.Error("ComponentSettings shares its include with the manifest")
	}
}

// S-0437: ComponentErrors refuses a glob that is not valid, a setup that is
// not a clean path inside the sub-project or names no file under it, and the
// key on a kind whose components flai does not read; it checks no file with
// no file system.
func TestComponentErrors(t *testing.T) {
	fsys := fstest.MapFS{
		"web/src/app.css":   {Data: []byte("@import 'tailwindcss';\n")},
		"web/src/lib/x.css": {Data: nil},
	}
	for _, c := range []struct {
		name, project, want string
	}{
		{"valid", "    components:\n      include: [\"src/lib/**/*.svelte\", src/Button.svelte]\n      exclude: [\"src/lib/internal/**\"]\n      setup: src/app.css\n", ""},
		{"the key with defaults alone", "    components: {}\n", ""},
		{"an unclosed bracket", "    components:\n      include: [\"src/[ab.svelte\"]\n",
			`projects[0].components.include (sub-project "web") pattern "src/[ab.svelte" is not a valid glob, such as an unclosed [ or a trailing \; fix it or remove it`},
		{"** inside a segment", "    components:\n      exclude: [\"src/**x\"]\n",
			`projects[0].components.exclude (sub-project "web") pattern "src/**x" has ** inside a segment`},
		{"an empty glob", "    components:\n      include: [\"\"]\n",
			`projects[0].components.include (sub-project "web") pattern "" is empty; write a glob relative to the sub-project's folder, such as src/lib/**/*.svelte, or remove it`},
		{"an absolute glob", "    components:\n      include: [\"/src/*.svelte\"]\n",
			`projects[0].components.include (sub-project "web") pattern "/src/*.svelte" is an absolute path; write it relative to the sub-project's folder`},
		{"a glob that leaves", "    components:\n      exclude: [\"../api/**\"]\n",
			`projects[0].components.exclude (sub-project "web") pattern "../api/**" has a .. segment, which leaves the sub-project`},
		{"a setup that names no file", "    components:\n      setup: src/missing.css\n",
			`projects[0].components.setup (sub-project "web") "src/missing.css" names no file under web; write the path of a module of the sub-project, relative to its folder, such as src/app.css, or remove it`},
		{"a setup that is a folder", "    components:\n      setup: src/lib\n",
			`projects[0].components.setup (sub-project "web") "src/lib" names no file under web`},
		{"a setup that leaves", "    components:\n      setup: ../api/main.go\n",
			`projects[0].components.setup (sub-project "web") "../api/main.go" leaves the sub-project; write a module inside it, relative to its root, such as src/app.css`},
		{"an absolute setup", "    components:\n      setup: /web/src/app.css\n",
			`projects[0].components.setup (sub-project "web") "/web/src/app.css" is an absolute path; write it relative to the sub-project's folder, such as src/app.css`},
		{"a setup not clean", "    components:\n      setup: src//app.css\n",
			`projects[0].components.setup (sub-project "web") "src//app.css" is not a clean path`},
	} {
		t.Run(c.name, func(t *testing.T) {
			m := componentsManifest(t, "  - name: web\n    path: web\n    kind: sveltekit\n"+c.project)
			errs := m.ComponentErrors(fsys)
			switch {
			case c.want == "" && len(errs) > 0:
				t.Errorf("want none, got %q", errs)
			case c.want != "" && (len(errs) != 1 || !strings.HasPrefix(errs[0], c.want)):
				t.Errorf("want one beginning %q, got %q", c.want, errs)
			}
		})
	}
	m := componentsManifest(t, "  - name: web\n    path: web\n    kind: sveltekit\n  - name: api\n    path: api\n    kind: go\n    components:\n      include: [\"src/[\"]\n")
	want := []string{`projects[1].components (sub-project "api") is set on a sub-project of kind "go", whose components flai does not read; remove it, or set kind to sveltekit or svelte if the sub-project is one`}
	if errs := m.ComponentErrors(fsys); !reflect.DeepEqual(errs, want) {
		t.Errorf("the key on an unsupported kind:\n got %q\nwant %q", errs, want)
	}
	m = componentsManifest(t, "  - name: web\n    path: web\n    kind: sveltekit\n    components:\n      setup: src/missing.css\n")
	if errs := m.ComponentErrors(nil); len(errs) != 0 {
		t.Errorf("with no file system the setup's file is not checked: %q", errs)
	}
	if p := splitProblem(m.ComponentErrors(fsys)[0]); p.Field != "projects[0].components.setup" {
		t.Errorf("the problem's field: %+v", p)
	}
}

// S-0483: a sub-project's install key, on a sub-project of any kind, loads
// from the manifest as its command and paths; none leaves it nil.
func TestInstallLoads(t *testing.T) {
	m := componentsManifest(t, `  - name: web
    path: web
    kind: sveltekit
    install:
      command: ["../scripts/web-install.sh", "--if-needed"]
      paths: ["package.json", "pnpm-lock.yaml"]
  - name: api
    path: api
    kind: go
    install:
      command: [go, mod, download]
      paths: [go.sum]
  - {name: site, path: site, kind: svelte}
`)
	for i, want := range []*Install{
		{Command: []string{"../scripts/web-install.sh", "--if-needed"}, Paths: []string{"package.json", "pnpm-lock.yaml"}},
		{Command: []string{"go", "mod", "download"}, Paths: []string{"go.sum"}},
		nil,
	} {
		if got := m.Projects[i].Install; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: install %+v, want %+v", m.Projects[i].Name, got, want)
		}
	}
	if errs := m.InstallErrors(); len(errs) != 0 {
		t.Errorf("valid install keys: %q", errs)
	}
}

// S-0483: InstallErrors refuses an install key with no command, no program
// first, no paths, or a glob components would refuse, each by the
// sub-project's index, field, and name; a sub-project with no install key
// has none.
func TestInstallErrors(t *testing.T) {
	for _, c := range []struct {
		name, project string
		want          []string
	}{
		{"valid", "    install:\n      command: [pnpm, install]\n      paths: [package.json, \"src/**/*.lock\"]\n", nil},
		{"unset", "", nil},
		{"no command", "    install:\n      paths: [package.json]\n",
			[]string{`projects[0].install.command (sub-project "web") is empty; write the program and its arguments as a list`}},
		{"an empty command", "    install:\n      command: []\n      paths: [package.json]\n",
			[]string{`projects[0].install.command (sub-project "web") is empty`}},
		{"an empty first element", "    install:\n      command: [\" \", install]\n      paths: [package.json]\n",
			[]string{`projects[0].install.command (sub-project "web") has no program first`}},
		{"no paths", "    install:\n      command: [pnpm, install]\n",
			[]string{`projects[0].install.paths (sub-project "web") is empty; write the globs of the files whose change needs the install, relative to the sub-project's folder, such as package.json`}},
		{"an empty key", "    install: {}\n",
			[]string{`projects[0].install.command (sub-project "web") is empty`, `projects[0].install.paths (sub-project "web") is empty`}},
		{"an unclosed bracket", "    install:\n      command: [pnpm, install]\n      paths: [\"pnpm-[.yaml\"]\n",
			[]string{`projects[0].install.paths (sub-project "web") pattern "pnpm-[.yaml" is not a valid glob`}},
		{"a glob that leaves", "    install:\n      command: [pnpm, install]\n      paths: [\"../pnpm-lock.yaml\"]\n",
			[]string{`projects[0].install.paths (sub-project "web") pattern "../pnpm-lock.yaml" has a .. segment, which leaves the sub-project; write a glob inside its folder, relative to it, such as package.json`}},
		{"an empty glob", "    install:\n      command: [pnpm, install]\n      paths: [\"\"]\n",
			[]string{`projects[0].install.paths (sub-project "web") pattern "" is empty; write a glob relative to the sub-project's folder, such as package.json, or remove it`}},
	} {
		t.Run(c.name, func(t *testing.T) {
			m := componentsManifest(t, "  - name: web\n    path: web\n    kind: go\n"+c.project)
			errs := m.InstallErrors()
			if len(errs) != len(c.want) {
				t.Fatalf("want %d beginning %q, got %q", len(c.want), c.want, errs)
			}
			for i, w := range c.want {
				if !strings.HasPrefix(errs[i], w) {
					t.Errorf("error %d: got %q, want one beginning %q", i, errs[i], w)
				}
			}
		})
	}
	m := componentsManifest(t, "  - name: web\n    path: web\n    kind: go\n    install:\n      paths: [package.json]\n")
	if p := splitProblem(m.InstallErrors()[0]); p.Field != "projects[0].install.command" {
		t.Errorf("the problem's field: %+v", p)
	}
}

// S-0483: InstallMatches says whether a changed path under the sub-project's
// folder matches one of its install globs, read from that folder; a path of
// another sub-project, another file, and a sub-project with no install key
// do not.
func TestInstallMatches(t *testing.T) {
	m := componentsManifest(t, `  - name: flaiover
    path: flaiover
    kind: sveltekit
    install:
      command: ["../scripts/flaiover-install.sh", "--if-needed"]
      paths: ["package.json", "pnpm-lock.yaml", "packages/**/package.json", "pnpm-[.yaml"]
  - {name: flai, path: flai, kind: go}
`)
	flaiover, flai := m.Projects[0], m.Projects[1]
	for _, c := range []struct {
		name    string
		changed []string
		want    bool
	}{
		{"its lock file", []string{"flai/main.go", "flaiover/pnpm-lock.yaml"}, true},
		{"its package.json", []string{"./flaiover/package.json"}, true},
		{"a ** glob, no segment between", []string{"flaiover/packages/package.json"}, true},
		{"a ** glob, segments between", []string{"flaiover/packages/ui/kit/package.json"}, true},
		{"another file of it", []string{"flaiover/src/app.css", "flaiover/package-lock.json"}, false},
		{"its file under another folder", []string{"flaiover/src/package.json"}, false},
		{"another sub-project's lock file", []string{"flai/pnpm-lock.yaml", "pnpm-lock.yaml"}, false},
		{"a folder whose name begins with its own", []string{"flaiover2/pnpm-lock.yaml"}, false},
		{"nothing changed", nil, false},
	} {
		if got := flaiover.InstallMatches(c.changed); got != c.want {
			t.Errorf("%s %q: got %v, want %v", c.name, c.changed, got, c.want)
		}
	}
	if flai.InstallMatches([]string{"flai/go.sum", "flai/pnpm-lock.yaml"}) {
		t.Error("a sub-project with no install key matches")
	}
}
