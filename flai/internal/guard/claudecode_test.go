package guard

import (
	"bytes"
	"reflect"
	"slices"
	"testing"
)

// claudeCodeReader is the reader flai guard --harness claude-code reads with.
func claudeCodeReader(t *testing.T) Reader {
	t.Helper()
	r, ok := ReaderFor(ClaudeCode)
	if !ok {
		t.Fatalf("no reader for %s among %v", ClaudeCode, Harnesses())
	}
	return r
}

// S-0353: Claude Code's reader maps a PreToolUse's tool onto the Call's
// kind, Bash onto a shell command line, each file-writing tool onto the file
// it writes, and a flai MCP tool onto its bare name and arguments, in the
// guard's role and story; a call with an agent_id is a sub-agent's.
func TestClaudeCodesReaderMapsAToolCallOntoACall(t *testing.T) {
	gd := Guard{Role: RolePlan, Story: "S-0001"}
	in := func(c Call) Call { c.Session, c.Role, c.Story = "s1", RolePlan, "S-0001"; return c }
	for _, tc := range []struct {
		input string
		want  Call
	}{
		{`{"session_id":"s1","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"make test"}}`,
			in(Call{Kind: Shell, Name: "Bash", Command: "make test"})},
		{`{"session_id":"s1","tool_name":"Edit","tool_input":{"file_path":"a.go","old_string":"x"}}`,
			in(Call{Kind: FileWrite, Name: "Edit", Path: "a.go"})},
		{`{"session_id":"s1","tool_name":"MultiEdit","tool_input":{"file_path":"b.go"}}`,
			in(Call{Kind: FileWrite, Name: "MultiEdit", Path: "b.go"})},
		{`{"session_id":"s1","tool_name":"Write","tool_input":{"file_path":".claude/settings.json"}}`,
			in(Call{Kind: FileWrite, Name: "Write", Path: ".claude/settings.json"})},
		{`{"session_id":"s1","tool_name":"NotebookEdit","tool_input":{"notebook_path":"n.ipynb"}}`,
			in(Call{Kind: FileWrite, Name: "NotebookEdit", Path: "n.ipynb"})},
		{`{"session_id":"s1","tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"review"}}`,
			in(Call{Kind: FlaiTool, Name: "mcp__flai__item_move", Tool: "item_move", Args: Args{ID: "S-1", To: "review", Fields: []string{"id", "to"}}})},
		{`{"session_id":"s1","tool_name":"Read","tool_input":{"file_path":"a.go"}}`,
			in(Call{Kind: Other, Name: "Read"})},
		{`{"session_id":"s1","tool_name":"Bash","tool_input":{"command":"flai board"},"agent_id":"a1","agent_type":"verifier"}`,
			in(Call{Subagent: true, Agent: "verifier", Kind: Shell, Name: "Bash", Command: "flai board"})},
	} {
		h, err := claudeCodeReader(t).Read([]byte(tc.input), gd)
		if err != nil {
			t.Fatalf("%s: %v", tc.input, err)
		}
		if h.SubAgent != nil || !reflect.DeepEqual(h.Call, tc.want) {
			t.Errorf("%s:\ngot  %+v (sub-agent change %+v)\nwant %+v", tc.input, h.Call, h.SubAgent, tc.want)
		}
	}
}

// S-0353: Claude Code's reader maps a SubagentStart and a SubagentStop onto
// a change to the record of running sub-agents, with the session and the
// sub-agent's ID and type.
func TestClaudeCodesReaderMapsASubAgentsStartAndStopOntoTheRecord(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  SubAgentChange
	}{
		{`{"session_id":"s1","hook_event_name":"SubagentStart","agent_id":"a1","agent_type":"explorer"}`, SubAgentChange{Session: "s1", ID: "a1", Type: "explorer"}},
		{`{"session_id":"s1","hook_event_name":"SubagentStop","agent_id":"a1","agent_type":"explorer"}`, SubAgentChange{Session: "s1", ID: "a1", Type: "explorer", Stopped: true}},
	} {
		h, err := claudeCodeReader(t).Read([]byte(tc.input), g)
		if err != nil {
			t.Fatalf("%s: %v", tc.input, err)
		}
		if h.SubAgent == nil || *h.SubAgent != tc.want {
			t.Errorf("%s: got %+v, want %+v", tc.input, h.SubAgent, tc.want)
		}
	}
}

// S-0353: input Claude Code's reader cannot read is an error, which flai
// guard warns of and lets through.
func TestClaudeCodesReaderFailsOnInputItCannotRead(t *testing.T) {
	for _, input := range []string{`not json`, `{"tool_name":7}`, `{"tool_name":"Bash","tool_input":"x"}`} {
		if _, err := claudeCodeReader(t).Read([]byte(input), g); err == nil {
			t.Errorf("%s: read", input)
		}
	}
}

// S-0353, S-0354: Claude Code's reader answers an allowed call with exit 0
// and nothing on stderr, a refusal with exit 2 and the reason on stderr, and
// an ask as it does an allowed call, so that Claude Code's own permission
// handling calls permission_prompt, which asks (ADR-0130).
func TestClaudeCodesReaderAnswersTheVerdictWithItsExit(t *testing.T) {
	for _, tc := range []struct {
		v      Verdict
		code   int
		stderr string
	}{
		{Verdict{}, 0, ""},
		{Verdict{Outcome: Deny, Reason: "it reads"}, 2, "it reads\n"},
		{Verdict{Outcome: Ask, Reason: "ask the operator"}, 0, ""},
	} {
		var stderr bytes.Buffer
		if code := claudeCodeReader(t).Answer(tc.v, &stderr); code != tc.code || stderr.String() != tc.stderr {
			t.Errorf("%+v: exit %d, stderr %q; want %d, %q", tc.v, code, stderr.String(), tc.code, tc.stderr)
		}
	}
}

// S-0354: the story's agent's Write of its worktree's .claude/settings.json,
// read from Claude Code's hook input, is asked, and Claude Code's entry
// refuses nothing for it: Claude Code calls permission_prompt as before.
func TestClaudeCodesReaderPassesAStoryAgentsProtectedWriteOnToAsk(t *testing.T) {
	gd := Guard{Commands: g.Commands, Story: "S-0001", Served: true}
	file := "/home/op/proj/.flai-cache/worktrees/S-0001/.claude/settings.json"
	input := `{"session_id":"s1","hook_event_name":"PreToolUse","tool_name":"Write","tool_input":{"file_path":"` + file + `","content":"{}"}}`
	h, err := claudeCodeReader(t).Read([]byte(input), gd)
	if err != nil {
		t.Fatal(err)
	}
	if v := Decide(gd, h.Call); v.Outcome != Ask {
		t.Errorf("decided %+v, want an ask", v)
	}
	if why := gd.Check(fileEdit("Write", file)); why != "" {
		t.Errorf("an ask reported as a refusal: %q", why)
	}
}

// S-0353: claude-code is the one harness with a reader, and a harness with
// none has no reader.
func TestReadersAreByHarness(t *testing.T) {
	if got := Harnesses(); !slices.Equal(got, []string{ClaudeCode}) {
		t.Errorf("harnesses %v", got)
	}
	if _, ok := ReaderFor("codex"); ok {
		t.Error("a reader for codex")
	}
}
