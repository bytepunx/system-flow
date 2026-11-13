package guard

import (
	"encoding/json"
	"strings"
	"testing"
)

// the main checkout's task file I-0352's Write replaced
const i0352Task = "/home/op/project/wip/kanban/tasks/T-3079-the-guard-reads-a-write.md"

// task3079 is content that is T-3079.
const task3079 = "---\nid: T-3079\ntype: task\ntitle: The guard reads a write\nstatus: backlog\nparent: S-0665\n---\n\n# T-3079 The guard reads a write\n\n## Work\n\nRewritten whole.\n"

// claudeWrite is a Claude Code Write hook input of content to path, by a
// sub-agent when agent is set.
func claudeWrite(t *testing.T, path, content, agent string) []byte {
	t.Helper()
	in := map[string]any{"session_id": "s1", "hook_event_name": "PreToolUse", "tool_name": "Write",
		"tool_input": map[string]string{"file_path": path, "content": content}}
	if agent != "" {
		in["agent_id"], in["agent_type"] = "a1", agent
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// I-0352: S-0665's agent replaced the main checkout's T-3079 file with a
// Write of PLACEHOLDER, and every Repo.List failed until it wrote the file
// back. The guard now refuses that Write, and any other that leaves a work
// item file without that item's front matter, and points to Edit and flai
// edit (ADR-0298).
func TestTheGuardRefusesAWriteThatReplacesAWorkItemWithWhatIsNotIt(t *testing.T) {
	g := Guard{Story: "S-0665", Served: true}
	h, err := claudeCodeReader(t).Read(claudeWrite(t, i0352Task, "PLACEHOLDER", ""), g)
	if err != nil {
		t.Fatal(err)
	}
	v := Decide(g, h.Call)
	if v.Outcome != Deny {
		t.Fatalf("I-0352's Write of PLACEHOLDER passed: %+v", v)
	}
	for _, want := range []string{"Write", i0352Task, "T-3079", "front matter", "Edit", "flai edit", "ADR-0298"} {
		if !strings.Contains(v.Reason, want) {
			t.Errorf("the refusal does not name %q: %s", want, v.Reason)
		}
	}

	for _, tc := range []struct {
		name, path, content, agent string
	}{
		{"no front matter, from a sub-agent", i0352Task, "PLACEHOLDER", "general-purpose"},
		{"empty", i0352Task, "", ""},
		{"front matter that does not parse", i0352Task, "---\nid: [T-3079\n---\n", ""},
		{"another item's id", i0352Task, strings.Replace(task3079, "id: T-3079", "id: T-3080", 1), ""},
		{"another type", i0352Task, strings.Replace(task3079, "type: task", "type: story", 1), ""},
		{"a story's worktree", "/home/op/project/.flai-cache/worktrees/S-0665/wip/kanban/stories/S-0665-a-story.md", "# S-0665\n", ""},
		{"the archive", "/home/op/project/wip/archive/kanban/epics/E-0001-an-epic.md", "PLACEHOLDER", ""},
		{"a relative path", "wip/kanban/tasks/T-0001-x.md", "PLACEHOLDER", ""},
	} {
		h, err := claudeCodeReader(t).Read(claudeWrite(t, tc.path, tc.content, tc.agent), g)
		if err != nil {
			t.Fatal(err)
		}
		if v := Decide(g, h.Call); v.Outcome != Deny {
			t.Errorf("%s: the Write passed: %+v", tc.name, v)
		}
	}
}

// ADR-0298: a Write that keeps the item's own id and type, an edit of part
// of an item file, and a Write to a file that is no work item pass.
func TestTheGuardLetsAWriteThatKeepsTheItemAndEveryEditPass(t *testing.T) {
	g := Guard{Story: "S-0665", Served: true}
	for _, tc := range []struct {
		name  string
		input string
	}{
		{"the item rewritten whole", string(claudeWrite(t, i0352Task, task3079, ""))},
		{"an id in any zero padding", string(claudeWrite(t, i0352Task, strings.Replace(task3079, "id: T-3079", "id: T-03079", 1), ""))},
		{"an Edit of the item", `{"session_id":"s1","tool_name":"Edit","tool_input":{"file_path":"` + i0352Task + `","old_string":"a","new_string":"PLACEHOLDER"}}`},
		{"the board", string(claudeWrite(t, "/home/op/project/wip/kanban/board.md", "PLACEHOLDER", ""))},
		{"a narrative", string(claudeWrite(t, "/home/op/project/wip/agents/S-0665.md", "PLACEHOLDER", ""))},
		{"a design document", string(claudeWrite(t, "/home/op/project/design/system/tasks/T-1-x.md", "PLACEHOLDER", ""))},
		{"a task's name in the stories folder", string(claudeWrite(t, "/home/op/project/wip/kanban/stories/T-0001-x.md", "PLACEHOLDER", ""))},
		{"an issue", string(claudeWrite(t, "/home/op/project/design/issues/I-0352-x.md", "PLACEHOLDER", ""))},
	} {
		h, err := claudeCodeReader(t).Read([]byte(tc.input), g)
		if err != nil {
			t.Fatal(err)
		}
		if v := Decide(g, h.Call); v.Outcome == Deny {
			t.Errorf("%s: refused: %s", tc.name, v.Reason)
		}
	}
}

// ADR-0298: OpenCode's write replaces its file whole as Claude Code's Write
// does, and is refused the same way; its edit is not.
func TestTheGuardRefusesOpenCodesWriteOfAWorkItemAsClaudeCodes(t *testing.T) {
	g := Guard{Story: "S-0665", Served: true}
	for _, tc := range []struct {
		input string
		deny  bool
	}{
		{`{"tool":"write","sessionID":"ses_1","args":{"filePath":"` + i0352Task + `","content":"PLACEHOLDER"},"session":{"id":"ses_1"}}`, true},
		{`{"tool":"edit","sessionID":"ses_1","args":{"filePath":"` + i0352Task + `","oldString":"a","newString":"b"},"session":{"id":"ses_1"}}`, false},
	} {
		h, err := openCodeReader(t).Read([]byte(tc.input), g)
		if err != nil {
			t.Fatal(err)
		}
		if got := Decide(g, h.Call).Outcome == Deny; got != tc.deny {
			t.Errorf("%s: refused %v, want %v", tc.input, got, tc.deny)
		}
	}
}
