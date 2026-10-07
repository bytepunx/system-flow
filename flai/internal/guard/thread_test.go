package guard

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// recommendation marks the orchestrator's thread_reply a recommendation.
func recommendation(e Event) Event {
	e.ToolInput.Recommendation = true
	return e
}

// openers knows TH-0001, which the orchestrator opened, TH-0002, which the
// operator opened, TH-0003, which a story's planner opened, and TH-0004,
// which an epic's planner opened; it cannot tell of any other.
func openers(id string) (string, error) {
	switch id {
	case "TH-0001":
		return workitem.ActivityOrchestrator, nil
	case "TH-0002":
		return "alex", nil
	case "TH-0003":
		return "planner-S-0328", nil
	case "TH-0004":
		return "planner-E-0016", nil
	}
	return "", errors.New(id + " not found")
}

// threadGuard is the orchestrator's guard under p, telling threads apart by
// openers.
func threadGuard(p manifest.Permissions) Guard {
	gd := orchestrator(p)
	gd.Opener = openers
	return gd
}

// The forms of the orchestrator's thread calls, by MCP and on the command
// line, on the thread id.
func threadCalls(t *testing.T, id string) map[string][]Event {
	return map[string][]Event{
		"recommendation": {
			call(t, "thread_reply", `{"id":"`+id+`","text":"S-0002 next","recommendation":true}`),
			call(t, "thread_reply", `{"id":"`+id+`","text":"S-0002 next","recommendation":true,"source":"design/system/workflow.md#Pull order"}`),
			bash("", "flai thread reply "+id+" --recommend 'S-0002 next'"),
			bash("", "scripts/flai.sh --config c.json thread reply --recommend=true "+id+" 'S-0002 next' --json"),
		},
		"sourced answer": {
			call(t, "thread_reply", `{"id":"`+id+`","text":"S-0002 next","source":"design/adrs/0090-x.md"}`),
			bash("", "flai thread reply --source 'design/system/workflow.md#Pull order' "+id+" 'S-0002 next'"),
			bash("", "flai thread reply "+id+" --source=design/adrs/0090-x.md 'S-0002 next'"),
		},
		"unsourced answer": {
			call(t, "thread_reply", `{"id":"`+id+`","text":"S-0002 next"}`),
			call(t, "thread_reply", `{"id":"`+id+`","text":"S-0002 next","recommendation":false,"source":"  "}`),
			bash("", "flai thread reply "+id+" 'S-0002 next'"),
			bash("", "flai thread reply --recommend=false --source= "+id+" 'S-0002 next'"),
		},
		"resolve": {
			call(t, "thread_resolve", `{"id":"`+id+`","reason":"settled"}`),
			bash("", "flai thread resolve "+id+" --reason settled"),
		},
	}
}

// S-0220: under each value of answer_threads, the orchestrator follows up
// and resolves on a thread it opened but never recommends or answers there;
// on another's it replies only while answer_threads is on, as a
// recommendation, or with autonomous as an answer citing its source, and
// never resolves it. Each refusal names the permission or the parameter that
// would allow the call.
func TestTheOrchestratorsThreadCallsFollowWhoOpenedTheThreadAndAnswerThreads(t *testing.T) {
	const (
		needs     = "orchestration.permissions.answer_threads"
		offWhy    = ": it needs " + needs + ", which is off. "
		recommend = "give recommendation true (--recommend)"
		source    = "give source <path> or <path>#<heading> (--source on the command line), or post it as a recommendation for the operator to confirm with recommendation true (--recommend)"
	)
	never := func(why string) string {
		return ": the orchestrator never does it, whatever its permissions: " + why + ". "
	}
	// verdict is what the guard should say: pass, or the words the refusal
	// holds and the permission it names.
	type verdict struct {
		pass  bool
		words string
		needs string
	}
	ok := verdict{pass: true}
	for _, c := range []struct {
		mode  string
		own   map[string]verdict
		other map[string]verdict
	}{
		{"", nil, map[string]verdict{
			"recommendation":   {words: offWhy, needs: needs},
			"sourced answer":   {words: offWhy, needs: needs},
			"unsourced answer": {words: offWhy, needs: needs},
		}},
		{manifest.AnswerOff, nil, map[string]verdict{
			"recommendation":   {words: offWhy, needs: needs},
			"sourced answer":   {words: offWhy, needs: needs},
			"unsourced answer": {words: offWhy, needs: needs},
		}},
		{manifest.AnswerRecommend, nil, map[string]verdict{
			"recommendation":   ok,
			"sourced answer":   {words: "answer_threads is recommend, so it replies to another's thread only as a recommendation for the operator to confirm: " + recommend + "; autonomous would let it answer citing its source. ", needs: needs},
			"unsourced answer": {words: "answer_threads is recommend, so it replies to another's thread only as a recommendation for the operator to confirm: " + recommend + "; autonomous would let it answer only citing a source. "},
		}},
		{manifest.AnswerAutonomous, nil, map[string]verdict{
			"recommendation":   ok,
			"sourced answer":   ok,
			"unsourced answer": {words: ": it answers another's thread only citing what the answer rests on, so " + source + ". "},
		}},
	} {
		c.own = map[string]verdict{
			"recommendation":   {words: never(answersOwn)},
			"sourced answer":   {words: never(answersOwn)},
			"unsourced answer": ok,
			"resolve":          ok,
		}
		c.other["resolve"] = verdict{words: never(resolvesOwn)}
		gd := threadGuard(manifest.Permissions{AnswerThreads: c.mode})
		for _, side := range []struct {
			name     string
			id       string
			verdicts map[string]verdict
		}{{"opened by the orchestrator", "TH-0001", c.own}, {"opened by another", "TH-0002", c.other}, {"that cannot be read", "TH-0099", c.other}} {
			for form, events := range threadCalls(t, side.id) {
				want := side.verdicts[form]
				for _, e := range events {
					r := gd.Decide(e)
					switch {
					case want.pass && r != (Refusal{}):
						t.Errorf("answer_threads %q, %s on a thread %s: %s refused: %+v", c.mode, form, side.name, described(e), r)
					case want.pass:
					case !strings.HasPrefix(r.Why, "the orchestrator cannot ") || !strings.Contains(r.Why, want.words) || !strings.HasSuffix(r.Why, askOperator) || r.Needs != want.needs || r.Call == "":
						t.Errorf("answer_threads %q, %s on a thread %s: %s: %+v, want %q naming %q", c.mode, form, side.name, described(e), r, want.words, want.needs)
					}
				}
			}
		}
	}
}

// S-0220: the orchestrator never confirms a recommendation, which is the
// operator's alone, nor writes to a thread as another, whatever its
// permissions and whoever opened the thread.
func TestTheOrchestratorNeverConfirmsOrWritesAsAnother(t *testing.T) {
	for _, mode := range []string{"", manifest.AnswerOff, manifest.AnswerRecommend, manifest.AnswerAutonomous} {
		gd := threadGuard(manifest.Permissions{AnswerThreads: mode})
		for why, cmds := range map[string][]string{
			confirms: {
				"flai thread confirm TH-0001",
				"flai thread confirm TH-0002",
				"flai thread confirm --by orchestrator TH-0002",
			},
			writesAsSelf: {
				"flai thread reply --by alex --recommend TH-0002 'S-0002 next'",
				"flai thread reply TH-0001 --by=alex 'follow-up'",
				"flai thread resolve TH-0001 --by alex",
			},
		} {
			for _, cmd := range cmds {
				if r := gd.Decide(bash("", cmd)); !strings.Contains(r.Why, ": the orchestrator never does it, whatever its permissions: "+why+". ") || r.Needs != "" || r.Call == "" {
					t.Errorf("answer_threads %q: %q: %+v", mode, cmd, r)
				}
			}
		}
		if r := gd.Decide(bash("", "flai thread resolve --by orchestrator TH-0001")); r.Why != "" {
			t.Errorf("answer_threads %q: resolve as itself refused: %s", mode, r.Why)
		}
	}
}

// S-0220: the refusals of thread calls say what the call would do, and
// carry the call with its marks for the orchestrator's activity document.
func TestTheOrchestratorsThreadRefusalsSayTheCall(t *testing.T) {
	for _, c := range []struct {
		p    manifest.Permissions
		e    Event
		why  string
		call string
	}{
		{manifest.Permissions{}, call(t, "thread_reply", `{"id":"TH-0002","text":"x"}`),
			"the orchestrator cannot reply to TH-0002: it needs orchestration.permissions.answer_threads, which is off. " + askOperator,
			"thread_reply TH-0002"},
		{allOn, call(t, "thread_reply", `{"id":"TH-0001","text":"x","recommendation":true,"source":" design/adrs/0090-x.md#Decision "}`),
			"the orchestrator cannot recommend an answer to TH-0001: the orchestrator never does it, whatever its permissions: " + answersOwn + ". " + askOperator,
			"thread_reply TH-0001 recommendation source design/adrs/0090-x.md#Decision"},
		{allOn, call(t, "thread_resolve", `{"id":"TH-0002"}`),
			"the orchestrator cannot resolve TH-0002: the orchestrator never does it, whatever its permissions: " + resolvesOwn + ". " + askOperator,
			"thread_resolve TH-0002"},
		{allOn, call(t, "thread_reply", `{"text":"x"}`),
			"the orchestrator cannot reply to a thread: it answers another's thread only citing what the answer rests on, so give source <path> or <path>#<heading> (--source on the command line), or post it as a recommendation for the operator to confirm with recommendation true (--recommend). " + askOperator,
			"thread_reply"},
		{manifest.Permissions{AnswerThreads: manifest.AnswerRecommend}, bash("", "cd /w && flai thread reply TH-0002 --source design/adrs/0090-x.md 'x'"),
			`the orchestrator cannot run "flai thread reply TH-0002 --source design/adrs/0090-x.md x": orchestration.permissions.answer_threads is recommend, so it replies to another's thread only as a recommendation for the operator to confirm: give recommendation true (--recommend); autonomous would let it answer citing its source. ` + askOperator,
			"flai thread reply TH-0002 --source design/adrs/0090-x.md x"},
	} {
		if r := threadGuard(c.p).Decide(c.e); r.Why != c.why || r.Call != c.call {
			t.Errorf("%s:\n got %q, call %q\nwant %q, call %q", described(c.e), r.Why, r.Call, c.why, c.call)
		}
	}
}

// S-0220: without Opener every thread is another's, so the orchestrator
// neither follows up nor resolves while answer_threads is off.
func TestWithoutAnOpenerEveryThreadIsAnothers(t *testing.T) {
	gd := orchestrator(manifest.Permissions{})
	if r := gd.Decide(call(t, "thread_reply", `{"id":"TH-0001","text":"x"}`)); r.Needs != "orchestration.permissions.answer_threads" {
		t.Errorf("reply: %+v", r)
	}
	if r := gd.Decide(call(t, "thread_resolve", `{"id":"TH-0001"}`)); !strings.Contains(r.Why, resolvesOwn) {
		t.Errorf("resolve: %+v", r)
	}
}

// S-0328, ADR-0119: on a thread a story's planner opened, plan_backlog_stories
// lets the orchestrator answer, with or without a source, recommend, and
// resolve, whatever answer_threads is. While it is off the thread is held
// as another's under answer_threads, but resolving it, and replying while
// answer_threads is off, are refused for want of plan_backlog_stories, and
// every other refusal names it too. An epic's planner's thread is another's
// either way.
func TestTheOrchestratorSettlesAStoryPlannersThreadsWithPlanBacklogStories(t *testing.T) {
	const needs = "orchestration.permissions.plan_backlog_stories"
	suffix := "; on a thread a story's planner opened, " + needs + ", which is off, would let it reply so. " + askOperator
	story, other, epic := threadCalls(t, "TH-0003"), threadCalls(t, "TH-0002"), threadCalls(t, "TH-0004")
	for _, mode := range []string{"", manifest.AnswerOff, manifest.AnswerRecommend, manifest.AnswerAutonomous} {
		on := threadGuard(manifest.Permissions{AnswerThreads: mode, PlanBacklogStories: true})
		off := threadGuard(manifest.Permissions{AnswerThreads: mode})
		answering := off.Permissions.Allows(manifest.PermitAnswerThreads)
		for form, events := range story {
			for i, e := range events {
				if r := on.Decide(e); r != (Refusal{}) {
					t.Errorf("answer_threads %q, plan_backlog_stories on: %s %s refused: %+v", mode, form, described(e), r)
				}
				r, held := off.Decide(e), off.Decide(other[form][i])
				switch {
				case form == "resolve" || !answering:
					if !strings.HasPrefix(r.Why, "the orchestrator cannot ") || !strings.HasSuffix(r.Why, ": it needs "+needs+", which is off. "+askOperator) || r.Needs != needs || r.Call == "" {
						t.Errorf("answer_threads %q, plan_backlog_stories off: %s %s: %+v", mode, form, described(e), r)
					}
				case held == (Refusal{}):
					if r != (Refusal{}) {
						t.Errorf("answer_threads %q, plan_backlog_stories off: %s %s refused, unlike on another's thread: %+v", mode, form, described(e), r)
					}
				default:
					want := strings.TrimSuffix(strings.ReplaceAll(held.Why, "TH-0002", "TH-0003"), ". "+askOperator) + suffix
					if r.Why != want || r.Needs != needs || r.Call != strings.ReplaceAll(held.Call, "TH-0002", "TH-0003") {
						t.Errorf("answer_threads %q, plan_backlog_stories off: %s %s:\n got %+v\nwant %q", mode, form, described(e), r, want)
					}
				}
			}
		}
		for form, events := range epic {
			for i, e := range events {
				for _, gd := range []Guard{on, off} {
					r, held := gd.Decide(e), gd.Decide(other[form][i])
					if strings.ReplaceAll(r.Why, "TH-0004", "TH-0002") != held.Why || r.Needs != held.Needs {
						t.Errorf("answer_threads %q, plan_backlog_stories %v: %s on an epic's planner's thread %s: %+v, want it held as another's: %+v", mode, gd.Permissions.PlanBacklogStories, form, described(e), r, held)
					}
				}
			}
		}
		for why, cmds := range map[string][]string{
			confirms:     {"flai thread confirm TH-0003"},
			writesAsSelf: {"flai thread reply --by planner-S-0328 TH-0003 'yes'", "flai thread resolve TH-0003 --by=alex"},
		} {
			for _, cmd := range cmds {
				if r := on.Decide(bash("", cmd)); !strings.Contains(r.Why, ": the orchestrator never does it, whatever its permissions: "+why+". ") || r.Needs != "" {
					t.Errorf("answer_threads %q, plan_backlog_stories on: %q: %+v", mode, cmd, r)
				}
			}
		}
		if r := on.Decide(bash("", "flai thread resolve --by orchestrator TH-0003")); r.Why != "" {
			t.Errorf("answer_threads %q: resolve as itself refused: %s", mode, r.Why)
		}
	}
}

// S-0328: a story's planner is planner- and a story's ID, in any padding;
// no other opener is one.
func TestAStorysPlannerIsPlannerAndAStorysID(t *testing.T) {
	for who, want := range map[string]bool{
		"planner-S-0328": true,
		"planner-s-1":    true,
		"planner-E-0016": false,
		"planner-T-0001": false,
		"planner-S-":     false,
		"planner":        false,
		"S-0328":         false,
		"orchestrator":   false,
		"":               false,
	} {
		if got := storyPlanner(who); got != want {
			t.Errorf("%q: %v, want %v", who, got, want)
		}
	}
}

// S-0220: ThreadOpener reads who opened a thread from the project's
// threads, in any padding of its ID, and fails for a thread there is not.
func TestThreadOpenerReadsTheProjectsThreads(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"design/system", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte("# Plan\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 6, 9, 0, 0, 0, time.UTC)
	for _, by := range []string{workitem.ActivityOrchestrator, "alex"} {
		if _, err := threads.New(r, threads.NewOptions{Title: "Which next?", On: "design/system/plan.md", Author: by, Text: "Which story next?", Now: now}); err != nil {
			t.Fatal(err)
		}
	}
	opener := ThreadOpener(r)
	for id, want := range map[string]string{"TH-0001": workitem.ActivityOrchestrator, "th-2": "alex"} {
		if got, err := opener(id); err != nil || got != want {
			t.Errorf("%s: %q, %v; want %q", id, got, err, want)
		}
	}
	if _, err := opener("TH-0003"); err == nil {
		t.Error("TH-0003: no error")
	}
}
