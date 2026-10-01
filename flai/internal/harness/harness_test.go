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
	if !slices.Contains(argv, "--strict-mcp-config") {
		t.Errorf("the agent's MCP servers are flai's alone: %q", argv)
	}
	// the server is told the agent's name as an argument, which no environment
	// a harness sets for it can replace (S-0114, I-0037)
	if err := json.Unmarshal([]byte(raw), &mcp); err != nil || mcp.MCPServers["flai"].Command != "/usr/local/bin/flai" || !slices.Equal(mcp.MCPServers["flai"].Args, []string{"mcp", "--agent", "agent-S-0104"}) {
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
	if strings.Contains(p, "--cat") {
		t.Errorf("an agent with a story primes with its pack, not every convention:\n%s", p)
	}
	for _, want := range []string{"agent-S-0104", "prime your session with flai prime --story S-0104 (or the flai MCP tool prime)", "A brief is not the document", "doc_get and its heading", "before relying on it or changing what it describes", "doc_search", "flai stream open S-0104", "no other story", "thread_open", "wait_for_events", "flai move S-0104 review", "commit everything outstanding in the worktree", "refused while anything is uncommitted", "flai block S-0104"} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt lacks %q", want)
		}
	}
}

// S-0116: a restarted agent is told how its last one ended and to go on
// from the story's narrative.
func TestARestartedAgentIsToldHowTheLastOneEnded(t *testing.T) {
	if p := Prompt(req(nil)); !strings.Contains(p, "because it entered ready.") {
		t.Errorf("first start: %s", p)
	}
	r := req(nil)
	r.Restart = "ended (exit 1) with S-0104 in in-progress"
	p := Prompt(r)
	for _, want := range []string{"because the operator restarted it: its last agent ended (exit 1) with S-0104 in in-progress.", "Current state and Next steps", "rather than starting over", "flai stream open S-0104", "flai move S-0104 review"} {
		if !strings.Contains(p, want) {
			t.Errorf("restart prompt lacks %q:\n%s", want, p)
		}
	}
}

// S-0182, I-0050: an agent the operator started is told so, and what the
// start went past, so that it does not take the start for flai serve's own.
func TestAnAgentTheOperatorStartedIsToldWhatItWentPast(t *testing.T) {
	r := req(nil)
	r.Started = true
	if p := Prompt(r); !strings.Contains(p, "because the operator started it now, from the story's page or with flai serve agent start.") || strings.Contains(p, "entered ready") || strings.Contains(p, "past what") {
		t.Errorf("a start past nothing: %s", p)
	}
	hold := "held (overlap): touches flai/cmd, which holds flai/cmd/prime.go that S-0138 (in progress) touches; starts when S-0138 is accepted, cancelled, or sent back"
	r.Past = []string{hold, "the in-progress limit was full"}
	p := Prompt(r)
	for _, want := range []string{
		"because the operator started it now, from the story's page or with flai serve agent start, past what kept flai serve from starting it: " + hold + "; the in-progress limit was full.",
		"pull S-0104 as they asked",
		"do not narrow them only to clear the hold",
		"flai stream open S-0104",
	} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt lacks %q:\n%s", want, p)
		}
	}
	if strings.Contains(p, "entered ready") {
		t.Errorf("an operator's start is not one for entering ready:\n%s", p)
	}
}

// S-0104: an agent that ended waiting for an answer is started again in the
// session it had, told what was answered.
func TestAnAnsweredAgentGoesOnInItsSession(t *testing.T) {
	a := &manifest.Agent{Harness: ClaudeCode, Model: "claude-haiku-4-5"}
	r := req(a)
	r.Session = "5f0c7d0e-1b2a-4c3d-8e9f-0a1b2c3d4e5f"
	st, _ := (claudeCode{}).Start(r, Host{})
	if after(st.Argv, "--session-id") != r.Session || slices.Contains(st.Argv, "--resume") {
		t.Fatalf("first start: %q", st.Argv)
	}
	r.Answered = "TH-0001"
	st, _ = (claudeCode{}).Start(r, Host{})
	if after(st.Argv, "--resume") != r.Session || slices.Contains(st.Argv, "--session-id") {
		t.Fatalf("resumed: %q", st.Argv)
	}
	if p := st.Argv[2]; !strings.Contains(p, "answered your question TH-0001 on S-0104") || !strings.Contains(p, "wait_for_events") {
		t.Errorf("resume prompt: %s", p)
	}
	cmd, _ := (command{}).Start(r, Host{Program: "run-agent"})
	if !slices.Contains(cmd.Env, "FLAI_ANSWERED=TH-0001") {
		t.Errorf("command env: %q", cmd.Env)
	}
}

// S-0140: an agent started to commit a story's worktree is told to do that
// and nothing else; the operator's command is told the worktree.
func TestAnAgentStartedToCommitIsToldToDoOnlyThat(t *testing.T) {
	r := req(nil)
	r.Commit = "/repo/.flai-cache/worktrees/S-0104"
	p := Prompt(r)
	for _, want := range []string{"agent-S-0104", "work left uncommitted in story S-0104's worktree, /repo/.flai-cache/worktrees/S-0104", "Do only this", "on story/S-0104", "Do not move S-0104", "flai stream log S-0104", "thread_open on S-0104"} {
		if !strings.Contains(p, want) {
			t.Errorf("commit prompt lacks %q:\n%s", want, p)
		}
	}
	if strings.Contains(p, "flai move S-0104 review") {
		t.Errorf("a commit run is not told to work the story to review:\n%s", p)
	}
	st, err := command{}.Start(r, Host{Program: "run"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(st.Env, "\n"), "FLAI_COMMIT=/repo/.flai-cache/worktrees/S-0104") {
		t.Errorf("env: %v", st.Env)
	}
}

// S-0175, ADR-0059: the claude-code prompt says when to delegate, what to
// give a sub-agent, and to verify before review; the operator's command gets
// no prompt, and an answered or commit run is not told again.
func TestThePromptHandsNoisyWorkToSubAgents(t *testing.T) {
	want := []string{"hand noisy work to sub-agents with the Agent tool", "to the explorer", "to the verifier", "the worktree's path, S-0104 and the task's ID, the question", "a summary with paths and lines, not raw output", "a question it returns for the designer is yours to ask with thread_open", "Before you move S-0104 to review, have a fresh verifier check the worktree's diff against the acceptance criteria and the conventions"}
	r := req(&manifest.Agent{Harness: ClaudeCode})
	st, err := (claudeCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	restarted := r
	restarted.Restart = "ended (exit 1)"
	for _, p := range []string{st.Argv[2], Prompt(restarted)} {
		for _, w := range want {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
	}
	answered, commit := r, r
	answered.Answered, commit.Commit = "TH-0001", "/w"
	for _, p := range []string{Prompt(answered), Prompt(commit)} {
		if strings.Contains(p, "sub-agents") {
			t.Errorf("told again:\n%s", p)
		}
	}
	cmd, err := (command{}).Start(r, Host{Program: "run-agent", Args: []string{"{story}"}})
	if err != nil {
		t.Fatal(err)
	}
	if all := strings.Join(append(cmd.Argv, cmd.Env...), "\n"); strings.Contains(all, "sub-agent") || strings.Contains(all, "explorer") {
		t.Errorf("the operator's command is told about sub-agents: %s", all)
	}
}
