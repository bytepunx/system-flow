// Package metrics is the reference implementation of design/system/metrics.md.
// flaiover ports these definitions and is fixture-tested against this output.
package metrics

import (
	"math"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Options select and window the computation.
type Options struct {
	Now   time.Time
	Since time.Duration // window for completed items, default 30 days
	Type  string        // item type to aggregate, default story
	By    string        // optional grouping: nature, type, parent
}

// ItemMetrics are the per-item derived values.
type ItemMetrics struct {
	ID        string             `json:"id"`
	Type      string             `json:"type"`
	Nature    string             `json:"nature"`
	Title     string             `json:"title"`
	Status    string             `json:"status"`
	Parent    string             `json:"parent,omitempty"`
	Created   string             `json:"created"`
	Committed string             `json:"committed,omitempty"`
	Started   string             `json:"started,omitempty"`
	Completed string             `json:"completed,omitempty"`
	LeadTime  *float64           `json:"lead_time_seconds,omitempty"`
	CycleTime *float64           `json:"cycle_time_seconds,omitempty"`
	QueueTime *float64           `json:"queue_time_seconds,omitempty"`
	Blocked   float64            `json:"blocked_seconds"`
	InState   map[string]float64 `json:"time_in_state_seconds"`
	Estimate  *float64           `json:"estimate_seconds,omitempty"`
	EstError  *float64           `json:"estimate_error,omitempty"`
	Age       *float64           `json:"age_seconds,omitempty"` // active items: now - started
}

// Distribution summarises a set of durations in seconds.
type Distribution struct {
	Count int     `json:"count"`
	P50   float64 `json:"p50_seconds"`
	P85   float64 `json:"p85_seconds"`
	Max   float64 `json:"max_seconds"`
	Mean  float64 `json:"mean_seconds"`
}

// Summary holds the aggregates for one group.
type Summary struct {
	Group            string             `json:"group,omitempty"`
	Completed        int                `json:"completed"`
	Cancelled        int                `json:"cancelled"`
	CancellationRate float64            `json:"cancellation_rate"`
	ThroughputWeek   float64            `json:"throughput_per_week"`
	WIP              int                `json:"wip"`
	CycleTime        Distribution       `json:"cycle_time"`
	LeadTime         Distribution       `json:"lead_time"`
	QueueTime        Distribution       `json:"queue_time"`
	FlowEfficiency   float64            `json:"flow_efficiency"`
	InStateShare     map[string]float64 `json:"time_in_state_share"`
}

// WeekBucket is throughput for one ISO week.
type WeekBucket struct {
	Week  string         `json:"week"` // 2026-W38
	Start string         `json:"start"`
	Done  int            `json:"done"`
	ByKey map[string]int `json:"by_nature"`
}

// DayPoint is one day of a series.
type DayPoint struct {
	Date   string         `json:"date"`
	Scope  int            `json:"scope,omitempty"`
	Done   int            `json:"done,omitempty"`
	Counts map[string]int `json:"counts,omitempty"`
}

// AgingItem is an active item compared with the cycle time p85.
type AgingItem struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Status   string  `json:"status"`
	Age      float64 `json:"age_seconds"`
	OverP85  bool    `json:"over_p85"`
	Blocked  bool    `json:"blocked"`
	Nature   string  `json:"nature"`
	Parent   string  `json:"parent,omitempty"`
	AgeHuman string  `json:"age"`
}

// Report is the full output of Compute.
type Report struct {
	GeneratedAt string                `json:"generated_at"`
	Type        string                `json:"type"`
	WindowDays  float64               `json:"window_days"`
	WindowStart string                `json:"window_start"`
	Summary     Summary               `json:"summary"`
	Groups      []Summary             `json:"groups,omitempty"`
	Items       []ItemMetrics         `json:"items"`
	Throughput  []WeekBucket          `json:"throughput"`
	Burnup      map[string][]DayPoint `json:"burnup"`
	CFD         []DayPoint            `json:"cfd"`
	Aging       []AgingItem           `json:"aging"`
}

// Compute derives every metric from the items.
func Compute(all []*workitem.Item, opt Options) *Report {
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}
	opt.Now = opt.Now.UTC()
	if opt.Since <= 0 {
		opt.Since = 30 * 24 * time.Hour
	}
	if opt.Type == "" {
		opt.Type = workitem.Story
	}
	start := opt.Now.Add(-opt.Since)

	var items []*workitem.Item
	for _, it := range all {
		if it.Type == opt.Type {
			items = append(items, it)
		}
	}
	rep := &Report{
		GeneratedAt: opt.Now.Format(workitem.TimeFormat), Type: opt.Type,
		WindowDays: opt.Since.Hours() / 24, WindowStart: start.Format(workitem.TimeFormat),
		Burnup: map[string][]DayPoint{},
	}
	perItem := map[string]ItemMetrics{}
	for _, it := range items {
		m := Derive(it, opt.Now)
		perItem[it.ID] = m
		rep.Items = append(rep.Items, m)
	}
	inWindow := func(it *workitem.Item) bool {
		c := it.FirstAt(workitem.Done)
		if c.IsZero() {
			c = it.FirstAt(workitem.Cancelled)
		}
		return !c.IsZero() && !c.Before(start) && !c.After(opt.Now)
	}
	rep.Summary = summarise("", items, perItem, inWindow, opt)
	if opt.By != "" {
		groups := map[string][]*workitem.Item{}
		var keys []string
		for _, it := range items {
			k := groupKey(it, opt.By)
			if _, ok := groups[k]; !ok {
				keys = append(keys, k)
			}
			groups[k] = append(groups[k], it)
		}
		sort.Strings(keys)
		for _, k := range keys {
			rep.Groups = append(rep.Groups, summarise(k, groups[k], perItem, inWindow, opt))
		}
	}
	rep.Throughput = throughput(items, start, opt.Now)
	rep.Burnup["all"] = burnup(items, opt.Now)
	parents := map[string]bool{}
	for _, it := range items {
		if it.Parent != "" && !parents[it.Parent] {
			parents[it.Parent] = true
			var sub []*workitem.Item
			for _, x := range items {
				if x.Parent == it.Parent {
					sub = append(sub, x)
				}
			}
			rep.Burnup[it.Parent] = burnup(sub, opt.Now)
		}
	}
	rep.CFD = cfd(items, opt.Now)
	rep.Aging = aging(items, perItem, rep.Summary.CycleTime.P85)
	return rep
}

// Derive computes the per-item values for one item.
func Derive(it *workitem.Item, now time.Time) ItemMetrics {
	m := ItemMetrics{ID: it.ID, Type: it.Type, Nature: it.Nature, Title: it.Title, Status: it.Status, Parent: it.Parent, Created: it.Created, InState: map[string]float64{}}
	created, _ := time.Parse(workitem.TimeFormat, it.Created)
	committed := it.FirstAt(workitem.Ready)
	started := it.FirstAt(workitem.InProgress)
	completed := it.FirstAt(workitem.Done)
	if completed.IsZero() {
		completed = it.FirstAt(workitem.Cancelled)
	}
	set := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format(workitem.TimeFormat)
	}
	m.Committed, m.Started, m.Completed = set(committed), set(started), set(completed)

	// state intervals: backlog from created, then each transition to the next or now
	prevState, prevAt := workitem.Backlog, created
	for _, tr := range it.Transitions {
		at, _ := time.Parse(workitem.TimeFormat, tr.At)
		m.InState[prevState] += at.Sub(prevAt).Seconds()
		prevState, prevAt = tr.To, at
	}
	if !it.Closed() {
		m.InState[prevState] += now.Sub(prevAt).Seconds()
	}
	for _, b := range it.Blocked {
		from, _ := time.Parse(workitem.TimeFormat, b.From)
		until := now
		if b.Until != "" {
			until, _ = time.Parse(workitem.TimeFormat, b.Until)
		}
		if until.After(from) {
			m.Blocked += until.Sub(from).Seconds()
		}
	}
	if !completed.IsZero() {
		m.LeadTime = secs(completed.Sub(created))
		if !started.IsZero() {
			m.CycleTime = secs(completed.Sub(started))
		}
	}
	if !started.IsZero() && !committed.IsZero() && !started.Before(committed) {
		m.QueueTime = secs(started.Sub(committed))
	}
	if it.Estimate != "" {
		if d, err := time.ParseDuration(it.Estimate); err == nil && d > 0 {
			m.Estimate = secs(d)
			if m.CycleTime != nil {
				e := (*m.CycleTime - d.Seconds()) / d.Seconds()
				m.EstError = &e
			}
		}
	}
	if !it.Closed() && !started.IsZero() {
		m.Age = secs(now.Sub(started))
	}
	return m
}

func summarise(group string, items []*workitem.Item, per map[string]ItemMetrics, inWindow func(*workitem.Item) bool, opt Options) Summary {
	s := Summary{Group: group, InStateShare: map[string]float64{}}
	var cycle, lead, queue []float64
	var eff []float64
	inState := map[string]float64{}
	var leadTotal float64
	for _, it := range items {
		if it.Status == workitem.InProgress || it.Status == workitem.Review {
			s.WIP++
		}
		if !inWindow(it) {
			continue
		}
		if it.Status == workitem.Cancelled {
			s.Cancelled++
			continue
		}
		s.Completed++
		m := per[it.ID]
		if m.LeadTime != nil {
			lead = append(lead, *m.LeadTime)
			leadTotal += *m.LeadTime
			for st, v := range m.InState {
				inState[st] += v
			}
		}
		if m.CycleTime != nil {
			cycle = append(cycle, *m.CycleTime)
			if *m.CycleTime > 0 {
				eff = append(eff, (*m.CycleTime-m.Blocked) / *m.CycleTime)
			}
		}
		if m.QueueTime != nil {
			queue = append(queue, *m.QueueTime)
		}
	}
	if s.Completed+s.Cancelled > 0 {
		s.CancellationRate = float64(s.Cancelled) / float64(s.Completed+s.Cancelled)
	}
	s.ThroughputWeek = float64(s.Completed) / (opt.Since.Hours() / 24 / 7)
	s.CycleTime, s.LeadTime, s.QueueTime = distribution(cycle), distribution(lead), distribution(queue)
	if len(eff) > 0 {
		s.FlowEfficiency = mean(eff)
	}
	if leadTotal > 0 {
		for st, v := range inState {
			s.InStateShare[st] = v / leadTotal
		}
	}
	return s
}

func distribution(v []float64) Distribution {
	if len(v) == 0 {
		return Distribution{}
	}
	sorted := append([]float64{}, v...)
	sort.Float64s(sorted)
	return Distribution{Count: len(v), P50: percentile(sorted, 50), P85: percentile(sorted, 85), Max: sorted[len(sorted)-1], Mean: mean(v)}
}

// percentile uses nearest-rank on a sorted slice.
func percentile(sorted []float64, p float64) float64 {
	rank := int(math.Ceil(p / 100 * float64(len(sorted))))
	if rank < 1 {
		rank = 1
	}
	return sorted[rank-1]
}

func mean(v []float64) float64 {
	var t float64
	for _, x := range v {
		t += x
	}
	return t / float64(len(v))
}

func secs(d time.Duration) *float64 {
	s := d.Seconds()
	return &s
}

func groupKey(it *workitem.Item, by string) string {
	switch by {
	case "nature":
		return it.Nature
	case "type":
		return it.Type
	case "parent":
		if it.Parent == "" {
			return "(none)"
		}
		return it.Parent
	}
	return ""
}

func throughput(items []*workitem.Item, start, now time.Time) []WeekBucket {
	buckets := map[string]*WeekBucket{}
	for _, it := range items {
		done := it.FirstAt(workitem.Done)
		if done.IsZero() || done.Before(start) || done.After(now) {
			continue
		}
		y, w := done.ISOWeek()
		key := weekKey(y, w)
		b, ok := buckets[key]
		if !ok {
			// Monday of that ISO week
			monday := done.AddDate(0, 0, -((int(done.Weekday()) + 6) % 7))
			b = &WeekBucket{Week: key, Start: monday.Format("2006-01-02"), ByKey: map[string]int{}}
			buckets[key] = b
		}
		b.Done++
		b.ByKey[it.Nature]++
	}
	var out []WeekBucket
	for _, b := range buckets {
		out = append(out, *b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Week < out[j].Week })
	return out
}

func weekKey(y, w int) string {
	return time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006") + "-W" + pad2(w)
}

func pad2(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

func dayRange(items []*workitem.Item, now time.Time) (time.Time, time.Time) {
	var first time.Time
	for _, it := range items {
		c, _ := time.Parse(workitem.TimeFormat, it.Created)
		if first.IsZero() || c.Before(first) {
			first = c
		}
	}
	if first.IsZero() {
		return time.Time{}, time.Time{}
	}
	first = time.Date(first.Year(), first.Month(), first.Day(), 0, 0, 0, 0, time.UTC)
	last := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return first, last
}

// stateAt returns the item's state at the end of day t, or "" if not yet created.
func stateAt(it *workitem.Item, endOfDay time.Time) string {
	created, _ := time.Parse(workitem.TimeFormat, it.Created)
	if created.After(endOfDay) {
		return ""
	}
	state := workitem.Backlog
	for _, tr := range it.Transitions {
		at, _ := time.Parse(workitem.TimeFormat, tr.At)
		if at.After(endOfDay) {
			break
		}
		state = tr.To
	}
	return state
}

func burnup(items []*workitem.Item, now time.Time) []DayPoint {
	first, last := dayRange(items, now)
	if first.IsZero() {
		return nil
	}
	var out []DayPoint
	for d := first; !d.After(last); d = d.AddDate(0, 0, 1) {
		end := d.Add(24*time.Hour - time.Second)
		p := DayPoint{Date: d.Format("2006-01-02")}
		for _, it := range items {
			st := stateAt(it, end)
			if st == "" || st == workitem.Cancelled {
				continue
			}
			p.Scope++
			if st == workitem.Done {
				p.Done++
			}
		}
		out = append(out, p)
	}
	return out
}

func cfd(items []*workitem.Item, now time.Time) []DayPoint {
	first, last := dayRange(items, now)
	if first.IsZero() {
		return nil
	}
	var out []DayPoint
	for d := first; !d.After(last); d = d.AddDate(0, 0, 1) {
		end := d.Add(24*time.Hour - time.Second)
		p := DayPoint{Date: d.Format("2006-01-02"), Counts: map[string]int{}}
		for _, st := range workitem.States {
			p.Counts[st] = 0
		}
		for _, it := range items {
			if st := stateAt(it, end); st != "" {
				p.Counts[st]++
			}
		}
		out = append(out, p)
	}
	return out
}

func aging(items []*workitem.Item, per map[string]ItemMetrics, p85 float64) []AgingItem {
	var out []AgingItem
	for _, it := range items {
		if it.Status != workitem.InProgress && it.Status != workitem.Review {
			continue
		}
		m := per[it.ID]
		age := 0.0
		if m.Age != nil {
			age = *m.Age
		}
		out = append(out, AgingItem{ID: it.ID, Title: it.Title, Status: it.Status, Age: age, OverP85: p85 > 0 && age > p85,
			Blocked: it.IsBlocked(), Nature: it.Nature, Parent: it.Parent, AgeHuman: Human(age)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Age > out[j].Age })
	return out
}

// Human renders seconds as a short duration: 45m, 3h, 2d5h.
func Human(seconds float64) string {
	d := time.Duration(seconds * float64(time.Second))
	switch {
	case d < time.Minute:
		return "0m"
	case d < time.Hour:
		return itoa(int(d.Minutes())) + "m"
	case d < 24*time.Hour:
		h := int(d.Hours())
		mm := int(d.Minutes()) % 60
		if mm == 0 {
			return itoa(h) + "h"
		}
		return itoa(h) + "h" + itoa(mm) + "m"
	default:
		days := int(d.Hours()) / 24
		h := int(d.Hours()) % 24
		if h == 0 {
			return itoa(days) + "d"
		}
		return itoa(days) + "d" + itoa(h) + "h"
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
