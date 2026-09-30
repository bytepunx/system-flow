//go:build !windows

package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0170: flai serve agent stop ends a story's running agent, journals it,
// and says why when there is none to stop.
func TestServeAgentStopEndsAStorysAgent(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	if _, errOut, code := runIn(t, root, "serve", "agent", "stop", "S-1"); code == 0 || !strings.Contains(errOut, "rule: flai serve has started no agent for S-0001") {
		t.Errorf("no agent: %d %s", code, errOut)
	}
	// an agent flai serve started, detached, and reaped by whoever started it
	agent := exec.Command("sleep", "30")
	serve.Detach(agent)
	if err := agent.Start(); err != nil {
		t.Fatal(err)
	}
	waited := make(chan struct{})
	go func() { _ = agent.Wait(); close(waited) }()
	t.Cleanup(func() { _ = agent.Process.Kill(); <-waited })
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	dir := serve.DirFor(cfg)
	_ = os.MkdirAll(string(dir), 0o700)
	state := map[string]serve.AgentState{mainRootOf(repo): {Stories: map[string]*serve.AgentRun{
		"S-0001": {Story: "S-0001", Agent: "agent-S-0001", Command: "sleep", PID: agent.Process.Pid, Started: time.Now().UTC().Format(time.RFC3339)},
	}}}
	data, _ := json.Marshal(state)
	if err := os.WriteFile(filepath.Join(string(dir), "agents.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "stop", "S-1")
	if code != 0 || out != "stopped S-0001's agent agent-S-0001 (pid "+strconv.Itoa(agent.Process.Pid)+")\n" {
		t.Fatalf("stop: %d %q %s", code, out, errOut)
	}
	select {
	case <-waited:
	case <-time.After(5 * time.Second):
		t.Error("the agent's process still runs")
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	if !strings.Contains(js, "stopped S-0001's agent agent-S-0001: ended its process group") {
		t.Errorf("the stop is not journalled: %s", js)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "stop", "S-1"); code == 0 || !strings.Contains(errOut, "rule: S-0001's agent is not running: it ended") || !strings.Contains(errOut, "stopped by the operator") {
		t.Errorf("twice: %d %s", code, errOut)
	}
}
