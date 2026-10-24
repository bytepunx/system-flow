package mcpserver

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// S-0431: an operator's own agent can have the curator started on the host,
// with a since or none, and the host judges whether it may.
func TestTheCuratorIsStartedThroughTheServer(t *testing.T) {
	type call struct{ root, since, by string }
	var got []call
	running := false
	f := setupWith(t, func(o *Options) {
		o.Curations = func(_ context.Context, root, since, by string) (CurateStarted, error) {
			got = append(got, call{root, since, by})
			if running {
				return CurateStarted{}, errors.New("the curator is already running for this project (since 2026-10-01, started by asked, pid 42, started 2026-10-23T10:00:00Z; log x.log); one curator runs per project at a time")
			}
			running = true
			return CurateStarted{Since: since, Agent: "curator", Harness: "claude-code", Command: "claude", PID: 42, Session: "abc", Started: "2026-10-23T10:00:00Z", Trigger: "asked"}, nil
		}
	})
	out, failed := f.call(t, "curate", map[string]any{"since": " 2026-10-01 "})
	if failed != "" || out["since"] != "2026-10-01" || out["agent"] != "curator" || out["pid"] != float64(42) || out["session"] != "abc" || out["trigger"] != "asked" {
		t.Fatalf("a since: %v %s", out, failed)
	}
	if _, failed := f.call(t, "curate", map[string]any{}); !strings.Contains(failed, "the curator is already running for this project") {
		t.Errorf("a second run: %q", failed)
	}
	running = false
	if out, failed := f.call(t, "curate", map[string]any{}); failed != "" || out["since"] != "" {
		t.Errorf("no since: %v %s", out, failed)
	}
	root := projectRoot(f.repo)
	if len(got) != 3 || got[0] != (call{root, "2026-10-01", "claude"}) || got[2].since != "" {
		t.Errorf("asked: %+v", got)
	}
}

func TestAServerThatCannotStartTheCuratorSaysHow(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "curate", map[string]any{"since": "2026-10-01"}); !strings.Contains(failed, "this flai mcp cannot start the curator; on the host run flai curate --since 2026-10-01") {
		t.Errorf("without a starter: %q", failed)
	}
}
