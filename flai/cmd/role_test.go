package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/serve"
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

// S-0417: flai role run starts a role the project defines, on an item its
// file takes, while the host action its file names is on, records the run
// under roles by name and item where flai serve tracks agents, and says why
// when it will not; role.run does the same for a dashboard, and both are
// journalled under the role's action.
func TestRoleRunStartsADefinedRole(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("starts true as the agent's command")
	}
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	cfg := os.Getenv("FLAI_CONFIG")
	inProcess(t)
	for _, title := range []string{"Slice", "Gone", "Other"} {
		runIn(t, root, "story", "new", title)
	}
	runIn(t, root, "move", "S-0002", "cancelled", "--reason", "not needed")
	refusals := []struct {
		args []string
		want string
	}{
		{[]string{"role", "run", "reviewer", "S-1"}, "rule: the review host action, which the role reviewer's file design/workflows/roles/reviewer.md names, is off for this project: flai serve enable review"},
		{[]string{"role", "run", "planner", "S-1"}, "rule: planner is a built-in role, which flai starts with flai plan <item>"},
		{[]string{"role", "run", "auditor"}, "rule: the project defines no role auditor; the roles it defines are reviewer"},
		{[]string{"role", "run", "reviewer", "S-1", "--input", "focus"}, "is not name=value: give each input as its name"},
	}
	for _, r := range refusals {
		if _, errOut, code := runIn(t, root, r.args...); code == 0 || !strings.Contains(errOut, r.want) {
			t.Errorf("%v: %d %s", r.args, code, errOut)
		}
	}
	out, _, code := runIn(t, root, "hostapi", "role.run", `{"name":"reviewer","item":"S-0001","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a60"}`)
	if code == 0 || !strings.Contains(out, `the host action \"review\" is not enabled`) || !strings.Contains(out, "flai serve enable review") {
		t.Errorf("role.run, action off: %d %s", code, out)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "review"); code != 0 {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	for _, r := range []struct {
		args []string
		want string
	}{
		{[]string{"role", "run", "reviewer"}, "rule: the role reviewer runs on an item, a story, as its file design/workflows/roles/reviewer.md says; give it one"},
		{[]string{"role", "run", "reviewer", "S-2"}, "rule: S-0002 is cancelled; the role reviewer runs on an item still open"},
		{[]string{"role", "run", "reviewer", "S-1"}, "rule: the role reviewer's agent names no harness"},
	} {
		if _, errOut, code := runIn(t, root, r.args...); code == 0 || !strings.Contains(errOut, r.want) {
			t.Errorf("%v: %d %s", r.args, code, errOut)
		}
	}
	runIn(t, root, "serve", "agent", "set", "--", "true")
	out, errOut, code := runIn(t, root, "role", "run", "reviewer", "S-1", "--input", "focus=docs only")
	if code != 0 || !strings.HasPrefix(out, "started true (command) to run the role reviewer on S-0001 as reviewer-S-0001 (pid ") {
		t.Fatalf("role run: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "flai serve is not running") {
		t.Errorf("not told that no flai serve settles the run: %s", errOut)
	}
	out, _, code = runIn(t, root, "hostapi", "role.run", `{"name":"reviewer","item":"S-0003","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a61"}`)
	var started struct {
		Data struct {
			Role, Item, Action, Agent, Log, Session, Trigger string
			PID                                              int
		}
	}
	if err := json.Unmarshal([]byte(out), &started); code != 0 || err != nil || started.Data.Role != "reviewer" || started.Data.Item != "S-0003" || started.Data.Action != "review" ||
		started.Data.Agent != "reviewer-S-0003" || started.Data.PID == 0 || started.Data.Log == "" || started.Data.Trigger != serve.TriggerAsked {
		t.Fatalf("role.run: %d %s", code, out)
	}
	var all map[string]serve.AgentState
	data, _ := os.ReadFile(filepath.Join(string(serve.DirFor(cfg)), "agents.json"))
	if err := json.Unmarshal(data, &all); err != nil {
		t.Fatal(err)
	}
	for _, st := range all {
		if st.Roles["reviewer"]["S-0001"] == nil || st.Roles["reviewer"]["S-0003"] == nil || len(st.Stories) != 0 || len(st.Plans) != 0 {
			t.Errorf("state = %+v, want both runs under roles by name and item, and no story's or planner's", st)
		}
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "role.run"`, `"action": "review"`, `"outcome": "disabled"`, "started true to run the role reviewer on S-0003 as reviewer-S-0003"} {
		if !strings.Contains(js, want) {
			t.Errorf("the journal has no %q: %s", want, js)
		}
	}
}

// S-0417: flai activity log logs an activity of a role the project defines
// in its activity document, as the MCP tool activity_log does, and refuses a
// kind that is no strategic agent, naming those there are.
func TestActivityLogTakesADefinedRole(t *testing.T) {
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	runIn(t, root, "story", "new", "Slice")
	out, errOut, code := runIn(t, root, "activity", "log", "reviewer", "  reviewed\n S-0001 ", "--item", "S-1", "--json")
	var logged struct {
		Entry struct {
			Summary string
			Items   []string
		}
		Activity struct {
			Kind           string
			TasksCompleted int `json:"tasks_completed"`
		}
	}
	if err := json.Unmarshal([]byte(out), &logged); code != 0 || err != nil || logged.Entry.Summary != "reviewed S-0001" || !slices.Equal(logged.Entry.Items, []string{"S-0001"}) ||
		logged.Activity.Kind != "reviewer" || logged.Activity.TasksCompleted != 1 {
		t.Fatalf("activity log: %d %s %s", code, out, errOut)
	}
	if doc, err := os.ReadFile(filepath.Join(root, "wip", "agents", "reviewer.md")); err != nil || !strings.Contains(string(doc), "reviewed S-0001") {
		t.Errorf("the reviewer's activity document: %s %v", doc, err)
	}
	if out, _, code := runIn(t, root, "activity", "log", "reviewer", "checked it again"); code != 0 || !strings.HasPrefix(out, "logged in wip/agents/reviewer.md at ") || !strings.Contains(out, "2 activities") {
		t.Errorf("activity log, as text: %d %s", code, out)
	}
	if _, errOut, code := runIn(t, root, "activity", "log", "auditor", "audited"); code == 0 || !strings.Contains(errOut, "is not a strategic agent: use one of planner, orchestrator, analyzer, curator, reviewer") {
		t.Errorf("an unknown kind: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "activity", "log", "reviewer", "  "); code == 0 || !strings.Contains(errOut, "an activity needs a summary") {
		t.Errorf("no summary: %d %s", code, errOut)
	}
}

// S-0417: flai mcp's role_run starts a role the project defines as flai
// role run does, under the host action its file names, passes its refusals
// through, journals the agent that asked under that action, and refuses an
// agent flai serve started.
func TestTheMCPServerStartsADefinedRoleUnderItsAction(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("starts true as the agent's command")
	}
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	t.Setenv("FLAI_STARTED_BY", "") // the operator's own agent, whoever runs the tests
	runIn(t, root, "story", "new", "Slice")
	a := &app{out: &bytes.Buffer{}, errOut: &bytes.Buffer{}, cwd: root, clock: func() time.Time { return time.Date(2026, 10, 28, 10, 0, 0, 0, time.UTC) }}
	ctx := context.Background()
	focus := []harness.Input{{Name: "focus", Value: "docs"}}
	if _, err := a.mcpRoleRun(ctx, root, "reviewer", "S-1", focus, "agent-ops"); err == nil || !strings.Contains(err.Error(), "the review host action, which the role reviewer's file design/workflows/roles/reviewer.md names, is off for this project: flai serve enable review") {
		t.Errorf("action off: %v", err)
	}
	if _, err := a.mcpRoleRun(ctx, root, "auditor", "", nil, "agent-ops"); err == nil || !strings.Contains(err.Error(), "the project defines no role auditor; the roles it defines are reviewer") {
		t.Errorf("no such role: %v", err)
	}
	runIn(t, root, "serve", "enable", "review")
	runIn(t, root, "serve", "agent", "set", "--", "true")
	got, err := a.mcpRoleRun(ctx, root, "reviewer", "S-1", focus, "agent-ops")
	if err != nil || got.Role != "reviewer" || got.Item != "S-0001" || got.Action != "review" || got.Agent != "reviewer-S-0001" || got.PID == 0 || got.Command != "true" || got.Log == "" || got.Trigger != serve.TriggerAsked {
		t.Fatalf("role_run: %+v %v", got, err)
	}
	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	if _, err := a.mcpRoleRun(ctx, root, "reviewer", "S-1", nil, "agent-S-0001"); err == nil || !strings.Contains(err.Error(), "an agent flai serve started does not start a role") {
		t.Errorf("an agent flai serve started: %v", err)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "mcp.role_run"`, `"action": "review"`, `"by": "agent-ops"`, `"outcome": "disabled"`,
		"agent-ops asked to run the role reviewer on S-1: the review host action", "agent-ops asked to run the role reviewer on S-1: started true as reviewer-S-0001"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
	if strings.Contains(js, "auditor") {
		t.Errorf("a role the project does not define is journalled under no action: %s", js)
	}
}
