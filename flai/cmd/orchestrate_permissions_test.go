//go:build !windows

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// orchestration is a fixture board the orchestrator's calls run on, each
// through flai guard and then, when the guard lets it through, through flai
// (S-0219, T-0901). Built afresh for each case, it holds:
//
//   - E-0001, an epic in the backlog with no stories, and E-0002, an open
//     epic whose one story, S-0001, is done: the planner's candidates;
//     E-0003, open with stories that are not, is not one
//   - S-0002, a complete draft, and S-0003, a draft that lacks a forecast
//     and a cost of delay
//   - S-0004 and S-0005, the two candidates to go to ready, in that order by
//     cod; S-0006, held by its overlap with S-0007, in progress
//   - S-0008, S-0009, and S-0010 in ready, whose limit is 4, one short;
//     S-0008, of the lowest cost of delay, placed first by hand an hour ago
type orchestration struct {
	t        *testing.T
	root     string
	manifest []byte // the manifest as the fixture leaves it, without orchestration
}

// orchestrationNow is the time the orchestrator's calls run at, runIn's.
var orchestrationNow = time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)

// via is how the orchestrator makes a call: a flai command line, or flai's
// MCP tool.
type via string

const (
	viaBash via = "Bash"
	viaMCP  via = "MCP"
)

func newOrchestration(t *testing.T) *orchestration {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	t.Setenv("FLAI_ROLE", "")
	t.Setenv("FLAI_STARTED_BY", "")
	root := tempProject(t)
	run := func(args ...string) {
		t.Helper()
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	file := func(id string) string {
		t.Helper()
		files, _ := filepath.Glob(filepath.Join(root, "wip/kanban/*", id+"-*.md"))
		if len(files) != 1 {
			t.Fatalf("%s: %v", id, files)
		}
		return files[0]
	}
	rewrite := func(id, old, new string) {
		t.Helper()
		path := file(id)
		s, _ := os.ReadFile(path)
		if err := os.WriteFile(path, []byte(strings.Replace(string(s), old, new, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// a story made with a goal and a criterion, the touches given, and, as
	// figures says, a cost of delay value and a forecast duration
	story := func(title, epic, touches string, draft bool, figures ...string) {
		t.Helper()
		args := []string{"story", "new", title, "--epic", epic, "--touches", touches}
		if draft {
			args = append(args, "--draft")
		}
		out, errOut, code := runIn(t, root, args...)
		if code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		id := strings.Fields(out)[0]
		rewrite(id, "## Goal\n", "## Goal\n\nDo it.\n")
		rewrite(id, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] it works\n")
		if len(figures) > 0 {
			run(append([]string{"edit", id}, figures...)...)
		}
	}
	planned := func(value string) []string {
		return []string{"--cost-of-delay-value", value, "--forecast-duration", "2h"}
	}

	run("epic", "new", "Unplanned")
	run("epic", "new", "Shipped")
	run("epic", "new", "Work")
	story("Shipped slice", "E-0002", "docs/shipped", false)
	// the moves to done are the operator's, beside the point here
	rewrite("S-0001", "status: backlog\n", "status: done\n")
	rewrite("E-0002", "status: backlog\n", "status: in-progress\n")
	story("Complete draft", "E-0003", "docs/complete", true, append(planned("300"), "--forecast-delivery", "2026-09-20T12:00:00Z")...)
	story("Incomplete draft", "E-0003", "docs/incomplete", true)
	story("First candidate", "E-0003", "scripts", false, planned("500")...)
	story("Second candidate", "E-0003", "design", false, planned("400")...)
	story("Held", "E-0003", "flai/cmd", false, planned("900")...)
	story("Under way", "E-0003", "flai", false, planned("50")...)
	story("Hand placed", "E-0003", "examples/a", false, planned("100")...)
	story("Low", "E-0003", "examples/b", false, planned("200")...)
	story("High", "E-0003", "examples/c", false, planned("300")...)
	run("move", "S-0007", "ready")
	run("move", "S-0007", "in-progress")
	run("move", "S-0009", "ready")
	run("move", "S-0010", "ready")
	run("move", "S-0008", "ready")
	run("board", "limit", "ready", "4")
	if _, errOut, code := runInAt(t, root, orchestrationNow.Add(-time.Hour), "order", "S-0008", "--top", "--placed-by", "alex"); code != 0 {
		t.Fatalf("placing S-0008 by hand: %s", errOut)
	}
	m, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	o := &orchestration{t: t, root: root, manifest: m}
	if got := o.ready(); !reflect.DeepEqual(got, []string{"S-0008", "S-0009", "S-0010"}) {
		t.Fatalf("the fixture's ready column: %v", got)
	}
	return o
}

// permit gives the orchestrator the permissions on, and no other, by the
// cod policy, and makes the session the orchestrator's from here on.
func (o *orchestration) permit(on ...string) {
	o.t.Helper()
	block := "orchestration:\n  policy: cod\n"
	if len(on) > 0 {
		block += "  permissions:\n"
	}
	for _, p := range on {
		block += "    " + p + ": true\n"
	}
	if err := os.WriteFile(filepath.Join(o.root, "system-flow.yaml"), append(append([]byte{}, o.manifest...), block...), 0o644); err != nil {
		o.t.Fatal(err)
	}
	o.t.Setenv("FLAI_ROLE", "orchestrate")
	o.t.Setenv("FLAI_AGENT", "orchestrator")
}

// allBut is the four permissions this test covers but off.
func allBut(off string) []string {
	var on []string
	for _, p := range []string{manifest.PermitPlanBacklogEpics, manifest.PermitFinalizeDrafts, manifest.PermitPromoteToReady, manifest.PermitOrderReady} {
		if p != off {
			on = append(on, p)
		}
	}
	return on
}

// call makes the orchestrator's call through flai guard, as the hook before
// Bash and flai's MCP tools does, and, when the guard lets it through,
// through flai: the guard's refusal, else what flai printed or returned, or
// why flai refused it.
func (o *orchestration) call(tool string, input map[string]any) (refused, out, failed string) {
	o.t.Helper()
	name := tool
	if tool != string(viaBash) {
		name = "mcp__flai__" + tool
	}
	event, _ := json.Marshal(map[string]any{"tool_name": name, "tool_input": input})
	if _, errOut, code := runStdin(o.t, o.root, string(event), "guard"); code != 0 {
		if code != exitGuardRefused || errOut == "" {
			o.t.Fatalf("guard on %s: exit %d %q", event, code, errOut)
		}
		return errOut, "", ""
	}
	if tool == string(viaBash) {
		words := strings.Fields(input["command"].(string))
		if words[0] != "flai" {
			o.t.Fatalf("not a flai command: %v", words)
		}
		stdout, errOut, code := runIn(o.t, o.root, words[1:]...)
		if code != 0 {
			return "", stdout, errOut
		}
		return "", stdout, ""
	}
	return "", "", o.mcp(tool, input)
}

// mcp calls flai's MCP tool as the orchestrator's flai mcp serves it, and
// says why it failed, or "".
func (o *orchestration) mcp(tool string, input map[string]any) string {
	o.t.Helper()
	repo, err := workitem.Open(o.root)
	if err != nil {
		o.t.Fatal(err)
	}
	clock := func() time.Time { return orchestrationNow }
	a := &app{out: &bytes.Buffer{}, errOut: &bytes.Buffer{}, cwd: o.root, clock: clock}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct, st := mcp.NewInMemoryTransports()
	if _, err := mcpserver.New(mcpserver.Options{Repo: repo, Agent: "orchestrator", Version: "test", Now: clock, Plans: a.mcpPlan}).Connect(ctx, st, nil); err != nil {
		o.t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "orchestrator", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		o.t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: input})
	if err != nil {
		o.t.Fatalf("%s: %v", tool, err)
	}
	failed := ""
	if res.IsError {
		for _, c := range res.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				failed += tc.Text
			}
		}
		if failed == "" {
			failed = "failed"
		}
	}
	return failed
}

// passes makes a call the guard must let through and flai carry out, and
// returns what flai printed.
func (o *orchestration) passes(tool string, input map[string]any) string {
	o.t.Helper()
	refused, out, failed := o.call(tool, input)
	if refused != "" || failed != "" {
		o.t.Fatalf("%s %v: refused %q, failed %q", tool, input, refused, failed)
	}
	return out
}

// flaiRefuses makes a call the guard must let through and flai refuse,
// saying what the returned func is given.
func (o *orchestration) flaiRefuses(tool string, input map[string]any) func(want string) {
	return func(want string) {
		o.t.Helper()
		refused, out, failed := o.call(tool, input)
		if refused != "" || !strings.Contains(failed, want) {
			o.t.Errorf("%s %v: refused by the guard %q, by flai %q (out %q), want flai to say %q", tool, input, refused, failed, out, want)
		}
	}
}

// guardRefuses makes a call the guard must refuse, saying want, and checks
// it changed nothing but the orchestrator's activity document, which logs
// the refusal with needs, the permission named, or "" for a call the
// orchestrator never makes; want and needs are what the returned func is
// given.
func (o *orchestration) guardRefuses(tool string, input map[string]any) func(want, needs string) {
	return func(want, needs string) {
		o.t.Helper()
		o.guardRefusal(tool, input, want, needs)
	}
}

func (o *orchestration) guardRefusal(tool string, input map[string]any, want, needs string) {
	o.t.Helper()
	activity := filepath.Join(o.root, "wip", "agents", "orchestrator.md")
	before := treeOf(o.t, o.root)
	refused, out, failed := o.call(tool, input)
	if !strings.Contains(refused, want) {
		o.t.Errorf("%s %v: guard said %q, want %q (flai: %q %q)", tool, input, refused, want, out, failed)
		return
	}
	after := treeOf(o.t, o.root)
	delete(before, activity)
	delete(after, activity)
	if !reflect.DeepEqual(after, before) {
		o.t.Errorf("%s %v: a refused call changed the board", tool, input)
	}
	doc, err := workitem.ReadActivity(activity)
	if err != nil {
		o.t.Fatal(err)
	}
	if len(doc.Refusals) == 0 || doc.Refusals[len(doc.Refusals)-1].Needs != needs {
		o.t.Errorf("%s %v: refusals logged %+v, want the last to need %q", tool, input, doc.Refusals, needs)
	}
}

// offRefusal is what the guard says of a call while permission is off.
func offRefusal(permission string) (want, needs string) {
	needs = "orchestration.permissions." + permission
	return "it needs " + needs + ", which is off", needs
}

// item is an item of the fixture as it stands.
func (o *orchestration) item(id string) *workitem.Item {
	o.t.Helper()
	repo, err := workitem.Open(o.root)
	if err != nil {
		o.t.Fatal(err)
	}
	it, err := repo.Get(id)
	if err != nil {
		o.t.Fatal(err)
	}
	return it
}

// ready is the ready column in its pull order.
func (o *orchestration) ready() (ids []string) {
	o.t.Helper()
	out, errOut, code := runIn(o.t, o.root, "board", "--json")
	if code != 0 {
		o.t.Fatalf("board: %s", errOut)
	}
	var v struct {
		Columns map[string][]struct {
			ID string `json:"id"`
		} `json:"columns"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		o.t.Fatal(err)
	}
	for _, c := range v.Columns["ready"] {
		ids = append(ids, c.ID)
	}
	return ids
}

// planCall, finalizeCall, and promoteCall are the orchestrator's calls for
// its first three permissions, made as v says.
func planCall(v via, id string) (string, map[string]any) {
	if v == viaMCP {
		return "plan", map[string]any{"id": id}
	}
	return string(viaBash), map[string]any{"command": "flai plan " + id}
}

func finalizeCall(v via, id string) (string, map[string]any) {
	if v == viaMCP {
		return "item_edit", map[string]any{"id": id, "draft": false}
	}
	return string(viaBash), map[string]any{"command": "flai edit " + id + " --no-draft"}
}

func promoteCall(v via, id string) (string, map[string]any) {
	if v == viaMCP {
		return "item_move", map[string]any{"id": id, "to": "ready"}
	}
	return string(viaBash), map[string]any{"command": "flai move " + id + " ready"}
}

// bash is a flai command line as the orchestrator's Bash call.
func bash(command string) (string, map[string]any) {
	return string(viaBash), map[string]any{"command": command}
}

// S-0219, T-0901: on the fixture board, the orchestrator with
// plan_backlog_epics asks for the planner on a candidate epic, by flai plan
// and the MCP tool plan alike, and flai refuses an epic that is not one;
// without it the guard refuses both naming the permission, flai would too,
// and nothing changes; flai plan --candidates it runs either way, and a plan
// for a story never.
func TestTheOrchestratorPlansACandidateEpicOnlyWithPlanBacklogEpics(t *testing.T) {
	const notOne = "the orchestrator asks for the planner on an epic flai plan --candidates lists alone, and E-0003 is not one"
	t.Run("on", func(t *testing.T) {
		o := newOrchestration(t)
		for _, args := range [][]string{{"serve", "enable", "plan"}, {"serve", "agent", "set", "--", "true"}} {
			if _, errOut, code := runIn(t, o.root, args...); code != 0 {
				t.Fatalf("flai %v: %s", args, errOut)
			}
		}
		o.permit(manifest.PermitPlanBacklogEpics)
		out := o.passes(bash("flai plan --candidates"))
		if !strings.Contains(out, "  E-0001  Unplanned\n") || !strings.Contains(out, "  E-0002  Shipped\n") || strings.Contains(out, "E-0003") {
			t.Errorf("the candidates:\n%s", out)
		}
		o.flaiRefuses(planCall(viaBash, "E-0003"))(notOne)
		if out := o.passes(planCall(viaBash, "E-0001")); !strings.Contains(out, "to plan E-0001 as planner-E-0001") {
			t.Errorf("flai plan on a candidate: %s", out)
		}
		o.flaiRefuses(planCall(viaMCP, "E-0003"))(notOne)
		o.passes(planCall(viaMCP, "E-0002"))
		js, _, _ := runIn(t, o.root, "serve", "journal", "--json")
		if !strings.Contains(js, "orchestrator asked to plan E-0002: started true as planner-E-0002") {
			t.Errorf("the MCP tool plan did not start the planner on E-0002: %s", js)
		}
		o.guardRefuses(planCall(viaMCP, "S-0004"))("the orchestrator never does it", "")
	})
	for _, v := range []via{viaBash, viaMCP} {
		t.Run("off by "+string(v), func(t *testing.T) {
			o := newOrchestration(t)
			o.permit(allBut(manifest.PermitPlanBacklogEpics)...)
			want, needs := offRefusal(manifest.PermitPlanBacklogEpics)
			o.guardRefuses(planCall(v, "E-0001"))(want, needs)
			before := treeOf(t, o.root)
			if out := o.passes(bash("flai plan --candidates")); !strings.Contains(out, "  E-0001  Unplanned\n") {
				t.Errorf("flai plan --candidates, which reads: %s", out)
			}
			// flai, past the guard, holds the orchestrator to the permission too
			if _, errOut, code := runIn(t, o.root, "plan", "E-0001"); code == 0 || !strings.Contains(errOut, "only with orchestration.permissions.plan_backlog_epics, which is off") {
				t.Errorf("flai plan past the guard: exit %d %s", code, errOut)
			}
			if !reflect.DeepEqual(treeOf(t, o.root), before) {
				t.Error("the reads, or flai's refusal, changed the board")
			}
		})
	}
}

// S-0219, T-0901: on the fixture board, the orchestrator with finalize_drafts
// finalizes the complete draft, by flai edit --no-draft and the MCP tool
// item_edit alike, and is named as who finalized it, while flai refuses the
// incomplete one naming what it lacks; without it the guard refuses
// naming the permission and nothing changes; an edit that does more than
// finalize it never makes.
func TestTheOrchestratorFinalizesACompleteDraftOnlyWithFinalizeDrafts(t *testing.T) {
	for _, v := range []via{viaBash, viaMCP} {
		t.Run("on by "+string(v), func(t *testing.T) {
			o := newOrchestration(t)
			o.permit(manifest.PermitFinalizeDrafts)
			out := o.passes(bash("flai promote --drafts"))
			if !strings.Contains(out, "S-0002  complete") || !strings.Contains(out, "S-0003  incomplete") {
				t.Errorf("the drafts:\n%s", out)
			}
			o.flaiRefuses(finalizeCall(v, "S-0003"))("S-0003 is not complete, so the orchestrator does not finalize it (flai promote --drafts): no forecast duration; no forecast delivery; no cost of delay value")
			if it := o.item("S-0003"); !it.Draft || it.Finalized != nil {
				t.Errorf("S-0003 finalized: %+v", it.Finalized)
			}
			o.passes(finalizeCall(v, "S-0002"))
			if it := o.item("S-0002"); it.Draft || it.Finalized == nil || it.Finalized.By != "orchestrator" {
				t.Errorf("S-0002 finalized by %+v, draft %v", it.Finalized, it.Draft)
			}
			if v == viaMCP {
				o.guardRefusal("item_edit", map[string]any{"id": "S-0003", "draft": false, "title": "Renamed"}, "the orchestrator never does it", "")
			} else {
				o.guardRefuses(bash("flai edit S-0003 --no-draft --title Renamed"))("the orchestrator never does it", "")
			}
		})
		t.Run("off by "+string(v), func(t *testing.T) {
			o := newOrchestration(t)
			o.permit(allBut(manifest.PermitFinalizeDrafts)...)
			want, needs := offRefusal(manifest.PermitFinalizeDrafts)
			o.guardRefuses(finalizeCall(v, "S-0002"))(want, needs)
			if it := o.item("S-0002"); !it.Draft || it.Finalized != nil {
				t.Errorf("S-0002 finalized: %+v", it.Finalized)
			}
			before := treeOf(t, o.root)
			if out := o.passes(bash("flai promote --drafts")); !strings.Contains(out, "S-0002  complete") {
				t.Errorf("flai promote --drafts, which reads: %s", out)
			}
			// flai, past the guard, holds the orchestrator to the permission too
			if _, errOut, code := runIn(t, o.root, "edit", "S-0002", "--no-draft"); code == 0 || !strings.Contains(errOut, "only with orchestration.permissions.finalize_drafts, which is off") {
				t.Errorf("flai edit --no-draft past the guard: exit %d %s", code, errOut)
			}
			if !reflect.DeepEqual(treeOf(t, o.root), before) {
				t.Error("flai promote --drafts changed the board")
			}
		})
	}
}

// S-0219, T-0901: on the fixture board, the orchestrator with
// promote_to_ready moves the first candidate to ready, by flai move and the
// MCP tool item_move alike, and flai refuses the second at the ready
// column's limit, and the held story and the drafts whenever; without it the
// guard refuses naming the permission and nothing changes.
func TestTheOrchestratorPromotesACandidateOnlyWithPromoteToReady(t *testing.T) {
	const notCandidate = " is not a candidate to go to ready (flai promote --candidates): "
	for _, v := range []via{viaBash, viaMCP} {
		t.Run("on by "+string(v), func(t *testing.T) {
			o := newOrchestration(t)
			o.permit(manifest.PermitPromoteToReady)
			var got struct {
				Candidates []struct{ ID string } `json:"candidates"`
			}
			if err := json.Unmarshal([]byte(o.passes(bash("flai promote --candidates --json"))), &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Candidates) != 2 || got.Candidates[0].ID != "S-0004" || got.Candidates[1].ID != "S-0005" {
				t.Fatalf("the candidates: %+v", got.Candidates)
			}
			o.flaiRefuses(promoteCall(v, "S-0006"))("S-0006" + notCandidate + "held (overlap)")
			o.flaiRefuses(promoteCall(v, "S-0002"))("S-0002" + notCandidate + "draft: finalize it first")
			o.flaiRefuses(promoteCall(v, "S-0003"))("S-0003" + notCandidate + "draft: finalize it first")
			o.passes(promoteCall(v, "S-0004"))
			o.flaiRefuses(promoteCall(v, "S-0005"))("the ready column is at its WIP limit (4 of 4): the orchestrator moves no story to ready until one leaves it")
			o.flaiRefuses(promoteCall(v, "S-0006"))("S-0006" + notCandidate + "held (overlap)")
			for id, want := range map[string]string{"S-0004": workitem.Ready, "S-0005": workitem.Backlog, "S-0006": workitem.Backlog, "S-0002": workitem.Backlog, "S-0003": workitem.Backlog} {
				if it := o.item(id); it.Status != want {
					t.Errorf("%s is %s, want %s", id, it.Status, want)
				}
			}
			for _, id := range []string{"S-0002", "S-0003"} {
				if it := o.item(id); !it.Draft {
					t.Errorf("%s finalized by a move", id)
				}
			}
		})
		t.Run("off by "+string(v), func(t *testing.T) {
			o := newOrchestration(t)
			o.permit(allBut(manifest.PermitPromoteToReady)...)
			want, needs := offRefusal(manifest.PermitPromoteToReady)
			o.guardRefuses(promoteCall(v, "S-0004"))(want, needs)
			if it := o.item("S-0004"); it.Status != workitem.Backlog {
				t.Errorf("S-0004 is %s", it.Status)
			}
			before := treeOf(t, o.root)
			if out := o.passes(bash("flai promote --candidates")); !strings.Contains(out, "S-0004") {
				t.Errorf("flai promote --candidates, which reads: %s", out)
			}
			// flai, past the guard, holds the orchestrator to the permission too
			if _, errOut, code := runIn(t, o.root, "move", "S-0004", "ready"); code == 0 || !strings.Contains(errOut, "only with orchestration.permissions.promote_to_ready, which is off") {
				t.Errorf("flai move past the guard: exit %d %s", code, errOut)
			}
			if !reflect.DeepEqual(treeOf(t, o.root), before) {
				t.Error("flai promote --candidates changed the board")
			}
		})
	}
}

// S-0219, T-0901: on the fixture board, the orchestrator with order_ready
// orders the ready column by its policy with flai order --by cod --apply,
// and the story placed by hand an hour ago keeps its place; without it the
// guard refuses naming the permission and nothing changes; flai order --by
// without --apply it runs either way, and a placement by hand never.
func TestTheOrchestratorOrdersReadyOnlyWithOrderReady(t *testing.T) {
	placedByHand := map[string]workitem.Placed{"S-0008": {By: "alex", At: "2026-09-15T20:00:00Z"}}
	placed := func(o *orchestration) map[string]workitem.Placed {
		o.t.Helper()
		repo, err := workitem.Open(o.root)
		if err != nil {
			o.t.Fatal(err)
		}
		b, err := repo.LoadBoard()
		if err != nil {
			o.t.Fatal(err)
		}
		return b.Placed
	}
	const byHand = "the orchestrator never does it, whatever its permissions: it orders the ready column by its policy"
	t.Run("on", func(t *testing.T) {
		o := newOrchestration(t)
		o.permit(manifest.PermitOrderReady)
		o.guardRefuses(bash("flai order S-0010 --top"))(byHand, "")
		out := o.passes(bash("flai order --by cod --apply"))
		if !strings.Contains(out, "(kept: placed by alex at 2026-09-15T20:00:00Z)") || !strings.Contains(out, "board.md's ready order is now this one") {
			t.Errorf("flai order --by cod --apply:\n%s", out)
		}
		if got, want := o.ready(), []string{"S-0008", "S-0010", "S-0009"}; !reflect.DeepEqual(got, want) {
			t.Errorf("ready: %v, want %v: S-0008 kept first, the rest by cod", got, want)
		}
		if got := placed(o); !reflect.DeepEqual(got, placedByHand) {
			t.Errorf("placed: %+v, want %+v", got, placedByHand)
		}
	})
	t.Run("off", func(t *testing.T) {
		o := newOrchestration(t)
		o.permit(allBut(manifest.PermitOrderReady)...)
		want, needs := offRefusal(manifest.PermitOrderReady)
		o.guardRefuses(bash("flai order --by cod --apply"))(want, needs)
		o.guardRefuses(bash("flai order S-0010 --top"))(byHand, "")
		before := treeOf(t, o.root)
		if out := o.passes(bash("flai order --by cod")); !strings.Contains(out, "ready by cod:") || strings.Contains(out, "is now this one") {
			t.Errorf("flai order --by cod, which reads:\n%s", out)
		}
		// flai, past the guard, holds the orchestrator to the permission too
		if _, errOut, code := runIn(t, o.root, "order", "--by", "cod", "--apply"); code == 0 || !strings.Contains(errOut, "only with orchestration.permissions.order_ready, which is off") {
			t.Errorf("flai order --apply past the guard: exit %d %s", code, errOut)
		}
		if !reflect.DeepEqual(treeOf(t, o.root), before) {
			t.Error("flai order --by cod changed the board")
		}
		if got, want := o.ready(), []string{"S-0008", "S-0009", "S-0010"}; !reflect.DeepEqual(got, want) {
			t.Errorf("ready: %v, want %v, as it was", got, want)
		}
		if got := placed(o); !reflect.DeepEqual(got, placedByHand) {
			t.Errorf("placed: %+v, want %+v", got, placedByHand)
		}
	})
}
