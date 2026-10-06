//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
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
	return lab.l.analyze(context.Background(), lab.cfg, nil, focus, TriggerAsked)
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
	refusedFor("risk", "the analyzer's agent names no harness (analysis.agent or agent in system-flow.yaml)")
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
