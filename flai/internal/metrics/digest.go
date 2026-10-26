package metrics

import (
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Digest is, per ISO week of the window, the threads opened to the operator
// against the stories completed and how many of those threads asked what a
// digest entry already answered, with the digest's entries counted by status
// and the curator's runs (S-0434, ADR-0200, ADR-0204).
type Digest struct {
	Weeks   []DigestWeek  `json:"weeks"`
	Entries DigestEntries `json:"entries"`
	Curator DigestCurator `json:"curator"`
}

// DigestWeek is one ISO week of the digest section. PerStory is
// ThreadsToOperator over StoriesCompleted, not rounded, and nil in a week
// with no completed story, which has no rate. AskedAgain is nil in a week no
// curator run has reached, where no repeat could have been counted, and set,
// 0 included, in every other.
type DigestWeek struct {
	Week              string   `json:"week"`
	Start             string   `json:"start"`
	ThreadsToOperator int      `json:"threads_to_operator"`
	StoriesCompleted  int      `json:"stories_completed"`
	PerStory          *float64 `json:"per_story,omitempty"`
	AskedAgain        *int     `json:"asked_again,omitempty"`
	CuratorRuns       int      `json:"curator_runs"`
}

// DigestCurator is the curator's runs over its whole log, not the window,
// and when the newest ended, in workitem.TimeFormat, empty before the first
// (ADR-0204).
type DigestCurator struct {
	Runs    int    `json:"runs"`
	LastRun string `json:"last_run,omitempty"`
}

// DigestEntries is the digest's entries now with each status, over the whole
// digest.
type DigestEntries struct {
	Active    int `json:"active"`
	Compacted int `json:"compacted"`
	Absorbed  int `json:"absorbed"`
}

// flaiAuthors are the authors flai writes its own notices as, such as the
// upgrade thread, whose threads are not questions to the operator.
var flaiAuthors = map[string]bool{"flai": true, "flai serve": true}

// digestSection counts, in the ISO weeks from the one that holds start to
// the one that holds now, the threads opened to the operator, each in the
// week that holds its first entry and only from start to now; the stories
// among all completed in each week, cancelled ones left out; and the repeats
// among those threads, from the entries' asked_again and date, in the weeks a
// curator run has reached alone: those whose Monday is at or before the
// newest run's end. It counts the entries by status, and the curator's runs,
// the entries of its activity document among docs, each in the week that
// holds its end. owner is the manifest's owner, whose threads, like flai's,
// are not opened to the operator.
func digestSection(all []*workitem.Item, ths []*threads.Thread, entries []digest.Entry, docs []*workitem.Activity, owner string, start, now time.Time) *Digest {
	out := &Digest{Weeks: []DigestWeek{}}
	at := map[string]int{}
	var mondays []time.Time
	for i, m := 0, monday(start); !m.After(now); i, m = i+1, m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		at[weekKey(y, w)] = i
		mondays = append(mondays, m)
		out.Weeks = append(out.Weeks, DigestWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02")})
	}
	week := func(t time.Time) int {
		t = t.UTC()
		if t.IsZero() || t.Before(start) || t.After(now) {
			return -1
		}
		y, w := t.ISOWeek()
		if i, ok := at[weekKey(y, w)]; ok {
			return i
		}
		return -1
	}
	for _, it := range all {
		if it.Type != workitem.Story || it.Status != workitem.Done {
			continue
		}
		if i := week(it.CompletedAt()); i >= 0 {
			out.Weeks[i].StoriesCompleted++
		}
	}
	// answered is, by each thread's canonical ID, the earliest day of an
	// entry that lists it in asked_again.
	answered := map[string]time.Time{}
	for _, e := range entries {
		switch e.Status {
		case "active":
			out.Entries.Active++
		case "compacted":
			out.Entries.Compacted++
		case "absorbed":
			out.Entries.Absorbed++
		}
		day, err := time.Parse(time.DateOnly, strings.TrimSpace(e.Date))
		if err != nil {
			continue
		}
		for _, id := range e.AskedAgain {
			id = threads.CanonicalID(id)
			if d, ok := answered[id]; !ok || day.Before(d) {
				answered[id] = day
			}
		}
	}
	// newest is when the curator's newest run ended: the later of its
	// last_run and its newest entry's end; zero when it has never run.
	var newest time.Time
	for _, d := range docs {
		if d == nil || d.Kind != workitem.ActivityCurator {
			continue
		}
		out.Curator.Runs += len(d.Entries)
		if t, err := time.Parse(workitem.TimeFormat, d.LastRun); err == nil && t.After(newest) {
			newest = t.UTC()
		}
		for _, e := range d.Entries {
			if e.At.After(newest) {
				newest = e.At.UTC()
			}
			if i := week(e.At); i >= 0 {
				out.Weeks[i].CuratorRuns++
			}
		}
	}
	if !newest.IsZero() {
		out.Curator.LastRun = newest.Format(workitem.TimeFormat)
	}
	repeats := make([]int, len(out.Weeks))
	owner = strings.TrimSpace(owner)
	for _, th := range ths {
		es := th.Entries()
		if len(es) == 0 {
			continue
		}
		author := strings.TrimSpace(es[0].Author)
		if flaiAuthors[author] || (owner != "" && author == owner) {
			continue
		}
		opened, err := time.Parse(workitem.TimeFormat, es[0].At)
		if err != nil {
			continue
		}
		i := week(opened)
		if i < 0 {
			continue
		}
		wk := &out.Weeks[i]
		wk.ThreadsToOperator++
		opened = opened.UTC()
		if d, ok := answered[threads.CanonicalID(th.ID)]; ok && d.Before(time.Date(opened.Year(), opened.Month(), opened.Day(), 0, 0, 0, 0, time.UTC)) {
			repeats[i]++
		}
	}
	for i := range out.Weeks {
		wk := &out.Weeks[i]
		if !newest.IsZero() && !mondays[i].After(newest) {
			n := repeats[i]
			wk.AskedAgain = &n
		}
		if wk.StoriesCompleted > 0 {
			r := float64(wk.ThreadsToOperator) / float64(wk.StoriesCompleted)
			wk.PerStory = &r
		}
	}
	return out
}
