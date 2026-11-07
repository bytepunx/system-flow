---
id: TH-0842
title: "Plan of S-0569: four timing tests, four tasks in one layer, then a joint stress run and I-0247's close"
anchor:
  path: wip/kanban/stories/S-0569-testservesaystheflaiisolderthantheproject-and.md
  item: S-0569
status: resolved
participants: [planner-S-0569, orchestrator]
created: 2026-11-07T03:49:03Z
updated: 2026-11-07T03:53:56Z
---

# TH-0842 Plan of S-0569: four timing tests, four tasks in one layer, then a joint stress run and I-0247's close

On wip/kanban/stories/S-0569-testservesaystheflaiisolderthantheproject-and.md.

## Entries

### 2026-11-07T03:49:03Z planner-S-0569
Planned S-0569 with five tasks in two layers. The prompt asked me to read a thread named `asked`; no such thread exists and S-0569 had no thread or earlier plan, so I planned it from scratch.

## Proposed solution

I-0247 now has five instances, and they name four tests in three packages, each failing on a loaded host (load average 37 to 76) for its own reason:

| Test | Cause found in the code | Task |
|------|-------------------------|------|
| `TestWaitForEventsAcrossAFolder` | Races itself: sleeps 150 ms, then moves the story before the call may be holding | T-2522 |
| `TestServeSaysTheFlaiIsOlderThanTheProject` | Asks `project.info` once after the handshake and got an empty answer | T-2525 |
| `TestServeStartsADefinedRoleOnAnEventAndOnASchedule` | Waits for the right thing, but runs out of `waitFor`'s 30 s | T-2528 |
| `TestHostRunsServeAndTheMCPServersAndTakesThemWithIt` | Runs out of its 20 s `waitUntil` for a state file with a 10 s staleness check | T-2529 |

Each task reproduces its failure under `GOMAXPROCS=2 -race` beside a CPU load before it changes anything, as S-0548's T-2362 did for I-0226.

## Layers

- Layer 1, run together, no path in common: T-2522 (`flai/internal/mcpserver/folder_test.go`), T-2525 (`flai/internal/serve/outdated_test.go`), T-2528 (`serve_test.go`, `agents_test.go`), T-2529 (`flai/cmd/host_integration_test.go`).
- Layer 2: T-2533, after all four. It runs them together ten times at load 40 or more, runs the integration tier, and closes I-0247.

## Figures

- Forecast 3h, delivery 2026-11-07T12:30Z: raised from flai's 18m for the stress runs.
- Cost of delay 52.50 USD a week, from the inputs as set.

## Assumptions

- The fixes are in the tests. T-2525 and T-2529 say to fix flai's code instead, and widen touches, if the empty `project.info` or the late state file turns out to be a defect.
- Out of scope: capping the integration tier's parallelism (`-p` in `scripts/integration.sh`) or holding close-outs while the host is loaded. Either would cover the other `timeout_seconds: 3` tests in `flai/internal/mcpserver` too. Say if you want either added as a task.
- The cost of delay inputs (21m per cycle) came from three occurrences. I-0247 has five now. Changing them is your call; I left them as they are.

Source: design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md

### 2026-11-07T03:53:08Z orchestrator
**Plan approved, as drafted.** One task per failing test of I-0247, each in layer 1 with no path in common, then T-2533's joint stress run and the issue's close. That covers both of S-0569's criteria. The 3h forecast stands on the stress runs, as S-0548's T-2362 shows. The 52.50 USD a week value is flai's figure from the inputs as set.

**Not added:**
- Capping the integration tier's parallelism.
- Holding close-outs on a loaded host.

Either would widen the story's scope, and that is the operator's call. If they want it, it is a story of its own from an issue. The 21m input stays as flai set it; raising it for the two newer occurrences is also the operator's call.

Nothing in the plan waits on an answer, so I resolve this thread and judge the draft for finalizing.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T03:53:56Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories, as drafted; the out-of-scope options are left to the operator as a story of their own.
