package serve

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// remindLab is a project owned by alex with flai serve's thread reminder
// over it: a clock, the host's channel, and a recording sender in place of
// the webhook and the command.
type remindLab struct {
	t     *testing.T
	root  string
	repo  *workitem.Repo
	dir   Dir
	at    time.Time
	ch    config.Notify
	fail  error // what the channel answers a send with
	sent  []Reminder
	logs  syncBuffer
	r     *reminder
	story string
}

// newRemindLab is the project with threads.remind_after set to after, ""
// for unset, and one story owned by alex.
func newRemindLab(t *testing.T, after string) *remindLab {
	t.Helper()
	lab := &remindLab{t: t, root: t.TempDir(), dir: DirFor(filepath.Join(t.TempDir(), "config.json")),
		at: time.Date(2026, 10, 20, 8, 0, 0, 0, time.UTC), ch: config.Notify{URL: "https://push.example/hook/s3cret"}}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "design", "docs"} {
		if err := os.MkdirAll(filepath.Join(lab.root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	lab.remindAfter(after)
	repo, err := workitem.Open(lab.root)
	if err != nil {
		t.Fatal(err)
	}
	lab.repo = repo
	lab.story = lab.newStory("Harbour lights")
	lab.restart()
	return lab
}

// remindAfter writes the manifest with threads.remind_after set to after,
// "" for unset.
func (lab *remindLab) remindAfter(after string) {
	lab.t.Helper()
	m := "version: 1\nname: Harbour\nkey: harbour\nowner: alex\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if after != "" {
		m += "threads:\n  remind_after: " + after + "\n"
	}
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// restart is a new flai serve over the same serve folder: it keeps in
// memory nothing of the last, only what it recorded.
func (lab *remindLab) restart() {
	lab.r = &reminder{e: Entry{Key: "harbour", Name: "Harbour", Root: lab.root}, dir: lab.dir,
		now: func() time.Time { return lab.at }, log: slog.New(slog.NewTextHandler(&lab.logs, nil)),
		channel: func() config.Notify { return lab.ch },
		send: func(_ context.Context, ch config.Notify, r Reminder) error {
			if ch != lab.ch {
				lab.t.Errorf("sent through %+v, want the host's channel %+v", ch, lab.ch)
			}
			lab.sent = append(lab.sent, r)
			return lab.fail
		},
		looks: &sync.WaitGroup{}}
}

// look is a look of the reminder; it answers the threads reminded by it.
func (lab *remindLab) look() []string {
	lab.t.Helper()
	before := len(lab.sent)
	lab.r.look(context.Background())
	var out []string
	for _, r := range lab.sent[before:] {
		out = append(out, r.Thread)
	}
	return out
}

func (lab *remindLab) after(d time.Duration) { lab.at = lab.at.Add(d) }

func (lab *remindLab) newStory(title string) string {
	lab.t.Helper()
	epic, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic " + title, Owner: "alex", Now: lab.at})
	if err != nil {
		lab.t.Fatal(err)
	}
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: epic.ID, Owner: "alex", Now: lab.at})
	if err != nil {
		lab.t.Fatal(err)
	}
	return st.ID
}

// finish sets the item's status to status, as accepting or cancelling it
// would leave it.
func (lab *remindLab) finish(id, status string) {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	it.Status = status
	if err := lab.repo.Save(it); err != nil {
		lab.t.Fatal(err)
	}
}

// ask opens a thread on on by author, now.
func (lab *remindLab) ask(on, author, title string) string {
	lab.t.Helper()
	th, err := threads.New(lab.repo, threads.NewOptions{Title: title, On: on, Author: author, Text: "Eight or nine?", Now: lab.at})
	if err != nil {
		lab.t.Fatal(err)
	}
	return th.ID
}

// recommend has the orchestrator recommend an answer on id, now.
func (lab *remindLab) recommend(id string) {
	lab.t.Helper()
	if _, err := threads.ReplyWith(lab.repo, id, "orchestrator", "Nine, as the design says.\n\nIt is the port the dashboard uses.", lab.at,
		threads.Marks{Recommendation: true, Source: threads.Source{Path: "design"}}); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *remindLab) reply(id, author, text string) {
	lab.t.Helper()
	if _, err := threads.Reply(lab.repo, id, author, text, lab.at); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *remindLab) confirm(id string) {
	lab.t.Helper()
	if _, err := threads.Confirm(lab.repo, id, "alex", lab.at); err != nil {
		lab.t.Fatal(err)
	}
}

// recorded is what reminders.json holds of the project, by thread.
func (lab *remindLab) recorded() map[string]reminded {
	lab.t.Helper()
	got, err := lab.dir.projectReminders(lab.root)
	if err != nil {
		lab.t.Fatal(err)
	}
	return got
}

func sameThreads(got []string, want ...string) bool {
	return strings.Join(got, ",") == strings.Join(want, ",")
}

// (1) A pending recommendation past remind_after is reminded once, then
// again after remind_every, not before.
func TestAPendingRecommendationIsRemindedAfterRemindAfterAndAgainEveryRemindEvery(t *testing.T) {
	lab := newRemindLab(t, "4h")
	id := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	lab.after(time.Minute)
	lab.recommend(id)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded at once of %v; nothing is due before threads.remind_after", got)
	}
	lab.after(4*time.Hour - time.Second)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded of %v a second before threads.remind_after", got)
	}
	lab.after(time.Second)
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("reminded of %v at threads.remind_after, want %s", got, id)
	}
	for _, d := range []time.Duration{time.Minute, time.Hour, 22 * time.Hour} {
		lab.after(d)
		if got := lab.look(); len(got) != 0 {
			t.Fatalf("reminded again of %v within threads.remind_every (24h)", got)
		}
	}
	lab.after(time.Hour - time.Minute)
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("reminded of %v 24h after the first reminder, want %s again", got, id)
	}
	// a remind_every the operator sets while flai serve runs is read at the next look
	m, _ := os.ReadFile(filepath.Join(lab.root, "system-flow.yaml"))
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), append(m, []byte("  remind_every: 2h\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.after(2 * time.Hour)
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("reminded of %v 2h after the last under remind_every 2h, want %s", got, id)
	}
	r := lab.sent[0]
	if r.Project != "Harbour" || r.Thread != id || r.Title != "Which port?" || r.Item != lab.story ||
		r.Href != "/items/"+lab.story+"?thread="+id || r.Since.Format(time.RFC3339) != "2026-10-20T08:01:00Z" ||
		r.AwaitingSeconds != int64(4*time.Hour/time.Second) {
		t.Fatalf("the reminder is %+v", r)
	}
	if rec := r.Recommendation; rec == nil || rec.Author != "orchestrator" || rec.FirstLine != "Nine, as the design says." || rec.Source == nil || rec.Source.Path != "design" {
		t.Fatalf("the reminder's recommendation is %+v", r.Recommendation)
	}
	if !strings.Contains(lab.logs.String(), "awaiting_hours=4 ") || !strings.Contains(lab.logs.String(), "thread="+id) {
		t.Fatalf("the serve log does not record the reminder with its thread and hours awaited:\n%s", lab.logs.String())
	}
	if strings.Contains(lab.logs.String(), "s3cret") {
		t.Fatalf("the serve log carries the webhook's URL:\n%s", lab.logs.String())
	}
}

// (2) The operator's answer or confirmation stops the reminders, and a new
// question starts the count over from the new entry.
func TestTheOperatorsAnswerStopsTheRemindersAndANewQuestionStartsTheCountOver(t *testing.T) {
	lab := newRemindLab(t, "4h")
	confirmed := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	lab.recommend(confirmed)
	answered := lab.ask(lab.story, "builder-"+lab.story, "Which colour?")
	lab.after(5 * time.Hour)
	if got := lab.look(); !sameThreads(got, confirmed, answered) {
		t.Fatalf("reminded of %v, want both threads", got)
	}
	lab.confirm(confirmed)
	lab.reply(answered, "alex", "Blue.")
	lab.after(24 * time.Hour)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded of %v after the operator confirmed and answered", got)
	}
	if got := lab.recorded(); len(got) != 0 {
		t.Fatalf("reminders.json still holds %v of threads that no longer await the operator", got)
	}
	lab.after(time.Hour)
	lab.reply(answered, "builder-"+lab.story, "Which blue?")
	asked := lab.at
	lab.after(4*time.Hour - time.Minute)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded of %v before threads.remind_after of the new question", got)
	}
	lab.after(time.Minute)
	if got := lab.look(); !sameThreads(got, answered) {
		t.Fatalf("reminded of %v at threads.remind_after of the new question, want %s", got, answered)
	}
	if got := lab.sent[len(lab.sent)-1].Since; !got.Equal(asked) {
		t.Fatalf("the reminder is for the entry at %s, want the new question's, %s", got, asked)
	}
	// a new entry while it is reminded starts the count over too
	lab.after(time.Hour)
	lab.reply(answered, "builder-"+lab.story, "Navy or sky?")
	lab.after(3 * time.Hour)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded of %v 3h after the newest entry", got)
	}
	lab.after(time.Hour)
	if got := lab.look(); !sameThreads(got, answered) {
		t.Fatalf("reminded of %v 4h after the newest entry, want %s", got, answered)
	}
}

// (3) A thread awaiting another agent, one on a done or cancelled item, and
// one younger than the threshold are not reminded; one on a document is.
func TestOnlyThreadsAwaitingTheOperatorOnOpenWorkPastTheThresholdAreReminded(t *testing.T) {
	lab := newRemindLab(t, "4h")
	toAgent := lab.ask(lab.story, "alex", "Is the port free?")
	done, cancelled := lab.newStory("Lighthouse"), lab.newStory("Buoy")
	lab.finish(done, workitem.Done)
	lab.finish(cancelled, workitem.Cancelled)
	onDone := lab.ask(done, "builder-"+done, "Which lamp?")
	onCancelled := lab.ask(cancelled, "builder-"+cancelled, "Which bell?")
	onDocument := lab.ask("system-flow.yaml", "flai serve", "Upgrade flai")
	resolved := lab.ask(lab.story, "builder-"+lab.story, "Which tide?")
	if _, err := threads.Resolve(lab.repo, resolved, "alex", "settled", lab.at); err != nil {
		t.Fatal(err)
	}
	lab.after(3 * time.Hour)
	young := lab.ask(lab.story, "builder-"+lab.story, "Which harbour?")
	lab.after(time.Hour)
	got := lab.look()
	if !sameThreads(got, onDocument) {
		t.Fatalf("reminded of %v, want %s on a document alone (not %s awaiting an agent, %s on a done story, %s on a cancelled one, %s resolved, %s an hour old)",
			got, onDocument, toAgent, onDone, onCancelled, resolved, young)
	}
	if r := lab.sent[0]; r.Item != "" || r.Href != "/threads?thread="+onDocument || r.Recommendation != nil {
		t.Fatalf("the reminder of a thread on a document is %+v", r)
	}
}

// (4) With threads.remind_after unset nothing is sent.
func TestUnsetRemindAfterSendsNothing(t *testing.T) {
	lab := newRemindLab(t, "")
	id := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	lab.recommend(id)
	lab.after(30 * 24 * time.Hour)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded of %v with threads.remind_after unset", got)
	}
	lab.remindAfter("0")
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded of %v with threads.remind_after 0", got)
	}
	if _, err := os.Stat(lab.dir.remindersFile()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("reminders.json written with reminders off: %v", err)
	}
	// set while flai serve runs, it is read at the next look
	lab.remindAfter("4h")
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("reminded of %v once threads.remind_after is set, want %s", got, id)
	}
}

// (5) A reminder writes nothing on the thread.
func TestAReminderLeavesTheThreadFileByteIdentical(t *testing.T) {
	lab := newRemindLab(t, "4h")
	id := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	lab.recommend(id)
	th, err := threads.Get(lab.repo, id)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(th.Path)
	if err != nil {
		t.Fatal(err)
	}
	lab.after(5 * time.Hour)
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("reminded of %v, want %s", got, id)
	}
	after, err := os.ReadFile(th.Path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("the thread file changed with the reminder:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if again, err := threads.Get(lab.repo, id); err != nil || again.PendingRecommendation() == nil || again.Status != th.Status {
		t.Fatalf("the thread no longer awaits the operator's confirmation after the reminder: %+v, %v", again, err)
	}
}

// (6) What was reminded survives a restart: a new reminder over the same
// serve folder repeats nothing.
func TestARestartDoesNotRepeatAReminder(t *testing.T) {
	lab := newRemindLab(t, "4h")
	id := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	lab.recommend(id)
	lab.after(5 * time.Hour)
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("reminded of %v, want %s", got, id)
	}
	lab.restart()
	lab.after(time.Hour)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("a restarted flai serve reminded again of %v within threads.remind_every", got)
	}
	if got := lab.recorded()[id]; !got.Sent.Equal(lab.at.Add(-time.Hour)) {
		t.Fatalf("reminders.json holds %+v of %s", got, id)
	}
	lab.after(23 * time.Hour)
	lab.restart()
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("a restarted flai serve reminded of %v 24h after the reminder, want %s", got, id)
	}
}

// (7) With no channel set the reminder is logged, and that none is set is
// warned of once per project per run.
func TestWithNoChannelTheReminderIsLoggedAndTheWarningGivenOncePerRun(t *testing.T) {
	lab := newRemindLab(t, "4h")
	lab.ch = config.Notify{}
	first := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	second := lab.ask(lab.story, "builder-"+lab.story, "Which colour?")
	lab.after(5 * time.Hour)
	lab.look()
	lab.after(24 * time.Hour)
	lab.look()
	if len(lab.sent) != 0 {
		t.Fatalf("sent %d reminders with no channel set", len(lab.sent))
	}
	logs := lab.logs.String()
	warnings, logged := 0, map[string]int{}
	for line := range strings.Lines(logs) {
		switch {
		case strings.Contains(line, "level=WARN") && strings.Contains(line, "hint="):
			warnings++
		case strings.Contains(line, "level=INFO") && strings.Contains(line, "awaiting_hours="):
			for _, id := range []string{first, second} {
				if strings.Contains(line, "thread="+id+" ") {
					logged[id]++
				}
			}
		}
	}
	if warnings != 1 || logged[first] != 2 || logged[second] != 2 {
		t.Fatalf("warned %d times that no channel is set, want 1; logged %v, want each thread twice:\n%s", warnings, logged, logs)
	}
	// a channel the operator sets while flai serve runs is used at the next look
	lab.ch = config.Notify{Command: "notify-send flai"}
	lab.after(24 * time.Hour)
	if got := lab.look(); !sameThreads(got, first, second) {
		t.Fatalf("reminded through the channel set while flai serve runs of %v, want both", got)
	}
}

// A send that fails is logged with its error, and is not tried again
// before threads.remind_every: one attempt, as ADR-0205 says.
func TestAFailedSendIsLoggedAndNotRetriedUntilRemindEvery(t *testing.T) {
	lab := newRemindLab(t, "4h")
	id := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	lab.fail = errors.New("post the reminder of " + id + " to the webhook at push.example: it answered 502 Bad Gateway")
	lab.after(5 * time.Hour)
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("tried %v, want %s", got, id)
	}
	if logs := lab.logs.String(); !strings.Contains(logs, "level=WARN") || !strings.Contains(logs, "502 Bad Gateway") || !strings.Contains(logs, "thread="+id) {
		t.Fatalf("the failed send is not in the serve log with its thread and error:\n%s", logs)
	}
	lab.after(time.Minute)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("tried again at the next look: %v", got)
	}
}

// (8) I-0151: TH-0436's recommendation awaited the operator's confirmation
// for 25 hours, unseen off the dashboard. With threads.remind_after 4h, a
// recommendation pending that long is reminded at the first look.
func TestARecommendationPending25HoursIsRemindedAtTheFirstLook(t *testing.T) {
	lab := newRemindLab(t, "4h")
	id := lab.ask(lab.story, "builder-"+lab.story, "Run flai host upgrade on the host?")
	lab.after(40 * time.Second)
	lab.recommend(id)
	lab.after(25 * time.Hour)
	if got := lab.look(); !sameThreads(got, id) {
		t.Fatalf("a recommendation pending 25h reminded %v, want %s", got, id)
	}
	r := lab.sent[0]
	if r.AwaitingSeconds != int64(25*time.Hour/time.Second) || r.Recommendation == nil || r.Recommendation.Author != "orchestrator" {
		t.Fatalf("the reminder is %+v", r)
	}
	lab.after(time.Minute)
	if got := lab.look(); len(got) != 0 {
		t.Fatalf("reminded again a minute later: %v", got)
	}
}

// A look in the background is waited for, and none starts once serve stops.
func TestSoonLooksInTheBackgroundAndNotOnceStopped(t *testing.T) {
	lab := newRemindLab(t, "4h")
	id := lab.ask(lab.story, "builder-"+lab.story, "Which port?")
	lab.after(5 * time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	lab.r.soon(ctx)
	lab.r.looks.Wait()
	if len(lab.sent) != 1 || lab.sent[0].Thread != id {
		t.Fatalf("the background look sent %+v, want %s", lab.sent, id)
	}
	cancel()
	lab.after(24 * time.Hour)
	lab.r.soon(ctx)
	lab.r.looks.Wait()
	if len(lab.sent) != 1 {
		t.Fatalf("a look started after serve stopped and sent %+v", lab.sent[1:])
	}
}
