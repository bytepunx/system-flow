---
id: TH-0841
title: "Plan of S-0571: a defined role runs once per ready story under each-ready-story; 4 tasks in 2 layers, 1h, 25.00 USD a week"
anchor:
  path: wip/kanban/stories/S-0571-flai-serve-does-not-start-a-defined-role-that-takes-items-on-a-schedule-entry.md
  item: S-0571
status: resolved
participants: [planner-S-0571, agent-S-0571]
created: 2026-11-07T03:47:33Z
updated: 2026-11-07T03:51:06Z
---

# TH-0841 Plan of S-0571: a defined role runs once per ready story under each-ready-story; 4 tasks in 2 layers, 1h, 25.00 USD a week

On wip/kanban/stories/S-0571-flai-serve-does-not-start-a-defined-role-that-takes-items-on-a-schedule-entry.md.

## Entries

### 2026-11-07T03:47:33Z planner-S-0571
S-0571 is planned. It had no plan before. The start prompt asked me to read a thread named `asked`, but no such thread exists and no thread was ever on S-0571. That looks like an unfilled field in flai serve's replan prompt.

## Proposed fix

I-0262 offers two fixes. I planned the first: `flai serve` starts a defined role that takes stories once for each ready story, in pull order, when an `each-ready-story` schedule comes round. A held start is kept per story. This is what ADR-0250 §5 already says the condition means, so no new ADR is planned. The second fix, `flai check` reporting the entry, would only hide the gap. If you would rather have the check-only fix, say so and I will replan T-2511 down to a check.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2511 role scheduler starts once per ready story, with a held start per story, and the test that reproduces I-0262 | `flai/internal/serve/role_schedule.go`, `role_schedule_test.go` | none |
| 1 | T-2512 `flai check` reports `each-ready-story` on a role whose `item` does not list `story` | `flai/internal/roles/roles.go`, `roles_test.go` | none |
| 2 | T-2515 design, roles README and its template copy, user and operator guides, template changelog | six docs files | T-2511, T-2512 |
| 2 | T-2517 close I-0262 with a reason | the issue file and `design/issues/summary.md` | T-2511, T-2512 |

The tasks within each layer share no path, so they can run together.

## Figures

- Touches: 12 files, no folder, no story `after`.
- Forecast 1h, raised from flai's 24m. S-0554 covered similar scope in 77m.
- Cost of delay value 25.00 USD a week, from flai's inputs (10m a cycle). Unchanged.

## Assumptions

- `roles.EachReadyStory` and the board's pull order work for a defined role as they do for the planner.
- T-2512 is a small addition beyond I-0262's words: without it, a role that takes only epics under `each-ready-story` would still start nothing.

### 2026-11-07T03:50:33Z agent-S-0571
Plan accepted as drafted: the first fix, one start per ready story under `each-ready-story`, with a held start per story, and T-2512's check for a role whose `item` does not list `story`. No ADR: ADR-0250 §5 already says what the condition means.

One change of how, not what: I work T-2512 myself first, then hand T-2511 to a sub-agent, because `serve` imports `roles` and the two would otherwise build against each other's half-done edits in one worktree. T-2517 I close myself after T-2511 and T-2512.

On the `asked` thread: in this tree the planner's `start` renders `.AskedOn` only for an orchestrator trigger, and this run's trigger was `asked`, so the field was empty here. The installed flai is 1.83.0, older than the tree. I record nothing for it in S-0571.

Source: design/adrs/0250-a-role-file-s-start-prompt-is-a-text-template-in-its-front-matter-s-start-field.md § 5. The trigger conditions

### 2026-11-07T03:51:06Z agent-S-0571
Resolved: Plan accepted by S-0571's agent; tasks worked as drafted.
