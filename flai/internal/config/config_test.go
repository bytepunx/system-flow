package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCreatesDefaultsOnFirstRun(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	cfg, created, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !created {
		t.Fatal("expected created=true on first run")
	}
	if cfg.Template.Repo != Default().Template.Repo || cfg.Dashboard.Port != 4242 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file not written: %v", err)
	}
	cfg2, created, err := Load(path)
	if err != nil || created {
		t.Fatalf("second Load: created=%v err=%v", created, err)
	}
	if cfg2 != cfg {
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
	if back != cfg {
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
