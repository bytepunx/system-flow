package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// S-0184: status probes the published port and says running, not answering,
// or gone, with docker's HEALTHCHECK verdict beside it.
func TestDashboardStatusSaysRunningNotAnsweringOrGone(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true}, probeAddr: "0.0.0.0:4242", health: map[string]string{"flaiover": "healthy"}}

	var probed []string
	answers := true
	status := func(args ...string) string {
		t.Helper()
		a := &app{cwd: root, runner: f, healthProbe: func(url string) bool { probed = append(probed, url); return answers }}
		out, errOut, code := runWithApp(t, a, append([]string{"dashboard", "status"}, args...)...)
		if code != 0 {
			t.Fatalf("status: %d %s", code, errOut)
		}
		return out
	}

	out := status()
	if !strings.Contains(out, "flaiover running at") || !strings.Contains(out, "docker: healthy") {
		t.Errorf("running: %s", out)
	}
	if len(probed) != 1 || probed[0] != "http://127.0.0.1:4242/_health" {
		t.Errorf("a port published on every interface is probed on loopback: %v", probed)
	}
	var js struct {
		Running bool   `json:"running"`
		State   string `json:"state"`
		Health  string `json:"health"`
	}
	if err := json.Unmarshal([]byte(status("--json")), &js); err != nil || !js.Running || js.State != "running" || js.Health != "healthy" {
		t.Errorf("running --json: %+v %v", js, err)
	}

	answers, f.health["flaiover"] = false, "unhealthy"
	out = status()
	if !strings.Contains(out, "flaiover not answering at") || !strings.Contains(out, "docker: unhealthy") || !strings.Contains(out, "flai dashboard restart") {
		t.Errorf("not answering: %s", out)
	}
	js = struct {
		Running bool   `json:"running"`
		State   string `json:"state"`
		Health  string `json:"health"`
	}{}
	if err := json.Unmarshal([]byte(status("--json")), &js); err != nil || !js.Running || js.State != "not-answering" {
		t.Errorf("not answering --json: %+v %v", js, err)
	}

	delete(f.running, "flaiover")
	probed = nil
	out = status()
	if !strings.Contains(out, "flaiover gone") || !strings.Contains(out, "flai dashboard") {
		t.Errorf("gone: %s", out)
	}
	if len(probed) != 0 {
		t.Errorf("no container, nothing to probe: %v", probed)
	}
	js.Health = ""
	if err := json.Unmarshal([]byte(status("--json")), &js); err != nil || js.Running || js.State != "gone" || js.Health != "" {
		t.Errorf("gone --json: %+v %v", js, err)
	}
}

func TestProbeHostProbesAWildcardOnThisHost(t *testing.T) {
	for in, want := range map[string]string{
		"0.0.0.0:4242":     "127.0.0.1:4242",
		"[::]:4242":        "127.0.0.1:4242",
		"127.0.0.1:4242":   "127.0.0.1:4242",
		"192.168.1.5:4242": "192.168.1.5:4242",
		"nonsense":         "nonsense",
	} {
		if got := probeHost(in); got != want {
			t.Errorf("probeHost(%q) = %q, want %q", in, got, want)
		}
	}
}
