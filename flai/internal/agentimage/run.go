package agentimage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// DefaultStall is how long a step may write nothing before it is stopped (ADR-0295 §5).
const DefaultStall = 10 * time.Minute

// RetryDelays are the waits before each further try of a step that failed on the network.
var RetryDelays = []time.Duration{30 * time.Second, 2 * time.Minute}

// tailSize is how much of a step's last output Run keeps to say why it failed.
const tailSize = 8 << 10

// killWait is how long SystemExec waits, after it kills a program, for the
// program's children to let go of its output before it returns.
const killWait = 5 * time.Second

// Exec runs program with args in dir, with env added to the process's environment,
// writing its standard output and standard error to out as they come, until it
// exits or ctx ends (then it is killed); its error is the exit's.
type Exec func(ctx context.Context, dir string, env []string, out io.Writer, program string, args ...string) error

// SystemExec is the Exec of os/exec.
func SystemExec(ctx context.Context, dir string, env []string, out io.Writer, program string, args ...string) error {
	cmd := exec.CommandContext(ctx, program, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.WaitDelay = killWait
	return cmd.Run()
}

// Step is one docker command of a pull or a build: a name the report uses
// (such as "build flai-agent-bedrock:1.2.3"), the folder it runs in, extra
// environment, and the arguments after the docker program.
type Step struct {
	Name string
	Dir  string
	Env  []string
	Args []string
}

// Runner runs steps with Program (the docker client, "docker" when empty)
// through Exec (SystemExec when nil), writing their output to Out (discarded when nil),
// stopping a step silent for Stall (DefaultStall when zero), and waiting
// Delays (RetryDelays when nil) before each further try, with Sleep
// (a timer honouring ctx when nil).
type Runner struct {
	Program string
	Exec    Exec
	Out     io.Writer
	Stall   time.Duration
	Delays  []time.Duration
	Sleep   func(ctx context.Context, d time.Duration) error
}

// Run runs step, trying it again after each of the delays while it fails on the network.
func (r Runner) Run(ctx context.Context, step Step) error {
	out := r.Out
	if out == nil {
		out = io.Discard
	}
	delays := r.Delays
	if delays == nil {
		delays = RetryDelays
	}
	sleep := r.Sleep
	if sleep == nil {
		sleep = timer
	}
	for try := 1; ; try++ {
		tail, err := r.once(ctx, step, out)
		if err == nil {
			return nil
		}
		var stall *StallError
		if errors.As(err, &stall) {
			return err
		}
		if ctx.Err() != nil {
			return fmt.Errorf("%s was stopped before it finished: %w", step.Name, ctx.Err())
		}
		said := lastLine(tail)
		if said == "" {
			said = firstLine(err.Error())
		}
		network := Network(tail)
		if !network || try > len(delays) {
			return &StepError{Step: step.Name, Tries: try, Said: said, Network: network, Err: err}
		}
		wait := delays[try-1]
		fmt.Fprintf(out, "agent image: %s failed on the network (%s); trying again in %s (try %d of %d)\n", step.Name, said, short(wait), try+1, len(delays)+1)
		if err := sleep(ctx, wait); err != nil {
			return fmt.Errorf("wait to try %s again: %w", step.Name, err)
		}
	}
}

// once runs step one time, writing its output to out, and is the tail of
// that output and the step's error: a *StallError when it wrote nothing
// for the Runner's stall limit and was stopped.
func (r Runner) once(ctx context.Context, step Step, out io.Writer) (string, error) {
	program := r.Program
	if program == "" {
		program = "docker"
	}
	run := r.Exec
	if run == nil {
		run = SystemExec
	}
	limit := r.Stall
	if limit == 0 {
		limit = DefaultStall
	}
	stepCtx, stop := context.WithCancel(ctx)
	defer stop()
	w := &watch{out: out, last: time.Now()}
	w.mu.Lock()
	w.timer = time.AfterFunc(limit, func() { w.check(limit, stop) })
	w.mu.Unlock()
	err := run(stepCtx, step.Dir, step.Env, w, program, step.Args...)
	stalled, tail := w.finish()
	if err != nil && stalled && ctx.Err() == nil {
		return tail, &StallError{Step: step.Name, Limit: limit}
	}
	return tail, err
}

// watch is the writer a step writes to: it passes each write on to out,
// one at a time, keeps the last tailSize bytes, and notes when the last
// write came, so that a timer can stop a step that has gone silent.
type watch struct {
	mu      sync.Mutex
	out     io.Writer
	tail    []byte
	last    time.Time
	timer   *time.Timer
	stalled bool
	done    bool
}

// Write passes p on to the watch's out and keeps it in the tail.
func (w *watch) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.last = time.Now()
	w.tail = append(w.tail, p...)
	if over := len(w.tail) - tailSize; over > 0 {
		w.tail = append(w.tail[:0], w.tail[over:]...)
	}
	return w.out.Write(p)
}

// check is the timer's: it stops the step with stop when nothing was
// written for limit, and otherwise sets the timer again for the rest of it.
func (w *watch) check(limit time.Duration, stop context.CancelFunc) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		return
	}
	if idle := time.Since(w.last); idle < limit {
		w.timer.Reset(limit - idle)
		return
	}
	w.stalled = true
	stop()
}

// finish ends the watch once its step has returned, and is whether the
// step was stopped for its silence and the tail of its output.
func (w *watch) finish() (bool, string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.done = true
	w.timer.Stop()
	return w.stalled, string(w.tail)
}

// StallError is a step stopped for writing nothing for Limit.
type StallError struct {
	Step  string
	Limit time.Duration
}

// Error says which step was stopped, after how long, and what to do.
func (e *StallError) Error() string {
	return fmt.Sprintf("%s wrote nothing for %s and was stopped: check the network and that the engine answers (docker info), then run it again", e.Step, short(e.Limit))
}

// StepError is a step that failed: its name, how many times it ran, the
// last line of its output that says why (the last non-empty line), whether that names a network failure, and the error.
type StepError struct {
	Step    string
	Tries   int
	Said    string
	Network bool
	Err     error
}

// Error says which step failed, why, and what to do.
func (e *StepError) Error() string {
	if e.Network {
		return fmt.Sprintf("%s failed (%s) after %s; check the connection to the registry, then run it again", e.Step, e.Said, tries(e.Tries))
	}
	return fmt.Sprintf("%s failed (%s); its output is above and in the log", e.Step, e.Said)
}

// Unwrap is the step's own error.
func (e *StepError) Unwrap() error { return e.Err }

// networkFailures are what a step's output says when it failed on the
// network (ADR-0295 §5): a connection reset, a timeout, a DNS failure, a
// TLS handshake failure, or a registry's 5xx, the status always named as
// one so that a number in a digest or a size is never taken for it.
var networkFailures = regexp.MustCompile(`(?i)` + strings.Join([]string{
	`connection reset`,
	`\btimeout\b`,
	`\btimed out\b`,
	`\bdeadline exceeded\b`,
	`\btls handshake\b`,
	`\bno such host\b`,
	`\btemporary failure in name resolution\b`,
	`\bserver misbehaving\b`,
	`\bcould not resolve\b`,
	`\b5\d\d (?:internal server error|bad gateway|service unavailable|gateway time-?out)\b`,
	`\bstatus(?: code)?:?\s*5\d\d\b`,
	`\bhttp(?:/\d(?:\.\d)?)?:?\s+5\d\d\b`,
}, "|"))

// Network reports whether output, a step's last output, names a network failure ADR-0295 §5 retries.
func Network(output string) bool {
	return networkFailures.MatchString(output)
}

// timer waits d, or until ctx ends, when it is ctx's error.
func timer(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// lastLine is the last line of s that is not blank, trimmed.
func lastLine(s string) string {
	lines := strings.FieldsFunc(s, func(r rune) bool { return r == '\n' || r == '\r' })
	for i := len(lines) - 1; i >= 0; i-- {
		if line := strings.TrimSpace(lines[i]); line != "" {
			return line
		}
	}
	return ""
}

// short is d as time.Duration writes it, less a trailing zero unit: 30s,
// 2m, 10m, 1h.
func short(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

// tries is n tries in words: "1 try", "3 tries".
func tries(n int) string {
	if n == 1 {
		return "1 try"
	}
	return fmt.Sprintf("%d tries", n)
}
