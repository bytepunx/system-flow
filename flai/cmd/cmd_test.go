package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = run(args, &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestVersionPlainAndJSON(t *testing.T) {
	out, _, code := runCLI(t, "version")
	if code != 0 || !strings.HasPrefix(out, "flai dev") || !strings.Contains(out, "commit:") || !strings.Contains(out, "built:") {
		t.Fatalf("plain version: code=%d out=%q", code, out)
	}
	out, _, code = runCLI(t, "version", "--json")
	if code != 0 {
		t.Fatalf("json version exit %d", code)
	}
	var v map[string]string
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	for _, k := range []string{"version", "commit", "date", "go"} {
		if v[k] == "" {
			t.Errorf("missing %s in %v", k, v)
		}
	}
}

func TestConfigFirstRunGetSet(t *testing.T) {
	t.Setenv("FLAI_CONFIG", "")
	path := filepath.Join(t.TempDir(), "cfg.json")

	out, errOut, code := runCLI(t, "config", "path", "--config", path)
	if code != 0 || strings.TrimSpace(out) != path {
		t.Fatalf("config path: %q %q %d", out, errOut, code)
	}
	if _, err := os.Stat(path); err == nil {
		t.Fatal("config path must not create the file")
	}

	out, errOut, code = runCLI(t, "config", "get", "--config", path)
	if code != 0 || !strings.Contains(errOut, "created") || !strings.Contains(out, "\"template\"") {
		t.Fatalf("first run get: %q %q %d", out, errOut, code)
	}

	out, _, code = runCLI(t, "config", "set", "dashboard.port", "9999", "--config", path)
	if code != 0 || strings.TrimSpace(out) != "dashboard.port = 9999" {
		t.Fatalf("set: %q %d", out, code)
	}
	out, _, code = runCLI(t, "config", "get", "dashboard.port", "--config", path)
	if code != 0 || strings.TrimSpace(out) != "9999" {
		t.Fatalf("get after set: %q %d", out, code)
	}
	out, _, code = runCLI(t, "--config", path, "config", "get", "template")
	if code != 0 || !strings.Contains(out, "\"repo\"") {
		t.Fatalf("get object: %q %d", out, code)
	}

	t.Setenv("FLAI_CONFIG", path)
	out, _, code = runCLI(t, "config", "get", "dashboard.port")
	if code != 0 || strings.TrimSpace(out) != "9999" {
		t.Fatalf("env config: %q %d", out, code)
	}

	_, errOut, code = runCLI(t, "config", "set", "dashboard.port", "abc", "--config", path)
	if code == 0 || !strings.Contains(errOut, "integer") {
		t.Fatalf("bad set should fail: %q %d", errOut, code)
	}
	_, errOut, code = runCLI(t, "config", "get", "nope", "--config", path)
	if code == 0 || !strings.Contains(errOut, "unknown key") {
		t.Fatalf("unknown key should fail: %q %d", errOut, code)
	}
}

func TestHelpListsCommands(t *testing.T) {
	out, _, code := runCLI(t, "--help")
	if code != 0 {
		t.Fatalf("help exit %d", code)
	}
	for _, c := range []string{"version", "config"} {
		if !strings.Contains(out, c) {
			t.Errorf("help missing %s:\n%s", c, out)
		}
	}
}
