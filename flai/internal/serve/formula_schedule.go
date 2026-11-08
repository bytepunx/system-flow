package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A formula's instance, made on the schedule its file names (S-0422,
// ADR-0238).
//
// A formula scheduler per served project looks whenever the project's
// launcher does, as the analysis scheduler does (S-0223). At each look it
// reads the project's formula files, and follows the schedule of each that
// names one. Each time a formula's schedule comes round with the formula host
// action on, it makes the formula's instance by running flai formula run
// <name> --ready --autocommit in the project's main checkout, with no inputs,
// through the flai flai serve runs as, so that the instance is committed in
// ready and the launcher starts its story's agent as any ready story's. The
// instance it made is journalled under the formula action, with the schedule
// as its trigger ("schedule daily"). A run flai formula run refuses, as for a
// formula flai check refuses, is logged and not tried again until the
// schedule next comes round. Times missed between two looks come round once,
// and a schedule set or changed, or a formula added, comes round first at
// its next time from then.
//
// What it has seen is kept in memory only, as the analysis scheduler's is.
// While the formula host action is off the schedules make nothing and keep
// up, so that turning it on does not act on the past.

// formulaRunTimeout is how long a scheduled flai formula run may take before
// it is stopped: the project's loop waits for it.
const formulaRunTimeout = 2 * time.Minute

// formulaSaid is how many of the last lines flai formula run wrote on
// standard error a refusal's log line gives.
const formulaSaid = 5

// outputRunner runs name with args in dir, with env as its whole
// environment, until it ends or ctx does, and answers what it wrote on
// standard output and on standard error, with an error when it could not run
// or exited non-zero.
type outputRunner func(ctx context.Context, dir string, env []string, name string, args ...string) (stdout, stderr []byte, err error)

// runOutput is the outputRunner that runs the process: ctx ending interrupts
// it, and kills it ten seconds later, as runCommand does.
func runOutput(ctx context.Context, dir string, env []string, name string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir, cmd.Env = dir, env
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 10 * time.Second
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	// what it wrote under dir's item folders is seen by this process's next list (S-0606)
	workitem.ItemsChanged(dir)
	return stdout.Bytes(), stderr.Bytes(), err
}

// formulaScheduler makes the instances of one project's formulas on the
// schedules their files name (S-0422).
type formulaScheduler struct {
	o       Options
	e       Entry
	starter *launcher
	// run runs flai formula run, and exe is this flai's executable.
	run outputRunner
	exe func() (string, error)

	mu sync.Mutex
	// follows are the schedules followed, by formula name.
	follows map[string]*cadence
	said    map[string]bool // the problems warned of
}

// newFormulaScheduler is the formula scheduler for project e, logging and
// journalling through starter, the project's launcher. Times before now do
// not come round.
func newFormulaScheduler(o Options, e Entry, starter *launcher) *formulaScheduler {
	s := &formulaScheduler{o: o, e: e, starter: starter, run: runOutput, exe: os.Executable, follows: map[string]*cadence{}, said: map[string]bool{}}
	now := o.Now().UTC()
	if list, err := s.scheduled(); err == nil {
		for _, f := range list {
			s.follow(f, now)
		}
	}
	return s
}

// look makes the instance of each formula whose schedule has come round
// since the last look, while the formula host action is on. Formulas no file
// schedules any longer are no longer followed.
func (s *formulaScheduler) look(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.o.Now().UTC()
	list, err := s.scheduled()
	if err != nil {
		s.warnOnce("formula scheduler could not read the formula files", err)
		return
	}
	on := s.o.Host.Enabled != nil && s.o.Host.Enabled(hostapi.ActionFormula, s.e.Root)
	follows := map[string]*cadence{}
	for _, f := range list {
		c := s.follow(f, now)
		follows[f.Name] = c
		s.lookAt(ctx, f, c, on, now)
	}
	s.follows = follows
}

// lookAt makes f's instance when its schedule has come round for c by now
// while on, the formula host action, is. While it is off the schedule keeps
// up. A schedule that does not parse is warned of once and makes nothing.
func (s *formulaScheduler) lookAt(ctx context.Context, f formulas.Formula, c *cadence, on bool, now time.Time) {
	sched, err := f.Scheduled()
	if err != nil {
		s.warnOnce("formula schedule not valid", err, "formula", f.Name, "path", f.Path)
		return
	}
	if !on {
		c.mark = now
		return
	}
	if !c.due(sched, now) {
		return
	}
	s.make(ctx, f, scheduleTrigger(sched), now)
}

// make runs flai formula run for f, with the schedule's trigger, and
// journals the instance it made; a run that fails is logged, at info when
// flai formula run refused it.
func (s *formulaScheduler) make(ctx context.Context, f formulas.Formula, trigger string, now time.Time) {
	flai, err := s.exe()
	if err != nil {
		s.starter.warn("scheduled formula not run", "formula", f.Name, "trigger", trigger, "err", fmt.Errorf("this flai's executable is not known: %w", err))
		return
	}
	rctx, cancel := context.WithTimeout(ctx, formulaRunTimeout)
	defer cancel()
	stdout, stderr, err := s.run(rctx, s.e.Root, acceptEnv(os.Environ()), flai, "formula", "run", f.Name, "--ready", "--autocommit", "--json")
	if err != nil {
		reason := lastLines(string(stderr), formulaSaid)
		if reason == "" {
			reason = err.Error()
		}
		s.starter.log("scheduled formula not run", "formula", f.Name, "trigger", trigger, "reason", reason)
		return
	}
	var made struct {
		Story       string `json:"story"`
		Status      string `json:"status"`
		Commit      string `json:"commit"`
		CommitError string `json:"commit_error"`
	}
	if err := json.Unmarshal(stdout, &made); err != nil || made.Story == "" {
		made.Story, made.Status = "a story flai formula run did not name", "unknown"
	}
	detail := fmt.Sprintf("made %s, the instance of formula %s, in %s, started by %s", made.Story, f.Name, made.Status, trigger)
	switch {
	case made.Commit != "":
		detail += "; committed " + made.Commit
	case made.CommitError != "":
		detail += "; not committed: " + firstLine(made.CommitError, made.CommitError)
	}
	if s.starter.record != nil {
		s.starter.record(hostapi.Entry{At: now.Format(time.RFC3339), Action: hostapi.ActionFormula, Method: "serve.formula", Project: s.e.Key, Root: s.e.Root,
			By: "flai serve", Outcome: "done", Detail: detail})
	}
	s.starter.log("scheduled formula instance made", "formula", f.Name, "trigger", trigger, "story", made.Story, "status", made.Status)
}

// scheduled are the project's formulas whose file names a schedule, each
// name once; a file that does not parse is passed over, as flai check
// reports it.
func (s *formulaScheduler) scheduled() ([]formulas.Formula, error) {
	repo, err := workitem.Open(s.e.Root)
	if err != nil {
		return nil, err
	}
	all, _, err := formulas.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		return nil, err
	}
	var out []formulas.Formula
	seen := map[string]bool{}
	for _, f := range all {
		if strings.TrimSpace(f.Schedule) == "" || f.Name == "" || seen[f.Name] {
			continue
		}
		seen[f.Name] = true
		out = append(out, f)
	}
	return out, nil
}

// follow is how f's schedule is followed: as before, or from now when it was
// not, and with its schedule seen at now, so that a schedule set or changed
// comes round first at its next time from now.
func (s *formulaScheduler) follow(f formulas.Formula, now time.Time) *cadence {
	c := s.follows[f.Name]
	if c == nil {
		c = &cadence{mark: now}
		s.follows[f.Name] = c
	}
	c.see(f.Schedule, now)
	return c
}

// warnOnce warns of err with msg, and args, the first time it is seen.
func (s *formulaScheduler) warnOnce(msg string, err error, args ...any) {
	if key := msg + "\x00" + err.Error(); !s.said[key] {
		s.said[key] = true
		s.starter.warn(msg, append(args, "err", err)...)
	}
}
