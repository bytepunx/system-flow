---
id: S-0560
type: story
nature: remediation
title: TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80
status: in-progress
owner: alex
created: 2026-11-02T03:55:02Z
updated: 2026-11-07T08:34:41Z
transitions:
  - to: ready
    at: 2026-11-07T08:22:47Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T08:34:41Z
    by: agent-S-0560
tags: [flai, tests]
topics: [testing]
touches: [flai/internal/mcpserver/work.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/timing_test.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/orchestrator_threads_test.go, flai/internal/serve/candidates.go, flai/internal/serve/candidates_test.go, flai/internal/hostapi/import.go, design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md, design/issues/summary.md]
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
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 47
          cache_read: 374104
          cache_write: 6602
          cost: 0.0941
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-11-02T03:55:02Z
  value: 15
  by: planner-S-0560
  at: 2026-11-07T08:21:35Z
forecast:
  duration: 1h30m
  delivery: 2026-11-07T14:26:00Z
  basis: "flai forecast's 26m from size 15 raised to 1h30m for two product fixes, 19 tests moved, and stress runs, as S-0569's five comparable tasks took 1h18m; delivery is flai's 13:22Z plus the 1h04m added."
  by: planner-S-0560
  at: 2026-11-07T08:21:35Z
finalized:
  by: orchestrator
  at: 2026-11-07T08:22:33Z
---
# S-0560 TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80

## Goal

This story remediates [I-0246](../../../design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md), "TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0246 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0246 is closed with `flai issue close I-0246 --reason` saying what fixed it

## Tasks
- T-2555 wait_for_work and the folder's waits take their snapshot before their first check, so a change made between the two wakes the wait
- T-2556 An import.run answer reaches the dashboard before flai serve halts the candidate's connection, so TestARepositoryIsOfferedImportedAndThenServed gets no EOF
- T-2557 The mcpserver tests that change a project under a held wait wait for its heartbeat instead of sleeping 150 ms, and hold up to 30 s
- T-2558 The tests I-0246 names pass together under -race -count=10 with GOMAXPROCS=2 beside a CPU load, and I-0246 is closed with what fixed each

## Notes

Cost of delay inputs set by flai from I-0246. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T00:55:15Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Touches, none declared before this plan:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/mcpserver/work.go`, `flai/internal/mcpserver/folder.go` | co-change and layout | `flai touches suggest` from the tests the instances name lists both; each takes its snapshot after its first check |
| `flai/internal/mcpserver/timing_test.go` | layout | Holds the wait timing tests, so it takes the new reproducing test |
| `flai/internal/mcpserver/work_test.go`, `flai/internal/mcpserver/folder_test.go` | I-0246's instances | They name these files; `folder_test.go` holds `heldWait` |
| `flai/internal/mcpserver/server_test.go` | co-change and layout | The fixture the held tests use, and seven of them |
| `flai/internal/mcpserver/messages_test.go`, `flai/internal/mcpserver/orchestrator_threads_test.go` | layout | Five more tests with the same 150 ms sleep |
| `flai/internal/serve/candidates_test.go` | I-0246's instances | S-0555's EOF |
| `flai/internal/serve/candidates.go`, `flai/internal/hostapi/import.go` | layout | Where a candidate is halted and where `imported` is called |
| `design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md`, `design/issues/summary.md` | goal and criteria | Closing the issue writes both |

Left out on purpose: `flai/internal/mcpserver/server.go`, whose `waitForEvents` already takes its snapshot first; `flai/internal/serve/serve_test.go`, which S-0569 fixed; and the docs, since no command or setting changes.

Forecast: `flai forecast` gave 26m from size 15. It is raised to 1h30m for two product fixes, 19 tests moved, and stress runs of several minutes each, as S-0569's five comparable tasks took 1h18m of agent work. The delivery is flai's 13:22Z plus the 1h04m added.

Cost of delay: the value stands at `flai cod`'s 15.00 USD a week, from the 6m input. That input counts one occurrence, and I-0246 now has nine. The input is the operator's, so it stands, and TH-0858 raises it.

No `after`: this story reads nothing another open story writes. S-0569's fixes are already on main.
