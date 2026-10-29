package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// roleFile is a valid role file, an analyzer of its own named name, with
// front matter edited by edit: each key's line replaced by the given lines,
// or removed when they are empty.
func roleFile(name string, edit map[string]string) string {
	lines := []string{
		"---",
		"title: Reviewer",
		"updated: 2026-10-28",
		"status: active",
		"name: " + name,
		"role: review",
		"definition:",
		"  description: Reviews what is asked.",
		"  tools: [Read, Grep, Write, mcp__flai__board, mcp__flai__inbox, mcp__flai__issue_new]",
		"action: review",
		"trigger: [asked]",
		"item: none",
		"topics: [analysis]",
		"charge: issues",
		"tools: [inbox, issue_new]",
		"commands: [flai issue new]",
		"files: [\"{design}/reviews/\"]",
		"---",
		"",
		"You are the reviewer.",
		"",
		roles.Marker,
		"",
		"## Project additions",
		"",
	}
	var out []string
	for _, l := range lines {
		key, _, _ := strings.Cut(l, ":")
		if repl, ok := edit[key]; ok {
			if repl != "" {
				out = append(out, repl)
			}
			continue
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

// rolesProject is a project whose design folder's role files are files.
func rolesProject(t *testing.T, files map[string]string) *workitem.Repo {
	t.Helper()
	repo := planningProject(t, "")
	dir := filepath.Join(repo.Root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

// roleFindings are the findings flai check's roles rule gives, without the
// other rules.
func roleFindings(repo *workitem.Repo) []Finding {
	c := &checker{repo: repo, res: &Result{}}
	c.roles()
	sortFindings(c.res.Findings)
	return c.res.Findings
}

// S-0415: the roles a convention may list beside flai's own are the verbs
// the project's role files define; files that do not parse define none.
func TestDefinedRolesAreTheRoleFilesVerbs(t *testing.T) {
	repo := rolesProject(t, map[string]string{"auditor.md": roleFile("auditor", nil)})
	c := &checker{repo: repo, res: &Result{}}
	if got := strings.Join(c.definedRoles(), ","); got != "review" {
		t.Errorf("defined roles: %q", got)
	}
	broken := rolesProject(t, map[string]string{"reviewer.md": "---\nrole: [\n---\n"})
	c = &checker{repo: broken, res: &Result{}}
	if got := c.definedRoles(); len(got) != 0 {
		t.Errorf("defined roles of a file that does not parse: %q", got)
	}
}

// S-0414: valid role files, and the index beside them, are not reported, a
// definition's flai MCP read tool among them, and flai check runs the rule.
func TestValidRoleFilesAreClean(t *testing.T) {
	repo := rolesProject(t, map[string]string{
		"README.md":   "# Roles\n\nNo front matter: the index is not a role.\n",
		"reviewer.md": roleFile("reviewer", nil),
	})
	if got := roleFindings(repo); len(got) != 0 {
		t.Errorf("valid role files: %+v", got)
	}
	for _, code := range []string{roles.CodeFile, roles.CodeGrant, roles.CodeDefinition} {
		if got := findings(t, repo, code); len(got) != 0 {
			t.Errorf("flai check on valid role files, %s: %+v", code, got)
		}
	}
}

// S-0535 (ADR-0211): a trigger naming each new event, and each new trigger
// condition on its event, passes flai check, and to-backlog still qualifies
// item_move beside the {event: moved} it now qualifies.
func TestNewTriggerEventsAndConditionsAreClean(t *testing.T) {
	repo := rolesProject(t, map[string]string{"reviewer.md": roleFile("reviewer", map[string]string{
		"trigger": "trigger: [asked, {event: moved}, {event: thread}, {event: planned}, {event: ready-empty}, {event: release-due}, " +
			"{event: moved, when: to-review}, {event: moved, when: to-backlog}, {event: edited, when: draft-changed}, {event: thread, when: awaits-orchestrator}]",
		"  tools": "  tools: [Read, Grep, Write, mcp__flai__board, mcp__flai__inbox, mcp__flai__issue_new, mcp__flai__item_move]",
		"tools":   "tools: [inbox, issue_new, {call: item_move, when: to-backlog}]",
	})})
	if got := roleFindings(repo); len(got) != 0 {
		t.Errorf("the new events and conditions: %+v", got)
	}
	if got := findings(t, repo, roles.CodeFile); len(got) != 0 {
		t.Errorf("flai check on the new events and conditions: %+v", got)
	}
}

// S-0535 (ADR-0211): each new trigger condition on an event it does not
// qualify is roles.file on the trigger's line.
func TestNewTriggerConditionsOnTheWrongEvent(t *testing.T) {
	for _, tc := range []struct{ trigger, says string }{
		{"{event: edited, when: to-review}", "it qualifies an {event} entry for moved"},
		{"{event: accepted, when: to-backlog}", "it qualifies an {event} entry for moved"},
		{"{event: moved, when: draft-changed}", "it qualifies an {event} entry for edited"},
		{"{event: thread-opened, when: awaits-orchestrator}", "it qualifies an {event} entry for thread"},
	} {
		repo := rolesProject(t, map[string]string{"reviewer.md": roleFile("reviewer", map[string]string{"trigger": "trigger: [asked, " + tc.trigger + "]"})})
		got := roleFindings(repo)
		if len(got) != 1 || got[0].Level != Error || got[0].Rule != roles.CodeFile || got[0].Line != 11 || !strings.Contains(got[0].Message, tc.says) {
			t.Errorf("%s: %+v, want one roles.file error on line 11 saying %q", tc.trigger, got, tc.says)
		}
	}
}

// S-0536 (ADR-0209 §1): continuous, alone or beside other entries, is a
// roles.file error on its entry's line.
func TestContinuousTriggerIsAFinding(t *testing.T) {
	for _, tc := range []struct {
		trigger string
		line    int
	}{
		{"trigger: [continuous]", 11},
		{"trigger:\n  - asked\n  - continuous", 13},
	} {
		repo := rolesProject(t, map[string]string{"reviewer.md": roleFile("reviewer", map[string]string{"trigger": tc.trigger})})
		got := findings(t, repo, roles.CodeFile)
		if len(got) != 1 || got[0].Level != Error || got[0].Path != "design/workflows/roles/reviewer.md" || got[0].Line != tc.line || !strings.Contains(got[0].Message, "no longer flai's (ADR-0209 §1)") {
			t.Errorf("%q: %+v, want one roles.file error on line %d saying continuous is no longer flai's", tc.trigger, got, tc.line)
		}
	}
}

// S-0414: flai check reports roles.file, roles.grant, and roles.definition as
// errors on the line of the field each is about, in the role file.
func TestRoleFileProblemsAreErrorsOnTheirLine(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit map[string]string
		code string
		line int
		says string
	}{
		{"a status outside the allowed ones", map[string]string{"status": "status: live"}, roles.CodeFile, 4, `status is "live"`},
		{"a field the format does not have", map[string]string{"item": "item: none\ncolour: red"}, roles.CodeFile, 13, "does not parse"},
		{"an unknown condition", map[string]string{"tools": "tools: [inbox, issue_new, {call: thread_reply, when: sometimes}]"}, roles.CodeFile, 15, `condition "sometimes"`},
		{"a gated call named plainly", map[string]string{"tools": "tools: [inbox, issue_new, item_move]"}, roles.CodeGrant, 15, "item_move named plainly"},
		{"a call no setting allows", map[string]string{"commands": "commands: [flai issue new, flai push]"}, roles.CodeGrant, 16, "flai push is never granted"},
		{"a flai MCP tool the rules refuse", map[string]string{"tools": "tools: [inbox]"}, roles.CodeDefinition, 9, "mcp__flai__issue_new"},
		{"an edit tool and no files", map[string]string{"files": "files: []"}, roles.CodeDefinition, 9, "definition.tools lists Write"},
		{"no definition", map[string]string{"definition": "", "  description": "", "  tools": ""}, roles.CodeDefinition, 1, "needs definition"},
	} {
		repo := rolesProject(t, map[string]string{"reviewer.md": roleFile("reviewer", tc.edit)})
		got := roleFindings(repo)
		if len(got) != 1 {
			t.Errorf("%s: want one finding, got %+v", tc.name, got)
			continue
		}
		f := got[0]
		if f.Level != Error || f.Rule != tc.code || f.Path != "design/workflows/roles/reviewer.md" || f.Line != tc.line || !strings.Contains(f.Message, tc.says) {
			t.Errorf("%s: %+v, want an error %s on line %d saying %q", tc.name, f, tc.code, tc.line, tc.says)
		}
		if all := findings(t, repo, tc.code); len(all) != 1 {
			t.Errorf("%s: flai check gives %+v under %s", tc.name, all, tc.code)
		}
	}
}

// S-0414: a role file that does not load leaves the others checked, and two
// roles with one name are one finding, on the second's name line.
func TestRoleFilesAreCheckedTogether(t *testing.T) {
	repo := rolesProject(t, map[string]string{
		"broken.md":   "no front matter\n",
		"reviewer.md": roleFile("reviewer", map[string]string{"role": "role: plan", "action": "action: plan"}),
	})
	got := roleFindings(repo)
	if len(got) != 2 {
		t.Fatalf("want the broken file and the planner's action: %+v", got)
	}
	if got[0].Path != "design/workflows/roles/broken.md" || got[0].Rule != roles.CodeFile || got[0].Line != 1 {
		t.Errorf("the file that does not load: %+v", got[0])
	}
	if got[1].Path != "design/workflows/roles/reviewer.md" || got[1].Line != 10 || !strings.Contains(got[1].Message, "action plan is the planner's") {
		t.Errorf("the other file is still checked: %+v", got[1])
	}
}

// S-0414: the role files this repository ships give flai check no finding.
func TestShippedRoleFilesAreClean(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo")
	}
	repo, err := workitem.Open("../../..")
	if err != nil {
		t.Skip("monorepo not present")
	}
	list, _, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil || len(list) == 0 {
		t.Fatalf("the repository's role files: %d, %v", len(list), err)
	}
	for _, f := range roleFindings(repo) {
		t.Errorf("%s:%d: %s: %s", f.Path, f.Line, f.Rule, f.Message)
	}
}
