//go:build !windows

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// curatorRun is the curator run recorded in the serve state under cfg, or
// nil when there is none.
func curatorRun(t *testing.T, cfg string) *serve.AgentRun {
	t.Helper()
	for _, st := range serve.DirFor(cfg).AgentStates() {
		if st.Curator != nil {
			return st.Curator
		}
	}
	return nil
}

// endCurator ends the curator run recorded under cfg, a child of the test
// process that the start handed over, and reaps it, so that it no longer
// runs.
func endCurator(t *testing.T, cfg string) {
	t.Helper()
	run := curatorRun(t, cfg)
	if run == nil || run.PID <= 0 {
		t.Fatal("no curator run to end")
	}
	_ = syscall.Kill(run.PID, syscall.SIGKILL)
	var ws syscall.WaitStatus
	_, _ = syscall.Wait4(run.PID, &ws, 0, nil)
}

// S-0431: flai curate starts the curator for the project while the curate
// action is on, with a since or none, records the run where flai serve
// tracks agents, and says why when it will not; curate.run does the same for
// a dashboard.
func TestCurateStartsTheCuratorForTheProject(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	inProcess(t)
	t.Cleanup(func() {
		if run := curatorRun(t, cfg); run != nil && run.PID > 0 {
			_ = syscall.Kill(run.PID, syscall.SIGKILL)
		}
	})
	if _, errOut, code := runIn(t, root, "curate"); code == 0 || !strings.Contains(errOut, "rule: the curate host action is off for this project: flai serve enable curate") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	out, _, code := runIn(t, root, "hostapi", "curate.run", `{"since":"2026-10-01","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a50"}`)
	if code == 0 || !strings.Contains(out, `the host action \"curate\" is not enabled`) || !strings.Contains(out, "flai serve enable curate") {
		t.Errorf("curate.run, action off: %d %s", code, out)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "curate"); code != 0 || !strings.Contains(out, "curate enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	badSince := "the curator takes a since as a YYYY-MM-DD day in UTC or an RFC 3339 time, such as 2026-10-01 or 2026-10-01T09:00:00Z, and "
	if _, errOut, code := runIn(t, root, "curate", "--since", "last week"); code == 0 || !strings.Contains(errOut, "rule: "+badSince) || !strings.Contains(errOut, "last week") {
		t.Errorf("a bad since: %d %s", code, errOut)
	}
	if out, _, code := runIn(t, root, "hostapi", "curate.run", `{"since":"last week","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a51"}`); code == 0 || !strings.Contains(out, badSince) {
		t.Errorf("curate.run, a bad since: %d %s", code, out)
	}
	if _, errOut, code := runIn(t, root, "curate"); code == 0 || !strings.Contains(errOut, "rule: the curator's agent names no harness") {
		t.Errorf("nothing to start it with: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "agent", "set", "--", "sleep", "30")

	out, errOut, code := runIn(t, root, "curate", "--since", "2026-10-01")
	if code != 0 || !strings.HasPrefix(out, "started sleep (command) to curate, since 2026-10-01, as curator (pid ") || !strings.Contains(out, "; log ") {
		t.Fatalf("a since: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "flai serve is not running") {
		t.Errorf("not told that no flai serve settles the run: %s", errOut)
	}
	if _, errOut, code := runIn(t, root, "curate"); code == 0 || !strings.Contains(errOut, "rule: the curator is already running for this project (since 2026-10-01, started by asked, pid ") {
		t.Errorf("a run already going: %d %s", code, errOut)
	}
	endCurator(t, cfg)

	out, _, code = runIn(t, root, "curate", "--json")
	var said map[string]any
	if err := json.Unmarshal([]byte(out), &said); code != 0 || err != nil {
		t.Fatalf("no since, --json: %d %v %s", code, err, out)
	}
	for k, want := range map[string]any{"since": "", "agent": "curator", "harness": "command", "command": "sleep", "trigger": "asked"} {
		if said[k] != want {
			t.Errorf("--json %s = %v, want %v: %s", k, said[k], want, out)
		}
	}
	for _, k := range []string{"pid", "log", "session", "started"} {
		if v, ok := said[k]; !ok || v == "" || v == float64(0) {
			t.Errorf("--json has no %s: %s", k, out)
		}
	}
	endCurator(t, cfg)

	out, _, code = runIn(t, root, "hostapi", "curate.run", `{"since":"2026-10-02T09:00:00Z","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a52"}`)
	var started struct {
		Data struct {
			Since, Agent, Log, Session, Trigger string
			PID                                 int
		}
	}
	if err := json.Unmarshal([]byte(out), &started); code != 0 || err != nil || started.Data.Since != "2026-10-02T09:00:00Z" || started.Data.Agent != "curator" || started.Data.PID == 0 || started.Data.Log == "" || started.Data.Session == "" || started.Data.Trigger != "asked" {
		t.Fatalf("curate.run: %d %s", code, out)
	}
	if run := curatorRun(t, cfg); run == nil || run.Since != "2026-10-02T09:00:00Z" || run.PID != started.Data.PID {
		t.Errorf("the newest curator run recorded = %+v, want curate.run's", run)
	}
	endCurator(t, cfg)

	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"outcome": "disabled"`, `"method": "serve.curate"`, `"action": "curate"`, "to curate, since 2026-10-01, as curator",
		"to curate, everything, its first run, as curator", "started sleep to curate, since 2026-10-02T09:00:00Z, as curator"} {
		if !strings.Contains(js, want) {
			t.Errorf("the journal has no %q: %s", want, js)
		}
	}
}

// S-0431: flai mcp's curate starts the curator as flai curate does, under
// the same curate host action, journals the agent that asked, and refuses an
// agent flai serve started.
func TestTheMCPServerStartsTheCuratorUnderTheCurateAction(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	t.Setenv("FLAI_STARTED_BY", "") // the operator's own agent, whoever runs the tests
	root := tempProject(t)
	t.Cleanup(func() {
		if run := curatorRun(t, cfg); run != nil && run.PID > 0 {
			_ = syscall.Kill(run.PID, syscall.SIGKILL)
		}
	})
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root, clock: func() time.Time { return time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC) }}
	ctx := context.Background()
	if _, err := a.mcpCurate(ctx, root, "2026-10-01", "agent-ops"); err == nil || !strings.Contains(err.Error(), "the curate host action is off for this project: flai serve enable curate") {
		t.Errorf("action off: %v", err)
	}
	runIn(t, root, "serve", "enable", "curate")
	runIn(t, root, "serve", "agent", "set", "--", "sleep", "30")
	if _, err := a.mcpCurate(ctx, root, "soon", "agent-ops"); err == nil || !strings.Contains(err.Error(), `"soon" is neither`) {
		t.Errorf("a bad since: %v", err)
	}
	got, err := a.mcpCurate(ctx, root, "2026-10-01", "agent-ops")
	if err != nil || got.Since != "2026-10-01" || got.Agent != "curator" || got.PID == 0 || got.Command != "sleep" || got.Log == "" || got.Session == "" || got.Trigger != "asked" {
		t.Fatalf("a since: %+v %v", got, err)
	}
	if _, err := a.mcpCurate(ctx, root, "", "agent-ops"); err == nil || !strings.Contains(err.Error(), "the curator is already running for this project") {
		t.Errorf("a run already going: %v", err)
	}
	endCurator(t, cfg)
	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	if _, err := a.mcpCurate(ctx, root, "", "agent-S-0001"); err == nil || !strings.Contains(err.Error(), "an agent flai serve started does not start the curator") {
		t.Errorf("an agent flai serve started: %v", err)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "mcp.curate"`, `"by": "agent-ops"`, `"outcome": "disabled"`, "agent-ops asked to curate, since 2026-10-01: started sleep as curator",
		"agent-ops asked to curate, since soon: the curator takes a since", "agent-ops asked to curate, since its last run: the curator is already running"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
}
