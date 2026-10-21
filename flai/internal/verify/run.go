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

// Run runs the selected tiers in order, each in its Dir below root, and
// stops at the first that fails: one that exits non-zero, cannot start, or,
// for gofmt-list, lists a file. The tiers after it are not reached. The
// failing tier's output is parsed into findings, capped across the run at
// opts.Max. A tier that Select marked covered is not run: it is skipped,
// naming the tier that covers it, in its place in the list, whether or not a
// tier before it failed, and never fails the run. With nothing selected the
// run passes. With opts.LogDir each tier that runs writes its whole output
// to its own file there, which its result names. When ctx ends the tier
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
		if s.CoveredBy != "" {
			tr.State, tr.CoveredBy = Skipped, s.CoveredBy
		}
		if tr.State == Skipped || !res.Passed {
			res.Tiers = append(res.Tiers, tr)
			continue
		}
		out := runTier(ctx, root, s, opts)
		tr.DurationMS = out.took.Round(time.Millisecond).Milliseconds()
		tr.Duration = human(out.took)
		tr.Log = out.log
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
