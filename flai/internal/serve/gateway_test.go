//go:build !windows

package serve

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The gateway check (S-0351): a headless claude -p, started with the
// claude-code adapter's own arguments and a provider entry for OpenRouter,
// in a scratch project with flai guard registered and auto-approve on. It
// runs only where OPENROUTER_API_KEY is set and claude is on PATH, and
// spends real money on that key, capped by max_budget_usd.

const (
	// gatewayKeyEnv holds the OpenRouter key the check runs on.
	gatewayKeyEnv = "OPENROUTER_API_KEY"
	// gatewayBaseURL is OpenRouter's Anthropic-compatible API, as a provider's base_url.
	gatewayBaseURL = "https://openrouter.ai/api"
	// gatewayHaikuEnv and gatewayHaikuDefault are the slug the haiku alias resolves to.
	gatewayHaikuEnv     = "FLAI_TEST_OPENROUTER_HAIKU"
	gatewayHaikuDefault = "anthropic/claude-haiku-4.5"
	// gatewayBudgetEnv and gatewayBudgetDefault are the run's max_budget_usd,
	// which Claude Code holds to its own count: it knows no price for an
	// OpenRouter slug and counts three to four times what the gateway
	// charges, so the default lets a run the gateway charges about $0.06 for
	// finish.
	gatewayBudgetEnv     = "FLAI_TEST_GATEWAY_BUDGET"
	gatewayBudgetDefault = "1.00"
	// gatewayAgent is who the scratch session works as.
	gatewayAgent = "flai-gateway-check"
	// gatewayTitle is the title the sub-agent is asked to give the story.
	gatewayTitle = "Edited by a sub-agent"
	// gatewayEditTool is the work-item write the sub-agent is asked to make.
	gatewayEditTool = "mcp__flai__item_edit"
	// gatewayLookups bounds the wait for OpenRouter's generation records.
	gatewayLookups = 30 * time.Second
	// gatewayKeyLag bounds the wait for the key's usage to count the run:
	// OpenRouter updates it some time after the calls.
	gatewayKeyLag = 2 * time.Minute
	// keySettled is how long the key's usage stays put before the run is
	// taken as wholly counted in it.
	keySettled = 30 * time.Second
	// keyReads is how many reads of the key one reading of its spend takes.
	keyReads = 5
)

// gatewayMissing names what the check needs and this host lacks, "" when
// nothing: the key's variable, claude on PATH, or both.
func gatewayMissing(key string, claudeErr error) string {
	var missing []string
	if key == "" {
		missing = append(missing, gatewayKeyEnv+" is not set")
	}
	if claudeErr != nil {
		missing = append(missing, "claude is not on PATH")
	}
	return strings.Join(missing, " and ")
}

// gatewayUsage is a message's usage as the stream reports it.
type gatewayUsage struct {
	Input         int `json:"input_tokens"`
	Output        int `json:"output_tokens"`
	CacheCreation int `json:"cache_creation_input_tokens"`
	CacheRead     int `json:"cache_read_input_tokens"`
	ServerToolUse *struct {
		WebSearch int `json:"web_search_requests"`
		WebFetch  int `json:"web_fetch_requests"`
	} `json:"server_tool_use"`
}

// gatewayLine is the part of a stream-json line the check reads.
type gatewayLine struct {
	Type       string   `json:"type"`
	Subtype    string   `json:"subtype"`
	Model      string   `json:"model"`
	Version    string   `json:"claude_code_version"`
	Tools      []string `json:"tools"`
	MCPServers []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"mcp_servers"`
	Betas           []string `json:"betas"`
	ParentToolUseID string   `json:"parent_tool_use_id"`
	Message         *struct {
		ID      string          `json:"id"`
		Model   string          `json:"model"`
		Content json.RawMessage `json:"content"`
		Usage   *gatewayUsage   `json:"usage"`
	} `json:"message"`
	IsError           bool            `json:"is_error"`
	NumTurns          int             `json:"num_turns"`
	DurationMS        int64           `json:"duration_ms"`
	DurationAPIMS     int64           `json:"duration_api_ms"`
	TotalCostUSD      float64         `json:"total_cost_usd"`
	Result            string          `json:"result"`
	Usage             json.RawMessage `json:"usage"`
	ModelUsage        json.RawMessage `json:"modelUsage"`
	PermissionDenials json.RawMessage `json:"permission_denials"`
}

// gatewayBlock is one block of a message's content.
type gatewayBlock struct {
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
	ToolUseID string          `json:"tool_use_id"`
	Content   json.RawMessage `json:"content"`
	IsError   bool            `json:"is_error"`
}

// gatewayMessage is one model response: its ID, model, the tool use it was
// made under (empty for the main session's), and its last usage.
type gatewayMessage struct {
	ID     string
	Model  string
	Parent string
	Usage  gatewayUsage
}

// gatewayCall is one tool call and its answer.
type gatewayCall struct {
	ID       string
	Name     string
	Input    json.RawMessage
	Parent   string // the tool use of the sub-agent that made it; empty for the main session's
	Answered bool
	Result   string
	IsError  bool
}

// gatewayRun is what a run's stream says.
type gatewayRun struct {
	Init     *gatewayLine
	Messages []gatewayMessage // in the order first seen, one per message ID
	Calls    []*gatewayCall   // in the order made
	Result   *gatewayLine
	Other    []string // lines that are not stream-json
}

// parseGatewayStream reads a run's output: stream-json lines and any others.
func parseGatewayStream(out string) gatewayRun {
	var run gatewayRun
	seen := map[string]int{}
	calls := map[string]*gatewayCall{}
	for line := range strings.Lines(out) {
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}
		var l gatewayLine
		if text[0] != '{' || json.Unmarshal([]byte(text), &l) != nil || l.Type == "" {
			run.Other = append(run.Other, text)
			continue
		}
		switch {
		case l.Type == "system" && l.Subtype == "init" && run.Init == nil:
			run.Init = &l
		case l.Type == "result":
			run.Result = &l
		case (l.Type == "assistant" || l.Type == "user") && l.Message != nil:
			if l.Type == "assistant" && l.Message.ID != "" {
				m := gatewayMessage{ID: l.Message.ID, Model: l.Message.Model, Parent: l.ParentToolUseID}
				if l.Message.Usage != nil {
					m.Usage = *l.Message.Usage
				}
				if i, ok := seen[m.ID]; ok {
					run.Messages[i] = m
				} else {
					seen[m.ID] = len(run.Messages)
					run.Messages = append(run.Messages, m)
				}
			}
			var blocks []gatewayBlock
			if json.Unmarshal(l.Message.Content, &blocks) != nil {
				continue
			}
			for _, b := range blocks {
				switch b.Type {
				case "tool_use":
					c := &gatewayCall{ID: b.ID, Name: b.Name, Input: b.Input, Parent: l.ParentToolUseID}
					calls[b.ID] = c
					run.Calls = append(run.Calls, c)
				case "tool_result":
					if c, ok := calls[b.ToolUseID]; ok {
						c.Answered, c.Result, c.IsError = true, resultText(b.Content), b.IsError
					}
				}
			}
		}
	}
	return run
}

// guardRefused is a sub-agent's call of gatewayEditTool that flai guard
// refused, and the refusal Claude Code handed back; nil when there is none.
func (r gatewayRun) guardRefused() *gatewayCall {
	short := strings.TrimPrefix(gatewayEditTool, "mcp__flai__")
	for _, c := range r.Calls {
		if c.Parent != "" && c.Name == gatewayEditTool && c.IsError && strings.Contains(c.Result, "a sub-agent (") && strings.Contains(c.Result, "cannot call "+short) {
			return c
		}
	}
	return nil
}

// wrote is the main session's call of tool on file, nil when it made none.
func (r gatewayRun) wrote(tool, file string) *gatewayCall {
	for _, c := range r.Calls {
		var in struct {
			FilePath string `json:"file_path"`
		}
		if c.Parent == "" && c.Name == tool && json.Unmarshal(c.Input, &in) == nil && in.FilePath == file {
			return c
		}
	}
	return nil
}

// resolvedTo says whether the run shows the model named slug answering:
// the session's model, a message's, or a model it reports usage for.
func (r gatewayRun) resolvedTo(slug string) bool {
	if r.Init != nil && r.Init.Model == slug {
		return true
	}
	for _, m := range r.Messages {
		if m.Model == slug {
			return true
		}
	}
	if r.Result != nil {
		var usage map[string]json.RawMessage
		if json.Unmarshal(r.Result.ModelUsage, &usage) == nil {
			if _, ok := usage[slug]; ok {
				return true
			}
		}
	}
	return false
}

// TestGatewayStreamIsJudged covers the check's reading of a stream with no
// network: a sub-agent's refused work-item write, the main session's write,
// the model the alias resolved to, the messages by ID, and the end.
func TestGatewayStreamIsJudged(t *testing.T) {
	stream := strings.Join([]string{
		`{"type":"system","subtype":"init","model":"anthropic/claude-haiku-4.5","tools":["Write","Agent","mcp__flai__item_edit"],"mcp_servers":[{"name":"flai","status":"connected"}]}`,
		`{"type":"assistant","parent_tool_use_id":null,"message":{"id":"gen-1","model":"anthropic/claude-haiku-4.5","content":[{"type":"tool_use","id":"tu1","name":"Write","input":{"file_path":"/w/.claude/x.md","content":"x"}}],"usage":{"input_tokens":3,"output_tokens":1}}}`,
		`{"type":"assistant","parent_tool_use_id":null,"message":{"id":"gen-1","model":"anthropic/claude-haiku-4.5","content":[{"type":"tool_use","id":"tu2","name":"Agent","input":{}}],"usage":{"input_tokens":3,"output_tokens":9,"cache_read_input_tokens":4}}}`,
		`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"tu1","content":"File created"}]}}`,
		`{"type":"assistant","parent_tool_use_id":"tu2","message":{"id":"gen-2","model":"anthropic/claude-haiku-4.5","content":[{"type":"tool_use","id":"tu3","name":"mcp__flai__item_edit","input":{"id":"S-0001"}}]}}`,
		`{"type":"user","parent_tool_use_id":"tu2","message":{"content":[{"type":"tool_result","tool_use_id":"tu3","is_error":true,"content":[{"type":"text","text":"PreToolUse:mcp__flai__item_edit hook error: a sub-agent (general-purpose) cannot call item_edit: it reads"}]}]}}`,
		`stand-in: not stream-json`,
		`{"type":"result","subtype":"success","is_error":false,"num_turns":3,"total_cost_usd":0.01,"modelUsage":{"anthropic/claude-haiku-4.5":{"inputTokens":6}}}`,
	}, "\n")
	run := parseGatewayStream(stream)

	if run.Init == nil || run.Result == nil || run.Result.Subtype != "success" || run.Result.NumTurns != 3 {
		t.Fatalf("init %+v, result %+v; want both, the result a success of 3 turns", run.Init, run.Result)
	}
	if len(run.Messages) != 2 || run.Messages[0].ID != "gen-1" || run.Messages[0].Usage.Output != 9 || run.Messages[0].Usage.CacheRead != 4 || run.Messages[1].Parent != "tu2" {
		t.Errorf("messages %+v, want gen-1 with its last usage and gen-2 under tu2", run.Messages)
	}
	if c := run.guardRefused(); c == nil || c.ID != "tu3" {
		t.Errorf("the sub-agent's refused %s is not found: %+v", gatewayEditTool, run.Calls)
	}
	if c := run.wrote("Write", "/w/.claude/x.md"); c == nil || !c.Answered || c.IsError {
		t.Errorf("the main session's Write is not found answered: %+v", c)
	}
	if !run.resolvedTo("anthropic/claude-haiku-4.5") || run.resolvedTo("claude-haiku-4-5") {
		t.Error("the model the alias resolved to is not told apart")
	}
	if len(run.Other) != 1 {
		t.Errorf("other lines %q, want the one that is not stream-json", run.Other)
	}

	// the main session's refused call, or a sub-agent's that passed, is not the guard's refusal
	mainCall := strings.ReplaceAll(stream, `"parent_tool_use_id":"tu2","message":{"id":"gen-2"`, `"parent_tool_use_id":null,"message":{"id":"gen-2"`)
	if c := parseGatewayStream(mainCall).guardRefused(); c != nil {
		t.Errorf("the main session's refused call counts as a sub-agent's: %+v", c)
	}
	passed := strings.ReplaceAll(stream, `"is_error":true,`, ``)
	if c := parseGatewayStream(passed).guardRefused(); c != nil {
		t.Errorf("a call that passed counts as refused: %+v", c)
	}
}

func TestGatewayMissingNamesWhatIsMissing(t *testing.T) {
	notFound := exec.ErrNotFound
	for _, c := range []struct {
		key  string
		err  error
		want string
	}{
		{"k", nil, ""},
		{"", nil, "OPENROUTER_API_KEY is not set"},
		{"k", notFound, "claude is not on PATH"},
		{"", notFound, "OPENROUTER_API_KEY is not set and claude is not on PATH"},
	} {
		if got := gatewayMissing(c.key, c.err); got != c.want {
			t.Errorf("gatewayMissing(%q, %v) = %q, want %q", c.key, c.err, got, c.want)
		}
	}
}

// S-0351: a claude-code agent runs through OpenRouter with the adapter's own
// arguments and a provider entry for it, flai guard registered, and
// auto-approve on: its protected Write goes through permission_prompt, its
// sub-agent's work-item write is refused by the guard, the haiku alias
// resolves through the provider's models, and the run ends with a result.
// It logs what the gateway charged beside what Claude Code counted.
func TestGatewayOpenRouter(t *testing.T) {
	if testing.Short() {
		t.Skip("runs claude -p through OpenRouter, a paid call")
	}
	key := os.Getenv(gatewayKeyEnv)
	claude, claudeErr := exec.LookPath("claude")
	if missing := gatewayMissing(key, claudeErr); missing != "" {
		t.Skip("the gateway check needs OpenRouter's key and claude: " + missing)
	}
	slug := envOr(gatewayHaikuEnv, gatewayHaikuDefault)
	budget := envOr(gatewayBudgetEnv, gatewayBudgetDefault)
	ctx := context.Background()

	version, err := exec.CommandContext(ctx, claude, "--version").Output()
	if err != nil {
		t.Fatalf("claude --version: %v", err)
	}
	t.Logf("claude --version: %s (%s)", strings.TrimSpace(string(version)), claude)

	flai := buildFlai(t)
	settings := guardSettings(t, flai)
	s, err := newClaudeScratch(execx.System{}, time.Now(), map[string]string{".claude/settings.json": settings})
	if err != nil {
		t.Fatalf("the scratch project could not be made: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(s.dir) })
	repo, err := workitem.Open(s.root)
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Get(s.story)
	if err != nil {
		t.Fatal(err)
	}
	title := story.Title

	provider := &manifest.Provider{API: manifest.APIAnthropicMessages, BaseURL: gatewayBaseURL, KeyEnv: gatewayKeyEnv, Models: map[string]string{"haiku": slug}}
	req := harness.Request{
		Story: s.story, Root: s.root, Project: "flai-check", Name: gatewayAgent, Flai: flai,
		Agent:        &manifest.Agent{Harness: harness.ClaudeCode, Model: "haiku", Provider: "openrouter", Config: map[string]string{"max_budget_usd": budget}},
		ProviderName: "openrouter", Provider: provider,
	}
	spec, err := harness.Adapters[harness.ClaudeCode].Start(req, harness.ClaudeCodeHost(harness.Host{}))
	if err != nil {
		t.Fatalf("the claude-code adapter refused the start: %v", err)
	}
	keys, err := providerKeys(req, spec, os.LookupEnv)
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("flai checked Claude Code %s through OpenRouter and permission_prompt.", strings.Fields(string(version))[0])
	prompt := gatewayPrompt(s.target, want, s.story)
	argv := slices.Clone(spec.Argv)
	i := slices.Index(argv, "-p")
	if i < 0 || i+1 >= len(argv) {
		t.Fatalf("the adapter's arguments have no -p: %q", spec.Argv)
	}
	argv[i+1] = prompt
	shown := slices.Clone(argv)
	shown[i+1] = "<the test's prompt>"
	t.Logf("argv: %q", shown)
	t.Logf("prompt: %s", prompt)

	env, dropped := withoutClaudeSession(agentEnv(os.Environ()))
	added := append(slices.Clone(spec.Env), "FLAI_STORY="+s.story, "FLAI_AGENT="+gatewayAgent, "FLAI_STARTED_BY=flai-serve", "FLAI_CONFIG="+s.config)
	t.Logf("env set: %q", added)
	t.Logf("env set from the host, names only: %q", envNames(keys))
	t.Logf("env dropped from this process's, names only (the session go test may run under): %q", dropped)
	env = append(append(env, added...), keys...)

	before, beforeErr := keyNow(ctx, key)
	if beforeErr != nil {
		t.Logf("key before: unread: %v", beforeErr)
	}
	started := time.Now()
	out, runErr := runClaudeWith(ctx, s.root, argv, env)
	t.Logf("claude ran %s; err: %v", time.Since(started).Round(time.Second), runErr)
	run := parseGatewayStream(out)
	logGatewayRun(t, run)
	lookups := lookUpGenerations(ctx, t, key, run.Messages)
	after, afterErr := keyAfter(ctx, key, before, beforeErr)
	if afterErr != nil {
		t.Logf("key after: unread: %v", afterErr)
	}
	if beforeErr == nil && afterErr == nil {
		t.Logf("key usage before $%.6f, after $%.6f, delta $%.6f (daily before $%.6f, after $%.6f); limit %s, remaining %s", before.Usage, after.Usage, after.Usage-before.Usage, before.UsageDaily, after.UsageDaily, dollars(after.Limit), dollars(after.LimitRemaining))
	}
	cost := -1.0
	if run.Result != nil {
		cost = run.Result.TotalCostUSD
	}
	t.Logf("gateway's cost: sum of generations $%.6f (%d of %d found), key delta %s; Claude Code's total_cost_usd $%.6f", lookups.total, lookups.found, len(run.Messages), keyDelta(before, after, beforeErr, afterErr), cost)

	// the protected Write, through permission_prompt
	if got, err := os.ReadFile(s.target); err != nil || strings.TrimSpace(string(got)) != want {
		t.Errorf("%s holds %q (%v), want %q; Claude Code said:\n%s", s.target, got, err, want, claudeSaid(out))
	}
	if c := run.wrote("Write", s.target); c == nil || !c.Answered || c.IsError {
		t.Errorf("the main session's Write of %s did not pass: %+v", s.target, c)
	}
	// the sub-agent's work-item write, refused by flai guard
	if c := run.guardRefused(); c == nil {
		var tried []string
		for _, c := range run.Calls {
			tried = append(tried, fmt.Sprintf("%s (parent %q, error %v): %s", c.Name, c.Parent, c.IsError, cut(c.Result, 300)))
		}
		t.Errorf("no sub-agent's %s was refused by flai guard; the calls:\n%s", gatewayEditTool, strings.Join(tried, "\n"))
	}
	if edited, err := repo.Get(s.story); err != nil {
		t.Errorf("the story could not be read after the run: %v", err)
	} else if edited.Title != title {
		t.Errorf("the story's title is %q, want %q unchanged: a work-item write got through", edited.Title, title)
	}
	// the alias, resolved through the provider's models
	if !run.resolvedTo(slug) {
		t.Errorf("no model in the stream is %s, the slug haiku resolves to", slug)
	}
	// the end
	if run.Result == nil {
		t.Errorf("the run ended with no result; Claude Code said:\n%s", claudeSaid(out))
	}
}

// envOr is the variable's value, or def when it is unset or empty.
func envOr(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return def
}

// buildFlai builds this tree's flai into a temporary folder.
func buildFlai(t *testing.T) string {
	t.Helper()
	exe := filepath.Join(t.TempDir(), "flai")
	build := exec.Command("go", "build", "-o", exe, "github.com/bytepunx/system-flow/flai")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return exe
}

// guardSettings is the template's .claude/settings.json with flai guard run
// from flai, this tree's build, rather than the installed flai.
func guardSettings(t *testing.T, flai string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "template", "root", ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("the template's settings: %v", err)
	}
	if strings.ContainsAny(flai, `'"\`) {
		t.Fatalf("the built flai's path %s cannot be quoted in the hooks", flai)
	}
	settings := string(data)
	if n := strings.Count(settings, "flai guard"); n != 4 {
		t.Fatalf("the template's settings run flai guard %d times, want 4 (PreToolUse twice, SubagentStart, SubagentStop):\n%s", n, settings)
	}
	return strings.ReplaceAll(settings, "flai guard", "'"+flai+"' guard")
}

// gatewayPrompt asks for the protected Write and for a sub-agent's
// work-item write.
func gatewayPrompt(target, content, story string) string {
	return fmt.Sprintf(`This is flai's check of Claude Code through a gateway. Do these two steps, in order, and nothing else.

1. Use the Write tool once to create the file %[1]s with exactly this content, on one line: %[2]s
Do not read the file first.

2. Use the Agent tool (named Task in older versions) once, with subagent_type general-purpose and model haiku, and give the sub-agent exactly this prompt: "Call the flai MCP tool %[4]s once, with id %[3]s and title %[5]q. If that tool is not loaded yet, load it with ToolSearch first. Do not retry it and use no other tool. Reply with the exact text the call returned, word for word."

Do not call any flai MCP tool yourself, and do not retry either step. When both are done, reply with one line saying what the Write returned and what the sub-agent reported.`, target, content, story, gatewayEditTool, gatewayTitle)
}

// withoutClaudeSession is env without the variables a Claude Code session
// gives the processes it runs, which a check run from one would otherwise
// pass on, as flai serve's agents never have them; and their names.
func withoutClaudeSession(env []string) (kept, dropped []string) {
	for _, kv := range env {
		name, _, _ := strings.Cut(kv, "=")
		if name == "CLAUDECODE" || name == "CLAUDE_PID" || name == "CLAUDE_EFFORT" || strings.HasPrefix(name, "CLAUDE_CODE_") {
			dropped = append(dropped, name)
			continue
		}
		kept = append(kept, kv)
	}
	return kept, dropped
}

// envNames are the names of env's variables.
func envNames(env []string) []string {
	out := make([]string, 0, len(env))
	for _, kv := range env {
		name, _, _ := strings.Cut(kv, "=")
		out = append(out, name)
	}
	return out
}

// logGatewayRun logs what the stream says, for the story's write-up.
func logGatewayRun(t *testing.T, run gatewayRun) {
	t.Helper()
	if in := run.Init; in != nil {
		var flaiTools []string
		for _, tool := range in.Tools {
			if strings.HasPrefix(tool, "mcp__flai__") {
				flaiTools = append(flaiTools, tool)
			}
		}
		t.Logf("init: model %s, claude_code_version %s, %d tools (%d of them flai's MCP tools; %s present: %v), mcp_servers %+v, betas %q", in.Model, in.Version, len(in.Tools), len(flaiTools), gatewayEditTool, slices.Contains(in.Tools, gatewayEditTool), in.MCPServers, in.Betas)
		t.Logf("init tools: %q", in.Tools)
	} else {
		t.Log("init: none in the stream")
	}
	var cacheSeen, serverTools bool
	for _, m := range run.Messages {
		who := "main"
		if m.Parent != "" {
			who = "sub-agent under " + m.Parent
		}
		u := m.Usage
		t.Logf("message %s (%s): model %s, input %d, output %d, cache_creation_input_tokens %d, cache_read_input_tokens %d", m.ID, who, m.Model, u.Input, u.Output, u.CacheCreation, u.CacheRead)
		cacheSeen = cacheSeen || u.CacheCreation > 0 || u.CacheRead > 0
		if u.ServerToolUse != nil && (u.ServerToolUse.WebSearch > 0 || u.ServerToolUse.WebFetch > 0) {
			serverTools = true
		}
	}
	t.Logf("prompt caching reached the model (cache tokens nonzero in a message): %v; server tools used: %v", cacheSeen, serverTools)
	for _, c := range run.Calls {
		who := "main"
		if c.Parent != "" {
			who = "sub-agent under " + c.Parent
		}
		t.Logf("call %s %s (%s): answered %v, error %v: %s", c.ID, c.Name, who, c.Answered, c.IsError, cut(oneLine(c.Result), 600))
	}
	if r := run.Result; r != nil {
		t.Logf("result: subtype %s, is_error %v, num_turns %d, duration_ms %d, duration_api_ms %d, total_cost_usd %.6f", r.Subtype, r.IsError, r.NumTurns, r.DurationMS, r.DurationAPIMS, r.TotalCostUSD)
		t.Logf("result usage: %s", compact(r.Usage))
		t.Logf("result modelUsage: %s", compact(r.ModelUsage))
		t.Logf("result permission_denials: %s", compact(r.PermissionDenials))
		t.Logf("result text: %s", cut(oneLine(r.Result), 600))
	} else {
		t.Log("result: none in the stream")
	}
	for _, o := range run.Other {
		t.Logf("not stream-json: %s", cut(o, 600))
	}
}

// compact is raw JSON on one line, or "none".
func compact(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "none"
	}
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return string(raw)
	}
	out, _ := json.Marshal(v)
	return string(out)
}

// openRouterKeyInfo is what OpenRouter's /api/v1/key says of the key's spend.
type openRouterKeyInfo struct {
	Usage          float64  `json:"usage"`
	UsageDaily     float64  `json:"usage_daily"`
	Limit          *float64 `json:"limit"`
	LimitRemaining *float64 `json:"limit_remaining"`
}

// openRouterKey reads the key's spend.
func openRouterKey(ctx context.Context, key string) (openRouterKeyInfo, error) {
	var body struct {
		Data openRouterKeyInfo `json:"data"`
	}
	_, err := openRouterGet(ctx, key, "/v1/key", &body)
	return body.Data, err
}

// keyNow is the key's spend as the highest of keyReads reads a second
// apart: OpenRouter answers /v1/key from copies that lag one another, and a
// key's usage only grows, so the highest is the newest.
func keyNow(ctx context.Context, key string) (openRouterKeyInfo, error) {
	var best openRouterKeyInfo
	var err error
	for i := range keyReads {
		if i > 0 {
			time.Sleep(time.Second)
		}
		info, rerr := openRouterKey(ctx, key)
		if rerr != nil {
			err = rerr
			continue
		}
		if info.Usage >= best.Usage {
			best = info
		}
		err = nil
	}
	return best, err
}

// keyAfter reads the key's spend once the run has settled in it: OpenRouter
// counts the calls some time after they are made, one by one, so it waits
// for the usage to move from before's and then to stay put for keySettled,
// at most gatewayKeyLag in all; it answers the highest usage read, or the
// last read when the usage never moved.
func keyAfter(ctx context.Context, key string, before openRouterKeyInfo, beforeErr error) (openRouterKeyInfo, error) {
	deadline := time.Now().Add(gatewayKeyLag)
	last, moved := before, time.Now()
	for {
		after, err := keyNow(ctx, key)
		switch {
		case beforeErr != nil:
			return after, err
		case err == nil && after.Usage > last.Usage:
			last, moved = after, time.Now()
		}
		if last.Usage > before.Usage && (time.Since(moved) >= keySettled || time.Now().After(deadline)) {
			return last, nil
		}
		if time.Now().After(deadline) {
			return after, err
		}
		time.Sleep(2 * time.Second)
	}
}

// dollars is an amount the gateway may leave out.
func dollars(v *float64) string {
	if v == nil {
		return "none"
	}
	return fmt.Sprintf("$%.6f", *v)
}

// keyDelta is the key's spend over the run, or why it is not known.
func keyDelta(before, after openRouterKeyInfo, beforeErr, afterErr error) string {
	if beforeErr != nil || afterErr != nil {
		return "unknown"
	}
	return fmt.Sprintf("$%.6f", after.Usage-before.Usage)
}

// generationTotals are what the generation lookups found.
type generationTotals struct {
	found int
	total float64
}

// lookUpGenerations reads OpenRouter's record of each message by its ID,
// waiting a little for one not yet recorded, and logs each; one never found
// is logged, not failed.
func lookUpGenerations(ctx context.Context, t *testing.T, key string, messages []gatewayMessage) generationTotals {
	t.Helper()
	var out generationTotals
	deadline := time.Now().Add(gatewayLookups)
	for _, m := range messages {
		for {
			var body struct {
				Data map[string]any `json:"data"`
			}
			code, err := openRouterGet(ctx, key, "/v1/generation?id="+url.QueryEscape(m.ID), &body)
			if err == nil && body.Data != nil {
				out.found++
				if cost, ok := body.Data["total_cost"].(float64); ok {
					out.total += cost
				}
				data, _ := json.Marshal(body.Data)
				t.Logf("generation %s: found: %s", m.ID, data)
				break
			}
			// only a generation not yet recorded is worth waiting for
			if code != http.StatusNotFound || time.Now().After(deadline) {
				t.Logf("generation %s: not found: %v", m.ID, err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}
	return out
}

// openRouterGet GETs path under the gateway's API with the key and decodes
// the body into v; it answers the status, and an error naming it and what
// the gateway said when it is not 200.
func openRouterGet(ctx context.Context, key, path string, v any) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayBaseURL+path, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("GET %s%s: %w", gatewayBaseURL, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		said, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return resp.StatusCode, fmt.Errorf("GET %s%s answered %d: %s", gatewayBaseURL, path, resp.StatusCode, oneLine(string(said)))
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return resp.StatusCode, fmt.Errorf("GET %s%s: %w", gatewayBaseURL, path, err)
	}
	return resp.StatusCode, nil
}
