package roles

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const orchestratorFile = `---
title: Orchestrator
updated: 2026-10-28
status: active
name: orchestrator
role: orchestrate
definition:
  description: Keeps the work of this project moving.
  tools: [Read, Grep, Bash, mcp__flai__prime, mcp__flai__inbox, mcp__flai__item_edit, mcp__flai__item_move, mcp__flai__thread_reply, mcp__flai__thread_resolve, mcp__flai__plan, mcp__flai__release_publish, mcp__flai__queue_add]
  model: inherit
action: orchestrate
trigger: [continuous]
item: none
topics: [orchestration]
works_story: false
charge: named
tools:
  - activity_log
  - inbox
  - thread_open
  - wait_for_events
  - {call: thread_reply, when: own-thread}
  - {call: thread_resolve, when: own-thread}
commands:
  - flai issue new
  - flai issue bump
  - flai queue
  - flai thread new
  - {call: flai thread reply, when: own-thread}
  - {call: flai thread resolve, when: own-thread}
files: []
permissions:
  - setting: orchestration.permissions.plan_backlog_epics
    tools: [{call: plan, when: plan-epic}]
    commands: [{call: flai plan, when: plan-epic}]
  - setting: orchestration.permissions.plan_backlog_stories
    tools:
      - {call: plan, when: plan-story}
      - {call: item_edit, when: cost-inputs}
      - {call: thread_reply, when: planner-thread}
      - {call: thread_resolve, when: planner-thread}
    commands:
      - {call: flai plan, when: plan-story}
      - {call: flai edit, when: cost-inputs}
      - {call: flai thread reply, when: planner-thread}
      - {call: flai thread resolve, when: planner-thread}
  - setting: orchestration.permissions.finalize_drafts
    tools: [{call: item_edit, when: finalize-draft}]
    commands: [{call: flai edit, when: finalize-draft}]
  - setting: orchestration.permissions.promote_to_ready
    tools: [{call: item_move, when: to-ready}]
    commands: [{call: flai move, when: to-ready}]
  - setting: orchestration.permissions.order_ready
    commands: [{call: flai order, when: by-policy}]
  - setting: orchestration.permissions.answer_threads
    tools: [{call: thread_reply, when: answer-threads}]
    commands: [{call: flai thread reply, when: answer-threads}]
  - setting: orchestration.permissions.accept_reviews
    tools: [queue_add]
    commands:
      - {call: flai accept, when: accept-as-self}
      - {call: flai move, when: accept-as-self}
      - {call: flai queue add, when: enqueue-as-self}
  - setting: orchestration.permissions.publish
    tools: [release_publish]
  - setting: orchestration.permissions.edit_subprojects
    commands: [flai subproject add, flai subproject set, flai subproject remove]
---

You are the orchestrator.

<!-- system-flow:end-of-baseline -->

## Project additions
`

const plannerFile = `---
title: Planner
updated: 2026-10-28
status: active
name: planner
role: plan
definition:
  description: Plans one epic or one story.
  tools: [Read, Grep, Bash, mcp__flai__prime, mcp__flai__item_new, mcp__flai__item_edit, mcp__flai__item_move]
action: plan
trigger:
  - asked
  - {event: edited, when: stale-figures}
  - {event: accepted, when: replan-agent}
  - {event: cancelled, when: replan-agent}
  - {event: reordered, when: replan-agent}
  - {schedule: planning.schedule, when: each-ready-story}
item: [epic, story]
topics: [planning]
charge: item
tools:
  - activity_log
  - inbox
  - {call: item_new, when: draft-only}
  - {call: item_edit, when: ungated-edit}
  - {call: item_move, when: to-backlog}
  - thread_open
  - thread_reply
  - wait_for_events
commands:
  - {call: flai story new, when: draft-only}
  - flai epic new
  - flai task new
  - {call: flai edit, when: ungated-edit}
  - flai touches *
  - flai thread new
  - flai thread reply
  - flai issue new
  - flai issue bump
  - {call: flai move, when: to-backlog}
---

You are the planner.

<!-- system-flow:end-of-baseline -->

## Project additions

Plan small.
`

// analyzerFile edits under one folder, so its definition lists Edit.
const analyzerFile = `---
title: Analyzer
updated: 2026-10-28
status: active
name: analyzer
role: analyze
definition:
  description: Analyzes the project.
  tools: [Read, Edit, Write, mcp__flai__issue_new]
action: analyze
trigger: [asked, {schedule: analysis.schedule}]
item: none
topics: [analysis]
charge: issues
tools: [issue_new]
files: ['{design}/analysis/']
---

You are the analyzer.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// vocab names every MCP tool and command the fixtures use.
var vocab = Vocabulary{
	Tools: []string{"activity_log", "inbox", "thread_open", "thread_reply", "thread_resolve", "wait_for_events", "plan", "item_edit", "item_move", "item_new",
		"queue_add", "release_publish", "task_done", "issue_new", "issue_bump", "subproject_edit", "shared_paths_edit", "story_start", "prime", "board"},
	Commands: []string{"issue new", "issue bump", "queue", "queue add", "queue hold", "queue remove", "queue resume", "thread new", "thread reply", "thread resolve", "thread confirm",
		"plan", "edit", "move", "order", "accept", "subproject add", "subproject set", "subproject remove", "story new", "story start", "epic new", "task new",
		"touches", "touches suggest", "release", "push"},
	Reads: []string{"prime", "board", "item_get", "doc_get"},
}

// project writes files under a temporary project's design/workflows/roles
// and returns its root.
func project(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
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

func load(t *testing.T, files map[string]string) ([]Role, []Problem) {
	t.Helper()
	rs, problems, err := Load(project(t, files), "design")
	if err != nil {
		t.Fatal(err)
	}
	return rs, problems
}

func TestValidRolesLoadWithEveryField(t *testing.T) {
	rs, problems := load(t, map[string]string{
		"orchestrator.md": orchestratorFile, "planner.md": plannerFile, "analyzer.md": analyzerFile,
		"README.md": "# Roles\n\nNo front matter here.\n", "notes.txt": "not a role",
	})
	if len(problems) != 0 {
		t.Fatalf("load problems: %v", problems)
	}
	if got := Validate(rs, "design", vocab); len(got) != 0 {
		t.Fatalf("Validate = %v, want none", got)
	}
	if len(rs) != 3 || rs[0].Name != "analyzer" || rs[1].Name != "orchestrator" || rs[2].Name != "planner" {
		t.Fatalf("roles = %v, want analyzer, orchestrator, planner sorted by path", rs)
	}
	o := rs[1]
	if o.Path != "design/workflows/roles/orchestrator.md" || o.Title != "Orchestrator" || o.Updated != "2026-10-28" || o.Status != "active" ||
		o.Role != "orchestrate" || o.Action != "orchestrate" || o.Charge != "named" || o.WorksStory {
		t.Errorf("orchestrator scalars = %+v", o)
	}
	if o.Definition.Description != "Keeps the work of this project moving." || o.Definition.Model != "inherit" || len(o.Definition.Tools) != 12 {
		t.Errorf("orchestrator definition = %+v", o.Definition)
	}
	if !slices.Equal(o.Trigger, []Trigger{{Kind: TriggerContinuous, Line: 12}}) {
		t.Errorf("orchestrator trigger = %+v", o.Trigger)
	}
	if len(o.Item) != 0 || !slices.Equal(o.Topics, []string{"orchestration"}) || len(o.Files) != 0 {
		t.Errorf("orchestrator item %v topics %v files %v", o.Item, o.Topics, o.Files)
	}
	if got := o.Tools[4]; got != (Call{Call: "thread_reply", When: "own-thread", Line: 22}) {
		t.Errorf("orchestrator tools[4] = %+v", got)
	}
	if got := o.Commands[2]; got != (Call{Call: "flai queue", Line: 27}) {
		t.Errorf("orchestrator commands[2] = %+v", got)
	}
	if len(o.Permissions) != 9 {
		t.Fatalf("orchestrator has %d permissions, want 9", len(o.Permissions))
	}
	pub := o.Permissions[7]
	if pub.Setting != "orchestration.permissions.publish" || pub.Line != 64 || !slices.Equal(pub.Tools, []Call{{Call: "release_publish", Line: 65}}) || len(pub.Commands) != 0 {
		t.Errorf("publish permission = %+v", pub)
	}
	if !strings.Contains(o.Baseline, "You are the orchestrator.") || strings.Contains(o.Baseline, Marker) {
		t.Errorf("baseline = %q", o.Baseline)
	}
	if !strings.HasPrefix(o.Additions, Marker) || !strings.Contains(o.Additions, "## Project additions") {
		t.Errorf("additions = %q", o.Additions)
	}

	p := rs[2]
	if !slices.Equal(p.Item, Kinds{"epic", "story"}) || p.Charge != "item" || p.Definition.Model != "inherit" {
		t.Errorf("planner item %v charge %s model %s", p.Item, p.Charge, p.Definition.Model)
	}
	wantTriggers := []Trigger{
		{Kind: TriggerAsked, Line: 12},
		{Kind: TriggerEvent, Event: "edited", When: "stale-figures", Line: 13},
		{Kind: TriggerEvent, Event: "accepted", When: "replan-agent", Line: 14},
		{Kind: TriggerEvent, Event: "cancelled", When: "replan-agent", Line: 15},
		{Kind: TriggerEvent, Event: "reordered", When: "replan-agent", Line: 16},
		{Kind: TriggerSchedule, Schedule: "planning.schedule", When: "each-ready-story", Line: 17},
	}
	if !slices.Equal(p.Trigger, wantTriggers) {
		t.Errorf("planner triggers = %+v", p.Trigger)
	}
	if !strings.Contains(p.Additions, "Plan small.") || len(p.Permissions) != 0 {
		t.Errorf("planner additions %q permissions %v", p.Additions, p.Permissions)
	}

	a := rs[0]
	if !slices.Equal(a.Files, []string{"{design}/analysis/"}) || a.Charge != "issues" || len(a.Commands) != 0 {
		t.Errorf("analyzer files %v charge %s commands %v", a.Files, a.Charge, a.Commands)
	}
}

func TestLineNumbersPointAtTheField(t *testing.T) {
	r, p := Parse("design/workflows/roles/planner.md", plannerFile)
	if p != nil {
		t.Fatal(p)
	}
	for key, want := range map[string]int{"title": 2, "name": 5, "definition": 7, "definition.description": 8, "definition.tools": 9, "action": 10, "trigger": 11, "item": 18, "tools": 21, "commands": 30} {
		if got := r.Lines[key]; got != want {
			t.Errorf("line of %s = %d, want %d", key, got, want)
		}
	}
}

func TestNoRolesFolderIsNoRoles(t *testing.T) {
	rs, problems, err := Load(t.TempDir(), "design")
	if err != nil || rs != nil || problems != nil {
		t.Fatalf("Load = %v, %v, %v; want nothing", rs, problems, err)
	}
}

func TestUnparseableFileIsAProblemAndOthersLoad(t *testing.T) {
	cases := []struct {
		name, content string
		line          int
		want          string
	}{
		{"no front matter", "# Just a body\n", 1, "no front matter"},
		{"bad yaml", "---\ntitle: X\nname: [unclosed\n---\n", 3, "does not parse"},
		{"unknown field", "---\ntitle: X\npermision: []\n---\n", 3, "does not parse"},
		{"call entry with unknown key", "---\ntools:\n  - {call: inbox, if: own-thread}\n---\n", 3, "does not parse"},
		{"trigger with event and schedule", "---\ntrigger:\n  - {event: accepted, schedule: digest.schedule}\n---\n", 3, "not both"},
		{"not a mapping", "---\n- a\n- b\n---\n", 2, "not a mapping"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rs, problems := load(t, map[string]string{"broken.md": tc.content, "planner.md": plannerFile})
			if len(rs) != 1 || rs[0].Name != "planner" {
				t.Errorf("roles = %v, want the planner alone", rs)
			}
			if len(problems) != 1 {
				t.Fatalf("problems = %v, want one", problems)
			}
			got := problems[0]
			if got.Code != CodeFile || got.Path != "design/workflows/roles/broken.md" || got.Line != tc.line || !strings.Contains(got.Message, tc.want) {
				t.Errorf("problem = %+v, want %s at line %d containing %q", got, CodeFile, tc.line, tc.want)
			}
		})
	}
}

// edit replaces old with new in content, failing when old is not there.
func edit(t *testing.T, content, old, new string) string {
	t.Helper()
	if !strings.Contains(content, old) {
		t.Fatalf("fixture lacks %q", old)
	}
	return strings.Replace(content, old, new, 1)
}

func TestValidateReportsEachProblem(t *testing.T) {
	cases := []struct {
		name     string
		file     string // fixture file name
		base     string // fixture content
		old, new string // the edit that breaks it
		code     string
		line     int
		want     string
	}{
		{"missing title", "planner.md", plannerFile, "title: Planner\n", "", CodeFile, 1, "needs title"},
		{"missing item", "planner.md", plannerFile, "item: [epic, story]\n", "", CodeFile, 1, "needs item"},
		{"missing trigger", "orchestrator.md", orchestratorFile, "trigger: [continuous]\n", "", CodeFile, 1, "needs trigger"},
		{"bad updated", "planner.md", plannerFile, "updated: 2026-10-28", "updated: yesterday", CodeFile, 3, "not a date"},
		{"bad status", "planner.md", plannerFile, "status: active", "status: live", CodeFile, 4, "status is \"live\""},
		{"name not stem", "plotter.md", plannerFile, "title: Planner", "title: Plotter", CodeFile, 5, "file is plotter.md"},
		{"reserved name", "explorer.md", strings.Replace(plannerFile, "action: plan", "action: explore-code", 1), "name: planner", "name: explorer", CodeFile, 5, "sub-agent"},
		{"reserved role", "planner.md", plannerFile, "role: plan", "role: verify", CodeFile, 6, "story's agent's or a sub-agent's"},
		{"reserved action", "planner.md", plannerFile, "action: plan", "action: push", CodeFile, 10, "no role's"},
		{"another role's action", "planner.md", plannerFile, "action: plan", "action: curate", CodeFile, 10, "curator's"},
		{"unknown trigger", "planner.md", plannerFile, "  - asked\n", "  - sometimes\n", CodeFile, 12, "trigger \"sometimes\""},
		{"unknown event", "planner.md", plannerFile, "{event: accepted, when: replan-agent}", "{event: merged}", CodeFile, 14, "event \"merged\""},
		{"unknown schedule", "planner.md", plannerFile, "schedule: planning.schedule", "schedule: roles.curator.schedule", CodeFile, 17, "schedule \"roles.curator.schedule\""},
		{"continuous not alone", "orchestrator.md", orchestratorFile, "trigger: [continuous]", "trigger: [continuous, asked]", CodeFile, 12, "only trigger"},
		{"trigger condition on wrong trigger", "planner.md", plannerFile, "{event: edited, when: stale-figures}", "{event: accepted, when: stale-figures}", CodeFile, 13, "does not qualify this trigger"},
		{"call condition on a trigger", "planner.md", plannerFile, "{event: edited, when: stale-figures}", "{event: edited, when: draft-only}", CodeFile, 13, "qualifies a call"},
		{"unknown item kind", "planner.md", plannerFile, "item: [epic, story]", "item: [epic, task]", CodeFile, 18, "item kind \"task\""},
		{"bad topic", "planner.md", plannerFile, "topics: [planning]", "topics: [\"two words\"]", CodeFile, 19, "not a word"},
		{"bad charge", "planner.md", plannerFile, "charge: item", "charge: everyone", CodeFile, 20, "charge is \"everyone\""},
		{"unknown condition", "planner.md", plannerFile, "{call: item_new, when: draft-only}", "{call: item_new, when: drafts-only}", CodeFile, 24, "condition \"drafts-only\" is not one"},
		{"condition on wrong call", "planner.md", plannerFile, "{call: item_move, when: to-backlog}", "{call: item_move, when: draft-only}", CodeFile, 26, "does not qualify item_move"},
		{"trigger condition on a call", "planner.md", plannerFile, "{call: item_move, when: to-backlog}", "{call: item_move, when: stale-figures}", CodeFile, 26, "qualifies a trigger"},
		{"under-folder on a call", "planner.md", plannerFile, "{call: item_move, when: to-backlog}", "{call: item_move, when: under-folder}", CodeFile, 26, "each files entry"},
		{"unknown tool", "planner.md", plannerFile, "  - wait_for_events\n", "  - wait_for_ever\n", CodeFile, 29, "not one of flai's MCP tools"},
		{"command under tools", "planner.md", plannerFile, "  - wait_for_events\n", "  - flai issue new\n", CodeFile, 29, "list it under commands"},
		{"tool under commands", "planner.md", plannerFile, "  - flai epic new\n", "  - inbox\n", CodeFile, 32, "not a flai command"},
		{"unknown command", "planner.md", plannerFile, "  - flai epic new\n", "  - flai epic burn\n", CodeFile, 32, "not a flai command"},
		{"wildcard on a command without subcommands", "planner.md", plannerFile, "  - flai epic new\n", "  - flai edit *\n", CodeFile, 32, "not a flai command"},
		{"git in commands", "planner.md", plannerFile, "  - flai epic new\n", "  - git commit\n", CodeGrant, 32, "names git"},
		{"plain gated tool outside a permission", "planner.md", plannerFile, "{call: item_move, when: to-backlog}", "item_move", CodeGrant, 26, "name it under to-backlog, or under to-ready in the permissions entry for roles.planner.permissions.promote_to_ready"},
		{"plain gated command outside a permission", "planner.md", plannerFile, "{call: flai edit, when: ungated-edit}", "flai edit", CodeGrant, 34, "named plainly allows flai edit"},
		{"wildcard reaching a gated call", "planner.md", plannerFile, "  - flai touches *\n", "  - flai subproject *\n", CodeGrant, 35, "allows flai subproject add"},
		{"gated condition outside a permission", "planner.md", plannerFile, "{call: item_edit, when: ungated-edit}", "{call: item_edit, when: finalize-draft}", CodeGrant, 25, "roles.planner.permissions.finalize_drafts"},
		{"gated condition under the wrong permission", "orchestrator.md", orchestratorFile, "tools: [{call: item_edit, when: finalize-draft}]", "tools: [{call: item_edit, when: cost-inputs}]", CodeGrant, 48, "gated by plan_backlog_stories but is named under finalize_drafts"},
		{"plain gated call under the wrong permission", "orchestrator.md", orchestratorFile, "    tools: [release_publish]", "    tools: [release_publish, queue_add]", CodeGrant, 65, "orchestration.permissions.accept_reviews"},
		{"never granted call", "orchestrator.md", orchestratorFile, "  - flai queue\n", "  - flai queue hold\n", CodeGrant, 27, "no setting allows"},
		{"never granted call inside a permission", "orchestrator.md", orchestratorFile, "    tools: [release_publish]", "    tools: [release_publish, shared_paths_edit]", CodeGrant, 65, "no setting allows"},
		{"wildcard reaching a gated call in a permission-less entry", "orchestrator.md", orchestratorFile, "  - flai queue\n", "  - flai queue *\n", CodeGrant, 27, "allows flai queue add"},
		{"wildcard reaching a never granted call", "orchestrator.md", orchestratorFile, "  - flai thread new\n", "  - flai thread *\n", CodeGrant, 28, "reaches flai thread confirm"},
		{"story agent call while works_story is false", "orchestrator.md", orchestratorFile, "  - wait_for_events\n", "  - story_start\n", CodeGrant, 21, "works_story is false"},
		{"accept with batch", "orchestrator.md", orchestratorFile, "{call: flai accept, when: accept-as-self}", "flai accept", CodeGrant, 61, "name it under accept-as-self"},
		{"another role's setting", "planner.md", plannerFile, "  - {call: flai move, when: to-backlog}\n", "  - {call: flai move, when: to-backlog}\npermissions:\n  - setting: roles.analyzer.permissions.publish\n    tools: [release_publish]\n", CodeGrant, 42, "roles.planner.permissions.publish"},
		{"orchestration setting in a role's file", "analyzer.md", analyzerFile, "files: ['{design}/analysis/']\n", "files: ['{design}/analysis/']\npermissions:\n  - setting: orchestration.permissions.publish\n    tools: [release_publish]\n", CodeGrant, 18, "roles.analyzer.permissions.publish"},
		{"unknown permission", "analyzer.md", analyzerFile, "files: ['{design}/analysis/']\n", "files: ['{design}/analysis/']\npermissions:\n  - setting: roles.analyzer.permissions.everything\n    tools: [release_publish]\n", CodeFile, 18, "permission \"everything\""},
		{"malformed setting", "analyzer.md", analyzerFile, "files: ['{design}/analysis/']\n", "files: ['{design}/analysis/']\npermissions:\n  - setting: publish\n    tools: [release_publish]\n", CodeFile, 18, "neither orchestration.permissions"},
		{"empty permission", "analyzer.md", analyzerFile, "files: ['{design}/analysis/']\n", "files: ['{design}/analysis/']\npermissions:\n  - setting: roles.analyzer.permissions.publish\n", CodeFile, 18, "allows nothing"},
		{"missing marker", "planner.md", plannerFile, "<!-- system-flow:end-of-baseline -->\n", "", CodeFile, 1, "no baseline marker"},
		{"files not a folder", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'{design}/analysis'", CodeFile, 16, "end it with /"},
		{"files the root", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'./'", CodeFile, 16, "project root"},
		{"files outside the project", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'../elsewhere/'", CodeFile, 16, "not a folder path"},
		{"files the roles folder", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'{design}/workflows/'", CodeFile, 16, "workflows/"},
		{"files inside the roles folder", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'design/workflows/roles/'", CodeFile, 16, "workflows/"},
		{"files holding the roles folder", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'{design}/'", CodeFile, 16, "workflows/"},
		{"files git", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'.git/hooks/'", CodeFile, 16, ".git/"},
		{"files claude", "analyzer.md", analyzerFile, "'{design}/analysis/'", "'.claude/'", CodeFile, 16, ".claude/"},
		{"definition missing", "planner.md", plannerFile, "definition:\n  description: Plans one epic or one story.\n  tools: [Read, Grep, Bash, mcp__flai__prime, mcp__flai__item_new, mcp__flai__item_edit, mcp__flai__item_move]\n", "", CodeDefinition, 1, "needs definition"},
		{"definition description missing", "planner.md", plannerFile, "  description: Plans one epic or one story.\n", "", CodeDefinition, 7, "needs description"},
		{"definition tools missing", "planner.md", plannerFile, "  tools: [Read, Grep, Bash, mcp__flai__prime, mcp__flai__item_new, mcp__flai__item_edit, mcp__flai__item_move]\n", "", CodeDefinition, 7, "needs tools"},
		{"definition lists a refused tool", "planner.md", plannerFile, "mcp__flai__item_move]", "mcp__flai__item_move, mcp__flai__release_publish]", CodeDefinition, 9, "mcp__flai__release_publish, which the role's rules refuse"},
		{"Edit with no files", "planner.md", plannerFile, "tools: [Read, Grep,", "tools: [Read, Edit, Grep,", CodeDefinition, 9, "lists Edit, but files names no folder"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rs, problems := load(t, map[string]string{tc.file: edit(t, tc.base, tc.old, tc.new)})
			if len(problems) != 0 {
				t.Fatalf("load problems: %v", problems)
			}
			got := Validate(rs, "design", vocab)
			if len(got) != 1 {
				t.Fatalf("Validate = %v, want one problem", got)
			}
			p := got[0]
			if p.Code != tc.code || p.Line != tc.line || !strings.Contains(p.Message, tc.want) || p.Path != "design/workflows/roles/"+tc.file {
				t.Errorf("problem = %s, want %s at line %d containing %q", p, tc.code, tc.line, tc.want)
			}
		})
	}
}

func TestValidateReportsDuplicateNamesAndVerbs(t *testing.T) {
	other := edit(t, edit(t, plannerFile, "name: planner", "name: replanner"), "action: plan", "action: replan")
	rs, _ := load(t, map[string]string{"planner.md": plannerFile, "replanner.md": other})
	got := Validate(rs, "design", vocab)
	if len(got) != 1 || got[0].Code != CodeFile || got[0].Line != 6 || !strings.Contains(got[0].Message, "role plan is also design/workflows/roles/planner.md's") {
		t.Errorf("Validate = %v, want the second file's duplicate role", got)
	}
}

func TestWorksStoryAllowsAStoryAgentsCalls(t *testing.T) {
	f := edit(t, edit(t, plannerFile, "charge: item\n", "charge: item\nworks_story: true\n"), "  - wait_for_events\n", "  - story_start\n  - task_done\n")
	rs, _ := load(t, map[string]string{"planner.md": f})
	if got := Validate(rs, "design", vocab); len(got) != 0 {
		t.Errorf("Validate = %v, want none", got)
	}
}

func TestNilVocabularySkipsNameChecks(t *testing.T) {
	f := edit(t, edit(t, plannerFile, "  - wait_for_events\n", "  - some_future_tool\n"), "mcp__flai__item_move]", "mcp__flai__item_move, mcp__flai__some_other]")
	rs, _ := load(t, map[string]string{"planner.md": f})
	if got := Validate(rs, "design", Vocabulary{}); len(got) != 0 {
		t.Errorf("Validate = %v, want none with a nil vocabulary", got)
	}
}

func TestDesignFolderFromTheManifest(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "docs", "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	f := edit(t, analyzerFile, "'{design}/analysis/'", "'docs/design/workflows/'")
	if err := os.WriteFile(filepath.Join(dir, "analyzer.md"), []byte(f), 0o644); err != nil {
		t.Fatal(err)
	}
	rs, problems, err := Load(root, "docs/design")
	if err != nil || len(problems) != 0 || len(rs) != 1 || rs[0].Path != "docs/design/workflows/roles/analyzer.md" {
		t.Fatalf("Load = %v, %v, %v", rs, problems, err)
	}
	got := Validate(rs, "docs/design", vocab)
	if len(got) != 1 || !strings.Contains(got[0].Message, "docs/design/workflows/") {
		t.Errorf("Validate = %v, want the files entry reaching docs/design/workflows/", got)
	}
}
