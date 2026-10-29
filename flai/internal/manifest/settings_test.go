package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

// settingsFixture is a manifest with comments in and around the strategic
// blocks, and keys flai manifest set does not write.
const settingsFixture = `version: 1
name: demo   # the project's name
key: d
layout:
  design: design
  docs: docs
  wip: wip
# the units of planning data
planning:                     # optional: units
  currency: EUR               # do not change
  cycle: 168h                 # a week
  agent:                      # the planner's agent
    model: claude-sonnet-5
    config:
      effort: medium
orchestration:
  policy: fifo    # the order
  release:
    policy: threshold
    count: 5      # five stories
  permissions:
    promote_to_ready: true
claims:
  shared:
    - design/adrs   # shared
`

// replaced is the fixture with old replaced by new once, failing when old is
// not in it.
func replaced(t *testing.T, old, new string) string {
	t.Helper()
	if !strings.Contains(settingsFixture, old) {
		t.Fatalf("the fixture has no %q", old)
	}
	return strings.Replace(settingsFixture, old, new, 1)
}

// S-0229: a key of each kind is written into its block, its comment kept in
// its column, and every other byte of the file is as it was.
func TestEditSettingsWritesEachKindAndKeepsTheRest(t *testing.T) {
	for _, c := range []struct {
		kind Kind
		set  Assignment
		want string
	}{
		{KindBoolean, Assignment{"orchestration.permissions.publish", "true"},
			replaced(t, "    promote_to_ready: true\n", "    promote_to_ready: true\n    publish: true\n")},
		{KindChoice, Assignment{"orchestration.policy", "wsjf"},
			replaced(t, "  policy: fifo    # the order\n", "  policy: wsjf    # the order\n")},
		{KindChoice, Assignment{"orchestration.permissions.answer_threads", "off"},
			replaced(t, "    promote_to_ready: true\n", "    promote_to_ready: true\n    answer_threads: \"off\"\n")},
		{KindNumber, Assignment{"planning.hour_rate", "95.50"},
			replaced(t, "      effort: medium\n", "      effort: medium\n  hour_rate: 95.5\n")},
		{KindNumber, Assignment{"orchestration.release.count", "12"},
			replaced(t, "    count: 5      # five stories\n", "    count: 12     # five stories\n")},
		{KindDuration, Assignment{"planning.cycle", "336h"},
			replaced(t, "  cycle: 168h                 # a week\n", "  cycle: 336h                 # a week\n")},
		{KindCron, Assignment{"analysis.schedule", "0 6 * * 1"},
			settingsFixture + "analysis:\n  schedule: \"0 6 * * 1\"\n"},
		{KindCron, Assignment{"planning.schedule", "daily"},
			replaced(t, "      effort: medium\n", "      effort: medium\n  schedule: daily\n")},
		{KindAgent, Assignment{"planning.agent", `{"harness":"claude-code","config":{"effort":"high"}}`},
			replaced(t, "    model: claude-sonnet-5\n    config:\n      effort: medium\n", "    harness: claude-code\n    config:\n      effort: high\n")},
		// S-0349: an agent names its provider, written after the model
		{KindAgent, Assignment{"planning.agent", `{"model":"claude-sonnet-5","provider":"litellm"}`},
			replaced(t, "    model: claude-sonnet-5\n    config:\n      effort: medium\n", "    model: claude-sonnet-5\n    provider: litellm\n")},
		{KindAgent, Assignment{"analysis.agent", `{"provider":"openrouter"}`},
			settingsFixture + "analysis:\n  agent:\n    provider: openrouter\n"},
		{KindText, Assignment{"orchestration.release.tag", "spring release"},
			replaced(t, "    count: 5      # five stories\n", "    count: 5      # five stories\n    tag: spring release\n")},
	} {
		t.Run(c.set.Key, func(t *testing.T) {
			if s, _ := SettingFor(c.set.Key); s.Kind != c.kind {
				t.Fatalf("%s is a %s, not a %s", c.set.Key, s.Kind, c.kind)
			}
			got, err := EditSettings([]byte(settingsFixture), []Assignment{c.set}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != c.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, c.want)
			}
		})
	}
}

// S-0229: several keys at once, into blocks that are missing, land in the
// order given.
func TestEditSettingsAddsMissingBlocks(t *testing.T) {
	base := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
	got, err := EditSettings([]byte(base), []Assignment{
		{"orchestration.release.whole_epics", "true"},
		{"orchestration.permissions.accept_reviews", "true"},
		{"analysis.agent", `{"model":"claude-sonnet-5"}`},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := base + "orchestration:\n  release:\n    whole_epics: true\n  permissions:\n    accept_reviews: true\nanalysis:\n  agent:\n    model: claude-sonnet-5\n"
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// S-0491: orchestration.permissions.edit_subprojects is a boolean off by
// default that flai manifest set writes, Load reads back, and Errors accepts.
func TestEditSettingsWritesEditSubprojects(t *testing.T) {
	if s, ok := SettingFor("orchestration.permissions.edit_subprojects"); !ok || s.Kind != KindBoolean || s.Default != "false" || s.Risk == "" {
		t.Fatalf("orchestration.permissions.edit_subprojects in the catalog: %+v, %v", s, ok)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"orchestration.permissions.edit_subprojects", "true"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "    edit_subprojects: true\n") {
		t.Fatalf("the key is not written:\n%s", got)
	}
	m := loaded(t, got)
	if !m.Orchestration.Permissions.Allows(PermitEditSubprojects) || len(m.Orchestration.Errors()) != 0 {
		t.Errorf("read back %+v, %v", m.Orchestration.Permissions, m.Orchestration.Errors())
	}
	if loaded(t, []byte(settingsFixture)).Orchestration.Permissions.Allows(PermitEditSubprojects) {
		t.Error("edit_subprojects is on unset")
	}
}

// S-0229: a bad value of each kind, a key outside the catalog, and a
// manifest left invalid as a whole are refused with the field and the
// reason, and the file is not written.
func TestEditSettingsRefusesWithTheFieldAndTheReason(t *testing.T) {
	for _, c := range []struct {
		name   string
		set    []Assignment
		unset  []string
		field  string
		reason string
	}{
		{"boolean", []Assignment{{"orchestration.permissions.publish", "yes"}}, nil, "orchestration.permissions.publish", `"yes" is not true or false`},
		{"choice", []Assignment{{"orchestration.policy", "fastest"}}, nil, "orchestration.policy", "is not an order policy"},
		{"choice of answers", []Assignment{{"orchestration.permissions.answer_threads", "always"}}, nil, "orchestration.permissions.answer_threads", "is not a way of answering threads"},
		{"number", []Assignment{{"orchestration.release.value", "-1"}}, nil, "orchestration.release.value", "-1 is not an amount of zero or more"},
		{"not a number", []Assignment{{"planning.hour_rate", "lots"}}, nil, "planning.hour_rate", `"lots" is not a number`},
		{"count", []Assignment{{"orchestration.release.count", "1.5"}}, nil, "orchestration.release.count", "is not a whole number"},
		{"negative count", []Assignment{{"orchestration.release.count", "-2"}}, nil, "orchestration.release.count", "-2 is not a number of zero or more"},
		{"duration", []Assignment{{"planning.cycle", "0s"}}, nil, "planning.cycle", "is not a duration longer than zero"},
		{"cron", []Assignment{{"analysis.schedule", "* * *"}}, nil, "analysis.schedule", "has 3 fields, not five"},
		{"agent", []Assignment{{"planning.agent", `{"harness":"Bad Harness"}`}}, nil, "planning.agent", "is not a name such as claude-code"},
		{"agent's provider", []Assignment{{"planning.agent", `{"provider":"Lite LLM"}`}}, nil, "planning.agent", `provider "Lite LLM" is not a name such as litellm`},
		{"not an agent", []Assignment{{"analysis.agent", "claude"}}, nil, "analysis.agent", "is not an agent"},
		{"agent with an unknown field", []Assignment{{"analysis.agent", `{"modelname":"x"}`}}, nil, "analysis.agent", "is not an agent"},
		{"text", []Assignment{{"orchestration.release.policy", "theme"}, {"orchestration.release.epic", "epic-7"}}, nil, "orchestration.release.epic", "is not an epic ID"},
		{"empty", []Assignment{{"planning.cycle", " "}}, nil, "planning.cycle", "is empty"},
		{"acceptance mode", []Assignment{{"acceptance.mode", "remote"}}, nil, "acceptance.mode", `"remote" is not an acceptance mode; write local`},
		{"forge", []Assignment{{"acceptance.forge", "gitlab"}}, nil, "acceptance.forge", `"gitlab" is not a forge flai reaches; write github`},
		{"merge method", []Assignment{{"acceptance.merge_method", "fast-forward"}}, nil, "acceptance.merge_method", `"fast-forward" is not a merge method; write merge`},
		{"digest schedule", []Assignment{{"digest.schedule", "weekly"}}, nil, "digest.schedule", `"weekly" has 1 fields, not five`},
		{"compaction window", []Assignment{{"digest.compact_after", "a month"}}, nil, "digest.compact_after", `"a month" is not a duration; write one such as 720h`},
		{"negative compaction window", []Assignment{{"digest.compact_after", "-24h"}}, nil, "digest.compact_after", `"-24h" is negative`},
		{"batch max zero", []Assignment{{"merge_queue.batch_max", "0"}}, nil, "merge_queue.batch_max", "0 is not a number of one or more; write the most enqueued stories one batch of the merge queue holds, such as 4, or remove it for 4"},
		{"batch max negative", []Assignment{{"merge_queue.batch_max", "-3"}}, nil, "merge_queue.batch_max", "-3 is not a number of one or more"},
		{"batch max not a number", []Assignment{{"merge_queue.batch_max", "four"}}, nil, "merge_queue.batch_max", `"four" is not a whole number; write one such as 5`},
		{"batch max not whole", []Assignment{{"merge_queue.batch_max", "2.5"}}, nil, "merge_queue.batch_max", `"2.5" is not a whole number`},
		{"unknown key", []Assignment{{"planning.colour", "blue"}}, nil, "planning.colour", "is not a setting flai manifest set writes"},
		{"unknown key unset", nil, []string{"dashboard.port"}, "dashboard.port", "is not a setting flai manifest set writes"},
		{"currency", []Assignment{{"planning.currency", "GBP"}}, nil, "planning.currency", "re-denominates every amount"},
		{"twice", []Assignment{{"planning.cycle", "1h"}}, []string{"planning.cycle"}, "planning.cycle", "more than once"},
		{"whole manifest", []Assignment{{"orchestration.release.policy", "theme"}}, nil, "orchestration.release", "is a theme with neither epic nor tag"},
	} {
		t.Run(c.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), File)
			if err := os.WriteFile(file, []byte(settingsFixture), 0o644); err != nil {
				t.Fatal(err)
			}
			err := WriteSettings(file, c.set, c.unset)
			var refused *RefusedError
			if !errors.As(err, &refused) {
				t.Fatalf("want a refusal, got %v", err)
			}
			i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == c.field })
			if i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
				t.Errorf("want %s: ...%s..., got %+v", c.field, c.reason, refused.Problems)
			}
			if data, _ := os.ReadFile(file); string(data) != settingsFixture {
				t.Errorf("the file changed:\n%s", data)
			}
		})
	}
}

// S-0392: each acceptance key is written into an acceptance block added at
// the end, set again in place, read back by Load, and unset so that its
// default applies; the rest of the file is as it was.
func TestEditSettingsWritesAcceptance(t *testing.T) {
	got, err := EditSettings([]byte(settingsFixture), []Assignment{
		{"acceptance.mode", "pull-request"},
		{"acceptance.forge", "github"},
		{"acceptance.merge_method", "squash"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "acceptance:\n  mode: pull-request\n  forge: github\n  merge_method: squash\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	if a := loaded(t, got).Acceptance; !a.PullRequest() || a.ForgeOrDefault() != ForgeGitHub || a.MergeMethodOrDefault() != MergeSquash {
		t.Errorf("read back %+v", a)
	}
	for _, method := range MergeMethods {
		again, err := EditSettings(got, []Assignment{{"acceptance.merge_method", method}}, nil)
		if err != nil || string(again) != strings.Replace(want, "merge_method: squash", "merge_method: "+method, 1) {
			t.Errorf("merge_method=%s: %v\n%s", method, err, again)
		}
	}
	local, err := EditSettings(got, []Assignment{{"acceptance.mode", "local"}}, nil)
	if err != nil || string(local) != strings.Replace(want, "mode: pull-request", "mode: local", 1) || loaded(t, local).Acceptance.PullRequest() {
		t.Errorf("mode=local, its forge and merge method kept: %v\n%s", err, local)
	}
	unset, err := EditSettings(got, nil, []string{"acceptance.mode", "acceptance.forge", "acceptance.merge_method"})
	if err != nil || string(unset) != settingsFixture+"acceptance:\n" {
		t.Errorf("unset: %v\n%s", err, unset)
	}
	if a := loaded(t, unset).Acceptance; a.ModeOrDefault() != AcceptLocal || a.MergeMethodOrDefault() != MergeMerge {
		t.Errorf("unset reads %+v", a)
	}
}

// S-0399: merge_queue.batch_max is written into a merge_queue block added at
// the end, set again in place, read back by Load, and unset so that its
// default, 4, applies; the rest of the file is as it was.
func TestEditSettingsWritesBatchMax(t *testing.T) {
	if s, ok := SettingFor("merge_queue.batch_max"); !ok || s.Kind != KindNumber || !s.Whole || s.Default != "4" {
		t.Fatalf("merge_queue.batch_max in the catalog: %+v, %v", s, ok)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"merge_queue.batch_max", "6"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "merge_queue:\n  batch_max: 6\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	if m := loaded(t, got); m.BatchMax() != 6 {
		t.Errorf("read back %d", m.BatchMax())
	}
	again, err := EditSettings(got, []Assignment{{"merge_queue.batch_max", "1"}}, nil)
	if err != nil || string(again) != settingsFixture+"merge_queue:\n  batch_max: 1\n" || loaded(t, again).BatchMax() != 1 {
		t.Errorf("batch_max=1: %v\n%s", err, again)
	}
	unset, err := EditSettings(got, nil, []string{"merge_queue.batch_max"})
	if err != nil || string(unset) != settingsFixture+"merge_queue:\n" || loaded(t, unset).BatchMax() != DefaultBatchMax {
		t.Errorf("unset: %v\n%s", err, unset)
	}
	var v SettingValue
	for _, sv := range loaded(t, got).SettingValues() {
		if sv.Key == "merge_queue.batch_max" {
			v = sv
		}
	}
	if !v.Set || v.Value != 6 {
		t.Errorf("its value: %+v", v)
	}
}

// ADR-0199: digest.schedule and digest.compact_after are written into a
// digest block added at the end, set again in place, 0 among them, read back
// by Load, and unset so that their defaults apply; the rest of the file is as
// it was.
func TestEditSettingsWritesTheDigestSettings(t *testing.T) {
	if s, ok := SettingFor("digest.schedule"); !ok || s.Kind != KindCron || s.Default != "" {
		t.Fatalf("digest.schedule in the catalog: %+v, %v", s, ok)
	}
	if s, ok := SettingFor("digest.compact_after"); !ok || s.Kind != KindDuration || s.Default != "720h" {
		t.Fatalf("digest.compact_after in the catalog: %+v, %v", s, ok)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"digest.schedule", "daily"}, {"digest.compact_after", "336h"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "digest:\n  schedule: daily\n  compact_after: 336h\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if s, err := m.Digest.DigestSchedule(); err != nil || s == nil || s.String() != "daily" {
		t.Errorf("schedule read back: %v, %v", s, err)
	}
	if d, err := m.Digest.CompactAfterDuration(); err != nil || d != 336*time.Hour {
		t.Errorf("window read back: %v, %v", d, err)
	}
	values := map[string]SettingValue{}
	for _, sv := range m.SettingValues() {
		values[sv.Key] = sv
	}
	if v := values["digest.schedule"]; !v.Set || v.Value != "daily" {
		t.Errorf("digest.schedule's value: %+v", v)
	}
	if v := values["digest.compact_after"]; !v.Set || v.Value != "336h" {
		t.Errorf("digest.compact_after's value: %+v", v)
	}

	off, err := EditSettings(got, []Assignment{{"digest.schedule", "0 6 * * 1"}, {"digest.compact_after", "0"}}, nil)
	if err != nil || string(off) != settingsFixture+"digest:\n  schedule: \"0 6 * * 1\"\n  compact_after: \"0\"\n" {
		t.Fatalf("set again: %v\n%s", err, off)
	}
	if d, err := loaded(t, off).Digest.CompactAfterDuration(); err != nil || d != 0 {
		t.Errorf("0 read back: %v, %v", d, err)
	}

	unset, err := EditSettings(got, nil, []string{"digest.schedule", "digest.compact_after"})
	if err != nil || string(unset) != settingsFixture+"digest:\n" {
		t.Fatalf("unset: %v\n%s", err, unset)
	}
	m = loaded(t, unset)
	if s, err := m.Digest.DigestSchedule(); err != nil || s != nil {
		t.Errorf("unset schedule: %v, %v", s, err)
	}
	if d, err := m.Digest.CompactAfterDuration(); err != nil || d != DefaultCompactAfter {
		t.Errorf("unset window: %v, %v", d, err)
	}
}

// ADR-0205: threads.remind_after and threads.remind_every are written into a
// threads block added at the end, as flai manifest set writes them, read back
// by Load, and unset so that their defaults apply; a value flai check
// refuses is refused naming the key.
func TestEditSettingsWritesTheThreadsSettings(t *testing.T) {
	if s, ok := SettingFor("threads.remind_after"); !ok || s.Kind != KindDuration || s.Default != "" {
		t.Fatalf("threads.remind_after in the catalog: %+v, %v", s, ok)
	}
	if s, ok := SettingFor("threads.remind_every"); !ok || s.Kind != KindDuration || s.Default != "24h" {
		t.Fatalf("threads.remind_every in the catalog: %+v, %v", s, ok)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"threads.remind_after", "4h"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := settingsFixture + "threads:\n  remind_after: 4h\n"; string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if d, err := m.Threads.RemindAfterDuration(); err != nil || d != 4*time.Hour {
		t.Errorf("remind_after read back: %v, %v", d, err)
	}
	if d, err := m.Threads.RemindEveryDuration(); err != nil || d != DefaultRemindEvery {
		t.Errorf("remind_every unset: %v, %v", d, err)
	}
	values := map[string]SettingValue{}
	for _, sv := range m.SettingValues() {
		values[sv.Key] = sv
	}
	if v := values["threads.remind_after"]; !v.Set || v.Value != "4h" {
		t.Errorf("threads.remind_after's value: %+v", v)
	}
	if v := values["threads.remind_every"]; v.Set {
		t.Errorf("threads.remind_every's value: %+v", v)
	}

	every, err := EditSettings(got, []Assignment{{"threads.remind_every", "12h"}}, nil)
	if err != nil || string(every) != settingsFixture+"threads:\n  remind_after: 4h\n  remind_every: 12h\n" {
		t.Fatalf("remind_every: %v\n%s", err, every)
	}
	if d, err := loaded(t, every).Threads.RemindEveryDuration(); err != nil || d != 12*time.Hour {
		t.Errorf("remind_every read back: %v, %v", d, err)
	}

	unset, err := EditSettings(every, nil, []string{"threads.remind_after", "threads.remind_every"})
	if err != nil || string(unset) != settingsFixture+"threads:\n" {
		t.Fatalf("unset: %v\n%s", err, unset)
	}
	m = loaded(t, unset)
	if d, err := m.Threads.RemindAfterDuration(); err != nil || d != 0 {
		t.Errorf("unset remind_after: %v, %v", d, err)
	}
	if d, err := m.Threads.RemindEveryDuration(); err != nil || d != DefaultRemindEvery {
		t.Errorf("unset remind_every: %v, %v", d, err)
	}

	for _, c := range []struct{ key, value, reason string }{
		{"threads.remind_after", "soon", "is not a duration"},
		{"threads.remind_after", "-1h", "is negative"},
		{"threads.remind_every", "0", "is not more than zero"},
	} {
		_, err := EditSettings([]byte(settingsFixture), []Assignment{{c.key, c.value}}, nil)
		var refused *RefusedError
		if !errors.As(err, &refused) {
			t.Errorf("%s=%s: want a refusal, got %v", c.key, c.value, err)
			continue
		}
		if i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == c.key }); i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
			t.Errorf("%s=%s: refusal %+v, want %q", c.key, c.value, refused.Problems, c.reason)
		}
	}
}

// ADR-0209: orchestration.coalesce, orchestration.run_timeout, and
// orchestration.schedule are written into the orchestration block, read back
// by Load, and unset so that their defaults apply; a value flai check refuses
// is refused naming the key.
func TestEditSettingsWritesTheOrchestrationRunSettings(t *testing.T) {
	for key, want := range map[string]struct {
		kind Kind
		def  string
	}{
		"orchestration.coalesce":    {KindDuration, "30s"},
		"orchestration.run_timeout": {KindDuration, "20m"},
		"orchestration.schedule":    {KindCron, ""},
	} {
		if s, ok := SettingFor(key); !ok || s.Kind != want.kind || s.Default != want.def || s.Meaning == "" || s.Risk != "" {
			t.Errorf("%s in the catalog: %+v, %v", key, s, ok)
		}
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{
		{"orchestration.coalesce", "2m"},
		{"orchestration.run_timeout", "1h"},
		{"orchestration.schedule", "0 * * * *"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := replaced(t, "    promote_to_ready: true\n", "    promote_to_ready: true\n  coalesce: 2m\n  run_timeout: 1h\n  schedule: \"0 * * * *\"\n")
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if d, err := m.Orchestration.CoalesceDuration(); err != nil || d != 2*time.Minute {
		t.Errorf("coalesce read back: %v, %v", d, err)
	}
	if d, err := m.Orchestration.RunTimeoutDuration(); err != nil || d != time.Hour {
		t.Errorf("run_timeout read back: %v, %v", d, err)
	}
	if s, err := m.Orchestration.OrchestrationSchedule(); err != nil || s == nil || s.String() != "0 * * * *" {
		t.Errorf("schedule read back: %v, %v", s, err)
	}
	values := map[string]SettingValue{}
	for _, sv := range m.SettingValues() {
		values[sv.Key] = sv
	}
	for key, want := range map[string]string{"orchestration.coalesce": "2m", "orchestration.run_timeout": "1h", "orchestration.schedule": "0 * * * *"} {
		if v := values[key]; !v.Set || v.Value != want {
			t.Errorf("%s's value: %+v", key, v)
		}
	}
	keys := []string{"orchestration.coalesce", "orchestration.run_timeout", "orchestration.schedule"}
	for _, v := range loaded(t, []byte(settingsFixture)).SettingValues() {
		if slices.Contains(keys, v.Key) && (v.Set || v.Value != nil) {
			t.Errorf("%s unset: %+v", v.Key, v)
		}
	}

	unset, err := EditSettings(got, nil, keys)
	if err != nil || string(unset) != settingsFixture {
		t.Fatalf("unset: %v\n%s", err, unset)
	}
	m = loaded(t, unset)
	if d, err := m.Orchestration.CoalesceDuration(); err != nil || d != DefaultCoalesce {
		t.Errorf("unset coalesce: %v, %v", d, err)
	}
	if d, err := m.Orchestration.RunTimeoutDuration(); err != nil || d != DefaultRunTimeout {
		t.Errorf("unset run_timeout: %v, %v", d, err)
	}
	if s, err := m.Orchestration.OrchestrationSchedule(); err != nil || s != nil {
		t.Errorf("unset schedule: %v, %v", s, err)
	}

	for _, c := range []struct{ key, value, reason string }{
		{"orchestration.coalesce", "soon", "is not a duration; write one such as 30s"},
		{"orchestration.coalesce", "0", "is not more than zero"},
		{"orchestration.coalesce", "-1m", "is not more than zero"},
		{"orchestration.run_timeout", "an hour", "is not a duration; write one such as 20m"},
		{"orchestration.run_timeout", "0s", "is not more than zero"},
		{"orchestration.schedule", "hourly", `"hourly" has 1 fields, not five`},
	} {
		_, err := EditSettings([]byte(settingsFixture), []Assignment{{c.key, c.value}}, nil)
		var refused *RefusedError
		if !errors.As(err, &refused) {
			t.Errorf("%s=%s: want a refusal, got %v", c.key, c.value, err)
			continue
		}
		if i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == c.key }); i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
			t.Errorf("%s=%s: refusal %+v, want %q", c.key, c.value, refused.Problems, c.reason)
		}
	}
}

// S-0502: orchestration.permissions.confirm_after is a duration, unset by
// default, that flai manifest set writes among the permissions, Load reads
// back, and unset turns off; a value flai check refuses is refused naming the
// key (ADR-0206).
func TestEditSettingsWritesConfirmAfter(t *testing.T) {
	const key = "orchestration.permissions.confirm_after"
	s, ok := SettingFor(key)
	if !ok || s.Kind != KindDuration || s.Default != "" || len(s.Values) != 0 || s.Meaning == "" || s.Risk == "" {
		t.Fatalf("%s in the catalog: %+v, %v", key, s, ok)
	}
	keys := settingKeys()
	if i := slices.Index(keys, key); i < 1 || keys[i-1] != "orchestration.permissions.answer_threads" {
		t.Errorf("%s is not listed after answer_threads: %v", key, keys)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{key, "48h"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := replaced(t, "    promote_to_ready: true\n", "    promote_to_ready: true\n    confirm_after: 48h\n"); string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if d, err := m.Orchestration.Permissions.ConfirmAfterDuration(); err != nil || d != 48*time.Hour || !m.Orchestration.Permissions.Allows(PermitConfirmAfter) {
		t.Errorf("confirm_after read back: %v, %v, allowed %v", d, err, m.Orchestration.Permissions.Allows(PermitConfirmAfter))
	}
	values := map[string]SettingValue{}
	for _, sv := range m.SettingValues() {
		values[sv.Key] = sv
	}
	if v := values[key]; !v.Set || v.Value != "48h" {
		t.Errorf("%s's value: %+v", key, v)
	}
	for _, v := range loaded(t, []byte(settingsFixture)).SettingValues() {
		if v.Key == key && (v.Set || v.Value != nil) {
			t.Errorf("%s unset: %+v", key, v)
		}
	}

	unset, err := EditSettings(got, nil, []string{key})
	if err != nil || string(unset) != settingsFixture {
		t.Fatalf("unset: %v\n%s", err, unset)
	}
	if loaded(t, unset).Orchestration.Permissions.Allows(PermitConfirmAfter) {
		t.Error("confirm_after is on unset")
	}

	for _, c := range []struct{ value, reason string }{
		{"soon", "is not a duration; write one such as 48h"},
		{"-1h", "is negative; write a duration such as 48h"},
	} {
		_, err := EditSettings([]byte(settingsFixture), []Assignment{{key, c.value}}, nil)
		var refused *RefusedError
		if !errors.As(err, &refused) {
			t.Errorf("%s=%s: want a refusal, got %v", key, c.value, err)
			continue
		}
		if i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == key }); i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
			t.Errorf("%s=%s: refusal %+v, want %q", key, c.value, refused.Problems, c.reason)
		}
	}
}

// S-0229: a block written on one line is refused rather than rewritten.
func TestEditSettingsRefusesAFlowBlock(t *testing.T) {
	data := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\norchestration: {policy: fifo}\n"
	_, err := EditSettings([]byte(data), []Assignment{{"orchestration.policy", "cod"}}, nil)
	var refused *RefusedError
	if !errors.As(err, &refused) || len(refused.Problems) != 1 || refused.Problems[0].Field != "orchestration" || !strings.Contains(refused.Problems[0].Reason, "written on one line") {
		t.Errorf("want orchestration refused as written on one line, got %v", err)
	}
}

// S-0229: unsetting a key removes it and its value, so its default applies;
// unsetting one that is not there changes nothing; an agent that sets
// nothing is an unset.
func TestEditSettingsUnsets(t *testing.T) {
	got, err := EditSettings([]byte(settingsFixture), nil, []string{"planning.cycle", "orchestration.permissions.promote_to_ready", "planning.agent", "planning.schedule"})
	if err != nil {
		t.Fatal(err)
	}
	want := strings.NewReplacer(
		"  cycle: 168h                 # a week\n", "",
		"  agent:                      # the planner's agent\n    model: claude-sonnet-5\n    config:\n      effort: medium\n", "",
		"    promote_to_ready: true\n", "",
	).Replace(settingsFixture)
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if m.Orchestration.Permissions.Allows(PermitPromoteToReady) || m.Planning.Agent != nil || m.Planning.Cycle != "" {
		t.Errorf("unset keys still read: %+v", m)
	}
	got, err = EditSettings([]byte(settingsFixture), []Assignment{{"planning.agent", "{}"}}, nil)
	if err != nil || strings.Contains(string(got), "claude-sonnet-5") {
		t.Errorf("an empty agent unsets it: %v\n%s", err, got)
	}
}

// S-0417: roles.<name>.agent, for a role's kebab-case name, is written into a
// roles block added at the end, read back as the role's agent merged over the
// project's, set again in place, and unset; a name not in kebab-case is no
// setting, and a built-in role's, or an agent flai check refuses, is refused
// naming the key.
func TestEditSettingsWritesARolesAgent(t *testing.T) {
	const key = "roles.code-reviewer.agent"
	if s, ok := SettingFor(key); !ok || s.Key != key || s.Kind != KindAgent || s.Default != "" || !strings.Contains(s.Meaning, "code-reviewer") {
		t.Fatalf("%s: %+v, %v", key, s, ok)
	}
	for _, k := range []string{"roles.Reviewer.agent", "roles..agent", "roles.a.b.agent", "roles.reviewer.permissions", "roles.reviewer"} {
		if _, ok := SettingFor(k); ok {
			t.Errorf("%s is a setting", k)
		}
	}
	if keys := settingKeys(); !strings.HasPrefix(keys[len(keys)-1], "roles.<name>.agent") {
		t.Errorf("roles.<name>.agent is not listed last: %v", keys)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{key, `{"model":"claude-sonnet-5","config":{"effort":"low"}}`}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "roles:\n  code-reviewer:\n    agent:\n      model: claude-sonnet-5\n      config:\n        effort: low\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if a := m.RoleAgent("code-reviewer"); a.Model != "claude-sonnet-5" || a.Config["effort"] != "low" {
		t.Errorf("read back: %v", a)
	}

	again, err := EditSettings(got, []Assignment{{key, `{"harness":"claude-code"}`}}, nil)
	if err != nil || string(again) != settingsFixture+"roles:\n  code-reviewer:\n    agent:\n      harness: claude-code\n" {
		t.Fatalf("set again: %v\n%s", err, again)
	}
	unset, err := EditSettings(got, nil, []string{key})
	if err != nil || string(unset) != settingsFixture+"roles:\n  code-reviewer:\n" {
		t.Fatalf("unset: %v\n%s", err, unset)
	}
	if a := loaded(t, unset).Roles["code-reviewer"].Agent; a != nil {
		t.Errorf("unset agent still reads: %v", a)
	}

	for _, c := range []struct{ set, field, reason string }{
		{`roles.Reviewer.agent={"model":"x"}`, "roles.Reviewer.agent", "is not a setting flai manifest set writes"},
		{`roles.planner.agent={"model":"claude-sonnet-5"}`, "roles.planner.agent", "is not read: the planner's agent is planning.agent; set that instead"},
		{`roles.code-reviewer.agent={"harness":"Bad Harness"}`, "roles.code-reviewer.agent", `harness "Bad Harness" is not a name such as claude-code`},
		{`roles.code-reviewer.agent=claude`, "roles.code-reviewer.agent", "is not an agent"},
	} {
		k, v, _ := strings.Cut(c.set, "=")
		_, err := EditSettings([]byte(settingsFixture), []Assignment{{k, v}}, nil)
		var refused *RefusedError
		if !errors.As(err, &refused) {
			t.Errorf("%s: want a refusal, got %v", c.set, err)
			continue
		}
		i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == c.field })
		if i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
			t.Errorf("%s: want %s: ...%s..., got %+v", c.set, c.field, c.reason, refused.Problems)
		}
	}
}

// S-0229: every value the catalog lists for a choice or a boolean, and every
// default, is one flai check accepts, and a value outside the list is one it
// refuses: the catalog and the validation cannot disagree.
func TestSettingsCatalogAgreesWithTheValidation(t *testing.T) {
	base := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
	// what a release policy needs beside it to be valid
	with := map[string][]Assignment{
		ReleaseThreshold: {{"orchestration.release.count", "1"}},
		ReleaseTheme:     {{"orchestration.release.epic", "E-0001"}},
	}
	seen := map[string]bool{}
	for _, s := range Settings() {
		if seen[s.Key] {
			t.Errorf("%s is in the catalog twice", s.Key)
		}
		seen[s.Key] = true
		if s.Meaning == "" {
			t.Errorf("%s has no meaning", s.Key)
		}
		if strings.HasPrefix(s.Key, "orchestration.permissions.") != (s.Risk != "") {
			t.Errorf("%s: a permission, and only one, has a risk: %q", s.Key, s.Risk)
		}
		values := slices.Clone(s.Values)
		if s.Default != "" {
			values = append(values, s.Default)
		}
		for _, v := range values {
			set := append([]Assignment{{s.Key, v}}, with[v]...)
			if _, err := EditSettings([]byte(base), set, nil); err != nil {
				t.Errorf("%s=%s is listed but refused: %v", s.Key, v, err)
			}
		}
		if s.Kind == KindChoice {
			if _, err := EditSettings([]byte(base), []Assignment{{s.Key, "none-of-these"}}, nil); err == nil {
				t.Errorf("%s=none-of-these is not listed but accepted", s.Key)
			}
		}
	}
	for _, name := range PermissionNames {
		if !seen["orchestration.permissions."+name] {
			t.Errorf("the permission %s is not in the catalog", name)
		}
	}
	if seen["planning.currency"] {
		t.Error("planning.currency is in the catalog")
	}
}

// S-0229: each setting's value in a manifest, and whether it is set.
func TestSettingValues(t *testing.T) {
	m := loaded(t, []byte(settingsFixture))
	got := map[string]SettingValue{}
	for _, v := range m.SettingValues() {
		got[v.Key] = v
	}
	if len(got) != len(Settings()) {
		t.Fatalf("%d values for %d settings", len(got), len(Settings()))
	}
	for key, want := range map[string]any{
		"planning.cycle":                             "168h",
		"orchestration.release.count":                5,
		"orchestration.permissions.promote_to_ready": true,
		"orchestration.policy":                       OrderFIFO,
	} {
		if v := got[key]; !v.Set || v.Value != want {
			t.Errorf("%s = %v (set %v), want %v", key, v.Value, v.Set, want)
		}
	}
	if a, ok := got["planning.agent"].Value.(*Agent); !ok || a.Model != "claude-sonnet-5" {
		t.Errorf("planning.agent = %v", got["planning.agent"].Value)
	}
	for _, key := range []string{"orchestration.permissions.publish", "planning.hour_rate", "analysis.schedule", "analysis.agent", "digest.schedule", "digest.compact_after", "acceptance.mode", "acceptance.forge", "acceptance.merge_method", "merge_queue.batch_max"} {
		if v := got[key]; v.Set || v.Value != nil {
			t.Errorf("%s is unset, got %v", key, v.Value)
		}
	}
	if d := got["planning.cycle"].Default; d != "168h" {
		t.Errorf("planning.cycle's default = %q, want 168h", d)
	}
	if d := got["planning.default_duration"].Default; d != "1h" {
		t.Errorf("planning.default_duration's default = %q, want 1h", d)
	}
	for key, want := range map[string]string{"digest.compact_after": "720h", "acceptance.mode": AcceptLocal, "acceptance.forge": ForgeGitHub, "acceptance.merge_method": MergeMerge, "merge_queue.batch_max": "4"} {
		if d := got[key].Default; d != want {
			t.Errorf("%s's default = %q, want %q", key, d, want)
		}
	}
}

// loaded is the manifest data as Load reads it.
func loaded(t *testing.T, data []byte) Manifest {
	t.Helper()
	file := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(file)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// S-0273: tests is written from JSON as a block, one tier an item, read back
// as given, replaced whole where it is with its key's comment kept, and
// removed whole; [] writes a list of none.
func TestEditSettingsWritesTests(t *testing.T) {
	value := `[{"name":"unit","command":["go","test","-json","{packages}"],"dir":"flai","paths":["flai/**","!flai/testdata/**"],"format":"go-test-json","all_command":["go","test","-json","./..."]},{"name":"smoke","command":["scripts/smoke.sh"],"all_only":true}]`
	block := "tests:\n" +
		"  - name: unit\n" +
		"    command: [go, test, \"-json\", \"{packages}\"]\n" +
		"    dir: flai\n" +
		"    paths: [\"flai/**\", \"!flai/testdata/**\"]\n" +
		"    format: go-test-json\n" +
		"    all_command: [go, test, \"-json\", \"./...\"]\n" +
		"  - name: smoke\n" +
		"    command: [scripts/smoke.sh]\n" +
		"    all_only: true\n"
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"tests", value}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := settingsFixture + block; string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if len(m.Tests) != 2 || m.Tests[0].Format != FormatGoTestJSON || !m.Tests[1].AllOnly || m.Tests[0].Paths[1] != "!flai/testdata/**" {
		t.Errorf("reads back: %+v", m.Tests)
	}
	v := m.SettingValues()[slices.IndexFunc(m.SettingValues(), func(v SettingValue) bool { return v.Key == "tests" })]
	if tiers, ok := v.Value.([]TestTier); !v.Set || !ok || !sameTiers(tiers, m.Tests) {
		t.Errorf("tests' setting value = %+v", v)
	}

	// a list written with its items at the key's indent, and a comment, is
	// replaced whole, and the keys after it are kept
	before := settingsFixture + "tests: # cheapest first\n- name: old\n  command: [make, test]\n  paths: [src]\nflai:\n  minimum: 1.27.0\n"
	got, err = EditSettings([]byte(before), []Assignment{{"tests", `[{"name":"all","command":["make","check"],"all_only":true}]`}}, nil)
	want := settingsFixture + "tests: # cheapest first\n  - name: all\n    command: [make, check]\n    all_only: true\nflai:\n  minimum: 1.27.0\n"
	if err != nil || string(got) != want {
		t.Errorf("replace: %v\ngot:\n%s\nwant:\n%s", err, got, want)
	}

	got, err = EditSettings([]byte(before), []Assignment{{"tests", "[]"}}, nil)
	if want := settingsFixture + "tests: [] # cheapest first\nflai:\n  minimum: 1.27.0\n"; err != nil || string(got) != want {
		t.Errorf("none: %v\ngot:\n%s\nwant:\n%s", err, got, want)
	}
	if m := loaded(t, got); m.Tests == nil || len(m.Tests) != 0 {
		t.Errorf("tests: [] reads back as %#v", m.Tests)
	}

	got, err = EditSettings([]byte(before), nil, []string{"tests"})
	if want := settingsFixture + "flai:\n  minimum: 1.27.0\n"; err != nil || string(got) != want {
		t.Errorf("unset: %v\ngot:\n%s\nwant:\n%s", err, got, want)
	}
}

// S-0342: a tier's covers is written as a list on one line and read back as
// given, and a list of tiers that differs only in it is not the same.
func TestEditSettingsWritesCovers(t *testing.T) {
	value := `[{"name":"short","command":["go","test","-short"],"paths":["**"]},{"name":"full","command":["go","test"],"paths":["**"],"covers":["short"]}]`
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"tests", value}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "tests:\n" +
		"  - name: short\n" +
		"    command: [go, test, \"-short\"]\n" +
		"    paths: [\"**\"]\n" +
		"  - name: full\n" +
		"    command: [go, test]\n" +
		"    paths: [\"**\"]\n" +
		"    covers: [short]\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if len(m.Tests) != 2 || m.Tests[0].Covers != nil || !slices.Equal(m.Tests[1].Covers, []string{"short"}) {
		t.Errorf("reads back: %+v", m.Tests)
	}
	uncovered := slices.Clone(m.Tests)
	uncovered[1].Covers = nil
	if sameTiers(m.Tests, uncovered) {
		t.Error("tiers that differ in covers are the same")
	}
}

// S-0405: a tier's queue is written, read back, and makes two tiers differ.
func TestEditSettingsWritesQueue(t *testing.T) {
	value := `[{"name":"smoke","command":["scripts/smoke.sh"],"paths":["**"],"all_only":true,"queue":true}]`
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"tests", value}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "tests:\n" +
		"  - name: smoke\n" +
		"    command: [scripts/smoke.sh]\n" +
		"    paths: [\"**\"]\n" +
		"    all_only: true\n" +
		"    queue: true\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if len(m.Tests) != 1 || !m.Tests[0].Queue {
		t.Errorf("reads back: %+v", m.Tests)
	}
	unqueued := slices.Clone(m.Tests)
	unqueued[0].Queue = false
	if sameTiers(m.Tests, unqueued) {
		t.Error("tiers that differ in queue are the same")
	}
}

// S-0273: a value that is not a list of tiers, and a tier the manifest's
// validation refuses, are refused with the field and the reason, and nothing
// is written.
func TestEditSettingsRefusesBadTests(t *testing.T) {
	for _, c := range []struct {
		name, value, field, reason string
	}{
		{"not a list", `{"name":"unit"}`, "tests", "is not a list of test tiers; write one as JSON"},
		{"not JSON", `[unit]`, "tests", "is not a list of test tiers"},
		{"an unknown field", `[{"name":"unit","command":["x"],"paths":["**"],"timeout":5}]`, "tests", "is not a list of test tiers"},
		{"a bad format", `[{"name":"unit","command":["x"],"paths":["**"],"format":"junit"}]`, "tests[0].format", `(tier "unit") "junit" is not a format flai reads`},
		{"no paths", `[{"name":"unit","command":["x"]}]`, "tests[0].paths", "is empty"},
		{"a name twice", `[{"name":"a","command":["x"],"paths":["**"]},{"name":"a","command":["y"],"paths":["**"]}]`, "tests[1].name", "is tests[0]'s name too"},
	} {
		t.Run(c.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), File)
			if err := os.WriteFile(file, []byte(settingsFixture), 0o644); err != nil {
				t.Fatal(err)
			}
			err := WriteSettings(file, []Assignment{{"tests", c.value}}, nil)
			var refused *RefusedError
			if !errors.As(err, &refused) {
				t.Fatalf("want a refusal, got %v", err)
			}
			i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == c.field })
			if i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
				t.Errorf("want %s: ...%s..., got %+v", c.field, c.reason, refused.Problems)
			}
			if data, _ := os.ReadFile(file); string(data) != settingsFixture {
				t.Errorf("the file changed:\n%s", data)
			}
		})
	}
}

// componentsFixture is a manifest with comments and blank lines around its
// sub-projects, one of each kind the components key concerns.
const componentsFixture = `version: 1
name: demo
layout:
  design: design
  docs: docs
  wip: wip
projects:                 # the releasable components
  - name: api
    path: api
    kind: go
    tags: [cli]
  # the app
  - name: web             # the dashboard
    path: web
    kind: sveltekit
    tags: [dashboard]   # its tag

  - name: ui
    path: packages/ui
    kind: svelte
dashboard:
  port: 4242
`

// componentsFiles are the files the fixture's setup modules name.
var componentsFiles = fstest.MapFS{
	"web/src/app.css":         {Data: []byte("@import 'tailwindcss';\n")},
	"web/src/theme.css":       {Data: nil},
	"packages/ui/src/ui.css":  {Data: nil},
	"packages/ui/src/lib/x.y": {Data: nil},
}

// S-0437: EditComponents adds a sub-project's components key after its last
// line, replaces it whole keeping its comment, and removes it, and every
// other line, the order of the keys, and the comments stay as they were.
func TestEditComponents(t *testing.T) {
	webBlock := "    components:\n      include: [\"src/lib/**/*.svelte\", \"src/stories/*.svelte\"]\n      exclude: [\"src/lib/internal/**\"]\n      setup: src/app.css\n"
	withWeb := strings.Replace(componentsFixture, "    tags: [dashboard]   # its tag\n", "    tags: [dashboard]   # its tag\n"+webBlock, 1)
	commented := strings.Replace(withWeb, "    components:\n", "    components:           # the catalogue\n", 1)
	for _, c := range []struct {
		name, data, project string
		set                 *Components
		want                string
	}{
		{"added to a sub-project in the middle", componentsFixture, "web",
			&Components{Include: []string{"src/lib/**/*.svelte", "src/stories/*.svelte"}, Exclude: []string{"src/lib/internal/**"}, Setup: "src/app.css"},
			withWeb},
		{"added to the last sub-project", componentsFixture, "ui", &Components{Setup: "src/ui.css"},
			strings.Replace(componentsFixture, "    kind: svelte\n", "    kind: svelte\n    components:\n      setup: src/ui.css\n", 1)},
		{"replaced whole, its comment kept", commented, "web", &Components{Exclude: []string{"src/lib/x.svelte"}},
			strings.Replace(commented, webBlock[len("    components:\n"):], "      exclude: [src/lib/x.svelte]\n", 1)},
		{"removed", commented, "web", nil, componentsFixture},
		{"removed where there is none", componentsFixture, "ui", nil, componentsFixture},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, err := EditComponents([]byte(c.data), componentsFiles, c.project, c.set)
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != c.want {
				t.Errorf("got:\n%s\nwant:\n%s", out, c.want)
			}
			m := loaded(t, out)
			i := slices.IndexFunc(m.Projects, func(p Project) bool { return p.Name == c.project })
			if !sameComponents(m.Projects[i].Components, c.set) {
				t.Errorf("reads back as %+v, want %+v", m.Projects[i].Components, c.set)
			}
		})
	}
}

// S-0437: a sub-project's entry whose keys are written after a lone dash, or
// further in than two spaces, is edited at its own indent.
func TestEditComponentsKeepsTheEntrysIndent(t *testing.T) {
	data := "version: 1\nname: demo\nprojects:\n-   name: web\n    path: web\n    kind: sveltekit\n-\n  name: ui\n  path: ui\n  kind: svelte\n"
	out, err := EditComponents([]byte(data), nil, "ui", &Components{Include: []string{"lib/*.svelte"}})
	want := strings.Replace(data, "  kind: svelte\n", "  kind: svelte\n  components:\n    include: [\"lib/*.svelte\"]\n", 1)
	if err != nil || string(out) != want {
		t.Errorf("after a lone dash: %v\n%s", err, out)
	}
	out, err = EditComponents([]byte(data), nil, "web", &Components{Setup: "src/app.css"})
	want = strings.Replace(data, "    kind: sveltekit\n", "    kind: sveltekit\n    components:\n      setup: src/app.css\n", 1)
	if err != nil || string(out) != want {
		t.Errorf("further in: %v\n%s", err, out)
	}
}

// S-0437: EditComponents refuses what ComponentErrors finds in the manifest
// as it would be, with the field and the reason, an unknown sub-project, a
// key that sets nothing, and an entry it cannot edit, and SetComponents then
// writes nothing.
func TestEditComponentsRefuses(t *testing.T) {
	for _, c := range []struct {
		name, project string
		set           *Components
		field, reason string
	}{
		{"a glob that does not parse", "web", &Components{Include: []string{"src/[ab.svelte"}}, "projects[1].components.include", `(sub-project "web") pattern "src/[ab.svelte" is not a valid glob`},
		{"a setup that names no file", "web", &Components{Setup: "src/missing.css"}, "projects[1].components.setup", `(sub-project "web") "src/missing.css" names no file under web`},
		{"a setup under another sub-project", "ui", &Components{Setup: "../../web/src/app.css"}, "projects[2].components.setup", "leaves the sub-project"},
		{"an unsupported kind", "api", &Components{Setup: "main.go"}, "projects[0].components", `(sub-project "api") is set on a sub-project of kind "go"`},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := EditComponents([]byte(componentsFixture), componentsFiles, c.project, c.set)
			var refused *RefusedError
			if !errors.As(err, &refused) {
				t.Fatalf("want a refusal, got %v", err)
			}
			i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == c.field })
			if i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
				t.Errorf("want %s: ...%s..., got %+v", c.field, c.reason, refused.Problems)
			}
		})
	}
	for _, c := range []struct {
		name, data, project string
		set                 *Components
		want                string
	}{
		{"an unknown sub-project", componentsFixture, "nope", &Components{Setup: "x.css"}, `system-flow.yaml has no sub-project "nope"; name one of api, web, ui`},
		{"no sub-projects", "version: 1\nname: demo\n", "web", nil, `system-flow.yaml has no sub-project "web", and names none under projects`},
		{"nothing to set", componentsFixture, "web", &Components{}, `nothing to set for sub-project "web"'s components`},
		{"an entry on one line", "version: 1\nname: demo\nprojects:\n  - {name: web, path: web, kind: sveltekit}\n", "web", &Components{Setup: "x.css"}, "its entry under projects is written on one line"},
		{"projects on one line", "version: 1\nname: demo\nprojects: [{name: web, path: web, kind: sveltekit}]\n", "web", &Components{Setup: "x.css"}, "projects is written on one line"},
		{"an entry that begins with components", "version: 1\nname: demo\nprojects:\n  - components: {setup: x.css}\n    name: web\n    path: web\n    kind: sveltekit\n", "web", nil, "begins with components"},
	} {
		t.Run(c.name, func(t *testing.T) {
			if _, err := EditComponents([]byte(c.data), nil, c.project, c.set); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("want ...%s..., got %v", c.want, err)
			}
		})
	}
	root := t.TempDir()
	file := filepath.Join(root, File)
	if err := os.WriteFile(file, []byte(componentsFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SetComponents(file, "web", &Components{Setup: "src/app.css"}); err == nil {
		t.Error("SetComponents wrote a setup that names no file under the manifest's folder")
	}
	if data, _ := os.ReadFile(file); string(data) != componentsFixture {
		t.Errorf("the file changed on a refusal:\n%s", data)
	}
	if err := os.MkdirAll(filepath.Join(root, "web", "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "web", "src", "app.css"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SetComponents(file, "web", &Components{Setup: "src/app.css"}); err != nil {
		t.Fatal(err)
	}
	want := strings.Replace(componentsFixture, "    tags: [dashboard]   # its tag\n", "    tags: [dashboard]   # its tag\n    components:\n      setup: src/app.css\n", 1)
	if data, _ := os.ReadFile(file); string(data) != want {
		t.Errorf("SetComponents wrote:\n%s", data)
	}
}
