---
id: T-0868
type: task
nature: remediation
title: wait_for_events and wait_for_work honour the requested timeout up to a 30-minute cap instead of five minutes
status: done
parent: S-0244
owner: alex
created: 2026-10-05T03:15:40Z
updated: 2026-10-05T05:01:19Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:45Z
    by: agent-S-0244
  - to: in-progress
    at: 2026-10-05T04:41:45Z
    by: agent-S-0244
  - to: done
    at: 2026-10-05T05:01:19Z
    by: agent-S-0244
stream: S-0244
tags: [flai, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/timing.go, flai/internal/mcpserver/timing_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go]
usage:
  source: log
  seconds: 931
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 103
      output: 36234
      cache_read: 4555050
      cache_write: 156678
      cost: 2.5572
---
# T-0868 wait_for_events and wait_for_work honour the requested timeout up to a 30-minute cap instead of five minutes

## Work

This is I-0059's remediation 3. The agents asked `wait_for_events` for 600 to 1800 s and got 300 s. Every return was a model turn that re-read the whole context, about 0.05 USD a tick at 150k cached tokens.

- Raise the cap from five minutes to 30 minutes. It lives in the `maxWait` default in `mcpserver/server.go` (line 108) and `folder.go` (line 190), and is applied in `work.go` and `server.go`. Name the cap once, as a constant both use.
- A requested `timeout_seconds` up to the cap is honoured. A request over the cap gets the cap. Leaving it out keeps a shorter default, so that a caller that never asked is not held for half an hour. Keep the default at 60 s for `wait_for_events`, as its schema says, and at 5 minutes for `wait_for_work`.
- Before settling the figure, find the longest call the claude-code harness's MCP client lets a tool hold, over stdio and over HTTP (`MCP_TOOL_TIMEOUT` and any client default). If that limit is shorter than 30 minutes, either take the cap down to it or note in the narrative that the harness must raise the limit. The harness task after this one sets it.
- Update the schema texts that state the cap (`WorkIn.TimeoutSeconds` in `work.go`, the `wait_for_events` description) and the comment on the HTTP server in `cmd/mcp_http.go`, which says "up to five minutes".

It waits for nothing and shares no path with the overlap task, so the two run together.

## Done when

- [ ] A Go test with the injected clock shows a `wait_for_events` call asking for 1800 s held past 300 s, and one asking for 3600 s cut to the cap; the same holds for `wait_for_work`
- [ ] A call without `timeout_seconds` keeps its default
- [ ] The narrative records the client's tool-call limit found, and whether the harness must raise it
- [ ] `scripts/flai-test.sh` passes

## Notes
