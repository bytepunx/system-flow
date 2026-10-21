package metrics

import (
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Claims is how many items were in progress and how many stories were held
// each day against the limit, the time held per week by reason, and how far
// each story's touches were from what its commits changed (S-0205, S-0214).
type Claims struct {
	Limit int         `json:"limit,omitempty"`
	Days  []ClaimDay  `json:"days"`
	Weeks []ClaimWeek `json:"weeks"`
	// Drift is nil when git could not be read.
	Drift *[]TouchDrift `json:"drift,omitempty"`
}

// ClaimDay is the number of items in progress and of stories held at the end
// of a day.
type ClaimDay struct {
	Date       string `json:"date"`
	InProgress int    `json:"in_progress"`
	Held       int    `json:"held"`
}

// ClaimWeek is the time stories were held in one ISO week, by reason, and of
// the stories with touches drift completed in it, those whose touches were
// exact (ADR-0113). Stories, Exact, and ExactShare are nil when git could not
// be read, and ExactShare when no such story was completed in the week.
type ClaimWeek struct {
	Week        string     `json:"week"`
	Start       string     `json:"start"`
	HeldSeconds HeldByCode `json:"held_seconds"`
	Stories     *int       `json:"stories,omitempty"`
	Exact       *int       `json:"exact,omitempty"`
	ExactShare  *float64   `json:"exact_share,omitempty"`
}

// HeldByCode is held time in seconds by the reason each hold is named by,
// every reason present.
type HeldByCode struct {
	Overlap   float64 `json:"overlap"`
	After     float64 `json:"after"`
	NoTouches float64 `json:"no-touches"`
}

// add counts seconds under the reason code.
func (h *HeldByCode) add(code string, seconds float64) {
	switch code {
	case workitem.HoldOverlap:
		h.Overlap += seconds
	case workitem.HoldAfter:
		h.After += seconds
	case workitem.HoldNoTouches:
		h.NoTouches += seconds
	}
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

// claims lays out the items in progress and the stories held by the day, from
// the one that holds the window's start to today; the time held by reason and
// the share of stories whose touches were exact by the ISO week, from the one
// that holds the window's start to this one; and the touches drift of the
// stories among the items. replay is the hold rules' replay over the same
// days and weeks.
func claims(items, all []*workitem.Item, start time.Time, opt Options, replay holdReplay) Claims {
	out := Claims{Limit: opt.WIPLimit, Days: []ClaimDay{}, Weeks: []ClaimWeek{}}
	for i, d := 0, firstDay(start); !d.After(opt.Now); i, d = i+1, d.AddDate(0, 0, 1) {
		end := d.Add(24*time.Hour - time.Second)
		p := ClaimDay{Date: d.Format("2006-01-02"), Held: replay.days[i]}
		for _, it := range items {
			if stateAt(it, end) == workitem.InProgress {
				p.InProgress++
			}
		}
		out.Days = append(out.Days, p)
	}
	at := map[string]int{}
	for i, m := 0, monday(start); !m.After(opt.Now); i, m = i+1, m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		at[weekKey(y, w)] = i
		out.Weeks = append(out.Weeks, ClaimWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02"), HeldSeconds: replay.weeks[i]})
	}
	if opt.Commits == nil {
		return out
	}
	drift := touchDrift(items, all, opt.Projects, opt.Commits)
	out.Drift = &drift
	stories, exact := make([]int, len(out.Weeks)), make([]int, len(out.Weeks))
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	for _, d := range drift {
		it := byID[d.ID]
		done := it.CompletedAt()
		if it.Status != workitem.Done || done.Before(start) || done.After(opt.Now) {
			continue
		}
		y, w := done.ISOWeek()
		i := at[weekKey(y, w)]
		stories[i]++
		if d.OutsideCount == 0 && d.UnchangedCount == 0 {
			exact[i]++
		}
	}
	for i := range out.Weeks {
		b := &out.Weeks[i]
		b.Stories, b.Exact = &stories[i], &exact[i]
		if stories[i] > 0 {
			share := float64(exact[i]) / float64(stories[i])
			b.ExactShare = &share
		}
	}
	return out
}

// firstDay is the start of the UTC day that holds t.
func firstDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
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

// holdReplay is what the hold rules held, replayed: the seconds each story
// was held, the stories held at the end of each day from the one that holds
// the window's start to today, and the seconds held by reason in each ISO
// week from the one that holds the window's start to this one (ADR-0113).
type holdReplay struct {
	seconds map[string]*float64
	days    []int
	weeks   []HeldByCode
	// judged is how many times the holds were judged anew (I-0155).
	judged int
}

// replayHolds replays the hold rules for each story among items that was ever
// in ready. They are replayed at each creation and transition of a story or
// task, over the states of the time and today's touches and after; nothing
// changes until the next one, and nothing after now counts. A hold counts
// under the reason it is named by, the code Holds.Of gives it. An overlap
// inside the manifest's shared paths holds nothing (ADR-0096). The holds are
// judged anew only after an event that can change one, of a story or of a
// task of a story in ready or in progress, over only the items they read, and
// with what did not change remembered (I-0155).
func replayHolds(items, all []*workitem.Item, start time.Time, opt Options) holdReplay {
	// Held time is whole seconds, as timestamps are: a hold still open runs
	// to now's last whole second.
	now := opt.Now.Truncate(time.Second)
	out := holdReplay{seconds: map[string]*float64{}}
	day0, week0 := firstDay(start), monday(start)
	for d := day0; !d.After(now); d = d.AddDate(0, 0, 1) {
		out.days = append(out.days, 0)
	}
	for m := week0; !m.After(now); m = m.AddDate(0, 0, 7) {
		out.weeks = append(out.weeks, HeldByCode{})
	}
	for _, it := range items {
		if it.Type == workitem.Story && !it.FirstAt(workitem.Ready).IsZero() {
			out.seconds[it.ID] = new(float64)
		}
	}
	if len(out.seconds) == 0 {
		return out
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
	// dayEnd is the moment a day's held stories are counted at: its last
	// second, today's at now.
	dayEnd := func(i int) time.Time {
		end := day0.AddDate(0, 0, i).Add(24*time.Hour - time.Second)
		if end.After(now) {
			return now
		}
		return end
	}
	// What a hold reads, kept as the events apply: the stories created, in
	// creation order and by ID, the last created of each; the tasks created, by
	// parent as written, as NewHolds keys them; and how many stories of each ID
	// are in ready or in progress.
	var stories []*workitem.Item
	byID := map[string]*workitem.Item{}
	tasks := map[string][]*workitem.Item{}
	active := map[string]int{}
	isActive := func(status string) bool { return status == workitem.Ready || status == workitem.InProgress }
	// A story's move changes a hold only into or out of a state a hold reads,
	// and a task's only into or out of done or cancelled, the states Claim and
	// NewHolds read.
	reads := func(status string) bool { return isActive(status) || status == workitem.Done }
	closing := func(status string) string {
		if status == workitem.Done || status == workitem.Cancelled {
			return status
		}
		return ""
	}
	// held is the judgement in force: each tracked story held, by its code.
	type heldStory struct {
		seconds *float64
		code    string
	}
	var held []heldStory
	// memo keeps what the holds work out, claims and how one holds another,
	// from one judgement to the next. A story's claim is forgotten when a task
	// of it is created or enters or leaves done or cancelled: all Claim reads
	// of them.
	memo := workitem.NewHoldsMemo()
	day := 0
	for k := 0; k < len(events) && !events[k].at.After(now); {
		at := events[k].at
		judge := false
		for ; k < len(events) && events[k].at.Equal(at); k++ {
			c := &copies[events[k].i]
			from, to := c.Status, events[k].to
			switch {
			case to == "" && c.Type == workitem.Story:
				stories = append(stories, c)
				byID[c.ID] = c
				judge = true
			case to == "":
				tasks[c.Parent] = append(tasks[c.Parent], c)
				memo.Forget(c.Parent)
				judge = judge || active[c.Parent] > 0
			case c.Type == workitem.Story:
				if isActive(from) {
					active[c.ID]--
				}
				if isActive(to) {
					active[c.ID]++
				}
				c.Status = to
				judge = judge || reads(from) || reads(to)
			default:
				c.Status = to
				if closing(from) != closing(to) {
					memo.Forget(c.Parent)
					judge = judge || active[c.Parent] > 0
				}
			}
		}
		if judge {
			held = held[:0]
			var ready, in []*workitem.Item
			for _, s := range stories {
				if s.Status == workitem.Ready && out.seconds[s.ID] != nil {
					ready = append(ready, s)
				}
			}
			if len(ready) > 0 {
				// Only the stories in ready or in progress and their tasks,
				// and the stories the ready ones name in after, are read.
				withTasks := map[string]bool{}
				for _, s := range stories {
					if !isActive(s.Status) {
						continue
					}
					in = append(in, s)
					if !withTasks[s.ID] {
						withTasks[s.ID] = true
						for _, t := range tasks[s.ID] {
							if t.Status != workitem.Cancelled {
								in = append(in, t)
							}
						}
					}
				}
				for _, s := range ready {
					for _, e := range s.After {
						if a := byID[workitem.CanonicalID(e)]; a != nil && !isActive(a.Status) {
							in = append(in, a)
						}
					}
				}
				h := memo.NewHolds(in, opt.Projects).WithShared(opt.Shared)
				out.judged++
				for _, s := range ready {
					if hold := h.Of(s); hold != nil {
						held = append(held, heldStory{seconds: out.seconds[s.ID], code: hold.Code})
					}
				}
			}
		}
		// The states hold from at up to next; until is next, not after now.
		until, next := now, now.Add(time.Second)
		if k < len(events) && !events[k].at.After(now) {
			until, next = events[k].at, events[k].at
		}
		byCode := map[string]int{}
		for _, s := range held {
			*s.seconds += until.Sub(at).Seconds()
			byCode[s.code]++
		}
		for code, n := range byCode {
			for i := range out.weeks {
				wk := week0.AddDate(0, 0, 7*i)
				from, to := laterOf(at, wk), earlierOf(until, wk.AddDate(0, 0, 7))
				if to.After(from) {
					out.weeks[i].add(code, float64(n)*to.Sub(from).Seconds())
				}
			}
		}
		for day < len(out.days) && dayEnd(day).Before(at) {
			day++ // counted under earlier states, or before any item was
		}
		for ; day < len(out.days) && dayEnd(day).Before(next); day++ {
			for _, n := range byCode {
				out.days[day] += n
			}
		}
	}
	return out
}
