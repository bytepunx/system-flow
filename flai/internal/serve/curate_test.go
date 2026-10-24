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

	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// curateLab is an agentLab whose command is a curator stub: it writes what it
// was given to curate.txt, prints stream-curate as its output when there is
// one, copies the entry files staged under digest-stage into design/digest,
// as flai digest new and edit would have written them, waits while held
// until released, and exits with the code in exit-curate when there is one.
// The curate action is on.
func curateLab(t *testing.T) *agentLab {
	t.Helper()
	lab := newAgentLab(t)
	stub := filepath.Join(lab.outDir, "curate-agent")
	d := lab.outDir
	script := "#!/bin/sh\n" +
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: ${FLAI_STORY-unset}\"; echo \"role: $FLAI_ROLE\"; echo \"item: ${FLAI_ITEM-unset}\"; echo \"since: ${FLAI_SINCE-unset}\"; echo \"session: $FLAI_SESSION\"; echo \"by: $FLAI_STARTED_BY\"; } > \"" + d + "/curate.txt\"\n" +
		"[ -f \"" + d + "/stream-curate\" ] && cat \"" + d + "/stream-curate\"\n" +
		"if [ -d \"" + d + "/digest-stage\" ]; then mkdir -p design/digest; cp \"" + d + "\"/digest-stage/*.md design/digest/; fi\n" +
		"while [ -f \"" + d + "/hold\" ] && [ ! -f \"" + d + "/release-curate\" ]; do sleep 0.05; done\n" +
		"touch \"" + d + "/ended-$$\"\n" +
		"[ -f \"" + d + "/exit-curate\" ] && exit \"$(cat \"" + d + "/exit-curate\")\"\n" +
		"exit 0\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg.Command, lab.cfg.Curate = []string{stub, "story={story}"}, true
	return lab
}

func (lab *agentLab) curate(since string) (*AgentRun, error) {
	return Curate(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, since, TriggerAsked)
}

// curateHere has the lab's launcher start the curator, as the operator
// asked, and wait for it, as a serving flai would.
func (lab *agentLab) curateHere(since string) bool {
	lab.l.mu.Lock()
	defer lab.l.mu.Unlock()
	return lab.l.curate(context.Background(), lab.cfg, nil, since, TriggerAsked)
}

// curatorGiven is what the curator stub was given, once it has written it.
func (lab *agentLab) curatorGiven() string {
	lab.t.Helper()
	out := filepath.Join(lab.outDir, "curate.txt")
	waitFor(lab.t, "the curator stub ran", func() bool {
		data, _ := os.ReadFile(out)
		return strings.Contains(string(data), "by: ")
	})
	data, _ := os.ReadFile(out)
	return string(data)
}

// curatorRuns has the stub print a stream of two calls from at, a minute
// apart, and final as its last word.
func (lab *agentLab) curatorRuns(final string, at time.Time) {
	lab.t.Helper()
	session := "s-" + at.Format("150405")
	end := strings.Replace(streamFinal(final, 2000, 0.5), `"session_id":"s"`, `"session_id":"`+session+`"`, 1)
	stream := strings.Join([]string{streamCall(session, session+"-m1", at, 999), streamCall(session, session+"-m2", at.Add(time.Minute), 999), end}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-curate"), []byte(stream), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// curatorActivity is the curator's activity document once it holds n entries.
func (lab *agentLab) curatorActivity(n int) *workitem.Activity {
	lab.t.Helper()
	waitFor(lab.t, "the curator's activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityCurator)
		return err == nil && len(doc.Entries) == n
	})
	doc, _ := lab.repo.Activity(workitem.ActivityCurator)
	return doc
}

// S-0431: the curator is started for the project in its main checkout, as
// the curator with FLAI_ROLE=curate and no story or item, logged as the
// curator's run, and recorded as the project's curator, marked with its role,
// apart from the orchestrator's, the analyzer's, and the stories' runs, with
// its since and what started it; a second run is refused while the first
// runs, naming it; once its process is gone, the next look settles it.
func TestTheCuratorIsStartedForTheProject(t *testing.T) {
	t.Setenv("FLAI_ITEM", "E-0999") // the session that ran flai curate
	lab := curateLab(t)
	lab.hold()
	run, err := lab.curate("2026-10-01")
	if err != nil {
		t.Fatal(err)
	}
	if run.Agent != "curator" || run.Role != "curate" || run.Since != "2026-10-01" || run.Trigger != "asked" || run.Story != "" || run.Item != "" || run.Focus != "" || run.PID == 0 || run.Session == "" || run.Harness != "command" || run.Command != "curate-agent" {
		t.Errorf("run = %+v", run)
	}
	if data, err := os.ReadFile(lab.o.Dir.agents()); err != nil || !strings.Contains(string(data), `"curator": {`) || !strings.Contains(string(data), `"role": "curate"`) || !strings.Contains(string(data), `"since": "2026-10-01"`) {
		t.Errorf("serve/agents.json does not record the curator's run (%v):\n%s", err, data)
	}
	logs, err := lab.o.Dir.ActivityLogs("t", workitem.ActivityCurator)
	if err != nil || len(logs) != 1 || logs[0] != run.Log || !strings.HasPrefix(filepath.Base(run.Log), "t-curator-") {
		t.Errorf("curator logs = %v (%v), want the run's %s", logs, err, run.Log)
	}
	got := lab.curatorGiven()
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{"args: story=\n", "dir: " + real + "\n", "agent: curator\n", "story: unset\n", "role: curate\n", "item: unset\n", "since: 2026-10-01\n", "session: 2", "by: flai-serve\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the curator was given\n%s\nwithout %q", got, want)
		}
	}
	if st := lab.state(); !st.Curator.same(run) || st.Orchestrator != nil || st.Analyzer != nil || st.Running != nil || len(st.Stories) != 0 {
		t.Errorf("state = %+v, want the run as the project's curator alone", st)
	}
	j := lab.entries()
	if len(j) != 1 || j[0].Action != "curate" || j[0].Method != "serve.curate" || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started curate-agent (command) to curate, since 2026-10-01, as curator (pid ") {
		t.Errorf("journal: %+v", j)
	}
	_, err = lab.curate("")
	var no *Refused
	if !errors.As(err, &no) || !strings.Contains(no.Why, "the curator is already running for this project (since 2026-10-01, started by asked, pid ") || !strings.Contains(no.Why, "started "+run.Started) || !strings.Contains(no.Why, run.Log) {
		t.Errorf("a second run: %v", err)
	}
	if st := lab.state(); !st.Curator.same(run) || len(lab.entries()) != 1 {
		t.Errorf("the refused run was recorded: %+v %+v", st.Curator, lab.entries())
	}
	// a story's agent is not held back by it
	story := lab.ready("Beside")
	lab.l.look(context.Background(), false)
	waitFor(t, "the story's agent started beside the curator", func() bool { return lab.run(story) != nil })
	lab.release("curate")
	lab.release(story)
}

// curatorSettled records the project's curator run as ended, as the serving
// flai does once a handed-over run's process is gone.
func (lab *agentLab) curatorSettled() {
	lab.t.Helper()
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		ended := *s.Curator
		ended.Ended, ended.Outcome = time.Now().UTC().Format(time.RFC3339), OutcomeWorked
		s.put(&ended)
	})
}

// S-0431: a curator run whose process is gone while no launcher waits for
// it, as one flai curate started is, is settled at the next look, worked,
// apart from the orchestrator's run.
func TestACuratorRunGoneIsSettledAtALook(t *testing.T) {
	lab := curateLab(t)
	gone := exec.Command("true")
	if err := gone.Run(); err != nil {
		t.Fatal(err)
	}
	orchestrator := &AgentRun{Agent: "orchestrator", Started: "2026-10-03T09:00:00Z", Ended: "2026-10-03T09:30:00Z", Outcome: OutcomeWorked}
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(orchestrator)
		s.put(&AgentRun{Agent: "curator", Role: "curate", Since: "2026-10-01", PID: gone.Process.Pid, Started: "2026-10-03T10:00:00Z", Session: "s", Trigger: "asked"})
	})
	lab.l.look(context.Background(), false)
	if r := lab.state().Curator; r.Ended == "" || r.Outcome != OutcomeWorked || r.Exit != nil || r.Why != "" || r.Since != "2026-10-01" {
		t.Errorf("settled = %+v, want worked, with no exit seen", r)
	}
	if o := lab.state().Orchestrator; o == nil || o.Ended != orchestrator.Ended || o.Role != "" || o.Agent != "orchestrator" {
		t.Errorf("the orchestrator's run = %+v, want it untouched", o)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if !strings.Contains(logged, `msg="curator ended" component=serve project=t pid=`) {
		t.Errorf("the end is not logged:\n%s", logged)
	}
}

// S-0431: a run given no since takes in from the time of the curator's last
// activity, and from the beginning on its first run, when it is given none.
func TestTheCuratorsSinceIsItsLastActivity(t *testing.T) {
	lab := curateLab(t)
	run, err := lab.curate("")
	if err != nil {
		t.Fatal(err)
	}
	if run.Since != "" {
		t.Errorf("a first run's since = %q, want none", run.Since)
	}
	if got := lab.curatorGiven(); !strings.Contains(got, "since: unset\n") || !strings.Contains(lab.entries()[0].Detail, "to curate, everything, its first run, as curator") {
		t.Errorf("a first run was given a since:\n%s\n%+v", got, lab.entries())
	}
	lab.curatorSettled()
	_ = os.Remove(filepath.Join(lab.outDir, "curate.txt"))
	for _, at := range []time.Time{runStart.Add(-time.Hour), runStart.Add(2 * time.Hour)} {
		if _, err := lab.repo.AppendActivity(workitem.ActivityCurator, workitem.ActivityEntry{At: at, Summary: "Took in two entries."}); err != nil {
			t.Fatal(err)
		}
	}
	run, err = lab.curate("")
	if err != nil {
		t.Fatal(err)
	}
	want := runStart.Add(2 * time.Hour).Format(time.RFC3339)
	if run.Since != want {
		t.Errorf("since = %q, want the last activity's %s", run.Since, want)
	}
	if got := lab.curatorGiven(); !strings.Contains(got, "since: "+want+"\n") {
		t.Errorf("the curator was given\n%s\nwithout since %s", got, want)
	}
	lab.curatorSettled()
	// a since given wins over the last activity, an RFC 3339 time as well as a day
	run, err = lab.curate("2026-09-30T08:00:00+02:00")
	if err != nil || run.Since != "2026-09-30T08:00:00+02:00" {
		t.Errorf("run = %+v (%v), want the since given", run, err)
	}
}

// S-0431: a curator run on a project with archived stories and threads that
// ends cleanly worked, and records and logs as its items the digest entries
// whose files it wrote or changed since it started, those it left aside
// unnamed; the entries it wrote pass flai check's digest rules. Its activity
// names the run's last word, what started it, its seconds and cost, and
// charges no item. A run that wrote no entry worked all the same.
func TestACuratorRunThatEndsLogsTheEntriesItWrote(t *testing.T) {
	lab := curateLab(t)
	// a story done and archived, and a thread the operator answered
	done := lab.backlog("Settled", nil)
	it, err := lab.repo.Get(done)
	if err != nil {
		t.Fatal(err)
	}
	it.Status = workitem.Done
	if err := lab.repo.Save(it); err != nil {
		t.Fatal(err)
	}
	items, _ := lab.repo.List(false)
	plan, err := lab.repo.PlanArchive(items, []string{done})
	if err != nil {
		t.Fatal(err)
	}
	if err := lab.repo.Archive(plan); err != nil {
		t.Fatal(err)
	}
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: lab.epic.ID, Author: "builder-" + done, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	env := digest.EnvFor(lab.repo)
	write := func(title string, sources ...string) digest.Result {
		t.Helper()
		res, err := digest.New(lab.repo, execx.System{}, digest.NewOptions{Title: title, Kind: "decision", Sources: sources, Topics: []string{"cli"}, Body: "Decided so.", Now: runStart, Env: env})
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	folder := filepath.Join(lab.root, "design", "digest")
	stage := filepath.Join(lab.outDir, "digest-stage")
	if err := os.MkdirAll(stage, 0o755); err != nil {
		t.Fatal(err)
	}
	hourAgo := time.Now().Add(-time.Hour)
	// two entries from before the run: the first the run merges an answer
	// into, the second it leaves alone
	merged, untouched := write("The port is eight.", done), write("Logs go to serve/agents.", done)
	before, err := os.ReadFile(filepath.Join(lab.root, merged.Path))
	if err != nil {
		t.Fatal(err)
	}
	asked := []string{th.ID}
	if _, err := digest.Edit(lab.repo, merged.ID, digest.EditOptions{AskedAgain: &asked, Env: env}); err != nil {
		t.Fatal(err)
	}
	stageFile := func(rel string) {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(lab.root, rel))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(stage, filepath.Base(rel)), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	stageFile(merged.Path)
	if err := os.WriteFile(filepath.Join(lab.root, merged.Path), before, 0o644); err != nil {
		t.Fatal(err)
	}
	// two the run writes, moved out of the folder until the stub copies them in
	written := []digest.Result{write("The curator writes through flai digest.", done, th.ID), write("Threads answered once are not asked again.", done, th.ID)}
	for _, res := range written {
		stageFile(res.Path)
		if err := os.Remove(filepath.Join(lab.root, res.Path)); err != nil {
			t.Fatal(err)
		}
	}
	for _, rel := range []string{merged.Path, untouched.Path} {
		if err := os.Chtimes(filepath.Join(lab.root, rel), hourAgo, hourAgo); err != nil {
			t.Fatal(err)
		}
	}
	final := "Wrote D-0003 and D-0004, and merged TH into D-0001."
	lab.curatorRuns(final, runStart)
	if !lab.curateHere("2026-10-01") {
		t.Fatalf("not started: %+v", lab.state().Curator)
	}
	waitFor(t, "the curator run ended", func() bool { r := lab.state().Curator; return r != nil && r.Ended != "" })
	run := lab.state().Curator
	if run.Outcome != OutcomeWorked || run.Exit == nil || *run.Exit != 0 || run.Why != "" || strings.Join(run.Entries, " ") != "D-0001 D-0003 D-0004" || run.Since != "2026-10-01" {
		t.Errorf("run = %+v, want worked with exit 0 and the entries D-0001, D-0003, and D-0004", run)
	}
	doc := lab.curatorActivity(1)
	if e := doc.Entries[0]; e.Summary != final || e.Trigger != "asked" || strings.Join(e.Items, " ") != "D-0001 D-0003 D-0004" || e.Cost != 0.5 || e.Seconds != 61 {
		t.Errorf("activity = %+v, want the run's last word, asked, the three entries, its 0.5 USD, and its 61 s", e)
	}
	if data, _ := os.ReadFile(doc.Path); !strings.Contains(string(data), "\n- Summary: "+final+"\n- Trigger: asked\n- Items: D-0001, D-0003, D-0004\n- Seconds: 61\n- Cost: 0.5000 USD") {
		t.Errorf("the entry does not name the entries, the seconds, and the cost:\n%s", data)
	}
	// the entries pass the check flai check makes of the digest
	files, err := os.ReadDir(folder)
	if err != nil || len(files) != 4 {
		t.Fatalf("the digest folder holds %d files (%v), want four", len(files), err)
	}
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(folder, f.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if _, problems := digest.Validate(f.Name(), string(data), digest.EnvFor(lab.repo)); len(problems) != 0 {
			t.Errorf("%s: %v", f.Name(), problems)
		}
	}
	if e, err := digest.Get(lab.repo, "D-0001"); err != nil || strings.Join(e.AskedAgain, " ") != th.ID {
		t.Errorf("the merged entry = %+v (%v), want %s asked again", e, err, th.ID)
	}
	// no item is charged: the archived story's usage is as it was
	if got, err := lab.repo.Get(done); err != nil || got.Usage != nil {
		t.Errorf("%s = %+v (%v), want no usage charged", done, got, err)
	}
	// a run that found nothing to take in worked, naming no entry
	if err := os.RemoveAll(stage); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if err := os.Chtimes(filepath.Join(folder, f.Name()), hourAgo, hourAgo); err != nil {
			t.Fatal(err)
		}
	}
	lab.curatorRuns("Nothing new to take in.", runStart.Add(time.Hour))
	if !lab.curateHere("") {
		t.Fatalf("not started: %+v", lab.state().Curator)
	}
	waitFor(t, "the second run ended", func() bool { r := lab.state().Curator; return r != nil && r.Since != "2026-10-01" && r.Ended != "" })
	if r := lab.state().Curator; r.Outcome != OutcomeWorked || len(r.Entries) != 0 || r.Why != "" {
		t.Errorf("run = %+v, want worked with no entry", r)
	}
	doc = lab.curatorActivity(2)
	if e := doc.Entries[1]; e.Summary != "Nothing new to take in." || len(e.Items) != 0 || e.Cost != 0.5 {
		t.Errorf("activity = %+v, want its last word, no item, and its cost", e)
	}
}

// S-0431: a curator run that exits with a failure failed, with why, and
// still names the entries it changed.
func TestACuratorRunThatFailsIsRecordedAsFailed(t *testing.T) {
	lab := curateLab(t)
	res, err := digest.New(lab.repo, execx.System{}, digest.NewOptions{Title: "Half done.", Kind: "decision", Sources: []string{lab.epic.ID}, Body: "Begun.", Now: runStart})
	if err != nil {
		t.Fatal(err)
	}
	stage := filepath.Join(lab.outDir, "digest-stage")
	if err := os.MkdirAll(stage, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(lab.root, res.Path), filepath.Join(stage, filepath.Base(res.Path))); err != nil {
		t.Fatal(err)
	}
	lab.curatorRuns("Gave up.", runStart)
	if err := os.WriteFile(filepath.Join(lab.outDir, "exit-curate"), []byte("3"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !lab.curateHere("") {
		t.Fatalf("not started: %+v", lab.state().Curator)
	}
	waitFor(t, "the curator run ended", func() bool { r := lab.state().Curator; return r != nil && r.Ended != "" })
	if r := lab.state().Curator; r.Outcome != OutcomeFailed || r.Exit == nil || *r.Exit != 3 || r.Why != "ended (exit 3)" || strings.Join(r.Entries, " ") != res.ID {
		t.Errorf("run = %+v, want failed with exit 3, naming %s", r, res.ID)
	}
	if e := lab.curatorActivity(1).Entries[0]; e.Summary != "Gave up." || strings.Join(e.Items, " ") != res.ID {
		t.Errorf("activity = %+v, want the run's last word and %s", e, res.ID)
	}
}

// S-0431: the curator is refused while the curate action is off, for a since
// in neither form, and when nothing can start it, and a refusal is not
// recorded.
func TestTheCuratorIsRefusedWithWhy(t *testing.T) {
	lab := curateLab(t)
	refusedFor := func(since, want string) {
		t.Helper()
		_, err := lab.curate(since)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("curate %q: %v, want refused for %q", since, err, want)
		}
	}
	for _, since := range []string{"yesterday", "2026-10-1", "2026-10-01 09:00", "2026-13-01"} {
		refusedFor(since, `the curator takes a since as a YYYY-MM-DD day in UTC or an RFC 3339 time, such as 2026-10-01 or 2026-10-01T09:00:00Z, and "`+since+`" is neither`)
	}
	lab.cfg.Command = nil
	refusedFor("", "the curator's agent names no harness (digest.agent or agent in system-flow.yaml), and no command is set on the host")
	lab.cfg.Curate = false
	refusedFor("2026-10-01", "the curate host action is off for this project: flai serve enable curate")
	if st := lab.state(); st.Curator != nil || len(lab.entries()) != 0 {
		t.Errorf("a refusal was recorded: %+v %+v", st.Curator, lab.entries())
	}
}

// S-0431, ADR-0196: a curator run, which names no story, item, or focus, is
// told from the orchestrator's by its role: its start and its end leave the
// orchestrator's record, and the analyzer's, as they were.
func TestACuratorRunLeavesTheOrchestratorsRecordAlone(t *testing.T) {
	lab := curateLab(t)
	orchestrator := &AgentRun{Agent: "orchestrator", Started: "2026-10-03T09:00:00Z", Ended: "2026-10-03T09:30:00Z", Outcome: OutcomeWorked, Session: "o"}
	analyzer := &AgentRun{Agent: "analyzer", Focus: "risk", Started: "2026-10-03T09:10:00Z", Ended: "2026-10-03T09:20:00Z", Outcome: OutcomeWorked}
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(orchestrator)
		s.put(analyzer)
	})
	lab.curatorRuns("Done.", runStart)
	if !lab.curateHere("") {
		t.Fatalf("not started: %+v", lab.state().Curator)
	}
	waitFor(t, "the curator run ended", func() bool { r := lab.state().Curator; return r != nil && r.Ended != "" })
	st := lab.state()
	if o := st.Orchestrator; o == nil || o.Agent != "orchestrator" || o.Started != orchestrator.Started || o.Ended != orchestrator.Ended || o.Role != "" || len(st.PastOrchestrators) != 0 {
		t.Errorf("the orchestrator's run = %+v, past %v, want it untouched", o, st.PastOrchestrators)
	}
	if a := st.Analyzer; a == nil || a.Started != analyzer.Started || len(st.PastAnalyzers) != 0 {
		t.Errorf("the analyzer's run = %+v, want it untouched", a)
	}
	if c := st.Curator; c.Agent != "curator" || c.Role != "curate" || !c.curates() || c.orchestrates() || c.analyzes() {
		t.Errorf("the curator's run = %+v, want it marked as the curator's alone", c)
	}
	if st.Running != nil || len(st.Stories) != 0 || len(st.Plans) != 0 {
		t.Errorf("state = %+v, want the curator's run held apart from the stories' and the plans'", st)
	}
}
