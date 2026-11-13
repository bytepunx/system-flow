package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// maxKlusterWait is the most wait_seconds a kluster tool waits for its job.
const maxKlusterWait = 45

// The descriptions of the tools kluster_deploy and kluster_teardown.
var (
	klusterDeployDescription   = klusterJobDescription(kluster.KindDeploy, "build, load, and deploy the session's own story at its branch's head into its kluster namespace")
	klusterTeardownDescription = klusterJobDescription(kluster.KindTeardown, "tear the session's own story's kluster namespace down")
)

// klusterJobDescription is the description of the tool that asks flai serve
// for the kluster job of kind, what being the job in words.
func klusterJobDescription(kind, what string) string {
	return fmt.Sprintf("Ask flai serve to %s as a job (ADR-0292), as flai kluster %s S-nnnn --job does: the tool runs nothing itself, and it is for the story's agent, for its own story alone, FLAI_STORY, else the story in an agent name of the form agent-S-nnnn. story is optional and is refused when it is not the session's own; a session with no story of its own names it. "+
		"It answers at once with answer, flai serve's answer: job, the job's id; kind; story; ahead, the jobs of the queue before it; started, true when this request queued it, false when a %s job of the story was already queued or running and is answered instead, queuing nothing; log, where its output goes; and, for a story in review whose deploy passed at its head, recorded true with the record's commit and outcome, nothing queued. "+
		"wait_seconds, 0 to %d, then waits that long for the job to end and answers job too, its state (id, outcome passed, failed, refused, or interrupted, commit, failed_step, exit_code, error, each step with its outcome, exit code, and log, and log), and ended, whether it ended in the wait; more than %d is refused. "+
		"While the job has not ended, do the work that needs no namespace and hold wait_for_events for the kluster event of your story, which reports the job, its kind, its commit, its outcome, the failed step, its exit code, and the log, and is not told to end while the job waits or runs; flai kluster deployments shows it too. A failure is yours to fix: read the step's log at the path it names, change the code, commit, and ask again."+
		"Refused, with the reason and nothing queued: for another story, while the kluster host action is off (the operator's flai serve enable kluster), with no kluster.deploy set, over kluster.agent_namespaces, while the story's worktree has uncommitted changes (commit first), and for a story not in progress. When flai serve is not running, it takes no request within ten seconds and the call says so: ask the operator.",
		what, kind, kind, maxKlusterWait, maxKlusterWait)
}

// KlusterJobIn asks flai serve for a kluster job of the session's own story.
type KlusterJobIn struct {
	Project     string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story       string `json:"story,omitempty" jsonschema:"the story (S-nnnn, any zero padding): the session's own story, which is the default; another is refused"`
	WaitSeconds int    `json:"wait_seconds,omitempty" jsonschema:"seconds to wait for the job to end, 0 to 45; 0, the default, answers at once"`
}

func (in KlusterJobIn) project() string { return in.Project }

// KlusterJobOut is flai serve's answer to a kluster request and, with
// wait_seconds, the job's state and whether it ended in the wait, as flai
// kluster deploy and teardown --job --json print them.
type KlusterJobOut struct {
	Answer kluster.Answer    `json:"answer"`
	Job    *kluster.AskedJob `json:"job,omitempty"`
	Ended  *bool             `json:"ended,omitempty"`
}

// klusterDeploy asks flai serve to deploy the session's own story.
func (s *server) klusterDeploy(ctx context.Context, _ *mcp.CallToolRequest, in KlusterJobIn) (*mcp.CallToolResult, KlusterJobOut, error) {
	out, err := s.klusterJob(ctx, kluster.KindDeploy, in)
	return nil, out, err
}

// klusterTeardown asks flai serve to tear the session's own story's
// namespace down.
func (s *server) klusterTeardown(ctx context.Context, _ *mcp.CallToolRequest, in KlusterJobIn) (*mcp.CallToolResult, KlusterJobOut, error) {
	out, err := s.klusterJob(ctx, kluster.KindTeardown, in)
	return nil, out, err
}

// klusterJob asks flai serve for the kluster job of kind for the story in,
// or the session's own, and answers its answer, then, with wait_seconds,
// the job's state once it ended or the wait ran out. A refusal and a request
// flai serve did not take are errors.
func (s *server) klusterJob(ctx context.Context, kind string, in KlusterJobIn) (KlusterJobOut, error) {
	if in.WaitSeconds < 0 || in.WaitSeconds > maxKlusterWait {
		return KlusterJobOut{}, fmt.Errorf("wait_seconds %d is not 0 to %d; give at most %d, and hold wait_for_events for the kluster event to wait longer", in.WaitSeconds, maxKlusterWait, maxKlusterWait)
	}
	story, err := s.klusterStory(kind, in.Story)
	if err != nil {
		return KlusterJobOut{}, err
	}
	root := projectRoot(s.repo)
	ans, err := kluster.Ask(ctx, root, kluster.Request{Story: story, Kind: kind, By: s.agent}, s.klusterWait)
	if err != nil {
		return KlusterJobOut{}, err
	}
	if ans.Refused() {
		return KlusterJobOut{}, errors.New(ans.Reason)
	}
	out := KlusterJobOut{Answer: ans}
	// A story in review may be answered with the review's own deploy, which
	// keeps no job state to wait on.
	if in.WaitSeconds == 0 || ans.Job == "" || ans.Recorded || !kluster.HasAskedJob(root, ans.Story, ans.Job) {
		return out, nil
	}
	j, ended, err := kluster.Await(ctx, root, ans.Story, ans.Job, time.Duration(in.WaitSeconds)*time.Second, s.klusterWait.Poll)
	if err != nil {
		return KlusterJobOut{}, err
	}
	out.Job, out.Ended = &j, &ended
	return out, nil
}

// agentStoryName is the name flai serve gives a story's agent, agent-S-nnnn.
var agentStoryName = regexp.MustCompile(`^agent-(S-\d+)$`)

// ownStory is the session's own story, canonical: FLAI_STORY, else the story
// in the agent's name of the form agent-S-nnnn; "" with neither.
func (s *server) ownStory() string {
	if s.story != "" {
		return workitem.CanonicalID(s.story)
	}
	if m := agentStoryName.FindStringSubmatch(s.agent); m != nil {
		return workitem.CanonicalID(m[1])
	}
	return ""
}

// klusterStory is the story a kluster request of kind is for: given, which
// must be the session's own when it has one, else the session's own.
func (s *server) klusterStory(kind, given string) (string, error) {
	own := s.ownStory()
	given = strings.TrimSpace(given)
	switch {
	case given == "" && own == "":
		return "", fmt.Errorf("name the story whose kluster %s to ask for, such as S-0042: this session has no story of its own", kind)
	case given == "":
		return own, nil
	case own != "" && workitem.CanonicalID(given) != own:
		return "", fmt.Errorf("this session works %s, and a story's agent asks flai serve for its own story's kluster %s alone (ADR-0292 §1), so nothing was asked for %s; leave story out, or leave %s's namespace to its own agent", own, kind, workitem.CanonicalID(given), workitem.CanonicalID(given))
	}
	return workitem.CanonicalID(given), nil
}
