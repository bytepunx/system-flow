package kluster

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/filelock"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Cmd is one program for Exec.Command to run: in Dir, with Env laid over
// flai's own environment, Stdin as its standard input when not nil, and its
// standard output and error both to Out.
type Cmd struct {
	Dir   string
	Env   []string
	Stdin io.Reader
	Out   io.Writer
	Name  string
	Args  []string
}

// Command runs cmd, a program named with a path from cmd.Dir run from there.
func (System) Command(ctx context.Context, c Cmd) error {
	name := c.Name
	if strings.ContainsRune(name, '/') && !filepath.IsAbs(name) {
		name = filepath.Join(c.Dir, name)
	}
	cmd := exec.CommandContext(ctx, name, c.Args...)
	cmd.Dir = c.Dir
	cmd.Env = append(os.Environ(), c.Env...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Out
	cmd.Stderr = c.Out
	return cmd.Run()
}

// ExitCode is the exit code a program's failure carries, -1 when it carries
// none, as when the program could not start or was killed.
func ExitCode(err error) int {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}

// The steps of a deploy, in order, and of a teardown.
const (
	StepBuild    = "build"
	StepLoad     = "load"
	StepDeploy   = "deploy"
	StepTeardown = "teardown"
)

// The outcomes of a deployment and of each of its steps.
const (
	OutcomeRunning = "running"
	OutcomePassed  = "passed"
	OutcomeFailed  = "failed"
	OutcomeSkipped = "skipped"
)

// Dir is where flai keeps what it knows of the cluster run from root.
func Dir(root string) string { return filepath.Join(root, ".flai-cache", "kluster") }

// RecordPath is the file of the deployment record of namespace ns.
func RecordPath(root, ns string) string {
	return filepath.Join(Dir(root), "deployments", ns+".json")
}

// LogDir is the folder of the logs of namespace ns's last job.
func LogDir(root, ns string) string { return filepath.Join(Dir(root), "logs", ns) }

// LogPath is the log of one step of namespace ns's last job.
func LogPath(root, ns, step string) string { return filepath.Join(LogDir(root, ns), step+".log") }

// StatePath is namespace ns's state folder, FLAI_KLUSTER_STATE.
func StatePath(root, ns string) string { return filepath.Join(Dir(root), "state", ns) }

// Target is what a job deploys or tears down: a namespace, the story it is
// for (empty for the default namespace), and the checkout at its commit.
type Target struct {
	// Story is the story's ID, such as S-0578; empty for the default
	// namespace, which main deploys into.
	Story string
	// Namespace is the Kubernetes namespace, s-0578 or the default.
	Namespace string
	// Prefix starts the names a namespace makes in shared services: s0578_,
	// or empty for the default namespace.
	Prefix string
	// Dir is the checkout the commands run in.
	Dir string
	// Commit is the checkout's head, the whole hash.
	Commit string
}

var storyID = regexp.MustCompile(`^[Ss]-([0-9]+)$`)

// StoryNamespace is the namespace of story, its ID in lower case with the
// zero padding, such as s-0578.
func StoryNamespace(story string) (string, error) {
	if !storyID.MatchString(story) {
		return "", fmt.Errorf("%q is not a story ID; give one such as S-0578", story)
	}
	return strings.ToLower(story), nil
}

// StoryTarget is story's namespace, deployed from its worktree dir at commit.
func StoryTarget(story, dir, commit string) (Target, error) {
	ns, err := StoryNamespace(story)
	if err != nil {
		return Target{}, err
	}
	return Target{Story: strings.ToUpper(story), Namespace: ns, Prefix: strings.ReplaceAll(ns, "-", "") + "_", Dir: dir, Commit: commit}, nil
}

// DefaultTarget is the default namespace m names, deployed from dir at
// commit.
func DefaultTarget(m manifest.Manifest, dir, commit string) (Target, error) {
	ns := m.KlusterNamespace()
	if ns == "" {
		return Target{}, errors.New("the project's kluster: block names no default namespace and its name makes no DNS label; set kluster.namespace in system-flow.yaml and run flai check")
	}
	return Target{Namespace: ns, Dir: dir, Commit: commit}, nil
}

var hexCommit = regexp.MustCompile(`^[0-9a-f]{12,}$`)

// Tag is the tag of the images deployed into the target: the namespace, or
// main for the default namespace, a hyphen, and the commit's first 12 hex
// digits, such as s-0578-0123456789ab.
func (t Target) Tag() string {
	head := t.Commit
	if len(head) > 12 {
		head = head[:12]
	}
	if t.Story == "" {
		return "main-" + head
	}
	return t.Namespace + "-" + head
}

// Env is the environment every command run against the target gets, as
// design/system/kluster.md § The steps and their variables names it, but
// FLAI_KLUSTER_IMAGES, which only build gets.
func (c Cluster) Env(t Target) []string {
	return []string{
		"FLAI_STORY=" + t.Story,
		"FLAI_KLUSTER_NAME=" + c.Name,
		"FLAI_KLUSTER_NAMESPACE=" + t.Namespace,
		"FLAI_KLUSTER_PREFIX=" + t.Prefix,
		"FLAI_KLUSTER_COMMIT=" + t.Commit,
		"FLAI_KLUSTER_TAG=" + t.Tag(),
		"FLAI_KLUSTER_STATE=" + StatePath(c.Root, t.Namespace),
		"KUBECONFIG=" + KubeconfigPath(c.Root),
		"FLAI_KLUSTER_CONTEXT=" + c.Context(),
	}
}

// Step is one step of a job, as its record keeps it.
type Step struct {
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at,omitzero"`
	Outcome   string    `json:"outcome"`
	// ExitCode is the step's exit code when it failed, -1 when its program
	// gave none.
	ExitCode int `json:"exit_code,omitempty"`
	// Log is the file the step's output went to.
	Log string `json:"log,omitempty"`
}

// Record is a namespace's last deployment, as flai keeps it under
// .flai-cache/kluster/deployments and the dashboard reads it.
type Record struct {
	// Story is the story deployed, empty for the default namespace.
	Story     string `json:"story"`
	Namespace string `json:"namespace"`
	Commit    string `json:"commit"`
	Tag       string `json:"tag"`
	// Checkout is the folder the commands ran in.
	Checkout  string    `json:"checkout"`
	Images    []string  `json:"images"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at,omitzero"`
	// Outcome is running while the deployment runs, then passed or failed.
	Outcome string `json:"outcome"`
	// FailedStep, ExitCode, and Error say why a failed deployment failed.
	FailedStep string `json:"failed_step,omitempty"`
	ExitCode   int    `json:"exit_code,omitempty"`
	Error      string `json:"error,omitempty"`
	Steps      []Step `json:"steps"`
}

// ReadRecord reads namespace ns's deployment record; ok is false when there
// is none.
func ReadRecord(root, ns string) (r Record, ok bool, err error) {
	data, err := os.ReadFile(RecordPath(root, ns))
	if errors.Is(err, os.ErrNotExist) {
		return Record{}, false, nil
	}
	if err != nil {
		return Record{}, false, fmt.Errorf("read the deployment record of %s: %w", ns, err)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return Record{}, false, fmt.Errorf("read the deployment record %s: %w; delete it, and deploy %s again", RecordPath(root, ns), err, ns)
	}
	return r, true, nil
}

// writeRecord writes r as its namespace's record, whole or not at all.
func writeRecord(root string, r Record) error {
	path := RecordPath(root, r.Namespace)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("make the folder of the deployment records: %w", err)
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("encode the deployment record of %s: %w", r.Namespace, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write the deployment record of %s: %w", r.Namespace, err)
	}
	return nil
}

// StepError is a step that failed: which, its exit code, and its log.
type StepError struct {
	Namespace string
	Step      string
	ExitCode  int
	Log       string
	Err       error
}

func (e *StepError) Error() string {
	code := "no exit code"
	if e.ExitCode >= 0 {
		code = "exit code " + strconv.Itoa(e.ExitCode)
	}
	return fmt.Sprintf("the %s step for the namespace %s failed with %s: %v; read its log %s, fix the cause, and run it again", e.Step, e.Namespace, code, e.Err, e.Log)
}

func (e *StepError) Unwrap() error { return e.Err }

// Job is the job that holds the project's kluster lock: what it does, and to
// which namespace.
type Job struct {
	// Kind is deploy, teardown, or promote.
	Kind      string    `json:"kind"`
	Story     string    `json:"story"`
	Namespace string    `json:"namespace"`
	StartedAt time.Time `json:"started_at"`
	PID       int       `json:"pid"`
}

// ErrBusy says another job holds the project's kluster lock.
var ErrBusy = errors.New("another kluster job runs for this project")

// lockWait rides over a list's brief look at the lock.
const lockWait = 200 * time.Millisecond

func lockPath(root string) string { return filepath.Join(Dir(root), "job.lock") }
func jobPath(root string) string  { return filepath.Join(Dir(root), "job.json") }

// acquire takes the project's kluster lock for j, so that one job runs at a
// time, and names j in the job file while it holds it. A lock another job
// holds is ErrBusy, naming that job.
func acquire(root string, j Job) (release func(), err error) {
	unlock, err := filelock.Lock(lockPath(root), lockWait)
	if err != nil {
		if held, ok := readJob(root); ok {
			return nil, fmt.Errorf("%w: a %s of the namespace %s, started %s by process %d; wait for it to end, then run this again", ErrBusy, held.Kind, held.Namespace, held.StartedAt.UTC().Format(time.RFC3339), held.PID)
		}
		return nil, fmt.Errorf("%w (%w); wait for it to end, then run this again", ErrBusy, err)
	}
	data, err := json.Marshal(j)
	if err == nil {
		err = atomicfile.WriteFile(jobPath(root), data, 0o644)
	}
	if err != nil {
		unlock()
		return nil, fmt.Errorf("name the running kluster job in %s: %w", jobPath(root), err)
	}
	return func() {
		_ = os.Remove(jobPath(root))
		unlock()
	}, nil
}

func readJob(root string) (Job, bool) {
	data, err := os.ReadFile(jobPath(root))
	if err != nil {
		return Job{}, false
	}
	var j Job
	if json.Unmarshal(data, &j) != nil {
		return Job{}, false
	}
	return j, true
}

// Running is the kluster job running for the project run from root, if any.
// A job file left by a job that ended without removing it is not running.
func Running(root string) (Job, bool) {
	j, ok := readJob(root)
	if !ok {
		return Job{}, false
	}
	unlock, err := filelock.Lock(lockPath(root), time.Nanosecond)
	if err != nil {
		return j, true
	}
	unlock()
	return Job{}, false
}

// Deploy deploys story, whose status is status, from its worktree dir at
// commit, into its namespace: build, load, and deploy, as
// design/system/kluster.md § The steps and their variables describes, each
// step logged, the outcome recorded, and progress lines written to out. It
// refuses a story that is not in progress or in review, and runs while no
// other kluster job runs for the project.
func (c Cluster) Deploy(ctx context.Context, m manifest.Manifest, story, status, dir, commit string, out io.Writer) (Record, error) {
	if status != "in-progress" && status != "review" {
		return Record{}, fmt.Errorf("%s is %s, and flai deploys only a story in progress or in review into its namespace; move it to in-progress or review first", story, status)
	}
	t, err := StoryTarget(story, dir, commit)
	if err != nil {
		return Record{}, err
	}
	return c.deploy(ctx, m, "deploy", t, out)
}

// deploy runs a deploy of t under the project's lock, recorded as kind.
func (c Cluster) deploy(ctx context.Context, m manifest.Manifest, kind string, t Target, out io.Writer) (Record, error) {
	if m.Kluster == nil {
		return Record{}, ErrNoBlock
	}
	if len(m.Kluster.Deploy) == 0 || strings.TrimSpace(m.Kluster.Deploy[0]) == "" {
		return Record{}, errors.New("the kluster: block of system-flow.yaml names no deploy command; write kluster.deploy, the argument list that deploys into a namespace, and run flai check")
	}
	if !hexCommit.MatchString(t.Commit) {
		return Record{}, fmt.Errorf("%q is not a commit's whole hash; give the head of the checkout %s", t.Commit, t.Dir)
	}
	if _, err := os.Stat(KubeconfigPath(c.Root)); err != nil {
		return Record{}, fmt.Errorf("the cluster %s has no kubeconfig at %s; bring it up with flai kluster up, then run this again", c.Name, KubeconfigPath(c.Root))
	}
	release, err := acquire(c.Root, Job{Kind: kind, Story: t.Story, Namespace: t.Namespace, StartedAt: time.Now().UTC(), PID: os.Getpid()})
	if err != nil {
		return Record{}, err
	}
	defer release()

	if timeout := m.KlusterTimeout(); timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	if err := os.RemoveAll(LogDir(c.Root, t.Namespace)); err != nil {
		return Record{}, fmt.Errorf("clear the logs of %s's last job: %w", t.Namespace, err)
	}
	for _, d := range []string{LogDir(c.Root, t.Namespace), StatePath(c.Root, t.Namespace)} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return Record{}, fmt.Errorf("make the folder %s for the namespace %s: %w", d, t.Namespace, err)
		}
	}
	r := Record{Story: t.Story, Namespace: t.Namespace, Commit: t.Commit, Tag: t.Tag(), Checkout: t.Dir, Images: []string{}, StartedAt: time.Now().UTC(), Outcome: OutcomeRunning, Steps: []Step{}}
	if err := writeRecord(c.Root, r); err != nil {
		return r, err
	}
	err = c.steps(ctx, m.Kluster, t, &r, out)
	r.EndedAt = time.Now().UTC()
	r.Outcome = OutcomePassed
	if err != nil {
		r.Outcome, r.Error = OutcomeFailed, err.Error()
		var se *StepError
		if errors.As(err, &se) {
			r.FailedStep, r.ExitCode = se.Step, se.ExitCode
		}
	}
	if werr := writeRecord(c.Root, r); werr != nil && err == nil {
		err = werr
	}
	return r, err
}

// steps runs build, load, and deploy for t, recording each in r.
func (c Cluster) steps(ctx context.Context, k *manifest.Kluster, t Target, r *Record, out io.Writer) error {
	env := c.Env(t)
	images := filepath.Join(LogDir(c.Root, t.Namespace), "images")
	if len(k.Build) == 0 {
		r.Steps = append(r.Steps, Step{Name: StepBuild, StartedAt: time.Now().UTC(), EndedAt: time.Now().UTC(), Outcome: OutcomeSkipped},
			Step{Name: StepLoad, StartedAt: time.Now().UTC(), EndedAt: time.Now().UTC(), Outcome: OutcomeSkipped})
	} else {
		if err := os.WriteFile(images, nil, 0o644); err != nil {
			return fmt.Errorf("make the images file %s for the build of %s: %w", images, t.Namespace, err)
		}
		buildEnv := append(append([]string{}, env...), "FLAI_KLUSTER_IMAGES="+images)
		if err := c.step(ctx, r, t, StepBuild, out, func(log io.Writer) error {
			return c.Exec.Command(ctx, Cmd{Dir: t.Dir, Env: buildEnv, Out: log, Name: k.Build[0], Args: k.Build[1:]})
		}); err != nil {
			return err
		}
		list, err := readImages(images)
		if err != nil {
			return err
		}
		r.Images = list
		if err := c.step(ctx, r, t, StepLoad, out, func(log io.Writer) error {
			return c.load(ctx, list, log)
		}); err != nil {
			return err
		}
	}
	return c.step(ctx, r, t, StepDeploy, out, func(log io.Writer) error {
		return c.Exec.Command(ctx, Cmd{Dir: t.Dir, Env: env, Out: log, Name: k.Deploy[0], Args: k.Deploy[1:]})
	})
}

// step runs one step with its output to its log, records it in r, and
// writes r so that a job cut short leaves what it did.
func (c Cluster) step(ctx context.Context, r *Record, t Target, name string, out io.Writer, run func(log io.Writer) error) error {
	path := LogPath(c.Root, t.Namespace, name)
	s := Step{Name: name, StartedAt: time.Now().UTC(), Outcome: OutcomeRunning, Log: path}
	fmt.Fprintf(out, "%s: %s, log %s\n", t.Namespace, name, path)
	f, err := os.Create(path)
	if err != nil {
		return &StepError{Namespace: t.Namespace, Step: name, ExitCode: -1, Log: path, Err: fmt.Errorf("make its log: %w", err)}
	}
	err = run(f)
	if cerr := f.Close(); cerr != nil && err == nil {
		err = cerr
	}
	s.EndedAt, s.Outcome = time.Now().UTC(), OutcomePassed
	if err != nil {
		if ctx.Err() != nil {
			err = fmt.Errorf("%w (the job's time, kluster.timeout, ran out or it was stopped)", err)
		}
		s.Outcome, s.ExitCode = OutcomeFailed, ExitCode(err)
		err = &StepError{Namespace: t.Namespace, Step: name, ExitCode: s.ExitCode, Log: path, Err: err}
	}
	r.Steps = append(r.Steps, s)
	if werr := writeRecord(c.Root, *r); werr != nil && err == nil {
		err = werr
	}
	return err
}

// readImages reads the image references build listed, one a line, blank
// lines passed over.
func readImages(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read the images the build listed in %s: %w", path, err)
	}
	list := []string{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" {
			list = append(list, l)
		}
	}
	return list, nil
}

// Nodes are the cluster's server and agent containers, which k3d labels
// k3d.cluster=<name> and k3d.role.
func (c Cluster) Nodes(ctx context.Context, log io.Writer) ([]string, error) {
	var buf bytes.Buffer
	err := c.Exec.Command(ctx, Cmd{Dir: c.Root, Out: &buf, Name: "docker", Args: []string{"ps", "--filter", "label=k3d.cluster=" + c.Name, "--format", `{{.Names}} {{.Label "k3d.role"}}`}})
	_, _ = log.Write(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("list the nodes of the cluster %s with docker ps: %w", c.Name, err)
	}
	var nodes []string
	for _, l := range strings.Split(buf.String(), "\n") {
		f := strings.Fields(l)
		if len(f) == 2 && (f[1] == "server" || f[1] == "agent") {
			nodes = append(nodes, f[0])
		}
	}
	return nodes, nil
}

// load loads each image into every node of the cluster: docker save into a
// file, then that file into docker exec -i <node> ctr -n k8s.io images
// import - on each node.
func (c Cluster) load(ctx context.Context, images []string, log io.Writer) error {
	if len(images) == 0 {
		fmt.Fprintln(log, "the build listed no image; nothing to load")
		return nil
	}
	if _, err := c.Exec.LookPath("docker"); err != nil {
		return errors.New("docker is not on PATH, and flai loads the built images into the cluster with it; install Docker, then deploy again")
	}
	nodes, err := c.Nodes(ctx, log)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("no node of the cluster %s runs; bring it up with flai kluster up, then deploy again", c.Name)
	}
	tmp, err := os.MkdirTemp("", "flai-kluster-load-")
	if err != nil {
		return fmt.Errorf("make a folder to save the images in: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	for i, image := range images {
		tar := filepath.Join(tmp, strconv.Itoa(i)+".tar")
		fmt.Fprintf(log, "saving %s\n", image)
		if err := c.Exec.Command(ctx, Cmd{Dir: c.Root, Out: log, Name: "docker", Args: []string{"save", "-o", tar, image}}); err != nil {
			return fmt.Errorf("save the image %s with docker save: %w", image, err)
		}
		for _, node := range nodes {
			fmt.Fprintf(log, "loading %s into %s\n", image, node)
			if err := c.importInto(ctx, tar, node, log); err != nil {
				return fmt.Errorf("load the image %s into the node %s: %w", image, node, err)
			}
		}
	}
	return nil
}

func (c Cluster) importInto(ctx context.Context, tar, node string, log io.Writer) error {
	f, err := os.Open(tar)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return c.Exec.Command(ctx, Cmd{Dir: c.Root, Stdin: f, Out: log, Name: "docker", Args: []string{"exec", "-i", node, "ctr", "-n", "k8s.io", "images", "import", "-"}})
}
