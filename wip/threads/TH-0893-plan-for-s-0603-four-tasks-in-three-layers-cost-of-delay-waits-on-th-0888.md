---
id: TH-0893
title: "Plan for S-0603: four tasks in three layers; cost of delay waits on TH-0888"
anchor:
  path: wip/kanban/stories/S-0603-agent-status-keeps-agents-json-parsed-and-reads-the-conversations-once-a-call.md
  item: S-0603
status: resolved
participants: [planner-S-0603, orchestrator]
created: 2026-11-07T23:50:22Z
updated: 2026-11-07T23:53:34Z
---

# TH-0893 Plan for S-0603: four tasks in three layers; cost of delay waits on TH-0888

On wip/kanban/stories/S-0603-agent-status-keeps-agents-json-parsed-and-reads-the-conversations-once-a-call.md.

## Entries

### 2026-11-07T23:50:22Z planner-S-0603
S-0603 now has four tasks, 11 file touches, and a 1h forecast. Its cost of delay value waits on one input.

## Tasks and layers

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2621 `AgentStates` keeps `agents.json` parsed, by size and modification time, `agentsMu` held only to read or replace it | `flai/internal/serve/agents.go`, `agents_test.go` | none |
| 1 | T-2629 `Host.Agent` takes the request's context; `agent.status` records `agent.states`, `agent.activity`, `agent.roles` | `flai/internal/hostapi/writes.go`, `hostapi.go`, `writes_test.go`, `flai/cmd/serve_actions.go`, `serve_actions_test.go` | none |
| 2 | T-2623 `Activity` reads the conversations once, for every live run and for `held`; same answer, a test with three live runs | `flai/internal/serve/agents.go`, `agents_test.go`, `flai/internal/messages/messages.go`, `flai/internal/workitem/hold.go`, `share.go` | T-2621, same files |
| 3 | T-2638 Measure 13 projects at once, before and after, on the operator's host; record it in `design/system/server-performance.md` | `design/system/server-performance.md` | T-2621, T-2623, T-2629 |

## Assumptions

- "Reads the conversations once" covers `held` too, as the goal says: one read feeds `runActivity` and `Repo.Holds`. `Repo.Holds` keeps its signature for its other callers.
- `Host.Agent` may change its signature: only `flai/cmd/serve_actions.go` and one test stub use it.
- The measurement runs as S-0596's did, in an in-process program under `/tmp`, writing nothing.
- No `after` on S-0602 (cause 13), though both change `flai/internal/messages/messages.go`: neither reads the other's output, and the trial merge judges the overlap.

## Figures

- Forecast 1h, delivery 2026-11-08T01:23:00Z. `flai forecast` said 22m. I raised it for the concurrency tests, the three packages, the 13-call measurement on a loaded host, and a close-out that runs the integration and `flaiover-flai` tiers.

## Question: cost of delay input

S-0603 and E-0033 have no cost of delay inputs, so `flai cod` refuses and the story has no value.

Recommended: `penalty_per_week: 100` on S-0603, as on TH-0888. It is the highest after S-0598's 150: this cause is about 30% of flai serve's wall time. No revenue rides on it.

Alternative: confirm TH-0888's whole table, which sets this figure with the others. I work out the value once an input is set.

### 2026-11-07T23:52:47Z planner-S-0603
The orchestrator set `penalty_per_week: 100` on S-0603 at 2026-11-07T23:51:56Z. `flai cod` works it out to 100 USD a week, and I set that as the value. The plan above stands: T-2621 and T-2629, then T-2623, then T-2638. Nothing waits on the operator now.

### 2026-11-07T23:53:07Z orchestrator
Plan approved. The four tasks in three layers fit S-0603's criteria. T-2623 waits on T-2621 for the shared `agents.go`, and T-2638 measures all three.

I set `penalty_per_week: 100` on S-0603, as you recommend. It is the figure on TH-0888. It is the story's first input, and E-0033 has none. I set no value: you work it out on your next run.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T23:53:34Z orchestrator
Resolved: Plan approved; penalty_per_week 100 set on S-0603, and planner-S-0603 set its value.
