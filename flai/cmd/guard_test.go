package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

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
		{`{"tool_name":"mcp__flai__thread_reply","tool_input":{"thread":"TH-1","text":"I recommend S-2"}}`, 0, ""},
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
	_, errOut, code := runStdin(t, root, `{"tool_name":"Bash","tool_input":{"command":"flai push"}}`, "guard")
	if code != 2 || !strings.Contains(errOut, "it needs orchestration.permissions.publish") || !strings.Contains(errOut, "refusal not logged") {
		t.Errorf("unlogged refusal: code %d, stderr %q", code, errOut)
	}
	// outside a project every permission is off
	_, errOut, code = runStdin(t, t.TempDir(), `{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"ready"}}`, "guard")
	if code != 2 || !strings.Contains(errOut, "it needs orchestration.permissions.promote_to_ready") || !strings.Contains(errOut, "permissions are taken as off") {
		t.Errorf("no project: code %d, stderr %q", code, errOut)
	}
}

// S-0208, TH-0096: this repository's settings and the template's run the
// guard before Bash and flai's MCP tools in every session, and before Edit,
// Write, and NotebookEdit in a planner session alone.
func TestTheSettingsRunTheGuardOnThePlannersEdits(t *testing.T) {
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
		if c := commands["Edit|Write|NotebookEdit"]; !strings.HasPrefix(c, `[ "$FLAI_ROLE" = plan ] || exit 0; `) || !strings.Contains(c, " guard 2>&1") {
			t.Errorf("%s: a planner's file edits are not guarded, or a story's are: %q", path, c)
		}
		if len(commands) != 2 {
			t.Errorf("%s: hooks %v", path, commands)
		}
	}
}
