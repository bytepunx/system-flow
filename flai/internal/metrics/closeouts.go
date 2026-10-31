package metrics

import (
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// The serve journal's close-out entries, as flai serve writes them for a
// close-out job (ADR-0229 §6).
const (
	closeOutAction = "close-out"
	closeOutMethod = "verify.job"
	// closeOutFromRecord begins the detail of an answer from the record.
	closeOutFromRecord = "from the record"
)

// JournalEntry is a line of flai serve's host journal, journal.jsonl, as far
// as the close-outs read it (ADR-0229 §6).
type JournalEntry struct {
	At      time.Time `json:"at"`
	Action  string    `json:"action"`
	Method  string    `json:"method"`
	Root    string    `json:"root"`
	Outcome string    `json:"outcome"`
	Detail  string    `json:"detail,omitempty"`
	// DurationMS is how long a close-out job ran, nil when the entry has none.
	DurationMS *int64 `json:"duration_ms,omitempty"`
	Verdict    string `json:"verdict,omitempty"`
}

// CloseOutJournal is the serve journal the close-outs are counted from: its
// entries, every project's and action's, and Root, the main checkout of the
// project whose close-outs count.
type CloseOutJournal struct {
	Root    string
	Entries []JournalEntry
}

// CloseOuts is, per ISO week of the window, the close-outs flai serve ran as
// jobs by verdict, those answered from the record apart, and the p50 and p85
// of the durations of those that ran (S-0546, ADR-0229).
type CloseOuts struct {
	Weeks []CloseOutWeek `json:"weeks"`
}

// CloseOutWeek is the close-outs that ended in one ISO week. Passed, Stopped,
// and Failed leave out those answered from the record, which Recorded counts
// whatever their verdict; Failed counts interrupted jobs too. P50 and P85 are
// in seconds, nil when no close-out of the week ran.
type CloseOutWeek struct {
	Week     string   `json:"week"`
	Start    string   `json:"start"`
	Passed   int      `json:"passed"`
	Stopped  int      `json:"stopped"`
	Failed   int      `json:"failed"`
	Recorded int      `json:"recorded"`
	P50      *float64 `json:"p50_seconds,omitempty"`
	P85      *float64 `json:"p85_seconds,omitempty"`
}

// closeOuts counts, in the ISO weeks from the one that holds start to the one
// that holds now, the close-out end entries of j's project, each in the week
// that holds its time and only from start to now.
func closeOuts(j CloseOutJournal, start, now time.Time) *CloseOuts {
	out := &CloseOuts{Weeks: []CloseOutWeek{}}
	at := map[string]int{}
	for i, m := 0, monday(start); !m.After(now); i, m = i+1, m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		at[weekKey(y, w)] = i
		out.Weeks = append(out.Weeks, CloseOutWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02")})
	}
	root := filepath.Clean(j.Root)
	ran := map[int][]int64{}
	for _, e := range j.Entries {
		if e.Action != closeOutAction || e.Method != closeOutMethod || filepath.Clean(e.Root) != root {
			continue
		}
		if e.Outcome != "done" && e.Outcome != "failed" {
			continue
		}
		t := e.At.UTC()
		if t.IsZero() || t.Before(start) || t.After(now) {
			continue
		}
		y, w := t.ISOWeek()
		i, ok := at[weekKey(y, w)]
		if !ok {
			continue
		}
		wk := &out.Weeks[i]
		if strings.HasPrefix(e.Detail, closeOutFromRecord) {
			wk.Recorded++
			continue
		}
		switch e.Verdict {
		case "passed":
			wk.Passed++
		case "stopped":
			wk.Stopped++
		case "failed", "interrupted":
			wk.Failed++
		}
		if e.Verdict != "interrupted" && e.DurationMS != nil {
			ran[i] = append(ran[i], *e.DurationMS)
		}
	}
	for i, ms := range ran {
		slices.Sort(ms)
		out.Weeks[i].P50, out.Weeks[i].P85 = msPercentile(ms, 50), msPercentile(ms, 85)
	}
	return out
}

// msPercentile is the nearest-rank p-th percentile of sorted milliseconds,
// in seconds, not rounded.
func msPercentile(sorted []int64, p int) *float64 {
	rank := max((p*len(sorted)+99)/100, 1)
	s := float64(sorted[rank-1]) / 1000
	return &s
}
