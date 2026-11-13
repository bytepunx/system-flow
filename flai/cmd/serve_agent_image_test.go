package cmd

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/config"
)

// S-0759, ADR-0295 point 7: flai serve agent image command names the
// project's build command, prints it bare, and --clear brings flai's own
// pull and build back; auto off and on turn the automatic image job off and
// back to its default; each writes the project of the working directory or
// of --project and leaves another project as it was; and outside a project,
// with no --project, it refuses, naming the flag.
func TestServeAgentImageSettings(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	a, b := tempProject(t), tempProject(t)
	imageSettingOf := func(root string) config.ProjectImage {
		t.Helper()
		cfg, _, err := (&app{configPath: configOf(t)}).loadConfig()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Agent.Project(root).Image
	}
	run := func(dir string, args ...string) (string, string, int) {
		t.Helper()
		return runIn(t, dir, append([]string{"serve", "agent", "image"}, args...)...)
	}
	if out, errOut, code := run(a); code != 0 || !strings.Contains(out, "image command: none, so flai's own pull and build make the agent images") || !strings.Contains(out, "image auto: on (default)") {
		t.Errorf("bare, nothing set: %d %q %s", code, out, errOut)
	}
	out, errOut, code := run(a, "command", "--", "scripts/agent-image.sh", "--layer", "{x}")
	if code != 0 || !strings.Contains(out, `image command: "scripts/agent-image.sh" "--layer" "{x}" (this project)`) || !strings.Contains(out, "FLAI_IMAGE_*") || strings.Contains(out, "image auto") {
		t.Errorf("command set: %d %q %s", code, out, errOut)
	}
	if got := imageSettingOf(a); !reflect.DeepEqual(got.Command, []string{"scripts/agent-image.sh", "--layer", "{x}"}) || got.Auto != nil {
		t.Errorf("kept as agent.projects.<root>.image.command: %+v", got)
	}
	if out, _, code := run(a, "command"); code != 0 || !strings.Contains(out, `image command: "scripts/agent-image.sh" "--layer" "{x}"`) {
		t.Errorf("command bare prints it: %d %q", code, out)
	}
	if _, errOut, code := run(a, "command", "--clear", "--", "x"); code == 0 || !strings.Contains(errOut, "not both") {
		t.Errorf("a program and --clear: %d %s", code, errOut)
	}
	if _, errOut, code := run(a, "command", "--", " "); code == 0 || !strings.Contains(errOut, "the program is empty") {
		t.Errorf("an empty program: %d %s", code, errOut)
	}
	if got := imageSettingOf(a); len(got.Command) != 3 {
		t.Errorf("a refusal changes nothing: %+v", got)
	}
	if got := imageSettingOf(b); !got.IsZero() {
		t.Errorf("project b is unchanged: %+v", got)
	}
	if out, errOut, code := run(a, "auto", "off"); code != 0 || !strings.Contains(out, "image auto: off, so a start refuses a missing or stale agent image; flai image ensure") || strings.Contains(out, "image command") {
		t.Errorf("auto off: %d %q %s", code, out, errOut)
	}
	if got := imageSettingOf(a); got.Auto == nil || *got.Auto || len(got.Command) != 3 {
		t.Errorf("auto off is kept as false beside the command: %+v", got)
	}
	if out, _, code := run(a, "auto"); code != 0 || !strings.Contains(out, "image auto: off") {
		t.Errorf("auto bare prints it: %d %q", code, out)
	}
	if _, errOut, code := run(a, "auto", "maybe"); code == 0 || !strings.Contains(errOut, "neither on nor off") {
		t.Errorf("auto with another word: %d %s", code, errOut)
	}
	js, errOut, code := run(a, "--json")
	var shown map[string]any
	if err := json.Unmarshal([]byte(js), &shown); err != nil || code != 0 || shown["auto"] != false || shown["auto_set"] != true || !reflect.DeepEqual(shown["command"], []any{"scripts/agent-image.sh", "--layer", "{x}"}) || shown["project"] == nil {
		t.Errorf("--json: %d %v %s %s", code, err, js, errOut)
	}
	if out, errOut, code := run(a, "auto", "on"); code != 0 || !strings.Contains(out, "image auto: on (default)") || imageSettingOf(a).Auto != nil {
		t.Errorf("auto on removes the key: %d %q %s %+v", code, out, errOut, imageSettingOf(a))
	}
	if out, errOut, code := run(a, "command", "--clear"); code != 0 || !strings.Contains(out, "image command: none, so flai's own pull and build") {
		t.Errorf("--clear: %d %q %s", code, out, errOut)
	}
	if cfg, _, _ := (&app{configPath: configOf(t)}).loadConfig(); cfg.Agent.Projects != nil {
		t.Errorf("a's entry, now empty, is removed: %+v", cfg.Agent.Projects)
	}
	// --project from b's directory writes a's, after the sub-command or
	// before it, and b stays as it was
	if _, errOut, code := run(b, "command", "--project", a, "--", "build-a"); code != 0 {
		t.Fatalf("--project: %s", errOut)
	}
	if _, errOut, code := run(b, "--project", a, "auto", "off"); code != 0 {
		t.Fatalf("--project before the sub-command: %s", errOut)
	}
	if got := imageSettingOf(a); !reflect.DeepEqual(got.Command, []string{"build-a"}) || got.Auto == nil || *got.Auto {
		t.Errorf("--project names the project written: %+v", got)
	}
	if got := imageSettingOf(b); !got.IsZero() {
		t.Errorf("project b is unchanged: %+v", got)
	}
	if out, _, _ := run(b); !strings.Contains(out, "image command: none") || !strings.Contains(out, "image auto: on (default)") {
		t.Errorf("b's own are shown in b: %q", out)
	}
	// outside a project, with no --project, every form refuses, naming the flag
	outside := t.TempDir()
	for _, args := range [][]string{{}, {"command"}, {"command", "--", "x"}, {"command", "--clear"}, {"auto"}, {"auto", "off"}} {
		if _, errOut, code := run(outside, args...); code == 0 || !strings.Contains(errOut, "--project <root>") || strings.Contains(errOut, "--all-projects") {
			t.Errorf("%q outside a project: %d %s", args, code, errOut)
		}
	}
	if _, errOut, code := run(a, "auto", "off", "--all-projects"); code == 0 || !strings.Contains(errOut, "all-projects") {
		t.Errorf("--all-projects is no flag of it: %d %s", code, errOut)
	}
}

// S-0759: flai serve agent show prints the project's image settings, as text
// and as JSON, and outside a project prints none.
func TestServeAgentShowImageSettings(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	if out, errOut, code := runIn(t, root, "serve", "agent", "show"); code != 0 || !strings.Contains(out, "image command: none") || !strings.Contains(out, "image auto: on (default)") {
		t.Errorf("show, nothing set: %d %q %s", code, out, errOut)
	}
	runIn(t, root, "serve", "agent", "image", "command", "--", "scripts/agent-image.sh")
	runIn(t, root, "serve", "agent", "image", "auto", "off")
	if out, errOut, code := runIn(t, root, "serve", "agent", "show"); code != 0 || !strings.Contains(out, `image command: "scripts/agent-image.sh" (this project)`) || !strings.Contains(out, "image auto: off") {
		t.Errorf("show: %d %q %s", code, out, errOut)
	}
	js, errOut, code := runIn(t, root, "serve", "agent", "show", "--json")
	var shown struct {
		Image struct {
			Command []string `json:"command"`
			Auto    *bool    `json:"auto"`
			AutoSet bool     `json:"auto_set"`
		} `json:"image"`
	}
	if err := json.Unmarshal([]byte(js), &shown); err != nil || code != 0 || !reflect.DeepEqual(shown.Image.Command, []string{"scripts/agent-image.sh"}) || shown.Image.Auto == nil || *shown.Image.Auto || !shown.Image.AutoSet {
		t.Errorf("show --json: %d %v %s %s", code, err, js, errOut)
	}
	if out, _, code := runIn(t, t.TempDir(), "serve", "agent", "show"); code != 0 || strings.Contains(out, "image command") || strings.Contains(out, "image auto") {
		t.Errorf("outside a project show has no image settings: %d %q", code, out)
	}
}
