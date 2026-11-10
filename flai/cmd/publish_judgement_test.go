package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// publishFormulaFile is this repository's own publish formula, the file the
// template ships and the orchestrator starts under publish (ADR-0240).
const publishFormulaFile = "../../design/workflows/formulas/publish.md"

// stepCommand is the command line the Work of a formula step's task runs,
// split into arguments by sh as the instance's agent's shell splits it, with
// the leading flai dropped.
func stepCommand(t *testing.T, task string) []string {
	t.Helper()
	_, work, ok := strings.Cut(task, "## Work\n\nRun `")
	line, _, closed := strings.Cut(work, "`")
	if !ok || !closed {
		t.Fatalf("the task's Work runs no command:\n%s", task)
	}
	out, err := exec.Command("sh", "-c", "set -- "+line+"; for a do printf '%s\\000' \"$a\"; done").Output()
	if err != nil {
		t.Fatalf("sh splits %q: %v", line, err)
	}
	args := strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00")
	if len(args) < 2 || args[0] != "flai" {
		t.Fatalf("the step does not run flai: %q", args)
	}
	return args[1:]
}

// evaluateTask is the task of the publish instance's evaluate step that
// flai formula run publish, with args, makes in root.
func evaluateTask(t *testing.T, root string, args ...string) string {
	t.Helper()
	js, errOut, code := runIn(t, root, append([]string{"formula", "run", "publish", "--json"}, args...)...)
	if code != 0 {
		t.Fatalf("formula run publish %v: %d %s", args, code, errOut)
	}
	var got formulaRunResult
	if err := json.Unmarshal([]byte(js), &got); err != nil {
		t.Fatalf("formula run --json: %v %s", err, js)
	}
	for _, task := range got.Tasks {
		if task.Step == "evaluate" {
			return itemFile(t, root, task.ID)
		}
	}
	t.Fatalf("the instance has no evaluate step: %+v", got.Tasks)
	return ""
}

// I-0291, ADR-0273: under judgement the publish formula's evaluate step
// stopped every instance, so the publish permission published nothing. With
// the orchestrator's reasoning given as the formula's judgement input, the
// step's command, as its task writes it and a shell splits it, meets the
// policy; with none it still stops the run.
func TestPublishFormulaUnderJudgement(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not installed")
	}
	root, _ := evaluateProject(t)
	setReleasePolicy(t, root, "    policy: judgement\n")
	formula, err := os.ReadFile(publishFormulaFile)
	if err != nil {
		t.Fatal(err)
	}
	for dir, files := range map[string]map[string]string{
		"design/workflows/formulas": {"publish.md": string(formula)},
		"scripts":                   {"fetch.sh": "#!/bin/sh\n"},
	} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
		mode := os.FileMode(0o644)
		if dir == "scripts" {
			mode = 0o755
		}
		for name, content := range files {
			if err := os.WriteFile(filepath.Join(root, dir, name), []byte(content), mode); err != nil {
				t.Fatal(err)
			}
		}
	}

	args := stepCommand(t, evaluateTask(t, root))
	_, errOut, code := runIn(t, root, args...)
	if code != 2 || !strings.Contains(errOut, "release policy judgement not met") {
		t.Errorf("with no judgement the evaluate step %v exits 2, stopping the run: %d %s", args, code, errOut)
	}

	reason := "both stories are whole and tested, so the batch is coherent"
	task := evaluateTask(t, root, "--input", "judgement="+reason)
	if !strings.Contains(task, `--judged "`+reason+`"`) {
		t.Errorf("the evaluate step's Work names the judgement:\n%s", task)
	}
	args = stepCommand(t, task)
	out, errOut, code := runIn(t, root, args...)
	if code != 0 || !strings.HasPrefix(out, "release policy judgement: met\n") || !strings.Contains(out, reason) {
		t.Errorf("with a judgement the evaluate step %v exits 0, so the publish step runs: %d %s %s", args, code, out, errOut)
	}
}
