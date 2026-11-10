---
id: TH-0903
title: "Plan of E-0033: 14 stories, one per cause 8 to 21, with 66 tasks"
anchor:
  path: wip/kanban/epics/E-0033-flai-stays-responsive-at-600-stories-checked-writes-kept-parses-and-concurrent.md
  item: E-0033
status: resolved
participants: [planner-E-0033, orchestrator, flai serve, alex]
created: 2026-11-08T00:03:50Z
updated: 2026-11-10T08:48:59Z
---

# TH-0903 Plan of E-0033: 14 stories, one per cause 8 to 21, with 66 tasks

On wip/kanban/epics/E-0033-flai-stays-responsive-at-600-stories-checked-writes-kept-parses-and-concurrent.md.

## Entries

### 2026-11-08T00:03:50Z planner-E-0033
Every story of E-0033 now has touches, a forecast, Planning notes, and tasks. The tasks were written by each story's own planner run, which flai serve started at 23:45. Each of those runs opened a plan thread on its story. My explorers' plans matched theirs, so I wrote no tasks of my own.

## Stories, tasks, and layers

| Story | Cause | State | Forecast | Layers |
|-------|------:|-------|---------:|--------|
| S-0598 | 8 | in progress | 45m | T-2620 → T-2622, T-2625 → T-2635, T-2640 |
| S-0599 | 9 | in progress | 2h30m | T-2628, T-2632 → T-2641 → T-2647 → T-2655 |
| S-0600 | 10 | in progress | 1h30m | T-2674, T-2678 → T-2682, T-2685 → T-2686, T-2688 |
| S-0601 | 11 | backlog | 40m | T-2633, T-2643 → T-2648 |
| S-0602 | 13 | backlog | 1h30m | T-2660 → T-2667, T-2671, T-2673, T-2676, T-2680 → T-2684 |
| S-0603 | 12 | backlog | 1h | T-2621, T-2629 → T-2623 → T-2638 |
| S-0604 | 14 | draft | 1h | T-2659, T-2666, T-2670 → T-2677 → T-2681 |
| S-0605 | 15 | backlog | 30m | T-2639 → T-2645 → T-2654 |
| S-0606 | 16 | ready | 45m | T-2615, T-2616 → T-2617, T-2618 → T-2619 |
| S-0607 | 17 | draft | 1h | T-2637, T-2656 → T-2646, T-2650 → T-2663 |
| S-0608 | 18 | draft | 1h | T-2662, T-2668 → T-2672, T-2675, T-2679 → T-2683 → T-2687 |
| S-0609 | 19 | draft | 1h | T-2636, T-2642 → T-2649 → T-2658 → T-2664 |
| S-0610 | 20 | draft | 1h30m | T-2634 → T-2644 → T-2653 → T-2661 |
| S-0611 | 21 | draft | 3h | T-2689, T-2690, T-2692 → T-2691 → T-2693 |

The last task of each story measures its cause before and after and records it in `design/system/server-performance.md`.

## after

- S-0605 and S-0606 name S-0596. They read the section of `server-performance.md` that S-0596 wrote. S-0596 is done, so neither is held.
- No other story names another. Several share files (`check.go`, `store.go`, `hostapi.go`, `metrics.go`, `serve_actions.go`), and the trial merge judges those.

## Cost of delay

- Inputs are set on S-0598 (by you), and on S-0601, S-0602, S-0603, S-0605, S-0607, and S-0610 (by the orchestrator). I wrote the values `flai cod` works out for S-0601 (50), S-0602 (50), S-0605 (20), S-0607 (20), and S-0610 (80).
- S-0599, S-0600, S-0604, S-0608, S-0609, and S-0611 have no inputs yet. They wait on TH-0888, where the orchestrator recommended my table.

## Proposals: nothing added, nothing dropped

1. **S-0608 after S-0606, or a different key.** S-0608's T-2668 counts a generation in `flai/internal/workitem/store.go`. S-0606's T-2615 rewrites that store around a kept snapshot, after which a list no longer stats files. If S-0608 lands first, its count goes stale under S-0606 and no text conflict shows it. Recommended: give S-0608 `after: [S-0606]`. The alternative is to key the kept replay on a fingerprint of the item fields `metrics` reads, which needs no `after`.
2. **Split S-0611.** It is 3h over flai and flaiover. Recommended: split it into the flai half (the send time in `request answered`, and `agent.status`'s finished-runs parameter) and the dashboard half (asking for open threads and open runs), with the second after the first, since it reads the parameter.
3. **S-0601's kept props key.** It keys on the component file and the svelte version. A skeleton also depends on the type declarations the component imports, so a kept answer could go stale when one changes. Recommended: add those declarations' stamps to the key in T-2643.
4. **S-0609's fix contract.** Treating a `fix_command` that exits 0 as "nothing remains" narrows ADR-0139 and ADR-0140, so it needs a new ADR. Recommended: T-2642 writes it.

## Assumptions

- Each cause from 8 to 21 is one story, as S-0596 proposed. No cause needs a second.
- Stories the operator has finalized keep their words. I changed only the cost of delay values above.
- New files named in touches, such as `check/itemwrite.go`, are the planners' predictions. `task done` widens touches to what is actually changed.

### 2026-11-08T00:26:00Z orchestrator (recommendation)
**Recommended:** keep S-0611 whole (proposal 2: no split). Proposals 1, 3, and 4 are already in the plans and need nothing from you. Once you confirm, I finalize S-0611. TH-0904 holds it for this answer.

1. **S-0608 after S-0606: done.** planner-S-0608 gave S-0608 `after: [S-0606]`. It also rewrote T-2668 so the generation moves when S-0606's snapshot is replaced. S-0608 is finalized and held until S-0606 is done.
2. **Split S-0611: no.** `work-management.md` sizes a story as one demonstrable increment for one to a few sessions. 3h fits that. S-0611 fixes one cause (21), and criterion 4 measures both answers once. Its layers already run T-2690 before T-2691 reads the `finished` parameter. A split would add an `after` and a second close-out. Each close-out runs `integration` and `flaiover`, which average about 238 s and 176 s (cause 20).
3. **S-0601's props key: done.** T-2643 already keys each entry on the stamps of the files `props.mjs` read for the component's skeletons. `props.mjs` is in S-0601's touches for that reason.
4. **S-0609's ADR: yes, and it needs no scope change.** ADR-0140 § Decision runs each fixed tier's `command` after the fix. Criterion 1 drops that run when the fix's exit already says what remains. Under `decisions.md`, the story's agent writes an ADR refining ADR-0140. I noted this for it on TH-0905, and I finalized S-0609.

**Alternative to 2:** split S-0611 into a flai half (T-2689, T-2690) and a dashboard half (T-2691 to T-2693), with the second after the first. Its planner then re-plans the tasks across the two drafts.

Source: design/conventions/work-management.md

### 2026-11-10T00:26:05Z flai serve
Confirmed the recommendation of 2026-11-08T00:26:00Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/conventions/work-management.md

### 2026-11-10T08:48:59Z alex
Resolved: S-0636 was accepted
