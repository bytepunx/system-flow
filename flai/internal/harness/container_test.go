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
	"github.com/bytepunx/system-flow/flai/internal/roles"
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
	linkedWorktree(t, root, "S-0508")
	r := harness.Request{Story: "S-0508", Root: root, Project: "flow", Agent: a, Name: "agent-S-0508", Flai: "/usr/local/bin/flai", Session: "3f2a-b1"}
	return r, harness.Host{Login: filepath.Join(home, ".claude")}
}

// linkedWorktree makes story's worktree under the project at root as git
// worktree add leaves it: a .git file naming the story's registration in the
// main checkout's .git, whose commondir names that .git, beside the main
// checkout's wip and system-flow.yaml, and returns the registration.
func linkedWorktree(t *testing.T, root, story string) string {
	t.Helper()
	reg := filepath.Join(root, ".git", "worktrees", story)
	worktree := filepath.Join(root, ".flai-cache", "worktrees", story)
	for _, d := range []string{reg, worktree, filepath.Join(root, "wip")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: flow\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reg, "commondir"), []byte("../..\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: "+reg+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return reg
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

// The container mounts eight paths: the worktree; the main checkout's .git
// read-only and the story's registration under it read-write, after it, so
// git works in the worktree and can write nothing of the repository's but
// the story's own index, HEAD, and logs (ADR-0246, I-0215); the main
// checkout's wip and system-flow.yaml read-only and its .flai-cache/verify
// read-write, made on the host when absent, so flai verify runs in the
// container against the main checkout's state and records its report where
// the host reads it (ADR-0253, I-0218); the cache; and the login read-only.
// The mounts are an allow-list: nothing else of the main checkout, nothing
// under .flai-cache but the story's worktree, its agent folder, and verify,
// and no Docker socket.
func TestTheContainerRunsDockerWithEightMountsFiveOfTheMainCheckoutsAndVariablesByName(t *testing.T) {
	a := &manifest.Agent{Harness: harness.Container, Model: "claude-opus-5-5", Provider: "gateway", Config: map[string]string{"inner": harness.ClaudeCode, "effort": "high"}}
	r, host := containerReq(t, a)
	r.ProviderName, r.Provider = "gateway", &manifest.Provider{API: manifest.APIAnthropicMessages, BaseURL: "https://gw.example", KeyEnv: "GATEWAY_KEY"}
	_, ad, err := harness.For(a, false)
	if err != nil {
		t.Fatal(err)
	}
	verify := filepath.Join(r.Root, ".flai-cache", "verify")
	if _, err := os.Stat(verify); !os.IsNotExist(err) {
		t.Fatalf("the verify folder should be absent before the start: %v", err)
	}
	st, err := ad.Start(r, host)
	if err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(verify); err != nil || !info.IsDir() || info.Mode().Perm() != 0o755 {
		t.Errorf("the start should make the verify folder on the host, 0755, before docker makes it as root: %v %v", info, err)
	}
	worktree := filepath.Join(r.Root, ".flai-cache", "worktrees", "S-0508")
	cache := filepath.Join(r.Root, ".flai-cache", "agents", "S-0508")
	common, reg := filepath.Join(r.Root, ".git"), filepath.Join(r.Root, ".git", "worktrees", "S-0508")
	wip, flow := filepath.Join(r.Root, "wip"), filepath.Join(r.Root, "system-flow.yaml")
	user := strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid())
	want := []string{"docker", "run", "--rm", "-i", "--name", "flai-agent-S-0508-3f2a-b1", "--user", user,
		"--network", "bridge", "--add-host", "host.docker.internal:host-gateway",
		"-v", worktree + ":" + worktree, "-w", worktree, "-v", common + ":" + common + ":ro", "-v", reg + ":" + reg,
		"-v", wip + ":" + wip + ":ro", "-v", flow + ":" + flow + ":ro", "-v", verify + ":" + verify,
		"-v", cache + ":" + cache, "-v", host.Login + ":/run/claude-login:ro",
		"-e", "FLAI_AGENT", "-e", "FLAI_SESSION", "-e", "FLAI_STORY", "-e", "FLAI_STARTED_BY", "-e", "FLAI_MCP_URL", "-e", "FLAI_MCP_TOKEN",
		"-e", "FLAI_CONFIG", "-e", "FLAI_CACHE_DIR", "-e", "GIT_CONFIG_COUNT", "-e", "GIT_CONFIG_KEY_0", "-e", "GIT_CONFIG_VALUE_0",
		"-e", "GATEWAY_KEY", "flai-agent-claude-code:local"}
	if st.Harness != harness.Container || !slices.Equal(st.Argv, want) {
		t.Fatalf("harness %q argv\n%q\nwant\n%q", st.Harness, st.Argv, want)
	}
	allowed := map[string]string{
		worktree:   worktree + ":" + worktree,
		common:     common + ":" + common + ":ro",
		reg:        reg + ":" + reg,
		wip:        wip + ":" + wip + ":ro",
		flow:       flow + ":" + flow + ":ro",
		verify:     verify + ":" + verify,
		cache:      cache + ":" + cache,
		host.Login: host.Login + ":/run/claude-login:ro",
	}
	mounts := 0
	for i, arg := range st.Argv {
		if i > 0 && st.Argv[i-1] == "-e" && strings.Contains(arg, "=") {
			t.Errorf("-e %s carries a value", arg)
		}
		if i > 0 && st.Argv[i-1] == "-v" {
			mounts++
			src, _, _ := strings.Cut(arg, ":")
			if allowed[src] != arg {
				t.Errorf("mounts %s, which is not one of the eight the container's mounts allow, in its mode", arg)
			}
		}
		if strings.Contains(arg, "docker.sock") {
			t.Errorf("mounts the Docker socket: %s", arg)
		}
	}
	if mounts != len(allowed) {
		t.Errorf("%d mounts, want %d: the worktree, its two git mounts, the main checkout's wip, system-flow.yaml, and verify folder, the cache, and the login", mounts, len(allowed))
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
	// Docker makes the main checkout's folder as root, as its mounts'
	// parent, so git is told it is safe there, by command-line configuration.
	wantEnv := []string{"FLAI_CONFIG=" + cache + "/flai/config.json", "FLAI_CACHE_DIR=" + cache + "/flai/cache",
		"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=safe.directory", "GIT_CONFIG_VALUE_0=" + r.Root}
	if !slices.Equal(st.Env, wantEnv) {
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
		{"no story", func(r *harness.Request, _ *harness.Host) { r.Story = "" }, "names neither a story nor a role: start it for one"},
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

// I-0215: a worktree whose registration is absent fails the start with the
// reason, naming the registration, rather than starting a container whose
// worktree is not a git repository; one whose registration names a path with
// a colon is refused, as a mount cannot carry it.
func TestTheContainerRefusesAWorktreeWhoseGitItCannotMount(t *testing.T) {
	r, host := containerReq(t, &manifest.Agent{Harness: harness.Container})
	reg := filepath.Join(r.Root, ".git", "worktrees", "S-0508")
	if err := os.RemoveAll(reg); err != nil {
		t.Fatal(err)
	}
	st, err := harness.Adapters[harness.Container].Start(r, host)
	if err == nil || !strings.Contains(err.Error(), "has no git for the container") || !strings.Contains(err.Error(), reg) {
		t.Errorf("start %v, want the missing registration %s named", err, reg)
	}
	if st.Argv != nil || st.Harness != "" || st.Stdin != nil {
		t.Errorf("start %+v, want none", st)
	}

	r, host = containerReq(t, &manifest.Agent{Harness: harness.Container})
	colon := linkedWorktree(t, filepath.Join(t.TempDir(), "a:b"), "S-0508")
	worktree := filepath.Join(r.Root, ".flai-cache", "worktrees", "S-0508")
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: "+colon+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if st, err := harness.Adapters[harness.Container].Start(r, host); err == nil || !strings.Contains(err.Error(), "holds a colon") || st.Argv != nil {
		t.Errorf("a registration under a colon: start %q %v", st.Argv, err)
	}
}

// ADR-0253: a verify folder the host already has is mounted as it is, and
// a main checkout the container cannot run flai verify against, one with no
// system-flow.yaml or no wip, fails the start naming the path, and makes no
// verify folder.
func TestTheContainerMountsTheMainCheckoutOrRefusesWhatItLacks(t *testing.T) {
	r, host := containerReq(t, &manifest.Agent{Harness: harness.Container})
	verify := filepath.Join(r.Root, ".flai-cache", "verify")
	if err := os.MkdirAll(verify, 0o755); err != nil {
		t.Fatal(err)
	}
	report := filepath.Join(verify, "S-0508.json")
	if err := os.WriteFile(report, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, err := harness.Adapters[harness.Container].Start(r, host)
	if err != nil || !slices.Contains(st.Argv, verify+":"+verify) {
		t.Fatalf("start %q %v, want the verify folder mounted read-write", st.Argv, err)
	}
	if _, err := os.Stat(report); err != nil {
		t.Errorf("the start lost the report the folder held: %v", err)
	}

	for _, missing := range []string{"system-flow.yaml", "wip"} {
		r, host := containerReq(t, &manifest.Agent{Harness: harness.Container})
		gone := filepath.Join(r.Root, missing)
		if err := os.RemoveAll(gone); err != nil {
			t.Fatal(err)
		}
		st, err := harness.Adapters[harness.Container].Start(r, host)
		if err == nil || !strings.Contains(err.Error(), "has no main checkout for the container") || !strings.Contains(err.Error(), gone) {
			t.Errorf("without %s: start %v, want it named", missing, err)
		}
		if st.Argv != nil {
			t.Errorf("without %s: start %q, want none", missing, st.Argv)
		}
		if _, err := os.Stat(filepath.Join(r.Root, ".flai-cache", "verify")); !os.IsNotExist(err) {
			t.Errorf("without %s: a refused start made the verify folder: %v", missing, err)
		}
	}
}

// A worktree whose .git is a folder carries its repository, so the container
// mounts no git and no main checkout beside it.
func TestTheContainerMountsNoGitBesideAWorktreeThatCarriesItsOwn(t *testing.T) {
	r, host := containerReq(t, &manifest.Agent{Harness: harness.Container})
	worktree := filepath.Join(r.Root, ".flai-cache", "worktrees", "S-0508")
	if err := os.Remove(filepath.Join(worktree, ".git")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(worktree, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	st, err := harness.Adapters[harness.Container].Start(r, host)
	if err != nil {
		t.Fatal(err)
	}
	whole := strings.Join(st.Argv, " ")
	if n := strings.Count(strings.Join(st.Argv, "\x00"), "\x00-v\x00"); n != 3 || strings.Contains(whole, filepath.Join(r.Root, ".git")) ||
		strings.Contains(whole, filepath.Join(r.Root, "wip")) || strings.Contains(whole, filepath.Join(r.Root, "system-flow.yaml")) {
		t.Errorf("argv %q, want the worktree, the cache, and the login alone", st.Argv)
	}
}

// roleReq is a start of the strategic agent in role on the container
// harness, in a main checkout of its own that holds the four roles'
// definitions, with a home of its own and the host's login, the token, and a
// provider's key set where a leak would show.
func roleReq(t *testing.T, role string) (harness.Request, harness.Host) {
	t.Helper()
	r, host := containerReq(t, &manifest.Agent{Harness: harness.Container, Model: "claude-opus-5-5", Provider: "gateway"})
	r.ProviderName, r.Provider = "gateway", &manifest.Provider{API: manifest.APIAnthropicMessages, BaseURL: "https://gw.example", KeyEnv: "GATEWAY_KEY"}
	dir := filepath.Join(r.Root, ".claude", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"planner", "orchestrator", "analyzer", "curator"} {
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte("---\nname: "+name+"\ndescription: The "+name+".\n---\nYou are the "+name+".\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r.Story, r.Role, r.Name = "", role, "agent-"+role
	return r, host
}

// roleRuns are a run of each built-in strategic role, as flai serve starts
// it, and the definition its inner session runs as.
var roleRuns = []struct {
	role, agent string
	change      func(*harness.Request)
}{
	{"plan", "planner", func(r *harness.Request) { r.Item, r.AskedOn = "E-0029", "TH-0042" }},
	{"orchestrate", "orchestrator", func(r *harness.Request) {
		r.Reasons, r.PrimeReasons, r.Queue = []string{"moved S-0123 to review: a title"}, []string{"review:S-0123"}, true
	}},
	{"analyze", "analyzer", func(r *harness.Request) { r.Focus = "risk" }},
	{"curate", "curator", func(r *harness.Request) { r.Since = "2026-10-01" }},
}

// S-0509, ADR-0247: a strategic agent's container mounts the main checkout
// read-write as its working directory, its .git/hooks read-only, the role's
// cache folder, and the login read-only, and nothing else; it is named for
// the role, and passed FLAI_ROLE by name rather than FLAI_STORY.
func TestAStrategicAgentsContainerMountsTheCheckout(t *testing.T) {
	for _, c := range roleRuns {
		r, host := roleReq(t, c.role)
		c.change(&r)
		st, err := harness.Adapters[harness.Container].Start(r, host)
		if err != nil {
			t.Fatalf("%s: %v", c.role, err)
		}
		root, hooks := r.Root, filepath.Join(r.Root, ".git", "hooks")
		cache := filepath.Join(r.Root, ".flai-cache", "agents", c.role)
		user := strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid())
		want := []string{"docker", "run", "--rm", "-i", "--name", "flai-agent-" + c.role + "-3f2a-b1", "--user", user,
			"--network", "bridge", "--add-host", "host.docker.internal:host-gateway",
			"-v", root + ":" + root, "-w", root, "-v", hooks + ":" + hooks + ":ro", "-v", cache + ":" + cache, "-v", host.Login + ":/run/claude-login:ro",
			"-e", "FLAI_AGENT", "-e", "FLAI_SESSION", "-e", "FLAI_ROLE", "-e", "FLAI_STARTED_BY", "-e", "FLAI_MCP_URL", "-e", "FLAI_MCP_TOKEN",
			"-e", "FLAI_CONFIG", "-e", "FLAI_CACHE_DIR", "-e", "GATEWAY_KEY", "flai-agent-claude-code:local"}
		if st.Harness != harness.Container || !slices.Equal(st.Argv, want) {
			t.Errorf("%s: harness %q argv\n%q\nwant\n%q", c.role, st.Harness, st.Argv, want)
		}
		if slices.Contains(st.Argv, "FLAI_STORY") {
			t.Errorf("%s: the container is passed FLAI_STORY", c.role)
		}
		wantEnv := []string{"FLAI_CONFIG=" + cache + "/flai/config.json", "FLAI_CACHE_DIR=" + cache + "/flai/cache", "FLAI_ROLE=" + c.role}
		if !slices.Equal(st.Env, wantEnv) {
			t.Errorf("%s: env %q, want %q", c.role, st.Env, wantEnv)
		}
		if !maps.Equal(st.EnvFrom, map[string]string{"GATEWAY_KEY": "GATEWAY_KEY"}) {
			t.Errorf("%s: EnvFrom %v, want the key's variable copied by its own name", c.role, st.EnvFrom)
		}
		whole := strings.Join(st.Argv, " ") + strings.Join(st.Env, " ") + string(st.Stdin)
		for _, secret := range []string{"token-value-never-shown", "key-value-never-shown"} {
			if strings.Contains(whole, secret) {
				t.Errorf("%s: the start carries a secret's value: %s", c.role, secret)
			}
		}
	}
}

// S-0509, ADR-0247 §3: pilot's start request for a strategic agent carries
// its role and its run's item, asked-on thread, focus, since, and reasons,
// is rooted in the main checkout with the role's cache, and is one pilot
// checks and its inner harness starts.
func TestAStrategicAgentsContainerHandsPilotItsRun(t *testing.T) {
	for _, c := range roleRuns {
		r, host := roleReq(t, c.role)
		c.change(&r)
		st, err := harness.Adapters[harness.Container].Start(r, host)
		if err != nil {
			t.Fatalf("%s: %v", c.role, err)
		}
		pr, err := pilot.Read(bytes.NewReader(st.Stdin))
		if err != nil {
			t.Fatalf("%s: pilot reads %s: %v", c.role, st.Stdin, err)
		}
		cache := filepath.Join(r.Root, ".flai-cache", "agents", c.role)
		env := map[string]string{"FLAI_MCP_URL": "http://host.docker.internal:4243/mcp", "FLAI_MCP_TOKEN": "t"}
		if err := pilot.Check(pr, func(n string) string { return env[n] }, func(p string) bool { return p == cache }); err != nil {
			t.Errorf("%s: pilot refuses %s: %v", c.role, st.Stdin, err)
		}
		if pr.Role != c.role || pr.Story != "" || pr.Root != r.Root || pr.Cache != cache || pr.Inner != harness.ClaudeCode || pr.Flai != "" ||
			pr.Item != r.Item || pr.AskedOn != r.AskedOn || pr.Focus != r.Focus || pr.Since != r.Since || pr.Queue != r.Queue ||
			!slices.Equal(pr.Reasons, r.Reasons) || !slices.Equal(pr.PrimeReasons, r.PrimeReasons) || pr.Session != "3f2a-b1" || pr.Name != "agent-"+c.role {
			t.Errorf("%s: request %+v", c.role, pr)
		}
		inner, err := harness.Adapters[pr.Inner].Start(pr, harness.Host{})
		if err != nil {
			t.Fatalf("%s: the inner harness refuses the request pilot reads: %v", c.role, err)
		}
		if i := slices.Index(inner.Argv, "--agent"); i < 0 || i+1 >= len(inner.Argv) || inner.Argv[i+1] != c.agent {
			t.Errorf("%s: the inner start has no --agent %s: %q", c.role, c.agent, inner.Argv)
		}
		if !slices.Contains(inner.Env, "FLAI_ROLE="+c.role) {
			t.Errorf("%s: the inner start's env %q lacks FLAI_ROLE=%s", c.role, inner.Env, c.role)
		}
	}
}

func TestAStrategicAgentsContainerRefusesWhatItCannotRun(t *testing.T) {
	for _, c := range []struct {
		name   string
		change func(*harness.Request, *harness.Host)
		says   string
	}{
		{"neither a story nor a role", func(r *harness.Request, _ *harness.Host) { r.Role = "" }, "names neither a story nor a role: start it for one"},
		{"a story and a role", func(r *harness.Request, _ *harness.Host) { r.Story = "S-0509" }, "not both: start it for one"},
		{"a role that is not a strategic one", func(r *harness.Request, _ *harness.Host) { r.Role = "review" }, `"review" is none of them: start one of them, or run this agent on another harness`},
		{"a role the project defines", func(r *harness.Request, _ *harness.Host) {
			r.Role, r.RoleFile = "review", &roles.Role{Name: "reviewer", Role: "review"}
		}, "run it on another harness until a story carries its file to pilot"},
		{"a login that is the checkout", func(r *harness.Request, h *harness.Host) { h.Login = r.Root }, "which a container never mounts"},
		{"a login that is the home", func(_ *harness.Request, h *harness.Host) { h.Login = filepath.Dir(h.Login) }, "which a container never mounts"},
		{"an analyzer given an item", func(r *harness.Request, _ *harness.Host) { r.Role, r.Item = "analyze", "E-0029" }, "inner harness claude-code refuses the start"},
		{"the operator's arguments", func(_ *harness.Request, h *harness.Host) { h.Args = []string{"-v", "/:/host"} }, "takes no arguments of the operator's"},
	} {
		r, host := roleReq(t, "plan")
		r.Item = "E-0029"
		c.change(&r, &host)
		if _, err := harness.Adapters[harness.Container].Start(r, host); err == nil || !strings.Contains(err.Error(), c.says) {
			t.Errorf("%s: %v, want %q", c.name, err, c.says)
		}
	}
}

func TestTheContainersNameAndCacheAreTheRoles(t *testing.T) {
	if got := harness.ContainerName("plan", "abc"); got != "flai-agent-plan-abc" {
		t.Errorf("name %q", got)
	}
	if got := harness.ContainerName("orchestrate", ""); got != "flai-agent-orchestrate" {
		t.Errorf("name without a session %q", got)
	}
	if got := harness.ContainerCache("/p/flow", "curate"); got != "/p/flow/.flai-cache/agents/curate" {
		t.Errorf("cache %q", got)
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

// S-0618: the docker run of a check in an image runs the image's claude as
// its entrypoint, as the user, on the host's network, with the scratch
// folder mounted at its host path as the working directory's parent, the
// login read-only, and the variables by name alone, then the claude
// arguments, and nothing else.
func TestClaudeCodeImageCheckRunsTheImagesClaudeOverTheScratchFolder(t *testing.T) {
	scratch := t.TempDir()
	workdir := filepath.Join(scratch, "project")
	login := t.TempDir()
	argv, err := harness.ClaudeCodeImageCheck(harness.Host{Program: "podman", Network: "none", Login: login}, "img:local", scratch, workdir, []string{"FLAI_AGENT", "GATEWAY_KEY"}, []string{"-p", "write it"})
	if err != nil {
		t.Fatal(err)
	}
	user := strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid())
	want := []string{"podman", "run", "--rm", "--user", user, "--network", "none", "--add-host", "host.docker.internal:host-gateway",
		"-v", scratch + ":" + scratch, "-w", workdir, "-v", login + ":" + harness.ContainerLogin + ":ro",
		"-e", "FLAI_AGENT", "-e", "GATEWAY_KEY", "--entrypoint", harness.ImageClaude, "img:local", "-p", "write it"}
	if !slices.Equal(argv, want) {
		t.Errorf("argv\n%q\nwant\n%q", argv, want)
	}
	// docker and bridge when the host names neither
	argv, err = harness.ClaudeCodeImageCheck(harness.Host{Login: login}, "img:local", scratch, scratch, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if argv[0] != "docker" || !slices.Contains(argv, "bridge") {
		t.Errorf("argv %q, want docker on the bridge network", argv)
	}
}

// S-0618: the check's docker run refuses what a story's container refuses,
// each saying what to do.
func TestClaudeCodeImageCheckRefusesWhatAContainerCannotRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	scratch := t.TempDir()
	login := t.TempDir()
	for name, c := range map[string]struct {
		host             harness.Host
		image, dir, work string
		names            []string
		want             string
	}{
		"operator's arguments": {harness.Host{Login: login, Args: []string{"--privileged"}}, "img", scratch, scratch, nil, "takes no arguments"},
		"no image":             {harness.Host{Login: login}, "", scratch, scratch, nil, "no image is set"},
		"no login":             {harness.Host{}, "img", scratch, scratch, nil, "--login <path>"},
		"the home as login":    {harness.Host{Login: home}, "img", scratch, scratch, nil, "which a container never mounts"},
		"a relative folder":    {harness.Host{Login: login}, "img", "scratch", "scratch", nil, "not an absolute path"},
		"outside the folder":   {harness.Host{Login: login}, "img", scratch, filepath.Dir(scratch), nil, "is not below its folder"},
		"a colon":              {harness.Host{Login: login}, "img", scratch + ":x", scratch + ":x", nil, "holds a colon"},
		"a value for a name":   {harness.Host{Login: login}, "img", scratch, scratch, []string{"KEY=value"}, "is not a variable's name"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := harness.ClaudeCodeImageCheck(c.host, c.image, c.dir, c.work, c.names, nil)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("err %v, want one saying %q", err, c.want)
			}
		})
	}
}

// commandRun is the close-out's flai verify of S-0508 in its container
// (I-0295), in the project containerReq makes, with the token and the key set
// where a leak would show.
func commandRun(t *testing.T) (harness.CommandRun, harness.Host) {
	t.Helper()
	r, host := containerReq(t, &manifest.Agent{Harness: harness.Container})
	return harness.CommandRun{Root: r.Root, Story: "S-0508", Purpose: "closeout", ID: "20261110T060000Z", Image: "flai-agent-sf-claude-code:local",
		Argv: []string{harness.ImageFlai, "verify", "S-0508", "--record-issues"}}, host
}

// I-0295: a command flai runs for a story whose agent runs in a container
// runs in a container of its own with the story's six mounts, the worktree
// the working directory, and its program the entrypoint, and with no login,
// cache folder, home, ~/.flai, Docker socket, MCP variable, or key.
func TestContainerCommandRunsTheCommandWithTheStorysSixMounts(t *testing.T) {
	c, host := commandRun(t)
	c.Env = []string{"CLOSE_OUT_STORY"}
	verify := filepath.Join(c.Root, ".flai-cache", "verify")
	got, err := harness.ContainerCommand(c, host)
	if err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(verify); err != nil || !info.IsDir() {
		t.Errorf("the run should make the verify folder on the host before docker makes it as root: %v", err)
	}
	worktree := filepath.Join(c.Root, ".flai-cache", "worktrees", "S-0508")
	common, reg := filepath.Join(c.Root, ".git"), filepath.Join(c.Root, ".git", "worktrees", "S-0508")
	wip, flow := filepath.Join(c.Root, "wip"), filepath.Join(c.Root, "system-flow.yaml")
	user := strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid())
	name := "flai-closeout-S-0508-20261110T060000Z"
	want := []string{"docker", "run", "--rm", "--name", name, "--user", user, "--network", "bridge",
		"-v", worktree + ":" + worktree, "-w", worktree, "-v", common + ":" + common + ":ro", "-v", reg + ":" + reg,
		"-v", wip + ":" + wip + ":ro", "-v", flow + ":" + flow + ":ro", "-v", verify + ":" + verify,
		"-e", "GIT_CONFIG_COUNT", "-e", "GIT_CONFIG_KEY_0", "-e", "GIT_CONFIG_VALUE_0", "-e", "CLOSE_OUT_STORY",
		"--entrypoint", "/usr/local/bin/flai", "flai-agent-sf-claude-code:local", "verify", "S-0508", "--record-issues"}
	if !slices.Equal(got.Argv, want) {
		t.Fatalf("argv\n%q\nwant\n%q", got.Argv, want)
	}
	if got.Name != name {
		t.Errorf("name %q, want %q", got.Name, name)
	}
	if wantEnv := []string{"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=safe.directory", "GIT_CONFIG_VALUE_0=" + c.Root}; !slices.Equal(got.Env, wantEnv) {
		t.Errorf("env %q, want %q", got.Env, wantEnv)
	}
	home, _ := os.UserHomeDir()
	whole := strings.Join(got.Argv, " ")
	for _, never := range []string{host.Login, "/run/claude-login", filepath.Join(c.Root, ".flai-cache", "agents"), filepath.Join(home, ".flai"), "docker.sock",
		"FLAI_MCP_URL", "FLAI_MCP_TOKEN", "GATEWAY_KEY", "token-value-never-shown", "key-value-never-shown", "host.docker.internal", " -i "} {
		if strings.Contains(whole, never) {
			t.Errorf("the run carries %q: %s", never, whole)
		}
	}
	for i, arg := range got.Argv {
		if i > 0 && got.Argv[i-1] == "-e" && strings.Contains(arg, "=") {
			t.Errorf("-e %s carries a value", arg)
		}
		if i > 0 && got.Argv[i-1] == "-v" && (arg == home+":"+home || strings.HasPrefix(arg, home+":")) {
			t.Errorf("mounts the home: %s", arg)
		}
	}
	host.Program, host.Network = "/usr/bin/podman", "none"
	if got, err := harness.ContainerCommand(c, host); err != nil || got.Argv[0] != "/usr/bin/podman" || !slices.Contains(got.Argv, "none") {
		t.Errorf("the host's program and network: %q %v", got.Argv, err)
	}
}

// A command run is refused what Start refuses of the story's mounts, and what
// names no container, image, command, or variable, and makes nothing then.
func TestContainerCommandRefusesWhatStartRefuses(t *testing.T) {
	for _, c := range []struct {
		name   string
		change func(*harness.CommandRun, *harness.Host)
		says   string
	}{
		{"no registration", func(c *harness.CommandRun, _ *harness.Host) {
			if err := os.RemoveAll(filepath.Join(c.Root, ".git", "worktrees", "S-0508")); err != nil {
				t.Fatal(err)
			}
		}, "has no git for the container"},
		{"no main checkout's manifest", func(c *harness.CommandRun, _ *harness.Host) {
			if err := os.Remove(filepath.Join(c.Root, "system-flow.yaml")); err != nil {
				t.Fatal(err)
			}
		}, "has no main checkout for the container"},
		{"no main checkout's wip", func(c *harness.CommandRun, _ *harness.Host) {
			if err := os.RemoveAll(filepath.Join(c.Root, "wip")); err != nil {
				t.Fatal(err)
			}
		}, "has no main checkout for the container"},
		{"a registration under a colon", func(c *harness.CommandRun, _ *harness.Host) {
			colon := linkedWorktree(t, filepath.Join(t.TempDir(), "a:b"), "S-0508")
			if err := os.WriteFile(filepath.Join(c.Root, ".flai-cache", "worktrees", "S-0508", ".git"), []byte("gitdir: "+colon+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}, "holds a colon"},
		{"the operator's arguments", func(_ *harness.CommandRun, h *harness.Host) { h.Args = []string{"-v", "/:/host"} }, "takes no arguments of the operator's"},
		{"a relative root", func(c *harness.CommandRun, _ *harness.Host) { c.Root = "flow" }, "not an absolute path"},
		{"no story", func(c *harness.CommandRun, _ *harness.Host) { c.Story = "" }, "is not a story's ID"},
		{"a story that is a path", func(c *harness.CommandRun, _ *harness.Host) { c.Story = "../S-0508" }, "is not a story's ID"},
		{"an id that names no container", func(c *harness.CommandRun, _ *harness.Host) { c.ID = "a b" }, "name no container"},
		{"no image", func(c *harness.CommandRun, _ *harness.Host) { c.Image = "" }, "no image is given"},
		{"no command", func(c *harness.CommandRun, _ *harness.Host) { c.Argv = nil }, "no command is given"},
		{"a variable that is no name", func(c *harness.CommandRun, _ *harness.Host) { c.Env = []string{"K=v"} }, "is not a variable's name"},
	} {
		run, host := commandRun(t)
		c.change(&run, &host)
		got, err := harness.ContainerCommand(run, host)
		if err == nil || !strings.Contains(err.Error(), c.says) {
			t.Errorf("%s: %v, want %q", c.name, err, c.says)
		}
		if got.Argv != nil || got.Name != "" {
			t.Errorf("%s: %+v, want none", c.name, got)
		}
		if _, err := os.Stat(filepath.Join(run.Root, ".flai-cache", "verify")); err == nil {
			t.Errorf("%s: a refused run made the verify folder", c.name)
		}
	}
}
