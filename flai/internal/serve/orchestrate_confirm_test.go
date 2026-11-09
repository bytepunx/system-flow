//go:build !windows

package serve

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// queuedLines are the lines of the reasons the lab's orchestrator has queued.
func (lab *orchestrateLab) queuedLines() []string {
	var out []string
	for _, r := range lab.o.Dir.reasonQueue(lab.root).Reasons {
		out = append(out, r.line())
	}
	return out
}

// I-0284, S-0622: a planner asks for a cost of delay input on a story in
// progress, the orchestrator recommends one, and the operator confirms it,
// or flai serve does after confirm_after. The confirmation must reach the
// orchestrator as a thread reason, so that it writes the input the operator
// confirmed rather than leaving it waiting for the operator to write.
func TestAConfirmedCostOfDelayRecommendationIsAThreadReasonForTheOrchestrator(t *testing.T) {
	const after = 10 * time.Second // orchestration.permissions.confirm_after, inside the coalescing window
	for _, tc := range []struct {
		name    string
		confirm func(lab *orchestrateLab, id string, recommended time.Time) time.Time
	}{
		{"by the operator", func(lab *orchestrateLab, id string, recommended time.Time) time.Time {
			at := recommended.Add(2 * time.Second)
			if _, err := threads.Confirm(lab.repo, id, "alex", at); err != nil {
				lab.t.Fatal(err)
			}
			return at
		}},
		{"by flai serve after confirm_after", func(lab *orchestrateLab, id string, recommended time.Time) time.Time {
			at := recommended.Add(after)
			if _, err := threads.ConfirmByTimeout(lab.repo, id, after, at); err != nil {
				lab.t.Fatal(err)
			}
			return at
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lab := newOrchestrateLab(t)
			lab.history()
			story, _ := lab.inProgress("Waits on its input")
			base := time.Now().UTC().Truncate(time.Second)
			first := lab.firstRun(base)

			planner := "planner-" + story
			title := story + " needs a cost of delay input"
			th, err := threads.New(lab.repo, threads.NewOptions{Title: title, On: story, Author: planner,
				Text: "No cost of delay inputs: what is `time_lost_per_cycle`?", Now: base.Add(time.Second)})
			if err != nil {
				t.Fatal(err)
			}
			lab.at = base.Add(2 * time.Second)
			lab.look()
			toOperator := "thread " + th.ID + " awaits the operator: " + title
			toOrchestrator := "thread " + th.ID + " awaits the orchestrator: " + title
			if got := lab.queuedLines(); !slices.Equal(got, []string{toOperator}) {
				t.Fatalf("the planner's question queued %q, want only %q", got, toOperator)
			}

			recommended := base.Add(3 * time.Second)
			if _, err := threads.ReplyWith(lab.repo, th.ID, workitem.ActivityOrchestrator,
				"**Recommendation:** set `time_lost_per_cycle: 20m` on "+story+".", recommended, threads.Marks{Recommendation: true}); err != nil {
				t.Fatal(err)
			}
			lab.at = recommended.Add(time.Second)
			lab.look()
			if got := lab.queuedLines(); slices.Contains(got, toOrchestrator) {
				t.Fatalf("the orchestrator's own recommendation, unconfirmed, queued %q", got)
			}

			confirmed := tc.confirm(lab, th.ID, recommended)
			lab.at = confirmed.Add(time.Second)
			lab.look()
			if r := lab.orchestrator(); r.Session != first.Session {
				t.Fatalf("started inside the coalescing window: %+v", r)
			}
			var reason *orchestrateReason
			for _, r := range lab.o.Dir.reasonQueue(lab.root).Reasons {
				if r.Kind == inbox.Thread && r.ID == th.ID && r.To == "you" {
					reason = &r
				}
			}
			if reason == nil {
				t.Fatalf("the confirmation queued no reason for the orchestrator: %q", lab.queuedLines())
			}
			if got := reason.line(); got != toOrchestrator {
				t.Errorf("the reason reads %q, want %q", got, toOrchestrator)
			}

			lab.at = base.Add(2*time.Second + 2*after + time.Minute)
			lab.look()
			run := lab.orchestrator()
			if run.Session == first.Session {
				t.Fatalf("no run started once the window passed: %+v", run)
			}
			if got := lab.reasonsGiven(run.PID); !slices.Contains(got, toOrchestrator) {
				t.Errorf("FLAI_REASONS =\n%s\nwant it to hold %q", strings.Join(got, "\n"), toOrchestrator)
			}
			lab.ended(run.Session)
		})
	}
}
