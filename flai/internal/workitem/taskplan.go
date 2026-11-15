package workitem

import (
	"fmt"
	"sort"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// A story's task plan (S-0176): what each task waits for, read from its
// after: and the states of its story's tasks, and the layers of tasks that
// can run at once. flai show, the board, the move warning, and the host read
// it from here. It lives in workitem, not a package of its own, because the
// board's cards are built here and a package that imported workitem could
// not be called from them.

// Task states in a plan: the plan's own words, which flai show, the board,
// and the host give as they are, read from the kind of the state each task
// is in (ADR-0301 §2), not from its name.
const (
	PlanReady      = "ready"       // in a queue state, and waits for no task that is open
	PlanWaiting    = "waiting"     // in a queue state, and a task in its after: is open
	PlanInProgress = "in-progress" // in an active or a waiting state
	PlanDone       = "done"        // in the accepted state
	PlanCancelled  = "cancelled"   // in the cancelled state
)

// PlanTask is one task of a plan: its state, its after: as written, and
// the tasks of its after: that are open while it waits.
type PlanTask struct {
	ID         string   `json:"id"`
	State      string   `json:"state"`
	After      []string `json:"after,omitempty"`
	WaitingFor []string `json:"waiting_for,omitempty"`
}

// Plan is a story's tasks in ID order, and its layers: the open and done
// tasks grouped by the longest chain of after: steps before each, layer 0
// waiting for none, task IDs in ID order within a layer. A cancelled task is
// in no layer and holds no one, and neither is a task in a cycle or one
// that waits on a cycle: flai check reports the cycle.
type Plan struct {
	Tasks  []PlanTask `json:"tasks"`
	Layers [][]string `json:"layers"`
}

// TaskSummary is a story's tasks on its board card: how many are in each
// state, cancelled ones uncounted, and how many layers the plan has.
type TaskSummary struct {
	Ready      int `json:"ready"`
	Waiting    int `json:"waiting"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
	Layers     int `json:"layers"`
}

// PlanOf is PlanOfIn under the baseline workflow, for a caller that has no
// repository's workflow to hand.
func PlanOf(items []*Item, story string) *Plan {
	return PlanOfIn(baselineFlow, items, story)
}

// PlanOfIn is the plan of story's tasks among items, or nil when it has none,
// each task's state read from the kind w gives its status (ADR-0301 §2). An
// after: entry that names no task of the story is left out of both the
// waiting and the layers: flai check reports it.
func PlanOfIn(w workflow.Workflow, items []*Item, story string) *Plan {
	var tasks []*Item
	for _, it := range items {
		if it.Type == Task && it.Parent == story {
			tasks = append(tasks, it)
		}
	}
	if len(tasks) == 0 {
		return nil
	}
	sort.Slice(tasks, func(i, j int) bool { return lessID(tasks[i].ID, tasks[j].ID) })
	byID := map[string]*Item{}
	for _, t := range tasks {
		byID[t.ID] = t
	}
	p := &Plan{Tasks: make([]PlanTask, 0, len(tasks)), Layers: [][]string{}}
	for _, t := range tasks {
		pt := PlanTask{ID: t.ID, After: t.After}
		pt.State, pt.WaitingFor = taskState(w, t, byID)
		p.Tasks = append(p.Tasks, pt)
	}
	p.Layers = layers(w, tasks, byID)
	return p
}

// taskState is a task's state in its story's plan, read from the kind of its
// status in w, and the tasks it waits for when it waits: a task in the
// accepted state is done, one in the cancelled state cancelled, one in an
// active or a waiting state in progress, and any other has not started
// (ADR-0301 §2).
func taskState(w workflow.Workflow, t *Item, byID map[string]*Item) (string, []string) {
	switch kind := w.Kind(t.Status); {
	case kind == workflow.KindAccepted:
		return PlanDone, nil
	case kind == workflow.KindCancelled:
		return PlanCancelled, nil
	case underwayKind(kind):
		return PlanInProgress, nil
	}
	if waits := openAfter(w, t, byID); len(waits) > 0 {
		return PlanWaiting, waits
	}
	return PlanReady, nil
}

// openAfter is the tasks of t's story named in its after: that are in no
// state of a kind that closes an item, accepted or cancelled (ADR-0301 §2),
// in the order after: names them.
func openAfter(w workflow.Workflow, t *Item, byID map[string]*Item) []string {
	var out []string
	seen := map[string]bool{}
	for _, e := range t.After {
		id := CanonicalID(e)
		other, ok := byID[id]
		if !ok || id == t.ID || seen[id] || closedKind(w.Kind(other.Status)) {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

// layers groups the tasks that are in no state of the cancelled kind in w by
// the longest chain of after: steps before each. A task on a cycle, or after
// one, has no depth.
func layers(w workflow.Workflow, tasks []*Item, byID map[string]*Item) [][]string {
	cancelled := func(t *Item) bool { return w.Kind(t.Status) == workflow.KindCancelled }
	const visiting, cyclic = -1, -2
	depth := map[string]int{}
	var walk func(t *Item) int
	walk = func(t *Item) int {
		if d, ok := depth[t.ID]; ok {
			if d == visiting {
				return cyclic
			}
			return d
		}
		depth[t.ID] = visiting
		d := 0
		for _, e := range t.After {
			other, ok := byID[CanonicalID(e)]
			if !ok || cancelled(other) {
				continue
			}
			od := walk(other)
			if od == cyclic {
				d = cyclic
				break
			}
			d = max(d, od+1)
		}
		depth[t.ID] = d
		return d
	}
	out := [][]string{}
	for _, t := range tasks {
		if cancelled(t) {
			continue
		}
		d := walk(t)
		if d < 0 {
			continue
		}
		for len(out) <= d {
			out = append(out, nil)
		}
		out[d] = append(out[d], t.ID)
	}
	return out
}

// Summary is the plan as a story's board card counts it.
func (p *Plan) Summary() *TaskSummary {
	s := &TaskSummary{Layers: len(p.Layers)}
	for _, t := range p.Tasks {
		switch t.State {
		case PlanReady:
			s.Ready++
		case PlanWaiting:
			s.Waiting++
		case PlanInProgress:
			s.InProgress++
		case PlanDone:
			s.Done++
		}
	}
	return s
}

// Waits is why a task that is not started waits, in the voice of a story's
// hold, or "" when it waits for nothing open: flai move warns with it and
// moves the task all the same (S-0176).
func (p *Plan) Waits(task string) string {
	for _, t := range p.Tasks {
		if t.ID != task || t.State != PlanWaiting {
			continue
		}
		return fmt.Sprintf("waiting (after): waits for %s; ready to start when %s %s done or cancelled", and(p.labelled(t.WaitingFor)), and(t.WaitingFor), oneOrMany(len(t.WaitingFor), "is", "are"))
	}
	return ""
}

// labelled is each of ids with its state in the plan, the plan's in-progress
// read as two words: T-0002 (in progress).
func (p *Plan) labelled(ids []string) []string {
	state := map[string]string{}
	for _, t := range p.Tasks {
		state[t.ID] = t.State
	}
	out := make([]string, len(ids))
	for i, id := range ids {
		s := state[id]
		if s == PlanInProgress {
			s = "in progress"
		}
		out[i] = id + " (" + s + ")"
	}
	return out
}
