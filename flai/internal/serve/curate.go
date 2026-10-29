package serve

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The curator, started for a project on the operator's word (S-0431,
// ADR-0196), and on its own after acceptances and on digest.schedule
// (S-0432, ADR-0199; curate_schedule.go).
//
// flai curate, the host API's curate.run, which runs it, and the MCP tool
// curate all start it through Curate: in the project's main checkout, with
// the project's digest agent (digest.agent over agent) and the harness's
// curator request, as curator, with FLAI_ROLE=curate and the run's since. It
// takes in what the archive and the threads gained since then as digest
// entries, written through flai digest new and edit alone, and ends. One
// runs per project at a time. Like the analyzer's, its run is recorded in
// serve/agents.json, under curator, the newest alone, marked with its role so
// that it is never taken for the orchestrator's, and its log beside the
// others as the curator's (Dir.ActivityLogs). The serving flai settles the
// run once its process is gone, and logs its activity in wip/agents/curator.md,
// naming the entries whose files the run wrote or changed, beside the run's
// seconds and cost. It charges no item: digest entries are not work items. It
// is no story's agent: the in-progress limit does not count it, and it holds
// no story back. Each start is journalled with what started it, its trigger.

// sinceForms are the forms a curator's since takes, as flai digest harvest
// --since reads it: a day in UTC, or an RFC 3339 time.
var sinceForms = []string{time.DateOnly, time.RFC3339}

// Curate starts the curator in project e, taking in what the archive and the
// threads gained since since, a YYYY-MM-DD day in UTC or an RFC 3339 time,
// or, when since is empty, since the curator's last activity, or everything
// on its first run. It returns its run as recorded, its since the one it was
// started with, with trigger as what started it: TriggerAsked for the
// operator's asking, from a shell, the dashboard, or their own agent. It is
// refused, as a *Refused, while the curate host action is off for the
// project, for a since in neither form, while a curator runs for the
// project, naming it, and when nothing can start it. It does not wait for
// the curator: the serving flai settles the run once the process is gone.
func Curate(ctx context.Context, o Options, e Entry, since, trigger string) (*AgentRun, error) {
	cs, err := curateCheck(o, e, since)
	if err != nil {
		return nil, err
	}
	l := newLauncher(o, e)
	l.handOver = true
	l.curate(ctx, cs.cfg, cs.agent, cs.since, trigger, "")
	run := o.Dir.AgentStates()[e.Root].Curator
	if run == nil {
		return nil, fmt.Errorf("the curator run was not recorded in %s", o.Dir.agents())
	}
	if run.Error != "" {
		return run, errors.New(run.Why)
	}
	return run, nil
}

// curateStart is what the curator is started with once the checks pass: the
// operator's say, the curator's agent, and the run's since.
type curateStart struct {
	cfg   AgentConfig
	agent *manifest.Agent
	since string
}

// curateCheck makes the checks before the curator is started in project e,
// and returns what it is started with: since, or when it is empty the time
// of the curator's last activity, or none on its first run. It is refused,
// as a *Refused, while the curate host action is off, for a since that is
// neither a YYYY-MM-DD day nor an RFC 3339 time, while a curator runs for
// the project, and when nothing can start it; a project, or a curator's
// activity document, that cannot be read is an error.
func curateCheck(o Options, e Entry, since string) (curateStart, error) {
	var cfg AgentConfig
	if o.Agent != nil {
		cfg = o.Agent(e.Root)
	}
	if !cfg.Curate {
		return curateStart{}, refused("the curate host action is off for this project: flai serve enable curate")
	}
	if since != "" && !sinceReads(since) {
		return curateStart{}, refused("the curator takes a since as a YYYY-MM-DD day in UTC or an RFC 3339 time, such as 2026-10-01 or 2026-10-01T09:00:00Z, and %q is neither", since)
	}
	if run := o.Dir.AgentStates()[e.Root].Curator; run.running() {
		return curateStart{}, refused("the curator is already running for this project (%s, %s, pid %d, started %s; log %s); one curator runs per project at a time", sinceSaid(run.Since), startedBy(run), run.PID, run.Started, run.Log)
	}
	repo, err := workitem.Open(e.Root)
	if err != nil {
		return curateStart{}, err
	}
	agent := repo.Manifest.DigestAgent()
	if (agent == nil || agent.Harness == "") && cfg.host(harness.Command).Program == "" {
		return curateStart{}, refused("the curator's agent names no harness (digest.agent or agent in system-flow.yaml), and no command is set on the host (flai serve agent set -- <program> [args...])")
	}
	if since == "" {
		if since, err = lastCurated(repo); err != nil {
			return curateStart{}, err
		}
	}
	return curateStart{cfg: cfg, agent: agent, since: since}, nil
}

// sinceReads says whether since is in one of sinceForms.
func sinceReads(since string) bool {
	return slices.ContainsFunc(sinceForms, func(form string) bool {
		_, err := time.Parse(form, since)
		return err == nil
	})
}

// sinceSaid says what a curator run takes in, for its journal entry and a
// refusal that names it.
func sinceSaid(since string) string {
	if since == "" {
		return "everything, its first run"
	}
	return "since " + since
}

// lastCurated is the time of the newest entry in the curator's activity
// document, as an RFC 3339 time in UTC, or "" when it has none.
func lastCurated(repo *workitem.Repo) (string, error) {
	doc, err := repo.Activity(workitem.ActivityCurator)
	if err != nil {
		return "", fmt.Errorf("cannot read when the curator last ran: %w", err)
	}
	var last time.Time
	for _, en := range doc.Entries {
		if en.At.After(last) {
			last = en.At
		}
	}
	if last.IsZero() {
		return "", nil
	}
	return last.UTC().Format(time.RFC3339), nil
}

// curate starts the curator with agent, taking in from since, or everything
// when it is empty, records the run, marked with its role, with what started
// it, trigger, journals it naming the trigger, and says whether it started.
// With resetAt, when the usage limit the run before it ended on reset, the
// run starts that one again: it is marked a retry, and its journal entry says
// so (ADR-0214); resetAt is empty for every other start.
func (l *launcher) curate(ctx context.Context, cfg AgentConfig, agent *manifest.Agent, since, trigger, resetAt string) bool {
	now := l.now().UTC()
	run := &AgentRun{Agent: workitem.ActivityCurator, Role: conventions.RoleCurate, Since: since, Started: now.Format(time.RFC3339), Session: newSession(), Trigger: trigger, Retry: resetAt != ""}
	by := startedBy(run)
	if run.Retry {
		by += "; started again after the usage limit reset at " + resetAt
	}
	if agent != nil {
		run.Model = agent.Model
	}
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionCurate, Method: "serve.curate", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) bool {
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.told(run)
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("the curator, %s, %s, could not be started: %s", sinceSaid(since), by, run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("curator could not be started", "since", since, "trigger", run.Trigger, "harness", run.Harness, "err", run.Error)
		return false
	}
	name, adapter, err := harness.For(agent, cfg.host(harness.Command).Program != "")
	run.Harness = name
	if err != nil {
		return fail(err)
	}
	spec, keys, err := l.startOf(cfg, run, adapter, name, harness.Request{Role: conventions.RoleCurate, Since: since, Root: l.entry.Root, Project: l.entry.Key, Agent: agent,
		Name: run.Agent, Flai: cfg.Flai, Session: run.Session})
	if err != nil {
		return fail(err)
	}
	cmd, out, err := l.spawn(ctx, run, spec.Argv, workitem.ActivityCurator, now, append(spec.Env, keys...))
	if err != nil {
		return fail(err)
	}
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) to curate, %s, as %s (pid %d), %s; log %s", run.Command, run.Harness, sinceSaid(since), run.Agent, run.PID, by, run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("curator started", "since", since, "trigger", run.Trigger, "harness", run.Harness, "command", run.Command, "pid", run.PID)
	l.await(cmd, out, run.PID, func(code int) { l.curateEnded(run, &code) })
	return true
}

// curateEnded records a curator run as ended, with its exit code when it was
// seen, and the digest entries it wrote or changed: failed on the harness's
// usage limit when its log ends on it, with why naming the limit and its
// reset time, and the reset time recorded, whether its exit was seen or not
// (ADR-0214), raised to the operator on a thread when it was a retry or the
// harness named no reset time; failed on a failed exit, with why; worked
// otherwise, whether it wrote an entry or found nothing to take in, an exit
// nobody saw included.
// It logs the run's activity in the curator's activity document, with what
// started it, the first line of the run's final text as its summary, and the
// entries as its items, which charge nothing. Entries that cannot be looked
// for, and an activity that cannot be logged, are warned of: the activity is
// logged naming no entry, and the run stays recorded as it ended.
func (l *launcher) curateEnded(run *AgentRun, exit *int) {
	ended := *run
	ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), exit
	entries, err := changedEntries(l.entry.Root, run.Started)
	if err != nil {
		l.warn("curator's entries not looked for", "err", err)
	}
	ended.Entries = entries
	lim, said := runLimit(run, l.warn)
	switch {
	case lim != nil:
		ended.Outcome, ended.Why, ended.ResetsAt = OutcomeFailed, limitWhy(lim), limitResets(lim)
		if run.Retry || ended.ResetsAt == "" {
			ended.LimitThread = l.raiseCurateLimit(&ended, said)
		}
	case exit != nil && *exit != 0:
		ended.Outcome, ended.Why = OutcomeFailed, fmt.Sprintf("ended (exit %d)", *exit)
	default:
		ended.Outcome, ended.Why = OutcomeWorked, ""
	}
	l.ended(&ended)
	args := []any{"pid", run.PID, "since", run.Since, "outcome", ended.Outcome, "entries", len(entries)}
	if exit != nil {
		args = append(args, "exit", *exit)
	}
	if lim != nil {
		args = append(args, "limit", lim.Error, "status", lim.Status, "resets_at", ended.ResetsAt)
	}
	l.log("curator ended", args...)
	// A run on the usage limit spent nothing, so logRunEnd logs no activity
	// for it, and the next run's since keeps pointing at the last activity:
	// the acceptances it never took in are taken in by the next (ADR-0214).
	if _, err := logRunEnd(l.dir, l.entry.Root, l.entry.Key, workitem.ActivityCurator, run.Trigger, "", entries); err != nil {
		l.warn("curator activity not logged", "err", err)
	}
}

// raiseCurateLimit opens a thread to the operator on the manifest, by flai
// serve, for run, a curator run that ended on the harness's usage limit as a
// retry, or on a limit with no reset time, said the harness's text at its end
// (ADR-0214 §7), journals the opening as the start is journalled, and returns
// the thread's ID. A run already recorded with its thread keeps it, so that
// it is opened once. A thread that cannot be opened is warned of and
// journalled failed, and "" is returned: flai serve schedules no retry for
// the run either way.
func (l *launcher) raiseCurateLimit(run *AgentRun, said string) string {
	if cur := l.dir.AgentStates()[l.entry.Root].Curator; cur.same(run) && cur.LimitThread != "" {
		return cur.LimitThread
	}
	title, what := "The curator ended on the harness's usage limit with no reset time; start it when the limit allows", ", which named no reset time"
	if run.Retry {
		title, what = "The curator ended on the harness's usage limit twice; start it when the limit allows", " again, as flai serve's retry"
	}
	at := l.now().UTC()
	entry := hostapi.Entry{At: at.Format(time.RFC3339), Action: hostapi.ActionCurate, Method: "serve.curate", Project: l.entry.Key, Root: l.entry.Root, By: limitAuthor}
	repo, err := workitem.Open(l.entry.Root)
	var th *threads.Thread
	if err == nil {
		th, err = threads.New(repo, threads.NewOptions{Title: title, On: manifest.File, Author: limitAuthor, Text: curateLimitText(run, said), Now: at})
	}
	if err != nil {
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("could not open a thread on %s for the curator's run started at %s, which ended on the harness's usage limit%s: %v", manifest.File, run.Started, what, err)
		if l.record != nil {
			l.record(entry)
		}
		l.warn("thread on the curator's usage limit not opened", "started", run.Started, "err", err.Error())
		return ""
	}
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("opened %s: the curator's run started at %s ended on the harness's usage limit%s; flai serve schedules no retry for it", th.ID, run.Started, what)
	if l.record != nil {
		l.record(entry)
	}
	l.log("curator's usage limit raised", "started", run.Started, "resets_at", run.ResetsAt, "thread", th.ID)
	return th.ID
}

// curateLimitText is the text of the thread raiseCurateLimit opens for run:
// what happened, the run, said the harness's text at its end, and what
// starts the curator again.
func curateLimitText(run *AgentRun, said string) string {
	lead := "The curator ended on the harness's usage limit, and the harness named no reset time."
	if run.Retry {
		lead = "The curator ended on the harness's usage limit twice in a row: its retry, started once the limit reset, ended on it again."
	}
	text := fmt.Sprintf("%s No digest entry is written until the curator runs.\n\nThe run started at %s, %s. It %s. Its log is %s.\n\n", lead, run.Started, startedBy(run), run.Why, run.Log)
	if said != "" {
		fence := "```"
		for strings.Contains(said, fence) {
			fence += "`"
		}
		text += fmt.Sprintf("The harness said:\n\n%[1]stext\n%[2]s\n%[1]s\n\n", fence, said)
	}
	text += "flai serve schedules no retry for this limit. `flai curate` starts the curator at any time."
	if run.ResetsAt != "" {
		return text + " The next acceptance, or `digest.schedule` coming round, starts it once the limit resets at " + run.ResetsAt + "."
	}
	return text + " The next acceptance, or `digest.schedule` coming round, starts it too."
}

// entryFile is a digest entry's file name, D-nnnn-<slug>.md, its number
// captured.
var entryFile = regexp.MustCompile(`^D-(\d+)-.*\.md$`)

// changedEntries are the IDs of the digest entries in the project at root
// whose files changed in the second started, a curator run's start, or
// after, in the order of their numbers: of the entry files directly in the
// design folder's digest/, its index aside. None when none changed since, or
// there is no folder.
func changedEntries(root, started string) ([]string, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	since, err := time.Parse(time.RFC3339, started)
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(digest.Dir(repo.Manifest))))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	type changed struct {
		n  int
		id string
	}
	var found []changed
	for _, f := range files {
		m := entryFile.FindStringSubmatch(f.Name())
		if !f.Type().IsRegular() || m == nil || digest.Skipped(f.Name()) {
			continue
		}
		info, err := f.Info()
		if err != nil || info.ModTime().Before(since) {
			continue
		}
		n, _ := strconv.Atoi(m[1])
		found = append(found, changed{n: n, id: digest.CanonicalID("D-" + m[1])})
	}
	slices.SortFunc(found, func(a, b changed) int { return a.n - b.n })
	ids := make([]string, 0, len(found))
	for _, c := range found {
		if !slices.Contains(ids, c.id) {
			ids = append(ids, c.id)
		}
	}
	if len(ids) == 0 {
		return nil, nil
	}
	return ids, nil
}
