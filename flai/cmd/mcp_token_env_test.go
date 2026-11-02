package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// S-0506: in a run flai serve handed the HTTP MCP server's token as
// FLAI_MCP_TOKEN, flai mcp token prints that token, the one in force, and
// reads and writes no file, in a project or in a folder that is none, as in a
// container with no .flai-cache; --rotate ignores it and rotates the
// project's file; with it unset, the file's token is printed as before.
func TestMCPTokenPrintsTheTokenHandedByEnvironment(t *testing.T) {
	const handed = "env-token-sentinel-S0506"
	t.Setenv(serve.MCPTokenEnv, handed)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	file := filepath.Join(root, ".flai-cache", serve.MCPTokenFile)

	for _, dir := range []string{root, t.TempDir()} {
		out, errOut, code := runIn(t, dir, "mcp", "token")
		if code != 0 || strings.TrimSpace(out) != handed {
			t.Errorf("in %s: %d %q %s", dir, code, out, errOut)
		}
	}
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Errorf("a token handed by environment wrote the project's file: %v", err)
	}
	out, _, code := runIn(t, root, "--json", "mcp", "token")
	var got struct {
		Token string `json:"token"`
		Env   string `json:"env"`
		File  string `json:"file"`
	}
	if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil || got.Token != handed || got.Env != serve.MCPTokenEnv || got.File != "" {
		t.Errorf("--json: %d %v %s", code, err, out)
	}

	rotated, _, code := runIn(t, root, "mcp", "token", "--rotate")
	rotated = strings.TrimSpace(rotated)
	if code != 0 || rotated == handed || len(rotated) < 40 {
		t.Fatalf("--rotate rotates the project's file, whatever the environment holds: %d %q", code, rotated)
	}
	if data, err := os.ReadFile(file); err != nil || strings.TrimSpace(string(data)) != rotated {
		t.Errorf("the rotated token is not the file's: %v", err)
	}

	t.Setenv(serve.MCPTokenEnv, "")
	if out, _, code := runIn(t, root, "mcp", "token"); code != 0 || strings.TrimSpace(out) != rotated {
		t.Errorf("unset, the file's token is printed: %d %q", code, out)
	}
}
