---
id: T-0908
type: task
nature: feature
title: A pending recommendation is not an answer to the agent that asked, and the designer's inbox lists it to confirm
status: done
parent: S-0220
owner: alex
created: 2026-10-05T04:48:02Z
updated: 2026-10-06T06:40:44Z
transitions:
  - to: ready
    at: 2026-10-06T06:19:41Z
    by: agent-S-0220
  - to: in-progress
    at: 2026-10-06T06:19:42Z
    by: agent-S-0220
  - to: done
    at: 2026-10-06T06:40:44Z
    by: agent-S-0220
stream: S-0220
tags: [flai]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/hostapi/people.go, flai/internal/hostapi/people_test.go]
after: [T-0902]
usage:
  source: log
  seconds: 1180
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 91
      output: 34469
      cache_read: 4266666
      cache_write: 133015
      cost: 2.3068
---
# T-0908 A pending recommendation is not an answer to the agent that asked, and the designer's inbox lists it to confirm

## Work

Three readers take the last entry by someone else as the answer. Each must treat a pending recommendation, which T-0894 marks, as no answer, and an autonomous answer as one, so that with `autonomous` the agent goes on and with `recommend` it waits for the operator:

- The MCP inbox's `awaiting` (`flai/internal/mcpserver/server.go`, lines 187-192) and the change it reports to the opener (`flai/internal/mcpserver/work.go`, line 83): a thread whose last entry is a pending recommendation stays `awaiting: other` for the agent that asked.
- `answered` in `flai/internal/serve/agents.go` (line 1047), which decides when flai serve starts a story's agent again after it ended asking: a pending recommendation does not start it; the operator's confirmation, or an answer, does.
- The designer's inbox, `inbox.designer` (`flai/internal/hostapi/people.go`, line 260): a thread with a pending recommendation is listed as awaiting the designer, with the recommendation's text and source, so that the dashboard can offer to confirm it.

It waits for T-0902, which also changes `flai/internal/mcpserver/server.go`. It runs with T-0910 and T-0912, whose paths it does not share.

## Done when

- an MCP test finds a thread with a pending recommendation awaiting `other` for its opener, and awaiting `you` once the operator confirms it
- a serve test finds a run that ended asking not started again on a recommendation, and started again on its confirmation and on an orchestrator's answer
- a host API test finds the recommendation and its source in `inbox.designer`
- `go test ./internal/mcpserver/ ./internal/serve/ ./internal/hostapi/` passes

## Notes
