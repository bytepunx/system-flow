package kluster

// The hand-off of a story's kluster job to flai serve (ADR-0292 §1, §5): the
// files a story's agent writes to ask flai serve to deploy or tear down its
// story's namespace, and those flai serve writes to answer it and to keep the
// job's state, under the project's .flai-cache/kluster/requests, as
// internal/closeout is for a close-out. The asker, the MCP tool or flai
// kluster deploy --job, writes S-nnnn.<kind>.request.json; flai serve takes it
// by renaming it to S-nnnn.<kind>.taken.json and answers in
// S-nnnn.<kind>.answer.json with the job's id and the jobs ahead of it, or a
// refusal. The job's state is S-nnnn.job.json, one per story.

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

// The kinds of job a story's agent asks flai serve for.
const (
	// KindDeploy builds, loads, and deploys the story at its branch's head.
	KindDeploy = "deploy"
	// KindTeardown tears the story's namespace down.
	KindTeardown = "teardown"
)

// OutcomeRefused is a job flai serve queued and then ran nothing for, as the
// kluster host action turned off while it waited. A job of the queue ends
// OutcomePassed, OutcomeFailed, OutcomeRefused, or OutcomeInterrupted.
const OutcomeRefused = "refused"

// The suffixes of a story's files in the requests folder: its request, the
// request once flai serve took it, and flai serve's answer, each after the
// kind, and the story's job state.
const (
	RequestSuffix = ".request.json"
	TakenSuffix   = ".taken.json"
	AnswerSuffix  = ".answer.json"
	JobSuffix     = ".job.json"
)

// ErrServeNotTaking is a kluster request no flai serve took in time.
var ErrServeNotTaking = errors.New("flai serve is not running for this project, and a story's agent cannot run a kluster job itself: the operator runs flai kluster deploy or flai kluster teardown in a shell, or starts flai serve, and then you ask again")

// RequestsDir is where kluster requests, answers, and job states are kept for
// the project run from root.
func RequestsDir(root string) string { return filepath.Join(Dir(root), "requests") }

// RequestFile is story's file of kind with suffix, RequestSuffix,
// TakenSuffix, or AnswerSuffix, in the requests folder of the project run
// from root.
func RequestFile(root, story, kind, suffix string) string {
	return filepath.Join(RequestsDir(root), story+"."+kind+suffix)
}

// JobFile is the file of story's job state in the requests folder of the
// project run from root.
func JobFile(root, story string) string { return filepath.Join(RequestsDir(root), story+JobSuffix) }

// JobID is the id of a job of kind for story queued at queued, such as
// S-0660-deploy-20261113T073353Z.
func JobID(story, kind string, queued time.Time) string {
	return story + "-" + kind + "-" + queued.UTC().Format("20060102T150405Z")
}

// Request is a story's agent asking flai serve for a kluster job.
type Request struct {
	// Story is the story's ID, zero padded and upper case.
	Story string `json:"story"`
	// Kind is KindDeploy or KindTeardown.
	Kind string `json:"kind"`
	// By is the agent that asks, FLAI_AGENT of the session.
	By string `json:"by"`
	// At is when it was asked; Ask stamps it when zero.
	At time.Time `json:"at"`
}

// Answer is flai serve's answer to a kluster request.
type Answer struct {
	// Job is the job's id; "" when none is queued or found for the request.
	Job   string `json:"job,omitempty"`
	Kind  string `json:"kind"`
	Story string `json:"story"`
	// Ahead is how many jobs of the queue wait or run before it.
	Ahead int `json:"ahead"`
	// Started is whether this request queued the job; false when a job of the
	// story was waiting or running already and its id is answered.
	Started bool `json:"started"`
	// Recorded is an answer from the namespace's record for a story in review
	// whose deploy passed at its head: Commit and Outcome are the record's,
	// and nothing is queued.
	Recorded bool   `json:"recorded,omitempty"`
	Commit   string `json:"commit,omitempty"`
	Outcome  string `json:"outcome,omitempty"`
	// Log is where the job's output is, or the record's logs.
	Log string `json:"log,omitempty"`
	// Reason is why flai serve queued nothing for the request, saying what to
	// do; "" when it queued a job, found one, or answered from the record.
	Reason string `json:"reason,omitempty"`
}

// Refused reports whether flai serve queued nothing for the request.
func (a Answer) Refused() bool { return a.Reason != "" }

// Found reports whether the answer names a job of the story that was waiting
// or running already, which the request did not queue.
func (a Answer) Found() bool { return a.Job != "" && !a.Started }

// AskedJob is the state of a story's kluster job in flai serve's queue,
// S-nnnn.job.json: the last one asked for the story, by its agent or by the
// review pass. It is not Job, the job that holds the project's kluster lock
// while it runs.
type AskedJob struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Story string `json:"story"`
	// By is who asked: the story's agent, or flai serve for the review pass.
	By       string    `json:"by"`
	QueuedAt time.Time `json:"queued_at"`
	// StartedAt is when it started, once the jobs ahead of it ended.
	StartedAt time.Time `json:"started_at,omitzero"`
	// EndedAt and the rest are set once the job ended: Commit is the head it
	// built or the record's; Outcome is OutcomePassed, OutcomeFailed,
	// OutcomeRefused, or OutcomeInterrupted; FailedStep, ExitCode, and Error
	// say why one failed, ExitCode -1 when its program gave none; Steps are
	// each step's name, outcome, exit code, and log.
	EndedAt    time.Time `json:"ended_at,omitzero"`
	Commit     string    `json:"commit,omitempty"`
	Outcome    string    `json:"outcome,omitempty"`
	FailedStep string    `json:"failed_step,omitempty"`
	ExitCode   *int      `json:"exit_code,omitempty"`
	Error      string    `json:"error,omitempty"`
	Steps      []Step    `json:"steps,omitempty"`
	// Log is where the job's output is: the log of the step it failed at, or
	// the folder of its steps' logs.
	Log string `json:"log,omitempty"`
}

// Waiting reports whether the job is queued and has not started.
func (j AskedJob) Waiting() bool { return j.StartedAt.IsZero() && j.EndedAt.IsZero() }

// Running reports whether the job started and has not ended.
func (j AskedJob) Running() bool { return !j.StartedAt.IsZero() && j.EndedAt.IsZero() }

// Ended reports whether the job ended, whatever its outcome.
func (j AskedJob) Ended() bool { return !j.EndedAt.IsZero() }

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

// checkRequest refuses a request whose story is not a story's ID or whose
// kind flai serve does not run, and gives its story upper case.
func checkRequest(req Request) (Request, error) {
	if !storyID.MatchString(req.Story) {
		return req, fmt.Errorf("%q is not a story's ID, so no kluster job is asked for it; give one such as S-0578", req.Story)
	}
	req.Story = strings.ToUpper(req.Story)
	switch req.Kind {
	case KindDeploy, KindTeardown:
		return req, nil
	}
	return req, fmt.Errorf("%q is not a kind of kluster job flai serve runs for %s; ask for %s or %s", req.Kind, req.Story, KindDeploy, KindTeardown)
}

// Ask asks flai serve for req's kluster job and waits for its answer: the
// job's id and the jobs ahead of it, the record's outcome, or a refusal. A
// request flai serve does not take within w.Take is withdrawn and the error
// wraps ErrServeNotTaking; one it took but did not answer within w.Answer is
// an error saying where to look.
func Ask(ctx context.Context, root string, req Request, w Wait) (Answer, error) {
	req, err := checkRequest(req)
	if err != nil {
		return Answer{}, err
	}
	w = w.orDefaults()
	if req.At.IsZero() {
		req.At = time.Now().UTC()
	}
	what := req.Story + "'s " + req.Kind
	request, taken, answer := RequestFile(root, req.Story, req.Kind, RequestSuffix), RequestFile(root, req.Story, req.Kind, TakenSuffix), RequestFile(root, req.Story, req.Kind, AnswerSuffix)
	for _, stale := range []string{answer, taken} {
		if err := os.Remove(stale); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return Answer{}, fmt.Errorf("remove the stale %s before asking flai serve for %s: %w; check that the project's .flai-cache is writable", stale, what, err)
		}
	}
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return Answer{}, fmt.Errorf("write the request for %s: %w", what, err)
	}
	if err := os.MkdirAll(RequestsDir(root), 0o755); err != nil {
		return Answer{}, fmt.Errorf("write the request for %s: %w; check that the project's .flai-cache is writable", what, err)
	}
	if err := atomicfile.WriteFile(request, append(data, '\n'), 0o644); err != nil {
		return Answer{}, fmt.Errorf("write the request for %s at %s: %w; check that the project's .flai-cache is writable", what, request, err)
	}
	asked := time.Now()
	var tookAt time.Time
	tick := time.NewTicker(w.Poll)
	defer tick.Stop()
	for {
		if a, ok, err := readAnswer(answer, what); err != nil || ok {
			_ = os.Remove(answer)
			_ = os.Remove(taken)
			return a, err
		}
		switch {
		case tookAt.IsZero() && !exists(request):
			tookAt = time.Now()
		case tookAt.IsZero() && time.Since(asked) >= w.Take:
			if err := os.Remove(request); err == nil {
				return Answer{}, fmt.Errorf("flai serve did not take the request for %s within %s, so nothing is queued: %w", what, w.Take, ErrServeNotTaking)
			}
			tookAt = time.Now() // taken as the wait ran out
		case !tookAt.IsZero() && time.Since(tookAt) >= w.Answer:
			return Answer{}, fmt.Errorf("flai serve took the request for %s but did not answer within %s; its log says why, and a job it queued ends with a kluster event", what, w.Answer)
		}
		select {
		case <-ctx.Done():
			if tookAt.IsZero() && os.Remove(request) == nil {
				return Answer{}, fmt.Errorf("stopped waiting for flai serve to take the request for %s, and withdrew it: %w", what, ctx.Err())
			}
			return Answer{}, fmt.Errorf("stopped waiting for flai serve's answer to the request for %s, which it took: %w; a job it queued ends with a kluster event", what, ctx.Err())
		case <-tick.C:
		}
	}
}

// exists reports whether path names a file or folder that can be seen.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// readAnswer is the answer at path to the request for what, and false with
// no error while there is none.
func readAnswer(path, what string) (Answer, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Answer{}, false, nil
	}
	if err != nil {
		return Answer{}, false, fmt.Errorf("read flai serve's answer to the request for %s at %s: %w; check that the project's .flai-cache is readable", what, path, err)
	}
	var a Answer
	if err := json.Unmarshal(data, &a); err != nil {
		return Answer{}, false, fmt.Errorf("read flai serve's answer to the request for %s at %s: %w; ask again, and a job waiting or running is answered, not queued twice", what, path, err)
	}
	return a, true, nil
}

// WriteAnswer writes a in the requests folder of the project run from root,
// where the asker waits for it.
func WriteAnswer(root string, a Answer) error {
	path := RequestFile(root, a.Story, a.Kind, AnswerSuffix)
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("answer the request for %s's %s: %w", a.Story, a.Kind, err)
	}
	if err := os.MkdirAll(RequestsDir(root), 0o755); err != nil {
		return fmt.Errorf("answer the request for %s's %s: %w; check that the project's .flai-cache is writable", a.Story, a.Kind, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("answer the request for %s's %s at %s: %w; check that the project's .flai-cache is writable", a.Story, a.Kind, path, err)
	}
	return nil
}

// Pending lists the requests waiting in the requests folder of the project
// run from root for flai serve to take them, oldest first; none when there is
// no folder. A request that does not read is passed over: Take reports it.
func Pending(root string) ([]Request, error) {
	entries, err := os.ReadDir(RequestsDir(root))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list the kluster requests in %s: %w; check that the project's .flai-cache is readable", RequestsDir(root), err)
	}
	var out []Request
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), RequestSuffix) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(RequestsDir(root), e.Name()))
		var req Request
		if err != nil || json.Unmarshal(data, &req) != nil {
			if story, kind, ok := strings.Cut(strings.TrimSuffix(e.Name(), RequestSuffix), "."); ok {
				out = append(out, Request{Story: story, Kind: kind})
			}
			continue
		}
		out = append(out, req)
	}
	slices.SortStableFunc(out, func(a, b Request) int {
		if c := a.At.Compare(b.At); c != 0 {
			return c
		}
		return strings.Compare(a.Story+"."+a.Kind, b.Story+"."+b.Kind)
	})
	return out, nil
}

// Take takes story's request of kind for flai serve by renaming it to its
// taken file, so that the asker knows it was taken, and reads it. A request
// withdrawn first is false with no error.
func Take(root, story, kind string) (Request, bool, error) {
	taken := RequestFile(root, story, kind, TakenSuffix)
	if err := os.Rename(RequestFile(root, story, kind, RequestSuffix), taken); errors.Is(err, fs.ErrNotExist) {
		return Request{}, false, nil
	} else if err != nil {
		return Request{}, false, fmt.Errorf("take the request for %s's %s: %w; check that the project's .flai-cache is writable", story, kind, err)
	}
	data, err := os.ReadFile(taken)
	if err != nil {
		return Request{}, false, fmt.Errorf("read the request for %s's %s at %s: %w; check that the project's .flai-cache is readable", story, kind, taken, err)
	}
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return Request{}, false, fmt.Errorf("read the request for %s's %s at %s: %w; answer it refused, and its agent asks again", story, kind, taken, err)
	}
	return req, true, nil
}

// ReadAskedJob is story's job state in the project run from root and true,
// or false with no error when it has had none.
func ReadAskedJob(root, story string) (AskedJob, bool, error) {
	path := JobFile(root, story)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AskedJob{}, false, nil
	}
	if err != nil {
		return AskedJob{}, false, fmt.Errorf("read %s's kluster job at %s: %w; check that the project's .flai-cache is readable", story, path, err)
	}
	var j AskedJob
	if err := json.Unmarshal(data, &j); err != nil {
		return AskedJob{}, false, fmt.Errorf("read %s's kluster job at %s: %w; delete it and ask for the job again", story, path, err)
	}
	return j, true, nil
}

// WriteAskedJob writes j's state in the requests folder of the project run
// from root, in one step.
func WriteAskedJob(root string, j AskedJob) error {
	path := JobFile(root, j.Story)
	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return fmt.Errorf("write the state of kluster job %s: %w", j.ID, err)
	}
	if err := os.MkdirAll(RequestsDir(root), 0o755); err != nil {
		return fmt.Errorf("write the state of kluster job %s: %w; check that the project's .flai-cache is writable", j.ID, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write the state of kluster job %s at %s: %w; check that the project's .flai-cache is writable", j.ID, path, err)
	}
	return nil
}

// Await waits at most d for story's job id to end, looking at its state
// every poll (250ms when not above 0), and answers its state and whether it
// ended. d of 0 or less answers the state at once. A story with no job, or
// whose job is another by now, is an error, as is a ctx that ends first.
func Await(ctx context.Context, root, story, id string, d, poll time.Duration) (AskedJob, bool, error) {
	if poll <= 0 {
		poll = 250 * time.Millisecond
	}
	deadline := time.Now().Add(d)
	for {
		j, ok, err := ReadAskedJob(root, story)
		switch {
		case err != nil:
			return AskedJob{}, false, err
		case !ok:
			return AskedJob{}, false, fmt.Errorf("%s has no kluster job, so job %s is not waiting or running; ask for it again", story, id)
		case j.ID != id:
			return j, false, fmt.Errorf("%s's kluster job is %s now, not %s: a later request took its place or followed it; await %s, or read the namespace's record with flai kluster deployments", story, j.ID, id, j.ID)
		case j.Ended() || !time.Now().Before(deadline):
			return j, j.Ended(), nil
		}
		wait := min(poll, time.Until(deadline))
		select {
		case <-ctx.Done():
			return j, false, fmt.Errorf("stopped waiting for %s's kluster job %s, which is still %s: %w; a kluster event says when it ends", story, id, j.state(), ctx.Err())
		case <-time.After(wait):
		}
	}
}

// state says in a word whether j waits, runs, or ended.
func (j AskedJob) state() string {
	switch {
	case j.Ended():
		return "ended"
	case j.Running():
		return "running"
	}
	return "queued"
}
