package cmd

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
)

// publishLike is a formula like the publish formula (ADR-0271): a check, an
// evaluation that stops the run cleanly on exit 2, a note in the same layer
// as the evaluation, and a publish after the evaluation.
const publishLike = `---
title: "Publish {{version}}"
updated: 2026-11-09
status: active
name: publish-like
description: "Publish what is accepted."
nature: improvement
inputs:
  - name: version
    description: The version.
    default: next
steps:
  - id: check
    title: The tests pass
    run: {script: scripts/test.sh}
  - id: evaluate
    title: The release evaluation is met
    run: {flai: release --evaluate --require-met}
    after: [check]
    stop: [2]
  - id: note
    title: The board is read
    run: {flai: board}
    after: [check]
  - id: publish
    title: Publish {{version}}
    run: {flai: release}
    after: [evaluate]
---

# Publish

<!-- system-flow:end-of-baseline -->

## Project additions
`

// stopProject is a project with the publish-like formula, its instance
// S-0001 in progress with a narrative, the task of its check step, T-0001,
// done, and that of its evaluate step, T-0002, in progress.
func stopProject(t *testing.T) string {
	t.Helper()
	root := formulaProject(t, map[string]string{"publish-like.md": publishLike})
	for _, args := range [][]string{
		{"formula", "run", "publish-like", "--ready"},
		{"move", "S-0001", "in-progress"},
		{"stream", "open", "S-0001", "--no-branch"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"},
	} {
		if out, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %d %s %s", args, code, out, errOut)
		}
	}
	return root
}

// S-0621 (ADR-0271 §2): a step that exits with an exit its stop lists stops
// the run: its task goes to done, unblocked, with the exit and the output's
// last lines in its Notes; the step tasks not done, the one in its own layer
// too, are cancelled with the stop's reason; the criteria are ticked, the
// skipped ones saying so; and the story's Notes and narrative log the stop.
func TestFormulaStop(t *testing.T) {
	t.Parallel()
	root := stopProject(t)
	if _, errOut, code := runIn(t, root, "block", "T-0002", "--reason", "step evaluate exited 2"); code != 0 {
		t.Fatal(errOut)
	}
	var output strings.Builder
	for i := 1; i <= 45; i++ {
		output.WriteString("line " + strings.Repeat("x", i%3) + "\n")
	}
	output.WriteString("not met: 0 accepted stories not yet released\n")
	js, errOut, code := runStdin(t, root, output.String(), "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2", "--json")
	if code != 0 {
		t.Fatalf("formula stop: %d %s %s", code, js, errOut)
	}
	var got formulaStopResult
	if err := json.Unmarshal([]byte(js), &got); err != nil {
		t.Fatalf("formula stop --json: %v %s", err, js)
	}
	if got.Story != "S-0001" || got.Step != "evaluate" || got.Exit != 2 || got.Done != "T-0002" ||
		!slices.Equal(got.Cancelled, []string{"T-0003", "T-0004"}) || !slices.Equal(got.Ticked, []int{2, 3, 4}) ||
		!slices.Equal(got.Skipped, []int{3, 4}) || !got.Logged {
		t.Errorf("formula stop --json: %+v", got)
	}
	evaluate := itemFile(t, root, "T-0002")
	if !strings.Contains(evaluate, "\nstatus: done\n") || !strings.Contains(evaluate, "until: ") {
		t.Errorf("T-0002 is not done:\n%s", evaluate)
	}
	if !strings.Contains(evaluate, "## Notes\n\nStopped the run with exit 2:\n\n```text\n") ||
		!strings.Contains(evaluate, "not met: 0 accepted stories not yet released\n```\n") ||
		strings.Count(evaluate, "line ") != 39 {
		t.Errorf("T-0002's Notes lack the exit and the output's last 40 lines:\n%s", evaluate)
	}
	for _, id := range []string{"T-0003", "T-0004"} {
		task := itemFile(t, root, id)
		if !strings.Contains(task, "\nstatus: cancelled\n") || !strings.Contains(task, ": moved to cancelled: step evaluate stopped the run with exit 2\n") {
			t.Errorf("%s is not cancelled by the stop:\n%s", id, task)
		}
	}
	story := itemFile(t, root, "S-0001")
	for _, want := range []string{
		"## Acceptance criteria\n\n- [ ] The tests pass\n- [x] The release evaluation is met\n" +
			"- [x] The board is read (skipped: step evaluate stopped the run with exit 2)\n" +
			"- [x] Publish next (skipped: step evaluate stopped the run with exit 2)\n",
		"\nStep evaluate stopped the run with exit 2; skipped T-0003, T-0004.\n",
		"\nstatus: in-progress\n",
	} {
		if !strings.Contains(story, want) {
			t.Errorf("the story lacks %q:\n%s", want, story)
		}
	}
	narrative, err := os.ReadFile(filepath.Join(root, "wip", "agents", "S-0001.md"))
	if err != nil || !strings.Contains(string(narrative), "\nStep evaluate stopped the run with exit 2; skipped T-0003, T-0004.\n") {
		t.Errorf("the narrative does not log the stop: %v\n%s", err, narrative)
	}
	if out, errOut, code := runIn(t, root, "check", "--strict", "--story", "S-0001"); code != 0 {
		t.Errorf("check --strict after the stop: %d\n%s%s", code, out, errOut)
	}
	// the story goes to review once its tasks are closed
	if out, errOut, code := runIn(t, root, "criteria", "tick", "S-0001", "1"); code != 0 {
		t.Fatalf("criteria tick: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Errorf("move to review after the stop: %d %s %s", code, out, errOut)
	}
}

// S-0621 (ADR-0271 §2): the output comes from --output as well, and a step
// with nothing after it to skip stops the run with nothing cancelled.
func TestFormulaStopOutputFile(t *testing.T) {
	t.Parallel()
	root := stopProject(t)
	for _, id := range []string{"T-0003"} {
		for _, to := range []string{"ready", "in-progress", "done"} {
			if _, errOut, code := runIn(t, root, "move", id, to); code != 0 {
				t.Fatal(errOut)
			}
		}
	}
	file := filepath.Join(t.TempDir(), "out.txt")
	if err := os.WriteFile(file, []byte("nothing to release\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2", "--output", file)
	if code != 0 || !strings.HasPrefix(out, "S-0001 Step evaluate stopped the run with exit 2; skipped T-0004.\n  T-0002 → done (step evaluate, exit 2)\n  T-0004 → cancelled: step evaluate stopped the run with exit 2\n") {
		t.Fatalf("formula stop --output: %d %s\n%s", code, errOut, out)
	}
	if task := itemFile(t, root, "T-0002"); !strings.Contains(task, "```text\nnothing to release\n```") {
		t.Errorf("T-0002 lacks the output:\n%s", task)
	}
	if task := itemFile(t, root, "T-0003"); !strings.Contains(task, "\nstatus: done\n") {
		t.Errorf("T-0003, done before the stop, was moved:\n%s", task)
	}
}

// S-0621 (ADR-0271 §2): each refusal writes nothing.
func TestFormulaStopRefuses(t *testing.T) {
	t.Parallel()
	stop := []string{"formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2"}
	for name, c := range map[string]struct {
		setup func(t *testing.T, root string)
		args  []string
		want  string
	}{
		"not an instance": {
			setup: func(t *testing.T, root string) {
				if _, errOut, code := runIn(t, root, "story", "new", "Plain"); code != 0 {
					t.Fatal(errOut)
				}
			},
			args: []string{"formula", "stop", "S-0002", "--step", "evaluate", "--exit", "2"},
			want: "S-0002 is no formula's instance",
		},
		"not in progress": {
			setup: func(t *testing.T, root string) {
				if _, errOut, code := runIn(t, root, "formula", "run", "publish-like", "--ready"); code != 0 {
					t.Fatal(errOut)
				}
			},
			args: []string{"formula", "stop", "S-0002", "--step", "evaluate", "--exit", "2"},
			want: "S-0002 is ready, not in progress",
		},
		"formula gone": {
			setup: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "design", "workflows", "formulas", "publish-like.md")); err != nil {
					t.Fatal(err)
				}
			},
			want: "which cannot be read",
		},
		"formula with a problem": {
			setup: func(t *testing.T, root string) {
				p := filepath.Join(root, "design", "workflows", "formulas", "publish-like.md")
				if err := os.WriteFile(p, []byte(strings.Replace(publishLike, "    stop: [2]\n", "    stop: [2, 2]\n", 1)), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: "the formula publish-like has 1 problem",
		},
		"unknown step": {
			args: []string{"formula", "stop", "S-0001", "--step", "release", "--exit", "2"},
			want: "give one of check, evaluate, note, publish",
		},
		"exit 0": {
			args: []string{"formula", "stop", "S-0001", "--step", "evaluate", "--exit", "0"},
			want: "stops the run only with exit 2, not 0",
		},
		"exit not listed": {
			args: []string{"formula", "stop", "S-0001", "--step", "evaluate", "--exit", "1"},
			want: "stops the run only with exit 2, not 1",
		},
		"step with no stop": {
			args: []string{"formula", "stop", "S-0001", "--step", "note", "--exit", "2"},
			want: "step note of the formula publish-like names no stop",
		},
		"step task not in progress": {
			setup: func(t *testing.T, root string) {
				p := filepath.Join(root, "design", "workflows", "formulas", "publish-like.md")
				if err := os.WriteFile(p, []byte(strings.Replace(publishLike, "    run: {flai: board}\n", "    run: {flai: board}\n    stop: [3]\n", 1)), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			args: []string{"formula", "stop", "S-0001", "--step", "note", "--exit", "3"},
			want: "T-0003, the task of step note, is backlog, not in progress",
		},
		"after not done": {
			setup: func(t *testing.T, root string) {
				p := filepath.Join(root, "design", "workflows", "formulas", "publish-like.md")
				if err := os.WriteFile(p, []byte(strings.Replace(publishLike, "    run: {flai: release}\n", "    run: {flai: release}\n    stop: [4]\n", 1)), 0o644); err != nil {
					t.Fatal(err)
				}
				for _, to := range []string{"ready", "in-progress"} {
					if _, errOut, code := runIn(t, root, "move", "T-0004", to); code != 0 {
						t.Fatal(errOut)
					}
				}
			},
			args: []string{"formula", "stop", "S-0001", "--step", "publish", "--exit", "4"},
			want: "T-0004, the task of step publish, waits for T-0002, which is not done",
		},
		"a later step done": {
			setup: func(t *testing.T, root string) {
				for _, to := range []string{"ready", "in-progress", "done"} {
					if _, errOut, code := runIn(t, root, "move", "T-0004", to); code != 0 {
						t.Fatal(errOut)
					}
				}
			},
			want: "T-0004, the task of step publish, is done, and its layer comes after step evaluate's",
		},
		"tasks added": {
			setup: func(t *testing.T, root string) {
				if _, errOut, code := runIn(t, root, "task", "new", "--story", "S-0001", "Another"); code != 0 {
					t.Fatal(errOut)
				}
			},
			want: "S-0001 has 5 tasks, not one per step of the formula publish-like (4)",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root := stopProject(t)
			if c.setup != nil {
				c.setup(t, root)
			}
			before := snapshotItems(t, root)
			args := c.args
			if args == nil {
				args = stop
			}
			out, errOut, code := runStdin(t, root, "output\n", args...)
			if code == 0 || !strings.Contains(errOut, c.want) {
				t.Errorf("formula stop: %d, want %q in:\n%s%s", code, c.want, out, errOut)
			}
			if after := snapshotItems(t, root); !maps.Equal(before, after) {
				t.Errorf("a refused stop wrote something")
			}
		})
	}
}

// S-0621 (ADR-0271 §2): an agent flai serve started stops the run of the
// instance it works alone, and a strategic role stops none.
func TestFormulaStopServedAgentStopsItsOwnInstance(t *testing.T) {
	root := stopProject(t)
	t.Setenv(guard.StartedByEnv, guard.StartedByServe)
	t.Setenv("FLAI_AGENT", "agent-S-0009")
	if _, errOut, code := runStdin(t, root, "x\n", "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2"); code == 0 ||
		!strings.Contains(errOut, "stops the run of the instance it works alone, and this one works S-0009, not S-0001") {
		t.Errorf("another story's agent: %d %s", code, errOut)
	}
	t.Setenv("FLAI_ROLE", "orchestrate")
	if _, errOut, code := runStdin(t, root, "x\n", "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2"); code == 0 ||
		!strings.Contains(errOut, "the orchestrate role does not stop a formula's run") {
		t.Errorf("the orchestrator: %d %s", code, errOut)
	}
	t.Setenv("FLAI_ROLE", "")
	t.Setenv("FLAI_AGENT", "agent-S-0001")
	if out, errOut, code := runStdin(t, root, "x\n", "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2"); code != 0 {
		t.Errorf("the instance's agent: %d %s %s", code, out, errOut)
	}
}

// snapshotItems are the contents of the item files and narratives under root,
// by path.
func snapshotItems(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, pattern := range []string{"wip/kanban/*/*.md", "wip/agents/*.md"} {
		matches, _ := filepath.Glob(filepath.Join(root, pattern))
		for _, m := range matches {
			data, err := os.ReadFile(m)
			if err != nil {
				t.Fatal(err)
			}
			out[m] = string(data)
		}
	}
	return out
}

// S-0621: the helpers keep every other byte of a body.
func TestFormulaStopHelpers(t *testing.T) {
	t.Parallel()
	body := "# S-0001 X\n\n## Acceptance criteria\n\n- [ ] one\n- [x] two\n\n## Tasks\n- T-0001 one\n\n## Notes\n\nMade by flai.\n"
	if got := skipCriteria(body, map[int]string{2: " (skipped)"}); got != strings.Replace(body, "- [x] two\n", "- [x] two (skipped)\n", 1) {
		t.Errorf("skipCriteria:\n%s", got)
	}
	if got := appendNotes(body, "Stopped."); got != body+"\nStopped.\n" {
		t.Errorf("appendNotes at the end:\n%q", got)
	}
	mid := "# T\n\n## Notes\n\nfirst\n\n## After\n\ntext\n"
	if got := appendNotes(mid, "second"); got != "# T\n\n## Notes\n\nfirst\n\nsecond\n\n## After\n\ntext\n" {
		t.Errorf("appendNotes before a later section:\n%q", got)
	}
	if got := appendNotes("# T\n", "x"); got != "# T\n\n## Notes\n\nx\n" {
		t.Errorf("appendNotes with no Notes:\n%q", got)
	}
	if got := fenced("a ``` b"); got != "````text\na ``` b\n````" {
		t.Errorf("fenced: %q", got)
	}
	if got := lastLines("1\n2\n3\n", 2); got != "2\n3" {
		t.Errorf("lastLines: %q", got)
	}
}
