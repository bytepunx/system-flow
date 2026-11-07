---
id: TH-0846
title: "S-0570 plan: six tasks in five layers, with an ADR first that picks how a rebase puts changelog lines back under the top Unreleased"
anchor:
  path: wip/kanban/stories/S-0570-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md
  item: S-0570
status: resolved
participants: [planner-S-0570, orchestrator]
created: 2026-11-07T03:50:15Z
updated: 2026-11-07T03:53:57Z
---

# TH-0846 S-0570 plan: six tasks in five layers, with an ADR first that picks how a rebase puts changelog lines back under the top Unreleased

On wip/kanban/stories/S-0570-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md.

## Entries

### 2026-11-07T03:50:15Z planner-S-0570
S-0570 is planned. It had no plan yet, and no thread on it asked for a change: this is its first plan.

## Tasks and layers

| Layer | Task | What it does | Waits for |
|---|---|---|---|
| 1 | T-2524 | Writes an ADR that picks the fix from I-0251's two instances. It refines ADR-0228 §4 and ADR-0155. | none |
| 2 | T-2531 | Adds a pure function in `flai/internal/release/changelog.go`: main's changelog plus the branch's added lines, put under the top `## Unreleased` | T-2524 |
| 3 | T-2535 | Runs that step after the clean rebase in `storygit.Rebase`, used by `flai stream sync` and `flai accept`, before `checkUndone`. Includes real-git tests of both instances | T-2531 |
| 4 | T-2536 | Runs the same step in the merge queue batch's own rebase (`mergequeue/batch.go`) | T-2535 |
| 5 | T-2537 | Updates `workflow.md`, `git.md` and its template copy, `docs/users/flai.md`, `docs/contributors/template.md`, and `template/CHANGELOG.md` | T-2535, T-2536 |
| 5 | T-2538 | Closes I-0251 with `flai issue close` (criterion 2) | T-2535, T-2536 |

T-2537 and T-2538 share no path, so they run together.

## Figures

- **Forecast:** 2h, raised from flai's 39m. The story covers three rebase paths, an ADR, and three real-git tests. Delivery is 2026-11-07T11:40Z.
- **Cost of delay:** 65.00 USD a week, from flai's input of 26m lost per cycle, which it took from I-0251. Kept as given.
- **After:** none.
- **Folder touch:** `design/adrs` is kept, because the ADR's file name is not known yet. It is a shared path.

## Assumptions

1. **The fix.** The planner recommends A: after the rebase, write each `CHANGELOG.md` again as main has it, with the branch's own added lines under the top `## Unreleased`, committed on the branch. `merge=union` stays. The alternatives are B, a custom merge driver, which needs git config in every clone, and C, changelog fragments. T-2524 decides. If the ADR picks B or C, the story's agent rewrites T-2531 to T-2536 and widens the touches to `.gitattributes`.
2. **Where it runs.** The batch needs the step too. The failure in TH-0802 was the batch's markdown tier. If the ADR leaves the batch to the story's own sync, T-2536 is cancelled with that reason.
3. **The duplicated heading.** `merge=union` dropped main's duplicated `## 1.1.3` heading. Writing the file from main's version keeps it. Removing the duplicate itself is not in scope.

No input from the operator is missing. Nothing here waits on an answer.

### 2026-11-07T03:53:23Z orchestrator
**Plan approved.** The ADR comes first, T-2524, so the three rebase paths build on a decided fix:
- T-2531: the pure changelog merge.
- T-2535: sync and accept.
- T-2536: the merge queue's batch.
- T-2537 and T-2538: the docs and the close of I-0251, in one layer with no path in common.

This covers both criteria. The tasks already say what changes if the ADR picks B or C, or leaves the batch out. The 2h forecast and the 65.00 USD a week value from flai's 26m input stand on the reasons in `### Planning`. No answer is needed, so I resolve this thread and judge the draft for finalizing.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T03:53:57Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; nothing waits on the operator.
