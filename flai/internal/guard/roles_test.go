package guard

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
)

// reviewerFile is a role a project defines: the reviewer, role review. Its
// release_publish, story_start, flai queue *, and the orchestrator's publish
// setting are entries flai check reports as roles.grant, which the guard
// never grants.
const reviewerFile = `---
title: Reviewer
updated: 2026-10-28
status: active
name: reviewer
role: review
definition:
  description: Reviews.
  tools: [Read, Edit, mcp__flai__thread_open]
action: review
trigger: [asked]
item: none
tools:
  - thread_open
  - {call: item_move, when: to-backlog}
  - {call: thread_reply, when: own-thread}
  - release_publish
  - story_start
commands:
  - flai issue new
  - {call: flai move, when: to-backlog}
  - flai touches *
  - flai queue
  - flai queue *
files: ["{design}/reviews/"]
permissions:
  - setting: roles.reviewer.permissions.promote_to_ready
    tools: [{call: item_move, when: to-ready}]
    commands: [{call: flai move, when: to-ready}]
  - setting: roles.reviewer.permissions.answer_threads
    tools: [{call: thread_reply, when: answer-threads}]
  - setting: orchestration.permissions.publish
    tools: [release_publish]
---

The reviewer reviews.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// reviewerPath is where the reviewer's file is, under the design folder.
const reviewerPath = "design/workflows/roles/reviewer.md"

// reviewer is a guard in the reviewer's session in a project at root, its
// threads told apart by reviewerOpeners, under its own permissions p.
func reviewer(t *testing.T, root string, p manifest.Permissions) Guard {
	t.Helper()
	r, problem := roles.Parse(reviewerPath, reviewerFile)
	if problem != nil {
		t.Fatal(problem)
	}
	return Guard{Commands: append(slices.Clone(g.Commands), "queue"), Role: "review", RoleFile: &r, RolePermissions: p,
		Root: root, Design: "design", Opener: reviewerOpeners}
}

// reviewerOpeners knows TH-0001, which the reviewer opened, and TH-0002,
// which the operator opened.
func reviewerOpeners(id string) (string, error) {
	switch id {
	case "TH-0001":
		return "reviewer", nil
	case "TH-0002":
		return "alex", nil
	}
	return "", errors.New(id + " not found")
}

// roleTool, roleShell, and roleWrite are the reviewer's own call of a flai
// tool, a command line, and Write on a file.
func roleTool(name string, in Args) Call {
	return Call{Role: "review", Kind: FlaiTool, Tool: name, Args: in}
}

func roleShell(command string) Call { return Call{Role: "review", Kind: Shell, Command: command} }

func roleWrite(file string) Call {
	return Call{Role: "review", Kind: FileWrite, Name: "Write", Path: file}
}

// S-0416, ADR-0207: a role a project defines reads, makes each call its file
// lists under the entry's condition, and is refused every other, naming the
// role, the rule, its file, and thread_open.
func TestADefinedRoleIsHeldToItsFile(t *testing.T) {
	root := t.TempDir()
	rv := reviewer(t, root, manifest.Permissions{})
	for _, c := range []Call{
		roleTool("board", Args{}),
		roleTool("item_get", Args{ID: "S-0001"}),
		roleShell("flai board --json && git log -1 && flai show S-0001"),
		roleTool("thread_open", Args{ID: "S-0001"}),
		roleTool("item_move", Args{ID: "S-0001", To: Backlog}),
		roleTool("thread_reply", Args{ID: "TH-0001"}),
		roleShell("flai issue new x --class defect"),
		roleShell("flai move S-0001 backlog"),
		roleShell("flai touches S-0001 --add flai/x.go && flai touches suggest S-0001"),
		roleShell("flai queue --json"),
	} {
		if v := Decide(rv, c); v.Outcome != Allow {
			t.Errorf("%+v refused: %s", c, v.Reason)
		}
	}
	for _, c := range []struct {
		call    Call
		reason  string
		refused string
	}{
		{roleTool("item_new", Args{Type: "story"}), "the role reviewer cannot call item_new: its role file lists no such call. Ask the operator with thread_open on the item if it needs doing (" + reviewerPath + ", ADR-0207).", "item_new"},
		{roleTool("item_move", Args{ID: "S-0001", To: "review"}), "the role reviewer cannot move S-0001 to review: it may make it only under to-backlog (on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog), or to-ready (a move to ready and nowhere else).", "item_move S-0001 review"},
		{roleShell("flai board && flai move S-0001 review"), `the role reviewer cannot run "flai move S-0001 review": it may make it only under to-backlog`, "flai move S-0001 review"},
		{roleTool("thread_reply", Args{ID: "TH-0001", Recommendation: true}), "the role reviewer cannot call thread_reply: it may make it only under own-thread", "thread_reply TH-0001 recommendation"},
		{roleTool("thread_reply", Args{ID: "TH-0002"}), "the role reviewer cannot call thread_reply: it may make it only under own-thread", "thread_reply TH-0002"},
		{roleShell("flai epic new x"), `the role reviewer cannot run "flai epic new x": its role file lists no such call.`, "flai epic new x"},
		{roleShell("git commit -m x"), `the role reviewer cannot run "git commit -m x": it runs only git's reads.`, "git commit -m x"},
		{roleShell("flai queue hold"), `the role reviewer cannot run "flai queue hold": its role file lists no such call.`, "flai queue hold"},
	} {
		v := Decide(rv, c.call)
		if v.Outcome != Deny || !strings.HasPrefix(v.Reason, c.reason) || !strings.Contains(v.Reason, "thread_open") || v.Call != c.refused || v.Needs != "" {
			t.Errorf("%+v: %+v", c.call, v)
		}
	}
	// a sub-agent of the role is held as every sub-agent is
	sub := roleTool("item_move", Args{ID: "S-0001", To: Backlog})
	sub.Subagent, sub.Agent = true, "explorer"
	if v := Decide(rv, sub); !strings.Contains(v.Reason, "a sub-agent (explorer) cannot call item_move") {
		t.Errorf("sub-agent: %+v", v)
	}
}

// S-0416: a call a permissions entry lists is refused while the role's own
// setting is off, naming it, and passes once it is on; the orchestrator's
// settings grant the role nothing.
func TestADefinedRolesPermissionIsItsOwnSetting(t *testing.T) {
	root := t.TempDir()
	promote, answer := roleTool("item_move", Args{ID: "S-0001", To: Ready}), roleTool("thread_reply", Args{ID: "TH-0002", Recommendation: true})
	cli := roleShell("flai move S-0001 ready")
	all := manifest.Permissions{PlanBacklogEpics: true, PlanBacklogStories: true, FinalizeDrafts: true, PromoteToReady: true, OrderReady: true, AnswerThreads: manifest.AnswerAutonomous, AcceptReviews: true, Publish: true, EditSubprojects: true}
	off := reviewer(t, root, manifest.Permissions{})
	off.Permissions = all
	setting := "roles.reviewer.permissions." + manifest.PermitPromoteToReady
	for _, c := range []Call{promote, cli} {
		v := Decide(off, c)
		if v.Outcome != Deny || v.Needs != setting || !strings.Contains(v.Reason, "it needs "+setting+", which is off") {
			t.Errorf("off: %+v", v)
		}
	}
	if v := Decide(off, answer); v.Outcome != Deny || v.Needs != "roles.reviewer.permissions.answer_threads" {
		t.Errorf("answer off: %+v", v)
	}
	if v := Decide(off, roleTool(ReleasePublish, Args{})); v.Outcome != Deny || !strings.Contains(v.Reason, "lists no such call") {
		t.Errorf("the orchestrator's publish granted the role: %+v", v)
	}

	on := reviewer(t, root, manifest.Permissions{PromoteToReady: true, AnswerThreads: manifest.AnswerRecommend})
	for _, c := range []Call{promote, cli, answer} {
		if v := Decide(on, c); v.Outcome != Allow {
			t.Errorf("on: %+v refused: %s", c, v.Reason)
		}
	}
	// recommend lets it recommend, not answer citing a source
	sourced := roleTool("thread_reply", Args{ID: "TH-0002", Source: "design/system/overview.md"})
	if v := Decide(on, sourced); v.Outcome != Deny || !strings.Contains(v.Reason, "answer-threads") {
		t.Errorf("a sourced answer under recommend: %+v", v)
	}
	on.RolePermissions.AnswerThreads = manifest.AnswerAutonomous
	if v := Decide(on, sourced); v.Outcome != Allow {
		t.Errorf("a sourced answer under autonomous: %+v", v)
	}
	// a move past ready meets the condition, not the setting
	if v := Decide(on, roleTool("item_move", Args{ID: "S-0001", To: "done"})); v.Outcome != Deny || v.Needs != "" {
		t.Errorf("to done: %+v", v)
	}
}

// S-0416: the role edits a file under a folder its file names, with {design}
// the design folder, and no other, however the path is put; an entry flai
// check reports as roles.grant is never granted.
func TestADefinedRoleEditsUnderItsFolderAndNeverGetsAForbiddenCall(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{"design/reviews", "design/system"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(filepath.Join(root, "design", "system"), filepath.Join(root, "design", "reviews", "out")); err != nil {
		t.Fatal(err)
	}
	rv := reviewer(t, root, manifest.Permissions{Publish: true, AcceptReviews: true})
	for _, file := range []string{filepath.Join(root, "design", "reviews", "r.md"), "design/reviews/new/r.md"} {
		if v := Decide(rv, roleWrite(file)); v.Outcome != Allow {
			t.Errorf("%s refused: %s", file, v.Reason)
		}
	}
	for _, file := range []string{
		filepath.Join(root, "design", "system", "x.md"),
		filepath.Join(root, "design", "reviews"),
		root + "/design/reviews/../system/x.md",
		filepath.Join(root, "design", "reviews", "out", "x.md"),
	} {
		v := Decide(rv, roleWrite(file))
		if v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the role reviewer cannot use Write on "+file+": it edits only a file under design/reviews/.") || v.Call != "Write "+file {
			t.Errorf("%s: %+v", file, v)
		}
	}
	// with no design folder given, {design} names no folder
	rv.Design = ""
	if v := Decide(rv, roleWrite(filepath.Join(root, "design", "reviews", "r.md"))); v.Outcome != Deny || !strings.Contains(v.Reason, "names no folder it may edit") {
		t.Errorf("no design folder: %+v", v)
	}
	rv.Design = "design"
	for _, c := range []Call{roleTool(ReleasePublish, Args{}), roleTool("story_start", Args{ID: "S-0001"}), roleShell("flai queue add S-0001 --by reviewer")} {
		if v := Decide(rv, c); v.Outcome != Deny || !strings.Contains(v.Reason, "lists no such call") {
			t.Errorf("never granted: %+v: %+v", c, v)
		}
	}
}

// S-0416: a role whose file is missing or does not parse reads and is
// refused every other call, naming the file or the role and flai check.
func TestADefinedRoleFailsClosed(t *testing.T) {
	root := t.TempDir()
	problem := reviewerPath + " does not parse: front matter does not parse; run flai check"
	broken := Guard{Commands: g.Commands, Role: "review", RoleProblem: problem, Root: root, Design: "design"}
	parsed := reviewer(t, root, manifest.Permissions{PromoteToReady: true})
	parsed.RoleProblem = problem
	missing := Guard{Commands: g.Commands, Role: "review", Root: root, Design: "design"}
	for _, tc := range []struct {
		name  string
		guard Guard
		why   string
	}{
		{"a file that does not parse", broken, problem},
		{"a problem beside a file", parsed, problem},
		{"no file", missing, "no role file defines the role review; run flai check"},
	} {
		for _, c := range []Call{roleTool("board", Args{}), roleShell("flai board && git status"), {Role: "review", Kind: Other, Name: "Read"}} {
			if v := Decide(tc.guard, c); v.Outcome != Allow {
				t.Errorf("%s: read %+v refused: %s", tc.name, c, v.Reason)
			}
		}
		for _, c := range []Call{roleTool("thread_open", Args{ID: "S-0001"}), roleTool("item_move", Args{ID: "S-0001", To: Ready}), roleShell("flai issue new x"), roleWrite(filepath.Join(root, "design", "reviews", "r.md"))} {
			v := Decide(tc.guard, c)
			if v.Outcome != Deny || !strings.Contains(v.Reason, ": "+tc.why+". ") || !strings.Contains(v.Reason, "flai check") || v.Needs != "" {
				t.Errorf("%s: %+v: %+v", tc.name, c, v)
			}
		}
	}
}

// S-0416: a flai command entry names the command with none of its
// subcommands, or in any form when it has none, a subcommand, or with * any
// form; the subcommands come from the command tree when given, and else
// from the reads, the file, and the calls no setting allows or a setting
// gates.
func TestADefinedRolesCommandEntriesMatchTheCommandTree(t *testing.T) {
	rv := reviewer(t, t.TempDir(), manifest.Permissions{})
	h := rv.roleHeld("review")
	for _, tc := range []struct {
		entry, cmd, sub string
		want            bool
	}{
		{"flai queue", "queue", "", true},
		{"flai queue", "queue", "hold", false},
		{"flai queue", "queue", "add", false},
		{"flai move", "move", "S-0001", true},
		{"flai issue new", "issue", "new", true},
		{"flai issue new", "issue", "bump", false},
		{"flai touches *", "touches", "S-0001", true},
		{"flai touches *", "touches", "suggest", true},
		{"flai touches *", "move", "S-0001", false},
		{"thread_open", "thread", "", false},
	} {
		if got := h.matches(tc.entry, tc.cmd, tc.sub); got != tc.want {
			t.Errorf("nil tree: %q on flai %s %s: %v", tc.entry, tc.cmd, tc.sub, got)
		}
	}
	rv.Subcommands = map[string][]string{"queue": {"add", "hold"}, "move": nil, "stream": {"diff", "sync"}}
	h = rv.roleHeld("review")
	if !h.matches("flai stream", "stream", "S-0001") || h.matches("flai stream", "stream", "sync") || !h.matches("flai move", "move", "sync") {
		t.Error("the command tree's subcommands are not the ones matched")
	}
}

// TestUngatedEditLeavesOutOnlyAFinalizingAndTheCostInputs: ungated-edit lets
// an edit set a cost of delay value but refuses one that finalizes a draft,
// sets a cost of delay input, or clears the cost of delay (ADR-0207 §5).
func TestUngatedEditLeavesOutOnlyAFinalizingAndTheCostInputs(t *testing.T) {
	file := strings.Replace(reviewerFile, "  - thread_open\n", "  - thread_open\n  - {call: item_edit, when: ungated-edit}\n", 1)
	file = strings.Replace(file, "  - flai issue new\n", "  - flai issue new\n  - {call: flai edit, when: ungated-edit}\n", 1)
	r, problem := roles.Parse(reviewerPath, file)
	if problem != nil {
		t.Fatal(problem)
	}
	gr := reviewer(t, t.TempDir(), manifest.Permissions{})
	gr.RoleFile = &r
	value := map[string]json.RawMessage{"value": json.RawMessage("12")}
	input := map[string]json.RawMessage{"revenue_per_week": json.RawMessage("500")}
	for _, in := range []Args{
		{ID: "S-0001", Fields: []string{"id", "title"}},
		{ID: "S-0001", Fields: []string{"cost_of_delay", "id"}, CostOfDelay: value},
	} {
		if v := Decide(gr, roleTool("item_edit", in)); v.Outcome != Allow {
			t.Errorf("item_edit %v: refused: %s", in.Fields, v.Reason)
		}
	}
	for _, in := range []Args{
		{ID: "S-0001", Fields: []string{"draft", "id"}},
		{ID: "S-0001", Fields: []string{"cost_of_delay", "id"}, CostOfDelay: input},
		{ID: "S-0001", Fields: []string{"clear_cost_of_delay", "id"}},
	} {
		if v := Decide(gr, roleTool("item_edit", in)); v.Outcome == Allow || !strings.Contains(v.Reason, "ungated-edit") {
			t.Errorf("item_edit %v: want a refusal under ungated-edit, got %+v", in.Fields, v)
		}
	}
	for _, c := range []string{"flai edit S-0001 --title x", "flai edit S-0001 --cost-of-delay-value 12"} {
		if v := Decide(gr, roleShell(c)); v.Outcome != Allow {
			t.Errorf("%s: refused: %s", c, v.Reason)
		}
	}
	for _, c := range []string{"flai edit S-0001 --no-draft", "flai edit S-0001 --revenue-per-week 500", "flai edit S-0001 --clear-cost-of-delay"} {
		if v := Decide(gr, roleShell(c)); v.Outcome == Allow || !strings.Contains(v.Reason, "ungated-edit") {
			t.Errorf("%s: want a refusal under ungated-edit, got %+v", c, v)
		}
	}
}

// builtInFile is flai's baseline role file of name, as the project's file
// at design/workflows/roles/<name>.md, with each of edits, old and new in
// turn, made to it.
func builtInFile(t *testing.T, name string, edits ...string) *roles.Role {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "roles", "baseline", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for i := 0; i+1 < len(edits); i += 2 {
		if !strings.Contains(content, edits[i]) {
			t.Fatalf("%s.md has no %q", name, edits[i])
		}
		content = strings.Replace(content, edits[i], edits[i+1], 1)
	}
	r, problem := roles.Parse("design/workflows/roles/"+name+".md", content)
	if problem != nil {
		t.Fatal(problem)
	}
	return &r
}

// S-0416, ADR-0207: the planner is held to its role file, the project's
// when there is one: a call the file leaves out is refused it with the
// file's reason, where its rule in Go allowed it, and a call the file adds
// passes, where its rule in Go refused it.
func TestThePlannerIsHeldToItsRoleFile(t *testing.T) {
	narrower := Guard{Commands: g.Commands, Role: RolePlan, Design: "design", RoleFile: builtInFile(t, "planner", "  - thread_reply\n", "")}
	reply := Call{Role: RolePlan, Kind: FlaiTool, Tool: "thread_reply", Args: Args{ID: "TH-0001"}}
	v := Decide(narrower, reply)
	want := "the role planner cannot call thread_reply: its role file lists no such call. Ask the operator with thread_open on the item if it needs doing (design/workflows/roles/planner.md, ADR-0207)."
	if v.Outcome != Deny || v.Reason != want || v.Call != "thread_reply TH-0001" {
		t.Errorf("thread_reply left out of the file: %+v", v)
	}
	if v := Decide(planGuard, reply); v.Outcome != Allow {
		t.Errorf("thread_reply under the baseline file: %+v", v)
	}

	wider := Guard{Commands: g.Commands, Role: RolePlan, Design: "design", RoleFile: builtInFile(t, "planner", "  - inbox\n", "  - inbox\n  - issue_new\n")}
	file := Call{Role: RolePlan, Kind: FlaiTool, Tool: "issue_new"}
	if v := Decide(wider, file); v.Outcome != Allow {
		t.Errorf("issue_new added to the file: %+v", v)
	}
	if v := Decide(planGuard, file); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the planner cannot call issue_new: ") {
		t.Errorf("issue_new under the baseline file: %+v", v)
	}
}

// S-0416, ADR-0207: with no file of its own, the planner is held to flai's
// baseline file: ungated-edit refuses an edit that finalizes a draft or sets
// a cost of delay input, which its rule in Go allowed by MCP, with the
// file's reason; a refusal its rule in Go makes too keeps its wording; and
// the commands it lists match with no command tree.
func TestThePlannerIsHeldToTheBaselineFile(t *testing.T) {
	input := map[string]json.RawMessage{"revenue_per_week": json.RawMessage("500")}
	for _, in := range []Args{
		{ID: "S-0001", Fields: []string{"draft", "id"}},
		{ID: "S-0001", Fields: []string{"cost_of_delay", "id"}, CostOfDelay: input},
	} {
		v := Decide(planGuard, Call{Role: RolePlan, Kind: FlaiTool, Tool: "item_edit", Args: in})
		if v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the role planner cannot call item_edit: it may make it only under ungated-edit") || !strings.Contains(v.Reason, "workflows/roles/planner.md, ADR-0207") {
			t.Errorf("item_edit %v: %+v", in.Fields, v)
		}
	}
	if v := Decide(planGuard, Call{Role: RolePlan, Kind: Shell, Command: "flai edit S-0001 --revenue-per-week 500"}); v.Outcome != Deny || !strings.Contains(v.Reason, "ungated-edit") {
		t.Errorf("flai edit with a cost of delay input: %+v", v)
	}
	for c, want := range map[string]string{
		"flai edit S-0001 --no-draft": "finalizing a draft is the operator's",
		"flai story new x":            drafts,
		"flai move S-0001 ready":      "it moves an item to backlog and no further",
	} {
		if v := Decide(planGuard, Call{Role: RolePlan, Kind: Shell, Command: c}); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the planner cannot run") || !strings.Contains(v.Reason, want) || v.Call != "" {
			t.Errorf("%s: %+v", c, v)
		}
	}
	if v := Decide(planGuard, Call{Role: RolePlan, Kind: FlaiTool, Tool: "item_new", Args: Args{Type: "story"}}); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the planner cannot create a story that is not a draft") {
		t.Errorf("item_new of a story that is not a draft: %+v", v)
	}
	for _, c := range []string{
		"flai touches S-0001 --add flai/x.go", "flai touches suggest S-0001", "flai edit S-0001 --title x",
		"flai move S-0001 backlog", "flai story new x --draft", "flai epic new x", "flai task new --story S-0001 x",
		"flai thread new S-0001 x", "flai thread reply TH-0001 x", "flai issue new x", "flai issue bump I-0001",
	} {
		if v := Decide(planGuard, Call{Role: RolePlan, Kind: Shell, Command: c}); v.Outcome != Allow {
			t.Errorf("%s refused: %s", c, v.Reason)
		}
	}
}

// S-0416: the analyzer and the curator are held to their role files too,
// the analyzer's files entry under the design folder its reports folder
// lies in when no design folder is given, and their refusals keep the
// wording of their rules in Go where those refuse the call too.
func TestTheAnalyzerAndTheCuratorAreHeldToTheirRoleFiles(t *testing.T) {
	root := t.TempDir()
	elsewhere := Guard{Commands: g.Commands, Role: RoleAnalyze, Root: root, Design: "design", RoleFile: builtInFile(t, "analyzer", `["{design}/analysis/"]`, `["{design}/findings/"]`)}
	write := func(file string) Call {
		return Call{Role: RoleAnalyze, Kind: FileWrite, Name: "Write", Path: filepath.Join(root, file)}
	}
	if v := Decide(elsewhere, write("design/findings/r.md")); v.Outcome != Allow {
		t.Errorf("a write under the file's folder: %+v", v)
	}
	// its rule in Go, given no reports folder, refuses the write too, and
	// words the refusal
	if v := Decide(elsewhere, write("design/analysis/r.md")); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the analyzer cannot use Write on ") {
		t.Errorf("a write under analysis/ the file leaves out: %+v", v)
	}
	elsewhere.Reports = "design/analysis"
	if v := Decide(elsewhere, write("design/analysis/r.md")); v.Outcome != Deny || !strings.Contains(v.Reason, "it edits only a file under design/findings/") {
		t.Errorf("a write under the reports folder the file leaves out: %+v", v)
	}
	reports := Guard{Commands: g.Commands, Role: RoleAnalyze, Root: root, Reports: "design/analysis"}
	if v := Decide(reports, write("design/analysis/r.md")); v.Outcome != Allow {
		t.Errorf("a write under the reports folder: %+v", v)
	}
	analyzer := Guard{Commands: g.Commands, Role: RoleAnalyze}
	for _, c := range []string{"flai issue new x --class defect", "flai issue bump I-0001"} {
		if v := Decide(analyzer, Call{Role: RoleAnalyze, Kind: Shell, Command: c}); v.Outcome != Allow {
			t.Errorf("%s refused: %s", c, v.Reason)
		}
	}
	if v := Decide(analyzer, Call{Role: RoleAnalyze, Kind: Shell, Command: "flai issue story I-0001"}); v.Outcome != Deny || !strings.Contains(v.Reason, authorsNoStories) {
		t.Errorf("flai issue story: %+v", v)
	}

	curator := Guard{Commands: append(slices.Clone(g.Commands), "digest"), Role: RoleCurate}
	for _, sub := range []string{"harvest", "new x", "edit D-0001", "absorb D-0001 --by design/system/x.md", "cited", "compact"} {
		if v := Decide(curator, Call{Role: RoleCurate, Kind: Shell, Command: "flai digest " + sub}); v.Outcome != Allow {
			t.Errorf("flai digest %s refused: %s", sub, v.Reason)
		}
	}
	if v := Decide(curator, Call{Role: RoleCurate, Kind: Shell, Command: "flai digest seed"}); v.Outcome != Deny || !strings.Contains(v.Reason, seedsUnjudged) {
		t.Errorf("flai digest seed: %+v", v)
	}
	curator.RoleFile = builtInFile(t, "curator", "  - flai digest compact\n", "")
	if v := Decide(curator, Call{Role: RoleCurate, Kind: Shell, Command: "flai digest compact"}); v.Outcome != Deny || !strings.HasPrefix(v.Reason, `the role curator cannot run "flai digest compact": its role file lists no such call.`) {
		t.Errorf("flai digest compact left out of the file: %+v", v)
	}
}

// S-0416: a built-in role whose project file does not parse reads and is
// refused every other call, naming the file and flai check, rather than
// falling back on the baseline file.
func TestABuiltInRoleWhoseFileDoesNotParseFailsClosed(t *testing.T) {
	problem := "design/workflows/roles/planner.md does not parse: front matter does not parse; run flai check"
	broken := Guard{Commands: g.Commands, Role: RolePlan, Design: "design", RoleProblem: problem}
	for _, c := range []Call{{Role: RolePlan, Kind: FlaiTool, Tool: "board"}, {Role: RolePlan, Kind: Shell, Command: "flai board && git status"}} {
		if v := Decide(broken, c); v.Outcome != Allow {
			t.Errorf("read %+v refused: %s", c, v.Reason)
		}
	}
	for _, c := range []Call{
		{Role: RolePlan, Kind: FlaiTool, Tool: "thread_open", Args: Args{ID: "S-0001"}},
		{Role: RolePlan, Kind: FlaiTool, Tool: "item_new", Args: Args{Type: "story", Draft: true}},
		{Role: RolePlan, Kind: Shell, Command: "flai story new x --draft"},
	} {
		v := Decide(broken, c)
		if v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the role planner cannot ") || !strings.Contains(v.Reason, ": "+problem+". ") {
			t.Errorf("%+v: %+v", c, v)
		}
	}
}

// orchestrating is a guard in an orchestrator session under permissions p,
// held to file, flai's baseline file when nil, its threads told apart by
// orchestratorOpeners.
func orchestrating(p manifest.Permissions, file *roles.Role) Guard {
	return Guard{Commands: append(slices.Clone(g.Commands), "plan", "queue", "subproject"), Role: RoleOrchestrate, Permissions: p,
		Design: "design", RoleFile: file, Opener: orchestratorOpeners, Agent: "system-flow"}
}

// orchestratorOpeners knows TH-0001, which the orchestrator opened, TH-0002,
// a story's planner's, and TH-0003, which the session named by the hook's
// FLAI_AGENT, system-flow, opened.
func orchestratorOpeners(id string) (string, error) {
	switch id {
	case "TH-0001":
		return "orchestrator", nil
	case "TH-0002":
		return "planner-S-0002", nil
	case "TH-0003":
		return "system-flow", nil
	}
	return "", errors.New("no such thread")
}

// S-0416, ADR-0207: the orchestrator is held to its role file: a call its
// rule in Go allowed is refused when the file leaves it out, with the file's
// reason, the call, and, when an entry under a setting that is off would
// allow it, that setting; and a call the file adds passes, where its rule in
// Go refused it.
func TestTheOrchestratorIsHeldToItsRoleFile(t *testing.T) {
	publish := Call{Role: RoleOrchestrate, Kind: FlaiTool, Tool: ReleasePublish}
	noPublish := orchestrating(allOn, builtInFile(t, "orchestrator", "  - setting: orchestration.permissions.publish\n    tools: [release_publish]\n", ""))
	want := "the role orchestrator cannot call release_publish: its role file lists no such call. Ask the operator with thread_open on the item if it needs doing (design/workflows/roles/orchestrator.md, ADR-0207)."
	if v := Decide(noPublish, publish); v.Outcome != Deny || v.Reason != want || v.Call != ReleasePublish || v.Needs != "" {
		t.Errorf("release_publish left out of the file, publish on: %+v", v)
	}
	if v := Decide(orchestrating(allOn, nil), publish); v.Outcome != Allow {
		t.Errorf("release_publish under the baseline file, publish on: %+v", v)
	}

	open := Call{Role: RoleOrchestrate, Kind: FlaiTool, Tool: "thread_open"}
	noOpen := orchestrating(allOn, builtInFile(t, "orchestrator", "  - thread_open\n", ""))
	if v := Decide(noOpen, open); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the role orchestrator cannot call thread_open: its role file lists no such call.") || v.Call != "thread_open" {
		t.Errorf("thread_open left out of the file: %+v", v)
	}
	gated := builtInFile(t, "orchestrator", "  - thread_open\n", "", "    tools: [release_publish]\n", "    tools: [release_publish, thread_open]\n")
	v := Decide(orchestrating(except(manifest.PermitPublish), gated), open)
	if v.Outcome != Deny || !strings.Contains(v.Reason, "cannot call thread_open: it needs orchestration.permissions.publish, which is off.") || v.Call != "thread_open" || v.Needs != "orchestration.permissions.publish" {
		t.Errorf("thread_open under publish, publish off: %+v", v)
	}
	if v := Decide(orchestrating(allOn, gated), open); v.Outcome != Allow {
		t.Errorf("thread_open under publish, publish on: %+v", v)
	}

	issue := Call{Role: RoleOrchestrate, Kind: FlaiTool, Tool: "issue_new"}
	if v := Decide(orchestrating(allOn, nil), issue); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the orchestrator cannot call issue_new: the orchestrator never does it") || v.Call != "issue_new" {
		t.Errorf("issue_new under the baseline file: %+v", v)
	}
	if v := Decide(orchestrating(manifest.Permissions{}, builtInFile(t, "orchestrator", "  - inbox\n", "  - inbox\n  - issue_new\n")), issue); v.Outcome != Allow {
		t.Errorf("issue_new added to the file: %+v", v)
	}
}

// S-0416: under flai's baseline file, the orchestrator makes each call the
// file lists while its setting is on, as its rules in Go allowed.
func TestTheOrchestratorsBaselineFileAllowsWhatItLists(t *testing.T) {
	input := map[string]json.RawMessage{"revenue_per_week": json.RawMessage("500")}
	tools := []Call{
		{Tool: "activity_log"}, {Tool: "inbox"}, {Tool: "thread_open"}, {Tool: "wait_for_events"},
		{Tool: "thread_reply", Args: Args{ID: "TH-0001"}}, {Tool: "thread_resolve", Args: Args{ID: "TH-0001"}},
		{Tool: "plan", Args: Args{ID: "E-0001"}}, {Tool: "plan", Args: Args{ID: "S-0001"}},
		{Tool: "item_edit", Args: Args{ID: "S-0001", Fields: []string{"cost_of_delay", "id"}, CostOfDelay: input}},
		{Tool: "thread_reply", Args: Args{ID: "TH-0002", Source: "design/x.md"}}, {Tool: "thread_resolve", Args: Args{ID: "TH-0002"}},
		{Tool: "item_edit", Args: Args{ID: "S-0001", Fields: []string{"draft", "id"}}},
		{Tool: "item_move", Args: Args{ID: "S-0001", To: Ready}},
		{Tool: "thread_reply", Args: Args{ID: "TH-0003", Recommendation: true}},
		{Tool: QueueAdd, Args: Args{Story: "S-0001"}}, {Tool: ReleasePublish}, {Tool: SubprojectEdit},
	}
	commands := []string{
		"flai issue new x", "flai issue bump I-0001", "flai queue", "flai thread new S-0001 x",
		"flai thread reply TH-0001 x --by orchestrator", "flai thread resolve TH-0001", "flai plan E-0001", "flai plan S-0001",
		"flai edit S-0001 --revenue-per-week 500", "flai thread reply TH-0002 x", "flai edit S-0001 --no-draft",
		"flai move S-0001 ready", "flai order --by wsjf --apply", "flai thread reply TH-0003 x --recommend",
		"flai accept S-0001 --by orchestrator", "flai move S-0001 done --by orchestrator", "flai queue add S-0001 --by orchestrator",
		"flai subproject add web --kind go --autocommit",
	}
	for _, gr := range []Guard{orchestrating(allOn, nil), orchestrating(allOn, builtInFile(t, "orchestrator"))} {
		for _, c := range tools {
			c.Role, c.Kind = RoleOrchestrate, FlaiTool
			if v := Decide(gr, c); v.Outcome != Allow {
				t.Errorf("%s %+v refused: %s", c.Tool, c.Args, v.Reason)
			}
		}
		for _, c := range commands {
			if v := Decide(gr, Call{Role: RoleOrchestrate, Kind: Shell, Command: c}); v.Outcome != Allow {
				t.Errorf("%s refused: %s", c, v.Reason)
			}
		}
	}
}

// S-0416, S-0220, ADR-0093: the orchestrator's own name under its role file
// is orchestrator alone, not the hook's FLAI_AGENT: a thread that name opened
// is another's, and it neither accepts, enqueues, nor writes to a thread
// naming it.
func TestTheOrchestratorsOwnNameIsOrchestratorAlone(t *testing.T) {
	gr := orchestrating(allOn, nil)
	for c, want := range map[string]string{
		"flai accept S-0001 --by system-flow":          acceptsAsSelf,
		"flai move S-0001 done --by system-flow":       acceptsAsSelf,
		"flai queue add S-0001 --by system-flow":       enqueuesAsSelf,
		"flai thread reply TH-0001 x --by system-flow": writesAsSelf,
		"flai thread resolve TH-0003":                  resolvesOwn,
	} {
		if v := Decide(gr, Call{Role: RoleOrchestrate, Kind: Shell, Command: c}); v.Outcome != Deny || !strings.Contains(v.Reason, want) || v.Call == "" {
			t.Errorf("%s: %+v", c, v)
		}
	}
}

// S-0416: a defined role's call on a thread under own-thread names no author
// with --by but the role itself; the file refuses another, before ADR-0210's
// rule is reached.
func TestADefinedRolesThreadConditionsHoldItToItsOwnName(t *testing.T) {
	r, problem := roles.Parse(reviewerPath, strings.Replace(reviewerFile, "  - flai issue new\n", "  - flai issue new\n  - {call: flai thread reply, when: own-thread}\n", 1))
	if problem != nil {
		t.Fatal(problem)
	}
	gr := reviewer(t, t.TempDir(), manifest.Permissions{})
	gr.RoleFile = &r
	for _, c := range []string{"flai thread reply TH-0001 x", "flai thread reply TH-0001 x --by reviewer"} {
		if v := Decide(gr, roleShell(c)); v.Outcome != Allow {
			t.Errorf("%s refused: %s", c, v.Reason)
		}
	}
	if v := Decide(gr, roleShell("flai thread reply TH-0001 x --by alex")); v.Outcome != Deny || !strings.Contains(v.Reason, "only under own-thread") || strings.Contains(v.Reason, "ADR-0143") {
		t.Errorf("a reply on its own thread naming another: %+v", v)
	}
}
