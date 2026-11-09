package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// releaseCheck is ADR-0238's example formula with the reviewer role
// reviewerRole defines: a script step, then a flai step and a role step after
// it, an input with a default and one without.
const releaseCheck = `---
title: "Release check for {{version}}"
updated: 2026-11-01
status: active
name: release-check
description: "Check that the work on main is ready to publish as {{version}}."
nature: remediation
inputs:
  - name: version
    description: The version to check, such as 1.4.0.
  - name: area
    description: What the reviewer looks at.
    default: the changes since the last release
steps:
  - id: test
    title: The tests pass
    run: {script: scripts/test.sh --short}
  - id: evaluate
    title: The release evaluation agrees on {{version}}
    run: {flai: release --evaluate}
    after: [test]
  - id: risk
    title: The reviewer reports the risk of {{version}}
    run:
      role: reviewer
      prompt: Report the risk in {{area}} before {{version}} is published.
    after: [test]
---

# Release check

Run before a release is published.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// formulaProject is a project with the reviewer role, scripts/test.sh, and
// the formula files files.
func formulaProject(t *testing.T, files map[string]string) string {
	t.Helper()
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	for dir, list := range map[string]map[string]string{"scripts": {"test.sh": "#!/bin/sh\n"}, "design/workflows/formulas": files} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
		for name, content := range list {
			if err := os.WriteFile(filepath.Join(root, dir, name), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return root
}

// itemFiles are the story and task files under root's kanban.
func itemFiles(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	for _, dir := range []string{"stories", "tasks"} {
		matches, err := filepath.Glob(filepath.Join(root, "wip", "kanban", dir, "*.md"))
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, matches...)
	}
	return out
}

// itemFile is the content of the one item file of id under root's kanban.
func itemFile(t *testing.T, root, id string) string {
	t.Helper()
	matches, _ := filepath.Glob(filepath.Join(root, "wip", "kanban", "*", id+"-*.md"))
	if len(matches) != 1 {
		t.Fatalf("%s: %v", id, matches)
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0421 (ADR-0238 §7): flai formula list gives each formula's name, nature,
// step count, file, and title, and says a file that does not load on stderr;
// --json gives a list.
func TestFormulaList(t *testing.T) {
	t.Parallel()
	root := formulaProject(t, map[string]string{"README.md": "# Formulas\n", "release-check.md": releaseCheck, "broken.md": "no front matter\n"})
	out, errOut, code := runIn(t, root, "formula", "list")
	if code != 0 || out != "release-check  remediation  3 steps  design/workflows/formulas/release-check.md\n  Release check for {{version}}\n" {
		t.Errorf("formula list: %d %q %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "design/workflows/formulas/broken.md:1: formulas.file:") || !strings.Contains(errOut, "(left out)") {
		t.Errorf("the file that does not load is said on stderr: %s", errOut)
	}
	js, _, code := runIn(t, root, "formula", "list", "--json")
	var got []map[string]any
	if err := json.Unmarshal([]byte(js), &got); code != 0 || err != nil || len(got) != 1 {
		t.Fatalf("formula list --json: %d %v %s", code, err, js)
	}
	if got[0]["name"] != "release-check" || got[0]["nature"] != "remediation" || got[0]["steps"] != 3.0 || got[0]["path"] != "design/workflows/formulas/release-check.md" || got[0]["title"] != "Release check for {{version}}" {
		t.Errorf("formula list --json: %v", got[0])
	}
	empty := formulaProject(t, nil)
	if out, _, _ := runIn(t, empty, "formula", "list"); out != "no formulas\n" {
		t.Errorf("no formulas: %q", out)
	}
	if js, _, _ := runIn(t, empty, "formula", "list", "--json"); strings.TrimSpace(js) != "[]" {
		t.Errorf("no formulas, --json: %q", js)
	}
}

// S-0421 (ADR-0238 §7): flai formula show prints the inputs, the steps with
// what each runs and is after, and the layers; --json gives the formula,
// each step's kind, and its layers; an unknown name names the formulas, and
// one whose file does not load says why.
func TestFormulaShow(t *testing.T) {
	t.Parallel()
	root := formulaProject(t, map[string]string{"release-check.md": releaseCheck, "broken.md": "no front matter\n"})
	out, errOut, code := runIn(t, root, "formula", "show", "release-check")
	want := `release-check  remediation  design/workflows/formulas/release-check.md
  title:       Release check for {{version}}
  status:      active
  updated:     2026-11-01
  description: Check that the work on main is ready to publish as {{version}}.
inputs:
  version  required  The version to check, such as 1.4.0.
  area  default "the changes since the last release"  What the reviewer looks at.
steps:
  test: The tests pass
    script: scripts/test.sh --short
  evaluate: The release evaluation agrees on {{version}}
    flai:   flai release --evaluate
    after:  test
  risk: The reviewer reports the risk of {{version}}
    role:   reviewer
    prompt: Report the risk in {{area}} before {{version}} is published.
    after:  test
layers:
  1  test
  2  evaluate, risk
`
	if code != 0 || out != want {
		t.Errorf("formula show: %d %s\n%s", code, errOut, out)
	}
	js, _, code := runIn(t, root, "formula", "show", "release-check", "--json")
	var got struct {
		Name   string     `json:"name"`
		Layers [][]string `json:"layers"`
		Inputs []struct {
			Name    string  `json:"name"`
			Default *string `json:"default"`
		} `json:"inputs"`
		Steps []struct {
			ID   string       `json:"id"`
			Kind string       `json:"kind"`
			Run  formulas.Run `json:"run"`
		} `json:"steps"`
	}
	if err := json.Unmarshal([]byte(js), &got); code != 0 || err != nil {
		t.Fatalf("formula show --json: %d %v %s", code, err, js)
	}
	var kinds []string
	for _, s := range got.Steps {
		kinds = append(kinds, s.Kind)
	}
	if got.Name != "release-check" || !slices.Equal(kinds, []string{"script", "flai", "role"}) || got.Steps[2].Run.Role != "reviewer" ||
		len(got.Layers) != 2 || !slices.Equal(got.Layers[1], []string{"evaluate", "risk"}) || got.Inputs[0].Default != nil || got.Inputs[1].Default == nil {
		t.Errorf("formula show --json: %s", js)
	}
	if _, errOut, code := runIn(t, root, "formula", "show", "nope"); code == 0 || !strings.Contains(errOut, "no formula is named nope; the formulas are release-check") {
		t.Errorf("an unknown formula: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "formula", "show", "broken"); code == 0 || !strings.Contains(errOut, "formula broken does not load: design/workflows/formulas/broken.md:1: formulas.file:") {
		t.Errorf("a formula whose file does not load: %d %s", code, errOut)
	}
}

// S-0421 (ADR-0238 §7, §9): flai formula run makes one story in backlog with
// the formula's nature, formula: <name>, the filled title, the description
// as its goal, a criterion per step, and notes naming the file and the
// inputs, in the epic given; and a task per step with its Work, its Done
// when, and the tasks of its after as its after, so that the tasks' layers
// are the steps'. flai check --strict passes on what it made.
func TestFormulaRun(t *testing.T) {
	t.Parallel()
	root := formulaProject(t, map[string]string{"release-check.md": releaseCheck})
	if _, errOut, code := runIn(t, root, "epic", "new", "Releases"); code != 0 {
		t.Fatal(errOut)
	}
	js, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=1.4.0", "--epic", "E-0001", "--json")
	if code != 0 {
		t.Fatalf("formula run: %d %s %s", code, js, errOut)
	}
	var got formulaRunResult
	if err := json.Unmarshal([]byte(js), &got); err != nil {
		t.Fatalf("formula run --json: %v %s", err, js)
	}
	wantTasks := []formulaTask{{ID: "T-0001", Step: "test"}, {ID: "T-0002", Step: "evaluate"}, {ID: "T-0003", Step: "risk"}}
	if got.Story != "S-0001" || got.Status != "backlog" || got.Formula != "release-check" || !slices.Equal(got.Tasks, wantTasks) ||
		got.Inputs["version"] != "1.4.0" || got.Inputs["area"] != "the changes since the last release" || got.Commit != "" || got.Followed != nil {
		t.Errorf("formula run --json: %+v", got)
	}
	story := itemFile(t, root, "S-0001")
	for _, want := range []string{
		"type: story\nnature: remediation\ntitle: Release check for 1.4.0\nstatus: backlog\nparent: E-0001\n",
		"\nformula: release-check\n",
		"# S-0001 Release check for 1.4.0\n\n## Goal\n\nCheck that the work on main is ready to publish as 1.4.0.\n\n",
		"## Acceptance criteria\n\n- [ ] The tests pass\n- [ ] The release evaluation agrees on 1.4.0\n- [ ] The reviewer reports the risk of 1.4.0\n\n",
		"## Tasks\n- T-0001 The tests pass\n- T-0002 The release evaluation agrees on 1.4.0\n- T-0003 The reviewer reports the risk of 1.4.0\n",
		"## Notes\n\nMade by `flai formula run release-check` from `design/workflows/formulas/release-check.md`, with the inputs `area=the changes since the last release`, `version=1.4.0`.\n",
	} {
		if !strings.Contains(story, want) {
			t.Errorf("the story lacks %q:\n%s", want, story)
		}
	}
	if epic := itemFile(t, root, "E-0001"); !strings.Contains(epic, "- S-0001 Release check for 1.4.0") {
		t.Errorf("the epic does not list the story:\n%s", epic)
	}
	for id, want := range map[string][]string{
		"T-0001": {"nature: remediation\ntitle: The tests pass\nstatus: backlog\nparent: S-0001\n", "## Work\n\nRun `scripts/test.sh --short` from the repository root.\n\n## Done when\n\n`scripts/test.sh --short` exits 0.\n"},
		"T-0002": {"after: [T-0001]", "## Work\n\nRun `flai release --evaluate`.\n\n## Done when\n\n`flai release --evaluate` exits 0.\n"},
		"T-0003": {"after: [T-0001]", "## Work\n\nRun the role `reviewer` with this prompt:\n\n> Report the risk in the changes since the last release before 1.4.0 is published.\n\n## Done when\n\nThe role's run ends with its work done.\n"},
	} {
		task := itemFile(t, root, id)
		for _, w := range want {
			if !strings.Contains(task, w) {
				t.Errorf("%s lacks %q:\n%s", id, w, task)
			}
		}
	}
	if strings.Contains(itemFile(t, root, "T-0001"), "after:") {
		t.Errorf("T-0001 waits for nothing:\n%s", itemFile(t, root, "T-0001"))
	}
	show, errOut, code := runIn(t, root, "show", "S-0001", "--json")
	var plan struct {
		Plan struct {
			Layers [][]string `json:"layers"`
		} `json:"plan"`
	}
	if err := json.Unmarshal([]byte(show), &plan); code != 0 || err != nil {
		t.Fatalf("show S-0001: %d %v %s", code, err, errOut)
	}
	if want := [][]string{{"T-0001"}, {"T-0002", "T-0003"}}; len(plan.Plan.Layers) != 2 || !slices.Equal(plan.Plan.Layers[0], want[0]) || !slices.Equal(plan.Plan.Layers[1], want[1]) {
		t.Errorf("the tasks' layers are not the steps': %v", plan.Plan.Layers)
	}
	// the fixture is no whole project, so the check is scoped to the story,
	// which takes in its tasks
	if out, errOut, code := runIn(t, root, "check", "--strict", "--story", "S-0001"); code != 0 {
		t.Errorf("check --strict on what formula run made: %d\n%s%s", code, out, errOut)
	}
	out, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=2.0.0", "--input", "area=docs")
	if code != 0 || !strings.HasPrefix(out, "S-0002 Release check for 2.0.0\n  wip/kanban/stories/S-0002-") ||
		!strings.Contains(out, "\n  T-0004 The tests pass (test)\n  T-0005 The release evaluation agrees on 2.0.0 (evaluate), after T-0004\n  T-0006 The reviewer reports the risk of 2.0.0 (risk), after T-0004\n") {
		t.Errorf("formula run: %d %s\n%s", code, errOut, out)
	}
	if story := itemFile(t, root, "S-0002"); strings.Contains(story, "parent: E") || !strings.Contains(story, "`area=docs`, `version=2.0.0`") {
		t.Errorf("a story without --epic, with both inputs given:\n%s", story)
	}
}

// S-0423: a step with in: main runs in the project's main checkout. flai
// formula show prints its in, and --json carries it; flai formula run writes
// into its task's Work that it runs from the main checkout, by its absolute
// path, while a step in the worktree, by in: worktree or by none, keeps the
// Work it had.
func TestFormulaStepInMain(t *testing.T) {
	t.Parallel()
	inMain := strings.Replace(releaseCheck, "    run: {flai: release --evaluate}\n    after: [test]\n", "    run: {flai: release --evaluate}\n    after: [test]\n    in: main\n", 1)
	inMain = strings.Replace(inMain, "    run: {script: scripts/test.sh --short}\n", "    run: {script: scripts/test.sh --short}\n    in: worktree\n", 1)
	if inMain == releaseCheck {
		t.Fatal("the fixture did not change")
	}
	root := formulaProject(t, map[string]string{"release-check.md": inMain})
	out, errOut, code := runIn(t, root, "formula", "show", "release-check")
	if code != 0 || strings.Contains(errOut, formulas.CodeFile) || !strings.Contains(out, "  evaluate: The release evaluation agrees on {{version}}\n    flai:   flai release --evaluate\n    after:  test\n    in:     main\n  risk:") ||
		!strings.Contains(out, "  test: The tests pass\n    script: scripts/test.sh --short\n  evaluate:") {
		t.Errorf("formula show: %d %s\n%s", code, errOut, out)
	}
	js, _, code := runIn(t, root, "formula", "show", "release-check", "--json")
	var shown struct {
		Steps []map[string]any `json:"steps"`
	}
	if err := json.Unmarshal([]byte(js), &shown); code != 0 || err != nil || len(shown.Steps) != 3 {
		t.Fatalf("formula show --json: %d %v %s", code, err, js)
	}
	if shown.Steps[0]["in"] != "worktree" || shown.Steps[1]["in"] != "main" {
		t.Errorf("formula show --json: the steps' in: %v", shown.Steps)
	}
	if _, ok := shown.Steps[2]["in"]; ok {
		t.Errorf("formula show --json: a step with no in has one: %v", shown.Steps[2])
	}
	if js, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=1.4.0", "--json"); code != 0 {
		t.Fatalf("formula run: %d %s %s", code, js, errOut)
	}
	checkout := codeSpan(root)
	for id, want := range map[string]string{
		"T-0001": "## Work\n\nRun `scripts/test.sh --short` from the repository root.\n\n## Done when\n\n`scripts/test.sh --short` exits 0.\n",
		"T-0002": "## Work\n\nRun `flai release --evaluate` from the project's main checkout, " + checkout + ", not from this story's worktree.\n\n## Done when\n\n`flai release --evaluate` exits 0 in the main checkout.\n",
		"T-0003": "## Work\n\nRun the role `reviewer` with this prompt:\n\n> Report the risk in the changes since the last release before 1.4.0 is published.\n\n## Done when\n\nThe role's run ends with its work done.\n",
	} {
		if task := itemFile(t, root, id); !strings.Contains(task, want) {
			t.Errorf("%s lacks %q:\n%s", id, want, task)
		}
	}
	if out, errOut, code := runIn(t, root, "check", "--strict", "--story", "S-0001"); code != 0 {
		t.Errorf("check --strict on what formula run made: %d\n%s%s", code, out, errOut)
	}
}

// S-0621 (ADR-0271 §1): a flai or script step's stop names the exits that
// stop the instance cleanly. flai formula show prints it, and --json carries
// it; flai formula run writes into the step task's Done when that each stop
// exit stops the run, and the flai formula stop on the instance that records
// it, while a step with no stop keeps the Done when it had.
func TestFormulaStepStop(t *testing.T) {
	t.Parallel()
	stops := strings.Replace(releaseCheck, "    run: {flai: release --evaluate}\n    after: [test]\n", "    run: {flai: release --evaluate}\n    after: [test]\n    in: main\n    stop: [2]\n", 1)
	stops = strings.Replace(stops, "    run: {script: scripts/test.sh --short}\n", "    run: {script: scripts/test.sh --short}\n    stop: [3, 4]\n", 1)
	if stops == releaseCheck {
		t.Fatal("the fixture did not change")
	}
	root := formulaProject(t, map[string]string{"release-check.md": stops})
	out, errOut, code := runIn(t, root, "formula", "show", "release-check")
	if code != 0 || strings.Contains(errOut, formulas.CodeFile) ||
		!strings.Contains(out, "  test: The tests pass\n    script: scripts/test.sh --short\n    stop:   3, 4\n  evaluate:") ||
		!strings.Contains(out, "    after:  test\n    in:     main\n    stop:   2\n  risk:") {
		t.Errorf("formula show: %d %s\n%s", code, errOut, out)
	}
	js, _, code := runIn(t, root, "formula", "show", "release-check", "--json")
	var shown struct {
		Steps []struct {
			Stop []int `json:"stop"`
		} `json:"steps"`
	}
	if err := json.Unmarshal([]byte(js), &shown); code != 0 || err != nil || len(shown.Steps) != 3 {
		t.Fatalf("formula show --json: %d %v %s", code, err, js)
	}
	if !slices.Equal(shown.Steps[0].Stop, []int{3, 4}) || !slices.Equal(shown.Steps[1].Stop, []int{2}) || shown.Steps[2].Stop != nil {
		t.Errorf("formula show --json: the steps' stop: %s", js)
	}
	if js, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=1.4.0", "--json"); code != 0 {
		t.Fatalf("formula run: %d %s %s", code, js, errOut)
	}
	for id, want := range map[string]string{
		"T-0001": "## Done when\n\n`scripts/test.sh --short` exits 0, or exits 3 or 4, each of which stops the run: record it with `flai formula stop S-0001 --step test --exit <exit>`, giving it the step's output.\n",
		"T-0002": "## Done when\n\n`flai release --evaluate` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0001 --step evaluate --exit 2`, giving it the step's output.\n",
		"T-0003": "## Done when\n\nThe role's run ends with its work done.\n",
	} {
		if task := itemFile(t, root, id); !strings.Contains(task, want) {
			t.Errorf("%s lacks %q:\n%s", id, want, task)
		}
	}
	if out, errOut, code := runIn(t, root, "check", "--strict", "--story", "S-0001"); code != 0 {
		t.Errorf("check --strict on what formula run made: %d\n%s%s", code, out, errOut)
	}
}

// S-0423: a step in main runs in the main checkout of a story worktree, the
// one its .git file names when it is not mounted here, or the project root.
func TestFormulaMainCheckout(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		repo workitem.Repo
		want string
	}{
		{"the main checkout", workitem.Repo{Root: "/p", MainRoot: "/p"}, "/p"},
		{"a story worktree", workitem.Repo{Root: "/p/.flai-cache/worktrees/S-1", MainRoot: "/p"}, "/p"},
		{"a worktree mounted alone", workitem.Repo{Root: "/w/S-1", MainRoot: "/w/S-1", MainAbsent: "/p"}, "/p"},
		{"no main checkout known", workitem.Repo{Root: "/p"}, "/p"},
	} {
		if got := formulaMainCheckout(&tc.repo); got != tc.want {
			t.Errorf("%s: %s, want %s", tc.name, got, tc.want)
		}
	}
}

// S-0421 (ADR-0238 §7): flai formula run refuses, before it writes anything,
// a malformed --input, one given twice, an undeclared input, a missing
// required one, an unknown formula, and one flai check reports a problem in;
// and when an item it makes is refused, it removes the story and the tasks
// it made and restores the epic.
func TestFormulaRunRefuses(t *testing.T) {
	t.Parallel()
	invalid := strings.Replace(strings.Replace(releaseCheck, "name: release-check", "name: invalid", 1), "role: reviewer", "role: nobody", 1)
	untitled := strings.Replace(strings.Replace(releaseCheck, "name: release-check", "name: untitled", 1), "title: The release evaluation agrees on {{version}}", `title: "..."`, 1)
	root := formulaProject(t, map[string]string{"release-check.md": releaseCheck, "invalid.md": invalid, "untitled.md": untitled})
	if _, errOut, code := runIn(t, root, "epic", "new", "Releases"); code != 0 {
		t.Fatal(errOut)
	}
	epicWas := itemFile(t, root, "E-0001")
	for _, tc := range []struct {
		name string
		args []string
		says string
	}{
		{"malformed", []string{"release-check", "--input", "version"}, `is not name=value: give each input as its name`},
		{"twice", []string{"release-check", "--input", "version=1", "--input", "version=2"}, "--input version is given twice"},
		{"undeclared", []string{"release-check", "--input", "version=1", "--input", "colour=red"}, "formula release-check declares no input colour"},
		{"missing", []string{"release-check"}, "formula release-check needs input version, which has no default; its inputs are version (required), area"},
		{"unknown", []string{"nope"}, "no formula is named nope; the formulas are"},
		{"invalid", []string{"invalid", "--input", "version=1"}, "formula invalid has 1 problem, so nothing was made"},
		{"refused item", []string{"untitled", "--input", "version=1", "--epic", "E-0001"}, "making the task for step evaluate of formula untitled, so the story S-0001 and the tasks made before it were removed: a title is required"},
	} {
		_, errOut, code := runIn(t, root, append([]string{"formula", "run"}, tc.args...)...)
		if code == 0 || !strings.Contains(errOut, tc.says) {
			t.Errorf("%s: %d %s, want %q", tc.name, code, errOut, tc.says)
		}
		if made := itemFiles(t, root); len(made) != 0 {
			t.Errorf("%s: wrote %v", tc.name, made)
		}
	}
	if epic := itemFile(t, root, "E-0001"); epic != epicWas {
		t.Errorf("the epic was not restored:\n%s", epic)
	}
}

// S-0422: flai formula run --ready moves the instance to ready once it and
// its tasks are made, as flai move does: the transition recorded, the epic
// following, the board placing it in its order; --json says the status, and
// the output names the move.
func TestFormulaRunReady(t *testing.T) {
	t.Parallel()
	root := formulaProject(t, map[string]string{"release-check.md": releaseCheck})
	if _, errOut, code := runIn(t, root, "epic", "new", "Releases"); code != 0 {
		t.Fatal(errOut)
	}
	js, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=1.4.0", "--epic", "E-0001", "--ready", "--json")
	if code != 0 {
		t.Fatalf("formula run --ready: %d %s %s", code, js, errOut)
	}
	var got formulaRunResult
	if err := json.Unmarshal([]byte(js), &got); err != nil {
		t.Fatalf("formula run --ready --json: %v %s", err, js)
	}
	if got.Story != "S-0001" || got.Status != "ready" || len(got.Tasks) != 3 || got.Followed == nil || got.Followed.ID != "E-0001" || got.Followed.To != "ready" {
		t.Errorf("formula run --ready --json: %+v", got)
	}
	story := itemFile(t, root, "S-0001")
	if !strings.Contains(story, "\nstatus: ready\n") || !strings.Contains(story, "  - to: ready\n") {
		t.Errorf("the story is not ready, with its transition:\n%s", story)
	}
	for _, id := range []string{"T-0001", "T-0002", "T-0003"} {
		if !strings.Contains(itemFile(t, root, id), "\nstatus: backlog\n") {
			t.Errorf("%s left backlog:\n%s", id, itemFile(t, root, id))
		}
	}
	if epic := itemFile(t, root, "E-0001"); !strings.Contains(epic, "\nstatus: ready\n") {
		t.Errorf("the epic did not follow the story:\n%s", epic)
	}
	if board, err := os.ReadFile(filepath.Join(root, "wip", "kanban", "board.md")); err != nil || !strings.Contains(string(board), "S-0001") {
		t.Errorf("the board does not order the story: %v\n%s", err, board)
	}
	out, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=2.0.0", "--ready")
	if code != 0 || !strings.HasPrefix(out, "S-0002 Release check for 2.0.0\n") || !strings.Contains(out, "\n  S-0002 → ready\n") {
		t.Errorf("formula run --ready: %d %s\n%s", code, errOut, out)
	}
	if out, errOut, code := runIn(t, root, "check", "--strict", "--story", "S-0002"); code != 0 {
		t.Errorf("check --strict on the ready instance: %d\n%s%s", code, out, errOut)
	}
}

// I-0282 (S-0620): the instance claims its own story file, so a story in
// progress, such as the one whose agent ran the formula, does not hold it as
// a story that declares no touches.
func TestFormulaRunInstanceIsNotHeldByAStoryInProgress(t *testing.T) {
	t.Parallel()
	root := formulaProject(t, map[string]string{"release-check.md": releaseCheck})
	if _, errOut, code := runIn(t, root, "story", "new", "Starts the formula", "--touches", "flai/cmd"); code != 0 {
		t.Fatal(errOut)
	}
	matches, _ := filepath.Glob(filepath.Join(root, "wip", "kanban", "stories", "S-0001-*.md"))
	if len(matches) != 1 {
		t.Fatalf("S-0001: %v", matches)
	}
	body, _ := os.ReadFile(matches[0])
	if err := os.WriteFile(matches[0], []byte(strings.Replace(string(body), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] done\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{"ready", "in-progress"} {
		if _, errOut, code := runIn(t, root, "move", "S-0001", to); code != 0 {
			t.Fatalf("move S-0001 %s: %s", to, errOut)
		}
	}
	js, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=1.4.0", "--ready", "--json")
	if code != 0 {
		t.Fatalf("formula run --ready: %d %s %s", code, js, errOut)
	}
	var run formulaRunResult
	if err := json.Unmarshal([]byte(js), &run); err != nil {
		t.Fatalf("formula run --ready --json: %v %s", err, js)
	}
	if run.Story != "S-0002" || run.Status != "ready" {
		t.Fatalf("formula run --ready --json: %+v", run)
	}
	if want := "\ntouches: [" + run.Path + "]\n"; !strings.Contains(itemFile(t, root, "S-0002"), want) {
		t.Errorf("the instance does not touch its own story file %s alone:\n%s", run.Path, itemFile(t, root, "S-0002"))
	}
	out, errOut, code := runIn(t, root, "board", "--json")
	if code != 0 {
		t.Fatal(errOut)
	}
	var view struct {
		Columns map[string][]struct {
			ID   string `json:"id"`
			Held *struct {
				Code   string `json:"code"`
				Reason string `json:"reason"`
			} `json:"held"`
		} `json:"columns"`
	}
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatal(err)
	}
	ready := view.Columns["ready"]
	if len(ready) != 1 || ready[0].ID != "S-0002" || ready[0].Held != nil {
		t.Errorf("the instance is held by the story in progress: %+v", ready)
		for _, r := range ready {
			if r.Held != nil {
				t.Logf("%s %s: %s", r.ID, r.Held.Code, r.Held.Reason)
			}
		}
	}
}

// S-0423, ADR-0240: outside the orchestrator's role, the operator or a
// story's agent runs any formula, publish among them, whatever
// orchestration.permissions.publish says; only the orchestrator is held to it
// (testTheOrchestratorRunsThePublishFormulaOnlyWithPublish).
func TestFormulaRunUnderNoRole(t *testing.T) {
	t.Parallel()
	root := formulaProject(t, map[string]string{"publish.md": formulaFile("publish"), "import.md": formulaFile("import")})
	for _, name := range []string{"publish", "import"} {
		if out, errOut, code := runIn(t, root, "formula", "run", name, "--ready"); code != 0 || !strings.Contains(out, " → ready\n") {
			t.Errorf("formula run %s --ready: %d %s\n%s", name, code, errOut, out)
		}
	}
}

// S-0422: flai formula run --ready --autocommit, as flai serve runs it in
// the main checkout, commits the story, its tasks, its epic, the board, and
// the stream index in the one commit, leaving nothing uncommitted.
func TestFormulaRunReadyAutocommit(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := formulaProject(t, map[string]string{"release-check.md": releaseCheck})
	if _, errOut, code := runIn(t, root, "epic", "new", "Releases"); code != 0 {
		t.Fatal(errOut)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "olive@example.invalid")
	gitIn(t, root, "config", "user.name", "Olive")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	js, errOut, code := runIn(t, root, "formula", "run", "release-check", "--input", "version=1.4.0", "--epic", "E-0001", "--ready", "--autocommit", "--trailer", "Co-Authored-By: flai <flai@localhost>", "--json")
	if code != 0 {
		t.Fatalf("formula run --ready --autocommit: %d %s %s", code, js, errOut)
	}
	var got formulaRunResult
	if err := json.Unmarshal([]byte(js), &got); err != nil || got.Status != "ready" || got.Commit == "" || got.CommitError != "" {
		t.Fatalf("formula run --ready --autocommit --json: %v %s", err, js)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("the run left files uncommitted: %q", st)
	}
	files := gitIn(t, root, "show", "--name-only", "--format=%B", "HEAD")
	for _, want := range []string{
		"chore: [S-0001] create story: Release check for 1.4.0, from formula release-check, ready\n",
		"Co-Authored-By: flai <flai@localhost>",
		"wip/kanban/stories/S-0001-", "wip/kanban/tasks/T-0003-", "wip/kanban/epics/E-0001-", "wip/kanban/board.md",
	} {
		if !strings.Contains(files, want) {
			t.Errorf("the commit lacks %q:\n%s", want, files)
		}
	}
	if story := gitIn(t, root, "show", "HEAD:"+got.Path); !strings.Contains(story, "\nstatus: ready\n") {
		t.Errorf("the committed story is not ready:\n%s", story)
	}
}

// S-0422: with --ready, a missing or undeclared input is refused before
// anything is written: no item, no board, and no move.
func TestFormulaRunReadyRefusesBeforeWriting(t *testing.T) {
	t.Parallel()
	root := formulaProject(t, map[string]string{"release-check.md": releaseCheck})
	for _, tc := range []struct {
		name string
		args []string
		says string
	}{
		{"missing", []string{"release-check", "--ready"}, "formula release-check needs input version, which has no default"},
		{"undeclared", []string{"release-check", "--input", "version=1", "--input", "colour=red", "--ready", "--autocommit"}, "formula release-check declares no input colour"},
	} {
		_, errOut, code := runIn(t, root, append([]string{"formula", "run"}, tc.args...)...)
		if code == 0 || !strings.Contains(errOut, tc.says) {
			t.Errorf("%s: %d %s, want %q", tc.name, code, errOut, tc.says)
		}
		if made := itemFiles(t, root); len(made) != 0 {
			t.Errorf("%s: wrote %v", tc.name, made)
		}
		if _, err := os.Stat(filepath.Join(root, "wip", "kanban", "board.md")); !os.IsNotExist(err) {
			t.Errorf("%s: the board was written: %v", tc.name, err)
		}
	}
}

// S-0421: the tasks are made in file order, save that a step comes after
// each step its after names.
func TestFormulaCreationOrder(t *testing.T) {
	t.Parallel()
	f := formulas.Formula{Steps: []formulas.Step{{ID: "a", After: []string{"c"}}, {ID: "b"}, {ID: "c"}, {ID: "d", After: []string{"b"}}}}
	var got []string
	for _, s := range creationOrder(f) {
		got = append(got, s.ID)
	}
	if want := []string{"b", "c", "a", "d"}; !slices.Equal(got, want) {
		t.Errorf("creation order %v, want %v", got, want)
	}
}
