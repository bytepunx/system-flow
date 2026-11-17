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
func testStatsReportsTheMergeQueue(t *testing.T) {
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	inEnv(t, "LOG_FORMAT", "json")
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
	// week 40: a batch of two fails and is bisected, in one step, to
	// S-0002, sent back after an hour; S-0001 lands alone after three: two
	// verification runs for one story landed (S-0552)
	events := []mergequeue.Event{
		ev("2026-09-29T09:00:00Z", enq, "", "S-0001"),
		ev("2026-09-29T09:00:00Z", enq, "", "S-0002"),
		ev("2026-09-29T09:10:00Z", built, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T09:40:00Z", verified, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T09:41:00Z", bisected, "B1", "S-0001", "S-0002"),
		ev("2026-09-29T09:50:00Z", mergequeue.EventStepVerified, "B1", "S-0001"),
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
		Batches: 2, BatchSizeMean: 1.5, BatchSizeMax: 2, Bisections: 1, Verifications: 2, VerificationsPerLanded: 2}
	if got := rep.Queue.Weeks[4]; got != want {
		t.Errorf("week 40 = %+v, want %+v", got, want)
	}
	if !strings.Contains(out, `"verifications": 2,`) || !strings.Contains(out, `"verifications_per_landed": 2`) {
		t.Errorf("stats --json does not carry the verifications:\n%s", out)
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
		"  batches        built 2 · mean size 1.5 · largest 2 · bisections 1 · verifications 2 · 2.0 per story landed\n" +
		"  by week:\n" +
		"    2026-W40 (2026-09-28)  stories landed 1 · sent back 1 · mean wait 2h · longest 3h; batches built 2 · mean size 1.5 · largest 2 · bisections 1 · verifications 2 · 2.0 per story landed\n"
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
	t.Parallel()
	var out bytes.Buffer
	printQueueStats(&app{configPath: configOf(t), out: &out}, &metrics.Queue{Weeks: []metrics.QueueWeek{{Week: "2026-W40", Start: "2026-09-28"}}})
	if out.Len() != 0 {
		t.Errorf("an empty queue prints:\n%s", out.String())
	}
	rep := &metrics.Report{
		WindowStart: "2026-09-03T20:00:00Z", Type: workitem.Story, WindowDays: 30,
		Waiting: metrics.Waiting{Weeks: []metrics.WaitWeek{{Week: "2026-W40", Start: "2026-09-28", Items: 1}}},
		Queue:   &metrics.Queue{Weeks: []metrics.QueueWeek{{Week: "2026-W40", Start: "2026-09-28", Landed: 1}}},
	}
	printSummary(&app{configPath: configOf(t), out: &out}, rep, false)
	text := out.String()
	waiting, queue := strings.Index(text, "\nwaiting, over 1 completed:\n"), strings.Index(text, "\nmerge queue:\n")
	if waiting < 0 || queue < waiting {
		t.Errorf("the queue's section is not after the waiting:\n%s", text)
	}
}

// S-0552: the batches' text names the verification runs, and them per story
// landed only when a story landed; ADR-0282: and the batches dropped only
// when one was.
func TestQueueBatchesTextNamesTheVerificationsPerStoryLanded(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		w    metrics.QueueWeek
		want string
	}{
		{metrics.QueueWeek{Landed: 2, Batches: 1, BatchSizeMean: 3, BatchSizeMax: 3, Bisections: 1, Verifications: 3, VerificationsPerLanded: 1.5},
			"built 1 · mean size 3.0 · largest 3 · bisections 1 · verifications 3 · 1.5 per story landed"},
		{metrics.QueueWeek{Batches: 1, BatchSizeMean: 1, BatchSizeMax: 1, Verifications: 1},
			"built 1 · mean size 1.0 · largest 1 · bisections 0 · verifications 1"},
		{metrics.QueueWeek{Batches: 3, BatchSizeMean: 1, BatchSizeMax: 1, Verifications: 3, Landed: 1, VerificationsPerLanded: 3, Dropped: 2},
			"built 3 · mean size 1.0 · largest 1 · bisections 0 · verifications 3 · 3.0 per story landed · dropped 2"},
	} {
		if got := queueBatchesText(tc.w); got != tc.want {
			t.Errorf("queueBatchesText(%+v) = %q, want %q", tc.w, got, tc.want)
		}
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
func testStatsReportsTheDigest(t *testing.T) {
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	inEnv(t, "LOG_FORMAT", "json")
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
	t.Parallel()
	var out bytes.Buffer
	printDigest(&app{configPath: configOf(t), out: &out}, nil)
	if out.Len() != 0 {
		t.Errorf("no digest prints:\n%s", out.String())
	}
	r := 7.0 / 4
	third := 1.0 / 3
	two, none := 2, 0
	printDigest(&app{configPath: configOf(t), out: &out}, &metrics.Digest{
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
	printDigest(&app{configPath: configOf(t), out: &out}, &metrics.Digest{Weeks: []metrics.DigestWeek{}})
	if want := "\ndigest:\n  entries  active 0 · compacted 0 · absorbed 0\n  curator  runs 0 · never run\n"; out.String() != want {
		t.Errorf("printDigest with no curator run:\n%s\nwant\n%s", out.String(), want)
	}
}

// TestLoggingJSON runs the tests of flai stats that read its JSON log events,
// as together runs them: LOG_FORMAT json, set once here, so that they run in
// parallel with each other (S-0548).
func TestLoggingJSON(t *testing.T) {
	t.Setenv("LOG_FORMAT", "json")
	together(t,
		testStatsReportsPlanningWaitingAndClaims,
		testStatsReportsCoordination,
		testStatsReportsTheMergeQueue,
		testStatsReportsTheDigest,
	)
}

// S-0551: flai stats prints the queue time as the ready wait, and under it
// the share of the window with each number of stories in progress at once
// through the board's limit, with the mean while busy; --json and the
// dashboard's read carry them under concurrency and ready_wait (ADR-0231).
func TestStatsReportsConcurrencyAndTheReadyWait(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("LOG_FORMAT", "json")
	root := tempProject(t)
	board := "---\ntitle: Board\nstatus: active\nwip_limits:\n  in-progress: 3\n---\n\n# Board\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// on the 3rd of October, S-0001 waits 2h in ready and is in progress
	// from 08:00 to 14:00, then in review; S-0002 waits 4h and is in progress
	// from 10:00 to 12:00: 1 at once for 4h, 2 for 2h, 0 for the other 18h
	for _, moves := range [][]workitem.Transition{
		{{To: workitem.Ready, At: "2026-10-03T06:00:00Z"}, {To: workitem.InProgress, At: "2026-10-03T08:00:00Z"},
			{To: workitem.Review, At: "2026-10-03T14:00:00Z"}, {To: workitem.Done, At: "2026-10-03T15:00:00Z"}},
		{{To: workitem.Ready, At: "2026-10-03T06:00:00Z"}, {To: workitem.InProgress, At: "2026-10-03T10:00:00Z"},
			{To: workitem.Done, At: "2026-10-03T12:00:00Z"}},
	} {
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Runs", Owner: "t", Now: time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)})
		if err != nil {
			t.Fatal(err)
		}
		it.Transitions, it.Status = moves, workitem.Done
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)

	out, errOut, code := runInAt(t, root, at, "stats", "--since", "1d")
	wantText := "  ready wait   p50 2h · p85 4h · max 4h · mean 3h (n=2)\n" +
		"  in progress  0 75% · 1 17% · 2 8% · 3 0% · mean 1.3 while busy\n"
	if code != 0 || !strings.Contains(out, wantText) {
		t.Errorf("text: %d %s, want it to hold %q, got\n%s", code, errOut, wantText, out)
	}

	out, errOut, code = runInAt(t, root, at, "stats", "--since", "1d", "--json")
	if code != 0 {
		t.Fatalf("stats --json: %s", errOut)
	}
	var rep struct {
		Concurrency map[string]any `json:"concurrency"`
		ReadyWait   map[string]any `json:"ready_wait"`
	}
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatal(err)
	}
	h := 3600.0
	wantConc := map[string]any{
		"window_seconds": 24 * h, "limit": 3.0,
		"seconds":   []any{18 * h, 4 * h, 2 * h, 0.0},
		"share":     []any{0.75, 4.0 / 24, 2.0 / 24, 0.0},
		"mean_busy": 8.0 / 6,
	}
	if !reflect.DeepEqual(rep.Concurrency, wantConc) {
		t.Errorf("concurrency = %v, want %v", rep.Concurrency, wantConc)
	}
	if want := map[string]any{"count": 2.0, "p50_seconds": 2 * h, "p85_seconds": 4 * h}; !reflect.DeepEqual(rep.ReadyWait, want) {
		t.Errorf("ready_wait = %v, want %v", rep.ReadyWait, want)
	}
	// the dashboard's read answers the same
	if _, rerr := sameAnswer(t, root, hostapi.Host{}, "stats.get", `{}`, "stats"); rerr != nil {
		t.Errorf("stats.get: %+v", rerr)
	}
}

// timedWorkflow is the baseline workflow file flai carries with two timers
// of its own after the baseline's: the wait from the first review to
// acceptance, and the time in progress and in review together.
func timedWorkflow(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "internal", "workflow", "baseline", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	last := "  - {name: review, title: Review time, sum: [review]}\n"
	if !strings.Contains(string(data), last) {
		t.Fatalf("the baseline workflow file has no line %q", last)
	}
	return strings.Replace(string(data), last, last+
		"  - {name: acceptance-wait, title: Acceptance wait, from: {first: review}, to: completed}\n"+
		"  - {name: in-work, title: In work, sum: [in-progress, review]}\n", 1)
}

// S-0752 criteria 2, 3, and 6: a workflow's own timers, added to its file
// after the items moved, read their recorded transitions: flai stats --json
// lists them, gives each item's seconds and each summary's distribution, and
// flai stats prints them after the ready wait; on the baseline there are
// none (ADR-0317).
func TestStatsReportsTheWorkflowsOwnTimers(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("LOG_FORMAT", "json")
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)
	project := func(t *testing.T, flow string) string {
		root := tempProject(t)
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		tr := func(to, at string) workitem.Transition { return workitem.Transition{To: to, At: at} }
		for _, s := range []struct {
			created time.Time
			status  string
			moves   []workitem.Transition
		}{
			// sent back from review once, then accepted: 6h from the first
			// review to acceptance, 12h in progress and in review
			{time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC), workitem.Done, []workitem.Transition{
				tr(workitem.Ready, "2026-10-02T06:00:00Z"), tr(workitem.InProgress, "2026-10-02T08:00:00Z"),
				tr(workitem.Review, "2026-10-02T14:00:00Z"), tr(workitem.InProgress, "2026-10-02T15:00:00Z"),
				tr(workitem.Review, "2026-10-02T17:00:00Z"), tr(workitem.Done, "2026-10-02T20:00:00Z")}},
			// done weeks ago, within the window: 2d and 3d
			{time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC), workitem.Done, []workitem.Transition{
				tr(workitem.Ready, "2026-09-10T01:00:00Z"), tr(workitem.InProgress, "2026-09-11T00:00:00Z"),
				tr(workitem.Review, "2026-09-12T00:00:00Z"), tr(workitem.Done, "2026-09-14T00:00:00Z")}},
			// open in progress since 10:00: 10h in work until now, never in
			// review, and not completed, so in no distribution
			{time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC), workitem.InProgress, []workitem.Transition{
				tr(workitem.Ready, "2026-10-03T09:30:00Z"), tr(workitem.InProgress, "2026-10-03T10:00:00Z")}},
		} {
			it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Runs", Owner: "t", Now: s.created})
			if err != nil {
				t.Fatal(err)
			}
			it.Transitions, it.Status = s.moves, s.status
			if err := repo.Save(it); err != nil {
				t.Fatal(err)
			}
		}
		if flow != "" {
			if err := os.MkdirAll(filepath.Join(root, "design", "workflows"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "design", "workflows", "board.md"), []byte(flow), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	h, d := 3600.0, 86400.0

	root := project(t, timedWorkflow(t))
	out, errOut, code := runInAt(t, root, at, "stats", "--json")
	if code != 0 {
		t.Fatalf("stats --json: %s", errOut)
	}
	var rep struct {
		Timers []metrics.TimerDefinition `json:"timers"`
		Items  []struct {
			ID     string             `json:"id"`
			Timers map[string]float64 `json:"timers"`
		} `json:"items"`
		Summary struct {
			Timers map[string]metrics.Distribution `json:"timers"`
		} `json:"summary"`
	}
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatal(err)
	}
	every := []string{workitem.Epic, workitem.Story, workitem.Task}
	wantDefs := []metrics.TimerDefinition{
		{Name: "acceptance-wait", Title: "Acceptance wait", Types: every, Counts: "accepted"},
		{Name: "in-work", Title: "In work", Types: every, Counts: "accepted"},
	}
	if !reflect.DeepEqual(rep.Timers, wantDefs) {
		t.Errorf("timers = %+v, want %+v", rep.Timers, wantDefs)
	}
	items := map[string]map[string]float64{}
	for _, it := range rep.Items {
		items[it.ID] = it.Timers
	}
	wantItems := map[string]map[string]float64{
		"S-0001": {"acceptance-wait": 6 * h, "in-work": 12 * h},
		"S-0002": {"acceptance-wait": 2 * d, "in-work": 3 * d},
		"S-0003": {"in-work": 10 * h},
	}
	if !reflect.DeepEqual(items, wantItems) {
		t.Errorf("items' timers = %v, want %v", items, wantItems)
	}
	wantSummary := map[string]metrics.Distribution{
		"acceptance-wait": {Count: 2, P50: 6 * h, P85: 2 * d, Max: 2 * d, Mean: (6*h + 2*d) / 2},
		"in-work":         {Count: 2, P50: 12 * h, P85: 3 * d, Max: 3 * d, Mean: (12*h + 3*d) / 2},
	}
	if !reflect.DeepEqual(rep.Summary.Timers, wantSummary) {
		t.Errorf("summary.timers = %+v, want %+v", rep.Summary.Timers, wantSummary)
	}

	out, errOut, code = runInAt(t, root, at, "stats", "--by", "nature")
	wantText := "  ready wait   p50 2h · p85 23h · max 23h · mean 12h30m (n=2)\n" +
		"  Acceptance wait p50 6h · p85 2d · max 2d · mean 1d3h (n=2)\n" +
		"  In work      p50 12h · p85 3d · max 3d · mean 1d18h (n=2)\n"
	if code != 0 || strings.Count(out, wantText) != 2 {
		t.Errorf("text: %d %s, want the summary and its one group to hold %q, got\n%s", code, errOut, wantText, out)
	}

	base := project(t, "")
	out, errOut, code = runInAt(t, base, at, "stats", "--json", "--by", "nature")
	if code != 0 {
		t.Fatalf("baseline stats --json: %s", errOut)
	}
	var raw struct {
		Timers  json.RawMessage  `json:"timers"`
		Summary map[string]any   `json:"summary"`
		Groups  []map[string]any `json:"groups"`
		Items   []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		t.Fatal(err)
	}
	if raw.Timers != nil {
		t.Errorf("baseline: timers = %s, want none", raw.Timers)
	}
	for _, m := range append(append([]map[string]any{raw.Summary}, raw.Groups...), raw.Items...) {
		if v, ok := m["timers"]; ok {
			t.Errorf("baseline: %v holds timers %v", m["id"], v)
		}
	}
	if len(raw.Groups) == 0 || len(raw.Items) != 3 {
		t.Errorf("baseline: %d groups and %d items, want some and 3", len(raw.Groups), len(raw.Items))
	}
	out, errOut, code = runInAt(t, base, at, "stats")
	if code != 0 || strings.Contains(out, "Acceptance wait") || strings.Contains(out, "In work") ||
		!strings.Contains(out, "  ready wait   p50 2h · p85 23h · max 23h · mean 12h30m (n=2)\n") {
		t.Errorf("baseline text: %d %s, want the ready wait and no timer of the workflow's own, got\n%s", code, errOut, out)
	}
}
