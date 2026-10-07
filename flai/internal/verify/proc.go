package verify

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"
)

// grace is how long a stopped command has to end on its own before it and
// everything it started are killed.
const grace = 2 * time.Second

// waitDelay is how long a command that has exited may leave its output open,
// held by something it started, before the run stops reading it.
const waitDelay = 5 * time.Second

// OS runs commands as processes, each in a process group of its own, so that
// stopping one reaches everything it started (as serve's runOneCheck does).
type OS struct{}

// Run runs argv in dir as a process, with env added to the environment flai
// has; when ctx ends it asks the process group to stop and, if it has not
// within a grace period, kills it.
func (OS) Run(ctx context.Context, dir string, argv, env []string, stdout, stderr io.Writer) (int, error) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	cmd.WaitDelay = waitDelay
	ownGroup(cmd)
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return exitStatus(cmd, err)
	case <-ctx.Done():
	}
	_ = terminateGroup(pid)
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case err := <-done:
		// what the process started may outlive it: kill what is left
		_ = killGroup(pid)
		return exitStatus(cmd, err)
	case <-timer.C:
	}
	_ = killGroup(pid)
	return exitStatus(cmd, <-done)
}

// exitStatus is the exit status Wait's error carries, or the error when it
// says the process could not be waited for.
func exitStatus(cmd *exec.Cmd, err error) (int, error) {
	var xe *exec.ExitError
	switch {
	case err == nil:
		return 0, nil
	case errors.As(err, &xe):
		return xe.ExitCode(), nil
	case errors.Is(err, exec.ErrWaitDelay) && cmd.ProcessState != nil:
		return cmd.ProcessState.ExitCode(), nil
	}
	return 0, err
}
