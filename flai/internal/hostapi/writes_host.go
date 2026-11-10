package hostapi

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// hostSpecs are the writes that publish, run the dashboard and the host, and start agents and checks.
func hostSpecs() map[string]spec {
	return map[string]spec{
		// publish.run: the push host action (ADR-0031, S-0078), and the one way
		// the dashboard sends accepted work to the remote (ADR-0067): flai
		// release --pending as the operator, applying, tagging, and pushing
		// everything accepted and unreleased. Never forced; when the remote has
		// moved it refuses (exit 3) and says to fetch and merge. Publisher runs
		// the same for the orchestrator's release_publish (S-0222).
		"publish.run": {action: ActionPush, exits: publishExits, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return slices.Clone(publishArgs), "", nil
		}},

		// dashboard.status: a read of what flai dashboard status already
		// reports, so the board can show the version running and who else it
		// serves without a host action (S-0081).
		"dashboard.status": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "status"}, "", nil
		}),

		// dashboard.check: a read that pulls the configured image and
		// compares it to what is running, changing nothing (S-0081).
		"dashboard.check": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "check"}, "", nil
		}),

		// dashboard.versions: a read of the published dashboard releases,
		// newest first, marking the one running (S-0298).
		"dashboard.versions": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "versions"}, "", nil
		}),

		// dashboard.restart and dashboard.upgrade: the host action (S-0081),
		// running the same commands the operator's own shell does. Both
		// stream progress: they take real Docker time, unlike a move or an
		// edit. The dashboard never names an image. dashboard.upgrade may
		// name a tag, a bare X.Y.Z, which flai deploys only when it is a
		// published dashboard release (ADR-0117 §3, --published); without
		// one it upgrades to what this host's own configuration names. A
		// chosen release keeps running through restarts until an upgrade
		// without a tag or a start after a stop (ADR-0118).
		"dashboard.restart": {action: ActionDashboard, progress: true, describe: describeDashboardRestart, detachTimeout: 90 * time.Second, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "restart"}, "", nil
		}},
		"dashboard.upgrade": {action: ActionDashboard, progress: true, detachTimeout: 6 * time.Minute, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			tag, e := upgradeRelease(raw, "tag")
			if e != nil {
				return nil, "", e
			}
			if tag == "" {
				return []string{"dashboard", "upgrade"}, "", nil
			}
			return []string{"dashboard", "upgrade", "--published", "--tag=" + tag}, "", nil
		}, asked: func(raw json.RawMessage) func(any, *channel.Error) (string, string) {
			tag, _ := upgradeRelease(raw, "tag")
			return describeDashboardUpgrade(tag)
		}},

		// dashboard.stop: the host action, the same command the operator's
		// own flai dashboard stop runs: unregisters this project, and stops
		// the container only when it was the last one registered.
		"dashboard.stop": {action: ActionDashboard, describe: describeDashboardStop, detachTimeout: 60 * time.Second, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "stop"}, "", nil
		}},

		// host.status and host.check: reads of flai host (S-0106): the host,
		// and each process it keeps; whether a newer flai is published.
		"host.status": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"host", "status"}, "", nil
		}),
		"host.check": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"host", "check"}, "", nil
		}),
		// host.versions: a read of the published flai releases, newest first,
		// marking the one installed and any below a project's flai.minimum
		// (S-0298).
		"host.versions": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"host", "versions"}, "", nil
		}),

		// host.process and host.upgrade: the host action (S-0106). Both are
		// detached: restarting flai serve, or the host restarting on a new
		// flai, ends the connection the request came on. host.upgrade may
		// name a version, a bare X.Y.Z, which the host installs only when it
		// is a published flai release (ADR-0117 §3); without one it installs
		// the newest.
		"host.process": {action: ActionHost, describe: describeHost(""), detachTimeout: 60 * time.Second, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Process string `json:"process"`
				Action  string `json:"action"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Process != "serve" && in.Process != "mcp" && in.Process != "all" {
				return nil, "", bad("process must be serve, mcp, or all")
			}
			if in.Action != "start" && in.Action != "stop" && in.Action != "restart" {
				return nil, "", bad("action must be start, stop, or restart")
			}
			return []string{"host", in.Action, in.Process}, "", nil
		}, say: func(raw json.RawMessage) string {
			var in struct{ Process, Action string }
			_ = json.Unmarshal(raw, &in)
			return in.Action + " " + in.Process
		}},
		"host.upgrade": {action: ActionHost, detachTimeout: 6 * time.Minute, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			version, e := upgradeRelease(raw, "version")
			if e != nil {
				return nil, "", e
			}
			if version == "" {
				return []string{"host", "upgrade"}, "", nil
			}
			return []string{"host", "upgrade", "--version=" + version}, "", nil
		}, asked: func(raw json.RawMessage) func(any, *channel.Error) (string, string) {
			version, _ := upgradeRelease(raw, "version")
			return describeHost(version)
		}},

		// checks.status: a read of what flai checks status already reports:
		// the current or last run for a story (S-0082).
		"checks.status": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"checks", "status", in.ID}, "", nil
		}),

		// checks.tail: a read, progress-streaming, of new lines in a run's
		// log since a byte offset, waiting a little for more while the run
		// is active (S-0082). Each line arrives as a log event, the same
		// machinery push and dashboard already stream progress with.
		"checks.tail": {reads: true, progress: true, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID   string `json:"id"`
				From int    `json:"from"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if in.From < 0 {
				return nil, "", bad("from must not be negative")
			}
			return []string{"checks", "tail", in.ID, "--from=" + strconv.Itoa(in.From), "--wait=20"}, "", nil
		}},

		// checks.run and checks.cancel: the host action (S-0082), running the
		// same commands the operator's own shell does. run streams progress
		// and is detached generously: the real time limit is the operator's
		// own (flai serve checks timeout), enforced by flai checks run
		// itself, not by this request surviving that long.
		"checks.run": {action: ActionChecks, progress: true, describe: describeChecksRun, detachTimeout: 2 * time.Hour, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"checks", "run", in.ID}, "", nil
		}},
		"checks.cancel": {action: ActionChecks, describe: describeChecksRun, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"checks", "cancel", in.ID}, "", nil
		}},
		// test.run: flai test, the project's test and lint tiers for paths,
		// in a story's worktree or the main checkout (S-0273), under the same
		// host action as checks.run, since it too runs the project's commands
		// on the host. It is detached and streams progress as checks.run is;
		// flai test has no time limit of its own, so the ceiling is the
		// detach timeout. A tier that failed (exit 1) is an answer; exit 2,
		// flai test unable to answer, is an error with flai's words.
		"test.run": {action: ActionChecks, progress: true, describe: describeTestRun, detachTimeout: 2 * time.Hour,
			dir: testDir, answers: []int{1}, judge: testResult, build: testArgs},
		// verify.run: flai verify, a story's close-out checks in its worktree
		// (S-0270), under the checks host action as test.run is, since it
		// runs the project's tiers on the host. flai verify finds the story's
		// worktree itself, so it runs in the project. It is detached and
		// streams progress as test.run is. A step that failed (exit 1) is an
		// answer; exit 2, flai verify unable to answer, is an error with
		// flai's words.
		"verify.run": {action: ActionChecks, progress: true, describe: describeVerifyRun, detachTimeout: 2 * time.Hour,
			answers: []int{1}, judge: verifyResult, build: verifyArgs},
		// verify.status: a read of the story's last verification, as flai
		// verify --last prints it from .flai-cache/verify: the report, or
		// null when the story has none. It runs nothing, so, as
		// checks.status, it is not gated.
		"verify.status": {reads: true, judge: verifyStored, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needStory(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"verify", in.ID, "--last"}, "", nil
		}},
		// agent.restart: a new agent for a story whose agent dropped or
		// failed (S-0116, ADR-0043); flai serve agent restart judges whether
		// it may, and says why not.
		"agent.restart": agentNow("restart", "restarted", ""),
		// agent.start: a ready story's agent now, whatever the launcher's
		// rules say about when (S-0115); flai serve agent start judges
		// whether it may, and says why not.
		"agent.start": agentNow("start", "started", ""),
		// agent.commit: an agent to commit what a story in review left
		// uncommitted in its worktree (S-0140); flai serve agent commit
		// judges whether it may, and says why not.
		"agent.commit": agentNow("commit", "started", " to commit what its worktree holds"),
		// agent.stop: the story's agent ended, on the operator's word from
		// the activity page (S-0170); flai serve agent stop judges whether
		// there is one to stop, and says why not.
		"agent.stop": {action: ActionAgent, describe: describeAgentStop, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			id, e := storyOf(raw)
			if e != nil {
				return nil, "", e
			}
			return []string{"serve", "agent", "stop", id}, "", nil
		}},
		// plan.run: the planner for an epic or a story, on the operator's
		// word from the item's page (S-0208); flai plan judges whether it
		// may, and says why not.
		"plan.run": {action: ActionPlan, describe: describePlan, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if strings.HasPrefix(in.ID, "T-") {
				return nil, "", bad("%s is a task; the planner plans an epic or a story", in.ID)
			}
			return []string{"plan", in.ID}, "", nil
		}},
		// analyze.run: the analyzer for the project, with a focus or none,
		// on the operator's word (S-0223); flai analyze judges whether it
		// may, and says why not.
		"analyze.run": {action: ActionAnalyze, describe: describeAnalyze, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Focus string `json:"focus"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Focus == "" {
				return []string{"analyze"}, "", nil
			}
			if !slices.Contains(harness.Focuses, in.Focus) {
				return nil, "", bad("the analyzer takes the focus %s, or none for all of them, and %q is none of them", strings.Join(harness.Focuses, ", "), in.Focus)
			}
			return []string{"analyze", "--focus=" + in.Focus}, "", nil
		}},
		// curate.run: the curator for the project, from a since or from its
		// last run, on the operator's word (S-0431, ADR-0196); flai curate
		// judges whether it may, and says why not.
		"curate.run": {action: ActionCurate, describe: describeCurate, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Since string `json:"since"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Since == "" {
				return []string{"curate"}, "", nil
			}
			if !curateSinceReads(in.Since) {
				return nil, "", bad("the curator takes a since as a YYYY-MM-DD day in UTC or an RFC 3339 time, such as 2026-10-01 or 2026-10-01T09:00:00Z, and %q is neither", in.Since)
			}
			return []string{"curate", "--since=" + in.Since}, "", nil
		}},
		// role.run: a role the project defines (S-0417, ADR-0207), on an item
		// its file takes or on none, with inputs, on the operator's word,
		// behind the host action its file names; flai role run judges
		// whether it may, and says why not.
		"role.run": {actionOf: roleAction, describe: describeRole, build: roleRunArgs},
		// orchestrate.stop: the orchestrator's run ended and the orchestrator
		// held stopped, so that flai serve does not start it again while the
		// orchestrate action stays on, on the operator's word from its page
		// (S-0228); orchestrate.start lifts the hold and starts it. flai serve
		// orchestrate stop and start judge whether they may, and say why not.
		"orchestrate.stop":  orchestrateNow("stop"),
		"orchestrate.start": orchestrateNow("start"),
	}
}
