package formulas

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

const releaseFile = `---
title: Release {{version}}
updated: 2026-11-01
status: active
name: release
description: Cuts {{version}} and announces it.
nature: feature
inputs:
  - name: version
    description: The version to release.
  - name: channel
    description: Where to announce it.
    default: stable
steps:
  - id: check
    title: Check the tree
    run:
      script: scripts/check.sh --strict
  - id: notes
    title: Write the notes for {{ version }}
    run:
      role: reviewer
      prompt: Draft the release notes for {{version}} on {{channel}}.
  - id: tag
    title: Tag {{version}}
    run:
      flai: release publish {{version}}
    after: [check, notes]
  - id: announce
    title: Announce
    run:
      flai: release announce --channel {{channel}}
    after: [tag]
  - id: lint
    title: Lint
    run:
      script: scripts/lint-md.sh
---

# Release

Cuts a release.
`

const releasePath = "design/workflows/formulas/release.md"

// names are the role and script names the fixtures' steps name.
var names = Names{
	Roles:  []string{"planner", "analyzer", "reviewer"},
	Script: func(p string) bool { return p == "scripts/check.sh" || p == "scripts/lint-md.sh" },
}

func parse(t *testing.T, rel, content string) Formula {
	t.Helper()
	f, p := Parse(rel, content)
	if p != nil {
		t.Fatalf("Parse(%s): %v", rel, p)
	}
	return f
}

func TestValidFormulaParsesWithEveryKind(t *testing.T) {
	f := parse(t, releasePath, releaseFile)
	if f.Name != "release" || f.Title != "Release {{version}}" || f.Nature != "feature" || f.Description != "Cuts {{version}} and announces it." || f.Path != releasePath {
		t.Errorf("front matter = %+v", f)
	}
	if !strings.Contains(f.Body, "Cuts a release.") {
		t.Errorf("body = %q", f.Body)
	}
	if len(f.Inputs) != 2 || !f.Inputs[0].Required() || f.Inputs[1].Required() || *f.Inputs[1].Default != "stable" || f.Inputs[0].Line != 9 {
		t.Errorf("inputs = %+v", f.Inputs)
	}
	kinds := []string{}
	for _, s := range f.Steps {
		kinds = append(kinds, s.Kind())
	}
	if want := []string{KindScript, KindRole, KindFlai, KindFlai, KindScript}; !reflect.DeepEqual(kinds, want) {
		t.Errorf("kinds = %v, want %v", kinds, want)
	}
	notes := f.Steps[1]
	if notes.Run.Role != "reviewer" || notes.Run.Prompt != "Draft the release notes for {{version}} on {{channel}}." || notes.Line != 19 || notes.Lines["run"] != 21 {
		t.Errorf("notes step = %+v", notes)
	}
	if got := f.Steps[0].Run.ScriptPath(); got != "scripts/check.sh" {
		t.Errorf("ScriptPath = %q", got)
	}
	if got := f.Steps[0].After; got == nil || len(got) != 0 {
		t.Errorf("an after not given = %#v, want empty", got)
	}
	if problems := Validate([]Formula{f}, names); len(problems) != 0 {
		t.Errorf("Validate on a valid formula: %v", problems)
	}
	want := [][]string{{"check", "notes", "lint"}, {"tag"}, {"announce"}}
	if got := Layers(f); !reflect.DeepEqual(got, want) {
		t.Errorf("Layers = %v, want %v", got, want)
	}
}

func TestDefaultsFillIn(t *testing.T) {
	f := parse(t, "design/workflows/formulas/x.md", "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nfuture_field: ignored\n---\n")
	if f.Nature != DefaultNature || f.Inputs == nil || f.Steps == nil {
		t.Errorf("defaults = %+v", f)
	}
	if f.Lines["future_field"] != 7 {
		t.Errorf("lines = %v", f.Lines)
	}
}

func TestParseFailuresAreProblems(t *testing.T) {
	for _, tc := range []struct {
		name, content, says string
		line                int
	}{
		{"no front matter", "# x\n", "front matter", 1},
		{"bad yaml", "---\ntitle: [x\n---\n", "does not parse", 2},
		{"not a mapping", "---\n- a\n---\n", "not a mapping", 2},
		{"run a string", "---\ntitle: X\nsteps:\n  - id: a\n    title: A\n    run: flai check\n---\n", "run is a mapping", 6},
		{"step a string", "---\ntitle: X\nsteps:\n  - a\n---\n", "a steps entry is a mapping", 4},
		{"input a string", "---\ntitle: X\ninputs:\n  - a\n---\n", "an inputs entry is a mapping", 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, p := Parse("design/workflows/formulas/x.md", tc.content)
			if p == nil || p.Code != CodeFile || p.Line != tc.line || !strings.Contains(p.Message, tc.says) {
				t.Errorf("Parse = %+v, want line %d saying %q", p, tc.line, tc.says)
			}
		})
	}
}

// problems are the messages and lines Validate reports for the formula file
// content at design/workflows/formulas/x.md.
func problems(t *testing.T, content string) []Problem {
	t.Helper()
	return Validate([]Formula{parse(t, "design/workflows/formulas/x.md", content)}, names)
}

func TestValidateReportsEachProblem(t *testing.T) {
	// head is the fields every case but its own keeps, on lines 2 to 6.
	const head = "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\n"
	const oneStep = "steps: [{id: a, title: A, run: {flai: check}}]\n---\n"
	step := func(lines ...string) string {
		return head + "steps:\n  - id: a\n    title: A\n" + strings.Join(lines, "\n") + "\n---\n"
	}
	without := func(field string) string {
		return strings.Replace(head, field+"\n", "", 1) + oneStep
	}
	for _, tc := range []struct {
		name, content, says string
		line                int
	}{
		{"title missing", without("title: X"), "needs title", 1},
		{"updated missing", without("updated: 2026-11-01"), "needs updated", 1},
		{"status missing", without("status: active"), "needs status", 1},
		{"name missing", without("name: x"), "needs name", 1},
		{"description missing", without("description: d"), "needs description", 1},
		{"steps missing", head + "---\n", "needs steps", 1},
		{"steps empty", head + "steps: []\n---\n", "needs steps", 7},
		{"updated not a date", strings.Replace(head, "2026-11-01", "soon", 1) + oneStep, `updated is "soon", not a date`, 3},
		{"status invalid", strings.Replace(head, "active", "live", 1) + oneStep, `status is "live"`, 4},
		{"name not the stem", strings.Replace(head, "name: x", "name: y", 1) + oneStep, "file is x.md", 5},
		{"name not kebab", strings.Replace(head, "name: x", "name: X_y", 1) + oneStep, "not kebab-case", 5},
		{"nature invalid", head + "nature: chore\n" + oneStep, `nature is "chore"`, 7},
		{"run missing", step(), "needs run", 8},
		{"run none", step("    run:", "      prompt: p"), "none of flai, script, or role", 10},
		{"run empty flai", step("    run:", `      flai: ""`), "none of flai, script, or role", 10},
		{"run two kinds", step("    run:", "      flai: check", "      script: scripts/check.sh"), "gives flai and script", 10},
		{"script outside scripts", step("    run:", "      script: bin/x.sh"), "not under scripts/", 10},
		{"script escaping scripts", step("    run:", "      script: scripts/../x.sh"), "not under scripts/", 10},
		{"script absolute", step("    run:", "      script: /scripts/x.sh"), "not under scripts/", 10},
		{"script missing", step("    run:", "      script: scripts/gone.sh --all"), "scripts/gone.sh does not exist", 10},
		{"role undefined", step("    run:", "      role: deployer", "      prompt: p"), "role deployer is no role file's", 10},
		{"role without prompt", step("    run:", "      role: reviewer"), "no prompt", 10},
		{"after no step", step("    run: {flai: check}", "    after: [b]"), `after "b", which is no step's id`, 11},
		{"in neither place", step("    run: {flai: check}", "    in: elsewhere"), `step a's in is "elsewhere", which is none of worktree, main`, 11},
		{"stop on a role step", step("    run: {role: reviewer, prompt: p}", "    stop: [2]"), "step a names stop, but a role step ends worked, asked, or failed with no exit code; give stop to a flai or script step only", 11},
		{"stop exit 0", step("    run: {flai: check}", "    stop: [0]"), "step a's stop lists exit 0; a stop exit is from 1 to 255, as 0 is the step's success", 11},
		{"stop exit over 255", step("    run: {script: scripts/check.sh}", "    stop: [256]"), "step a's stop lists exit 256", 11},
		{"stop exit negative", step("    run: {flai: check}", "    stop: [-1]"), "step a's stop lists exit -1", 11},
		{"stop exit twice", step("    run: {flai: check}", "    stop: [2, 3, 2]"), "step a's stop lists exit 2 twice; keep one", 11},
		{"stop a number", step("    run: {flai: check}", "    stop: 2"), "step a's stop is not a list of exit codes; give the exits that stop the run cleanly, such as stop: [2]", 11},
		{"stop of words", step("    run: {flai: check}", "    stop: [two]"), "step a's stop is not a list of exit codes", 11},
		{"stop of a huge number", step("    run: {flai: check}", "    stop: [99999999999999999999]"), "step a's stop is not a list of exit codes", 11},
		{"stop a mapping", step("    run: {flai: check}", "    stop: {exit: 2}"), "step a's stop is not a list of exit codes", 11},
		{"step id missing", head + "steps:\n  - title: A\n    run: {flai: check}\n---\n", "step needs id", 8},
		{"step id not kebab", head + "steps:\n  - id: A_1\n    title: A\n    run: {flai: check}\n---\n", "not kebab-case", 8},
		{"step title missing", head + "steps:\n  - id: a\n    run: {flai: check}\n---\n", "step a needs title", 8},
		{"title uses undeclared", strings.Replace(head, "title: X", "title: X {{who}}", 1) + oneStep, "title uses {{who}}", 2},
		{"description uses undeclared", strings.Replace(head, "description: d", "description: d {{who}}", 1) + oneStep, "description uses {{who}}", 6},
		{"step title uses undeclared", head + "steps:\n  - id: a\n    title: A {{ who }}\n    run: {flai: check}\n---\n", "step a's title uses {{who}}", 9},
		{"flai uses undeclared", step(`    run: {flai: "check {{who}}"}`), "step a's run uses {{who}}", 10},
		{"script uses undeclared", step(`    run: {script: "scripts/check.sh {{who}}"}`), "step a's run uses {{who}}", 10},
		{"prompt uses undeclared", step(`    run: {role: reviewer, prompt: "plan {{who}}"}`), "step a's run uses {{who}}", 10},
		{"input name invalid", head + "inputs:\n  - {name: Who, description: d}\n" + oneStep, `input name "Who"`, 8},
		{"input name missing", head + "inputs:\n  - description: d\n" + oneStep, "needs name", 8},
		{"input description missing", head + "inputs:\n  - name: who\n" + oneStep, "input who needs description", 8},
		{"schedule not parsed", head + "schedule: every day\n" + oneStep, `schedule "every day" has 2 fields, not five`, 7},
		{"schedule out of range", head + "schedule: \"0 25 * * *\"\n" + oneStep, `schedule "0 25 * * *"`, 7},
		{"schedule with inputs", head + "schedule: daily\ninputs:\n  - {name: who, description: d}\n" + oneStep, "schedule starts the formula with no inputs, as nobody is there to give them, but who has no default; give it a default, or remove the schedule", 7},
		{"schedule with an input of no default among defaults", head + "schedule: daily\ninputs:\n  - {name: who, description: d, default: me}\n  - {name: what, description: d}\n  - {name: when, description: d, default: \"\"}\n" + oneStep, "but what has no default;", 7},
		{"schedule with inputs of no default", head + "schedule: daily\ninputs:\n  - {name: who, description: d}\n  - {name: what, description: d, default: it}\n  - {name: when, description: d}\n" + oneStep, "but who, when have no default; give each a default", 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := problems(t, tc.content)
			if len(got) != 1 || got[0].Code != CodeFile || got[0].Line != tc.line || !strings.Contains(got[0].Message, tc.says) {
				t.Errorf("Validate = %v, want one on line %d saying %q", got, tc.line, tc.says)
			}
		})
	}
}

// S-0423: a step's in says where it runs: main is the project's main
// checkout, and worktree or no in at all is the instance's worktree; the line
// of its in key is kept, and any other value is refused on that line, naming
// the two.
func TestStepIn(t *testing.T) {
	const head = "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nsteps:\n  - id: a\n    title: A\n    run: {flai: check}\n"
	for _, tc := range []struct {
		name, line, in string
		main           bool
	}{
		{"main", "    in: main\n", InMain, true},
		{"worktree", "    in: worktree\n", InWorktree, false},
		{"absent", "", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			content := head + tc.line + "---\n"
			s := parse(t, "design/workflows/formulas/x.md", content).Steps[0]
			if s.In != tc.in || s.InMain() != tc.main {
				t.Errorf("in = %q, InMain = %v, want %q, %v", s.In, s.InMain(), tc.in, tc.main)
			}
			if tc.line != "" && s.Lines["in"] != 11 {
				t.Errorf("the in key's line = %d, want 11", s.Lines["in"])
			}
			if got := problems(t, content); len(got) != 0 {
				t.Errorf("Validate: %v", got)
			}
		})
	}
	got := problems(t, head+"    in: elsewhere\n---\n")
	if len(got) != 1 || got[0].Code != CodeFile || got[0].Line != 11 || !strings.Contains(got[0].Message, `step a's in is "elsewhere", which is none of worktree, main; give worktree`) {
		t.Errorf("in: elsewhere = %v, want one problem on line 11 naming worktree and main", got)
	}
	filled := parse(t, "design/workflows/formulas/x.md", head+"    in: main\n---\n").Filled(nil)
	if !filled.Steps[0].InMain() {
		t.Errorf("Filled dropped the step's in: %+v", filled.Steps[0])
	}
}

// S-0621 (ADR-0271 §1): a flai or script step's stop is the list of exits
// that stop the instance cleanly, in either YAML list form, with the line of
// its key kept; Stops says whether an exit is one of them, 0 never; Filled
// keeps it; and a step with no stop has none.
func TestStepStop(t *testing.T) {
	const head = "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nsteps:\n  - id: a\n    title: A\n    run: {flai: check}\n"
	for _, tc := range []struct {
		name, lines string
		stop        []int
	}{
		{"flow", "    stop: [2]\n", []int{2}},
		{"block", "    stop:\n      - 2\n      - 255\n", []int{2, 255}},
		{"absent", "", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			content := head + tc.lines + "---\n"
			s := parse(t, "design/workflows/formulas/x.md", content).Steps[0]
			if !slices.Equal(s.Stop, tc.stop) {
				t.Errorf("stop = %v, want %v", s.Stop, tc.stop)
			}
			if tc.lines != "" && s.Lines["stop"] != 11 {
				t.Errorf("the stop key's line = %d, want 11", s.Lines["stop"])
			}
			if got := problems(t, content); len(got) != 0 {
				t.Errorf("Validate: %v", got)
			}
		})
	}
	s := parse(t, "design/workflows/formulas/x.md", head+"    stop: [2, 3]\n---\n").Filled(nil).Steps[0]
	if !s.Stops(2) || !s.Stops(3) || s.Stops(0) || s.Stops(1) || !slices.Equal(s.Stop, []int{2, 3}) {
		t.Errorf("Filled's step: Stop %v, Stops(2) %v, Stops(0) %v", s.Stop, s.Stops(2), s.Stops(0))
	}
	if (Step{}).Stops(0) {
		t.Error("a step with no stop stops on exit 0")
	}
}

// S-0422: a formula's schedule is read as planning.schedule is, a cron
// expression in UTC or daily, and none is no schedule.
func TestScheduleParses(t *testing.T) {
	const head = "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\n"
	const oneStep = "steps: [{id: a, title: A, run: {flai: check}}]\n---\n"
	for _, spec := range []string{"daily", "0 6 * * 1-5"} {
		content := head + "schedule: \"" + spec + "\"\n" + oneStep
		if got := problems(t, content); len(got) != 0 {
			t.Errorf("schedule %s: %v", spec, got)
		}
		sched, err := parse(t, "design/workflows/formulas/x.md", content).Scheduled()
		if err != nil || sched == nil || sched.String() != spec {
			t.Errorf("Scheduled() for %s = %v, %v", spec, sched, err)
		}
	}
	if sched, err := parse(t, "design/workflows/formulas/x.md", head+oneStep).Scheduled(); sched != nil || err != nil {
		t.Errorf("no schedule: %v, %v", sched, err)
	}
}

// S-0424: a scheduled formula may declare inputs when each carries a default,
// an empty one included, since the run its schedule starts takes them.
func TestAScheduledFormulaMayDeclareInputsWithDefaults(t *testing.T) {
	content := "---\ntitle: X {{focus}}\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nschedule: daily\n" +
		"inputs:\n  - {name: focus, description: d, default: all}\n  - {name: note, description: d, default: \"\"}\n" +
		"steps: [{id: a, title: A, run: {flai: \"analyze --focus {{focus}} {{note}}\"}}]\n---\n"
	if got := problems(t, content); len(got) != 0 {
		t.Errorf("Validate = %v, want none", got)
	}
}

func TestValidateReportsDuplicates(t *testing.T) {
	got := problems(t, "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\ninputs:\n  - {name: who, description: d}\n  - {name: who, description: d}\nsteps:\n  - id: a\n    title: A\n    run: {flai: check}\n  - id: a\n    title: B\n    run: {flai: check}\n---\n")
	if len(got) != 2 || got[0].Line != 9 || !strings.Contains(got[0].Message, "input who is declared twice") || got[1].Line != 14 || !strings.Contains(got[1].Message, "step id a is also an earlier step's") {
		t.Errorf("Validate = %v", got)
	}
}

func TestValidateReportsACycleNamingItsSteps(t *testing.T) {
	got := problems(t, `---
title: X
updated: 2026-11-01
status: active
name: x
description: d
steps:
  - id: a
    title: A
    run: {flai: check}
    after: [c]
  - id: b
    title: B
    run: {flai: check}
    after: [a]
  - id: c
    title: C
    run: {flai: check}
    after: [b]
  - id: d
    title: D
    run: {flai: check}
    after: [d]
  - id: e
    title: E
    run: {flai: check}
---
`)
	if len(got) != 2 {
		t.Fatalf("Validate = %v, want two cycles", got)
	}
	if got[0].Line != 11 || !strings.Contains(got[0].Message, "steps a, b, c are after each other in a cycle: a after c after b after a") {
		t.Errorf("first cycle = %v", got[0])
	}
	if got[1].Line != 23 || !strings.Contains(got[1].Message, "steps d are after each other in a cycle: d after d") {
		t.Errorf("second cycle = %v", got[1])
	}
	f := parse(t, "design/workflows/formulas/x.md", "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nsteps:\n  - {id: a, title: A, run: {flai: x}, after: [b]}\n  - {id: b, title: B, run: {flai: x}, after: [a]}\n  - {id: c, title: C, run: {flai: x}, after: [a]}\n  - {id: e, title: E, run: {flai: x}}\n---\n")
	if want := [][]string{{"e"}}; !reflect.DeepEqual(Layers(f), want) {
		t.Errorf("Layers past a cycle = %v, want %v", Layers(f), want)
	}
}

func TestValidateSortsByPathAndLine(t *testing.T) {
	a := parse(t, "design/workflows/formulas/b.md", "---\nname: b\n---\n")
	b := parse(t, "design/workflows/formulas/a.md", "---\nname: a\ntitle: A\n---\n")
	got := Validate([]Formula{a, b}, names)
	if len(got) != 9 || got[0].Path != "design/workflows/formulas/a.md" || got[len(got)-1].Path != "design/workflows/formulas/b.md" {
		t.Errorf("Validate = %v", got)
	}
}

func TestNoRolesRefusesEveryRole(t *testing.T) {
	f := parse(t, "design/workflows/formulas/x.md", "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nsteps: [{id: a, title: A, run: {role: reviewer, prompt: p}}]\n---\n")
	got := Validate([]Formula{f}, Names{})
	if len(got) != 1 || !strings.Contains(got[0].Message, "the project has none") {
		t.Errorf("Validate with no roles = %v", got)
	}
}

// S-0424: a role step that names the planner, the orchestrator, or the
// curator is reported, though the project has the role's file, since a role
// step cannot start it; one that names the analyzer is not, as it is started
// for a role step through its own start.
func TestARoleStepNamingABuiltinRoleOtherThanTheAnalyzerIsReported(t *testing.T) {
	all := Names{Roles: []string{"planner", "orchestrator", "analyzer", "curator"}}
	for _, role := range []string{"planner", "orchestrator", "curator"} {
		f := parse(t, "design/workflows/formulas/x.md", "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nsteps: [{id: a, title: A, run: {role: "+role+", prompt: p}}]\n---\n")
		want := "step a's role " + role + " is a built-in role a role step cannot start; only the analyzer is started for a role step, through flai analyze's start"
		if got := Validate([]Formula{f}, all); len(got) != 1 || got[0].Code != CodeFile || got[0].Line != 7 || got[0].Message != want {
			t.Errorf("Validate on a step naming the %s = %v, want one on line 7 saying %q", role, got, want)
		}
	}
	f := parse(t, "design/workflows/formulas/x.md", "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nsteps: [{id: a, title: A, run: {role: analyzer, prompt: p}}]\n---\n")
	if got := Validate([]Formula{f}, all); len(got) != 0 {
		t.Errorf("Validate on a step naming the analyzer = %v, want none", got)
	}
}

func TestNilScriptSkipsTheExistenceCheck(t *testing.T) {
	f := parse(t, "design/workflows/formulas/x.md", "---\ntitle: X\nupdated: 2026-11-01\nstatus: active\nname: x\ndescription: d\nsteps: [{id: a, title: A, run: {script: ./scripts/gone.sh}}]\n---\n")
	if got := Validate([]Formula{f}, Names{}); len(got) != 0 {
		t.Errorf("Validate with no Script = %v", got)
	}
}

func TestLayersOrderLikeAPlan(t *testing.T) {
	f := parse(t, "design/workflows/formulas/x.md", `---
title: X
name: x
description: d
steps:
  - {id: d, title: D, run: {flai: x}, after: [b, c]}
  - {id: c, title: C, run: {flai: x}, after: [a]}
  - {id: b, title: B, run: {flai: x}}
  - {id: a, title: A, run: {flai: x}, after: [gone]}
  - {id: e, title: E, run: {flai: x}, after: [a]}
---
`)
	want := [][]string{{"b", "a"}, {"c", "e"}, {"d"}}
	if got := Layers(f); !reflect.DeepEqual(got, want) {
		t.Errorf("Layers = %v, want %v", got, want)
	}
	if got := Layers(Formula{}); len(got) != 0 {
		t.Errorf("Layers of no steps = %v", got)
	}
}

func TestUsesAndFill(t *testing.T) {
	for _, tc := range []struct {
		text, filled string
		uses         []string
	}{
		{"plain", "plain", nil},
		{"{{a}} and {{ b }} and {{a}}", "1 and 2 and 1", []string{"a", "b"}},
		{"{{  a}}{{c}}", "1{{c}}", []string{"a", "c"}},
		{"{{ }} {a} {{a b}}", "{{ }} {a} {{a b}}", nil},
	} {
		if got := Uses(tc.text); !reflect.DeepEqual(got, tc.uses) {
			t.Errorf("Uses(%q) = %v, want %v", tc.text, got, tc.uses)
		}
		if got := Fill(tc.text, map[string]string{"a": "1", "b": "2"}); got != tc.filled {
			t.Errorf("Fill(%q) = %q, want %q", tc.text, got, tc.filled)
		}
	}
}

func TestFilledFillsTitlesAndRuns(t *testing.T) {
	f := parse(t, releasePath, releaseFile)
	got := f.Filled(map[string]string{"version": "1.2.0", "channel": "beta"})
	if got.Title != "Release 1.2.0" || got.Description != "Cuts 1.2.0 and announces it." || got.Steps[1].Title != "Write the notes for 1.2.0" || got.Steps[1].Run.Prompt != "Draft the release notes for 1.2.0 on beta." || got.Steps[2].Run.Flai != "release publish 1.2.0" || got.Steps[3].Run.Flai != "release announce --channel beta" {
		t.Errorf("Filled = %+v", got)
	}
	if f.Title != "Release {{version}}" || f.Steps[1].Title != "Write the notes for {{ version }}" {
		t.Errorf("Filled changed the formula: %+v", f)
	}
}

func TestResolve(t *testing.T) {
	f := parse(t, releasePath, releaseFile)
	got, err := f.Resolve(map[string]string{"version": "1.2.0"})
	if err != nil || !reflect.DeepEqual(got, map[string]string{"version": "1.2.0", "channel": "stable"}) {
		t.Errorf("Resolve with the default = %v, %v", got, err)
	}
	got, err = f.Resolve(map[string]string{"version": "1.2.0", "channel": "beta"})
	if err != nil || got["channel"] != "beta" {
		t.Errorf("Resolve over the default = %v, %v", got, err)
	}
	for _, tc := range []struct {
		name  string
		given map[string]string
		says  []string
	}{
		{"missing", map[string]string{}, []string{"formula release needs input version, which has no default", "its inputs are version (required), channel", "--input <name>=<value>"}},
		{"undeclared", map[string]string{"version": "1", "zone": "x", "area": "y"}, []string{"formula release declares no input area, zone;", "its inputs are version (required), channel"}},
		{"both", map[string]string{"zone": "x"}, []string{"declares no input zone, and formula release needs input version"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := f.Resolve(tc.given)
			if err == nil || got != nil {
				t.Fatalf("Resolve = %v, %v, want an error", got, err)
			}
			for _, s := range tc.says {
				if !strings.Contains(err.Error(), s) {
					t.Errorf("error %q does not say %q", err, s)
				}
			}
		})
	}
	none := Formula{Name: "bare"}
	if _, err := none.Resolve(map[string]string{"x": "1"}); err == nil || !strings.Contains(err.Error(), "its inputs are none") {
		t.Errorf("Resolve on no inputs = %v", err)
	}
}

func TestFind(t *testing.T) {
	fs := []Formula{{Name: "release"}, {Name: "hotfix"}}
	if f, err := Find(fs, "hotfix"); err != nil || f.Name != "hotfix" {
		t.Errorf("Find = %v, %v", f, err)
	}
	if _, err := Find(fs, "deploy"); err == nil || !strings.Contains(err.Error(), "the formulas are release, hotfix") {
		t.Errorf("Find of no formula = %v", err)
	}
	if _, err := Find(nil, "deploy"); err == nil || !strings.Contains(err.Error(), "workflows/formulas/deploy.md") {
		t.Errorf("Find among none = %v", err)
	}
}

func TestLoadReadsTheFolder(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "design", "workflows", "formulas")
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{
		"README.md":  "# Formulas\n",
		"release.md": releaseFile,
		"broken.md":  "---\ntitle: [x\n---\n",
		"notes.txt":  "not a formula",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	fs, ps, err := Load(root, "design")
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 1 || fs[0].Path != releasePath {
		t.Errorf("Load formulas = %+v", fs)
	}
	if len(ps) != 1 || ps[0].Path != "design/workflows/formulas/broken.md" || ps[0].Line != 2 {
		t.Errorf("Load problems = %v", ps)
	}
	fs, ps, err = Load(t.TempDir(), "design")
	if fs != nil || ps != nil || err != nil {
		t.Errorf("Load of no folder = %v, %v, %v", fs, ps, err)
	}
}

func TestScriptOnDisk(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "scripts", "dir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "check.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	exists := ScriptOnDisk(root)
	if !exists("scripts/check.sh") || exists("scripts/gone.sh") || exists("scripts/dir") {
		t.Errorf("ScriptOnDisk: check.sh %v, gone.sh %v, dir %v", exists("scripts/check.sh"), exists("scripts/gone.sh"), exists("scripts/dir"))
	}
}

// adrExample is ADR-0238's example formula, release-check.md.
const adrExample = `---
title: "Release check for {{version}}"
updated: 2026-11-01
status: active
name: release-check
description: "Check that the work on main is ready to publish as {{version}}: the tests pass, flai's release evaluation agrees, and the analyzer reports the risk."
nature: improvement
inputs:
  - name: version
    description: The version to check, such as 1.4.0.
  - name: area
    description: What the analyzer looks at.
    default: the changes since the last release
steps:
  - id: test
    title: The tests pass
    run: {script: scripts/test.sh}
  - id: evaluate
    title: The release evaluation agrees on {{version}}
    run: {flai: release --evaluate}
    after: [test]
  - id: risk
    title: The analyzer reports the risk of {{version}}
    run:
      role: analyzer
      prompt: Report the risk in {{area}} before {{version}} is published.
    after: [test]
---

# Release check

<!-- system-flow:end-of-baseline -->

## Project additions
`

func TestTheADRsExampleIsValid(t *testing.T) {
	f := parse(t, "design/workflows/formulas/release-check.md", adrExample)
	if got := Validate([]Formula{f}, Names{Roles: []string{"analyzer"}}); len(got) != 0 {
		t.Errorf("Validate on ADR-0238's example: %v", got)
	}
	if want := [][]string{{"test"}, {"evaluate", "risk"}}; !reflect.DeepEqual(Layers(f), want) {
		t.Errorf("Layers = %v, want %v", Layers(f), want)
	}
	values, err := f.Resolve(map[string]string{"version": "1.4.0"})
	if err != nil {
		t.Fatal(err)
	}
	if got := f.Filled(values); got.Title != "Release check for 1.4.0" || got.Steps[2].Run.Prompt != "Report the risk in the changes since the last release before 1.4.0 is published." {
		t.Errorf("Filled = %+v", got)
	}
}
