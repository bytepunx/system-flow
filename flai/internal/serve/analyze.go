package serve

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/analysis"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The analyzer, started for a project on the operator's word or on a
// schedule (S-0223).
//
// flai analyze, the host API's analyze.run, which runs it, the
// MCP tool analyze, and analysis.schedule all start it through Analyze: in
// the project's main checkout, with the project's analysis agent
// (analysis.agent over agent) and the harness's analyzer request, as
// analyzer, with FLAI_ROLE=analyze and a focus or none. It looks for what its
// focus names, writes one report under design/analysis, and ends. One runs
// per project at a time. Like the planner's, its run is recorded in
// serve/agents.json, under analyzer, the newest alone, and its log beside
// the others as the analyzer's (Dir.ActivityLogs). The serving flai settles
// the run once its process is gone, and logs its activity in
// wip/agents/analyzer.md, naming the report it wrote, the newest file under
// design/analysis changed since it started, beside the run's seconds and
// cost, and naming the issues that name the report, which share the cost
// (S-0227). It is no story's agent: the in-progress limit does not count it,
// and it holds no story back. The agent of a formula's instance starts it
// too, for a role step of its formula that names it, through RunRole and
// analyzeStep, with Analyze's checks; that run's start and end are agent
// events on the instance (S-0424).

// Analyze starts the analyzer in project e, looking for focus, one of
// harness.Focuses, or for all of them when focus is empty, and returns its
// run as recorded, with trigger as what started it: TriggerAsked for the
// operator's asking, from a shell, the dashboard, or their own agent, and
// the schedule's own for a scheduled run. It is refused, as a *Refused,
// while the analyze host action is off for the project, for a focus the
// analyzer does not take, while an analyzer runs for the project, naming
// it, when nothing can start it, and while the budget state of the
// analyzer's provider holds its start, naming the provider and until, with
// nothing recorded (S-0534). It does not wait for the analyzer: the serving
// flai settles the run once the process is gone.
func Analyze(ctx context.Context, o Options, e Entry, focus, trigger string) (*AgentRun, error) {
	return analyzeFor(ctx, o, e, focus, trigger, "")
}

// analyzeFor is Analyze for instance, the formula's instance whose agent
// started the run for a role step, or "" for any other start.
func analyzeFor(ctx context.Context, o Options, e Entry, focus, trigger, instance string) (*AgentRun, error) {
	as, err := analyzeCheck(o, e, focus)
	if err != nil {
		return nil, err
	}
	l := newLauncher(o, e)
	l.handOver = true
	if err := l.budgetRefusal(as.cfg, as.agent); err != nil {
		return nil, err
	}
	if _, held := l.analyzeOn(ctx, as.cfg, as.agent, focus, trigger, instance); held != nil {
		return nil, l.budgetRefused(*held)
	}
	run := o.Dir.AgentStates()[e.Root].Analyzer
	if run == nil {
		return nil, fmt.Errorf("the analyzer run was not recorded in %s", o.Dir.agents())
	}
	if run.Error != "" {
		return run, errors.New(run.Why)
	}
	return run, nil
}

// analyzeStep starts the analyzer for a role step of a formula that names it,
// as RunRole does when the agent of story, the formula's instance, asks for
// it on item with trigger, FormulaTrigger's (S-0424, ADR-0239 §2). It goes
// through the analyzer's own start, with Analyze's checks, under the analyze
// host action, and records story as the run's instance, so that its start
// and its end are agent events on the instance. Its focus is the last input
// named focus among inputs, else the focus input the instance's Notes say
// flai formula run took (instanceFocus), empty for all three. An item that
// is not story is refused.
func analyzeStep(ctx context.Context, o Options, e Entry, story, item string, inputs []harness.Input, trigger string) (*AgentRun, error) {
	if workitem.CanonicalID(item) != story {
		return nil, refused("the analyzer runs for a role step of a formula on its instance %s alone, and was asked for %s: give the item %s", story, onItem(item), story)
	}
	focus, given := "", false
	for _, in := range inputs {
		if in.Name == "focus" {
			focus, given = strings.TrimSpace(in.Value), true
		}
	}
	if !given {
		repo, err := workitem.Open(e.Root)
		if err != nil {
			return nil, err
		}
		it, err := repo.Get(story)
		if err != nil {
			return nil, fmt.Errorf("reading %s, the instance whose role step asked for the analyzer: %w", story, err)
		}
		focus = instanceFocus(it.Body)
	}
	return analyzeFor(ctx, o, e, focus, trigger, story)
}

// madeBy is the line of an instance's Notes that flai formula run writes,
// naming the inputs it took, and focusTaken the focus input among them.
var (
	madeBy     = regexp.MustCompile("(?m)^Made by `flai formula run [^`]+` from .*$")
	focusTaken = regexp.MustCompile("`focus=([^`]*)`")
)

// instanceFocus is the focus input an instance's body says flai formula run
// took, on the line of its Notes that names the inputs, or "" when the line
// names none.
func instanceFocus(body string) string {
	line := madeBy.FindString(body)
	_, took, ok := strings.Cut(line, " with the inputs ")
	if !ok {
		return ""
	}
	if m := focusTaken.FindStringSubmatch(took); m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// analyzeStart is what the analyzer is started with once the checks pass:
// the operator's say and the analyzer's agent.
type analyzeStart struct {
	cfg   AgentConfig
	agent *manifest.Agent
}

// analyzeCheck makes the checks before the analyzer is started in project e,
// for Analyze and for a start flai serve makes itself, and returns what it is
// started with. It is refused, as a *Refused, while the analyze host action
// is off, for a focus that is not empty and not one of harness.Focuses, while
// an analyzer runs for the project, and when nothing can start it; a project
// that cannot be read is an error.
func analyzeCheck(o Options, e Entry, focus string) (analyzeStart, error) {
	var cfg AgentConfig
	if o.Agent != nil {
		cfg = o.Agent(e.Root)
	}
	if !cfg.Analyze {
		return analyzeStart{}, refused("the analyze host action is off for this project: flai serve enable analyze")
	}
	if focus != "" && !slices.Contains(harness.Focuses, focus) {
		return analyzeStart{}, refused("the analyzer takes the focus %s, or none for all of them, and %q is none of them", strings.Join(harness.Focuses, ", "), focus)
	}
	if run := o.Dir.AgentStates()[e.Root].Analyzer; run.running() {
		return analyzeStart{}, refused("the analyzer is already running for this project (focus %s, %s, pid %d, started %s; log %s); one analyzer runs per project at a time", run.Focus, startedBy(run), run.PID, run.Started, run.Log)
	}
	repo, err := workitem.Open(e.Root)
	if err != nil {
		return analyzeStart{}, err
	}
	// flai's built-in default when nothing names a harness and no command is
	// set (ADR-0261)
	agent := cfg.startAgent(repo.Manifest.AnalysisAgent())
	if err := cfg.refuseUnstartable("the analyzer", agent); err != nil {
		return analyzeStart{}, err
	}
	return analyzeStart{cfg: cfg, agent: agent}, nil
}

// startedBy says what started run, for a refusal that names it.
func startedBy(run *AgentRun) string {
	if run.Trigger == "" {
		return "trigger not recorded"
	}
	return "started by " + run.Trigger
}

// analyze starts the analyzer with agent, looking for focus, or for all of
// it when focus is empty, records the run with what started it, trigger,
// journals it, and says whether it started. Its start, or the start it could
// not make, is an agent event on the analyzer's role (S-0532), or on the
// formula's instance it was started for (analyzeOn). A start the
// budget state of its provider holds is not made and nothing is recorded:
// held is that state, nil when it was not held (spawnGatedWith, S-0534).
func (l *launcher) analyze(ctx context.Context, cfg AgentConfig, agent *manifest.Agent, focus, trigger string) (started bool, held *Budget) {
	return l.analyzeOn(ctx, cfg, agent, focus, trigger, "")
}

// analyzeOn is analyze for instance, the formula's instance whose agent
// started the run for a role step, recorded as the run's Instance, so that
// its start and its end are agent events on the instance (S-0424); "" for
// any other start.
func (l *launcher) analyzeOn(ctx context.Context, cfg AgentConfig, agent *manifest.Agent, focus, trigger, instance string) (started bool, held *Budget) {
	now := l.now().UTC()
	run := &AgentRun{Agent: workitem.ActivityAnalyzer, Focus: focus, Started: now.Format(time.RFC3339), Session: newSession(), Trigger: trigger, Instance: instance}
	on := ""
	if instance != "" {
		on = ", on " + instance
	}
	if run.Focus == "" {
		run.Focus = harness.AllFocus
	}
	if agent != nil {
		run.Model = agent.Model
	}
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionAnalyze, Method: "serve.analyze", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) (bool, *Budget) {
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.told(run)
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("the analyzer, focus %s%s, could not be started: %s", run.Focus, on, run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("analyzer could not be started", "focus", run.Focus, "harness", run.Harness, "err", run.Error)
		appendAgentEvent(l.entry.Root, run, inbox.AgentFailed, run.Why, l.warn)
		return false, nil
	}
	name, adapter, err := harness.For(agent, cfg.host(harness.Command).Program != "")
	run.Harness = name
	if err != nil {
		return fail(err)
	}
	spec, keys, err := l.startOf(cfg, run, adapter, name, harness.Request{Role: conventions.RoleAnalyze, Focus: focus, Root: l.entry.Root, Project: l.entry.Key, Agent: agent,
		Name: run.Agent, Flai: cfg.Flai, Session: run.Session})
	if w := (*imageWait)(nil); errors.As(err, &w) {
		// held for its image job, with what started it, and made again once
		// the job ends (holdImage, S-0751): no failure is recorded
		l.holdImage(ctx, run, w, func(ctx context.Context) bool {
			_, held := l.analyzeOn(ctx, cfg, agent, focus, trigger, instance)
			return held != nil
		})
		return true, nil
	}
	if err != nil {
		return fail(err)
	}
	cmd, out, held, err := l.spawnGatedWith(ctx, run, spec.Argv, workitem.ActivityAnalyzer, now, append(spec.Env, keys...), spec.Stdin)
	if held != nil {
		return false, held
	}
	if err != nil {
		return fail(err)
	}
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s to analyze, focus %s%s, as %s (pid %d); log %s", run.startedSaid(), run.Focus, on, run.Agent, run.PID, run.Log)
	if instance != "" {
		entry.Detail += ", started by " + run.Trigger
	}
	if l.record != nil {
		l.record(entry)
	}
	l.log("analyzer started", "focus", run.Focus, "trigger", run.Trigger, "instance", run.Instance, "harness", run.Harness, "command", run.Command, "pid", run.PID)
	appendAgentEvent(l.entry.Root, run, inbox.AgentStarted, triggerWhy(run)+", focus "+run.Focus, l.warn)
	l.await(cmd, out, run.PID, func(code int) { l.analyzeEnded(run, &code) })
	return true, nil
}

// analyzeEnded records an analyzer run as ended, with its exit code when it
// was seen, and the report it wrote: failed on a failed exit, and on an exit
// that left no report; worked otherwise, an exit nobody saw included. It
// logs the run's activity in the analyzer's activity document, with what
// started it (ADR-0084), a summary that names the report, or says there is
// none, and as its items the issues that name the report, so that its cost
// is split evenly between them (S-0227). A report, or issues naming it, that
// cannot be looked for, and an activity that cannot be logged, are warned
// of: the activity is logged naming no issue, and the run stays recorded as
// it ended. The end is an agent event on the analyzer's role, or on the
// formula's instance the run was started for (S-0424), appended once, by the
// end that recorded the run as ended first (S-0532). A run whose log
// ends on the harness's usage limit is failed on it, with why naming the
// limit, whether its exit was seen or not, and marked so (limitEnded); every
// end feeds its provider's budget state once (budgetEnd, S-0534).
func (l *launcher) analyzeEnded(run *AgentRun, exit *int) {
	ended := *run
	ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), exit
	code := "an exit code nobody saw"
	if exit != nil {
		code = fmt.Sprintf("exit %d", *exit)
	}
	dir, report, err := newestReport(l.entry.Root, run.Started)
	if err != nil {
		l.warn("analyzer's report not looked for", "err", err)
	}
	ended.Report = report
	lim, _ := runLimit(run, l.warn)
	ended.limitEnded(lim)
	switch {
	case lim != nil:
		ended.Outcome, ended.Why = OutcomeFailed, limitWhy(lim)
	case exit != nil && *exit != 0:
		ended.Outcome, ended.Why = OutcomeFailed, fmt.Sprintf("ended (%s)", code)
	case report == "":
		ended.Outcome, ended.Why = OutcomeFailed, fmt.Sprintf("ended (%s) without writing a report under %s", code, dir)
	default:
		ended.Outcome, ended.Why = OutcomeWorked, ""
	}
	first := l.ended(&ended)
	if first {
		appendAgentEnd(l.entry.Root, &ended, l.warn)
	}
	l.budgetEnd(&ended, lim, first)
	args := []any{"pid", run.PID, "focus", run.Focus, "outcome", ended.Outcome, "report", report}
	if exit != nil {
		args = append(args, "exit", *exit)
	}
	l.log("analyzer ended", args...)
	var named []string
	if report != "" {
		repo, err := workitem.Open(l.entry.Root)
		if err == nil {
			named, err = issues.NamingReport(repo, report)
		}
		if err != nil {
			named = nil
			l.warn("issues naming the analyzer's report not looked for", "report", report, "err", err)
		}
	}
	say := func(final string) string { return reportSaid(final, report, dir) }
	if _, err := logRunEndSaying(l.dir, l.entry.Root, l.entry.Key, workitem.ActivityAnalyzer, run.Log, run.Trigger, say, named); err != nil {
		l.warn("analyzer activity not logged", "err", err)
	}
}

// reportSaid is an analyzer run's activity summary: its final line, which
// its prompt has name the report, with the report it wrote added when the
// line does not name it, or with dir, where it wrote none, said.
func reportSaid(final, report, dir string) string {
	switch {
	case report == "":
		return fmt.Sprintf("%s (no report written under %s)", final, dir)
	case strings.Contains(final, report):
		return final
	}
	return fmt.Sprintf("%s (report %s)", final, report)
}

// newestReport is the report an analyzer run that started at started wrote
// in the project at root: of the Markdown files directly in the reports
// folder, the folder's README.md, its index, aside, the one changed last,
// if it changed in the second the run started or after. It returns the
// folder and the report, both relative to root, and "" for the report when
// none changed since, or there is no folder.
func newestReport(root, started string) (dir, report string, err error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return "design/analysis", "", err
	}
	dir = analysis.Dir(repo.Manifest)
	since, err := time.Parse(time.RFC3339, started)
	if err != nil {
		return dir, "", err
	}
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if errors.Is(err, os.ErrNotExist) {
		return dir, "", nil
	}
	if err != nil {
		return dir, "", err
	}
	var newest time.Time
	for _, de := range entries {
		name := de.Name()
		if !de.Type().IsRegular() || !strings.HasSuffix(name, ".md") || name == analysis.Index {
			continue
		}
		info, err := de.Info()
		if err != nil {
			continue
		}
		if at := info.ModTime(); !at.Before(since) && (report == "" || at.After(newest)) {
			newest, report = at, path.Join(dir, name)
		}
	}
	return dir, report, nil
}
