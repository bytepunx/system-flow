package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/protected"
	"github.com/bytepunx/system-flow/flai/internal/roles"
)

// openCodeAgent is a story's agent on opencode with model.
func openCodeAgent(model string) *manifest.Agent {
	return &manifest.Agent{Harness: OpenCode, Model: model}
}

// S-0406: the opencode adapter is the harness a story names as opencode,
// under the name protected keeps OpenCode's list by.
func TestAStoryOnOpenCodeGetsTheOpenCodeAdapter(t *testing.T) {
	name, ad, err := For(openCodeAgent("anthropic/claude-opus-5-5"), false)
	if err != nil || name != OpenCode || ad != (openCode{}) {
		t.Fatalf("For(opencode) = %q, %T, %v; want opencode's adapter", name, ad, err)
	}
	if OpenCode != protected.OpenCode {
		t.Errorf("harness.OpenCode %q, protected.OpenCode %q: the protected list is kept under another name", OpenCode, protected.OpenCode)
	}
	if h := (openCode{}).DefaultHost(); h.Program != "opencode" || h.Args == nil || len(h.Args) != 0 {
		t.Errorf("default host %+v, want opencode with no arguments, --auto least of all", h)
	}
}

// S-0406: a new run is opencode run with the operator's arguments, its
// events as JSON, the story's model as provider/model, and the prompt last
// as its message; it names no session, since OpenCode makes its own, and its
// environment keeps OpenCode from updating itself and gives it flai's MCP
// server, and without deny_protected no permission.
func TestOpenCodeRunsHeadlessWithTheStorysModelAndJSON(t *testing.T) {
	r := req(openCodeAgent("anthropic/claude-opus-5-5"))
	r.Session = "6f1c2d3e-4f5a-4b6c-8d7e-9f0a1b2c3d4e"
	s, err := (openCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"opencode", "run", "--format", "json", "--model", "anthropic/claude-opus-5-5", mustPrompt(t, r)}
	if !slices.Equal(s.Argv, want) {
		t.Errorf("argv %q,\nwant %q", s.Argv, want)
	}
	content, err := openCodeConfigContent(r)
	if err != nil {
		t.Fatal(err)
	}
	if s.Harness != OpenCode || !slices.Equal(s.Env, []string{"OPENCODE_DISABLE_AUTOUPDATE=1", "OPENCODE_CONFIG_CONTENT=" + content}) || s.EnvFrom != nil || s.Files != nil || s.Stdin != nil {
		t.Errorf("start %+v, want opencode's with OPENCODE_DISABLE_AUTOUPDATE=1 and flai's MCP server alone in its environment, no OPENCODE_PERMISSION, and no file", s)
	}

	s, err = (openCode{}).Start(r, Host{Program: "/opt/bin/opencode", Args: []string{"--agent", "build"}})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"/opt/bin/opencode", "run", "--agent", "build", "--format", "json", "--model", "anthropic/claude-opus-5-5", mustPrompt(t, r)}
	if !slices.Equal(s.Argv, want) {
		t.Errorf("the operator's host: argv %q,\nwant %q", s.Argv, want)
	}

	// OpenCode splits at the first /, so the model ID may hold more.
	r.Agent = openCodeAgent("openrouter/anthropic/claude-opus-5-5")
	if s, err = (openCode{}).Start(r, Host{}); err != nil || after(s.Argv, "--model") != "openrouter/anthropic/claude-opus-5-5" {
		t.Errorf("a model ID with a /: %q, %v", s.Argv, err)
	}
	// No model leaves OpenCode's own.
	r.Agent = &manifest.Agent{Harness: OpenCode}
	if s, err = (openCode{}).Start(r, Host{}); err != nil || slices.Contains(s.Argv, "--model") {
		t.Errorf("no model: %q, %v; want no --model", s.Argv, err)
	}
}

// S-0406: a run that goes on in its session, after an answer or a send back,
// passes the ses_ ID OpenCode made before its message; a run that does not
// go on passes none, whatever session it names.
func TestOpenCodeResumesTheSessionItIsGiven(t *testing.T) {
	const ses = "ses_0b1c2d3e4f5aAbCdEfGhIjKlMn"
	answered := req(openCodeAgent("anthropic/claude-opus-5-5"))
	answered.Session, answered.Answered = ses, "TH-0001"
	sentBack := req(openCodeAgent("anthropic/claude-opus-5-5"))
	sentBack.Session, sentBack.SentBack = ses, &SentBack{Thread: "TH-0002", What: "check failed", Resumed: true}
	for _, r := range []Request{answered, sentBack} {
		s, err := (openCode{}).Start(r, Host{})
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"opencode", "run", "--format", "json", "--model", "anthropic/claude-opus-5-5", "--session", ses, mustPrompt(t, r)}
		if !slices.Equal(s.Argv, want) {
			t.Errorf("argv %q,\nwant %q", s.Argv, want)
		}
	}

	fresh := req(openCodeAgent("anthropic/claude-opus-5-5"))
	fresh.Session = ses
	if s, err := (openCode{}).Start(fresh, Host{}); err != nil || slices.Contains(s.Argv, "--session") {
		t.Errorf("a new run: %q, %v; want no --session", s.Argv, err)
	}
}

// S-0406: a resumed run whose session is not one OpenCode made is refused,
// since opencode run --session would not find it, saying how to start anew.
func TestOpenCodeRefusesToResumeASessionItDidNotMake(t *testing.T) {
	for _, session := range []string{"6f1c2d3e-4f5a-4b6c-8d7e-9f0a1b2c3d4e", "ses_", "--auto"} {
		r := req(openCodeAgent("anthropic/claude-opus-5-5"))
		r.Session, r.Answered = session, "TH-0001"
		_, err := (openCode{}).Start(r, Host{})
		if err == nil || !strings.Contains(err.Error(), "by the ses_ ID OpenCode made for it") || !strings.Contains(err.Error(), "flai serve agent restart S-0104") {
			t.Errorf("session %q: %v", session, err)
		}
	}
}

// S-0406: OpenCode splits the model at its first /, so a model without a
// provider before it, or a model after it, is refused, saying the form.
func TestOpenCodeRefusesAModelItCannotSplit(t *testing.T) {
	for _, model := range []string{"claude-opus-5-5", "anthropic/"} {
		_, err := (openCode{}).Start(req(openCodeAgent(model)), Host{})
		if err == nil || !strings.Contains(err.Error(), "opencode takes the model as provider/model") || !strings.Contains(err.Error(), fmt.Sprintf("the agent's model is %q; name its provider and model so, such as anthropic/claude-opus-5-5", model)) {
			t.Errorf("model %q: %v", model, err)
		}
	}
}

// S-0406: what opencode cannot do yet is refused at start, each saying what
// to do: a config key other than effort, a strategic role, a provider
// (S-0410), and a start in a container.
func TestOpenCodeRefusesWhatItCannotDoYet(t *testing.T) {
	config := req(&manifest.Agent{Harness: OpenCode, Model: "anthropic/claude-opus-5-5", Config: map[string]string{"max_budget_usd": "5"}})
	plan := planReq("E-0001", openCodeAgent("anthropic/claude-opus-5-5"))
	defined := Request{Root: "/p/flow", Project: "flow", Agent: openCodeAgent("anthropic/claude-opus-5-5"), Name: "triage", Role: "triage",
		RoleFile: &roles.Role{Name: "triage", Role: "triage", Path: "design/workflows/roles/triage.md"}}
	named := req(&manifest.Agent{Harness: OpenCode, Model: "anthropic/claude-opus-5-5", Provider: "litellm"})
	resolved := withProvider(req(openCodeAgent("anthropic/claude-opus-5-5")))
	contained := req(openCodeAgent("anthropic/claude-opus-5-5"))
	contained.MCPURLEnv, contained.MCPTokenEnv, contained.Cache = "FLAI_MCP_URL", "FLAI_MCP_TOKEN", "/cache"
	for _, c := range []struct {
		what string
		r    Request
		says string
	}{
		{"config", config, `opencode takes no config key "max_budget_usd"; it takes effort`},
		{"planner", plan, "opencode cannot run the strategic role plan: flai runs a strategic agent from its definition, and has none for opencode yet; run the role on claude-code"},
		{"defined role", defined, "opencode cannot run the strategic role triage"},
		{"provider named", named, "opencode cannot call the agent's provider litellm yet"},
		{"provider resolved", resolved, "leave provider out of the story's agent and name the model as provider/model, or run the story on claude-code"},
		{"container", contained, "opencode cannot run in a container yet"},
	} {
		if s, err := (openCode{}).Start(c.r, Host{}); err == nil || !strings.Contains(err.Error(), c.says) {
			t.Errorf("%s: %q, %v; want a refusal saying %q", c.what, s.Argv, err, c.says)
		}
	}
}

// S-0407: a story whose agent names roles starts on opencode only in a
// project that keeps the flai-guard plugin, which guards its sub-agents; the
// refusal names the plugin and flai upgrade. A story without roles needs none.
func TestOpenCodeRunsRolesOnlyWhereTheGuardPluginIs(t *testing.T) {
	root := t.TempDir()
	agent := openCodeAgent("anthropic/claude-opus-5-5")
	agent.Roles = map[string]manifest.Role{manifest.RoleExplore: {Model: "anthropic/claude-haiku-5-5"}}
	r := req(agent)
	r.Root = root
	_, err := (openCode{}).Start(r, Host{})
	if err == nil || !strings.Contains(err.Error(), OpenCodeGuardPlugin) || !strings.Contains(err.Error(), "flai upgrade") {
		t.Fatalf("roles without the plugin: %v; want refused naming %s and flai upgrade", err, OpenCodeGuardPlugin)
	}
	r.Agent = openCodeAgent("anthropic/claude-opus-5-5")
	if _, err := (openCode{}).Start(r, Host{}); err != nil {
		t.Errorf("no roles, no plugin: %v; want started", err)
	}
	plugin := filepath.Join(root, filepath.FromSlash(OpenCodeGuardPlugin))
	if err := os.MkdirAll(filepath.Dir(plugin), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(plugin, []byte("export const FlaiGuard = async () => ({})\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r.Agent = agent
	if _, err := (openCode{}).Start(r, Host{}); err != nil {
		t.Errorf("roles with the plugin: %v; want started", err)
	}
}
