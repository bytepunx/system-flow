---
id: S-0583
type: story
nature: remediation
title: Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container
status: done
parent: E-0029
owner: alex
created: 2026-11-07T04:46:50Z
updated: 2026-11-07T07:41:03Z
transitions:
  - to: ready
    at: 2026-11-07T07:10:32Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T07:10:40Z
    by: agent-S-0583
  - to: review
    at: 2026-11-07T07:39:35Z
    by: agent-S-0583
  - to: done
    at: 2026-11-07T07:41:03Z
    by: orchestrator
tags: []
topics: [mcp, harness]
touches: [flai/internal/mcphttp/mcphttp.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/pilot/run_test.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, design/adrs, design/adrs/README.md, docs/operators/index.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1751
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 3
      hand_edits: 3
      work: 56
  models:
    - model: claude-opus-5-5
      input: 190
      output: 47347
      cache_read: 24441231
      cache_write: 362903
      cost: 8.6022
  strategic:
    - kind: planner
      seconds: 375
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 44
          output: 31956
          cache_read: 3744272
          cache_write: 227154
          cost: 3.2054
    - kind: orchestrator
      seconds: 76
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 32
          output: 246
          cache_read: 1161680
          cache_write: 149799
          cost: 0.3231
cost_of_delay:
  inputs:
    time_lost_per_cycle: 7m
    by: flai
    at: 2026-11-07T04:46:50Z
  value: 17.5
  by: planner-S-0583
  at: 2026-11-07T07:07:55Z
forecast:
  duration: 50m
  delivery: 2026-11-07T13:16:00Z
  basis: "flai forecast's 30m (median 105 s per unit of size over 55 done remediation stories on claude-opus-5-5, large band, size 17) plus 20m for T-2541's and T-2544's calls, which wait about 11m and 4m on the clock; delivery is flai's 12:56Z moved by those 20m."
  by: planner-S-0583
  at: 2026-11-07T07:07:55Z
finalized:
  by: alex
  at: 2026-11-07T07:02:21Z
---
# S-0583 Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container

## Goal

This story remediates [I-0228](../../../design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md), "Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container". The issue recommends this solution:

S-0506 fixes `wait_for_events` and `wait_for_work`: over HTTP the server holds one call at most 45 seconds, the bound `permission_prompt` holds there, and then answers `timed_out` for the agent to call again ([ADR-0244](../adrs/0244-over-http-one-wait-for-events-or-wait-for-work-call-holds-at-most-the-45.md)). `verify` has answered at once since S-0546 made the close-out a `flai serve` job, whose `verified` event the agent waits for within that bound ([ADR-0229](../adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)). Still open: a `task_done` or `story_start` slower than 60 seconds over HTTP still meets Claude Code's limit while the server finishes it. A per-server `timeout` in the MCP file S-0508's container adapter writes may lift that limit; it does not replace the bound.

## Acceptance criteria
- [x] The cause I-0228 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0228 is closed with `flai issue close I-0228 --reason` saying what fixed it

## Tasks
- T-2541 Measure how long a per-server timeout in the MCP file holds a call in the agent container, and record the figure in an ADR refining ADR-0244
- T-2542 The MCP file flai writes for its HTTP server, and the one flai mcp status prints, name the per-server timeout the ADR sets
- T-2543 The operator and user docs and the design say what the per-server timeout lifts, and that an agent image must be rebuilt to carry it
- T-2544 A flai test call slower than 60 seconds answers over HTTP under the new MCP file, as I-0228's third instance could not, and I-0228 is closed

## Notes

Cost of delay inputs set by flai from I-0228. time_lost_per_cycle 7m: 3m per occurrence × 3 occurrences ÷ 1.3 cycles of 168h (first reported 2026-10-29T00:56:36Z, 9.2 days before this story).

### Planning

**Touches.**

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/harness/adapters.go` | design, layout | `httpMCPConfig` writes the MCP file the Goal names |
| `flai/internal/mcphttp/mcphttp.go` | layout | Holds the one constant, beside `AgentHeader`, which `adapters.go` already imports |
| `flai/cmd/mcp_http.go` | layout | `mcpClientConfig` prints the sample a hand-written client copies |
| `flai/internal/harness/harness_test.go` | co-change, layout | `flai touches suggest`, 44%; it asserts the MCP file today |
| `flai/cmd/mcp_http_test.go` | co-change | `flai touches suggest`, 10% |
| `flai/internal/pilot/run_test.go` | layout | pilot writes the file in the container |
| `design/adrs`, `design/adrs/README.md` | co-change, design | ADR-0244 left the per-server `timeout` open; `flai touches suggest` lists the README, 8% |
| `docs/operators/index.md`, `docs/operators/runbooks/agent-container.md` | design | They state the limit and describe `mcp.json` |
| `docs/users/flai.md`, `design/system/flai-cli.md` | co-change, design | `flai touches suggest`, 15% and 12%; both state the limit |
| `design/system/agent-adapters.md` | design | It states the hold figures and the container's MCP |
| the I-0228 issue file, `design/issues/summary.md` | design | The second criterion closes the issue |

The story declared no touches. The suggestion's other paths are left out: `flai/internal/mcpserver` and `flai/cmd/mcp.go` change no hold, and `flai/internal/serve` and `docs/operators/settings.md` gain nothing, since the figure is a constant and no setting.

`design/adrs` is kept as a folder because the story adds one ADR whose number and slug `flai adr new` gives when it is written, so no task can name that file yet.

**After.** None. S-0506 and S-0508, whose output this story reads, are done.

**Forecast.** 50m. `flai forecast` gives 30m at size 17; 20m is added because the calls T-2541 and T-2544 make wait about 15m on the clock, which no size carries. The delivery is flai's 2026-11-07T12:56:00Z moved by those 20m.

**Cost of delay.** 17.50 USD a week, as `flai cod` works it out from the 7m input at 150 USD an hour. It stands: nothing in the plan changes that input.

### Accepted by the orchestrator

- Verified: b33657be64376df628785ff33b2bf343298a40fd
- At: 2026-11-07T07:41:03Z

Verdict: pass. Both criteria are met; flai verify S-0583 passed every step at b33657be6437, the branch head.
- 1: flai/internal/mcphttp/mcphttp.go, flai/internal/harness/adapters.go, flai/cmd/mcp_http.go, flai/internal/harness/harness_test.go, flai/internal/pilot/run_test.go, flai/cmd/mcp_http_test.go, design/adrs/0258-the-mcp-file-flai-writes-for-its-http-server-names-a-per-server-timeout-of.md
- 2: design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md, design/issues/summary.md
