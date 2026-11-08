---
id: S-0611
type: story
nature: improvement
title: The dashboard asks threads.list for open threads and agent.status for open runs, and the request timing counts the send, so a thread list is under 100 KB, not 3.3 MB
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:42:08Z
updated: 2026-11-08T01:49:10Z
transitions: []
tags: [dashboard, cli]
topics: [logging]
touches: [flaiover/src/lib/server/repo.ts, flai/internal/channel/channel.go, flai/internal/hostapi/people.go, flai/internal/channel/channel_test.go, flai/internal/perf/perf.go, flai/internal/perf/perf_test.go, docs/users/flai.md, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/contract_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/system/flaiover-dashboard.md, flaiover/src/routes/api/threads/+server.ts, flaiover/src/routes/api/threads/threads.test.ts, flaiover/src/lib/server/repo.test.ts, flaiover/src/lib/server/repo-channel.test.ts, flaiover/src/lib/server/repo-forget.test.ts, flaiover/src/routes/api/planner/+server.ts, flaiover/src/routes/api/planner/planner.test.ts, flaiover/src/routes/api/orchestrator/+server.ts, flaiover/src/routes/api/orchestrator/orchestrator.test.ts, flaiover/src/routes/api/analyzer/+server.ts, flaiover/src/routes/api/analyzer/analyzer.test.ts, "flaiover/src/routes/api/roles/[name]/+server.ts", "flaiover/src/routes/api/roles/[name]/roles.test.ts", "flaiover/src/routes/api/items/[id]/plan/+server.ts", "flaiover/src/routes/api/items/[id]/plan/plan.test.ts", flaiover/src/routes/api/host-agent/host-agent.test.ts, design/system/server-performance.md]
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
    - kind: planner
      seconds: 757
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 272
          output: 959
          cache_read: 11746859
          cache_write: 680899
          cost: 0.8208
        - model: claude-opus-5-5
          input: 1763
          output: 76638
          cache_read: 148222754
          cache_write: 2193579
          cost: 41.0331
    - kind: orchestrator
      seconds: 47
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 28
          cache_read: 370064
          cache_write: 9065
          cost: 0.0968
draft: true
cost_of_delay:
  inputs:
    penalty_per_week: 30
    by: orchestrator
    at: 2026-11-08T00:04:14Z
  value: 30
  by: planner-S-0611
  at: 2026-11-08T00:04:59Z
forecast:
  duration: 3h
  delivery: 2026-11-08T06:51:00Z
  basis: "Its own forecast of 3h; 8th in the pull order with an in-progress limit of 5, behind S-0599, S-0601, S-0603, S-0606, S-0607, S-0604, S-0609, S-0602, S-0610, S-0635, S-0241 and S-0608."
  by: flai
  at: 2026-11-08T01:49:10Z
---
# S-0611 The dashboard asks threads.list for open threads and agent.status for open runs, and the request timing counts the send, so a thread list is under 100 KB, not 3.3 MB

## Goal

Cause 21 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): the dashboard asks `threads.list` with `all: true` (`flaiover/src/lib/server/repo.ts`), 3.28 MB for 884 threads of which 12 are open, and `agent.status` sends 437 finished runs and 288 planner runs, 392 KB. Each holds its connection's one write lock while it is written (`flai/internal/channel/channel.go`), and the time is not in `duration_ms`, which is logged before the send.

## Acceptance criteria

- [ ] The dashboard asks `threads.list` for open threads, and for the resolved ones only on a page that shows them.
- [ ] `agent.status` sends the open and live runs, and the finished ones only to a page that shows them, through a parameter the host API documents.
- [ ] The `request answered` event gives the send's time in its own field, and `docs/users/flai.md` § Request timing says so.
- [ ] On this repository a board's `threads.list` is under 100 KB and `agent.status` under 50 KB, and server-performance.md records the before and after.

## Tasks
- T-2689 The request answered event gives the send's time in its own field, send_ms, logged after the write
- T-2690 agent.status takes a finished parameter, and without it sends only the open and live runs
- T-2691 The pages that list finished runs ask agent.status for them, and the rest ask for the open and live runs
- T-2692 The dashboard asks threads.list for open threads, and for the resolved ones only when a page asks with all=1
- T-2693 Measure threads.list and agent.status on this repository before and after, and record cause 21's figures in server-performance.md

## Notes

- These change what the host API answers by default, not how flai keeps state.

### Planning

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flaiover/src/lib/server/repo.ts` | declared | `threads()` and `threadsFor()` send `all: true` to flai, cached under keys without `all` |
| `flai/internal/channel/channel.go` | declared | `rec.Log` runs before `send`, which takes the connection's write lock `wmu` |
| `flai/internal/hostapi/people.go` | declared | Kept as declared. It holds `activity.get`, not `agent.status`, so no task names it; the story's agent narrows it if nothing changes there |
| `flai/internal/channel/channel_test.go` | co-change | 11 of 67 co-changing commits; its test of the event's fields gains the send's |
| `flai/internal/perf/perf.go`, `perf_test.go` | layout | `Answered` and `Recorder.Log` build the event's fields |
| `docs/users/flai.md` | co-change, criterion 3 | § Request timing's field table |
| `flai/internal/hostapi/hostapi.go` | design | `agent.status` takes no parameter today; its comment is the method's documentation |
| `flai/internal/hostapi/writes.go` | layout | `Host.Agent func(root string) any` gains the parameter |
| `flai/internal/hostapi/hostapi_test.go`, `writes_test.go`, `contract_test.go` | layout | Tests that call `agent.status` and the host contract |
| `flai/cmd/serve_actions.go`, `serve_actions_test.go` | layout | `host().Agent` builds the answer: `stories`, `plans`, `last`, and the past runs |
| `flai/internal/serve/agents.go`, `agents_test.go` | layout | `Activity` and `RoleStates`, where finished runs are told apart from waiting, held, and failed ones |
| `design/system/flaiover-dashboard.md` | co-change, criterion 2 | § API documents the parameter beside the runs the pages list |
| `flaiover/src/routes/api/threads/+server.ts`, `threads.test.ts` | layout | The one server caller of `threads()`; it filters resolved threads after asking for them |
| `flaiover/src/lib/server/repo.test.ts`, `repo-channel.test.ts`, `repo-forget.test.ts` | co-change | `repo-channel.test.ts` asserts the `all: true` params; `repo-forget.test.ts` the cache keys |
| `flaiover/src/routes/api/planner`, `orchestrator`, `analyzer`, `roles/[name]`, `items/[id]/plan` `+server.ts` and their tests | layout | The pages that list finished runs ask for them |
| `flaiover/src/routes/api/host-agent/host-agent.test.ts` | layout | `/api/host-agent` keeps the default; its test pins that |
| `design/system/server-performance.md` | design, criterion 4 | Cause 21's before and after, and its sentence that timing ends when the answer is ready to write |

Left out: the other co-change files (`flaiover/src/lib/server/agent.ts`, `flaiover/src/routes/items/[id]/+page.svelte`, `docs/users/flaiover.md`, `design/system/flai-cli.md`, `flai/internal/serve/serve.go`), which changed with these for other reasons. The item page keeps asking `threadsFor` with `all=1`: it shows resolved threads, and one item's list is small. No folder touch.

No `after`. S-0603 (cause 12) also changes `flai/cmd/serve_actions.go` and `flai/internal/serve/agents.go`, but this story reads nothing it writes; the trial merge judges the overlap. S-0602 (cause 13) makes each thread read cheaper, but the size of the answer is this story's alone.

Forecast: 3h, not the 13m `flai forecast` gives. Its 103 s a unit of size times size 7 counts 4 criteria and 3 touches. The plan has 5 tasks over 3 layers, about 30 files in Go and TypeScript, a measurement, and a close-out whose `integration` and `flaiover` tiers mean about 238 s and 176 s (server-performance.md, cause 20). Delivery 2026-11-08T11:00Z: it is 15th in the pull order with an in-progress limit of 5, so about two rounds of stories go first; the replanner moves it as they land.

Cost of delay: 30 USD a week, as `flai cod` works it out from the one input, `penalty_per_week: 30`, which the orchestrator set at 2026-11-08T00:04:14Z from TH-0888's table. It stands unadjusted: no revenue rides on the story, and the figure was scaled against S-0606's 30 for a cause of like weight.
