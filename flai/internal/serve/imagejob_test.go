package serve

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// imageJobCall is one call of the fake ensure: what it was given, and where
// the test answers it.
type imageJobCall struct {
	ctx    context.Context
	root   string
	images []agentimage.Image
	log    io.Writer
	answer chan imageJobAnswer
}

// imageJobAnswer is what the test has a call of the fake ensure return.
type imageJobAnswer struct {
	outcomes []agentimage.Outcome
	err      error
}

// pass has the call return outcomes, a way for each of its images.
func (c imageJobCall) pass(way agentimage.Way) []agentimage.Outcome {
	outcomes := make([]agentimage.Outcome, 0, len(c.images))
	for _, img := range c.images {
		outcomes = append(outcomes, agentimage.Outcome{Image: img, Way: way})
	}
	c.answer <- imageJobAnswer{outcomes: outcomes}
	return outcomes
}

// imageJobBuffer is a log the runner writes to from more than one goroutine.
type imageJobBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *imageJobBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// events are the log's events, decoded.
func (b *imageJobBuffer) events(t *testing.T) []map[string]any {
	t.Helper()
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []map[string]any
	sc := bufio.NewScanner(bytes.NewReader(b.buf.Bytes()))
	for sc.Scan() {
		var ev map[string]any
		if err := json.Unmarshal(sc.Bytes(), &ev); err != nil {
			t.Fatalf("log line %q is not JSON: %v", sc.Text(), err)
		}
		out = append(out, ev)
	}
	return out
}

// imageJobLab is an image job runner over a fake ensure, a ticking clock, a
// JSON log, and a journal, in a temporary folder.
type imageJobLab struct {
	t      *testing.T
	dir    string
	calls  chan imageJobCall
	runner *imageJobRunner
	stop   context.CancelFunc
	logs   *imageJobBuffer

	mu      sync.Mutex
	entries []hostapi.Entry
}

func newImageJobLab(t *testing.T) *imageJobLab {
	t.Helper()
	ctx, stop := context.WithCancel(context.Background())
	lab := &imageJobLab{t: t, dir: Dir(t.TempDir()).imageJobsDir(), calls: make(chan imageJobCall, 16), stop: stop, logs: &imageJobBuffer{}}
	var clockMu sync.Mutex
	clock := time.Date(2026, 11, 13, 9, 30, 0, 0, time.UTC)
	now := func() time.Time {
		clockMu.Lock()
		defer clockMu.Unlock()
		clock = clock.Add(250 * time.Millisecond)
		return clock
	}
	ensure := func(ctx context.Context, root string, images []agentimage.Image, log io.Writer) ([]agentimage.Outcome, error) {
		c := imageJobCall{ctx: ctx, root: root, images: images, log: log, answer: make(chan imageJobAnswer, 1)}
		_, _ = fmt.Fprintf(log, "fake ensure of %d images\n", len(images))
		lab.calls <- c
		select {
		case a := <-c.answer:
			return a.outcomes, a.err
		case <-ctx.Done():
			return nil, fmt.Errorf("pull stopped: %w", ctx.Err())
		}
	}
	lab.runner = newImageJobRunner(ctx, imageJobOptions{
		Dir: lab.dir, Ensure: ensure, Now: now,
		Logger: slog.New(slog.NewJSONHandler(lab.logs, nil)),
		Record: func(e hostapi.Entry) {
			lab.mu.Lock()
			defer lab.mu.Unlock()
			lab.entries = append(lab.entries, e)
		},
	})
	t.Cleanup(func() {
		stop()
		lab.runner.drain()
	})
	return lab
}

// need is what a start of the project at root asks for, its images named by refs.
func imageJobNeed(root string, refs ...string) imageNeed {
	images := make([]agentimage.Image, 0, len(refs))
	for _, ref := range refs {
		images = append(images, agentimage.Image{Kind: agentimage.KindHarness, Inner: "claude-code", Ref: ref, Source: agentimage.SourceDefault})
	}
	return imageNeed{Key: filepath.Base(root), Root: root, Images: images}
}

// nextCall is the next call of the fake ensure.
func (lab *imageJobLab) nextCall() imageJobCall {
	lab.t.Helper()
	select {
	case c := <-lab.calls:
		return c
	case <-time.After(5 * time.Second):
		lab.t.Fatal("the fake ensure was never called")
		return imageJobCall{}
	}
}

// noCall fails when the fake ensure has a call waiting.
func (lab *imageJobLab) noCall() {
	lab.t.Helper()
	select {
	case c := <-lab.calls:
		lab.t.Fatalf("the fake ensure was called for %v, want no call", c.images)
	default:
	}
}

// started reports whether job j has started.
func (lab *imageJobLab) started(j *imageJob) bool {
	lab.runner.mu.Lock()
	defer lab.runner.mu.Unlock()
	return !j.rec.Started.IsZero()
}

// journal is the image job entries journaled so far.
func (lab *imageJobLab) journal() []hostapi.Entry {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	var out []hostapi.Entry
	for _, e := range lab.entries {
		if e.Action == imageJobAction && e.Method == imageJobMethod {
			out = append(out, e)
		}
	}
	return out
}

// waitJob waits for j and fails the test if it takes long.
func imageJobWait(t *testing.T, j *imageJob) (imageJobResult, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := j.wait(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("image job %s never ended", j.ID())
	}
	return res, err
}

// readImageJob reads the record of job id from the lab's folder.
func (lab *imageJobLab) readImageJob(id string) imageJobRecord {
	lab.t.Helper()
	data, err := os.ReadFile(filepath.Join(lab.dir, id+imageJobSuffix))
	if err != nil {
		lab.t.Fatal(err)
	}
	var rec imageJobRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		lab.t.Fatal(err)
	}
	return rec
}

func TestImageJobRunsOneAtATime(t *testing.T) {
	lab := newImageJobLab(t)
	first := lab.runner.request(imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0"))
	call := lab.nextCall()
	second := lab.runner.request(imageJobNeed("/p/two", "flai-agent-claude-code:2.0.0"))
	if second.ID() == first.ID() {
		t.Fatalf("a request for other images joined job %s", first.ID())
	}
	lab.noCall()
	if lab.started(second) {
		t.Fatal("the second job started while the first ran")
	}
	call.pass(agentimage.WayPulled)
	if _, err := imageJobWait(t, first); err != nil {
		t.Fatalf("first job: %v", err)
	}
	next := lab.nextCall()
	if next.root != "/p/two" || next.images[0].Ref != "flai-agent-claude-code:2.0.0" {
		t.Fatalf("second ensure ran for %s %v, want /p/two's image", next.root, next.images)
	}
	next.pass(agentimage.WayBuilt)
	if _, err := imageJobWait(t, second); err != nil {
		t.Fatalf("second job: %v", err)
	}
}

func TestImageJobJoinsARequestForTheSameImages(t *testing.T) {
	lab := newImageJobLab(t)
	first := lab.runner.request(imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0", "p-layer:abc"))
	call := lab.nextCall()
	same := lab.runner.request(imageJobNeed("/p/one/", "p-layer:abc", "flai-agent-claude-code:1.0.0"))
	harnessOnly := lab.runner.request(imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0"))
	if same != first || harnessOnly != first {
		t.Fatalf("requests for images job %s makes got jobs %s and %s, want it", first.ID(), same.ID(), harnessOnly.ID())
	}
	want := call.pass(agentimage.WayPulled)
	for _, j := range []*imageJob{first, same, harnessOnly} {
		res, err := imageJobWait(t, j)
		if err != nil {
			t.Fatalf("joined job: %v", err)
		}
		if res.ID != first.ID() || len(res.Outcomes) != len(want) || res.Outcomes[0].Way != agentimage.WayPulled {
			t.Fatalf("joined result %+v, want job %s's %d pulled outcomes", res, first.ID(), len(want))
		}
	}
	lab.runner.drain()
	lab.noCall()
}

func TestImageJobQueuesARequestForOtherImagesInOrder(t *testing.T) {
	lab := newImageJobLab(t)
	a := lab.runner.request(imageJobNeed("/p/one", "a:1"))
	callA := lab.nextCall()
	b := lab.runner.request(imageJobNeed("/p/one", "b:1"))
	c := lab.runner.request(imageJobNeed("/p/one", "c:1"))
	again := lab.runner.request(imageJobNeed("/p/one", "b:1"))
	if again != b {
		t.Fatalf("a second request for b:1 got job %s, want the queued job %s", again.ID(), b.ID())
	}
	if a.ID() >= b.ID() || b.ID() >= c.ID() {
		t.Fatalf("job ids %s, %s, %s do not sort in the order they were asked for", a.ID(), b.ID(), c.ID())
	}
	callA.pass(agentimage.WayCurrent)
	for _, want := range []struct {
		ref string
		job *imageJob
	}{{"b:1", b}, {"c:1", c}} {
		call := lab.nextCall()
		if call.images[0].Ref != want.ref {
			t.Fatalf("ensure ran for %s, want %s next", call.images[0].Ref, want.ref)
		}
		call.pass(agentimage.WayBuilt)
		if _, err := imageJobWait(t, want.job); err != nil {
			t.Fatalf("job for %s: %v", want.ref, err)
		}
	}
}

func TestImageJobPasses(t *testing.T) {
	lab := newImageJobLab(t)
	j := lab.runner.request(imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0", "p-layer:abc"))
	if j.Log() != filepath.Join(lab.dir, j.ID()+imageJobLogSuffix) {
		t.Fatalf("log %s, want %s under the jobs' folder", j.Log(), j.ID()+imageJobLogSuffix)
	}
	queued := lab.readImageJob(j.ID())
	if queued.ID != j.ID() || queued.Root != "/p/one" || queued.Requested.IsZero() {
		t.Fatalf("record on request %+v, want its id, root, and request time", queued)
	}
	call := lab.nextCall()
	call.pass(agentimage.WayLayered)
	res, err := imageJobWait(t, j)
	if err != nil {
		t.Fatal(err)
	}
	if res.Duration <= 0 || res.Log != j.Log() || len(res.Outcomes) != 2 {
		t.Fatalf("result %+v, want a duration, the log, and two outcomes", res)
	}
	rec := lab.readImageJob(j.ID())
	if rec.Outcome != imageJobPassed || rec.Started.IsZero() || rec.Ended.IsZero() || rec.DurationMS <= 0 || len(rec.Ways) != 2 || rec.Ways[0].Way != agentimage.WayLayered {
		t.Fatalf("record %+v, want a passed job with its times, duration, and ways", rec)
	}
	if rec.Step != "" || rec.Cause != "" {
		t.Fatalf("a passed job's record has step %q and cause %q", rec.Step, rec.Cause)
	}
	log, err := os.ReadFile(j.Log())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(log), "fake ensure of 2 images") {
		t.Fatalf("log %q does not hold the ensure's output", log)
	}
	if last, ok := lab.runner.lastJob("/p/one"); !ok || last.ID != j.ID() || last.Outcome != imageJobPassed {
		t.Fatalf("last job %+v, %v, want %s passed", last, ok, j.ID())
	}
}

func TestImageJobFailureCarriesStepCauseAndLog(t *testing.T) {
	lab := newImageJobLab(t)
	need := imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0")
	j := lab.runner.request(need)
	pull := &agentimage.StepError{Step: "pull ghcr.io/bytepunx/flai-agent-claude-code", Tries: 3, Said: "connection reset by peer", Network: true, Err: errors.New("exit status 1")}
	build := &agentimage.StallError{Step: "build flai-agent-claude-code:1.0.0", Limit: 10 * time.Minute}
	lab.nextCall().answer <- imageJobAnswer{err: fmt.Errorf("make the agent image flai-agent-claude-code:1.0.0: its pull failed (%w), and the build on this host that stands in for it failed too: %w", pull, build)}
	_, err := imageJobWait(t, j)
	var failed *ImageJobError
	if !errors.As(err, &failed) {
		t.Fatalf("error %v, want an *ImageJobError", err)
	}
	if failed.ID != j.ID() || failed.Step != build.Step || failed.Log != j.Log() || failed.Stopped || !strings.Contains(failed.Cause, "connection reset by peer") {
		t.Fatalf("failure %+v, want job %s failed at the build's step with the cause and log", failed, j.ID())
	}
	if !errors.As(err, &build) {
		t.Fatal("the failure does not unwrap to the ensure's error")
	}
	msg := err.Error()
	for _, want := range []string{j.ID(), build.Step, j.Log(), "flai image ensure", "connection reset by peer"} {
		if !strings.Contains(msg, want) {
			t.Errorf("failure %q does not name %q", msg, want)
		}
	}
	if strings.Contains(msg, "make agent-images") {
		t.Errorf("failure %q names make agent-images", msg)
	}
	rec := lab.readImageJob(j.ID())
	if rec.Outcome != imageJobFailed || rec.Step != build.Step || rec.Cause == "" {
		t.Fatalf("record %+v, want a failed job with its step and cause", rec)
	}

	again := lab.runner.request(need)
	if again.ID() == j.ID() {
		t.Fatalf("a request after the failure joined the failed job %s", j.ID())
	}
	lab.nextCall().pass(agentimage.WayPulled)
	if _, err := imageJobWait(t, again); err != nil {
		t.Fatalf("the new job: %v", err)
	}
}

func TestImageJobStopsWithFlaiServe(t *testing.T) {
	lab := newImageJobLab(t)
	running := lab.runner.request(imageJobNeed("/p/one", "a:1"))
	call := lab.nextCall()
	queued := lab.runner.request(imageJobNeed("/p/two", "b:1"))
	lab.stop()
	select {
	case <-call.ctx.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("the running ensure's context was not cancelled")
	}
	for _, j := range []*imageJob{running, queued} {
		_, err := imageJobWait(t, j)
		var stopped *ImageJobError
		if !errors.As(err, &stopped) || !stopped.Stopped || !strings.Contains(err.Error(), "flai serve stopped") {
			t.Fatalf("job %s: error %v, want one saying flai serve stopped", j.ID(), err)
		}
	}
	late := lab.runner.request(imageJobNeed("/p/one", "a:1"))
	if _, err := imageJobWait(t, late); err == nil {
		t.Fatal("a request after flai serve stopped passed")
	}
	lab.runner.drain()
	lab.noCall()
	if rec := lab.readImageJob(queued.ID()); rec.Outcome != imageJobStopped || !rec.Started.IsZero() {
		t.Fatalf("queued job's record %+v, want it stopped without starting", rec)
	}
}

func TestImageJobWaiterThatStopsWaitingLeavesTheJob(t *testing.T) {
	lab := newImageJobLab(t)
	j := lab.runner.request(imageJobNeed("/p/one", "a:1"))
	call := lab.nextCall()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := j.wait(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("a waiter whose context ended got %v, want context.Canceled", err)
	}
	if call.ctx.Err() != nil {
		t.Fatal("a waiter that stopped waiting cancelled the job")
	}
	call.pass(agentimage.WayPulled)
	if _, err := imageJobWait(t, j); err != nil {
		t.Fatalf("the job's other waiter: %v", err)
	}
}

func TestImageJobJournalsItsStartAndEnd(t *testing.T) {
	lab := newImageJobLab(t)
	passed := lab.runner.request(imageJobNeed("/p/one", "a:1"))
	lab.nextCall().pass(agentimage.WayPulled)
	if _, err := imageJobWait(t, passed); err != nil {
		t.Fatal(err)
	}
	failed := lab.runner.request(imageJobNeed("/p/one", "b:1"))
	lab.nextCall().answer <- imageJobAnswer{err: &agentimage.StepError{Step: "build b:1", Tries: 1, Said: "no space left on device", Err: errors.New("exit status 1")}}
	if _, err := imageJobWait(t, failed); err == nil {
		t.Fatal("the failing job passed")
	}

	entries := lab.journal()
	if len(entries) != 4 {
		t.Fatalf("journal has %d image job entries, want 4: %+v", len(entries), entries)
	}
	for i, want := range []struct {
		job              *imageJob
		outcome, verdict string
	}{{passed, "started", ""}, {passed, "done", imageJobPassed}, {failed, "started", ""}, {failed, "failed", imageJobFailed}} {
		e := entries[i]
		if e.Outcome != want.outcome || e.Verdict != want.verdict || e.Project != "one" || e.Root != "/p/one" || e.By == "" || e.At == "" {
			t.Errorf("entry %d %+v, want outcome %q, verdict %q, project one at /p/one", i, e, want.outcome, want.verdict)
		}
		if !strings.Contains(e.Detail, want.job.ID()) || !strings.Contains(e.Detail, want.job.Log()) {
			t.Errorf("entry %d detail %q does not name job %s and its log", i, e.Detail, want.job.ID())
		}
		if want.verdict != "" && e.DurationMS <= 0 {
			t.Errorf("end entry %d has no duration", i)
		}
	}

	// a start is the event with the job's log and no outcome, an end the one with its outcome
	starts, ends := map[string]map[string]any{}, map[string]map[string]any{}
	for _, ev := range lab.logs.events(t) {
		id, _ := ev["job_id"].(string)
		_, hasLog := ev["log"]
		_, hasOutcome := ev["outcome"]
		switch {
		case hasOutcome:
			ends[id] = ev
		case hasLog:
			starts[id] = ev
		}
	}
	for _, j := range []*imageJob{passed, failed} {
		start, end := starts[j.ID()], ends[j.ID()]
		if start == nil || end == nil {
			t.Fatalf("job %s logged start %v and end %v, want both", j.ID(), start, end)
		}
		for _, field := range []string{"component", "job_id", "project", "root", "images", "log"} {
			if _, ok := start[field]; !ok {
				t.Errorf("job %s's start has no %s: %v", j.ID(), field, start)
			}
		}
		fields := []string{"component", "job_id", "project", "root", "images", "outcome", "ways", "duration_ms", "log"}
		if j == failed {
			fields = append(fields, "step", "cause")
		}
		for _, field := range fields {
			if _, ok := end[field]; !ok {
				t.Errorf("job %s's end has no %s: %v", j.ID(), field, end)
			}
		}
	}
	if end := ends[passed.ID()]; end["ways"] != "a:1=pulled" || end["outcome"] != imageJobPassed {
		t.Errorf("passed job's end %v, want its way and outcome", end)
	}
	if end := ends[failed.ID()]; end["step"] != "build b:1" || end["outcome"] != imageJobFailed || end["level"] != "WARN" {
		t.Errorf("failed job's end %v, want its step and outcome at warn", end)
	}
}

// S-0757: the job that runs shows what it says it does last, its log's last
// "agent image:" line, and once it has failed, the docker step that failed.
func TestImageJobShowsItsCurrentStepWhileItRunsAndTheFailedStepAfter(t *testing.T) {
	lab := newImageJobLab(t)
	j := lab.runner.request(imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0"))
	call := lab.nextCall()
	_, _ = fmt.Fprint(call.log, "agent image: flai-agent-claude-code:1.0.0 is missing; pulling claude-code of flai 1.0.0\nlayer 1/3 ")
	_, _ = fmt.Fprint(call.log, "done\nagent image: the pull of flai-agent-claude-code:1.0.0 failed (no route); building it on this host instead\n")
	running, queued, ended := lab.runner.jobsOf("/p/one")
	if running == nil || running.ID != j.ID() || running.State != hostapi.ImageJobRunning || running.Started == "" || running.Ended != "" {
		t.Fatalf("running: %+v, want job %s running", running, j.ID())
	}
	if want := "the pull of flai-agent-claude-code:1.0.0 failed (no route); building it on this host instead"; running.Step != want {
		t.Errorf("the running job's step is %q, want %q", running.Step, want)
	}
	if running.LogSize == 0 || len(queued) != 0 || len(ended) != 0 {
		t.Errorf("log size %d, queued %+v, ended %+v: want a log, nothing queued, nothing ended", running.LogSize, queued, ended)
	}
	step := "build flai-agent-claude-code:1.0.0"
	call.answer <- imageJobAnswer{err: &agentimage.StepError{Step: step, Tries: 1, Said: "no space left on device", Err: errors.New("exit status 1")}}
	if _, err := imageJobWait(t, j); err == nil {
		t.Fatal("the job passed, want it failed")
	}
	running, _, ended = lab.runner.jobsOf("/p/one")
	if running != nil || len(ended) != 1 {
		t.Fatalf("running %+v, ended %+v: want the job ended", running, ended)
	}
	if e := ended[0]; e.ID != j.ID() || e.State != hostapi.ImageJobEnded || e.Outcome != imageJobFailed || e.Step != step || !strings.Contains(e.Cause, "no space left on device") || e.Ended == "" {
		t.Errorf("the ended job: %+v, want it failed at %s with its cause", e, step)
	}
}

// S-0757 criterion 2: a start by the operator's ask is refused, naming the
// job, while a job of the project that makes the same images waits or runs,
// and queues a new job otherwise.
func TestImageJobStartIsRefusedWhileAJobForTheSameImagesRuns(t *testing.T) {
	lab := newImageJobLab(t)
	first, ahead, err := lab.runner.start(imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0", "p-layer:abc"))
	if err != nil || ahead != 0 {
		t.Fatalf("the first start: %v, %d ahead, want a job with none ahead", err, ahead)
	}
	call := lab.nextCall()
	for _, refs := range [][]string{{"p-layer:abc", "flai-agent-claude-code:1.0.0"}, {"flai-agent-claude-code:1.0.0"}} {
		j, _, err := lab.runner.start(imageJobNeed("/p/one", refs...))
		var busy *hostapi.ImageJobBusy
		if !errors.As(err, &busy) || j != nil {
			t.Fatalf("a start of %v while job %s runs: %v, want it refused", refs, first.ID(), err)
		}
		if busy.Job.ID != first.ID() || busy.Job.State != hostapi.ImageJobRunning || !strings.Contains(err.Error(), first.ID()) {
			t.Errorf("the refusal %q, job %+v, does not name the running job %s", err, busy.Job, first.ID())
		}
	}
	other, ahead, err := lab.runner.start(imageJobNeed("/p/two", "flai-agent-claude-code:1.0.0"))
	if err != nil || ahead != 1 {
		t.Fatalf("a start of another project's images: %v, %d ahead, want it queued behind the running job", err, ahead)
	}
	if _, _, err := lab.runner.start(imageJobNeed("/p/two", "flai-agent-claude-code:1.0.0")); err == nil {
		t.Error("a start of the images of a job that waits its turn was queued again")
	}
	_, queued, _ := lab.runner.jobsOf("/p/two")
	if len(queued) != 1 || queued[0].ID != other.ID() || queued[0].State != hostapi.ImageJobQueued || queued[0].Started != "" || queued[0].Requested == "" {
		t.Errorf("queued of /p/two: %+v, want job %s waiting", queued, other.ID())
	}
	call.pass(agentimage.WayPulled)
	if _, err := imageJobWait(t, first); err != nil {
		t.Fatal(err)
	}
	lab.nextCall().pass(agentimage.WayCurrent)
	if _, err := imageJobWait(t, other); err != nil {
		t.Fatal(err)
	}
	again, _, err := lab.runner.start(imageJobNeed("/p/one", "flai-agent-claude-code:1.0.0", "p-layer:abc"))
	if err != nil || again.ID() == first.ID() {
		t.Fatalf("a start once the job ended: %v, want a new job", err)
	}
	lab.nextCall().pass(agentimage.WayCurrent)
	if _, err := imageJobWait(t, again); err != nil {
		t.Fatal(err)
	}
}

// S-0757 criterion 1: the ended jobs are the last of each set of images of
// the project, newest first, read from the records, so that those a runner
// before flai serve's restart ended are among them; another project's are not.
func TestImageJobsOfKeepTheLastEndedOfEachSetOfImages(t *testing.T) {
	lab := newImageJobLab(t)
	run := func(root string, way agentimage.Way, refs ...string) *imageJob {
		t.Helper()
		j := lab.runner.request(imageJobNeed(root, refs...))
		lab.nextCall().pass(way)
		if _, err := imageJobWait(t, j); err != nil {
			t.Fatal(err)
		}
		return j
	}
	run("/p/one", agentimage.WayPulled, "a:1")
	both := run("/p/one", agentimage.WayBuilt, "a:1", "l:1")
	harness := run("/p/one", agentimage.WayCurrent, "a:1")
	run("/p/two", agentimage.WayPulled, "a:1")

	_, _, ended := lab.runner.jobsOf("/p/one")
	if len(ended) != 2 || ended[0].ID != harness.ID() || ended[1].ID != both.ID() {
		t.Fatalf("ended: %+v, want %s then %s", ended, harness.ID(), both.ID())
	}
	if e := ended[1]; e.Outcome != imageJobPassed || len(e.Ways) != 2 || e.Ways[0].Way != string(agentimage.WayBuilt) || e.LogSize == 0 || e.DurationMS == 0 {
		t.Errorf("the ended job of both images: %+v", e)
	}

	again := newImageJobRunner(context.Background(), imageJobOptions{Dir: lab.dir, Ensure: nil})
	running, queued, after := again.jobsOf("/p/one/")
	if running != nil || len(queued) != 0 || len(after) != 2 || after[0].ID != harness.ID() || after[1].ID != both.ID() {
		t.Errorf("after a restart: running %+v, queued %+v, ended %+v; want the same two ended", running, queued, after)
	}
	if _, _, none := again.jobsOf("/p/three"); len(none) != 0 {
		t.Errorf("a project with no job: %+v", none)
	}
}

// S-0757 criterion 1: a job's log is found by its id, the project's own
// alone, while it runs and after a restart.
func TestImageJobLogIsTheProjectsOwnJob(t *testing.T) {
	lab := newImageJobLab(t)
	j := lab.runner.request(imageJobNeed("/p/one", "a:1"))
	call := lab.nextCall()
	got, path, ok := lab.runner.jobLog("/p/one", j.ID())
	if !ok || got.State != hostapi.ImageJobRunning || path != j.Log() {
		t.Fatalf("the running job's log: %+v %q %v, want job %s running with its log %s", got, path, ok, j.ID(), j.Log())
	}
	if _, _, ok := lab.runner.jobLog("/p/two", j.ID()); ok {
		t.Error("another project read the job's log")
	}
	call.pass(agentimage.WayPulled)
	if _, err := imageJobWait(t, j); err != nil {
		t.Fatal(err)
	}
	again := newImageJobRunner(context.Background(), imageJobOptions{Dir: lab.dir})
	got, path, ok = again.jobLog("/p/one", j.ID())
	if !ok || got.State != hostapi.ImageJobEnded || got.Outcome != imageJobPassed || path != j.Log() {
		t.Errorf("the ended job's log after a restart: %+v %q %v", got, path, ok)
	}
	for _, id := range []string{"20261113T093000Z-999", "../" + j.ID(), ""} {
		if _, _, ok := again.jobLog("/p/one", id); ok {
			t.Errorf("job %q found", id)
		}
	}
}

// S-0757: the host API's view of the jobs answers the images as judged, or
// unjudged with why, and starts a job of every image the agents need, unless
// none is flai's to make.
func TestImageJobsHostAnswersTheImagesAndStartsAJob(t *testing.T) {
	lab := newImageJobLab(t)
	harness := agentimage.Image{Kind: agentimage.KindHarness, Inner: "claude-code", Ref: "a:1", Source: agentimage.SourceDefault}
	layer := agentimage.Image{Kind: agentimage.KindLayer, Inner: "claude-code", Ref: "l:1", Source: agentimage.SourceLayer, Dockerfile: "agent.Dockerfile"}
	judged := projectImages{images: []agentimage.Image{harness, layer}, flai: "1.0.0",
		statuses: []agentimage.Status{{Image: harness, State: agentimage.StateCurrent}, {Image: layer, State: agentimage.StateMissing}}}
	var judgeErr error
	h := &imageJobsHost{runner: lab.runner, judge: func(string) (projectImages, error) { return judged, judgeErr }}

	st := h.Images("/p/one")
	if st.Flai != "1.0.0" || len(st.Images) != 2 || st.Images[1].State != agentimage.StateMissing || st.Unjudged != "" || st.Running != nil || len(st.Queued) != 0 || len(st.Ended) != 0 {
		t.Errorf("the images with no job: %+v", st)
	}

	started, err := h.Ensure("one", "/p/one")
	if err != nil || started.Ahead != 0 || strings.Join(started.Images, ",") != "a:1,l:1" {
		t.Fatalf("a start: %+v %v, want a job of both images", started, err)
	}
	call := lab.nextCall()
	if len(call.images) != 2 || call.root != "/p/one" {
		t.Errorf("the job makes %v for %s, want both images of /p/one", call.images, call.root)
	}
	if _, err := h.Ensure("one", "/p/one"); !errors.As(err, new(*hostapi.ImageJobBusy)) {
		t.Errorf("a second start while it runs: %v, want it refused", err)
	}
	if st := h.Images("/p/one"); st.Running == nil || st.Running.ID != started.ID {
		t.Errorf("the images while the job runs: %+v", st)
	}
	call.pass(agentimage.WayLayered)

	judged.unjudged, judged.statuses = errors.New("docker is not on PATH"), nil
	st = h.Images("/p/one")
	if st.Unjudged != "docker is not on PATH" || len(st.Images) != 2 || st.Images[0].Ref != "a:1" || st.Images[0].State != "" {
		t.Errorf("unjudged images: %+v", st)
	}

	judgeErr = errors.New("no manifest")
	if st := h.Images("/p/one"); st.Unjudged != "no manifest" || len(st.Images) != 0 {
		t.Errorf("images that cannot be listed: %+v", st)
	}
	if _, err := h.Ensure("one", "/p/one"); err == nil || err.Error() != "no manifest" {
		t.Errorf("a start of images that cannot be listed: %v", err)
	}

	judgeErr = nil
	for _, images := range [][]agentimage.Image{nil, {{Kind: agentimage.KindHarness, Inner: "claude-code", Ref: "mine:1", Source: agentimage.SourceHost}}} {
		judged.images = images
		if _, err := h.Ensure("one", "/p/one"); !errors.Is(err, hostapi.ErrNoImages) {
			t.Errorf("a start with the images %v: %v, want ErrNoImages", images, err)
		}
	}
}

// judgeProjectImages lists the images the agents of a project on the
// container harness need and judges them through docker.
func TestJudgeProjectImagesListsWhatTheAgentsNeed(t *testing.T) {
	root := t.TempDir()
	manifestOn := func(harnessName string) {
		t.Helper()
		m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nagent:\n  harness: " + harnessName + "\n"
		if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	manifestOn("claude-code")
	got, err := judgeProjectImages(root, nil, imageDocker{})
	if err != nil || len(got.images) != 0 || got.statuses != nil {
		t.Fatalf("a project with no agent on the container harness: %+v %v, want no image", got, err)
	}
	manifestOn("container")
	got, err = judgeProjectImages(root, nil, imageDocker{})
	if err != nil || len(got.images) != 1 || got.images[0].Inner != "claude-code" || got.unjudged != nil {
		t.Fatalf("a project on the container harness: %+v %v, want its harness image judged", got, err)
	}
	if len(got.statuses) != 1 || got.statuses[0].State != agentimage.StateMissing {
		t.Errorf("statuses: %+v, want the image missing", got.statuses)
	}
	if _, err := judgeProjectImages(t.TempDir(), nil, imageDocker{}); err == nil {
		t.Error("a folder that is no project was judged")
	}
}
