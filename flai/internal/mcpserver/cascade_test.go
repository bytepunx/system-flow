package mcpserver

import (
	"testing"
	"time"
)

// S-0070: the designer cancels the epic the agent's story sits under. The
// agent learns of each cancellation with its cause, and its own cancellation
// through item_move reports what went with it.
func TestACascadedCancellationNamesItsCause(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	epic, _ := f.repo.Get("E-0001")
	if _, err := f.repo.TransitionAll(epic, "cancelled", "alex", "a different route", t0.Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	*f.clock = t0.Add(10 * time.Minute)
	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	got := changeSummaries(out, "changes")
	want := []string{
		"E-0001 Epic moved to cancelled by alex",
		"S-0001 Story was cancelled with E-0001 by alex",
		"T-0001 Task was cancelled with E-0001 by alex",
	}
	if len(got) != len(want) {
		t.Fatalf("changes: %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("change %d: %q, want %q", i, got[i], want[i])
		}
	}
	for _, c := range out["changes"].([]any) {
		m := c.(map[string]any)
		if m["id"] != "E-0001" && m["cause"] != "E-0001" {
			t.Errorf("%v has no cause", m)
		}
	}
}

func TestItemMoveToCancelledReportsWhatWentWithIt(t *testing.T) {
	f := setup(t)
	*f.clock = t0.Add(2 * time.Minute)
	out, failed := f.call(t, "item_move", map[string]any{"id": "S-0001", "to": "cancelled", "reason": "not needed"})
	if failed != "" {
		t.Fatal(failed)
	}
	with, _ := out["cancelled"].([]any)
	if len(with) != 1 || with[0].(map[string]any)["id"] != "T-0001" || with[0].(map[string]any)["from"] != "backlog" {
		t.Errorf("cancelled with it: %v", out["cancelled"])
	}
	task, _ := f.repo.Get("T-0001")
	if task.Status != "cancelled" {
		t.Errorf("T-0001 is %s", task.Status)
	}
}
