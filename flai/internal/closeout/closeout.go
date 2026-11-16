// Package closeout is the hand-off of a story's close-out to flai serve
// (S-0546, ADR-0229): the files a story's agent writes to ask for its
// close-out as a job, and those flai serve writes to answer it and to keep the
// job's state, under the project's .flai-cache/verify/jobs. The asker, the MCP
// tool verify or flai verify --job, writes S-nnnn.request.json; flai serve
// takes it by renaming it to S-nnnn.taken.json and answers in
// S-nnnn.answer.json with the job's id, or the verdict the verify record
// holds, as the merge pass takes a merge request. The job's state is
// S-nnnn.job.json and its output S-nnnn.log. Both sides import this package,
// so neither imports the other.
package closeout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// The suffixes of a story's files in the jobs folder: its request, the
// request once flai serve took it, flai serve's answer, the job's state, and
// the job's output.
const (
	RequestSuffix = ".request.json"
	TakenSuffix   = ".taken.json"
	AnswerSuffix  = ".answer.json"
	JobSuffix     = ".job.json"
	LogSuffix     = ".log"
)

// The verdicts of a close-out (ADR-0229 §3).
const (
	// VerdictPassed is a close-out that exited 0.
	VerdictPassed = "passed"
	// VerdictStopped is one that did not, at the step the verify record names.
	VerdictStopped = "stopped"
	// VerdictFailed is one that did not, at no step the verify record names.
	VerdictFailed = "failed"
	// VerdictInterrupted is a job flai serve was running when it stopped.
	VerdictInterrupted = "interrupted"
)

// KindResolve is the kind of a Rebase and resolve job: the story's branch
// rebased onto the main branch, then its close-out run at the rebased head,
// the story sent back to its agent only when the rebase stops on conflicts
// or a step of the close-out fails (S-0804). A close-out's kind is "".
const KindResolve = "resolve"

// The ends of a resolve job, as its Outcome says.
const (
	// OutcomeStayed is a resolve job whose rebase was clean and whose
	// close-out passed: the story stays in review.
	OutcomeStayed = "stayed in review"
	// OutcomeSentBack is one that sent the story back to in-progress.
	OutcomeSentBack = "sent back"
)

// ErrServeNotTaking is a close-out request no flai serve took in time.
var ErrServeNotTaking = errors.New("flai serve is not running for this project: run the close-out script, or flai verify where the project has none, yourself")

// Request is a story's agent asking flai serve for its close-out.
type Request struct {
	// Story is the story's ID, zero padded.
	Story string `json:"story"`
	// Kind is the job asked for: KindResolve, or "" for a close-out.
	Kind string `json:"kind,omitempty"`
	// By is the agent that asks, FLAI_AGENT of the job.
	By string `json:"by"`
	// Message is the commit message the close-out commits with; "" for none.
	Message string `json:"message,omitempty"`
	// All and Fresh are flai verify's --all and --fresh.
	All   bool `json:"all,omitempty"`
	Fresh bool `json:"fresh,omitempty"`
	// At is when it was asked; Ask stamps it when zero.
	At time.Time `json:"at"`
}

// Answer is flai serve's answer to a close-out request.
type Answer struct {
	// Job is the job's id; "" when none runs for the request.
	Job   string `json:"job,omitempty"`
	Story string `json:"story"`
	// Started is whether this request started the job; false when a job of
	// the story was running and its id is answered.
	Started bool `json:"started"`
	// Recorded is an answer from the verify record: Verdict, StoppedAt,
	// Commit, and RanAt are its run's, and nothing runs.
	Recorded  bool      `json:"recorded,omitempty"`
	Verdict   string    `json:"verdict,omitempty"`
	StoppedAt string    `json:"stopped_at,omitempty"`
	Commit    string    `json:"commit,omitempty"`
	RanAt     time.Time `json:"ran_at,omitzero"`
	// Log is where the job's output is kept.
	Log string `json:"log,omitempty"`
	// Reason is why flai serve ran nothing for the request, saying what to
	// do; "" when it started a job, found one, or answered from the record.
	Reason string `json:"reason,omitempty"`
}

// Refused reports whether flai serve ran nothing for the request.
func (a Answer) Refused() bool { return a.Reason != "" }

// Running reports whether the answer names a job of the story that was
// running already, which the request did not start.
func (a Answer) Running() bool { return a.Job != "" && !a.Started }

// Job is the state of a story's close-out or resolve job, S-nnnn.job.json.
type Job struct {
	ID    string `json:"id"`
	Story string `json:"story"`
	// Kind is the job's kind: KindResolve, or "" for a close-out.
	Kind    string    `json:"kind,omitempty"`
	By      string    `json:"by"`
	Message string    `json:"message,omitempty"`
	All     bool      `json:"all,omitempty"`
	Fresh   bool      `json:"fresh,omitempty"`
	Started time.Time `json:"started"`
	// Head is the worktree's head when the job's run began, written once it
	// is read, before the close-out starts; "" until then, when it could not
	// be read, and in a job written before S-0799.
	Head string `json:"head,omitempty"`
	// Ended and the rest are set once the job ended: ExitCode is the
	// close-out's, -1 when it could not run or was interrupted; Commit is
	// the worktree's head at the end; DurationMS is how long it ran.
	Ended      time.Time `json:"ended,omitzero"`
	ExitCode   *int      `json:"exit_code,omitempty"`
	Verdict    string    `json:"verdict,omitempty"`
	StoppedAt  string    `json:"stopped_at,omitempty"`
	Commit     string    `json:"commit,omitempty"`
	DurationMS int64     `json:"duration_ms,omitempty"`
	// Outcome is how a resolve job ended, OutcomeStayed or OutcomeSentBack;
	// "" for a close-out, for a resolve job while it runs, and for one
	// interrupted.
	Outcome string `json:"outcome,omitempty"`
	// Thread is the send-back thread a resolve job opened; "" when it sent
	// nothing back.
	Thread string `json:"thread,omitempty"`
	// Conflicts are the paths a resolve job's rebase stopped on, which it
	// aborted.
	Conflicts []string `json:"conflicts,omitempty"`
	// Reason says why a resolve job sent the story back, or what it did when
	// the story stayed in review.
	Reason string `json:"reason,omitempty"`
}

// Running reports whether the job has not ended.
func (j Job) Running() bool { return j.Ended.IsZero() }

// The steps of a close-out as its log shows them, in the order they run,
// named as the close-out's last line names the step it stopped at; after
// StepSync the next round begins again at StepVerify (ADR-0156).
const (
	// StepVerify is flai verify, from the rebase to the last tier.
	StepVerify = "flai verify"
	// StepCommit is the commit of what is outstanding.
	StepCommit = "commit"
	// StepSyncCheck is flai verify --sync-only.
	StepSyncCheck = "the sync check"
	// StepSync is flai stream sync once main moved, before the next round.
	StepSync = "flai stream sync"
)

// The steps of a resolve job, in the order they run, as its progress lines
// name them.
const (
	// StepRebase is the rebase of the story's branch onto the main branch.
	StepRebase = "rebase"
	// StepCloseOut is the story's close-out at the rebased head.
	StepCloseOut = "close-out"
)

// ProgressEnv names the variable that, set to 1, has flai verify write a
// progress line to stderr as each of its steps and tiers starts and ends,
// for a job's log.
const ProgressEnv = "FLAI_VERIFY_PROGRESS"

// The sources of the progress lines a job's log holds: flai verify's, with
// ProgressEnv set, and the resolve job's own.
const (
	SourceVerify  = "verify"
	SourceResolve = "resolve"
)

// The kinds of what a progress line says started or ended: a step, or a
// test or lint tier.
const (
	LineStep = "step"
	LineTier = "tier"
)

// StartedLine is the progress line source writes as story's step or tier,
// by kind, named name starts at at: "<source>: <story> <at> <kind> started
// <name>", at in RFC 3339 UTC, with no newline.
func StartedLine(source, story, kind, name string, at time.Time) string {
	return fmt.Sprintf("%s: %s %s %s started %s", source, story, at.UTC().Format(time.RFC3339), kind, name)
}

// EndedLine is the progress line source writes as story's step or tier ends
// at at with outcome, one word such as passed, failed, skipped, reused, or
// queued: "<source>: <story> <at> <kind> ended <outcome> <name>", with no
// newline.
func EndedLine(source, story, kind, outcome, name string, at time.Time) string {
	return fmt.Sprintf("%s: %s %s %s ended %s %s", source, story, at.UTC().Format(time.RFC3339), kind, outcome, name)
}

// StepProgress is a step or tier of a job as its log's progress lines show
// it.
type StepProgress struct {
	Name string `json:"name"`
	// Kind is LineStep or LineTier.
	Kind string `json:"kind"`
	// Source is who ran it: SourceVerify or SourceResolve.
	Source string `json:"source"`
	// Started is when it started; zero for a tier that ended without
	// running, as one reused, skipped, or left to the merge queue.
	Started time.Time `json:"started,omitzero"`
	// Ended is when it ended, and Outcome how; zero and "" while it runs.
	Ended   time.Time `json:"ended,omitzero"`
	Outcome string    `json:"outcome,omitempty"`
}

// Progress is how far a story's running close-out or resolve job has got.
type Progress struct {
	// ID is the job's id.
	ID string `json:"id"`
	// Started is when the job started.
	Started time.Time `json:"started"`
	// Head is the worktree's head when its run began; "" while the job has
	// not written it yet and in a job written before S-0799.
	Head string `json:"head,omitempty"`
	// Kind is the job's kind: KindResolve, or "" for a close-out.
	Kind string `json:"kind,omitempty"`
	// Step is the step the job is on: StepRebase while a resolve job's
	// rebase runs, else the close-out's, one of StepVerify, StepCommit,
	// StepSyncCheck, and StepSync, or "" when its log cannot say.
	Step string `json:"step,omitempty"`
	// Tier is the tier running, the last its log says started and has not
	// ended; "" when none does.
	Tier string `json:"tier,omitempty"`
}

// Running is the progress of story's close-out or resolve job under the
// project's cache dir and true while the job runs, or false with no error
// when the story has had no job or its job ended. The step and the tier are
// read from the job's log (readLog); Status gives its steps too. The step is that of the
// lines the close-out script, and flai serve's close-out of a container
// story, print as each step begins, and StepVerify throughout flai verify,
// whose own steps and tiers are the progress lines it writes with
// ProgressEnv set. Step is "" when the log has no such line yet: before the
// verify starts, when the log is not there, between a resolve job's rebase
// and its close-out, and for the whole of a close-out that is flai serve's
// own flai verify, run where the worktree has no close-out script.
func Running(cacheDir, story string) (Progress, bool, error) {
	j, ok, err := ReadJob(cacheDir, story)
	if err != nil || !ok || !j.Running() {
		return Progress{}, false, err
	}
	lp, err := readJobLog(cacheDir, j)
	if err != nil {
		return Progress{}, false, err
	}
	return Progress{ID: j.ID, Started: j.Started, Head: j.Head, Kind: j.Kind, Step: lp.step, Tier: lp.tier}, true, nil
}

// JobStatus is a story's last job, running or ended, with what its log says
// of its progress and where the log is: what the dashboard follows.
type JobStatus struct {
	Job
	// Step is the step it is on, as Progress has it; for an ended job, the
	// step it was on as it ended.
	Step string `json:"step,omitempty"`
	// Tier is the tier running; "" when none does, and once the job ended.
	Tier string `json:"tier,omitempty"`
	// Steps are the steps and tiers its log's progress lines say started or
	// ended, in the order they began: by their start, or by their end for
	// one that never started.
	Steps []StepProgress `json:"steps,omitempty"`
	// Log is where the job's output is kept, LogPath.
	Log string `json:"log"`
}

// Status is story's last close-out or resolve job under the project's cache
// dir, running or ended, with its progress read from its log as Running
// reads it, and true; or false with no error when the story has had no job.
func Status(cacheDir, story string) (JobStatus, bool, error) {
	j, ok, err := ReadJob(cacheDir, story)
	if err != nil || !ok {
		return JobStatus{}, false, err
	}
	lp, err := readJobLog(cacheDir, j)
	if err != nil {
		return JobStatus{}, false, err
	}
	st := JobStatus{Job: j, Step: lp.step, Steps: lp.steps, Log: LogPath(cacheDir, j.Story)}
	if j.Running() {
		st.Tier = lp.tier
	}
	return st, true, nil
}

// readJobLog is what job j's log under the project's cache dir says of its
// progress; nothing when the log is not there.
func readJobLog(cacheDir string, j Job) (logProgress, error) {
	path := LogPath(cacheDir, j.Story)
	data, err := os.ReadFile(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return logProgress{}, nil
	case err != nil:
		return logProgress{}, fmt.Errorf("read the log of job %s at %s: %w; check that the project's .flai-cache is readable", j.ID, path, err)
	}
	return readLog(string(data), j.Story), nil
}

// logProgress is what a job's log says of its progress, as Progress has it.
type logProgress struct {
	step, tier string
	steps      []StepProgress
}

// readLog is what story's job's log says of its progress: each progress
// line of story's, in order, a step or tier begun by its start, or by its
// end when it has none open of its source, kind, and name; the tier running,
// the last begun that has not ended; and the step, StepRebase while a
// resolve job's rebase runs, else as the close-out's lines say (stepAfter).
func readLog(log, story string) logProgress {
	var p logProgress
	rebasing := false
	for line := range strings.Lines(log) {
		line = strings.TrimRight(line, "\r\n")
		m, ok := parseLine(line, story)
		if !ok {
			p.step = stepAfter(p.step, line, story)
			continue
		}
		p.add(m)
		if m.Source == SourceResolve && m.Kind == LineStep && m.Name == StepRebase {
			rebasing = m.Ended.IsZero()
		}
	}
	if rebasing {
		p.step = StepRebase
	}
	for _, s := range slices.Backward(p.steps) {
		if s.Kind == LineTier && !s.Started.IsZero() && s.Ended.IsZero() {
			p.tier = s.Name
			break
		}
	}
	return p
}

// add puts the step or tier a progress line says started or ended on p: a
// start begins one; an end ends the last of its source, kind, and name that
// is open, or begins one that ended without a start.
func (p *logProgress) add(m StepProgress) {
	if m.Ended.IsZero() {
		p.steps = append(p.steps, m)
		return
	}
	for i := len(p.steps) - 1; i >= 0; i-- {
		s := &p.steps[i]
		if s.Ended.IsZero() && s.Source == m.Source && s.Kind == m.Kind && s.Name == m.Name {
			s.Ended, s.Outcome = m.Ended, m.Outcome
			return
		}
	}
	p.steps = append(p.steps, m)
}

// parseLine reads a progress line of story's, as StartedLine and EndedLine
// write it, into the step or tier it says started, with Started set, or
// ended, with Ended and Outcome set; false for any other line.
func parseLine(line, story string) (StepProgress, bool) {
	source, rest, ok := strings.Cut(line, ": ")
	if !ok || (source != SourceVerify && source != SourceResolve) {
		return StepProgress{}, false
	}
	f := strings.SplitN(rest, " ", 5)
	if len(f) < 5 || f[0] != story || (f[2] != LineStep && f[2] != LineTier) {
		return StepProgress{}, false
	}
	at, err := time.Parse(time.RFC3339, f[1])
	if err != nil {
		return StepProgress{}, false
	}
	s := StepProgress{Kind: f[2], Source: source}
	switch f[3] {
	case "started":
		s.Name, s.Started = f[4], at
	case "ended":
		s.Outcome, s.Name, _ = strings.Cut(f[4], " ")
		s.Ended = at
	default:
		return StepProgress{}, false
	}
	return s, s.Name != "" && (s.Ended.IsZero() || s.Outcome != "")
}

// stepAfter is the step story's close-out is on after line of its log, from
// step before it: the step a line that begins one names, StepCommit once
// flai verify said it passed every step, or step when line says neither.
func stepAfter(step, line, story string) string {
	switch {
	case strings.HasPrefix(line, "close-out: flai verify "+story+" --sync-only"):
		return StepSyncCheck
	case strings.HasPrefix(line, "close-out: flai verify "+story+" "):
		return StepVerify
	case line == "close-out: commit":
		return StepCommit
	case strings.HasPrefix(line, "close-out: main moved during the run; flai stream sync "+story):
		return StepSync
	case step == StepVerify && strings.HasPrefix(line, "verify: "+story+" passed every step"):
		// the commit logs a line only when something is outstanding
		return StepCommit
	}
	return step
}

// Wait is how long Ask waits for flai serve.
type Wait struct {
	// Take is how long flai serve has to take the request: 10s.
	Take time.Duration
	// Answer is how long it has to answer once it took it: 10s.
	Answer time.Duration
	// Poll is how often the files are looked at: 250ms.
	Poll time.Duration
}

func (w Wait) orDefaults() Wait {
	if w.Take <= 0 {
		w.Take = 10 * time.Second
	}
	if w.Answer <= 0 {
		w.Answer = 10 * time.Second
	}
	if w.Poll <= 0 {
		w.Poll = 250 * time.Millisecond
	}
	return w
}

// Dir is where close-out requests, answers, jobs, and logs are, under the
// project's cache dir.
func Dir(cacheDir string) string { return filepath.Join(cacheDir, "verify", "jobs") }

// File is story's file of the given suffix in dir, the jobs folder.
func File(dir, story, suffix string) string { return filepath.Join(dir, story+suffix) }

// LogPath is where the output of story's close-out job is kept, under the
// project's cache dir.
func LogPath(cacheDir, story string) string { return File(Dir(cacheDir), story, LogSuffix) }

// Ask asks flai serve for req.Story's close-out and waits for its answer:
// the job's id, or the verdict the verify record holds. A request flai serve
// does not take within w.Take is withdrawn and the error wraps
// ErrServeNotTaking; one it took but did not answer within w.Answer is an
// error saying where to look.
func Ask(ctx context.Context, cacheDir string, req Request, w Wait) (Answer, error) {
	w = w.orDefaults()
	if req.At.IsZero() {
		req.At = time.Now().UTC()
	}
	dir := Dir(cacheDir)
	request, taken, answer := File(dir, req.Story, RequestSuffix), File(dir, req.Story, TakenSuffix), File(dir, req.Story, AnswerSuffix)
	for _, stale := range []string{answer, taken} {
		if err := os.Remove(stale); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return Answer{}, fmt.Errorf("remove the stale %s before asking flai serve for %s's close-out: %w; check that the project's .flai-cache is writable", stale, req.Story, err)
		}
	}
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return Answer{}, fmt.Errorf("write the request for %s's close-out: %w", req.Story, err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Answer{}, fmt.Errorf("write the request for %s's close-out: %w; check that the project's .flai-cache is writable", req.Story, err)
	}
	if err := atomicfile.WriteFile(request, append(data, '\n'), 0o644); err != nil {
		return Answer{}, fmt.Errorf("write the request for %s's close-out at %s: %w; check that the project's .flai-cache is writable", req.Story, request, err)
	}
	asked := time.Now()
	var tookAt time.Time
	tick := time.NewTicker(w.Poll)
	defer tick.Stop()
	for {
		if a, ok, err := readAnswer(answer, req.Story); err != nil || ok {
			_ = os.Remove(answer)
			_ = os.Remove(taken)
			return a, err
		}
		switch {
		case tookAt.IsZero() && !exists(request):
			tookAt = time.Now()
		case tookAt.IsZero() && time.Since(asked) >= w.Take:
			if err := os.Remove(request); err == nil {
				return Answer{}, fmt.Errorf("flai serve did not take the request for %s's close-out within %s, so nothing runs: %w", req.Story, w.Take, ErrServeNotTaking)
			}
			tookAt = time.Now() // taken as the wait ran out
		case !tookAt.IsZero() && time.Since(tookAt) >= w.Answer:
			return Answer{}, fmt.Errorf("flai serve took the request for %s's close-out but did not answer within %s; its log says why, and a job it started ends with a verified event", req.Story, w.Answer)
		}
		select {
		case <-ctx.Done():
			if tookAt.IsZero() && os.Remove(request) == nil {
				return Answer{}, fmt.Errorf("stopped waiting for flai serve to take the request for %s's close-out, and withdrew it: %w", req.Story, ctx.Err())
			}
			return Answer{}, fmt.Errorf("stopped waiting for flai serve's answer to the request for %s's close-out, which it took: %w; a job it started ends with a verified event", req.Story, ctx.Err())
		case <-tick.C:
		}
	}
}

// exists reports whether path names a file or folder that can be seen.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// readAnswer is the answer at path, and false with no error while there is none.
func readAnswer(path, story string) (Answer, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Answer{}, false, nil
	}
	if err != nil {
		return Answer{}, false, fmt.Errorf("read flai serve's answer to the request for %s's close-out at %s: %w; check that the project's .flai-cache is readable", story, path, err)
	}
	var a Answer
	if err := json.Unmarshal(data, &a); err != nil {
		return Answer{}, false, fmt.Errorf("read flai serve's answer to the request for %s's close-out at %s: %w; ask again, and a job running is answered, not started twice", story, path, err)
	}
	return a, true, nil
}

// WriteAnswer writes a in dir, the jobs folder, where the asker waits for it.
func WriteAnswer(dir string, a Answer) error {
	path := File(dir, a.Story, AnswerSuffix)
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("answer the request for %s's close-out: %w", a.Story, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("answer the request for %s's close-out at %s: %w; check that the project's .flai-cache is writable", a.Story, path, err)
	}
	return nil
}

// ReadJob is story's close-out job under the project's cache dir and true,
// or false with no error when it has had none.
func ReadJob(cacheDir, story string) (Job, bool, error) {
	return ReadJobAt(File(Dir(cacheDir), story, JobSuffix))
}

// ReadJobAt is the close-out job whose state is at path and true, or false
// with no error when there is none.
func ReadJobAt(path string) (Job, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Job{}, false, nil
	}
	if err != nil {
		return Job{}, false, fmt.Errorf("read the close-out job at %s: %w; check that the project's .flai-cache is readable", path, err)
	}
	var j Job
	if err := json.Unmarshal(data, &j); err != nil {
		return Job{}, false, fmt.Errorf("read the close-out job at %s: %w; delete it and ask for the close-out again", path, err)
	}
	return j, true, nil
}

// WriteJob writes j's state in dir, the jobs folder, in one step.
func WriteJob(dir string, j Job) error {
	path := File(dir, j.Story, JobSuffix)
	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return fmt.Errorf("write the state of close-out job %s: %w", j.ID, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write the state of close-out job %s at %s: %w; check that the project's .flai-cache is writable", j.ID, path, err)
	}
	return nil
}
