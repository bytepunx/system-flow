package mcpserver

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// orchestratorWith is the server of the orchestrator flai serve started with
// reasons in FLAI_REASONS, none when reasons is blank, its cursor set by a
// first inbox so that a wait starts with nothing behind it.
func orchestratorWith(t *testing.T, reasons string) (*fixture, *deadlines) {
	t.Helper()
	t.Setenv(harness.ReasonsEnv, reasons)
	f, d := endFixture(t, "", guard.RoleOrchestrate)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	return f, d
}

// ADR-0209 §4: reasons that arrive during a run start the next run, so the
// orchestrator flai serve started for reasons is told to end, with why
// reasons done, at once rather than held.
func TestWaitForEventsEndsTheOrchestratorWhoseReasonsAreDone(t *testing.T) {
	f, d := orchestratorWith(t, "thread TH-0001 awaits the orchestrator\nS-0002 moved to review")
	out := f.endsAtOnce(t, d)
	if out["why"] != "reasons done" {
		t.Errorf("why is reasons done: %v", out["why"])
	}
	if events, ok := out["events"].([]any); !ok || len(events) != 0 {
		t.Errorf("events is an empty list: %v", out["events"])
	}
	// the tool and the server's instructions say so
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range res.Tools {
		if tool.Name == "wait_for_events" && !strings.Contains(tool.Description, "why reasons done") {
			t.Errorf("wait_for_events's description does not say when the orchestrator ends: %s", tool.Description)
		}
	}
	if in := f.cs.InitializeResult().Instructions; !strings.Contains(in, "gets end with why reasons done") {
		t.Errorf("the instructions do not say when the orchestrator ends: %s", in)
	}
}

// Events already behind the orchestrator's cursor come with end, so none is
// lost.
func TestWaitForEventsReportsWhatIsBehindTheOrchestratorsCursorWithEnd(t *testing.T) {
	f, d := orchestratorWith(t, "ready-empty")
	*f.clock = t0.Add(10 * time.Minute)
	s := f.readyStory(t, "While it ran", t0.Add(5*time.Minute))
	out := f.endsAtOnce(t, d)
	if out["why"] != "reasons done" {
		t.Errorf("why is reasons done: %v", out["why"])
	}
	if got := changeSummaries(out, "events"); len(got) != 1 || !strings.Contains(got[0], s.ID) {
		t.Errorf("events behind the cursor come with end: %v", got)
	}
}

// An orchestrator with no reasons, as one run by hand, or with only blank
// ones, holds as before.
func TestWaitForEventsHoldsAnOrchestratorWithNoReasons(t *testing.T) {
	for name, reasons := range map[string]string{"empty": "", "blank": " \n\t"} {
		t.Run(name, func(t *testing.T) {
			f, d := orchestratorWith(t, reasons)
			f.holds(t, d)
		})
	}
}

// FLAI_REASONS is the orchestrator's: a story's agent with reasons in its
// environment and nothing pending holds, and with its question pending is
// told why in its own terms.
func TestWaitForEventsReasonsLeaveAStorysAgentAlone(t *testing.T) {
	t.Run("holds", func(t *testing.T) {
		t.Setenv(harness.ReasonsEnv, "ready-empty")
		f, d := storysAgent(t)
		if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
			t.Fatal(failed)
		}
		f.holds(t, d)
	})
	t.Run("ends on its question", func(t *testing.T) {
		t.Setenv(harness.ReasonsEnv, "ready-empty")
		f, d := storysAgent(t)
		th := f.ask(t, f.story.ID)
		why, _ := f.endsAtOnce(t, d)["why"].(string)
		if why == "reasons done" || !strings.Contains(why, th.ID) {
			t.Errorf("why names the story's question, not reasons: %q", why)
		}
	})
}
