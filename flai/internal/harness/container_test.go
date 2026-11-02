package harness_test

import (
	"bytes"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/pilot"
)

// containerReq is a start of S-0508's agent on the container harness in a
// project at a folder of its own, with a home of its own, and the host's
// login, the token, and a provider's key set where a leak would show.
func containerReq(t *testing.T, a *manifest.Agent) (harness.Request, harness.Host) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("FLAI_MCP_TOKEN", "token-value-never-shown")
	t.Setenv("GATEWAY_KEY", "key-value-never-shown")
	root := filepath.Join(t.TempDir(), "flow")
	r := harness.Request{Story: "S-0508", Root: root, Project: "flow", Agent: a, Name: "agent-S-0508", Flai: "/usr/local/bin/flai", Session: "3f2a-b1"}
	return r, harness.Host{Login: filepath.Join(home, ".claude")}
}

func TestForSelectsTheContainerWithItsInnersCapabilitiesAndReader(t *testing.T) {
	line := []byte(`{"type":"assistant","timestamp":"2026-10-12T10:00:00Z","session_id":"s","message":{"id":"m1","model":"claude-opus-5-5",` +
		`"content":[],"usage":{"input_tokens":1,"output_tokens":1}}}`)
	want, _ := harness.Adapters[harness.ClaudeCode].Reader().Read(line)
	for _, a := range []*manifest.Agent{{Harness: harness.Container}, {Harness: harness.Container, Config: map[string]string{"inner": harness.ClaudeCode}}} {
		name, ad, err := harness.For(a, false)
		if err != nil || name != harness.Container {
			t.Fatalf("%+v: For %q %v", a, name, err)
		}
		if got := ad.Capabilities(); got != harness.Adapters[harness.ClaudeCode].Capabilities() {
			t.Errorf("%+v: capabilities %+v, want claude-code's", a, got)
		}
		if got, ok := ad.Reader().Read(line); !ok || len(got) != 1 || len(want) != 1 || got[0].MessageID != want[0].MessageID {
			t.Errorf("%+v: reader read %+v, want claude-code's %+v", a, got, want)
		}
	}
	if !slices.Contains(harness.Names(), harness.Container) {
		t.Errorf("Names %q lacks container", harness.Names())
	}
	if got := harness.Adapters[harness.Container].Capabilities(); got != harness.Adapters[harness.ClaudeCode].Capabilities() {
		t.Errorf("the plain container adapter's capabilities %+v, want claude-code's, the default inner", got)
	}
}

func TestForRefusesAnInnerHarnessAContainerCannotRun(t *testing.T) {
	for _, inner := range []string{"cursor", harness.Container, harness.Command} {
		a := &manifest.Agent{Harness: harness.Container, Config: map[string]string{"inner": inner}}
		if _, _, err := harness.For(a, true); err == nil || !strings.Contains(err.Error(), `cannot run the inner harness "`+inner+`"; set config inner to one of claude-code`) {
			t.Errorf("inner %s: %v", inner, err)
		}
		r, host := containerReq(t, a)
		if _, err := harness.Adapters[harness.Container].Start(r, host); err == nil || !strings.Contains(err.Error(), "set config inner to one of claude-code") {
			t.Errorf("inner %s: start %v", inner, err)
		}
	}
}

func TestTheContainerRunsDockerWithThreeMountsAndVariablesByName(t *testing.T) {
	a := &manifest.Agent{Harness: harness.Container, Model: "claude-opus-5-5", Provider: "gateway", Config: map[string]string{"inner": harness.ClaudeCode, "effort": "high"}}
	r, host := containerReq(t, a)
	r.ProviderName, r.Provider = "gateway", &manifest.Provider{API: manifest.APIAnthropicMessages, BaseURL: "https://gw.example", KeyEnv: "GATEWAY_KEY"}
	_, ad, err := harness.For(a, false)
	if err != nil {
		t.Fatal(err)
	}
	st, err := ad.Start(r, host)
	if err != nil {
		t.Fatal(err)
	}
	worktree := filepath.Join(r.Root, ".flai-cache", "worktrees", "S-0508")
	cache := filepath.Join(r.Root, ".flai-cache", "agents", "S-0508")
	user := strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid())
	want := []string{"docker", "run", "--rm", "-i", "--name", "flai-agent-S-0508-3f2a-b1", "--user", user,
		"--network", "bridge", "--add-host", "host.docker.internal:host-gateway",
		"-v", worktree + ":" + worktree, "-w", worktree, "-v", cache + ":" + cache, "-v", host.Login + ":/run/claude-login:ro",
		"-e", "FLAI_AGENT", "-e", "FLAI_SESSION", "-e", "FLAI_STORY", "-e", "FLAI_STARTED_BY", "-e", "FLAI_MCP_URL", "-e", "FLAI_MCP_TOKEN",
		"-e", "FLAI_CONFIG", "-e", "FLAI_CACHE_DIR", "-e", "GATEWAY_KEY", "flai-agent-claude-code:local"}
	if st.Harness != harness.Container || !slices.Equal(st.Argv, want) {
		t.Fatalf("harness %q argv\n%q\nwant\n%q", st.Harness, st.Argv, want)
	}
	home, _ := os.UserHomeDir()
	for i, arg := range st.Argv {
		if i > 0 && st.Argv[i-1] == "-e" && strings.Contains(arg, "=") {
			t.Errorf("-e %s carries a value", arg)
		}
		if i > 0 && st.Argv[i-1] == "-v" {
			src, _, _ := strings.Cut(arg, ":")
			for _, never := range []string{home, filepath.Join(home, ".flai"), r.Root, filepath.Join(r.Root, ".git")} {
				if src == never {
					t.Errorf("mounts %s", src)
				}
			}
		}
		if strings.Contains(arg, "docker.sock") {
			t.Errorf("mounts the Docker socket: %s", arg)
		}
	}
	if n := strings.Count(strings.Join(st.Argv, "\x00"), "\x00-v\x00"); n != 3 {
		t.Errorf("%d mounts, want 3", n)
	}
	whole := strings.Join(st.Argv, " ") + strings.Join(st.Env, " ") + string(st.Stdin) + strings.Join(slices.Collect(maps.Values(st.EnvFrom)), " ")
	for _, secret := range []string{"token-value-never-shown", "key-value-never-shown"} {
		if strings.Contains(whole, secret) {
			t.Errorf("the start carries a secret's value: %s", secret)
		}
	}
	if !maps.Equal(st.EnvFrom, map[string]string{"GATEWAY_KEY": "GATEWAY_KEY"}) {
		t.Errorf("EnvFrom %v, want the key's variable copied by its own name", st.EnvFrom)
	}
	if wantEnv := []string{"FLAI_CONFIG=" + cache + "/flai/config.json", "FLAI_CACHE_DIR=" + cache + "/flai/cache"}; !slices.Equal(st.Env, wantEnv) {
		t.Errorf("env %q, want %q", st.Env, wantEnv)
	}
	if len(st.Files) != 0 {
		t.Errorf("files %v, want none: pilot writes the inner harness's", st.Files)
	}
}

func TestTheContainersStandardInputIsARequestPilotRuns(t *testing.T) {
	a := &manifest.Agent{Harness: harness.Container, Model: "claude-opus-5-5", Config: map[string]string{"inner": harness.ClaudeCode, "max_budget_usd": "5"}}
	r, host := containerReq(t, a)
	r.Answered = "TH-0001"
	st, err := harness.Adapters[harness.Container].Start(r, host)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := pilot.Read(bytes.NewReader(st.Stdin))
	if err != nil {
		t.Fatalf("pilot reads %s: %v", st.Stdin, err)
	}
	cache := filepath.Join(r.Root, ".flai-cache", "agents", "S-0508")
	env := map[string]string{"FLAI_MCP_URL": "http://host.docker.internal:4243/mcp", "FLAI_MCP_TOKEN": "t"}
	if err := pilot.Check(pr, func(n string) string { return env[n] }, func(p string) bool { return p == cache }); err != nil {
		t.Fatalf("pilot refuses %s: %v", st.Stdin, err)
	}
	worktree := filepath.Join(r.Root, ".flai-cache", "worktrees", "S-0508")
	if pr.Inner != harness.ClaudeCode || pr.Cache != cache || pr.MCPURLEnv != "FLAI_MCP_URL" || pr.MCPTokenEnv != "FLAI_MCP_TOKEN" || pr.MCPHostEnv != "" ||
		pr.Root != worktree || pr.Flai != "" || pr.Story != "S-0508" || pr.Name != "agent-S-0508" || pr.Session != "3f2a-b1" || pr.Answered != "TH-0001" {
		t.Errorf("request %+v", pr)
	}
	if pr.Agent == nil || pr.Agent.Harness != harness.ClaudeCode || pr.Agent.Model != "claude-opus-5-5" || !maps.Equal(pr.Agent.Config, map[string]string{"max_budget_usd": "5"}) {
		t.Errorf("agent %+v: want the inner harness's, its config without inner", pr.Agent)
	}
	if a.Config["inner"] != harness.ClaudeCode {
		t.Errorf("the story's agent lost its inner: %v", a.Config)
	}
	if _, err := harness.Adapters[pr.Inner].Start(pr, harness.Host{}); err != nil {
		t.Errorf("the inner harness refuses the request pilot reads: %v", err)
	}
}

func TestTheContainersImageNetworkAndProgramAreTheHosts(t *testing.T) {
	r, host := containerReq(t, &manifest.Agent{Harness: harness.Container})
	host.Program, host.Network, host.Images = "/usr/bin/podman", "host", map[string]string{harness.ClaudeCode: "registry.example/agent:1"}
	st, err := harness.Adapters[harness.Container].Start(r, host)
	if err != nil {
		t.Fatal(err)
	}
	if st.Argv[0] != "/usr/bin/podman" || st.Argv[len(st.Argv)-1] != "registry.example/agent:1" || !slices.Contains(st.Argv, "host") {
		t.Errorf("argv %q", st.Argv)
	}
	if harness.DefaultImages[harness.ClaudeCode] != "flai-agent-claude-code:local" {
		t.Errorf("default images %v", harness.DefaultImages)
	}
	if got := harness.Adapters[harness.Container].DefaultHost(); got.Program != "docker" || got.Args != nil {
		t.Errorf("default host %+v", got)
	}
}

func TestTheContainerRefusesWhatItCannotRun(t *testing.T) {
	for _, c := range []struct {
		name   string
		change func(*harness.Request, *harness.Host)
		says   string
	}{
		{"no login", func(_ *harness.Request, h *harness.Host) { h.Login = "" }, "flai serve agent harness container --login <path>"},
		{"no story", func(r *harness.Request, _ *harness.Host) { r.Story = "" }, "a strategic agent's mounts are S-0509's"},
		{"the operator's arguments", func(_ *harness.Request, h *harness.Host) { h.Args = []string{"-v", "/:/host"} }, "takes no arguments of the operator's"},
		{"a login that is the home", func(_ *harness.Request, h *harness.Host) { h.Login = filepath.Dir(h.Login) }, "which a container never mounts"},
		{"a login above the project", func(r *harness.Request, h *harness.Host) { h.Login = filepath.Dir(r.Root) }, "which a container never mounts"},
		{"a login that is the root", func(_ *harness.Request, h *harness.Host) { h.Login = "/" }, "which a container never mounts"},
		{"a relative login", func(_ *harness.Request, h *harness.Host) { h.Login = ".claude" }, "not an absolute path"},
		{"a key that is no variable", func(r *harness.Request, _ *harness.Host) {
			r.ProviderName, r.Provider = "gw", &manifest.Provider{API: manifest.APIAnthropicMessages, BaseURL: "https://gw.example", KeyEnv: "K=v"}
		}, "is not a variable's name"},
		{"a config the inner harness refuses", func(r *harness.Request, _ *harness.Host) {
			r.Agent.Config = map[string]string{"permission_mode": "bypassPermissions"}
		}, `inner harness claude-code refuses the start: claude-code takes no config key "permission_mode"`},
		{"a session that names no container", func(r *harness.Request, _ *harness.Host) { r.Session = "a b" }, "not a session's name"},
	} {
		r, host := containerReq(t, &manifest.Agent{Harness: harness.Container})
		c.change(&r, &host)
		if _, err := harness.Adapters[harness.Container].Start(r, host); err == nil || !strings.Contains(err.Error(), c.says) {
			t.Errorf("%s: %v, want %q", c.name, err, c.says)
		}
	}
}

func TestTheContainersNameAndCacheAreTheStorys(t *testing.T) {
	if got := harness.ContainerName("S-0508", "abc"); got != "flai-agent-S-0508-abc" {
		t.Errorf("name %q", got)
	}
	if got := harness.ContainerName("S-0508", ""); got != "flai-agent-S-0508" {
		t.Errorf("name without a session %q", got)
	}
	if got := harness.ContainerCache("/p/flow", "S-0508"); got != "/p/flow/.flai-cache/agents/S-0508" {
		t.Errorf("cache %q", got)
	}
}
