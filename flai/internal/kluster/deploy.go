package kluster

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/execx"
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

// The outcomes of a deployment and of each of its steps. A step is pending
// from the job's start, when the record lists every step the job will run,
// until it starts (ADR-0263).
const (
	OutcomePending = "pending"
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

// ComponentStep is the build step of the kluster: block's component name,
// build-<name> (ADR-0263).
func ComponentStep(name string) string { return StepBuild + "-" + name }

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
	head := shortCommit(t.Commit)
	if t.Story == "" {
		return "main-" + head
	}
	return t.Namespace + "-" + head
}

// Env is the environment every command run against the target gets, as
// design/system/kluster.md § The steps and their variables names it, but
// FLAI_KLUSTER_IMAGES and FLAI_KLUSTER_COMPONENT, which only a build gets.
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

// Step is one step of a job, as its record keeps it: pending, with no start,
// until it runs.
type Step struct {
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at,omitzero"`
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
	// ChecksCommit, ChecksOutcome, and ChecksLog are the commit, the outcome,
	// and the log of the last checks run against the namespace since its
	// deploy: passed, failed, timed-out, or cancelled. All are empty until a
	// run ends, and a deploy clears them (ADR-0254 §5a, §8, ADR-0292 §6).
	ChecksCommit  string `json:"checks_commit,omitempty"`
	ChecksOutcome string `json:"checks_outcome,omitempty"`
	ChecksLog     string `json:"checks_log,omitempty"`
	// ChecksAskedBy is who asked for that checks run: the story's agent, or
	// flai serve for the review pass; empty for a run from a shell or the
	// dashboard, and in a record written before it was kept (S-0663).
	ChecksAskedBy string `json:"checks_asked_by,omitempty"`
	// AskedBy is who asked for the record's last job: the story's agent that
	// asked flai serve, or flai serve for the review pass; empty for a job run
	// in a shell, and in a record written before it was kept (ADR-0292 §6).
	AskedBy string `json:"asked_by,omitempty"`
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

// recordsLockPath is the lock every write of a namespace's record takes, so
// that a checks outcome kept in a record never overwrites a deploy's write
// to it, nor a deploy's the outcome's half done.
func recordsLockPath(root string) string { return filepath.Join(Dir(root), "deployments.lock") }

// recordsLockWait is the most a write of a record waits for another's.
const recordsLockWait = 10 * time.Second

// lockRecords takes the lock of the records' writes.
func lockRecords(root string) (unlock func(), err error) {
	unlock, err = filelock.Lock(recordsLockPath(root), recordsLockWait)
	if err != nil {
		return nil, fmt.Errorf("take the lock of the deployment records: %w", err)
	}
	return unlock, nil
}

// writeRecord writes r as its namespace's record, whole or not at all, under
// the records' lock.
func writeRecord(root string, r Record) error {
	unlock, err := lockRecords(root)
	if err != nil {
		return err
	}
	defer unlock()
	return putRecord(root, r)
}

// putRecord writes r as its namespace's record, whole or not at all; the
// caller holds the records' lock.
func putRecord(root string, r Record) error {
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

// RecordChecks keeps in namespace ns's record the outcome of a checks run
// against it at commit, which ended outcome, its log, and by, who asked for
// it. It refuses a namespace with no record, and one deployed again at
// another commit since the run started, whose record the outcome does not
// describe.
func RecordChecks(root, ns, commit, outcome, log, by string) error {
	unlock, err := lockRecords(root)
	if err != nil {
		return err
	}
	defer unlock()
	r, ok, err := ReadRecord(root, ns)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the namespace %s has no deployment record to keep its checks run at %s, which ended %s, in; deploy it, then run its checks again", ns, shortCommit(commit), outcome)
	}
	if r.Commit != commit {
		return fmt.Errorf("the namespace %s was deployed again, at %s, while its checks ran at %s, so their outcome %s is not kept; run its checks again at %s", ns, shortCommit(r.Commit), shortCommit(commit), outcome, shortCommit(r.Commit))
	}
	r.ChecksCommit, r.ChecksOutcome, r.ChecksLog, r.ChecksAskedBy = commit, outcome, log, by
	return putRecord(root, r)
}

// StoryHead is the commit at the head of story's branch, refs/heads/story/<ID>,
// in the checkout root, as git resolves it there: the head a story's deploy
// deploys and its acceptance is gated at. A branch git cannot resolve is an
// error, with an empty head.
func StoryHead(git execx.Runner, root, story string) (string, error) {
	ref := "refs/heads/story/" + strings.ToUpper(story)
	out, err := git.Run(root, "git", "rev-parse", "--verify", "--quiet", ref+"^{commit}")
	head := strings.TrimSpace(out)
	if err != nil || head == "" {
		return "", fmt.Errorf("resolve the head of %s's branch %s in %s with git rev-parse: %w; open the story's branch with flai stream open %s", story, ref, root, cmp.Or(err, errors.New("no commit")), story)
	}
	return head, nil
}

// AcceptGate says whether story may be accepted while the kluster host action
// is on for the project run from root (ADR-0254 §5a): ok when its namespace's
// record says its deploy passed at head, the story branch's head, and, while
// checksOn says the checks host action is on, a checks run against it passed
// there too. A record's commit is at head as AtHead judges it with git in
// root and m's layout (ADR-0292 §3): head itself, or a commit followed only
// by changes under the design, docs, and wip folders. Otherwise reason names
// the last job, its commit, and its outcome, and for a commit not at head
// why. Only a record that cannot be read, or a diff git cannot take, is an
// error.
func AcceptGate(git execx.Runner, root string, m manifest.Manifest, story, head string, checksOn bool) (ok bool, reason string, err error) {
	ns, err := StoryNamespace(story)
	if err != nil {
		return false, "", err
	}
	story = strings.ToUpper(story)
	if head == "" {
		return false, fmt.Sprintf("%s has no branch story/%s, so its namespace %s cannot be deployed at its head", story, story, ns), nil
	}
	r, found, err := ReadRecord(root, ns)
	if err != nil {
		return false, "", err
	}
	if !found {
		return false, fmt.Sprintf("%s's namespace %s was never deployed: no deploy of it is recorded", story, ns), nil
	}
	deployed, err := AtHead(git, root, m, r.Commit, head)
	if err != nil {
		return false, "", err
	}
	checked := deployed
	if checksOn && r.ChecksCommit != r.Commit {
		if checked, err = AtHead(git, root, m, r.ChecksCommit, head); err != nil {
			return false, "", err
		}
	}
	job, outcome := lastJob(root, r)
	switch {
	case !deployed.AtHead:
		return false, fmt.Sprintf("the last %s of %s's namespace %s was at %s, not at the story's head %s, and it %s: %s", job, story, ns, shortCommit(r.Commit), shortCommit(head), outcome, deployed.Why), nil
	case r.Outcome != OutcomePassed:
		return false, fmt.Sprintf("the %s of %s's namespace %s at %s %s", job, story, ns, shortCommit(head), outcome), nil
	case !checksOn:
		return true, "", nil
	case r.ChecksCommit == "":
		return false, fmt.Sprintf("the deploy of %s's namespace %s at %s passed, and no checks run against it has ended there yet", story, ns, shortCommit(head)), nil
	case !checked.AtHead:
		return false, fmt.Sprintf("the deploy of %s's namespace %s at %s passed, and no checks run against it has ended there yet: its last, at %s, ended %s, and %s", story, ns, shortCommit(head), shortCommit(r.ChecksCommit), r.ChecksOutcome, checked.Why), nil
	case r.ChecksOutcome != OutcomePassed:
		return false, fmt.Sprintf("the deploy of %s's namespace %s at %s passed, and the checks run against it there ended %s", story, ns, shortCommit(head), r.ChecksOutcome), nil
	}
	return true, "", nil
}

// lastJob names the last job r records, a deploy or a teardown, and its
// outcome: passed, failed at a step, running, or interrupted when no job runs
// for its namespace any more.
func lastJob(root string, r Record) (job, outcome string) {
	job = "deploy"
	if r.FailedStep == StepTeardown {
		job = StepTeardown
	}
	switch r.Outcome {
	case OutcomeFailed:
		if r.FailedStep == "" {
			return job, "failed before its first step"
		}
		return job, "failed at the " + r.FailedStep + " step"
	case OutcomeRunning:
		if j, running := Running(root); running && j.Namespace == r.Namespace {
			return job, "is still running"
		}
		return job, "was interrupted before it recorded its outcome"
	}
	return job, r.Outcome
}

// shortCommit is a commit's first 12 hex digits, as a tag names it.
func shortCommit(commit string) string {
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
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
	// By is who asked for the job, as the record's AskedBy; empty for a job
	// run in a shell, and in a job file written before it was kept.
	By string `json:"by,omitempty"`
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
// other kluster job runs for the project. Neither its record nor its job says
// who asked; DeployBy's do.
func (c Cluster) Deploy(ctx context.Context, m manifest.Manifest, story, status, dir, commit string, out io.Writer) (Record, error) {
	return c.DeployBy(ctx, m, "", story, status, dir, commit, out)
}

// DeployBy deploys as Deploy does, for by, who asked for it: the story's
// agent, or flai serve for the review pass. The record's AskedBy and the
// running job's By name by; empty names no one, as Deploy does.
func (c Cluster) DeployBy(ctx context.Context, m manifest.Manifest, by, story, status, dir, commit string, out io.Writer) (Record, error) {
	if status != "in-progress" && status != "review" {
		return Record{}, fmt.Errorf("%s is %s, and flai deploys only a story in progress or in review into its namespace; move it to in-progress or review first", story, status)
	}
	t, err := StoryTarget(story, dir, commit)
	if err != nil {
		return Record{}, err
	}
	return c.deploy(ctx, m, "deploy", by, t, out)
}

// ready says why a job cannot run for the project: no kluster: block, no
// deploy command when deploy is set, or no kubeconfig for the cluster.
func (c Cluster) ready(m manifest.Manifest, deploy bool) error {
	if m.Kluster == nil {
		return ErrNoBlock
	}
	if deploy && (len(m.Kluster.Deploy) == 0 || strings.TrimSpace(m.Kluster.Deploy[0]) == "") {
		return errors.New("the kluster: block of system-flow.yaml names no deploy command; write kluster.deploy, the argument list that deploys into a namespace, and run flai check")
	}
	if _, err := os.Stat(KubeconfigPath(c.Root)); err != nil {
		return fmt.Errorf("the cluster %s has no kubeconfig at %s; bring it up with flai kluster up, then run this again", c.Name, KubeconfigPath(c.Root))
	}
	return nil
}

// Ready says why a job cannot run for the project, as Deploy, Teardown, and
// Promote would refuse it before they start: no kluster: block, no deploy
// command when deploy is set, or no kubeconfig for the cluster. The
// dashboard's host methods ask it before they start a job in the background.
func (c Cluster) Ready(m manifest.Manifest, deploy bool) error { return c.ready(m, deploy) }

// withTimeout bounds ctx by kluster.timeout, the most one job may take.
func withTimeout(ctx context.Context, m manifest.Manifest) (context.Context, context.CancelFunc) {
	if timeout := m.KlusterTimeout(); timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return context.WithCancel(ctx)
}

// deploy runs a deploy of t under the project's lock, recorded as kind and
// asked for by by.
func (c Cluster) deploy(ctx context.Context, m manifest.Manifest, kind, by string, t Target, out io.Writer) (Record, error) {
	if err := c.ready(m, true); err != nil {
		return Record{}, err
	}
	if !hexCommit.MatchString(t.Commit) {
		return Record{}, fmt.Errorf("%q is not a commit's whole hash; give the head of the checkout %s", t.Commit, t.Dir)
	}
	release, err := acquire(c.Root, Job{Kind: kind, Story: t.Story, Namespace: t.Namespace, StartedAt: time.Now().UTC(), PID: os.Getpid(), By: by})
	if err != nil {
		return Record{}, err
	}
	defer release()
	ctx, cancel := withTimeout(ctx, m)
	defer cancel()
	return c.runDeploy(ctx, m, by, t, out)
}

// runDeploy deploys t for by, its logs cleared first and its record written
// as it goes; the caller holds the project's lock and bounds ctx.
func (c Cluster) runDeploy(ctx context.Context, m manifest.Manifest, by string, t Target, out io.Writer) (Record, error) {
	if err := os.RemoveAll(LogDir(c.Root, t.Namespace)); err != nil {
		return Record{}, fmt.Errorf("clear the logs of %s's last job: %w", t.Namespace, err)
	}
	for _, d := range []string{LogDir(c.Root, t.Namespace), StatePath(c.Root, t.Namespace)} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return Record{}, fmt.Errorf("make the folder %s for the namespace %s: %w", d, t.Namespace, err)
		}
	}
	r := Record{Story: t.Story, Namespace: t.Namespace, Commit: t.Commit, Tag: t.Tag(), Checkout: t.Dir, Images: []string{}, StartedAt: time.Now().UTC(), Outcome: OutcomeRunning, Steps: plan(m.Kluster, time.Now().UTC()), AskedBy: by}
	if err := writeRecord(c.Root, r); err != nil {
		return r, err
	}
	err := c.steps(ctx, m.Kluster, t, &r, out)
	r.EndedAt = time.Now().UTC()
	r.Outcome = OutcomePassed
	if err != nil {
		r.Outcome, r.Error = OutcomeFailed, err.Error()
		var se *StepError
		if errors.As(err, &se) {
			r.FailedStep, r.ExitCode = se.Step, se.ExitCode
		}
	}
	skipPending(&r)
	if werr := writeRecord(c.Root, r); werr != nil && err == nil {
		err = werr
	}
	return r, err
}

// plan is every step a deploy by k will run, in order, each pending but a
// build and a load k gives nothing to do, which are skipped at now: one build
// step per component, or build, then load, then deploy (ADR-0263).
func plan(k *manifest.Kluster, now time.Time) []Step {
	var steps []Step
	switch {
	case len(k.Components) > 0:
		for _, comp := range k.Components {
			steps = append(steps, Step{Name: ComponentStep(comp.Name), Outcome: OutcomePending})
		}
		steps = append(steps, Step{Name: StepLoad, Outcome: OutcomePending})
	case len(k.Build) > 0:
		steps = append(steps, Step{Name: StepBuild, Outcome: OutcomePending}, Step{Name: StepLoad, Outcome: OutcomePending})
	default:
		steps = append(steps, Step{Name: StepBuild, StartedAt: now, EndedAt: now, Outcome: OutcomeSkipped},
			Step{Name: StepLoad, StartedAt: now, EndedAt: now, Outcome: OutcomeSkipped})
	}
	return append(steps, Step{Name: StepDeploy, Outcome: OutcomePending})
}

// skipPending records as skipped each step of r still pending when its job
// ended: the steps after the one that failed.
func skipPending(r *Record) {
	for i := range r.Steps {
		if r.Steps[i].Outcome == OutcomePending {
			r.Steps[i].Outcome = OutcomeSkipped
		}
	}
}

// build is one build a deploy runs: its step's name, its command, its images
// file, and the component it builds, empty for the block's build.
type build struct {
	step, component, images string
	cmd                     []string
}

// builds are the builds k gives a deploy into t: one per component, each
// with its own images file, or the block's build alone.
func (c Cluster) builds(k *manifest.Kluster, t Target) []build {
	dir := LogDir(c.Root, t.Namespace)
	if len(k.Components) > 0 {
		var out []build
		for _, comp := range k.Components {
			out = append(out, build{step: ComponentStep(comp.Name), component: comp.Name, images: filepath.Join(dir, "images-"+comp.Name), cmd: comp.Build})
		}
		return out
	}
	if len(k.Build) > 0 {
		return []build{{step: StepBuild, images: filepath.Join(dir, "images"), cmd: k.Build}}
	}
	return nil
}

// steps runs the builds, the load, and deploy for t, recording each in r.
// For a story's namespace the deploy step, once its command passes, also
// writes the namespace's own kubeconfig, so that a failure to write it fails
// the deploy step, in its log (ADR-0292 §2).
func (c Cluster) steps(ctx context.Context, k *manifest.Kluster, t Target, r *Record, out io.Writer) error {
	env := c.Env(t)
	if bs := c.builds(k, t); len(bs) > 0 {
		list := []string{}
		for _, b := range bs {
			if err := os.WriteFile(b.images, nil, 0o644); err != nil {
				return fmt.Errorf("make the images file %s for the %s step of %s: %w", b.images, b.step, t.Namespace, err)
			}
			buildEnv := append(append([]string{}, env...), "FLAI_KLUSTER_IMAGES="+b.images)
			if b.component != "" {
				buildEnv = append(buildEnv, "FLAI_KLUSTER_COMPONENT="+b.component)
			}
			if err := c.step(ctx, r, t, b.step, out, func(log io.Writer) error {
				return c.Exec.Command(ctx, Cmd{Dir: t.Dir, Env: buildEnv, Out: log, Name: b.cmd[0], Args: b.cmd[1:]})
			}); err != nil {
				return err
			}
			got, err := readImages(b.images)
			if err != nil {
				return err
			}
			list = append(list, got...)
			r.Images = list
		}
		if err := c.step(ctx, r, t, StepLoad, out, func(log io.Writer) error {
			return c.load(ctx, list, log)
		}); err != nil {
			return err
		}
	}
	return c.step(ctx, r, t, StepDeploy, out, func(log io.Writer) error {
		if err := c.Exec.Command(ctx, Cmd{Dir: t.Dir, Env: env, Out: log, Name: k.Deploy[0], Args: k.Deploy[1:]}); err != nil {
			return err
		}
		if t.Story == "" {
			return nil
		}
		return c.writeAgentKubeconfig(ctx, t, log)
	})
}

// step runs one step with its output to its log. It records the step in r
// as running, in place of the step's pending entry when r plans it, and
// writes r, so that a reader sees which step runs; then it records the
// step's outcome and writes r again, so that a job cut short leaves what it
// did.
func (c Cluster) step(ctx context.Context, r *Record, t Target, name string, out io.Writer, run func(log io.Writer) error) error {
	path := LogPath(c.Root, t.Namespace, name)
	at := -1
	for i := len(r.Steps) - 1; i >= 0; i-- {
		if r.Steps[i].Name == name && r.Steps[i].Outcome == OutcomePending {
			at = i
			break
		}
	}
	if at < 0 {
		r.Steps = append(r.Steps, Step{})
		at = len(r.Steps) - 1
	}
	r.Steps[at] = Step{Name: name, StartedAt: time.Now().UTC(), Outcome: OutcomeRunning, Log: path}
	fmt.Fprintf(out, "%s: %s, log %s\n", t.Namespace, name, path)
	f, err := os.Create(path)
	if err == nil {
		err = writeRecord(c.Root, *r)
		if err != nil {
			_ = f.Close()
		}
	} else {
		err = &StepError{Namespace: t.Namespace, Step: name, ExitCode: -1, Log: path, Err: fmt.Errorf("make its log: %w", err)}
	}
	if err != nil {
		r.Steps[at].EndedAt, r.Steps[at].Outcome, r.Steps[at].ExitCode = time.Now().UTC(), OutcomeFailed, -1
		return err
	}
	err = run(f)
	if cerr := f.Close(); cerr != nil && err == nil {
		err = cerr
	}
	s := &r.Steps[at]
	s.EndedAt, s.Outcome = time.Now().UTC(), OutcomePassed
	if err != nil {
		if ctx.Err() != nil {
			err = fmt.Errorf("%w (the job's time, kluster.timeout, ran out or it was stopped)", err)
		}
		s.Outcome, s.ExitCode = OutcomeFailed, ExitCode(err)
		err = &StepError{Namespace: t.Namespace, Step: name, ExitCode: s.ExitCode, Log: path, Err: err}
	}
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

// ErrNotDeployed says a story has no deployment record, so there is nothing
// to tear down.
var ErrNotDeployed = errors.New("no deployment is recorded for it, so there is nothing to tear down")

// Teardown tears story's namespace down: the block's teardown command, when
// it names one, run in dir with the namespace's variables at the commit its
// record names, then kubectl delete namespace, which passes over a namespace
// already gone, both logged to the namespace's teardown log. dir is the
// story's worktree, or the main checkout when the worktree is gone; empty
// means the main checkout. On success the record, the state folder, the
// logs, and the namespace's kubeconfig are removed, and the record as it was
// is returned. A failure keeps them, the record marked failed at the teardown
// step, and is a StepError. A story with no record is ErrNotDeployed and
// runs nothing. It runs while no other kluster job runs for the project.
// Neither its job nor a record it keeps says who asked; TeardownBy's do.
func (c Cluster) Teardown(ctx context.Context, m manifest.Manifest, story, dir string, out io.Writer) (Record, error) {
	return c.TeardownBy(ctx, m, "", story, dir, out)
}

// TeardownBy tears down as Teardown does, for by, who asked for it. The
// running job's By names by, and so does the AskedBy of a record a failed
// teardown keeps, whose last job it is; empty names no one, as Teardown does.
func (c Cluster) TeardownBy(ctx context.Context, m manifest.Manifest, by, story, dir string, out io.Writer) (Record, error) {
	if err := c.ready(m, false); err != nil {
		return Record{}, err
	}
	t, err := StoryTarget(story, dir, "")
	if err != nil {
		return Record{}, err
	}
	if t.Dir == "" {
		t.Dir = c.Root
	}
	release, err := acquire(c.Root, Job{Kind: "teardown", Story: t.Story, Namespace: t.Namespace, StartedAt: time.Now().UTC(), PID: os.Getpid(), By: by})
	if err != nil {
		return Record{}, err
	}
	defer release()
	r, ok, err := ReadRecord(c.Root, t.Namespace)
	if err != nil {
		return Record{}, err
	}
	if !ok {
		return Record{}, fmt.Errorf("%s: %w", t.Story, ErrNotDeployed)
	}
	t.Commit = r.Commit
	ctx, cancel := withTimeout(ctx, m)
	defer cancel()
	if err := os.MkdirAll(LogDir(c.Root, t.Namespace), 0o755); err != nil {
		return r, fmt.Errorf("make the folder %s for the namespace %s: %w", LogDir(c.Root, t.Namespace), t.Namespace, err)
	}
	kept := r
	kept.AskedBy = by
	kept.Steps = append(slices.Clone(r.Steps), Step{Name: StepTeardown, Outcome: OutcomePending})
	if err := writeRecord(c.Root, kept); err != nil {
		return r, err
	}
	err = c.step(ctx, &kept, t, StepTeardown, out, func(log io.Writer) error {
		return c.teardown(ctx, m.Kluster, t, log)
	})
	if err != nil {
		skipPending(&kept)
		kept.EndedAt, kept.Outcome, kept.Error = time.Now().UTC(), OutcomeFailed, err.Error()
		var se *StepError
		if errors.As(err, &se) {
			kept.FailedStep, kept.ExitCode = se.Step, se.ExitCode
		}
		if werr := writeRecord(c.Root, kept); werr != nil {
			return kept, errors.Join(err, werr)
		}
		return kept, err
	}
	for _, p := range []string{RecordPath(c.Root, t.Namespace), StatePath(c.Root, t.Namespace), LogDir(c.Root, t.Namespace), AgentKubeconfigPath(c.Root, t.Story)} {
		if err := os.RemoveAll(p); err != nil {
			return r, fmt.Errorf("remove %s, left by the namespace %s once it was torn down: %w; remove it by hand", p, t.Namespace, err)
		}
	}
	return r, nil
}

// teardown runs the block's teardown command, when it names one, then
// deletes t's namespace, with their output to log.
func (c Cluster) teardown(ctx context.Context, k *manifest.Kluster, t Target, log io.Writer) error {
	env := c.Env(t)
	if len(k.Teardown) > 0 && strings.TrimSpace(k.Teardown[0]) != "" {
		if err := c.Exec.Command(ctx, Cmd{Dir: t.Dir, Env: env, Out: log, Name: k.Teardown[0], Args: k.Teardown[1:]}); err != nil {
			return err
		}
	}
	if _, err := c.Exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf("kubectl is not on PATH, and flai deletes the namespace %s with it; install kubectl, then tear it down again", t.Namespace)
	}
	fmt.Fprintf(log, "deleting the namespace %s\n", t.Namespace)
	return c.Exec.Command(ctx, Cmd{Dir: c.Root, Env: env, Out: log, Name: "kubectl", Args: []string{"--context", c.Context(), "delete", "namespace", t.Namespace, "--wait", "--ignore-not-found"}})
}

// MainWorktree is the worktree of main's head a promotion builds and deploys
// in, which flai keeps under the cluster's folder.
func MainWorktree(root string) string { return filepath.Join(Dir(root), "main") }

// Promote deploys the head of the main branch branch into the default
// namespace: it makes, or moves, a detached worktree of that head at
// MainWorktree, then builds, loads, and deploys there as a story's deploy
// does, recorded under the default namespace. The worktree is prepared under
// the project's lock, so two promotions never race on it.
func (c Cluster) Promote(ctx context.Context, m manifest.Manifest, branch string, out io.Writer) (Record, error) {
	if err := c.ready(m, true); err != nil {
		return Record{}, err
	}
	ns := m.KlusterNamespace()
	if _, err := DefaultTarget(m, "", ""); err != nil {
		return Record{}, err
	}
	release, err := acquire(c.Root, Job{Kind: "promote", Namespace: ns, StartedAt: time.Now().UTC(), PID: os.Getpid()})
	if err != nil {
		return Record{}, err
	}
	defer release()
	ctx, cancel := withTimeout(ctx, m)
	defer cancel()
	commit, err := c.checkoutMain(ctx, branch, out)
	if err != nil {
		return Record{}, err
	}
	t, err := DefaultTarget(m, MainWorktree(c.Root), commit)
	if err != nil {
		return Record{}, err
	}
	return c.runDeploy(ctx, m, "", t, out)
}

// checkoutMain resolves branch's head in the main checkout and puts the
// worktree at MainWorktree there, detached, making it when it is absent.
func (c Cluster) checkoutMain(ctx context.Context, branch string, out io.Writer) (string, error) {
	git := func(what string, args ...string) (string, error) {
		var buf bytes.Buffer
		if err := c.Exec.Command(ctx, Cmd{Dir: c.Root, Out: &buf, Name: "git", Args: args}); err != nil {
			return "", fmt.Errorf("%s with git %s in %s: %w\n%s", what, strings.Join(args, " "), c.Root, err, strings.TrimSpace(buf.String()))
		}
		return strings.TrimSpace(buf.String()), nil
	}
	head, err := git("resolve the head of "+branch, "rev-parse", "--verify", branch+"^{commit}")
	if err != nil {
		return "", err
	}
	if !hexCommit.MatchString(head) {
		return "", fmt.Errorf("git rev-parse gave %q for the head of %s in %s, not a commit's hash; check the branch, then promote again", head, branch, c.Root)
	}
	wt := MainWorktree(c.Root)
	fmt.Fprintf(out, "checkout of %s at %s in %s\n", branch, head, wt)
	if _, err := os.Stat(filepath.Join(wt, ".git")); err == nil {
		if _, err := git("move the worktree "+wt+" to "+head, "-C", wt, "checkout", "--detach", "--force", head); err != nil {
			return "", err
		}
		return head, nil
	}
	if err := os.RemoveAll(wt); err != nil {
		return "", fmt.Errorf("remove %s, a folder that is no worktree, before making the worktree of main there: %w", wt, err)
	}
	if _, err := git("prune the worktrees", "worktree", "prune"); err != nil {
		return "", err
	}
	if _, err := git("make the worktree "+wt+" of "+head, "worktree", "add", "--detach", "--force", wt, head); err != nil {
		return "", err
	}
	return head, nil
}

// OutcomeInterrupted is a record left running by a job that no longer runs:
// it ended without recording its outcome.
const OutcomeInterrupted = "interrupted"

// Deployment is a namespace's record as the list shows it: whether it is the
// default namespace, the job running against it now, if any, and where its
// last job's output is; and, as List fills them for a story's namespace, the
// story's column and its branch's head.
type Deployment struct {
	Record
	Default bool `json:"default"`
	Job     *Job `json:"job,omitempty"`
	// Log is where the last job's output is, as List fills it: the log of
	// the step it failed at, else the folder of its steps' logs; empty from
	// Deployments (S-0663).
	Log string `json:"log,omitempty"`
	// Column is the story's column on the board, such as in-progress or
	// review; empty for the default namespace, for a story the board does not
	// hold, and from Deployments (S-0663).
	Column string `json:"column,omitempty"`
	// Head is the story branch's head now, and HeadError why it could not be
	// resolved; AtHead says whether Commit is Head. All are empty, and AtHead
	// false, for the default namespace and from Deployments (S-0663).
	Head      string `json:"head,omitempty"`
	HeadError string `json:"head_error,omitempty"`
	AtHead    bool   `json:"at_head"`
}

// Deployments lists the namespaces' records kept under root, the default
// namespace, defaultNS, first, then the stories' by their number, then any
// other by name. A record whose namespace a running job acts on carries that
// job; a record that still says running while no job runs for its namespace
// is shown interrupted. None is an empty list. List gives each story's
// column and head besides.
func Deployments(root, defaultNS string) ([]Deployment, error) {
	entries, err := os.ReadDir(filepath.Join(Dir(root), "deployments"))
	if errors.Is(err, os.ErrNotExist) {
		return []Deployment{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list the deployment records under %s: %w", filepath.Join(Dir(root), "deployments"), err)
	}
	job, running := Running(root)
	list := []Deployment{}
	for _, e := range entries {
		ns, ok := strings.CutSuffix(e.Name(), ".json")
		if !ok || e.IsDir() {
			continue
		}
		r, ok, err := ReadRecord(root, ns)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		d := Deployment{Record: r, Default: ns == defaultNS}
		if running && job.Namespace == ns {
			j := job
			d.Job = &j
		} else if d.Outcome == OutcomeRunning {
			d.Outcome = OutcomeInterrupted
		}
		list = append(list, d)
	}
	slices.SortFunc(list, func(a, b Deployment) int {
		ra, rb := rank(a), rank(b)
		if ra != rb {
			return ra - rb
		}
		if ra == 1 {
			if c := storyNumber(a.Namespace) - storyNumber(b.Namespace); c != 0 {
				return c
			}
		}
		return strings.Compare(a.Namespace, b.Namespace)
	})
	return list, nil
}

// rank orders the default namespace first, then stories', then others.
func rank(d Deployment) int {
	switch {
	case d.Default:
		return 0
	case storyNumber(d.Namespace) >= 0:
		return 1
	}
	return 2
}

// storyNumber is the number of a story namespace such as s-0578, -1 for any
// other.
func storyNumber(ns string) int {
	m := storyID.FindStringSubmatch(ns)
	if m == nil {
		return -1
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return -1
	}
	return n
}

// lastLog is where r's last job's output is: the log of the step it failed
// at, else the folder of its steps' logs under root.
func lastLog(root string, r Record) string {
	if r.FailedStep != "" {
		for i := len(r.Steps) - 1; i >= 0; i-- {
			if s := r.Steps[i]; s.Name == r.FailedStep && s.Log != "" {
				return s.Log
			}
		}
	}
	return LogDir(root, r.Namespace)
}

// inProgress is the column of a story whose namespace counts against
// kluster.agent_namespaces (ADR-0292 §4).
const inProgress = "in-progress"

// ListOptions is what List needs beyond the records under the root.
type ListOptions struct {
	// DefaultNamespace is the namespace main is promoted into.
	DefaultNamespace string
	// Limit is kluster.agent_namespaces, the most stories in progress that
	// hold a namespace at once.
	Limit int
	// Column is a story's column on the board, empty for one the board does
	// not hold; nil leaves every column empty and counts no namespace.
	Column func(story string) string
	// Git resolves a story branch's head in the root; nil leaves every head
	// empty.
	Git execx.Runner
}

// Listing is the list of the namespaces flai kluster deployments and the
// Deployments page show (S-0663): each namespace's record with its story's
// column and head, the stories in progress that hold a namespace against
// the limit, and the deploy requests refused for it, newest first.
type Listing struct {
	Deployments []Deployment `json:"deployments"`
	// Live is how many stories in progress hold a namespace, Holding which,
	// by their number: a story with a record, or with a deploy its agent
	// asked for that waits or runs (ADR-0292 §4).
	Live    int      `json:"live"`
	Holding []string `json:"holding"`
	// Limit is kluster.agent_namespaces.
	Limit int `json:"limit"`
	// Refusals are the deploy requests refused for the limit, newest first,
	// at most RefusalsKept.
	Refusals []Refusal `json:"refusals"`
}

// List lists the namespaces' records kept under root, as Deployments does,
// each story's with its column and its branch's head, and counts the stories
// in progress that hold a namespace against o.Limit: those with a record, and
// those with no record yet whose agent asked for a deploy that waits or runs.
// An asked-for job left by a flai serve that stopped counts until the next
// flai serve ends it interrupted.
func List(root string, o ListOptions) (Listing, error) {
	list, err := Deployments(root, o.DefaultNamespace)
	if err != nil {
		return Listing{}, err
	}
	column := o.Column
	if column == nil {
		column = func(string) string { return "" }
	}
	out := Listing{Deployments: list, Holding: []string{}, Limit: o.Limit}
	recorded := map[string]bool{}
	for i := range list {
		d := &list[i]
		d.Log = lastLog(root, d.Record)
		if d.Story == "" {
			continue
		}
		story := strings.ToUpper(d.Story)
		recorded[story] = true
		d.Column = column(story)
		if d.Column == inProgress {
			out.Holding = append(out.Holding, story)
		}
		if o.Git == nil {
			continue
		}
		if head, err := StoryHead(o.Git, root, story); err != nil {
			d.HeadError = err.Error()
		} else {
			d.Head, d.AtHead = head, head == d.Commit
		}
	}
	asked, err := askedDeploys(root)
	if err != nil {
		return Listing{}, err
	}
	for _, story := range asked {
		if !recorded[story] && column(story) == inProgress {
			out.Holding = append(out.Holding, story)
		}
	}
	slices.SortFunc(out.Holding, func(a, b string) int {
		return cmp.Or(storyNumber(a)-storyNumber(b), strings.Compare(a, b))
	})
	out.Live = len(out.Holding)
	if out.Refusals, err = Refusals(root); err != nil {
		return Listing{}, err
	}
	return out, nil
}

// askedDeploys are the stories whose agent asked for a deploy that waits or
// runs, as their job states under root say.
func askedDeploys(root string) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(RequestsDir(root), "*"+JobSuffix))
	if err != nil {
		return nil, fmt.Errorf("list the kluster jobs stories' agents asked for under %s: %w", RequestsDir(root), err)
	}
	var out []string
	for _, path := range paths {
		story := strings.TrimSuffix(filepath.Base(path), JobSuffix)
		j, ok, err := ReadAskedJob(root, story)
		if err != nil {
			return nil, err
		}
		if ok && j.Kind == KindDeploy && !j.Ended() {
			out = append(out, strings.ToUpper(j.Story))
		}
	}
	return out, nil
}

// Refusal is a deploy request flai serve refused because kluster.agent_namespaces
// stories in progress held a namespace already (ADR-0292 §4, S-0663).
type Refusal struct {
	Story string `json:"story"`
	// By is the agent that asked.
	By string    `json:"by"`
	At time.Time `json:"at"`
	// Limit is kluster.agent_namespaces when it was refused.
	Limit int `json:"limit"`
	// Reason is the refusal's words, naming the stories that held one.
	Reason string `json:"reason"`
}

// RefusalsKept is how many refusals for the limit flai keeps, the newest.
const RefusalsKept = 10

// RefusalsPath is the file of the refusals for the limit kept under root.
func RefusalsPath(root string) string { return filepath.Join(Dir(root), "refusals.json") }

func refusalsLockPath(root string) string { return filepath.Join(Dir(root), "refusals.lock") }

// KeepRefusal keeps r among the refusals for the limit under root, the
// newest first, dropping all but the RefusalsKept newest; the file is
// written whole or not at all.
func KeepRefusal(root string, r Refusal) error {
	unlock, err := filelock.Lock(refusalsLockPath(root), recordsLockWait)
	if err != nil {
		return fmt.Errorf("take the lock of the kluster refusals: %w", err)
	}
	defer unlock()
	kept, err := Refusals(root)
	if err != nil {
		return err
	}
	kept = append([]Refusal{r}, kept...)
	if len(kept) > RefusalsKept {
		kept = kept[:RefusalsKept]
	}
	data, err := json.MarshalIndent(kept, "", "  ")
	if err != nil {
		return fmt.Errorf("encode the kluster refusals: %w", err)
	}
	if err := atomicfile.WriteFile(RefusalsPath(root), append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write the kluster refusals %s: %w; check that the project's .flai-cache is writable", RefusalsPath(root), err)
	}
	return nil
}

// Refusals are the refusals for the limit kept under root, the newest first;
// none is an empty list.
func Refusals(root string) ([]Refusal, error) {
	data, err := os.ReadFile(RefusalsPath(root))
	if errors.Is(err, os.ErrNotExist) {
		return []Refusal{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read the kluster refusals: %w", err)
	}
	out := []Refusal{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("read the kluster refusals %s: %w; delete it, and the next refusal writes it again", RefusalsPath(root), err)
	}
	return out, nil
}
