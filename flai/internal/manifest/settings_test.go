package manifest

import (
	"errors"
	"fmt"
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
		// S-0692: analysis.formula is flai manifest set's to write
		{KindText, Assignment{"analysis.formula", "analysis-cadence"},
			settingsFixture + "analysis:\n  formula: analysis-cadence\n"},
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
		{"formula's name", []Assignment{{"analysis.formula", "Analysis Cadence"}}, nil, "analysis.formula", `"Analysis Cadence" is not a formula's name such as analysis-cadence`},
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
		{"coalesce window", []Assignment{{"merge_queue.coalesce", "soon"}}, nil, "merge_queue.coalesce", `"soon" is not a duration; write how long the merge queue waits for a batch to fill, such as 3m`},
		{"negative coalesce window", []Assignment{{"merge_queue.coalesce", "-1m"}}, nil, "merge_queue.coalesce", `"-1m" is negative`},
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

// S-0552: merge_queue.coalesce is written into a merge_queue block added at
// the end, 0 among its values, read back by Load, and unset so that its
// default, 3m, applies; the rest of the file is as it was.
func TestEditSettingsWritesMergeQueueCoalesce(t *testing.T) {
	if s, ok := SettingFor("merge_queue.coalesce"); !ok || s.Kind != KindDuration || s.Default != "3m" {
		t.Fatalf("merge_queue.coalesce in the catalog: %+v, %v", s, ok)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"merge_queue.coalesce", "90s"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := settingsFixture + "merge_queue:\n  coalesce: 90s\n"; string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	if d := loaded(t, got).MergeQueueCoalesce(); d != 90*time.Second {
		t.Errorf("read back %v", d)
	}
	off, err := EditSettings(got, []Assignment{{"merge_queue.coalesce", "0"}}, nil)
	if err != nil || string(off) != settingsFixture+"merge_queue:\n  coalesce: \"0\"\n" || loaded(t, off).MergeQueueCoalesce() != 0 {
		t.Errorf("coalesce=0: %v\n%s", err, off)
	}
	unset, err := EditSettings(got, nil, []string{"merge_queue.coalesce"})
	if err != nil || string(unset) != settingsFixture+"merge_queue:\n" || loaded(t, unset).MergeQueueCoalesce() != DefaultMergeQueueCoalesce {
		t.Errorf("unset: %v\n%s", err, unset)
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
	if keys := settingKeys(); !strings.HasPrefix(keys[len(keys)-2], "roles.<name>.agent") {
		t.Errorf("roles.<name>.agent is not listed before roles.<name>.schedule, last: %v", keys)
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

// S-0418: roles.<name>.schedule, for a role's kebab-case name, is a cron
// setting written into a roles block, read back as the role's schedule, and
// unset; it is listed last, and a value flai check refuses, as one that
// never comes round or a built-in role's, is refused naming the key.
func TestEditSettingsWritesARolesSchedule(t *testing.T) {
	const key = "roles.code-reviewer.schedule"
	if s, ok := SettingFor(key); !ok || s.Key != key || s.Kind != KindCron || !strings.Contains(s.Meaning, "code-reviewer") {
		t.Fatalf("%s: %+v, %v", key, s, ok)
	}
	for _, k := range []string{"roles.Reviewer.schedule", "roles..schedule", "roles.a.b.schedule"} {
		if _, ok := SettingFor(k); ok {
			t.Errorf("%s is a setting", k)
		}
	}
	if keys := settingKeys(); keys[len(keys)-1] != "roles.<name>.schedule" {
		t.Errorf("roles.<name>.schedule is not listed last: %v", keys)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{key, "0 6 * * 1"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sched, err := loaded(t, got).RoleSchedule("code-reviewer"); err != nil || sched == nil || sched.String() != "0 6 * * 1" {
		t.Errorf("read back: %v, %v\n%s", sched, err, got)
	}
	unset, err := EditSettings(got, nil, []string{key})
	if err != nil {
		t.Fatal(err)
	}
	if v := loaded(t, unset).Roles["code-reviewer"].Schedule; v != "" {
		t.Errorf("unset schedule still reads: %q", v)
	}
	for _, c := range []struct{ set, reason string }{
		{"roles.code-reviewer.schedule=0 0 31 2 *", "never comes round"},
		{"roles.planner.schedule=daily", "is not read"},
	} {
		k, v, _ := strings.Cut(c.set, "=")
		_, err := EditSettings([]byte(settingsFixture), []Assignment{{k, v}}, nil)
		var refused *RefusedError
		if !errors.As(err, &refused) {
			t.Errorf("%s: want a refusal, got %v", c.set, err)
			continue
		}
		if i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == k }); i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
			t.Errorf("%s: want %s: ...%s..., got %+v", c.set, k, c.reason, refused.Problems)
		}
	}
}

// S-0420: roles.<name>.permissions.<p>, for a role's kebab-case name and a
// permission a role file may declare, is a boolean, or answer_threads a
// choice, written into the role's permissions block and read back as the
// role's permission; an answer mode flai check refuses, the orchestrator's,
// and confirm_after are refused naming the key, and a permission that is none
// is no setting. It is listed before roles.<name>.agent.
func TestEditSettingsWritesARolesPermission(t *testing.T) {
	const key = "roles.code-reviewer.permissions.publish"
	if s, ok := SettingFor(key); !ok || s.Key != key || s.Kind != KindBoolean || s.Default != "false" || !strings.Contains(s.Meaning, "code-reviewer") {
		t.Fatalf("%s: %+v, %v", key, s, ok)
	}
	if s, ok := SettingFor("roles.code-reviewer.permissions.answer_threads"); !ok || s.Kind != KindChoice || !slices.Equal(s.Values, AnswerModes) || s.Default != AnswerOff {
		t.Errorf("answer_threads: %+v, %v", s, ok)
	}
	for _, k := range []string{"roles.code-reviewer.permissions.fly", "roles.Reviewer.permissions.publish", "roles..permissions.publish", "roles.code-reviewer.permissions."} {
		if _, ok := SettingFor(k); ok {
			t.Errorf("%s is a setting", k)
		}
	}
	if keys := settingKeys(); !strings.HasPrefix(keys[len(keys)-3], "roles.<name>.permissions.<p>") {
		t.Errorf("roles.<name>.permissions.<p> is not listed before roles.<name>.agent: %v", keys)
	}
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{key, "true"}, {"roles.code-reviewer.permissions.answer_threads", "recommend"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := loaded(t, got)
	if p := m.Roles["code-reviewer"].Permissions; !p.Allows(PermitPublish) || p.AnswerMode() != AnswerRecommend || p.Allows(PermitAcceptReviews) {
		t.Errorf("read back: %+v\n%s", p, got)
	}
	unset, err := EditSettings(got, nil, []string{key})
	if err != nil {
		t.Fatal(err)
	}
	if loaded(t, unset).Roles["code-reviewer"].Permissions.Allows(PermitPublish) {
		t.Errorf("unset permission still reads:\n%s", unset)
	}
	for _, c := range []struct{ set, reason string }{
		{"roles.code-reviewer.permissions.answer_threads=always", "is not a way of answering threads"},
		{"roles.orchestrator.permissions.publish=true", "set orchestration.permissions.publish instead"},
		{"roles.code-reviewer.permissions.confirm_after=1h", "confirms only the orchestrator's recommendations"},
		{"roles.code-reviewer.permissions.publish=yes", "is not true or false"},
	} {
		k, v, _ := strings.Cut(c.set, "=")
		_, err := EditSettings([]byte(settingsFixture), []Assignment{{k, v}}, nil)
		var refused *RefusedError
		if !errors.As(err, &refused) {
			t.Errorf("%s: want a refusal, got %v", c.set, err)
			continue
		}
		if i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return strings.HasPrefix(p.Field, k) }); i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
			t.Errorf("%s: want %s: ...%s..., got %+v", c.set, k, c.reason, refused.Problems)
		}
	}
}

// S-0420: a defined role's settings are its agent, the schedules its
// triggers name, and the permissions its file declares, each with its value
// in the manifest, a key named twice once; a built-in role has none here.
func TestRoleSettingValues(t *testing.T) {
	m := Manifest{Roles: map[string]RoleSettings{"reviewer": {Schedule: "daily", Permissions: Permissions{Publish: true}}},
		Analysis: Analysis{Schedule: "0 6 * * *"}}
	got := m.RoleSettingValues("reviewer", []string{"roles.reviewer.schedule", "analysis.schedule", "roles.reviewer.schedule"}, []string{PermitPublish, PermitAnswerThreads, "fly"})
	var keys []string
	for _, v := range got {
		keys = append(keys, fmt.Sprintf("%s=%v/%v", v.Key, v.Value, v.Set))
	}
	want := []string{"roles.reviewer.agent=<nil>/false", "roles.reviewer.schedule=daily/true", "analysis.schedule=0 6 * * */true",
		"roles.reviewer.permissions.publish=true/true", "roles.reviewer.permissions.answer_threads=<nil>/false"}
	if !slices.Equal(keys, want) {
		t.Errorf("got %v\nwant %v", keys, want)
	}
	if got := m.RoleSettingValues("planner", []string{"planning.schedule"}, nil); got != nil {
		t.Errorf("a built-in role's: %v", got)
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
			if strings.HasPrefix(s.Key, "kluster.") && s.Key != "kluster.config" {
				// a kluster: block needs its config beside any other key
				set = append(set, Assignment{"kluster.config", "kluster.yaml"})
			}
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
	for _, key := range []string{"orchestration.permissions.publish", "planning.hour_rate", "analysis.schedule", "analysis.formula", "analysis.agent", "digest.schedule", "digest.compact_after", "acceptance.mode", "acceptance.forge", "acceptance.merge_method", "merge_queue.batch_max", "merge_queue.coalesce"} {
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
	for key, want := range map[string]string{"digest.compact_after": "720h", "acceptance.mode": AcceptLocal, "acceptance.forge": ForgeGitHub, "acceptance.merge_method": MergeMerge, "merge_queue.batch_max": "4", "merge_queue.coalesce": "3m"} {
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

// S-0548: a tier's manual is written, read back, and makes two tiers differ,
// and a manual tier needs no paths.
func TestEditSettingsWritesManual(t *testing.T) {
	value := `[{"name":"gateway","command":["scripts/gateway-smoke.sh"],"manual":true}]`
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"tests", value}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "tests:\n" +
		"  - name: gateway\n" +
		"    command: [scripts/gateway-smoke.sh]\n" +
		"    manual: true\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if len(m.Tests) != 1 || !m.Tests[0].Manual {
		t.Errorf("reads back: %+v", m.Tests)
	}
	automatic := slices.Clone(m.Tests)
	automatic[0].Manual = false
	if sameTiers(m.Tests, automatic) {
		t.Error("tiers that differ in manual are the same")
	}
}

// S-0610: a tier's parallel is written, read back, and makes two tiers
// differ.
func TestEditSettingsWritesParallel(t *testing.T) {
	value := `[{"name":"integration","command":["scripts/integration.sh"],"all_only":true,"parallel":true}]`
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"tests", value}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "tests:\n" +
		"  - name: integration\n" +
		"    command: [scripts/integration.sh]\n" +
		"    all_only: true\n" +
		"    parallel: true\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if len(m.Tests) != 1 || !m.Tests[0].Parallel {
		t.Errorf("reads back: %+v", m.Tests)
	}
	serial := slices.Clone(m.Tests)
	serial[0].Parallel = false
	if sameTiers(m.Tests, serial) {
		t.Error("tiers that differ in parallel are the same")
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
		// S-0547: when only beside queue.
		{"when without queue", `[{"name":"unit","command":["x"],"paths":["**"],"when":["flai/**"]}]`, "tests[0].when", `(tier "unit") is set on a tier without queue: true`},
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

// S-0547: a queue tier's when is written after queue, read back, and makes
// two tiers differ.
func TestEditSettingsWritesWhen(t *testing.T) {
	value := `[{"name":"integration","command":["scripts/integration.sh"],"all_only":true,"queue":true,"when":["flai/internal/serve/**","!flai/internal/serve/testdata/**"]}]`
	got, err := EditSettings([]byte(settingsFixture), []Assignment{{"tests", value}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := settingsFixture + "tests:\n" +
		"  - name: integration\n" +
		"    command: [scripts/integration.sh]\n" +
		"    all_only: true\n" +
		"    queue: true\n" +
		"    when: [\"flai/internal/serve/**\", \"!flai/internal/serve/testdata/**\"]\n"
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	m := loaded(t, got)
	if len(m.Tests) != 1 || !slices.Equal(m.Tests[0].When, []string{"flai/internal/serve/**", "!flai/internal/serve/testdata/**"}) {
		t.Errorf("reads back: %+v", m.Tests)
	}
	again, err := EditSettings(got, []Assignment{{"tests", value}}, nil)
	if err != nil || string(again) != want {
		t.Errorf("setting the same tiers again: %v\n%s", err, again)
	}
	unwhen := slices.Clone(m.Tests)
	unwhen[0].When = []string{"flai/internal/serve/**"}
	if sameTiers(m.Tests, unwhen) {
		t.Error("tiers that differ in when are the same")
	}
}

// klusterFixture is a manifest with a kluster: block written by hand, with
// comments, and keys around it.
const klusterFixture = `version: 1
name: demo   # the project's name
key: d
layout:
  design: design
  docs: docs
  wip: wip
kluster:                          # the project's cluster
  config: integration/kluster.yaml   # its config
  timeout: 45m
claims:
  shared:
    - design/adrs   # shared
`

// S-0586: each key of the kluster: block is written, a command as a flow
// list, its comment kept, and every other byte of the file as it was; []
// unsets a command.
func TestEditSettingsWritesTheKlusterBlock(t *testing.T) {
	out, err := EditSettings([]byte(klusterFixture), []Assignment{
		{"kluster.deploy", `["scripts/kluster-deploy.sh", "--wait", "two words"]`},
		{"kluster.build", `["scripts/kluster-build.sh"]`},
		{"kluster.teardown", `["scripts/kluster-teardown.sh"]`},
		{"kluster.namespace", "shop"},
		{"kluster.timeout", "1h"},
		{"kluster.config", "kluster.yaml"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.NewReplacer(
		"  config: integration/kluster.yaml   # its config\n  timeout: 45m\n",
		"  config: kluster.yaml               # its config\n  timeout: 1h\n  deploy: [scripts/kluster-deploy.sh, \"--wait\", \"two words\"]\n  build: [scripts/kluster-build.sh]\n  teardown: [scripts/kluster-teardown.sh]\n  namespace: shop\n",
	).Replace(klusterFixture)
	if got, want := strings.Split(string(out), "\n"), strings.Split(want, "\n"); !slices.Equal(got, want) {
		for i := range max(len(got), len(want)) {
			if i >= len(got) || i >= len(want) || got[i] != want[i] {
				t.Fatalf("line %d differs:\ngot  %q\nwant %q\nwhole:\n%s", i+1, got[min(i, len(got)-1)], want[min(i, len(want)-1)], out)
			}
		}
	}
	path := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if k := m.Kluster; k == nil || !slices.Equal(k.Deploy, []string{"scripts/kluster-deploy.sh", "--wait", "two words"}) || k.Namespace != "shop" || k.Timeout != "1h" || k.Config != "kluster.yaml" {
		t.Errorf("reads back as %+v", m.Kluster)
	}

	unset, err := EditSettings(out, []Assignment{{"kluster.deploy", "[]"}}, []string{"kluster.build"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(unset), "deploy:") || strings.Contains(string(unset), "build:") || !strings.Contains(string(unset), "teardown: [scripts/kluster-teardown.sh]") {
		t.Errorf("[] and --unset remove the commands:\n%s", unset)
	}
}

// S-0595, ADR-0263: kluster.components is written as a block list, one
// component an item with its name and its build, reads back as given, and []
// removes it; a component that is not of its kind, or that the block's
// validation faults, is refused naming its field.
func TestEditSettingsWritesKlusterComponents(t *testing.T) {
	out, err := EditSettings([]byte(klusterFixture), []Assignment{
		{"kluster.components", `[{"name":"api","build":["scripts/build.sh","api"]},{"name":"web","build":["scripts/build.sh","two words"]}]`},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Replace(klusterFixture, "  timeout: 45m\n",
		"  timeout: 45m\n  components:\n    - name: api\n      build: [scripts/build.sh, api]\n    - name: web\n      build: [scripts/build.sh, \"two words\"]\n", 1)
	if string(out) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", out, want)
	}
	path := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	wantList := []KlusterComponent{{Name: "api", Build: []string{"scripts/build.sh", "api"}}, {Name: "web", Build: []string{"scripts/build.sh", "two words"}}}
	if m.Kluster == nil || !sameComponentList(m.Kluster.Components, wantList) {
		t.Errorf("reads back as %+v", m.Kluster)
	}

	again, err := EditSettings(out, []Assignment{{"kluster.components", `[{"name":"api","build":["scripts/build.sh"]}]`}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Replace(klusterFixture, "  timeout: 45m\n", "  timeout: 45m\n  components:\n    - name: api\n      build: [scripts/build.sh]\n", 1); string(again) != want {
		t.Errorf("a second write replaces the list:\n%s", again)
	}

	unset, err := EditSettings(out, []Assignment{{"kluster.components", "[]"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(unset) != klusterFixture {
		t.Errorf("[] removes the list:\n%s", unset)
	}

	for _, c := range []struct {
		name, value, field, want string
	}{
		{"not a list", `{"name":"api"}`, "kluster.components", "is not a list of components"},
		{"unknown key", `[{"name":"api","build":["b.sh"],"image":"x"}]`, "kluster.components", "is not a list of components"},
		{"build not a list", `[{"name":"api","build":"b.sh"}]`, "kluster.components", "is not a list of components"},
		{"name not a label", `[{"name":"API","build":["b.sh"]}]`, "kluster.components[1].name", "is not a DNS label"},
		{"no build", `[{"name":"api"}]`, "kluster.components[1].build", "is required"},
		{"repeated name", `[{"name":"api","build":["b.sh"]},{"name":"api","build":["c.sh"]}]`, "kluster.components[2].name", "is the name of kluster.components[1] too"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := EditSettings([]byte(klusterFixture), []Assignment{{"kluster.components", c.value}}, nil)
			var r *RefusedError
			if !errors.As(err, &r) {
				t.Fatalf("err = %v, want a refusal", err)
			}
			if !slices.ContainsFunc(r.Problems, func(p Problem) bool { return p.Field == c.field && strings.Contains(p.Reason, c.want) }) {
				t.Errorf("problems = %+v, want %s with %q", r.Problems, c.field, c.want)
			}
		})
	}
	withBuild := strings.Replace(klusterFixture, "  timeout: 45m\n", "  timeout: 45m\n  build: [scripts/build.sh]\n", 1)
	_, err = EditSettings([]byte(withBuild), []Assignment{{"kluster.components", `[{"name":"api","build":["b.sh"]}]`}}, nil)
	var r *RefusedError
	both := func(p Problem) bool {
		return p.Field == "kluster.components" && strings.Contains(p.Reason, "kluster.build are both set")
	}
	if !errors.As(err, &r) || !slices.ContainsFunc(r.Problems, both) {
		t.Errorf("components beside build: %v", err)
	}
	if _, err := EditSettings([]byte(withBuild), []Assignment{{"kluster.components", `[{"name":"api","build":["b.sh"]}]`}}, []string{"kluster.build"}); err != nil {
		t.Errorf("components with build unset in the same edit: %v", err)
	}
}

// S-0586: setting kluster.config on a manifest with no block adds one at the
// end, the rest of the file byte for byte as it was.
func TestEditSettingsAddsTheKlusterBlock(t *testing.T) {
	out, err := EditSettings([]byte(settingsFixture), []Assignment{{"kluster.config", "kluster.yaml"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := settingsFixture + "kluster:\n  config: kluster.yaml\n"; string(out) != want {
		t.Errorf("got:\n%s\nwant:\n%s", out, want)
	}
}

// S-0586: a kluster value of the wrong kind, or one the block's validation
// faults, is refused naming its field, and nothing is written.
func TestEditSettingsRefusesBadKlusterValues(t *testing.T) {
	for _, c := range []struct {
		name  string
		data  string
		set   []Assignment
		field string
		want  string
	}{
		{"command not a list", klusterFixture, []Assignment{{"kluster.deploy", "scripts/kluster-deploy.sh"}}, "kluster.deploy", "is not an argument list"},
		{"command not text", klusterFixture, []Assignment{{"kluster.deploy", `[{"a": 1}]`}}, "kluster.deploy", "is not an argument list"},
		{"command with no program", klusterFixture, []Assignment{{"kluster.build", `["", "--all"]`}}, "kluster.build", "has no program first"},
		{"namespace not a label", klusterFixture, []Assignment{{"kluster.namespace", "Auth_Star"}}, "kluster.namespace", "is not a DNS label"},
		{"namespace story form", klusterFixture, []Assignment{{"kluster.namespace", "s-0042"}}, "kluster.namespace", "has a story namespace's form"},
		{"timeout", klusterFixture, []Assignment{{"kluster.timeout", "soon"}}, "kluster.timeout", "is not a duration"},
		{"config outside", klusterFixture, []Assignment{{"kluster.config", "../kluster.yaml"}}, "kluster.config", "leaves the repository"},
		{"deploy with no block", settingsFixture, []Assignment{{"kluster.deploy", `["scripts/kluster-deploy.sh"]`}}, "kluster.config", "is required"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := EditSettings([]byte(c.data), c.set, nil)
			var r *RefusedError
			if !errors.As(err, &r) {
				t.Fatalf("err = %v, want a refusal", err)
			}
			if !slices.ContainsFunc(r.Problems, func(p Problem) bool { return p.Field == c.field && strings.Contains(p.Reason, c.want) }) {
				t.Errorf("problems = %+v, want %s with %q", r.Problems, c.field, c.want)
			}
		})
	}
}
