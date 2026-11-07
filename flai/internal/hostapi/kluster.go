package hostapi

// The project's kluster (S-0580, ADR-0254): kluster.status, a read, reports
// the cluster the manifest's kluster: block names and the up or down flai
// serve runs for it; kluster.up and kluster.down, writes under the kluster
// host action, start kluster up or kluster down in flai serve's own process
// and answer at once, since either takes minutes. Host.Kluster runs them, one
// per project at a time, and keeps the last one's outcome and log for
// kluster.status, as Host.Components keeps its dev servers' state.
//
// Deployments (S-0581): kluster.deployments, a read, lists the namespaces'
// deployment records, main's first, with each story's title and head, each
// step's log, and how to reach the namespace for a demo; kluster.deploy,
// kluster.teardown, and kluster.promote, writes under the same host action,
// run what flai kluster deploy, teardown, and promote run, in the background
// under kluster's own lock, and answer the job at once.

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
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The jobs Host.Kluster runs: kluster up and kluster down of the cluster
// (S-0580), and a story's deploy and teardown and main's promotion into the
// default namespace (S-0581).
const (
	KlusterUp       = "up"
	KlusterDown     = "down"
	KlusterDeploy   = "deploy"
	KlusterTeardown = "teardown"
	KlusterPromote  = "promote"
)

// deploysTo says whether action is one of the jobs that deploy into a
// namespace, which kluster.deployments reports, rather than an up or a down,
// which kluster.status reports.
func deploysTo(action string) bool {
	return action == KlusterDeploy || action == KlusterTeardown || action == KlusterPromote
}

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
	// Ready says why a deploy (deploy true), a teardown, or a promotion
	// cannot start: no block, no deploy command, or no kubeconfig.
	Ready(m manifest.Manifest, deploy bool) error
	// Deploy, Teardown, and Promote run what flai kluster deploy, teardown,
	// and promote run, under kluster's lock, their progress written to out.
	Deploy(ctx context.Context, m manifest.Manifest, story, status, dir, commit string, out io.Writer) (kluster.Record, error)
	Teardown(ctx context.Context, m manifest.Manifest, story, dir string, out io.Writer) (kluster.Record, error)
	Promote(ctx context.Context, m manifest.Manifest, branch string, out io.Writer) (kluster.Record, error)
	// Services lists the services of every namespace of the cluster, as
	// kubectl get services reports them.
	Services(ctx context.Context) ([]KlusterService, error)
}

// systemCluster is kluster.Cluster as a KlusterCluster.
type systemCluster struct{ kluster.Cluster }

// Named is the cluster as its config names it.
func (c systemCluster) Named() kluster.Status {
	return kluster.Status{Name: c.Name, Profile: c.Profile, Provider: c.Provider, Context: c.Context()}
}

// Services runs kubectl get services across the cluster's namespaces with
// flai's kubeconfig and the cluster's context, and reads its JSON.
func (c systemCluster) Services(ctx context.Context) ([]KlusterService, error) {
	if _, err := c.Exec.LookPath("kubectl"); err != nil {
		return nil, errors.New("kubectl is not on the host's PATH, and flai lists a namespace's services with it; install kubectl to see them")
	}
	var out strings.Builder
	if err := c.Exec.Run(ctx, c.Root, &out, "kubectl", "--kubeconfig", kluster.KubeconfigPath(c.Root), "--context", c.Context(), "get", "services", "--all-namespaces", "--output", "json"); err != nil {
		return nil, fmt.Errorf("kubectl get services of the cluster %s failed: %w", c.Name, err)
	}
	return parseServices([]byte(out.String()))
}

// parseServices reads kubectl get services --output json: each service's
// namespace, name, and ports.
func parseServices(data []byte) ([]KlusterService, error) {
	var list struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Spec struct {
				Ports []KlusterPort `json:"ports"`
			} `json:"spec"`
		} `json:"items"`
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("read kubectl get services' JSON: %w", err)
	}
	out := make([]KlusterService, 0, len(list.Items))
	for _, it := range list.Items {
		ports := it.Spec.Ports
		if ports == nil {
			ports = []KlusterPort{}
		}
		out = append(out, KlusterService{Namespace: it.Metadata.Namespace, Name: it.Metadata.Name, Ports: ports})
	}
	return out, nil
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
	// Action is up or down, or deploy, teardown, or promote.
	Action string `json:"action"`
	// Cluster is the cluster's name.
	Cluster string `json:"cluster"`
	// Story is the story a deploy or a teardown is of; absent otherwise.
	Story string `json:"story,omitempty"`
	// Namespace is the namespace a deploy, a teardown, or a promotion acts
	// on; absent for an up or a down.
	Namespace string `json:"namespace,omitempty"`
	// StartedAt is when it started, in RFC 3339, UTC.
	StartedAt string `json:"started_at"`
	// EndedAt is when it ended; absent while it runs.
	EndedAt string `json:"ended_at,omitempty"`
	// Outcome is done or failed once it ended; absent while it runs.
	Outcome string `json:"outcome,omitempty"`
	// Error is why it failed, kluster's or flai's own sentence.
	Error string `json:"error,omitempty"`
	// Log is the file kluster's output goes to, from the project's root: for
	// a deploy, a teardown, or a promotion, its progress lines.
	Log string `json:"log"`
	// Logs is the folder each step's own log goes to, from the project's
	// root, for a deploy, a teardown, or a promotion.
	Logs string `json:"logs,omitempty"`
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
	last    map[string]KlusterJob // the last up or down
	// lastDeploy is the last deploy, teardown, or promotion.
	lastDeploy map[string]KlusterJob
	wg         sync.WaitGroup
}

// NewKlusterJobs is a KlusterJobs whose jobs end, failed, when ctx ends.
func NewKlusterJobs(ctx context.Context, o KlusterJobsOptions) *KlusterJobs {
	if o.Now == nil {
		o.Now = time.Now
	}
	return &KlusterJobs{ctx: ctx, now: o.Now, logger: o.Logger, running: map[string]KlusterJob{}, last: map[string]KlusterJob{}, lastDeploy: map[string]KlusterJob{}}
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
	return k.StartJob(root, KlusterJob{Action: action, Cluster: c.Named().Name}, func(ctx context.Context, out io.Writer) error {
		if action == KlusterDown {
			return c.Down(ctx, out)
		}
		return c.Up(ctx, out)
	}, timeout)
}

// StartJob starts job, whose Action, Cluster, and for a deploy, a teardown,
// or a promotion Story and Namespace are set, for the project at root in the
// background: run runs within timeout with its output to the job's log, and
// the job, its start and log filled in, is answered at once. While another
// job runs for the project it is refused with a *KlusterRunningError.
func (k *KlusterJobs) StartJob(root string, job KlusterJob, run func(ctx context.Context, out io.Writer) error, timeout time.Duration) (KlusterJob, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if r, ok := k.running[root]; ok {
		return KlusterJob{}, &KlusterRunningError{Job: r}
	}
	job.StartedAt, job.Log = k.now().UTC().Format(time.RFC3339), klusterLog(job.Action)
	if job.Namespace != "" {
		job.Logs = klusterRel(root, kluster.LogDir(root, job.Namespace))
	}
	path := filepath.Join(root, filepath.FromSlash(job.Log))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return KlusterJob{}, fmt.Errorf("make the folder for kluster %s's log of the cluster %s: %w", job.Action, job.Cluster, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return KlusterJob{}, fmt.Errorf("open kluster %s's log of the cluster %s: %w", job.Action, job.Cluster, err)
	}
	k.running[root] = job
	k.wg.Add(1)
	go k.run(root, job, run, f, timeout)
	return job, nil
}

// run runs job with run, its output to f, and records how it ended.
func (k *KlusterJobs) run(root string, job KlusterJob, run func(context.Context, io.Writer) error, f *os.File, timeout time.Duration) {
	defer k.wg.Done()
	ctx, cancel := context.WithTimeout(k.ctx, timeout)
	defer cancel()
	err := run(ctx, f)
	if cerr := f.Close(); err == nil && cerr != nil {
		err = fmt.Errorf("close kluster %s's log of the cluster %s: %w", job.Action, job.Cluster, cerr)
	}
	job.EndedAt, job.Outcome = k.now().UTC().Format(time.RFC3339), "done"
	if err != nil {
		job.Outcome, job.Error = "failed", klusterRelative(root, err.Error())
	}
	k.mu.Lock()
	delete(k.running, root)
	if deploysTo(job.Action) {
		k.lastDeploy[root] = job
	} else {
		k.last[root] = job
	}
	k.mu.Unlock()
	if k.logger != nil {
		k.logger.Info("kluster job ended", "component", "hostapi", "action", job.Action, "cluster", job.Cluster, "namespace", job.Namespace, "outcome", job.Outcome, "root", root)
	}
}

// State is the job that runs for the project at root, if any, whichever it
// is, and the last up or down to end there, if any.
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

// DeployState is the deploy, teardown, or promotion that runs for the
// project at root, if any, and the last of them to end there, if any.
func (k *KlusterJobs) DeployState(root string) (running, last *KlusterJob) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if r, ok := k.running[root]; ok && deploysTo(r.Action) {
		running = &r
	}
	if l, ok := k.lastDeploy[root]; ok {
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
	if j.Namespace != "" {
		what := "the namespace " + j.Namespace
		if j.Story != "" {
			what = j.Story + " in the namespace " + j.Namespace
		}
		return "done", fmt.Sprintf("started kluster %s of %s in the cluster %s, its progress in %s and its steps' logs in %s", j.Action, what, j.Cluster, j.Log, j.Logs)
	}
	return "done", fmt.Sprintf("started kluster %s of the cluster %s, its output in %s", j.Action, j.Cluster, j.Log)
}

// KlusterPort is one port of a service, as kubectl reports it.
type KlusterPort struct {
	Name     string `json:"name,omitempty"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol,omitempty"`
}

// KlusterService is one service of a namespace: its name, the host it
// answers on inside the cluster, its ports, and the command that forwards
// them to the host for a demo.
type KlusterService struct {
	// Namespace is the service's namespace; left out of kluster.deployments,
	// which lists services under their namespace.
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
	// Host is the name the service answers on inside the cluster,
	// <name>.<namespace>.svc.cluster.local.
	Host  string        `json:"host,omitempty"`
	Ports []KlusterPort `json:"ports"`
	// PortForward is the command, run in the project's main checkout, that
	// forwards each of its ports to the same port on the host.
	PortForward string `json:"port_forward,omitempty"`
}

// KlusterAccess is how to reach a namespace for a demo: through a
// port-forward, as ADR-0254 decides until a project gives each namespace a
// host of its own. Its commands are run in the project's main checkout,
// where kubeconfig is.
type KlusterAccess struct {
	// Kubeconfig is flai's kubeconfig of the cluster, from the project's root.
	Kubeconfig string `json:"kubeconfig"`
	// Context is the cluster's kube context.
	Context string `json:"context"`
	// Namespace is the namespace to reach.
	Namespace string `json:"namespace"`
	// List is the command that lists the namespace's services.
	List string `json:"list"`
	// Services are the namespace's services, each with its port-forward.
	Services []KlusterService `json:"services"`
	// ServicesError says why the services could not be read.
	ServicesError string `json:"services_error,omitempty"`
}

// KlusterStep is one step of a deployment, its log from the project's root
// with the end of it.
type KlusterStep struct {
	kluster.Step
	// LogTail is the end of the step's log, at most 4 KB, cut to whole lines.
	LogTail string `json:"log_tail,omitempty"`
}

// KlusterDeployment is one namespace's deployment as kluster.deployments
// answers it: its record, with its paths from the project's root, and what
// the dashboard needs beside it.
type KlusterDeployment struct {
	kluster.Deployment
	// Checkout, Error, and Steps stand in for the record's, with the project's
	// root left out of their paths.
	Checkout string        `json:"checkout"`
	Error    string        `json:"error,omitempty"`
	Steps    []KlusterStep `json:"steps"`
	// Title and Status are the story's, for a story's namespace.
	Title  string `json:"title,omitempty"`
	Status string `json:"status,omitempty"`
	// Head is the story branch's head now, or main's for the default
	// namespace, so that a deployment at an older commit shows; HeadError
	// says why it could not be read.
	Head      string `json:"head,omitempty"`
	HeadError string `json:"head_error,omitempty"`
	// Access is how to reach the namespace for a demo.
	Access KlusterAccess `json:"access"`
}

// KlusterReview is a story in review with no deployment recorded.
type KlusterReview struct {
	Story     string `json:"story"`
	Title     string `json:"title"`
	Head      string `json:"head,omitempty"`
	HeadError string `json:"head_error,omitempty"`
}

// KlusterDeployments is what kluster.deployments answers.
type KlusterDeployments struct {
	// Enabled, Enable, Block, and NoBlock are as kluster.status answers them.
	Enabled bool   `json:"enabled"`
	Enable  string `json:"enable"`
	Block   bool   `json:"block"`
	NoBlock string `json:"no_block,omitempty"`
	// StatusError is why the block's config could not be read.
	StatusError string `json:"status_error,omitempty"`
	// Cluster is the cluster's name, and DefaultNamespace the namespace
	// main is promoted into.
	Cluster          string `json:"cluster,omitempty"`
	DefaultNamespace string `json:"default_namespace,omitempty"`
	// Deployments are the namespaces' records, the default namespace's first,
	// then the stories' by number.
	Deployments []KlusterDeployment `json:"deployments"`
	// Review are the stories in review with no deployment recorded.
	Review []KlusterReview `json:"review"`
	// Running is the deploy, teardown, or promotion that runs now, whoever
	// started it; its log only when this flai serve did.
	Running *KlusterJob `json:"running,omitempty"`
	// Last is the last deploy, teardown, or promotion this flai serve ran.
	Last *KlusterJob `json:"last,omitempty"`
}

// klusterHead resolves rev to a commit in dir with git; tests replace it.
var klusterHead = func(ctx context.Context, dir, rev string) (string, error) {
	out, err := execx.Timed(ctx, execx.System{}).Run(dir, "git", "rev-parse", "--verify", "--quiet", rev+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("git rev-parse %s: %w", rev, err)
	}
	return strings.TrimSpace(out), nil
}

// klusterMainBranch is the project's main branch; tests replace it.
var klusterMainBranch = func(ctx context.Context, root string) (string, error) {
	return storygit.MainBranch(execx.Timed(ctx, execx.System{}), root)
}

// klusterServicesTimeout bounds the kubectl call kluster.deployments makes.
const klusterServicesTimeout = 10 * time.Second

// klusterRel is path from root, in slashes, or path with root left out
// when it is not under root.
func klusterRel(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil && filepath.IsLocal(rel) {
		return filepath.ToSlash(rel)
	}
	return klusterRelative(root, path)
}

// accessCommands are the kubectl arguments every command of a namespace's
// access starts with, from the project's main checkout.
func accessCommands(context, ns string) string {
	return fmt.Sprintf("kubectl --kubeconfig %s --context %s --namespace %s", filepath.ToSlash(filepath.Join(".flai-cache", "kluster", "kubeconfig")), context, ns)
}

// klusterAccess is how to reach ns for a demo, with the services listed
// in services, or servicesErr.
func klusterAccess(context, ns string, services []KlusterService, servicesErr string) KlusterAccess {
	base := accessCommands(context, ns)
	a := KlusterAccess{Kubeconfig: filepath.ToSlash(filepath.Join(".flai-cache", "kluster", "kubeconfig")), Context: context, Namespace: ns, List: base + " get services", Services: []KlusterService{}, ServicesError: servicesErr}
	for _, s := range services {
		if s.Namespace != ns {
			continue
		}
		s.Namespace = ""
		s.Host = s.Name + "." + ns + ".svc.cluster.local"
		if len(s.Ports) > 0 {
			pairs := make([]string, 0, len(s.Ports))
			for _, p := range s.Ports {
				pairs = append(pairs, fmt.Sprintf("%d:%d", p.Port, p.Port))
			}
			s.PortForward = base + " port-forward service/" + s.Name + " " + strings.Join(pairs, " ")
		}
		a.Services = append(a.Services, s)
	}
	return a
}

// klusterDeployments answers kluster.deployments: each namespace's record,
// the stories in review with none, and the deploy that runs.
func klusterDeployments(ctx context.Context, host Host, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
	out := KlusterDeployments{Enabled: host.enabled(ActionKluster, p.Root), Enable: EnableCommand(ActionKluster), Deployments: []KlusterDeployment{}, Review: []KlusterReview{}}
	if host.Kluster != nil {
		running, last := host.Kluster.DeployState(p.Root)
		out.Running, out.Last = withTail(p.Root, running), withTail(p.Root, last)
	}
	m, e := projectManifest(p)
	if e != nil {
		return nil, e
	}
	repo, err := workitem.Open(p.Root)
	if err != nil {
		return nil, failed(fmt.Errorf("cannot read the project's work items to list its deployments: %s", hostRelative(p.Root, err.Error())))
	}
	review, err := reviewStories(repo)
	if err != nil {
		return nil, failed(fmt.Errorf("cannot list the stories in review: %s", hostRelative(p.Root, err.Error())))
	}
	if m.Kluster == nil {
		out.NoBlock = kluster.ErrNoBlock.Error()
		out.Review = reviewRows(ctx, p.Root, review, nil)
		return out, nil
	}
	out.Block = true
	out.DefaultNamespace = m.KlusterNamespace()
	c, err := openKluster(p.Root, m)
	if err != nil {
		out.StatusError = klusterRelative(p.Root, err.Error())
		out.Review = reviewRows(ctx, p.Root, review, nil)
		return out, nil
	}
	named := c.Named()
	out.Cluster = named.Name
	if out.Running == nil {
		if j, ok := kluster.Running(p.Root); ok {
			out.Running = &KlusterJob{Action: j.Kind, Cluster: named.Name, Story: j.Story, Namespace: j.Namespace, StartedAt: j.StartedAt.UTC().Format(time.RFC3339)}
		}
	}
	list, err := kluster.Deployments(p.Root, out.DefaultNamespace)
	if err != nil {
		return nil, failed(errors.New(klusterRelative(p.Root, err.Error())))
	}
	var services []KlusterService
	servicesErr := ""
	if len(list) > 0 {
		if _, err := os.Stat(kluster.KubeconfigPath(p.Root)); err != nil {
			servicesErr = fmt.Sprintf("the cluster %s has no kubeconfig at .flai-cache/kluster/kubeconfig; bring it up on the Klusters page or with flai kluster up", named.Name)
		} else {
			sctx, cancel := context.WithTimeout(ctx, klusterServicesTimeout)
			services, err = c.Services(sctx)
			cancel()
			if err != nil {
				servicesErr = klusterRelative(p.Root, err.Error())
			}
		}
	}
	deployed := map[string]bool{}
	mainHead, mainErr := "", ""
	for _, d := range list {
		row := KlusterDeployment{Deployment: d, Checkout: klusterRel(p.Root, d.Checkout), Error: klusterRelative(p.Root, d.Error), Steps: make([]KlusterStep, 0, len(d.Steps))}
		for _, st := range d.Steps {
			step := KlusterStep{Step: st}
			if st.Log != "" {
				step.LogTail = tail(st.Log, klusterLogTail)
				step.Log = klusterRel(p.Root, st.Log)
			}
			row.Steps = append(row.Steps, step)
		}
		switch {
		case d.Story != "":
			deployed[d.Story] = true
			if it, err := repo.Get(d.Story); err == nil {
				row.Title, row.Status = it.Title, it.Status
			}
			row.Head, row.HeadError = storyHead(ctx, p.Root, d.Story)
		case d.Default:
			if mainHead == "" && mainErr == "" {
				mainHead, mainErr = defaultHead(ctx, p.Root)
			}
			row.Head, row.HeadError = mainHead, mainErr
		}
		row.Access = klusterAccess(named.Context, d.Namespace, services, servicesErr)
		out.Deployments = append(out.Deployments, row)
	}
	out.Review = reviewRows(ctx, p.Root, review, deployed)
	return out, nil
}

// reviewStories are the project's stories in review, in ID order.
func reviewStories(repo *workitem.Repo) ([]*workitem.Item, error) {
	all, err := repo.List(false)
	if err != nil {
		return nil, err
	}
	var out []*workitem.Item
	for _, it := range all {
		if it.Type == workitem.Story && it.Status == workitem.Review {
			out = append(out, it)
		}
	}
	return out, nil
}

// reviewRows are the stories in review deployed does not name, each with its
// branch's head.
func reviewRows(ctx context.Context, root string, stories []*workitem.Item, deployed map[string]bool) []KlusterReview {
	out := []KlusterReview{}
	for _, it := range stories {
		if deployed[it.ID] {
			continue
		}
		head, headErr := storyHead(ctx, root, it.ID)
		out = append(out, KlusterReview{Story: it.ID, Title: it.Title, Head: head, HeadError: headErr})
	}
	return out
}

// storyHead is the head of story's branch in the project's main checkout,
// or why it could not be read.
func storyHead(ctx context.Context, root, story string) (string, string) {
	head, err := klusterHead(ctx, root, "refs/heads/story/"+story)
	if err != nil {
		return "", klusterRelative(root, err.Error())
	}
	return head, ""
}

// defaultHead is the head of the project's main branch, or why it could not
// be read.
func defaultHead(ctx context.Context, root string) (string, string) {
	branch, err := klusterMainBranch(ctx, root)
	if err != nil {
		return "", klusterRelative(root, err.Error())
	}
	head, err := klusterHead(ctx, root, "refs/heads/"+branch)
	if err != nil {
		return "", klusterRelative(root, err.Error())
	}
	return head, ""
}

// noDeployJobs refuses a deploy, a teardown, or a promotion asked of a flai
// that runs no kluster jobs, naming the command to run in a shell instead.
func noDeployJobs(command string) *channel.Error {
	return &channel.Error{Code: Rule, Message: "this flai runs no kluster jobs: flai serve runs them, so ask through the dashboard flai serve answers, or run " + command + " in a shell on the host"}
}

// klusterStory reads the story params of kluster.deploy and kluster.teardown.
func klusterStory(raw json.RawMessage) (string, *channel.Error) {
	in, e := decode[struct {
		Story string `json:"story"`
	}](raw)
	if e != nil {
		return "", e
	}
	story := workitem.CanonicalID(in.Story)
	if !storyID.MatchString(story) {
		return "", bad("%q is not a story's ID; give one such as S-0578", in.Story)
	}
	return story, nil
}

// deployJob is a write's cluster and manifest, refused as a rule when the
// project has no block, its config cannot be read, a deploy cannot start,
// or kluster's lock is held by a job flai serve did not start here.
func deployJob(host Host, p channel.Project, command string, deploy bool) (KlusterCluster, manifest.Manifest, *channel.Error) {
	if host.Kluster == nil {
		return nil, manifest.Manifest{}, noDeployJobs(command)
	}
	m, e := projectManifest(p)
	if e != nil {
		return nil, m, e
	}
	if m.Kluster == nil {
		return nil, m, &channel.Error{Code: Rule, Message: kluster.ErrNoBlock.Error()}
	}
	c, err := openKluster(p.Root, m)
	if err != nil {
		return nil, m, &channel.Error{Code: Rule, Message: klusterRelative(p.Root, err.Error())}
	}
	if err := c.Ready(m, deploy); err != nil {
		return nil, m, &channel.Error{Code: Rule, Message: klusterRelative(p.Root, err.Error())}
	}
	if j, ok := kluster.Running(p.Root); ok {
		what := "the namespace " + j.Namespace
		if j.Story != "" {
			what = j.Story + " in the namespace " + j.Namespace
		}
		return nil, m, &channel.Error{Code: Rule, Message: fmt.Sprintf("kluster %s of %s runs since %s; wait for it to end, as kluster.deployments says, then ask again", j.Kind, what, j.StartedAt.UTC().Format(time.RFC3339))}
	}
	return c, m, nil
}

// startDeployJob starts job on Host.Kluster, its run given the job's own
// context, and answers it, refused as a rule while another job runs.
func startDeployJob(host Host, p channel.Project, m manifest.Manifest, job KlusterJob, run func(context.Context, io.Writer) error) (any, *channel.Error) {
	started, err := host.Kluster.StartJob(p.Root, job, run, m.KlusterTimeout())
	var running *KlusterRunningError
	switch {
	case errors.As(err, &running):
		return nil, &channel.Error{Code: Rule, Message: running.Error()}
	case err != nil:
		return nil, failed(errors.New(klusterRelative(p.Root, err.Error())))
	}
	return written(started)
}

// klusterDeployDo answers kluster.deploy: the story's deploy from its
// worktree at its head, started, refused as flai kluster deploy refuses it.
func klusterDeployDo(ctx context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
	story, e := klusterStory(raw)
	if e != nil {
		return nil, e
	}
	c, m, e := deployJob(host, p, "flai kluster deploy "+story, true)
	if e != nil {
		return nil, e
	}
	repo, err := workitem.Open(p.Root)
	if err != nil {
		return nil, failed(fmt.Errorf("cannot read the project's work items to deploy %s: %s", story, hostRelative(p.Root, err.Error())))
	}
	it, err := repo.Get(story)
	if err != nil {
		return nil, &channel.Error{Code: Rule, Message: hostRelative(p.Root, err.Error())}
	}
	if it.Status != workitem.InProgress && it.Status != workitem.Review {
		return nil, &channel.Error{Code: Rule, Message: fmt.Sprintf("%s is %s, and flai deploys only a story in progress or in review into its namespace; move it to in-progress or review first", it.ID, it.Status)}
	}
	wt := repo.WorktreePath(it.ID)
	if info, err := os.Stat(wt); err != nil || !info.IsDir() {
		return nil, &channel.Error{Code: Rule, Message: fmt.Sprintf("%s has no worktree at %s to build and deploy from; open it with flai stream open %s, then deploy again", it.ID, klusterRel(p.Root, wt), it.ID)}
	}
	commit, err := klusterHead(ctx, wt, "HEAD")
	if err != nil {
		return nil, failed(fmt.Errorf("read the head of %s's worktree %s: %s", it.ID, klusterRel(p.Root, wt), klusterRelative(p.Root, err.Error())))
	}
	ns, _ := kluster.StoryNamespace(it.ID)
	id, status := it.ID, it.Status
	return startDeployJob(host, p, m, KlusterJob{Action: KlusterDeploy, Cluster: c.Named().Name, Story: id, Namespace: ns}, func(ctx context.Context, out io.Writer) error {
		_, err := c.Deploy(ctx, m, id, status, wt, commit, out)
		return err
	})
}

// klusterTeardownDo answers kluster.teardown: the story's namespace torn
// down, started, refused when no deployment of it is recorded.
func klusterTeardownDo(_ context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
	story, e := klusterStory(raw)
	if e != nil {
		return nil, e
	}
	c, m, e := deployJob(host, p, "flai kluster teardown "+story, false)
	if e != nil {
		return nil, e
	}
	ns, err := kluster.StoryNamespace(story)
	if err != nil {
		return nil, bad("%s", err.Error())
	}
	if _, ok, err := kluster.ReadRecord(p.Root, ns); err != nil {
		return nil, failed(errors.New(klusterRelative(p.Root, err.Error())))
	} else if !ok {
		return nil, &channel.Error{Code: Rule, Message: fmt.Sprintf("%s has no deployment recorded; nothing to tear down", story)}
	}
	repo, err := workitem.Open(p.Root)
	if err != nil {
		return nil, failed(fmt.Errorf("cannot read the project's work items to tear down %s: %s", story, hostRelative(p.Root, err.Error())))
	}
	dir := ""
	if wt := repo.WorktreePath(story); dirIsThere(wt) {
		dir = wt
	}
	return startDeployJob(host, p, m, KlusterJob{Action: KlusterTeardown, Cluster: c.Named().Name, Story: story, Namespace: ns}, func(ctx context.Context, out io.Writer) error {
		_, err := c.Teardown(ctx, m, story, dir, out)
		return err
	})
}

// klusterPromoteDo answers kluster.promote: main's head deployed into the
// default namespace, started.
func klusterPromoteDo(ctx context.Context, host Host, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
	c, m, e := deployJob(host, p, "flai kluster promote", true)
	if e != nil {
		return nil, e
	}
	if _, err := kluster.DefaultTarget(m, "", ""); err != nil {
		return nil, &channel.Error{Code: Rule, Message: err.Error()}
	}
	branch, err := klusterMainBranch(ctx, p.Root)
	if err != nil {
		return nil, failed(fmt.Errorf("find the main branch of the project's main checkout to promote: %s", klusterRelative(p.Root, err.Error())))
	}
	return startDeployJob(host, p, m, KlusterJob{Action: KlusterPromote, Cluster: c.Named().Name, Namespace: m.KlusterNamespace()}, func(ctx context.Context, out io.Writer) error {
		_, err := c.Promote(ctx, m, branch, out)
		return err
	})
}

// dirIsThere says whether path is a folder.
func dirIsThere(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// klusterSpecs are the kluster methods. They run no flai command: build
// takes no params, and do reads the cluster or asks Host.Kluster for a job.
// kluster.status and kluster.deployments are reads, answered whatever the
// host action says; the writes are refused while it is off.
func klusterSpecs() map[string]spec {
	none := func(channel.Project, json.RawMessage) ([]string, string, *channel.Error) { return nil, "", nil }
	story := func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		_, e := klusterStory(raw)
		return nil, "", e
	}
	return map[string]spec{
		// kluster.status: the cluster, whether it runs, its age, and the up or
		// down flai serve runs or ran, with the end of its log.
		"kluster.status": {reads: true, build: none, do: klusterStatus},
		// kluster.up: kluster up of the project's cluster, started.
		"kluster.up": {action: ActionKluster, describe: describeKlusterJob, build: none, do: klusterStart(KlusterUp)},
		// kluster.down: kluster down of the project's cluster, started.
		"kluster.down": {action: ActionKluster, describe: describeKlusterJob, build: none, do: klusterStart(KlusterDown)},
		// kluster.deployments: each namespace's deployment, main's first, the
		// stories in review with none, and the deploy that runs (S-0581).
		"kluster.deployments": {reads: true, build: none, do: klusterDeployments},
		// kluster.deploy: a story's deploy into its namespace, started.
		"kluster.deploy": {action: ActionKluster, describe: describeKlusterJob, build: story, do: klusterDeployDo},
		// kluster.teardown: a story's namespace torn down, started.
		"kluster.teardown": {action: ActionKluster, describe: describeKlusterJob, build: story, do: klusterTeardownDo},
		// kluster.promote: main's head deployed into the default namespace,
		// started.
		"kluster.promote": {action: ActionKluster, describe: describeKlusterJob, build: none, do: klusterPromoteDo},
	}
}
