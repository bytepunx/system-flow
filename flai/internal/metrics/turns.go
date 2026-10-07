package metrics

import (
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Turns is the turns of the stories' own agents in the window, per class
// (S-0293): summed over the window, per day of it, and per story. They come
// from each story's usage.turns, which flai serve measures from its agent's
// logs; every story counts, archived or not and whatever its status or the
// report's type.
type Turns struct {
	// Classes are the classes of a turn, in the order they are written.
	Classes []string `json:"classes"`
	// Total is the sum over Days.
	Total TurnCounts `json:"total"`
	// Days are every UTC day of the window, from its start's to now's, oldest
	// first, each with the turns of every story on it.
	Days []TurnsDay `json:"days"`
	// Stories are those with a turn on a day of the window, by canonical ID,
	// each with its turns on those days.
	Stories []StoryTurns `json:"stories"`
}

// TurnCounts counts turns by class, and in all as Turns.
type TurnCounts struct {
	Turns      int `json:"turns"`
	Ceremony   int `json:"ceremony"`
	TestRuns   int `json:"test_runs"`
	EmptyWakes int `json:"empty_wakes"`
	HandEdits  int `json:"hand_edits"`
	Work       int `json:"work"`
}

// TurnsDay is the stories' turns on one UTC day, as 2006-01-02.
type TurnsDay struct {
	Day string `json:"day"`
	TurnCounts
}

// StoryTurns is one story's turns over the days of the window.
type StoryTurns struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	TurnCounts
}

// Count is the turns of class, one of usage.TurnClasses; 0 for another.
func (c TurnCounts) Count(class string) int { return c.classes().Count(class) }

// classes is c's counts by class, as a day of usage holds them.
func (c TurnCounts) classes() usage.TurnDay {
	return usage.TurnDay{Ceremony: c.Ceremony, TestRuns: c.TestRuns, EmptyWakes: c.EmptyWakes, HandEdits: c.HandEdits, Work: c.Work}
}

// add adds one day's turns to c.
func (c *TurnCounts) add(d usage.TurnDay) {
	c.Turns += d.Total()
	c.Ceremony += d.Ceremony
	c.TestRuns += d.TestRuns
	c.EmptyWakes += d.EmptyWakes
	c.HandEdits += d.HandEdits
	c.Work += d.Work
}

// turns sums the turns the stories of all carry on the UTC days from
// start's to now's, per day, per story, and over the window.
func turns(all []*workitem.Item, start, now time.Time) Turns {
	out := Turns{Classes: slices.Clone(usage.TurnClasses), Days: []TurnsDay{}, Stories: []StoryTurns{}}
	at := map[string]int{}
	for d := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC); !d.After(now); d = d.AddDate(0, 0, 1) {
		day := d.Format(usage.DayFormat)
		at[day] = len(out.Days)
		out.Days = append(out.Days, TurnsDay{Day: day})
	}
	for _, it := range all {
		if it.Type != workitem.Story || it.Usage == nil {
			continue
		}
		s := StoryTurns{ID: it.ID, Title: it.Title, Status: it.Status}
		for _, d := range it.Usage.Turns {
			i, ok := at[d.Day]
			if !ok {
				continue
			}
			out.Days[i].add(d)
			s.add(d)
		}
		if s.Turns > 0 {
			out.Stories = append(out.Stories, s)
		}
	}
	for _, d := range out.Days {
		out.Total.add(d.classes())
	}
	slices.SortFunc(out.Stories, func(a, b StoryTurns) int {
		return strings.Compare(workitem.CanonicalID(a.ID), workitem.CanonicalID(b.ID))
	})
	return out
}
