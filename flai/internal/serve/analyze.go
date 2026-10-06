package serve

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/analysis"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The analyzer, started for a project on the operator's word or on a
// schedule (S-0223).
//
// flai analyze, the Analyze button through analyze.run, which runs it, the
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
// cost. It is no story's agent: the in-progress limit does not count it, and
// it holds no story back.

// Analyze starts the analyzer in project e, looking for focus, one of
// harness.Focuses, or for all of them when focus is empty, and returns its
// run as recorded, with trigger as what started it: TriggerAsked for the
// operator's asking, from a shell, the dashboard, or their own agent, and
// the schedule's own for a scheduled run. It is refused, as a *Refused,
// while the analyze host action is off for the project, for a focus the
// analyzer does not take, while an analyzer runs for the project, naming
// it, and when nothing can start it. It does not wait for the analyzer: the
// serving flai settles the run once the process is gone.
func Analyze(ctx context.Context, o Options, e Entry, focus, trigger string) (*AgentRun, error) {
	as, err := analyzeCheck(o, e, focus)
	if err != nil {
		return nil, err
	}
	l := newLauncher(o, e)
	l.handOver = true
	l.analyze(ctx, as.cfg, as.agent, focus, trigger)
	run := o.Dir.AgentStates()[e.Root].Analyzer
	if run == nil {
		return nil, fmt.Errorf("the analyzer run was not recorded in %s", o.Dir.agents())
	}
	if run.Error != "" {
		return run, errors.New(run.Why)
	}
	return run, nil
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
	agent := repo.Manifest.AnalysisAgent()
	if (agent == nil || agent.Harness == "") && cfg.host(harness.Command).Program == "" {
		return analyzeStart{}, refused("the analyzer's agent names no harness (analysis.agent or agent in system-flow.yaml), and no command is set on the host (flai serve agent set -- <program> [args...])")
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
// journals it, and says whether it started.
func (l *launcher) analyze(ctx context.Context, cfg AgentConfig, agent *manifest.Agent, focus, trigger string) bool {
	now := l.now().UTC()
	run := &AgentRun{Agent: workitem.ActivityAnalyzer, Focus: focus, Started: now.Format(time.RFC3339), Session: newSession(), Trigger: trigger}
	if run.Focus == "" {
		run.Focus = harness.AllFocus
	}
	if agent != nil {
		run.Model = agent.Model
	}
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionAnalyze, Method: "serve.analyze", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) bool {
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.told(run)
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("the analyzer, focus %s, could not be started: %s", run.Focus, run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("analyzer could not be started", "focus", run.Focus, "harness", run.Harness, "err", run.Error)
		return false
	}
	name, adapter, err := harness.For(agent, cfg.host(harness.Command).Program != "")
	run.Harness = name
	if err != nil {
		return fail(err)
	}
	spec, err := adapter.Start(harness.Request{Role: conventions.RoleAnalyze, Focus: focus, Root: l.entry.Root, Project: l.entry.Key, Agent: agent,
		Name: run.Agent, Flai: cfg.Flai, Session: run.Session}, cfg.host(name))
	if err != nil {
		return fail(err)
	}
	cmd, out, err := l.spawn(ctx, run, spec.Argv, workitem.ActivityAnalyzer, now, spec.Env)
	if err != nil {
		return fail(err)
	}
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) to analyze, focus %s, as %s (pid %d); log %s", run.Command, run.Harness, run.Focus, run.Agent, run.PID, run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("analyzer started", "focus", run.Focus, "trigger", run.Trigger, "harness", run.Harness, "command", run.Command, "pid", run.PID)
	l.await(cmd, out, run.PID, func(code int) { l.analyzeEnded(run, &code) })
	return true
}

// analyzeEnded records an analyzer run as ended, with its exit code when it
// was seen, and the report it wrote: failed on a failed exit, and on an exit
// that left no report; worked otherwise, an exit nobody saw included. It
// logs the run's activity in the analyzer's activity document, with what
// started it (ADR-0084) and a summary that names the report, or says there
// is none. A report that cannot be looked for, and an activity that cannot
// be logged, are warned of, and the run stays recorded as it ended.
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
	switch {
	case exit != nil && *exit != 0:
		ended.Outcome, ended.Why = OutcomeFailed, fmt.Sprintf("ended (%s)", code)
	case report == "":
		ended.Outcome, ended.Why = OutcomeFailed, fmt.Sprintf("ended (%s) without writing a report under %s", code, dir)
	default:
		ended.Outcome, ended.Why = OutcomeWorked, ""
	}
	l.ended(&ended)
	args := []any{"pid", run.PID, "focus", run.Focus, "outcome", ended.Outcome, "report", report}
	if exit != nil {
		args = append(args, "exit", *exit)
	}
	l.log("analyzer ended", args...)
	say := func(final string) string { return reportSaid(final, report, dir) }
	if _, err := logRunEndSaying(l.dir, l.entry.Root, l.entry.Key, workitem.ActivityAnalyzer, run.Trigger, say, nil); err != nil {
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
