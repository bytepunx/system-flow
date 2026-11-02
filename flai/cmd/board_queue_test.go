package cmd

import (
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/verify"
)

// S-0549: the queue line under the review column names the preflight's hold,
// its files each once, and its thread once flai serve has opened it
// (ADR-0234).
func TestQueueSummaryNamesThePreflightHold(t *testing.T) {
	q := &mergequeue.Queue{Eligible: []mergequeue.Entry{{Story: "S-0001"}}}
	if got, want := queueSummary(q), "queue: S-0001; batch: none"; got != want {
		t.Errorf("no hold: got %q, want %q", got, want)
	}
	q.Preflight = &mergequeue.PreflightHold{Findings: []verify.FileFinding{
		{Tier: "markdown", Path: "wip/kanban/tasks/T-0001-a.md", Line: 3, Rule: "MD051"},
		{Tier: "markdown", Path: "wip/kanban/tasks/T-0001-a.md", Line: 9, Rule: "MD051"},
		{Tier: "go-vet"},
	}}
	if got, want := queueSummary(q), "queue: S-0001; batch: none; held by the preflight on wip/kanban/tasks/T-0001-a.md, go-vet"; got != want {
		t.Errorf("held, no thread: got %q, want %q", got, want)
	}
	q.Preflight.Thread = "TH-0001"
	if got, want := queueSummary(q), "queue: S-0001; batch: none; held by the preflight on wip/kanban/tasks/T-0001-a.md, go-vet, thread TH-0001"; got != want {
		t.Errorf("held with its thread: got %q, want %q", got, want)
	}
}
