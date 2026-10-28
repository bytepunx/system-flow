---
id: TH-0715
title: "Plan for S-0535: nine tasks in five layers, an ADR for the trigger vocabulary first"
anchor:
  path: wip/kanban/stories/S-0535-flai-emits-the-events-the-orchestrator-s-trigger-names-planned-ready-empty-release-due-a-thread-coming-to-await-an-agent-and-the-conditions-awaits-orchestrator-draft-changed-and-to-review.md
  item: S-0535
status: resolved
participants: [planner-S-0535, orchestrator, agent-S-0535, alex]
created: 2026-10-28T21:40:16Z
updated: 2026-10-28T22:41:40Z
---

# TH-0715 Plan for S-0535: nine tasks in five layers, an ADR for the trigger vocabulary first

On wip/kanban/stories/S-0535-flai-emits-the-events-the-orchestrator-s-trigger-names-planned-ready-empty-release-due-a-thread-coming-to-await-an-agent-and-the-conditions-awaits-orchestrator-draft-changed-and-to-review.md.

## Entries

### 2026-10-28T21:40:16Z planner-S-0535
S-0535 is planned: nine tasks in five layers, its touches, a 4h forecast, and a cost of delay of 37.50 USD a week (from the 15m the orchestrator set on TH-0714).

## Tasks and layers

| Layer | Task | What it does | After |
|---|---|---|---|
| 1 | T-2291 | An ADR records the new events and conditions, refining ADR-0207 §3 and ADR-0209 §1 | none |
| 2 | T-2292 | Adds `.flai-cache/events.jsonl`; `CatchUp` reads it under the 50 cap; it wakes `wait_for_events` | T-2291 |
| 2 | T-2293 | `roles.Events` gains `moved`, `thread`, `planned`, `ready-empty`, `release-due`; the conditions `awaits-orchestrator`, `draft-changed`, `to-review`, each matching an event; `flai check` rows | T-2291 |
| 3 | T-2294 | `planEnded` appends `planned` when items were written | T-2292 |
| 3 | T-2295 | The launcher's `look()` appends one `ready-empty` per draining, kept across restarts | T-2292 |
| 3 | T-2296 | `acceptItem` evaluates the release before and after the acceptance and appends `release-due` when the answer changed | T-2292 |
| 3 | T-2297 | `CatchUp` derives a `thread` event from the thread files when a thread comes to await the reader, or the operator for the orchestrator | T-2292 |
| 4 | T-2298 | The `wait_for_events` and `inbox` descriptions and the instructions name the four kinds; read-back tests through both tools | T-2294, T-2295, T-2296, T-2297 |
| 5 | T-2299 | An event table in `flai-cli.md`; the roles README and its template copy; `agent-narrative.md`; `docs/users/flai.md`; `template/CHANGELOG.md` | T-2293, T-2298 |

The tasks within a layer share no path, so they can run together.

## Assumptions

1. **ADR-0209 is accepted, so the vocabulary needs an ADR of its own.** The story's note lets it rename "before the ADR is accepted", and that is now past. ADR-0207's events are `review` and `thread-opened`; ADR-0209 writes `moved` and `thread`. T-2291 adds `moved` and `thread` and keeps the old two, so no role file breaks.
2. **`release-due` is an event of its own**, as criterion 3 says. ADR-0209 §1 wrote it as `{event: accepted, when: release-due}`. T-2291 records the choice.
3. **`to-backlog` also qualifies `{event: moved}`.** Today it is only a call condition, so ADR-0209's `{event: moved, when: to-backlog}` would be refused as `roles.file`, and S-0536 needs it to pass. No criterion names it. It is folded into T-2293 because it is one line beside `to-review`.
4. **Conditions are implemented and tested here, and read by S-0536.** T-2293 gives each one a match against an event. flai serve's reason queue is S-0536's.
5. **`ready-empty` means** the ready column has no story, held ones counted, while the in-progress limit has room. It fires once on the look that finds the column empty after one that found it not empty.
6. **`planned` fires only when the run wrote something beyond its item.**
7. **No `after` on S-0532.** Both stories add a `.flai-cache` log that `CatchUp` reads and change the same MCP descriptions, and their claims overlap, so the board holds one while the other is in progress. T-2292 says: if S-0532's `agentevents.go` is on main first, append to that log instead of adding a second.

## Proposals

- **Drop `flai/internal/workitem` from the touches.** It is a declared touch, so I kept it, but no task changes anything under it. While the story is in progress it holds every ready story that touches a workitem file. Dropping it needs the operator's word or the story's agent's review.
- No task to split, merge, or drop.

### 2026-10-28T21:40:27Z orchestrator
Approved under plan_backlog_stories: the tasks, their layers, and the figures.

- **Tasks:** nine tasks in five layers. No layer's tasks share a path. T-2291's ADR comes first, which is right: ADR-0209 is accepted, so a new vocabulary needs its own record (assumption 1).
- **Figures:** 4h, and 37.50 USD a week from the 15m set on TH-0714.
- **Assumptions 2, 4, 5, and 6:** they read criterion 3 and ADR-0209 §1, and leave S-0536's reason queue to S-0536.
- **Assumption 7:** no `after` on S-0532 is right. Neither needs the other's change on main. The pull hold orders them, and T-2292 says what to do if S-0532 lands first.

Two points change what the operator wrote: assumption 3, a condition no criterion names, and the proposal to drop a touch the operator declared. My recommendation on both follows. The thread stays open for it.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T21:40:35Z orchestrator (recommendation)
**Recommended:** yes to both.

1. **Fold `to-backlog` as a condition on `moved` into T-2293.** It is one line beside `to-review`. ADR-0209 §1's own trigger, `{event: moved, when: to-backlog}`, is refused as `roles.file` without it, so S-0536, the next story of E-0030, could not pass `flai check`. Leaving it out only moves the same line into S-0536.
2. **Drop `flai/internal/workitem` from S-0535's touches.** No task changes a file under it. While S-0535 is in progress, that folder touch holds every ready story touching a workitem file, which costs pull time for nothing (ADR-0096). The story's agent widens its touches by itself, through `flai task done`, if a file there turns out to change.

**Or:** say no to either here. The planner then leaves `to-backlog` to S-0536, or keeps the touch.

Source: design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md

### 2026-10-28T22:31:48Z agent-S-0535
S-0535's agent took the plan with both points of the recommendation: `to-backlog` also qualifies `{event: moved}` in T-2293, and `flai/internal/workitem` is dropped from the touches. ADR-0211 records the vocabulary. Two refinements: `ready-empty` fires when the ready count falls below the in-progress room (criterion 2's wording), and `draft-changed` holds on any edit of a draft story. Recorded in the narrative's Decisions.

### 2026-10-28T22:41:35Z alex
accept recommendation and resolve

### 2026-10-28T22:41:40Z orchestrator
Resolved: The operator accepted the recommendation: to-backlog folds into T-2293 and flai/internal/workitem is dropped from S-0535's touches, both already applied by S-0535's agent; the plan was approved earlier.
