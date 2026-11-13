package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// checksProject is a project with a real worktree directory for one story
// (RunChecks only needs the directory to exist, not a real git worktree)
// and, on the host's own config, one named check that always passes.
func checksProject(t *testing.T) (root, story string) {
	t.Helper()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root = tempProject(t)
	story = "S-0001"
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.WorktreePath(story), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--name", "ok", "--", "true"); code != 0 {
		t.Fatal("set")
	}
	return root, story
}

func TestChecksRunRefusesSomethingThatIsNotAStoryID(t *testing.T) {
	t.Parallel()
	root := tempProject(t)
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	if _, errOut, code := runIn(t, root, "checks", "run", "not-a-story"); code == 0 || !strings.Contains(errOut, "is not a story's ID") {
		t.Errorf("run: %d %s", code, errOut)
	}
}

func TestChecksStatusBeforeAnyRun(t *testing.T) {
	t.Parallel()
	root, story := checksProject(t)
	out, _, code := runIn(t, root, "checks", "status", story, "--json")
	if code != 0 {
		t.Fatalf("status: %s", out)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got["running"] != false || got["started"] != "" || got["outcome"] != nil {
		t.Errorf("nothing has run yet: %+v", got)
	}
}

func TestChecksRunEndToEndThroughTheCLI(t *testing.T) {
	t.Parallel()
	root, story := checksProject(t)
	out, errOut, code := runIn(t, root, "checks", "run", story, "--json")
	if code != 0 {
		t.Fatalf("run: %d %s %s", code, out, errOut)
	}
	var run map[string]any
	if err := json.Unmarshal([]byte(out), &run); err != nil {
		t.Fatal(err)
	}
	if run["outcome"] != "passed" || run["running"] != false {
		t.Fatalf("run: %+v", run)
	}
	steps, ok := run["steps"].([]any)
	if !ok || len(steps) != 1 {
		t.Fatalf("steps: %+v", run["steps"])
	}

	// status now shows the same outcome, and tail can read the log from 0
	statusOut, _, code := runIn(t, root, "checks", "status", story, "--json")
	if code != 0 || !strings.Contains(statusOut, `"outcome": "passed"`) {
		t.Errorf("status after: %d %s", code, statusOut)
	}
	tailOut, tailErr, code := runIn(t, root, "checks", "tail", story, "--from", "0", "--wait", "0", "--json")
	if code != 0 {
		t.Fatalf("tail: %d %s %s", code, tailOut, tailErr)
	}
	var tailResult map[string]any
	if err := json.Unmarshal([]byte(tailOut), &tailResult); err != nil {
		t.Fatal(err)
	}
	if tailResult["running"] != false {
		t.Errorf("tail: %+v", tailResult)
	}
	if off, ok := tailResult["offset"].(float64); !ok || off <= 0 {
		t.Errorf("tail should have read something: %+v", tailResult)
	}

	// a second run is fine now that the first has ended
	if _, errOut, code := runIn(t, root, "checks", "run", story); code != 0 {
		t.Errorf("a second run once the first ended: %d %s", code, errOut)
	}
}

func TestChecksRunRefusesWithoutAWorktree(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--name", "ok", "--", "true"); code != 0 {
		t.Fatal("set")
	}
	if _, errOut, code := runIn(t, root, "checks", "run", "S-0099"); code == 0 || !strings.Contains(errOut, "no worktree") || !strings.Contains(errOut, "flai stream open S-0099") {
		t.Errorf("run without a worktree: %d %s", code, errOut)
	}
}

func TestChecksCancelRefusesWhenNothingIsRunning(t *testing.T) {
	t.Parallel()
	root, story := checksProject(t)
	if _, errOut, code := runIn(t, root, "checks", "cancel", story); code == 0 || !strings.Contains(errOut, "no checks run is active") {
		t.Errorf("cancel: %d %s", code, errOut)
	}
}

// checksJob is S-0007's checks job, ended with outcome when it is not empty:
// in its namespace s-0007, with the checks unit and e2e, e2e failing with
// exit code 3 when outcome is failed.
func checksJob(root string, started bool, outcome string) kluster.AskedJob {
	j := askedJob(kluster.KindChecks, started, outcome)
	j.Namespace = "s-0007"
	if outcome == "" {
		return j
	}
	log := filepath.Join(root, ".flai-cache", "serve", "checks", "S-0007.log")
	zero, three := 0, 3
	j.Log = log
	j.Checks = &kluster.ChecksResult{Outcome: outcome, Commit: klusterHead, Log: log, Checks: []kluster.Check{
		{Name: "unit", Outcome: kluster.OutcomePassed, ExitCode: &zero, Log: log},
		{Name: "e2e", Outcome: kluster.OutcomePassed, ExitCode: &zero, Log: log},
	}}
	if outcome == kluster.OutcomeFailed {
		j.Checks.Checks[1].Outcome, j.Checks.Checks[1].ExitCode = kluster.OutcomeFailed, &three
		j.FailedStep, j.ExitCode, j.Error = "e2e", &three, "the checks run against the namespace ended failed, at e2e"
	}
	return j
}

// queuedChecks is flai serve's answer that it queued the checks job.
func queuedChecks(job kluster.AskedJob, log string) func(kluster.Request) kluster.Answer {
	return func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: job.ID, Kind: r.Kind, Story: r.Story, Ahead: 1, Started: true, Log: log}
	}
}

// S-0661 criterion 1, ADR-0292 §1: flai checks run --job runs nothing here:
// it asks flai serve for the checks of the session's own story and prints
// the job's id, the jobs ahead, where the log goes, and how to read it later.
func TestChecksRunJobAsksFlaiServe(t *testing.T) {
	root := jobProject(t)
	running := checksJob(root, false, "")
	dir := filepath.Join(root, ".flai-cache", "serve", "checks")
	got := serveKluster(t, root, "S-0007", kluster.KindChecks, queuedChecks(running, dir), running)
	out, errOut, code := runKluster(t, root, false, "", "checks", "run", "S-7", "--job")
	if code != 0 {
		t.Fatalf("checks run --job: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{"kluster checks job " + running.ID + " queued for S-0007, 1 job ahead of it", "Wait for its kluster event (wait_for_events), or read it later with flai checks status S-0007", "its output goes to " + dir} {
		if !strings.Contains(out, want) {
			t.Errorf("checks run --job prints %q: %s", want, out)
		}
	}
	if req := <-got; req.Story != "S-0007" || req.Kind != kluster.KindChecks || req.By != "agent-S-0007" {
		t.Errorf("the request flai serve read: %+v", req)
	}
}

// S-0661 criterion 1, ADR-0292 §5: --wait prints a passing run's outcome,
// commit, namespace, and each check's name, outcome, exit code, and log, and
// exits 0; --json prints the answer, the job with its checks, and ended.
func TestChecksRunJobWaitsForAPassingRun(t *testing.T) {
	root := jobProject(t)
	running, passed := checksJob(root, true, ""), checksJob(root, true, kluster.OutcomePassed)
	log := passed.Log
	serveKluster(t, root, "S-0007", kluster.KindChecks, queuedChecks(running, ""), running, passed)
	out, errOut, code := runKluster(t, root, false, "", "checks", "run", "S-0007", "--job", "--wait", "5s")
	if code != 0 {
		t.Fatalf("a passing run exits 0: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{"S-0007's checks job " + running.ID + ": passed", "commit  0123456789ab\n", "namespace s-0007\n", "checks  passed at 0123456789ab, log " + log,
		"unit     passed, exit code 0, log " + log, "e2e      passed, exit code 0, log " + log, "log     " + log} {
		if !strings.Contains(out, want) {
			t.Errorf("--wait prints %q: %s", want, out)
		}
	}

	serveKluster(t, root, "S-0007", kluster.KindChecks, queuedChecks(running, ""), running, passed)
	out, errOut, code = runKluster(t, root, false, "", "checks", "run", "S-0007", "--job", "--wait", "5s", "--json")
	var res struct {
		Answer kluster.Answer    `json:"answer"`
		Job    *kluster.AskedJob `json:"job"`
		Ended  *bool             `json:"ended"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil {
		t.Fatalf("--job --wait --json: %d %s %s", code, out, errOut)
	}
	if res.Answer.Job != running.ID || res.Ended == nil || !*res.Ended || res.Job == nil || res.Job.Checks == nil || len(res.Job.Checks.Checks) != 2 ||
		res.Job.Checks.Checks[1].Name != "e2e" || res.Job.Checks.Checks[1].Log != log || res.Job.Namespace != "s-0007" {
		t.Errorf("--json: %s", out)
	}
}

// S-0661 criterion 1: --wait prints a failing run's failed check, its exit
// code, its log, and the error, and exits 1.
func TestChecksRunJobWaitsForAFailingRun(t *testing.T) {
	root := jobProject(t)
	running, failed := checksJob(root, true, ""), checksJob(root, true, kluster.OutcomeFailed)
	log := failed.Log
	serveKluster(t, root, "S-0007", kluster.KindChecks, queuedChecks(running, ""), running, failed)
	out, errOut, code := runKluster(t, root, false, "", "checks", "run", "S-0007", "--job", "--wait", "5s")
	for _, want := range []string{"S-0007's checks job " + running.ID + ": failed", "checks  failed at 0123456789ab, log " + log, "unit     passed, exit code 0, log " + log,
		"e2e      failed, exit code 3, log " + log, "failed at e2e, exit code 3", "error   the checks run against the namespace ended failed, at e2e"} {
		if !strings.Contains(out, want) {
			t.Errorf("--wait prints %q: %s", want, out)
		}
	}
	if code == 0 {
		t.Errorf("a failed run exits 1: %s %s", out, errOut)
	}
}

// S-0661 criterion 2: a namespace not deployed at the story branch's head is
// refused with flai serve's reason, saying to deploy first, and exits 1, its
// answer with --json; so is an answer from the record that did not pass.
func TestChecksRunJobRefusedWhenNotDeployedAtTheHead(t *testing.T) {
	root := jobProject(t)
	reason := "S-0007's namespace s-0007 was deployed at 89abcdef0123, not at the story branch's head 0123456789ab, so S-0007's checks have nothing deployed at its head to run against and nothing is queued: deploy it first with flai kluster deploy S-0007 --job"
	refuse := func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Kind: r.Kind, Story: r.Story, Reason: reason}
	}
	serveKluster(t, root, "S-0007", kluster.KindChecks, refuse)
	out, errOut, code := runKluster(t, root, false, "", "checks", "run", "S-0007", "--job", "--wait", "5s")
	if code == 0 || !strings.Contains(errOut, reason) || out != "" {
		t.Errorf("not deployed at the head: %d %q %s", code, out, errOut)
	}
	serveKluster(t, root, "S-0007", kluster.KindChecks, refuse)
	out, errOut, code = runKluster(t, root, false, "", "checks", "run", "S-0007", "--job", "--json")
	var res map[string]map[string]any
	if code == 0 || json.Unmarshal([]byte(out), &res) != nil || res["answer"]["reason"] != reason || res["answer"]["kind"] != kluster.KindChecks {
		t.Errorf("a refusal --json prints the answer and exits 1: %d %s %s", code, out, errOut)
	}

	serveKluster(t, root, "S-0007", kluster.KindChecks, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Kind: r.Kind, Story: r.Story, Recorded: true, Commit: klusterHead, Outcome: kluster.OutcomeFailed, Log: "checks/S-0007.log"}
	})
	out, errOut, code = runKluster(t, root, false, "", "checks", "run", "S-0007", "--job")
	if code == 0 || !strings.Contains(out, "S-0007's checks is answered from its namespace's record, which failed at 0123456789ab; nothing was queued") || !strings.Contains(out, "log     checks/S-0007.log") {
		t.Errorf("an answer from the record that failed: %d %s %s", code, out, errOut)
	}
}

// S-0661 criterion 1, ADR-0292 §1: a session with a story of its own is
// refused another story's checks before any request is written, and --wait
// goes with --job alone.
func TestChecksRunJobRefusesAnotherStory(t *testing.T) {
	root := jobProject(t)
	t.Setenv("FLAI_STORY", "S-0007")
	_, errOut, code := runKluster(t, root, false, "", "checks", "run", "S-0008", "--job")
	if code == 0 || !strings.Contains(errOut, "this session works S-0007") || !strings.Contains(errOut, "nothing was asked for S-0008") {
		t.Errorf("checks run --job for another story: %d %s", code, errOut)
	}
	_, errOut, code = runKluster(t, root, false, "", "checks", "run", "S-0007", "--wait", "1m")
	if code == 0 || !strings.Contains(errOut, "--wait waits for the job flai serve runs, so it goes with --job") {
		t.Errorf("--wait without --job: %d %s", code, errOut)
	}
	if _, err := os.Stat(kluster.RequestsDir(root)); !os.IsNotExist(err) {
		t.Errorf("a refusal wrote a request: %v", err)
	}
}

// S-0661 criterion 1, ADR-0292 §5: a run that outlasts --wait is said to be
// still running, with how to read it later, and exits 0.
func TestChecksRunJobWaitRunsOut(t *testing.T) {
	root := jobProject(t)
	running := checksJob(root, true, "")
	serveKluster(t, root, "S-0007", kluster.KindChecks, queuedChecks(running, ""), running)
	out, errOut, code := runKluster(t, root, false, "", "checks", "run", "S-0007", "--job", "--wait", "100ms")
	if code != 0 || !strings.Contains(out, "It is still running after 100ms: wait for its kluster event (wait_for_events), or read it later with flai checks status S-0007") {
		t.Errorf("a run that outlasts the wait: %d %s %s", code, out, errOut)
	}
}
