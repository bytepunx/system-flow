package hostapi

// The project's agent images and flai serve's image jobs (S-0757, ADR-0295
// §5): images.get, a read, answers the images the project's agents need with
// their state, as flai image status judges them, and the image jobs of the
// project, the one that runs, those queued, and the last ended of each set of
// images; images.log, a read, answers at most 64 KB of a job's log from an
// offset, as kluster.log answers a step's (ADR-0263 §4); and images.ensure, a
// write under the agent host action, asks flai serve for a new job that makes
// the project's images, refused while a job that makes them all waits or
// runs. Host.Images answers them; flai serve's image job runner is the real
// one, and tests give their own.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// The states of an image job as the images methods answer it.
const (
	ImageJobQueued  = "queued"
	ImageJobRunning = "running"
	ImageJobEnded   = "ended"
)

// imageLogChunk is the most of a job's log one images.log answers.
const imageLogChunk = 64 << 10

// imageJobID is an image job's id as flai serve issues it: the UTC stamp it
// was asked for at and a counter, such as 20261113T093000Z-001. Anything
// else would reach the folder of the jobs as a path of its own.
var imageJobID = regexp.MustCompile(`^\d{8}T\d{6}Z-\d{3,}$`)

// ErrNoImages is Host.Images.Ensure's refusal for a project whose agents need
// no image flai makes: none runs on the container harness and the manifest
// names no agent.image.layer, or every image they need is the operator's own.
var ErrNoImages = errors.New("the project's agents need no agent image flai makes: none runs on the container harness and its manifest names no agent.image.layer, or each image they run in is the operator's own, named with flai serve agent harness container --image, which flai never makes")

// ImageJobs is flai serve's image jobs as the images methods read and start
// them for a project, by its main checkout's root.
type ImageJobs interface {
	// Images is the project's agent images with their state, and its image
	// jobs. Images that cannot be listed or judged are answered unjudged,
	// with why, beside the jobs.
	Images(root string) ImagesState
	// JobLog is the project's job id and the path of its log; false when the
	// project has no job id.
	JobLog(root, id string) (ImageJob, string, bool)
	// Ensure queues a new job that makes the project's agent images, key
	// naming the project in the journal: an *ImageJobBusy while a job of the
	// project that makes them all waits or runs, and ErrNoImages when its
	// agents need none flai makes.
	Ensure(key, root string) (ImageStarted, error)
}

// ImagesState is images.get's answer.
type ImagesState struct {
	// Flai is the flai version the images must hold, "" when no image on
	// this host was judged.
	Flai string `json:"flai"`
	// Images are the images the project's agents need, each with its state,
	// as flai image status reports them; with Unjudged set, each with no
	// state.
	Images []agentimage.Status `json:"images"`
	// Unjudged is why the images' state could not be told, such as docker
	// not on flai serve's PATH; "" when it was.
	Unjudged string `json:"unjudged,omitempty"`
	// Running is the job of the project that runs, nil when none does.
	Running *ImageJob `json:"running"`
	// Queued are the jobs of the project that wait their turn, in order.
	Queued []ImageJob `json:"queued"`
	// Ended are the last ended job of each set of images of the project,
	// newest first.
	Ended []ImageJob `json:"ended"`
}

// ImageJob is one image job: its id, the refs of its images, its state, the
// way it left each image, its current step, when it was asked for, started,
// and ended, how it ended, why it failed, how long it ran, and its log's size.
type ImageJob struct {
	ID     string   `json:"id"`
	Images []string `json:"images"`
	// State is queued, running, or ended.
	State string `json:"state"`
	// Ways are the way it left each image, once it has ended.
	Ways []ImageWay `json:"ways"`
	// Step is, while it runs, what it says it does last, its log's last
	// "agent image:" line, and once it has failed, the docker step that
	// failed; "" otherwise.
	Step      string `json:"step,omitempty"`
	Requested string `json:"requested"`
	Started   string `json:"started,omitempty"`
	Ended     string `json:"ended,omitempty"`
	// Outcome is passed, failed, or stopped, once it has ended.
	Outcome string `json:"outcome,omitempty"`
	// Cause is why it failed or was stopped.
	Cause      string `json:"cause,omitempty"`
	DurationMS int64  `json:"duration_ms,omitempty"`
	// LogSize is its log's size, 0 before it starts.
	LogSize int64 `json:"log_size"`
}

// ImageWay is the way an image job left one image: pulled, built, layered,
// current, left, or own.
type ImageWay struct {
	Image string `json:"image"`
	Way   string `json:"way"`
}

// ImageStarted is images.ensure's answer: the job queued, and how many jobs
// wait or run before it, across every project flai serve serves.
type ImageStarted struct {
	ImageJob
	Ahead int `json:"ahead"`
}

// ImageJobBusy is Host.Images.Ensure's refusal while a job of the project
// that makes every image a new job would make waits or runs: that job.
type ImageJobBusy struct {
	Job ImageJob
}

// Error says which job makes the images already, and what to do instead.
func (e *ImageJobBusy) Error() string {
	at := "runs since " + e.Job.Started
	if e.Job.State == ImageJobQueued {
		at = "waits its turn since " + e.Job.Requested
	}
	return fmt.Sprintf("the image job %s, which makes %s, %s, so no other job is started for the same images: follow its output with images.log, and ask again once it has ended", e.Job.ID, strings.Join(e.Job.Images, ", "), at)
}

// ImageLog is images.log's answer: a part of a job's log, read from Offset,
// with where to read from next, as kluster.log answers a step's.
type ImageLog struct {
	ID string `json:"id"`
	// State and Outcome are the job's, Outcome "" until it has ended.
	State   string `json:"state"`
	Outcome string `json:"outcome,omitempty"`
	// Exists says whether the log is there; a queued job has none yet.
	Exists bool `json:"exists"`
	// Offset is where Text starts: the offset asked, or 0 after a reset.
	Offset int64 `json:"offset"`
	// Next is the offset to ask from next: Offset and Text's length.
	Next int64 `json:"next"`
	// Size is the log's size when it was read.
	Size int64 `json:"size"`
	// Text is at most 64 KB of the log from Offset, ending on a whole
	// character.
	Text string `json:"text"`
	// Reset is true when the log was shorter than the offset asked, and
	// Text then starts at 0.
	Reset bool `json:"reset,omitempty"`
}

// imageLogParams are images.log's params.
type imageLogParams struct {
	ID     string `json:"id"`
	Offset int64  `json:"offset"`
}

// imageLogIn reads and checks images.log's params.
func imageLogIn(raw json.RawMessage) (imageLogParams, *channel.Error) {
	in, e := decode[imageLogParams](raw)
	if e != nil {
		return in, e
	}
	switch {
	case !imageJobID.MatchString(in.ID):
		return in, bad("%q is not an image job's id; give one images.get lists, such as 20261113T093000Z-001", in.ID)
	case in.Offset < 0:
		return in, bad("%d is not an offset; give 0 to read from the start, or the next of the last read", in.Offset)
	}
	return in, nil
}

// noImageJobs refuses an images method asked of a flai that runs no image jobs.
func noImageJobs() *channel.Error {
	return &channel.Error{Code: Rule, Message: "this flai runs no image jobs: flai serve runs them, so ask through the dashboard flai serve answers, or run flai image status or flai image ensure in the project on the host"}
}

// imagesGet answers images.get: the project's agent images and its image
// jobs, whatever the host action says.
func imagesGet(_ context.Context, host Host, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
	if host.Images == nil {
		return nil, noImageJobs()
	}
	st := host.Images.Images(p.Root)
	st.Unjudged = klusterRelative(p.Root, st.Unjudged)
	if st.Images == nil {
		st.Images = []agentimage.Status{}
	}
	if st.Running != nil {
		job := imageJobShown(p.Root, *st.Running)
		st.Running = &job
	}
	st.Queued = imageJobsShown(p.Root, st.Queued)
	st.Ended = imageJobsShown(p.Root, st.Ended)
	return st, nil
}

// imageJobsShown are jobs as the dashboard is shown them, never nil.
func imageJobsShown(root string, jobs []ImageJob) []ImageJob {
	out := make([]ImageJob, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, imageJobShown(root, j))
	}
	return out
}

// imageJobShown is j with no path of the host's under root in its step or
// its cause, and its lists never nil.
func imageJobShown(root string, j ImageJob) ImageJob {
	j.Step, j.Cause = klusterRelative(root, j.Step), klusterRelative(root, j.Cause)
	if j.Images == nil {
		j.Images = []string{}
	}
	if j.Ways == nil {
		j.Ways = []ImageWay{}
	}
	return j
}

// imagesLogDo answers images.log: at most 64 KB of a job's log from an
// offset, whatever the host action says.
func imagesLogDo(_ context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
	in, e := imageLogIn(raw)
	if e != nil {
		return nil, e
	}
	if host.Images == nil {
		return nil, noImageJobs()
	}
	job, path, ok := host.Images.JobLog(p.Root, in.ID)
	if !ok {
		return nil, &channel.Error{Code: NotFound, Message: fmt.Sprintf("this project has no image job %s: give one images.get lists", in.ID)}
	}
	out, err := readImageLog(path, in.Offset)
	if err != nil {
		return nil, failed(fmt.Errorf("read the log of the image job %s: %w; check that flai serve's state folder can be read", in.ID, err))
	}
	out.ID, out.State, out.Outcome = job.ID, job.State, job.Outcome
	return out, nil
}

// readImageLog reads at most imageLogChunk bytes of the log at path from
// offset, ending on a whole character; a log not there is none, read from 0.
func readImageLog(path string, offset int64) (ImageLog, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return imageLogOf(absentLog(offset)), nil
	}
	if err != nil {
		return ImageLog{}, err
	}
	defer func() { _ = f.Close() }()
	st, err := f.Stat()
	if err != nil {
		return ImageLog{}, err
	}
	part, err := readOpenLog(f, st.Size(), offset)
	if err != nil {
		return ImageLog{}, err
	}
	return imageLogOf(part), nil
}

// imageLogOf is a part of a job's log as images.log answers it, without the
// job's id and state.
func imageLogOf(part LogRead) ImageLog {
	return ImageLog{Exists: part.Exists, Offset: part.Offset, Next: part.Next, Size: part.Size, Text: part.Text, Reset: part.Reset}
}

// wholeChars is b without a character cut short at its end, so the next
// read starts where that character does.
func wholeChars(b []byte) []byte {
	for i := 1; i <= utf8.UTFMax && i <= len(b); i++ {
		if utf8.RuneStart(b[len(b)-i]) {
			if !utf8.FullRune(b[len(b)-i:]) {
				return b[:len(b)-i]
			}
			return b
		}
	}
	return b
}

// imagesEnsureDo answers images.ensure: a new image job for the project's
// agent images, queued in flai serve, refused while one that makes them
// waits or runs, or while the project's agents need none.
func imagesEnsureDo(_ context.Context, host Host, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
	if host.Images == nil {
		return nil, noImageJobs()
	}
	started, err := host.Images.Ensure(p.Key, p.Root)
	var busy *ImageJobBusy
	switch {
	case errors.As(err, &busy):
		job := imageJobShown(p.Root, busy.Job)
		return nil, &channel.Error{Code: Rule, Message: klusterRelative(p.Root, busy.Error()), Data: map[string]any{"job": job}}
	case errors.Is(err, ErrNoImages):
		return nil, &channel.Error{Code: Rule, Message: ErrNoImages.Error() + "; so no image job is started"}
	case err != nil:
		return nil, failed(errors.New(klusterRelative(p.Root, err.Error())))
	}
	started.ImageJob = imageJobShown(p.Root, started.ImageJob)
	return written(started)
}

// describeImagesEnsure says for the journal which job images.ensure queued.
func describeImagesEnsure(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var s ImageStarted
	_ = json.Unmarshal(w.Data, &s)
	return "done", fmt.Sprintf("image job %s for %s: queued with %d ahead", s.ID, strings.Join(s.Images, ", "), s.Ahead)
}

// imageSpecs are the images methods. They run no flai command: build checks
// the params, and do asks Host.Images. images.get and images.log are reads,
// answered whatever the host action says; images.ensure is refused while the
// agent host action is off.
func imageSpecs() map[string]spec {
	none := func(channel.Project, json.RawMessage) ([]string, string, *channel.Error) { return nil, "", nil }
	return map[string]spec{
		// images.get: the project's agent images with their state, and its
		// image jobs, running, queued, and the last ended of each set.
		"images.get": {reads: true, build: none, do: imagesGet},
		// images.log: at most 64 KB of an image job's log from an offset, the
		// dashboard's live feed.
		"images.log": {reads: true, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			_, e := imageLogIn(raw)
			return nil, "", e
		}, do: imagesLogDo},
		// images.ensure: a new image job for the project's agent images.
		"images.ensure": {action: ActionAgent, describe: describeImagesEnsure, build: none, do: imagesEnsureDo},
	}
}
