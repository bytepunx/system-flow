package serve

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What agents spent (S-0143, ADR-0051). flai serve keeps the log of every
// agent it starts, named for the project's key, the story, and the time it
// started. From those logs it measures a story whole, and each of the
// story's tasks from the calls of the sub-agents started for it and its
// share of the rest over the intervals it was in progress (S-0230), and
// writes their usage into the front matter: when an agent it started ends,
// and when a task of a story whose agent runs enters done without having
// been measured. flai serve agent usage does the same on the operator's word.

// AgentLogs are the logs flai serve keeps of the agents it started for
// story in the project named key, oldest first.
func (d Dir) AgentLogs(key, story string) ([]string, error) {
	if !StoryID.MatchString(story) {
		return nil, fmt.Errorf("%q is not a story's ID", story)
	}
	paths, err := filepath.Glob(filepath.Join(string(d), "agents", key+"-"+story+"-*.log"))
	sort.Strings(paths)
	return paths, err
}

// LoggedStories are the stories of the project named key that flai serve
// has kept a log for, in order of ID.
func (d Dir) LoggedStories(key string) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(string(d), "agents", key+"-S-*.log"))
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []string
	for _, p := range paths {
		name := filepath.Base(p)[len(key)+1:]
		for i := 2; i < len(name); i++ {
			if name[i] == '-' {
				if id := name[:i]; StoryID.MatchString(id) && !seen[id] {
					seen[id] = true
					out = append(out, id)
				}
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return workitem.CanonicalID(out[i]) < workitem.CanonicalID(out[j]) })
	return out, nil
}

// Rates are what each model's tokens cost across every log flai serve
// keeps, to price what no result reported.
func (d Dir) Rates() (usage.Rates, error) {
	paths, err := filepath.Glob(filepath.Join(string(d), "agents", "*.log"))
	if err != nil {
		return nil, err
	}
	return usage.ReadRates(paths...)
}

// Measured is what the logs say a story and its tasks spent.
type Measured struct {
	Story string                  `json:"story"`
	Logs  []string                `json:"logs"`
	Usage *usage.Usage            `json:"usage,omitempty"`
	Tasks map[string]*usage.Usage `json:"tasks"`
	// Changed are the items whose usage was written, their epic's among them.
	Changed []string `json:"changed"`
}

// Measure measures story in the project at root, named key, from the logs
// of the agents flai serve started for it: the story from all of them, each
// task from its sub-agents' calls and its share of the rest over the
// intervals it was in progress. With write, it records what differs from
// what the items say, and sums the story's epic again. A task that has been
// done without anything measured for it is written as having spent nothing,
// so that it is not measured at every look.
func Measure(d Dir, root, key, story string, write bool) (*Measured, error) {
	logs, err := d.AgentLogs(key, story)
	if err != nil {
		return nil, err
	}
	out := &Measured{Story: story, Logs: logs, Tasks: map[string]*usage.Usage{}, Changed: []string{}}
	if logs == nil {
		out.Logs = []string{}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	if _, err := repo.Get(story); err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return out, nil
	}
	rec, err := usage.Read(logs...)
	if err != nil {
		return nil, err
	}
	rates, err := d.Rates()
	if err != nil {
		return nil, err
	}
	out.Usage = rec.Total(rates)
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	tasks := workitem.Children(items, story)
	spans := map[string][]usage.Span{}
	for _, t := range tasks {
		spans[t.ID] = inProgress(t)
	}
	measured := rec.Tasks(spans, rates)
	for _, t := range tasks {
		u := measured[t.ID]
		if u == nil && t.Status == workitem.Done && len(spans[t.ID]) > 0 {
			u = &usage.Usage{Source: usage.SourceLog, Models: []usage.Model{}}
		}
		if u != nil {
			out.Tasks[t.ID] = u
		}
	}
	if !write {
		return out, nil
	}
	put := func(id string, u *usage.Usage) error {
		// read again just before writing: an agent may have changed it
		fresh, err := repo.Get(id)
		if err != nil {
			return err
		}
		if u == nil {
			return nil
		}
		// what strategic agents spent on it is no agent's: keep it (ADR-0083)
		u = usage.WithStrategic(u, fresh.Usage)
		if fresh.Usage != nil && usage.Same(fresh.Usage, u) {
			return nil
		}
		fresh.Usage = u
		if err := repo.Save(fresh); err != nil {
			return err
		}
		out.Changed = append(out.Changed, id)
		return nil
	}
	ids := make([]string, 0, len(out.Tasks))
	for id := range out.Tasks {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := put(id, out.Tasks[id]); err != nil {
			return out, err
		}
	}
	if err := put(story, out.Usage); err != nil {
		return out, err
	}
	it, err := repo.Get(story)
	if err != nil {
		return out, err
	}
	up, err := repo.RollUp(it)
	out.Changed = append(out.Changed, up...)
	return out, err
}

// inProgress are the intervals an item was in progress, the last open when
// it still is.
func inProgress(it *workitem.Item) []usage.Span {
	var out []usage.Span
	for i, tr := range it.Transitions {
		if tr.To != workitem.InProgress {
			continue
		}
		from, err := time.Parse(workitem.TimeFormat, tr.At)
		if err != nil {
			continue
		}
		s := usage.Span{From: from}
		if i+1 < len(it.Transitions) {
			s.To, _ = time.Parse(workitem.TimeFormat, it.Transitions[i+1].At)
		}
		out = append(out, s)
	}
	return out
}

// measure measures story and says so in flai serve's log.
func (l *launcher) measure(story string) {
	m, err := Measure(l.dir, l.entry.Root, l.entry.Key, story, true)
	if err != nil {
		l.warn("usage not measured", "story", story, "err", err)
		return
	}
	if len(m.Changed) > 0 {
		l.log("usage measured", "story", story, "logs", len(m.Logs), "tokens", m.Usage.Tokens(), "cost", fmt.Sprintf("%.2f", m.Usage.Cost()), "changed", len(m.Changed))
	}
}

// measureDone measures the story of each agent that runs when one of its
// tasks is done and has not been measured.
func (l *launcher) measureDone(st AgentState) {
	var running []string
	for id, run := range st.Stories {
		if run.live() {
			running = append(running, id)
		}
	}
	if len(running) == 0 {
		return
	}
	repo, err := workitem.Open(l.entry.Root)
	if err != nil {
		return
	}
	items, err := repo.List(false)
	if err != nil {
		return
	}
	sort.Strings(running)
	for _, id := range running {
		for _, t := range workitem.Children(items, id) {
			if t.Status == workitem.Done && t.Usage == nil {
				l.measure(id)
				break
			}
		}
	}
}
