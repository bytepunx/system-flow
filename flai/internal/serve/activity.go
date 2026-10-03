package serve

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What the strategic agents did (S-0206, ADR-0079). flai serve keeps the log
// of each run of the planner, the orchestrator, and the analyzer as it does a
// story's agent's, named for the project's key, the kind, and the time it
// started. An activity ends when the agent reports it or when its run ends;
// its span runs from the later of the newest run's start and the document's
// last entry to its end, and its cost is what the kind's logs say was spent,
// apportioned to that span as a task's is to its intervals (ADR-0051), so a
// session resumed in a later run is charged only its share.

// ActivityLogs are the logs flai serve keeps of the runs of the strategic
// agent kind in the project named key, oldest first.
func (d Dir) ActivityLogs(key, kind string) ([]string, error) {
	if !workitem.IsActivityKind(kind) {
		return nil, fmt.Errorf("%q is not a strategic agent: use one of %s", kind, strings.Join(workitem.ActivityKinds, ", "))
	}
	// the start time follows the kind, so a project whose key ends in a kind
	// does not lend its stories' logs
	paths, err := filepath.Glob(filepath.Join(string(d), "agents", key+"-"+kind+"-[0-9]*.log"))
	sort.Strings(paths)
	return paths, err
}

// Logged is an activity flai logged and the document it was logged in, its
// totals accrued.
type Logged struct {
	Entry    workitem.ActivityEntry `json:"entry"`
	Activity *workitem.Activity     `json:"activity"`
	Logs     []string               `json:"logs"`
}

// LogActivity logs an activity of the strategic agent kind in the project at
// root, named key, that ended at end, measured from the kind's logs. An
// activity outside any run flai serve logged, such as a planner a person
// runs by hand, is logged with no seconds and no cost, so that it is still
// recorded.
func LogActivity(d Dir, root, key, kind, summary string, items []string, end time.Time) (*Logged, error) {
	m, err := readActivity(d, root, key, kind)
	if err != nil {
		return nil, err
	}
	end = end.UTC().Truncate(time.Second)
	e := workitem.ActivityEntry{At: end, Summary: summary, Items: items}
	if run, ok := m.newest(); ok {
		// no later than the second after the run's last event: a run that
		// died without its end logged is not charged the time since
		to := end
		if last := runEnd(run); to.After(last) {
			to = last
		}
		if s, inRun := activitySpan(run, m.last, to); inRun {
			e.Seconds, e.Cost, e.Estimated = m.measure(s)
		}
	}
	return m.append(e)
}

// LogRunEnd logs the newest run of the strategic agent kind in the project
// at root, named key, as having ended: the time in it since the last entry
// is one activity, with the run's final text as its summary. It logs
// nothing, and returns nil, when nothing was spent in that time.
func LogRunEnd(d Dir, root, key, kind string) (*Logged, error) {
	m, err := readActivity(d, root, key, kind)
	if err != nil {
		return nil, err
	}
	run, ok := m.newest()
	if !ok {
		return nil, nil
	}
	s, inRun := activitySpan(run, m.last, runEnd(run))
	if !inRun {
		return nil, nil
	}
	u := m.spent(s)
	if u == nil || len(u.Models) == 0 {
		return nil, nil
	}
	return m.append(workitem.ActivityEntry{
		At: s.To, Summary: firstLine(run.Result, "run ended"),
		Seconds: seconds(s), Cost: u.Cost(), Estimated: u.Estimated,
	})
}

// activityMeasure is what an activity is measured from: the kind's logs and
// the document's last entry.
type activityMeasure struct {
	repo  *workitem.Repo
	kind  string
	logs  []string
	rec   *usage.Record
	rates usage.Rates
	// last is when the document's newest entry ended; zero before the first.
	last time.Time
}

func readActivity(d Dir, root, key, kind string) (*activityMeasure, error) {
	logs, err := d.ActivityLogs(key, kind)
	if err != nil {
		return nil, err
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	doc, err := repo.Activity(kind)
	if err != nil {
		return nil, fmt.Errorf("cannot measure the %s's activity: %w", kind, err)
	}
	m := &activityMeasure{repo: repo, kind: kind, logs: logs}
	if m.logs == nil {
		m.logs = []string{}
	}
	for _, e := range doc.Entries {
		if e.At.After(m.last) {
			m.last = e.At
		}
	}
	if len(logs) == 0 {
		return m, nil
	}
	if m.rec, err = usage.Read(logs...); err != nil {
		return nil, err
	}
	if m.rates, err = d.Rates(); err != nil {
		return nil, err
	}
	return m, nil
}

// newest is the run that started last; false when there is none.
func (m *activityMeasure) newest() (usage.Run, bool) {
	if m.rec == nil || len(m.rec.Runs) == 0 {
		return usage.Run{}, false
	}
	run := m.rec.Runs[0]
	for _, r := range m.rec.Runs[1:] {
		if r.Start.After(run.Start) {
			run = r
		}
	}
	return run, true
}

// activitySpan is the span of an activity in run that ended at end: from the
// later of the run's start and the last entry's end, to end. It is not in
// the run when the run has nothing after the last entry: the run was logged
// already, and the activity happened outside any run flai serve logged.
func activitySpan(run usage.Run, last, end time.Time) (usage.Span, bool) {
	from := run.Start.UTC()
	if last.After(from) {
		if !run.End.After(last) {
			return usage.Span{}, false
		}
		from = last
	}
	if end.Before(from) {
		end = from
	}
	return usage.Span{From: from, To: end}, true
}

// spent is the usage apportioned to the span; nil when nothing was spent.
func (m *activityMeasure) spent(s usage.Span) *usage.Usage {
	const id = "activity"
	return m.rec.Tasks(map[string][]usage.Span{id: {s}}, m.rates)[id]
}

// measure is the span's wall-clock seconds and the cost apportioned to it.
func (m *activityMeasure) measure(s usage.Span) (secs int64, cost float64, estimated bool) {
	u := m.spent(s)
	if u == nil || len(u.Models) == 0 {
		return seconds(s), 0, false
	}
	return seconds(s), u.Cost(), u.Estimated
}

func (m *activityMeasure) append(e workitem.ActivityEntry) (*Logged, error) {
	doc, err := m.repo.AppendActivity(m.kind, e)
	if err != nil {
		return nil, err
	}
	return &Logged{Entry: doc.Entries[len(doc.Entries)-1], Activity: doc, Logs: m.logs}, nil
}

// runEnd is the second after run's last event, so that a span to it holds
// the run's last calls and the next activity begins after them.
func runEnd(run usage.Run) time.Time { return run.End.UTC().Truncate(time.Second).Add(time.Second) }

// seconds is the span's length in whole seconds.
func seconds(s usage.Span) int64 { return int64(s.To.Sub(s.From) / time.Second) }

// firstLine is the first line of text that is not blank, trimmed, or
// fallback when there is none.
func firstLine(text, fallback string) string {
	for l := range strings.Lines(text) {
		if l = strings.TrimSpace(l); l != "" {
			return l
		}
	}
	return fallback
}
