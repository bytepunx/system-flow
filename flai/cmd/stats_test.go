package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0404: flai stats prints the merge queue's section after the waiting,
// from the queue's event log, and --json and the dashboard's read carry it
// under queue, with each story's queue wait under items; with no log, there
// is neither.
func TestStatsReportsTheMergeQueue(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("LOG_FORMAT", "json")
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, title := range []string{"Lands", "Sent back"} {
		if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Owner: "t", Now: time.Date(2026, 9, 28, 9, 0, 0, 0, time.UTC)}); err != nil {
			t.Fatal(err)
		}
	}
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)

	// no log: no section and no queue waits
	out, errOut, code := runInAt(t, root, at, "stats", "--json")
	if code != 0 || strings.Contains(out, `"queue"`) || strings.Contains(out, `"queue_wait_seconds"`) {
		t.Errorf("stats --json without a log: %d %s\n%s", code, errOut, out)
	}
	if out, _, _ = runInAt(t, root, at, "stats"); strings.Contains(out, "merge queue") {
		t.Errorf("stats without a log prints the queue:\n%s", out)
	}

	enq, built, verified, bisected := mergequeue.EventEnqueued, mergequeue.EventBatchBuilt, mergequeue.EventVerified, mergequeue.EventBisected
	ev := func(at, kind, batch string, stories ...string) mergequeue.Event {
		t.Helper()
		when, err := time.Parse(workitem.TimeFormat, at)
		if err != nil {
			t.Fatal(err)
		}
		return mergequeue.Event{At: when, Event: kind, Batch: batch, Stories: stories}
	}
	// week 40: a batch of two fails and is bisected to S-0002, sent back
	// after an hour; S-0001 lands alone after three
	events := []mergequeue.Event{
		ev("2026-09-29T09:00:00Z", enq, "", "S-0001"),
		ev("2026-09-29T09:00:00Z", enq, "", "S-0002"),
		ev("2026-09-29T09:10:00Z", built, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T09:40:00Z", verified, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T09:41:00Z", bisected, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T10:00:00Z", mergequeue.EventSentBack, "B1", "S-0002"),
		ev("2026-09-29T11:00:00Z", built, "B2", "S-0001"),
		ev("2026-09-29T12:00:00Z", mergequeue.EventLanded, "B2", "S-0001"),
	}
	if err := mergequeue.AppendEvents(repo.CacheDir(), events); err != nil {
		t.Fatal(err)
	}

	out, errOut, code = runInAt(t, root, at, "stats", "--json")
	if code != 0 {
		t.Fatalf("stats --json: %s", errOut)
	}
	var rep struct {
		Queue *metrics.Queue        `json:"queue"`
		Items []metrics.ItemMetrics `json:"items"`
	}
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatal(err)
	}
	// the 30 days to the 3rd of October span ISO weeks 36 to 40
	if rep.Queue == nil || len(rep.Queue.Weeks) != 5 || rep.Queue.Weeks[0].Week != "2026-W36" {
		t.Fatalf("queue: %+v", rep.Queue)
	}
	want := metrics.QueueWeek{Week: "2026-W40", Start: "2026-09-28", Landed: 1, SentBack: 1, WaitMean: 2 * 3600, WaitMax: 3 * 3600,
		Batches: 2, BatchSizeMean: 1.5, BatchSizeMax: 2, Bisections: 1}
	if got := rep.Queue.Weeks[4]; got != want {
		t.Errorf("week 40 = %+v, want %+v", got, want)
	}
	waits := map[string]float64{}
	for _, m := range rep.Items {
		if m.QueueWait != nil {
			waits[m.ID] = *m.QueueWait
		}
	}
	if len(waits) != 2 || waits["S-0001"] != 3*3600 || waits["S-0002"] != 3600 {
		t.Errorf("items' queue waits = %v", waits)
	}

	out, _, code = runInAt(t, root, at, "stats")
	wantText := "\nmerge queue:\n" +
		"  stories        landed 1 · sent back 1 · longest wait 3h\n" +
		"  batches        built 2 · mean size 1.5 · largest 2 · bisections 1\n" +
		"  by week:\n" +
		"    2026-W40 (2026-09-28)  stories landed 1 · sent back 1 · mean wait 2h · longest 3h; batches built 2 · mean size 1.5 · largest 2 · bisections 1\n"
	if code != 0 || !strings.Contains(out, wantText) {
		t.Errorf("text: %d, want it to hold %q, got\n%s", code, wantText, out)
	}
	// the dashboard's read answers the same
	if _, rerr := sameAnswer(t, root, hostapi.Host{}, "stats.get", `{}`, "stats"); rerr != nil {
		t.Errorf("stats.get: %+v", rerr)
	}

	// a log that cannot be read is warned of, and the section is absent
	if err := os.WriteFile(mergequeue.EventsPath(repo.CacheDir()), []byte("{not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runInAt(t, root, at, "stats", "--json")
	if code != 0 || strings.Contains(out, `"queue"`) || !strings.Contains(errOut, "merge queue's event log") {
		t.Errorf("an unreadable log: %d\n%s\n%s", code, out, errOut)
	}
}

// S-0404: the queue's section comes after the waiting, and before the
// turns, and a queue that counts nothing in the window prints nothing.
func TestPrintQueueStatsComesAfterTheWaitingAndSkipsAnEmptyQueue(t *testing.T) {
	var out bytes.Buffer
	printQueueStats(&app{out: &out}, &metrics.Queue{Weeks: []metrics.QueueWeek{{Week: "2026-W40", Start: "2026-09-28"}}})
	if out.Len() != 0 {
		t.Errorf("an empty queue prints:\n%s", out.String())
	}
	rep := &metrics.Report{
		WindowStart: "2026-09-03T20:00:00Z", Type: workitem.Story, WindowDays: 30,
		Waiting: metrics.Waiting{Weeks: []metrics.WaitWeek{{Week: "2026-W40", Start: "2026-09-28", Items: 1}}},
		Queue:   &metrics.Queue{Weeks: []metrics.QueueWeek{{Week: "2026-W40", Start: "2026-09-28", Landed: 1}}},
	}
	printSummary(&app{out: &out}, rep, false)
	text := out.String()
	waiting, queue := strings.Index(text, "\nwaiting, over 1 completed:\n"), strings.Index(text, "\nmerge queue:\n")
	if waiting < 0 || queue < waiting {
		t.Errorf("the queue's section is not after the waiting:\n%s", text)
	}
}

// digestEntry is a digest entry with the ID, status, and date given that
// lists the threads given in asked_again.
func digestEntry(id, status, date string, askedAgain ...string) string {
	return "---\nid: " + id + "\ntitle: Left it is.\nkind: answer\nsources: []\ntopics: []\npaths: []\n" +
		"date: " + date + "\nstatus: " + status + "\nabsorbed_by:\nasked_again: [" + strings.Join(askedAgain, ", ") + "]\ncited: []\n---\n\n# " + id + " Left it is.\n\nLeft.\n"
}

// S-0434: flai stats prints the digest's section last, a line for each week
// of the window, one for the entries, and one for the curator, and --json
// carries it under digest, per_story absent in a week with no completed
// story and asked_again in a week no curator run has reached (ADR-0204,
// I-0207); with no digest folder the section is still there, with no
// repeats and no entries; and a digest that cannot be read is warned of and
// leaves it out.
func TestStatsReportsTheDigest(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("LOG_FORMAT", "json")
	root := tempProject(t)
	manifestPath := filepath.Join(root, "system-flow.yaml")
	m, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, append(m, "owner: alex\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// two stories done in week 39, none in week 40
	var stories []string
	for _, done := range []string{"2026-09-23T09:00:00Z", "2026-09-24T09:00:00Z"} {
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Done " + done[:10], Owner: "t", Now: time.Date(2026, 9, 14, 9, 0, 0, 0, time.UTC)})
		if err != nil {
			t.Fatal(err)
		}
		it.Transitions = []workitem.Transition{{To: workitem.InProgress, At: "2026-09-21T09:00:00Z", By: "t"}, {To: workitem.Done, At: done, By: "t"}}
		it.Status = workitem.Done
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
		stories = append(stories, it.ID)
	}
	open := func(author string, at time.Time) string {
		t.Helper()
		th, err := threads.New(repo, threads.NewOptions{Title: "Which way", On: stories[0], Author: author, Text: "Left?", Now: at})
		if err != nil {
			t.Fatal(err)
		}
		return th.ID
	}
	// week 39: three threads opened to the operator, and one the owner
	// opened, which is not
	agent := open("agent", time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC))
	planner := open("planner", time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC))
	open("alex", time.Date(2026, 9, 23, 11, 0, 0, 0, time.UTC))
	story := open(stories[1], time.Date(2026, 9, 24, 10, 0, 0, 0, time.UTC))
	// week 40: one, in a week with no completed story
	open("orchestrator", time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC))
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)
	// the 14 days to the 3rd of October span ISO weeks 38 to 40

	// no digest folder and no curator run (I-0207): the section, with no
	// entries and no repeat count in any week, not 0
	out, errOut, code := runInAt(t, root, at, "stats", "--since", "14d")
	wantText := "\ndigest:\n" +
		"  2026-W38  2026-09-14  threads 0 · stories 0 · – per story · asked again – · curator runs 0\n" +
		"  2026-W39  2026-09-21  threads 3 · stories 2 · 1.50 per story · asked again – · curator runs 0\n" +
		"  2026-W40  2026-09-28  threads 1 · stories 0 · – per story · asked again – · curator runs 0\n" +
		"  entries  active 0 · compacted 0 · absorbed 0\n" +
		"  curator  runs 0 · never run\n"
	if code != 0 || !strings.HasSuffix(out, wantText) {
		t.Errorf("no digest folder and no curator run: %d %s, want it to end in %q, got\n%s", code, errOut, wantText, out)
	}

	// a curator run in week 39 reaches weeks 38 and 39, not 40: with no
	// digest folder, those read no repeats
	if _, err := repo.AppendActivity(workitem.ActivityCurator, workitem.ActivityEntry{At: time.Date(2026, 9, 25, 5, 0, 0, 0, time.UTC), Summary: "Curated the digest"}); err != nil {
		t.Fatal(err)
	}
	week38 := "  2026-W38  2026-09-14  threads 0 · stories 0 · – per story · asked again 0 · curator runs 0\n"
	week40 := "  2026-W40  2026-09-28  threads 1 · stories 0 · – per story · asked again – · curator runs 0\n"
	curator := "  curator  runs 1 · last 2026-09-25T05:00:00Z\n"
	out, errOut, code = runInAt(t, root, at, "stats", "--since", "14d")
	wantText = "\ndigest:\n" + week38 +
		"  2026-W39  2026-09-21  threads 3 · stories 2 · 1.50 per story · asked again 0 · curator runs 1\n" + week40 +
		"  entries  active 0 · compacted 0 · absorbed 0\n" + curator
	if code != 0 || !strings.HasSuffix(out, wantText) {
		t.Errorf("no digest folder: %d %s, want it to end in %q, got\n%s", code, errOut, wantText, out)
	}

	// an entry dated before the agent's thread lists it, a repeat; one of
	// the planner's thread's own day lists it, not a repeat; an absorbed
	// one before the story's thread lists it, a repeat
	dir := filepath.Join(root, "design", "digest")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, doc := range map[string]string{
		"README.md":              "# Digest\n",
		"D-0001-left-it-is.md":   digestEntry("D-0001", "active", "2026-09-20", agent),
		"D-0002-left-again.md":   digestEntry("D-0002", "compacted", "2026-09-23", planner),
		"D-0003-left-at-last.md": digestEntry("D-0003", "absorbed", "2026-09-21", story),
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(doc), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out, errOut, code = runInAt(t, root, at, "stats", "--since", "14d")
	wantText = "\ndigest:\n" + week38 +
		"  2026-W39  2026-09-21  threads 3 · stories 2 · 1.50 per story · asked again 2 · curator runs 1\n" + week40 +
		"  entries  active 1 · compacted 1 · absorbed 1\n" + curator
	if code != 0 || !strings.HasSuffix(out, wantText) {
		t.Errorf("with repeats: %d %s, want it to end in %q, got\n%s", code, errOut, wantText, out)
	}

	out, errOut, code = runInAt(t, root, at, "stats", "--since", "14d", "--json")
	if code != 0 {
		t.Fatalf("stats --json: %s", errOut)
	}
	var rep struct {
		Digest *struct {
			Weeks   []map[string]any      `json:"weeks"`
			Entries metrics.DigestEntries `json:"entries"`
			Curator map[string]any        `json:"curator"`
		} `json:"digest"`
	}
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatal(err)
	}
	if rep.Digest == nil || len(rep.Digest.Weeks) != 3 {
		t.Fatalf("digest: %+v", rep.Digest)
	}
	want := []map[string]any{
		{"week": "2026-W38", "start": "2026-09-14", "threads_to_operator": 0.0, "stories_completed": 0.0, "asked_again": 0.0, "curator_runs": 0.0},
		{"week": "2026-W39", "start": "2026-09-21", "threads_to_operator": 3.0, "stories_completed": 2.0, "per_story": 1.5, "asked_again": 2.0, "curator_runs": 1.0},
		{"week": "2026-W40", "start": "2026-09-28", "threads_to_operator": 1.0, "stories_completed": 0.0, "curator_runs": 0.0},
	}
	for i, w := range want {
		if got := rep.Digest.Weeks[i]; !reflect.DeepEqual(got, w) {
			t.Errorf("digest.weeks[%d] = %v, want %v", i, got, w)
		}
	}
	if want := (metrics.DigestEntries{Active: 1, Compacted: 1, Absorbed: 1}); rep.Digest.Entries != want {
		t.Errorf("digest.entries = %+v, want %+v", rep.Digest.Entries, want)
	}
	if want := (map[string]any{"runs": 1.0, "last_run": "2026-09-25T05:00:00Z"}); !reflect.DeepEqual(rep.Digest.Curator, want) {
		t.Errorf("digest.curator = %v, want %v", rep.Digest.Curator, want)
	}
	// the dashboard's read answers the same
	if _, rerr := sameAnswer(t, root, hostapi.Host{}, "stats.get", `{}`, "stats"); rerr != nil {
		t.Errorf("stats.get: %+v", rerr)
	}

	// an entry that cannot be read is warned of, and the section is absent
	if err := os.WriteFile(filepath.Join(dir, "D-0004-bad.md"), []byte("---\nid: [\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runInAt(t, root, at, "stats", "--json")
	if code != 0 || strings.Contains(out, `"digest"`) || !strings.Contains(errOut, "cannot read the digest") {
		t.Errorf("an unreadable digest: %d\n%s\n%s", code, out, errOut)
	}
	if out, _, _ = runInAt(t, root, at, "stats"); strings.Contains(out, "\ndigest:") {
		t.Errorf("an unreadable digest prints the section:\n%s", out)
	}
}

// S-0434: printDigest prints nothing without a digest, gives the rate to two
// decimals, prints an en dash for asked again in a week no curator run has
// reached, and the curator's runs and last run (ADR-0204).
func TestPrintDigestRoundsTheRateAndSkipsNoDigest(t *testing.T) {
	var out bytes.Buffer
	printDigest(&app{out: &out}, nil)
	if out.Len() != 0 {
		t.Errorf("no digest prints:\n%s", out.String())
	}
	r := 7.0 / 4
	third := 1.0 / 3
	two, none := 2, 0
	printDigest(&app{out: &out}, &metrics.Digest{
		Weeks: []metrics.DigestWeek{
			{Week: "2026-W42", Start: "2026-10-12", ThreadsToOperator: 2, AskedAgain: &none},
			{Week: "2026-W43", Start: "2026-10-19", ThreadsToOperator: 7, StoriesCompleted: 4, PerStory: &r, AskedAgain: &two, CuratorRuns: 1},
			{Week: "2026-W44", Start: "2026-10-26", ThreadsToOperator: 1, StoriesCompleted: 3, PerStory: &third},
		},
		Entries: metrics.DigestEntries{Active: 120, Compacted: 3, Absorbed: 5},
		Curator: metrics.DigestCurator{Runs: 3, LastRun: "2026-10-26T05:00:00Z"},
	})
	want := "\ndigest:\n" +
		"  2026-W42  2026-10-12  threads 2 · stories 0 · – per story · asked again 0 · curator runs 0\n" +
		"  2026-W43  2026-10-19  threads 7 · stories 4 · 1.75 per story · asked again 2 · curator runs 1\n" +
		"  2026-W44  2026-10-26  threads 1 · stories 3 · 0.33 per story · asked again – · curator runs 0\n" +
		"  entries  active 120 · compacted 3 · absorbed 5\n" +
		"  curator  runs 3 · last 2026-10-26T05:00:00Z\n"
	if out.String() != want {
		t.Errorf("printDigest:\n%s\nwant\n%s", out.String(), want)
	}
	out.Reset()
	printDigest(&app{out: &out}, &metrics.Digest{Weeks: []metrics.DigestWeek{}})
	if want := "\ndigest:\n  entries  active 0 · compacted 0 · absorbed 0\n  curator  runs 0 · never run\n"; out.String() != want {
		t.Errorf("printDigest with no curator run:\n%s\nwant\n%s", out.String(), want)
	}
}
