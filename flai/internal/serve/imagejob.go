package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// The image job (ADR-0295 §5): flai serve's host job that pulls or builds the
// agent images a start finds missing or stale, one at a time per host, while
// every start that needs them waits for it.
const (
	// imageJobAction and imageJobMethod are the journal's action and method for an image job.
	imageJobAction = "image"
	imageJobMethod = "image.job"
	// imageJobsFolder is the folder of the image jobs under flai serve's state folder.
	imageJobsFolder = "image-jobs"
	// imageJobSuffix and imageJobLogSuffix end an image job's record and its log.
	imageJobSuffix    = ".job.json"
	imageJobLogSuffix = ".log"
)

// The outcomes of an image job's end.
const (
	imageJobPassed  = "passed"
	imageJobFailed  = "failed"
	imageJobStopped = "stopped"
)

// imageJobsDir is the folder of the image jobs' records and logs, beside dashboard-digests.
func (d Dir) imageJobsDir() string { return filepath.Join(string(d), imageJobsFolder) }

// ImageEnsure makes the agent images of the project at root that are missing
// or stale, writing what it does to log, and is what it did with each: the
// work of one image job, flai image ensure's for the images named.
type ImageEnsure func(ctx context.Context, root string, images []agentimage.Image, log io.Writer) ([]agentimage.Outcome, error)

// imageNeed is what a start asks an image job for: the project's key, for the
// journal, its main checkout's root, and the images it runs in.
type imageNeed struct {
	Key    string
	Root   string
	Images []agentimage.Image
}

// imageJobRecord is an image job as image-jobs/<id>.job.json keeps it: its
// id, its project, the refs of its images, when it was asked for, started,
// and ended, how it ended, the way each image took, its log, and, when it
// failed, the step that failed and why.
type imageJobRecord struct {
	ID         string        `json:"id"`
	Project    string        `json:"project,omitempty"`
	Root       string        `json:"root"`
	Images     []string      `json:"images"`
	Requested  time.Time     `json:"requested"`
	Started    time.Time     `json:"started,omitzero"`
	Ended      time.Time     `json:"ended,omitzero"`
	Outcome    string        `json:"outcome,omitempty"`
	Ways       []imageJobWay `json:"ways,omitempty"`
	Log        string        `json:"log"`
	Step       string        `json:"step,omitempty"`
	Cause      string        `json:"cause,omitempty"`
	DurationMS int64         `json:"duration_ms,omitempty"`
}

// imageJobWay is the way an image job left one image.
type imageJobWay struct {
	Image string         `json:"image"`
	Way   agentimage.Way `json:"way"`
}

// imageJobResult is what a passed image job made: its id, what it did with
// each image, its log, and how long it ran.
type imageJobResult struct {
	ID       string
	Outcomes []agentimage.Outcome
	Log      string
	Duration time.Duration
}

// ImageJobError is an image job that failed, or that flai serve stopped: its
// id, the project's root, its images, the step that failed ("" when no
// docker step did), the cause, its log, and the error itself.
type ImageJobError struct {
	ID      string
	Root    string
	Images  []string
	Step    string
	Cause   string
	Log     string
	Stopped bool
	Err     error
}

// Error says which job failed, at which step, why, where its log is, and how to see it again.
func (e *ImageJobError) Error() string {
	images := strings.Join(e.Images, ", ")
	if e.Stopped {
		return fmt.Sprintf("the image job %s, which makes %s for %s, was stopped as flai serve stopped before it ended; its log is %s; start again once flai serve runs, or run flai image ensure in the project by hand", e.ID, images, e.Root, e.Log)
	}
	return fmt.Sprintf("the image job %s, which makes %s for %s, failed%s: %s; read its log %s, mend what it names, and run flai image ensure in the project to see it again", e.ID, images, e.Root, atStep(e.Step), e.Cause, e.Log)
}

// Unwrap is the job's own error.
func (e *ImageJobError) Unwrap() error { return e.Err }

// imageJob is one image job, queued, running, or ended. Its id, project,
// images, and log are fixed when it is asked for; its record is the runner's
// to change, under the runner's lock; its result and error are set before
// done is closed and never change after.
type imageJob struct {
	id   string
	need imageNeed
	refs []string
	log  string

	rec imageJobRecord

	done   chan struct{}
	result imageJobResult
	err    error
}

// ID is the job's id, known as soon as it is asked for.
func (j *imageJob) ID() string { return j.id }

// Log is the path of the job's log.
func (j *imageJob) Log() string { return j.log }

// Done is closed when the job ends.
func (j *imageJob) Done() <-chan struct{} { return j.done }

// wait waits for the job to end and is what it made, or an *ImageJobError
// when it failed or flai serve stopped it. A waiter whose ctx ends first
// stops waiting; the job goes on for the others.
func (j *imageJob) wait(ctx context.Context) (imageJobResult, error) {
	select {
	case <-j.done:
		return j.result, j.err
	case <-ctx.Done():
		return imageJobResult{ID: j.id, Log: j.log}, fmt.Errorf("wait for the image job %s, which makes %s for %s: %w; the job goes on, and its log is %s", j.id, strings.Join(j.refs, ", "), j.need.Root, ctx.Err(), j.log)
	}
}

// covers reports whether the job makes every image of refs for the project at root.
func (j *imageJob) covers(root string, refs []string) bool {
	if j.need.Root != root {
		return false
	}
	for _, ref := range refs {
		if !slices.Contains(j.refs, ref) {
			return false
		}
	}
	return true
}

// imageJobOptions configure an image job runner.
type imageJobOptions struct {
	// Dir is the folder of the jobs' records and logs, Dir.imageJobsDir.
	Dir string
	// Ensure does each job's work.
	Ensure ImageEnsure
	// Now is the clock; time.Now when nil.
	Now func() time.Time
	// Logger takes the jobs' log events; discarded when nil.
	Logger *slog.Logger
	// Record writes an entry of the host's journal; nil journals nothing.
	Record func(hostapi.Entry)
}

// imageJobRunner runs the image jobs of every project flai serve serves, one
// at a time, in the order they were asked for, until its context ends.
type imageJobRunner struct {
	ctx     context.Context
	o       imageJobOptions
	wg      sync.WaitGroup
	mu      sync.Mutex
	running *imageJob
	queue   []*imageJob
	last    map[string]imageJobRecord // the last ended job of each project, by root
	issued  int
}

// newImageJobRunner is a runner whose jobs run until ctx, flai serve's own, ends.
func newImageJobRunner(ctx context.Context, o imageJobOptions) *imageJobRunner {
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.Logger == nil {
		o.Logger = slog.New(slog.DiscardHandler)
	}
	return &imageJobRunner{ctx: ctx, o: o, last: map[string]imageJobRecord{}}
}

// request is the job that makes need's images: the running or queued job of
// need's project that makes them all, which it joins, or else a new job,
// queued behind those asked for before it. A request once flai serve has
// stopped is a job ended stopped, which never runs.
func (r *imageJobRunner) request(need imageNeed) *imageJob {
	need.Root = filepath.Clean(need.Root)
	refs := imageJobRefs(need.Images)
	r.mu.Lock()
	defer r.mu.Unlock()
	if j := r.covering(need.Root, refs); j != nil {
		r.o.Logger.Info("image job joined", "component", "serve", "root", need.Root, "project", need.Key, "job_id", j.id, "images", strings.Join(refs, ","))
		return j
	}
	at := r.o.Now().UTC()
	id := r.newID(at)
	j := &imageJob{id: id, need: need, refs: refs, log: r.path(id, imageJobLogSuffix), done: make(chan struct{})}
	j.rec = imageJobRecord{ID: id, Project: need.Key, Root: need.Root, Images: refs, Requested: at, Log: j.log}
	r.write(j)
	if r.ctx.Err() != nil {
		r.endUnrun(j)
		return j
	}
	r.queue = append(r.queue, j)
	r.o.Logger.Info("image job queued", "component", "serve", "root", need.Root, "project", need.Key, "job_id", id, "images", strings.Join(refs, ","), "ahead", len(r.queue)-1, "running", r.running != nil)
	r.next()
	return j
}

// lastJob is the record of the last job of the project at root that ended, if any.
func (r *imageJobRunner) lastJob(root string) (imageJobRecord, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.last[filepath.Clean(root)]
	return rec, ok
}

// drain waits for the running job's work to return.
func (r *imageJobRunner) drain() { r.wg.Wait() }

// covering is the running or queued job of root that makes every image of refs, or nil.
func (r *imageJobRunner) covering(root string, refs []string) *imageJob {
	if r.running != nil && r.running.covers(root, refs) {
		return r.running
	}
	for _, j := range r.queue {
		if j.covers(root, refs) {
			return j
		}
	}
	return nil
}

// newID is a new job's id: the UTC stamp of at and a counter, sortable, and
// never one whose record is already in the folder.
func (r *imageJobRunner) newID(at time.Time) string {
	stamp := at.Format("20060102T150405Z")
	for {
		r.issued++
		id := fmt.Sprintf("%s-%03d", stamp, r.issued)
		if _, err := os.Stat(r.path(id, imageJobSuffix)); err != nil {
			return id
		}
	}
}

// path is the file of job id that ends in suffix.
func (r *imageJobRunner) path(id, suffix string) string {
	return filepath.Join(r.o.Dir, id+suffix)
}

// next starts the first queued job when none runs, and ends every queued job
// stopped once flai serve has stopped. The lock is held.
func (r *imageJobRunner) next() {
	if r.ctx.Err() != nil {
		for _, j := range r.queue {
			r.endUnrun(j)
		}
		r.queue = nil
		return
	}
	if r.running != nil || len(r.queue) == 0 {
		return
	}
	j := r.queue[0]
	r.queue = r.queue[1:]
	r.running = j
	started := r.o.Now().UTC()
	j.rec.Started = started
	r.write(j)
	images := strings.Join(j.refs, ", ")
	r.journal(j, hostapi.Entry{At: started.Format(time.RFC3339), Outcome: "started",
		Detail: fmt.Sprintf("image job %s for %s: started, log %s", j.id, images, j.log)})
	r.o.Logger.Info("image job started", "component", "serve", "root", j.need.Root, "project", j.need.Key, "job_id", j.id, "images", strings.Join(j.refs, ","), "log", j.log)
	r.wg.Add(1)
	go r.work(j)
}

// work runs job j's ensure, its output in its log, then ends it and starts the next.
func (r *imageJobRunner) work(j *imageJob) {
	defer r.wg.Done()
	outcomes, err := r.ensure(j)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.end(j, outcomes, err)
	r.running = nil
	r.next()
}

// ensure runs j's ensure with flai serve's context, writing to j's log.
func (r *imageJobRunner) ensure(j *imageJob) ([]agentimage.Outcome, error) {
	if err := os.MkdirAll(r.o.Dir, 0o700); err != nil {
		return nil, fmt.Errorf("make the folder of the image job's log %s: %w; check that flai serve's state folder can be written", r.o.Dir, err)
	}
	f, err := os.OpenFile(j.log, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, fmt.Errorf("create the image job's log %s: %w; check that flai serve's state folder can be written", j.log, err)
	}
	_, _ = fmt.Fprintf(f, "image job %s: making %s for %s\n", j.id, strings.Join(j.refs, ", "), j.need.Root)
	outcomes, err := r.o.Ensure(r.ctx, j.need.Root, slices.Clone(j.need.Images), f)
	if err != nil {
		_, _ = fmt.Fprintf(f, "image job %s: failed: %v\n", j.id, err)
	} else {
		_, _ = fmt.Fprintf(f, "image job %s: passed\n", j.id)
	}
	if cerr := f.Close(); cerr != nil && err == nil {
		err = fmt.Errorf("close the image job's log %s: %w; check that flai serve's state folder can be written", j.log, cerr)
	}
	return outcomes, err
}

// end ends j, run, with what its ensure did: it records, journals, and logs
// the end, keeps it as its project's last, and tells its waiters. The lock is held.
func (r *imageJobRunner) end(j *imageJob, outcomes []agentimage.Outcome, err error) {
	ended := r.o.Now().UTC()
	j.rec.Ended = ended
	j.rec.DurationMS = ended.Sub(j.rec.Started).Milliseconds()
	for _, o := range outcomes {
		j.rec.Ways = append(j.rec.Ways, imageJobWay{Image: o.Image.Ref, Way: o.Way})
	}
	switch {
	case err == nil:
		j.rec.Outcome = imageJobPassed
	case r.ctx.Err() != nil:
		j.rec.Outcome, j.rec.Cause = imageJobStopped, err.Error()
		j.err = &ImageJobError{ID: j.id, Root: j.need.Root, Images: j.refs, Cause: j.rec.Cause, Log: j.log, Stopped: true, Err: err}
	default:
		j.rec.Outcome, j.rec.Step, j.rec.Cause = imageJobFailed, failedStep(err), err.Error()
		j.err = &ImageJobError{ID: j.id, Root: j.need.Root, Images: j.refs, Step: j.rec.Step, Cause: j.rec.Cause, Log: j.log, Err: err}
	}
	r.write(j)
	duration := time.Duration(j.rec.DurationMS) * time.Millisecond
	ways := imageJobWays(j.rec.Ways)
	detail := fmt.Sprintf("image job %s for %s: %s in %s", j.id, strings.Join(j.refs, ", "), j.rec.Outcome, duration.Round(time.Second))
	if ways != "" {
		detail += " (" + strings.ReplaceAll(ways, ",", ", ") + ")"
	}
	if err != nil {
		detail += fmt.Sprintf("%s: %s", atStep(j.rec.Step), j.rec.Cause)
	}
	r.journal(j, hostapi.Entry{At: ended.Format(time.RFC3339), Outcome: imageJobJournalOutcome(j.rec.Outcome), Verdict: j.rec.Outcome,
		DurationMS: j.rec.DurationMS, Detail: detail + ", log " + j.log})
	level, attrs := r.o.Logger.Info, []any{"component", "serve", "root", j.need.Root, "project", j.need.Key, "job_id", j.id,
		"images", strings.Join(j.refs, ","), "outcome", j.rec.Outcome, "ways", ways, "duration_ms", j.rec.DurationMS, "log", j.log}
	if j.rec.Outcome == imageJobFailed {
		level, attrs = r.o.Logger.Warn, append(attrs, "step", j.rec.Step, "cause", j.rec.Cause)
	}
	level("image job ended", attrs...)
	j.result = imageJobResult{ID: j.id, Outcomes: outcomes, Log: j.log, Duration: duration}
	r.last[j.need.Root] = j.rec
	close(j.done)
}

// endUnrun ends j stopped without running it, as flai serve stopped before
// its turn came. The lock is held.
func (r *imageJobRunner) endUnrun(j *imageJob) {
	j.rec.Ended, j.rec.Outcome = r.o.Now().UTC(), imageJobStopped
	r.write(j)
	cause := r.ctx.Err()
	j.err = &ImageJobError{ID: j.id, Root: j.need.Root, Images: j.refs, Cause: cause.Error(), Log: j.log, Stopped: true, Err: cause}
	r.o.Logger.Info("image job not run as flai serve stopped", "component", "serve", "root", j.need.Root, "project", j.need.Key, "job_id", j.id, "images", strings.Join(j.refs, ","))
	j.result = imageJobResult{ID: j.id, Log: j.log}
	r.last[j.need.Root] = j.rec
	close(j.done)
}

// write writes j's record; a failure loses the record, not the job, and is logged. The lock is held.
func (r *imageJobRunner) write(j *imageJob) {
	path := r.path(j.id, imageJobSuffix)
	err := os.MkdirAll(r.o.Dir, 0o700)
	if err == nil {
		var data []byte
		data, err = json.MarshalIndent(j.rec, "", "  ")
		if err == nil {
			err = atomicfile.WriteFile(path, append(data, '\n'), 0o600)
		}
	}
	if err != nil {
		r.o.Logger.Warn("image job's record not written", "component", "serve", "root", j.need.Root, "job_id", j.id, "path", path, "err", err.Error())
	}
}

// journal records an image job entry of j in the host journal.
func (r *imageJobRunner) journal(j *imageJob, e hostapi.Entry) {
	if r.o.Record == nil {
		return
	}
	e.Action, e.Method, e.Project, e.Root, e.By = imageJobAction, imageJobMethod, j.need.Key, j.need.Root, "flai serve"
	r.o.Record(e)
}

// imageJobJournalOutcome is the journal outcome of an image job's end: done for a pass, failed otherwise.
func imageJobJournalOutcome(outcome string) string {
	if outcome == imageJobPassed {
		return "done"
	}
	return "failed"
}

// imageJobRefs are the refs of images, sorted, each once, without the empty ref of an image with none.
func imageJobRefs(images []agentimage.Image) []string {
	refs := make([]string, 0, len(images))
	for _, img := range images {
		if img.Ref != "" {
			refs = append(refs, img.Ref)
		}
	}
	slices.Sort(refs)
	return slices.Compact(refs)
}

// imageJobWays is ways as image=way pairs, joined by commas.
func imageJobWays(ways []imageJobWay) string {
	pairs := make([]string, 0, len(ways))
	for _, w := range ways {
		pairs = append(pairs, w.Image+"="+string(w.Way))
	}
	return strings.Join(pairs, ",")
}

// failedStep is the docker step err names as failing: the
// *agentimage.StepError or *agentimage.StallError of the last error err
// joins that holds one, which is the build's when a build stood in for a
// failed pull; "" when it names none.
func failedStep(err error) string {
	var joined interface{ Unwrap() []error }
	if errors.As(err, &joined) {
		errs := joined.Unwrap()
		for i := len(errs) - 1; i >= 0; i-- {
			if step := failedStep(errs[i]); step != "" {
				return step
			}
		}
	}
	var stepErr *agentimage.StepError
	if errors.As(err, &stepErr) {
		return stepErr.Step
	}
	var stall *agentimage.StallError
	if errors.As(err, &stall) {
		return stall.Step
	}
	return ""
}
