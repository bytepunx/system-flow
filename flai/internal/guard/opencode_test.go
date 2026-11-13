package guard

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// openCodeReader is the reader flai guard --harness opencode reads with.
func openCodeReader(t *testing.T) Reader {
	t.Helper()
	r, ok := ReaderFor(OpenCode)
	if !ok {
		t.Fatalf("no reader for %s among %v", OpenCode, Harnesses())
	}
	return r
}

// S-0407: OpenCode's reader maps a tool.execute.before's tool onto the
// Call's kind, bash onto a shell command line, edit, write, and multiedit
// onto the file they write, and flai_<tool> onto flai's tool and its
// arguments, in the guard's role and story; a call from a session with a
// parent is a sub-agent's, named by its agent.
func TestOpenCodesReaderMapsAToolCallOntoACall(t *testing.T) {
	gd := Guard{Role: RolePlan, Story: "S-0001"}
	in := func(c Call) Call { c.Session, c.Role, c.Story = "ses_1", RolePlan, "S-0001"; return c }
	const root = `"session":{"id":"ses_1","agent":"build"}`
	const child = `"session":{"id":"ses_1","parentID":"ses_0","agent":"explorer"}`
	for _, tc := range []struct {
		input string
		want  Call
	}{
		{`{"event":"tool.execute.before","tool":"bash","sessionID":"ses_1","callID":"c1","args":{"command":"make test"},` + root + `}`,
			in(Call{Kind: Shell, Name: "bash", Command: "make test"})},
		{`{"tool":"edit","sessionID":"ses_1","args":{"filePath":"/w/a.go","oldString":"x","newString":"y"},` + root + `}`,
			in(Call{Kind: FileWrite, Name: "edit", Path: "/w/a.go"})},
		{`{"tool":"write","sessionID":"ses_1","args":{"filePath":".claude/settings.json","content":"{}"},` + root + `}`,
			in(Call{Kind: FileWrite, Name: "write", Path: ".claude/settings.json", Whole: true, Content: "{}"})},
		{`{"tool":"multiedit","sessionID":"ses_1","args":{"filePath":"b.go"},` + root + `}`,
			in(Call{Kind: FileWrite, Name: "multiedit", Path: "b.go"})},
		{`{"tool":"read","sessionID":"ses_1","args":{"filePath":"a.go"},` + root + `}`,
			in(Call{Kind: Other, Name: "read"})},
		{`{"tool":"task","sessionID":"ses_1","args":{"subagent_type":"explorer","prompt":"find"},` + root + `}`,
			in(Call{Kind: Other, Name: "task"})},
		{`{"tool":"flai_item_move","sessionID":"ses_1","args":{"id":"S-1","to":"review"},` + child + `}`,
			in(Call{Kind: FlaiTool, Name: "flai_item_move", Tool: "item_move", Subagent: true, Agent: "explorer", Args: Args{ID: "S-1", To: "review", Fields: []string{"id", "to"}}})},
		{`{"tool":"flai_inbox","sessionID":"ses_1",` + root + `}`,
			in(Call{Kind: FlaiTool, Name: "flai_inbox", Tool: "inbox"})},
		{`{"tool":"bash","sessionID":"ses_1","args":{"command":"ls"},` + child + `}`,
			in(Call{Kind: Shell, Name: "bash", Command: "ls", Subagent: true, Agent: "explorer"})},
	} {
		h, err := openCodeReader(t).Read([]byte(tc.input), gd)
		if err != nil {
			t.Fatalf("%s: %v", tc.input, err)
		}
		if h.SubAgent != nil || len(h.Also) != 0 || !reflect.DeepEqual(h.Call, tc.want) {
			t.Errorf("%s:\n got %+v\nwant %+v", tc.input, h, tc.want)
		}
	}
	if _, err := openCodeReader(t).Read([]byte(`not json`), gd); err == nil {
		t.Error("input that is not JSON was read")
	}
}

// S-0407: an apply_patch writes each file its patch names, once each, the
// first as the Call and the rest in Also.
func TestOpenCodesReaderReadsEachFileAPatchWrites(t *testing.T) {
	patch := "*** Begin Patch\n*** Add File: docs/a.md\n+x\n*** Update File: flai/b.go\n@@\n-y\n+z\n*** Move to: flai/c.go\n*** Delete File: .claude/settings.json\n*** Update File: docs/a.md\n*** End Patch\n"
	if got, want := PatchPaths(patch), []string{"docs/a.md", "flai/b.go", "flai/c.go", ".claude/settings.json"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("PatchPaths = %v, want %v", got, want)
	}
	input := `{"tool":"apply_patch","sessionID":"ses_1","args":{"patchText":` + quote(patch) + `},"session":{"id":"ses_1"}}`
	h, err := openCodeReader(t).Read([]byte(input), Guard{})
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, c := range h.Calls() {
		if c.Kind != FileWrite || c.Name != "apply_patch" {
			t.Errorf("call %+v is not apply_patch's file write", c)
		}
		paths = append(paths, c.Path)
	}
	if want := []string{"docs/a.md", "flai/b.go", "flai/c.go", ".claude/settings.json"}; !reflect.DeepEqual(paths, want) {
		t.Errorf("paths %v, want %v", paths, want)
	}
	// a patch naming no file is one write with no path
	if h, _ := openCodeReader(t).Read([]byte(`{"tool":"apply_patch","sessionID":"ses_1","args":{"patchText":"nothing"}}`), Guard{}); h.Call.Kind != FileWrite || h.Call.Path != "" || len(h.Also) != 0 {
		t.Errorf("empty patch: %+v", h)
	}
}

// quote is s as a JSON string.
func quote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// S-0407: a session.created whose session has a parent is a sub-agent
// starting in its parent's session, and a task's tool.execute.after the
// sub-agent of its child session stopping, as SubagentStart and SubagentStop
// are; a root session's creation and any other tool's end ask nothing.
func TestOpenCodesReaderRecordsASubAgentsSession(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  *SubAgentChange
	}{
		{`{"event":"session.created","session":{"id":"ses_2","parentID":"ses_1","agent":"explorer"}}`,
			&SubAgentChange{Session: "ses_1", ID: "ses_2", Type: "explorer"}},
		{`{"event":"tool.execute.after","tool":"task","sessionID":"ses_1","session":{"id":"ses_2","parentID":"ses_1","agent":"explorer"}}`,
			&SubAgentChange{Session: "ses_1", ID: "ses_2", Type: "explorer", Stopped: true}},
		{`{"event":"session.created","session":{"id":"ses_1","agent":"build"}}`, nil},
		{`{"event":"tool.execute.after","tool":"bash","sessionID":"ses_1","session":{"id":"ses_1"}}`, nil},
	} {
		h, err := openCodeReader(t).Read([]byte(tc.input), Guard{Story: "S-0001"})
		if err != nil {
			t.Fatalf("%s: %v", tc.input, err)
		}
		if !reflect.DeepEqual(h.SubAgent, tc.want) {
			t.Errorf("%s: sub-agent %+v, want %+v", tc.input, h.SubAgent, tc.want)
		}
		if tc.want == nil {
			if h.Call.Kind != Other || h.Call.Subagent || Decide(Guard{Story: "S-0001"}, h.Call).Outcome != Allow {
				t.Errorf("%s: asks %+v, want a call that passes", tc.input, h.Call)
			}
		}
	}
}

// S-0407: the guard decides OpenCode's calls as it does Claude Code's: a
// sub-agent's flai write is refused and the story's own agent's let through,
// a sub-agent's write to a protected path is refused while auto-approve is
// off, and a patch is refused when any file it writes is.
func TestOpenCodesCallsAreDecidedAsClaudeCodesAre(t *testing.T) {
	gd := Guard{Story: "S-0001"}
	decide := func(input string) Verdict {
		t.Helper()
		h, err := openCodeReader(t).Read([]byte(input), gd)
		if err != nil {
			t.Fatal(err)
		}
		return DecideAll(gd, h.Calls())
	}
	const child = `"session":{"id":"ses_2","parentID":"ses_1","agent":"explorer"}`
	const root = `"session":{"id":"ses_1","agent":"build"}`
	move := `{"tool":"flai_item_move","sessionID":"ses_2","args":{"id":"S-1","to":"review"},` + child + `}`
	if v := decide(move); v.Outcome != Deny || !strings.Contains(v.Reason, "a sub-agent (explorer) cannot call item_move") {
		t.Errorf("sub-agent's item_move: %+v", v)
	}
	ccMove := `{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"review"},"agent_type":"explorer","agent_id":"a1"}`
	if got, want := decide(move).Reason, gd.Check(mustEvent(t, ccMove)); got != want {
		t.Errorf("opencode refuses %q, claude-code %q", got, want)
	}
	if v := decide(`{"tool":"flai_item_move","sessionID":"ses_1","args":{"id":"S-1","to":"review"},` + root + `}`); v.Outcome != Allow {
		t.Errorf("story's agent's item_move: %+v", v)
	}
	if v := decide(`{"tool":"write","sessionID":"ses_2","args":{"filePath":"/w/.claude/settings.json"},` + child + `}`); v.Outcome != Deny || !strings.Contains(v.Reason, "a sub-agent (explorer) cannot use write on /w/.claude/settings.json") {
		t.Errorf("sub-agent's protected write: %+v", v)
	}
	if v := decide(`{"tool":"write","sessionID":"ses_2","args":{"filePath":"/w/flai/a.go"},` + child + `}`); v.Outcome != Allow {
		t.Errorf("sub-agent's write: %+v", v)
	}
	patch := "*** Begin Patch\n*** Update File: flai/a.go\n*** Update File: .claude/settings.json\n*** End Patch\n"
	if v := decide(`{"tool":"apply_patch","sessionID":"ses_2","args":{"patchText":` + quote(patch) + `},` + child + `}`); v.Outcome != Deny || !strings.Contains(v.Reason, "on .claude/settings.json") {
		t.Errorf("sub-agent's patch with a protected file: %+v", v)
	}
}

// mustEvent is the Claude Code hook input in.
func mustEvent(t *testing.T, in string) Event {
	t.Helper()
	var e Event
	if err := json.Unmarshal([]byte(in), &e); err != nil {
		t.Fatal(err)
	}
	return e
}

// S-0407: DecideAll answers a refusal over an Ask, an Ask over a pass, and
// passes when every call does.
func TestDecideAllAnswersTheFirstRefusalThenTheFirstAsk(t *testing.T) {
	gd := Guard{Story: "S-0001"}
	wt := "/p/.flai-cache/worktrees/S-0001/"
	write := func(p string) Call { return Call{Kind: FileWrite, Name: "write", Story: "S-0001", Path: wt + p} }
	sub := func(c Call) Call { c.Subagent, c.Agent = true, "explorer"; return c }
	if v := DecideAll(gd, []Call{write("a.go"), write(".claude/x.md")}); v.Outcome != Ask {
		t.Errorf("story's agent's protected write: %+v", v)
	}
	if v := DecideAll(gd, []Call{write(".claude/x.md"), sub(write(".mcp.json"))}); v.Outcome != Deny {
		t.Errorf("a refusal after an ask: %+v", v)
	}
	if v := DecideAll(gd, []Call{write("a.go"), write("b.go")}); v.Outcome != Allow {
		t.Errorf("plain writes: %+v", v)
	}
}

// S-0407: OpenCode's answer refuses with exit 2 and the reason on standard
// error, and passes an Ask and an Allow with exit 0 and nothing written.
func TestOpenCodesAnswer(t *testing.T) {
	for _, tc := range []struct {
		v    Verdict
		code int
		err  string
	}{
		{Verdict{Outcome: Deny, Reason: "no"}, OpenCodeRefuses, "no\n"},
		{Verdict{Outcome: Ask, Reason: "ask"}, 0, ""},
		{Verdict{}, 0, ""},
	} {
		var stderr bytes.Buffer
		if code := openCodeReader(t).Answer(tc.v, &stderr); code != tc.code || stderr.String() != tc.err {
			t.Errorf("%+v: code %d, stderr %q", tc.v, code, stderr.String())
		}
	}
}
