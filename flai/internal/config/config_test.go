package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

func TestLoadCreatesDefaultsOnFirstRun(t *testing.T) {
	t.Setenv(CacheEnvVar, "/tmp/flai-test-cache")
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	cfg, created, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !created {
		t.Fatal("expected created=true on first run")
	}
	if cfg.Template.Repo != Default().Template.Repo || cfg.Dashboard.Port != 4242 || cfg.CacheDir != "/tmp/flai-test-cache" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file not written: %v", err)
	}
	cfg2, created, err := Load(path)
	if err != nil || created {
		t.Fatalf("second Load: created=%v err=%v", created, err)
	}
	if !reflect.DeepEqual(cfg2, cfg) {
		t.Fatalf("round trip mismatch:\n%+v\n%+v", cfg, cfg2)
	}
}

func TestSetAndGetRoundTrip(t *testing.T) {
	cfg := Default()
	cases := map[string]string{
		"template.repo":   "git@github.com:me/my-template.git",
		"template.ref":    "feature/x",
		"dashboard.image": "ghcr.io/me/flaiover",
		"dashboard.tag":   "0.2.0",
		"dashboard.port":  "8080",
		"cache_dir":       "/tmp/cache",
		"author":          "alex",
	}
	for k, v := range cases {
		var err error
		cfg, err = cfg.Set(k, v)
		if err != nil {
			t.Fatalf("Set %s: %v", k, err)
		}
	}
	if cfg.Dashboard.Port != 8080 || cfg.Template.Ref != "feature/x" {
		t.Fatalf("Set did not apply: %+v", cfg)
	}
	got, err := cfg.Get("dashboard.port")
	if err != nil {
		t.Fatal(err)
	}
	if got.(float64) != 8080 {
		t.Fatalf("Get dashboard.port = %v", got)
	}
	if _, err := cfg.Get("template"); err != nil {
		t.Fatalf("Get nested object: %v", err)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := Save(path, cfg); err != nil {
		t.Fatal(err)
	}
	back, _, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, cfg) {
		t.Fatalf("save/load mismatch:\n%+v\n%+v", cfg, back)
	}
}

// S-0239: dashboard.allow_unsigned is off on a fresh host and config set turns
// it on and off again as a bool.
func TestAllowUnsigned(t *testing.T) {
	cfg := Default()
	if got, err := cfg.Get("dashboard.allow_unsigned"); err != nil || got != false {
		t.Fatalf("fresh dashboard.allow_unsigned = %v, %v; want false", got, err)
	}
	on, err := cfg.Set("dashboard.allow_unsigned", "true")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := on.Get("dashboard.allow_unsigned"); err != nil || got != true || !on.Dashboard.AllowUnsigned {
		t.Fatalf("after set true: %v, %v, field %v", got, err, on.Dashboard.AllowUnsigned)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := Save(path, on); err != nil {
		t.Fatal(err)
	}
	back, _, err := Load(path)
	if err != nil || !back.Dashboard.AllowUnsigned {
		t.Fatalf("save/load lost allow_unsigned: %+v, %v", back.Dashboard, err)
	}
	off, err := back.Set("dashboard.allow_unsigned", "false")
	if err != nil || off.Dashboard.AllowUnsigned {
		t.Fatalf("after set false: %+v, %v", off.Dashboard, err)
	}
	if _, err := cfg.Set("dashboard.allow_unsigned", "yes please"); err == nil {
		t.Fatal("expected error for a value that is not a bool")
	}
}

// S-0489: dashboard.allowed_hosts is empty on a fresh host, config set takes a
// list of host names as written, which save and load keep, and refuses an
// entry that is no host name with an optional port, naming it.
func TestAllowedHosts(t *testing.T) {
	cfg := Default()
	if got, err := cfg.Get("dashboard.allowed_hosts"); err != nil || got != "" {
		t.Fatalf("fresh dashboard.allowed_hosts = %v, %v; want empty", got, err)
	}
	const list = "Flow.Example.com, board.example.net:8443"
	set, err := cfg.Set("dashboard.allowed_hosts", list)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := set.Get("dashboard.allowed_hosts"); err != nil || got != list || set.Dashboard.AllowedHosts != list {
		t.Fatalf("after set: %v, %v, field %q", got, err, set.Dashboard.AllowedHosts)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := Save(path, set); err != nil {
		t.Fatal(err)
	}
	back, _, err := Load(path)
	if err != nil || back.Dashboard.AllowedHosts != list {
		t.Fatalf("save/load lost allowed_hosts: %+v, %v", back.Dashboard, err)
	}
	cleared, err := back.Set("dashboard.allowed_hosts", "")
	if err != nil || cleared.Dashboard.AllowedHosts != "" {
		t.Fatalf("after set empty: %+v, %v", cleared.Dashboard, err)
	}
	for _, bad := range []string{"https://flow.example.com", "flow.example.com/board", "flow example.com", "flow.example.com:http", "flow.example.com:0", "flow.example.com:70000", ":8443", "::1", "[::1", "[::1]x", "[]:80", "me@flow.example.com"} {
		_, err := cfg.Set("dashboard.allowed_hosts", "ok.example.com, "+bad)
		if err == nil || !strings.Contains(err.Error(), "dashboard.allowed_hosts") || !strings.Contains(err.Error(), strings.ToLower(bad)) || strings.Contains(err.Error(), `"ok.example.com"`) {
			t.Errorf("set %q: %v; want a refusal naming that entry alone", bad, err)
		}
	}
	for _, good := range []string{"localhost", "localhost:4242", "192.168.1.20:4242", "[::1]", "[::1]:4242", "board-1.example.co.uk"} {
		if _, err := cfg.Set("dashboard.allowed_hosts", good); err != nil {
			t.Errorf("set %q: %v; want it taken", good, err)
		}
	}
	if _, err := Parse([]byte(`{"dashboard":{"allowed_hosts":"https://flow.example.com"}}`)); err == nil {
		t.Error("a configuration file with a scheme in allowed_hosts loads")
	}
}

// S-0489: the list flai passes the container is each entry trimmed and
// lower-cased, the empty ones dropped, joined with commas and no spaces.
func TestNormalizeAllowedHosts(t *testing.T) {
	for _, c := range []struct{ raw, want string }{
		{"", ""},
		{" , ,", ""},
		{"flow.example.com", "flow.example.com"},
		{" Flow.Example.COM ", "flow.example.com"},
		{"a.example, B.example:8443", "a.example,b.example:8443"},
		{",a.example,, b.example ,", "a.example,b.example"},
	} {
		if got := NormalizeAllowedHosts(c.raw); got != c.want {
			t.Errorf("NormalizeAllowedHosts(%q) = %q, want %q", c.raw, got, c.want)
		}
	}
	if got := AllowedHostList(" , "); got != nil {
		t.Errorf("AllowedHostList of no entry = %#v, want nil", got)
	}
}

func TestSetRejectsBadInput(t *testing.T) {
	cfg := Default()
	if _, err := cfg.Set("dashboard.port", "http"); err == nil {
		t.Fatal("expected error for non-integer port")
	}
	if _, err := cfg.Set("nope.key", "x"); err == nil {
		t.Fatal("expected error for unknown key")
	}
	if _, err := cfg.Set("template", "x"); err == nil {
		t.Fatal("expected error when setting an object")
	}
	if _, err := cfg.Get("dashboard.nope"); err == nil {
		t.Fatal("expected error for unknown nested key")
	}
}

func TestParseRejectsUnknownKeys(t *testing.T) {
	_, err := Parse([]byte(`{"template":{"repo":"x","ref":"y"},"typo":1}`))
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestLoadReportsParseErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err := Load(path)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatal("parse error should not look like missing file")
	}
}

func TestResolvePath(t *testing.T) {
	t.Setenv(EnvVar, "")
	if got := ResolvePath("/x/y.json"); got != "/x/y.json" {
		t.Fatalf("override: %s", got)
	}
	t.Setenv(EnvVar, "/env/config.json")
	if got := ResolvePath(""); got != "/env/config.json" {
		t.Fatalf("env: %s", got)
	}
	if got := ResolvePath("/flag.json"); got != "/flag.json" {
		t.Fatalf("flag beats env: %s", got)
	}
	t.Setenv(EnvVar, "")
	if got := ResolvePath(""); filepath.Base(got) != "config.json" || filepath.Base(filepath.Dir(got)) != ".flai" {
		t.Fatalf("default: %s", got)
	}
	home, _ := os.UserHomeDir()
	if got := ExpandHome("~/a"); got != filepath.Join(home, "a") {
		t.Fatalf("ExpandHome: %s", got)
	}
}

func TestKeysMatchSchema(t *testing.T) {
	cfg := Default()
	for _, k := range Keys() {
		if _, err := cfg.Get(k); err != nil {
			t.Errorf("Keys() lists %s but Get fails: %v", k, err)
		}
	}
}

// S-0078: host actions are off by default and enabled by name, per project.
func TestHostActions(t *testing.T) {
	c := Default()
	if c.ActionEnabled("push", "/a") || c.HostActions != nil {
		t.Fatal("nothing is enabled by default")
	}
	c = c.WithAction("push", "/a", true)
	if !c.ActionEnabled("push", "/a") || c.ActionEnabled("push", "/b") || c.ActionEnabled("agent", "/a") {
		t.Errorf("enabled for one project and one action: %+v", c.HostActions)
	}
	c = c.WithAction("push", "/a", true) // twice is once
	if len(c.HostActions["push"]) != 1 {
		t.Errorf("not listed twice: %+v", c.HostActions)
	}
	all := c.WithAction("push", AllProjects, true)
	if !all.ActionEnabled("push", "/b") {
		t.Error("every project")
	}
	if off := all.WithAction("push", AllProjects, false); off.ActionEnabled("push", "/a") || off.HostActions != nil {
		t.Errorf("disabling for all leaves none: %+v", off.HostActions)
	}
	if off := all.WithAction("push", "/a", false); !off.ActionEnabled("push", "/a") {
		t.Error("with all projects enabled, one project's entry going changes nothing, and status must say so")
	}
	for _, k := range Keys() {
		if strings.HasPrefix(k, "host_actions") {
			t.Error("flai config set cannot enable an action: only flai serve enable does")
		}
	}
	if _, err := c.Set("host_actions.push", "*"); err == nil {
		t.Error("flai config set refuses the key")
	}
}

// S-0294 (ADR-0108): agent.auto_restarts is 2 while unset, 0 turns the
// automatic restart off and stays apart from unset through a save and a
// load, and a negative count is refused when set and when loaded.
func TestAgentAutoRestarts(t *testing.T) {
	c := Default()
	if got := c.Agent.AutoRestartLimit(); got != 2 || c.Agent.AutoRestarts != nil {
		t.Errorf("unset: limit %d, field %v; want 2 and nil", got, c.Agent.AutoRestarts)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	for _, n := range []int{0, 3} {
		agent, err := c.Agent.WithAutoRestarts(n)
		if err != nil {
			t.Fatalf("WithAutoRestarts(%d): %v", n, err)
		}
		c.Agent = agent
		if err := Save(path, c); err != nil {
			t.Fatal(err)
		}
		back, _, err := Load(path)
		if err != nil {
			t.Fatal(err)
		}
		if back.Agent.AutoRestarts == nil || back.Agent.AutoRestartLimit() != n {
			t.Errorf("set to %d, loaded back as %v", n, back.Agent.AutoRestarts)
		}
	}
	data, err := json.Marshal(Default())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "auto_restarts") {
		t.Errorf("unset is not written: %s", data)
	}
	if _, err := c.Agent.WithAutoRestarts(-1); err == nil || !strings.Contains(err.Error(), "0 or more") {
		t.Errorf("a negative count is refused, saying what to give: %v", err)
	}
	if _, err := Parse([]byte(`{"agent":{"auto_restarts":-1}}`)); err == nil || !strings.Contains(err.Error(), "agent.auto_restarts is -1") {
		t.Errorf("a negative count in the file is refused on load: %v", err)
	}
	got, err := Parse([]byte(`{"agent":{"auto_restarts":0}}`))
	if err != nil || got.Agent.AutoRestarts == nil || got.Agent.AutoRestartLimit() != 0 {
		t.Errorf("0 read from the file is off, not unset: %v %v", got.Agent.AutoRestarts, err)
	}
}

// S-0349 (ADR-0129): agent.providers round-trips whole through a save and a
// load; an entry with an unknown api, a base_url that is not an http or https
// URL, a key_env that is not a variable's name, or a field left out is
// refused when set and when loaded, saying what sets it; and an entry is
// removed alone.
func TestAgentProviders(t *testing.T) {
	whole := manifest.Provider{API: manifest.APIOpenAIChat, BaseURL: "https://litellm.example.com", KeyEnv: "LITELLM_API_KEY",
		Models: map[string]string{"opus": "anthropic/claude-opus-5-5"}}
	c := Default()
	agent, err := c.Agent.WithProvider("litellm", whole)
	if err != nil {
		t.Fatalf("a whole entry: %v", err)
	}
	if c.Agent.Providers != nil {
		t.Errorf("WithProvider copies the map: %+v", c.Agent.Providers)
	}
	c.Agent = agent
	path := filepath.Join(t.TempDir(), "config.json")
	if err := Save(path, c); err != nil {
		t.Fatal(err)
	}
	back, _, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back.Agent.Providers, map[string]manifest.Provider{"litellm": whole}) {
		t.Errorf("round trip: %+v", back.Agent.Providers)
	}
	for _, bad := range []struct {
		name, want string
		p          manifest.Provider
	}{
		{"unknown api", `agent.providers.litellm.api "bedrock" is not an API flai speaks`, manifest.Provider{API: "bedrock", BaseURL: whole.BaseURL, KeyEnv: whole.KeyEnv}},
		{"base_url not http", `agent.providers.litellm.base_url "ftp://x" is not an absolute http or https URL`, manifest.Provider{API: whole.API, BaseURL: "ftp://x", KeyEnv: whole.KeyEnv}},
		{"key_env not a name", `agent.providers.litellm.key_env is not the name of an environment variable`, manifest.Provider{API: whole.API, BaseURL: whole.BaseURL, KeyEnv: "sk-123 456"}},
		{"no key_env", "agent.providers.litellm has no key_env; set it with flai serve agent provider litellm --key-env", manifest.Provider{API: whole.API, BaseURL: whole.BaseURL}},
		{"no api", "agent.providers.litellm has no api; set it with flai serve agent provider litellm --api", manifest.Provider{BaseURL: whole.BaseURL, KeyEnv: whole.KeyEnv}},
	} {
		if _, err := back.Agent.WithProvider("litellm", bad.p); err == nil || !strings.Contains(err.Error(), bad.want) {
			t.Errorf("%s is refused when set, saying so: %v", bad.name, err)
		}
		data, err := json.Marshal(map[string]any{"agent": map[string]any{"providers": map[string]manifest.Provider{"litellm": bad.p}}})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Parse(data); err == nil || !strings.Contains(err.Error(), bad.want) {
			t.Errorf("%s in the file is refused on load: %v", bad.name, err)
		}
	}
	if !reflect.DeepEqual(back.Agent.Providers["litellm"], whole) {
		t.Errorf("a refused entry changes nothing: %+v", back.Agent.Providers)
	}
	if got := back.Agent.WithoutProvider("litellm"); got.Providers != nil || len(back.Agent.Providers) != 1 {
		t.Errorf("removed, and the map copied: %+v %+v", got.Providers, back.Agent.Providers)
	}
}

// S-0392: a forge's token is read from GITHUB_TOKEN while the host names no
// variable; forges round-trips through a save and a load, and is left out of
// the file while unset; an unknown forge or a token_env that is not a
// variable's name is refused when set and when loaded, without quoting it;
// an entry that sets nothing is removed; and forges stays out of Keys.
func TestForges(t *testing.T) {
	c := Default()
	if got := c.ForgeTokenEnv(manifest.ForgeGitHub); got != "GITHUB_TOKEN" {
		t.Errorf("unset, github's token is read from %q; want GITHUB_TOKEN", got)
	}
	if got := c.ForgeTokenEnv("gitlab"); got != "" {
		t.Errorf("a forge flai does not call has no variable: %q", got)
	}
	if got := KnownForges(); !reflect.DeepEqual(got, []string{"github"}) {
		t.Errorf("known forges: %v", got)
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "forges") {
		t.Errorf("unset is not written: %s", data)
	}
	set, err := c.WithForge(manifest.ForgeGitHub, ForgeHost{TokenEnv: "FLAI_GH_TOKEN"})
	if err != nil {
		t.Fatalf("WithForge: %v", err)
	}
	if c.Forges != nil {
		t.Errorf("WithForge copies the map: %+v", c.Forges)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := Save(path, set); err != nil {
		t.Fatal(err)
	}
	back, _, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back.Forges, map[string]ForgeHost{"github": {TokenEnv: "FLAI_GH_TOKEN"}}) || back.ForgeTokenEnv(manifest.ForgeGitHub) != "FLAI_GH_TOKEN" {
		t.Errorf("round trip: %+v", back.Forges)
	}
	if file, _ := os.ReadFile(path); !strings.Contains(string(file), `"forges": {`) || !strings.Contains(string(file), `"token_env": "FLAI_GH_TOKEN"`) {
		t.Errorf("written as forges.github.token_env: %s", file)
	}
	const pasted = "ghp_pasted0123456789"
	for _, bad := range []struct {
		name, forge, env, want string
	}{
		{"unknown forge", "gitlab", "GITLAB_TOKEN", `forges.gitlab: flai calls no forge "gitlab"; it calls github`},
		{"a token in its place", manifest.ForgeGitHub, pasted + " x", "forges.github.token_env is not the name of an environment variable; write one such as GITHUB_TOKEN"},
		{"a leading digit", manifest.ForgeGitHub, "1TOKEN", "flai serve forge github --token-env <NAME>"},
	} {
		if _, err := back.WithForge(bad.forge, ForgeHost{TokenEnv: bad.env}); err == nil || !strings.Contains(err.Error(), bad.want) || strings.Contains(err.Error(), pasted) {
			t.Errorf("%s is refused when set, saying so: %v", bad.name, err)
		}
		data, err := json.Marshal(map[string]any{"forges": map[string]ForgeHost{bad.forge: {TokenEnv: bad.env}}})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Parse(data); err == nil || !strings.Contains(err.Error(), bad.want) {
			t.Errorf("%s in the file is refused on load: %v", bad.name, err)
		}
	}
	if cleared, err := back.WithForge(manifest.ForgeGitHub, ForgeHost{}); err != nil || cleared.Forges != nil {
		t.Errorf("an entry that sets nothing is removed: %+v %v", cleared.Forges, err)
	}
	if got := back.WithoutForge(manifest.ForgeGitHub); got.Forges != nil || got.ForgeTokenEnv(manifest.ForgeGitHub) != "GITHUB_TOKEN" || len(back.Forges) != 1 {
		t.Errorf("removed, back to the default, and the map copied: %+v %+v", got.Forges, back.Forges)
	}
	for _, k := range Keys() {
		if strings.HasPrefix(k, "forges") {
			t.Error("flai config set cannot reach forges: only flai serve forge does")
		}
	}
}

// S-0392: the token is read through getenv only when the function is called,
// from the variable the host names at the time the function was made, and an
// unset or empty variable is refused naming it and the setting that renames
// it, never a value.
func TestForgeTokenFunc(t *testing.T) {
	const token = "ghp_sentinel0123456789"
	env := map[string]string{}
	var asked []string
	getenv := func(name string) string {
		asked = append(asked, name)
		return env[name]
	}
	c := Default()
	read := c.ForgeTokenFunc(manifest.ForgeGitHub, getenv)
	if len(asked) != 0 {
		t.Errorf("nothing is read before the call: %v", asked)
	}
	_, err := read()
	if err == nil || !strings.Contains(err.Error(), "the github token is read from the variable GITHUB_TOKEN, which is not set") || !strings.Contains(err.Error(), "flai serve forge github --token-env <NAME>") {
		t.Errorf("an unset variable is refused, naming it and the setting: %v", err)
	}
	env["GITHUB_TOKEN"] = "  \n"
	if _, err := read(); err == nil {
		t.Error("a variable holding only space is refused")
	}
	env["GITHUB_TOKEN"] = token + "\n"
	if got, err := read(); err != nil || got != token {
		t.Errorf("read at the call: %q %v", got, err)
	}
	if !reflect.DeepEqual(asked, []string{"GITHUB_TOKEN", "GITHUB_TOKEN", "GITHUB_TOKEN"}) {
		t.Errorf("read once per call, from GITHUB_TOKEN: %v", asked)
	}
	named, err := c.WithForge(manifest.ForgeGitHub, ForgeHost{TokenEnv: "FLAI_GH_TOKEN"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = named.ForgeTokenFunc(manifest.ForgeGitHub, getenv)()
	if err == nil || !strings.Contains(err.Error(), "FLAI_GH_TOKEN") || strings.Contains(err.Error(), token) {
		t.Errorf("the host's variable is read, not the default's, and its error holds no token: %v", err)
	}
	env["FLAI_GH_TOKEN"] = token
	if got, err := named.ForgeTokenFunc(manifest.ForgeGitHub, getenv)(); err != nil || got != token {
		t.Errorf("the host's variable: %q %v", got, err)
	}
	if _, err := c.ForgeTokenFunc("gitlab", getenv)(); err == nil || !strings.Contains(err.Error(), `flai calls no forge "gitlab"`) {
		t.Errorf("a forge flai does not call is refused: %v", err)
	}
}

// S-0355 (ADR-0131): a harness's guard none and deny_protected round-trip
// through a save and a load, an entry that sets nothing is removed, and a
// guard other than none is refused when set, saved, and loaded, saying what
// sets it.
func TestHarnessGuardAndDenyProtected(t *testing.T) {
	c := Default()
	agent, err := c.Agent.WithHarness("command", HarnessHost{Guard: GuardNone, DenyProtected: true})
	if err != nil {
		t.Fatalf("guard none: %v", err)
	}
	if c.Agent.Harnesses != nil {
		t.Errorf("WithHarness copies the map: %+v", c.Agent.Harnesses)
	}
	c.Agent = agent
	path := filepath.Join(t.TempDir(), "config.json")
	if err := Save(path, c); err != nil {
		t.Fatal(err)
	}
	back, _, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := back.Agent.Harnesses["command"]; got != (HarnessHost{Guard: GuardNone, DenyProtected: true}) {
		t.Errorf("round trip: %+v", back.Agent.Harnesses)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"guard": "none"`) || !strings.Contains(string(data), `"deny_protected": true`) {
		t.Errorf("written as agent.harnesses.command.guard and deny_protected: %s", data)
	}
	if cleared, err := back.Agent.WithHarness("command", HarnessHost{}); err != nil || cleared.Harnesses != nil {
		t.Errorf("an entry that sets nothing is removed: %+v %v", cleared.Harnesses, err)
	}
	want := `agent.harnesses.claude-code.guard is "off"`
	how := `flai serve agent harness claude-code --guard none, or clear it with --guard ""`
	if _, err := back.Agent.WithHarness("claude-code", HarnessHost{Guard: "off"}); err == nil || !strings.Contains(err.Error(), want) || !strings.Contains(err.Error(), how) {
		t.Errorf("a guard other than none is refused when set, saying what sets it: %v", err)
	}
	if _, err := Parse([]byte(`{"agent":{"harnesses":{"claude-code":{"guard":"off"}}}}`)); err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("a guard other than none in the file is refused on load: %v", err)
	}
	bad := back
	bad.Agent.Harnesses = map[string]HarnessHost{"claude-code": {Guard: "off"}}
	if err := Save(path, bad); err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("a guard other than none is refused on save: %v", err)
	}
	if again, _, err := Load(path); err != nil || again.Agent.Harnesses["command"].Guard != GuardNone {
		t.Errorf("a refused save leaves the file: %+v %v", again.Agent.Harnesses, err)
	}
}

// S-0082: checks round-trips as a list of named argument lists, and, like
// agent, stays out of Keys() and unreachable by flai config set.
func TestChecksRoundTrip(t *testing.T) {
	c := Default()
	if len(c.Checks.Commands) != 0 || c.Checks.TimeoutMinutes != 0 {
		t.Fatalf("nothing named by default: %+v", c.Checks)
	}
	c.Checks.Commands = []NamedCommand{{Name: "flai", Command: []string{"scripts/flai-test.sh"}}}
	c.Checks.TimeoutMinutes = 20
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Checks.Commands) != 1 || got.Checks.Commands[0].Name != "flai" || got.Checks.TimeoutMinutes != 20 {
		t.Errorf("round trip: %+v", got.Checks)
	}
	for _, k := range Keys() {
		if strings.HasPrefix(k, "checks") {
			t.Error("flai config set cannot reach checks: only flai serve checks does")
		}
	}
	if _, err := Default().Set("checks.timeout_minutes", "5"); err == nil {
		t.Error("flai config set refuses the key, at least while nothing is named yet")
	}
}
