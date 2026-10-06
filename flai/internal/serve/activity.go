package serve

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/issues"
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
// session resumed in a later run is charged only its share. A planner's
// activity is also charged, that same share, to the item its run was started
// for and every item above it, under usage.strategic (S-0225, ADR-0083); an
// orchestrator's is split evenly between the work items its entry names and
// charged to each and every item above it (S-0226, ADR-0095); an analyzer's
// is split evenly between the issues its entry names and charged to each
// (S-0227). An analyzer run's end names the report the run wrote, or says it
// wrote none (S-0223), and names the issues that name the report, so its
// cost is split between them (S-0227).

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
	// Planned is the item a planner's activity was charged to (ADR-0083),
	// and Shared the work items an orchestrator's activity named, or the
	// issues an analyzer's named, that each took a share of it (ADR-0095,
	// S-0227). Charged are the items or issues whose usage the charge
	// changed, each item charged first and then upward, each once. All are
	// empty when nothing was charged.
	Planned string   `json:"planned,omitempty"`
	Shared  []string `json:"shared,omitempty"`
	Charged []string `json:"charged,omitempty"`
}

// LogActivity logs an activity of the strategic agent kind in the project at
// root, named key, that ended at end, measured from the kind's logs, and
// charges a planner's to the item it planned, an orchestrator's to the work
// items it names, and an analyzer's to the issues it names. An activity
// outside any run flai serve logged, such as a planner a person runs by
// hand, is logged with no seconds and no cost, so that it is still recorded,
// and charges nothing.
// When the charge fails the activity is logged all the same, and is returned
// with the error.
func LogActivity(d Dir, root, key, kind, summary string, items []string, end time.Time) (*Logged, error) {
	m, err := readActivity(d, root, key, kind)
	if err != nil {
		return nil, err
	}
	end = end.UTC().Truncate(time.Second)
	e := workitem.ActivityEntry{At: end, Summary: summary, Items: items}
	run, ok := m.newest()
	var spent *usage.Usage
	if ok {
		// no later than the second after the run's last event: a run that
		// died without its end logged is not charged the time since
		to := end
		if last := runEnd(run); to.After(last) {
			to = last
		}
		if s, inRun := activitySpan(run, m.last, to); inRun {
			spent = m.spent(s)
			e.Seconds = seconds(s)
			if spent != nil && len(spent.Models) > 0 {
				e.Cost, e.Estimated = spent.Cost(), spent.Estimated
			}
		}
	}
	logged, err := m.append(e)
	if err != nil {
		return nil, err
	}
	return logged, m.charge(logged, run, spent)
}

// LogRunEnd logs the newest run of the strategic agent kind in the project
// at root, named key, as having ended: the time in it since the last entry
// is one activity, with the run's final text as its summary and items as
// the items it names (S-0209), and is charged as LogActivity's is. It logs
// nothing, and returns nil, when nothing was spent in that time.
func LogRunEnd(d Dir, root, key, kind string, items []string) (*Logged, error) {
	return logRunEnd(d, root, key, kind, "", "", items)
}

// logRunEnd is LogRunEnd with the entry's trigger, what started the run
// (ADR-0084), empty for none, and its summary in place of the run's final
// text, such as why it was stopped (S-0218), empty for that text.
func logRunEnd(d Dir, root, key, kind, trigger, summary string, items []string) (*Logged, error) {
	return logRunEndSaying(d, root, key, kind, trigger, func(final string) string {
		if summary != "" {
			return summary
		}
		return final
	}, items)
}

// logRunEndSaying is logRunEnd with the entry's summary made by say from the
// first line of the run's final text, "run ended" when it has none: the
// analyzer's names the report the run wrote, or says it wrote none (S-0223).
func logRunEndSaying(d Dir, root, key, kind, trigger string, say func(final string) string, items []string) (*Logged, error) {
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
	logged, err := m.append(workitem.ActivityEntry{
		At: s.To, Summary: say(firstLine(run.Result, "run ended")), Trigger: trigger, Items: items,
		Seconds: seconds(s), Cost: u.Cost(), Estimated: u.Estimated,
	})
	if err != nil {
		return nil, err
	}
	return logged, m.charge(logged, run, u)
}

// activityMeasure is what an activity is measured from: the kind's logs and
// the document's last entry.
type activityMeasure struct {
	dir   Dir
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
	m := &activityMeasure{dir: d, repo: repo, kind: kind, logs: logs}
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

// charge charges u, what an activity in run spent, with the entry's
// seconds, so that the items and the activity document are two views of one
// spend: a planner's to the item it planned, an orchestrator's to the work
// items its entry names, an analyzer's to the issues its entry names. It
// charges nothing for a u that spent nothing. It notes what it charged on
// logged.
func (m *activityMeasure) charge(logged *Logged, run usage.Run, u *usage.Usage) error {
	if u == nil || len(u.Models) == 0 {
		return nil
	}
	c := u.Clone()
	c.Seconds = logged.Entry.Seconds
	switch m.kind {
	case workitem.ActivityPlanner:
		return m.chargePlanned(logged, run, c)
	case workitem.ActivityOrchestrator:
		return m.chargeNamed(logged, c)
	case workitem.ActivityAnalyzer:
		return m.chargeIssues(logged, c)
	}
	return nil
}

// chargePlanned charges u to the item run was started for, as
// serve/agents.json records it under plans, and to every item above it
// (ADR-0083). It charges nothing for a run no item's newest planner run is,
// such as one a person ran by hand.
func (m *activityMeasure) chargePlanned(logged *Logged, run usage.Run, u *usage.Usage) error {
	item := m.planned(run)
	if item == "" {
		return nil
	}
	changed, err := m.repo.ChargeStrategic(item, m.kind, u)
	if len(changed) > 0 {
		logged.Planned, logged.Charged = item, changed
	}
	if err != nil {
		return fmt.Errorf("the %s's activity was logged, but its cost could not be charged to %s: %w", m.kind, item, err)
	}
	return nil
}

// chargeNamed splits u evenly between the work items the entry names and
// charges each share to its item and every item above it (ADR-0095). It
// charges nothing when the entry names no work item. A share that cannot be
// charged is left out, the others charged all the same.
func (m *activityMeasure) chargeNamed(logged *Logged, u *usage.Usage) error {
	items := m.workItems(logged.Entry.Items)
	if len(items) == 0 {
		return nil
	}
	var errs []error
	for i, share := range u.Split(len(items)) {
		changed, err := m.repo.ChargeStrategic(items[i], m.kind, share)
		if len(changed) > 0 {
			logged.Shared = append(logged.Shared, items[i])
			for _, id := range changed {
				if !slices.Contains(logged.Charged, id) {
					logged.Charged = append(logged.Charged, id)
				}
			}
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", items[i], err))
		}
	}
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("the %s's activity was logged, but not all its cost could be charged to the items it named: %w", m.kind, err)
	}
	return nil
}

// workItems are the epics, stories, and tasks named that exist, archived
// ones included, each once by its own ID, in the order first named. Anything
// else named, such as a thread, an issue, or an ID no item has, is left out.
// An item that cannot be read is kept, so that its share's charge fails and
// says why.
func (m *activityMeasure) workItems(named []string) []string {
	var items []string
	for _, id := range named {
		if workitem.TypeOfID(workitem.CanonicalID(id)) == "" {
			continue
		}
		it, err := m.repo.Get(id)
		switch {
		case errors.Is(err, workitem.ErrNotFound):
			continue
		case err != nil:
			id = workitem.CanonicalID(id)
		default:
			id = it.ID
		}
		if !slices.Contains(items, id) {
			items = append(items, id)
		}
	}
	return items
}

// chargeIssues splits u evenly between the issues the entry names and
// charges each share to its issue, under its usage.strategic (S-0227). It
// charges nothing when the entry names no issue, and no work item the entry
// names. A share that cannot be charged is left out, the others charged all
// the same.
func (m *activityMeasure) chargeIssues(logged *Logged, u *usage.Usage) error {
	ids := m.issues(logged.Entry.Items)
	if len(ids) == 0 {
		return nil
	}
	var errs []error
	for i, share := range u.Split(len(ids)) {
		is, err := issues.ChargeStrategic(m.repo, ids[i], m.kind, share)
		if is != nil {
			logged.Shared = append(logged.Shared, is.ID)
			logged.Charged = append(logged.Charged, is.ID)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("the %s's activity was logged, but not all its cost could be charged to the issues it named: %w", m.kind, err)
	}
	return nil
}

// issueID is an issue's ID as workitem.CanonicalID writes it.
var issueID = regexp.MustCompile(`^I-\d+$`)

// issues are the issues named that exist, each once by its own ID, in the
// order first named. Anything else named, such as a work item, a thread, or
// an ID no issue has, is left out, as is an issue that cannot be read, since
// the issues package does not tell it from one that does not exist.
func (m *activityMeasure) issues(named []string) []string {
	var ids []string
	for _, id := range named {
		id = workitem.CanonicalID(id)
		if !issueID.MatchString(id) {
			continue
		}
		is, err := issues.Get(m.repo, id)
		if err != nil {
			continue
		}
		if !slices.Contains(ids, is.ID) {
			ids = append(ids, is.ID)
		}
	}
	return ids
}

// planned is the item whose newest planner run kept its log in run's, or ""
// when there is none. The log is told by its name, which holds the
// project's key, so that the form a root or the serve folder is written in
// does not matter.
func (m *activityMeasure) planned(run usage.Run) string {
	name := filepath.Base(run.Path)
	for _, st := range m.dir.AgentStates() {
		for item, r := range st.Plans {
			if r != nil && r.Log != "" && filepath.Base(r.Log) == name {
				return item
			}
		}
	}
	return ""
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
