package serve

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// statLine is a /proc/<pid>/stat line whose command name holds a space and a
// parenthesis, in state state, session 4242, and started at tick 987654.
func statLine(state string) []byte {
	fields := []string{"4242 (flai (mcp) x)", state, "1", "4242", "4242"}
	for len(fields) < 20 {
		fields = append(fields, "0")
	}
	fields = append(fields, "987654", "0", "0")
	return []byte(strings.Join(fields, " ") + "\n")
}

func TestParseStatReadsSessionStartAndState(t *testing.T) {
	for state, want := range map[string]bool{"R": false, "S": false, "D": false, "T": false, "Z": true, "X": true} {
		sid, start, ended, err := parseStat("fixture", statLine(state))
		if err != nil {
			t.Fatalf("state %s: %v", state, err)
		}
		if sid != 4242 || start != 987654 || ended != want {
			t.Errorf("state %s: got session %d, start %d, ended %v; want 4242, 987654, %v", state, sid, start, ended, want)
		}
	}
}

func TestParseStatRefusesAShortOrNamelessLine(t *testing.T) {
	for name, line := range map[string]string{"no name": "4242 S 1 4242", "too few fields": "4242 (x) S 1 4242 4242"} {
		if _, _, _, err := parseStat("fixture", []byte(line)); err == nil {
			t.Errorf("%s: parsed %q without an error", name, line)
		}
	}
}

// TestAZombieIsNeitherAliveNorOwned reproduces I-0299: a child that has
// exited and that its parent has not reaped still answers signal 0, and was
// counted as running.
func TestAZombieIsNeitherAliveNorOwned(t *testing.T) {
	cmd := exec.Command("true")
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid
	defer func() { _ = cmd.Wait() }()
	start := Started(pid)
	deadline := time.Now().Add(10 * time.Second)
	for {
		stat, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
		if err != nil {
			t.Fatal(err)
		}
		if _, _, ended, err := parseStat("stat", stat); err == nil && ended {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("pid %d did not become a zombie within 10s", pid)
		}
		time.Sleep(10 * time.Millisecond)
	}
	if Alive(pid) {
		t.Error("Alive counts an unreaped zombie as alive")
	}
	if Owns(pid, start) {
		t.Error("Owns counts an unreaped zombie as the process it started")
	}
}

func TestALiveDetachedProcessIsAliveAndOwned(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()
	if !Alive(pid) || !Owns(pid, Started(pid)) {
		t.Errorf("a running detached process is not alive and owned: Alive %v, Owns %v", Alive(pid), Owns(pid, Started(pid)))
	}
}
