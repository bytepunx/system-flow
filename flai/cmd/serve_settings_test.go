package cmd

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0506, S-0732: flai serve settings prints the host settings settings.get
// gives, its mcp block naming the addresses flai serve mcp bind sets, or
// auto, the default, and what auto resolves to now, as text and as JSON, and
// never the MCP token.
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
	native := mcpDocker{engine: "Ubuntu 24.04 LTS"}
	run := func(args ...string) (string, string, int) {
		t.Helper()
		return runWithApp(t, &app{configPath: configOf(t), cwd: root, runner: native}, args...)
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
	check := func(bind string, set bool, binds []any, text string) {
		t.Helper()
		out, errOut, code := run("serve", "settings")
		if code != 0 || !strings.Contains(out, text) || !strings.HasPrefix(out, "host settings for t (") {
			t.Errorf("text: %d %q %s", code, out, errOut)
		}
		js, errOut, code := run("serve", "settings", "--json")
		if mcp := mcpOf(js); code != 0 || mcp["bind"] != bind || mcp["bind_set"] != set || !reflect.DeepEqual(mcp["binds"], binds) || mcp["listen"] != nil || mcp["running"] != false {
			t.Errorf("--json: %d %v %s", code, mcp, errOut)
		}
		for _, printed := range []string{out, js} {
			if strings.Contains(printed, token) {
				t.Errorf("the token is printed: %s", printed)
			}
		}
	}
	check("auto", false, []any{"auto"},
		"mcp: not running; bind auto (default; flai serve mcp bind <address> [<address>...] sets others)\n")
	if _, errOut, code := run("serve", "mcp", "bind", "0.0.0.0"); code != 0 {
		t.Fatalf("bind: %d %s", code, errOut)
	}
	check("0.0.0.0", true, []any{"0.0.0.0"}, "mcp: not running; bind 0.0.0.0 (set; flai serve mcp clear goes back to the default, auto)\n")
	if _, errOut, code := run("serve", "mcp", "bind", "10.0.0.5", "auto"); code != 0 {
		t.Fatalf("bind two: %d %s", code, errOut)
	}
	check("10.0.0.5", true, []any{"10.0.0.5", "auto"},
		"mcp: not running; bind 10.0.0.5 auto (set; flai serve mcp clear goes back to the default, auto)\n")
	if _, errOut, code := runIn(t, t.TempDir(), "serve", "settings"); code == 0 || !strings.Contains(errOut, "run flai new or flai import first") {
		t.Errorf("outside a project it says what to do: %d %s", code, errOut)
	}
}

// S-0759, ADR-0295 point 7: flai serve settings shows the project's image
// build command and whether flai serve makes its images by itself, as text
// and in the agent block of its JSON.
func TestServeSettingsShowImageSettings(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	settings := func(args ...string) (string, string, int) {
		t.Helper()
		return runWithApp(t, &app{configPath: configOf(t), cwd: root, runner: mcpDocker{engine: "Ubuntu 24.04 LTS"}}, append([]string{"serve", "settings"}, args...)...)
	}
	imageOf := func() map[string]any {
		t.Helper()
		js, errOut, code := settings("--json")
		var got struct {
			Agent struct {
				Image map[string]any `json:"image"`
			} `json:"agent"`
		}
		if err := json.Unmarshal([]byte(js), &got); err != nil || code != 0 {
			t.Fatalf("--json: %d %v %s %s", code, err, js, errOut)
		}
		return got.Agent.Image
	}
	if out, errOut, code := settings(); code != 0 || !strings.Contains(out, "  image: command none, flai's own pull and build, auto on\n") {
		t.Errorf("text, nothing set: %d %q %s", code, out, errOut)
	}
	if got := imageOf(); !reflect.DeepEqual(got, map[string]any{"command": []any{}, "auto": true, "auto_set": false}) {
		t.Errorf("--json, nothing set: %v", got)
	}
	runIn(t, root, "serve", "agent", "image", "command", "--", "scripts/agent-image.sh", "--layer")
	runIn(t, root, "serve", "agent", "image", "auto", "off")
	if out, errOut, code := settings(); code != 0 || !strings.Contains(out, "  image: command scripts/agent-image.sh --layer, auto off\n") {
		t.Errorf("text: %d %q %s", code, out, errOut)
	}
	if got := imageOf(); !reflect.DeepEqual(got, map[string]any{"command": []any{"scripts/agent-image.sh", "--layer"}, "auto": false, "auto_set": true}) {
		t.Errorf("--json: %v", got)
	}
}
