package mcpserver

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// stories reads the ids, in order, of a list of ranked stories in a tool's
// answer.
func stories(v any) (ids []string) {
	list, _ := v.([]any)
	for _, s := range list {
		m, _ := s.(map[string]any)
		id, _ := m["id"].(string)
		ids = append(ids, id)
	}
	return ids
}

// S-0217: order_by_policy computes the ready column's order by a policy, the
// project's when none is given, and writes nothing.
func TestOrderByPolicyOrdersTheReadyColumn(t *testing.T) {
	f := setup(t)
	bare := f.readyStory(t, "Bare", t0)
	valued := f.readyStory(t, "Valued", t0)
	if _, failed := f.call(t, "item_edit", map[string]any{"id": valued.ID, "cost_of_delay": map[string]any{"value": 300}, "forecast": map[string]any{"duration": "3h"}}); failed != "" {
		t.Fatal(failed)
	}
	board := filepath.Join(f.repo.KanbanDir(), workitem.BoardFile)
	before, _ := os.ReadFile(board)

	out, failed := f.call(t, "order_by_policy", map[string]any{"policy": "wsjf"})
	if failed != "" {
		t.Fatal(failed)
	}
	got := stories(out["stories"])
	if out["policy"] != "wsjf" || out["applied"] != false || strings.Join(got, " ") != valued.ID+" "+bare.ID {
		t.Fatalf("by wsjf, the valued story first and the bare one after: %v", out)
	}
	list := out["stories"].([]any)
	if first := list[0].(map[string]any); first["figure"] != 100.0 || first["unit"] != "per hour" {
		t.Errorf("wsjf is 300 over 3 hours: %v", first)
	}
	if second := list[1].(map[string]any); second["missing"] != "cost of delay value and forecast duration" {
		t.Errorf("the bare story names what it lacks: %v", second)
	}
	if out, _ := f.call(t, "order_by_policy", map[string]any{}); out["policy"] != "fifo" || strings.Join(stories(out["stories"]), " ") != bare.ID+" "+valued.ID {
		t.Errorf("no policy is the project's, fifo when it sets none: %v", out)
	}
	if _, failed := f.call(t, "order_by_policy", map[string]any{"policy": "random"}); !strings.Contains(failed, `unknown order policy "random"`) {
		t.Errorf("an unknown policy is refused: %q", failed)
	}
	if after, _ := os.ReadFile(board); string(after) != string(before) {
		t.Error("order_by_policy wrote board.md")
	}
}

// S-0217: promote_candidates lists the backlog stories that could go to
// ready and why each other one cannot, capped by limit.
func TestPromoteCandidatesListsTheBacklog(t *testing.T) {
	f := setup(t)
	backlog := func(title string, ready bool) *workitem.Item {
		t.Helper()
		s, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: f.story.Parent, Owner: "alex", Touches: []string{"docs/" + strings.ToLower(title)}, Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		if ready {
			data, _ := os.ReadFile(s.Path)
			body := strings.Replace(string(data), "## Goal\n", "## Goal\n\nDo it.\n", 1)
			body = strings.Replace(body, "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)
			_ = os.WriteFile(s.Path, []byte(body), 0o644)
			if _, failed := f.call(t, "item_edit", map[string]any{"id": s.ID, "cost_of_delay": map[string]any{"value": 50}, "forecast": map[string]any{"duration": "1h"}}); failed != "" {
				t.Fatal(failed)
			}
		}
		return s
	}
	one, two, bare := backlog("One", true), backlog("Two", true), backlog("Bare", false)

	out, failed := f.call(t, "promote_candidates", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["policy"] != "fifo" || strings.Join(stories(out["candidates"]), " ") != one.ID+" "+two.ID {
		t.Errorf("the candidates by fifo: %v", out)
	}
	others, _ := out["others"].([]any)
	if len(others) != 1 || others[0].(map[string]any)["id"] != bare.ID || !strings.Contains(strings.Join(toStrings(others[0].(map[string]any)["reasons"]), "; "), "no goal") {
		t.Errorf("the bare story with its reasons: %v", out["others"])
	}
	if out, _ := f.call(t, "promote_candidates", map[string]any{"limit": 1}); strings.Join(stories(out["candidates"]), " ") != one.ID {
		t.Errorf("limit 1 lists one candidate: %v", out)
	}
	if _, failed := f.call(t, "promote_candidates", map[string]any{"limit": -1}); !strings.Contains(failed, "limit is a number of candidates") {
		t.Errorf("a negative limit is refused: %q", failed)
	}
}

// S-0217: release_evaluate says whether the release policy is met, reading
// git for what is released, and says so when it cannot read git.
func TestReleaseEvaluateWeighsWhatIsPending(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	if _, failed := setup(t).call(t, "release_evaluate", map[string]any{}); !strings.Contains(failed, "on the host run flai release --evaluate") {
		t.Errorf("without a runner it says how to on the host: %q", failed)
	}

	f := setupWith(t, func(o *Options) { o.Runner = execx.System{} })
	root := f.repo.Root
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"commit", "-q", "-m", "init"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	out, failed := f.call(t, "release_evaluate", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["policy"] != "judgement" || out["met"] != false || out["count"] != 0.0 || out["reason"] == "" {
		t.Errorf("judgement, with nothing pending, is not met: %v", out)
	}
	count := 1
	f.repo.Manifest.Orchestration.Release.Policy, f.repo.Manifest.Orchestration.Release.Count = "threshold", &count
	if out, failed := f.call(t, "release_evaluate", map[string]any{}); failed != "" || out["policy"] != "threshold" || out["met"] != false || out["count_threshold"] != 1.0 ||
		!strings.Contains(out["reason"].(string), "no accepted story is waiting for a release") {
		t.Errorf("a threshold is not met by nothing pending: %v %s", out, failed)
	}
}
