package mcpserver

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// hostAccept stands in for the flai accept the tool runs on the host, which
// is flai's command and tested there: it records each command line and its
// standard input, and answers with ran.
type hostAccept struct {
	mu      sync.Mutex
	args    [][]string
	stdin   []string
	journal []hostapi.Entry
	ran     hostapi.Ran
}

func (h *hostAccept) run(_ context.Context, r hostapi.Run) (hostapi.Ran, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.args, h.stdin = append(h.args, r.Args), append(h.stdin, r.Stdin)
	return h.ran, nil
}

// acceptFixture is setup served to the caller in role, with accept_reviews
// as permit says in the manifest on disk, which the acceptance reads there,
// and in the server's, and the tool accept running host's flai.
func acceptFixture(t *testing.T, role string, permit bool) (*fixture, *hostAccept) {
	t.Helper()
	host := &hostAccept{ran: hostapi.Ran{Stdout: []byte(`{"id":"S-0001","dry_run":false}`)}}
	f := setupWith(t, func(o *Options) {
		o.Role, o.Agent = role, workitem.ActivityOrchestrator
		o.Accept = &hostapi.Accepter{Run: host.run, Now: func() time.Time { return t0 }, Host: hostapi.Host{
			Record: func(e hostapi.Entry) { host.mu.Lock(); host.journal = append(host.journal, e); host.mu.Unlock() },
		}}
	})
	if permit {
		path := filepath.Join(f.repo.Root, "system-flow.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, []byte("orchestration:\n  permissions:\n    accept_reviews: true\n")...), 0o644); err != nil {
			t.Fatal(err)
		}
		f.repo.Manifest.Orchestration.Permissions.AcceptReviews = true
	}
	return f, host
}

// ADR-0278 §1: accept runs flai accept --by orchestrator on the host, with
// the commit verified and the evidence on standard input, as accept.run
// builds the orchestrator's, journals it under mcp.accept, and answers
// flai's result.
func TestAcceptRunsTheOrchestratorsAcceptanceOnTheHost(t *testing.T) {
	t.Parallel()
	f, host := acceptFixture(t, guard.RoleOrchestrate, true)
	evidence := "Verdict: pass\n- 1: `flai/x.go`"
	out, failed := f.call(t, "accept", map[string]any{"story": "S-1", "verified": "abc1234", "evidence": evidence})
	if failed != "" {
		t.Fatal(failed)
	}
	want := []string{"accept", "S-0001", "--by=orchestrator", "--verified=abc1234", "--evidence=-", "--json"}
	if len(host.args) != 1 || !reflect.DeepEqual(host.args[0], want) || host.stdin[0] != evidence {
		t.Fatalf("flai ran %v with %q, want %v with the evidence", host.args, host.stdin, want)
	}
	if out["story"] != "S-0001" || out["result"].(map[string]any)["id"] != "S-0001" {
		t.Errorf("the answer is flai's result: %v", out)
	}
	if len(host.journal) != 1 || host.journal[0].Method != "mcp.accept" || host.journal[0].Outcome != "done" || host.journal[0].By != workitem.ActivityOrchestrator {
		t.Errorf("the acceptance is journalled under mcp.accept as the orchestrator's: %+v", host.journal)
	}
}

// ADR-0278 §1: with dry_run accept runs flai accept --by orchestrator
// --dry-run, verified and evidence optional, journals nothing, and answers
// the preview with its blockers.
func TestAcceptDryRunAnswersThePreview(t *testing.T) {
	t.Parallel()
	f, host := acceptFixture(t, guard.RoleOrchestrate, true)
	host.ran.Stdout = []byte(`{"id":"S-0001","dry_run":true,"blockers":[{"code":"unticked","message":"criterion 1 is not ticked"}]}`)
	out, failed := f.call(t, "accept", map[string]any{"story": "S-0001", "dry_run": true})
	if failed != "" {
		t.Fatal(failed)
	}
	want := []string{"accept", "S-0001", "--by=orchestrator", "--dry-run", "--json"}
	if len(host.args) != 1 || !reflect.DeepEqual(host.args[0], want) || host.stdin[0] != "" {
		t.Fatalf("flai ran %v, want %v", host.args, want)
	}
	if out["dry_run"] != true || len(out["result"].(map[string]any)["blockers"].([]any)) != 1 {
		t.Errorf("the answer is the preview and its blockers: %v", out)
	}
	if len(host.journal) != 0 {
		t.Errorf("a dry run changes nothing and is not journalled: %+v", host.journal)
	}
}

// ADR-0278 §1, ADR-0225: flai on the host reads the merge-queue host action,
// and while the queue batches refuses the orchestrator's acceptance, naming
// queue_add; the tool answers that refusal as its error, so the orchestrator
// enqueues the story instead.
func TestAcceptAnswersTheHostsRefusalWhileTheQueueBatches(t *testing.T) {
	t.Parallel()
	f, host := acceptFixture(t, guard.RoleOrchestrate, true)
	host.ran = hostapi.Ran{Exit: 1, Events: []map[string]any{{"level": "FATAL", "err": "rule: the merge-queue host action is on and the queue batches, so the merge queue lands S-0001, not the orchestrator: enqueue S-0001 with flai queue add S-0001 --by orchestrator --verified <commit> --evidence <file>, or the MCP tool queue_add"}}}
	_, failed := f.call(t, "accept", map[string]any{"story": "S-0001", "verified": "abc1234", "evidence": "Verdict: pass\n- 1: x"})
	if !strings.Contains(failed, "the merge-queue host action is on") || !strings.Contains(failed, "queue_add") {
		t.Errorf("the refusal is the host's flai's, naming queue_add: %q", failed)
	}
	if len(host.journal) != 1 || host.journal[0].Outcome != "failed" {
		t.Errorf("the refused acceptance is journalled as failed: %+v", host.journal)
	}
}

// ADR-0278 §2: accept is refused, running nothing, to anyone but the
// orchestrator, while accept_reviews is off, and for a verified that is not
// a commit's name.
func TestAcceptIsRefusedOutsideTheOrchestratorsPermission(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		why    string
		role   string
		permit bool
		args   map[string]any
		want   string
	}{
		{"a story's agent", "", true, map[string]any{"story": "S-0001", "verified": "abc1234", "evidence": "Verdict: pass"}, "accept is the orchestrator's"},
		{"the planner", guard.RolePlan, true, map[string]any{"story": "S-0001", "dry_run": true}, "accept is the orchestrator's"},
		{"accept_reviews off", guard.RoleOrchestrate, false, map[string]any{"story": "S-0001", "verified": "abc1234", "evidence": "Verdict: pass"}, "accept_reviews"},
		{"a verified that is not a commit", guard.RoleOrchestrate, true, map[string]any{"story": "S-0001", "verified": "HEAD; rm", "evidence": "Verdict: pass"}, "is not one"},
		{"a dry run's verified that is not a commit", guard.RoleOrchestrate, true, map[string]any{"story": "S-0001", "verified": "--yes", "dry_run": true}, "is not one"},
	} {
		f, host := acceptFixture(t, c.role, c.permit)
		_, failed := f.call(t, "accept", c.args)
		if !strings.Contains(failed, c.want) {
			t.Errorf("%s: refused with %q, want %q", c.why, failed, c.want)
		}
		if len(host.args) != 0 || len(host.journal) != 0 {
			t.Errorf("%s: flai ran %v and journalled %+v", c.why, host.args, host.journal)
		}
	}
}
