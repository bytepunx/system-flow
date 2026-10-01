package guard

import (
	"strings"
	"testing"
)

// g knows flai's commands as the cmd package gives them.
var g = Guard{Commands: []string{"accept", "adr", "archive", "block", "board", "check", "doc", "edit", "guard", "help", "issue", "move", "prime", "release", "show", "stats", "story", "stream", "task", "thread", "touches", "unblock", "version"}}

func bash(agent, cmd string) Event {
	e := Event{ToolName: "Bash", AgentType: agent}
	if agent != "" {
		e.AgentID = "a1b2"
	}
	e.ToolInput.Command = cmd
	return e
}

func TestTheStorysAgentIsNeverRefused(t *testing.T) {
	main := Event{ToolName: "mcp__flai__item_move", AgentType: "planner"} // a session started with --agent has a type, not an ID
	for _, e := range []Event{main, bash("", "flai move S-0001 review && git commit -m x")} {
		if why := g.Check(e); why != "" {
			t.Errorf("%+v refused: %s", e, why)
		}
	}
}

func TestASubAgentReadsFlaiOverMCP(t *testing.T) {
	for _, tool := range MCPReads {
		if why := g.Check(Event{ToolName: MCPPrefix + tool, AgentID: "a1", AgentType: "explorer"}); why != "" {
			t.Errorf("%s refused: %s", tool, why)
		}
	}
	for _, tool := range []string{"item_move", "item_edit", "item_new", "inbox", "wait_for_work", "wait_for_events", "thread_open", "thread_reply", "thread_resolve"} {
		why := g.Check(Event{ToolName: MCPPrefix + tool, AgentID: "a1", AgentType: "verifier"})
		if !strings.Contains(why, "a sub-agent (verifier) cannot call "+tool) || !strings.Contains(why, "final message") {
			t.Errorf("%s: %q", tool, why)
		}
	}
	if why := g.Check(Event{ToolName: MCPPrefix + "item_move", AgentID: "a1"}); !strings.Contains(why, "a sub-agent (unnamed)") {
		t.Errorf("no agent type: %q", why)
	}
	if why := g.Check(Event{ToolName: "Read", AgentID: "a1", AgentType: "explorer"}); why != "" {
		t.Errorf("Read refused: %s", why)
	}
}

func TestASubAgentRunsChecksButNotWrites(t *testing.T) {
	allowed := []string{
		"make test",
		"make flai-test",
		"cd flai && go test -race -short ./... 2>&1 | tail -60",
		"ls flai internal",
		"scripts/flai.sh check --strict",
		"flai board --json",
		"flai --config /tmp/c.json show S-0001",
		"flai doc search 'sub-agents'",
		"flai thread show TH-0042",
		"flai stream diff S-0175",
		"flai stream --help",
		"git -C /w diff main...HEAD --stat",
		"git log --oneline -5 && git status --short",
		"git",
		"grep -rn 'git commit' .",
		"echo 'flai move S-1 review'",
		"FOO=1 go vet ./...",
		"bash -c 'go test ./...'",
	}
	for _, c := range allowed {
		if why := g.Check(bash("verifier", c)); why != "" {
			t.Errorf("%q refused: %s", c, why)
		}
	}
	refused := []string{
		"flai move S-0175 review",
		"cd /w && scripts/flai.sh block S-1 --reason x",
		"flai board limit in-progress 3",
		"flai thread new --on S-1 'q'",
		"flai --config c.json accept S-1",
		"FLAI_AGENT=x env flai stream log S-1 hi",
		"echo $(flai issue bump I-1)",
		"git commit -am wip",
		"git -C /w stash push -u",
		"true; git checkout main",
		"bash -c 'git reset --hard'",
		"sh -c \"flai move S-1 done\"",
		// found by the verifier before review
		"HOME=/tmp flai move S-1 done",
		"GIT_DIR=/w/.git git commit -m x",
		"env -i git commit -m x",
		"sudo -u me git commit -m x",
		"timeout 30 git commit",
		"git ls-files | xargs git rm",
		`find . -exec git checkout {} ;`,
		"bash -lc 'git commit -m x'",
		"sh -ec 'flai move S-1 done'",
		"eval git push",
	}
	for _, c := range refused {
		why := g.Check(bash("verifier", c))
		if !strings.Contains(why, "a sub-agent (verifier) cannot run") || !strings.Contains(why, "ADR-0060") {
			t.Errorf("%q: %q", c, why)
		}
	}
}

func TestCommands(t *testing.T) {
	got := commands(`a "b c" 'd'; e|f && g $(h i) ` + "`j`")
	want := [][]string{{"a", "b c", "d"}, {"e"}, {"f"}, {"g"}, {"h", "i"}, {"j"}}
	if len(got) != len(want) {
		t.Fatalf("got %q", got)
	}
	for i := range want {
		if strings.Join(got[i], "|") != strings.Join(want[i], "|") {
			t.Errorf("%d: got %q want %q", i, got[i], want[i])
		}
	}
}
