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
		{roleTool("item_move", Args{ID: "S-0001", To: "review"}), "the role reviewer cannot move S-0001 to review: it may make it only under to-backlog (a move to backlog and nowhere else), or to-ready (a move to ready and nowhere else).", "item_move S-0001 review"},
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
