package mcpserver

import (
	"context"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/perf"
)

// waits are the tools that hold a call open by design; they are never slow.
var waits = map[string]bool{"wait_for_work": true, "wait_for_events": true}

// LongestWait caps one wait_for_events or wait_for_work call (I-0059): every
// return is a model turn that reads the agent's whole context again.
const LongestWait = 30 * time.Minute

const (
	eventsWait = time.Minute     // wait_for_events without timeout_seconds
	workWait   = 5 * time.Minute // wait_for_work without timeout_seconds
	heartbeat  = time.Minute     // how often a held wait tells the client it is alive
)

// keepAlive sends the client a progress notification on the call's progress
// token every beat until stop is called. Claude Code drops a call that sends
// nothing for 30 minutes over stdio and 5 over HTTP, whatever it asked for
// (I-0059), and a progress notification starts that count again. A call
// without a progress token gets nothing.
func keepAlive(ctx context.Context, req *mcp.CallToolRequest, beat time.Duration) (stop func()) {
	if req == nil || req.Session == nil || req.Params == nil || beat <= 0 {
		return func() {}
	}
	token := req.Params.GetProgressToken()
	if token == nil {
		return func() {}
	}
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		tick := time.NewTicker(beat)
		defer tick.Stop()
		for n := 1; ; n++ {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				_ = req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{ProgressToken: token, Progress: float64(n), Message: "still waiting"})
			}
		}
	}()
	return func() { cancel(); <-done }
}

// holdFor is how long a wait that asked for seconds is held: def when it
// did not ask, and never longer than longest.
func holdFor(seconds int, def, longest time.Duration) time.Duration {
	if seconds <= 0 {
		return min(def, longest)
	}
	if time.Duration(seconds) > longest/time.Second {
		return longest
	}
	return time.Duration(seconds) * time.Second
}

// timing times each request the server receives, beneath stdio or HTTP,
// and logs one "request answered" event for it (S-0152). A tool call is
// named by its tool.
func timing(logger *slog.Logger, now func() time.Time, slow time.Duration) mcp.Middleware {
	if now == nil {
		now = time.Now
	}
	if slow <= 0 {
		slow = perf.Slow()
	}
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			ctx, rec := perf.StartAt(ctx, now)
			res, err := next(ctx, method, req)
			a := perf.Answered{Transport: "mcp", Method: method}
			if p, ok := req.GetParams().(*mcp.CallToolParamsRaw); ok && p != nil {
				a.Method, a.Wait = p.Name, waits[p.Name]
			}
			if r, ok := res.(*mcp.CallToolResult); ok && r != nil {
				for _, c := range r.Content {
					if t, ok := c.(*mcp.TextContent); ok {
						a.Bytes += len(t.Text)
					}
				}
				if r.IsError && err == nil {
					a.Err = "the tool answered an error"
				}
			}
			if err != nil {
				a.Err = err.Error()
			}
			rec.Log(logger, a, slow)
			return res, err
		}
	}
}
