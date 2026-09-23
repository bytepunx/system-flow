package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// S-0101: flai dashboard in a folder with no system-flow.yaml, anywhere up
// the tree, starts the shared dashboard and flai serve without registering a
// project, names the folder for import, and records the dashboard so that
// flai serve offers it the repositories to import with no project served.
func TestDashboardInAFolderWithNoProject(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	folder := t.TempDir()
	dir := serve.DirFor(cfg)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}

	out, errOut, code := runWith(t, folder, f, "dashboard")
	if code != 0 {
		t.Fatalf("dashboard: %d %s", code, errOut)
	}
	if !strings.Contains(out, "flaiover running at") || !strings.Contains(out, folder+" is named for import") || !strings.Contains(out, "no project here, so none is registered") {
		t.Errorf("output: %s", out)
	}
	if projects, _ := dir.Projects(); len(projects) != 0 {
		t.Errorf("registered: %+v", projects)
	}
	c, _, err := config.Load(cfg)
	if err != nil || len(c.ImportRoots) != 1 || c.ImportRoots[0] != folder {
		t.Errorf("import roots: %v %v", c.ImportRoots, err)
	}
	if dbs, _ := dir.Dashboards(); len(dbs) != 1 || dbs[0].URL == "" || dbs[0].KeyFile != filepath.Join(string(dir), "dashboard.agent-key") {
		t.Errorf("dashboards: %+v", dbs)
	}

	// again: already running, already named, nothing doubled
	out, _, code = runWith(t, folder, f, "dashboard")
	if code != 0 || !strings.Contains(out, "already runs") {
		t.Errorf("again: %d %s", code, out)
	}
	if c, _, _ := config.Load(cfg); len(c.ImportRoots) != 1 {
		t.Errorf("named twice: %v", c.ImportRoots)
	}

	// status and token work there too
	out, errOut, code = runWith(t, folder, f, "dashboard", "status")
	if code != 0 || !strings.Contains(out, "running at") || !strings.Contains(out, "no project here") {
		t.Errorf("status: %d %s %s", code, out, errOut)
	}
	js, _, code := runWith(t, folder, f, "dashboard", "token", "--json")
	var tok struct {
		Token string `json:"token"`
	}
	if code != 0 || json.Unmarshal([]byte(js), &tok) != nil || tok.Token == "" {
		t.Errorf("token: %d %s", code, js)
	}

	// stopped with no project left: the container stops and the dashboard is forgotten
	out, errOut, code = runWith(t, folder, f, "dashboard", "stop")
	if code != 0 || !strings.Contains(out, "stopped flaiover") {
		t.Errorf("stop: %d %s %s", code, out, errOut)
	}
	if dbs, _ := dir.Dashboards(); len(dbs) != 0 {
		t.Errorf("still recorded: %+v", dbs)
	}
}

// A manifest that does not parse is still an error, not "no project here".
func TestDashboardStillRefusesABrokenManifest(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	folder := t.TempDir()
	if err := os.WriteFile(filepath.Join(folder, "system-flow.yaml"), []byte("version: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, folder, f, "dashboard"); code == 0 || !strings.Contains(errOut, "system-flow.yaml") {
		t.Errorf("broken manifest: %d %s", code, errOut)
	}
}
