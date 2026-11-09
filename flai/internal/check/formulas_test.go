package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// formulaFile is a valid formula file, release: a script step, a role step
// naming the reviewer roleFile defines, and a flai step after both, with an
// input the title, a prompt, and a command line use. Its lines are numbered
// in the comments the tests' lines refer to.
const formulaFile = `---
title: Release {{version}}
updated: 2026-11-01
status: active
name: release
description: Cuts a release.
inputs:
  - name: version
    description: The version to release.
steps:
  - id: check
    title: Check the tree
    run:
      script: scripts/check.sh --strict
  - id: notes
    title: Write the notes
    run:
      role: reviewer
      prompt: Draft the notes for {{version}}.
  - id: tag
    title: Tag {{version}}
    run:
      flai: release publish {{version}}
    after: [check, notes]
---

# Release

Cuts a release.
`

// Lines of formulaFile: the check step's run is on 13, the notes step's run
// on 17, the tag step's id on 20, its run on 22, and its after on 24.

const formulaPath = "design/workflows/formulas/release.md"

// formulasProject is a project with the reviewer role file, scripts/check.sh,
// and the design folder's formula files files.
func formulasProject(t *testing.T, files map[string]string) *workitem.Repo {
	t.Helper()
	repo := rolesProject(t, map[string]string{"reviewer.md": roleFile("reviewer", nil)})
	writeFormulaFiles(t, repo.Root, files)
	scripts := filepath.Join(repo.Root, "scripts")
	if err := os.MkdirAll(scripts, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scripts, "check.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return repo
}

// writeFormulaFiles writes files into the formula files' folder under root.
func writeFormulaFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	dir := filepath.Join(root, "design", "workflows", "formulas")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// formulaFindings are the findings flai check's formulas rule gives, without
// the other rules.
func formulaFindings(repo *workitem.Repo) []Finding {
	c := &checker{repo: repo, res: &Result{}}
	c.formulas()
	sortFindings(c.res.Findings)
	return c.res.Findings
}

// S-0421: a valid formula, and the index beside it, are not reported, and
// flai check runs the rule.
func TestValidFormulaIsClean(t *testing.T) {
	repo := formulasProject(t, map[string]string{
		"README.md":  "# Formulas\n\nNo front matter: the index is not a formula.\n",
		"release.md": formulaFile,
	})
	if got := formulaFindings(repo); len(got) != 0 {
		t.Errorf("a valid formula: %+v", got)
	}
	if got := findings(t, repo, formulas.CodeFile); len(got) != 0 {
		t.Errorf("flai check on a valid formula: %+v", got)
	}
}

// S-0423: a step may say in: main or in: worktree, and flai check reports
// neither.
func TestAStepsInMainOrWorktreeIsClean(t *testing.T) {
	content := strings.Replace(formulaFile, "    after: [check, notes]\n", "    after: [check, notes]\n    in: main\n", 1)
	content = strings.Replace(content, "      script: scripts/check.sh --strict\n", "      script: scripts/check.sh --strict\n    in: worktree\n", 1)
	repo := formulasProject(t, map[string]string{"release.md": content})
	if got := findings(t, repo, formulas.CodeFile); len(got) != 0 {
		t.Errorf("flai check on steps in main and in the worktree: %+v", got)
	}
}

// S-0421: a project with no formula files' folder reports nothing.
func TestNoFormulasFolderReportsNothing(t *testing.T) {
	repo := planningProject(t, "")
	if got := formulaFindings(repo); len(got) != 0 {
		t.Errorf("no formulas folder: %+v", got)
	}
	if got := findings(t, repo, formulas.CodeFile); len(got) != 0 {
		t.Errorf("flai check with no formulas folder: %+v", got)
	}
}

// S-0421 (ADR-0238 §6): each way a formula falls short is one formulas.file
// error in its file, on the line of what it is about, through flai check.
func TestFormulaProblemsAreErrorsOnTheirLine(t *testing.T) {
	for _, tc := range []struct {
		name, old, new string
		line           int
		says           string
	}{
		{"no front matter", formulaFile, "# Release\n", 1, "front matter"},
		{"a steps entry that does not parse", "  - id: check\n    title: Check the tree\n    run:\n      script: scripts/check.sh --strict\n", "  - check\n", 11, "a steps entry is a mapping"},
		{"a missing field", "description: Cuts a release.\n", "", 1, "front matter needs description"},
		{"a run of none of the three kinds", "      flai: release publish {{version}}\n", "      prompt: Tag it.\n", 22, "gives none of flai, script, or role"},
		{"a run of two kinds", "      flai: release publish {{version}}\n", "      flai: release publish {{version}}\n      script: scripts/check.sh\n", 22, "gives flai and script"},
		{"a script not under scripts/", "script: scripts/check.sh", "script: tools/check.sh", 13, "is not under scripts/"},
		{"a script that does not exist", "script: scripts/check.sh", "script: scripts/missing.sh", 13, "scripts/missing.sh does not exist"},
		{"a role no file defines", "role: reviewer", "role: auditor", 17, "role auditor is no role file's"},
		{"a role named by its verb", "role: reviewer", "role: review", 17, "role review is no role file's; name one of reviewer"},
		{"an after naming no step", "after: [check, notes]", "after: [check, publish]", 24, `after "publish", which is no step's id`},
		{"a cycle", "      script: scripts/check.sh --strict\n", "      script: scripts/check.sh --strict\n    after: [tag]\n", 15, "in a cycle: check after tag after check"},
		{"a duplicate step id", "  - id: tag\n", "  - id: notes\n", 20, "step id notes is also an earlier step's"},
		{"an input a prompt uses undeclared", "for {{version}}.", "for {{release}}.", 17, "uses {{release}}, which the formula's inputs do not declare"},
		{"an input a title uses undeclared", "title: Tag {{version}}", "title: Tag {{tag}}", 21, "uses {{tag}}, which the formula's inputs do not declare"},
		{"an in that is neither place", "    after: [check, notes]\n", "    after: [check, notes]\n    in: elsewhere\n", 25, `step tag's in is "elsewhere", which is none of worktree, main`},
		{"a stop on a role step", "      prompt: Draft the notes for {{version}}.\n", "      prompt: Draft the notes for {{version}}.\n    stop: [2]\n", 20, "step notes names stop, but a role step ends worked, asked, or failed with no exit code"},
		{"a stop exit out of range", "    after: [check, notes]\n", "    after: [check, notes]\n    stop: [0]\n", 25, "step tag's stop lists exit 0; a stop exit is from 1 to 255"},
		{"a stop that is no list", "    after: [check, notes]\n", "    after: [check, notes]\n    stop: 2\n", 25, "step tag's stop is not a list of exit codes"},
		{"a stop exit twice", "    after: [check, notes]\n", "    after: [check, notes]\n    stop: [2, 2]\n", 25, "step tag's stop lists exit 2 twice"},
		{"a schedule on a formula with inputs", "description: Cuts a release.\n", "description: Cuts a release.\nschedule: daily\n", 7, "schedule starts the formula with no inputs, as nobody is there to give them, but version has no default"},
	} {
		if !strings.Contains(formulaFile, tc.old) {
			t.Fatalf("%s: the fixture has no %q", tc.name, tc.old)
		}
		repo := formulasProject(t, map[string]string{"release.md": strings.Replace(formulaFile, tc.old, tc.new, 1)})
		got := formulaFindings(repo)
		if len(got) != 1 {
			t.Errorf("%s: want one finding, got %+v", tc.name, got)
			continue
		}
		f := got[0]
		if f.Level != Error || f.Rule != formulas.CodeFile || f.Path != formulaPath || f.Line != tc.line || !strings.Contains(f.Message, tc.says) {
			t.Errorf("%s: %+v, want a %s error on %s:%d saying %q", tc.name, f, formulas.CodeFile, formulaPath, tc.line, tc.says)
		}
		if all := findings(t, repo, formulas.CodeFile); len(all) != 1 {
			t.Errorf("%s: flai check gives %+v under %s", tc.name, all, formulas.CodeFile)
		}
	}
}

// scheduledFile is a formula with no inputs that names a schedule, on line
// 7, for flai serve to make its instance on.
const scheduledFile = `---
title: Nightly check
updated: 2026-11-01
status: active
name: nightly
description: Checks the tree every night.
schedule: "0 2 * * *"
steps:
  - id: check
    title: Check the tree
    run:
      script: scripts/check.sh --strict
---

# Nightly
`

// S-0422: a formula with no inputs may name a schedule, which flai check
// reads as planning.schedule is read, and one it cannot parse is a
// formulas.file error on the schedule's line.
func TestAFormulasScheduleIsChecked(t *testing.T) {
	const path = "design/workflows/formulas/nightly.md"
	repo := formulasProject(t, map[string]string{"nightly.md": scheduledFile})
	if got := formulaFindings(repo); len(got) != 0 {
		t.Errorf("a valid schedule: %+v", got)
	}
	repo = formulasProject(t, map[string]string{"nightly.md": strings.Replace(scheduledFile, `"0 2 * * *"`, "nightly", 1)})
	got := findings(t, repo, formulas.CodeFile)
	if len(got) != 1 || got[0].Level != Error || got[0].Path != path || got[0].Line != 7 || !strings.Contains(got[0].Message, `schedule "nightly" has 1 fields, not five`) {
		t.Errorf("a schedule that does not parse: %+v", got)
	}
}

// S-0424: a scheduled formula whose every input carries a default passes, and
// one with an input of no default is a formulas.file error on the schedule's
// line naming that input.
func TestAScheduledFormulasInputsNeedDefaults(t *testing.T) {
	defaulted := strings.Replace(formulaFile, "description: Cuts a release.\n", "description: Cuts a release.\nschedule: daily\n", 1)
	defaulted = strings.Replace(defaulted, "    description: The version to release.\n", "    description: The version to release.\n    default: next\n", 1)
	repo := formulasProject(t, map[string]string{"release.md": defaulted})
	if got := findings(t, repo, formulas.CodeFile); len(got) != 0 {
		t.Errorf("inputs that all carry a default: %+v", got)
	}
	twoInputs := strings.Replace(defaulted, "steps:\n", "  - name: channel\n    description: Where to announce it.\nsteps:\n", 1)
	repo = formulasProject(t, map[string]string{"release.md": twoInputs})
	got := findings(t, repo, formulas.CodeFile)
	if len(got) != 1 || got[0].Level != Error || got[0].Path != formulaPath || got[0].Line != 7 ||
		!strings.Contains(got[0].Message, "but channel has no default; give it a default, or remove the schedule") {
		t.Errorf("an input with no default: %+v", got)
	}
}

// S-0421: a formula file that does not parse leaves the others checked.
func TestFormulaFilesAreCheckedTogether(t *testing.T) {
	repo := formulasProject(t, map[string]string{
		"broken.md":  "no front matter\n",
		"release.md": strings.Replace(formulaFile, "role: reviewer", "role: auditor", 1),
	})
	got := formulaFindings(repo)
	if len(got) != 2 {
		t.Fatalf("want the broken file and the release's role: %+v", got)
	}
	if got[0].Path != "design/workflows/formulas/broken.md" || got[0].Line != 1 {
		t.Errorf("the file that does not parse: %+v", got[0])
	}
	if got[1].Path != formulaPath || got[1].Line != 17 {
		t.Errorf("the other file is still checked: %+v", got[1])
	}
}

// S-0421: scoped to a story, a formula finding is held to the paths the
// story changes, as a role file's is: outside it unless the story changes
// the formula file or a folder that holds it.
func TestScopedCheckHoldsAFormulaToTheStorysPaths(t *testing.T) {
	root := scopeProject(t)
	writeFormulaFiles(t, root, map[string]string{"release.md": strings.Replace(formulaFile, "after: [check, notes]", "after: [check, publish]", 1)})
	for _, tc := range []struct {
		name    string
		changed []string
		outside bool
	}{
		{"not changed", nil, true},
		{"the file changed", []string{formulaPath}, false},
		{"its folder changed", []string{"design/workflows/"}, false},
	} {
		_, scoped := runScoped(t, root, tc.changed)
		var got []Finding
		for _, f := range scoped.Findings {
			if f.Rule == formulas.CodeFile && f.Path == formulaPath && f.Line == 24 {
				got = append(got, f)
			}
		}
		if len(got) != 1 || got[0].Outside != tc.outside {
			t.Errorf("%s: %+v, want one finding on the after line, outside %v", tc.name, got, tc.outside)
		}
	}
}

// analysisFormulaProject is a formulasProject whose manifest names formula
// as analysis.formula, with files as its formula files.
func analysisFormulaProject(t *testing.T, formula string, files map[string]string) *workitem.Repo {
	t.Helper()
	repo := formulasProject(t, files)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nanalysis:\n  schedule: daily\n  formula: " + formula + "\n"
	if err := os.WriteFile(filepath.Join(repo.Root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	again, err := workitem.Open(repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	return again
}

// S-0424: analysis.formula names a formula a file defines whose every input
// carries a default; one no file defines, or one with an input of no
// default, is a manifest.analysis error on the manifest's analysis line. A
// formula whose file does not parse is left to its formulas.file error.
func TestAnalysisFormulaIsHeldToTheFormulaFiles(t *testing.T) {
	defaulted := strings.Replace(formulaFile, "    description: The version to release.\n", "    description: The version to release.\n    default: next\n", 1)
	for _, tc := range []struct {
		name, formula string
		files         map[string]string
		says          string
	}{
		{"a formula whose inputs all carry a default", "release", map[string]string{"release.md": defaulted}, ""},
		{"a formula no file defines", "cadence", map[string]string{"release.md": defaulted}, "analysis.formula: no formula is named cadence; the formulas are release"},
		{"a formula with an input of no default", "release", map[string]string{"release.md": formulaFile},
			"analysis.formula names release, whose input version has no default, but analysis.schedule makes its instance with no inputs; give version a default in design/workflows/formulas/release.md"},
		{"a formula whose file does not parse", "release", map[string]string{"release.md": "no front matter\n"}, ""},
	} {
		repo := analysisFormulaProject(t, tc.formula, tc.files)
		got := findings(t, repo, "manifest.analysis")
		if tc.says == "" {
			if len(got) != 0 {
				t.Errorf("%s: %+v", tc.name, got)
			}
			continue
		}
		if len(got) != 1 || got[0].Level != Error || got[0].Path != "system-flow.yaml" || got[0].Line != 8 || !strings.Contains(got[0].Message, tc.says) {
			t.Errorf("%s: %+v, want one manifest.analysis error on system-flow.yaml:8 saying %q", tc.name, got, tc.says)
		}
	}
	repo := analysisFormulaProject(t, "release", map[string]string{"release.md": "no front matter\n"})
	if got := findings(t, repo, formulas.CodeFile); len(got) != 1 {
		t.Errorf("the file that does not parse: %+v", got)
	}
}
