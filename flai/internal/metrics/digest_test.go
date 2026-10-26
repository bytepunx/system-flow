package metrics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

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

// again is an asked_again value.
func again(v int) *int { return &v }

// activityDoc is the activity document of the kind given whose last_run is
// lastRun and whose entries ended at the times given, in workitem.TimeFormat.
func activityDoc(t *testing.T, kind, lastRun string, ends ...string) *workitem.Activity {
	t.Helper()
	d := &workitem.Activity{Kind: kind, LastRun: lastRun, Entries: []workitem.ActivityEntry{}}
	for _, s := range ends {
		at, err := time.Parse(workitem.TimeFormat, s)
		if err != nil {
			t.Fatal(err)
		}
		d.Entries = append(d.Entries, workitem.ActivityEntry{At: at, Summary: "Curated the digest"})
	}
	return d
}

// curatorRanLast is a curator that ran once, in week 38, after the
// fixtures' every thread and entry, so that it has reached every week.
func curatorRanLast(t *testing.T) []*workitem.Activity {
	return []*workitem.Activity{activityDoc(t, workitem.ActivityCurator, "2026-09-15T05:00:00Z", "2026-09-15T05:00:00Z")}
}

// S-0434: per ISO week of the window, the threads opened to the operator,
// the stories completed, their ratio, the repeats, and the curator's runs,
// with the entries counted by status and the curator's runs in all.
func TestDigestCountsEachValuePerWeek(t *testing.T) {
	got := digestSection(digestItems(), digestThreads(), digestEntries(), curatorRanLast(t), digestOwner, coordStart, coordNow)
	want := &Digest{
		Weeks: []DigestWeek{
			{Week: "2026-W35", Start: "2026-08-24", AskedAgain: again(0)},
			{Week: "2026-W36", Start: "2026-08-31", ThreadsToOperator: 3, StoriesCompleted: 2, PerStory: ratio(1.5), AskedAgain: again(2)},
			{Week: "2026-W37", Start: "2026-09-07", StoriesCompleted: 1, PerStory: ratio(0), AskedAgain: again(0)},
			{Week: "2026-W38", Start: "2026-09-14", ThreadsToOperator: 1, AskedAgain: again(0), CuratorRuns: 1},
		},
		Entries: DigestEntries{Active: 2, Compacted: 1, Absorbed: 1},
		Curator: DigestCurator{Runs: 1, LastRun: "2026-09-15T05:00:00Z"},
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
	got := digestSection(items, ths, nil, nil, digestOwner, coordStart, coordNow)
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

// S-0434: a project with no digest still counts its weeks, with every
// entry count 0 and, in each week a curator run has reached, no repeat, so
// the weeks before the digest are measured too.
func TestDigestWithoutEntriesCountsTheWeeks(t *testing.T) {
	got := digestSection(digestItems(), digestThreads(), nil, curatorRanLast(t), digestOwner, coordStart, coordNow)
	if got.Entries != (DigestEntries{}) {
		t.Errorf("entries = %+v, want all 0", got.Entries)
	}
	want := []DigestWeek{
		{Week: "2026-W35", Start: "2026-08-24", AskedAgain: again(0)},
		{Week: "2026-W36", Start: "2026-08-31", ThreadsToOperator: 3, StoriesCompleted: 2, PerStory: ratio(1.5), AskedAgain: again(0)},
		{Week: "2026-W37", Start: "2026-09-07", StoriesCompleted: 1, PerStory: ratio(0), AskedAgain: again(0)},
		{Week: "2026-W38", Start: "2026-09-14", ThreadsToOperator: 1, AskedAgain: again(0), CuratorRuns: 1},
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
	opt := Options{Now: coordNow, Since: coordNow.Sub(coordStart), Type: workitem.Task, Threads: digestThreads(), Entries: digestEntries(), Owner: digestOwner, Activities: curatorRanLast(t)}
	rep := Compute(digestItems(), opt)
	if rep.Digest == nil || len(rep.Digest.Weeks) != 4 || rep.Digest.Weeks[1].StoriesCompleted != 2 || !reflect.DeepEqual(rep.Digest.Weeks[1].AskedAgain, again(2)) || rep.Digest.Entries.Active != 2 || rep.Digest.Curator.Runs != 1 {
		t.Errorf("digest of a report of tasks = %+v", rep.Digest)
	}
	b, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"digest":{"weeks":[`) || !strings.Contains(string(b), `"entries":{"active":2,"compacted":1,"absorbed":1},"curator":{"runs":1,"last_run":"2026-09-15T05:00:00Z"}`) {
		t.Errorf("report JSON lacks the digest section: %s", b)
	}
	opt.Entries, opt.DigestUnreadable = nil, true
	if rep := Compute(digestItems(), opt); rep.Digest != nil {
		t.Errorf("an unreadable digest: digest = %+v, want absent", rep.Digest)
	}
}

// I-0207: a project whose curator has never run reads no repeat count in any
// week, rather than 0 as if it were a measurement, and its curator no run,
// though threads were opened to the operator and the digest has entries. An
// activity document of another kind is no curator run.
func TestDigestLeavesAskedAgainOutWithoutACuratorRunOfI0207(t *testing.T) {
	entries := []digest.Entry{
		{ID: "D-0001", Status: "active", Date: "2026-09-01", AskedAgain: []string{}},
		{ID: "D-0002", Status: "compacted", Date: "2026-09-02"},
	}
	docs := []*workitem.Activity{activityDoc(t, workitem.ActivityPlanner, "2026-09-15T05:00:00Z", "2026-09-15T05:00:00Z")}
	got := digestSection(digestItems(), digestThreads(), entries, docs, digestOwner, coordStart, coordNow)
	if got.Curator != (DigestCurator{}) {
		t.Errorf("curator = %+v, want no run and no last_run", got.Curator)
	}
	threadsSeen := 0
	for _, wk := range got.Weeks {
		threadsSeen += wk.ThreadsToOperator
		if wk.AskedAgain != nil {
			t.Errorf("week %s: asked_again %d, want absent", wk.Week, *wk.AskedAgain)
		}
		if wk.CuratorRuns != 0 {
			t.Errorf("week %s: curator_runs %d, want 0", wk.Week, wk.CuratorRuns)
		}
	}
	if threadsSeen != 4 {
		t.Errorf("threads opened to the operator = %d, want 4", threadsSeen)
	}
	b, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "asked_again") || !strings.Contains(string(b), `"curator":{"runs":0}`) {
		t.Errorf("digest JSON = %s, want no asked_again and a curator with no run", b)
	}
}

// ADR-0204: the curator's runs are counted in the weeks they ended in and
// over its whole log, its newest end is the later of last_run and its newest
// entry's, and a week has asked_again, 0 included, only when its Monday is
// at or before that end.
func TestDigestCountsTheCuratorsRunsAndTheWeeksTheyReached(t *testing.T) {
	docs := []*workitem.Activity{
		activityDoc(t, workitem.ActivityPlanner, "2026-09-15T05:00:00Z", "2026-09-15T05:00:00Z"),
		// before the window, in week 36, and in week 37, with a last_run
		// later than the newest entry, on Thursday of week 37
		activityDoc(t, workitem.ActivityCurator, "2026-09-10T05:00:00Z", "2026-08-20T05:00:00Z", "2026-09-02T05:00:00Z", "2026-09-08T05:00:00Z"),
	}
	got := digestSection(digestItems(), digestThreads(), digestEntries(), docs, digestOwner, coordStart, coordNow)
	if want := (DigestCurator{Runs: 3, LastRun: "2026-09-10T05:00:00Z"}); got.Curator != want {
		t.Errorf("curator = %+v, want %+v", got.Curator, want)
	}
	want := []DigestWeek{
		{Week: "2026-W35", Start: "2026-08-24", AskedAgain: again(0)},
		{Week: "2026-W36", Start: "2026-08-31", ThreadsToOperator: 3, StoriesCompleted: 2, PerStory: ratio(1.5), AskedAgain: again(2), CuratorRuns: 1},
		{Week: "2026-W37", Start: "2026-09-07", StoriesCompleted: 1, PerStory: ratio(0), AskedAgain: again(0), CuratorRuns: 1},
		{Week: "2026-W38", Start: "2026-09-14", ThreadsToOperator: 1},
	}
	if !reflect.DeepEqual(got.Weeks, want) {
		gotJSON, _ := json.Marshal(got.Weeks)
		wantJSON, _ := json.Marshal(want)
		t.Errorf("weeks =\n%s\nwant\n%s", gotJSON, wantJSON)
	}
	b, err := json.Marshal(got.Weeks)
	if err != nil {
		t.Fatal(err)
	}
	if s := string(b); !strings.Contains(s, `"start":"2026-09-07","threads_to_operator":0,"stories_completed":1,"per_story":0,"asked_again":0,"curator_runs":1}`) ||
		!strings.Contains(s, `"start":"2026-09-14","threads_to_operator":1,"stories_completed":0,"curator_runs":0}`) {
		t.Errorf("weeks JSON = %s, want asked_again 0 in week 37 and none in week 38", s)
	}
	// a run that ended at a Monday's midnight reaches that week
	docs = []*workitem.Activity{activityDoc(t, workitem.ActivityCurator, "", "2026-09-14T00:00:00Z")}
	got = digestSection(digestItems(), digestThreads(), digestEntries(), docs, digestOwner, coordStart, coordNow)
	if wk := got.Weeks[3]; !reflect.DeepEqual(wk.AskedAgain, again(0)) || wk.CuratorRuns != 1 || got.Curator.LastRun != "2026-09-14T00:00:00Z" {
		t.Errorf("a run at week 38's Monday: week 38 = %+v, curator = %+v", wk, got.Curator)
	}
}
