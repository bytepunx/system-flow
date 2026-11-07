//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// analyzeLab is an agentLab whose command is an analyzer stub: it writes
// what it was given to analyze.txt, prints stream-analyze as its output when
// there is one, writes the report named in report under design/analysis, and
// the folder's index after it, when there is one, waits while held until
// released, and exits with the code in exit-analyze when there is one. The
// analyze action is on.
func analyzeLab(t *testing.T) *agentLab {
	t.Helper()
	lab := newAgentLab(t)
	stub := filepath.Join(lab.outDir, "analyze-agent")
	d := lab.outDir
	script := "#!/bin/sh\n" +
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: ${FLAI_STORY-unset}\"; echo \"role: $FLAI_ROLE\"; echo \"item: ${FLAI_ITEM-unset}\"; echo \"focus: ${FLAI_FOCUS-unset}\"; echo \"session: $FLAI_SESSION\"; echo \"by: $FLAI_STARTED_BY\"; } > \"" + d + "/analyze.txt\"\n" +
		"[ -f \"" + d + "/stream-analyze\" ] && cat \"" + d + "/stream-analyze\"\n" +
		"if [ -f \"" + d + "/report\" ]; then mkdir -p design/analysis; echo '# Report' > \"design/analysis/$(cat \"" + d + "/report\")\"; echo '# Analysis' > design/analysis/README.md; fi\n" +
		"while [ -f \"" + d + "/hold\" ] && [ ! -f \"" + d + "/release-analyze\" ]; do sleep 0.05; done\n" +
		"touch \"" + d + "/ended-$$\"\n" +
		"[ -f \"" + d + "/exit-analyze\" ] && exit \"$(cat \"" + d + "/exit-analyze\")\"\n" +
		"exit 0\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg.Command, lab.cfg.Analyze = []string{stub, "story={story}"}, true
	return lab
}

func (lab *agentLab) analyze(focus, trigger string) (*AgentRun, error) {
	return Analyze(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, focus, trigger)
}

// analyzeHere has the lab's launcher start the analyzer, as the operator
// asked, and wait for it, as a serving flai would.
func (lab *agentLab) analyzeHere(focus string) bool {
	lab.l.mu.Lock()
	defer lab.l.mu.Unlock()
	started, _ := lab.l.analyze(context.Background(), lab.cfg, nil, focus, TriggerAsked)
	return started
}

// analyzerGiven is what the analyzer stub was given, once it has written it.
func (lab *agentLab) analyzerGiven() string {
	lab.t.Helper()
	out := filepath.Join(lab.outDir, "analyze.txt")
	waitFor(lab.t, "the analyzer stub ran", func() bool {
		data, _ := os.ReadFile(out)
		return strings.Contains(string(data), "by: ")
	})
	data, _ := os.ReadFile(out)
	return string(data)
}

// analyzerRuns has the stub write the report named, and print a stream of
// two calls from at, a minute apart, and final as its last word.
func (lab *agentLab) analyzerRuns(report, final string, at time.Time) {
	lab.t.Helper()
	// each run a session of its own, as each run is
	session := "s-" + at.Format("150405")
	end := strings.Replace(streamFinal(final, 2000, 0.5), `"session_id":"s"`, `"session_id":"`+session+`"`, 1)
	stream := strings.Join([]string{streamCall(session, session+"-m1", at, 999), streamCall(session, session+"-m2", at.Add(time.Minute), 999), end}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-analyze"), []byte(stream), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	_ = os.Remove(filepath.Join(lab.outDir, "report"))
	if report != "" {
		if err := os.WriteFile(filepath.Join(lab.outDir, "report"), []byte(report), 0o644); err != nil {
			lab.t.Fatal(err)
		}
	}
}

// S-0223: the analyzer is started for the project in its main checkout, as
// the analyzer with its focus, logged as the analyzer's run, and recorded as
// the project's analyzer, apart from the orchestrator's and the stories'
// runs, with what started it; a second run is refused while the first runs,
// naming it.
func TestTheAnalyzerIsStartedForTheProject(t *testing.T) {
	t.Setenv("FLAI_ITEM", "E-0999") // the session that ran flai analyze
	lab := analyzeLab(t)
	lab.hold()
	run, err := lab.analyze("risk", TriggerAsked)
	if err != nil {
		t.Fatal(err)
	}
	if run.Agent != "analyzer" || run.Focus != "risk" || run.Trigger != "asked" || run.Story != "" || run.Item != "" || run.PID == 0 || run.Session == "" || run.Harness != "command" || run.Command != "analyze-agent" {
		t.Errorf("run = %+v", run)
	}
	if data, err := os.ReadFile(lab.o.Dir.agents()); err != nil || !strings.Contains(string(data), `"analyzer": {`) || !strings.Contains(string(data), `"focus": "risk"`) {
		t.Errorf("serve/agents.json does not record the analyzer's run (%v):\n%s", err, data)
	}
	logs, err := lab.o.Dir.ActivityLogs("t", workitem.ActivityAnalyzer)
	if err != nil || len(logs) != 1 || logs[0] != run.Log {
		t.Errorf("analyzer logs = %v (%v), want the run's %s", logs, err, run.Log)
	}
	got := lab.analyzerGiven()
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{"args: story=\n", "dir: " + real + "\n", "agent: analyzer\n", "story: unset\n", "role: analyze\n", "item: unset\n", "focus: risk\n", "session: 2", "by: flai-serve\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the analyzer was given\n%s\nwithout %q", got, want)
		}
	}
	if st := lab.state(); !st.Analyzer.same(run) || st.Orchestrator != nil || st.Running != nil || len(st.Stories) != 0 {
		t.Errorf("state = %+v, want the run as the project's analyzer alone", st)
	}
	j := lab.entries()
	if len(j) != 1 || j[0].Action != hostapi.ActionAnalyze || j[0].Method != "serve.analyze" || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started analyze-agent (command) to analyze, focus risk, as analyzer (pid ") {
		t.Errorf("journal: %+v", j)
	}
	_, err = lab.analyze("", TriggerAsked)
	var no *Refused
	if !errors.As(err, &no) || !strings.Contains(no.Why, "the analyzer is already running for this project (focus risk, started by asked, pid ") || !strings.Contains(no.Why, "started "+run.Started) {
		t.Errorf("a second run: %v", err)
	}
	if st := lab.state(); !st.Analyzer.same(run) || len(lab.entries()) != 1 {
		t.Errorf("the refused run was recorded: %+v %+v", st.Analyzer, lab.entries())
	}
	// a story's agent is not held back by it
	story := lab.ready("Beside")
	lab.l.look(context.Background(), false)
	waitFor(t, "the story's agent started beside the analyzer", func() bool { return lab.run(story) != nil })
	lab.release("analyze")
	lab.release(story)
}

// S-0223: a run asked for no focus looks for all of them: the harness is
// given none, and the run records all.
func TestAnAnalyzerRunWithNoFocusLooksForAll(t *testing.T) {
	lab := analyzeLab(t)
	run, err := lab.analyze("", "schedule daily")
	if err != nil {
		t.Fatal(err)
	}
	if run.Focus != "all" || run.Trigger != "schedule daily" {
		t.Errorf("run = %+v, want focus all and the schedule as its trigger", run)
	}
	if got := lab.analyzerGiven(); !strings.Contains(got, "focus: unset\n") || !strings.Contains(got, "role: analyze\n") {
		t.Errorf("the analyzer was given a focus:\n%s", got)
	}
}

// S-0223: an analyzer run that ends is recorded with its exit, how it went,
// and the report it wrote, the newest file under design/analysis changed
// since it started, its index aside; its activity is logged in the
// analyzer's document from its log, naming the report beside the run's
// seconds and cost, with what started it. A run that wrote none failed, and
// its activity says so.
func TestAnAnalyzerRunThatEndsLogsItsReportAndItsCost(t *testing.T) {
	lab := analyzeLab(t)
	// a report from an earlier run, changed before this one started
	old := filepath.Join(lab.root, "design", "analysis", "2026-10-01-intent.md")
	if err := os.MkdirAll(filepath.Dir(old), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(old, []byte("# Old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hourAgo := time.Now().Add(-time.Hour)
	if err := os.Chtimes(old, hourAgo, hourAgo); err != nil {
		t.Fatal(err)
	}
	report := "design/analysis/2026-10-06-risk.md"
	lab.analyzerRuns("2026-10-06-risk.md", "Wrote "+report+" with two findings.", runStart)
	if !lab.analyzeHere("risk") {
		t.Fatalf("not started: %+v", lab.state().Analyzer)
	}
	waitFor(t, "the analyzer run ended", func() bool { r := lab.state().Analyzer; return r != nil && r.Ended != "" })
	run := lab.state().Analyzer
	if run.Outcome != OutcomeWorked || run.Exit == nil || *run.Exit != 0 || run.Report != report || run.Log == "" {
		t.Errorf("run = %+v, want worked with exit 0 and its report", run)
	}
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityAnalyzer)
		return err == nil && len(doc.Entries) == 1
	})
	doc, _ := lab.repo.Activity(workitem.ActivityAnalyzer)
	if e := doc.Entries[0]; e.Summary != "Wrote "+report+" with two findings." || e.Trigger != "asked" || e.Cost != 0.5 || e.Seconds != 61 || len(e.Items) != 0 {
		t.Errorf("activity = %+v, want the run's summary naming its report, asked, its 0.5 USD, its 61 s, and no item", e)
	}
	if data, _ := os.ReadFile(doc.Path); !strings.Contains(string(data), "\n- Summary: Wrote "+report+" with two findings.\n- Trigger: asked\n- Items: none\n- Seconds: 61\n- Cost: 0.5000 USD") {
		t.Errorf("the entry does not name the report, the seconds, and the cost:\n%s", data)
	}
	// a run whose last word does not name its report has it named after it
	lab.analyzerRuns("2026-10-06-all.md", "Done.", runStart.Add(time.Hour))
	if !lab.analyzeHere("") {
		t.Fatalf("not started: %+v", lab.state().Analyzer)
	}
	waitFor(t, "the second run ended", func() bool { r := lab.state().Analyzer; return r != nil && r.Focus == "all" && r.Ended != "" })
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityAnalyzer)
		return err == nil && len(doc.Entries) == 2
	})
	doc, _ = lab.repo.Activity(workitem.ActivityAnalyzer)
	if e := doc.Entries[1]; e.Summary != "Done. (report design/analysis/2026-10-06-all.md)" || e.Cost != 0.5 {
		t.Errorf("activity = %+v, want the report named after the run's last word", e)
	}
	// a run that wrote no report failed, and says so; the reports before it
	// were changed before it started
	for _, name := range []string{"2026-10-06-risk.md", "2026-10-06-all.md", "README.md"} {
		if err := os.Chtimes(filepath.Join(lab.root, "design", "analysis", name), hourAgo, hourAgo); err != nil {
			t.Fatal(err)
		}
	}
	lab.analyzerRuns("", "Nothing to say.", runStart.Add(2*time.Hour))
	if !lab.analyzeHere("intent") {
		t.Fatalf("not started: %+v", lab.state().Analyzer)
	}
	waitFor(t, "the third run ended", func() bool { r := lab.state().Analyzer; return r != nil && r.Focus == "intent" && r.Ended != "" })
	if r := lab.state().Analyzer; r.Outcome != OutcomeFailed || r.Report != "" || r.Why != "ended (exit 0) without writing a report under design/analysis" {
		t.Errorf("run = %+v, want failed for want of a report", r)
	}
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityAnalyzer)
		return err == nil && len(doc.Entries) == 3
	})
	doc, _ = lab.repo.Activity(workitem.ActivityAnalyzer)
	if e := doc.Entries[2]; e.Summary != "Nothing to say. (no report written under design/analysis)" || e.Cost != 0.5 || e.Seconds != 61 {
		t.Errorf("activity = %+v, want it to say it wrote no report, with its cost and seconds", e)
	}
}

// S-0227: an analyzer run's end names the issues that name the report it
// wrote, and its cost is split evenly between them, under each one's
// analyzer entry; an issue naming another report takes no share. A run whose
// report no issue names names no issue and charges none.
func TestAnAnalyzerRunIsChargedToTheIssuesNamingItsReport(t *testing.T) {
	lab := analyzeLab(t)
	report := "design/analysis/2026-10-06-risk.md"
	var ids []string
	for _, f := range []struct{ title, report string }{
		{"Slow review", report}, {"Old finding", "design/analysis/2026-10-01-intent.md"}, {"Flaky test", report},
	} {
		is, err := issues.New(lab.repo, issues.NewOptions{Title: f.title, Class: "defect", Report: f.report, Now: runStart})
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, is.ID)
	}
	analyzed := func(id string) *usage.Strategic {
		is, err := issues.Get(lab.repo, id)
		if err != nil || is.Usage == nil {
			return nil
		}
		for i, s := range is.Usage.Strategic {
			if s.Kind == workitem.ActivityAnalyzer {
				return &is.Usage.Strategic[i]
			}
		}
		return nil
	}
	lab.analyzerRuns("2026-10-06-risk.md", "Wrote "+report+".", runStart)
	if !lab.analyzeHere("risk") {
		t.Fatalf("not started: %+v", lab.state().Analyzer)
	}
	waitFor(t, "its cost charged", func() bool { return analyzed(ids[2]) != nil })
	doc, _ := lab.repo.Activity(workitem.ActivityAnalyzer)
	if len(doc.Entries) != 1 {
		t.Fatalf("activity = %+v, want one entry", doc.Entries)
	}
	if e := doc.Entries[0]; strings.Join(e.Items, " ") != ids[0]+" "+ids[2] || e.Cost != 0.5 || e.Seconds != 61 {
		t.Errorf("activity = %+v, want it to name %s and %s, with its 0.5 USD and 61 s", e, ids[0], ids[2])
	}
	// 2 input, 300 output, and 2000 cache read tokens, and 61 seconds, in
	// two shares, the remainder to the first named
	for i, w := range map[int]int64{0: 31, 2: 30} {
		if s := analyzed(ids[i]); s == nil || s.Tokens() != 1+150+1000 || s.Seconds != w || s.Cost() != 0.25 {
			t.Errorf("%s's analyzer usage = %+v, want 1151 tokens over %d s, 0.25 USD", ids[i], s, w)
		}
	}
	if s := analyzed(ids[1]); s != nil {
		t.Errorf("%s, naming another report, was charged %+v", ids[1], s)
	}
	// a report no issue names: the entry names none, and no issue is charged
	lab.analyzerRuns("2026-10-06-all.md", "Done.", runStart.Add(time.Hour))
	if !lab.analyzeHere("") {
		t.Fatalf("not started: %+v", lab.state().Analyzer)
	}
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityAnalyzer)
		return err == nil && len(doc.Entries) == 2
	})
	doc, _ = lab.repo.Activity(workitem.ActivityAnalyzer)
	if e := doc.Entries[1]; len(e.Items) != 0 || e.Cost != 0.5 {
		t.Errorf("activity = %+v, want its cost and no item", e)
	}
	if s := analyzed(ids[0]); s == nil || s.Seconds != 31 || s.Cost() != 0.25 {
		t.Errorf("%s's analyzer usage = %+v, want the first run's share alone", ids[0], s)
	}
}

// S-0223: an analyzer run that exits with a failure failed, whatever it
// wrote.
func TestAnAnalyzerRunThatFailsIsRecordedAsFailed(t *testing.T) {
	lab := analyzeLab(t)
	lab.analyzerRuns("2026-10-06-risk.md", "Gave up.", runStart)
	if err := os.WriteFile(filepath.Join(lab.outDir, "exit-analyze"), []byte("3"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !lab.analyzeHere("risk") {
		t.Fatalf("not started: %+v", lab.state().Analyzer)
	}
	waitFor(t, "the analyzer run ended", func() bool { r := lab.state().Analyzer; return r != nil && r.Ended != "" })
	if r := lab.state().Analyzer; r.Outcome != OutcomeFailed || r.Exit == nil || *r.Exit != 3 || r.Why != "ended (exit 3)" || r.Report != "design/analysis/2026-10-06-risk.md" {
		t.Errorf("run = %+v, want failed with exit 3", r)
	}
}

// S-0223: an analyzer run whose process is gone while no launcher waits for
// it, as one flai analyze started is, is settled at the next look, apart
// from the orchestrator's run.
func TestAnAnalyzerRunGoneIsSettledAtALook(t *testing.T) {
	lab := analyzeLab(t)
	gone := exec.Command("true")
	if err := gone.Run(); err != nil {
		t.Fatal(err)
	}
	orchestrator := &AgentRun{Agent: "orchestrator", Started: "2026-10-03T09:00:00Z", Ended: "2026-10-03T09:30:00Z", Outcome: OutcomeWorked}
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(orchestrator)
		s.put(&AgentRun{Agent: "analyzer", Focus: "bottlenecks", PID: gone.Process.Pid, Started: "2026-10-03T10:00:00Z", Session: "s", Trigger: "asked"})
	})
	lab.l.look(context.Background(), false)
	r := lab.state().Analyzer
	if r.Ended == "" || r.Outcome != OutcomeFailed || r.Exit != nil || r.Why != "ended (an exit code nobody saw) without writing a report under design/analysis" {
		t.Errorf("settled = %+v, want failed, with no exit seen and no report", r)
	}
	if o := lab.state().Orchestrator; o == nil || o.Ended != orchestrator.Ended || o.Focus != "" {
		t.Errorf("the orchestrator's run = %+v, want it untouched", o)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if !strings.Contains(logged, `msg="analyzer ended" component=serve project=t pid=`) {
		t.Errorf("the end is not logged:\n%s", logged)
	}
}

// S-0223: the analyzer is refused while the analyze action is off, for a
// focus it does not take, and when nothing can start it, and a refusal is
// not recorded.
func TestTheAnalyzerIsRefusedWithWhy(t *testing.T) {
	lab := analyzeLab(t)
	refusedFor := func(focus, want string) {
		t.Helper()
		_, err := lab.analyze(focus, TriggerAsked)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("analyze %q: %v, want refused for %q", focus, err, want)
		}
	}
	refusedFor("velocity", `the analyzer takes the focus bottlenecks, intent, risk, or none for all of them, and "velocity" is none of them`)
	refusedFor("all", `and "all" is none of them`)
	lab.cfg.Command = nil
	lab.onTheHostsCommand()
	refusedFor("risk", "the analyzer's harness is the host's command, and none is set for this project: set it on the dashboard's Harnesses page, /settings/harnesses, or with flai serve agent set -- <program> [args...]")
	lab.cfg.Analyze = false
	refusedFor("risk", "the analyze host action is off for this project: flai serve enable analyze")
	if st := lab.state(); st.Analyzer != nil || len(lab.entries()) != 0 {
		t.Errorf("a refusal was recorded: %+v %+v", st.Analyzer, lab.entries())
	}
}

// S-0223: the summary of an analyzer run's activity names its report, once.
func TestAnAnalyzerActivityNamesItsReport(t *testing.T) {
	for _, c := range []struct{ final, report, want string }{
		{"Wrote design/analysis/2026-10-06-risk.md.", "design/analysis/2026-10-06-risk.md", "Wrote design/analysis/2026-10-06-risk.md."},
		{"run ended", "design/analysis/2026-10-06-risk.md", "run ended (report design/analysis/2026-10-06-risk.md)"},
		{"run ended", "", "run ended (no report written under design/analysis)"},
	} {
		if got := reportSaid(c.final, c.report, "design/analysis"); got != c.want {
			t.Errorf("reportSaid(%q, %q) = %q, want %q", c.final, c.report, got, c.want)
		}
	}
}

// S-0532: an analyzer run appends one agent event on the analyzer's role at
// its start, saying what started it and its focus, and one at its end, ended
// for a run that wrote its report and failed for one that wrote none, for
// the words its record gives; a start that could not be made is failed, for
// its error.
func TestAnAnalyzerRunAppendsAnAgentEventAtItsStartAndItsEnd(t *testing.T) {
	role := func(lab *agentLab, want ...[2]string) {
		t.Helper()
		lab.agentEventsAre(workitem.ActivityAnalyzer, inbox.RoleType, workitem.ActivityAnalyzer, workitem.ActivityAnalyzer, conventions.RoleAnalyze, want...)
	}
	t.Run("a report is ended", func(t *testing.T) {
		lab := analyzeLab(t)
		lab.analyzerRuns("2026-10-06-risk.md", "Wrote the risk report.", runStart)
		if !lab.analyzeHere("risk") {
			t.Fatalf("not started: %+v", lab.state().Analyzer)
		}
		role(lab, [2]string{inbox.AgentStarted, "the operator asked for it, focus risk"}, [2]string{inbox.AgentEnded, "ended with its work done"})
	})
	t.Run("no report is failed", func(t *testing.T) {
		lab := analyzeLab(t)
		lab.analyzerRuns("", "Found nothing.", runStart)
		if !lab.analyzeHere("") {
			t.Fatalf("not started: %+v", lab.state().Analyzer)
		}
		waitFor(t, "the analyzer run ended", func() bool { r := lab.state().Analyzer; return r != nil && r.Ended != "" })
		r := lab.state().Analyzer
		if r.Outcome != OutcomeFailed || !strings.Contains(r.Why, "(exit 0) without writing a report") {
			t.Fatalf("run = %+v, want failed without a report", r)
		}
		role(lab, [2]string{inbox.AgentStarted, "the operator asked for it, focus " + harness.AllFocus}, [2]string{inbox.AgentFailed, r.Why})
	})
	t.Run("a start that could not be made is failed", func(t *testing.T) {
		lab := analyzeLab(t)
		lab.cfg.Command = []string{filepath.Join(lab.outDir, "no-such-program")}
		if lab.analyzeHere("risk") {
			t.Fatalf("started: %+v", lab.state().Analyzer)
		}
		r := lab.state().Analyzer
		if r == nil || !strings.HasPrefix(r.Why, "could not be started: ") {
			t.Fatalf("not recorded as failed to start: %+v", r)
		}
		role(lab, [2]string{inbox.AgentFailed, r.Why})
	})
}

// instance is a story, Cadence, the lab's repository holds as an instance
// of the formula analysis-cadence that flai formula run made with the focus
// input focus, as its Notes say.
func (lab *agentLab) instance(focus string) string {
	lab.t.Helper()
	id := lab.backlog("Cadence", nil)
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(it.Path)
	notes := "\nMade by `flai formula run analysis-cadence` from `design/workflows/formulas/analysis-cadence.md`, with the inputs `focus=" + focus + "`.\n"
	if err := os.WriteFile(it.Path, append(data, notes...), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	return id
}

// S-0424: the agent of a formula's instance that asks, through RunRole, for
// the analyzer a role step of its formula names starts the analyzer through
// its own start, as the project's analyzer, with the instance's focus, the
// focus input its Notes record, or the run's own focus input when given, and
// the instance as the run's, under the analyze action, with the formula's
// trigger.
func TestAFormulaRoleStepStartsTheAnalyzerOnItsInstance(t *testing.T) {
	t.Run("the instance's focus", func(t *testing.T) {
		lab := analyzeLab(t)
		lab.hold()
		story := lab.instance("risk")
		trigger := FormulaTrigger("analysis-cadence", story)
		run, err := RunRole(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, "analyzer", story, []harness.Input{{Name: "prompt", Value: "Analyze."}}, trigger)
		if err != nil {
			t.Fatal(err)
		}
		if run.Agent != "analyzer" || run.Focus != "risk" || run.Instance != story || run.Item != "" || run.Story != "" || run.RoleName != "" || run.Trigger != trigger || run.PID == 0 {
			t.Errorf("run = %+v, want the analyzer's, focus risk, on %s", run, story)
		}
		if got := lab.analyzerGiven(); !strings.Contains(got, "focus: risk\n") || !strings.Contains(got, "role: analyze\n") || !strings.Contains(got, "agent: analyzer\n") {
			t.Errorf("the analyzer was given\n%s", got)
		}
		if st := lab.state(); !st.Analyzer.same(run) || st.Analyzer.Instance != story || len(st.Roles) != 0 || len(st.Plans) != 0 {
			t.Errorf("state = %+v, want the run as the project's analyzer alone", st)
		}
		if data, err := os.ReadFile(lab.o.Dir.agents()); err != nil || !strings.Contains(string(data), `"instance": "`+story+`"`) {
			t.Errorf("serve/agents.json does not record the run's instance (%v):\n%s", err, data)
		}
		j := lab.entries()
		if len(j) != 1 || j[0].Action != hostapi.ActionAnalyze || j[0].Method != "serve.analyze" || !strings.Contains(j[0].Detail, "to analyze, focus risk, on "+story+", as analyzer (pid ") || !strings.HasSuffix(j[0].Detail, "started by "+trigger) {
			t.Errorf("journal: %+v", j)
		}
		lab.release("analyze")
	})
	t.Run("the run's focus input", func(t *testing.T) {
		lab := analyzeLab(t)
		story := lab.instance("risk")
		run, err := RunRole(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, "analyzer", story, []harness.Input{{Name: "focus", Value: "intent"}}, FormulaTrigger("analysis-cadence", story))
		if err != nil || run.Focus != "intent" || run.Instance != story {
			t.Fatalf("run = %+v (%v), want focus intent on %s", run, err, story)
		}
	})
	t.Run("no focus is all three", func(t *testing.T) {
		lab := analyzeLab(t)
		story := lab.instance("")
		run, err := RunRole(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, "analyzer", story, nil, FormulaTrigger("analysis-cadence", story))
		if err != nil || run.Focus != harness.AllFocus || run.Instance != story {
			t.Fatalf("run = %+v (%v), want focus all on %s", run, err, story)
		}
		if got := lab.analyzerGiven(); !strings.Contains(got, "focus: unset\n") {
			t.Errorf("the analyzer was given a focus:\n%s", got)
		}
	})
}

// S-0424: the analyzer is refused to the operator's own asking through
// RunRole, naming flai analyze; to a role step on an item not its instance;
// with the analyze action off, as Analyze refuses it; and the planner, the
// orchestrator, and the curator are refused to a role step. A refusal records
// nothing.
func TestTheAnalyzerThroughRunRoleIsRefusedWithWhy(t *testing.T) {
	lab := analyzeLab(t)
	story, other := lab.instance("risk"), lab.backlog("Other", nil)
	trigger := FormulaTrigger("analysis-cadence", story)
	refusedFor := func(name, item, trigger, want string) {
		t.Helper()
		_, err := RunRole(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, name, item, nil, trigger)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("run %s on %q by %q: %v, want refused for %q", name, item, trigger, err, want)
		}
	}
	refusedFor("analyzer", story, TriggerAsked, "analyzer is a built-in role, which flai starts with flai analyze, not as a role the project defines")
	refusedFor("analyzer", other, trigger, "the analyzer runs for a role step of a formula on its instance "+story+" alone, and was asked for on "+other+": give the item "+story)
	refusedFor("analyzer", "", trigger, "and was asked for on no item")
	for name, how := range map[string]string{"planner": "flai plan <item>", "orchestrator": "flai serve enable orchestrate", "curator": "flai curate"} {
		refusedFor(name, story, trigger, name+" is a built-in role a role step cannot start: only the analyzer is started for a role step, through flai analyze's start; flai starts the "+name+" with "+how)
	}
	refusedFor("analyzer", story, trigger+"; schedule daily", "analyzer is a built-in role, which flai starts with flai analyze")
	lab.cfg.Analyze = false
	refusedFor("analyzer", story, trigger, "the analyze host action is off for this project: flai serve enable analyze")
	if st := lab.state(); st.Analyzer != nil || len(lab.entries()) != 0 {
		t.Errorf("a refusal was recorded: %+v %+v", st.Analyzer, lab.entries())
	}
}

// S-0424: the analyzer run a formula instance's role step started appends
// its agent events on the instance, by the analyzer, role analyze: started
// with its trigger and focus, then ended for a run that wrote its report and
// failed for one that did not, which the instance's agent is told of; the
// analyzer's other runs keep theirs on the analyzer's role.
func TestAFormulaAnalyzerRunsEndIsAnAgentEventOnTheInstance(t *testing.T) {
	onInstance := func(lab *agentLab, story, title string, want ...[2]string) {
		t.Helper()
		lab.agentEventsAre(story, workitem.Story, title, workitem.ActivityAnalyzer, conventions.RoleAnalyze, want...)
		if got := inbox.AgentEventsOf(lab.repo, workitem.ActivityAnalyzer, time.Time{}); len(got) != 0 {
			t.Errorf("agent events on the analyzer's role: %+v, want none", got)
		}
	}
	start := func(lab *agentLab, story, focus string) {
		t.Helper()
		lab.l.mu.Lock()
		started, _ := lab.l.analyzeOn(context.Background(), lab.cfg, nil, focus, FormulaTrigger("analysis-cadence", story), story)
		lab.l.mu.Unlock()
		if !started {
			t.Fatalf("not started: %+v", lab.state().Analyzer)
		}
		waitFor(t, "the analyzer run ended", func() bool { r := lab.state().Analyzer; return r != nil && r.Ended != "" })
	}
	t.Run("a report is ended", func(t *testing.T) {
		lab := analyzeLab(t)
		story := lab.instance("risk")
		lab.analyzerRuns("2026-11-02-risk.md", "Wrote the risk report.", runStart)
		start(lab, story, "risk")
		onInstance(lab, story, "Cadence", [2]string{inbox.AgentStarted, "started by formula analysis-cadence " + story + ", focus risk"}, [2]string{inbox.AgentEnded, "ended with its work done"})
		items, err := lab.repo.List(true)
		if err != nil {
			t.Fatal(err)
		}
		events, _, err := inbox.CatchUp(lab.repo, inbox.Reader{Agent: "agent-" + story}, story, false, time.Now(), items)
		if err != nil {
			t.Fatal(err)
		}
		var ends []inbox.Event
		for _, e := range events {
			if e.Kind == inbox.Agent && e.ID == story && e.To != inbox.AgentStarted {
				ends = append(ends, e)
			}
		}
		if len(ends) != 1 || ends[0].By != workitem.ActivityAnalyzer || ends[0].To != inbox.AgentEnded {
			t.Errorf("agent-%s is told of %+v, want the analyzer's one end", story, ends)
		}
		if r := lab.state().Analyzer; r.Instance != story || r.Outcome != OutcomeWorked || r.Report != "design/analysis/2026-11-02-risk.md" {
			t.Errorf("run = %+v, want worked on %s with its report", r, story)
		}
	})
	t.Run("no report is failed", func(t *testing.T) {
		lab := analyzeLab(t)
		story := lab.instance("")
		lab.analyzerRuns("", "Found nothing.", runStart)
		start(lab, story, "")
		r := lab.state().Analyzer
		if r.Outcome != OutcomeFailed || !strings.Contains(r.Why, "(exit 0) without writing a report") {
			t.Fatalf("run = %+v, want failed without a report", r)
		}
		onInstance(lab, story, "Cadence", [2]string{inbox.AgentStarted, "started by formula analysis-cadence " + story + ", focus " + harness.AllFocus}, [2]string{inbox.AgentFailed, r.Why})
	})
}

// S-0534: an analyzer run whose log ends on the recorded usage-limit refusal
// is failed on the limit and sets the budget state until the reset the API
// reported; the probe that starts at the reset and ends on the limit again
// extends it, doubling the backoff, the reset being past.
func TestAnAnalyzerRunOnTheUsageLimitSetsTheBudgetStateAndAProbeOnItExtendsIt(t *testing.T) {
	lab := newAgentLab(t)
	at := limitReset.Add(-time.Hour)
	lab.budgetClock(&at)
	one := 1
	run := &AgentRun{Agent: "analyzer", Harness: harness.ClaudeCode, Focus: harness.AllFocus, Trigger: TriggerAsked, Started: at.Format(time.RFC3339), Log: limitResult(t, true)}
	lab.l.analyzeEnded(run, &one)
	lab.heldUntil(limitReset, true, "5m0s")
	if r := lab.state().Analyzer; r.Outcome != OutcomeFailed || !r.UsageLimited || r.ResetsAt != limitResetsAt || r.Why != limitWhyText {
		t.Errorf("analyzer run = %+v; want failed on the limit", r)
	}

	at = limitReset
	probe := &AgentRun{Agent: "analyzer", Harness: harness.ClaudeCode, Focus: harness.AllFocus, Trigger: TriggerAsked, Started: at.Format(time.RFC3339), Log: limitResult(t, true)}
	lab.probing(probe)
	lab.l.analyzeEnded(probe, &one)
	lab.heldUntil(limitReset.Add(2*budgetBackoffFirst), false, "10m0s")
	if got := lab.budgetChanges(); !slices.Equal(got, []string{"set the budget state", "started a probe of the budget state", "extended the budget state"}) {
		t.Errorf("budget changes journalled: %q", got)
	}
}

// S-0534: flai analyze is refused while the budget state of the analyzer's
// provider holds, naming the provider and until, with no run recorded and
// nothing journalled; cleared, it starts.
func TestAnAnalyzeIsRefusedWhileTheBudgetStateHoldsItsProvider(t *testing.T) {
	lab := analyzeLab(t)
	until := anHour()
	lab.budgetHeldOn(harness.Command, until)
	if _, err := lab.analyze("", TriggerAsked); refusal(err) != heldFor(until) {
		t.Errorf("analyze: %v; want refused with %q", err, heldFor(until))
	}
	if r := lab.state().Analyzer; r != nil || len(lab.entries()) != 0 {
		t.Fatalf("a held start recorded %+v, journalled %+v", r, lab.entries())
	}
	lab.budgetCleared()
	if _, err := lab.analyze("", TriggerAsked); err != nil {
		t.Errorf("once the state cleared: %v", err)
	}
}

// S-0509, ADR-0247: an analyzer whose agent names the container harness is
// started in a container of its role, mounting the main checkout, and pilot's
// start request names its focus.
func TestAnAnalyzerOnTheContainerHarnessIsStartedInAContainerOfItsRole(t *testing.T) {
	lab, _ := roleContainerLab(t)
	lab.l.mu.Lock()
	started, _ := lab.l.analyze(context.Background(), lab.cfg, &manifest.Agent{Harness: harness.Container}, "risk", TriggerAsked)
	lab.l.mu.Unlock()
	if !started {
		t.Fatalf("not started: %+v", lab.state().Analyzer)
	}
	if req := lab.roleContainerStarted(conventions.RoleAnalyze, lab.state().Analyzer); req.Focus != "risk" {
		t.Errorf("the start request's focus is %q, not risk", req.Focus)
	}
}
