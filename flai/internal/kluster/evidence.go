package kluster

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// HeadJudgement is how a record's commit stands to a story branch's head,
// as ADR-0292 §3 judges it: at the head when it is the head, or when every
// path that changed from it to the head is under the manifest's design,
// docs, or wip folder.
type HeadJudgement struct {
	// Commit is the commit judged, and Head the head it was judged against.
	Commit string `json:"commit"`
	Head   string `json:"head"`
	// AtHead says the commit counts as the head; Exact that it is the head.
	AtHead bool `json:"at_head"`
	Exact  bool `json:"exact,omitempty"`
	// Outside is up to maxOutside of the paths outside the design, docs,
	// and wip folders that changed from Commit to Head, when that is why
	// it is not at the head.
	Outside []string `json:"outside,omitempty"`
	// Why says why the commit is not at the head; empty when it is.
	Why string `json:"why,omitempty"`
}

// maxOutside is the most paths a HeadJudgement names in Outside.
const maxOutside = 5

// AtHead judges whether commit, a record's, is at head, a story branch's,
// with git in the checkout root (ADR-0292 §3): it is when it is head, or
// when commit is in head's history and every path that differs between
// them is under the folders m's layout names design, docs, and wip. An
// empty commit or head, or a commit git does not find in head's history,
// is not at the head; only a diff git cannot take is an error.
func AtHead(git execx.Runner, root string, m manifest.Manifest, commit, head string) (HeadJudgement, error) {
	j := HeadJudgement{Commit: commit, Head: head}
	switch {
	case commit == "":
		j.Why = "no commit is recorded"
		return j, nil
	case head == "":
		j.Why = "the story has no branch head to judge it against"
		return j, nil
	case commit == head:
		j.AtHead, j.Exact = true, true
		return j, nil
	}
	if _, err := git.Run(root, "git", "merge-base", "--is-ancestor", commit, head); err != nil {
		j.Why = fmt.Sprintf("%s is not in the history of the head %s: the branch was rebased or reset since", shortCommit(commit), shortCommit(head))
		return j, nil //nolint:nilerr // exit 1 is no ancestor and 128 no such commit: either way not at the head
	}
	out, err := git.Run(root, "git", "diff", "--name-only", "--no-renames", "-z", commit, head)
	if err != nil {
		return HeadJudgement{}, fmt.Errorf("list the paths that changed from %s to the head %s with git diff in %s, to tell whether its deploy stands at the head: %w; check that git can read the repository there", shortCommit(commit), shortCommit(head), root, err)
	}
	folders := headFolders(m)
	var outside []string
	for p := range strings.SplitSeq(out, "\x00") {
		if p = strings.TrimSpace(p); p != "" && !under(p, folders) {
			outside = append(outside, p)
		}
	}
	if len(outside) == 0 {
		j.AtHead = true
		return j, nil
	}
	j.Outside = outside[:min(len(outside), maxOutside)]
	j.Why = fmt.Sprintf("%d %s outside %s changed from %s to the head %s, such as %s", len(outside), plural(len(outside), "path", "paths"), joinFolders(folders), shortCommit(commit), shortCommit(head), strings.Join(j.Outside, ", "))
	return j, nil
}

// headFolders are the folders whose changes leave a record at the head:
// the layout's design, docs, and wip, each its key when the layout names
// none, cleaned of a leading ./ and a trailing /.
func headFolders(m manifest.Manifest) []string {
	var folders []string
	for _, key := range []string{"design", "docs", "wip"} {
		f := strings.TrimSpace(m.Layout[key])
		if f == "" {
			f = key
		}
		if f = path.Clean(strings.ReplaceAll(f, `\`, "/")); f != "." && f != "/" {
			folders = append(folders, strings.TrimPrefix(f, "/"))
		}
	}
	return folders
}

// under says whether the repository path p is in one of folders.
func under(p string, folders []string) bool {
	for _, f := range folders {
		if p == f || strings.HasPrefix(p, f+"/") {
			return true
		}
	}
	return false
}

// joinFolders names folders as a sentence does: "design, docs, and wip".
func joinFolders(folders []string) string {
	switch len(folders) {
	case 0:
		return "no folder"
	case 1:
		return folders[0]
	case 2:
		return folders[0] + " and " + folders[1]
	}
	return strings.Join(folders[:len(folders)-1], ", ") + ", and " + folders[len(folders)-1]
}

// plural is one when n is 1 and many otherwise.
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// DeployEvidence is a namespace's last deploy as its record keeps it, and
// how its commit stands to the story branch's head.
type DeployEvidence struct {
	Commit string `json:"commit"`
	// Outcome is running, passed, or failed; FailedStep and ExitCode say
	// where and how a failed one failed, and Error what it said.
	Outcome    string `json:"outcome"`
	FailedStep string `json:"failed_step,omitempty"`
	ExitCode   int    `json:"exit_code,omitempty"`
	Error      string `json:"error,omitempty"`
	// Job and Said name the record's last job, a deploy or a teardown, and
	// its outcome in words: "failed at the build step", "is still running".
	Job  string `json:"job"`
	Said string `json:"said"`
	// AskedBy is who asked for the job: a story's agent, flai serve, or
	// empty for one run in a shell.
	AskedBy   string    `json:"asked_by,omitempty"`
	StartedAt time.Time `json:"started_at,omitzero"`
	EndedAt   time.Time `json:"ended_at,omitzero"`
	// Steps are the job's steps, each with its outcome and log path.
	Steps []Step `json:"steps,omitempty"`
	// Head is how Commit stands to the head, and PassedAtHead says the
	// deploy passed there.
	Head         HeadJudgement `json:"head"`
	PassedAtHead bool          `json:"passed_at_head"`
}

// ChecksEvidence is the last checks run against a namespace since its
// deploy, as its record keeps it, and how its commit stands to the head.
type ChecksEvidence struct {
	Commit string `json:"commit"`
	// Outcome is passed, failed, timed-out, or cancelled.
	Outcome string `json:"outcome"`
	// Log is the run's log.
	Log string `json:"log,omitempty"`
	// Head is how Commit stands to the head, and PassedAtHead says the run
	// passed there.
	Head         HeadJudgement `json:"head"`
	PassedAtHead bool          `json:"passed_at_head"`
}

// EvidenceCommands are what a story's agent deploys and checks its
// namespace with, each with the story's ID written in (ADR-0292 §1).
type EvidenceCommands struct {
	Deploy     string `json:"deploy"`      // the CLI command that asks flai serve to deploy it
	DeployTool string `json:"deploy_tool"` // the MCP tool that does the same
	Checks     string `json:"checks"`      // the CLI command that asks flai serve to run its checks
	ChecksTool string `json:"checks_tool"` // the MCP tool that does the same
}

// Evidence is what a story's kluster namespace shows of its change: its
// last deploy and checks run, and whether each passed at the head of the
// story's branch (ADR-0292 §3, §6). The move to review, the close-out, the
// review page, and an acceptance read it.
type Evidence struct {
	// Applies is false while the manifest has no kluster: block, the kluster
	// host action is off, or the block names no deploy; Why then says which.
	Applies bool   `json:"applies"`
	Why     string `json:"why,omitempty"`
	Story   string `json:"story"`
	// Namespace is the story's namespace, and Head its branch's head, empty
	// when the branch cannot be resolved.
	Namespace string `json:"namespace,omitempty"`
	Head      string `json:"head,omitempty"`
	// ChecksOn says the checks host action is on, so that the checks must
	// pass at the head too.
	ChecksOn bool `json:"checks_on"`
	// Deploy is the last deploy, nil when none is recorded; Checks is the
	// last checks run against it, nil when none has ended.
	Deploy *DeployEvidence `json:"deploy,omitempty"`
	Checks *ChecksEvidence `json:"checks,omitempty"`
	// Passed says the deploy passed at the head, and the checks too while
	// ChecksOn.
	Passed bool `json:"passed"`
	// Commands are what deploy and check the namespace, while it Applies.
	Commands *EvidenceCommands `json:"commands,omitempty"`
}

// ReadEvidence reads story's namespace's record under root, the main
// checkout, and judges it against the head of the story's branch there with
// git. klusterOn and checksOn say whether the kluster and checks host
// actions are on for the project. A branch git cannot resolve leaves Head
// empty and nothing at the head; only an invalid story ID, a record that
// cannot be read, or a diff git cannot take is an error.
func ReadEvidence(git execx.Runner, root string, m manifest.Manifest, story string, klusterOn, checksOn bool) (Evidence, error) {
	ns, err := StoryNamespace(story)
	if err != nil {
		return Evidence{}, err
	}
	story = strings.ToUpper(story)
	e := Evidence{Story: story, Namespace: ns, ChecksOn: checksOn}
	switch {
	case m.Kluster == nil:
		e.Why = "the manifest has no kluster: block"
		return e, nil
	case !klusterOn:
		e.Why = "the kluster host action is off for this project"
		return e, nil
	case len(m.Kluster.Deploy) == 0:
		e.Why = "the manifest's kluster: block names no deploy"
		return e, nil
	}
	e.Applies = true
	e.Commands = &EvidenceCommands{
		Deploy:     "flai kluster deploy " + story + " --job",
		DeployTool: "kluster_deploy",
		Checks:     "flai checks run " + story + " --job",
		ChecksTool: "kluster_checks",
	}
	r, found, err := ReadRecord(root, ns)
	if err != nil {
		return Evidence{}, fmt.Errorf("read %s's kluster evidence: %w", story, err)
	}
	e.Head, _ = StoryHead(git, root, story)
	if !found {
		return e, nil
	}
	job, said := lastJob(root, r)
	d := &DeployEvidence{
		Commit: r.Commit, Outcome: r.Outcome, FailedStep: r.FailedStep, ExitCode: r.ExitCode, Error: r.Error,
		Job: job, Said: said, AskedBy: r.AskedBy, StartedAt: r.StartedAt, EndedAt: r.EndedAt, Steps: r.Steps,
	}
	if d.Head, err = AtHead(git, root, m, r.Commit, e.Head); err != nil {
		return Evidence{}, err
	}
	d.PassedAtHead = d.Head.AtHead && r.Outcome == OutcomePassed && job == "deploy"
	e.Deploy = d
	if r.ChecksOutcome != "" {
		c := &ChecksEvidence{Commit: r.ChecksCommit, Outcome: r.ChecksOutcome, Log: r.ChecksLog, Head: d.Head}
		if r.ChecksCommit != r.Commit {
			if c.Head, err = AtHead(git, root, m, r.ChecksCommit, e.Head); err != nil {
				return Evidence{}, err
			}
		}
		c.PassedAtHead = c.Head.AtHead && r.ChecksOutcome == OutcomePassed
		e.Checks = c
	}
	e.Passed = d.PassedAtHead && (!checksOn || e.Checks != nil && e.Checks.PassedAtHead)
	return e, nil
}

// Summary is e in one sentence, as the move to review warns with it and
// the close-out prints it: what was deployed and checked at the head, or
// why not.
func (e Evidence) Summary() string {
	who := fmt.Sprintf("%s's namespace %s", e.Story, e.Namespace)
	d, c := e.Deploy, e.Checks
	switch {
	case !e.Applies:
		return fmt.Sprintf("no kluster evidence applies to %s: %s", e.Story, e.Why)
	case e.Head == "":
		return fmt.Sprintf("%s has no branch story/%s, so %s cannot be deployed at its head", e.Story, e.Story, who)
	case d == nil:
		return fmt.Sprintf("%s is not deployed at the head %s: no deploy of it is recorded", who, shortCommit(e.Head))
	case !d.Head.AtHead:
		return fmt.Sprintf("%s is not deployed at the head %s: its last %s, at %s, %s, and %s", who, shortCommit(e.Head), d.Job, shortCommit(d.Commit), d.Said, d.Head.Why)
	case !d.PassedAtHead:
		return fmt.Sprintf("%s is not deployed at the head %s: its last %s there%s %s", who, shortCommit(e.Head), d.Job, d.Head.since(), d.Said)
	case c == nil && e.ChecksOn:
		return fmt.Sprintf("%s is deployed at the head %s%s, and no checks run against it has ended there", who, shortCommit(e.Head), d.Head.since())
	case c == nil:
		return fmt.Sprintf("%s is deployed at the head %s%s; the checks host action is off", who, shortCommit(e.Head), d.Head.since())
	case !c.Head.AtHead:
		return fmt.Sprintf("%s is deployed at the head %s%s, and its last checks run, at %s, ended %s, and %s", who, shortCommit(e.Head), d.Head.since(), shortCommit(c.Commit), c.Outcome, c.Head.Why)
	case !c.PassedAtHead:
		return fmt.Sprintf("%s is deployed at the head %s%s, and its checks run there ended %s", who, shortCommit(e.Head), d.Head.since(), c.Outcome)
	}
	return fmt.Sprintf("%s is deployed and checked at the head %s%s", who, shortCommit(e.Head), d.Head.since())
}

// Next is what the story's agent runs to make e pass at the head, naming
// the command and the MCP tool; empty when e passed or does not apply.
func (e Evidence) Next() string {
	switch {
	case !e.Applies || e.Passed || e.Head == "":
		return ""
	case e.Deploy == nil || !e.Deploy.PassedAtHead:
		return fmt.Sprintf("deploy it at the head with %s, or the MCP tool %s, which runs its checks once the deploy passes", e.Commands.Deploy, e.Commands.DeployTool)
	}
	return fmt.Sprintf("run its checks at the head with %s, or the MCP tool %s", e.Commands.Checks, e.Commands.ChecksTool)
}

// since says, for a judgement at the head but not exactly it, the commit it
// was at; empty for the head itself.
func (j HeadJudgement) since() string {
	if !j.AtHead || j.Exact {
		return ""
	}
	return fmt.Sprintf(" (at %s, followed only by changes under the design, docs, and wip folders)", shortCommit(j.Commit))
}
