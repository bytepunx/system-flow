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
		"leave the whole suite, the lint, and flai check to the verifier rather than running them yourself as well",
		"have one fresh verifier run the whole suite, the lint, and flai check in the worktree, through the project's close-out script where it has one",
		"Fix what it finds yourself, never through a sub-agent, commit, and have one more fresh verifier run the same and confirm the fixes",
		"A verifier's passing run is the story's run before review: do not repeat it",
	} {
		if !strings.Contains(p, w) {
			t.Errorf("prompt lacks %q:\n%s", w, p)
		}
	}
}

// S-0176: the story's agent plans its tasks in layers as it writes them,
// records why in Decisions, and hands each task to a task sub-agent, a layer
// at once when its tasks are long, whose work it reviews and commits itself;
// an answered or commit run is not told again.
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
			"A task sub-agent edits only what its task touches, runs only its own tests, and never commits, syncs, or writes through flai",
			"Name the task's ID in each task sub-agent's description, so that flai measures the task by its calls",
			"Review each one's work yourself, fix what falls short, commit it, sync and test as above, and move the task",
			"only you commit, sync the stream, move items, and talk to the designer",
		} {
			if !strings.Contains(p, w) {
				t.Errorf("prompt lacks %q:\n%s", w, p)
			}
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
			"then run the tests for what the task changed and commit any fix they need",
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
			"End with a one-line summary of what you changed, on which items: flai serve logs the run's activity",
		}
	}
	for _, c := range []struct {
		id, flag string
		want     []string
	}{
		{"E-0016", "epic", []string{"If it has no stories, draft the stories that deliver its outcome", "create each as a draft in the backlog", "If it has stories, revisit each one not done or cancelled against the epic's outcome"}},
		{"S-0208", "story", []string{"It is a story: enrich it with its predicted touches, a forecast, and a cost of delay value worked out from the operator's inputs"}},
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
	if p := Prompt(planReq("S-0208", nil)); strings.Contains(p, "draft the stories") {
		t.Errorf("a story's planner is told to draft an epic's stories:\n%s", p)
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
	if want := []string{"FLAI_ROLE=plan", "FLAI_ITEM=E-0016"}; !slices.Equal(st.Env, want) {
		t.Errorf("env %q, want %q", st.Env, want)
	}
	// flai guard judges the planner's file edits, which the project's hook
	// does not match
	var settings struct {
		Hooks map[string][]struct {
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Type    string `json:"type"`
				Command string `json:"command"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal([]byte(after(st.Argv, "--settings")), &settings); err != nil {
		t.Fatalf("--settings: %v in %v", err, st.Argv)
	}
	if pre := settings.Hooks["PreToolUse"]; len(pre) != 1 || pre[0].Matcher != "Edit|Write|NotebookEdit" || len(pre[0].Hooks) != 1 || pre[0].Hooks[0].Type != "command" ||
		pre[0].Hooks[0].Command != `out=$('/usr/local/bin/flai' guard 2>&1); [ $? -eq 2 ] || exit 0; echo "$out" >&2; exit 2` {
		t.Errorf("settings: %+v", settings)
	}
	if got := shQuote("/opt/it's/flai"); got != `'/opt/it'\''s/flai'` {
		t.Errorf("quoted %s", got)
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
		`flai cannot start an agent in role "orchestrate"`: func(r *Request) { r.Role = "orchestrate" },
		"was given story S-0104 to work":                   func(r *Request) { r.Story = "S-0104" },
		`"T-0784" is neither`:                              func(r *Request) { r.Item = "T-0784" },
		`"" is neither`:                                    func(r *Request) { r.Item = "" },
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

// The template's own definitions read as --agents takes them.
func TestTheTemplatesDefinitionsRead(t *testing.T) {
	for _, def := range []string{"explorer", "verifier", "planner"} {
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
