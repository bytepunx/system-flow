//go:build !windows

package verify

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestOSRunsTiersAsProcessesAndStopsAtTheFailingOne(t *testing.T) {
	dir := t.TempDir()
	res, err := Run(context.Background(), dir, []Selected{
		sel("pass", "", "sh", "-c", "echo fine"),
		sel("fail", "", "sh", "-c", "echo out; echo wrong >&2; exit 3"),
		sel("after", "", "sh", "-c", "touch reached"),
	}, RunOptions{})
	if err != nil || res.Passed {
		t.Fatalf("run %+v, %v", res, err)
	}
	pass, fail, after := res.Tiers[0], res.Tiers[1], res.Tiers[2]
	if pass.State != Passed || *pass.ExitCode != 0 || pass.Duration == "" {
		t.Errorf("pass %+v", pass)
	}
	// stdout and stderr are separate pipes, read in no fixed order
	if fail.State != Failed || *fail.ExitCode != 3 || len(fail.Findings) != 1 || !strings.Contains(fail.Findings[0].Message, "out") ||
		!strings.Contains(fail.Findings[0].Message, "wrong") || !strings.HasSuffix(fail.Findings[0].Message, "\nexit status 3") {
		t.Errorf("fail %+v", fail)
	}
	if _, err := os.Stat(filepath.Join(dir, "reached")); after.State != NotReached || err == nil {
		t.Errorf("after %+v ran: %v", after, err)
	}
}

func TestOSAddsTheEnvironmentToTheOneTheTierInherits(t *testing.T) {
	t.Setenv("FLAI_VERIFY_INHERITED", "kept")
	res, err := Run(context.Background(), t.TempDir(), []Selected{
		sel("env", "", "sh", "-c", `test "$CLOSE_OUT_STORY" = S-0001 && test "$FLAI_VERIFY_INHERITED" = kept`),
	}, RunOptions{Env: []string{"CLOSE_OUT_STORY=S-0001"}})
	if err != nil || !res.Passed {
		t.Errorf("the tier did not see the story and the inherited environment: %+v, %v", res, err)
	}
}

func TestCancellingTheContextKillsTheTiersProcessGroup(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "pid")
	ctx := cancelOnceWritten(t, pidFile)
	start := time.Now()
	res, err := Run(ctx, dir, []Selected{sel("sleep", "", "sh", "-c", "sleep 30 & echo $! > pid; wait")}, RunOptions{})
	if !errors.Is(err, context.Canceled) || res.Tiers[0].State != Failed {
		t.Fatalf("run %+v, %v", res, err)
	}
	if took := time.Since(start); took > grace {
		t.Errorf("the run took %s to stop", took)
	}
	gone(t, pidFile)
}

func TestAProcessGroupThatIgnoresTerminateIsKilled(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "pid")
	ctx := cancelOnceWritten(t, pidFile)
	start := time.Now()
	res, err := Run(ctx, dir, []Selected{sel("stubborn", "", "sh", "-c", `trap "" TERM; sleep 30 & echo $! > pid; wait; wait`)}, RunOptions{})
	if !errors.Is(err, context.Canceled) || res.Tiers[0].State != Failed {
		t.Fatalf("run %+v, %v", res, err)
	}
	if took := time.Since(start); took < grace || took > 2*grace {
		t.Errorf("the run took %s to stop, not the grace period of %s", took, grace)
	}
	gone(t, pidFile)
}

// cancelOnceWritten is a context cancelled once a line is written to file,
// as the tier's command does when it has started what it runs.
func cancelOnceWritten(t *testing.T, file string) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go func() {
		for ctx.Err() == nil {
			if data, err := os.ReadFile(file); err == nil && strings.HasSuffix(string(data), "\n") {
				cancel()
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return ctx
}

// gone fails the test unless the process whose PID the file holds has
// stopped within a few seconds; it may linger as a zombie until reaped.
func gone(t *testing.T, pidFile string) {
	t.Helper()
	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatal(err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		if syscall.Kill(pid, 0) != nil {
			return
		}
	}
	_ = syscall.Kill(pid, syscall.SIGKILL)
	t.Errorf("process %d that the tier started is still running", pid)
}
