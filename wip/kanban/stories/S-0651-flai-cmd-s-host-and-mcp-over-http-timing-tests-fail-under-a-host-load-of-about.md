---
id: S-0651
type: story
nature: improvement
title: flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200
status: ready
owner: alex
created: 2026-11-10T05:44:44Z
updated: 2026-11-10T07:12:18Z
transitions:
  - to: ready
    at: 2026-11-10T07:08:50Z
    by: orchestrator
tags: [cli]
topics: [testing]
touches: [flai/cmd/mcp_http_test.go, flai/internal/channel/channel_test.go, design/issues/I-0309-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 31
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 41
          cache_read: 279469
          cache_write: 5647
          cost: 0.0767
cost_of_delay:
  inputs:
    time_lost_per_cycle: 56m
    by: flai
    at: 2026-11-10T05:44:44Z
  value: 70
  by: planner-S-0651
  at: 2026-11-10T06:59:13Z
forecast:
  duration: 1h
  delivery: 2026-11-10T09:12:00Z
  basis: "Its own forecast of 1h; 3rd in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647, S-0649, S-0650 and S-0655."
  by: flai
  at: 2026-11-10T07:12:18Z
finalized:
  by: orchestrator
  at: 2026-11-10T07:08:13Z
---
# S-0651 flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200

## Goal

This story remediates [I-0309](../../../design/issues/I-0309-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md), "flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0309 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0309 is closed with `flai issue close I-0309 --reason` saying what fixed it

## Tasks
- T-2839 TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp gives its client a deadline far above the server's bound
- T-2842 TestBacksOffWhenTheDashboardSaysTheProjectIsHeld tells an ordinary close from the held back-off by a wide gap
- T-2846 Confirm I-0309's tests under load and close I-0309

## Notes

Cost of delay inputs set by flai from I-0309. time_lost_per_cycle 56m: 14m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-11-08T00:27:14Z, 2.2 days before this story; under one cycle counts as one).

### Planning

What is left of I-0309 on main, read from its four instances and the code:

- The two import tests, `TestImportWithAHostRegistersTheProject` and `TestImportWithServeStoppedSaysToStartIt`, are fixed on main already: S-0606's T-2709 stamps `fakeRunning`'s state an hour ahead (`flai/cmd/import_serve_test.go`), so it outlasts `host.StaleAfter`. Nothing to do there but confirm them under load.
- `TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp` (`flai/cmd/mcp_http_test.go`) scales both the server's bound and the client's deadline down by 20: the server holds at most 2.25 s and the client gives up at 3 s, a 0.75 s margin that a loaded host eats. Proposed fix: keep the static check that `HTTPLongestWait` is under Claude Code's 60 s, and test the behaviour with a short server bound and a client deadline far above it but far below the 600 s asked, so that only a wait held for the 600 s fails.
- `TestBacksOffWhenTheDashboardSaysTheProjectIsHeld` (`flai/internal/channel/channel_test.go`), from the fourth instance, tells an ordinary close's retry, at most 40 ms, from the held back-off, at least 300 ms, by a 250 ms bound; under load an ordinary close waited 252 ms. Proposed fix: widen the gap, a longer `HeldBackoff` and a bound on the ordinary close well under its lower half. It is outside flai/cmd, but I-0309 records it, and closing the issue without it would lose the instance.

Tasks, in two layers: T-2839 and T-2842 first, together, since they change one file each and share none; then T-2846, which runs all four tests under load and closes the issue with a reason naming both fixes.

Touches, file by file:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/mcp_http_test.go` | design | the issue's first and fourth instances name the test at line 357 |
| `flai/internal/channel/channel_test.go` | design | the issue's fourth instance names the test at line 297 |
| `design/issues/I-0309-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md` | design | the second criterion closes it |
| `design/issues/summary.md` | co-change (4 of 25 commits) | `flai issue close` regenerates it |

Left out from `flai touches suggest`: `flai/internal/channel/channel.go` (48%) and `flai/cmd/mcp_http.go` (28%), since both fixes are in the tests and the code's timing is right; the rest are under 20% and not reached. No folder touch.

Forecast: `flai forecast` gave 3m from size 2 with no touches, and 8m from size 6 once the touches were set, at 72 s per unit over 25 done improvement stories. Both are below one close-out run. Set to 1h: two test edits of about 15 minutes, a reproduction under load of about 15 minutes, and a close-out whose integration tier ran 12 to 14 minutes in each instance. flai's delivery follows from that duration and the pull order.

Cost of delay: `flai cod` gives 140 USD a week from 56m per cycle at 150 USD an hour. Set to 70 USD a week: two of the four occurrences (2026-11-08T01:32:32Z and 2026-11-09T08:25:21Z) failed only on the import tests S-0606 has since fixed, so the time still lost is half the inputs'. The inputs stay as flai set them.

No `after`: the fix reads nothing another open story writes.
