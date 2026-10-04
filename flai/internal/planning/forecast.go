package planning

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A story's size is the checkboxes of its acceptance criteria and its own
// touches counted together, at least 1. A size up to SmallMax is in the
// small band, up to MediumMax in the medium band, and above it in the large.
const (
	SmallMax  = 6
	MediumMax = 12
)

// MinHistory is how many done stories a rung of the match needs for its
// medians to be used; below it the match falls back a rung.
const MinHistory = 3

// Result is a story's forecast and what it rests on.
type Result struct {
	ID string `json:"id"`
	// Duration is the work the story is expected to take, a Go duration
	// without zero tails such as 1h5m.
	Duration        string `json:"duration"`
	DurationSeconds int64  `json:"duration_seconds"`
	// Delivery is when the story is expected to be done, a UTC timestamp.
	Delivery string `json:"delivery"`
	// Basis is what the forecast rests on, in one sentence.
	Basis string `json:"basis"`
	// Default says too little history matched, so the duration is
	// planning.default_duration.
	Default  bool    `json:"default"`
	Size     int     `json:"size"`
	Criteria int     `json:"criteria"`
	Touches  int     `json:"touches"`
	History  History `json:"history"`
	// Position is the story's place in the pull order, ready stories then
	// backlog ones, from 1; 0 when it is in progress or in review.
	Position int `json:"position"`
	// Limit is the in-progress limit the pull order was played out with; 0
	// means none.
	Limit int `json:"limit"`
	// Ahead are the stories in progress and those before it in the pull
	// order, as they were played out.
	Ahead []Ahead `json:"ahead"`
}

// History is the done stories the duration was worked out from.
type History struct {
	// Match is what they share with the story, such as "nature and size
	// band"; "none" when too few matched.
	Match   string `json:"match"`
	Stories int    `json:"stories"`
	// SecondsPerUnit is their median agent seconds per unit of size.
	SecondsPerUnit float64 `json:"seconds_per_unit"`
	// BusyFactor is their median of seconds in progress per agent second,
	// at least 1: how long a story holds its place in progress.
	BusyFactor float64 `json:"busy_factor"`
	// CycleFactor is their median of cycle time per agent second, at least
	// 1: how long from the start of progress to done.
	CycleFactor float64 `json:"cycle_factor"`
}

// Ahead is a story played out before the one forecast.
type Ahead struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	// Duration is the work it was played out with.
	Duration string `json:"duration"`
	// FromForecast says the duration is the story's own forecast.duration
	// rather than one worked out from history.
	FromForecast bool   `json:"from_forecast"`
	Start        string `json:"start"`
	Delivery     string `json:"delivery"`
}

// Forecast works out how long story id will take from the done stories in
// items with usage, and when it will be delivered by playing out the board's
// pull order and in-progress limit from now. items includes archived ones.
func Forecast(items []*workitem.Item, board *workitem.Board, plan manifest.Planning, id string, now time.Time) (Result, error) {
	return playOut(items, board, plan, id, now, false)
}

// Replay is Forecast for a story that keeps its own forecast.duration: only
// its delivery and basis are worked out again. A story without a valid one
// replays as Forecast gives it.
func Replay(items []*workitem.Item, board *workitem.Board, plan manifest.Planning, id string, now time.Time) (Result, error) {
	return playOut(items, board, plan, id, now, true)
}

// playOut forecasts story id, played out with its own forecast.duration when
// keep is set and it has a valid one, else with its estimate.
func playOut(items []*workitem.Item, board *workitem.Board, plan manifest.Planning, id string, now time.Time, keep bool) (Result, error) {
	fallback, err := plan.FallbackDuration()
	if err != nil {
		return Result{}, fmt.Errorf("cannot forecast %s: %w", id, err)
	}
	id = workitem.CanonicalID(id)
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	target := byID[id]
	switch {
	case target == nil:
		return Result{}, fmt.Errorf("cannot forecast %s: there is no such item; check the ID with flai board", id)
	case target.Type != workitem.Story:
		return Result{}, fmt.Errorf("cannot forecast %s: it is a %s, and only a story has a forecast", id, target.Type)
	case target.Archived:
		return Result{}, fmt.Errorf("cannot forecast %s: it is archived; forecast a story that is still open", id)
	case target.Closed():
		return Result{}, fmt.Errorf("cannot forecast %s: it is %s; forecast a story that is still open", id, target.Status)
	}
	now = now.UTC()
	f := newForecaster(items, fallback)
	est := f.estimate(target)
	dur, basis, def := est.duration, est.basis, est.def
	if d, ok := ownDuration(target); keep && ok {
		dur, basis, def = d, "Its own forecast of "+FormatDuration(d), false
	}
	var order []string
	limit := 0
	if board != nil {
		order, limit = board.Order, max(board.WIPLimits[workitem.InProgress], 0)
	}
	res := Result{
		ID: id, Duration: FormatDuration(dur), DurationSeconds: int64(dur / time.Second),
		Default: def, Size: est.size, Criteria: est.criteria, Touches: est.touches,
		History: est.history, Limit: limit, Ahead: []Ahead{},
	}
	if target.Status == workitem.InProgress || target.Status == workitem.Review {
		start := startedAt(target, now)
		res.Delivery = stamp(later(now, start.Add(scale(dur, est.cycle))))
		res.Basis = fmt.Sprintf("%s; %s since it started at %s, so delivery counts from then.", basis, phrase(target.Status), stamp(start))
		return res, nil
	}

	delivered := map[string]time.Time{}
	var releases []time.Time
	for _, sid := range workitem.PullSequence(order, items, workitem.InProgress) {
		it := byID[sid]
		d, fromForecast := f.durationOf(it)
		e := f.estimate(it)
		start := startedAt(it, now)
		releases = append(releases, later(now, start.Add(scale(d, e.busy))))
		delivered[sid] = later(now, start.Add(scale(d, e.cycle)))
		res.Ahead = append(res.Ahead, Ahead{ID: sid, Status: it.Status, Duration: FormatDuration(d), FromForecast: fromForecast, Start: stamp(start), Delivery: stamp(delivered[sid])})
	}
	free := newLanes(limit, releases, now)
	queue := append(workitem.PullSequence(order, items, workitem.Ready), workitem.PullSequence(order, items, workitem.Backlog)...)
	for i, sid := range queue {
		it := byID[sid]
		e := f.estimate(it)
		d, fromForecast := dur, false
		if sid != id {
			d, fromForecast = f.durationOf(it)
		}
		lane, start := free.earliest()
		for _, dep := range it.After {
			if t, ok := delivered[dep]; ok && t.After(start) {
				start = t
			}
		}
		free.hold(lane, start.Add(scale(d, e.busy)))
		delivered[sid] = start.Add(scale(d, e.cycle))
		if sid == id {
			res.Position = i + 1
			res.Delivery = stamp(delivered[sid])
			res.Basis = fmt.Sprintf("%s; %s in the pull order with %s, %s.", basis, ordinal(i+1), limitPhrase(limit), behind(res.Ahead))
			break
		}
		res.Ahead = append(res.Ahead, Ahead{ID: sid, Status: it.Status, Duration: FormatDuration(d), FromForecast: fromForecast, Start: stamp(start), Delivery: stamp(delivered[sid])})
	}
	return res, nil
}

// FormatDuration writes d as a Go duration without zero tails: 2h, 1h5m,
// 45m.
func FormatDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

// sample is what one done story says about how long work takes.
type sample struct {
	nature, model, band  string
	perUnit, busy, cycle float64
}

// estimate is a story's duration worked out from history, and on what.
type estimate struct {
	size, criteria, touches int
	duration                time.Duration
	busy, cycle             float64
	def                     bool
	history                 History
	basis                   string
}

type forecaster struct {
	samples   []sample
	fallback  time.Duration
	estimates map[string]estimate
}

func newForecaster(items []*workitem.Item, fallback time.Duration) *forecaster {
	f := &forecaster{fallback: fallback, estimates: map[string]estimate{}}
	for _, it := range items {
		if s, ok := sampleOf(it); ok {
			f.samples = append(f.samples, s)
		}
	}
	return f
}

// sampleOf reads a done story with usage and a cycle time; false for any
// other item.
func sampleOf(it *workitem.Item) (sample, bool) {
	if it.Type != workitem.Story || it.Status != workitem.Done || it.Usage == nil || it.Usage.Seconds <= 0 {
		return sample{}, false
	}
	started, completed := it.FirstAt(workitem.InProgress), it.CompletedAt()
	if started.IsZero() || !completed.After(started) {
		return sample{}, false
	}
	used := float64(it.Usage.Seconds)
	_, _, size := sizeOf(it)
	return sample{
		nature: it.Nature, model: modelOf(it), band: band(size),
		perUnit: used / float64(size), busy: inProgressSeconds(it) / used, cycle: completed.Sub(started).Seconds() / used,
	}, true
}

// rung is one step of the match ladder: what a sample must share with the
// story, and how the stories that do are described.
type rung struct {
	match string
	keep  func(s, k sample) bool
	over  func(n int, k sample) string
}

var ladder = []rung{
	{"nature, model, and size band", func(s, k sample) bool { return s.nature == k.nature && s.model == k.model && s.band == k.band },
		func(n int, k sample) string {
			return fmt.Sprintf("%d done %s stories %s in the %s band", n, k.nature, onModel(k.model), k.band)
		}},
	{"nature and size band", func(s, k sample) bool { return s.nature == k.nature && s.band == k.band },
		func(n int, k sample) string {
			return fmt.Sprintf("%d done %s stories in the %s band", n, k.nature, k.band)
		}},
	{"nature", func(s, k sample) bool { return s.nature == k.nature },
		func(n int, k sample) string { return fmt.Sprintf("%d done %s stories", n, k.nature) }},
	{"all done stories", func(s, k sample) bool { return true },
		func(n int, k sample) string { return fmt.Sprintf("%d done stories", n) }},
}

// estimate works out a story's duration from the first rung of the ladder
// with MinHistory samples, or gives the fallback when none has.
func (f *forecaster) estimate(it *workitem.Item) estimate {
	if e, ok := f.estimates[it.ID]; ok {
		return e
	}
	criteria, touches, size := sizeOf(it)
	e := estimate{size: size, criteria: criteria, touches: touches}
	sized := fmt.Sprintf("size %d (%d criteria, %d touches)", size, criteria, touches)
	key := sample{nature: it.Nature, model: modelOf(it), band: band(size)}
	for _, r := range ladder {
		var perUnit, busy, cycle []float64
		for _, s := range f.samples {
			if r.keep(s, key) {
				perUnit, busy, cycle = append(perUnit, s.perUnit), append(busy, s.busy), append(cycle, s.cycle)
			}
		}
		if len(perUnit) < MinHistory {
			continue
		}
		spu := median(perUnit)
		e.duration = time.Duration(math.Ceil(spu*float64(size)/60-1e-9)) * time.Minute
		e.busy, e.cycle = math.Max(1, median(busy)), math.Max(1, median(cycle))
		e.history = History{Match: r.match, Stories: len(perUnit), SecondsPerUnit: spu, BusyFactor: e.busy, CycleFactor: e.cycle}
		e.basis = fmt.Sprintf("Median %.0f s per unit of size over %s, times %s", spu, r.over(len(perUnit), key), sized)
		f.estimates[it.ID] = e
		return e
	}
	e.duration, e.busy, e.cycle, e.def = f.fallback, 1, 1, true
	e.history = History{Match: "none", Stories: len(f.samples), BusyFactor: 1, CycleFactor: 1}
	e.basis = fmt.Sprintf("Not enough history, %d done stories with usage where %d are needed, so planning.default_duration of %s for %s", len(f.samples), MinHistory, FormatDuration(f.fallback), sized)
	f.estimates[it.ID] = e
	return e
}

// durationOf is the duration a story ahead is played out with: its own
// forecast.duration when it has a valid one, else its estimate.
func (f *forecaster) durationOf(it *workitem.Item) (time.Duration, bool) {
	if d, ok := ownDuration(it); ok {
		return d, true
	}
	return f.estimate(it).duration, false
}

// ownDuration is a story's own forecast.duration; false when it has none or
// it is not a positive duration.
func ownDuration(it *workitem.Item) (time.Duration, bool) {
	if it.Forecast != nil && it.Forecast.Duration != "" {
		if d, err := time.ParseDuration(it.Forecast.Duration); err == nil && d > 0 {
			return d, true
		}
	}
	return 0, false
}

// sizeOf is a story's criteria, its own touches, and its size, at least 1.
func sizeOf(it *workitem.Item) (criteria, touches, size int) {
	criteria, touches = workitem.CriteriaCount(it.Body), len(it.Touches)
	return criteria, touches, max(1, criteria+touches)
}

func band(size int) string {
	switch {
	case size <= SmallMax:
		return "small"
	case size <= MediumMax:
		return "medium"
	}
	return "large"
}

func modelOf(it *workitem.Item) string {
	if it.Agent == nil {
		return ""
	}
	return it.Agent.Model
}

func onModel(model string) string {
	if model == "" {
		return "with no model recorded"
	}
	return "on " + model
}

// inProgressSeconds is the time an item spent in progress, over every spell.
func inProgressSeconds(it *workitem.Item) float64 {
	var total time.Duration
	var since time.Time
	for _, tr := range it.Transitions {
		at, err := time.Parse(workitem.TimeFormat, tr.At)
		if err != nil {
			continue
		}
		if !since.IsZero() {
			total += at.Sub(since)
			since = time.Time{}
		}
		if tr.To == workitem.InProgress {
			since = at
		}
	}
	return total.Seconds()
}

// startedAt is when an item first went in progress, as metrics counts cycle
// time from; now when it has no such transition.
func startedAt(it *workitem.Item, now time.Time) time.Time {
	if t := it.FirstAt(workitem.InProgress); !t.IsZero() {
		return t
	}
	return now
}

func median(v []float64) float64 {
	s := append([]float64(nil), v...)
	sort.Float64s(s)
	n := len(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2
}

func scale(d time.Duration, factor float64) time.Duration {
	return time.Duration(float64(d) * factor)
}

func later(a, b time.Time) time.Time {
	if b.After(a) {
		return b
	}
	return a
}

// stamp writes t as a UTC timestamp, rounded up to the minute.
func stamp(t time.Time) string {
	m := t.UTC().Truncate(time.Minute)
	if m.Before(t) {
		m = m.Add(time.Minute)
	}
	return m.Format(workitem.TimeFormat)
}

// lanes are the places in progress the pull order is played out over, each
// with when it is next free; none means there is no limit.
type lanes struct {
	free []time.Time
	now  time.Time
}

// newLanes gives limit lanes, those the stories in progress hold freed as
// they finish. Over the limit, a lane frees only once enough of them have
// finished to bring the count under it, so the lanes are the last limit of
// them to finish.
func newLanes(limit int, busy []time.Time, now time.Time) *lanes {
	if limit <= 0 {
		return &lanes{now: now}
	}
	free := append([]time.Time(nil), busy...)
	sort.Slice(free, func(i, j int) bool { return free[i].Before(free[j]) })
	if len(free) > limit {
		free = free[len(free)-limit:]
	}
	for len(free) < limit {
		free = append(free, now)
	}
	return &lanes{free: free, now: now}
}

// earliest is the lane that frees first and when; -1 and now with no limit.
func (l *lanes) earliest() (int, time.Time) {
	if len(l.free) == 0 {
		return -1, l.now
	}
	i := 0
	for j, t := range l.free {
		if t.Before(l.free[i]) {
			i = j
		}
	}
	return i, l.free[i]
}

// hold marks a lane busy until a time; a no-op with no limit.
func (l *lanes) hold(i int, until time.Time) {
	if i >= 0 {
		l.free[i] = until
	}
}

func phrase(status string) string {
	if status == workitem.Review {
		return "in review"
	}
	return "in progress"
}

func limitPhrase(limit int) string {
	if limit <= 0 {
		return "no in-progress limit"
	}
	return fmt.Sprintf("an in-progress limit of %d", limit)
}

func behind(ahead []Ahead) string {
	if len(ahead) == 0 {
		return "with nothing ahead of it"
	}
	ids := make([]string, len(ahead))
	for i, a := range ahead {
		ids[i] = a.ID
	}
	if len(ids) == 1 {
		return "behind " + ids[0]
	}
	return "behind " + strings.Join(ids[:len(ids)-1], ", ") + " and " + ids[len(ids)-1]
}

// ordinal writes n as 1st, 2nd, 3rd, 4th, 11th, 21st.
func ordinal(n int) string {
	suffix := "th"
	if n%100 < 11 || n%100 > 13 {
		switch n % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%d%s", n, suffix)
}
