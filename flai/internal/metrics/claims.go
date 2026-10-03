package metrics

import (
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Claims is how many items were in progress each day against the limit, and
// how far each story's touches were from what its commits changed (S-0205).
type Claims struct {
	Limit int        `json:"limit,omitempty"`
	Days  []ClaimDay `json:"days"`
	// Drift is nil when git could not be read.
	Drift *[]TouchDrift `json:"drift,omitempty"`
}

// ClaimDay is the number of items in progress at the end of a day.
type ClaimDay struct {
	Date       string `json:"date"`
	InProgress int    `json:"in_progress"`
}

// TouchDrift compares the files a story's commits changed with its touches.
type TouchDrift struct {
	ID             string   `json:"id"`
	Committed      []string `json:"committed"`
	Outside        []string `json:"outside"`
	Unchanged      []string `json:"unchanged"`
	OutsideCount   int      `json:"outside_count"`
	UnchangedCount int      `json:"unchanged_count"`
}

// claims lays out the items in progress by the day, from the one that holds
// the window's start to today, and the touches drift of the stories among
// them.
func claims(items, all []*workitem.Item, start time.Time, opt Options) Claims {
	out := Claims{Limit: opt.WIPLimit, Days: []ClaimDay{}}
	for d := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC); !d.After(opt.Now); d = d.AddDate(0, 0, 1) {
		end := d.Add(24*time.Hour - time.Second)
		p := ClaimDay{Date: d.Format("2006-01-02")}
		for _, it := range items {
			if stateAt(it, end) == workitem.InProgress {
				p.InProgress++
			}
		}
		out.Days = append(out.Days, p)
	}
	if opt.Commits != nil {
		drift := touchDrift(items, all, opt.Projects, opt.Commits)
		out.Drift = &drift
	}
	return out
}

// touchDrift is, in ID order, each story among items with committed files:
// those under none of its touches and its touches under which none falls. Its
// touches are its own and its tasks' that were not cancelled, read as a claim
// reads them.
func touchDrift(items, all []*workitem.Item, projects []manifest.Project, commits map[string][]string) []TouchDrift {
	tasks := map[string][]*workitem.Item{}
	for _, it := range all {
		if it.Type == workitem.Task && it.Status != workitem.Cancelled {
			p := workitem.CanonicalID(it.Parent)
			tasks[p] = append(tasks[p], it)
		}
	}
	reader := workitem.NewHolds(nil, projects)
	out := []TouchDrift{}
	for _, it := range items {
		id := workitem.CanonicalID(it.ID)
		files := commits[id]
		if it.Type != workitem.Story || len(files) == 0 {
			continue
		}
		touches := slices.Clone(it.Touches)
		for _, t := range tasks[id] {
			touches = append(touches, t.Touches...)
		}
		claim := reader.Claim(&workitem.Item{Touches: touches})
		d := TouchDrift{ID: it.ID, Committed: slices.Compact(slices.Sorted(slices.Values(files))), Outside: []string{}, Unchanged: []string{}}
		used := map[string]bool{}
		for _, f := range d.Committed {
			in := false
			for _, t := range claim {
				if f == t || strings.HasPrefix(f, t+"/") {
					used[t], in = true, true
				}
			}
			if !in {
				d.Outside = append(d.Outside, f)
			}
		}
		for _, t := range claim {
			if !used[t] {
				d.Unchanged = append(d.Unchanged, t)
			}
		}
		sort.Strings(d.Unchanged)
		d.OutsideCount, d.UnchangedCount = len(d.Outside), len(d.Unchanged)
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// heldSeconds is, for each story among items that was ever in ready, the
// seconds the hold rules held it there. They are replayed at each creation and
// transition of a story or task, over the states of the time and today's
// touches and after; nothing changes until the next one.
func heldSeconds(items, all []*workitem.Item, projects []manifest.Project, now time.Time) map[string]*float64 {
	held := map[string]*float64{}
	for _, it := range items {
		if it.Type == workitem.Story && !it.FirstAt(workitem.Ready).IsZero() {
			held[it.ID] = new(float64)
		}
	}
	if len(held) == 0 {
		return held
	}
	type event struct {
		at time.Time
		i  int    // the item's copy
		to string // the state entered, "" for its creation
	}
	var copies []workitem.Item
	var events []event
	for _, it := range all {
		if it.Type != workitem.Story && it.Type != workitem.Task {
			continue
		}
		c := *it
		c.Status, c.Archived = workitem.Backlog, false
		i := len(copies)
		copies = append(copies, c)
		created, _ := time.Parse(workitem.TimeFormat, it.Created)
		events = append(events, event{at: created, i: i})
		for _, tr := range it.Transitions {
			at, _ := time.Parse(workitem.TimeFormat, tr.At)
			events = append(events, event{at: at, i: i, to: tr.To})
		}
	}
	sort.SliceStable(events, func(a, b int) bool { return events[a].at.Before(events[b].at) })
	var snapshot, ready []*workitem.Item
	for k := 0; k < len(events) && events[k].at.Before(now); {
		at := events[k].at
		for ; k < len(events) && events[k].at.Equal(at); k++ {
			e := events[k]
			if e.to == "" {
				snapshot = append(snapshot, &copies[e.i])
			} else {
				copies[e.i].Status = e.to
			}
		}
		until := now
		if k < len(events) && events[k].at.Before(now) {
			until = events[k].at
		}
		ready = ready[:0]
		for _, it := range snapshot {
			if held[it.ID] != nil && it.Status == workitem.Ready {
				ready = append(ready, it)
			}
		}
		if len(ready) == 0 {
			continue
		}
		h := workitem.NewHolds(snapshot, projects)
		for _, s := range ready {
			if h.Of(s) != nil {
				*held[s.ID] += until.Sub(at).Seconds()
			}
		}
	}
	return held
}
