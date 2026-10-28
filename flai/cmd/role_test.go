package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// reviewerRole is a role file with a condition on a trigger and a call, and
// a permission.
const reviewerRole = `---
title: Reviewer
updated: 2026-10-28
status: active
name: reviewer
role: review
definition:
  description: Reviews what is asked.
  tools: [Read, mcp__flai__board, mcp__flai__inbox]
action: review
trigger:
  - asked
  - {schedule: planning.schedule, when: each-ready-story}
item: [story]
topics: [analysis]
charge: item
tools: [inbox, {call: thread_reply, when: own-thread}]
commands: [flai issue new]
files: ["{design}/reviews/"]
permissions:
  - setting: roles.reviewer.permissions.publish
    tools: [release_publish]
---

You are the reviewer.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// rolesProject is a project whose design folder's role files are files.
func rolesProject(t *testing.T, files map[string]string) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	dir := filepath.Join(root, "design", "workflows", "roles")
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

// S-0414: flai role list prints each role that loads with its file, verb,
// action, triggers, items, topics, charge, calls, and permissions, a call or
// trigger under a condition as (when condition), and says a file that does
// not load on stderr; --json gives the roles with their paths.
func TestRoleList(t *testing.T) {
	root := rolesProject(t, map[string]string{"README.md": "# Roles\n", "reviewer.md": reviewerRole, "broken.md": "no front matter\n"})
	out, errOut, code := runIn(t, root, "role", "list")
	if code != 0 {
		t.Fatalf("role list: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{
		"reviewer  review  design/workflows/roles/reviewer.md\n",
		"  action:      review\n",
		"  trigger:     asked, schedule planning.schedule (when each-ready-story)\n",
		"  item:        story\n",
		"  topics:      analysis\n",
		"  charge:      item\n",
		"  tools:       inbox, thread_reply (when own-thread)\n",
		"  commands:    flai issue new\n",
		"  files:       {design}/reviews/\n",
		"  permissions:\n    roles.reviewer.permissions.publish → release_publish\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("role list lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "broken") || !strings.Contains(errOut, "design/workflows/roles/broken.md:1: roles.file:") {
		t.Errorf("the file that does not load is left out and said on stderr: %s / %s", out, errOut)
	}
	js, _, code := runIn(t, root, "role", "list", "--json")
	var got []map[string]any
	if err := json.Unmarshal([]byte(js), &got); code != 0 || err != nil || len(got) != 1 {
		t.Fatalf("role list --json: %d %v %s", code, err, js)
	}
	if got[0]["path"] != "design/workflows/roles/reviewer.md" || got[0]["name"] != "reviewer" || got[0]["charge"] != "item" {
		t.Errorf("role list --json: %v", got[0])
	}
	empty := rolesProject(t, nil)
	if out, _, _ := runIn(t, empty, "role", "list"); out != "no roles\n" {
		t.Errorf("no roles: %q", out)
	}
	if js, _, _ := runIn(t, empty, "role", "list", "--json"); strings.TrimSpace(js) != "[]" {
		t.Errorf("no roles, --json: %q", js)
	}
}

// S-0414: flai role show prints a role, its definition, and its body; --json
// gives the role; an unknown name is an error naming the roles there are.
func TestRoleShow(t *testing.T) {
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	out, errOut, code := runIn(t, root, "role", "show", "reviewer")
	if code != 0 {
		t.Fatalf("role show: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{
		"reviewer  review  design/workflows/roles/reviewer.md\n",
		"  tools:       inbox, thread_reply (when own-thread)\n",
		"  title:       Reviewer\n",
		"  works_story: false\n",
		"  description: Reviews what is asked.\n",
		"  tools:       Read, mcp__flai__board, mcp__flai__inbox\n",
		"  model:       inherit\n",
		"\nYou are the reviewer.\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("role show lacks %q:\n%s", want, out)
		}
	}
	js, _, _ := runIn(t, root, "role", "show", "reviewer", "--json")
	var r roles.Role
	if err := json.Unmarshal([]byte(js), &r); err != nil || r.Path != "design/workflows/roles/reviewer.md" || r.Definition.Model != "inherit" || len(r.Permissions) != 1 || r.Tools[1].When != "own-thread" {
		t.Errorf("role show --json: %v %+v", err, r)
	}
	if _, errOut, code := runIn(t, root, "role", "show", "planner"); code == 0 || !strings.Contains(errOut, "no role planner; the roles are reviewer") {
		t.Errorf("an unknown role: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, rolesProject(t, nil), "role", "show", "planner"); code == 0 || !strings.Contains(errOut, "no role planner: the project has no role that loads") {
		t.Errorf("no roles: %d %s", code, errOut)
	}
}

// S-0414: flai check holds role files to the guard's reads but not to the
// MCP server's tools and the command tree, which import it; this holds the
// shipped role files to both, so that a renamed command or MCP tool that a
// baseline role names fails here.
func TestShippedRolesNameFlaisToolsAndCommands(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo's role files")
	}
	repo, err := workitem.Open("../..")
	if err != nil {
		t.Skip("monorepo not present")
	}
	list, problems, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil || len(problems) != 0 || len(list) == 0 {
		t.Fatalf("the shipped role files: %d roles, %v, %v", len(list), problems, err)
	}
	v := check.RolesVocabulary()
	v.Commands = commandPaths(newRootCmdWith(&app{out: &bytes.Buffer{}, errOut: &bytes.Buffer{}}))
	v.Tools = mcpToolNames(t)
	if !slices.Contains(v.Commands, "story new") || !slices.Contains(v.Tools, "inbox") {
		t.Fatalf("the vocabulary misses flai story new or inbox: %v %v", v.Commands, v.Tools)
	}
	for _, p := range roles.Validate(list, repo.Manifest.Layout["design"], v) {
		t.Error(p)
	}
}

// commandPaths are every command under root, as "edit" and "story new".
func commandPaths(root *cobra.Command) []string {
	var out []string
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			out = append(out, strings.TrimPrefix(sub.CommandPath(), root.Name()+" "))
			walk(sub)
		}
	}
	walk(root)
	return out
}

// mcpToolNames are the tools flai mcp serves a project: the server adds
// every tool whatever its options.
func mcpToolNames(t *testing.T) []string {
	t.Helper()
	root := rolesProject(t, nil)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct, st := mcp.NewInMemoryTransports()
	if _, err := mcpserver.New(mcpserver.Options{Repo: repo, Agent: "tester", Version: "test"}).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "tester", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()
	res, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, tool := range res.Tools {
		out = append(out, tool.Name)
	}
	return out
}
