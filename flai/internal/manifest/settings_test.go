package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
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
		{"not an agent", []Assignment{{"analysis.agent", "claude"}}, nil, "analysis.agent", "is not an agent"},
		{"agent with an unknown field", []Assignment{{"analysis.agent", `{"modelname":"x"}`}}, nil, "analysis.agent", "is not an agent"},
		{"text", []Assignment{{"orchestration.release.policy", "theme"}, {"orchestration.release.epic", "epic-7"}}, nil, "orchestration.release.epic", "is not an epic ID"},
		{"empty", []Assignment{{"planning.cycle", " "}}, nil, "planning.cycle", "is empty"},
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
	for _, key := range []string{"orchestration.permissions.publish", "planning.hour_rate", "analysis.schedule", "analysis.agent"} {
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
