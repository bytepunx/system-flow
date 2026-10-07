package verify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Proc runs a tier's command; OS is the real one.
type Proc interface {
	// Run runs argv in dir, with env, KEY=value entries, added to the
	// environment it inherits, writing what it prints to stdout and stderr,
	// and answers its exit status. An error says it could not start or its
	// output could not be read. When ctx ends it stops the command and
	// everything the command started, and answers.
	Run(ctx context.Context, dir string, argv, env []string, stdout, stderr io.Writer) (int, error)
}

// Run runs the selected tiers in order, each in its Dir below root, and
// stops at the first that fails: one that exits non-zero, cannot start, or,
// for gofmt-list, lists a file. The tiers after it are not reached. The
// failing tier's output is parsed into findings, capped across the run at
// opts.Max. With nothing selected the run passes. When ctx ends the tier
// running is stopped and failed, and Run returns ctx's error with the result.
func Run(ctx context.Context, root string, selected []Selected, opts RunOptions) (Result, error) {
	if opts.FS == nil {
		opts.FS = os.DirFS(root)
	}
	if opts.Proc == nil {
		opts.Proc = OS{}
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	budget := opts.Max
	if budget <= 0 {
		budget = DefaultMax
	}
	res := Result{Passed: true, Tiers: make([]TierResult, 0, len(selected))}
	for _, s := range selected {
		tr := TierResult{Name: s.Tier.Name, Command: s.Argv, Dir: s.Tier.Dir, State: NotReached}
		if !res.Passed {
			res.Tiers = append(res.Tiers, tr)
			continue
		}
		out := runTier(ctx, root, s, opts)
		tr.DurationMS = out.took.Round(time.Millisecond).Milliseconds()
		tr.Duration = human(out.took)
		if out.err == nil {
			tr.ExitCode = &out.exit
		}
		tr.State = Passed
		if out.failed(s.Tier.Format) {
			tr.State = Failed
			res.Passed = false
			all := parse(s.Tier, out, env{root: root, dir: cleanDir(s.Tier.Dir), fsys: opts.FS})
			n := min(budget, len(all))
			tr.Findings, tr.Omitted = all[:n], len(all)-n
			budget -= n
		}
		res.Tiers = append(res.Tiers, tr)
	}
	return res, ctx.Err()
}

// output is what one tier's command did.
type output struct {
	stdout, stderr, combined string
	exit                     int
	err                      error
	took                     time.Duration
}

// failed reports whether the command failed, as a tier of format reads it.
func (o output) failed(format string) bool {
	if o.err != nil || o.exit != 0 {
		return true
	}
	return format == FormatGofmtList && len(bytes.TrimSpace([]byte(o.stdout))) > 0
}

// runTier runs one selected tier and keeps what it printed.
func runTier(ctx context.Context, root string, s Selected, opts RunOptions) output {
	var stdout, stderr bytes.Buffer
	var combined lockedBuffer
	dir := filepath.Join(root, filepath.FromSlash(s.Tier.Dir))
	start := opts.Now()
	var exit int
	var err error
	if len(s.Argv) == 0 {
		err = fmt.Errorf("tier %s has no command; give it one in the manifest", s.Tier.Name)
	} else {
		exit, err = opts.Proc.Run(ctx, dir, s.Argv, opts.Env, io.MultiWriter(&stdout, &combined), io.MultiWriter(&stderr, &combined))
	}
	if err == nil && ctx.Err() != nil {
		err = fmt.Errorf("stopped: %w", ctx.Err())
	}
	return output{
		stdout: stdout.String(), stderr: stderr.String(), combined: combined.String(),
		exit: exit, err: err, took: opts.Now().Sub(start),
	}
}

// lockedBuffer is a buffer two copies, of stdout and of stderr, write at
// once.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// human is a duration as a person reads it: milliseconds under a second,
// tenths of a second above.
func human(d time.Duration) string {
	if d < time.Second {
		return d.Round(time.Millisecond).String()
	}
	return d.Round(100 * time.Millisecond).String()
}
