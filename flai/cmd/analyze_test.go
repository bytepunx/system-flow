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

// analyzerRun is the analyzer run recorded in the serve state under cfg, or
// nil when there is none.
func analyzerRun(t *testing.T, cfg string) *serve.AgentRun {
	t.Helper()
	for _, st := range serve.DirFor(cfg).AgentStates() {
		if st.Analyzer != nil {
			return st.Analyzer
		}
	}
	return nil
}

// endAnalyzer ends the analyzer run recorded under cfg, a child of the test
// process that the start handed over, and reaps it, so that it no longer
// runs.
func endAnalyzer(t *testing.T, cfg string) {
	t.Helper()
	run := analyzerRun(t, cfg)
	if run == nil || run.PID <= 0 {
		t.Fatal("no analyzer run to end")
	}
	_ = syscall.Kill(run.PID, syscall.SIGKILL)
	var ws syscall.WaitStatus
	_, _ = syscall.Wait4(run.PID, &ws, 0, nil)
}

// S-0223: flai analyze starts the analyzer for the project while the analyze
// action is on, with a focus or none, records the run where flai serve
// tracks agents, and says why when it will not; analyze.run does the same
// for a dashboard.
func TestAnalyzeStartsTheAnalyzerForTheProject(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	inProcess(t)
	t.Cleanup(func() {
		if run := analyzerRun(t, cfg); run != nil && run.PID > 0 {
			_ = syscall.Kill(run.PID, syscall.SIGKILL)
		}
	})
	if _, errOut, code := runIn(t, root, "analyze"); code == 0 || !strings.Contains(errOut, "rule: the analyze host action is off for this project: flai serve enable analyze") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	out, _, code := runIn(t, root, "hostapi", "analyze.run", `{"focus":"risk","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a40"}`)
	if code == 0 || !strings.Contains(out, `the host action \"analyze\" is not enabled`) || !strings.Contains(out, "flai serve enable analyze") {
		t.Errorf("analyze.run, action off: %d %s", code, out)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "analyze"); code != 0 || !strings.Contains(out, "analyze enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	// standard error is JSON here, its quotes escaped
	badFocus := `the analyzer takes the focus bottlenecks, intent, risk, or none for all of them, and \"style\" is none of them`
	if _, errOut, code := runIn(t, root, "analyze", "--focus", "style"); code == 0 || !strings.Contains(errOut, "rule: "+badFocus) {
		t.Errorf("a bad focus: %d %s", code, errOut)
	}
	if out, _, code := runIn(t, root, "hostapi", "analyze.run", `{"focus":"style","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a41"}`); code == 0 || !strings.Contains(out, badFocus) {
		t.Errorf("analyze.run, a bad focus: %d %s", code, out)
	}
	if _, errOut, code := runIn(t, root, "analyze"); code == 0 || !strings.Contains(errOut, "rule: the analyzer's agent names no harness") {
		t.Errorf("nothing to start it with: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "agent", "set", "--", "sleep", "30")

	out, errOut, code := runIn(t, root, "analyze", "--focus", "risk")
	if code != 0 || !strings.HasPrefix(out, "started sleep (command) to analyze, focus risk, as analyzer (pid ") || !strings.Contains(out, "; log ") {
		t.Fatalf("a focus: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "flai serve is not running") {
		t.Errorf("not told that no flai serve settles the run: %s", errOut)
	}
	if _, errOut, code := runIn(t, root, "analyze"); code == 0 || !strings.Contains(errOut, "rule: the analyzer is already running for this project (focus risk, started by asked, pid ") {
		t.Errorf("a run already going: %d %s", code, errOut)
	}
	endAnalyzer(t, cfg)

	out, _, code = runIn(t, root, "analyze", "--json")
	var said map[string]any
	if err := json.Unmarshal([]byte(out), &said); code != 0 || err != nil {
		t.Fatalf("no focus, --json: %d %v %s", code, err, out)
	}
	for k, want := range map[string]any{"focus": "all", "agent": "analyzer", "harness": "command", "command": "sleep", "trigger": "asked"} {
		if said[k] != want {
			t.Errorf("--json %s = %v, want %v: %s", k, said[k], want, out)
		}
	}
	for _, k := range []string{"pid", "log", "session", "started"} {
		if v, ok := said[k]; !ok || v == "" || v == float64(0) {
			t.Errorf("--json has no %s: %s", k, out)
		}
	}
	endAnalyzer(t, cfg)

	out, _, code = runIn(t, root, "hostapi", "analyze.run", `{"focus":"intent","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a42"}`)
	var started struct {
		Data struct {
			Focus, Agent, Log, Session, Trigger string
			PID                                 int
		}
	}
	if err := json.Unmarshal([]byte(out), &started); code != 0 || err != nil || started.Data.Focus != "intent" || started.Data.Agent != "analyzer" || started.Data.PID == 0 || started.Data.Log == "" || started.Data.Session == "" || started.Data.Trigger != "asked" {
		t.Fatalf("analyze.run: %d %s", code, out)
	}
	if run := analyzerRun(t, cfg); run == nil || run.Focus != "intent" || run.PID != started.Data.PID {
		t.Errorf("the newest analyzer run recorded = %+v, want analyze.run's", run)
	}
	endAnalyzer(t, cfg)

	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"outcome": "disabled"`, `"method": "serve.analyze"`, "to analyze, focus risk, as analyzer", "started sleep to analyze, focus intent, as analyzer"} {
		if !strings.Contains(js, want) {
			t.Errorf("the journal has no %q: %s", want, js)
		}
	}
}

// S-0223: flai mcp's analyze starts the analyzer as flai analyze does, under
// the same analyze host action, journals the agent that asked, and refuses an
// agent flai serve started.
func TestTheMCPServerStartsTheAnalyzerUnderTheAnalyzeAction(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	t.Setenv("FLAI_STARTED_BY", "") // the operator's own agent, whoever runs the tests
	root := tempProject(t)
	t.Cleanup(func() {
		if run := analyzerRun(t, cfg); run != nil && run.PID > 0 {
			_ = syscall.Kill(run.PID, syscall.SIGKILL)
		}
	})
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root, clock: func() time.Time { return time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC) }}
	ctx := context.Background()
	if _, err := a.mcpAnalyze(ctx, root, "risk", "agent-ops"); err == nil || !strings.Contains(err.Error(), "the analyze host action is off for this project: flai serve enable analyze") {
		t.Errorf("action off: %v", err)
	}
	runIn(t, root, "serve", "enable", "analyze")
	runIn(t, root, "serve", "agent", "set", "--", "sleep", "30")
	if _, err := a.mcpAnalyze(ctx, root, "style", "agent-ops"); err == nil || !strings.Contains(err.Error(), `"style" is none of them`) {
		t.Errorf("a bad focus: %v", err)
	}
	got, err := a.mcpAnalyze(ctx, root, "risk", "agent-ops")
	if err != nil || got.Focus != "risk" || got.Agent != "analyzer" || got.PID == 0 || got.Command != "sleep" || got.Log == "" || got.Session == "" || got.Trigger != "asked" {
		t.Fatalf("a focus: %+v %v", got, err)
	}
	if _, err := a.mcpAnalyze(ctx, root, "", "agent-ops"); err == nil || !strings.Contains(err.Error(), "the analyzer is already running for this project") {
		t.Errorf("a run already going: %v", err)
	}
	endAnalyzer(t, cfg)
	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	if _, err := a.mcpAnalyze(ctx, root, "", "agent-S-0001"); err == nil || !strings.Contains(err.Error(), "an agent flai serve started does not start the analyzer") {
		t.Errorf("an agent flai serve started: %v", err)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "mcp.analyze"`, `"by": "agent-ops"`, `"outcome": "disabled"`, "agent-ops asked to analyze, focus risk: started sleep as analyzer", `agent-ops asked to analyze, focus style: the analyzer takes the focus`, "agent-ops asked to analyze, focus all: the analyzer is already running"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
}
