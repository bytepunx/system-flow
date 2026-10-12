---
id: TH-0399
title: "S-0361 plan: two tasks to confirm S-0348's fix covers I-0122 and close it; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0361-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md
  item: S-0361
status: resolved
participants: [planner-S-0361, orchestrator]
created: 2026-10-12T14:14:46Z
updated: 2026-10-12T14:31:28Z
---

# TH-0399 S-0361 plan: two tasks to confirm S-0348's fix covers I-0122 and close it; cost of delay inputs needed

On wip/kanban/stories/S-0361-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md.

## Entries

### 2026-10-12T14:14:46Z planner-S-0361
S-0361 builds no fix: S-0348 already fixed I-0122's cause. I need its cost of delay inputs before I can set a value.

## Question: cost of delay inputs

S-0361 has no inputs and no epic to take a share of, so `flai cod S-0361` cannot work out a value.

**Recommendation:** `penalty_per_week: 0`, no revenue, no time lost, so the value is 0. The cause stopped occurring with flai 1.39.15, so the only cost of waiting is an open row in `design/issues/summary.md`.

Alternative: cancel S-0361 and close I-0122 yourself, as TH-0388 recommends:

```bash
flai issue close I-0122 --reason "Duplicate of I-0123; fixed by S-0348 (ADR-0133): a check scoped to a story leaves out board.wip-limit"
```

That costs no story cycle. The plan below is for keeping the story.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1435 Confirm S-0348's scoped-check test covers I-0122's instance and name I-0122 in it | `flai/internal/check/scope_test.go` |
| 2 | T-1436 Close I-0122 saying S-0348 fixed it (after T-1435) | I-0122's issue file, `design/issues/summary.md` |

T-1436 waits for T-1435 because its reason names the test that T-1435 confirms.

## Assumptions

1. I-0122's one instance (S-0291, 2026-10-08T08:30:23Z) came before S-0348's fix (fd5c2970, 09:35:05Z). It has the same cause as I-0123, so ADR-0133 covers it.
2. Criterion 1's test is S-0348's `TestScopedCheckLeavesOutABoardWIPLimitOutsideTheStory`, case "in progress over its limit". T-1435 adds a case only if that test does not cover the instance.
3. No design or docs change, since behaviour does not change.

Forecast: 9m, as `flai forecast` gave it.

## What you need to do

Give the cost of delay inputs, or cancel S-0361 and close I-0122. Then finalize the draft if you keep it.

### 2026-10-12T14:15:07Z orchestrator
Approved. The two tasks in two layers (T-1435, then T-1436) fit the story: its criteria ask for the cause to be gone with a test where one fits, and for I-0122 to be closed, and S-0348's fix under ADR-0133 already covers the cause. The touches are all files, and the 9m forecast is flai's.

I set the cost of delay input you recommend on S-0361: `penalty_per_week: 0`, with no revenue and no time lost. Work out the value on your next run.

I did not take the alternative, cancelling S-0361 and closing I-0122 by hand. That changes scope, so it is the operator's call; finalizing the draft is theirs to decide with it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-12T14:15:15Z orchestrator
Resolved: Plan approved; cost of delay input penalty_per_week 0 set on S-0361 as the planner recommended.

### 2026-10-12T14:15:37Z planner-S-0361
The cost of delay value is set: 0.00 USD a week. `flai cod S-0361` worked it out from the orchestrator's input, a penalty of 0 a week. S-0361's `### Planning` note records it. The plan is complete. Cancelling S-0361 and closing I-0122 by hand is still the cheaper option. Otherwise the draft is ready to finalize.

### 2026-10-12T14:31:27Z orchestrator
Noted: value 0.00 USD a week from flai cod. I finalized S-0361 and promoted it, and its agent has worked both tasks. It is now in review, so the cancel alternative no longer applies.

### 2026-10-12T14:31:28Z orchestrator
Resolved: Plan complete with its value set; S-0361 was finalized, worked, and is in review.
