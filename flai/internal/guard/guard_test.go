package guard

import (
	"strings"
	"testing"
)

func bash(agent, cmd string) Event {
	e := Event{ToolName: "Bash", AgentType: agent}
	e.ToolInput.Command = cmd
	return e
}

func TestTheStorysAgentIsNeverRefused(t *testing.T) {
	for _, e := range []Event{{ToolName: "mcp__flai__item_move"}, bash("", "flai move S-0001 review && git commit -m x")} {
		if why := Check(e); why != "" {
			t.Errorf("%+v refused: %s", e, why)
		}
	}
}

func TestASubAgentReadsFlaiOverMCP(t *testing.T) {
	for _, tool := range MCPReads {
		if why := Check(Event{ToolName: MCPPrefix + tool, AgentType: "explorer"}); why != "" {
			t.Errorf("%s refused: %s", tool, why)
		}
	}
	for _, tool := range []string{"item_move", "item_edit", "item_new", "inbox", "wait_for_work", "wait_for_events", "thread_open", "thread_reply", "thread_resolve"} {
		why := Check(Event{ToolName: MCPPrefix + tool, AgentType: "verifier"})
		if !strings.Contains(why, "a sub-agent (verifier) cannot call "+tool) || !strings.Contains(why, "final message") {
			t.Errorf("%s: %q", tool, why)
		}
	}
	if why := Check(Event{ToolName: "Read", AgentType: "explorer"}); why != "" {
		t.Errorf("Read refused: %s", why)
	}
}

func TestASubAgentRunsChecksButNotWrites(t *testing.T) {
	allowed := []string{
		"make test",
		"cd flai && go test -race -short ./... 2>&1 | tail -60",
		"scripts/flai.sh check --strict",
		"flai board --json",
		"flai --config /tmp/c.json show S-0001",
		"flai doc search 'sub-agents'",
		"flai thread show TH-0042",
		"flai stream diff S-0175",
		"git -C /w diff main...HEAD --stat",
		"git log --oneline -5 && git status --short",
		"git",
		"echo 'flai move S-1 review'",
		"FOO=1 go vet ./...",
	}
	for _, c := range allowed {
		if why := Check(bash("verifier", c)); why != "" {
			t.Errorf("%q refused: %s", c, why)
		}
	}
	refused := map[string]string{
		"flai move S-0175 review":                       "flai move S-0175 review",
		"cd /w && scripts/flai.sh block S-1 --reason x": "scripts/flai.sh block S-1 --reason x",
		"flai board limit in-progress 3":                "flai board limit in-progress 3",
		"flai thread new --on S-1 'q'":                  "flai thread new --on S-1 q",
		"flai --config c.json accept S-1":               "flai --config c.json accept S-1",
		"FLAI_AGENT=x env flai stream log S-1 hi":       "FLAI_AGENT=x env flai stream log S-1 hi",
		"echo $(flai issue bump I-1)":                   "flai issue bump I-1",
		"git commit -am wip":                            "git commit -am wip",
		"git -C /w stash push -u":                       "git -C /w stash push -u",
		"true; git checkout main":                       "git checkout main",
		"bash -c 'git reset --hard'":                    "bash -c git reset --hard",
		"sh -c \"flai move S-1 done\"":                  "sh -c flai move S-1 done",
	}
	for c, words := range refused {
		why := Check(bash("verifier", c))
		if !strings.Contains(why, "a sub-agent (verifier) cannot run \""+words+"\"") || !strings.Contains(why, "ADR-0060") {
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
