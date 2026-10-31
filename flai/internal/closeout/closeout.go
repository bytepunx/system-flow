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

// ErrServeNotTaking is a close-out request no flai serve took in time.
var ErrServeNotTaking = errors.New("flai serve is not running for this project: run the close-out script, or flai verify where the project has none, yourself")

// Request is a story's agent asking flai serve for its close-out.
type Request struct {
	// Story is the story's ID, zero padded.
	Story string `json:"story"`
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

// Job is the state of a story's close-out job, S-nnnn.job.json.
type Job struct {
	ID      string    `json:"id"`
	Story   string    `json:"story"`
	By      string    `json:"by"`
	Message string    `json:"message,omitempty"`
	All     bool      `json:"all,omitempty"`
	Fresh   bool      `json:"fresh,omitempty"`
	Started time.Time `json:"started"`
	// Ended and the rest are set once the job ended: ExitCode is the
	// close-out's, -1 when it could not run or was interrupted; Commit is
	// the worktree's head at the end; DurationMS is how long it ran.
	Ended      time.Time `json:"ended,omitzero"`
	ExitCode   *int      `json:"exit_code,omitempty"`
	Verdict    string    `json:"verdict,omitempty"`
	StoppedAt  string    `json:"stopped_at,omitempty"`
	Commit     string    `json:"commit,omitempty"`
	DurationMS int64     `json:"duration_ms,omitempty"`
}

// Running reports whether the job has not ended.
func (j Job) Running() bool { return j.Ended.IsZero() }

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
