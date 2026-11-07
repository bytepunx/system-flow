package hostapi

// The project's kluster (S-0580, ADR-0254): kluster.status, a read, reports
// the cluster the manifest's kluster: block names and the up or down flai
// serve runs for it; kluster.up and kluster.down, writes under the kluster
// host action, start kluster up or kluster down in flai serve's own process
// and answer at once, since either takes minutes. Host.Kluster runs them, one
// per project at a time, and keeps the last one's outcome and log for
// kluster.status, as Host.Components keeps its dev servers' state.

import (
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
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// KlusterUp and KlusterDown are the two jobs Host.Kluster runs.
const (
	KlusterUp   = "up"
	KlusterDown = "down"
)

// klusterLogTail is the most of a job's log kluster.status answers, its end.
const klusterLogTail = 4096

// KlusterCluster is a project's cluster as the kluster methods run it: what
// its config names, and kluster up, down, and status. kluster.Cluster is the
// real one; tests give their own.
type KlusterCluster interface {
	// Named is the cluster as its config names it: name, profile, provider,
	// and kube context, with nothing read from kluster.
	Named() kluster.Status
	// Up brings the cluster up, its output written to out.
	Up(ctx context.Context, out io.Writer) error
	// Down takes the cluster down, its output written to out.
	Down(ctx context.Context, out io.Writer) error
	// Status reads the cluster's status from kluster status.
	Status(ctx context.Context) (kluster.Status, error)
}

// systemCluster is kluster.Cluster as a KlusterCluster.
type systemCluster struct{ kluster.Cluster }

// Named is the cluster as its config names it.
func (c systemCluster) Named() kluster.Status {
	return kluster.Status{Name: c.Name, Profile: c.Profile, Provider: c.Provider, Context: c.Context()}
}

// OpenKluster is the cluster m's kluster: block names, run from root, the
// project's main checkout, with this host's kluster; kluster.ErrNoBlock when
// m has no block.
func OpenKluster(root string, m manifest.Manifest) (KlusterCluster, error) {
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return nil, err
	}
	return systemCluster{c}, nil
}

// openKluster opens a project's cluster for the kluster methods; tests
// replace it.
var openKluster = OpenKluster

// KlusterJob is one kluster up or down flai serve ran, or runs, for a
// project, as kluster.status, kluster.up, and kluster.down answer it.
type KlusterJob struct {
	// Action is up or down.
	Action string `json:"action"`
	// Cluster is the cluster's name.
	Cluster string `json:"cluster"`
	// StartedAt is when it started, in RFC 3339, UTC.
	StartedAt string `json:"started_at"`
	// EndedAt is when it ended; absent while it runs.
	EndedAt string `json:"ended_at,omitempty"`
	// Outcome is done or failed once it ended; absent while it runs.
	Outcome string `json:"outcome,omitempty"`
	// Error is why it failed, kluster's or flai's own sentence.
	Error string `json:"error,omitempty"`
	// Log is the file kluster's output goes to, from the project's root.
	Log string `json:"log"`
	// LogTail is the end of that file, at most 4 KB, as kluster.status reads
	// it; absent from the writes' answers.
	LogTail string `json:"log_tail,omitempty"`
}

// KlusterStatus is what kluster.status answers.
type KlusterStatus struct {
	// Enabled is whether the kluster host action is on for the project.
	Enabled bool `json:"enabled"`
	// Enable is what the operator runs on the host to turn it on.
	Enable string `json:"enable"`
	// Block is whether the project's manifest has a kluster: block.
	Block bool `json:"block"`
	// NoBlock is flai's sentence naming the block to add, when Block is false.
	NoBlock string `json:"no_block,omitempty"`
	// Cluster is the cluster: what its config names and, when kluster status
	// answered, whether it exists, whether it runs, and its age. Absent when
	// there is no block or its config cannot be read.
	Cluster *kluster.Status `json:"cluster,omitempty"`
	// StatusError is why the cluster could not be read: its config, or
	// kluster status, or kluster missing from PATH.
	StatusError string `json:"status_error,omitempty"`
	// Running is the up or down that runs now; absent when none runs.
	Running *KlusterJob `json:"running,omitempty"`
	// Last is the last up or down that ended since flai serve started.
	Last *KlusterJob `json:"last,omitempty"`
}

// KlusterJobsOptions configure NewKlusterJobs.
type KlusterJobsOptions struct {
	// Now is the clock; nil is time.Now.
	Now func() time.Time
	// Logger takes a line as each job ends; nil logs nothing.
	Logger *slog.Logger
}

// KlusterJobs runs each project's kluster up and kluster down in the
// background, one per project at a time, and keeps the running one and the
// last one to end. The zero value is not usable: make one with
// NewKlusterJobs.
type KlusterJobs struct {
	ctx     context.Context
	now     func() time.Time
	logger  *slog.Logger
	mu      sync.Mutex
	running map[string]KlusterJob // by project root
	last    map[string]KlusterJob
	wg      sync.WaitGroup
}

// NewKlusterJobs is a KlusterJobs whose jobs end, failed, when ctx ends.
func NewKlusterJobs(ctx context.Context, o KlusterJobsOptions) *KlusterJobs {
	if o.Now == nil {
		o.Now = time.Now
	}
	return &KlusterJobs{ctx: ctx, now: o.Now, logger: o.Logger, running: map[string]KlusterJob{}, last: map[string]KlusterJob{}}
}

// KlusterRunningError is Start's refusal of a job while Job runs for the
// project.
type KlusterRunningError struct{ Job KlusterJob }

// Error says which job runs, since when, and what to do.
func (e *KlusterRunningError) Error() string {
	return fmt.Sprintf("kluster %s of the cluster %s runs since %s; wait for it to end, as kluster.status says, then ask again", e.Job.Action, e.Job.Cluster, e.Job.StartedAt)
}

// klusterLog is the file a job's output goes to, from the project's root.
func klusterLog(action string) string {
	return filepath.ToSlash(filepath.Join(".flai-cache", "kluster", action+".log"))
}

// Start starts action, up or down, of c for the project at root in the
// background, within timeout, kluster's output written to the job's log, and
// answers the job at once. While another job runs for the project it is
// refused with a *KlusterRunningError.
func (k *KlusterJobs) Start(root, action string, c KlusterCluster, timeout time.Duration) (KlusterJob, error) {
	name := c.Named().Name
	k.mu.Lock()
	defer k.mu.Unlock()
	if r, ok := k.running[root]; ok {
		return KlusterJob{}, &KlusterRunningError{Job: r}
	}
	job := KlusterJob{Action: action, Cluster: name, StartedAt: k.now().UTC().Format(time.RFC3339), Log: klusterLog(action)}
	path := filepath.Join(root, filepath.FromSlash(job.Log))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return KlusterJob{}, fmt.Errorf("make the folder for kluster %s's log of the cluster %s: %w", action, name, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return KlusterJob{}, fmt.Errorf("open kluster %s's log of the cluster %s: %w", action, name, err)
	}
	k.running[root] = job
	k.wg.Add(1)
	go k.run(root, job, c, f, timeout)
	return job, nil
}

// run runs job with c, its output to f, and records how it ended.
func (k *KlusterJobs) run(root string, job KlusterJob, c KlusterCluster, f *os.File, timeout time.Duration) {
	defer k.wg.Done()
	ctx, cancel := context.WithTimeout(k.ctx, timeout)
	defer cancel()
	var err error
	if job.Action == KlusterDown {
		err = c.Down(ctx, f)
	} else {
		err = c.Up(ctx, f)
	}
	if cerr := f.Close(); err == nil && cerr != nil {
		err = fmt.Errorf("close kluster %s's log of the cluster %s: %w", job.Action, job.Cluster, cerr)
	}
	job.EndedAt, job.Outcome = k.now().UTC().Format(time.RFC3339), "done"
	if err != nil {
		job.Outcome, job.Error = "failed", klusterRelative(root, err.Error())
	}
	k.mu.Lock()
	delete(k.running, root)
	k.last[root] = job
	k.mu.Unlock()
	if k.logger != nil {
		k.logger.Info("kluster job ended", "component", "hostapi", "action", job.Action, "cluster", job.Cluster, "outcome", job.Outcome, "root", root)
	}
}

// State is the job that runs for the project at root, if any, and the last
// one to end there, if any.
func (k *KlusterJobs) State(root string) (running, last *KlusterJob) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if r, ok := k.running[root]; ok {
		running = &r
	}
	if l, ok := k.last[root]; ok {
		last = &l
	}
	return running, last
}

// Wait waits until every job has ended.
func (k *KlusterJobs) Wait() { k.wg.Wait() }

// withTail is job with the end of its log, read from root; nil stays nil.
func withTail(root string, job *KlusterJob) *KlusterJob {
	if job == nil {
		return nil
	}
	job.LogTail = tail(filepath.Join(root, filepath.FromSlash(job.Log)), klusterLogTail)
	return job
}

// tail is the last n bytes of the file at path, from the first line that
// starts within them; "" when it cannot be read.
func tail(path string, n int64) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer func() { _ = f.Close() }()
	st, err := f.Stat()
	if err != nil {
		return ""
	}
	from := max(st.Size()-n, 0)
	buf := make([]byte, st.Size()-from)
	if _, err := f.ReadAt(buf, from); err != nil && !errors.Is(err, io.EOF) {
		return ""
	}
	s := string(buf)
	if from > 0 {
		if i := strings.IndexByte(s, '\n'); i >= 0 {
			s = s[i+1:]
		}
	}
	return s
}

// klusterRelative is msg with the project's root left out of the paths it
// names, as hostRelative does, and the root itself, which kluster's errors
// name as where it ran, said as the project's main checkout.
func klusterRelative(root, msg string) string {
	return strings.ReplaceAll(hostRelative(root, msg), root, "the project's main checkout")
}

// noKlusterJobs refuses an up or a down asked of a flai that runs no kluster
// jobs.
func noKlusterJobs() *channel.Error {
	return &channel.Error{Code: Rule, Message: "this flai runs no kluster jobs: flai serve runs them, so ask through the dashboard flai serve answers, or run flai kluster up or flai kluster down in a shell on the host"}
}

// projectManifest is the project's manifest, or the failure to read it.
func projectManifest(p channel.Project) (manifest.Manifest, *channel.Error) {
	m, err := manifest.Load(filepath.Join(p.Root, manifest.File))
	if err != nil {
		return manifest.Manifest{}, failed(fmt.Errorf("cannot read the project's manifest to find its kluster: block: %s", hostRelative(p.Root, err.Error())))
	}
	return m, nil
}

// klusterStatus answers kluster.status: the host action, the block, the
// cluster as kluster status reports it, and the running and last jobs.
func klusterStatus(ctx context.Context, host Host, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
	out := KlusterStatus{Enabled: host.enabled(ActionKluster, p.Root), Enable: EnableCommand(ActionKluster)}
	if host.Kluster != nil {
		running, last := host.Kluster.State(p.Root)
		out.Running, out.Last = withTail(p.Root, running), withTail(p.Root, last)
	}
	m, e := projectManifest(p)
	if e != nil {
		return nil, e
	}
	if m.Kluster == nil {
		out.NoBlock = kluster.ErrNoBlock.Error()
		return out, nil
	}
	out.Block = true
	c, err := openKluster(p.Root, m)
	if err != nil {
		out.StatusError = klusterRelative(p.Root, err.Error())
		return out, nil
	}
	named := c.Named()
	out.Cluster = &named
	s, err := c.Status(ctx)
	if err != nil {
		out.StatusError = klusterRelative(p.Root, err.Error())
		return out, nil
	}
	out.Cluster = &s
	return out, nil
}

// klusterStart is the do of kluster.up and kluster.down: the project's
// cluster, then its job started, refused when the block is missing, its
// config cannot be read, or a job runs.
func klusterStart(action string) func(context.Context, Host, channel.Project, json.RawMessage) (any, *channel.Error) {
	return func(_ context.Context, host Host, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
		if host.Kluster == nil {
			return nil, noKlusterJobs()
		}
		m, e := projectManifest(p)
		if e != nil {
			return nil, e
		}
		if m.Kluster == nil {
			return nil, &channel.Error{Code: Rule, Message: kluster.ErrNoBlock.Error()}
		}
		// a config that cannot be read is the operator's to fix, as flai
		// kluster up says
		c, err := openKluster(p.Root, m)
		if err != nil {
			return nil, &channel.Error{Code: Rule, Message: klusterRelative(p.Root, err.Error())}
		}
		job, err := host.Kluster.Start(p.Root, action, c, m.KlusterTimeout())
		var running *KlusterRunningError
		switch {
		case errors.As(err, &running):
			return nil, &channel.Error{Code: Rule, Message: running.Error()}
		case err != nil:
			return nil, failed(errors.New(klusterRelative(p.Root, err.Error())))
		}
		return written(job)
	}
}

// describeKlusterJob says for the journal which job was started, and where
// its output goes.
func describeKlusterJob(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var j KlusterJob
	_ = json.Unmarshal(w.Data, &j)
	return "done", fmt.Sprintf("started kluster %s of the cluster %s, its output in %s", j.Action, j.Cluster, j.Log)
}

// klusterSpecs are the kluster methods. They run no flai command: build
// takes no params, and do reads the cluster or asks Host.Kluster for a job.
// kluster.status is a read, answered whatever the host action says; the
// writes are refused while it is off.
func klusterSpecs() map[string]spec {
	none := func(channel.Project, json.RawMessage) ([]string, string, *channel.Error) { return nil, "", nil }
	return map[string]spec{
		// kluster.status: the cluster, whether it runs, its age, and the up or
		// down flai serve runs or ran, with the end of its log.
		"kluster.status": {reads: true, build: none, do: klusterStatus},
		// kluster.up: kluster up of the project's cluster, started.
		"kluster.up": {action: ActionKluster, describe: describeKlusterJob, build: none, do: klusterStart(KlusterUp)},
		// kluster.down: kluster down of the project's cluster, started.
		"kluster.down": {action: ActionKluster, describe: describeKlusterJob, build: none, do: klusterStart(KlusterDown)},
	}
}
