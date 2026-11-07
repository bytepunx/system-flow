---
id: S-0569
type: story
nature: remediation
title: TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve
status: done
owner: alex
created: 2026-11-02T21:41:39Z
updated: 2026-11-07T05:34:34Z
transitions:
  - to: ready
    at: 2026-11-07T03:56:09Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T04:15:08Z
    by: agent-S-0569
  - to: review
    at: 2026-11-07T05:32:57Z
    by: agent-S-0569
  - to: done
    at: 2026-11-07T05:34:34Z
    by: orchestrator
tags: [flai, tests]
topics: [testing]
touches: [flai/internal/serve/outdated_test.go, flai/internal/serve/serve_test.go, flai/internal/serve/agents_test.go, flai/internal/mcpserver/folder_test.go, flai/cmd/host_integration_test.go, design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md, design/issues/summary.md, design/issues/I-0271-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md, design/issues/I-0273-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4688
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      hand_edits: 4
      work: 32
  models:
    - model: claude-opus-5-5
      input: 358
      output: 94688
      cache_read: 44726735
      cache_write: 1504297
      cost: 19.2249
  strategic:
    - kind: planner
      seconds: 319
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 80
          output: 248
          cache_read: 2252864
          cache_write: 252769
          cost: 0.1658
        - model: claude-opus-5-5
          input: 408
          output: 30015
          cache_read: 23111255
          cache_write: 823148
          cost: 7.7876
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 13
          cache_read: 167036
          cache_write: 1708
          cost: 0.041
cost_of_delay:
  inputs:
    time_lost_per_cycle: 21m
    by: flai
    at: 2026-11-02T21:41:39Z
  value: 52.5
  by: planner-S-0569
  at: 2026-11-07T03:48:36Z
forecast:
  duration: 3h
  delivery: 2026-11-07T12:30:00Z
  basis: "flai forecast's 18m from size 9 raised to 3h for four tests to reproduce and time under GOMAXPROCS=2 -race stress runs of several minutes each, then a joint run and the integration tier, as S-0548's T-2362 took 42m for three such tests."
  by: planner-S-0569
  at: 2026-11-07T03:48:36Z
finalized:
  by: orchestrator
  at: 2026-11-07T03:54:12Z
---
# S-0569 TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve

## Goal

This story remediates [I-0247](../../../design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md), "TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0247 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0247 is closed with `flai issue close I-0247 --reason` saying what fixed it

## Tasks
- T-2522 TestWaitForEventsAcrossAFolder makes its move only once the wait_for_events call is holding, and passes under -race -count=10 with GOMAXPROCS=2
- T-2525 TestServeSaysTheFlaiIsOlderThanTheProject finds why project.info answered nil on a loaded host and waits for an answer that holds the project
- T-2528 TestServeStartsADefinedRoleOnAnEventAndOnASchedule waits for the watcher's run for the move within a bound a host at load 76 does not reach
- T-2529 TestHostRunsServeAndTheMCPServersAndTakesThemWithIt waits for the MCP server's state within a bound that outlasts its staleness check on a loaded host
- T-2533 The four I-0247 tests pass together under -race -count=10 with GOMAXPROCS=2 beside a CPU load, and I-0247 is closed with what fixed each

## Notes

Cost of delay inputs set by flai from I-0247. time_lost_per_cycle 21m: 7m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-02T01:00:19Z, 0.9 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0247's five instances: four tests in three packages fail on a loaded host (load average 37 to 76), each for its own reason, so each gets its own task. `TestWaitForEventsAcrossAFolder` races itself: it sleeps 150 ms and moves the story before the call may be holding (T-2522). `TestServeSaysTheFlaiIsOlderThanTheProject` asks `project.info` once, right after the handshake, and got an empty answer (T-2525). `TestServeStartsADefinedRoleOnAnEventAndOnASchedule` waits correctly but runs out of `waitFor`'s 30 s (T-2528). `TestHostRunsServeAndTheMCPServersAndTakesThemWithIt` runs out of its 20 s `waitUntil` for a state file with a 10 s staleness check (T-2529). Each task reproduces its failure under `GOMAXPROCS=2 -race` beside a CPU load before it changes anything, as S-0548's T-2362 did for I-0226. T-2533 runs the four together and closes the issue.

Touches: the five test files come from the issue's instances and the explorer's read of the code (layout). `flai touches suggest` found no declared touches to start from, then, given them, listed `flai/internal/mcpserver/server.go`, `folder.go`, and `flai/internal/serve/agents.go` as co-changes; they are left out, since the plan changes tests, and T-2522 or T-2525 widens its touches if a test hook or a defect fix in flai's code turns out to be needed. The issue file and `design/issues/summary.md` are what `flai issue close` writes (layout). No folder touch.

Forecast: `flai forecast` gave 18m, 115 s per unit of size over 37 remediation stories at size 9. Raised to 3h: size counts criteria and touches, not the stress runs, each of several minutes under `-race`, that four tasks need before and after their change, nor the joint run and the integration tier in T-2533. S-0548's T-2362 took 42m for three such tests. Delivery is flai's 09:39 moved by the same 2h42m.

Cost of delay: `flai cod` gives 52.50 USD a week from the 21m lost per 168h cycle at 150 USD an hour. The value stands. The inputs were set from three occurrences; I-0247 now has five, which would raise them, and that is the operator's to change.

After: none. S-0569 reads no other story's output.

### Accepted by the orchestrator

- Verified: 55c53a18b25071f1534b06c29a55b15afa2f9b9f
- At: 2026-11-07T05:34:34Z

Verdict: pass — both criteria are met by test-only race fixes inside the declared touches, I-0247 is closed with a per-test remediation, and the two new issues are recorded findings; verified at 55c53a18b25071f1534b06c29a55b15afa2f9b9f.
- 1: flai/internal/mcpserver/folder_test.go, flai/internal/serve/outdated_test.go, flai/internal/serve/serve_test.go, flai/cmd/host_integration_test.go
- 2: design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md, design/issues/summary.md
