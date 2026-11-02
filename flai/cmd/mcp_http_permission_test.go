//go:build !windows

package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// ADR-0243: the server over HTTP reads the auto-approve host action as flai
// mcp over stdio does, and narrows it to the agent flai serve records as
// running for the story: by its name, for that story, in that project, while
// the process it started is the one alive.
func TestMCPOverHTTPAutoApprovesOnlyTheAgentFlaiServeRunsForTheStory(t *testing.T) {
	t.Parallel()
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	a := &app{configPath: cfg, out: &bytes.Buffer{}, errOut: &bytes.Buffer{}}
	if a.mcpServedAgent("/p", "agent-S-0001", "S-0001") {
		t.Error("with no state flai serve wrote, no agent is flai serve's")
	}
	opt := a.mcpHTTPOptions(nil, "agent-S-0001", nil, "")
	if opt.AutoApprove == nil || opt.ServedAgent == nil {
		t.Fatal("over HTTP the server reads the auto-approve host action and narrows it to the agents flai serve started")
	}

	agent := exec.Command("sleep", "30")
	serve.Detach(agent) // the leader of its own session, as flai serve starts an agent
	if err := agent.Start(); err != nil {
		t.Fatal(err)
	}
	waited := make(chan struct{})
	go func() { _ = agent.Wait(); close(waited) }()
	t.Cleanup(func() { _ = agent.Process.Kill(); <-waited })
	now := time.Now().UTC().Format(time.RFC3339)
	pid := agent.Process.Pid
	state := map[string]serve.AgentState{"/p": {Stories: map[string]*serve.AgentRun{
		"S-0001": {Story: "S-0001", Agent: "agent-S-0001", PID: pid, Started: now},
		"S-0002": {Story: "S-0002", Agent: "agent-S-0002", PID: pid, Started: now, Ended: now},
		"S-0003": {Story: "S-0003", Agent: "agent-S-0003", Started: now, Error: "no harness"},
	}}}
	dir := serve.DirFor(cfg)
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(state)
	if err := os.WriteFile(filepath.Join(string(dir), "agents.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		why                string
		root, agent, story string
		want               bool
	}{
		{"the agent flai serve runs for the story", "/p", "agent-S-0001", "S-0001", true},
		{"another name, claimed by any holder of the token", "/p", "claude", "S-0001", false},
		{"the agent of another story", "/p", "agent-S-0001", "S-0002", false},
		{"a run that ended", "/p", "agent-S-0002", "S-0002", false},
		{"a run that could not be started", "/p", "agent-S-0003", "S-0003", false},
		{"another project", "/q", "agent-S-0001", "S-0001", false},
	} {
		if got := a.mcpServedAgent(c.root, c.agent, c.story); got != c.want {
			t.Errorf("%s: flai serve's is %v, want %v", c.why, got, c.want)
		}
	}
	_ = agent.Process.Kill()
	<-waited
	if a.mcpServedAgent("/p", "agent-S-0001", "S-0001") {
		t.Error("once the process flai serve started has gone, its name is no longer flai serve's agent")
	}
}
