package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// lockedBuffer is a log destination a test reads while the server writes.
type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

// answered are the "request answered" events logged, by method.
func (l *lockedBuffer) answered(t *testing.T) map[string]map[string]any {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	out := map[string]map[string]any{}
	for _, line := range strings.Split(l.b.String(), "\n") {
		if !strings.Contains(line, `"request answered"`) {
			continue
		}
		var e map[string]any
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatal(err)
		}
		out[e["method"].(string)] = e
	}
	return out
}

func TestEachToolCallIsTimedBeneathTheTransport(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var logs lockedBuffer
	srv := New(Options{Repo: repo, Agent: "claude", Version: "test", MaxWait: 50 * time.Millisecond, Poll: 10 * time.Millisecond,
		Logger: slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})), Slow: time.Nanosecond})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test-agent", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	for name, args := range map[string]any{"who_touches": map[string]any{"path": "flai"}, "wait_for_events": map[string]any{"timeout_seconds": 1}} {
		if res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args}); err != nil || res.IsError {
			t.Fatalf("%s: %v %+v", name, err, res)
		}
	}
	got := logs.answered(t)
	who := got["who_touches"]
	if who == nil || who["level"] != "INFO" || who["transport"] != "mcp" || who["duration_ms"] == nil || who["bytes"] == nil || who["err"] != nil {
		t.Errorf("who_touches: %v", who)
	}
	// a wait is long by design and is never logged as slow
	if w := got["wait_for_events"]; w == nil || w["level"] != "DEBUG" {
		t.Errorf("wait_for_events: %v", w)
	}
	// a request that is not a tool call is named by its method
	if got["tools/call"] != nil || len(got) < 3 {
		t.Errorf("methods timed: %v", got)
	}
}
