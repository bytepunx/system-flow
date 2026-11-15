package serve

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// renamedFlow is the board the serve tests run against a workflow other than
// the baseline (S-0743): in-progress is doing and review is awaiting, neither
// with an alias, and testing, of the active kind, sits between them. A pass
// that reads a state's name in place of its kind or a transition's effect
// deploys, pushes, queues, or sends back the wrong story under it.
const renamedFlow = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: doing, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: awaiting, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: doing, effects: [pull]}
  - {from: doing, to: testing}
  - {from: testing, to: awaiting}
  - {from: awaiting, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: awaiting, to: done, types: [task]}
  - {from: doing, to: done, types: [task]}
  - {from: awaiting, to: doing, effects: [send-back]}
  - {from: testing, to: doing}
  - {from: doing, to: ready}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, doing, testing], to: cancelled, effects: [cancel]}
  - {from: awaiting, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`

// renamedWorkflow is renamedFlow parsed, for Options.Workflow.
func renamedWorkflow(t *testing.T) *workflow.Workflow {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", renamedFlow)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %s", p)
	}
	return &w
}

// S-0743: openRepo gives the repository the options' workflow, and the
// baseline when they give none, and flow answers the same.
func TestOpenRepoRunsTheOptionsWorkflow(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := openRepo(Options{}, root)
	if err != nil {
		t.Fatal(err)
	}
	if repo.Workflow != nil || !repo.Flow().AwaitsAcceptance(workflow.TypeStory, "review") || !(Options{}).flow().AwaitsAcceptance(workflow.TypeStory, "review") {
		t.Error("with no workflow in the options, the repository does not run the baseline")
	}
	o := Options{Workflow: renamedWorkflow(t)}
	repo, err = openRepo(o, root)
	if err != nil {
		t.Fatal(err)
	}
	if repo.Workflow != o.Workflow || !repo.Flow().AwaitsAcceptance(workflow.TypeStory, "awaiting") || !o.flow().AwaitsAcceptance(workflow.TypeStory, "awaiting") {
		t.Error("the repository does not run the options' workflow")
	}
	if _, err := openRepo(o, t.TempDir()); err == nil {
		t.Error("openRepo opened a folder with no manifest")
	}
}

// S-0755, ADR-0307 §3: with no workflow in the options, openRepo runs the
// project's own workflow file; a file flai cannot run refuses the open and is
// logged once for each distinct finding, however many looks open it; and an
// open that succeeds again lets the same finding be told again.
func TestOpenRepoRunsTheProjectsFileAndLogsARefusalOnce(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	board := filepath.Join(root, "design", "workflows", "board.md")
	if err := os.MkdirAll(filepath.Dir(board), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(text string) {
		t.Helper()
		if err := os.WriteFile(board, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var logged bytes.Buffer
	o := Options{Logger: slog.New(slog.NewTextHandler(&logged, nil))}
	write(renamedFlow)
	repo, err := openRepo(o, root)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(repo.Flow().Names(), "testing") || !repo.Flow().AwaitsAcceptance(workflow.TypeStory, "awaiting") {
		t.Errorf("openRepo runs %s; want the project's file", repo.Flow().Names())
	}
	refusals := func() int { return strings.Count(logged.String(), "level=WARN") }
	noPull := strings.Replace(renamedFlow, "{from: ready, to: doing, effects: [pull]}", "{from: ready, to: doing}", 1)
	write(noPull)
	for range 3 {
		if _, err := openRepo(o, root); err == nil || !strings.Contains(err.Error(), workflow.CodeNoPull) {
			t.Fatalf("a file with no pull opens with %v; want the %s refusal", err, workflow.CodeNoPull)
		}
	}
	if n := refusals(); n != 1 || !strings.Contains(logged.String(), "root="+root) || !strings.Contains(logged.String(), workflow.CodeNoPull) {
		t.Errorf("three looks at one finding logged %d refusals:\n%s\nwant one naming the root and the finding", n, logged.String())
	}
	write(strings.Replace(renamedFlow, ", effects: [accept]", "", 1))
	if _, err := openRepo(o, root); err == nil || !strings.Contains(err.Error(), workflow.CodeNoAccept) {
		t.Fatalf("a file with no accept opens with %v; want the %s refusal", err, workflow.CodeNoAccept)
	}
	if n := refusals(); n != 2 {
		t.Errorf("another finding logged %d refusals in all; want 2", n)
	}
	write(renamedFlow)
	if _, err := openRepo(o, root); err != nil {
		t.Fatalf("after the fix: %v", err)
	}
	write(noPull)
	if _, err := openRepo(o, root); err == nil {
		t.Fatal("the file broken again opened")
	}
	if n := refusals(); n != 3 {
		t.Errorf("the file broken again after a fix logged %d refusals in all; want 3", n)
	}
}
