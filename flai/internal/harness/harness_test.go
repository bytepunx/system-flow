package harness

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

func req(a *manifest.Agent) Request {
	return Request{Story: "S-0104", Root: "/p/flow", Project: "flow", Agent: a, Name: "agent-S-0104", Flai: "/usr/local/bin/flai"}
}

// after is the value that follows flag in argv, empty when flag is not there.
func after(argv []string, flag string) string {
	i := slices.Index(argv, flag)
	if i < 0 || i+1 >= len(argv) {
		return ""
	}
	return argv[i+1]
}

func TestClaudeCodeRunsHeadlessWithTheStorysModelAndFlaisMCP(t *testing.T) {
	a := &manifest.Agent{Harness: ClaudeCode, Model: "claude-opus-5-5", Config: map[string]string{"effort": "high", "max_budget_usd": "12.50"}}
	name, ad, err := For(a, false)
	if err != nil || name != ClaudeCode {
		t.Fatalf("For: %q %v", name, err)
	}
	st, err := ad.Start(req(a), Host{})
	if err != nil {
		t.Fatal(err)
	}
	argv := st.Argv
	if argv[0] != "claude" || argv[1] != "-p" || !strings.Contains(argv[2], "work story S-0104") {
		t.Fatalf("program and prompt: %q", argv[:3])
	}
	for flag, want := range map[string]string{"--model": "claude-opus-5-5", "--effort": "high", "--max-budget-usd": "12.50",
		"--output-format": "stream-json", "--permission-mode": "acceptEdits", "--allowedTools": "Bash,mcp__flai", "--name": "flow S-0104"} {
		if got := after(argv, flag); got != want {
			t.Errorf("%s = %q, want %q (%q)", flag, got, want, argv)
		}
	}
	raw := after(argv, "--mcp-config")
	var mcp struct {
		MCPServers map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(raw), &mcp); err != nil || mcp.MCPServers["flai"].Command != "/usr/local/bin/flai" || !slices.Equal(mcp.MCPServers["flai"].Args, []string{"mcp"}) {
		t.Fatalf("mcp config %s: %v", raw, err)
	}

	// the operator's program and arguments replace the defaults; the story's stay
	st, err = ad.Start(req(a), Host{Program: "/opt/claude", Args: []string{"--permission-mode", "bypassPermissions"}})
	if err != nil {
		t.Fatal(err)
	}
	if st.Argv[0] != "/opt/claude" || slices.Contains(st.Argv, "--allowedTools") {
		t.Fatalf("host: %q", st.Argv)
	}
	if got := after(st.Argv, "--permission-mode"); got != "bypassPermissions" {
		t.Fatalf("permission mode %q", got)
	}
	// no model: the harness's own default
	st, _ = ad.Start(req(&manifest.Agent{Harness: ClaudeCode}), Host{})
	if slices.Contains(st.Argv, "--model") {
		t.Fatalf("model with none given: %q", st.Argv)
	}
}

func TestAStoryNamesTunablesNeverFlags(t *testing.T) {
	for _, c := range []struct {
		config map[string]string
		says   string
	}{
		{map[string]string{"permission_mode": "bypassPermissions"}, `takes no config key "permission_mode"; it takes effort, fallback_model, max_budget_usd`},
		{map[string]string{"effort": "--dangerously-skip-permissions"}, "must be low, medium, high, xhigh, or max"},
		{map[string]string{"max_budget_usd": "lots"}, "an amount of dollars"},
		{map[string]string{"fallback_model": "-x"}, "a model's name"},
	} {
		a := &manifest.Agent{Harness: ClaudeCode, Config: c.config}
		if _, err := (claudeCode{}).Start(req(a), Host{}); err == nil || !strings.Contains(err.Error(), c.says) {
			t.Errorf("%v: %v", c.config, err)
		}
	}
}

func TestWhichAdapterAStoryGets(t *testing.T) {
	if _, _, err := For(&manifest.Agent{Harness: "cursor"}, true); err == nil || !strings.Contains(err.Error(), `cannot start harness "cursor"; it can start claude-code, command`) {
		t.Fatalf("unknown: %v", err)
	}
	if name, _, err := For(nil, true); err != nil || name != Command {
		t.Fatalf("no agent, a command set: %q %v", name, err)
	}
	if name, _, err := For(&manifest.Agent{Model: "m"}, true); err != nil || name != Command {
		t.Fatalf("a model, no harness: %q %v", name, err)
	}
	if _, _, err := For(nil, false); err == nil || !strings.Contains(err.Error(), "names no harness") {
		t.Fatalf("nothing to start: %v", err)
	}
	if _, _, err := For(&manifest.Agent{Harness: Command}, false); err == nil || !strings.Contains(err.Error(), "none is set") {
		t.Fatalf("command harness with no command: %v", err)
	}
}

func TestTheOperatorsCommandGetsTheModelAndTheConfig(t *testing.T) {
	a := &manifest.Agent{Model: "claude-sonnet-5", Config: map[string]string{"anything": "goes here"}}
	st, err := (command{}).Start(req(a), Host{Program: "run-agent", Args: []string{"--story={story}", "{root}", "--model", "{model}", "{harness}"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"run-agent", "--story=S-0104", "/p/flow", "--model", "claude-sonnet-5", "command"}; !slices.Equal(st.Argv, want) {
		t.Fatalf("argv %q, want %q", st.Argv, want)
	}
	if want := []string{"FLAI_MODEL=claude-sonnet-5", "FLAI_HARNESS=command", `FLAI_AGENT_CONFIG={"anything":"goes here"}`}; !slices.Equal(st.Env, want) {
		t.Fatalf("env %q", st.Env)
	}
	if _, err := (command{}).Start(req(nil), Host{}); err == nil {
		t.Fatal("no command set: started")
	}
}

func TestThePromptKeepsTheAgentToItsStoryAndTheInbox(t *testing.T) {
	p := Prompt(req(nil))
	for _, want := range []string{"agent-S-0104", "flai stream open S-0104", "no other story", "thread_open", "wait_for_events", "flai move S-0104 review", "flai block S-0104"} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt lacks %q", want)
		}
	}
}
