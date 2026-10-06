package mcpserver

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// S-0223: an operator's own agent can have the analyzer started on the host,
// with a focus or none, and the host judges whether it may; a focus the
// analyzer does not take never reaches the host.
func TestTheAnalyzerIsStartedThroughTheServer(t *testing.T) {
	type call struct{ root, focus, by string }
	var got []call
	running := false
	f := setupWith(t, func(o *Options) {
		o.Analyses = func(_ context.Context, root, focus, by string) (AnalyzeStarted, error) {
			got = append(got, call{root, focus, by})
			if running {
				return AnalyzeStarted{}, errors.New("the analyzer is already running for this project (focus risk, started by asked, pid 42, started 2026-10-06T10:00:00Z; log x.log); one analyzer runs per project at a time")
			}
			running = true
			look := focus
			if look == "" {
				look = "all"
			}
			return AnalyzeStarted{Focus: look, Agent: "analyzer", Harness: "claude-code", Command: "claude", PID: 42, Session: "abc", Started: "2026-10-06T10:00:00Z", Trigger: "asked"}, nil
		}
	})
	out, failed := f.call(t, "analyze", map[string]any{"focus": "risk"})
	if failed != "" || out["focus"] != "risk" || out["agent"] != "analyzer" || out["pid"] != float64(42) || out["session"] != "abc" || out["trigger"] != "asked" {
		t.Fatalf("a focus: %v %s", out, failed)
	}
	if _, failed := f.call(t, "analyze", map[string]any{}); !strings.Contains(failed, "the analyzer is already running for this project") {
		t.Errorf("a second run: %q", failed)
	}
	running = false
	if out, failed := f.call(t, "analyze", map[string]any{}); failed != "" || out["focus"] != "all" {
		t.Errorf("no focus: %v %s", out, failed)
	}
	if _, failed := f.call(t, "analyze", map[string]any{"focus": "style"}); !strings.Contains(failed, `the analyzer takes the focus bottlenecks, intent, risk, or none for all of them, and "style" is none of them`) {
		t.Errorf("a bad focus: %q", failed)
	}
	root := projectRoot(f.repo)
	if len(got) != 3 || got[0] != (call{root, "risk", "claude"}) || got[2].focus != "" {
		t.Errorf("asked: %+v", got)
	}
}

func TestAServerThatCannotStartTheAnalyzerSaysHow(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "analyze", map[string]any{"focus": "intent"}); !strings.Contains(failed, "this flai mcp cannot start the analyzer; on the host run flai analyze --focus intent") {
		t.Errorf("without a starter: %q", failed)
	}
}
