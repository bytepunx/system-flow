package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/statsread"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0220, end to end: the orchestrator's thread calls go as Claude Code
// makes them in its session, first through flai guard's decision, built from
// the PreToolUse hook's input with the project's
// orchestration.permissions.answer_threads, and then, where the guard allows
// them, through the MCP server it serves the orchestrator in role
// orchestrate. A story's agent waits on a thread it asked on its story; the
// orchestrator recommends an answer for the operator to confirm, answers it
// itself citing its source, or escalates what it cannot source. Each
// scenario checks the thread's status, who it awaits, the orchestrator's
// decision log (its activity document, wip/agents/orchestrator.md), and the
// wait as flai stats counts it.

// syncedClock is a now that the tests move while a held wait_for_work reads
// it.
type syncedClock struct {
	mu sync.Mutex
	at time.Time
}

func (c *syncedClock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.at
}

func (c *syncedClock) set(at time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.at = at
}

// noGitRunner finds no git, so flai stats leaves touches drift out, as it
// does on a host without git.
type noGitRunner struct{}

func (noGitRunner) Run(string, string, ...string) (string, error) {
	return "", errors.New("no git")
}

func (noGitRunner) RunInput(string, string, string, ...string) (string, error) {
	return "", errors.New("no git")
}

func (noGitRunner) LookPath(string) (string, error) { return "", errors.New("git not on PATH") }

// answering is a project whose story's agent, claude, asked a question on
// its story's thread, and an orchestrator session on the same project.
type answering struct {
	asker  *fixture // the story's agent's session
	orch   *fixture // the orchestrator's session, in role orchestrate
	root   string
	clock  *syncedClock
	thread string // the question the story's agent asked
}

// The times of each scenario: the story is in progress from t0; its agent
// asks at asked and, idle, moves the story to review at idle; the
// orchestrator replies at replied, the operator confirms at confirmed, and
// accepts the story at accepted; flai stats runs at measured.
var (
	askedAt     = t0.Add(time.Minute)
	idleAt      = t0.Add(2 * time.Minute)
	repliedAt   = t0.Add(5 * time.Minute)
	confirmedAt = t0.Add(10 * time.Minute)
	acceptedAt  = t0.Add(30 * time.Minute)
	measuredAt  = t0.Add(time.Hour)
)

// The source the orchestrator cites, and the line it is logged with.
const (
	answerSource = "design/system/plan.md#Shape"
	answerCited  = "design/system/plan.md § Shape"
)

// answeringFixture sets the project's answer_threads to mode, starts the
// orchestrator's MCP server beside the story's agent's, which logs the
// orchestrator's activities in its activity document as flai mcp does, and
// has the story's agent ask on its story and go idle.
func answeringFixture(t *testing.T, mode string) *answering {
	t.Helper()
	clock := &syncedClock{at: askedAt}
	t.Setenv("FLAI_ROLE", "")
	asker := setupWith(t, func(o *Options) { o.Now = clock.now })
	root := projectRoot(asker.repo)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\norchestration:\n  permissions:\n    answer_threads: " + mode + "\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := repo.Manifest.Orchestration.Permissions.AnswerThreads; got != mode {
		t.Fatalf("answer_threads: %q, want %q", got, mode)
	}
	// flai serve sets the orchestrator's role in its session, and flai mcp
	// inherits it
	t.Setenv("FLAI_ROLE", guard.RoleOrchestrate)
	logs := serve.Dir(t.TempDir())
	opt := Options{Repo: repo, Agent: workitem.ActivityOrchestrator, Version: "test", Now: clock.now, Poll: 20 * time.Millisecond, MaxWait: 3 * time.Second,
		Activities: func(_ context.Context, root, kind, summary string, items []string, _ string) (ActivityLogged, error) {
			logged, err := serve.LogActivity(logs, root, repo.Manifest.Key, kind, summary, items, clock.now())
			if err != nil {
				return ActivityLogged{}, err
			}
			return ActivityLogged{Entry: ActivityEntry{Summary: logged.Entry.Summary, Items: logged.Entry.Items}}, nil
		}}
	a := &answering{asker: asker, orch: &fixture{repo: repo, cs: connectServer(t, opt), story: asker.story, task: asker.task}, root: root, clock: clock}

	out, failed := asker.call(t, "thread_open", map[string]any{"on": asker.story.ID, "title": "Which port?", "text": "Eight or nine?"})
	if failed != "" {
		t.Fatal(failed)
	}
	a.thread = out["id"].(string)
	clock.set(idleAt)
	story, _ := asker.repo.Get(asker.story.ID)
	if _, err := asker.repo.Transition(story, workitem.Review, "claude", "", idleAt); err != nil {
		t.Fatal(err)
	}
	return a
}

// connectServer starts an MCP server with opt and connects a client to it.
func connectServer(t *testing.T, opt Options) *mcp.ClientSession {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	if _, err := New(opt).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "orchestrator", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

// decide is flai guard's decision on the orchestrator's call e, made as the
// hook makes it in the project: with the manifest's permissions and the
// project's threads to tell who opened one, and a refusal logged under
// ## Refusals in the orchestrator's activity document.
func (a *answering) decide(t *testing.T, e guard.Event) guard.Refusal {
	t.Helper()
	repo, err := workitem.Open(a.root)
	if err != nil {
		t.Fatal(err)
	}
	gd := guard.Guard{Commands: []string{"stats", "thread"}, Role: guard.RoleOrchestrate, Permissions: repo.Manifest.Orchestration.Permissions, Opener: guard.ThreadOpener(repo)}
	r := gd.Decide(e)
	if r.Why != "" && r.Call != "" {
		if _, err := repo.AppendRefusal(workitem.ActivityOrchestrator, workitem.ActivityRefusal{At: a.clock.now(), Call: r.Call, Needs: r.Needs}); err != nil {
			t.Fatal(err)
		}
	}
	return r
}

// tool makes the orchestrator's call of the MCP tool with args as Claude
// Code makes it: the guard's PreToolUse decision on the hook's input, and the
// tool only when the guard lets it through. refused is the guard's refusal,
// and failed the tool's error.
func (a *answering) tool(t *testing.T, name string, args map[string]any) (out map[string]any, refused guard.Refusal, failed string) {
	t.Helper()
	in, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}
	var e guard.Event
	if err := json.Unmarshal([]byte(`{"tool_name":"`+guard.MCPPrefix+name+`","tool_input":`+string(in)+`}`), &e); err != nil {
		t.Fatal(err)
	}
	if r := a.decide(t, e); r.Why != "" {
		return nil, r, ""
	}
	out, failed = a.orch.call(t, name, args)
	return out, guard.Refusal{}, failed
}

// refusedFor checks that the guard refused a call with the words and the
// permission given, and that nothing reached the thread.
func (a *answering) refusedFor(t *testing.T, what string, out map[string]any, r guard.Refusal, words, needs string) {
	t.Helper()
	if out != nil || !strings.Contains(r.Why, words) || r.Needs != needs || r.Call == "" {
		t.Errorf("%s: %+v (out %v), want refused with %q needing %q", what, r, out, words, needs)
	}
}

// awaits checks the thread's status on disk and as the story's agent sees it
// in its inbox: whether it awaits that agent, and whether a recommendation
// awaits the operator's confirmation.
func (a *answering) awaits(t *testing.T, status string, asker, pending bool) {
	t.Helper()
	th, err := threads.Get(a.asker.repo, a.thread)
	if err != nil {
		t.Fatal(err)
	}
	if th.Status != status || (th.PendingRecommendation() != nil) != pending {
		t.Errorf("%s: status %s, pending %v; want %s, pending %v", a.thread, th.Status, th.PendingRecommendation(), status, pending)
	}
	out, _ := a.asker.call(t, "inbox", map[string]any{"all": true})
	var sum map[string]any
	for _, s := range out["threads"].([]any) {
		if s := s.(map[string]any); s["id"] == a.thread {
			sum = s
		}
	}
	want, wantN := "other", 0.0
	if asker {
		want, wantN = "you", 1
	}
	if sum == nil || sum["awaiting"] != want || out["awaiting_you"] != wantN || (sum["pending_recommendation"] != nil) != pending || sum["status"] != status {
		t.Errorf("the asker's inbox: %v; want %s awaiting %s, pending %v", out, a.thread, want, pending)
	}
}

// decisions is the orchestrator's activity document: its activities and the
// calls the guard refused it.
func (a *answering) decisions(t *testing.T) *workitem.Activity {
	t.Helper()
	repo, err := workitem.Open(a.root)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := repo.Activity(workitem.ActivityOrchestrator)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

// loggedOnce checks that the orchestrator's activity document holds one
// activity, summary, on the story.
func (a *answering) loggedOnce(t *testing.T, summary string) {
	t.Helper()
	doc := a.decisions(t)
	if len(doc.Entries) != 1 || doc.Entries[0].Summary != summary || strings.Join(doc.Entries[0].Items, ",") != a.asker.story.ID {
		t.Errorf("decision log: %+v, want one entry %q on %s", doc.Entries, summary, a.asker.story.ID)
	}
}

// refusalsLogged checks the calls the guard refused, as its activity
// document logs them, in order: each call's start and the permission it
// needs.
func (a *answering) refusalsLogged(t *testing.T, want ...[2]string) {
	t.Helper()
	got := a.decisions(t).Refusals
	ok := len(got) == len(want)
	for i := 0; ok && i < len(want); i++ {
		ok = strings.HasPrefix(got[i].Call, want[i][0]) && got[i].Needs == want[i][1]
	}
	if !ok {
		t.Errorf("refusals logged: %+v, want %v", got, want)
	}
}

// stats accepts the story and reads its wait as flai stats does, returning
// the week it was completed in and the story's own metrics.
func (a *answering) stats(t *testing.T) (metrics.WaitWeek, metrics.ItemMetrics) {
	t.Helper()
	repo := a.asker.repo
	for _, st := range []string{workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done} {
		task, _ := repo.Get(a.asker.task.ID)
		if _, err := repo.Transition(task, st, "claude", "", acceptedAt); err != nil {
			t.Fatal(err)
		}
	}
	story, _ := repo.Get(a.asker.story.ID)
	if _, err := repo.Transition(story, workitem.Done, "alex", "", acceptedAt); err != nil {
		t.Fatal(err)
	}
	window, err := metrics.ParseWindow(metrics.DefaultWindow)
	if err != nil {
		t.Fatal(err)
	}
	items, opt, err := statsread.Read(noGitRunner{}, repo, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	opt.Now, opt.Since, opt.Type, opt.Bucket = measuredAt, window, workitem.Story, metrics.DefaultBucket
	rep := metrics.Compute(items, opt)
	var week metrics.WaitWeek
	for _, w := range rep.Waiting.Weeks {
		if w.Items > 0 {
			week = w
		}
	}
	var item metrics.ItemMetrics
	for _, m := range rep.Items {
		if m.ID == a.asker.story.ID {
			item = m
		}
	}
	if week.Items != 1 || item.ID == "" {
		t.Fatalf("the accepted story is not measured: weeks %+v, items %+v", rep.Waiting.Weeks, rep.Items)
	}
	return week, item
}

// waited checks the story's wait on its thread: the seconds it waited while
// in progress, the part the orchestrator ended, and the week's counts of the
// waits the orchestrator and a confirmation ended.
func waited(t *testing.T, week metrics.WaitWeek, item metrics.ItemMetrics, total, orchestrator float64, byOrchestrator, byConfirmation metrics.WaitCount) {
	t.Helper()
	if item.WaitThreads == nil || *item.WaitThreads != total || item.WaitThreadsOrchestrator == nil || *item.WaitThreadsOrchestrator != orchestrator {
		shown := func(v *float64) any {
			if v == nil {
				return "absent"
			}
			return *v
		}
		t.Errorf("the story's waits: threads %v, orchestrator %v; want %v, %v", shown(item.WaitThreads), shown(item.WaitThreadsOrchestrator), total, orchestrator)
	}
	if week.Threads.Total != total || week.Threads.Orchestrator != byOrchestrator || week.Threads.Confirmed != byConfirmation {
		t.Errorf("the week's thread waits: %+v; want total %v, orchestrator %+v, confirmed %+v", week.Threads, total, byOrchestrator, byConfirmation)
	}
}

// The story is in progress from t0 until its agent goes idle, so each wait
// on its thread, from askedAt, counts that long.
var inProgressWait = idleAt.Sub(askedAt).Seconds()

// S-0220: with answer_threads recommend, the orchestrator's answer is
// refused and its recommendation citing a source is posted: the thread still
// awaits the operator, the agent that asked is not woken, and the decision
// log names the source. The operator's confirmation, which the orchestrator
// may not give, makes it the answer and wakes the agent, and flai stats
// counts the wait as ended by a confirmation.
func TestTheOrchestratorRecommendsAndTheOperatorConfirms(t *testing.T) {
	a := answeringFixture(t, manifest.AnswerRecommend)
	done := a.asker.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	a.clock.set(repliedAt)

	out, r, _ := a.tool(t, "thread_reply", map[string]any{"id": a.thread, "text": "Nine.", "source": answerSource})
	a.refusedFor(t, "a sourced answer under recommend", out, r, "answer_threads is recommend, so it replies to another's thread only as a recommendation", "orchestration.permissions.answer_threads")
	out, r, _ = a.tool(t, "thread_reply", map[string]any{"id": a.thread, "text": "Nine."})
	a.refusedFor(t, "an unsourced answer under recommend", out, r, "only as a recommendation for the operator to confirm", "")
	a.awaits(t, "open", false, false)

	out, r, failed := a.tool(t, "thread_reply", map[string]any{"id": a.thread, "text": "Nine, as the plan's shape has it.", "recommendation": true, "source": answerSource})
	if r.Why != "" || failed != "" || out["status"] != "open" {
		t.Fatalf("the recommendation: %v, refused %+v, failed %s", out, r, failed)
	}
	select {
	case woke := <-done:
		t.Fatalf("a recommendation woke the agent that asked: %v", woke)
	case <-time.After(600 * time.Millisecond):
	}
	a.awaits(t, "open", false, true)
	a.loggedOnce(t, "Recommended an answer on "+a.thread+", citing "+answerCited)

	if r := a.decide(t, guard.Event{ToolName: "Bash", ToolInput: guard.Input{Command: "flai thread confirm " + a.thread}}); !strings.Contains(r.Why, "confirming a recommendation is the operator's alone") {
		t.Errorf("the orchestrator's confirmation: %+v", r)
	}
	a.refusalsLogged(t,
		[2]string{"thread_reply " + a.thread + " source", "orchestration.permissions.answer_threads"},
		[2]string{"thread_reply " + a.thread, ""},
		[2]string{"flai thread confirm " + a.thread, ""})

	// the operator confirms it, as flai thread confirm does
	a.clock.set(confirmedAt)
	if _, err := threads.Confirm(a.asker.repo, a.thread, "alex", confirmedAt); err != nil {
		t.Fatal(err)
	}
	woke := answered(t, done)
	if ths := woke["threads"].([]any); woke["reason"] != WorkThread || len(ths) != 1 || ths[0].(map[string]any)["id"] != a.thread {
		t.Errorf("the confirmation wakes the agent that asked: %v", woke)
	}
	a.awaits(t, "answered", true, false)

	week, item := a.stats(t)
	waited(t, week, item, inProgressWait, 0, metrics.WaitCount{}, metrics.WaitCount{Count: 1, Total: inProgressWait})
}

// S-0220: with answer_threads autonomous, the orchestrator answers citing its
// source: the thread is answered, the agent that asked is woken, the
// decision log names the source, and flai stats counts the wait as ended by
// the orchestrator.
func TestTheOrchestratorAnswersCitingItsSource(t *testing.T) {
	a := answeringFixture(t, manifest.AnswerAutonomous)
	done := a.asker.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	a.clock.set(repliedAt)

	out, r, failed := a.tool(t, "thread_reply", map[string]any{"id": a.thread, "text": "Nine, as the plan's shape has it.", "source": answerSource})
	if r.Why != "" || failed != "" || out["status"] != "answered" {
		t.Fatalf("the answer: %v, refused %+v, failed %s", out, r, failed)
	}
	if entries := out["entry_list"].([]any); entries[len(entries)-1].(map[string]any)["recommendation"] != false {
		t.Errorf("an answer is no recommendation: %v", entries[len(entries)-1])
	}
	woke := answered(t, done)
	if ths := woke["threads"].([]any); woke["reason"] != WorkThread || len(ths) != 1 || ths[0].(map[string]any)["id"] != a.thread {
		t.Errorf("the answer wakes the agent that asked: %v", woke)
	}
	a.awaits(t, "answered", true, false)
	a.loggedOnce(t, "Answered "+a.thread+", citing "+answerCited)
	a.refusalsLogged(t)

	week, item := a.stats(t)
	waited(t, week, item, inProgressWait, inProgressWait, metrics.WaitCount{Count: 1, Total: inProgressWait}, metrics.WaitCount{})
}

// S-0220: with answer_threads autonomous, an answer the orchestrator cannot
// source is refused, and it escalates it as a recommendation, which leaves the
// thread awaiting the operator and the agent that asked asleep. It may
// resolve neither that thread nor answer one it opened itself; each refusal
// is logged. flai stats counts no wait ended, by the orchestrator or by a
// confirmation.
func TestTheOrchestratorEscalatesWhatItCannotSource(t *testing.T) {
	a := answeringFixture(t, manifest.AnswerAutonomous)
	done := a.asker.held(t, 1)
	time.Sleep(150 * time.Millisecond)
	a.clock.set(repliedAt)

	out, r, _ := a.tool(t, "thread_reply", map[string]any{"id": a.thread, "text": "Nine."})
	a.refusedFor(t, "an unsourced answer under autonomous", out, r, "it answers another's thread only citing what the answer rests on", "")
	a.awaits(t, "open", false, false)

	out, r, failed := a.tool(t, "thread_reply", map[string]any{"id": a.thread, "text": "Nine.", "recommendation": true})
	if r.Why != "" || failed != "" || out["status"] != "open" {
		t.Fatalf("the escalation: %v, refused %+v, failed %s", out, r, failed)
	}
	if woke := answered(t, done); woke["reason"] != "" || woke["timed_out"] != true || len(woke["threads"].([]any)) != 0 {
		t.Errorf("an escalation woke the agent that asked: %v", woke)
	}
	a.awaits(t, "open", false, true)
	a.loggedOnce(t, "Recommended an answer on "+a.thread+", citing no source")

	out, r, _ = a.tool(t, "thread_resolve", map[string]any{"id": a.thread, "reason": "escalated"})
	a.refusedFor(t, "resolving the asker's thread", out, r, "it resolves only a thread it opened", "")
	a.awaits(t, "open", false, true)

	own, r, failed := a.tool(t, "thread_open", map[string]any{"on": "design/system/plan.md", "title": "Is the shape final?", "text": "The plan names one shape."})
	if r.Why != "" || failed != "" {
		t.Fatalf("the orchestrator's own thread: refused %+v, failed %s", r, failed)
	}
	ownID := own["id"].(string)
	for what, args := range map[string]map[string]any{
		"an answer on its own thread":        {"id": ownID, "text": "Yes.", "source": answerSource},
		"a recommendation on its own thread": {"id": ownID, "text": "Yes.", "recommendation": true},
	} {
		out, r, _ := a.tool(t, "thread_reply", args)
		a.refusedFor(t, what, out, r, "it never recommends or answers its own question", "")
	}
	if th, err := threads.Get(a.asker.repo, ownID); err != nil || len(th.Entries()) != 1 || th.Status != "open" {
		t.Errorf("its own thread took a reply: %+v, %v", th, err)
	}
	a.loggedOnce(t, "Recommended an answer on "+a.thread+", citing no source")
	a.refusalsLogged(t,
		[2]string{"thread_reply " + a.thread, ""},
		[2]string{"thread_resolve " + a.thread, ""},
		[2]string{"thread_reply " + ownID, ""},
		[2]string{"thread_reply " + ownID, ""})

	week, item := a.stats(t)
	waited(t, week, item, inProgressWait, 0, metrics.WaitCount{}, metrics.WaitCount{})
}
