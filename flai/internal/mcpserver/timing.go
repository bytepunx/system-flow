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
