package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// guardedRole is a role file a project defines: the reviewer, role review,
// which may move a story to ready only while its own setting is on.
const guardedRole = `---
title: Reviewer
updated: 2026-10-28
status: active
name: reviewer
role: review
definition:
  description: Reviews.
  tools: [Read, Edit, mcp__flai__thread_open, mcp__flai__item_move]
action: review
trigger: [asked]
item: none
tools:
  - thread_open
  - {call: item_move, when: to-backlog}
commands:
  - flai issue new
  - flai digest
files: ["{design}/reviews/"]
permissions:
  - setting: roles.reviewer.permissions.promote_to_ready
    tools: [{call: item_move, when: to-ready}]
  - setting: roles.reviewer.permissions.publish
    tools: [release_publish]
---

The reviewer reviews.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// guardedProject is a temporary project whose design folder is plans, with
// files, by name, under its workflows/roles, and block added to its
// manifest.
func guardedProject(t *testing.T, files map[string]string, block string) string {
	t.Helper()
	root := tempProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	data, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, []byte(strings.Replace(string(data), "design: design", "design: plans", 1)+block), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "plans", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// S-0416, ADR-0207: with FLAI_ROLE naming a role a project's file defines,
// flai guard holds the session's own calls to that file: from the
// manifest's design folder, under the role's own roles.<name>.permissions,
// with its command entries matched against flai's command tree.
func TestGuardHoldsADefinedRoleToItsFile(t *testing.T) {
	t.Setenv("FLAI_ROLE", "review")
	root := guardedProject(t, map[string]string{"reviewer.md": guardedRole}, "roles:\n  reviewer:\n    permissions:\n      promote_to_ready: true\norchestration:\n  permissions:\n    publish: true\n")
	edit := func(file string) string {
		in, _ := json.Marshal(map[string]any{"tool_name": "Write", "tool_input": map[string]string{"file_path": file, "content": ""}})
		return string(in)
	}
	for _, c := range []struct {
		in   string
		code int
		err  string
	}{
		{`{"tool_name":"mcp__flai__board","tool_input":{}}`, 0, ""},
		{`{"tool_name":"mcp__flai__thread_open","tool_input":{"on":"S-1"}}`, 0, ""},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"backlog"}}`, 0, ""},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"ready"}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"flai board --json && flai issue new x --class defect && flai digest"}}`, 0, ""},
		{edit(filepath.Join(root, "plans", "reviews", "r.md")), 0, ""},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"review"}}`, 2, "the role reviewer cannot move S-1 to review: it may make it only under to-backlog"},
		{`{"tool_name":"mcp__flai__release_publish","tool_input":{"reason":"r"}}`, 2, "the role reviewer cannot call release_publish: it needs roles.reviewer.permissions.publish, which is off. Ask the operator with thread_open on the item if it needs doing (plans/workflows/roles/reviewer.md, ADR-0207)."},
		{`{"tool_name":"Bash","tool_input":{"command":"flai digest harvest"}}`, 2, `the role reviewer cannot run "flai digest harvest": its role file lists no such call`},
		{edit(filepath.Join(root, "design", "reviews", "r.md")), 2, "the role reviewer cannot use Write on " + filepath.Join(root, "design", "reviews", "r.md") + ": it edits only a file under plans/reviews/."},
		{`{"tool_name":"mcp__flai__item_new","tool_input":{"type":"epic"},"agent_type":"explorer","agent_id":"a1"}`, 2, "a sub-agent (explorer) cannot call item_new"},
	} {
		_, errOut, code := runStdin(t, root, c.in, "guard")
		if code != c.code || (c.err == "" && errOut != "") || !strings.Contains(errOut, c.err) {
			t.Errorf("%s: code %d, stderr %q", c.in, code, errOut)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "wip", "agents", "orchestrator.md")); err == nil {
		t.Error("a defined role's refusal was logged as the orchestrator's")
	}
}

// S-0416: a role whose file is missing or does not parse, or whose project
// cannot be read, reads and is refused every other call, naming the file or
// the folder and flai check.
func TestGuardRefusesADefinedRoleWithoutAFile(t *testing.T) {
	t.Setenv("FLAI_ROLE", "scout")
	write := `{"tool_name":"mcp__flai__thread_open","tool_input":{"on":"S-1"}}`
	read := `{"tool_name":"Bash","tool_input":{"command":"flai board --json"}}`
	for _, c := range []struct {
		name string
		dir  string
		err  string
	}{
		{"no file", guardedProject(t, map[string]string{"reviewer.md": guardedRole}, ""), "no role file under plans/workflows/roles has role scout; run flai check"},
		{"a file that does not parse", guardedProject(t, map[string]string{"reviewer.md": guardedRole, "scout.md": "---\nname: [scout\n---\n"}, ""), "plans/workflows/roles/scout.md does not parse: "},
		{"no project", t.TempDir(), "the project could not be read for the role file of scout"},
	} {
		_, errOut, code := runStdin(t, c.dir, write, "guard")
		if code != 2 || !strings.Contains(errOut, "the role scout cannot call thread_open: ") || !strings.Contains(errOut, c.err) || !strings.Contains(errOut, "run flai check") {
			t.Errorf("%s: code %d, stderr %q", c.name, code, errOut)
		}
		if _, errOut, code := runStdin(t, c.dir, read, "guard"); code != 0 || strings.Contains(errOut, "cannot") {
			t.Errorf("%s: a read: code %d, stderr %q", c.name, code, errOut)
		}
	}
}
