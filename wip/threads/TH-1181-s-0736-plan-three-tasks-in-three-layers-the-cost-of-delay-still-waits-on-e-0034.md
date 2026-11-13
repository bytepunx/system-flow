---
id: TH-1181
title: "S-0736 plan: three tasks in three layers; the cost of delay still waits on E-0034's inputs"
anchor:
  path: wip/kanban/stories/S-0736-an-adr-fixes-the-workflow-file-s-format-and-place-under-design-workflows-its.md
  item: S-0736
status: open
participants: [planner-S-0736, orchestrator]
created: 2026-11-13T20:30:01Z
updated: 2026-11-13T20:31:18Z
---

# TH-1181 S-0736 plan: three tasks in three layers; the cost of delay still waits on E-0034's inputs

On wip/kanban/stories/S-0736-an-adr-fixes-the-workflow-file-s-format-and-place-under-design-workflows-its.md.

## Entries

### 2026-11-13T20:30:01Z planner-S-0736
**Needs you:** a cost of delay input for E-0034, or for S-0736. Without one, `flai cod S-0736` refuses, and the draft cannot be finalized.

**Recommended:** `time_lost_per_cycle: 2h` on E-0034, as asked on TH-1097. That is 300 USD a week, and S-0736's value is its share by forecast.

```sh
flai edit E-0034 --time-lost-per-cycle 2h
```

Alternatives: `penalty_per_week: 25` if nobody waits on the epic, or a `revenue_per_week` if a project or client does.

## Plan

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3187 Write the ADR, proposed | `design/adrs/README.md` | none |
| 2 | T-3193 Settle the open choices with you, accept the ADR, mark ADR-0003 and ADR-0004 | `design/adrs/README.md`, `design/adrs/0003-work-item-hierarchy.md`, `design/adrs/0004-workflow-states-and-transitions.md` | T-3187 |
| 3 | T-3200 Describe the workflow file in the design and link the ADR | `design/system/workflow.md`, `design/system/work-hierarchy.md`, `design/workflows/README.md` | T-3193 |

Each task reads the one before, so they run one at a time. I checked all three again against the nine criteria. None needs a split, merge, or drop, and none was changed.

## Figures

- Forecast 45m, kept over the 28m `flai forecast` gives. The ADR is about as large as ADR-0207, and the ADR gets rewritten after your answers. The delivery comes from flai's replan.
- Touches: the six declared files, kept. `flai touches suggest` also lists `flai-cli.md`, the guides, `metrics.md`, and the template. I left them out because the stories that build the workflow change them.
- No `after`: S-0736 reads no other story's output.

## Assumptions

- The ADR's own file is not a touch. It gets its number from `flai adr new` when the story runs, and the index touch covers the claim.
- T-3193's open choices go on a thread of the story's agent, not this one.

### 2026-11-13T20:31:18Z orchestrator (recommendation)
**Recommended:** answer or confirm my recommendation on TH-1097: `time_lost_per_cycle: 2h` on E-0034. One answer there covers this thread, TH-1113, and the 18 other E-0034 drafts.

```sh
flai edit E-0034 --time-lost-per-cycle 2h
```

**The plan fits the story.** The three tasks, one per layer, cover the nine criteria. Their touches are the story's six files. The 45m forecast and the lack of an `after` are explained in the Notes. Nothing needs changing.

**Why I do not set the input myself:**

- The planner's figure goes on the epic, so S-0736's value is its share by forecast. ADR-0224 lets me set an epic's input only once you decide it on a thread.
- Setting 2h on S-0736 alone would give this 45m story the whole epic's 300 USD a week. That does not match the planner's reasoning.

**What happens next:** once you answer, I set the input under `set_cost_of_delay`. The planner's next run writes each story's value, and I finalize S-0736 when `flai promote --drafts` lists it complete. Today it lacks only a cost of delay value.

**Alternatives:** `penalty_per_week: 25` on E-0034 if nobody waits on the epic, or a `revenue_per_week` if a project or client does.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0224-orchestration-permissions-set-cost-of-delay-off-by-default-lets-the.md § Decision
