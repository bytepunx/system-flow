package cmd

import (
	"strings"
	"testing"
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
		{`{"tool_name":"Bash","tool_input":{"command":"flai move S-1 review"},"agent_type":"Explore"}`, 2, "a sub-agent (Explore) cannot run \"flai move S-1 review\""},
		{`{"tool_name":"Bash","tool_input":{"command":"make test"},"agent_type":"verifier"}`, 0, ""},
		{`{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"review"}}`, 0, ""},
		{`not json`, 0, "hook input is not a tool call"},
	} {
		_, errOut, code := runStdin(t, dir, c.in, "guard")
		if code != c.code || (c.err == "" && errOut != "") || !strings.Contains(errOut, c.err) {
			t.Errorf("%s: code %d, stderr %q", c.in, code, errOut)
		}
	}
}
