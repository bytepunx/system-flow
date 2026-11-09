---
id: TH-0940
title: "S-0622 plan: three tasks in two layers; the cause is fixed by the manifest, so the story pins it with tests and closes I-0284"
anchor:
  path: wip/kanban/stories/S-0622-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md
  item: S-0622
status: resolved
participants: [planner-S-0622, orchestrator]
created: 2026-11-09T18:06:07Z
updated: 2026-11-09T18:07:18Z
---

# TH-0940 S-0622 plan: three tasks in two layers; the cause is fixed by the manifest, so the story pins it with tests and closes I-0284

On wip/kanban/stories/S-0622-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md.

## Entries

### 2026-11-09T18:06:07Z planner-S-0622
**Proposed solution:** no code change. I-0284 had one instance, S-0547 and S-0549 waiting about 37 h each on TH-0784 and TH-0789. Its cause was configuration: `orchestration.permissions` had no `set_cost_of_delay`. You turned the permission on in c76b23a2d at 2026-11-07T09:17:55Z, an hour after the analysis. Since then no confirmed input has waited.

The code path holds, as read:

1. A confirmation sets the thread `answered` (`threads.go` `confirm`).
2. The orchestrator took part in the thread, so the thread now awaits it (`inbox.go` `Summarize`).
3. Its trigger, `thread` / `awaits-orchestrator`, gives it a `thread:TH-nnnn` reason.
4. `OrchestratorSetsCostOfDelay` admits the input on any open story while the permission is on.

No test covers steps 1 to 3 together, nor an in-progress story's input. The story adds those two regression tests and closes the issue.

## Tasks

| Layer | Task | Touches | After |
|------:|------|---------|-------|
| 1 | T-2763 Test that a confirmation of a cost of delay recommendation on a planner's thread gives the orchestrator a thread reason, by your confirm and by flai serve's after `confirm_after` | `flai/internal/serve/orchestrate_confirm_test.go` (new) | none |
| 1 | T-2764 Test that the orchestrator sets an in-progress story's input (the fixture's S-0007) only while `set_cost_of_delay` is on, by `flai edit` and `item_edit` | `flai/cmd/orchestrate_permissions_test.go` | none |
| 2 | T-2765 Close I-0284 with `flai issue close --reason`, naming c76b23a2d and both tests | `design/issues/I-0284-….md`, `design/issues/summary.md` | T-2763, T-2764 |

T-2763 and T-2764 share no path, so they run together. T-2765 waits for both, because its reason names their tests.

## Figures

- **Touches:** four files, no folder. `flai touches suggest`'s co-change list names the permission's implementation, which this story does not change, so none of it is kept.
- **Forecast:** 45m, raised from flai's 3m, which counts two criteria and no touches. Delivery 2026-11-10T12:00:00Z, because review is full at 6 of 5.
- **Cost of delay:** 11,250 USD a week, from the input flai carried over from I-0284: `time_lost_per_cycle` 75h. That figure is W44's two waits, measured before the permission was on. It is not mine to change. If you want the value to match today's cost, which is near zero, set a smaller input, for example `flai edit S-0622 --time-lost-per-cycle 10m`. As it stands, it ranks this story above everything under a cod or WSJF policy.
- **After:** none.

## Assumptions

- "The cause no longer occurs" is met by the manifest change plus tests that fail if the path breaks. Turning `set_cost_of_delay` on in the template's default manifest is out of scope: ADR-0224 keeps it off by default.
- If T-2763's test fails, the task fixes the reason path, and its touches widen as it closes.
- Not covered, and not proposed as a story yet: a plain answer from you, rather than a confirmation, on a planner's thread the orchestrator never wrote on starts no orchestrator run. In practice the orchestrator recommends within seconds, so it is a participant before you answer.

Nothing to split, merge, or drop.

### 2026-11-09T18:07:07Z orchestrator
**Approved:** the plan fits S-0622. Its agent works it as written once you finalize the story.

- **Criterion 1:** T-2763 and T-2764 pin both halves of I-0284's instance: the confirmation reaching the orchestrator, and an in-progress story's input. Both are new cases on existing fixtures.
- **Criterion 2:** T-2765 closes I-0284, after both tests, since its reason names them. Its check that `set_cost_of_delay: true` still holds on main is right: the manifest has it, from c76b23a2d (2026-11-07T09:17:55Z), as you said.
- **Layers and touches:** T-2763 and T-2764 share no path, so layer 1 runs together. Four files, no folder, each named by a task.
- **Forecast:** 45m over flai's 3m holds, for the reason you gave.
- **Assumption on scope:** leaving the template's default off is right. ADR-0224 keeps `set_cost_of_delay` off by default.

**The cost of delay input is not mine to change.** The 75h `time_lost_per_cycle` is an input that already exists, not a backlog story's first. TH-0941 already puts it to the operator, recommending a smaller input such as 1h, and it holds S-0622's finalization until they answer. Once they decide, I set it under `set_cost_of_delay`, and you work out the value on your next run.

**Not covered:** a plain answer on a planner's thread the orchestrator never wrote on. I agree no story is needed now.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T18:07:18Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and forecast fit S-0622. The cost of delay input question stays with the operator on TH-0941.
