---
id: TH-0351
title: "S-0323 plan: tasks, layers, and assumptions"
anchor:
  path: wip/kanban/stories/S-0323-flai-check-finds-narrative-state-outside-the-story-at-close-out.md
  item: S-0323
status: resolved
participants: [planner-S-0323, orchestrator]
created: 2026-10-08T00:27:42Z
updated: 2026-10-08T00:27:50Z
---

# TH-0351 S-0323 plan: tasks, layers, and assumptions

On wip/kanban/stories/S-0323-flai-check-finds-narrative-state-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T00:27:42Z planner-S-0323
S-0323 is planned: 4 tasks in 3 layers, 14 touches, forecast 24m, cost of delay 37.5 USD a week. It waits for S-0318.

## Remedy proposed

A check scoped to a story leaves out a `narrative.state` finding on another story's narrative, so a close-out records none. This follows ADR-0122 for `item.archive` and S-0318's leave-out for `markdown.*`.

All three of I-0109's instances are another story in progress, just started, whose narrative still holds the template's placeholder. Only that story's agent can write it, and its own close-out stops on it in the `narrative` step.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1313 ADR refining ADR-0085 and S-0318's ADR | none | `design/adrs` |
| 2 | T-1314 `ScopeToStory` leaves the finding out, with tests that reproduce I-0109 | T-1313 | `flai/internal/check/scope.go`, `scope_test.go`, `flai/cmd/check.go`, `check_test.go`, `docs/users/flai-reference.md` |
| 3 | T-1315 close I-0109 | T-1314 | the issue file and `design/issues/summary.md` |
| 3 | T-1316 convention, template copy, changelog, design, users' guide | T-1314 | 6 files |

T-1315 and T-1316 share no path, so they run together.

## Assumptions

- The leave-out covers `narrative.state` only. Other findings on another story's narrative (`narrative.front-matter`, `narrative.stream`, `narrative.section`) stay recorded, because they point at a hand edit or a flai defect. The broader option the orchestrator raised on TH-0346, every finding on another open story's narrative, is recorded in T-1313's ADR as the alternative.
- `after: S-0318` on the story. Both change `ScopeToStory` and the same convention line, and T-1314 reuses S-0318's test for another open story's narrative.
- `design/adrs` stays a folder touch, because `flai adr new` allocates the number when it runs.
- The plan covers I-0109 only. TH-0346, merging S-0325 (I-0111) in, is still the operator's.

## Proposals

- If the operator merges S-0325 into S-0323 on TH-0346, add one task in layer 3 beside T-1315: close I-0111 with `flai issue close`, touching its issue file and `design/issues/summary.md`, after T-1314. Give it an `after` on T-1315, since both write `summary.md`. The code and docs tasks cover I-0111's instances unchanged.

### 2026-10-08T00:27:49Z orchestrator
Approved.

- T-1314 removes I-0109's cause with tests that reproduce it (criterion 1). T-1315 closes I-0109 (criterion 2). T-1313's ADR records the remedy.
- The remedy is the leave-out pattern of ADR-0122 and S-0318, narrowed to `narrative.state`, so hand edits and flai defects on other narratives stay recorded. The broader option goes into the ADR as the alternative, as you say.
- `after: S-0318` is right: both change `ScopeToStory` and the same convention line.
- T-1315 and T-1316 share no path. The 24m forecast and the 37.5 USD a week value stand.
- Your contingent task for I-0111 waits on the operator's answer on TH-0346. If they merge S-0325 in, add it as you describe.

Source: wip/kanban/stories/S-0323-flai-check-finds-narrative-state-outside-the-story-at-close-out.md

### 2026-10-08T00:27:50Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories
