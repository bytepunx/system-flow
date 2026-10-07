//go:build !windows

package serve

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// agentSession is a story's agent's flai mcp, in process: a fresh server on
// the lab's project, as the agent's harness starts one with each session.
// armed says whether a held wait_for_events armed its deadline, that is,
// whether it held rather than answered at once.
type agentSession struct {
	cs    *mcp.ClientSession
	armed *atomic.Bool
}

// session starts agent's flai mcp on the lab's project. Its deadline fires
// as soon as it is armed, so a wait that holds times out at once.
func (lab *agentLab) session(agent string) *agentSession {
	lab.t.Helper()
	repo, err := workitem.Open(lab.root)
	if err != nil {
		lab.t.Fatal(err)
	}
	armed := &atomic.Bool{}
	srv := mcpserver.New(mcpserver.Options{Repo: repo, Agent: agent, Version: "test", Poll: 20 * time.Millisecond, MaxWait: time.Minute,
		After: func(time.Duration) <-chan time.Time {
			armed.Store(true)
			c := make(chan time.Time, 1)
			c <- time.Now()
			return c
		}})
	ctx, cancel := context.WithCancel(context.Background())
	lab.t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		lab.t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: agent, Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		lab.t.Fatal(err)
	}
	lab.t.Cleanup(func() { _ = cs.Close() })
	return &agentSession{cs: cs, armed: armed}
}

// call invokes a tool and decodes its structured result into out.
func (a *agentSession) call(t *testing.T, name string, args, out any) {
	t.Helper()
	res, err := a.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if res.IsError {
		t.Fatalf("%s failed: %+v", name, res.Content)
	}
	data, _ := json.Marshal(res.StructuredContent)
	if err := json.Unmarshal(data, out); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
}

// S-0272: an agent that flai serve started, holding wait_for_events with its
// own question open on its story and no task in progress, is told to end at
// once; it ends asking, and serve starts it again, as itself in its session,
// when the question is answered, with the answer in its first inbox, and
// its wait_for_events holds again.
func TestAnAgentThatWaitForEventsEndsOnAQuestionIsStartedAgainOnTheAnswerWithItInItsFirstInbox(t *testing.T) {
	for _, env := range []string{"FLAI_STORY", "FLAI_AGENT", "FLAI_ROLE", "FLAI_SESSION", "FLAI_ANSWERED", "FLAI_STARTED_BY"} {
		t.Setenv(env, "")
	}
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	id := lab.ready("Asks through wait_for_events")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: first.Agent, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	// the agent, as serve starts it, holds wait_for_events: it is told to end
	t.Setenv("FLAI_STORY", id)
	t.Setenv("FLAI_AGENT", first.Agent)
	asking := lab.session(first.Agent)
	var told mcpserver.WaitOut
	began := time.Now()
	asking.call(t, "wait_for_events", map[string]any{"timeout_seconds": 60}, &told)
	if !told.End || !strings.Contains(told.Why, th.ID) || !strings.Contains(told.Why, id) {
		t.Fatalf("told to end, naming the question: %+v", told)
	}
	if asking.armed.Load() || told.TimedOut || time.Since(began) > 5*time.Second {
		t.Fatalf("it held rather than answered at once (%s): %+v", time.Since(began), told)
	}

	// it ends, asking; nothing starts it again while the question is open
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	if r := lab.run(id); r.Outcome != OutcomeAsked || r.Thread != th.ID {
		t.Fatalf("it ended asking: %+v", r)
	}
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	lab.l.look(ctx, false)
	if lab.run(id).live() {
		t.Fatal("started again before the answer")
	}

	// the designer answers: serve starts the same agent again in its session
	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again", func() bool { return lab.run(id).live() })
	again := lab.run(id)
	if again.Agent != first.Agent || again.Session != first.Session || again.Answered != th.ID {
		t.Fatalf("the same agent, in its session, for the answer: %+v, first %+v", again, first)
	}

	// its first inbox, in a session of its own, holds the answer
	answered := lab.session(again.Agent)
	var in mcpserver.InboxOut
	answered.call(t, "inbox", map[string]any{}, &in)
	var got *mcpserver.ThreadSummary
	for i := range in.Threads {
		if in.Threads[i].ID == th.ID {
			got = &in.Threads[i]
		}
	}
	if got == nil || got.Awaiting != "you" || got.LastBy != "alex" || got.LastEntry != "Nine." || got.Story != id || in.AwaitingYou < 1 {
		t.Fatalf("the answer in its first inbox: %+v in %+v", got, in)
	}
	// answered, the question no longer ends its wait: it holds
	var held mcpserver.WaitOut
	answered.call(t, "wait_for_events", map[string]any{"timeout_seconds": 60}, &held)
	if held.End || held.Why != "" || !held.TimedOut || !answered.armed.Load() {
		t.Fatalf("told to end with its question answered: %+v", held)
	}

	lab.release(id)
	waitFor(t, "it ends again", func() bool { return !lab.run(id).live() })
}
