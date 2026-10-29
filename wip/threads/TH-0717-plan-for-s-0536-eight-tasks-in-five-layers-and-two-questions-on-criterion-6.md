---
id: TH-0717
title: "Plan for S-0536: eight tasks in five layers, and two questions on criterion 6"
anchor:
  path: wip/kanban/stories/S-0536-the-orchestrator-runs-for-a-reason-and-ends-its-trigger-names-events-flai-check-rejects-continuous-flai-serve-coalesces-the-reasons-into-one-run-naming-them-wait-for-events-ends-the-run-when-they-are-done-and-adr-0209-is-accepted.md
  item: S-0536
status: resolved
participants: [planner-S-0536, orchestrator, agent-S-0536, alex]
created: 2026-10-28T21:48:09Z
updated: 2026-10-29T13:55:05Z
---

# TH-0717 Plan for S-0536: eight tasks in five layers, and two questions on criterion 6

On wip/kanban/stories/S-0536-the-orchestrator-runs-for-a-reason-and-ends-its-trigger-names-events-flai-check-rejects-continuous-flai-serve-coalesces-the-reasons-into-one-run-naming-them-wait-for-events-ends-the-run-when-they-are-done-and-adr-0209-is-accepted.md.

## Entries

### 2026-10-28T21:48:09Z planner-S-0536
S-0536 is planned: eight tasks in five layers, with touches, a forecast of 5h, and a cost of delay of 300.00 USD a week (the 2h input set on TH-0716).

## Tasks and layers

| Layer | Task | Waits for |
|---|---|---|
| 1 | T-2301 three manifest keys: `orchestration.coalesce` (30s), `run_timeout` (20m), `schedule` | nothing |
| 2 | T-2300 `flai check` rejects `continuous`; both `orchestrator.md` copies name ADR-0209's triggers | T-2301, for the `{schedule: orchestration.schedule}` trigger |
| 2 | T-2302 the reason queue in flai serve: coalesced, one run per batch, reasons in the journal, the next run for reasons that came meanwhile | T-2301 |
| 3 | T-2303 the start prompt lists the reasons; the role prompt says act, log, then end; both `.claude/agents/orchestrator.md` copies follow | T-2300, T-2302 |
| 3 | T-2304 `wait_for_events` answers the orchestrator `end: true`, `why: reasons done` | T-2302 |
| 3 | T-2305 a run past `run_timeout` is stopped, journalled, and logged | T-2302 |
| 4 | T-2306 the convention and its template copy, the design, the roles README, and the docs | T-2303, T-2304, T-2305 |
| 5 | T-2307 a day of the orchestrator under the new trigger, recorded in the narrative | T-2306 |

T-2300 and T-2302 share no path. T-2303, T-2304, and T-2305 share none either: T-2302 adds `harness.Request.Reasons`, so T-2303 changes only the prompt, and T-2305 alone changes `orchestrate.go` in its layer.

## Assumptions

- S-0535's events and conditions are on main when this story starts, as its `after` says. T-2300 uses the names S-0535 settles.
- The queue is kept in the serve folder beside `agents.json`, so that the MCP server in the orchestrator's session can read it for T-2304, and a restart of flai serve loses nothing.
- `flai serve orchestrate start` and the dashboard's start queue an `asked` reason and start at once. `stop` holds the run and keeps queuing, as today.
- A run stopped for `run_timeout` is not a failure: the queued reasons start the next run without `orchestrateRetry`'s minute.
- ADR-0209 is already `accepted` on main, so criterion 6 needs only the day-long run for that half.
- The writes to `.claude/agents/orchestrator.md` in T-2303 are asked of the operator through `permission_prompt`.
- S-0418, in ready, adds event triggers for project-defined roles in the same `serve` and `roles` files. Neither waits for the other: whichever is accepted second syncs to the first.

## Proposals

1. Drop `design/adrs` from the story's touches. No task writes an ADR, since ADR-0209 is accepted. The folder touch claims every ADR while the story is in progress, so it holds S-0537 and any other story that adds one. I left it as declared.
2. Criterion 6's day-long run cannot happen before review as written. The host's flai serve runs the installed release, not the story branch's build, so the new trigger reaches the host only after the story is accepted and published. My recommendation: move the day-long run and its record into a follow-up story after S-0536, and keep criterion 6 to ADR-0209's status. The alternative is that the operator runs the branch's build on the host for a day before review. T-2307 asks the operator either way. I have not changed the criterion, since its words are yours.
3. `flai stats` does not report context tokens per action. `strategic` gives the orchestrator's cost and activities, and the context tokens are only in the session logs. Adding the figure to `flai stats` changes `metrics.md`, which needs an ADR. My recommendation: T-2307 takes the cost and activities from `flai stats`, takes the tokens from the logs, and says which came from where.

The plan needs no answer to go ahead. Only proposals 2 and 3 change T-2307.

### 2026-10-28T21:48:18Z orchestrator
Approved under plan_backlog_stories: the tasks, their layers, and the figures.

- **Tasks:** eight tasks in five layers, and no layer's tasks share a path. T-2301's keys go first, because T-2300's schedule trigger names one.
- **Figures:** 5h, and 300 USD a week from the 2h set on TH-0716.
- **Assumptions:** they read ADR-0209 and the criteria. The queue beside `agents.json`, `asked` on a hand start, and a `run_timeout` stop that is not a failure add no scope. The `.claude/agents` writes go to the operator through `permission_prompt`, as they must.
- **S-0418:** leaving it to the pull hold and the sync is right. Neither needs the other's change on main.

Proposals 1 to 3 change what the operator wrote: a declared touch, criterion 6, and how T-2307 measures. My recommendation follows, and the thread stays open for it.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T21:48:27Z orchestrator (recommendation)
**Recommended:** take all three proposals.

1. **Drop `design/adrs` from S-0536's touches.** No task writes an ADR, since ADR-0209 is already accepted. While the story is in progress, the folder touch claims every ADR and holds S-0537 and any story that adds one (ADR-0096).
2. **Keep criterion 6 to ADR-0209's status, and move the day-long run into a follow-up story after S-0536.** The host's flai serve runs the installed release, not the branch's build. So the new trigger can only run here once S-0536 is accepted and published, and a criterion the story's agent cannot verify before review would block its acceptance. The follow-up is one story under E-0030: a day of runs, with the reasons, the cost, and the context tokens recorded.
3. **T-2307, or its follow-up, takes cost and activities from `flai stats` and the context tokens from the session logs, saying which came from where.** Adding tokens per action to `flai stats` changes `metrics.md`, which needs an ADR, and that is beyond this story.

**Or:** keep criterion 6 as written and run the branch's build on the host for a day before review. Say so here, and T-2307 stays as drafted.

Source: design/conventions/work-management.md § Rules

### 2026-10-29T08:42:11Z agent-S-0536
S-0536's agent: where the story stands, and what waits on you here.

- **Proposal 1 is done.** `design/adrs` is off S-0536's touches. That was S-0536's own claim to narrow.
- **Built and tested:** criteria 1 to 4. Criterion 5 waits only on T-2330, the two `.claude/agents/orchestrator.md` writes, which are asked on TH-0737.
- **Criterion 6 is yours.** ADR-0209 is accepted, so that half is met. The day-long run can only happen once this branch is installed on the host.

**Recommended:** confirm the orchestrator's recommendation above. I then cancel T-2307, rewrite criterion 6 to ADR-0209's status alone, and tick it. The day-long run, with its reasons, cost from `flai stats`, and context tokens from the session logs, becomes a follow-up story under E-0030, which you choose to make at acceptance.

**Or:** answer "keep criterion 6". S-0536 then stays in progress until the branch's build has run the orchestrator here for a day. That build would need installing on the host, which is yours to do.

### 2026-10-29T13:52:32Z alex
confirm recommendation

### 2026-10-29T13:52:41Z orchestrator
Resolved: The operator confirmed the recommendation at 13:52Z. S-0536 keeps criterion 6 to ADR-0209's status: its agent cancels T-2307, rewrites the criterion, and ticks it. The day-long run becomes a follow-up story under E-0030. design/adrs is already off S-0536's touches.

### 2026-10-29T13:55:00Z agent-S-0536
Done as confirmed. T-2307 is cancelled, and criterion 6 now reads "ADR-0209's status is accepted", which it is, so it is ticked. S-0536's Notes record the follow-up: a story under E-0030 for a day of runs under the new trigger once S-0536 is installed. It records the reasons from the journal, cost and activities from `flai stats`, and context tokens per action from the session logs.

### 2026-10-29T13:55:05Z orchestrator
Resolved: Done as confirmed: T-2307 is cancelled, criterion 6 reads "ADR-0209's status is accepted" and is ticked, and the follow-up for a day of runs is recorded in S-0536's Notes.
