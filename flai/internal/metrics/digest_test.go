package metrics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// digestOwner is the manifest's owner in the digest fixtures.
const digestOwner = "alex"

// openedThread is the thread with the ID and status given whose entries are
// the pairs of time and author given; with none, it was created at
// 2026-09-02T10:00:00Z and has no entry.
func openedThread(id, status string, entries ...string) *threads.Thread {
	th := &threads.Thread{ID: id, Status: status, Created: "2026-09-02T10:00:00Z", Updated: "2026-09-02T10:00:00Z"}
	th.Body = "\n# " + id + " A question\n\n## Entries\n"
	for i := 0; i < len(entries); i += 2 {
		th.Body += "\n### " + entries[i] + " " + entries[i+1] + "\nSome text.\n"
	}
	return th
}

// digestItems are stories, a cancelled one, and a task, closed over the
// window from coordStart, a Wednesday of week 35 at noon, to coordNow, a
// Wednesday of week 38 at noon, the week of each noted beside it.
func digestItems() []*workitem.Item {
	return []*workitem.Item{
		// week 35, before the window's start: left out
		claimItem("S-0004", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-08-25T09:00:00Z"),
		// week 36: two stories done, one cancelled, and a task done
		claimItem("S-0001", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-09-01T09:00:00Z"),
		claimItem("S-0002", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-09-02T09:00:00Z"),
		claimItem("S-0003", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Cancelled, "2026-09-02T09:00:00Z"),
		claimItem("T-0001", workitem.Task, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-09-02T09:00:00Z"),
		// week 37: one story done
		claimItem("S-0005", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-09-09T09:00:00Z"),
		// still in progress
		claimItem("S-0006", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.InProgress, "2026-09-09T09:00:00Z"),
	}
}

// digestThreads are the threads of the fixtures, the week of each noted
// above it.
func digestThreads() []*threads.Thread {
	return []*threads.Thread{
		// week 36: opened by an agent, listed by two entries dated before
		// its day, one repeat
		openedThread("TH-0001", "resolved", "2026-09-03T10:00:00Z", "agent", "2026-09-03T11:00:00Z", digestOwner),
		// week 36: opened by the planner, listed by an entry of its own day
		// only, no repeat
		openedThread("TH-0002", "open", "2026-09-03T08:00:00Z", "planner"),
		// week 36: opened by the owner, by flai serve, and by flai, listed
		// by an entry before their day: none counts
		openedThread("TH-0003", "open", "2026-09-03T08:00:00Z", digestOwner),
		openedThread("TH-0004", "resolved", "2026-09-03T08:00:00Z", "flai serve"),
		openedThread("TH-0005", "resolved", "2026-09-03T08:00:00Z", "flai"),
		// week 36: with no entry, no opener: not counted
		openedThread("TH-0006", "open"),
		// week 36: opened by an agent, listed as th-7 by an absorbed entry
		// before its day, a repeat
		openedThread("TH-0007", "resolved", "2026-09-04T00:30:00Z", "S-0002"),
		// week 38: opened by the orchestrator, in a week with no completed
		// story
		openedThread("TH-0008", "open", "2026-09-15T10:00:00Z", "orchestrator"),
		// week 35, before the window's start: left out
		openedThread("TH-0009", "resolved", "2026-08-25T10:00:00Z", "agent"),
		// after now: left out
		openedThread("TH-0010", "open", "2026-09-17T10:00:00Z", "agent"),
	}
}

// digestEntries are the digest's entries of the fixtures: two active, one
// compacted, and one absorbed.
func digestEntries() []digest.Entry {
	return []digest.Entry{
		{ID: "D-0001", Status: "active", Date: "2026-09-01", AskedAgain: []string{"TH-0001", "TH-0003", "TH-0004", "TH-0005", "TH-0009"}},
		{ID: "D-0002", Status: "compacted", Date: "2026-09-02", AskedAgain: []string{"TH-0001", "TH-0006"}},
		{ID: "D-0003", Status: "active", Date: "2026-09-03", AskedAgain: []string{"TH-0002"}},
		{ID: "D-0004", Status: "absorbed", Date: "2026-08-30", AskedAgain: []string{"th-7", "TH-0010"}},
	}
}

// ratio is a per_story value.
func ratio(v float64) *float64 { return &v }

// S-0434: per ISO week of the window, the threads opened to the operator,
// the stories completed, their ratio, and the repeats, with the entries
// counted by status.
func TestDigestCountsEachValuePerWeek(t *testing.T) {
	got := digestSection(digestItems(), digestThreads(), digestEntries(), digestOwner, coordStart, coordNow)
	want := &Digest{
		Weeks: []DigestWeek{
			{Week: "2026-W35", Start: "2026-08-24"},
			{Week: "2026-W36", Start: "2026-08-31", ThreadsToOperator: 3, StoriesCompleted: 2, PerStory: ratio(1.5), AskedAgain: 2},
			{Week: "2026-W37", Start: "2026-09-07", StoriesCompleted: 1, PerStory: ratio(0)},
			{Week: "2026-W38", Start: "2026-09-14", ThreadsToOperator: 1},
		},
		Entries: DigestEntries{Active: 2, Compacted: 1, Absorbed: 1},
	}
	if !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.Marshal(got)
		wantJSON, _ := json.Marshal(want)
		t.Errorf("digest =\n%s\nwant\n%s", gotJSON, wantJSON)
	}
}

// S-0434: the rate is not rounded, and a week with threads and no completed
// story has none: per_story is absent from its JSON.
func TestDigestLeavesTheRateOutOfAWeekWithNoCompletedStory(t *testing.T) {
	items := []*workitem.Item{
		claimItem("S-0001", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-09-01T09:00:00Z"),
		claimItem("S-0002", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-09-01T10:00:00Z"),
		claimItem("S-0003", workitem.Story, "2026-08-20T09:00:00Z", nil, workitem.Done, "2026-09-01T11:00:00Z"),
	}
	ths := []*threads.Thread{
		openedThread("TH-0001", "open", "2026-09-01T08:00:00Z", "agent"),
		openedThread("TH-0002", "open", "2026-09-08T08:00:00Z", "agent"),
	}
	got := digestSection(items, ths, nil, digestOwner, coordStart, coordNow)
	if pr := got.Weeks[1].PerStory; pr == nil || *pr != 1.0/3 {
		t.Errorf("week 36 per_story = %v, want 1/3 unrounded", pr)
	}
	wk := got.Weeks[2]
	if wk.ThreadsToOperator != 1 || wk.StoriesCompleted != 0 || wk.PerStory != nil {
		t.Errorf("week 37 = %+v, want one thread, no story, and no rate", wk)
	}
	b, err := json.Marshal(wk)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "per_story") {
		t.Errorf("a week with no completed story carries per_story: %s", b)
	}
}

// S-0434: a project with no digest still counts its weeks, with no repeat
// and every entry count 0, so the weeks before the digest are measured too.
func TestDigestWithoutEntriesCountsTheWeeks(t *testing.T) {
	got := digestSection(digestItems(), digestThreads(), nil, digestOwner, coordStart, coordNow)
	if got.Entries != (DigestEntries{}) {
		t.Errorf("entries = %+v, want all 0", got.Entries)
	}
	want := []DigestWeek{
		{Week: "2026-W35", Start: "2026-08-24"},
		{Week: "2026-W36", Start: "2026-08-31", ThreadsToOperator: 3, StoriesCompleted: 2, PerStory: ratio(1.5)},
		{Week: "2026-W37", Start: "2026-09-07", StoriesCompleted: 1, PerStory: ratio(0)},
		{Week: "2026-W38", Start: "2026-09-14", ThreadsToOperator: 1},
	}
	if !reflect.DeepEqual(got.Weeks, want) {
		gotJSON, _ := json.Marshal(got.Weeks)
		wantJSON, _ := json.Marshal(want)
		t.Errorf("weeks =\n%s\nwant\n%s", gotJSON, wantJSON)
	}
}

// S-0434: Compute reports digest in every report, counting stories whatever
// the report's type, and leaves it out when the digest cannot be read.
func TestComputeReportsTheDigestUnlessItCannotBeRead(t *testing.T) {
	opt := Options{Now: coordNow, Since: coordNow.Sub(coordStart), Type: workitem.Task, Threads: digestThreads(), Entries: digestEntries(), Owner: digestOwner}
	rep := Compute(digestItems(), opt)
	if rep.Digest == nil || len(rep.Digest.Weeks) != 4 || rep.Digest.Weeks[1].StoriesCompleted != 2 || rep.Digest.Weeks[1].AskedAgain != 2 || rep.Digest.Entries.Active != 2 {
		t.Errorf("digest of a report of tasks = %+v", rep.Digest)
	}
	b, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"digest":{"weeks":[`) || !strings.Contains(string(b), `"entries":{"active":2,"compacted":1,"absorbed":1}`) {
		t.Errorf("report JSON lacks the digest section: %s", b)
	}
	opt.Entries, opt.DigestUnreadable = nil, true
	if rep := Compute(digestItems(), opt); rep.Digest != nil {
		t.Errorf("an unreadable digest: digest = %+v, want absent", rep.Digest)
	}
}
