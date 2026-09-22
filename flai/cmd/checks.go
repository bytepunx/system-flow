package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// flai checks: running the operator's named commands for a story in review,
// in the story's worktree, and reporting the outcome (S-0082, ADR-0029).

func newChecksCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "checks",
		Short: "Run, watch, and cancel checks for a story in review",
		Args:  cobra.NoArgs,
	}
	c.AddCommand(newChecksRunCmd(a), newChecksStatusCmd(a), newChecksCancelCmd(a), newChecksTailCmd(a))
	return c
}

func newChecksRunCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "run <story-id>",
		Short: "Run every named check, in order, in the story's worktree; refuses a second run while one is active",
		Long: `Runs each check named in flai serve checks set, or, when it names none, the
manifest's checks:, in order, in the story's worktree, stopping at the
first to fail. Bounded by the configured time limit (flai serve checks
timeout); flai checks cancel stops it early. This blocks until the run
ends, real minutes, not seconds; the outcome and duration stay with the
story (flai checks status) until it is accepted.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.runChecksRun(args[0])
		},
	}
}

func (a *app) runChecksRun(story string) error {
	if !serve.StoryID.MatchString(story) {
		return fmt.Errorf("%q is not a story's ID", story)
	}
	repo, err := a.project()
	if err != nil {
		return err
	}
	cfg, err := a.checksConfig(repo)
	if err != nil {
		return err
	}
	worktree := repo.WorktreePath(story)
	run, err := a.serveDir().RunChecks(story, mainRootOf(repo), worktree, cfg, a.now)
	if err != nil {
		return err
	}
	return a.printChecksRun(run)
}

func newChecksStatusCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "status <story-id>",
		Short: "The current or last run for a story: which check, pass or fail so far, outcome and duration once ended",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			run, ok := a.serveDir().ChecksState(args[0])
			if !ok {
				run = serve.ChecksRun{Story: args[0]}
			}
			return a.printChecksRun(run)
		},
	}
}

func newChecksCancelCmd(a *app) *cobra.Command {
	var grace int
	c := &cobra.Command{
		Use:   "cancel <story-id>",
		Short: "Stop an active run: terminate, then, if it has not gone, kill",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			d := time.Duration(grace) * time.Second
			if err := a.serveDir().CancelChecks(args[0], d, a.sleep); err != nil {
				return err
			}
			run, _ := a.serveDir().ChecksState(args[0])
			return a.printChecksRun(run)
		},
	}
	c.Flags().IntVar(&grace, "grace-seconds", 10, "how long to wait for the run to stop before giving up on reporting it (the kill itself has its own, shorter grace period)")
	return c
}

func (a *app) printChecksRun(run serve.ChecksRun) error {
	if a.jsonOut {
		return a.printJSON(run)
	}
	if run.Started == "" {
		fmt.Fprintf(a.out, "%s: no checks have run\n", run.Story)
		return nil
	}
	if run.Running {
		fmt.Fprintf(a.out, "%s: running %s (pid %d), started %s\n", run.Story, run.Current, run.PID, run.Started)
	} else {
		fmt.Fprintf(a.out, "%s: %s, started %s, ended %s\n", run.Story, run.Outcome, run.Started, run.Ended)
	}
	for _, s := range run.Steps {
		state := "ok"
		if s.Exit == nil {
			state = "did not start"
		} else if *s.Exit != 0 {
			state = fmt.Sprintf("exit %d", *s.Exit)
		}
		fmt.Fprintf(a.out, "  %s (%s): %s\n", s.Name, s.Command, state)
	}
	return nil
}

func newChecksTailCmd(a *app) *cobra.Command {
	var from, wait int
	c := &cobra.Command{
		Use:   "tail <story-id>",
		Short: "New output since --from, waiting for more while the run is active",
		Long: `Reads the run's log from --from, an offset in bytes into the file (0 for the
start), and prints each new complete line as it is found. While the run is
still active it polls for up to --wait seconds for more before returning;
the answer's offset is where to --from next time. Meant to be called in a
loop for as long as the review page (or a shell) wants to keep watching.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.runChecksTail(args[0], int64(from), time.Duration(wait)*time.Second)
		},
	}
	c.Flags().IntVar(&from, "from", 0, "byte offset into the log to read from")
	c.Flags().IntVar(&wait, "wait", 20, "seconds to wait for more output while the run is active")
	return c
}

func (a *app) runChecksTail(story string, from int64, wait time.Duration) error {
	deadline := a.now().Add(wait)
	for {
		run, _ := a.serveDir().ChecksState(story)
		next, err := a.tailOnce(run.Log, from)
		if err != nil {
			return err
		}
		from = next
		if !run.Running || a.now().After(deadline) {
			if a.jsonOut {
				return a.printJSON(map[string]any{"offset": from, "running": run.Running, "outcome": run.Outcome})
			}
			fmt.Fprintf(a.errOut, "offset %d, running %v\n", from, run.Running)
			return nil
		}
		a.sleep(500 * time.Millisecond)
	}
}

// tailOnce prints, as log events (so a host action with progress: true
// streams them with no new transport, S-0081's lesson reused), every
// complete line found from offset to the end of the log, and returns the
// new offset: the start of whatever incomplete line, if any, is left.
func (a *app) tailOnce(logPath string, from int64) (int64, error) {
	if logPath == "" {
		return from, nil
	}
	f, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return from, nil
		}
		return from, err
	}
	defer func() { _ = f.Close() }()
	if _, err := f.Seek(from, 0); err != nil {
		return from, err
	}
	r := bufio.NewReader(f)
	offset := from
	for {
		line, err := r.ReadString('\n')
		if line != "" && strings.HasSuffix(line, "\n") {
			a.logger().Info("line", "component", "checks", "text", strings.TrimSuffix(line, "\n"))
			offset += int64(len(line))
		}
		if err != nil {
			break // EOF, or an incomplete last line: leave offset before it
		}
	}
	return offset, nil
}
