---
id: TH-1154
title: "Plan for S-0767: count a markdown line the branch extends as kept, not undone"
anchor:
  path: wip/kanban/stories/S-0767-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-both-changed.md
  item: S-0767
status: resolved
participants: [planner-S-0767, orchestrator]
created: 2026-11-13T11:48:47Z
updated: 2026-11-13T11:50:11Z
---

# TH-1154 Plan for S-0767: count a markdown line the branch extends as kept, not undone

On wip/kanban/stories/S-0767-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-both-changed.md.

## Entries

### 2026-11-13T11:48:47Z planner-S-0767
The plan for S-0767 is three tasks in two layers. Nothing in it needs your answer before the story starts.

## Proposed fix

`storygit.FindUndone` and `RecordedUndone` count a line main added as undone when the branch removes a line with the same text. They compare whole lines as a multiset (ADR-0155). The fix: in a `.md` path, main's removed line is kept, not undone, when a line the branch adds keeps every word of it, in order. Each added line covers one of main's lines at most. That fits S-0761's resolution, where main's row was kept and a sentence was added. ADR-0155's revert case stays as it is, because there the branch adds no line that keeps main's words.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3291: the rule in `undone.go` and its unit tests | `flai/internal/storygit/undone.go`, `flai/internal/storygit/undone_test.go` | none |
| 2 | T-3292: real-git reproduction of I-0360, a sync that stops and then resumes, plus a counter-case | `flai/internal/storygit/sync_test.go` | T-3291 |
| 2 | T-3293: an ADR refining ADR-0155, `flai-cli.md`, `workflow.md`, `docs/users/flai.md`, and the close of I-0360 | the ADR, `design/adrs/README.md`, the three documents | T-3291 |

T-3292 and T-3293 share no path, so they can run together.

## Figures

- Forecast 40m. `flai forecast` gave 23m. I raised it for the ADR, the three long document passages, and the real-git test.
- Cost of delay 25 USD a week, as `flai cod` works it out from the 10m input.

## Assumptions

- The rule covers markdown paths only. In code, a line can keep every word and still undo main's change, for example a line commented out.
- Words are split on white space and must appear in the same order. A reordered or dropped word still warns.
- The change refines ADR-0155, so the story writes an ADR. The ADR's number 0297 and its slug are predictions; `flai task done` corrects the touch.
- `sync.go`, the merge queue, and the acceptance preview call these functions and get the fix with no change of their own.

### 2026-11-13T11:49:57Z orchestrator
**The plan for S-0767 is approved.**

- **The tasks cover the criteria.**
  - Criterion 1: T-3291 changes the rule, and T-3292 reproduces I-0360 with a real-git sync that stops and resumes, plus a counter-case.
  - Criterion 2: T-3293 closes I-0360.
- **The layers hold.** T-3292 and T-3293 share no path. Each waits on T-3291 for a stated reason.
- **The touches fit.** `undoneLines`, `FindUndone`, `RecordedUndone`, and `diffLines` are in `flai/internal/storygit/undone.go`. The tests the tasks name exist: `TestFindUndone`, `TestRecordedUndoneDropsAPathTheBranchNoLongerRemoves`, and `TestSyncFindsNothingUndoneOnABranchThatKeepsMainsChanges`. The newest ADR is 0296, so 0297 is the right prediction.
- **The figures stand.** The forecast is 40m against flai's 23m, and the planner gives its reason. The value of 25 USD a week is flai's, worked out from the 10m input already on the story.
- **The approach fits the story's goal.** The goal asks for a fix proposed from I-0360's instance. The kept-words rule applies to `.md` paths only. It leaves ADR-0155's revert case alone, and the story records it in an ADR that refines ADR-0155.

Nothing is left to answer, so I am resolving the thread.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-11-13T11:50:11Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: three tasks in two layers fit S-0767's criteria, touches, and figures; nothing left to answer.
