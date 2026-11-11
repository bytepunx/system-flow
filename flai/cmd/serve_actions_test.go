package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// inProcess runs the flai commands behind the host's write methods in this
// process: under go test the executable is the test binary, not flai.
func inProcess(t *testing.T) {
	t.Helper()
	was := hostapi.Commands
	t.Cleanup(func() { hostapi.Commands = was })
	t.Setenv("LOG_FORMAT", "json")
	hostapi.Commands = func(_ context.Context, r hostapi.Run) (hostapi.Ran, error) {
		out, errOut, code := runIn(t, r.Dir, r.Args...)
		ran := hostapi.Ran{Stdout: []byte(out), Exit: code}
		for _, line := range strings.Split(errOut, "\n") {
			var ev map[string]any
			if json.Unmarshal([]byte(line), &ev) == nil {
				ran.Events = append(ran.Events, ev)
			}
		}
		return ran, nil
	}
}

const publishRequest = `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a10"}`

// S-0078: a host action is off until the operator enables it by name on the
// host; asked for meanwhile it is refused and says what enables it; enabled,
// it runs; and every request is in the journal, refused ones too. The push
// action gates Publish, the one way accepted work reaches the remote
// (ADR-0067).
func TestHostActionPush(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	inProcess(t)
	head := func(dir string) string { return strings.TrimSpace(gitIn(t, dir, "rev-parse", "main")) }
	before := head(remote)
	// S-0087: acceptance itself never pushes or tags, whatever the push
	// action is set to; it only merges, archives, and commits locally.
	// Releasing waits for Publish, so acceptances batch into one release.
	out, _, code := runIn(t, root, "hostapi", "accept.run", `{"id":"S-0001","include_uncommitted":true,"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a09"}`)
	if code != 0 || head(remote) != before || head(root) == before {
		t.Fatalf("an acceptance with the action off: accepted here, nothing pushed: %d %s", code, out)
	}

	out, _, _ = runIn(t, root, "serve", "actions")
	if !strings.Contains(out, "push: off everywhere") || !strings.Contains(out, "off for this project; flai serve enable push") || !strings.Contains(out, "publish everything accepted") {
		t.Errorf("actions, before: %s", out)
	}
	if j, _, _ := runIn(t, root, "serve", "journal"); !strings.Contains(j, "no host action has been asked for") {
		t.Errorf("an empty journal: %s", j)
	}
	info, _, _ := runIn(t, root, "hostapi", "project.info")
	if !strings.Contains(info, `"host_actions"`) || !strings.Contains(info, `"push":false`) {
		t.Errorf("the dashboard is told the action is off: %s", info)
	}

	// off: refused, nothing pushed
	out, _, code = runIn(t, root, "hostapi", "publish.run", publishRequest)
	if code == 0 || !strings.Contains(out, "is not enabled for this project") || !strings.Contains(out, "flai serve enable push") || head(remote) != before {
		t.Fatalf("disabled: %d %s", code, out)
	}

	if _, errOut, code := runIn(t, root, "serve", "enable", "pull"); code == 0 || !strings.Contains(errOut, "there are: agent, analyze, auto-approve, auto-publish, auto-upgrade, checks, components, curate, dashboard, forge, formula, host, kluster, merge-queue, orchestrate, plan, push") {
		t.Errorf("an action there is not: %d %s", code, errOut)
	}
	out, _, code = runIn(t, root, "serve", "enable", "push")
	if code != 0 || !strings.Contains(out, "push enabled for t") || !strings.Contains(out, "flai serve journal") || !strings.Contains(out, "flai serve disable push") {
		t.Fatalf("enable: %d %s", code, out)
	}
	cfg, _ := os.ReadFile(configOf(t))
	if !strings.Contains(string(cfg), `"host_actions"`) || !strings.Contains(string(cfg), root) {
		t.Errorf("it is written to the host's configuration, for this project: %s", cfg)
	}
	// another project on the same host is not enabled by it
	other := tempProject(t)
	if info, _, _ := runIn(t, other, "hostapi", "project.info"); !strings.Contains(info, `"push":false`) {
		t.Errorf("per project: %s", info)
	}

	// on: published, the acceptance tagged and pushed with its tag
	out, _, code = runIn(t, root, "hostapi", "publish.run", publishRequest)
	if code != 0 || !strings.Contains(out, `"pushed":true`) || !strings.Contains(out, "cli/v1.1.0") || head(remote) != head(root) {
		t.Fatalf("enabled: %d %s", code, out)
	}
	if tags := gitIn(t, remote, "tag", "--list"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the release tag reached the remote: %s", tags)
	}
	out, _, _ = runIn(t, root, "hostapi", "publish.run", `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a11"}`)
	if !strings.Contains(out, "nothing pending") {
		t.Errorf("nothing left: %s", out)
	}

	// the journal has all three, in order, and says for whom
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	var entries []struct{ Action, Method, Outcome, Detail, By, Root, At string }
	if err := json.Unmarshal([]byte(js), &entries); err != nil || len(entries) != 3 {
		t.Fatalf("journal: %v %s", err, js)
	}
	for i, want := range []struct{ outcome, detail string }{{"disabled", ""}, {"done", "pushed with tags cli/v1.1.0"}, {"done", "nothing pushed: nothing pending"}} {
		e := entries[i]
		if e.Outcome != want.outcome || e.Detail != want.detail || e.Action != "push" || e.Method != "publish.run" || e.Root != root || e.By == "" || e.At == "" {
			t.Errorf("entry %d: %+v", i, e)
		}
	}
	text, _, _ := runIn(t, root, "serve", "journal", "-n", "1")
	if !strings.Contains(text, "nothing pushed") || strings.Contains(text, "disabled") {
		t.Errorf("journal -n 1: %s", text)
	}
	if st, err := os.Stat(filepath.Join(filepath.Dir(configOf(t)), "serve", "journal.jsonl")); err != nil || st.Mode().Perm() != 0o600 {
		t.Errorf("the journal is the operator's alone: %v %v", st, err)
	}

	// off again, at once
	if out, _, _ := runIn(t, root, "serve", "disable", "push"); !strings.Contains(out, "push disabled for t") {
		t.Errorf("disable: %s", out)
	}
	if out, _, code := runIn(t, root, "hostapi", "publish.run", `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a12"}`); code == 0 || !strings.Contains(out, "not enabled") {
		t.Errorf("disabled again: %d %s", code, out)
	}
	// every project, and one project's disable under it says so
	runIn(t, root, "serve", "enable", "push", "--all-projects")
	if out, _, _ := runIn(t, root, "serve", "disable", "push"); !strings.Contains(out, "still enabled") || !strings.Contains(out, "--all-projects") {
		t.Errorf("disable under all projects: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "disable", "push", "--all-projects"); !strings.Contains(out, "push disabled for every project") {
		t.Errorf("disable everywhere: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "push: off everywhere") {
		t.Errorf("actions, after: %s", out)
	}
}

// S-0079: the agent's command is the operator's, an argument list with no
// default, managed on the host and nowhere else.
func TestServeAgentCommand(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "agent")
	if !strings.Contains(out, "no command is set, so a story whose agent and project name no harness starts on flai's built-in default, claude-code") || !strings.Contains(out, "off for this project") {
		t.Errorf("by default: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "agent: off everywhere") || !strings.Contains(out, "whoever can move a story to ready") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "set", "--name", "builder", "--", "claude", "-p", "work on {story}; echo $HOME")
	if code != 0 || !strings.Contains(out, `command: "claude" "-p" "work on {story}; echo $HOME"`) || !strings.Contains(out, "never through a shell") {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	cfg, _ := os.ReadFile(configOf(t))
	if !strings.Contains(string(cfg), `"work on {story}; echo $HOME"`) || !strings.Contains(string(cfg), `"name": "builder"`) {
		t.Errorf("kept as an argument list, as written: %s", cfg)
	}
	if _, errOut, code := runIn(t, root, "config", "set", "agent.command", "rm -rf /"); code == 0 {
		t.Errorf("flai config set does not reach it: %s", errOut)
	}
	if _, _, code := runIn(t, root, "serve", "agent", "set"); code == 0 {
		t.Error("set needs a program")
	}
	runIn(t, root, "serve", "enable", "agent")
	js, _, _ := runIn(t, root, "serve", "agent", "show", "--json")
	if !strings.Contains(js, `"enabled_here": true`) || !strings.Contains(js, `"claude"`) {
		t.Errorf("show --json: %s", js)
	}
	got := (&app{configPath: configOf(t)}).agentConfig(root)
	if !got.Enabled || len(got.Command) != 3 || got.Name != "builder" {
		t.Errorf("what flai serve is given: %+v", got)
	}
	if out, _, _ := runIn(t, root, "serve", "agent", "clear"); !strings.Contains(out, "no command is set") {
		t.Errorf("clear: %s", out)
	}
	if got := (&app{configPath: configOf(t)}).agentConfig(root); len(got.Command) != 0 {
		t.Errorf("cleared: %+v", got)
	}
}

// S-0218: orchestrate is off until the operator enables it, flai serve is
// told at every look whether it is on, and agent.status carries the
// project's orchestrator run, null before the first.
func TestServeOrchestrateAction(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "orchestrate: off everywhere") || !strings.Contains(out, "orchestration.permissions") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{configPath: configOf(t)}).agentConfig(root).Orchestrate {
		t.Error("orchestrate is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{configPath: configOf(t)}).host().Agent(context.Background(), root, true))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"orchestrator":null`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "orchestrate"); code != 0 || !strings.Contains(out, "orchestrate enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{configPath: configOf(t)}).agentConfig(root).Orchestrate {
		t.Error("enabled, flai serve is not told")
	}
	run := `{"` + root + `":{"orchestrator":{"story":"","harness":"command","command":"run-agent","agent":"orchestrator","pid":42,"started":"2026-10-05T10:00:00Z","session":"s","log":"/l/t-orchestrator-20261005T100000Z.log"}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := status(); !strings.Contains(got, `"orchestrator":{"story":"","harness":"command","command":"run-agent","agent":"orchestrator","pid":42,"started":"2026-10-05T10:00:00Z"`) || !strings.Contains(got, `"running":null`) {
		t.Errorf("agent.status with a run: %s", got)
	}
	runIn(t, root, "serve", "disable", "orchestrate")
	if (&app{configPath: configOf(t)}).agentConfig(root).Orchestrate {
		t.Error("disabled, flai serve is still told it is on")
	}
}

// S-0603: agent.status records agent.state on the request's recorder, and
// inside it the three parts of the state the host builds: reading the agent
// states, the stories' activity, and the roles; and inside the activity its
// own parts (I-0313).
func TestAgentStatusRecordsItsPhases(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	// each look at the clock is a millisecond on, so a phase that holds
	// another is longer than it
	var ticks atomic.Int64
	at := time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)
	ctx, rec := perf.StartAt(context.Background(), func() time.Time {
		return at.Add(time.Duration(ticks.Add(1)) * time.Millisecond)
	})
	host := (&app{configPath: configOf(t), errOut: io.Discard}).host()
	res, e := hostapi.MethodsFor("test", nil, host)["agent.status"](ctx, channel.Project{Key: "t", Root: root}, json.RawMessage(`{}`))
	if e != nil {
		t.Fatal(e)
	}
	if out, _ := res.(map[string]any); out["state"] == nil {
		t.Fatalf("agent.status carries no state: %v", res)
	}
	got := map[string]perf.Phase{}
	for _, ph := range rec.Phases() {
		got[ph.Name] = ph
	}
	outer, ok := got["agent.state"]
	if !ok || outer.Count != 1 {
		t.Fatalf("agent.state not recorded once: %v", rec)
	}
	var inner time.Duration
	for _, name := range []string{"agent.states", "agent.activity", "agent.roles"} {
		ph, ok := got[name]
		if !ok || ph.Count != 1 {
			t.Errorf("%s not recorded once: %v", name, rec)
			continue
		}
		inner += ph.Total
	}
	if inner >= outer.Total {
		t.Errorf("the three phases (%v) are not inside agent.state (%v): %v", inner, outer.Total, rec)
	}
	// I-0313: inside agent.activity, the runs judged, the holds built, and
	// the stories begun elsewhere; the trial judge's phase, which needs a
	// story's branch to trial-merge, is shown in serve's
	// TestActivityRecordsItsPhases
	var parts time.Duration
	for _, name := range []string{"activity.conversations", "activity.runs", "activity.holds", "activity.elsewhere"} {
		ph, ok := got[name]
		if !ok || ph.Count != 1 {
			t.Errorf("%s not recorded once: %v", name, rec)
			continue
		}
		parts += ph.Total
	}
	if activity := got["agent.activity"]; parts >= activity.Total {
		t.Errorf("activity's phases (%v) are not inside agent.activity (%v): %v", parts, activity.Total, rec)
	}
}

// S-0534, ADR-0209 §5: agent.status's state carries the host's budget state,
// one entry per provider held, read from budget.json in the serve folder,
// and leaves it out once the state is cleared or the file cannot be read.
func TestAgentStatusCarriesTheHostsBudget(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	write := func(state string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "budget.json"), []byte(state), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	var warned bytes.Buffer
	status := func() map[string]any {
		t.Helper()
		out, ok := (&app{configPath: configOf(t), errOut: &warned}).host().Agent(context.Background(), root, true).(map[string]any)
		if !ok {
			t.Fatal("agent.status's state is not a map")
		}
		return out
	}
	if got, ok := status()["budget"]; ok {
		t.Errorf("budget with no budget.json: %v", got)
	}
	write(`{"codex":{"provider":"codex","since":"2026-10-28T07:55:00Z","until":"2026-10-28T08:00:00Z","backoff":"5m0s",
"probe":{"root":"/r","work":"S-0001","started":"2026-10-28T08:01:00Z","claimed":"2026-10-28T08:01:00Z"}},
"claude-code":{"provider":"claude-code","limit":"usage_limit_reached","status":429,"since":"2026-10-28T06:30:00+01:00","until":"2026-10-28T09:00:00Z","reset":"2026-10-28T09:00:00Z","backoff":"5m0s"}}`)
	got, err := json.Marshal(status()["budget"])
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"provider":"claude-code","limit":"usage_limit_reached","status":429,"since":"2026-10-28T05:30:00Z","until":"2026-10-28T09:00:00Z"},` +
		`{"provider":"codex","since":"2026-10-28T07:55:00Z","until":"2026-10-28T08:00:00Z","probing":true}]`
	if string(got) != want {
		t.Errorf("budget while held:\n got %s\nwant %s", got, want)
	}
	write(`{}`)
	if got, ok := status()["budget"]; ok {
		t.Errorf("budget once cleared: %v", got)
	}
	write(`not json`)
	if st := status(); st["budget"] != nil || st["command"] == nil {
		t.Errorf("an unreadable budget.json is none, and the rest of the state stands: %v", st)
	}
	if !strings.Contains(warned.String(), "budget state not read") {
		t.Errorf("an unreadable budget.json is not warned of: %s", warned.String())
	}
}

// S-0534: flai serve status names each held provider with its limit and
// until, says when a probe runs, and when until has passed with none.
func TestBudgetLine(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 10, 28, 8, 30, 0, 0, time.UTC)
	for _, c := range []struct {
		b    hostapi.Budget
		want string
	}{
		{hostapi.Budget{Provider: "claude-code", Limit: "usage_limit_reached", Status: 429, Until: "2026-10-28T09:00:00Z"},
			"budget: claude-code held (usage_limit_reached) until 2026-10-28T09:00:00Z"},
		{hostapi.Budget{Provider: "codex", Status: 429, Until: "2026-10-28T09:00:00Z"},
			"budget: codex held (HTTP 429) until 2026-10-28T09:00:00Z"},
		{hostapi.Budget{Provider: "codex", Until: "2026-10-28T08:00:00Z", Probing: true},
			"budget: codex held until 2026-10-28T08:00:00Z, probing"},
		{hostapi.Budget{Provider: "codex", Until: "2026-10-28T08:00:00Z"},
			"budget: codex held until 2026-10-28T08:00:00Z, passed; the next start probes"},
	} {
		if got := budgetLine(c.b, now); got != c.want {
			t.Errorf("budgetLine(%+v):\n got %s\nwant %s", c.b, got, c.want)
		}
	}
}

// S-0399, ADR-0173: merge-queue is off until the operator enables it, flai
// serve actions names it and what it means, the dashboard's settings page
// lists it, and enable and disable turn it on and off for the project.
func TestServeMergeQueueAction(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	enabled := func() bool {
		t.Helper()
		cfg, _, err := (&app{configPath: configOf(t)}).loadConfig()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.ActionEnabled(hostapi.ActionMergeQueue, root)
	}
	listed := func() map[string]any {
		t.Helper()
		got, err := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), root))
		if err != nil {
			t.Fatal(err)
		}
		var settings struct {
			Actions []map[string]any `json:"actions"`
		}
		if err := json.Unmarshal(got, &settings); err != nil {
			t.Fatal(err)
		}
		for _, a := range settings.Actions {
			if a["name"] == hostapi.ActionMergeQueue {
				return a
			}
		}
		return nil
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "merge-queue: off everywhere") || !strings.Contains(out, "merge_queue.batch_max") || !strings.Contains(out, "orchestration.permissions.accept_reviews") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if enabled() {
		t.Error("merge-queue is on by default")
	}
	if a := listed(); a == nil || a["here"] != false || a["everywhere"] != false {
		t.Errorf("the settings page lists it, off: %v", a)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "merge-queue"); code != 0 || !strings.Contains(out, "merge-queue enabled for t") || !strings.Contains(out, "flai serve disable merge-queue") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !enabled() {
		t.Error("enabled, the host's configuration does not say so")
	}
	if a := listed(); a == nil || a["here"] != true {
		t.Errorf("the settings page lists it, on: %v", a)
	}
	if out, errOut, code := runIn(t, root, "serve", "disable", "merge-queue"); code != 0 || !strings.Contains(out, "merge-queue disabled for t") {
		t.Fatalf("disable: %d %s %s", code, out, errOut)
	}
	if enabled() {
		t.Error("disabled, the host's configuration still has it on")
	}
}

// S-0422: formula is off until the operator enables it, flai serve actions
// names it and what it means, the dashboard's settings page lists it, and
// enable and disable turn it on and off for the project.
func TestServeFormulaAction(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	enabled := func() bool {
		t.Helper()
		cfg, _, err := (&app{configPath: configOf(t)}).loadConfig()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.ActionEnabled(hostapi.ActionFormula, root)
	}
	listed := func() map[string]any {
		t.Helper()
		got, err := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), root))
		if err != nil {
			t.Fatal(err)
		}
		var settings struct {
			Actions []map[string]any `json:"actions"`
		}
		if err := json.Unmarshal(got, &settings); err != nil {
			t.Fatal(err)
		}
		for _, a := range settings.Actions {
			if a["name"] == hostapi.ActionFormula {
				return a
			}
		}
		return nil
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "formula: off everywhere") || !strings.Contains(out, "flai formula run <name> --ready --autocommit") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if enabled() {
		t.Error("formula is on by default")
	}
	if a := listed(); a == nil || a["here"] != false || a["everywhere"] != false {
		t.Errorf("the settings page lists it, off: %v", a)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "formula"); code != 0 || !strings.Contains(out, "formula enabled for t") || !strings.Contains(out, "flai serve disable formula") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !enabled() {
		t.Error("enabled, the host's configuration does not say so")
	}
	if a := listed(); a == nil || a["here"] != true {
		t.Errorf("the settings page lists it, on: %v", a)
	}
	if out, errOut, code := runIn(t, root, "serve", "disable", "formula"); code != 0 || !strings.Contains(out, "formula disabled for t") {
		t.Fatalf("disable: %d %s %s", code, out, errOut)
	}
	if enabled() {
		t.Error("disabled, the host's configuration still has it on")
	}
}

// S-0577: kluster is off until the operator enables it, and the dashboard's
// settings page lists it; S-0586: enabling it for a project with no kluster:
// block writes one rather than refusing, while disable and --all-projects
// write nothing (ADR-0254).
func TestServeKlusterAction(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	bare := tempProject(t)
	enabled := func(root string) bool {
		t.Helper()
		cfg, _, err := (&app{configPath: configOf(t)}).loadConfig()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.ActionEnabled(hostapi.ActionKluster, root)
	}
	if out, _, _ := runIn(t, bare, "serve", "actions"); !strings.Contains(out, "kluster: off everywhere") {
		t.Errorf("actions names it, off: %s", out)
	}
	before, _ := os.ReadFile(filepath.Join(bare, "system-flow.yaml"))
	if out, errOut, code := runIn(t, bare, "serve", "disable", "kluster"); code != 0 || !strings.Contains(out, "kluster disabled for t") {
		t.Errorf("disable without a block: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, bare, "serve", "enable", "kluster", "--all-projects"); code != 0 || !strings.Contains(out, "kluster enabled for every project") {
		t.Errorf("enable --all-projects: %d %s %s", code, out, errOut)
	}
	if after, _ := os.ReadFile(filepath.Join(bare, "system-flow.yaml")); string(after) != string(before) {
		t.Errorf("disable and --all-projects wrote the manifest:\n%s", after)
	}
	if _, err := os.Stat(filepath.Join(bare, "kluster.yaml")); err == nil {
		t.Error("--all-projects wrote kluster.yaml")
	}
	if out, errOut, code := runIn(t, bare, "serve", "disable", "kluster", "--all-projects"); code != 0 {
		t.Errorf("disable --all-projects: %d %s %s", code, out, errOut)
	}

	root := tempProject(t)
	m := filepath.Join(root, "system-flow.yaml")
	b, _ := os.ReadFile(m)
	withBlock := string(b) + "kluster:\n  config: integration/kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n"
	writeFiles(t, root, map[string]string{
		"system-flow.yaml":         withBlock,
		"integration/kluster.yaml": "name: demo\n",
	})
	if out, errOut, code := runIn(t, root, "serve", "enable", "kluster"); code != 0 || !strings.Contains(out, "kluster enabled for t") || !strings.Contains(out, "flai serve disable kluster") || strings.Contains(out, "wrote") {
		t.Fatalf("enable with a block: %d %s %s", code, out, errOut)
	}
	if !enabled(root) {
		t.Error("enabled, the host's configuration does not say so")
	}
	if got, _ := os.ReadFile(m); string(got) != withBlock {
		t.Errorf("enable with a block changed the manifest:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(root, "kluster.yaml")); err == nil {
		t.Error("enable with a block wrote kluster.yaml")
	}
}

// klusterGitProject is a tempProject committed in a git repository of its
// own, with extra written before the first commit.
func klusterGitProject(t *testing.T, extra map[string]string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := tempProject(t)
	writeFiles(t, root, extra)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "start")
	return root
}

// S-0586: enabling kluster for a project with no kluster: block writes a
// block naming only config: kluster.yaml at the end of system-flow.yaml,
// every other line as it was, and kluster.yaml naming the cluster, k3d, and
// spire; commits those two files alone, with the trailer; turns the action
// on; and leaves flai check with no manifest.kluster finding. Enabled again,
// it writes nothing.
func TestServeEnableKlusterWritesTheBlock(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := klusterGitProject(t, nil)
	before, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	_ = os.WriteFile(filepath.Join(root, "notes.txt"), []byte("not mine\n"), 0o644)

	out, errOut, code := runIn(t, root, "serve", "enable", "kluster", "--trailer", "Co-Authored-By: someone <s@s>")
	if code != 0 || !strings.Contains(out, "wrote a kluster: block naming kluster.yaml") || !strings.Contains(out, "committed as") || !strings.Contains(out, "kluster enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	after, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if want := string(before) + "kluster:\n  config: kluster.yaml\n"; string(after) != want {
		t.Errorf("manifest:\n%s\nwant:\n%s", after, want)
	}
	cfg, _ := os.ReadFile(filepath.Join(root, "kluster.yaml"))
	for _, want := range []string{"\nname: t-local\n", "\nprovider: k3d\n", "\nprofile: spire\n"} {
		if !strings.Contains(string(cfg), want) {
			t.Errorf("kluster.yaml lacks %q:\n%s", want, cfg)
		}
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); files != "kluster.yaml\nsystem-flow.yaml" {
		t.Errorf("the commit holds the two files alone: %q", files)
	}
	if msg := gitIn(t, root, "log", "-1", "--format=%B"); !strings.Contains(msg, "kluster: block") || !strings.Contains(msg, "Co-Authored-By: someone <s@s>") {
		t.Errorf("commit message: %q", msg)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "?? notes.txt" {
		t.Errorf("status after the commit: %q", st)
	}
	conf, _, err := (&app{configPath: configOf(t)}).loadConfig()
	if err != nil || !conf.ActionEnabled(hostapi.ActionKluster, root) {
		t.Errorf("the action is on: %v", err)
	}
	if out, _, _ := runIn(t, root, "check", "--strict"); strings.Contains(out, "manifest.kluster") {
		t.Errorf("flai check finds the scaffold faulty: %s", out)
	}

	head := gitIn(t, root, "rev-parse", "HEAD")
	if out, errOut, code := runIn(t, root, "serve", "enable", "kluster"); code != 0 || strings.Contains(out, "wrote") {
		t.Errorf("enabled again: %d %s %s", code, out, errOut)
	}
	if again := gitIn(t, root, "rev-parse", "HEAD"); again != head {
		t.Error("enabled again, it committed")
	}
	if got, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml")); string(got) != string(after) {
		t.Errorf("enabled again, it wrote the manifest:\n%s", got)
	}
}

// S-0586: an existing kluster.yaml is kept and only the manifest committed,
// and dashboard.autocommit: false leaves what enabling wrote uncommitted;
// --json says what it wrote.
func TestServeEnableKlusterKeepsTheConfigAndHonoursAutocommit(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	kept := klusterGitProject(t, map[string]string{"kluster.yaml": "name: mine\nprofile: authstar\n"})
	out, errOut, code := runIn(t, kept, "serve", "enable", "kluster", "--json")
	if code != 0 {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	var got struct {
		Enabled bool `json:"enabled"`
		Kluster struct {
			Wrote  []string `json:"wrote"`
			Kept   bool     `json:"kept"`
			Commit string   `json:"commit"`
		} `json:"kluster"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if !got.Enabled || !got.Kluster.Kept || len(got.Kluster.Wrote) != 1 || got.Kluster.Wrote[0] != "system-flow.yaml" || got.Kluster.Commit == "" {
		t.Errorf("--json: %+v", got)
	}
	if cfg, _ := os.ReadFile(filepath.Join(kept, "kluster.yaml")); string(cfg) != "name: mine\nprofile: authstar\n" {
		t.Errorf("the existing kluster.yaml was changed:\n%s", cfg)
	}
	if files := gitIn(t, kept, "show", "--name-only", "--format=", "HEAD"); files != "system-flow.yaml" {
		t.Errorf("the commit holds the manifest alone: %q", files)
	}

	off := klusterGitProject(t, nil)
	b, _ := os.ReadFile(filepath.Join(off, "system-flow.yaml"))
	writeFiles(t, off, map[string]string{"system-flow.yaml": string(b) + "dashboard:\n  autocommit: false\n"})
	gitIn(t, off, "commit", "-q", "-am", "no autocommit")
	head := gitIn(t, off, "rev-parse", "HEAD")
	out, errOut, code = runIn(t, off, "serve", "enable", "kluster")
	if code != 0 || !strings.Contains(out, "not committed") {
		t.Fatalf("enable with autocommit off: %d %s %s", code, out, errOut)
	}
	if gitIn(t, off, "rev-parse", "HEAD") != head {
		t.Error("dashboard.autocommit: false, it committed")
	}
	if st := gitIn(t, off, "status", "--porcelain"); st != "M system-flow.yaml\n?? kluster.yaml" {
		t.Errorf("status: %q", st)
	}
}

// S-0223: analyze is off until the operator enables it, flai serve is told
// whether it is on, and agent.status carries the project's newest analyzer
// run, with its focus, trigger, outcome, report, and log, null before the
// first, apart from the orchestrator's.
func TestServeAnalyzeAction(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "analyze: off everywhere") || !strings.Contains(out, "analysis.agent") || !strings.Contains(out, "design/analysis") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{configPath: configOf(t)}).agentConfig(root).Analyze {
		t.Error("analyze is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{configPath: configOf(t)}).host().Agent(context.Background(), root, true))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"analyzer":null`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "analyze"); code != 0 || !strings.Contains(out, "analyze enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{configPath: configOf(t)}).agentConfig(root).Analyze {
		t.Error("enabled, flai serve is not told")
	}
	run := `{"` + root + `":{"analyzer":{"story":"","harness":"command","command":"run-agent","agent":"analyzer","pid":42,"started":"2026-10-06T10:00:00Z","ended":"2026-10-06T10:05:00Z","log":"/l/t-analyzer-20261006T100000Z.log","trigger":"asked","focus":"risk","report":"design/analysis/2026-10-06-risk.md","outcome":"worked"}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	got := status()
	for _, want := range []string{`"analyzer":{"story":"","harness":"command","command":"run-agent","agent":"analyzer","pid":42,"started":"2026-10-06T10:00:00Z","ended":"2026-10-06T10:05:00Z"`,
		`"log":"/l/t-analyzer-20261006T100000Z.log"`, `"trigger":"asked","focus":"risk","report":"design/analysis/2026-10-06-risk.md","outcome":"worked"`, `"orchestrator":null`, `"running":null`} {
		if !strings.Contains(got, want) {
			t.Errorf("agent.status with a run lacks %s: %s", want, got)
		}
	}
	runIn(t, root, "serve", "disable", "analyze")
	if (&app{configPath: configOf(t)}).agentConfig(root).Analyze {
		t.Error("disabled, flai serve is still told it is on")
	}
}

// S-0431: curate is off until the operator enables it, flai serve is told
// whether it is on, and agent.status carries the project's newest curator
// run, with its since, trigger, outcome, entries, and log, null before the
// first, apart from the analyzer's.
func TestServeCurateAction(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "curate: off everywhere") || !strings.Contains(out, "digest.agent") || !strings.Contains(out, "flai digest new and edit") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{configPath: configOf(t)}).agentConfig(root).Curate {
		t.Error("curate is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{configPath: configOf(t)}).host().Agent(context.Background(), root, true))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"curator":null`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "curate"); code != 0 || !strings.Contains(out, "curate enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{configPath: configOf(t)}).agentConfig(root).Curate {
		t.Error("enabled, flai serve is not told")
	}
	run := `{"` + root + `":{"curator":{"story":"","harness":"command","command":"run-agent","agent":"curator","role":"curate","pid":42,"started":"2026-10-23T10:00:00Z","ended":"2026-10-23T10:05:00Z","log":"/l/t-curator-20261023T100000Z.log","trigger":"asked","since":"2026-10-01","entries":["D-0007"],"outcome":"worked"}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	got := status()
	for _, want := range []string{`"curator":{"story":"","harness":"command","command":"run-agent","agent":"curator"`, `"pid":42`, `"log":"/l/t-curator-20261023T100000Z.log"`,
		`"trigger":"asked"`, `"since":"2026-10-01"`, `"entries":["D-0007"]`, `"outcome":"worked"`, `"analyzer":null`, `"orchestrator":null`, `"running":null`} {
		if !strings.Contains(got, want) {
			t.Errorf("agent.status with a run lacks %s: %s", want, got)
		}
	}
	runIn(t, root, "serve", "disable", "curate")
	if (&app{configPath: configOf(t)}).agentConfig(root).Curate {
		t.Error("disabled, flai serve is still told it is on")
	}
}

// S-0104: what each harness a story may name is on this host, and what its
// agent may do, are the operator's, set here and nowhere else.
func TestServeAgentHarness(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "agent", "harness")
	if !strings.Contains(out, `harness claude-code: "claude" "--permission-mode" "acceptEdits" "--allowedTools" "Bash,mcp__flai"`) {
		t.Errorf("the defaults: %s", out)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude", "--", "--permission-mode", "bypassPermissions")
	if code != 0 || !strings.Contains(out, `harness claude-code: "/opt/claude" "--permission-mode" "bypassPermissions"`) {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	got := (&app{configPath: configOf(t)}).agentConfig(root).Harnesses["claude-code"]
	if got.Program != "/opt/claude" || strings.Join(got.Args, " ") != "--permission-mode bypassPermissions" {
		t.Errorf("what flai serve is given: %+v", got)
	}
	// the program alone keeps the arguments; "--" alone means none
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/usr/bin/claude")
	if got := (&app{configPath: configOf(t)}).agentConfig(root).Harnesses["claude-code"]; got.Program != "/usr/bin/claude" || len(got.Args) != 2 {
		t.Errorf("program only: %+v", got)
	}
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--")
	if got := (&app{configPath: configOf(t)}).agentConfig(root).Harnesses["claude-code"]; len(got.Args) != 0 {
		t.Errorf("no arguments: %+v", got)
	}
	// setting the command keeps the harnesses, and it is a harness of its own
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	cfg := (&app{configPath: configOf(t)}).agentConfig(root)
	if cfg.Harnesses["command"].Program != "run-agent" || cfg.Harnesses["claude-code"].Program != "/usr/bin/claude" {
		t.Errorf("both: %+v", cfg.Harnesses)
	}
	runIn(t, root, "serve", "agent", "clear")
	if got := (&app{configPath: configOf(t)}).agentConfig(root).Harnesses["claude-code"]; got.Program != "/usr/bin/claude" {
		t.Errorf("clear removes only the command: %+v", got)
	}
	out, _, _ = runIn(t, root, "serve", "agent", "harness", "claude-code", "--reset")
	if !strings.Contains(out, `harness claude-code: "claude" "--permission-mode" "acceptEdits"`) {
		t.Errorf("reset: %s", out)
	}
	for _, bad := range [][]string{{"command", "--program", "x"}, {"cursor"}, {"--program", "x"}, {"claude-code", "--reset", "--program", "x"}} {
		if _, _, code := runIn(t, root, append([]string{"serve", "agent", "harness"}, bad...)...); code == 0 {
			t.Errorf("%q was taken", bad)
		}
	}
}

// S-0508: flai serve agent harness container keeps the image per inner
// harness, the network, and the login folder, ~/ resolved, in the host's
// configuration; flai serve agent harness and flai serve settings show them,
// the default image filled in; flai serve is given them; and the next --image
// keeps the images already set. Not parallel: it sets HOME for ~/.
func TestServeAgentHarnessContainer(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	login := filepath.Join(home, ".claude")
	if err := os.Mkdir(login, 0o700); err != nil {
		t.Fatal(err)
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "agent", "harness")
	if !strings.Contains(out, "harness container: \"docker\"\n  image claude-code: flai-agent-claude-code:local (default, built by make agent-images)\n  network: Docker's default bridge, the host reached through Docker's host gateway\n  login: not set; set the folder holding your Claude login with --login ~/.claude\n") {
		t.Errorf("the defaults: %s", out)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "harness", "container", "--image", "claude-code=flai-agent-claude-code:local", "--login", "~/.claude")
	if code != 0 || !strings.Contains(out, "  image claude-code: flai-agent-claude-code:local (set)\n") || !strings.Contains(out, "  login: "+login+", mounted read-only\n") {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	cfg, _, err := (&app{configPath: configOf(t)}).loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	want := config.HarnessHost{Images: map[string]string{"claude-code": "flai-agent-claude-code:local"}, Login: login}
	if got := cfg.Agent.Project(root).Harnesses[config.ContainerHarness]; !reflect.DeepEqual(got, want) {
		t.Errorf("kept in the host configuration, as the project's: %+v", got)
	}
	runIn(t, root, "serve", "agent", "harness", "container", "--network", "flai-agents")
	if got := (&app{configPath: configOf(t)}).agentConfig(root).Harnesses[config.ContainerHarness]; got.Network != "flai-agents" || got.Login != login || got.Images["claude-code"] != "flai-agent-claude-code:local" {
		t.Errorf("what flai serve is given, the images kept: %+v", got)
	}
	out, _, _ = runIn(t, root, "serve", "settings")
	if !strings.Contains(out, "harness container: docker") || !strings.Contains(out, "    image claude-code: flai-agent-claude-code:local\n") || !strings.Contains(out, "    network flai-agents, login "+login+"\n") {
		t.Errorf("flai serve settings: %s", out)
	}
	js, errOut, code := runIn(t, root, "serve", "settings", "--json")
	var settings struct {
		Agent struct {
			Harnesses map[string]struct {
				Images  map[string]string `json:"images"`
				Network string            `json:"network"`
				Login   string            `json:"login"`
			} `json:"harnesses"`
		} `json:"agent"`
	}
	if err := json.Unmarshal([]byte(js), &settings); err != nil || code != 0 {
		t.Fatalf("flai serve settings --json: %d %v %s %s", code, err, js, errOut)
	}
	if got := settings.Agent.Harnesses[config.ContainerHarness]; got.Images["claude-code"] != "flai-agent-claude-code:local" || got.Network != "flai-agents" || got.Login != login {
		t.Errorf("flai serve settings --json: %+v", got)
	}
	js, _, _ = runIn(t, root, "serve", "agent", "show", "--json")
	if err := json.Unmarshal([]byte(js), &settings.Agent); err != nil {
		t.Fatal(err)
	}
	if got := settings.Agent.Harnesses[config.ContainerHarness]; got.Network != "flai-agents" || got.Login != login || len(got.Images) != 1 {
		t.Errorf("flai serve agent show --json: %s", js)
	}
	runIn(t, root, "serve", "agent", "harness", "container", "--network", "", "--login", "")
	if got := (&app{configPath: configOf(t)}).agentConfig(root).Harnesses[config.ContainerHarness]; got.Network != "" || got.Login != "" || len(got.Images) != 1 {
		t.Errorf(`--network "" and --login "" clear only themselves: %+v`, got)
	}
	runIn(t, root, "serve", "agent", "harness", "container", "--reset")
	if cfg, _, _ := (&app{configPath: configOf(t)}).loadConfig(); cfg.Agent.Harnesses != nil || cfg.Agent.Projects != nil {
		t.Errorf("--reset removes the entry: %+v %+v", cfg.Agent.Harnesses, cfg.Agent.Projects)
	}
}

// S-0508: --image, --network, and --login are refused on another harness, a
// malformed or empty --image, one for an inner harness flai does not run in
// the container, a --login that is not a folder, and with --reset, each saying
// what to do and changing nothing.
func TestServeAgentHarnessContainerRefusals(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"claude-code", "--network", "flai-agents"}, "--image, --network, and --login are the container harness's, which claude-code is not"},
		{[]string{"command", "--login", "/tmp"}, "which command is not"},
		{[]string{"--image", "claude-code=x"}, "name the harness"},
		{[]string{"container", "--image", "claude-code"}, "is not <inner>=<image>: give the inner harness and its image"},
		{[]string{"container", "--image", "=x"}, "=x\\\" is not <inner>=<image>"},
		{[]string{"container", "--image", "claude-code="}, "agent.harnesses.container.images.claude-code is empty"},
		{[]string{"container", "--image", "codex=x:1"}, "inside the container, only claude-code"},
		{[]string{"container", "--login", filepath.Join(t.TempDir(), "missing")}, "is not a folder; give the folder holding your Claude login"},
		{[]string{"container", "--reset", "--network", "n"}, "--reset takes neither"},
	} {
		if _, errOut, code := runIn(t, root, append([]string{"serve", "agent", "harness"}, tc.args...)...); code == 0 || !strings.Contains(errOut, tc.want) {
			t.Errorf("%q: %d %s, want %q", tc.args, code, errOut, tc.want)
		}
	}
	if cfg, _, _ := (&app{configPath: configOf(t)}).loadConfig(); cfg.Agent.Harnesses != nil || cfg.Agent.Projects != nil {
		t.Errorf("a refusal changes nothing: %+v %+v", cfg.Agent.Harnesses, cfg.Agent.Projects)
	}
}

// S-0508: config names the container harness and its inner harnesses
// itself, since it cannot import harness; they are harness's.
func TestConfigContainerNamesAreTheHarnesss(t *testing.T) {
	t.Parallel()
	if config.ContainerHarness != harness.Container {
		t.Errorf("config.ContainerHarness %q, harness.Container %q", config.ContainerHarness, harness.Container)
	}
	inners := harness.ContainerInners()
	if got := config.ContainerInners(); !slices.Equal(got, inners) {
		t.Errorf("config.ContainerInners %v, harness's %v", got, inners)
	}
	for _, inner := range inners {
		if got := config.DefaultImage(inner); got != harness.DefaultImages[inner] {
			t.Errorf("default image of %s: config %q, harness %q", inner, got, harness.DefaultImages[inner])
		}
	}
}

// S-0355 (ADR-0131): a harness's guard none and deny_protected are set and
// cleared on the host, the command's among them, each flag changing only its
// own field; --reset clears both; the command's program stays flai serve agent
// set's; and a guard other than none is refused, naming what sets it.
func TestServeAgentHarnessGuardAndDenyProtected(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	harnesses := func() map[string]config.HarnessHost {
		cfg, _, err := (&app{configPath: configOf(t)}).loadConfig()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Agent.ForProject(root).Harnesses
	}
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	out, errOut, code := runIn(t, root, "serve", "agent", "harness", "command", "--guard", "none")
	if code != 0 || !strings.Contains(out, "harness command: the command above (this project)\n  guard: none, so a story's sub-agents run without flai guard") {
		t.Fatalf("guard none on the command: %d %s %s", code, out, errOut)
	}
	out, _, _ = runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected")
	if got := harnesses()["command"]; !reflect.DeepEqual(got, config.HarnessHost{Guard: config.GuardNone, DenyProtected: true}) {
		t.Errorf("--deny-protected keeps the guard: %+v", got)
	}
	if !strings.Contains(out, "  deny_protected: your arguments deny the protected paths") {
		t.Errorf("show says deny_protected: %s", out)
	}
	if got := (&app{configPath: configOf(t)}).agentConfig(root).Harnesses["command"]; got.Program != "run-agent" || strings.Join(got.Args, " ") != "{story}" {
		t.Errorf("the command's program is still flai serve agent set's: %+v", got)
	}
	out, _, _ = runIn(t, root, "serve", "agent", "show", "--json")
	var shown struct {
		Harnesses map[string]struct {
			Program       string
			Args          []string
			Guard         string `json:"guard"`
			DenyProtected bool   `json:"deny_protected"`
		} `json:"harnesses"`
	}
	if err := json.Unmarshal([]byte(out), &shown); err != nil {
		t.Fatalf("show --json: %v %s", err, out)
	}
	if got := shown.Harnesses["command"]; got.Program != "run-agent" || strings.Join(got.Args, " ") != "{story}" || got.Guard != "none" || !got.DenyProtected {
		t.Errorf("show --json gives the command's guard and deny_protected: %s", out)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--guard", "")
	if got := harnesses()["command"]; !reflect.DeepEqual(got, config.HarnessHost{DenyProtected: true}) {
		t.Errorf(`--guard "" clears only the guard: %+v`, got)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected=false")
	if _, ok := harnesses()["command"]; ok {
		t.Errorf("an entry that lifts nothing is removed: %+v", harnesses())
	}
	// on claude-code, beside its program; --reset clears both
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--guard", "none", "--deny-protected")
	if got := harnesses()["claude-code"]; got.Program != "/opt/claude" || got.Guard != config.GuardNone || !got.DenyProtected {
		t.Errorf("guard and deny_protected keep the program: %+v", got)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--guard", "none")
	runIn(t, root, "serve", "agent", "harness", "command", "--reset")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--reset")
	if got := harnesses(); got != nil {
		t.Errorf("--reset clears guard and deny_protected: %+v", got)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--guard", "off"); code == 0 || !strings.Contains(errOut, `flai serve agent harness claude-code --guard none, or clear it with --guard`) {
		t.Errorf("a guard other than none is refused, naming what sets it: %d %s", code, errOut)
	}
	for _, bad := range [][]string{{"command", "--program", "x"}, {"command", "--", "x"}, {"command", "--reset", "--guard", "none"}, {"--guard", "none"}, {"--deny-protected"}} {
		if _, _, code := runIn(t, root, append([]string{"serve", "agent", "harness"}, bad...)...); code == 0 {
			t.Errorf("%q was taken", bad)
		}
	}
	if _, errOut, _ := runIn(t, root, "serve", "agent", "harness", "command", "--program", "x"); !strings.Contains(errOut, "flai serve agent set") {
		t.Errorf("the command's program points at flai serve agent set: %s", errOut)
	}
	if got := harnesses(); got != nil {
		t.Errorf("a refusal changes nothing: %+v", got)
	}
}

// S-0355 (ADR-0131): flai serve is told at every look what the operator
// lifts on each harness, the command's among them, and whether the
// project's auto-approve host action is on.
func TestServeAgentConfigCarriesTheLiftsAndAutoApprove(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	if got := (&app{configPath: configOf(t)}).agentConfig(root); len(got.Lifts) != 0 || got.AutoApprove {
		t.Errorf("nothing lifted: %+v %v", got.Lifts, got.AutoApprove)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--guard", "none")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--deny-protected")
	got := (&app{configPath: configOf(t)}).agentConfig(root).Lifts
	if got["command"] != (serve.HarnessLifts{GuardNone: true}) || got["claude-code"] != (serve.HarnessLifts{DenyProtected: true}) {
		t.Errorf("lifts: %+v", got)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "auto-approve"); code != 0 {
		t.Fatal(errOut)
	}
	if !(&app{configPath: configOf(t)}).agentConfig(root).AutoApprove {
		t.Error("auto-approve is on for the project")
	}
}

// S-0349 (ADR-0129): a provider's entry is the operator's, set, shown,
// changed field by field, and removed on the host; a new entry is whole, a
// refused value changes nothing, clear keeps the providers, and nothing the
// command prints or writes holds the key, which no flag takes and the
// command never reads.
func TestServeAgentProvider(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	const key = "sk-never-shown-0123456789"
	t.Setenv("LITELLM_API_KEY", key)
	const masterKey = "sk-master-never-shown-0123456789"
	t.Setenv("LITELLM_MASTER_KEY", masterKey)
	root := tempProject(t)
	var seen strings.Builder
	run := func(args ...string) (string, string, int) {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"serve", "agent", "provider"}, args...)...)
		seen.WriteString(out + errOut)
		return out, errOut, code
	}
	saved := func() map[string]manifest.Provider {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Agent.Providers
	}
	if out, _, _ := run(); !strings.Contains(out, "no provider is set on this host; add one with flai serve agent provider <name> --api <api> --base-url <url> --key-env <variable>") {
		t.Errorf("none set: %s", out)
	}
	if out, _, _ := run("litellm"); !strings.Contains(out, "provider litellm has no entry on this host; add it with flai serve agent provider litellm --api") {
		t.Errorf("a name with no entry: %s", out)
	}
	if _, errOut, code := run("litellm", "--api", "openai-chat"); code == 0 || !strings.Contains(errOut, "has no base_url; set it with flai serve agent provider litellm --base-url") || saved() != nil {
		t.Errorf("a new entry must be whole: %d %s %+v", code, errOut, saved())
	}
	out, errOut, code := run("litellm", "--api", "openai-chat", "--base-url", "http://127.0.0.1:4000", "--key-env", "LITELLM_API_KEY", "--model", "opus=anthropic/claude-opus-5-5", "--model", "sonnet=anthropic/claude-sonnet-5")
	if code != 0 || !strings.Contains(out, "provider litellm: openai-chat at http://127.0.0.1:4000, its key in the variable LITELLM_API_KEY\n  model opus: anthropic/claude-opus-5-5\n  model sonnet: anthropic/claude-sonnet-5\n") {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	if out, _, _ := run("litellm"); !strings.Contains(out, "provider litellm: openai-chat at http://127.0.0.1:4000") {
		t.Errorf("print: %s", out)
	}
	// the flags given change those fields; --model given replaces the aliases
	run("litellm", "--base-url", "https://litellm.example.com", "--model", "haiku=anthropic/claude-haiku-4-5")
	want := manifest.Provider{API: "openai-chat", BaseURL: "https://litellm.example.com", KeyEnv: "LITELLM_API_KEY", Models: map[string]string{"haiku": "anthropic/claude-haiku-4-5"}}
	if got := saved(); !reflect.DeepEqual(got, map[string]manifest.Provider{"litellm": want}) {
		t.Errorf("update: %+v", got)
	}
	run("litellm", "--api", "anthropic-messages")
	want.API = "anthropic-messages"
	if got := saved()["litellm"]; !reflect.DeepEqual(got, want) {
		t.Errorf("--api alone keeps the rest: %+v", got)
	}
	if js, _, _ := run("litellm", "--json"); !strings.Contains(js, `"litellm": {`) || !strings.Contains(js, `"key_env": "LITELLM_API_KEY"`) || strings.Contains(js, "spend_key_env") {
		t.Errorf("--json: %s", js)
	}
	// S-0357: --spend-key-env names the variable of the key that reads the
	// spend log, shown where set, and an empty value clears it
	out, errOut, code = run("litellm", "--spend-key-env", "LITELLM_MASTER_KEY")
	if code != 0 || !strings.Contains(out, "its key in the variable LITELLM_API_KEY\n  spend log read with the key in the variable LITELLM_MASTER_KEY\n  model haiku:") {
		t.Errorf("--spend-key-env: %d %s %s", code, out, errOut)
	}
	if got, w := saved()["litellm"], (manifest.Provider{API: want.API, BaseURL: want.BaseURL, KeyEnv: want.KeyEnv, SpendKeyEnv: "LITELLM_MASTER_KEY", Models: want.Models}); !reflect.DeepEqual(got, w) {
		t.Errorf("--spend-key-env keeps the rest: %+v", got)
	}
	if js, _, _ := run("litellm", "--json"); !strings.Contains(js, `"spend_key_env": "LITELLM_MASTER_KEY"`) {
		t.Errorf("--json with spend_key_env: %s", js)
	}
	if out, _, _ := run("litellm", "--spend-key-env", ""); strings.Contains(out, "spend log") || !reflect.DeepEqual(saved()["litellm"], want) {
		t.Errorf("--spend-key-env \"\" clears it: %s %+v", out, saved()["litellm"])
	}
	for _, bad := range []struct {
		args []string
		want string
	}{
		{[]string{"litellm", "--api", "bedrock"}, "is not an API flai speaks"},
		{[]string{"litellm", "--base-url", "litellm.example.com"}, "is not an absolute http or https URL"},
		{[]string{"litellm", "--key-env", "litellm-key"}, "is not the name of an environment variable"},
		{[]string{"litellm", "--spend-key-env", masterKey}, "spend_key_env is not the name of an environment variable"},
		{[]string{"litellm", "--model", "opus"}, "is not alias=name"},
		{[]string{"--api", "openai-chat"}, "name the provider"},
		{[]string{"litellm", "--reset", "--api", "openai-chat"}, "--reset takes no other flag"},
		{[]string{"lite llm", "--api", "openai-chat", "--base-url", "https://x.example.com", "--key-env", "X"}, "is not a provider name"},
	} {
		if _, errOut, code := run(bad.args...); code == 0 || !strings.Contains(errOut, bad.want) {
			t.Errorf("%q is refused, saying why: %d %s", bad.args, code, errOut)
		}
	}
	if got := saved()["litellm"]; !reflect.DeepEqual(got, want) {
		t.Errorf("a refused value changes nothing: %+v", got)
	}
	runIn(t, root, "serve", "agent", "clear")
	if got := saved(); len(got) != 1 {
		t.Errorf("clear keeps the providers: %+v", got)
	}
	if out, _, _ := run("litellm", "--reset"); !strings.Contains(out, "provider litellm has no entry on this host") || saved() != nil {
		t.Errorf("reset: %s %+v", out, saved())
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(seen.String()+string(data), key) || strings.Contains(seen.String()+string(data), masterKey) {
		t.Errorf("the key is in what flai printed or wrote:\n%s\n%s", seen.String(), data)
	}
}

// S-0392: flai serve forge shows the variable each forge's token is read
// from, GITHUB_TOKEN by default, and whether it is set here; --token-env
// names another and --reset goes back; an unknown forge and a name that is
// not a variable's are refused, changing nothing; and nothing the command
// prints or writes holds the token, which no flag takes.
func TestServeForge(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	const token = "ghp_never-shown-0123456789"
	t.Setenv("FLAI_GH_TOKEN", token)
	root := tempProject(t)
	var seen strings.Builder
	run := func(args ...string) (string, string, int) {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"serve", "forge"}, args...)...)
		seen.WriteString(out + errOut)
		return out, errOut, code
	}
	saved := func() map[string]config.ForgeHost {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Forges
	}
	if out, _, code := run(); code != 0 || out != "forge github: its token in the variable GITHUB_TOKEN (default), not set in this environment\n" {
		t.Errorf("unset: %d %q", code, out)
	}
	out, errOut, code := run("github", "--token-env", "FLAI_GH_TOKEN")
	if code != 0 || out != "forge github: its token in the variable FLAI_GH_TOKEN, set in this environment\n" {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	if got := saved(); !reflect.DeepEqual(got, map[string]config.ForgeHost{"github": {TokenEnv: "FLAI_GH_TOKEN"}}) {
		t.Errorf("saved: %+v", got)
	}
	if out, _, _ := run("github"); !strings.Contains(out, "forge github: its token in the variable FLAI_GH_TOKEN, set") {
		t.Errorf("show one: %s", out)
	}
	if js, _, _ := run("--json"); !strings.Contains(js, `"github": {`) || !strings.Contains(js, `"token_env": "FLAI_GH_TOKEN"`) || !strings.Contains(js, `"token_set": true`) || !strings.Contains(js, `"default": false`) {
		t.Errorf("--json: %s", js)
	}
	for _, bad := range []struct {
		args []string
		want string
	}{
		{[]string{"gitlab", "--token-env", "GITLAB_TOKEN"}, "gitlab\\\"; it calls github"},
		{[]string{"gitlab"}, "flai calls no forge"},
		{[]string{"github", "--token-env", "gh-token"}, "forges.github.token_env is not the name of an environment variable"},
		{[]string{"github", "--token-env", token}, "never the token"},
		{[]string{"github", "--token-env", ""}, "--token-env needs the name of a variable, such as GITHUB_TOKEN"},
		{[]string{"--token-env", "X"}, "name the forge"},
		{[]string{"github", "--reset", "--token-env", "X"}, "--reset takes no --token-env"},
	} {
		if _, errOut, code := run(bad.args...); code == 0 || !strings.Contains(errOut, bad.want) {
			t.Errorf("%q is refused, saying why: %d %s", bad.args, code, errOut)
		}
	}
	if got := saved()["github"]; got.TokenEnv != "FLAI_GH_TOKEN" {
		t.Errorf("a refused value changes nothing: %+v", got)
	}
	// S-0393: flai serve is told the variable it keeps from every agent
	if got := (&app{configPath: configOf(t)}).agentConfig(root).ForgeTokenEnvs; len(got) != 1 || got[0] != "FLAI_GH_TOKEN" {
		t.Errorf("the forge tokens kept from the agents: %q, want FLAI_GH_TOKEN", got)
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if out, _, _ := run("github", "--reset"); out != "forge github: its token in the variable GITHUB_TOKEN (default), not set in this environment\n" || saved() != nil {
		t.Errorf("reset: %q %+v", out, saved())
	}
	if got := (&app{configPath: configOf(t)}).agentConfig(root).ForgeTokenEnvs; len(got) != 1 || got[0] != "GITHUB_TOKEN" {
		t.Errorf("the forge tokens kept from the agents once reset: %q, want GITHUB_TOKEN", got)
	}
	if strings.Contains(seen.String()+string(data), token) {
		t.Errorf("the token is in what flai printed or wrote:\n%s\n%s", seen.String(), data)
	}
}

// S-0294 (ADR-0108): --auto-restarts sets how many times flai serve restarts
// a story's agent on its own, alone without touching the command, the name,
// or the harnesses; show says it and whether it is the default; a negative
// count is refused; and clear keeps it.
func TestServeAgentAutoRestarts(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "agent", "show"); !strings.Contains(out, "automatic restarts: 2 (default)") {
		t.Errorf("unset, show says the default: %s", out)
	}
	runIn(t, root, "serve", "agent", "set", "--name", "builder", "--", "run-agent", "{story}")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude")
	saved := func() config.AgentStart {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Agent
	}
	for _, c := range []struct {
		n, said string
		want    int
	}{{"0", "automatic restarts: 0 (off)", 0}, {"3", "automatic restarts: 3\n", 3}} {
		out, errOut, code := runIn(t, root, "serve", "agent", "set", "--auto-restarts", c.n)
		if code != 0 || !strings.Contains(out, c.said) {
			t.Fatalf("--auto-restarts %s: %d %s %s", c.n, code, out, errOut)
		}
		got := saved()
		if got.AutoRestarts == nil || *got.AutoRestarts != c.want {
			t.Errorf("--auto-restarts %s is written: %v", c.n, got.AutoRestarts)
		}
		if p := got.ForProject(root); strings.Join(p.Command, " ") != "run-agent {story}" || p.Name != "builder" || p.Harnesses["claude-code"].Program != "/opt/claude" {
			t.Errorf("--auto-restarts %s alone keeps the project's command, name, and harnesses: %+v", c.n, got)
		}
	}
	if js, _, _ := runIn(t, root, "serve", "agent", "show", "--json"); !strings.Contains(js, `"auto_restarts": 3`) {
		t.Errorf("show --json: %s", js)
	}
	_, errOut, code := runIn(t, root, "serve", "agent", "set", "--auto-restarts", "-1")
	if code == 0 || !strings.Contains(errOut, "0 or more") {
		t.Errorf("a negative count is refused, saying what to give: %d %s", code, errOut)
	}
	if got := saved(); got.AutoRestarts == nil || *got.AutoRestarts != 3 {
		t.Errorf("a refused count changes nothing: %v", got.AutoRestarts)
	}
	runIn(t, root, "serve", "agent", "clear")
	if got := saved(); len(got.Command) != 0 || got.AutoRestarts == nil || *got.AutoRestarts != 3 {
		t.Errorf("clear removes the command and keeps the limit: %+v", got)
	}
}

// S-0082: the checks commands are the operator's, named argument lists with
// no default, several at once, managed on the host and nowhere else.
func TestServeChecksCommands(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "checks")
	if !strings.Contains(out, "no command is named here; the manifest's own checks: is used instead") ||
		!strings.Contains(out, "timeout: 15 minutes") || !strings.Contains(out, "off for this project") {
		t.Errorf("by default: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "checks: off everywhere") || !strings.Contains(out, "in a story's worktree") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--", "go", "test", "./..."); code == 0 {
		t.Error("set needs --name")
	}
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--name", "flai"); code == 0 {
		t.Error("set needs a program")
	}
	out, errOut, code := runIn(t, root, "serve", "checks", "set", "--name", "flai", "--", "scripts/flai-test.sh")
	if code != 0 || !strings.Contains(out, `flai: "scripts/flai-test.sh"`) {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	out, _, code = runIn(t, root, "serve", "checks", "set", "--name", "flaiover", "--", "bash", "-c", "cd flaiover; pnpm test")
	if code != 0 || !strings.Contains(out, `flai: "scripts/flai-test.sh"`) || !strings.Contains(out, `flaiover: "bash" "-c" "cd flaiover; pnpm test"`) {
		t.Fatalf("a second name adds, does not replace: %d %s", code, out)
	}
	out, _, code = runIn(t, root, "serve", "checks", "set", "--name", "flai", "--", "scripts/flai-test.sh", "--fast")
	if code != 0 || !strings.Contains(out, `flai: "scripts/flai-test.sh" "--fast"`) || strings.Count(out, "flai:") != 1 {
		t.Fatalf("the same name replaces: %d %s", code, out)
	}
	cfg, _ := os.ReadFile(configOf(t))
	if !strings.Contains(string(cfg), `"cd flaiover; pnpm test"`) || !strings.Contains(string(cfg), `"name": "flaiover"`) {
		t.Errorf("kept as an argument list, as written: %s", cfg)
	}
	if _, errOut, code := runIn(t, root, "config", "set", "checks.commands", "rm -rf /"); code == 0 {
		t.Errorf("flai config set does not reach it: %s", errOut)
	}
	out, _, code = runIn(t, root, "serve", "checks", "timeout", "20")
	if code != 0 || !strings.Contains(out, "timeout: 20 minutes") {
		t.Fatalf("timeout: %d %s", code, out)
	}
	if _, _, code := runIn(t, root, "serve", "checks", "timeout", "0"); code == 0 {
		t.Error("timeout must be positive")
	}
	runIn(t, root, "serve", "enable", "checks")
	js, _, _ := runIn(t, root, "serve", "checks", "show", "--json")
	if !strings.Contains(js, `"enabled_here": true`) || !strings.Contains(js, `"timeout_minutes": 20`) || !strings.Contains(js, `"flai"`) || !strings.Contains(js, `"flaiover"`) {
		t.Errorf("show --json: %s", js)
	}
	if out, _, _ := runIn(t, root, "serve", "checks", "clear", "flai"); strings.Contains(out, `flai:`) || !strings.Contains(out, `flaiover:`) {
		t.Errorf("clear one name: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "checks", "clear"); !strings.Contains(out, "no command is named here") {
		t.Errorf("clear with no name clears all: %s", out)
	}
}

// S-0438: components is a host action the operator enables as checks is, and
// components.idle a host setting flai serve components idle keeps, which
// flai serve's supervisor of the dev servers reads.
func TestServeComponentsActionAndIdle(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "components")
	if !strings.Contains(out, "idle: 30m (default)") || !strings.Contains(out, "max_servers: 3 (default): at most this many dev servers run at once per project") ||
		!strings.Contains(out, "components action is off for this project; flai serve enable components") {
		t.Errorf("by default: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "components: off everywhere") || !strings.Contains(out, "Vite dev server") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	for _, bad := range []string{"0", "-5m", "soon", "30", ""} {
		if _, errOut, code := runIn(t, root, "serve", "components", "idle", "--", bad); code == 0 || !strings.Contains(errOut, "positive duration") {
			t.Errorf("idle %q is refused: %d %s", bad, code, errOut)
		}
	}
	out, errOut, code := runIn(t, root, "serve", "components", "idle", "45m")
	if code != 0 || !strings.Contains(out, "idle: 45m: a dev server") {
		t.Fatalf("idle 45m: %d %s %s", code, out, errOut)
	}
	// S-0442: components.max_servers, a positive whole number, 3 by default
	for _, bad := range []string{"0", "-1", "two", "2.5", ""} {
		if _, errOut, code := runIn(t, root, "serve", "components", "max-servers", "--", bad); code == 0 || !strings.Contains(errOut, "positive whole number") {
			t.Errorf("max-servers %q is refused: %d %s", bad, code, errOut)
		}
	}
	out, errOut, code = runIn(t, root, "serve", "components", "max-servers", "5")
	if code != 0 || !strings.Contains(out, "max_servers: 5: at most") || !strings.Contains(out, "idle: 45m: a dev server") {
		t.Fatalf("max-servers 5: %d %s %s", code, out, errOut)
	}
	cfg, _ := os.ReadFile(configOf(t))
	if !strings.Contains(string(cfg), `"components"`) || !strings.Contains(string(cfg), `"idle": "45m"`) || !strings.Contains(string(cfg), `"max_servers": 5`) {
		t.Errorf("kept in the host's configuration: %s", cfg)
	}
	if out, _, code := runIn(t, root, "serve", "enable", "components"); code != 0 || !strings.Contains(out, "components enabled for") {
		t.Fatalf("enable components: %d %s", code, out)
	}
	js, _, _ := runIn(t, root, "serve", "components", "show", "--json")
	var shown struct {
		Idle        string `json:"idle"`
		IdleSeconds int    `json:"idle_seconds"`
		IdleSet     bool   `json:"idle_set"`
		MaxServers  int    `json:"max_servers"`
		MaxSet      bool   `json:"max_servers_set"`
		EnabledHere bool   `json:"enabled_here"`
	}
	if err := json.Unmarshal([]byte(js), &shown); err != nil || shown.Idle != "45m" || shown.IdleSeconds != 2700 || !shown.IdleSet || shown.MaxServers != 5 || !shown.MaxSet || !shown.EnabledHere {
		t.Errorf("show --json: %v %s", err, js)
	}
	if st, _, _ := runIn(t, root, "hostapi", "components.status"); !strings.Contains(st, `"enabled":true`) || !strings.Contains(st, `"enable":"flai serve enable components"`) {
		t.Errorf("the dashboard is told the action is on: %s", st)
	}

	// flai serve's own host has the supervisor, on components.idle and
	// components.max_servers as set, and stops it when asked
	var sink bytes.Buffer
	a := &app{configPath: configOf(t), out: &sink, errOut: &sink, cwd: root}
	host, stop := a.serveHost(t.Context())
	if host.Components == nil || host.Components.Idle() != 45*time.Minute || !host.Enabled(hostapi.ActionComponents, root) {
		t.Errorf("serveHost: %+v", host)
	}
	if most, ok := host.Components.(interface{ MaxServers() int }); !ok || most.MaxServers() != 5 {
		t.Errorf("serveHost's supervisor reads components.max_servers: %+v", host.Components)
	}
	stopped := make(chan struct{})
	go func() { stop(); close(stopped) }()
	select {
	case <-stopped:
	case <-time.After(30 * time.Second):
		t.Fatal("serveHost's stop did not return")
	}
	if _, err := host.Components.Start(t.Context(), root, "t", manifest.Project{Name: "ui", Path: "ui", Kind: manifest.ProjectKindSvelteKit}, "", ""); err == nil || !strings.Contains(err.Error(), "flai serve is stopping") {
		t.Errorf("once stopped, it starts no dev server: %v", err)
	}
	if out, _, _ := runIn(t, root, "serve", "disable", "components"); !strings.Contains(out, "components disabled for") {
		t.Errorf("disable components: %s", out)
	}
}

// S-0082: the host's own named commands are used when there are any; the
// manifest's checks: only when the host names none. Never a merge.
func TestChecksConfigResolvesHostOverManifest(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "serve", "checks", "show") // creates the config file, as flai itself would
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	a := &app{configPath: configOf(t)}

	// nothing named anywhere
	cc, err := a.checksConfig(repo)
	if err != nil || len(cc.Commands) != 0 || cc.Enabled || cc.Timeout != serve.DefaultChecksTimeout {
		t.Fatalf("nothing named: %+v %v", cc, err)
	}

	// the manifest names checks:, host config names none: the manifest's are used
	manifestYAML := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" +
		"checks:\n  - name: flai\n    command: [scripts/flai-test.sh]\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifestYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err = workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	cc, err = a.checksConfig(repo)
	if err != nil || len(cc.Commands) != 1 || cc.Commands[0].Name != "flai" {
		t.Fatalf("falls back to the manifest: %+v %v", cc, err)
	}

	// the host names a different one: the host's is used, not a merge
	runIn(t, root, "serve", "checks", "set", "--name", "flaiover", "--", "true")
	cc, err = a.checksConfig(repo)
	if err != nil || len(cc.Commands) != 1 || cc.Commands[0].Name != "flaiover" {
		t.Fatalf("the host's own commands, not merged with the manifest's: %+v %v", cc, err)
	}
	if cc.Enabled {
		t.Error("still off until enabled")
	}
	runIn(t, root, "serve", "enable", "checks")
	if cc, err := a.checksConfig(repo); err != nil || !cc.Enabled {
		t.Errorf("enabled: %+v %v", cc, err)
	}
}

// S-0105: the settings host action, the agent's name without a command, and
// what the dashboard's settings page is told.
func TestServeSettingsAreShownToTheDashboard(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	// S-0122: the served projects, without the dashboard's address or the credential's path
	if err := serve.DirFor(cfgPath).Register(serve.Entry{Key: "harbour", Name: "Harbour", Root: "/p/harbour", URL: "http://127.0.0.1:4242", KeyFile: "/p/harbour/k"}); err != nil {
		t.Fatal(err)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "settings: off everywhere") || !strings.Contains(out, "only a shell turns this off") {
		t.Errorf("actions names settings: %s", out)
	}
	// --attended-minutes is retired (S-0116): still accepted, it says so and does nothing
	if _, errOut, code := runIn(t, root, "serve", "agent", "set", "--name", "builder", "--attended-minutes", "9"); code != 0 || !strings.Contains(errOut, "attended-minutes") || !strings.Contains(errOut, "does nothing") {
		t.Fatalf("name without a command: %d %s", code, errOut)
	}
	cfg := (&app{configPath: configOf(t)}).agentConfig(root)
	if cfg.Name != "builder" || len(cfg.Command) != 0 {
		t.Errorf("name: %+v", cfg)
	}
	if loaded, _, _ := (&app{configPath: configOf(t)}).loadConfig(); loaded.Agent.AttendedMinutes != 0 {
		t.Errorf("the retired minutes were stored: %+v", loaded.Agent)
	}
	runIn(t, root, "serve", "enable", "settings")
	runIn(t, root, "serve", "checks", "set", "--name", "unit", "--", "go", "test", "./...")
	// S-0488: a temp directory's name may hold 4242, so the import root holds
	// it always, and only the registered entry's URL is refused below
	importRoot := filepath.Join(t.TempDir(), "harbour-4242")
	if err := os.Mkdir(importRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	runIn(t, root, "serve", "import", "add", importRoot)
	runIn(t, root, "agent", "set", "--harness", "claude-code", "--model", "claude-haiku-4-5")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx, rec := perf.Start(context.Background())
	got, _ := json.Marshal((&app{configPath: configOf(t)}).hostSettings(ctx, mainRootOf(repo)))
	// S-0605: settings.get says where its time goes
	phases := map[string]int{}
	for _, ph := range rec.Phases() {
		phases[ph.Name] = ph.Count
	}
	if want := map[string]int{"settings.config": 1, "settings.served": 1, "settings.project": 1}; !reflect.DeepEqual(phases, want) {
		t.Errorf("phases %v, want %v", phases, want)
	}
	for _, want := range []string{`"here":true,"means":"change this project's host settings`, `"here":false,"means":"publish when you press Publish`, `"name":"builder"`,
		`"claude-code":{"args":["--permission-mode","acceptEdits","--allowedTools","Bash,mcp__flai"],"ask_capable":true,"deny_protected":false,"guard":"","guard_capable":true,"program":"claude","set":false}`,
		`"command":{"args":[],"ask_capable":false,"deny_protected":false,"guard":"","guard_capable":false,"program":"","set":false}`,
		`{"name":"unit","command":["go","test","./..."]}`, `"timeout_minutes":15`, `"import_roots":["` + importRoot + `"]`, `"default_agent":{"harness":"claude-code","model":"claude-haiku-4-5"}`, `"mcp":{"bind":"127.0.0.1","bind_set":false,"running":false}`,
		`"projects":{"running":false,"served":[{"key":"harbour","name":"Harbour","root":"/p/harbour","from":"registry","state":"unavailable","reason":"`, `"unserved":[]}`} {
		if !strings.Contains(string(got), want) {
			t.Errorf("settings lack %s:\n%s", want, got)
		}
	}
	if strings.Contains(string(got), `"token"`) || strings.Contains(string(got), "127.0.0.1:4242") || strings.Contains(string(got), "/p/harbour/k") {
		t.Errorf("no token in the settings: %s", got)
	}
}

// S-0355, ADR-0131: settings.get gives each harness its guard and
// deny_protected; an entry holding only those keeps set false, so the page
// still says the program and arguments are the defaults; and the command is
// there with flai serve agent set's program and arguments.
func TestServeSettingsShowEachHarnessGuardAndDenyProtected(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	for _, args := range [][]string{{"claude-code", "--guard", "none"}, {"command", "--guard", "none", "--deny-protected"}} {
		if _, errOut, code := runIn(t, root, append([]string{"serve", "agent", "harness"}, args...)...); code != 0 {
			t.Fatalf("%q: %s", args, errOut)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), mainRootOf(repo)))
	for _, want := range []string{
		`"claude-code":{"args":["--permission-mode","acceptEdits","--allowedTools","Bash,mcp__flai"],"ask_capable":true,"deny_protected":false,"from":"project","guard":"none","guard_capable":true,"program":"claude","set":false}`,
		`"command":{"args":["{story}"],"ask_capable":false,"deny_protected":true,"from":"project","guard":"none","guard_capable":false,"program":"run-agent","set":false}`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("settings lack %s:\n%s", want, got)
		}
	}
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude")
	got, _ = json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), mainRootOf(repo)))
	if want := `"guard":"none","guard_capable":true,"program":"/opt/claude","set":true}`; !strings.Contains(string(got), want) {
		t.Errorf("a program beside the guard is set: want %s in\n%s", want, got)
	}
}

// S-0359, ADR-0129: settings.get's host carries the host's providers, each
// entry's fields and whether the variables its key_env and spend_key_env
// (S-0357) name are set and not empty in flai serve's environment, never a
// variable's value.
func TestServeSettingsShowTheProvidersNeverTheirKeys(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	const key = "sk-never-shown-0123456789"
	t.Setenv("LITELLM_API_KEY", key)
	t.Setenv("OPENROUTER_API_KEY", "")
	const masterKey = "sk-master-never-shown-0123456789"
	t.Setenv("LITELLM_MASTER_KEY", masterKey)
	t.Setenv("FLAI_TEST_UNSET_GATEWAY_KEY", "") // restored after the test, unset in it
	if err := os.Unsetenv("FLAI_TEST_UNSET_GATEWAY_KEY"); err != nil {
		t.Fatal(err)
	}
	root := tempProject(t)
	runIn(t, root, "serve", "agent", "provider") // creates the config, which an app with no logger cannot
	agent := func() map[string]any {
		t.Helper()
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), mainRootOf(repo)))
		if strings.Contains(string(got), key) || strings.Contains(string(got), masterKey) {
			t.Errorf("a key is in the settings: %s", got)
		}
		var out struct {
			Agent map[string]any `json:"agent"`
		}
		if err := json.Unmarshal(got, &out); err != nil {
			t.Fatal(err)
		}
		return out.Agent
	}
	if got, _ := json.Marshal(agent()["providers"]); string(got) != `{}` {
		t.Errorf("no provider set is an empty object: %s", got)
	}
	for _, args := range [][]string{
		{"litellm", "--api", "anthropic-messages", "--base-url", "http://127.0.0.1:4000", "--key-env", "LITELLM_API_KEY", "--spend-key-env", "LITELLM_MASTER_KEY", "--model", "opus=anthropic/claude-opus-5-5"},
		{"openrouter", "--api", "openai-chat", "--base-url", "https://openrouter.ai/api/v1", "--key-env", "OPENROUTER_API_KEY"},
		{"gateway", "--api", "openai-responses", "--base-url", "https://gw.example.com", "--key-env", "FLAI_TEST_UNSET_GATEWAY_KEY", "--spend-key-env", "FLAI_TEST_UNSET_GATEWAY_KEY"},
	} {
		if _, errOut, code := runIn(t, root, append([]string{"serve", "agent", "provider"}, args...)...); code != 0 {
			t.Fatalf("set %s: %s", args[0], errOut)
		}
	}
	got, _ := json.Marshal(agent()["providers"])
	want := `{"gateway":{"api":"openai-responses","base_url":"https://gw.example.com","key_env":"FLAI_TEST_UNSET_GATEWAY_KEY","key_set":false,"models":{},"spend_key_env":"FLAI_TEST_UNSET_GATEWAY_KEY","spend_key_set":false},` +
		`"litellm":{"api":"anthropic-messages","base_url":"http://127.0.0.1:4000","key_env":"LITELLM_API_KEY","key_set":true,"models":{"opus":"anthropic/claude-opus-5-5"},"spend_key_env":"LITELLM_MASTER_KEY","spend_key_set":true},` +
		`"openrouter":{"api":"openai-chat","base_url":"https://openrouter.ai/api/v1","key_env":"OPENROUTER_API_KEY","key_set":false,"models":{},"spend_key_env":"","spend_key_set":false}}`
	if string(got) != want {
		t.Errorf("providers:\n%s\nwant\n%s", got, want)
	}
}

// S-0295, ADR-0096: settings.get carries the shared paths as the manifest
// has them now; settings.shared changes them through flai shared under the
// settings action, its refusal a rule with flai's reason; and
// settings.shared_check answers what flai shared check does, for anyone.
func TestSettingsSharedThroughTheHostAPI(t *testing.T) {
	t.Setenv("FLAI_AGENT", "tester")
	root := sharedProject(t)
	inProcess(t)
	hostapiCall := func(method, params string) (string, int) {
		t.Helper()
		out, _, code := runIn(t, root, "hostapi", method, params)
		return out, code
	}
	shared := func() string {
		t.Helper()
		out, code := hostapiCall("settings.get", `{}`)
		var got struct {
			Host struct {
				Shared json.RawMessage `json:"shared"`
			} `json:"host"`
		}
		if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil {
			t.Fatalf("settings.get: %d %v %s", code, err, out)
		}
		return string(got.Host.Shared)
	}
	if got := shared(); got != `[]` {
		t.Errorf("settings.get with no shared paths: %s, want []", got)
	}

	before := readManifestShared(t, root)
	if out, code := hostapiCall("settings.shared", `{"action":"add","pattern":"design/adrs","request_id":"req-00000051"}`); code == 0 || !strings.Contains(out, "flai serve enable settings") || !strings.Contains(out, "-32012") {
		t.Errorf("with the settings action off: %d %s", code, out)
	}
	if got := readManifestShared(t, root); !reflect.DeepEqual(got, before) {
		t.Errorf("a refused change wrote the manifest: %v", got)
	}

	runIn(t, root, "serve", "enable", "settings")
	for i, c := range []struct{ params, want string }{
		{`{"action":"add","pattern":"design/adrs"`, `"data":{"added":["design/adrs"],"shared":["design/adrs"]}`},
		{`{"action":"add","pattern":"docs/users/*.md"`, `"data":{"added":["docs/users/*.md"],"shared":["design/adrs","docs/users/*.md"]}`},
		{`{"action":"add","pattern":"template"`, `"data":{"added":["template"],"shared":["design/adrs","docs/users/*.md","template"]}`},
		{`{"action":"remove","pattern":"template"`, `"data":{"removed":["template"],"shared":["design/adrs","docs/users/*.md"]}`},
	} {
		out, code := hostapiCall("settings.shared", c.params+`,"request_id":"req-0000006`+strconv.Itoa(i)+`"}`)
		if code != 0 || !strings.Contains(out, c.want) {
			t.Errorf("settings.shared %s: %d %s, want %s", c.params, code, out, c.want)
		}
	}
	if got := shared(); got != `["design/adrs","docs/users/*.md"]` {
		t.Errorf("settings.get after the changes: %s", got)
	}

	for i, c := range []struct{ params, why string }{
		{`{"action":"add","pattern":"design/adrs"`, "in the list already"},
		{`{"action":"add","pattern":"/design"`, "absolute"},
		{`{"action":"remove","pattern":"docs"`, "not in the list"},
	} {
		out, code := hostapiCall("settings.shared", c.params+`,"request_id":"req-0000007`+strconv.Itoa(i)+`"}`)
		if code == 0 || !strings.Contains(out, "-32011") || !strings.Contains(out, c.why) {
			t.Errorf("settings.shared %s: %d %s, want a rule saying %q", c.params, code, out, c.why)
		}
	}
	if got := readManifestShared(t, root); !reflect.DeepEqual(got, []string{"design/adrs", "docs/users/*.md"}) {
		t.Errorf("the refusals changed the manifest: %v", got)
	}

	// a read for anyone: with the settings action off again it still answers
	runIn(t, root, "serve", "disable", "settings")
	out, code := hostapiCall("settings.shared_check", `{"paths":["design/adrs/0001-x.md","docs/users","cli"],"story":"S-0001"}`)
	var got struct {
		Data []sharedCheck `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil {
		t.Fatalf("settings.shared_check: %d %v %s", code, err, out)
	}
	want := []sharedCheck{
		{Entry: "design/adrs/0001-x.md", Path: "design/adrs/0001-x.md", Shared: true, Pattern: "design/adrs"},
		{Entry: "docs/users", Path: "docs/users"},
		{Entry: "cli", Path: "flai"},
		{Entry: "docs/users/flai.md", Path: "docs/users/flai.md", Shared: true, Pattern: "docs/users/*.md", Story: "S-0001"},
		{Entry: "flai", Path: "flai", Story: "S-0001"},
		{Entry: "design/adrs", Path: "design/adrs", Shared: true, Pattern: "design/adrs", Story: "S-0001"},
		{Entry: "template/CHANGELOG.md", Path: "template/CHANGELOG.md", Story: "S-0001"},
	}
	if !reflect.DeepEqual(got.Data, want) {
		t.Errorf("settings.shared_check =\n%+v\nwant\n%+v", got.Data, want)
	}
	if out, code := hostapiCall("settings.shared_check", `{"story":"S-0099"}`); code == 0 || !strings.Contains(out, "-32004") {
		t.Errorf("a story that is not there: %d %s", code, out)
	}

	js, _, _ := runIn(t, root, "serve", "journal")
	if !strings.Contains(js, "flai shared add --autocommit -- design/adrs") || !strings.Contains(js, "flai shared remove --autocommit -- template") || strings.Contains(js, "shared_check") {
		t.Errorf("the journal names each change and no check:\n%s", js)
	}
}

// S-0211, ADR-0084: the settings page is told when flai serve plans again on
// its own: whether the plan host action is on, the replan policy and whether
// the manifest sets it, and the schedule as written with its next run; a
// value flai cannot read comes with its error.
func TestSettingsShowThePlanningTriggers(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	a := &app{configPath: configOf(t), errOut: &bytes.Buffer{}, clock: func() time.Time { return time.Date(2026, 10, 4, 9, 30, 0, 0, time.UTC) }}
	planning := func() string {
		t.Helper()
		got, _ := json.Marshal(a.hostSettings(context.Background(), mainRootOf(repo)).(map[string]any)["planning"])
		return string(got)
	}
	if got, want := planning(), `{"plan":false,"replan":"deterministic","replan_set":false,"schedule":""}`; got != want {
		t.Errorf("defaults:\n got %s\nwant %s", got, want)
	}
	m := filepath.Join(root, "system-flow.yaml")
	data, err := os.ReadFile(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(m, append(data, []byte("planning:\n  replan: agent\n  schedule: daily\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "plan"); code != 0 {
		t.Fatalf("enable plan: %d %s", code, errOut)
	}
	if got, want := planning(), `{"next":"2026-10-05T00:00:00Z","plan":true,"replan":"agent","replan_set":true,"schedule":"daily"}`; got != want {
		t.Errorf("set:\n got %s\nwant %s", got, want)
	}
	if err := os.WriteFile(m, append(data, []byte("planning:\n  replan: sometimes\n  schedule: \"61 * * * *\"\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	got := planning()
	for _, want := range []string{`"replan":"sometimes"`, `"replan_error":"planning.replan \"sometimes\" is not a replan policy`, `"schedule":"61 * * * *"`, `"schedule_error":"planning.schedule`} {
		if !strings.Contains(got, want) {
			t.Errorf("errors lack %s:\n%s", want, got)
		}
	}
	if strings.Contains(got, `"next"`) {
		t.Errorf("a schedule flai cannot read has no next run: %s", got)
	}
}

// S-0229: settings.get gives the strategic agents' settings, each with its
// value as the manifest has it or unset, its default, and what it means,
// with whether the settings action lets the dashboard change them and what
// turns it on; settings.manifest changes them through flai manifest set only
// with that action on, and flai's refusal reaches the caller as Refused with
// the field and the reason, the manifest left as it was.
func TestSettingsShowAndSaveTheStrategicAgentsSettings(t *testing.T) {
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	file := filepath.Join(root, manifest.File)
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	fixture := "orchestration:\n  policy: wsjf\n  permissions:\n    promote_to_ready: true\n  release:\n    policy: threshold\n    value: 500\nplanning:\n  agent:\n    model: claude-sonnet-5\n"
	if err := os.WriteFile(file, append(data, fixture...), 0o644); err != nil {
		t.Fatal(err)
	}
	inProcess(t)
	type setting struct {
		Key     string          `json:"key"`
		Kind    string          `json:"kind"`
		Values  []string        `json:"values"`
		Default string          `json:"default"`
		Meaning string          `json:"meaning"`
		Risk    string          `json:"risk"`
		Whole   bool            `json:"whole"`
		Value   json.RawMessage `json:"value"`
		Set     bool            `json:"set"`
	}
	type strategic struct {
		Editable bool      `json:"editable"`
		Enable   string    `json:"enable"`
		Settings []setting `json:"settings"`
	}
	read := func() (strategic, map[string]setting) {
		t.Helper()
		out, _, code := runIn(t, root, "hostapi", "settings.get", `{}`)
		var got struct {
			Host struct {
				Strategic strategic `json:"strategic"`
			} `json:"host"`
		}
		if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil {
			t.Fatalf("settings.get: %d %v %s", code, err, out)
		}
		byKey := map[string]setting{}
		for _, s := range got.Host.Strategic.Settings {
			byKey[s.Key] = s
		}
		return got.Host.Strategic, byKey
	}

	st, byKey := read()
	if st.Editable || st.Enable != "flai serve enable settings" || len(st.Settings) != len(manifest.Settings()) {
		t.Errorf("strategic: editable %v, enable %q, %d settings, want false, flai serve enable settings, %d", st.Editable, st.Enable, len(st.Settings), len(manifest.Settings()))
	}
	for _, c := range []struct {
		key, kind, value, def string
		set                   bool
	}{
		{"orchestration.permissions.promote_to_ready", "boolean", `true`, "false", true},
		{"orchestration.permissions.accept_reviews", "boolean", ``, "false", false},
		{"orchestration.permissions.answer_threads", "choice", ``, "off", false},
		{"orchestration.policy", "choice", `"wsjf"`, "fifo", true},
		{"orchestration.release.policy", "choice", `"threshold"`, "judgement", true},
		{"orchestration.release.value", "number", `500`, "", true},
		{"orchestration.release.count", "number", ``, "", false},
		{"planning.agent", "agent", `{"model":"claude-sonnet-5"}`, "", true},
		{"planning.replan", "choice", ``, "deterministic", false},
		{"analysis.schedule", "cron", ``, "", false},
		{"analysis.formula", "text", ``, "", false},
	} {
		s, ok := byKey[c.key]
		if !ok || s.Kind != c.kind || string(s.Value) != c.value || s.Default != c.def || s.Set != c.set || s.Meaning == "" {
			t.Errorf("%s: %+v (value %s), want kind %s, value %q, default %q, set %v, and a meaning", c.key, s, s.Value, c.kind, c.value, c.def, c.set)
		}
	}
	if s := byKey["orchestration.permissions.promote_to_ready"]; s.Risk == "" || !reflect.DeepEqual(s.Values, []string{"true", "false"}) {
		t.Errorf("a permission says its risk and its values: %+v", s)
	}
	if s := byKey["orchestration.release.count"]; !s.Whole {
		t.Errorf("a count is a whole number: %+v", s)
	}
	if s := byKey["orchestration.policy"]; !reflect.DeepEqual(s.Values, manifest.OrderPolicies) {
		t.Errorf("a choice lists its values: %+v", s)
	}

	save := func(params string) (string, int) {
		t.Helper()
		out, _, code := runIn(t, root, "hostapi", "settings.manifest", params)
		return out, code
	}
	written := func() string {
		t.Helper()
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	before := written()
	if out, code := save(`{"set":{"planning.replan":"agent"},"request_id":"req-00000091"}`); code == 0 || !strings.Contains(out, "-32012") || !strings.Contains(out, "flai serve enable settings") {
		t.Errorf("with the settings action off: %d %s", code, out)
	}
	runIn(t, root, "serve", "enable", "settings")
	if out, code := save(`{"set":{"planning.replan":"agent","orchestration.release.value":-1},"request_id":"req-00000092"}`); code == 0 || !strings.Contains(out, "-32010") ||
		!strings.Contains(out, `"field":"orchestration.release.value"`) || !strings.Contains(out, `"reason":"`) {
		t.Errorf("flai's refusal reaches the caller with the field and the reason: %d %s", code, out)
	}
	if got := written(); got != before {
		t.Errorf("a refused save wrote the manifest:\n%s", got)
	}

	if out, code := save(`{"set":{"planning.replan":"agent"},"unset":["orchestration.policy"],"request_id":"req-00000093"}`); code != 0 || !strings.Contains(out, `{"key":"planning.replan","value":"agent"}`) {
		t.Fatalf("a save: %d %s", code, out)
	}
	st, byKey = read()
	if !st.Editable {
		t.Error("with the settings action on, strategic says editable")
	}
	if s := byKey["planning.replan"]; string(s.Value) != `"agent"` || !s.Set {
		t.Errorf("planning.replan after the save: %+v (value %s)", s, s.Value)
	}
	if s := byKey["orchestration.policy"]; s.Set || len(s.Value) != 0 || s.Default != "fifo" {
		t.Errorf("orchestration.policy after it was unset: %+v (value %s)", s, s.Value)
	}
	if js, _, _ := runIn(t, root, "serve", "journal"); !strings.Contains(js, "flai manifest set --autocommit --unset=orchestration.policy -- planning.replan=agent") {
		t.Errorf("the journal names the change as the command it ran:\n%s", js)
	}
}

// S-0420: settings.get's strategic block lists each role the project defines
// that is not built-in, with its agent, the schedule its trigger names, and
// the permissions its file declares, each with its value and, for a
// permission, the calls the file lists under it; the strategic agents' keys
// stay as they were. settings.manifest writes each of them through flai
// manifest set only with the settings action on.
func TestSettingsShowAndSaveADefinedRolesSettings(t *testing.T) {
	scheduled := strings.Replace(reviewerRole, "{schedule: planning.schedule, when: each-ready-story}", "{schedule: roles.reviewer.schedule}", 1)
	root := rolesProject(t, map[string]string{"reviewer.md": scheduled})
	inProcess(t)
	type setting struct {
		Key    string          `json:"key"`
		Kind   string          `json:"kind"`
		Value  json.RawMessage `json:"value"`
		Set    bool            `json:"set"`
		Allows []string        `json:"allows"`
	}
	type role struct {
		Name     string    `json:"name"`
		Action   string    `json:"action"`
		Settings []setting `json:"settings"`
	}
	read := func() (bool, []role, int) {
		t.Helper()
		out, _, code := runIn(t, root, "hostapi", "settings.get", `{}`)
		var got struct {
			Host struct {
				Strategic struct {
					Editable bool              `json:"editable"`
					Settings []json.RawMessage `json:"settings"`
					Roles    []role            `json:"roles"`
				} `json:"strategic"`
			} `json:"host"`
		}
		if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil {
			t.Fatalf("settings.get: %d %v %s", code, err, out)
		}
		st := got.Host.Strategic
		return st.Editable, st.Roles, len(st.Settings)
	}
	editable, rs, n := read()
	if editable || n != len(manifest.Settings()) || len(rs) != 1 || rs[0].Name != "reviewer" || rs[0].Action != "review" {
		t.Fatalf("strategic: editable %v, %d settings, roles %+v; want false, %d, the reviewer alone", editable, n, rs, len(manifest.Settings()))
	}
	var keys []string
	for _, s := range rs[0].Settings {
		keys = append(keys, s.Key+":"+s.Kind)
	}
	if want := []string{"roles.reviewer.agent:agent", "roles.reviewer.schedule:cron", "roles.reviewer.permissions.publish:boolean"}; !slices.Equal(keys, want) {
		t.Errorf("the reviewer's settings %v, want %v", keys, want)
	}
	if s := rs[0].Settings[2]; s.Set || !slices.Equal(s.Allows, []string{"release_publish"}) {
		t.Errorf("the publish permission: %+v", s)
	}

	params := `{"set":{"roles.reviewer.agent":{"model":"claude-sonnet-5"},"roles.reviewer.schedule":"daily","roles.reviewer.permissions.publish":true},"request_id":"req-00000191"}`
	if out, _, code := runIn(t, root, "hostapi", "settings.manifest", params); code == 0 || !strings.Contains(out, "-32012") || !strings.Contains(out, "flai serve enable settings") {
		t.Errorf("with the settings action off: %d %s", code, out)
	}
	runIn(t, root, "serve", "enable", "settings")
	if out, _, code := runIn(t, root, "hostapi", "settings.manifest", strings.Replace(params, "req-00000191", "req-00000192", 1)); code != 0 {
		t.Fatalf("a save: %d %s", code, out)
	}
	editable, rs, _ = read()
	got := map[string]string{}
	for _, s := range rs[0].Settings {
		got[s.Key] = fmt.Sprintf("%s/%v", s.Value, s.Set)
	}
	if !editable || got["roles.reviewer.agent"] != `{"model":"claude-sonnet-5"}/true` || got["roles.reviewer.schedule"] != `"daily"/true` || got["roles.reviewer.permissions.publish"] != "true/true" {
		t.Errorf("after the save: editable %v, %v", editable, got)
	}
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil || m.Roles["reviewer"].Schedule != "daily" || !m.Roles["reviewer"].Permissions.Publish {
		t.Errorf("the manifest after the save: %+v, %v", m.Roles, err)
	}
}

// ADR-0067: auto-publish is the operator's shell tool, outside the workflow.
// The serve actions listing names it and serve enable turns it on, saying
// what it has flai push do; the dashboard is neither told of it nor can turn
// it on or off.
func TestAutoPublishIsTheShellsAlone(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "actions")
	if !strings.Contains(out, "auto-publish: off everywhere\n  has flai push --pending, run in a shell") || !strings.Contains(out, "push: off everywhere\n  lets a dashboard publish when you press Publish") {
		t.Errorf("actions lists auto-publish for the operator: %s", out)
	}
	out, errOut, code := runIn(t, root, "serve", "enable", "auto-publish")
	if code != 0 || !strings.Contains(out, "auto-publish enabled for t\n  it has flai push --pending") || strings.Contains(out, "journal") || !strings.Contains(out, "flai serve disable auto-publish") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if info, _, _ := runIn(t, root, "hostapi", "project.info"); strings.Contains(info, "auto-publish") || !strings.Contains(info, `"push":false`) {
		t.Errorf("project.info names auto-publish: %s", info)
	}
	runIn(t, root, "serve", "enable", "settings")
	out, _, code = runIn(t, root, "hostapi", "settings.action", `{"action":"auto-publish","on":false,"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a20"}`)
	if code == 0 || !strings.Contains(out, "flai serve disable auto-publish") || !strings.Contains(out, "in a shell on the host") {
		t.Errorf("settings.action refuses auto-publish, naming what to run: %d %s", code, out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg, _, _ := (&app{configPath: configOf(t)}).loadConfig(); !cfg.ActionEnabled("auto-publish", mainRootOf(repo)) {
		t.Error("the refusal left auto-publish as the shell set it")
	}
	got, _ := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), mainRootOf(repo)))
	if strings.Contains(string(got), "auto-publish") || !strings.Contains(string(got), `"name":"push"`) {
		t.Errorf("the settings page is told of auto-publish: %s", got)
	}
}

// S-0387, ADR-0158: auto-upgrade, which has flai serve upgrade the flai on
// the host, is the operator's shell tool: serve enable turns it on for the
// project, saying what it has flai serve do, and the dashboard is neither
// told of it nor can turn it on or off.
func TestAutoUpgradeIsTheShellsAlone(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, errOut, code := runIn(t, root, "serve", "enable", "auto-upgrade")
	if code != 0 || !strings.Contains(out, "auto-upgrade enabled for t\n  it has flai serve ask flai host") || strings.Contains(out, "journal") || !strings.Contains(out, "flai serve disable auto-upgrade") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if !(&app{configPath: configOf(t)}).host().Enabled(hostapi.ActionAutoUpgrade, mainRootOf(repo)) {
		t.Error("flai serve sees auto-upgrade on for the project")
	}
	if info, _, _ := runIn(t, root, "hostapi", "project.info"); strings.Contains(info, "auto-upgrade") {
		t.Errorf("project.info names auto-upgrade: %s", info)
	}
	runIn(t, root, "serve", "enable", "settings")
	out, _, code = runIn(t, root, "hostapi", "settings.action", `{"action":"auto-upgrade","on":false,"request_id":"5d2e3b61-8c4a-4e1f-9b62-1d3e5f7a9b31"}`)
	if code == 0 || !strings.Contains(out, "flai serve disable auto-upgrade") {
		t.Errorf("settings.action refuses auto-upgrade, naming what to run: %d %s", code, out)
	}
	got, _ := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), mainRootOf(repo)))
	if strings.Contains(string(got), "auto-upgrade") {
		t.Errorf("the settings page is told of auto-upgrade: %s", got)
	}
}

// ADR-0158: flai serve asks the host it runs under for the project's release
// by version, as flai host upgrade --version does, and hears whether the
// host restarts on it or why it refused; a serve under no host asks none.
func TestServeAsksItsHostForTheRelease(t *testing.T) {
	t.Parallel()
	var asked []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		asked = append(asked, r.Method+" "+r.URL.Path+" "+string(body)+" "+r.Header.Get("Authorization"))
		if strings.Contains(string(body), "9.9.9") {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"flai 9.9.9 is not published; published are 1.50.0"}`))
			return
		}
		_, _ = w.Write([]byte(`{"upgrade":{"installed":"1.50.0"},"restarting":true}`))
	}))
	t.Cleanup(srv.Close)
	upgrade := hostUpgrade(hostMCP{c: &host.Client{URL: srv.URL, Token: "host-token"}})
	if upgrade == nil {
		t.Fatal("a serve under a host has no way to ask it for an upgrade")
	}
	restarting, err := upgrade(context.Background(), "1.50.0")
	if err != nil || !restarting {
		t.Errorf("an installed release: %v %v", restarting, err)
	}
	if len(asked) != 1 || asked[0] != `POST /upgrade {"version":"1.50.0"} Bearer host-token` {
		t.Errorf("the host was asked %q", asked)
	}
	if restarting, err := upgrade(context.Background(), "9.9.9"); err == nil || restarting || !strings.Contains(err.Error(), "not published") {
		t.Errorf("a refused release: %v %v", restarting, err)
	}
	if hostUpgrade(nil) != nil {
		t.Error("a serve under no host asks one for an upgrade")
	}
}

// S-0116: flai serve agent restart says why it will not start an agent.
func TestServeAgentRestartSaysWhyItRefuses(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	if _, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-0001"); code == 0 || !strings.Contains(errOut, "rule: the agent host action is off for this project") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "enable", "agent")
	if _, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-1"); code == 0 || !strings.Contains(errOut, "S-0001") {
		t.Errorf("no such story: %d %s", code, errOut)
	}
}

// S-0118: flai serve agent restart queues the agent of a story in ready
// while the in-progress limit is full, and says so.
func TestServeAgentRestartQueuesWhenTheLimitIsFull(t *testing.T) {
	t.Parallel()
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfg)
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	for _, s := range []struct{ id, file string }{{"S-0001", "S-0001-slice.md"}, {"S-0002", "S-0002-busy.md"}} {
		runIn(t, root, "story", "new", strings.TrimSuffix(s.file[7:], ".md"), "--epic", "E-0001")
		file := filepath.Join(root, "wip/kanban/stories", s.file)
		body, _ := os.ReadFile(file)
		_ = os.WriteFile(file, []byte(strings.Replace(string(body), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
		if _, errOut, code := runIn(t, root, "move", s.id, "ready"); code != 0 {
			t.Fatal(errOut)
		}
	}
	runIn(t, root, "serve", "enable", "agent")
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected") // the command cannot ask (ADR-0131)
	if out, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code != 0 {
		t.Fatalf("start: %d %s %s", code, out, errOut)
	}
	// its agent failed with the story in ready; another story takes the only slot
	agents := filepath.Join(string(serve.DirFor(cfg)), "agents.json")
	var all map[string]serve.AgentState
	data, _ := os.ReadFile(agents)
	if err := json.Unmarshal(data, &all); err != nil {
		t.Fatal(err)
	}
	for key, st := range all {
		r := st.Stories["S-0001"]
		r.PID, r.Ended, r.Outcome, r.Why = 0, r.Started, serve.OutcomeFailed, "ended (exit 1) with S-0001 in ready"
		all[key] = st
	}
	data, _ = json.Marshal(all)
	_ = os.WriteFile(agents, data, 0o600)
	runIn(t, root, "move", "S-0002", "in-progress")
	_ = os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte("---\ntitle: Board\nstatus: active\nwip_limits:\n  ready: 5\n  in-progress: 1\n  review: 3\n---\n\n# Board\n"), 0o644)

	out, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-1")
	if code != 0 || !strings.HasPrefix(out, "queued another agent for S-0001: the in-progress limit, a full review, or a hold leaves no room") {
		t.Fatalf("restart: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-1"); code == 0 || !strings.Contains(errOut, "already queued") {
		t.Errorf("twice: %d %s", code, errOut)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	if !strings.Contains(js, "queued another agent for S-0001") {
		t.Errorf("the queue is not journalled: %s", js)
	}
}

// S-0115: flai serve agent start starts a ready story's agent now, recorded
// where flai serve tracks it, and says why when it will not.
func TestServeAgentStartStartsAReadyStorysAgent(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	onTheHostsCommand(t, root)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	file := filepath.Join(root, "wip/kanban/stories/S-0001-slice.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "rule: the agent host action is off for this project") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "enable", "agent")
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "rule: S-0001 is in backlog") {
		t.Errorf("in backlog: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "ready"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "rule: S-0001's harness is the host's command, and none is set for this project: set it on the dashboard's Harnesses page, /settings/harnesses, or with flai serve agent set -- <program> [args...]") {
		t.Errorf("nothing to start it with: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	// S-0355 (ADR-0131): the command cannot ask before a protected write, so
	// it starts only once the operator lifts that
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "harness command cannot ask before a protected write") || !strings.Contains(errOut, "flai serve agent harness command --deny-protected") {
		t.Errorf("started where nothing lifts asking: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected")
	out, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1")
	if code != 0 || !strings.HasPrefix(out, "started true (command) for S-0001 as agent-S-0001") {
		t.Fatalf("start: %d %s %s", code, out, errOut)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	if !strings.Contains(js, "started true (command) for S-0001") {
		t.Errorf("the start is not journalled: %s", js)
	}
}

// S-0140: flai serve agent commit starts an agent to commit what a story in
// review left uncommitted in its worktree, and says why when it will not.
func TestServeAgentCommitStartsAnAgentForAnUncommittedWorktree(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	file := filepath.Join(root, "wip/kanban/stories/S-0001-slice.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	if _, errOut, code := runIn(t, root, "serve", "agent", "commit", "S-1"); code == 0 || !strings.Contains(errOut, "rule: the agent host action is off for this project") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "enable", "agent")
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	for _, args := range [][]string{{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"task", "new", "Do", "--story", "S-0001"}, {"move", "S-0001", "review"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "commit", "S-1"); code == 0 || !strings.Contains(errOut, "rule: S-0001 has no worktree at .flai-cache/worktrees/S-0001") {
		t.Errorf("no worktree: %d %s", code, errOut)
	}
}

// S-0142: flai serve agent stream prints what a story's agent said and did,
// from the log flai serve gave it, and the host answers agent.stream from it.
func TestServeAgentStreamPrintsTheLog(t *testing.T) {
	t.Parallel()
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfg)
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "sf-S-0001.log")
	lines := `{"type":"system","subtype":"init","model":"m"}` + "\n" +
		`{"type":"assistant","message":{"content":[{"type":"text","text":"Reading.\nThen writing."},{"type":"tool_use","name":"Bash","input":{"description":"List files"}}]}}` + "\n" +
		`{"type":"user","message":{"content":[{"type":"tool_result","content":"nope","is_error":true}]}}` + "\n"
	if err := os.WriteFile(log, []byte(lines), 0o600); err != nil {
		t.Fatal(err)
	}
	dir := serve.DirFor(cfg)
	_ = os.MkdirAll(string(dir), 0o700)
	state := map[string]serve.AgentState{mainRootOf(repo): {Stories: map[string]*serve.AgentRun{
		"S-0001": {Story: "S-0001", Agent: "agent-S-0001", Started: "2026-09-29T05:00:00Z", Ended: "2026-09-29T05:01:00Z", Log: log, Outcome: serve.OutcomeWorked},
	}}}
	data, _ := json.Marshal(state)
	if err := os.WriteFile(filepath.Join(string(dir), "agents.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}

	out, errOut, code := runIn(t, root, "serve", "agent", "stream", "S-1")
	want := "S-0001: agent-S-0001, started 2026-09-29T05:00:00Z, ended 2026-09-29T05:01:00Z\n" +
		"session   session started (m)\n" +
		"text      Reading.\n          Then writing.\n" +
		"tool      Bash: List files\n" +
		"result!   nope\n"
	if code != 0 || out != want {
		t.Errorf("stream: %d %s\n%q\nwant\n%q", code, errOut, out, want)
	}
	// --follow ends with the agent, which has ended
	if out, _, code := runIn(t, root, "serve", "agent", "stream", "S-0001", "--follow", "--from", "0"); code != 0 || out != want {
		t.Errorf("follow an ended agent: %d %q", code, out)
	}
	first := strconv.Itoa(strings.Index(lines, "\n") + 1)
	js, _, _ := runIn(t, root, "serve", "agent", "stream", "S-0001", "--json", "--from", first)
	var got serve.StreamRead
	if err := json.Unmarshal([]byte(js), &got); err != nil || strconv.FormatInt(got.From, 10) != first || got.Next != int64(len(lines)) || len(got.Entries) != 3 || got.Running {
		t.Errorf("--json from %s: %v %s", first, err, js)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "stream", "S-0002"); code == 0 || !strings.Contains(errOut, "rule: flai serve has started no agent for S-0002") {
		t.Errorf("no agent: %d %s", code, errOut)
	}

	// the host answers agent.stream from the same state
	a := &app{configPath: configOf(t)}
	res, e := hostapi.MethodsFor("test", nil, a.host())["agent.stream"](context.Background(), channel.Project{Key: "t", Root: mainRootOf(repo)}, json.RawMessage(`{"story":"S-0001","after":`+first+`}`))
	if e != nil {
		t.Fatal(e)
	}
	if r := res.(*serve.StreamRead); strconv.FormatInt(r.From, 10) != first || len(r.Entries) != 3 {
		t.Errorf("agent.stream: %+v", r)
	}
}

// S-0177, ADR-0064: flai mcp's agent tools start a story's agent as flai
// serve agent restart does, under the same agent host action, and journal
// the agent that asked.
func TestTheMCPServerStartsAStorysAgentUnderTheAgentAction(t *testing.T) {
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_STARTED_BY", "") // the operator's own agent, whoever runs the tests
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	file := filepath.Join(root, "wip/kanban/stories/S-0001-slice.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	for _, to := range []string{"ready", "in-progress"} {
		if _, errOut, code := runIn(t, root, "move", "S-0001", to); code != 0 {
			t.Fatal(errOut)
		}
	}
	var out, errOut bytes.Buffer
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: root, clock: func() time.Time { return time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC) }}
	ctx := context.Background()
	if _, err := a.mcpAgents(ctx, "restart", root, "S-0001", "agent-ops"); err == nil || !strings.Contains(err.Error(), "the agent host action is off") {
		t.Errorf("action off: %v", err)
	}
	runIn(t, root, "serve", "enable", "agent")
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected") // the command cannot ask (ADR-0131)
	got, err := a.mcpAgents(ctx, "restart", root, "S-0001", "agent-ops")
	if err != nil || got.Story != "S-0001" || got.PID == 0 || got.Command != "true" {
		t.Fatalf("a story in progress with no agent here: %+v %v", got, err)
	}
	// an agent flai serve started does not start others
	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	if _, err := a.mcpAgents(ctx, "start", root, "S-0001", "agent-S-0002"); err == nil || !strings.Contains(err.Error(), "an agent flai serve started does not start agents") {
		t.Errorf("an agent flai serve started: %v", err)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "mcp.agent_restart"`, `"by": "agent-ops"`, `"outcome": "disabled"`, "agent-ops asked to restart S-0001's agent: started true"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
}

// orchestratorState writes the project's orchestrator run, and the
// analyzer's when there is one, as flai serve records them.
func orchestratorState(t *testing.T, cfg, root string, orchestrator, analyzer *serve.AgentRun) {
	t.Helper()
	dir := serve.DirFor(cfg)
	_ = os.MkdirAll(string(dir), 0o700)
	data, _ := json.Marshal(map[string]serve.AgentState{root: {Orchestrator: orchestrator, Analyzer: analyzer}})
	if err := os.WriteFile(filepath.Join(string(dir), "agents.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// S-0228: the stop holds the orchestrator stopped, here a run that failed
// and waits to be started again, and the start lifts the hold and starts
// it, each refused while the orchestrate action is off; the host's
// orchestrate.stop and orchestrate.start run them, and agent.status says
// the orchestrator is held.
func TestServeOrchestrateHoldsTheOrchestratorStopped(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfg)
	root := tempProject(t)
	onTheHostsCommand(t, root)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	main := mainRootOf(repo)
	for _, verb := range []string{"stop", "start"} {
		if _, errOut, code := runIn(t, root, "serve", "orchestrate", verb); code == 0 || !strings.Contains(errOut, "rule: the orchestrate host action is off for this project, so flai serve runs no orchestrator to "+verb) {
			t.Errorf("%s while off: %d %s", verb, code, errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "orchestrate"); code != 0 {
		t.Fatal(errOut)
	}
	// with no hold and no run going, start starts a run for asked (ADR-0209),
	// which this project, on the host's command with none set, cannot make
	if _, errOut, code := runIn(t, root, "serve", "orchestrate", "start"); code == 0 || !strings.Contains(errOut, "harness is the host's command, and none is set for this project") || !strings.Contains(errOut, "/settings/harnesses") {
		t.Errorf("start with no hold: %d %s", code, errOut)
	}
	failed := &serve.AgentRun{Agent: "orchestrator", Harness: "command", Started: "2026-09-15T20:59:00Z", Ended: "2026-09-15T20:59:00Z", Session: "s1",
		Error: "no harness", Outcome: serve.OutcomeFailed, Why: "could not be started: no harness"}
	orchestratorState(t, cfg, main, failed, nil)
	out, errOut, code := runIn(t, root, "serve", "orchestrate", "stop", "--json")
	var said struct {
		Held         bool           `json:"held"`
		Orchestrator serve.AgentRun `json:"orchestrator"`
	}
	if err := json.Unmarshal([]byte(out), &said); code != 0 || err != nil || !said.Held || !said.Orchestrator.Held || said.Orchestrator.Session != "s1" || said.Orchestrator.Outcome != serve.OutcomeFailed {
		t.Fatalf("stop: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "orchestrate", "stop"); code == 0 || !strings.Contains(errOut, "rule: the orchestrator is already held stopped") {
		t.Errorf("stop twice: %d %s", code, errOut)
	}
	a := &app{configPath: configOf(t), out: &bytes.Buffer{}, errOut: &bytes.Buffer{}}
	p := channel.Project{Key: "t", Root: main}
	res, e := hostapi.MethodsFor("test", nil, a.host())["agent.status"](context.Background(), p, json.RawMessage(`{}`))
	if b, _ := json.Marshal(res); e != nil || !strings.Contains(string(b), `"orchestrator":{`) || !strings.Contains(string(b), `"session":"s1","outcome":"failed","why":"could not be started: no harness","held":true}`) {
		t.Errorf("agent.status: %+v %s", e, b)
	}

	// through the host: a start that fails, as with no harness and no
	// command, lifts the hold as well, and a stop holds the failed run
	inProcess(t)
	m := hostapi.MethodsFor("test", nil, a.host())
	if _, e := m["orchestrate.start"](context.Background(), p, json.RawMessage(`{"request_id":"req-00000001"}`)); e == nil || !strings.Contains(e.Message, "harness is the host's command, and none is set for this project") {
		t.Fatalf("orchestrate.start: %+v", e)
	}
	run := a.serveDir().AgentStates()[main].Orchestrator
	if run == nil || run.Held || run.Session == "s1" || run.Outcome != serve.OutcomeFailed {
		t.Fatalf("after the start: %+v, want a new failed run, not held", run)
	}
	res, e = m["orchestrate.stop"](context.Background(), p, json.RawMessage(`{"request_id":"req-00000002"}`))
	w, _ := res.(hostapi.Written)
	if err := json.Unmarshal(w.Data, &said); e != nil || err != nil || !said.Held || !said.Orchestrator.Held || said.Orchestrator.Session != run.Session {
		t.Fatalf("orchestrate.stop: %s %+v", w.Data, e)
	}
	if r := a.serveDir().AgentStates()[main].Orchestrator; !r.Held || r.Session != run.Session {
		t.Errorf("after the stop: %+v, want the failed run held", r)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "orchestrate.start"`, `"method": "orchestrate.stop"`,
		"stopped the orchestrator orchestrator on the operator's word, and held it stopped while orchestrate stays on: it was not running"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
}

// S-0228: the host answers agent.stream with role from the analyzer's newest
// run, as flai serve recorded it, and is not found with no run of the role.
func TestTheHostAnswersARolesStream(t *testing.T) {
	t.Parallel()
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfg)
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "t-analyzer.log")
	if err := os.WriteFile(log, []byte(`{"type":"assistant","message":{"content":[{"type":"text","text":"Reading the metrics."}]}}`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	orchestratorState(t, cfg, mainRootOf(repo), nil, &serve.AgentRun{Agent: "analyzer", Focus: "risk", Started: "2026-09-15T20:00:00Z", Ended: "2026-09-15T20:05:00Z", Outcome: serve.OutcomeWorked, Log: log})
	a := &app{configPath: configOf(t)}
	m := hostapi.MethodsFor("test", nil, a.host())["agent.stream"]
	p := channel.Project{Key: "t", Root: mainRootOf(repo)}
	res, e := m(context.Background(), p, json.RawMessage(`{"role":"analyze"}`))
	if e != nil {
		t.Fatal(e)
	}
	if r := res.(*serve.StreamRead); r.Role != "analyze" || r.Agent != "analyzer" || r.Running || len(r.Entries) != 1 || r.Entries[0].Text != "Reading the metrics." {
		t.Errorf("agent.stream: %+v", r)
	}
	if _, e := m(context.Background(), p, json.RawMessage(`{"role":"orchestrate"}`)); e == nil || e.Code != hostapi.NotFound || !strings.Contains(e.Message, "flai serve has started no orchestrator for this project") {
		t.Errorf("no orchestrator run: %+v", e)
	}
}

// S-0489: settings.get's host block carries dashboard.allowed_hosts as the
// configuration names it, normalized as the container is given it: empty
// until it is set.
func TestServeSettingsShowTheAllowedHosts(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	allowed := func() any {
		t.Helper()
		got, err := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), root))
		if err != nil {
			t.Fatal(err)
		}
		var settings map[string]any
		if err := json.Unmarshal(got, &settings); err != nil {
			t.Fatal(err)
		}
		v, ok := settings["allowed_hosts"]
		if !ok {
			t.Fatalf("settings lack allowed_hosts: %s", got)
		}
		return v
	}
	if _, errOut, code := runIn(t, root, "config", "get", "dashboard.allowed_hosts"); code != 0 { // creates the config
		t.Fatalf("config get: %s", errOut)
	}
	if got := allowed(); got != "" {
		t.Errorf("allowed_hosts before it is set = %v, want empty", got)
	}
	if _, errOut, code := runIn(t, root, "config", "set", "dashboard.allowed_hosts", " Flow.Example.com, ,board.example.net:8443"); code != 0 {
		t.Fatalf("config set: %s", errOut)
	}
	if got := allowed(); got != "flow.example.com,board.example.net:8443" {
		t.Errorf("allowed_hosts = %v, want the normalized list", got)
	}
}

// S-0498 (ADR-0205): flai serve notify set keeps a webhook and a shell
// command on the host, changing only what it is given, show prints them,
// with --json too, clear removes them, and a URL that is not http or https,
// an empty command, and a set with neither are refused, changing nothing.
func TestServeNotify(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	run := func(args ...string) (string, string, int) {
		t.Helper()
		return runIn(t, root, append([]string{"serve", "notify"}, args...)...)
	}
	saved := func() config.Notify {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Notify
	}
	if out, _, code := run("show"); code != 0 || !strings.Contains(out, "no reminder channel is set") {
		t.Errorf("unset: %d %q", code, out)
	}
	if out, errOut, code := run("set", "--url", "https://ntfy.example/flai"); code != 0 || out != "url: https://ntfy.example/flai\n" {
		t.Fatalf("set --url: %d %q %s", code, out, errOut)
	}
	if out, _, code := run("show"); code != 0 || out != "url: https://ntfy.example/flai\n" {
		t.Errorf("show: %d %q", code, out)
	}
	const command = `curl -s -H "Title: flai" -d @- https://ntfy.sh/flai`
	if _, errOut, code := run("set", "--command", command); code != 0 {
		t.Fatalf("set --command: %d %s", code, errOut)
	}
	if got := saved(); got != (config.Notify{URL: "https://ntfy.example/flai", Command: command}) {
		t.Errorf("set changes only what it is given: %+v", got)
	}
	if js, _, _ := run("show", "--json"); !strings.Contains(js, `"url": "https://ntfy.example/flai"`) || !strings.Contains(js, `"command": "curl -s -H \"Title: flai\" -d @- https://ntfy.sh/flai"`) {
		t.Errorf("--json: %s", js)
	}
	for _, bad := range []struct {
		args []string
		want string
	}{
		{[]string{"set", "--url", "ftp://x"}, "notify.url's scheme is ftp, and flai posts reminders over http or https alone"},
		{[]string{"set", "--url", ""}, "--url needs an http or https webhook"},
		{[]string{"set", "--command", " "}, "--command needs a shell command"},
		{[]string{"set"}, "give --url <webhook>, --command <shell command>, or both"},
	} {
		if _, errOut, code := run(bad.args...); code == 0 || !strings.Contains(errOut, bad.want) {
			t.Errorf("%q is refused, saying why: %d %s", bad.args, code, errOut)
		}
	}
	if got := saved(); got.URL != "https://ntfy.example/flai" || got.Command != command {
		t.Errorf("a refused value changes nothing: %+v", got)
	}
	if out, _, code := run("clear"); code != 0 || !strings.Contains(out, "no reminder channel is set") || saved() != (config.Notify{}) {
		t.Errorf("clear: %d %q %+v", code, out, saved())
	}
	if data, _ := os.ReadFile(cfgPath); strings.Contains(string(data), "notify") {
		t.Errorf("a cleared channel is not written: %s", data)
	}
}

// S-0417: the host action a role the project defines names is one flai
// serve enable and disable take for the project, flai serve actions and the
// settings page list it, flai serve is told whether it is on, and
// agent.status carries the role under roles, by name, with its action's
// state and its run (S-0420).
func TestServeEnableTakesADefinedRolesAction(t *testing.T) {
	t.Parallel()
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	cfgPath := configOf(t)
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "review: off everywhere\n  lets a dashboard start the role reviewer, which design/workflows/roles/reviewer.md defines") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{configPath: configOf(t)}).agentConfig(root).RoleActions["review"] {
		t.Error("review is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{configPath: configOf(t)}).host().Agent(context.Background(), root, true))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"reviewer":{"name":"reviewer","role":"review","path":"design/workflows/roles/reviewer.md","item":["story"],"trigger":["asked","schedule"],"action":{"name":"review","enabled":false},"run":null,"running":[],"runs":[],"held":false}`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "audit"); code == 0 || !strings.Contains(errOut, `there is no host action \"audit\"`) || !strings.Contains(errOut, ", review, ") {
		t.Errorf("an action no role names: %d %s", code, errOut)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "review"); code != 0 || !strings.Contains(out, "review enabled for t\n  it lets a dashboard start the role reviewer") || !strings.Contains(out, "every use is journalled") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{configPath: configOf(t)}).agentConfig(root).RoleActions["review"] {
		t.Error("enabled, flai serve is not told")
	}
	settings, err := json.Marshal((&app{configPath: configOf(t)}).hostSettings(context.Background(), root))
	if err != nil || !strings.Contains(string(settings), `"here":true,"means":"start the role reviewer`) || !strings.Contains(string(settings), `"name":"review"`) {
		t.Errorf("the settings page does not list review on: %s %v", settings, err)
	}
	run := `{"` + root + `":{"roles":{"reviewer":{"S-0001":{"story":"","item":"S-0001","role":"review","role_name":"reviewer","action":"review","harness":"command","command":"run-agent","agent":"reviewer-S-0001","pid":42,"started":"2026-10-28T10:00:00Z","session":"s","log":"/l/t-reviewer-20261028T100000Z.log"}}}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := status(); !strings.Contains(got, `"action":{"name":"review","enabled":true},"run":{"story":"","item":"S-0001",`) || !strings.Contains(got, `"role_name":"reviewer"`) || !strings.Contains(got, `"agent":"reviewer-S-0001"`) {
		t.Errorf("agent.status with a run: %s", got)
	}
	// a role file gone, its action is still turned off
	if err := os.Remove(filepath.Join(root, "design", "workflows", "roles", "reviewer.md")); err != nil {
		t.Fatal(err)
	}
	if out, errOut, code := runIn(t, root, "serve", "disable", "review"); code != 0 || !strings.Contains(out, "review disabled for t") {
		t.Errorf("disable once its role is gone: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "review"); code == 0 || !strings.Contains(errOut, `there is no host action \"review\"`) {
		t.Errorf("enable once its role is gone: %d %s", code, errOut)
	}
}

// S-0420: agent.status carries every role the project runs under roles, by
// name, the built-in ones from flai's baseline when the project has no file
// for them: each with its action's state, its run under way, its past runs
// newest first, and held, the orchestrator's hold; plans, orchestrator,
// analyzer, curator, past_orchestrators, and past_analyzers stay beside it.
func TestAgentStatusCarriesEveryRoleByName(t *testing.T) {
	t.Parallel()
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	cfgPath := configOf(t)
	if _, errOut, code := runIn(t, root, "serve", "enable", "review"); code != 0 {
		t.Fatalf("enable review: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "plan"); code != 0 {
		t.Fatalf("enable plan: %d %s", code, errOut)
	}
	run := func(started, ended, item string, pid int, extra string) string {
		r := `{"story":"","item":"` + item + `","agent":"a","pid":` + strconv.Itoa(pid) + `,"started":"` + started + `"`
		if ended != "" {
			r += `,"ended":"` + ended + `","outcome":"done"`
		}
		return r + extra + `}`
	}
	defined := `,"role":"review","role_name":"reviewer","action":"review"`
	state := `{"` + root + `":{` +
		`"plans":{"S-0001":` + run("2026-10-28T09:00:00Z", "", "S-0001", 41, "") + `,"E-0001":` + run("2026-10-28T08:00:00Z", "2026-10-28T08:10:00Z", "E-0001", 40, "") + `},` +
		`"past_plans":[` + run("2026-10-27T08:00:00Z", "2026-10-27T08:10:00Z", "E-0001", 39, "") + `],` +
		`"orchestrator":` + run("2026-10-28T07:00:00Z", "2026-10-28T07:05:00Z", "", 30, `,"held":true`) + `,` +
		`"past_orchestrators":[` + run("2026-10-27T07:00:00Z", "2026-10-27T07:05:00Z", "", 29, "") + `],` +
		`"analyzer":` + run("2026-10-28T06:00:00Z", "", "", 20, `,"focus":"risk"`) + `,` +
		`"roles":{"reviewer":{"S-0001":` + run("2026-10-28T10:00:00Z", "", "S-0001", 42, defined) + `,"S-0003":` + run("2026-10-28T05:00:00Z", "2026-10-28T05:30:00Z", "S-0003", 43, defined) + `}},` +
		`"past_roles":{"reviewer":[` + run("2026-10-27T10:00:00Z", "2026-10-27T10:30:00Z", "S-0001", 44, defined) + `]}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(state), 0o600); err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal((&app{configPath: cfgPath}).host().Agent(context.Background(), root, true))
	if err != nil {
		t.Fatal(err)
	}
	type runOf struct {
		Item    string `json:"item"`
		PID     int    `json:"pid"`
		Started string `json:"started"`
	}
	var got struct {
		Roles map[string]struct {
			Name    string   `json:"name"`
			Role    string   `json:"role"`
			Path    string   `json:"path"`
			Item    []string `json:"item"`
			Trigger []string `json:"trigger"`
			Action  struct {
				Name    string `json:"name"`
				Enabled bool   `json:"enabled"`
			} `json:"action"`
			Run     *runOf  `json:"run"`
			Running []runOf `json:"running"`
			Runs    []runOf `json:"runs"`
			Held    bool    `json:"held"`
		} `json:"roles"`
		Plans             map[string]runOf `json:"plans"`
		Orchestrator      *runOf           `json:"orchestrator"`
		Analyzer          *runOf           `json:"analyzer"`
		PastOrchestrators []runOf          `json:"past_orchestrators"`
		PastAnalyzers     []runOf          `json:"past_analyzers"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if names := slices.Sorted(maps.Keys(got.Roles)); !slices.Equal(names, []string{"analyzer", "curator", "orchestrator", "planner", "reviewer"}) {
		t.Fatalf("roles %v, want the four built-in ones and reviewer: %s", names, raw)
	}
	pids := func(runs []runOf) []int {
		out := []int{}
		for _, r := range runs {
			out = append(out, r.PID)
		}
		return out
	}
	r := got.Roles["reviewer"]
	if r.Name != "reviewer" || r.Role != "review" || r.Path != "design/workflows/roles/reviewer.md" || !slices.Equal(r.Item, []string{"story"}) ||
		!slices.Equal(r.Trigger, []string{"asked", "schedule"}) || r.Action.Name != "review" || !r.Action.Enabled || r.Held {
		t.Errorf("reviewer = %+v", r)
	}
	if r.Run == nil || r.Run.PID != 42 || r.Run.Item != "S-0001" || !slices.Equal(pids(r.Running), []int{42}) || !slices.Equal(pids(r.Runs), []int{43, 44}) {
		t.Errorf("reviewer's runs: run %+v, running %v, runs %v; want 42 under way, then 43 and 44 ended, newest first", r.Run, pids(r.Running), pids(r.Runs))
	}
	p := got.Roles["planner"]
	if p.Role != "plan" || p.Action.Name != "plan" || !p.Action.Enabled || p.Run == nil || p.Run.PID != 41 || !slices.Equal(pids(p.Runs), []int{40, 39}) || p.Held {
		t.Errorf("planner = %+v", p)
	}
	o := got.Roles["orchestrator"]
	if o.Action.Name != "orchestrate" || o.Action.Enabled || o.Run != nil || !slices.Equal(pids(o.Runs), []int{30, 29}) || !o.Held || len(o.Item) != 0 {
		t.Errorf("orchestrator = %+v; want its two ended runs, and held", o)
	}
	if a := got.Roles["analyzer"]; a.Action.Name != "analyze" || a.Run == nil || a.Run.PID != 20 || len(a.Runs) != 0 || a.Held {
		t.Errorf("analyzer = %+v", a)
	}
	if c := got.Roles["curator"]; c.Action.Name != "curate" || c.Run != nil || len(c.Runs) != 0 || c.Running == nil || c.Runs == nil {
		t.Errorf("curator = %+v; want no runs, as empty lists", c)
	}
	if len(got.Plans) != 2 || got.Orchestrator == nil || got.Orchestrator.PID != 30 || got.Analyzer == nil || !slices.Equal(pids(got.PastOrchestrators), []int{29}) || got.PastAnalyzers == nil {
		t.Errorf("the fields an older dashboard reads: %s", raw)
	}
}

// S-0611: agent.status without finished answers the open and live runs: a
// live run, a waiting one, and the newest run of each open story or planned
// item, and leaves out the runs of done, cancelled, or archived work and every
// past run; with finished true it answers them all, in the same shape.
func TestAgentStatusAnswersFinishedRunsOnlyWhenAsked(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	root := tempProject(t)
	for _, args := range [][]string{
		{"epic", "new", "Open epic"}, {"epic", "new", "Done epic"},
		{"story", "new", "Live", "--epic", "E-0001"}, {"story", "new", "Failed", "--epic", "E-0001"},
		{"story", "new", "Worked and done", "--epic", "E-0001"}, {"story", "new", "Asked and cancelled", "--epic", "E-0001"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	// S-0003 and E-0002 are done and S-0004 cancelled, as acceptance and a
	// cancellation leave them before they are archived
	for id, status := range map[string]string{"S-0003": workitem.Done, "E-0002": workitem.Done, "S-0004": workitem.Cancelled} {
		folder := "stories"
		if strings.HasPrefix(id, "E-") {
			folder = "epics"
		}
		paths, _ := filepath.Glob(filepath.Join(root, "wip", "kanban", folder, id+"-*.md"))
		if len(paths) != 1 {
			t.Fatalf("%s: %v", id, paths)
		}
		data, _ := os.ReadFile(paths[0])
		if !strings.Contains(string(data), "\nstatus: backlog\n") {
			t.Fatalf("%s has no backlog status: %s", id, data)
		}
		if err := os.WriteFile(paths[0], []byte(strings.Replace(string(data), "\nstatus: backlog\n", "\nstatus: "+status+"\n", 1)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	ended := func(item, story string, pid int, outcome string) *serve.AgentRun {
		return &serve.AgentRun{Item: item, Story: story, Agent: "a", PID: pid, Started: "2026-11-09T10:00:00Z", Ended: "2026-11-09T10:30:00Z", Outcome: outcome}
	}
	live := &serve.AgentRun{Story: "S-0001", Agent: "agent-S-0001", PID: 11, Started: "2026-11-09T11:00:00Z"}
	asked := ended("", "S-0004", 14, serve.OutcomeAsked)
	asked.Thread = "TH-0001"
	st := serve.AgentState{
		Stories: map[string]*serve.AgentRun{
			"S-0001": live,
			"S-0002": ended("", "S-0002", 12, serve.OutcomeFailed),
			"S-0003": ended("", "S-0003", 13, serve.OutcomeWorked),
			"S-0004": asked,
			"S-0099": ended("", "S-0099", 19, serve.OutcomeWorked), // archived: no file outside the archive
		},
		Plans: map[string]*serve.AgentRun{
			"E-0001": ended("E-0001", "", 21, serve.OutcomeWorked),
			"E-0002": ended("E-0002", "", 22, serve.OutcomeWorked),
			"S-0098": {Item: "S-0098", Agent: "planner-S-0098", PID: 23, Started: "2026-11-09T11:00:00Z"}, // live, on an item no list holds
		},
		PastPlans:         []*serve.AgentRun{ended("E-0001", "", 24, serve.OutcomeWorked)},
		Orchestrator:      ended("", "", 31, serve.OutcomeWorked),
		PastOrchestrators: []*serve.AgentRun{ended("", "", 32, serve.OutcomeWorked)},
		Analyzer:          ended("", "", 41, serve.OutcomeWorked),
		PastAnalyzers:     []*serve.AgentRun{ended("", "", 42, serve.OutcomeWorked)},
	}
	st.Orchestrator.Started, st.PastOrchestrators[0].Started = "2026-11-09T09:00:00Z", "2026-11-08T09:00:00Z"
	st.Analyzer.Started, st.PastAnalyzers[0].Started = "2026-11-09T08:00:00Z", "2026-11-08T08:00:00Z"
	data, err := json.Marshal(map[string]serve.AgentState{root: st})
	if err != nil {
		t.Fatal(err)
	}
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	type answer struct {
		Stories map[string]struct {
			State string          `json:"state"`
			Run   *serve.AgentRun `json:"run"`
		} `json:"stories"`
		Plans             map[string]*serve.AgentRun `json:"plans"`
		Orchestrator      *serve.AgentRun            `json:"orchestrator"`
		Analyzer          *serve.AgentRun            `json:"analyzer"`
		PastOrchestrators []*serve.AgentRun          `json:"past_orchestrators"`
		PastAnalyzers     []*serve.AgentRun          `json:"past_analyzers"`
		Roles             map[string]struct {
			Run  *serve.AgentRun   `json:"run"`
			Runs []*serve.AgentRun `json:"runs"`
		} `json:"roles"`
	}
	host := (&app{configPath: cfgPath, errOut: io.Discard}).host()
	ask := func(params string) answer {
		t.Helper()
		res, e := hostapi.MethodsFor("test", nil, host)["agent.status"](context.Background(), channel.Project{Key: "t", Root: root}, json.RawMessage(params))
		if e != nil {
			t.Fatal(e)
		}
		raw, err := json.Marshal(res.(map[string]any)["state"])
		if err != nil {
			t.Fatal(err)
		}
		var got answer
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatal(err)
		}
		if got.PastOrchestrators == nil || got.PastAnalyzers == nil || got.Stories == nil || got.Plans == nil {
			t.Errorf("%s: a key the dashboard reads is missing or null: %s", params, raw)
		}
		for name, r := range got.Roles {
			if r.Runs == nil {
				t.Errorf("%s: role %s's runs is null: %s", params, name, raw)
			}
		}
		return got
	}

	open := ask(`{}`)
	if ids := slices.Sorted(maps.Keys(open.Stories)); !slices.Equal(ids, []string{"S-0001", "S-0002", "S-0004"}) {
		t.Errorf("stories without finished: %v, want the live S-0001, the failed S-0002 on an open story, and the waiting S-0004", ids)
	}
	if open.Stories["S-0004"].State != serve.ActivityWaiting || open.Stories["S-0002"].State != serve.ActivityFailed {
		t.Errorf("states: %+v", open.Stories)
	}
	if ids := slices.Sorted(maps.Keys(open.Plans)); !slices.Equal(ids, []string{"E-0001", "S-0098"}) {
		t.Errorf("plans without finished: %v, want E-0001's, on an open epic, and the live S-0098", ids)
	}
	if len(open.PastOrchestrators) != 0 || len(open.PastAnalyzers) != 0 {
		t.Errorf("past runs without finished: %v %v, want none", open.PastOrchestrators, open.PastAnalyzers)
	}
	for name, r := range open.Roles {
		if len(r.Runs) != 0 {
			t.Errorf("role %s's runs without finished: %d, want none", name, len(r.Runs))
		}
	}
	if open.Orchestrator == nil || open.Orchestrator.PID != 31 || open.Analyzer == nil || open.Analyzer.PID != 41 || open.Roles["planner"].Run == nil || open.Roles["planner"].Run.PID != 23 {
		t.Errorf("the newest orchestrator and analyzer runs and the planner's run under way are kept: %+v %+v %+v", open.Orchestrator, open.Analyzer, open.Roles["planner"].Run)
	}

	all := ask(`{"finished":true}`)
	if ids := slices.Sorted(maps.Keys(all.Stories)); !slices.Equal(ids, []string{"S-0001", "S-0002", "S-0003", "S-0004", "S-0099"}) {
		t.Errorf("stories with finished: %v, want every one", ids)
	}
	if ids := slices.Sorted(maps.Keys(all.Plans)); !slices.Equal(ids, []string{"E-0001", "E-0002", "S-0098"}) {
		t.Errorf("plans with finished: %v, want every one", ids)
	}
	if len(all.PastOrchestrators) != 1 || len(all.PastAnalyzers) != 1 || len(all.Roles["orchestrator"].Runs) != 2 || len(all.Roles["planner"].Runs) != 3 {
		t.Errorf("past runs with finished: %d orchestrators, %d analyzers, %d orchestrator role runs, %d planner role runs; want 1, 1, 2, 3",
			len(all.PastOrchestrators), len(all.PastAnalyzers), len(all.Roles["orchestrator"].Runs), len(all.Roles["planner"].Runs))
	}
}

// S-0534: flai serve status prints a line per provider the host holds out of
// budget, and its --json carries the same entries as agent.status.
func TestServeStatusShowsEachHeldProvider(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	cfgPath := filepath.Join(home, "cfg.json")
	useConfig(t, cfgPath)
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	state := `{"claude-code":{"provider":"claude-code","limit":"usage_limit_reached","status":429,"since":"2099-10-28T05:30:00Z","until":"2099-10-28T09:00:00Z","backoff":"5m0s"}}`
	if err := os.WriteFile(filepath.Join(dir, "budget.json"), []byte(state), 0o600); err != nil {
		t.Fatal(err)
	}
	out, _, _ := runIn(t, home, "serve", "status")
	if want := "budget: claude-code held (usage_limit_reached) until 2099-10-28T09:00:00Z\n"; !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	js, _, _ := runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(js), &st); err != nil || len(st.Budget) != 1 || st.Budget[0].Provider != "claude-code" || st.Budget[0].Until != "2099-10-28T09:00:00Z" {
		t.Errorf("json: %v %s", err, js)
	}
}

// S-0506: flai serve mcp bind keeps the address flai host starts each
// project's MCP server on, show prints it, with --json too, and the
// loopback default when it is unset, clear unsets it, and an address with a
// port or one that is no IP address or host name is refused, changing
// nothing.
func TestServeMCPBind(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	run := func(asJSON bool, args ...string) (string, error) {
		t.Helper()
		var out bytes.Buffer
		a := &app{configPath: cfgPath, out: &out, errOut: &bytes.Buffer{}, cwd: t.TempDir(), jsonOut: asJSON}
		c := newServeMCPCmd(a)
		c.SetArgs(args)
		c.SetOut(&out)
		c.SetErr(&bytes.Buffer{})
		c.SilenceUsage, c.SilenceErrors = true, true
		err := c.Execute()
		return out.String(), err
	}
	saved := func() string {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.MCP.Bind
	}
	if out, err := run(false, "show"); err != nil || !strings.Contains(out, "bind: 127.0.0.1 (default)") {
		t.Errorf("unset: %v %q", err, out)
	}
	if js, err := run(true); err != nil || !strings.Contains(js, `"bind": "127.0.0.1"`) || !strings.Contains(js, `"bind_set": false`) {
		t.Errorf("unset, --json: %v %s", err, js)
	}
	if out, err := run(false, "bind", "0.0.0.0"); err != nil || !strings.HasPrefix(out, "bind: 0.0.0.0: ") {
		t.Fatalf("bind: %v %q", err, out)
	}
	if got := saved(); got != "0.0.0.0" {
		t.Errorf("saved: %q", got)
	}
	if js, err := run(true, "show"); err != nil || !strings.Contains(js, `"bind": "0.0.0.0"`) || !strings.Contains(js, `"bind_set": true`) {
		t.Errorf("--json: %v %s", err, js)
	}
	for _, bad := range []struct {
		args []string
		want string
	}{
		{[]string{"bind", "0.0.0.0:4243"}, "names a port"},
		{[]string{"bind", "[::1]"}, "is bracketed"},
		{[]string{"bind", "http://example.com"}, "names a port, or is no IP address"},
		{[]string{"bind", "not a host"}, "is not a host name"},
		{[]string{"bind", " "}, "mcp.bind needs an address"},
		{[]string{"bind"}, "accepts 1 arg"},
	} {
		if _, err := run(false, bad.args...); err == nil || !strings.Contains(err.Error(), bad.want) {
			t.Errorf("%q is refused, saying why: %v", bad.args, err)
		}
	}
	if got := saved(); got != "0.0.0.0" {
		t.Errorf("a refused address changes nothing: %q", got)
	}
	if out, err := run(false, "clear"); err != nil || !strings.Contains(out, "bind: 127.0.0.1 (default)") || saved() != "" {
		t.Errorf("clear: %v %q %q", err, out, saved())
	}
	if data, _ := os.ReadFile(cfgPath); strings.Contains(string(data), `"mcp"`) {
		t.Errorf("a cleared address is not written: %s", data)
	}
}

// S-0587, ADR-0261: flai serve agent set and harness write the project in
// the working directory, or of --project, and leave another project as it
// was; with --all-projects they write the host-wide fallback every project
// without its own reads; clear removes only the project's own; show says
// whose each value is; and outside a project, with neither flag, they refuse
// naming both.
func TestServeAgentSettingsArePerProject(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	a, b := tempProject(t), tempProject(t)
	cfgOf := func(root string) serve.AgentConfig { return (&app{configPath: configOf(t)}).agentConfig(root) }
	if out, errOut, code := runIn(t, a, "serve", "agent", "set", "--name", "builder", "--", "run-a", "{story}"); code != 0 || !strings.Contains(out, `command: "run-a" "{story}" (this project)`) || !strings.Contains(out, "name: builder (this project)") {
		t.Fatalf("set in project a: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, a, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude-a"); code != 0 || !strings.Contains(out, `harness claude-code: "/opt/claude-a" "--permission-mode" "acceptEdits" "--allowedTools" "Bash,mcp__flai" (this project)`) {
		t.Fatalf("harness in project a: %d %s %s", code, out, errOut)
	}
	if got := cfgOf(a); strings.Join(got.Command, " ") != "run-a {story}" || got.Name != "builder" || got.Harnesses["claude-code"].Program != "/opt/claude-a" {
		t.Errorf("project a: %+v", got)
	}
	if got := cfgOf(b); len(got.Command) != 0 || got.Name != "" || got.Harnesses["claude-code"].Program != "claude" {
		t.Errorf("project b is unchanged: %+v", got)
	}
	// --project from b's directory sets a's, not b's
	if _, errOut, code := runIn(t, b, "serve", "agent", "harness", "claude-code", "--project", a, "--program", "/opt/claude-a2"); code != 0 {
		t.Fatalf("--project: %s", errOut)
	}
	if cfgOf(a).Harnesses["claude-code"].Program != "/opt/claude-a2" || cfgOf(b).Harnesses["claude-code"].Program != "claude" {
		t.Errorf("--project names the project written: a %+v, b %+v", cfgOf(a).Harnesses["claude-code"], cfgOf(b).Harnesses["claude-code"])
	}
	// --all-projects: every project without its own reads it, a keeps its own
	if out, errOut, code := runIn(t, b, "serve", "agent", "set", "--all-projects", "--", "run-all"); code != 0 || !strings.Contains(out, `command: "run-all" (every project)`) {
		t.Fatalf("--all-projects: %d %s %s", code, out, errOut)
	}
	runIn(t, b, "serve", "agent", "harness", "claude-code", "--all-projects", "--program", "/usr/local/bin/claude")
	if got := cfgOf(b); strings.Join(got.Command, " ") != "run-all" || got.Harnesses["claude-code"].Program != "/usr/local/bin/claude" {
		t.Errorf("project b reads every project's: %+v", got)
	}
	if got := cfgOf(a); strings.Join(got.Command, " ") != "run-a {story}" || got.Harnesses["claude-code"].Program != "/opt/claude-a2" {
		t.Errorf("project a keeps its own: %+v", got)
	}
	js, _, _ := runIn(t, b, "serve", "agent", "show", "--json")
	if !strings.Contains(js, `"command_from": "all-projects"`) || !strings.Contains(js, `"from": "all-projects"`) {
		t.Errorf("show --json says whose each value is: %s", js)
	}
	// clear and --reset in a remove a's own, so that it reads every project's
	runIn(t, a, "serve", "agent", "clear")
	runIn(t, a, "serve", "agent", "harness", "claude-code", "--reset")
	if got := cfgOf(a); strings.Join(got.Command, " ") != "run-all" || got.Name != "" || got.Harnesses["claude-code"].Program != "/usr/local/bin/claude" {
		t.Errorf("a reads every project's again: %+v", got)
	}
	if cfg, _, _ := (&app{configPath: configOf(t)}).loadConfig(); cfg.Agent.Projects != nil {
		t.Errorf("a's entry, now empty, is removed: %+v", cfg.Agent.Projects)
	}
	// outside a project, neither flag is refused; --project and --all-projects together too
	outside := t.TempDir()
	for _, args := range [][]string{
		{"serve", "agent", "set", "--", "x"},
		{"serve", "agent", "clear"},
		{"serve", "agent", "harness", "claude-code", "--program", "/x"},
	} {
		if _, errOut, code := runIn(t, outside, args...); code == 0 || !strings.Contains(errOut, "--project <root>") || !strings.Contains(errOut, "--all-projects") {
			t.Errorf("%q outside a project: %d %s", args, code, errOut)
		}
	}
	if _, errOut, code := runIn(t, a, "serve", "agent", "set", "--project", a, "--all-projects", "--", "x"); code == 0 || !strings.Contains(errOut, "not both") {
		t.Errorf("both flags: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, a, "serve", "agent", "set", "--project", outside, "--", "x"); code == 0 || !strings.Contains(errOut, "name the main checkout of a system-flow project") {
		t.Errorf("--project that is no project: %d %s", code, errOut)
	}
	// the restart limit stays host-wide, and needs no project
	if _, errOut, code := runIn(t, outside, "serve", "agent", "set", "--auto-restarts", "3"); code != 0 || cfgOf(a).AutoRestarts != 3 || cfgOf(b).AutoRestarts != 3 {
		t.Errorf("--auto-restarts is host-wide: %d %s", code, errOut)
	}
}

// onTheHostsCommand makes the project at root's default agent the host's
// command, so a start with no command set is refused rather than falling
// back to flai's built-in claude-code, which would run the host's real
// claude (ADR-0261).
func onTheHostsCommand(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, "system-flow.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, []byte("agent:\n  harness: command\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
}

// holdBudgets writes dir's budget.json with a state held for a day under each
// of keys, as a run that ended on its usage limit leaves it.
func holdBudgets(t *testing.T, dir serve.Dir, keys ...string) {
	t.Helper()
	states := map[string]serve.Budget{}
	for _, k := range keys {
		states[k] = serve.Budget{Provider: k, Limit: "usage_limit_reached", Status: 429, Since: time.Now().UTC(), Until: time.Now().UTC().Add(24 * time.Hour), Backoff: "5m0s"}
	}
	data, err := json.Marshal(states)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(string(dir), "budget.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// S-0720: a change of a harness's program, for a project or for every
// project, lifts the hold of the old program's budget state, and of a state
// kept under the harness's name alone from before S-0720: it is cleared,
// journaled under serve.budget with why, said, and sent as a notice. A state
// of another program on the harness stays, and a change of anything else,
// such as --guard, clears nothing.
func TestServeAgentHarnessProgramChangeLiftsTheOldBudget(t *testing.T) {
	t.Parallel()
	dir := budgetFolder(t, "")
	root := tempProject(t)
	notices := filepath.Join(t.TempDir(), "notices")
	if _, errOut, code := runIn(t, root, "serve", "notify", "set", "--command", "cat >> '"+notices+"'; echo >> '"+notices+"'"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "rmhclod"); code != 0 {
		t.Fatal(errOut)
	}
	holdBudgets(t, dir, "claude-code rmhclod", "claude-code other")
	out, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "meclod")
	if code != 0 {
		t.Fatalf("program change: %d %s", code, errOut)
	}
	why := "flai serve agent harness changed the claude-code harness's program from rmhclod to meclod for " + root
	if !strings.Contains(out, "cleared the budget state of claude-code rmhclod: "+why) {
		t.Errorf("the clear is said: %s", out)
	}
	if got := budgetKeys(t, dir); !slices.Equal(got, []string{"claude-code other"}) {
		t.Errorf("the old program's state is cleared, another program's stays: %q", got)
	}
	entries := budgetJournal(t, dir)
	if len(entries) != 1 || !strings.Contains(entries[0].Detail, why) || !strings.Contains(entries[0].Detail, "claude-code rmhclod") || entries[0].Root != root {
		t.Errorf("journaled under serve.budget with why: %+v", entries)
	}
	if sent, err := os.ReadFile(notices); err != nil || !strings.Contains(string(sent), `"state":"cleared"`) || !strings.Contains(string(sent), `"provider":"claude-code rmhclod"`) || !strings.Contains(string(sent), why) {
		t.Errorf("sent as a notice: %v %s", err, sent)
	}
	// a change of neither program, arguments, nor login clears nothing
	holdBudgets(t, dir, "claude-code meclod")
	if out, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--guard", "none"); code != 0 || strings.Contains(out, "cleared the budget state") {
		t.Fatalf("--guard none: %d %s %s", code, out, errOut)
	}
	if got := budgetKeys(t, dir); !slices.Equal(got, []string{"claude-code meclod"}) {
		t.Errorf("--guard none leaves it held: %q", got)
	}
	// the arguments are part of what a change lifts
	js, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--json", "--", "--permission-mode", "acceptEdits")
	var shown struct {
		Cleared []budgetLift `json:"budget_cleared"`
	}
	if err := json.Unmarshal([]byte(js), &shown); err != nil || code != 0 {
		t.Fatalf("--json: %d %v %s %s", code, err, js, errOut)
	}
	if len(shown.Cleared) != 1 || shown.Cleared[0].Key != "claude-code meclod" || !strings.Contains(shown.Cleared[0].Why, "arguments from") {
		t.Errorf("--json says the clear: %+v", shown.Cleared)
	}
	if got := budgetKeys(t, dir); len(got) != 0 {
		t.Errorf("the next start is not held: %q", got)
	}
	// a state of the old shape, under the harness's name alone, is cleared
	// by a change of every project's program, which is said
	holdBudgets(t, dir, "claude-code", "opencode opencode")
	if out, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--all-projects", "--program", "/usr/local/bin/claude"); code != 0 ||
		!strings.Contains(out, "cleared the budget state of claude-code: flai serve agent harness changed the claude-code harness's program from claude to /usr/local/bin/claude for every project") {
		t.Fatalf("--all-projects: %d %s %s", code, out, errOut)
	}
	if got := budgetKeys(t, dir); !slices.Equal(got, []string{"opencode opencode"}) {
		t.Errorf("the old-shape state is cleared, another harness's stays: %q", got)
	}
}

// S-0720: a change of the container harness's login lifts the hold of the
// old login's budget state. Not parallel: it sets HOME for ~/.
func TestServeAgentHarnessLoginChangeLiftsTheOldBudget(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	rmh, me := filepath.Join(home, ".claude-rmh"), filepath.Join(home, ".claude-me")
	for _, d := range []string{rmh, me} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	dir := budgetFolder(t, "")
	root := tempProject(t)
	if _, errOut, code := runIn(t, root, "serve", "agent", "harness", "container", "--login", "~/.claude-rmh"); code != 0 {
		t.Fatal(errOut)
	}
	old := "container docker login " + rmh
	holdBudgets(t, dir, old, "container docker login "+me)
	// the network is not part of what a start runs as
	if _, errOut, code := runIn(t, root, "serve", "agent", "harness", "container", "--network", "flai-agents"); code != 0 {
		t.Fatal(errOut)
	}
	if got := budgetKeys(t, dir); len(got) != 2 {
		t.Errorf("--network leaves both held: %q", got)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "harness", "container", "--login", "~/.claude-me")
	if code != 0 || !strings.Contains(out, "cleared the budget state of "+old+": flai serve agent harness changed the container harness's login from "+rmh+" to "+me) {
		t.Fatalf("login change: %d %s %s", code, out, errOut)
	}
	if got := budgetKeys(t, dir); !slices.Equal(got, []string{"container docker login " + me}) {
		t.Errorf("the old login's state is cleared: %q", got)
	}
	if entries := budgetJournal(t, dir); len(entries) != 1 || !strings.Contains(entries[0].Detail, "login from "+rmh) {
		t.Errorf("journaled: %+v", entries)
	}
}

// S-0720: flai serve agent set lifts the hold of the old command's budget
// state when it changes the command's program or arguments, and only then.
func TestServeAgentSetLiftsTheOldCommandsBudget(t *testing.T) {
	t.Parallel()
	dir := budgetFolder(t, "")
	root := tempProject(t)
	if _, errOut, code := runIn(t, root, "serve", "agent", "set", "--", "run-a", "{story}"); code != 0 {
		t.Fatal(errOut)
	}
	holdBudgets(t, dir, "command run-a")
	if out, errOut, code := runIn(t, root, "serve", "agent", "set", "--name", "builder"); code != 0 || strings.Contains(out, "cleared the budget state") {
		t.Fatalf("--name: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, root, "serve", "agent", "set", "--", "run-a", "{story}"); code != 0 || strings.Contains(out, "cleared the budget state") {
		t.Fatalf("the same command: %d %s %s", code, out, errOut)
	}
	if got := budgetKeys(t, dir); !slices.Equal(got, []string{"command run-a"}) {
		t.Errorf("neither changes it, so it stays held: %q", got)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "set", "--", "run-b", "{story}")
	if code != 0 || !strings.Contains(out, `cleared the budget state of command run-a: flai serve agent set changed the command harness's program from run-a to run-b for `+root) {
		t.Fatalf("program change: %d %s %s", code, out, errOut)
	}
	if got := budgetKeys(t, dir); len(got) != 0 {
		t.Errorf("the old command's state is cleared: %q", got)
	}
	if entries := budgetJournal(t, dir); len(entries) != 1 {
		t.Errorf("journaled: %+v", entries)
	}
}
