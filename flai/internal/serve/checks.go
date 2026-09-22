package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Running checks for a story in review (S-0082, ADR-0029): the operator's
// named commands (ResolveChecks), run in order in the story's worktree,
// stopping at the first failure. One run per story at a time. A run's
// process, and everything it spawned, is reachable to cancel: Detach gives
// each one its own process group, so TerminateGroup and, if it does not
// answer, KillGroup reach the whole tree, not only the direct child.

// CheckStep is one named command's own record within a run.
type CheckStep struct {
	Name    string `json:"name"`
	Command string `json:"command"` // argv[0]'s base name, not the full command
	Started string `json:"started,omitempty"`
	Ended   string `json:"ended,omitempty"`
	Exit    *int   `json:"exit,omitempty"`
}

// ChecksRun is one run of every configured check for one story.
type ChecksRun struct {
	Story   string      `json:"story"`
	Root    string      `json:"root"`
	Running bool        `json:"running"`
	PID     int         `json:"pid,omitempty"`     // of the step now running
	Current string      `json:"current,omitempty"` // its name
	Steps   []CheckStep `json:"steps"`
	// Outcome is passed, failed, timed-out, or cancelled, once Running is false.
	Outcome string `json:"outcome,omitempty"`
	Started string `json:"started"`
	Ended   string `json:"ended,omitempty"`
	Log     string `json:"log,omitempty"`
}

func (d Dir) checksDir() string                { return filepath.Join(string(d), "checks") }
func (d Dir) checksState(story string) string  { return filepath.Join(d.checksDir(), story+".json") }
func (d Dir) checksCancel(story string) string { return filepath.Join(d.checksDir(), story+".cancel") }

var checksMu sync.Mutex

// ChecksState reads the last known state of a story's checks, if any run
// has ever started.
func (d Dir) ChecksState(story string) (ChecksRun, bool) {
	checksMu.Lock()
	defer checksMu.Unlock()
	data, err := os.ReadFile(d.checksState(story))
	if err != nil {
		return ChecksRun{}, false
	}
	var run ChecksRun
	if json.Unmarshal(data, &run) != nil {
		return ChecksRun{}, false
	}
	return run, true
}

func (d Dir) writeChecksState(run ChecksRun) error {
	checksMu.Lock()
	defer checksMu.Unlock()
	if err := os.MkdirAll(d.checksDir(), 0o700); err != nil {
		return err
	}
	return d.write(d.checksState(run.Story), run)
}

// CancelChecks asks a story's active run to stop: it touches a sentinel
// file the run itself watches for and kills its own current step on
// (below), then waits up to grace for the state file to say it is no
// longer running. It never signals a process directly: only the run's own
// process, which alone writes this story's state, ever does that, so there
// is nothing here for two processes to race over.
func (d Dir) CancelChecks(story string, grace time.Duration, sleep func(time.Duration)) error {
	run, ok := d.ChecksState(story)
	if !ok || !run.Running {
		return fmt.Errorf("no checks run is active for %s", story)
	}
	if err := os.MkdirAll(d.checksDir(), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(d.checksCancel(story), nil, 0o600); err != nil {
		return err
	}
	deadline := time.Now().Add(grace)
	for time.Now().Before(deadline) {
		if run, ok := d.ChecksState(story); ok && !run.Running {
			return nil
		}
		sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("%s did not stop within %s; flai checks status %s says whether it has since", story, grace, story)
}

// RunChecks runs every configured command in order, in the worktree, for
// story, refusing if a run is already active for it. It blocks until the
// sequence ends, is cancelled, or cfg.Timeout passes, writing the state
// file as it goes so any process asking (flai checks status, flai checks
// tail, a concurrent cancel) sees it live, not only once this returns.
func (d Dir) RunChecks(story, root, worktree string, cfg ChecksConfig, now func() time.Time) (ChecksRun, error) {
	if existing, ok := d.ChecksState(story); ok && existing.Running && Alive(existing.PID) {
		return ChecksRun{}, fmt.Errorf("a checks run for %s is already active (pid %d); flai checks status %s says more", story, existing.PID, story)
	}
	if len(cfg.Commands) == 0 {
		return ChecksRun{}, fmt.Errorf("no checks are named for this project; flai serve checks set, or checks: in the manifest")
	}
	if _, err := os.Stat(worktree); err != nil {
		return ChecksRun{}, fmt.Errorf("%s has no worktree at %s; open one with flai stream open %s", story, worktree, story)
	}
	_ = os.Remove(d.checksCancel(story)) // a stale one from a run that ended on its own

	if err := os.MkdirAll(d.checksDir(), 0o700); err != nil {
		return ChecksRun{}, err
	}
	logPath := filepath.Join(d.checksDir(), fmt.Sprintf("%s-%s.log", story, now().UTC().Format("20060102T150405Z")))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return ChecksRun{}, err
	}
	defer func() { _ = logFile.Close() }()

	run := ChecksRun{Story: story, Root: root, Running: true, Started: now().UTC().Format(time.RFC3339), Log: logPath}
	if err := d.writeChecksState(run); err != nil {
		return ChecksRun{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	outcome := "passed"
stepLoop:
	for _, nc := range cfg.Commands {
		argv := Substitute(nc.Command, story, worktree)
		if len(argv) == 0 || strings.TrimSpace(argv[0]) == "" {
			outcome = "failed"
			run.Steps = append(run.Steps, CheckStep{Name: nc.Name, Started: now().UTC().Format(time.RFC3339), Ended: now().UTC().Format(time.RFC3339)})
			fmt.Fprintf(logFile, "=== %s: no program named ===\n", nc.Name)
			break
		}
		fmt.Fprintf(logFile, "=== %s: %s ===\n", nc.Name, strings.Join(argv, " "))
		step := CheckStep{Name: nc.Name, Command: filepath.Base(argv[0]), Started: now().UTC().Format(time.RFC3339)}
		onStart := func(pid int) {
			run.PID, run.Current = pid, nc.Name
			_ = d.writeChecksState(run)
		}
		exitCode, killedBy, startErr := runOneCheck(ctx, d.checksCancel(story), worktree, argv, logFile, onStart)
		step.Ended = now().UTC().Format(time.RFC3339)
		step.Exit = exitCode
		run.Steps = append(run.Steps, step)
		run.PID = 0
		switch {
		case startErr != nil:
			fmt.Fprintf(logFile, "=== %s: could not start: %s ===\n", nc.Name, startErr.Error())
			outcome = "failed"
			break stepLoop
		case killedBy != "":
			outcome = killedBy
			break stepLoop
		case exitCode == nil || *exitCode != 0:
			outcome = "failed"
			break stepLoop
		}
		_ = d.writeChecksState(run)
	}

	run.Running, run.Current, run.PID = false, "", 0
	run.Outcome = outcome
	run.Ended = now().UTC().Format(time.RFC3339)
	_ = os.Remove(d.checksCancel(story))
	if err := d.writeChecksState(run); err != nil {
		return run, err
	}
	return run, nil
}

// runOneCheck runs one command, in its own process group so cancelling it
// reaches anything it spawned. It ends the process (terminate, then kill)
// the moment ctx ends (a time limit) or the cancel file appears (an
// operator's cancel), and says which; otherwise it returns the exit code.
func runOneCheck(ctx context.Context, cancelFile, dir string, argv []string, out io.Writer, onStart func(pid int)) (exitCode *int, killedBy string, err error) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Stdout, cmd.Stderr = out, out
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		return nil, "", err
	}
	pid := cmd.Process.Pid
	onStart(pid)
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
watch:
	for {
		select {
		case waitErr := <-done:
			return exitCodeOf(waitErr), "", nil
		case <-ctx.Done():
			killedBy = "timed-out"
			break watch
		case <-ticker.C:
			if _, statErr := os.Stat(cancelFile); statErr == nil {
				killedBy = "cancelled"
				break watch
			}
		}
	}
	terminateThenKill(pid)
	return exitCodeOf(<-done), killedBy, nil
}

// terminateThenKill asks the process group to stop, and forces it if it
// has not gone within a short grace period.
func terminateThenKill(pid int) {
	_ = TerminateGroup(pid)
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if !Alive(pid) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	_ = KillGroup(pid)
}

func exitCodeOf(err error) *int {
	code := 0
	var xe *exec.ExitError
	switch {
	case errors.As(err, &xe):
		code = xe.ExitCode()
	case err != nil:
		code = -1
	}
	return &code
}
