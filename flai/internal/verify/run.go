package verify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
)

// roleEnv names the variable flai reads an agent's role from. Every tier
// runs with it set to the verify role, whatever role started the run, so a
// test that runs flai's writes in process is not refused as the
// orchestrator that ran flai verify (S-0311, TH-0260).
const roleEnv = "FLAI_ROLE"

// Proc runs a tier's command; OS is the real one.
type Proc interface {
	// Run runs argv in dir, with env, KEY=value entries, added to the
	// environment it inherits, writing what it prints to stdout and stderr,
	// and answers its exit status. An error says it could not start or its
	// output could not be read. When ctx ends it stops the command and
	// everything the command started, and answers.
	Run(ctx context.Context, dir string, argv, env []string, stdout, stderr io.Writer) (int, error)
}

// Run runs the selected tiers in list order, each in its Dir below root, and
// stops at the first that fails: one that exits non-zero, cannot start, or,
// for gofmt-list, lists a file. The tiers after it are not reached. A tier
// starts once the tiers it waits for have ended (waits): a tier marked
// Parallel waits only for the tiers before it that are not marked, beside
// the marked ones; any other waits for every tier before it, so with none
// marked the run is serial (ADR-0266). When a tier fails no tier after it in
// the list starts, a tier after it that is running is stopped, and a tier
// before it that is running finishes; every tier after the first in list
// order that failed is not reached, even one that had ended, so the result
// is a serial run's over the same outcomes. The failing tier's output is
// parsed into findings, capped across the run at opts.Max. A tier that
// Select marked covered is not run: it is skipped, naming the tier that
// covers it, in its place in the list, whether or not a tier before it
// failed, and never fails the run. With nothing selected the run passes.
// With opts.LogDir each tier that runs writes its whole output to its own
// file there, which its result names; a tier listed not reached keeps none.
// When ctx ends the tiers running are stopped, the first of them failed, and
// Run returns ctx's error with the result.
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
	outs, ran, first := schedule(ctx, root, selected, opts)
	res := Result{Passed: first == len(selected), Tiers: make([]TierResult, 0, len(selected))}
	for i, s := range selected {
		tr := TierResult{Name: s.Tier.Name, Command: s.Argv, Dir: s.Tier.Dir, State: NotReached}
		switch {
		case s.CoveredBy != "":
			tr.State, tr.CoveredBy = Skipped, s.CoveredBy
		case !ran[i]:
		case i > first:
			// it ran beside the tier that failed, which a serial run
			// would not have reached
			removeLog(root, outs[i].log)
		default:
			out := outs[i]
			tr.DurationMS = out.took.Round(time.Millisecond).Milliseconds()
			tr.Duration = human(out.took)
			tr.Log = out.log
			if out.err == nil {
				exit := out.exit
				tr.ExitCode = &exit
			}
			tr.State = Passed
			if i == first {
				tr.State = Failed
				all := parse(s.Tier, out, env{root: root, dir: cleanDir(s.Tier.Dir), fsys: opts.FS})
				n := min(budget, len(all))
				tr.Findings, tr.Omitted = all[:n], len(all)-n
			}
		}
		res.Tiers = append(res.Tiers, tr)
	}
	return res, ctx.Err()
}

// waits are, for each selected tier, the indexes of the tiers before it
// that it waits for (ADR-0266): none for a covered tier, which does not run,
// and no covered tier for any. A tier marked Parallel waits for the tiers
// before it that are not marked; any other tier waits for every tier before
// it.
func waits(selected []Selected) [][]int {
	out := make([][]int, len(selected))
	for i, s := range selected {
		if s.CoveredBy != "" {
			continue
		}
		for j, before := range selected[:i] {
			if before.CoveredBy == "" && (!s.Tier.Parallel || !before.Tier.Parallel) {
				out[i] = append(out[i], j)
			}
		}
	}
	return out
}

// ended is a tier's command come back.
type ended struct {
	i   int
	out output
}

// schedule runs the selected tiers that are not covered, each once the
// tiers it waits for have ended, and answers what each did, whether it ran,
// and the index of the first in list order that failed, len(selected) when
// none did. Once a tier fails, no tier after it starts and those after it
// that are running are stopped through their own context; the tiers before
// it go on, and one of them failing becomes the first. With no tier marked
// Parallel one tier runs at a time, in list order.
func schedule(ctx context.Context, root string, selected []Selected, opts RunOptions) ([]output, []bool, int) {
	n := len(selected)
	wait := waits(selected)
	outs := make([]output, n)
	ran := make([]bool, n)
	done := make([]bool, n)
	cancels := make([]context.CancelFunc, n)
	first := n
	endc := make(chan ended)
	running := 0
	for {
		for i, s := range selected {
			if ran[i] || s.CoveredBy != "" || i > first || !allDone(done, wait[i]) {
				continue
			}
			ran[i] = true
			tctx, cancel := context.WithCancel(ctx)
			cancels[i] = cancel
			running++
			go func(i int, s Selected) {
				endc <- ended{i, runTier(tctx, root, s, opts)}
			}(i, s)
		}
		if running == 0 {
			return outs, ran, first
		}
		e := <-endc
		running--
		cancels[e.i]()
		done[e.i], outs[e.i] = true, e.out
		if e.i < first && e.out.failed(selected[e.i].Tier.Format) {
			first = e.i
			for k := first + 1; k < n; k++ {
				if ran[k] && !done[k] {
					cancels[k]()
				}
			}
		}
	}
}

// allDone reports whether every tier of indexes has ended.
func allDone(done []bool, indexes []int) bool {
	for _, j := range indexes {
		if !done[j] {
			return false
		}
	}
	return true
}

// removeLog deletes the log file path names, below root unless absolute;
// none when path is empty. A file that cannot be removed is left.
func removeLog(root, path string) {
	if path == "" {
		return
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	_ = os.Remove(path)
}

// output is what one tier's command did.
type output struct {
	stdout, stderr, combined string
	exit                     int
	err                      error
	took                     time.Duration
	// log is the file holding combined, as TierResult.Log names it.
	log string
}

// failed reports whether the command failed, as a tier of format reads it.
func (o output) failed(format string) bool {
	if o.err != nil || o.exit != 0 {
		return true
	}
	return format == FormatGofmtList && len(bytes.TrimSpace([]byte(o.stdout))) > 0
}

// runTier runs one selected tier and keeps what it printed, in a log file
// too when opts.LogDir names a folder. A log that cannot be created or
// written is not named, and the tier's result stands.
func runTier(ctx context.Context, root string, s Selected, opts RunOptions) output {
	var stdout, stderr bytes.Buffer
	var combined lockedBuffer
	logPath, logFile := openLog(root, opts.LogDir, s.Tier.Name)
	if logFile != nil {
		combined.log = logFile
	}
	dir := filepath.Join(root, filepath.FromSlash(s.Tier.Dir))
	start := opts.Now()
	var exit int
	var err error
	// last of a name wins in exec, so the role beats the inherited one and
	// any opts.Env sets
	tierEnv := append(append([]string{}, opts.Env...), roleEnv+"="+conventions.RoleVerify)
	if len(s.Argv) == 0 {
		err = fmt.Errorf("tier %s has no command; give it one in the manifest", s.Tier.Name)
	} else {
		exit, err = opts.Proc.Run(ctx, dir, s.Argv, tierEnv, io.MultiWriter(&stdout, &combined), io.MultiWriter(&stderr, &combined))
	}
	if err == nil && ctx.Err() != nil {
		err = fmt.Errorf("stopped: %w", ctx.Err())
	}
	out := output{
		stdout: stdout.String(), stderr: stderr.String(), combined: combined.String(),
		exit: exit, err: err, took: opts.Now().Sub(start),
	}
	if logFile != nil && logFile.Close() == nil && combined.logErr == nil {
		out.log = logPath
	}
	return out
}

// unsafeInLogName is what a tier's name may not keep in its log file's name.
var unsafeInLogName = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// logName is the name of the log file of the tier named name: the name,
// with each character outside A-Z, a-z, 0-9, dot, underscore, and hyphen a
// hyphen, and .log.
func logName(name string) string {
	return unsafeInLogName.ReplaceAllString(name, "-") + ".log"
}

// openLog creates the log file of the tier named name in dir, below root
// unless dir is absolute, and answers its path as dir names it; none when
// dir is empty or the file cannot be created.
func openLog(root, dir, name string) (string, *os.File) {
	if dir == "" {
		return "", nil
	}
	path := filepath.Join(dir, logName(name))
	full := path
	if !filepath.IsAbs(full) {
		full = filepath.Join(root, full)
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", nil
	}
	f, err := os.Create(full)
	if err != nil {
		return "", nil
	}
	return path, f
}

// lockedBuffer is a buffer two copies, of stdout and of stderr, write at
// once, and that copies each write to log, when it has one, until a write
// to log fails, which logErr keeps.
type lockedBuffer struct {
	mu     sync.Mutex
	buf    bytes.Buffer
	log    io.Writer
	logErr error
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.log != nil && b.logErr == nil {
		_, b.logErr = b.log.Write(p)
	}
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
