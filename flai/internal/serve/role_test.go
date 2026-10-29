//go:build !windows

package serve

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// reviewerRole is a role a project defines (ADR-0207): reviewer, in role
// review, behind the review host action, on a story, with a baseline and the
// project's additions.
const reviewerRole = `---
title: Reviewer
updated: 2026-10-28
status: active
name: reviewer
role: review
definition:
  description: Reviews a story against the conventions, for flai serve.
  tools: [Read, Grep, mcp__flai__item_get]
  model: inherit
action: review
trigger:
  - asked
item: [story]
topics: [planning]
works_story: false
charge: none
tools: []
commands: []
files: []
permissions: []
---

You are the reviewer: read the story and say what it lacks.

<!-- system-flow:end-of-baseline -->

## Project additions

Review its touches too.
`

// auditorRole is a role a project defines that runs on no item: auditor, in
// role audit, behind the audit host action.
const auditorRole = `---
title: Auditor
updated: 2026-10-28
status: active
name: auditor
role: audit
definition:
  description: Audits the project, for flai serve.
  tools: [Read, Grep]
  model: inherit
action: audit
trigger:
  - asked
item: none
topics: [planning]
works_story: false
charge: none
tools: []
commands: []
files: []
permissions: []
---

You are the auditor: read the project and say what it lacks.

<!-- system-flow:end-of-baseline -->
`

// roleLab is an agentLab whose project defines reviewer and auditor in
// their role files, with their host actions on, and whose command is a role
// stub: it writes what it was given, its arguments last, to
// role-<FLAI_AGENT>.txt, prints stream-<FLAI_AGENT> as its output when there
// is one, waits while held until released, and exits with the code in
// exit-<FLAI_AGENT> when there is one.
func roleLab(t *testing.T) *agentLab {
	t.Helper()
	lab := newAgentLab(t)
	dir := filepath.Join(lab.root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"reviewer.md": reviewerRole, "auditor.md": auditorRole} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	stub := filepath.Join(lab.outDir, "role-agent")
	d := lab.outDir
	script := "#!/bin/sh\n" +
		"{ echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"role: $FLAI_ROLE\"; echo \"item: ${FLAI_ITEM-unset}\"; echo \"inputs: ${FLAI_INPUTS-unset}\"; echo \"story: ${FLAI_STORY-unset}\"; echo \"session: $FLAI_SESSION\"; echo \"by: $FLAI_STARTED_BY\"; echo \"args:\"; printf '%s\\n' \"$@\"; echo \"end of args\"; } > \"" + d + "/role-$FLAI_AGENT.txt\"\n" +
		"[ -f \"" + d + "/stream-$FLAI_AGENT\" ] && cat \"" + d + "/stream-$FLAI_AGENT\"\n" +
		"while [ -f \"" + d + "/hold\" ] && [ ! -f \"" + d + "/release-$FLAI_AGENT\" ]; do sleep 0.05; done\n" +
		"touch \"" + d + "/ended-$$\"\n" +
		"[ -f \"" + d + "/exit-$FLAI_AGENT\" ] && exit \"$(cat \"" + d + "/exit-$FLAI_AGENT\")\"\n" +
		"exit 0\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg.Command = []string{stub}
	lab.cfg.RoleActions = map[string]bool{"review": true, "audit": true}
	return lab
}

func (lab *agentLab) runRole(name, item string, inputs ...harness.Input) (*AgentRun, error) {
	return RunRole(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, name, item, inputs)
}

// roleHere has the lab's launcher start the role name on item, as the
// operator asked, and wait for it, as a serving flai would.
func (lab *agentLab) roleHere(name, item string) bool {
	lab.t.Helper()
	repo, err := workitem.Open(lab.root)
	if err != nil {
		lab.t.Fatal(err)
	}
	f, err := definedRole(repo, name)
	if err != nil {
		lab.t.Fatal(err)
	}
	lab.l.mu.Lock()
	defer lab.l.mu.Unlock()
	return lab.l.role(context.Background(), lab.cfg, f, item, repo.Manifest.RoleAgent(name), nil)
}

func (lab *agentLab) roleRun(name, item string) *AgentRun { return lab.state().Roles[name][item] }

// roleGiven is what the role stub was given as agent, once it has written
// it.
func (lab *agentLab) roleGiven(agent string) string {
	lab.t.Helper()
	out := filepath.Join(lab.outDir, "role-"+agent+".txt")
	waitFor(lab.t, "the role stub ran as "+agent, func() bool {
		data, _ := os.ReadFile(out)
		return strings.Contains(string(data), "end of args\n")
	})
	data, _ := os.ReadFile(out)
	return string(data)
}

// roleLog is the name of a defined role's run's log: the project's key, the
// role's name, and the second it started.
var roleLog = regexp.MustCompile(`^t-reviewer-\d{8}T\d{6}Z\.log$`)

// S-0417: a role the project defines is started from its role file, on its
// item, in the project's main checkout, as its verb, under its name and its
// item, with its inputs; it is recorded in serve/agents.json under roles, by
// its name and item, with a log of its own under its name, and journalled
// under its action; a second run on the item is refused while the first
// runs.
func TestADefinedRoleIsStartedFromItsFile(t *testing.T) {
	t.Setenv("FLAI_STORY", "S-0999") // the session that ran flai role run
	lab := roleLab(t)
	lab.hold()
	story := lab.backlog("Reviewed", nil)
	run, err := lab.runRole("reviewer", story, harness.Input{Name: "depth", Value: "thorough"})
	if err != nil {
		t.Fatal(err)
	}
	agent := "reviewer-" + story
	if run.Item != story || run.Story != "" || run.RoleName != "reviewer" || run.Role != "review" || run.Action != "review" || run.Agent != agent ||
		run.PID == 0 || run.Session == "" || run.Harness != "command" || run.Command != "role-agent" || run.Trigger != TriggerAsked {
		t.Errorf("run = %+v", run)
	}
	if !roleLog.MatchString(filepath.Base(run.Log)) {
		t.Errorf("log = %s, want <key>-<name>-<start>.log", run.Log)
	}
	got := lab.roleGiven(agent)
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{"dir: " + real + "\n", "agent: " + agent + "\n", "role: review\n", "item: " + story + "\n",
		`inputs: [{"name":"depth","value":"thorough"}]` + "\n", "story: unset\n", "session: 2", "by: flai-serve\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the role was given\n%s\nwithout %q", got, want)
		}
	}
	data, err := os.ReadFile(lab.o.Dir.agents())
	if err != nil {
		t.Fatal(err)
	}
	var recorded map[string]struct {
		Roles map[string]map[string]struct {
			RoleName string `json:"role_name"`
			Agent    string `json:"agent"`
			Log      string `json:"log"`
		} `json:"roles"`
	}
	if err := json.Unmarshal(data, &recorded); err != nil {
		t.Fatal(err)
	}
	if r := recorded[lab.root].Roles["reviewer"][story]; r.RoleName != "reviewer" || r.Agent != agent || r.Log != run.Log {
		t.Errorf("serve/agents.json does not record the run under roles, by name and item:\n%s", data)
	}
	if st := lab.state(); st.Stories[story] != nil || st.Running != nil || len(st.Plans) != 0 || st.Orchestrator != nil || st.Curator != nil {
		t.Errorf("state = %+v, want the run under roles alone", st)
	}
	j := lab.entries()
	if len(j) != 1 || j[0].Action != "review" || j[0].Method != MethodRole || j[0].Outcome != "done" ||
		!strings.Contains(j[0].Detail, "started role-agent (command) to run the role reviewer on "+story+" as "+agent+" (pid ") {
		t.Errorf("journal: %+v", j)
	}
	_, err = lab.runRole("reviewer", story)
	var no *Refused
	if !errors.As(err, &no) || !strings.Contains(no.Why, "the role reviewer is already running on "+story) || !strings.Contains(no.Why, "started "+run.Started) {
		t.Errorf("a second run: %v", err)
	}
	// another item's run starts beside it, with a log of its own
	other := lab.backlog("Also reviewed", nil)
	second, err := lab.runRole("reviewer", other)
	if err != nil {
		t.Fatalf("another item: %v", err)
	}
	if second.Log == run.Log || !strings.Contains(lab.roleGiven("reviewer-"+other), "item: "+other+"\n") {
		t.Errorf("second run = %+v, want a log of its own", second)
	}
	lab.release(agent)
	lab.release("reviewer-" + other)
}

// S-0417: a role that runs on no item is started under its name alone, with
// no FLAI_ITEM, and recorded under its name and the empty item; a second run
// is refused while it runs.
func TestADefinedRoleOnNoItemRunsUnderItsName(t *testing.T) {
	lab := roleLab(t)
	lab.hold()
	run, err := lab.runRole("auditor", "")
	if err != nil {
		t.Fatal(err)
	}
	if run.Item != "" || run.Agent != "auditor" || run.RoleName != "auditor" || run.Role != "audit" {
		t.Errorf("run = %+v", run)
	}
	got := lab.roleGiven("auditor")
	for _, want := range []string{"agent: auditor\n", "role: audit\n", "item: unset\n", "inputs: unset\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the role was given\n%s\nwithout %q", got, want)
		}
	}
	if r := lab.roleRun("auditor", ""); !r.same(run) {
		t.Errorf("recorded = %+v, want the run under auditor and no item", r)
	}
	if j := lab.entries(); len(j) != 1 || j[0].Action != "audit" || !strings.Contains(j[0].Detail, "to run the role auditor on no item as auditor (pid ") {
		t.Errorf("journal: %+v", j)
	}
	_, err = lab.runRole("auditor", "")
	var no *Refused
	if !errors.As(err, &no) || !strings.Contains(no.Why, "the role auditor is already running on no item") {
		t.Errorf("a second run: %v", err)
	}
	lab.release("auditor")
}

// S-0417, ADR-0207 §4: on claude-code, the run's prompt is its role file's
// body with its item and inputs, and its session runs as the definition flai
// generates from the file; its agent is roles.<name>.agent over the
// project's.
func TestADefinedRolesPromptIsItsFilesBody(t *testing.T) {
	lab := roleLab(t)
	m, _ := os.ReadFile(filepath.Join(lab.root, manifest.File))
	if err := os.WriteFile(filepath.Join(lab.root, manifest.File), append(m, "agent:\n  model: opus\nroles:\n  reviewer:\n    agent:\n      harness: claude-code\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.cfg.Command = nil
	lab.cfg.Flai = "/usr/local/bin/flai"
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: filepath.Join(lab.outDir, "role-agent")}}
	story := lab.backlog("Reviewed", nil)
	run, err := lab.runRole("reviewer", story, harness.Input{Name: "depth", Value: "thorough"})
	if err != nil {
		t.Fatal(err)
	}
	if run.Harness != harness.ClaudeCode || run.Model != "opus" {
		t.Errorf("run = %+v, want claude-code with the project's model", run)
	}
	got := lab.roleGiven("reviewer-" + story)
	for _, want := range []string{
		"You are reviewer-" + story + ", the reviewer role this project defines in design/workflows/roles/reviewer.md",
		"You are the reviewer: read the story and say what it lacks.",
		"Review its touches too.",
		"This run is on " + story,
		"- depth: thorough",
		"--agent\nreviewer\n",
		"role: review\n", "item: " + story + "\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the role was given\n%s\nwithout %q", got, want)
		}
	}
	if strings.Contains(got, "end-of-baseline") {
		t.Errorf("the prompt holds the marker:\n%s", got)
	}
}

// S-0417: a defined role's run that ends is recorded with its exit and how
// it went, journalled under its action, and its activity logged in its own
// activity document, naming its item and what started it, charging nothing.
func TestADefinedRoleRunThatEndsIsRecordedAndItsActivityLogged(t *testing.T) {
	lab := roleLab(t)
	story := lab.backlog("Reviewed", nil)
	agent := "reviewer-" + story
	stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamCall("s", "m2", runStart.Add(time.Minute), 999), streamFinal("Reviewed "+story+": two criteria lack a test.", 2000, 0.5)}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-"+agent), []byte(stream), 0o644); err != nil {
		t.Fatal(err)
	}
	if !lab.roleHere("reviewer", story) {
		t.Fatalf("not started: %+v", lab.roleRun("reviewer", story))
	}
	waitFor(t, "the role's run ended", func() bool { r := lab.roleRun("reviewer", story); return r != nil && r.Ended != "" })
	run := lab.roleRun("reviewer", story)
	if run.Outcome != OutcomeWorked || run.Exit == nil || *run.Exit != 0 || run.RoleName != "reviewer" || run.Action != "review" {
		t.Errorf("run = %+v, want worked with exit 0", run)
	}
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity("reviewer")
		return err == nil && len(doc.Entries) == 1
	})
	doc, _ := lab.repo.Activity("reviewer")
	if e := doc.Entries[0]; e.Summary != "Reviewed "+story+": two criteria lack a test." || e.Trigger != "asked" || strings.Join(e.Items, ",") != story || e.Cost != 0.5 || e.Seconds != 61 {
		t.Errorf("activity = %+v, want the run's summary, asked, its story, its 0.5 USD, and its 61 s", e)
	}
	if rel, _ := filepath.Rel(lab.root, doc.Path); rel != filepath.Join("wip", "agents", "reviewer.md") {
		t.Errorf("activity document = %s, want wip/agents/reviewer.md", doc.Path)
	}
	if it, err := lab.repo.Get(story); err != nil || it.Usage != nil {
		t.Errorf("story's usage = %+v (%v), want nothing charged", it.Usage, err)
	}
	waitFor(t, "its end journalled", func() bool { return len(lab.entries()) == 2 })
	if e := lab.entries()[1]; e.Action != "review" || e.Method != MethodRole || e.Outcome != "done" || !strings.Contains(e.Detail, "the role reviewer on "+story+", as "+agent+" (pid ") || !strings.HasSuffix(e.Detail, "ended, worked") {
		t.Errorf("end entry = %+v", e)
	}
	// its start and its end are agent events on its story, in its verb, as a
	// planner run's are (S-0532)
	var tos []string
	for _, e := range agentEvents(lab.repo) {
		if e.ID != story || e.By != agent || e.Fields[inbox.FieldRole] != "review" {
			t.Errorf("agent event = %+v, want on %s by %s in role review", e, story, agent)
		}
		tos = append(tos, e.To)
	}
	if strings.Join(tos, ",") != inbox.AgentStarted+","+inbox.AgentEnded {
		t.Errorf("agent events = %v, want started, ended", tos)
	}
}

// S-0417: a defined role's run that ended with its question on its item open
// is asked; one that exited with a failure, failed, and journalled so.
func TestADefinedRoleRunIsJudgedByItsQuestionAndItsExit(t *testing.T) {
	lab := roleLab(t)
	asks, fails := lab.backlog("Asks", nil), lab.backlog("Fails", nil)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which criteria count?", On: asks, Author: "reviewer-" + asks, Text: "The scope is unclear.", Now: lab.now})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lab.outDir, "exit-reviewer-"+fails), []byte("3"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{asks, fails} {
		if !lab.roleHere("reviewer", id) {
			t.Fatalf("%s not started: %+v", id, lab.roleRun("reviewer", id))
		}
		waitFor(t, id+"'s run ended", func() bool { r := lab.roleRun("reviewer", id); return r != nil && r.Ended != "" })
	}
	if r := lab.roleRun("reviewer", asks); r.Outcome != OutcomeAsked || r.Thread != th.ID || !strings.Contains(r.Why, "waiting for an answer to "+th.ID) {
		t.Errorf("asked = %+v", r)
	}
	if r := lab.roleRun("reviewer", fails); r.Outcome != OutcomeFailed || r.Exit == nil || *r.Exit != 3 || r.Why != "ended (exit 3) running the role reviewer on "+fails {
		t.Errorf("failed = %+v", r)
	}
	waitFor(t, "both ends journalled", func() bool { return len(lab.entries()) == 4 })
	var failed bool
	for _, e := range lab.entries() {
		if e.Outcome == "failed" && e.Action == "review" && strings.HasSuffix(e.Detail, "ended, failed: ended (exit 3) running the role reviewer on "+fails) {
			failed = true
		}
	}
	if !failed {
		t.Errorf("no failed end journalled: %+v", lab.entries())
	}
}

// S-0417: a defined role's run whose process is gone while no launcher waits
// for it, as one flai role run started is, is settled at the next look.
func TestADefinedRoleRunGoneIsSettledAtALook(t *testing.T) {
	lab := roleLab(t)
	gone := exec.Command("true")
	if err := gone.Run(); err != nil {
		t.Fatal(err)
	}
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(&AgentRun{Agent: "auditor", Role: "audit", RoleName: "auditor", Action: "audit", PID: gone.Process.Pid, Started: "2026-10-03T10:00:00Z", Session: "s", Trigger: TriggerAsked})
	})
	lab.l.look(context.Background(), false)
	if r := lab.roleRun("auditor", ""); r == nil || r.Ended == "" || r.Outcome != OutcomeWorked || r.Exit != nil {
		t.Errorf("settled = %+v, want worked with no exit seen", r)
	}
	if st := lab.state(); st.Orchestrator != nil || st.Curator != nil {
		t.Errorf("state = %+v, want the run under roles alone", st)
	}
	if !strings.Contains(lab.logText(), `msg="role ended" component=serve project=t role=auditor item="" pid=`) {
		t.Errorf("the end is not logged:\n%s", lab.logText())
	}
	if j := lab.entries(); len(j) != 1 || j[0].Action != "audit" || !strings.HasSuffix(j[0].Detail, "ended, worked") {
		t.Errorf("journal: %+v", j)
	}
}

// S-0417: a defined role is refused for a built-in role's name, a name no
// file gives, and a file that does not load; while its host action is off;
// for an item when its file takes none, for none when it takes some, and for
// an item of a kind it does not take, archived, done, or cancelled; and when
// nothing can start it. A refusal records and journals nothing.
func TestADefinedRoleIsRefusedWithWhy(t *testing.T) {
	lab := roleLab(t)
	if err := os.WriteFile(filepath.Join(lab.root, "design", "workflows", "roles", "broken.md"), []byte("---\nname: [broken\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status := func(id, to string) {
		t.Helper()
		it, err := lab.repo.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(it.Path)
		if err := os.WriteFile(it.Path, []byte(strings.Replace(string(data), "status: backlog", "status: "+to, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	done, cancelled, archived := lab.backlog("Done", nil), lab.backlog("Dropped", nil), lab.backlog("Archived", nil)
	status(done, workitem.Done)
	status(cancelled, workitem.Cancelled)
	it, err := lab.repo.Get(archived)
	if err != nil {
		t.Fatal(err)
	}
	into := lab.repo.ItemDir(workitem.Story, true)
	if err := os.MkdirAll(into, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(it.Path, filepath.Join(into, filepath.Base(it.Path))); err != nil {
		t.Fatal(err)
	}
	open := lab.backlog("Open", nil)
	refusedFor := func(name, item, want string) {
		t.Helper()
		_, err := lab.runRole(name, item)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("run %s %s: %v, want refused for %q", name, item, err, want)
		}
	}
	refusedFor("planner", open, "planner is a built-in role, which flai starts with flai plan <item>")
	refusedFor("curator", "", "curator is a built-in role, which flai starts with flai curate")
	refusedFor("nobody", open, "the project defines no role nobody; the roles it defines are auditor, reviewer")
	refusedFor("broken", "", "the role file design/workflows/roles/broken.md does not load")
	refusedFor("reviewer", "", "the role reviewer runs on an item, a story, as its file design/workflows/roles/reviewer.md says; give it one")
	refusedFor("reviewer", lab.epic.ID, `the role reviewer runs on a story, as its file design/workflows/roles/reviewer.md says, and "`+lab.epic.ID+`" is none`)
	refusedFor("auditor", open, "the role auditor takes no item, as its file design/workflows/roles/auditor.md says, and was given "+open+"; run it on none")
	refusedFor("reviewer", done, done+" is done; the role reviewer runs on an item still open")
	refusedFor("reviewer", cancelled, cancelled+" is cancelled")
	refusedFor("reviewer", archived, archived+" is archived; the role reviewer runs on an item still on the board")
	if _, err := lab.runRole("reviewer", "S-0999"); err == nil || !strings.Contains(err.Error(), "S-0999 not found") {
		t.Errorf("an item that does not exist: %v", err)
	}
	lab.cfg.Command = nil
	refusedFor("reviewer", open, "the role reviewer's agent names no harness (roles.reviewer.agent or agent in system-flow.yaml), and no command is set on the host")
	lab.cfg.RoleActions = map[string]bool{"audit": true}
	refusedFor("reviewer", open, "the review host action, which the role reviewer's file design/workflows/roles/reviewer.md names, is off for this project: flai serve enable review")
	if st := lab.state(); len(st.Roles) != 0 || len(lab.entries()) != 0 {
		t.Errorf("a refusal was recorded: %+v %+v", st.Roles, lab.entries())
	}
}

// S-0417: the agent notification of a defined role's run names its verb, its
// name, and its item when it ran on one; the others' say what they did.
func TestTheAgentNotificationNamesADefinedRole(t *testing.T) {
	for _, c := range []struct {
		run  AgentRun
		want map[string]string
	}{
		{AgentRun{Item: "S-0001", Role: "review", RoleName: "reviewer"}, map[string]string{"project": "t", "role": "review", "name": "reviewer", "item": "S-0001"}},
		{AgentRun{Role: "audit", RoleName: "auditor"}, map[string]string{"project": "t", "role": "audit", "name": "auditor"}},
		{AgentRun{Role: "curate", RoleName: "keeper"}, map[string]string{"project": "t", "role": "curate", "name": "keeper"}},
		{AgentRun{Item: "E-0001"}, map[string]string{"project": "t", "item": "E-0001"}},
		{AgentRun{}, map[string]string{"project": "t", "role": "orchestrate"}},
		{AgentRun{Focus: "risk"}, map[string]string{"project": "t", "role": "analyze"}},
		{AgentRun{Role: "curate"}, map[string]string{"project": "t", "role": "curate"}},
		{AgentRun{Story: "S-0002"}, map[string]string{"project": "t", "story": "S-0002"}},
	} {
		got := agentAbout("t", &c.run)
		if len(got) != len(c.want) {
			t.Errorf("about %+v = %v, want %v", c.run, got, c.want)
			continue
		}
		for k, v := range c.want {
			if got[k] != v {
				t.Errorf("about %+v = %v, want %v", c.run, got, c.want)
				break
			}
		}
	}
}

// S-0417: agents.json without roles still reads, and a defined role's run is
// kept under roles, apart from the planner's, the curator's, and the
// orchestrator's, whatever its verb and item.
func TestAgentStatesKeepDefinedRolesApart(t *testing.T) {
	d := Dir(t.TempDir())
	d.updateAgent("/p", func(s *AgentState) {
		s.put(&AgentRun{Item: "S-0001", Agent: "reviewer-S-0001", Role: "review", RoleName: "reviewer", Started: "2026-10-03T00:00:00Z"})
		s.put(&AgentRun{Agent: "keeper", Role: "curate", RoleName: "keeper", Started: "2026-10-03T00:00:00Z"})
	})
	st := d.AgentStates()["/p"]
	if st.Roles["reviewer"]["S-0001"] == nil || st.Roles["keeper"][""] == nil || st.Plans != nil || st.Curator != nil || st.Orchestrator != nil {
		t.Errorf("state = %+v, want both runs under roles alone", st)
	}
	run := st.Roles["reviewer"]["S-0001"]
	if !st.of(run).same(run) {
		t.Errorf("of = %+v, want the run", st.of(run))
	}
}
