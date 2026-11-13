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
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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
	// AskedBy is who asked for the run: the story's agent that asked flai
	// serve, flai serve for the review pass, or empty for a run from a shell
	// or the dashboard (ADR-0292 §6).
	AskedBy string `json:"asked_by,omitempty"`
	// Namespace and Commit are the kluster namespace the run ran against and
	// the commit it was deployed at, as its variables name them; empty for a
	// run with no namespace.
	Namespace string `json:"namespace,omitempty"`
	Commit    string `json:"commit,omitempty"`
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
	return d.RunChecksWith(context.Background(), story, root, worktree, cfg, nil, now)
}

// RunChecksWith is RunChecks with env added to each check's environment,
// over flai's own, such as a kluster namespace's variables (S-0579,
// ADR-0254 §8), and cut short, as cancelled, when ctx ends. With no env each
// check gets flai's environment as it is.
func (d Dir) RunChecksWith(ctx context.Context, story, root, worktree string, cfg ChecksConfig, env []string, now func() time.Time) (ChecksRun, error) {
	return d.RunChecksIn(ctx, story, root, worktree, cfg, env, nil, now)
}

// RunChecksIn is RunChecksWith with each check run in the story's container,
// in, when in is not nil, as ChecksContainer gives it for a story whose agent
// runs on the container harness (ADR-0274, I-0295): no check's program runs on
// the host then, and a container that cannot run fails the check. With in nil
// each check runs on the host, in the worktree, as RunChecksWith runs it.
func (d Dir) RunChecksIn(ctx context.Context, story, root, worktree string, cfg ChecksConfig, env []string, in *CheckContainer, now func() time.Time) (ChecksRun, error) {
	return d.RunChecksAs(ctx, "", story, root, worktree, cfg, env, in, now)
}

// RunChecksAs is RunChecksIn for by, who asked for the run, as its AskedBy
// keeps it.
func (d Dir) RunChecksAs(ctx context.Context, by, story, root, worktree string, cfg ChecksConfig, env []string, in *CheckContainer, now func() time.Time) (ChecksRun, error) {
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

	started := now().UTC()
	run := ChecksRun{Story: story, Root: root, Running: true, Started: started.Format(time.RFC3339), Log: logPath, AskedBy: by,
		Namespace: envValue(env, "FLAI_KLUSTER_NAMESPACE"), Commit: envValue(env, "FLAI_KLUSTER_COMMIT")}
	if err := d.writeChecksState(run); err != nil {
		return ChecksRun{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	if run.Namespace != "" {
		fmt.Fprintf(logFile, "=== against the kluster namespace %s ===\n", run.Namespace)
	}

	outcome := "passed"
	if in != nil {
		fmt.Fprintf(logFile, "=== in %s's container ===\n", story)
	}
stepLoop:
	for i, nc := range cfg.Commands {
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
		var exitCode *int
		var killedBy string
		var startErr error
		if in == nil {
			exitCode, killedBy, startErr = runOneCheck(ctx, d.checksCancel(story), worktree, argv, env, logFile, onStart, nil)
		} else {
			id := fmt.Sprintf("%s-%d", started.Format("20060102T150405Z"), i+1)
			exitCode, killedBy, startErr = in.run(ctx, d.checksCancel(story), root, story, id, worktree, argv, env, logFile, onStart)
		}
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

// CheckContainer is where the checks of a story whose agent runs on the
// container harness run (ADR-0274): the container harness's host and the
// story's image, or Err, why its container cannot run, which fails each check
// rather than run it on the host. Docker stops a check's container when the
// check is cancelled or runs out of time; nil is the system's.
type CheckContainer struct {
	Host   harness.Host
	Image  string
	Err    error
	Docker execx.Runner
}

// ChecksContainer is where story's checks run when its agent, resolved as flai
// serve starts it (ADR-0251, ADR-0261), runs on the container harness: its
// image, checked through docker as an agent's start checks it, or why it
// cannot run there. It is nil for a story whose agent runs on any other
// harness, or that cannot be read, whose checks run on the host as before.
func ChecksContainer(repo *workitem.Repo, story string, cfg AgentConfig, docker execx.Runner) *CheckContainer {
	it, err := repo.Get(story)
	if err != nil {
		return nil
	}
	agent := cfg.storyAgent(repo, it)
	if agent == nil || agent.Harness != harness.Container {
		return nil
	}
	host, image, err := containerImage(docker, mainCheckout(repo), cfg, agent)
	if err == nil && image == "" {
		err = fmt.Errorf("the container harness has no image for the inner harness %s of %s's agent: name one with flai serve agent harness container --image %s=<image>", innerOf(agent), story, innerOf(agent))
	}
	if err != nil {
		return &CheckContainer{Err: fmt.Errorf("%s's agent runs in a container, so its checks run there and nothing runs on the host: %w", story, err), Docker: docker}
	}
	return &CheckContainer{Host: host, Image: image, Docker: docker}
}

// run runs one check, argv, in story's container named for id, with env's
// names passed to it and their values in the docker client's environment, as
// runOneCheck runs it on the host: the docker client in its own process
// group, ended when ctx ends or the cancel file appears, after docker stop
// has stopped its container by name (ADR-0058).
func (c *CheckContainer) run(ctx context.Context, cancelFile, root, story, id, worktree string, argv, env []string, out io.Writer, onStart func(pid int)) (*int, string, error) {
	if c.Err != nil {
		return nil, "", c.Err
	}
	var names []string
	for _, kv := range env {
		if name, _, _ := strings.Cut(kv, "="); !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	cc, err := harness.ContainerCommand(harness.CommandRun{Root: root, Story: story, Purpose: "check", ID: id, Image: c.Image, Env: names, Argv: argv}, c.Host)
	if err != nil {
		return nil, "", err
	}
	docker := c.Docker
	if docker == nil {
		docker = execx.System{}
	}
	stop := func() { _, _ = docker.Run(worktree, cc.Argv[0], "stop", cc.Name) }
	return runOneCheck(ctx, cancelFile, worktree, cc.Argv, append(slices.Clone(env), cc.Env...), out, onStart, stop)
}

// runOneCheck runs one command, with env over flai's environment when there
// is any, in its own process group so cancelling it reaches anything it
// spawned. It ends the process (terminate, then kill) the moment ctx ends (a
// time limit, or flai serve stopping) or the cancel file appears (an
// operator's cancel), and says which, running stop first when there is one;
// otherwise it returns the exit code.
func runOneCheck(ctx context.Context, cancelFile, dir string, argv, env []string, out io.Writer, onStart func(pid int), stop func()) (exitCode *int, killedBy string, err error) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
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
			if errors.Is(ctx.Err(), context.Canceled) {
				killedBy = "cancelled"
			}
			break watch
		case <-ticker.C:
			if _, statErr := os.Stat(cancelFile); statErr == nil {
				killedBy = "cancelled"
				break watch
			}
		}
	}
	if stop != nil {
		stop()
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

// envValue is the value env, a list of NAME=value, gives name, the last
// when it gives it twice; empty when it gives none.
func envValue(env []string, name string) string {
	value := ""
	for _, kv := range env {
		if k, v, ok := strings.Cut(kv, "="); ok && k == name {
			value = v
		}
	}
	return value
}
