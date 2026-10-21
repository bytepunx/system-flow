//go:build !windows

package serve

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
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

// The gateway checks (S-0351, S-0356): a headless claude -p, started with
// the claude-code adapter's own arguments and a provider entry for a
// gateway, in a scratch project with flai guard registered and auto-approve
// on. The OpenRouter check runs only where OPENROUTER_API_KEY is set and
// claude is on PATH; the LiteLLM check, by the proxy's unified API and, when
// FLAI_TEST_LITELLM_PASSTHROUGH=1, by its Anthropic pass-through, only where
// LITELLM_BASE_URL and LITELLM_API_KEY are set and claude is on PATH. Each
// spends on its key, capped by max_budget_usd, but the LiteLLM check's
// default model is a free one.

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

	// litellmBaseEnv and litellmKeyEnv hold the LiteLLM proxy's root URL and
	// the virtual key the check runs on.
	litellmBaseEnv = "LITELLM_BASE_URL"
	litellmKeyEnv  = "LITELLM_API_KEY"
	// litellmMasterEnv holds the proxy's master key, the only key that reads
	// its spend log.
	litellmMasterEnv = "LITELLM_MASTER_KEY"
	// litellmHaikuEnv and litellmHaikuDefault are the slug the haiku alias
	// resolves to by the unified API: by default a free OpenRouter model, so
	// the check spends nothing (TH-0429). OpenRouter retired the free
	// nemotron-3-super, answering 404, so the default is nemotron-3-ultra
	// (I-0170, TH-0529).
	litellmHaikuEnv     = "FLAI_TEST_LITELLM_HAIKU"
	litellmHaikuDefault = "openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free"
	// litellmPassThroughEnv set to 1 runs the pass-through too: it forwards
	// the body to Anthropic as is, so it needs the proxy's upstream Anthropic
	// key and takes Anthropic's own model ID, litellmPassThroughHaiku.
	litellmPassThroughEnv   = "FLAI_TEST_LITELLM_PASSTHROUGH"
	litellmPassThroughHaiku = "claude-haiku-4-5"
	// litellmSpendLag bounds the wait for the key's spend to count the run:
	// the proxy writes spend in batches.
	litellmSpendLag = 2 * time.Minute
	// litellmSettled is how long the key's spend stays put before the run is
	// taken as wholly counted in it.
	litellmSettled = 15 * time.Second
	// litellmLogLag bounds the wait for the run's rows in the spend log.
	litellmLogLag = time.Minute
	// litellmLogSettled is how long the count of the run's rows in the spend
	// log stays put before they are taken as all there.
	litellmLogSettled = 6 * time.Second
)

// gatewayMissing names what the check needs and this host lacks, "" when
// nothing: each of the variables names that getenv finds empty, and claude
// on PATH.
func gatewayMissing(getenv func(string) string, claudeErr error, names ...string) string {
	var missing []string
	for _, name := range names {
		if getenv(name) == "" {
			missing = append(missing, name+" is not set")
		}
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
	SessionID  string   `json:"session_id"`
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

// totalCost is the run's total_cost_usd, or -1 when it ended with no result.
func (r gatewayRun) totalCost() float64 {
	if r.Result == nil {
		return -1
	}
	return r.Result.TotalCostUSD
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

// TestGatewayMissingNamesWhatIsMissing covers both checks' skips: each
// variable that is not set, claude not on PATH, and nothing missing.
func TestGatewayMissingNamesWhatIsMissing(t *testing.T) {
	notFound := exec.ErrNotFound
	openRouter := []string{gatewayKeyEnv}
	litellm := []string{litellmBaseEnv, litellmKeyEnv}
	for _, c := range []struct {
		env   map[string]string
		names []string
		err   error
		want  string
	}{
		{map[string]string{gatewayKeyEnv: "k"}, openRouter, nil, ""},
		{nil, openRouter, nil, "OPENROUTER_API_KEY is not set"},
		{map[string]string{gatewayKeyEnv: "k"}, openRouter, notFound, "claude is not on PATH"},
		{nil, openRouter, notFound, "OPENROUTER_API_KEY is not set and claude is not on PATH"},
		{map[string]string{litellmBaseEnv: "u", litellmKeyEnv: "k"}, litellm, nil, ""},
		{map[string]string{litellmBaseEnv: "u"}, litellm, notFound, "LITELLM_API_KEY is not set and claude is not on PATH"},
		{map[string]string{litellmKeyEnv: "k"}, litellm, nil, "LITELLM_BASE_URL is not set"},
		{nil, litellm, notFound, "LITELLM_BASE_URL is not set and LITELLM_API_KEY is not set and claude is not on PATH"},
	} {
		getenv := func(name string) string { return c.env[name] }
		if got := gatewayMissing(getenv, c.err, c.names...); got != c.want {
			t.Errorf("gatewayMissing(%v, %v, %q) = %q, want %q", c.env, c.err, c.names, got, c.want)
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
	if missing := gatewayMissing(os.Getenv, claudeErr, gatewayKeyEnv); missing != "" {
		t.Skip("the gateway check needs OpenRouter's key and claude: " + missing)
	}
	slug := envOr(gatewayHaikuEnv, gatewayHaikuDefault)
	ctx := context.Background()
	var before openRouterKeyInfo
	var beforeErr error
	checkGateway(t, claude, gatewayRoute{
		name: "OpenRouter", providerName: "openrouter", slug: slug,
		provider: &manifest.Provider{API: manifest.APIAnthropicMessages, BaseURL: gatewayBaseURL, KeyEnv: gatewayKeyEnv, Models: map[string]string{"haiku": slug}},
		before: func() {
			before, beforeErr = keyNow(ctx, key)
			if beforeErr != nil {
				t.Logf("key before: unread: %v", beforeErr)
			}
		},
		after: func(run gatewayRun, _ time.Time) {
			lookups := lookUpGenerations(ctx, t, key, run.Messages)
			after, afterErr := keyAfter(ctx, key, before, beforeErr)
			if afterErr != nil {
				t.Logf("key after: unread: %v", afterErr)
			}
			if beforeErr == nil && afterErr == nil {
				t.Logf("key usage before $%.6f, after $%.6f, delta $%.6f (daily before $%.6f, after $%.6f); limit %s, remaining %s", before.Usage, after.Usage, after.Usage-before.Usage, before.UsageDaily, after.UsageDaily, dollars(after.Limit), dollars(after.LimitRemaining))
			}
			t.Logf("gateway's cost: sum of generations $%.6f (%d of %d found), key delta %s; Claude Code's total_cost_usd $%.6f", lookups.total, lookups.found, len(run.Messages), keyDelta(before, after, beforeErr, afterErr), run.totalCost())
		},
	})
}

// S-0356: the same check through a LiteLLM proxy: by its unified API, which
// translates and maps the slug through the proxy's model list, and, when
// FLAI_TEST_LITELLM_PASSTHROUGH=1, by its Anthropic pass-through, which
// forwards the body to Anthropic as is and so takes Anthropic's own model
// ID and the proxy's upstream Anthropic key. After each run it logs the
// key's spend and, with the master key, the proxy's spend log rows for the
// run and whether they carry the run's session ID.
func TestGatewayLiteLLM(t *testing.T) {
	if testing.Short() {
		t.Skip("runs claude -p through a LiteLLM proxy, a call out of this host")
	}
	claude, claudeErr := exec.LookPath("claude")
	if missing := gatewayMissing(os.Getenv, claudeErr, litellmBaseEnv, litellmKeyEnv); missing != "" {
		t.Skip("the gateway check needs a LiteLLM proxy, its key, and claude: " + missing)
	}
	root := strings.TrimRight(os.Getenv(litellmBaseEnv), "/")
	key := os.Getenv(litellmKeyEnv)
	master := os.Getenv(litellmMasterEnv)
	slug := envOr(litellmHaikuEnv, litellmHaikuDefault)
	for _, c := range []struct {
		name, route, baseURL, slug string
		optIn, probe               bool
	}{
		{"unified", "LiteLLM (unified)", root, slug, false, true},
		{"pass-through", "LiteLLM (Anthropic pass-through)", root + "/anthropic", litellmPassThroughHaiku, true, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			if c.optIn && os.Getenv(litellmPassThroughEnv) != "1" {
				t.Skip("the pass-through reaches only Anthropic, a paid upstream: set " + litellmPassThroughEnv + "=1 where the proxy has an Anthropic key")
			}
			ctx := context.Background()
			// a model the proxy cannot reach fails claude with only a warning
			// it prints on passing runs too, so ask the proxy first (I-0171)
			if c.probe {
				if err := litellmProbe(ctx, root, key, c.slug); err != nil {
					t.Fatal(litellmProbeFailure(c.slug, litellmHaikuSource(os.Getenv), err))
				}
			}
			var before litellmKeyInfo
			var beforeErr error
			checkGateway(t, claude, gatewayRoute{
				name: c.route, providerName: "litellm", slug: c.slug,
				provider: &manifest.Provider{API: manifest.APIAnthropicMessages, BaseURL: c.baseURL, KeyEnv: litellmKeyEnv, Models: map[string]string{"haiku": c.slug}},
				before: func() {
					before, beforeErr = litellmKey(ctx, root, key)
					if beforeErr != nil {
						t.Logf("key before: unread: %v", beforeErr)
					}
				},
				after: func(run gatewayRun, started time.Time) {
					logLiteLLMSpend(ctx, t, root, key, master, run, started, before, beforeErr)
				},
			})
		})
	}
}

// gatewayRoute is one way through a gateway that the check runs claude by.
type gatewayRoute struct {
	name         string // as the prompt names it, e.g. "OpenRouter"
	providerName string
	provider     *manifest.Provider
	slug         string                                  // the model haiku resolves to
	before       func()                                  // reads the gateway's spend just before claude runs
	after        func(run gatewayRun, started time.Time) // reads it once the run is logged, before the judgements
}

// checkGateway runs claude, the binary at that path, through route in a
// scratch project with flai guard registered and auto-approve on, logs what
// the stream says and what route reads of the spend, and judges the run:
// the protected Write through permission_prompt, the sub-agent's work-item
// write refused by the guard, the haiku alias resolved to route.slug, and a
// result.
func checkGateway(t *testing.T, claude string, route gatewayRoute) {
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

	req := harness.Request{
		Story: s.story, Root: s.root, Project: "flai-check", Name: gatewayAgent, Flai: flai,
		Agent:        &manifest.Agent{Harness: harness.ClaudeCode, Model: "haiku", Provider: route.providerName, Config: map[string]string{"max_budget_usd": budget}},
		ProviderName: route.providerName, Provider: route.provider,
	}
	spec, err := harness.Adapters[harness.ClaudeCode].Start(req, harness.ClaudeCodeHost(harness.Host{}))
	if err != nil {
		t.Fatalf("the claude-code adapter refused the start: %v", err)
	}
	keys, err := providerKeys(req, spec, os.LookupEnv)
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("flai checked Claude Code %s through %s and permission_prompt", strings.Fields(string(version))[0], route.name)
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

	route.before()
	started := time.Now()
	out, runErr := runClaudeWith(ctx, s.root, argv, env)
	t.Logf("claude ran %s; err: %v", time.Since(started).Round(time.Second), runErr)
	run := parseGatewayStream(out)
	logGatewayRun(t, run)
	route.after(run, started)

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
	if !run.resolvedTo(route.slug) {
		t.Errorf("no model in the stream is %s, the slug haiku resolves to", route.slug)
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

2. Use the Agent tool (named Task in older versions) once, in the foreground with run_in_background false, with subagent_type general-purpose and model haiku, and give the sub-agent exactly this prompt: "Call the flai MCP tool %[4]s once, with id %[3]s and title %[5]q. If that tool is not loaded yet, load it with ToolSearch first. Do not retry it and use no other tool. Reply with the exact text the call returned, word for word."

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

// TestLiteLLMSpendLogIsJudged covers the check's reading of the proxy's
// spend log with no network: the rows since the run started in each time
// layout, the rows it cannot read, their metadata's keys, and where the
// run's session ID appears in them.
func TestLiteLLMSpendLogIsJudged(t *testing.T) {
	since := time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)
	raw := []json.RawMessage{
		json.RawMessage(`{"startTime":"2026-10-12T08:59:59.999999Z","model":"claude-haiku-4-5","spend":1,"session_id":"s-1"}`),
		json.RawMessage(`{"startTime":"2026-10-12T09:00:00","model":"claude-haiku-4-5","call_type":"anthropic_messages","spend":0.25,"prompt_tokens":10,"completion_tokens":2,"session_id":"s-1","metadata":{"user_api_key_alias":"k","headers":{"x-claude-code-session-id":"s-1"}}}`),
		json.RawMessage(`{"startTime":"2026-10-12 09:00:05.5+00:00","model":"anthropic/claude-haiku-4-5","spend":0.5,"session_id":null,"metadata":"{\"tags\":[\"x\"]}","messages":[{"note":"session s-1"}]}`),
		json.RawMessage(`{"startTime":"yesterday","spend":2}`),
		json.RawMessage(`[]`),
	}
	rows, unread := spendRowsSince(raw, since)
	if len(rows) != 2 || unread != 2 || rows[0].Spend != 0.25 || rows[0].PromptTokens != 10 || rows[1].Model != "anthropic/claude-haiku-4-5" {
		t.Fatalf("rows %+v, %d unread; want the two since %s and two unread", rows, unread, since)
	}
	if got := jsonKeys(rows[0].Metadata); !slices.Equal(got, []string{"headers", "user_api_key_alias"}) {
		t.Errorf("the first row's metadata keys are %q, want headers and user_api_key_alias", got)
	}
	if got := jsonKeys(rows[1].Metadata); !slices.Equal(got, []string{"tags"}) {
		t.Errorf("the second row's metadata keys, from a string, are %q, want tags", got)
	}
	field, where := sessionJoin(rows, "s-1")
	want := []string{"row 1: metadata.headers.x-claude-code-session-id", "row 1: session_id", "row 2: messages[0].note"}
	if !field || !slices.Equal(where, want) {
		t.Errorf("sessionJoin(s-1) = %v, %q; want true, %q", field, where, want)
	}
	if field, where := sessionJoin(rows, "s-2"); field || where != nil {
		t.Errorf("sessionJoin(s-2) = %v, %q; want no join", field, where)
	}
	if field, where := sessionJoin(rows, ""); field || where != nil {
		t.Errorf("sessionJoin of no session ID = %v, %q; want no join", field, where)
	}
}

// TestLiteLLMProbeNamesTheModel covers the probe of the unified case's model
// against a stand-in proxy, with no call out of this host (I-0171): a model
// OpenRouter retired, which the proxy answers 404 for, fails naming the
// slug, where it came from, the status, the upstream's message, and how to
// name another; a model it serves passes; a proxy that is down fails.
func TestLiteLLMProbeNamesTheModel(t *testing.T) {
	const slug = "openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free"
	retired := `{"error":{"message":"litellm.NotFoundError: AnthropicException - {\"error\":{\"message\":\"No endpoints found for nvidia/nemotron-3-super-120b-a12b:free.\",\"code\":404}}","type":null,"param":null,"code":"404"}}`
	type request struct {
		method, path, auth, version string
		model                       string
		maxTokens                   int
	}
	status, answer := http.StatusNotFound, retired
	sent := make(chan request, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Model     string `json:"model"`
			MaxTokens int    `json:"max_tokens"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		sent <- request{r.Method, r.URL.Path, r.Header.Get("Authorization"), r.Header.Get("anthropic-version"), body.Model, body.MaxTokens}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, answer)
	}))
	defer srv.Close()
	ctx := context.Background()

	err := litellmProbe(ctx, srv.URL, "sk-virtual", slug)
	if err == nil {
		t.Fatal("the probe of a retired model passed, want the proxy's 404")
	}
	if got, want := <-sent, (request{http.MethodPost, "/v1/messages", "Bearer sk-virtual", "2023-06-01", slug, 1}); got != want {
		t.Errorf("the probe sent %+v, want %+v", got, want)
	}
	set := func(string) string { return slug }
	unset := func(string) string { return "" }
	for _, c := range []struct {
		getenv func(string) string
		source string
	}{
		{set, litellmHaikuEnv + ", set in the environment or by .flai-cache/litellm/litellm.env"},
		{unset, "litellmHaikuDefault"},
	} {
		msg := litellmProbeFailure(slug, litellmHaikuSource(c.getenv), err)
		for _, want := range []string{slug, c.source, "answered 404", "No endpoints found for nvidia/nemotron-3-super-120b-a12b:free.", "LITELLM_HAIKU_MODEL=<slug> scripts/litellm.sh up", "set " + litellmHaikuEnv} {
			if !strings.Contains(msg, want) {
				t.Errorf("the probe's failure does not name %q:\n%s", want, msg)
			}
		}
	}

	status, answer = http.StatusOK, `{"type":"message","content":[]}`
	if err := litellmProbe(ctx, srv.URL, "sk-virtual", slug); err != nil {
		t.Errorf("the probe of a model the proxy serves failed: %v", err)
	}
	<-sent

	srv.Close()
	if err := litellmProbe(ctx, srv.URL, "sk-virtual", slug); err == nil || !strings.Contains(err.Error(), "POST "+srv.URL+"/v1/messages") {
		t.Errorf("the probe of a proxy that is down answered %v, want an error naming the POST", err)
	}
}

// logLiteLLMSpend logs the proxy's reading of a run's spend: the key's
// spend before and after, the spend log's rows since started when the master
// key is at hand and whether they carry the run's session ID, and the cost
// beside Claude Code's total_cost_usd.
func logLiteLLMSpend(ctx context.Context, t *testing.T, root, key, master string, run gatewayRun, started time.Time, before litellmKeyInfo, beforeErr error) {
	t.Helper()
	after, afterErr := litellmKeyAfter(ctx, root, key, before, beforeErr)
	if afterErr != nil {
		t.Logf("key after: unread: %v", afterErr)
	}
	delta := "unknown"
	if beforeErr == nil && afterErr == nil {
		delta = fmt.Sprintf("$%.6f", after.Spend-before.Spend)
		t.Logf("key spend before $%.6f, after $%.6f, delta %s; max_budget %s", before.Spend, after.Spend, delta, dollars(after.MaxBudget))
	}
	session := ""
	if run.Init != nil {
		session = run.Init.SessionID
	}
	t.Logf("the run's session ID (session_id on system/init): %q", session)
	logged := "unknown (not read)"
	if master == "" {
		t.Logf("spend log: not read: only the proxy's master key reads /spend/logs and %s is not set", litellmMasterEnv)
	} else if rows, ok := readSpendLog(ctx, t, root, master, key, started); ok {
		total := 0.0
		for i, r := range rows {
			total += r.Spend
			t.Logf("spend log row %d: startTime %s, model %q, model_group %q, custom_llm_provider %q, call_type %q, spend $%.6f, prompt_tokens %d, completion_tokens %d, session_id %q, request_tags %s, metadata keys %q", i+1, r.StartTime, r.Model, r.ModelGroup, r.Provider, r.CallType, r.Spend, r.PromptTokens, r.CompletionTokens, r.SessionID, compact(r.RequestTags), jsonKeys(r.Metadata))
		}
		logged = fmt.Sprintf("$%.6f (%d rows)", total, len(rows))
		field, where := sessionJoin(rows, session)
		t.Logf("session join: a row's session_id is the run's: %v; the run's session ID appears in the rows at: %q", field, where)
	}
	t.Logf("gateway's cost: key delta %s, spend log rows %s; Claude Code's total_cost_usd $%.6f", delta, logged, run.totalCost())
}

// litellmKeyInfo is what the proxy's /key/info says of the key's spend.
type litellmKeyInfo struct {
	Spend     float64  `json:"spend"`
	MaxBudget *float64 `json:"max_budget"`
}

// litellmKey reads the key's spend from the proxy at root.
func litellmKey(ctx context.Context, root, key string) (litellmKeyInfo, error) {
	var body struct {
		Info litellmKeyInfo `json:"info"`
	}
	_, err := litellmGet(ctx, root, key, "/key/info", nil, &body)
	return body.Info, err
}

// litellmProbe sends one message of one token for slug by the unified API
// of the proxy at root with key; it answers nil when the proxy answers 2xx,
// and otherwise an error naming the status and what the proxy said.
func litellmProbe(ctx context.Context, root, key, slug string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	target := root + "/v1/messages"
	body, err := json.Marshal(map[string]any{
		"model":      slug,
		"max_tokens": 1,
		"messages":   []map[string]string{{"role": "user", "content": "ping"}},
	})
	if err != nil {
		return fmt.Errorf("POST %s: %w", target, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("POST %s: %w", target, err)
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST %s: %w", target, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	said := litellmErrorMessage(data)
	if key != "" {
		said = strings.ReplaceAll(said, key, "<redacted>")
	}
	return fmt.Errorf("POST %s answered %d: %s", target, resp.StatusCode, said)
}

// litellmErrorMessage is the message of the proxy's error body,
// {"error":{"message":...}}, on one line, or the body on one line when it
// holds none.
func litellmErrorMessage(body []byte) string {
	var e struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &e) == nil && e.Error.Message != "" {
		return oneLine(e.Error.Message)
	}
	return oneLine(string(body))
}

// litellmHaikuSource says where the unified case's slug came from:
// litellmHaikuEnv, which the environment or the proxy's key file sets, the
// check cannot tell which, or the built-in default when it is unset.
func litellmHaikuSource(getenv func(string) string) string {
	if getenv(litellmHaikuEnv) != "" {
		return litellmHaikuEnv + ", set in the environment or by .flai-cache/litellm/litellm.env (scripts/env.sh sources it)"
	}
	return "the built-in default, litellmHaikuDefault in gateway_test.go, as " + litellmHaikuEnv + " is not set"
}

// litellmProbeFailure is what the check says when the probe of slug, which
// came from source, failed with err: what failed, on what, and how to name a
// model the proxy serves.
func litellmProbeFailure(slug, source string, err error) string {
	return fmt.Sprintf("the haiku alias's model %q, from %s, does not answer through the LiteLLM proxy, so claude -p would fail with no more than its unrecognized_model warning: %v; "+
		"a 404 is an upstream that does not serve the slug, as when OpenRouter retires a free model (I-0170): name one it serves with LITELLM_HAIKU_MODEL=<slug> scripts/litellm.sh up, which writes it to .flai-cache/litellm/litellm.env, or set %s for this run",
		slug, source, err, litellmHaikuEnv)
}

// litellmKeyAfter reads the key's spend once the run has settled in it: the
// proxy writes spend in batches, so it waits for the spend to move from
// before's and then to stay put for litellmSettled, at most litellmSpendLag
// in all; it answers the highest spend read, or the last read when the spend
// never moved.
func litellmKeyAfter(ctx context.Context, root, key string, before litellmKeyInfo, beforeErr error) (litellmKeyInfo, error) {
	deadline := time.Now().Add(litellmSpendLag)
	last, moved := before, time.Now()
	for {
		after, err := litellmKey(ctx, root, key)
		switch {
		case beforeErr != nil:
			return after, err
		case err == nil && after.Spend > last.Spend:
			last, moved = after, time.Now()
		}
		if last.Spend > before.Spend && (time.Since(moved) >= litellmSettled || time.Now().After(deadline)) {
			return last, nil
		}
		if time.Now().After(deadline) {
			return after, err
		}
		time.Sleep(2 * time.Second)
	}
}

// readSpendLog reads the proxy's spend log, by the master key, for the key's
// rows that started at or after since, waiting up to litellmLogLag for rows
// to appear and their count to stay put for litellmLogSettled; it logs what
// it read, a log it could not read, and no row found, and answers whether it
// read the log at all.
func readSpendLog(ctx context.Context, t *testing.T, root, master, key string, since time.Time) ([]litellmRow, bool) {
	t.Helper()
	// the proxy files a key's rows under its hex SHA-256
	sum := sha256.Sum256([]byte(key))
	query := url.Values{"api_key": {hex.EncodeToString(sum[:])}}
	deadline := time.Now().Add(litellmLogLag)
	var rows []litellmRow
	var err error
	read, all, unread, count, changed := false, 0, 0, -1, time.Now()
	for {
		var raw []json.RawMessage
		var code int
		if code, err = litellmGet(ctx, root, master, "/spend/logs", query, &raw); err == nil {
			read, all = true, len(raw)
			rows, unread = spendRowsSince(raw, since)
			if len(rows) != count {
				count, changed = len(rows), time.Now()
			}
			if count > 0 && time.Since(changed) >= litellmLogSettled {
				break
			}
		}
		// a refusal, unlike a row not yet written, does not pass with waiting
		if code >= 400 && code < 500 || time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		t.Logf("spend log: unread: %v", err)
	}
	if !read {
		return nil, false
	}
	t.Logf("spend log: %d rows for the key, %d of them since the run started, %d unreadable", all, len(rows), unread)
	if len(rows) == 0 {
		t.Log("spend log: no row for the run")
	}
	return rows, true
}

// litellmRow is one row of the proxy's spend log, as far as the check reads it.
type litellmRow struct {
	StartTime        string          `json:"startTime"`
	Model            string          `json:"model"`
	ModelGroup       string          `json:"model_group"`
	Provider         string          `json:"custom_llm_provider"`
	CallType         string          `json:"call_type"`
	Spend            float64         `json:"spend"`
	PromptTokens     int             `json:"prompt_tokens"`
	CompletionTokens int             `json:"completion_tokens"`
	SessionID        string          `json:"session_id"`
	RequestTags      json.RawMessage `json:"request_tags"`
	Metadata         json.RawMessage `json:"metadata"`
	whole            any             // the whole row, for finding the session ID in it
}

// spendRowsSince is the rows of raw that started at or after since, and how
// many rows it could not read.
func spendRowsSince(raw []json.RawMessage, since time.Time) (rows []litellmRow, unread int) {
	for _, r := range raw {
		var row litellmRow
		if json.Unmarshal(r, &row) != nil || json.Unmarshal(r, &row.whole) != nil {
			unread++
			continue
		}
		start, err := litellmTime(row.StartTime)
		if err != nil {
			unread++
			continue
		}
		if !start.Before(since) {
			rows = append(rows, row)
		}
	}
	return rows, unread
}

// litellmTime reads a spend log row's time: RFC 3339, or the same with a
// space for the T, or with no zone, which is taken as UTC.
func litellmTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.999999999", "2006-01-02 15:04:05.999999999Z07:00", "2006-01-02 15:04:05.999999999"} {
		if at, err := time.Parse(layout, s); err == nil {
			return at, nil
		}
	}
	return time.Time{}, fmt.Errorf("the spend log's time %q is in no layout it uses", s)
}

// jsonKeys is the sorted names of raw's keys when it is an object or a
// string holding one, nil otherwise.
func jsonKeys(raw json.RawMessage) []string {
	var obj map[string]json.RawMessage
	if json.Unmarshal(raw, &obj) != nil {
		var s string
		if json.Unmarshal(raw, &s) != nil || json.Unmarshal([]byte(s), &obj) != nil {
			return nil
		}
	}
	return slices.Sorted(maps.Keys(obj))
}

// sessionJoin says whether a row's session_id is id, and where in the rows
// id appears at all, as "row N: path" counting from 1; neither for no id.
func sessionJoin(rows []litellmRow, id string) (field bool, where []string) {
	if id == "" {
		return false, nil
	}
	for i, r := range rows {
		field = field || r.SessionID == id
		for _, p := range pathsHolding(r.whole, "", id) {
			where = append(where, fmt.Sprintf("row %d: %s", i+1, p))
		}
	}
	return field, where
}

// pathsHolding is the paths below path of the strings in v that hold s, in
// key order: keys joined by dots, list items by their index in brackets.
func pathsHolding(v any, path, s string) []string {
	var out []string
	switch v := v.(type) {
	case string:
		if strings.Contains(v, s) {
			out = append(out, path)
		}
	case map[string]any:
		for _, k := range slices.Sorted(maps.Keys(v)) {
			p := k
			if path != "" {
				p = path + "." + k
			}
			out = append(out, pathsHolding(v[k], p, s)...)
		}
	case []any:
		for i, item := range v {
			out = append(out, pathsHolding(item, fmt.Sprintf("%s[%d]", path, i), s)...)
		}
	}
	return out
}

// litellmGet GETs path, with query, under the proxy at root with bearer as
// the key, and decodes the body into v; it answers the status, and an error
// naming it and what the proxy said when it is not 200. Its errors name the
// path without the query and leave out the key and the query's values.
func litellmGet(ctx context.Context, root, bearer, path string, query url.Values, v any) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	target := root + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	// a *url.Error names the whole URL, query and all; its Err does not
	bare := func(err error) error {
		var ue *url.Error
		if errors.As(err, &ue) {
			return ue.Err
		}
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return 0, fmt.Errorf("GET %s%s: %w", root, path, bare(err))
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("GET %s%s: %w", root, path, bare(err))
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		said := string(data)
		secrets := []string{bearer}
		for _, vs := range query {
			secrets = append(secrets, vs...)
		}
		for _, secret := range secrets {
			if secret != "" {
				said = strings.ReplaceAll(said, secret, "<redacted>")
			}
		}
		return resp.StatusCode, fmt.Errorf("GET %s%s answered %d: %s", root, path, resp.StatusCode, oneLine(said))
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return resp.StatusCode, fmt.Errorf("GET %s%s: %w", root, path, err)
	}
	return resp.StatusCode, nil
}
