package workitem_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0181: release.FieldsFile lists the front-matter fields this flai reads.
// Publishing a flai release that changed it raises the project's
// flai.minimum, so the list must be the code's: when a field is added or
// removed, change the file with it.
func TestFieldsFileIsTheCode(t *testing.T) {
	data, err := os.ReadFile("front-matter-fields.txt")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" && !strings.HasPrefix(line, "#") {
			got = append(got, line)
		}
	}
	want := []string{
		"item: " + strings.Join(workitem.KnownFields(workitem.Item{}), " "),
		"thread: " + strings.Join(workitem.KnownFields(threads.Thread{}), " "),
		"issue: " + strings.Join(workitem.KnownFields(issues.Issue{}), " "),
		// a story's agent block (S-0189): an older flai would drop a key of it
		// that it does not know when it rewrote the story
		"item.agent: " + strings.Join(workitem.KnownFields(manifest.Agent{}), " "),
		"item.agent.roles: " + strings.Join(workitem.KnownFields(manifest.Role{}), " "),
		// the planning blocks (S-0199), likewise
		"item.cost_of_delay: " + strings.Join(workitem.KnownFields(workitem.CostOfDelay{}), " "),
		"item.cost_of_delay.inputs: " + strings.Join(workitem.KnownFields(workitem.CostInputs{}), " "),
		"item.forecast: " + strings.Join(workitem.KnownFields(workitem.Forecast{}), " "),
		"item.finalized: " + strings.Join(workitem.KnownFields(workitem.Finalized{}), " "),
		// the usage block, its models, and its strategic entries (S-0225),
		// likewise
		"item.usage: " + strings.Join(workitem.KnownFields(usage.Usage{}), " "),
		"item.usage.models: " + strings.Join(workitem.KnownFields(usage.Model{}), " "),
		"item.usage.strategic: " + strings.Join(workitem.KnownFields(usage.Strategic{}), " "),
		// and the counts of each day of its turns (S-0293)
		"item.usage.turns: " + strings.Join(workitem.KnownFields(usage.TurnDay{}), " "),
		// the types each type-restricted field is valid on (S-0176): an older
		// flai refuses an item that carries one on a type it does not allow
		"item.types: " + workitem.FieldTypes(),
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("%s is not the code's fields; write these lines:\n%s", release.FieldsFile, strings.Join(want, "\n"))
	}
	if !strings.HasSuffix(release.FieldsFile, "internal/workitem/front-matter-fields.txt") {
		t.Errorf("release.FieldsFile is %s", release.FieldsFile)
	}
}

const formulaStory = `---
id: S-0001
type: story
nature: feature
title: Release 1.4.0
status: backlog
owner: alex
created: 2026-11-02T09:00:00Z
updated: 2026-11-02T09:00:00Z
transitions: []
tags: []
draft: true
formula: release
---

# S-0001: Release 1.4.0
`

// S-0421: a story names the formula flai formula run made it from, and reads
// and saves back byte for byte, written after the planning fields, where an
// older flai writes back a key it does not know; only a story carries one.
func TestFormulaRoundTrips(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "wip", "kanban", "stories", "S-0001-release-1-4-0.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(formulaStory), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := workitem.ReadItem(path)
	if err != nil {
		t.Fatal(err)
	}
	if it.Formula != "release" || len(it.Unknown) > 0 {
		t.Errorf("read formula %q, unknown %v", it.Formula, it.Unknown)
	}
	if err := it.Validate(); err != nil {
		t.Errorf("a story with a formula: %v", err)
	}
	if err := r.Save(it); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); string(data) != formulaStory {
		t.Errorf("saved as:\n%s\nwant:\n%s", data, formulaStory)
	}

	task, err := workitem.ParseItem(strings.Replace(strings.Replace(formulaStory, "S-0001", "T-0001", 2), "type: story\n", "type: task\nparent: S-0001\n", 1))
	if err != nil {
		t.Fatal(err)
	}
	if err := task.Validate(); err == nil || !strings.Contains(err.Error(), "only a story carries a formula") {
		t.Errorf("a task with a formula: %v", err)
	}
}
