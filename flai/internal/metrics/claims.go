package metrics

import (
	"encoding/json"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
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
// days and weeks. The days are swept once, each item's state moved on as the
// days pass its transitions, and a day before every item is answered without
// reading them (S-0608).
func claims(items, all []*workitem.Item, tl timelines, start time.Time, opt Options, replay holdReplay) Claims {
	out := Claims{Limit: opt.WIPLimit, Days: []ClaimDay{}, Weeks: []ClaimWeek{}}
	states := newDayStates(items, tl)
	for i, d := 0, firstDay(start); !d.After(opt.Now); i, d = i+1, d.AddDate(0, 0, 1) {
		end := d.Add(24*time.Hour - time.Second)
		p := ClaimDay{Date: d.Format("2006-01-02"), Held: replay.days[i]}
		if !states.before(d.AddDate(0, 0, 1)) {
			for k := range items {
				if states.at(k, end) == workitem.InProgress {
					p.InProgress++
				}
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
	drift := driftFor(items, all, opt)
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
// reads them. The files and the touches under opt.FlaiFolders are left out
// on both sides, since flai writes them (ADR-0284): a story with no committed
// file outside them has no entry.
func touchDrift(items, all []*workitem.Item, opt Options) []TouchDrift {
	commits, folders := opt.Commits, opt.FlaiFolders
	tasks := map[string][]*workitem.Item{}
	for _, it := range all {
		if it.Type == workitem.Task && it.Status != workitem.Cancelled {
			p := workitem.CanonicalID(it.Parent)
			tasks[p] = append(tasks[p], it)
		}
	}
	reader := workitem.NewHolds(nil, opt.Projects)
	out := []TouchDrift{}
	for _, it := range items {
		id := workitem.CanonicalID(it.ID)
		if it.Type != workitem.Story {
			continue
		}
		files := slices.DeleteFunc(slices.Clone(commits[id]), func(f string) bool { return storygit.InFolders(f, folders) })
		if len(files) == 0 {
			continue
		}
		touches := slices.Clone(it.Touches)
		for _, t := range tasks[id] {
			touches = append(touches, t.Touches...)
		}
		claim := slices.DeleteFunc(slices.Clone(reader.Claim(&workitem.Item{Touches: touches})), func(t string) bool { return storygit.InFolders(t, folders) })
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
// in ready, over all, up to the request's now, and lays what they held out
// over the window from start: holdTimeline replays, and slice lays it out.
// Nothing is kept; holdsFor keeps the replay between requests (S-0608).
func replayHolds(items, all []*workitem.Item, start time.Time, opt Options) holdReplay {
	now := opt.Now.Truncate(time.Second)
	out := newHoldReplay(items, start, now)
	if len(out.seconds) == 0 {
		return out
	}
	return replayTimeline(all, now, opt).slice(out, start, now)
}

// holdsFor is replayHolds, its replay kept for opt.Root between requests while
// the item generation, the manifest's projects and shared paths, and the
// replay's horizon hold (S-0608): only the slicing to the window is done
// again. Without a root or a generation, as in flai stats, nothing is kept.
func holdsFor(items, all []*workitem.Item, start time.Time, opt Options) holdReplay {
	now := opt.Now.Truncate(time.Second)
	out := newHoldReplay(items, start, now)
	if len(out.seconds) == 0 {
		return out
	}
	t, _ := keptTimeline(all, now, opt)
	return t.slice(out, start, now)
}

// newHoldReplay is the replay's answer before any hold is counted: a zero for
// each story among items that was ever in ready, and for each day from the one
// that holds start to now and each ISO week from the one that holds start to
// this one.
func newHoldReplay(items []*workitem.Item, start, now time.Time) holdReplay {
	out := holdReplay{seconds: map[string]*float64{}}
	for d := firstDay(start); !d.After(now); d = d.AddDate(0, 0, 1) {
		out.days = append(out.days, 0)
	}
	for m := monday(start); !m.After(now); m = m.AddDate(0, 0, 7) {
		out.weeks = append(out.weeks, HeldByCode{})
	}
	for _, it := range items {
		if it.Type == workitem.Story && !it.FirstAt(workitem.Ready).IsZero() {
			out.seconds[it.ID] = new(float64)
		}
	}
	return out
}

// holdTimeline is the hold rules' replay over a history, whatever the window
// and now: the judgement in force from each time it changed, the last open
// until now (S-0608).
type holdTimeline struct {
	segments []holdSegment
	// judged is how many times the holds were judged anew (I-0155).
	judged int
	// applied is the time of the last event the replay applied, zero with
	// none; pending that of the first it left for being after its horizon,
	// zero when it applied every event.
	applied, pending time.Time
}

// holdSegment is the stories held from at until the next segment, by code.
type holdSegment struct {
	at   time.Time
	held []heldEntry
}

// heldEntry is a story held, by its ID, and the code of its hold.
type heldEntry struct {
	id, code string
}

// answers says whether the timeline answers for a request at now: every
// event up to now is applied, and none after it.
func (t *holdTimeline) answers(now time.Time) bool {
	return !now.Before(t.applied) && (t.pending.IsZero() || now.Before(t.pending))
}

// replayTimeline replays the hold rules for every story of all that was ever
// in ready. They are replayed at each creation and transition of a story or
// task not after now, over the states of the time and today's touches and
// after; nothing changes until the next one. A hold counts under the reason it
// is named by, the code Holds.Of gives it. An overlap inside the manifest's
// shared paths holds nothing (ADR-0096). The holds are judged anew only after
// an event that can change one, of a story or of a task of a story in ready or
// in progress, over only the items they read, and with what did not change
// remembered (I-0155). Whether a story is held does not depend on which other
// stories are tracked, so a request's stories are the timeline's filtered.
func replayTimeline(all []*workitem.Item, now time.Time, opt Options) *holdTimeline {
	out := &holdTimeline{}
	tracked := map[string]bool{}
	for _, it := range all {
		if it.Type == workitem.Story && !it.FirstAt(workitem.Ready).IsZero() {
			tracked[it.ID] = true
		}
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
	var held []heldEntry
	// memo keeps what the holds work out, claims and how one holds another,
	// from one judgement to the next. A story's claim is forgotten when a task
	// of it is created or enters or leaves done or cancelled: all Claim reads
	// of them.
	memo := workitem.NewHoldsMemo()
	k := 0
	for k < len(events) && !events[k].at.After(now) {
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
		out.applied = at
		if !judge {
			continue
		}
		var found []heldEntry
		var ready, in []*workitem.Item
		for _, s := range stories {
			if s.Status == workitem.Ready && tracked[s.ID] {
				ready = append(ready, s)
			}
		}
		if len(ready) > 0 {
			// Only the stories in ready or in progress and their tasks, and
			// the stories the ready ones name in after, are read.
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
					found = append(found, heldEntry{id: s.ID, code: hold.Code})
				}
			}
		}
		// A judgement that holds what the one in force holds starts nothing.
		if !slices.Equal(found, held) {
			held = found
			out.segments = append(out.segments, holdSegment{at: at, held: held})
		}
	}
	if k < len(events) {
		out.pending = events[k].at
	}
	return out
}

// slice counts into out, made by newHoldReplay for the window from start to
// now, what the timeline held of out's stories: the seconds each was held,
// the stories held at the end of each day, and the seconds held by code in
// each ISO week. A judgement holds from its time to the next one's, the last
// to now; a day's stories are counted at its last second, today's at now.
func (t *holdTimeline) slice(out holdReplay, start, now time.Time) holdReplay {
	out.judged = t.judged
	day0, week0 := firstDay(start), monday(start)
	dayEnd := func(i int) time.Time {
		end := day0.AddDate(0, 0, i).Add(24*time.Hour - time.Second)
		if end.After(now) {
			return now
		}
		return end
	}
	const week = 7 * 24 * time.Hour
	day := 0
	for i, seg := range t.segments {
		at := seg.at
		if at.After(now) {
			break
		}
		// The judgement holds from at up to next; until is next, not after now.
		until, next := now, now.Add(time.Second)
		if i+1 < len(t.segments) && !t.segments[i+1].at.After(now) {
			until, next = t.segments[i+1].at, t.segments[i+1].at
		}
		byCode := map[string]int{}
		for _, h := range seg.held {
			if s := out.seconds[h.id]; s != nil {
				*s += until.Sub(at).Seconds()
				byCode[h.code]++
			}
		}
		// Only the weeks from the one that holds at to the one that holds
		// until can take any of it.
		first := 0
		if at.After(week0) {
			first = int(at.Sub(week0) / week)
		}
		for code, n := range byCode {
			for w := first; w < len(out.weeks); w++ {
				wk := week0.AddDate(0, 0, 7*w)
				if !wk.Before(until) {
					break
				}
				from, to := laterOf(at, wk), earlierOf(until, wk.AddDate(0, 0, 7))
				if to.After(from) {
					out.weeks[w].add(code, float64(n)*to.Sub(from).Seconds())
				}
			}
		}
		for day < len(out.days) && dayEnd(day).Before(at) {
			day++ // counted under earlier judgements, or before any
		}
		for ; day < len(out.days) && dayEnd(day).Before(next); day++ {
			for _, n := range byCode {
				out.days[day] += n
			}
		}
	}
	return out
}

// kept holds, per repository root, the last hold timeline and touches drift
// Compute worked out, each with the keys it was worked out under, so that a
// request while they hold works out neither again (S-0608). A root holds one
// of each, replaced when its keys move.
var kept = struct {
	sync.Mutex
	holds map[string]keptHolds
	drift map[string]keptDrift
}{holds: map[string]keptHolds{}, drift: map[string]keptDrift{}}

// keptHolds is a hold timeline and its keys.
type keptHolds struct {
	generation uint64
	config     string
	timeline   *holdTimeline
}

// keptDrift is the touches drift of every story with commits and its keys.
type keptDrift struct {
	generation uint64
	commits    string
	config     string
	drift      []TouchDrift
}

// configKey is the manifest's projects and shared paths as the replay and the
// drift read them, and flai's folders the drift leaves out, "" when they cannot be written down, which keeps nothing.
func configKey(opt Options) string {
	b, err := json.Marshal(struct {
		Projects    []manifest.Project
		Shared      manifest.Claims
		FlaiFolders []string
	}{opt.Projects, opt.Shared, opt.FlaiFolders})
	if err != nil {
		return ""
	}
	return string(b)
}

// keptTimeline is the hold timeline of all up to now: the one kept for
// opt.Root while the generation and the manifest's projects and shared paths
// are those it was replayed under and it answers for now, else a new replay,
// kept in its place unless a newer generation's is kept already. replayed
// says whether it replayed. With no root, generation, or config key, it
// replays and keeps nothing.
func keptTimeline(all []*workitem.Item, now time.Time, opt Options) (t *holdTimeline, replayed bool) {
	config := configKey(opt)
	if opt.Root == "" || opt.ItemsGeneration == 0 || config == "" {
		return replayTimeline(all, now, opt), true
	}
	kept.Lock()
	k, ok := kept.holds[opt.Root]
	kept.Unlock()
	if ok && k.generation == opt.ItemsGeneration && k.config == config && k.timeline.answers(now) {
		return k.timeline, false
	}
	t = replayTimeline(all, now, opt)
	kept.Lock()
	if cur, ok := kept.holds[opt.Root]; !ok || cur.generation <= opt.ItemsGeneration {
		kept.holds[opt.Root] = keptHolds{generation: opt.ItemsGeneration, config: config, timeline: t}
	}
	kept.Unlock()
	return t, true
}

// driftFor is the touches drift of the stories among items, in ID order:
// touchDrift over every story of all, kept for opt.Root while the generation,
// the commits' key, and the manifest's projects are those it was worked out
// under, and filtered to the stories among items. The drift is shared with
// the requests that keep it, which only read it.
func driftFor(items, all []*workitem.Item, opt Options) []TouchDrift {
	every, _ := keptTouchDrift(all, opt)
	ids := map[string]bool{}
	for _, it := range items {
		if it.Type == workitem.Story {
			ids[it.ID] = true
		}
	}
	out := []TouchDrift{}
	for _, d := range every {
		if ids[d.ID] {
			out = append(out, d)
		}
	}
	return out
}

// keptTouchDrift is the touches drift of every story of all, in ID order: the
// one kept for opt.Root while its keys hold, else worked out again and kept
// in its place unless a newer generation's is kept already. worked says
// whether it was worked out. With no root, generation, commits' key, or
// config key, it is worked out and nothing is kept.
func keptTouchDrift(all []*workitem.Item, opt Options) (drift []TouchDrift, worked bool) {
	config := configKey(opt)
	if opt.Root == "" || opt.ItemsGeneration == 0 || opt.CommitsKey == "" || config == "" {
		return touchDrift(all, all, opt), true
	}
	kept.Lock()
	k, ok := kept.drift[opt.Root]
	kept.Unlock()
	if ok && k.generation == opt.ItemsGeneration && k.commits == opt.CommitsKey && k.config == config {
		return k.drift, false
	}
	drift = touchDrift(all, all, opt)
	kept.Lock()
	if cur, ok := kept.drift[opt.Root]; !ok || cur.generation <= opt.ItemsGeneration {
		kept.drift[opt.Root] = keptDrift{generation: opt.ItemsGeneration, commits: opt.CommitsKey, config: config, drift: drift}
	}
	kept.Unlock()
	return drift, true
}
