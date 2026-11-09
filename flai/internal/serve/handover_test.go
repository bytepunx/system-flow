package serve

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// TestAHandedOverAgentIsReaped reproduces I-0299: a launcher that hands its
// agent over, started from a parent that outlives the agent as flai mcp over
// stdio does, released the process without waiting for it, so the agent
// stayed a zombie of that parent once it ended. It is reaped now, and its end
// is still left to the serving flai.
func TestAHandedOverAgentIsReaped(t *testing.T) {
	if _, err := os.Stat("/proc/self/stat"); err != nil {
		t.Skip("needs /proc to tell a reaped process from a zombie")
	}
	out, err := os.Create(filepath.Join(t.TempDir(), "agent.log"))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("true")
	cmd.Stdout, cmd.Stderr = out, out
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid
	l := &launcher{handOver: true}
	ended := make(chan int, 1)
	l.await(cmd, out, pid, func(code int) { ended <- code })

	proc := "/proc/" + strconv.Itoa(pid)
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(proc); os.IsNotExist(err) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("the handed-over agent, pid %d, was not reaped within 10s: %s is still there", pid, proc)
		}
		time.Sleep(10 * time.Millisecond)
	}
	select {
	case code := <-ended:
		t.Errorf("a handed-over agent's end was recorded here, with code %d, rather than left to the serving flai", code)
	case <-time.After(50 * time.Millisecond):
	}
}
