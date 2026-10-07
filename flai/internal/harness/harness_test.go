package harness

import (
	"encoding/json"
	"os"
	"path/filepath"
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

// Headless, claude-code refuses a write under .claude/ unless a permission
// handler approves it; flai's permission_prompt is that handler, with the
// default host arguments and with the operator's own, ahead of them (S-0257).
func TestClaudeCodeAsksFlaisPermissionPrompt(t *testing.T) {
	a := &manifest.Agent{Harness: ClaudeCode}
	for _, host := range []Host{{}, {Program: "claude", Args: []string{"--permission-mode", "acceptEdits"}}} {
		st, err := (claudeCode{}).Start(req(a), host)
		if err != nil {
			t.Fatal(err)
		}
		if got := after(st.Argv, "--permission-prompt-tool"); got != "mcp__flai__permission_prompt" {
			t.Errorf("host %v: --permission-prompt-tool = %q (%q)", host, got, st.Argv)
		}
		if i, j := slices.Index(st.Argv, "--permission-prompt-tool"), slices.Index(st.Argv, "--permission-mode"); i < 0 || j < i {
			t.Errorf("host %v: the handler comes before the operator's arguments: %q", host, st.Argv)
		}
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

// S-0177, ADR-0064: an agent started for a story begun on another host is
// told where, when, and by whom, to reconcile rather than start over, and
// which threads to read.
func TestAnAgentForAStoryBegunElsewhereIsToldWhereItWasBegun(t *testing.T) {
	r := req(nil)
	r.Begun = &Begun{By: "agent-S-0104", At: "2026-10-01T07:40:00Z", Agent: "agent-S-0104", Host: "far-away", Threads: []string{"TH-0041"}}
	p := Prompt(r)
	for _, want := range []string{
		"because the operator started it here, and this host has had no agent for it: it was begun by agent-S-0104, who moved it to in-progress at 2026-10-01T07:40:00Z, on the host far-away.",
		"run flai stream open S-0104",
		"from this clone, else from the remote, else new from the main branch",
		"go on from what is committed rather than starting over",
		"TH-0041, on the story, has an entry by someone else since it was begun: read it with the flai MCP tool thread_get",
		"flai move S-0104 review",
	} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt lacks %q:\n%s", want, p)
		}
	}
	if strings.Contains(p, "entered ready") {
		t.Errorf("a story begun elsewhere did not just enter ready:\n%s", p)
	}
	r.Begun = &Begun{By: "alex", At: "2026-10-01T07:40:00Z", Agent: "agent-S-0104", Threads: []string{"TH-0041", "TH-0042"}}
	p = Prompt(r)
	for _, want := range []string{"begun by alex, who moved it to in-progress at 2026-10-01T07:40:00Z; its narrative names agent-S-0104, on another host.", "Its threads TH-0041, TH-0042 have entries"} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt lacks %q:\n%s", want, p)
		}
	}
	if p := Prompt(req(nil)); strings.Contains(p, "begun") {
		t.Errorf("a story that entered ready is not told it was begun elsewhere:\n%s", p)
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

// S-0272: a story's agent, started, answered, or committing, asks on a
// thread, writes the narrative's state, and ends rather than holding
// wait_for_events for the answer, and ends when wait_for_events answers end;
// flai serve starts it again on the answer. The planner, the orchestrator, and
// the analyzer, whom flai serve does not start again, still hold.
func TestAStorysAgentEndsOnAnOpenQuestion(t *testing.T) {
	ends := []string{
		"write the narrative's Current state and Next steps, saying what you asked",
		"and end: flai serve starts you again in this session when the thread is answered, and your first inbox holds the answer",
		"Do not hold the flai MCP tool wait_for_events for an answer; when it answers end: true",
	}
	answered := req(nil)
	answered.Answered = "TH-0001"
	commit := req(nil)
	commit.Commit = "/repo/.flai-cache/worktrees/S-0104"
	for name, r := range map[string]Request{"story": req(nil), "answered": answered, "commit": commit} {
		p := Prompt(r)
		for _, want := range ends {
			if !strings.Contains(p, want) {
				t.Errorf("the %s run's prompt lacks %q:\n%s", name, want, p)
			}
		}
		if strings.Contains(p, "until the thread has an answer") || strings.Contains(p, "If you end while the question is open") {
			t.Errorf("the %s run's prompt still holds wait_for_events for an answer:\n%s", name, p)
		}
	}
	p := Prompt(req(nil))
	for _, want := range []string{
		"ask with the flai MCP tool thread_open on S-0104, and go on with the work of S-0104 that does not wait on the answer. When nothing is left but the answer, write the narrative's Current state and Next steps, saying what you asked and what you will do with each answer, and end",
		"when it answers end: true, do as its why says: write the narrative's Current state and Next steps, and end.",
		// S-0285: still never ends nor waits on wait_for_events while a sub-agent runs
		"Never end your turn while a sub-agent runs in the background",
		"Never wait for a sub-agent with the flai MCP tool wait_for_events either",
	} {
		if !strings.Contains(p, want) {
			t.Errorf("the story run's prompt lacks %q:\n%s", want, p)
		}
	}
	if p := Prompt(commit); !strings.Contains(p, "ask with the flai MCP tool thread_open on S-0104, and commit what does not wait on the answer. Then write the narrative's Current state and Next steps, saying what you asked and what is left to commit, and end") {
		t.Errorf("the commit run's prompt:\n%s", p)
	}
	for name, c := range map[string]struct {
		r    Request
		hold string
	}{
		"planner":      {planReq("E-0016", nil), "and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on"},
		"analyzer":     {analyzeReq("risk", nil), "and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on"},
		"orchestrator": {orchestrateReq(nil), "Then hold the flai MCP tool wait_for_events, again each time it returns"},
	} {
		p := Prompt(c.r)
		if !strings.Contains(p, c.hold) {
			t.Errorf("the %s's prompt no longer holds wait_for_events:\n%s", name, p)
		}
		for _, w := range ends {
			if strings.Contains(p, w) {
				t.Errorf("the %s's prompt says %q, a story agent's rule:\n%s", name, w, p)
			}
		}
	}
}

// S-0175, ADR-0059: the claude-code prompt says when to delegate, what to
// give a sub-agent, and to verify before review; the operator's command gets
// no prompt, and an answered or commit run is not told again.
func TestThePromptHandsNoisyWorkToSubAgents(t *testing.T) {
	want := []string{"hand noisy work to sub-agents with the Agent tool", "to the explorer", "to the verifier", "the worktree's path, S-0104 and the task's ID, the question", "a summary with paths and lines, not raw output", "a question it returns for the designer is yours to ask with thread_open", "Before you move S-0104 to review, commit everything, run flai stream sync S-0104 again and resolve what it lists, then have one fresh verifier", "check the diff against the acceptance criteria and the conventions"}
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

// S-0195, ADR-0067: publishing is the operator's, so no prompt, however the
// agent was started, tells it to push accepted work.
func TestThePromptNeverTellsTheAgentToPush(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	restarted, answered, commit, begun, started := r, r, r, r, r
	restarted.Restart = "ended (exit 1)"
	answered.Answered = "TH-0001"
	commit.Commit = "/w"
	begun.Begun = &Begun{By: "alex", At: "2026-10-01T07:40:00Z", Agent: "agent-S-0104", Threads: []string{"TH-0041"}}
	started.Started, started.Past = true, []string{"the in-progress limit was full"}
	for _, x := range []Request{r, restarted, answered, commit, begun, started} {
		p := Prompt(x)
		for _, never := range []string{"flai push", "git push", "unpushed", "push --pending", "publish"} {
			if strings.Contains(p, never) {
				t.Errorf("prompt says %q:\n%s", never, p)
			}
		}
	}
}

// S-0198: an agent working its story, fresh or answered, records issues for
// it and leaves to the operator, at acceptance, which become stories; one
// started only to commit is not told.
func TestThePromptSaysTheOperatorTurnsIssuesIntoStories(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	answered, commit := r, r
	answered.Answered = "TH-0001"
	commit.Commit = "/w"
	for _, x := range []Request{r, answered} {
		p := Prompt(x)
		for _, w := range []string{
			"flai issue new, or flai issue bump when the issue exists: each instance names " + r.Story,
			"Make no story for them yourself: " + r.Story + "'s review page lists them, checked, so the operator chooses at acceptance which become backlog stories",
			"flai issue story or the flai MCP tool issue_story only when the operator asks, or when flai check warns issues.no-story",
		} {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
	}
	if strings.Contains(Prompt(commit), "flai issue") {
		t.Errorf("the commit prompt mentions issues:\n%s", Prompt(commit))
	}
}

// S-0189: the story's agent runs only the tests for what it changed, leaves
// the whole suite, lint, and check to one verifier before review and one
// more after its fixes, and fixes what a verifier finds itself.
func TestThePromptLeavesTheWholeSuiteToTheVerifier(t *testing.T) {
	p := Prompt(req(&manifest.Agent{Harness: ClaudeCode}))
	for _, w := range []string{
		"run only the tests for what you changed",
		"Leave the whole suite, the whole lint, and flai check to the verifier rather than running them yourself as well",
		"have one fresh verifier run the whole suite, the lint, and flai check in the worktree, through the project's close-out script where it has one",
		"Fix what it finds yourself, never through a sub-agent, commit, and have one more fresh verifier run the same and confirm the fixes",
		"A verifier's passing run is the story's run before review: do not repeat it",
	} {
		if !strings.Contains(p, w) {
			t.Errorf("prompt lacks %q:\n%s", w, p)
		}
	}
}

// S-0273: the story's agent runs the tests for what it changed itself, with
// flai test or the MCP tool test on the paths it changed, not the test and
// lint tools by hand and not through a sub-agent; the whole suite and lint
// stay the verifier's.
func TestThePromptRunsTheTestsForWhatChangedWithFlaiTest(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	restarted := r
	restarted.Restart = "ended (exit 1)"
	for _, p := range []string{Prompt(r), Prompt(restarted)} {
		for _, w := range []string{
			"run only the tests for what you changed, yourself, with flai test and the paths you changed, or the flai MCP tool test with them",
			"it runs the test and lint tiers the manifest's tests declare that those paths select, cheapest first, and answers pass or the first findings",
			"Do not run go test, vitest, golangci-lint, or gofmt by hand and read their logs, and do not hand that run to a sub-agent",
			"runs of the whole suite, the whole lint, and flai check, and long logs, to the verifier",
			"Leave the whole suite, the whole lint, and flai check to the verifier rather than running them yourself as well",
		} {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
		if strings.Contains(p, "test, lint, and flai check runs and long logs to the verifier") {
			t.Errorf("the prompt still hands every test run to the verifier:\n%s", p)
		}
	}
	answered, commit := r, r
	answered.Answered, commit.Commit = "TH-0001", "/w"
	for _, p := range []string{Prompt(answered), Prompt(commit)} {
		if strings.Contains(p, "flai test") {
			t.Errorf("told about flai test again:\n%s", p)
		}
	}
}

// S-0266: the story's agent has the verifier run the close-out once, in one
// command without a pipe or a file, and read its last line, and names in the
// verifier's prompt any step it already knows will stop; an answered or
// commit run is not told again.
func TestThePromptSaysHowTheVerifierRunsTheCloseOut(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	restarted := r
	restarted.Restart = "ended (exit 1)"
	want := []string{
		"through the project's close-out script where it has one",
		"Tell it to run the close-out once, in one command and without a pipe or a file",
		"read its last line, which names the outcome and the step it stopped at",
		"name in its prompt any step you already know will stop, and why",
		"so that it reports that stop and checks the steps after it rather than running the close-out again",
		"Fix what it finds yourself",
	}
	for _, p := range []string{Prompt(r), Prompt(restarted)} {
		at := 0
		for _, w := range want {
			i := strings.Index(p[at:], w)
			if i < 0 {
				t.Errorf("prompt lacks %q after what comes before it:\n%s", w, p)
				continue
			}
			at += i + len(w)
		}
	}
	answered, commit := r, r
	answered.Answered, commit.Commit = "TH-0001", "/w"
	for _, p := range []string{Prompt(answered), Prompt(commit)} {
		if strings.Contains(p, "run the close-out once") {
			t.Errorf("told how the verifier runs the close-out again:\n%s", p)
		}
	}
}

// S-0268: the story's agent makes independent edits and commands in one turn,
// and moves a task it has just written through ready to in-progress in one
// command, since flai move refuses it straight from backlog; an answered or
// commit run is not told again.
func TestThePromptSaysToBatchIndependentCalls(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	restarted := r
	restarted.Restart = "ended (exit 1)"
	want := []string{
		"Make independent edits and commands in one turn, as several tool calls in one message",
		"consecutive edits to one file, reads of files you already know, and commands that do not wait on each other",
		"Move a task you have just written to ready and in-progress in one command",
		"flai move T-nnnn ready && flai move T-nnnn in-progress",
		"flai move refuses a task straight from backlog to in-progress",
	}
	for _, p := range []string{Prompt(r), Prompt(restarted)} {
		for _, w := range want {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
	}
	answered, commit := r, r
	answered.Answered, commit.Commit = "TH-0001", "/w"
	for _, p := range []string{Prompt(answered), Prompt(commit)} {
		if strings.Contains(p, "in one turn") || strings.Contains(p, "ready && flai move") {
			t.Errorf("told to batch its calls again:\n%s", p)
		}
	}
}

// S-0176: the story's agent plans its tasks in layers as it writes them,
// records why in Decisions, and hands each task to a task sub-agent, a layer
// at once when its tasks are long, whose work it reviews and commits itself,
// launching each in the foreground and never waiting for one by ending its
// turn or with wait_for_events (S-0285); an answered or commit run is not told
// again.
func TestThePromptAsksForThePlan(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	restarted := r
	restarted.Restart = "ended (exit 1)"
	for _, p := range []string{Prompt(r), Prompt(restarted)} {
		for _, w := range []string{
			"Plan S-0104's tasks as you write them",
			"the paths it touches (--touches)",
			"the tasks of S-0104 it must wait for (--after)",
			"tasks with no after between them and no path in common form layers that can run together",
			"record the layers and why each task waits in the narrative's Decisions",
			"Work the plan layer by layer, handing each task to a task sub-agent",
			"A task sub-agent edits only what its task touches, runs only its own tests with flai test, and never commits, syncs, or writes through flai",
			"Name the task's ID in each task sub-agent's description, so that flai measures the task by its calls",
			// S-0285: a sub-agent is launched in the foreground, the
			// session is not left to end under a running one, and
			// wait_for_events is for a thread awaiting the designer.
			"Launch every sub-agent with the Agent tool's run_in_background set to false, a layer's in one message, so that each result comes back as the tool's result however long the sub-agent runs",
			"Never end your turn while a sub-agent runs in the background: Claude Code ends this session ten minutes after the turn ends, and the sub-agent with it",
			"Never wait for a sub-agent with the flai MCP tool wait_for_events either, which reports work items and threads, not sub-agents, and which flai guard refuses while one runs: wait_for_events is for a thread awaiting the designer",
			"Review each one's work yourself, fix what falls short, commit it, sync and test as above, and move the task",
			"only you commit, sync the stream, move items, and talk to the designer",
		} {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
		if strings.Contains(p, "through the harness's notice that it has finished") {
			t.Errorf("still told to wait for a background sub-agent's notice (S-0285):\n%s", p)
		}
	}
	answered, commit := r, r
	answered.Answered, commit.Commit = "TH-0001", "/w"
	for _, p := range []string{Prompt(answered), Prompt(commit)} {
		if strings.Contains(p, "layer") {
			t.Errorf("told the plan again:\n%s", p)
		}
	}
}

// S-0299, I-0093: the story's agent is told before it launches a layer that
// a write under a .claude/ folder is never a sub-agent's, that it makes such
// writes itself once the layer is back, and that it leaves them in
// .flai-cache/ with the cp commands on a thread when the operator may be
// away; an answered or commit run is not told again.
func TestThePromptKeepsClaudeWritesFromSubAgents(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	restarted := r
	restarted.Restart = "ended (exit 1)"
	for _, p := range []string{Prompt(r), Prompt(restarted)} {
		for _, w := range []string{
			"A write to a path Claude Code protects, such as a file in a .claude/ folder or .mcp.json, is never a sub-agent's: flai guard refuses it unless the operator has turned on auto-approve",
			"say so in the prompt of every sub-agent whose task changes such a file, and ask it to return the file's whole new content in its final message",
			"Make those writes yourself once the layer's sub-agents are back, never while a layer runs: permission_prompt opens a thread on S-0104 and holds the call until the operator answers",
			"write each whole file into the worktree's ignored .flai-cache/ folder instead, on a path with no .claude folder along it, open one thread on S-0104 with the exact cp commands that put each in place, and end rather than wait",
		} {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
		// Told before the agent launches anything: ahead of the review
		// of each task sub-agent's work.
		if strings.Index(p, ".claude/ folder") > strings.Index(p, "Review each one's work yourself") {
			t.Errorf("the .claude/ warning comes after the review of a layer:\n%s", p)
		}
	}
	answered, commit := r, r
	answered.Answered, commit.Commit = "TH-0001", "/w"
	for _, p := range []string{Prompt(answered), Prompt(commit)} {
		if strings.Contains(p, ".claude/") {
			t.Errorf("told about .claude/ writes again:\n%s", p)
		}
	}
}

// S-0197, ADR-0069: at each task the story's agent commits the task, syncs
// with flai stream sync and resolves what it lists, then runs the task's
// tests, in that order, never rebasing or merging by hand; and before review
// it commits everything and syncs again before the close-out run.
func TestThePromptSaysThePerTaskCycle(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	restarted := r
	restarted.Restart = "ended (exit 1)"
	for _, p := range []string{Prompt(r), Prompt(restarted)} {
		cycle := []string{
			"When a task is done, commit its changes, with its docs and work-item updates, on story/S-0104",
			"then run flai stream sync S-0104, and resolve each conflict it lists in the worktree, git add it, and git rebase --continue",
			"then run the task's tests with flai test and the paths it changed (or the flai MCP tool test with them) and commit any fix they need",
			"Before you move S-0104 to review, commit everything, run flai stream sync S-0104 again and resolve what it lists, then have one fresh verifier",
			"through the project's close-out script where it has one (it refuses a branch that does not contain the main branch)",
		}
		at := 0
		for _, w := range cycle {
			i := strings.Index(p[at:], w)
			if i < 0 {
				t.Errorf("prompt lacks %q after what comes before it:\n%s", w, p)
				continue
			}
			at += i + len(w)
		}
		for _, w := range []string{
			"refuses while anything is uncommitted: never start a rebase or merge by hand",
			"call the flai MCP tool inbox at every task transition",
		} {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
	}
	answered := r
	answered.Answered = "TH-0001"
	if p := Prompt(answered); !strings.Contains(p, "commit everything outstanding in the worktree, so that git status there is clean, run flai stream sync S-0104 again and resolve what it lists, then move S-0104 to review") {
		t.Errorf("an answered agent is not told to sync before review:\n%s", p)
	}
}

// S-0282, ADR-0089: every run working a story is told to tick each acceptance
// criterion through flai once it has verified it, never by editing the
// story's file; a full run is also told to ask each task sub-agent which
// criteria its task meets and to tick them after its review.
func TestThePromptSaysHowCriteriaAreTicked(t *testing.T) {
	r := req(&manifest.Agent{Harness: ClaudeCode})
	answered := r
	answered.Answered = "TH-0001"
	for _, p := range []string{Prompt(r), Prompt(answered)} {
		for _, w := range []string{
			"Tick each acceptance criterion with flai criteria tick S-0104 <n> (or the flai MCP tool criteria_tick) as soon as you have verified it, never by editing S-0104's file",
			"leave one you cannot verify here unticked and say why in S-0104's notes",
		} {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
		}
	}
	p := Prompt(r)
	for _, w := range []string{
		"ask each to name in its final message the acceptance criteria its task meets, by their numbers in flai criteria list S-0104",
		"and move the task, then tick the criteria your review verified it meets",
	} {
		if !strings.Contains(p, w) {
			t.Errorf("prompt lacks %q:\n%s", w, p)
		}
	}
}

// S-0189: a story whose agent gives its roles a model starts claude-code
// with --agents: the project's definition of each role's sub-agent, with the
// role's model over the definition's; what claude-code cannot run is
// refused; the operator's command is told the roles.
func TestClaudeCodeRunsEachRolesModel(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	verifier := "---\nname: verifier\ndescription: Runs the checks.\ntools: Read, Bash, mcp__flai__prime\nmodel: sonnet\n---\n\nYou are the verifier.\n"
	if err := os.WriteFile(filepath.Join(dir, "verifier.md"), []byte(verifier), 0o644); err != nil {
		t.Fatal(err)
	}
	r := req(&manifest.Agent{Harness: ClaudeCode, Model: "claude-opus-5-5", Roles: map[string]manifest.Role{
		"verify":  {Model: "claude-haiku-4-5"},
		"explore": {Harness: ClaudeCode},
	}})
	r.Root = root
	st, err := (claudeCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	var agents map[string]map[string]any
	if err := json.Unmarshal([]byte(after(st.Argv, "--agents")), &agents); err != nil {
		t.Fatalf("--agents: %v in %v", err, st.Argv)
	}
	v := agents["verifier"]
	if v["model"] != "claude-haiku-4-5" || v["description"] != "Runs the checks." || v["prompt"] != "You are the verifier." || v["name"] != nil {
		t.Errorf("verifier: %v", v)
	}
	if tools, _ := v["tools"].([]any); len(tools) != 3 || tools[2] != "mcp__flai__prime" {
		t.Errorf("tools: %v", v["tools"])
	}
	if _, ok := agents["explorer"]; ok {
		t.Errorf("a role with no model of its own was passed: %v", agents)
	}
	if i, j := slices.Index(st.Argv, "--agents"), slices.Index(st.Argv, "--permission-mode"); i < 0 || j < i {
		t.Errorf("--agents comes before the operator's arguments: %v", st.Argv)
	}

	// no roles, no flag
	plain := req(&manifest.Agent{Harness: ClaudeCode})
	plain.Root = root
	if st, _ := (claudeCode{}).Start(plain, Host{}); slices.Contains(st.Argv, "--agents") {
		t.Errorf("--agents with no roles: %v", st.Argv)
	}

	for _, c := range []struct {
		role manifest.Role
		want string
	}{
		{manifest.Role{Harness: "codex", Model: "m"}, "cannot run role verify on harness"},
		{manifest.Role{Model: "m", Config: map[string]string{"effort": "x"}}, "takes no config for role verify"},
	} {
		role, want := c.role, c.want
		bad := req(&manifest.Agent{Harness: ClaudeCode, Roles: map[string]manifest.Role{"verify": role}})
		bad.Root = root
		if _, err := (claudeCode{}).Start(bad, Host{}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%+v: %v", role, err)
		}
	}
	for name, want := range map[string]string{"plan": `no sub-agent for role "plan"`, "explore": "explorer, and its definition could not be read"} {
		bad := req(&manifest.Agent{Harness: ClaudeCode, Roles: map[string]manifest.Role{name: {Model: "m"}}})
		bad.Root = root
		if _, err := (claudeCode{}).Start(bad, Host{}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v", name, err)
		}
	}

	cmd, err := (command{}).Start(r, Host{Program: "run-agent"})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(cmd.Env, `FLAI_AGENT_ROLES={"explore":{"harness":"claude-code"},"verify":{"model":"claude-haiku-4-5"}}`) {
		t.Errorf("command env: %v", cmd.Env)
	}
}

func planReq(item string, a *manifest.Agent) Request {
	return Request{Role: "plan", Item: item, Root: "/p/flow", Project: "flow", Agent: a, Name: "planner-" + item, Flai: "/usr/local/bin/flai"}
}

// S-0208: the planner is asked to plan its item as its state calls for,
// through flai alone, to ask the operator on the item for a missing input,
// and to end with the summary flai serve logs; never to work a story.
func TestThePlannersPromptPlansItsItem(t *testing.T) {
	common := func(id, flag string) []string {
		return []string{
			"You are planner-" + id + ", the planner",
			"Plan " + id + ", and nothing else, as design/conventions/strategic-agents.md says under As the planner",
			"flai prime --role plan --" + flag + " " + id + " (or the flai MCP tool prime with role plan and " + flag + " " + id + ")",
			"A brief is not the document", "doc_get and its heading", "doc_search",
			"Call the flai MCP tool inbox", "Read " + id + " with item_get, and what it links with item_get and doc_get",
			"to the explorer with the Agent tool",
			"write only through flai: the flai MCP tools item_new and item_edit, or the flai CLI",
			"Never edit a file yourself, code or anything else, never move an item past backlog, and never finalize a draft",
			"Never overwrite the operator's inputs",
			"ask with the flai MCP tool thread_open on " + id + ", your recommended answer first, plan what needs no answer meanwhile, and hold the flai MCP tool wait_for_events",
			"flai serve logs the run's activity in wip/agents/planner.md",
		}
	}
	for _, c := range []struct {
		id, flag string
		want     []string
	}{
		{"E-0016", "epic", []string{
			"If it has no stories, draft the stories that deliver its outcome",
			"create each with draft true in the backlog (item_new's draft, or flai story new --draft)",
			"If it has stories, revisit each one not done or cancelled against the epic's outcome.",
			// S-0300: each story it drafts is enriched as a story's planner would
			"Enrich each story you draft, and again each one you revisit, as you would a story, S-nnnn being its ID: its predicted touches, a forecast, and a cost of delay value",
			"Run flai touches suggest S-nnnn, adding", "Run flai forecast S-nnnn and flai cod S-nnnn", "under a ### Planning heading",
			// S-0295: file-level touches, a folder kept only for files no task can name yet
			"Name files, not folders, in the touches you write for the story and for each of its tasks: keep a folder touch only where the story may add files there that no task can name yet",
			"make each one you write pass flai check --strict",
			// S-0300: and its tasks drafted, in layers, with no draft flag of their own
			"For each story you draft, and each one you revisit that is a draft with no tasks, draft the tasks that deliver its outcome",
			"A task carries no draft flag: it is a draft because its story is one",
			// S-0209: one thread holds the plan and what the planner proposes
			"Open one thread on E-0016 that summarises the plan: the stories, their order (their after), each story's tasks and their layers, and the assumptions you made",
			"In that same thread, propose each story you would split, merge, add, or drop, and create drafts for the additions only",
			"never cancel a finalized story or rewrite its words without asking",
			"End with a one-line summary that names by ID the stories and tasks you created and the stories you revisited: flai serve logs",
		}},
		{"S-0208", "story", []string{
			"It is a story: enrich it with its predicted touches, a forecast, and a cost of delay value worked out from the operator's inputs",
			// S-0210: from flai's three planning reads, reviewed, with its reasoning in the Notes
			"Run flai touches suggest S-0208, adding the paths its goal, criteria, and linked design name when it declares no touches",
			"predict its touches from what that lists, its goal and criteria, the design documents it links, and the code layout, keeping every touch it already declares",
			// S-0295 (ADR-0096): file-level touches, and why a folder was kept
			"Name files, not folders, in the touches you write for the story and for each of its tasks: keep a folder touch only where the story may add files there that no task can name yet, since a folder touch claims every file below it and, while the story is in progress, holds every ready story that touches one; record each folder touch you kept, and why, under the story's ### Planning heading",
			"Run flai forecast S-0208 and flai cod S-0208",
			"Review each figure, adjust it where you have a reason and state the reason",
			"write the touches, the forecast (duration, delivery, and basis), and the cost of delay value through flai: item_edit, or flai edit and flai touches",
			"In the story's Notes, under a ### Planning heading that is yours to rewrite, record where each touch came from (declared, co-change, design, or layout) and why each figure stands or was adjusted, and leave the rest of the Notes as it was. If it has no tasks",
			"Make each task you write pass flai check --strict and the markdown lint",
			"Open one thread on S-0208 that summarises the plan",
			"End with a one-line summary that names by ID the tasks you created and the tasks you revisited: flai serve logs",
		}},
	} {
		p := Prompt(planReq(c.id, nil))
		for _, w := range append(common(c.id, c.flag), c.want...) {
			if !strings.Contains(p, w) {
				t.Errorf("%s: prompt lacks %q:\n%s", c.id, w, p)
			}
		}
		for _, never := range []string{"flai stream open", "flai move", "worktree", "work story", "publish"} {
			if strings.Contains(p, never) {
				t.Errorf("%s: the planner's prompt says %q:\n%s", c.id, never, p)
			}
		}
	}
	if p := Prompt(planReq("S-0208", nil)); strings.Contains(p, "draft the stories") || strings.Contains(p, "split, merge, add, or drop") {
		t.Errorf("a story's planner is told to plan an epic's stories:\n%s", p)
	}
}

// S-0255: a story's planner drafts its tasks into the backlog in layers, or
// revisits the open ones without cancelling or rewriting them unasked, sums
// up the plan in one thread, and names the tasks it created and revisited in
// the summary flai serve logs. An epic's planner drafts the tasks of each
// story it drafts in the same words (S-0300), but is told nothing of
// revisiting a story's tasks or of the story's own thread and summary.
func TestAStorysPlannerDraftsAndRevisitsItsTasks(t *testing.T) {
	shared := []string{
		"draft the tasks that deliver its outcome, each with ## Work and ## Done when in its body, a nature, tags, touches (the paths it changes), and after (the tasks of the story it waits for)",
		"so that they form layers as work-management.md says",
		"create each in the backlog with the flai MCP tool item_new, type task and parent the story, or flai task new",
		"Tasks carry no topics: when a task reaches a topic the story lacks, add the topic to the story",
		"Make each task you write pass flai check --strict and the markdown lint: flai refuses one that does not, so fix what the refusal names",
	}
	storyOnly := []string{
		"If it has no tasks, draft the tasks that deliver its outcome",
		"If it has tasks, revisit each one not done or cancelled against the story's outcome: re-enrich its touches and after with item_edit, create the tasks the outcome still lacks",
		"propose in the plan's thread any task you would split, merge, or drop",
		"Never cancel a task, or rewrite the title, Work, or Done when of a task you did not write, without asking on that thread",
		"Open one thread on S-0255 that summarises the plan",
		"In the plan's thread, name the tasks, their order and layers, and the assumptions you made",
		"End with a one-line summary that names by ID the tasks you created and the tasks you revisited",
	}
	p := Prompt(planReq("S-0255", nil))
	for _, w := range append(slices.Clone(shared), storyOnly...) {
		if !strings.Contains(p, w) {
			t.Errorf("a story's planner prompt lacks %q:\n%s", w, p)
		}
	}
	e := Prompt(planReq("E-0016", nil))
	for _, w := range shared {
		if !strings.Contains(e, w) {
			t.Errorf("an epic's planner is not told %q:\n%s", w, e)
		}
	}
	for _, w := range storyOnly {
		if strings.Contains(e, w) {
			t.Errorf("an epic's planner is told %q:\n%s", w, e)
		}
	}
}

// S-0208: claude-code runs the planner's session as the project's planner
// definition, passed with --agents beside the explorer and the verifier, with
// the planner's model over the definition's, and tells it its role and item
// in its environment; what is not a planner's request is refused.
func TestClaudeCodeStartsThePlanner(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, def := range map[string]string{
		"planner":  "---\nname: planner\ndescription: Plans one item.\ntools: Read, mcp__flai__item_new\nmodel: inherit\n---\n\nYou are the planner.\n",
		"explorer": "---\nname: explorer\ndescription: Searches.\ntools: Read\nmodel: haiku\n---\n\nYou are the explorer.\n",
		"verifier": "---\nname: verifier\ndescription: Runs the checks.\ntools: Read, Bash\nmodel: sonnet\n---\n\nYou are the verifier.\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(def), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r := planReq("E-0016", &manifest.Agent{Harness: ClaudeCode, Model: "claude-opus-5-5", Roles: map[string]manifest.Role{
		"explore": {Model: "claude-haiku-4-5"},
		"verify":  {Model: "claude-sonnet-5"},
	}})
	r.Root = root
	st, err := (claudeCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	if p := st.Argv[2]; !strings.Contains(p, "asked for E-0016 to be planned") || !strings.Contains(p, "--role plan --epic E-0016") {
		t.Errorf("prompt: %s", p)
	}
	for flag, want := range map[string]string{"--agent": "planner", "--model": "claude-opus-5-5", "--name": "flow plan E-0016"} {
		if got := after(st.Argv, flag); got != want {
			t.Errorf("%s = %q, want %q (%q)", flag, got, want, st.Argv)
		}
	}
	var agents map[string]map[string]any
	if err := json.Unmarshal([]byte(after(st.Argv, "--agents")), &agents); err != nil {
		t.Fatalf("--agents: %v in %v", err, st.Argv)
	}
	for def, want := range map[string]string{"planner": "claude-opus-5-5", "explorer": "claude-haiku-4-5", "verifier": "claude-sonnet-5"} {
		if got := agents[def]; got["model"] != want || got["prompt"] == "" || got["name"] != nil {
			t.Errorf("%s: %v", def, got)
		}
	}
	if i, j := slices.Index(st.Argv, "--agent"), slices.Index(st.Argv, "--permission-mode"); i < 0 || j < i {
		t.Errorf("--agent comes before the operator's arguments: %v", st.Argv)
	}
	// the planner has flai's permission handler too, which denies anything
	// outside a story's worktree (S-0257)
	if got := after(st.Argv, "--permission-prompt-tool"); got != PermissionPromptTool {
		t.Errorf("planner: --permission-prompt-tool = %q (%v)", got, st.Argv)
	}
	if want := []string{"FLAI_ROLE=plan", "FLAI_ITEM=E-0016"}; !slices.Equal(st.Env, want) {
		t.Errorf("env %q, want %q", st.Env, want)
	}
	// the project's settings have flai guard judge the planner's file edits
	if slices.Contains(st.Argv, "--settings") {
		t.Errorf("a planner run is passed --settings: %v", st.Argv)
	}

	// a planner's agent with no model and no roles: the definition's model,
	// and the planner alone
	plain := planReq("S-0208", &manifest.Agent{Harness: ClaudeCode})
	plain.Root = root
	st, err = (claudeCode{}).Start(plain, Host{})
	if err != nil {
		t.Fatal(err)
	}
	agents = nil
	if err := json.Unmarshal([]byte(after(st.Argv, "--agents")), &agents); err != nil || len(agents) != 1 || agents["planner"]["model"] != "inherit" {
		t.Errorf("--agents %v: %v", agents, err)
	}
	if slices.Contains(st.Argv, "--model") || !slices.Contains(st.Env, "FLAI_ITEM=S-0208") {
		t.Errorf("argv %q, env %q", st.Argv, st.Env)
	}

	// a story's agent is not the planner
	story := req(&manifest.Agent{Harness: ClaudeCode})
	story.Root = root
	if st, _ := (claudeCode{}).Start(story, Host{}); slices.Contains(st.Argv, "--agent") || slices.Contains(st.Argv, "--settings") || st.Env != nil {
		t.Errorf("a story's agent: argv %q, env %q", st.Argv, st.Env)
	}

	missing := planReq("E-0016", nil)
	missing.Root = t.TempDir()
	bad := map[string]Request{"runs as the agent planner, and its definition could not be read": missing}
	for want, mod := range map[string]func(*Request){
		`flai cannot start an agent in role "review"`: func(r *Request) { r.Role = "review" },
		"was given story S-0104 to work":              func(r *Request) { r.Story = "S-0104" },
		`"T-0784" is neither`:                         func(r *Request) { r.Item = "T-0784" },
		`"" is neither`:                               func(r *Request) { r.Item = "" },
		"only the analyzer takes a focus":             func(r *Request) { r.Focus = "risk" },
	} {
		x := planReq("E-0016", nil)
		x.Root = root
		mod(&x)
		bad[want] = x
	}
	for want, x := range bad {
		if _, err := (claudeCode{}).Start(x, Host{}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%+v: %v, want %q", x, err, want)
		}
		if want != "runs as the agent planner, and its definition could not be read" {
			if _, err := (command{}).Start(x, Host{Program: "run"}); err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("command %+v: %v, want %q", x, err, want)
			}
		}
	}

	// the operator's command is told the role and the item, with no prompt
	cmd, err := (command{}).Start(planReq("E-0016", nil), Host{Program: "run-agent", Args: []string{"{story}"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"run-agent", ""}; !slices.Equal(cmd.Argv, want) {
		t.Errorf("argv %q, want %q", cmd.Argv, want)
	}
	if !slices.Contains(cmd.Env, "FLAI_ROLE=plan") || !slices.Contains(cmd.Env, "FLAI_ITEM=E-0016") {
		t.Errorf("command env: %q", cmd.Env)
	}
}

func orchestrateReq(a *manifest.Agent) Request {
	return Request{Role: "orchestrate", Root: "/p/flow", Project: "flow", Agent: a, Name: "orchestrator", Flai: "/usr/local/bin/flai"}
}

// S-0218: the orchestrator is asked to keep the project's work moving within
// the operator's permissions, with flai's figures, logging each decision and
// waiting on events between them without ending; never to edit a file, work
// a story, or work around flai guard.
func TestTheOrchestratorsPromptKeepsWorkMovingWithinItsPermissions(t *testing.T) {
	p := Prompt(orchestrateReq(nil))
	for _, w := range []string{
		"You are orchestrator, the orchestrator, started by flai serve on this host because the operator turned on the orchestrate host action for the project at /p/flow",
		"as design/conventions/strategic-agents.md says under As the orchestrator",
		"flai prime --role orchestrate (or the flai MCP tool prime with role orchestrate)",
		"A brief is not the document", "doc_get and its heading", "doc_search",
		"Call the flai MCP tool inbox, and read the board with the flai MCP tool board",
		"to the explorer with the Agent tool",
		"Act only within the permissions the operator sets in system-flow.yaml under orchestration.permissions, each off by default, and by its policy, orchestration.policy",
		"with the flai MCP tool plan, only while plan_backlog_epics is on", "only while finalize_drafts is on", "only while promote_to_ready is on",
		"only while order_ready is on", "only as answer_threads says", "only while accept_reviews is on", "publish a release only while publish is on",
		"if a permission is unclear, ask; do not act",
		"never do the arithmetic yourself",
		"flai order --by (the flai MCP tool order_by_policy)", "flai promote --candidates (promote_candidates)", "flai release --evaluate (release_evaluate)",
		"Log each action when you have taken it with the flai MCP tool activity_log, kind orchestrator: what you did, to which item, why, and the policy figure that justified it",
		"Then hold the flai MCP tool wait_for_events, again each time it returns",
		"Repeat without ending",
		"write only through flai",
		"Never edit code or documents, and never work a story yourself",
		"flai guard refuses a call outside your permissions and names the permission it needs: never work around a refusal",
		"When a decision needs the operator, ask with the flai MCP tool thread_open on the item it concerns, your recommended answer first",
	} {
		if !strings.Contains(p, w) {
			t.Errorf("the orchestrator's prompt lacks %q:\n%s", w, p)
		}
	}
	// it names a story's worktree only to hand it to the verifier (S-0221)
	for _, never := range []string{"flai stream open", "flai stream sync", "flai move", "Work it in the worktree", "As the planner", "End with"} {
		if strings.Contains(p, never) {
			t.Errorf("the orchestrator's prompt says %q:\n%s", never, p)
		}
	}
}

// S-0229: the orchestrator is told that the operator may change its
// permissions, policy, and release policy while it runs, that it reads them
// again before each decision, and that the guard's verdict on a call holds.
// No other agent is told it.
func TestTheOrchestratorsPromptSaysItsPermissionsMayChangeWhileItRuns(t *testing.T) {
	const w = "The operator may change your permissions, your policy, and the release policy, orchestration.release, while you run: read them again in system-flow.yaml in the main checkout before each decision rather than keep what you read at the start, and when flai guard's verdict on a call differs from what you read, the guard's verdict holds"
	if p := Prompt(orchestrateReq(nil)); !strings.Contains(p, w) {
		t.Errorf("the orchestrator's prompt lacks %q:\n%s", w, p)
	}
	others := map[string]string{"a story's agent": Prompt(req(nil)), "an epic's planner": Prompt(planReq("E-0016", nil)), "a story's planner": Prompt(planReq("S-0208", nil))}
	for who, o := range others {
		if strings.Contains(o, "The operator may change your permissions") {
			t.Errorf("%s's prompt says its permissions may change while it runs:\n%s", who, o)
		}
	}
}

// S-0219: the orchestrator is told, for each of plan_backlog_epics,
// finalize_drafts, promote_to_ready, and order_ready, which flai command
// gives it the work and what it does with it; to log each action with the
// policy figure that justified it; and that a refusal ends the attempt,
// which it logs and does not retry until something changes.
func TestTheOrchestratorsPromptSaysWhatEachPermissionDoes(t *testing.T) {
	p := Prompt(orchestrateReq(nil))
	for _, w := range []string{
		// plan_backlog_epics
		"While plan_backlog_epics is on, run flai plan --candidates and start the planner with the flai MCP tool plan for each epic it lists, one at a time, and for no other",
		// finalize_drafts
		"While finalize_drafts is on, run flai promote --drafts",
		"finalize a draft it lists as complete, and whose criteria, touches, forecast, and value you judge consistent, with the flai MCP tool item_edit giving only its id and draft false",
		"for any other draft, open one thread on the story saying what it lacks or what is inconsistent, once, and leave the draft as it is",
		// promote_to_ready
		"While promote_to_ready is on, run flai promote --candidates and move its candidates to ready in its order with the flai MCP tool item_move while the ready column's limit has room",
		"never a draft, and never a held story",
		// order_ready
		"While order_ready is on, run flai order --by <policy> --apply, with the policy orchestration.policy names, after each change to the ready column, yours or another's",
		"it keeps in place a story placed by hand within the last day",
		"Never place a story by hand yourself",
		// the figures, the log, and a refusal
		"the epics to plan from flai plan --candidates", "the drafts complete enough to finalize from flai promote --drafts",
		"the policy figure that justified it, as flai gave it: the story's value, its value over its duration, its forecast, or the candidate's rank",
		"A refusal from flai guard or from flai ends that attempt: log it with activity_log, with the refusal, and do not try it again until something it depends on changes",
	} {
		if !strings.Contains(p, w) {
			t.Errorf("the orchestrator's prompt lacks %q:\n%s", w, p)
		}
	}
}

// S-0220: the orchestrator is told what to do with the threads awaiting the
// operator under each value of answer_threads, which it reads afresh since
// the operator may change it while it runs: nothing while it is off, a
// sourced recommendation while it is recommend, and while it is autonomous
// a sourced answer, or a recommendation escalating to the operator when no
// source settles the question or it asks for the operator's judgement. It
// never resolves a thread it did not open, answers its own, or confirms a
// recommendation. No other agent is told any of it.
func TestTheOrchestratorsPromptSaysWhatToDoWithThreads(t *testing.T) {
	rules := map[string][]string{
		"each value": {
			"Reply on threads only as answer_threads says, and read it in system-flow.yaml in the main checkout each time before you act on threads: the operator may change it while you run",
			"Take from inbox the threads awaiting the operator, those whose status is open and whose last entry is by a story's agent, leaving out the threads you opened and those whose pending_recommendation is not null",
			"thread_reply logs the reply and its source in your decision log",
			"Never resolve a thread you did not open, never answer a thread you opened, and never confirm a recommendation: the operator does",
		},
		"off": {"While answer_threads is off, leave them alone"},
		"recommend": {
			"While it is recommend, reply to each with the flai MCP tool thread_reply, recommendation true, and a source: the ADR, design section, or convention your answer rests on, as <path> or <path>#<heading>, read with doc_get first",
		},
		"autonomous": {
			"While it is autonomous, answer with thread_reply and a source when a source settles the question",
			"post a recommendation instead, citing what it draws on, which escalates it to the operator, when no source settles it, or when it asks for the operator's judgement: a decision not yet recorded, a change of scope, or money, such as a cost of delay input, an estimate, or spend",
		},
	}
	p := Prompt(orchestrateReq(nil))
	others := map[string]string{"a story's agent": Prompt(req(nil)), "an epic's planner": Prompt(planReq("E-0016", nil)), "a story's planner": Prompt(planReq("S-0208", nil))}
	for value, ws := range rules {
		for _, w := range ws {
			if !strings.Contains(p, w) {
				t.Errorf("the orchestrator's prompt lacks, for %s, %q:\n%s", value, w, p)
			}
			for who, o := range others {
				if strings.Contains(o, w) {
					t.Errorf("%s's prompt says %q", who, w)
				}
			}
		}
	}
	for who, o := range others {
		if strings.Contains(o, "answer_threads") || strings.Contains(o, "pending_recommendation") {
			t.Errorf("%s's prompt names answer_threads or pending_recommendation:\n%s", who, o)
		}
	}
}

// S-0221: while accept_reviews is on, the orchestrator hands a story in
// review's worktree to its verifier, reads flai accept --dry-run's blockers
// at the commit verified, and accepts with --verified and evidence for every
// criterion, passed on standard input since it cannot write a file, logging
// the acceptance; otherwise it leaves the story in review and says on a
// thread what is missing. It never moves a story to done with item_move. No
// other agent is told any of it (ADR-0093).
func TestTheOrchestratorsPromptSaysHowItAcceptsAStory(t *testing.T) {
	p := Prompt(orchestrateReq(nil))
	for _, w := range []string{
		"While accept_reviews is on, take each story in review in turn",
		// the verifier
		"Hand its worktree, /p/flow/.flai-cache/worktrees/<S-nnnn>, to the verifier with the Agent tool",
		"it runs the tests, the lint, and flai check --strict, or the project's close-out script where it has one",
		"checks the diff against each acceptance criterion, naming for each the changed files that meet it, and names the commit it verified",
		// the preview
		"Run flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run with that commit, and read the blockers",
		// the acceptance and its evidence
		"With no blocker and every criterion matched to changed files, write the evidence, a Verdict: line from the verifier's report and one list item per criterion, - <n>: <files>",
		"run flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -, passing the evidence on standard input with a heredoc in the shell, since you cannot write a file",
		"then log the acceptance with activity_log, naming the story and the commit",
		// what is missing
		"Otherwise leave the story in review, open a thread on it with the flai MCP tool thread_open saying what is missing, each blocker and each criterion you could not check against the diff, and log that decision",
		"A commit added after the verifier's run makes flai refuse: verify again",
		"Never move a story to done with item_move, and never accept a story while accept_reviews is off",
	} {
		if !strings.Contains(p, w) {
			t.Errorf("the orchestrator's prompt lacks %q:\n%s", w, p)
		}
	}
	others := map[string]string{"a story's agent": Prompt(req(nil)), "an epic's planner": Prompt(planReq("E-0016", nil)), "a story's planner": Prompt(planReq("S-0208", nil))}
	for who, o := range others {
		for _, w := range []string{"accept_reviews", "--by orchestrator", "--verified", "--evidence"} {
			if strings.Contains(o, w) {
				t.Errorf("%s's prompt says %q", who, w)
			}
		}
	}
}

// S-0222: while publish is on, the orchestrator evaluates the release policy
// after each acceptance and publishes through release_publish when a
// threshold or theme is met, or under judgement when it judges the batch
// coherent and complete, never a batch whole_epics holds back; it logs each
// release from what release_publish returned, each decision not to publish,
// and each refusal, which it raises on a thread to the operator. It
// publishes no other way. No other agent is told any of it.
func TestTheOrchestratorsPromptSaysHowItPublishes(t *testing.T) {
	p := Prompt(orchestrateReq(nil))
	for _, w := range []string{
		"While publish is on, after each acceptance, yours or one wait_for_events reports as a story moved to done, call the flai MCP tool release_evaluate",
		// threshold and theme
		"Under the threshold or theme policy, when the evaluation is met, call the flai MCP tool release_publish with the figure that was met as its reason",
		// judgement
		"Under judgement, call release_publish only when you judge the unreleased work coherent and complete, with that reasoning as its reason",
		// whole_epics
		"Under any policy, never publish a batch the evaluation holds back under whole_epics, the stories it names in held_by_epic",
		// the log of a release and of a decision not to publish
		"Log each release with activity_log: the policy, its figures, the versions and tags, and the items bundled, all as release_publish returned them",
		"Log a decision not to publish too, with the evaluation's figures",
		// a refusal: its log and its thread
		"When release_publish refuses, log the refusal with activity_log",
		"open a thread to the operator with the flai MCP tool thread_open on the most recently accepted story of the batch, with the refusal's words and what fixes it",
		"do not try again until that thread is answered or the next acceptance",
		// off, and no other route
		"While publish is off, neither evaluate nor publish",
		"Publish only through release_publish: flai guard refuses you flai release other than --evaluate, flai push, git push, and git tag, whatever your permissions",
	} {
		if !strings.Contains(p, w) {
			t.Errorf("the orchestrator's prompt lacks %q:\n%s", w, p)
		}
	}
	for _, never := range []string{"flai release --pending", "flai push --pending"} {
		if strings.Contains(p, never) {
			t.Errorf("the orchestrator's prompt says %q:\n%s", never, p)
		}
	}
	others := map[string]string{"a story's agent": Prompt(req(nil)), "an epic's planner": Prompt(planReq("E-0016", nil)), "a story's planner": Prompt(planReq("S-0208", nil))}
	for who, o := range others {
		for _, w := range []string{"release_publish", "release_evaluate", "whole_epics", "held_by_epic"} {
			if strings.Contains(o, w) {
				t.Errorf("%s's prompt says %q", who, w)
			}
		}
	}
}

// S-0218: claude-code runs the orchestrator's session as the project's
// orchestrator definition, passed with --agents beside the explorer and the
// verifier, named for the project and the role, and tells it its role alone;
// a project with no definition, and an orchestrator given a story or an
// item, are refused.
func TestClaudeCodeStartsTheOrchestrator(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, def := range map[string]string{
		"orchestrator": "---\nname: orchestrator\ndescription: Keeps work moving.\ntools: Read, mcp__flai__item_move\nmodel: inherit\n---\n\nYou are the orchestrator.\n",
		"explorer":     "---\nname: explorer\ndescription: Searches.\ntools: Read\nmodel: haiku\n---\n\nYou are the explorer.\n",
		"verifier":     "---\nname: verifier\ndescription: Runs the checks.\ntools: Read, Bash\nmodel: sonnet\n---\n\nYou are the verifier.\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(def), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r := orchestrateReq(&manifest.Agent{Harness: ClaudeCode, Model: "claude-opus-5-5", Roles: map[string]manifest.Role{
		"explore": {Model: "claude-haiku-4-5"},
		"verify":  {Model: "claude-sonnet-5"},
	}})
	r.Root = root
	st, err := (claudeCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	if p := st.Argv[2]; !strings.Contains(p, "the orchestrator, started by flai serve") || !strings.Contains(p, "flai prime --role orchestrate") {
		t.Errorf("prompt: %s", p)
	}
	for flag, want := range map[string]string{"--agent": "orchestrator", "--model": "claude-opus-5-5", "--name": "flow orchestrate"} {
		if got := after(st.Argv, flag); got != want {
			t.Errorf("%s = %q, want %q (%q)", flag, got, want, st.Argv)
		}
	}
	var agents map[string]map[string]any
	if err := json.Unmarshal([]byte(after(st.Argv, "--agents")), &agents); err != nil {
		t.Fatalf("--agents: %v in %v", err, st.Argv)
	}
	for def, want := range map[string]string{"orchestrator": "claude-opus-5-5", "explorer": "claude-haiku-4-5", "verifier": "claude-sonnet-5"} {
		if got := agents[def]; got["model"] != want || got["prompt"] == "" || got["name"] != nil {
			t.Errorf("%s: %v", def, got)
		}
	}
	if _, ok := agents["planner"]; ok {
		t.Errorf("the orchestrator is passed the planner: %v", agents)
	}
	if got := after(st.Argv, "--permission-prompt-tool"); got != PermissionPromptTool {
		t.Errorf("orchestrator: --permission-prompt-tool = %q (%v)", got, st.Argv)
	}
	if want := []string{"FLAI_ROLE=orchestrate"}; !slices.Equal(st.Env, want) {
		t.Errorf("env %q, want %q", st.Env, want)
	}

	// an orchestrator with no agent of its own: the definition's model, and
	// the orchestrator alone
	plain := orchestrateReq(nil)
	plain.Root = root
	st, err = (claudeCode{}).Start(plain, Host{})
	if err != nil {
		t.Fatal(err)
	}
	agents = nil
	if err := json.Unmarshal([]byte(after(st.Argv, "--agents")), &agents); err != nil || len(agents) != 1 || agents["orchestrator"]["model"] != "inherit" {
		t.Errorf("--agents %v: %v", agents, err)
	}

	missing := orchestrateReq(nil)
	missing.Root = t.TempDir()
	if _, err := (claudeCode{}).Start(missing, Host{}); err == nil ||
		!strings.Contains(err.Error(), "the orchestrator's session runs as the agent orchestrator, and its definition could not be read") ||
		!strings.Contains(err.Error(), "flai upgrade adds it from the template") {
		t.Errorf("a project with no orchestrator.md: %v", err)
	}
	for want, mod := range map[string]func(*Request){
		"the orchestrator was given story S-0104 to work": func(r *Request) { r.Story = "S-0104" },
		"the orchestrator was given E-0016":               func(r *Request) { r.Item = "E-0016" },
	} {
		x := orchestrateReq(nil)
		x.Root = root
		mod(&x)
		for name, ad := range map[string]Adapter{ClaudeCode: claudeCode{}, Command: command{}} {
			if _, err := ad.Start(x, Host{Program: "run"}); err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s %+v: %v, want %q", name, x, err, want)
			}
		}
	}

	// the operator's command is told the role and no item, with no prompt
	cmd, err := (command{}).Start(orchestrateReq(nil), Host{Program: "run-agent", Args: []string{"{story}"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"run-agent", ""}; !slices.Equal(cmd.Argv, want) {
		t.Errorf("argv %q, want %q", cmd.Argv, want)
	}
	if !slices.Contains(cmd.Env, "FLAI_ROLE=orchestrate") || slices.ContainsFunc(cmd.Env, func(e string) bool {
		return strings.HasPrefix(e, "FLAI_ITEM=") || strings.HasPrefix(e, "FLAI_STORY=")
	}) {
		t.Errorf("command env: %q", cmd.Env)
	}
}

// The template's own definitions read as --agents takes them.
func TestTheTemplatesDefinitionsRead(t *testing.T) {
	for _, def := range []string{"explorer", "verifier", "planner", "orchestrator", "analyzer"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "..", "template", "root", ".claude", "agents", def+".md"))
		if err != nil {
			t.Fatal(err)
		}
		got, err := definition(data)
		if err != nil {
			t.Fatalf("%s: %v", def, err)
		}
		if got["model"] == nil || got["prompt"] == "" || len(got["tools"].([]string)) == 0 {
			t.Errorf("%s: %v", def, got)
		}
	}
	// a definition with no prompt, ending at its front matter's last line
	if got, err := definition([]byte("---\ndescription: d\n---")); err != nil || got["prompt"] != "" {
		t.Errorf("no prompt and no last newline: %v %v", got, err)
	}
	for _, bad := range []string{"description: d\n", "---\ndescription: d\n", "---\nname: x\n---\nbody\n"} {
		if _, err := definition([]byte(bad)); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

// The orchestrator's definition, the project's and the template's alike,
// says what it does with each permission S-0219 gives work to, logs each
// action with its policy figure, and ends an attempt on a refusal.
func TestTheOrchestratorsDefinitionSaysWhatEachPermissionDoes(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	var copies []string
	for _, p := range []string{filepath.Join(root, ".claude", "agents", "orchestrator.md"), filepath.Join(root, "template", "root", ".claude", "agents", "orchestrator.md")} {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		copies = append(copies, string(data))
		for _, want := range []string{
			"With `plan_backlog_epics`, run `flai plan --candidates`",
			"With `finalize_drafts`, run `flai promote --drafts`",
			"with `item_edit` giving only its id and `draft: false`",
			"With `promote_to_ready`, run `flai promote --candidates`",
			"Never a draft, never a held story.",
			"With `order_ready`, run `flai order --by <policy> --apply`",
			"Never place a story by hand yourself.",
			"With `accept_reviews`, take each story in review in turn.",
			"`flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run`",
			"`flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -`",
			"open a thread on it with `thread_open` saying what is missing",
			"Never move a story to done with `item_move`.",
			"With `publish`, after each acceptance",
			"call `release_evaluate`",
			"call `release_publish` with the figure that was met as its reason",
			"Under `judgement`, publish only when you judge the unreleased work coherent and complete",
			"never publish a batch the evaluation holds back under `whole_epics` (its `held_by_epic`)",
			"Log each release with `activity_log`: the policy, its figures, the versions and tags, and the items bundled",
			"Log a decision not to publish too, with the figures.",
			"When `release_publish` refuses, log the refusal with `activity_log`, and open a thread to the operator with `thread_open` on the most recently accepted story of the batch",
			"Without `publish`, neither evaluate nor publish.",
			"mcp__flai__release_evaluate, mcp__flai__release_publish",
			"the policy figure that justified it",
			"ends that attempt: log it with the refusal",
		} {
			if !strings.Contains(string(data), want) {
				t.Errorf("%s does not say %q", p, want)
			}
		}
	}
	if copies[0] != copies[1] {
		t.Error("the project's and the template's orchestrator.md differ")
	}
}

func analyzeReq(focus string, a *manifest.Agent) Request {
	return Request{Role: "analyze", Focus: focus, Root: "/p/flow", Project: "flow", Agent: a, Name: "analyzer", Flai: "/usr/local/bin/flai"}
}

// S-0223: the analyzer is asked for one report on its focus, or on all three
// with none: primed with its role, reading flai stats --json, the design with
// doc_search and doc_get, and the issues, handing search of the code to the
// explorer; the report's path, front matter, and findings with their evidence,
// severity, and impact; its entry in the README; each actionable finding filed
// as an issue, or the open one that records it bumped, with its class,
// impact, evidence, and the report, and linked from the report (S-0224);
// nothing else edited and no story authored; and a summary that names the
// report, for the run's log.
func TestTheAnalyzersPromptAsksForOneReport(t *testing.T) {
	common := func(focus string) []string {
		return []string{
			"You are analyzer, the analyzer, started by flai serve on this host because the operator asked for an analysis, focus " + focus + ", of the project at /p/flow",
			"as design/conventions/strategic-agents.md says under As the analyzer",
			"flai prime --role analyze (or the flai MCP tool prime with role analyze)",
			"Call the flai MCP tool inbox",
			"Read the metrics with flai stats --json, and take every figure from it",
			"doc_search, which finds sections by their words, and doc_get and its heading", "a brief is not the document",
			"Read the issues under design/issues, design/issues/summary.md first",
			"Hand wide search of the code to the explorer with the Agent tool",
			"Write one report, design/analysis/<date>-" + focus + ".md, where <date> is today's date in UTC as YYYY-MM-DD",
			"front matter with title, updated, status draft while you write it and active once it is done, focus " + focus + ", and the window its metrics cover, from and to, as dates",
			"one section per finding, with its evidence (the metric figures as flai gave them, the file paths, and the design sections quoted), its severity, and its estimated impact: the time it loses per cycle, or the revenue or penalty it puts at stake where the design states them",
			"Add the report to design/analysis/README.md",
			"File each actionable finding as an issue, and make no story of it",
			"List the open issues with flai issue list --json, or read design/issues/summary.md",
			"When one already records the finding, under whatever title, bump it with flai issue bump <id> --report design/analysis/<date>-" + focus + ".md and the finding's impact (or the flai MCP tool issue_bump)",
			`Otherwise file one with flai issue new "<title>" --class <class> --report design/analysis/<date>-` + focus + ".md --json (or the flai MCP tool issue_new)",
			"the class from the finding, defect, efficiency, or impression for a risk with no measured instance",
			"its impact as --time-lost-per-cycle (a duration such as 4h), or --revenue-per-week or --penalty-per-week (an amount in planning.currency)",
			"--evidence, the evidence behind them",
			"flai bumps an open issue of the same title rather than opening a second",
			"the issue's Remediation section links your report",
			"Link each issue you filed or bumped from its finding in the report, as [I-nnnn](../issues/<file>), the file name of the path flai returns",
			"Never run flai issue story or flai issue close",
			"Edit nothing else: no code, no design, no issue, and no work item, and author no stories; the issues you file and bump, flai writes",
			"flai guard refuses an edit outside design/analysis and any write to a work item: never work around a refusal",
			"ask with the flai MCP tool thread_open on your report, your recommended answer first",
			"hold the flai MCP tool wait_for_events",
			"End with a one-line summary that names the report, design/analysis/<date>-" + focus + ".md with its date",
			"logs the run's activity in wip/agents/analyzer.md with it",
		}
	}
	finding := map[string]string{
		"bottlenecks": "bottlenecks in the flow of work, from the cumulative flow, the time items spend in each state, the time they wait, and the holds on them",
		"intent":      "gaps between what design/system says and what the code does",
		"risk":        "technical and security risks",
	}
	for _, focus := range append([]string{""}, Focuses...) {
		p := Prompt(analyzeReq(focus, nil))
		want := append(common(focus), "Look for "+finding[focus], ", and for nothing else.")
		if focus == "" {
			want = append(common(AllFocus), "Look for three kinds of finding: ")
		}
		for _, w := range want {
			if !strings.Contains(p, w) {
				t.Errorf("focus %q: the analyzer's prompt lacks %q:\n%s", focus, w, p)
			}
		}
		for f, w := range finding {
			if looks := strings.Contains(p, w); looks != (focus == "" || focus == f) {
				t.Errorf("focus %q: the prompt looks for %s: %v\n%s", focus, f, looks, p)
			}
		}
		for _, never := range []string{"flai stream open", "flai move", "worktree", "item_new", "item_edit", "As the planner", "As the orchestrator", "publish"} {
			if strings.Contains(p, never) {
				t.Errorf("focus %q: the analyzer's prompt says %q:\n%s", focus, never, p)
			}
		}
	}
}

// S-0223: claude-code runs the analyzer's session as the project's analyzer
// definition, beside the explorer, named for the project, the role, and the
// focus, and tells it its role alone in its environment; the operator's
// command is told the role and the focus. An analyzer given a story, an
// item, or a focus it does not take is refused, as is a focus for a story.
func TestClaudeCodeStartsTheAnalyzer(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, def := range map[string]string{
		"analyzer": "---\nname: analyzer\ndescription: Writes a report.\ntools: Read, Write\nmodel: inherit\n---\n\nYou are the analyzer.\n",
		"explorer": "---\nname: explorer\ndescription: Searches.\ntools: Read\nmodel: haiku\n---\n\nYou are the explorer.\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(def), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r := analyzeReq("risk", &manifest.Agent{Harness: ClaudeCode, Model: "claude-opus-5-5", Roles: map[string]manifest.Role{"explore": {Model: "claude-haiku-4-5"}}})
	r.Root = root
	st, err := (claudeCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	if p := st.Argv[2]; !strings.Contains(p, "the analyzer, started by flai serve") || !strings.Contains(p, "design/analysis/<date>-risk.md") {
		t.Errorf("prompt: %s", p)
	}
	for flag, want := range map[string]string{"--agent": "analyzer", "--model": "claude-opus-5-5", "--name": "flow analyze risk", "--permission-prompt-tool": PermissionPromptTool} {
		if got := after(st.Argv, flag); got != want {
			t.Errorf("%s = %q, want %q (%q)", flag, got, want, st.Argv)
		}
	}
	var agents map[string]map[string]any
	if err := json.Unmarshal([]byte(after(st.Argv, "--agents")), &agents); err != nil {
		t.Fatalf("--agents: %v in %v", err, st.Argv)
	}
	if len(agents) != 2 || agents["analyzer"]["model"] != "claude-opus-5-5" || agents["explorer"]["model"] != "claude-haiku-4-5" {
		t.Errorf("--agents %v", agents)
	}
	if want := []string{"FLAI_ROLE=analyze"}; !slices.Equal(st.Env, want) {
		t.Errorf("env %q, want %q", st.Env, want)
	}

	// with no focus: all of them, and the definition alone
	plain := analyzeReq("", nil)
	plain.Root = root
	st, err = (claudeCode{}).Start(plain, Host{})
	if err != nil {
		t.Fatal(err)
	}
	if got := after(st.Argv, "--name"); got != "flow analyze" || !strings.Contains(st.Argv[2], "design/analysis/<date>-all.md") {
		t.Errorf("--name %q, prompt %s", got, st.Argv[2])
	}
	if want := []string{"FLAI_ROLE=analyze"}; !slices.Equal(st.Env, want) {
		t.Errorf("env %q, want %q", st.Env, want)
	}

	missing := analyzeReq("", nil)
	missing.Root = t.TempDir()
	if _, err := (claudeCode{}).Start(missing, Host{}); err == nil || !strings.Contains(err.Error(), "the analyzer's session runs as the agent analyzer, and its definition could not be read") {
		t.Errorf("a project with no analyzer.md: %v", err)
	}
	for want, mod := range map[string]func(*Request){
		"the analyzer was given story S-0104 to work": func(r *Request) { r.Story = "S-0104" },
		"the analyzer was given E-0016":               func(r *Request) { r.Item = "E-0016" },
		`the analyzer takes the focus bottlenecks, intent, risk, or none for all of them, and "all" is none of them`: func(r *Request) { r.Focus = "all" },
	} {
		x := analyzeReq("", nil)
		x.Root = root
		mod(&x)
		for name, ad := range map[string]Adapter{ClaudeCode: claudeCode{}, Command: command{}} {
			if _, err := ad.Start(x, Host{Program: "run"}); err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s %+v: %v, want %q", name, x, err, want)
			}
		}
	}
	story := req(nil)
	story.Focus = "risk"
	if _, err := (command{}).Start(story, Host{Program: "run"}); err == nil || !strings.Contains(err.Error(), "only the analyzer takes a focus") {
		t.Errorf("a story's agent with a focus: %v", err)
	}

	// the operator's command is told the role and the focus, with no prompt
	cmd, err := (command{}).Start(analyzeReq("intent", nil), Host{Program: "run-agent", Args: []string{"{story}"}})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(cmd.Env, "FLAI_ROLE=analyze") || !slices.Contains(cmd.Env, "FLAI_FOCUS=intent") || slices.ContainsFunc(cmd.Env, func(e string) bool {
		return strings.HasPrefix(e, "FLAI_ITEM=") || strings.HasPrefix(e, "FLAI_STORY=")
	}) {
		t.Errorf("command env: %q", cmd.Env)
	}
	if cmd, _ := (command{}).Start(analyzeReq("", nil), Host{Program: "run-agent"}); slices.ContainsFunc(cmd.Env, func(e string) bool { return strings.HasPrefix(e, "FLAI_FOCUS=") }) {
		t.Errorf("command env with no focus: %q", cmd.Env)
	}
}

// S-0223: the template's analyzer definition has the reads, its report's
// edits, the shell, the explorer, and flai's reads, activity log, threads,
// and events, and no item write; S-0224 gives it issue_new and issue_bump,
// and has it file each actionable finding as an issue, bump the open one
// that records it, and link each from its report, but make no story of one,
// so it has no issue_story; the template's settings run flai guard on
// the analyzer's file edits as they do on the planner's and the
// orchestrator's.
func TestTheTemplatesAnalyzerIsHeldToItsReport(t *testing.T) {
	root := filepath.Join("..", "..", "..", "template", "root", ".claude")
	data, err := os.ReadFile(filepath.Join(root, "agents", "analyzer.md"))
	if err != nil {
		t.Fatal(err)
	}
	def, err := definition(data)
	if err != nil {
		t.Fatal(err)
	}
	tools, _ := def["tools"].([]string)
	for _, want := range []string{"Read", "Grep", "Glob", "Edit", "Write", "Bash", "Agent",
		"mcp__flai__prime", "mcp__flai__inbox", "mcp__flai__board", "mcp__flai__item_get", "mcp__flai__doc_get", "mcp__flai__doc_search",
		"mcp__flai__thread_get", "mcp__flai__who_touches", "mcp__flai__activity_log", "mcp__flai__thread_open", "mcp__flai__thread_reply", "mcp__flai__wait_for_events",
		"mcp__flai__issue_new", "mcp__flai__issue_bump"} {
		if !slices.Contains(tools, want) {
			t.Errorf("analyzer.md lacks the tool %s: %v", want, tools)
		}
	}
	for _, never := range []string{"mcp__flai__item_new", "mcp__flai__item_edit", "mcp__flai__item_move", "mcp__flai__issue_story", "mcp__flai__plan", "NotebookEdit"} {
		if slices.Contains(tools, never) {
			t.Errorf("analyzer.md has the tool %s", never)
		}
	}
	prompt, _ := def["prompt"].(string)
	for _, want := range []string{"role `analyze`", "`flai stats --json`", "`design/analysis/<date>-<focus>.md`", "`design/analysis/README.md`", "author no stories",
		// S-0224: each actionable finding filed as an issue, a duplicate bumped, and each linked from the report
		"File each actionable finding as an issue, and make no story of it", "`flai issue list --json`",
		"bump it with `flai issue bump <id> --report <your report>`", "`flai issue new \"<title>\" --class <class> --report <your report> --json`",
		"`impression` for a risk with no measured instance", "`--time-lost-per-cycle`", "`--revenue-per-week` or `--penalty-per-week`", "`--evidence`",
		"`[I-nnnn](../issues/<file>)`", "Never run `flai issue story` or `flai issue close`"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("analyzer.md does not say %q", want)
		}
	}
	data, err = os.ReadFile(filepath.Join(root, "settings.json"))
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
		t.Fatal(err)
	}
	for _, h := range settings.Hooks["PreToolUse"] {
		if h.Matcher != "Edit|MultiEdit|Write|NotebookEdit" {
			continue
		}
		if c := h.Hooks[0].Command; !strings.HasPrefix(c, `[ "$FLAI_ROLE" = plan ] || [ "$FLAI_ROLE" = orchestrate ] || [ "$FLAI_ROLE" = analyze ] || [ -n "$FLAI_STORY" ] || exit 0; `) || !strings.Contains(c, "flai guard 2>&1") {
			t.Errorf("the analyzer's file edits are not guarded: %q", c)
		}
		return
	}
	t.Error("the template's settings have no Edit|MultiEdit|Write|NotebookEdit hook")
}
