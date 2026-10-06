package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0175, ADR-0060: flai guard refuses a sub-agent's writes with exit 2 and
// lets everything else through, the story's agent's calls and unreadable
// input included.
func TestGuard(t *testing.T) {
	dir := t.TempDir()
	for _, c := range []struct {
		in   string
		code int
		err  string
	}{
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"review"},"agent_type":"verifier","agent_id":"a1"}`, 2, "a sub-agent (verifier) cannot call item_move"},
		{`{"tool_name":"Bash","tool_input":{"command":"flai move S-1 review"},"agent_type":"Explore","agent_id":"a2"}`, 2, "a sub-agent (Explore) cannot run \"flai move S-1 review\""},
		{`{"tool_name":"Bash","tool_input":{"command":"make test"},"agent_type":"verifier","agent_id":"a3"}`, 0, ""},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"review"}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"timeout 9 flai accept S-1"},"agent_type":"verifier","agent_id":"a4"}`, 2, "cannot run \"timeout 9 flai accept S-1\""},
		{`not json`, 0, "hook input is not a tool call"},
	} {
		_, errOut, code := runStdin(t, dir, c.in, "guard")
		if code != c.code || (c.err == "" && errOut != "") || !strings.Contains(errOut, c.err) {
			t.Errorf("%s: code %d, stderr %q", c.in, code, errOut)
		}
	}
}

// S-0208: in a session flai serve starts with FLAI_ROLE=plan, flai guard
// holds the planner's own calls to planning.
func TestGuardHoldsThePlannerToPlanning(t *testing.T) {
	t.Setenv("FLAI_ROLE", "plan")
	dir := t.TempDir()
	for _, c := range []struct {
		in   string
		code int
		err  string
	}{
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"backlog"}}`, 0, ""},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"ready"}}`, 2, "the planner cannot move S-1 to ready"},
		{`{"tool_name":"Bash","tool_input":{"command":"flai story new --epic E-1 --draft x"}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"flai accept S-1"}}`, 2, "the planner cannot run \"flai accept S-1\""},
		{`{"tool_name":"Write","tool_input":{"file_path":"x.go","content":""}}`, 2, "the planner cannot use Write"},
		{`{"tool_name":"mcp__flai__item_new","tool_input":{"type":"story"},"agent_type":"explorer","agent_id":"a1"}`, 2, "a sub-agent (explorer) cannot call item_new"},
	} {
		_, errOut, code := runStdin(t, dir, c.in, "guard")
		if code != c.code || (c.err == "" && errOut != "") || !strings.Contains(errOut, c.err) {
			t.Errorf("%s: code %d, stderr %q", c.in, code, errOut)
		}
	}
}

// S-0223: in a session flai serve starts with FLAI_ROLE=analyze, flai guard
// holds the analyzer's own calls to reads and its report, under the
// manifest's design folder's analysis/, whatever the folder is called, and
// to filing and bumping issues but making no story of one (S-0224).
func TestGuardHoldsTheAnalyzerToItsReport(t *testing.T) {
	t.Setenv("FLAI_ROLE", "analyze")
	root := tempProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	data, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, []byte(strings.Replace(string(data), "design: design", "design: plans", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	edit := func(tool, file string) string {
		in, _ := json.Marshal(map[string]any{"tool_name": tool, "tool_input": map[string]string{"file_path": file, "content": ""}})
		return string(in)
	}
	for _, c := range []struct {
		in   string
		code int
		err  string
	}{
		{edit("Write", filepath.Join(root, "plans", "analysis", "2026-10-06-all.md")), 0, ""},
		{edit("Edit", filepath.Join(root, "plans", "analysis", "README.md")), 0, ""},
		{edit("Write", filepath.Join(root, "design", "analysis", "2026-10-06-all.md")), 2, "the analyzer cannot use Write on " + filepath.Join(root, "design", "analysis", "2026-10-06-all.md") + ": it edits only its report and the index, under plans/analysis/"},
		{edit("Edit", root+"/plans/analysis/../system/overview.md"), 2, "under plans/analysis/"},
		{`{"tool_name":"mcp__flai__item_new","tool_input":{"type":"story","draft":true}}`, 2, "the analyzer cannot call item_new: the analyzer authors no stories"},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"backlog"}}`, 2, "the analyzer cannot call item_move: the analyzer authors no stories"},
		{`{"tool_name":"mcp__flai__activity_log","tool_input":{"kind":"analyzer"}}`, 0, ""},
		{`{"tool_name":"mcp__flai__thread_open","tool_input":{"on":"S-1"}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"flai stats --json"}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"flai issue new x --class defect --report plans/analysis/2026-10-06-all.md"}}`, 0, ""},
		{`{"tool_name":"mcp__flai__issue_bump","tool_input":{"id":"I-1","report":"plans/analysis/2026-10-06-all.md"}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"flai issue story I-1"}}`, 2, `the analyzer cannot run "flai issue story I-1": the analyzer authors no stories`},
		{`{"tool_name":"mcp__flai__issue_story","tool_input":{"id":"I-1"}}`, 2, "the analyzer cannot call issue_story: the analyzer authors no stories"},
		{`{"tool_name":"Bash","tool_input":{"command":"flai issue close I-1 --reason x"}}`, 2, `the analyzer cannot run "flai issue close I-1 --reason x"`},
		{`{"tool_name":"mcp__flai__analyze","tool_input":{}}`, 2, "the analyzer cannot call analyze"},
		{`{"tool_name":"mcp__flai__item_new","tool_input":{"type":"story"},"agent_type":"explorer","agent_id":"a1"}`, 2, "a sub-agent (explorer) cannot call item_new"},
	} {
		_, errOut, code := runStdin(t, root, c.in, "guard")
		if code != c.code || (c.err == "" && errOut != "") || !strings.Contains(errOut, c.err) {
			t.Errorf("%s: code %d, stderr %q", c.in, code, errOut)
		}
	}
	// outside a project it may edit no file
	_, errOut, code := runStdin(t, t.TempDir(), edit("Write", filepath.Join(root, "plans", "analysis", "2026-10-06-all.md")), "guard")
	if code != 2 || !strings.Contains(errOut, "under the design folder's analysis/") || !strings.Contains(errOut, "the analyzer may edit no file") {
		t.Errorf("no project: code %d, stderr %q", code, errOut)
	}
}

// S-0218: in a session flai serve starts with FLAI_ROLE=orchestrate, flai
// guard holds the orchestrator's own calls to the permissions in the
// project's manifest, and logs each refusal in wip/agents/orchestrator.md.
func TestGuardHoldsTheOrchestratorToItsPermissions(t *testing.T) {
	t.Setenv("FLAI_ROLE", "orchestrate")
	root := tempProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	data, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, append(data, "orchestration:\n  permissions:\n    promote_to_ready: true\n    answer_threads: recommend\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		in   string
		code int
		err  string
	}{
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"ready"}}`, 0, ""},
		{`{"tool_name":"mcp__flai__thread_reply","tool_input":{"id":"TH-1","text":"I recommend S-2","recommendation":true}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"flai order --by wsjf && flai thread new --on S-1 'q' 'text'"}}`, 0, ""},
		{`{"tool_name":"Bash","tool_input":{"command":"flai accept S-1"}}`, 2, `the orchestrator cannot run "flai accept S-1": it needs orchestration.permissions.accept_reviews, which is off. Ask the operator with thread_open on the item`},
		{`{"tool_name":"Write","tool_input":{"file_path":"x.go","content":""}}`, 2, "the orchestrator cannot use Write: the orchestrator never does it"},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"ready"},"agent_type":"explorer","agent_id":"a1"}`, 2, "a sub-agent (explorer) cannot call item_move"},
	} {
		_, errOut, code := runStdin(t, root, c.in, "guard")
		if code != c.code || (c.err == "" && errOut != "") || !strings.Contains(errOut, c.err) {
			t.Errorf("%s: code %d, stderr %q", c.in, code, errOut)
		}
	}
	doc, err := workitem.ReadActivity(filepath.Join(root, "wip", "agents", "orchestrator.md"))
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 9, 19, 2, 0, 0, 0, time.UTC)
	want := []workitem.ActivityRefusal{{At: at, Call: "flai accept S-1", Needs: "orchestration.permissions.accept_reviews"}, {At: at, Call: "Write x.go"}}
	if !reflect.DeepEqual(doc.Refusals, want) || len(doc.Entries) != 0 {
		t.Errorf("the orchestrator's refusals: %+v, want %+v (a sub-agent's are not logged)", doc.Refusals, want)
	}

	// a refusal that cannot be logged is refused all the same
	if err := os.WriteFile(doc.Path, []byte("not an activity document\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, errOut, code := runStdin(t, root, `{"tool_name":"mcp__flai__release_publish","tool_input":{"reason":"r"}}`, "guard")
	if code != 2 || !strings.Contains(errOut, "it needs orchestration.permissions.publish") || !strings.Contains(errOut, "refusal not logged") {
		t.Errorf("unlogged refusal: code %d, stderr %q", code, errOut)
	}
	// outside a project every permission is off
	_, errOut, code = runStdin(t, t.TempDir(), `{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"ready"}}`, "guard")
	if code != 2 || !strings.Contains(errOut, "it needs orchestration.permissions.promote_to_ready") || !strings.Contains(errOut, "permissions are taken as off") {
		t.Errorf("no project: code %d, stderr %q", code, errOut)
	}
}

// I-0093, S-0299: flai guard refuses a sub-agent's write in a .claude/ folder
// with exit 2 while the project's auto-approve is off, so that it never waits
// on permission_prompt's thread, and lets it through once the operator
// enables auto-approve for the project; the story's agent's own write, and a
// sub-agent's outside a .claude/ folder, pass either way.
func TestGuardRefusesASubAgentsWriteUnderClaudeWhileAutoApproveIsOff(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_STORY", "S-0001")
	root := tempProject(t)
	file := filepath.Join(root, "template", "root", ".claude", "agents", "analyzer.md")
	write := func(agent, file string) string {
		e := map[string]any{"tool_name": "Write", "tool_input": map[string]string{"file_path": file, "content": "x"}}
		if agent != "" {
			e["agent_id"], e["agent_type"] = "a1", agent
		}
		in, _ := json.Marshal(e)
		return string(in)
	}
	_, errOut, code := runStdin(t, root, write("general-purpose", file), "guard")
	if code != 2 || !strings.Contains(errOut, "a sub-agent (general-purpose) cannot use Write on "+file) || !strings.Contains(errOut, "whole new content in your final message") {
		t.Errorf("auto-approve off: code %d, stderr %q", code, errOut)
	}
	for _, in := range []string{write("", file), write("general-purpose", filepath.Join(root, "flai", "x.go"))} {
		if _, errOut, code := runStdin(t, root, in, "guard"); code != 0 || errOut != "" {
			t.Errorf("%s: code %d, stderr %q", in, code, errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "auto-approve"); code != 0 {
		t.Fatalf("enable: %s", errOut)
	}
	if _, errOut, code := runStdin(t, root, write("general-purpose", file), "guard"); code != 0 || errOut != "" {
		t.Errorf("auto-approve on: code %d, stderr %q", code, errOut)
	}
}

// S-0208, TH-0096, S-0218: this repository's settings and the template's run
// the guard before Bash and flai's MCP tools in every session, and before
// Edit, Write, and NotebookEdit in a planner or orchestrator session alone.
func TestTheSettingsRunTheGuardOnTheStrategicAgentsEdits(t *testing.T) {
	for _, path := range []string{filepath.Join("..", "..", ".claude", "settings.json"), filepath.Join("..", "..", "template", "root", ".claude", "settings.json")} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var settings struct {
			Hooks map[string][]struct {
				Matcher string `json:"matcher"`
				Hooks   []struct {
					Command string `json:"command"`
				} `json:"hooks"`
			} `json:"hooks"`
		}
		if err := json.Unmarshal(data, &settings); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		commands := map[string]string{}
		for _, h := range settings.Hooks["PreToolUse"] {
			if len(h.Hooks) == 1 {
				commands[h.Matcher] = h.Hooks[0].Command
			}
		}
		if c := commands["Bash|mcp__flai__.*"]; !strings.Contains(c, " guard 2>&1") || strings.Contains(c, "FLAI_ROLE") {
			t.Errorf("%s: Bash and flai's MCP tools are not guarded in every session: %q", path, c)
		}
		if c := commands["Edit|Write|NotebookEdit"]; !strings.HasPrefix(c, `[ "$FLAI_ROLE" = plan ] || [ "$FLAI_ROLE" = orchestrate ] || [ "$FLAI_ROLE" = analyze ] || exit 0; `) || !strings.Contains(c, " guard 2>&1") {
			t.Errorf("%s: a strategic agent's file edits are not guarded, or a story's are: %q", path, c)
		}
		if len(commands) != 2 {
			t.Errorf("%s: hooks %v", path, commands)
		}
	}
}

// subAgentHook is a SubagentStart's or SubagentStop's input, as Claude Code
// 2.1.290 gives it, for the sub-agent id of type typ in session s1.
func subAgentHook(event, id, typ string) string {
	return `{"session_id":"s1","transcript_path":"/t.jsonl","cwd":"/tmp/x","prompt_id":"p1","agent_id":"` + id + `","agent_type":"` + typ + `","hook_event_name":"` + event + `","stop_hook_active":false}`
}

// waitHook is the story's agent's own wait_for_events in session s1.
const waitHook = `{"session_id":"s1","hook_event_name":"PreToolUse","tool_name":"mcp__flai__wait_for_events","tool_input":{"timeout_seconds":600}}`

// S-0285: the guard records each sub-agent SubagentStart reports as running
// in its session, in the main checkout's .flai-cache, and removes it on
// SubagentStop, saying nothing.
func TestGuardRecordsASessionsSubAgents(t *testing.T) {
	root := tempProject(t)
	dir := filepath.Join(root, ".flai-cache", "guard")
	for _, c := range []struct {
		in   string
		want []string
	}{
		{subAgentHook("SubagentStart", "ac73f87145abc4315", "general-purpose"), []string{"ac73f87145abc4315 general-purpose"}},
		{subAgentHook("SubagentStart", "b2", "verifier"), []string{"ac73f87145abc4315 general-purpose", "b2 verifier"}},
		{subAgentHook("SubagentStop", "ac73f87145abc4315", "general-purpose"), []string{"b2 verifier"}},
		{subAgentHook("SubagentStop", "b2", "verifier"), nil},
	} {
		out, errOut, code := runStdin(t, root, c.in, "guard")
		if code != 0 || out != "" || errOut != "" {
			t.Errorf("%s: code %d, stdout %q, stderr %q", c.in, code, out, errOut)
		}
		running, err := guard.Running(dir, "s1")
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, s := range running {
			got = append(got, s.ID+" "+s.Type)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("after %s: running %v, want %v", c.in, got, c.want)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "s1.json")); !os.IsNotExist(err) {
		t.Errorf("a record with no sub-agent running is left: %v", err)
	}
	// a session ID that is not a file name is ignored
	if _, errOut, code := runStdin(t, root, strings.Replace(subAgentHook("SubagentStart", "c3", "x"), `"s1"`, `"../s1"`, 1), "guard"); code != 0 || errOut != "" {
		t.Errorf("unsafe session: code %d, stderr %q", code, errOut)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("unsafe session recorded: %v", entries)
	}
	// outside a project it warns and fails open
	if _, errOut, code := runStdin(t, t.TempDir(), subAgentHook("SubagentStart", "c3", "x"), "guard"); code != 0 || !strings.Contains(errOut, "sub-agent not recorded") {
		t.Errorf("no project: code %d, stderr %q", code, errOut)
	}
}

// S-0285: in a story's agent's session the guard refuses the agent's own
// wait_for_events while a sub-agent of the session runs and no thread on the
// story or its tasks is open, naming the sub-agent and how to wait instead;
// a wait for the designer, with a thread open, passes.
func TestGuardRefusesTheStorysAgentsWaitOnItsSubAgents(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "agent-S-0001")
	t.Setenv("FLAI_ROLE", "")
	t.Setenv("FLAI_STORY", "S-0001")
	root := tempProject(t)
	for _, args := range [][]string{{"story", "new", "Slice"}, {"story", "new", "Other"}, {"task", "new", "--story", "S-0001", "Piece"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	wait := func(what string, code int, err string) {
		t.Helper()
		_, errOut, got := runStdin(t, root, waitHook, "guard")
		if got != code || (err == "" && errOut != "") || !strings.Contains(errOut, err) {
			t.Errorf("%s: code %d, stderr %q", what, got, errOut)
		}
	}
	hook := func(in string) {
		t.Helper()
		if _, errOut, code := runStdin(t, root, in, "guard"); code != 0 {
			t.Fatalf("%s: code %d, stderr %q", in, code, errOut)
		}
	}
	thread := func(args ...string) {
		t.Helper()
		if _, errOut, code := runIn(t, root, append([]string{"thread"}, args...)...); code != 0 {
			t.Fatalf("thread %v: %s", args, errOut)
		}
	}

	wait("no sub-agent started", 0, "")
	hook(subAgentHook("SubagentStart", "ac73f87145abc4315", "general-purpose"))
	wait("a sub-agent running", 2, "the story's agent cannot call wait_for_events while its sub-agents run (general-purpose ac73f87145abc4315) and no thread on S-0001 awaits the designer: wait_for_events reports work items and threads, not sub-agents, so it would run to its timeout. Wait for a sub-agent by launching it with the Agent tool's run_in_background set to false")
	thread("new", "--on", "S-0002", "Elsewhere?", "On another story.")
	wait("a thread open on another story", 2, "general-purpose ac73f87145abc4315")
	thread("new", "--on", "S-0001", "Which first?", "Say which.")
	wait("a thread open on the story", 0, "")
	thread("resolve", "TH-0002")
	wait("the story's thread resolved", 2, "general-purpose ac73f87145abc4315")
	thread("new", "--on", "T-0001", "Which piece?", "Say which.")
	wait("a thread open on the story's task", 0, "")
	thread("resolve", "TH-0003")
	wait("the task's thread resolved", 2, "general-purpose ac73f87145abc4315")

	// a sub-agent's own wait_for_events is refused as before
	_, errOut, code := runStdin(t, root, `{"session_id":"s1","hook_event_name":"PreToolUse","tool_name":"mcp__flai__wait_for_events","tool_input":{},"agent_id":"ac73f87145abc4315","agent_type":"general-purpose"}`, "guard")
	if code != 2 || !strings.Contains(errOut, "a sub-agent (general-purpose) cannot call wait_for_events") {
		t.Errorf("a sub-agent's wait: code %d, stderr %q", code, errOut)
	}
	// the orchestrator and a session outside a story wait as they did
	t.Setenv("FLAI_ROLE", "orchestrate")
	wait("the orchestrator", 0, "")
	t.Setenv("FLAI_ROLE", "")
	t.Setenv("FLAI_STORY", "")
	wait("outside a story's agent's session", 0, "")
	t.Setenv("FLAI_STORY", "S-0001")

	hook(subAgentHook("SubagentStop", "ac73f87145abc4315", "general-purpose"))
	wait("the sub-agent stopped", 0, "")
}
