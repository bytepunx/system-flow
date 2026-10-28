package serve

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// confirmLab is remindLab's project, owned by alex, with flai serve's
// confirmer over it in place of the reminder: a clock, and the orchestrate
// host action as a field.
type confirmLab struct {
	*remindLab
	orchestrate bool
	c           *confirmer
}

// newConfirmLab is the project with orchestration.permissions.confirm_after
// set to after, "" for unset, and the orchestrate host action on.
func newConfirmLab(t *testing.T, after string) *confirmLab {
	t.Helper()
	lab := &confirmLab{remindLab: newRemindLab(t, ""), orchestrate: true}
	lab.confirmAfter(after)
	lab.c = &confirmer{e: Entry{Key: "harbour", Name: "Harbour", Root: lab.root}, now: func() time.Time { return lab.at },
		log: slog.New(slog.NewTextHandler(&lab.logs, nil)), orchestrating: func() bool { return lab.orchestrate }}
	return lab
}

// confirmAfter writes the manifest with confirm_after set to after, "" for
// unset.
func (lab *confirmLab) confirmAfter(after string) {
	lab.t.Helper()
	m := "version: 1\nname: Harbour\nkey: harbour\nowner: alex\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if after != "" {
		m += "orchestration:\n  permissions:\n    answer_threads: recommend\n    confirm_after: " + after + "\n"
	}
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *confirmLab) look() { lab.c.look(context.Background()) }

// recommendAs has author recommend an answer on id, now, marked for the
// operator alone when operatorOnly.
func (lab *confirmLab) recommendAs(id, author string, operatorOnly bool) {
	lab.t.Helper()
	if _, err := threads.ReplyWith(lab.repo, id, author, "Nine, as the design says.", lab.at,
		threads.Marks{Recommendation: true, OperatorOnly: operatorOnly, Source: threads.Source{Path: "design"}}); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *confirmLab) thread(id string) *threads.Thread {
	lab.t.Helper()
	th, err := threads.Get(lab.repo, id)
	if err != nil {
		lab.t.Fatal(err)
	}
	return th
}

func (lab *confirmLab) file(id string) string {
	lab.t.Helper()
	data, err := os.ReadFile(lab.thread(id).Path)
	if err != nil {
		lab.t.Fatal(err)
	}
	return string(data)
}

func (lab *confirmLab) activity() *workitem.Activity {
	lab.t.Helper()
	a, err := lab.repo.Activity(workitem.ActivityOrchestrator)
	if err != nil {
		lab.t.Fatal(err)
	}
	return a
}

// I-0195: a planner's question the orchestrator recommended on within
// minutes waited on the operator for a day or more. With confirm_after set,
// flai serve confirms it once it has waited that long, and not before; the
// entry is its own, the thread answered, and the confirmation logged in the
// orchestrator's activity at no cost, once.
func TestARecommendationLeftUnansweredPastConfirmAfterIsConfirmedByFlaiServe(t *testing.T) {
	lab := newConfirmLab(t, "48h")
	id := lab.ask(lab.story, "planner-"+lab.story, "Which port?")
	lab.after(3 * time.Minute)
	lab.recommendAs(id, "orchestrator", false)
	recAt := lab.at.Format(workitem.TimeFormat)
	lab.after(48*time.Hour - time.Second)
	lab.look()
	if th := lab.thread(id); th.PendingRecommendation() == nil || len(lab.activity().Entries) != 0 {
		t.Fatalf("confirmed a second before confirm_after:\n%s", lab.file(id))
	}
	lab.after(time.Second)
	lab.look()
	th := lab.thread(id)
	es := th.Entries()
	last := es[len(es)-1]
	if last.Author != threads.TimeoutAuthor || !last.ConfirmsByTimeout() || th.Status != "answered" || th.PendingRecommendation() != nil {
		t.Fatalf("not confirmed by flai serve at confirm_after (status %s):\n%s", th.Status, lab.file(id))
	}
	if want := "Confirmed the recommendation of " + recAt + " orchestrator after 48h with no answer from the operator"; !strings.Contains(last.Text, want) {
		t.Errorf("the confirmation says %q, want it to say %q", last.Text, want)
	}
	a := lab.activity()
	if len(a.Entries) != 1 {
		t.Fatalf("the orchestrator's activity has %d entries, want the confirmation's alone", len(a.Entries))
	}
	e := a.Entries[0]
	if e.Seconds != 0 || e.Cost != 0 || a.AccruedSeconds != 0 || a.AccruedCost != 0 {
		t.Errorf("the confirmation was charged %ds and $%v (totals %ds, $%v), want nothing", e.Seconds, e.Cost, a.AccruedSeconds, a.AccruedCost)
	}
	if e.Trigger != ConfirmTrigger {
		t.Errorf("the activity's trigger is %q, want %q", e.Trigger, ConfirmTrigger)
	}
	// the orchestrator's next activity is measured from its own last entry,
	// not from flai serve's confirmation
	m, err := readActivity(Dir(t.TempDir()), lab.root, "harbour", workitem.ActivityOrchestrator)
	if err != nil {
		t.Fatalf("readActivity: %v", err)
	}
	if !m.last.IsZero() {
		t.Errorf("readActivity measures from %v, want the confirmation passed over", m.last)
	}
	if strings.Join(e.Items, ",") != id+","+lab.story || !e.At.Equal(lab.at) {
		t.Errorf("the activity names %v at %s, want %s and %s at %s", e.Items, e.At, id, lab.story, lab.at)
	}
	if want := "flai serve confirmed the orchestrator's recommendation of " + recAt + " on " + id + " after 48h with no answer from the operator (orchestration.permissions.confirm_after)"; e.Summary != want {
		t.Errorf("the activity says %q, want %q", e.Summary, want)
	}
	logs := lab.logs.String()
	for _, f := range []string{"thread=" + id, "item=" + lab.story, "waited_hours=48", "project=harbour", "component=serve"} {
		if !strings.Contains(logs, f) {
			t.Errorf("the log lacks %s:\n%s", f, logs)
		}
	}

	// a second look writes nothing
	before := lab.file(id)
	lab.after(time.Hour)
	lab.look()
	if after := lab.file(id); after != before || len(lab.activity().Entries) != 1 {
		t.Errorf("a second look wrote again:\n%s", after)
	}
	if strings.Contains(lab.logs.String(), "level=WARN") {
		t.Errorf("a look warned:\n%s", lab.logs.String())
	}
}

// ADR-0206 items 2 and 3: nothing is confirmed while the setting is unset,
// the orchestrate host action is off, or the thread does not qualify.
func TestNoRecommendationIsConfirmedThatDoesNotQualify(t *testing.T) {
	agent := func(lab *confirmLab) string { return "agent-" + lab.story }
	cases := []struct {
		name  string
		after string
		setup func(lab *confirmLab) string // makes the thread and answers its ID
	}{
		{"confirm_after unset", "", func(lab *confirmLab) string {
			id := lab.ask(lab.story, agent(lab), "Which port?")
			lab.recommendAs(id, "orchestrator", false)
			return id
		}},
		{"confirm_after 0", "0", func(lab *confirmLab) string {
			id := lab.ask(lab.story, agent(lab), "Which port?")
			lab.recommendAs(id, "orchestrator", false)
			return id
		}},
		{"orchestrate off", "1h", func(lab *confirmLab) string {
			lab.orchestrate = false
			id := lab.ask(lab.story, agent(lab), "Which port?")
			lab.recommendAs(id, "orchestrator", false)
			return id
		}},
		{"for the operator alone", "1h", func(lab *confirmLab) string {
			id := lab.ask(lab.story, agent(lab), "Buy the domain?")
			lab.recommendAs(id, "orchestrator", true)
			return id
		}},
		{"opened by the operator", "1h", func(lab *confirmLab) string {
			id := lab.ask(lab.story, "alex", "Which port?")
			lab.recommendAs(id, "orchestrator", false)
			return id
		}},
		{"opened by flai serve", "1h", func(lab *confirmLab) string {
			id := lab.ask(lab.story, threads.TimeoutAuthor, "Upgrade flai?")
			lab.recommendAs(id, "orchestrator", false)
			return id
		}},
		{"a permission prompt", "1h", func(lab *confirmLab) string {
			id := lab.ask(lab.story, agent(lab), "Allow Write .claude/settings.json?")
			lab.recommendAs(id, "orchestrator", false)
			return id
		}},
		{"recommended by another", "1h", func(lab *confirmLab) string {
			id := lab.ask(lab.story, agent(lab), "Which port?")
			lab.recommendAs(id, "planner-"+lab.story, false)
			return id
		}},
		{"the operator wrote since", "1h", func(lab *confirmLab) string {
			id := lab.ask(lab.story, agent(lab), "Which port?")
			lab.recommendAs(id, "orchestrator", false)
			lab.after(time.Minute)
			lab.reply(id, "alex", "Let me think about it.")
			lab.after(time.Minute)
			lab.reply(id, agent(lab), "Waiting on you.")
			return id
		}},
		{"on a story done", "1h", func(lab *confirmLab) string {
			id := lab.ask(lab.story, agent(lab), "Which port?")
			lab.recommendAs(id, "orchestrator", false)
			lab.finish(lab.story, workitem.Done)
			return id
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lab := newConfirmLab(t, tc.after)
			id := tc.setup(lab)
			before := lab.file(id)
			lab.after(72 * time.Hour)
			lab.look()
			if after := lab.file(id); after != before {
				t.Errorf("the thread was written:\n%s", after)
			}
			if n := len(lab.activity().Entries); n != 0 {
				t.Errorf("the orchestrator's activity has %d entries, want none", n)
			}
			if strings.Contains(lab.logs.String(), "level=WARN") {
				t.Errorf("a thread that does not qualify was warned of:\n%s", lab.logs.String())
			}
		})
	}
}

// A setting that is not a duration confirms nothing and is warned of once,
// not every minute.
func TestAConfirmAfterNotADurationIsWarnedOfOnce(t *testing.T) {
	lab := newConfirmLab(t, "two days")
	id := lab.ask(lab.story, "agent-"+lab.story, "Which port?")
	lab.recommendAs(id, "orchestrator", false)
	before := lab.file(id)
	for range 3 {
		lab.after(24 * time.Hour)
		lab.look()
	}
	if after := lab.file(id); after != before {
		t.Errorf("the thread was written:\n%s", after)
	}
	if n := strings.Count(lab.logs.String(), "level=WARN"); n != 1 {
		t.Errorf("warned %d times, want once:\n%s", n, lab.logs.String())
	}
}

// A confirmation that cannot be written is warned of once per thread and
// error, not at every look, and logs no activity.
func TestAConfirmationThatFailsIsWarnedOfOncePerThread(t *testing.T) {
	lab := newConfirmLab(t, "1h")
	id := lab.ask(lab.story, "agent-"+lab.story, "Which port?")
	lab.recommendAs(id, "orchestrator", false)
	dir := threads.Dir(lab.repo)
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	for range 3 {
		lab.after(time.Hour)
		lab.look()
	}
	logs := lab.logs.String()
	if n := strings.Count(logs, "level=WARN"); n != 1 || !strings.Contains(logs, "thread="+id) {
		t.Errorf("warned %d times, want once naming %s:\n%s", n, id, logs)
	}
	if n := len(lab.activity().Entries); n != 0 {
		t.Errorf("the orchestrator's activity has %d entries, want none", n)
	}
}
