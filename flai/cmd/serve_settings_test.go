package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0506: flai serve settings prints the host settings settings.get gives,
// its mcp block naming the bind address flai serve mcp bind sets, or the
// loopback default, as text and as JSON, and never the MCP token.
func TestServeSettings(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := mcpToken(repo, false)
	if err != nil || token == "" {
		t.Fatalf("token: %q %v", token, err)
	}
	mcpOf := func(js string) map[string]any {
		t.Helper()
		var got struct {
			MCP map[string]any `json:"mcp"`
		}
		if err := json.Unmarshal([]byte(js), &got); err != nil {
			t.Fatalf("--json is not JSON: %v\n%s", err, js)
		}
		return got.MCP
	}
	check := func(bind string, set bool, text string) {
		t.Helper()
		out, errOut, code := runIn(t, root, "serve", "settings")
		if code != 0 || !strings.Contains(out, text) || !strings.HasPrefix(out, "host settings for t (") {
			t.Errorf("text: %d %q %s", code, out, errOut)
		}
		js, errOut, code := runIn(t, root, "serve", "settings", "--json")
		if mcp := mcpOf(js); code != 0 || mcp["bind"] != bind || mcp["bind_set"] != set || mcp["running"] != false {
			t.Errorf("--json: %d %v %s", code, mcp, errOut)
		}
		for _, printed := range []string{out, js} {
			if strings.Contains(printed, token) {
				t.Errorf("the token is printed: %s", printed)
			}
		}
	}
	check("127.0.0.1", false, "mcp: not running; bind 127.0.0.1 (default; flai serve mcp bind <address> sets another)\n")
	if _, errOut, code := runIn(t, root, "serve", "mcp", "bind", "0.0.0.0"); code != 0 {
		t.Fatalf("bind: %d %s", code, errOut)
	}
	check("0.0.0.0", true, "mcp: not running; bind 0.0.0.0 (set; flai serve mcp clear goes back to the default)\n")
	if _, errOut, code := runIn(t, t.TempDir(), "serve", "settings"); code == 0 || !strings.Contains(errOut, "run flai new or flai import first") {
		t.Errorf("outside a project it says what to do: %d %s", code, errOut)
	}
}
