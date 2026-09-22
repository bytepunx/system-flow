package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
