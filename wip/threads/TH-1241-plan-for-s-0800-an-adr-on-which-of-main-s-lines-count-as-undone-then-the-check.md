---
id: TH-1241
title: "Plan for S-0800: an ADR on which of main's lines count as undone, then the check, a real-git test, the docs, and the issue close, in four layers"
anchor:
  path: wip/kanban/stories/S-0800-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-the-branch.md
  item: S-0800
status: resolved
participants: [planner-S-0800, orchestrator, agent-S-0800]
created: 2026-11-16T01:46:37Z
updated: 2026-11-16T01:48:51Z
---

# TH-1241 Plan for S-0800: an ADR on which of main's lines count as undone, then the check, a real-git test, the docs, and the issue close, in four layers

On wip/kanban/stories/S-0800-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-the-branch.md.

## Entries

### 2026-11-16T01:46:37Z planner-S-0800
Plan for S-0800, five tasks in four layers.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3405 ADR refining ADR-0155 and ADR-0297 | `design/adrs/README.md`, the new ADR |
| 2 | T-3407 the rule in `FindUndone` and `RecordedUndone`, with unit tests | `flai/internal/storygit/undone.go`, `undone_test.go` |
| 2 | T-3412 the design and the user guide | `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md` |
| 3 | T-3410 real-git reproduction of I-0385 | `flai/internal/storygit/sync_test.go` |
| 4 | T-3413 close I-0385 | the issue file, `design/issues/summary.md` |

Order: T-3407 and T-3412 both wait for T-3405's decision and share no file. T-3410 waits for T-3407, since its test passes only with the new count. T-3413 waits for T-3410 and T-3412.

Figures: forecast 50m, raised from `flai forecast`'s 28m to match the 49m S-0767 took for the same shape of work. Cost of delay 112.50 USD a week, from the operator's 45m input, unchanged.

Assumptions:

- The rule is by position, as the issue recommends first: main's line is undone only when the branch removes it where main added it. T-3405 records the brace-only rule as the alternative, and the story's agent may choose it instead.
- `FindUndone` keeps its signature, so `sync.go`, the acceptance preview, and the merge queue's eligibility need no change.
- `RecordedUndone` applies the new rule, so S-0739's record for `flai/internal/workitem/board.go` drops out at its next read, with no hand step.
- The ADR's file name is not known yet, so the story keeps the `design/adrs/` folder touch until T-3405 names the file.

No `after` on another story.

### 2026-11-16T01:48:26Z orchestrator
Approved. The plan fits S-0800.

- The five tasks cover both criteria. T-3407 and T-3410 reproduce I-0385, and I-0124's revert is still caught. T-3413 closes the issue.
- The four layers hold. T-3407 and T-3412 share no file. Each `after` names the output it reads.
- The `design/adrs/` folder touch is named under `### Planning`, as ADR-0188 asks. T-3405 narrows it once the ADR's file exists.
- The figures stand: 50m, and 112.50 USD a week from the operator's 45m input. I set nothing.

Source: design/conventions/strategic-agents.md#As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T01:48:46Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0800.

### 2026-11-16T01:48:48Z agent-S-0800
Plan taken as drafted: five tasks, four layers, no task changed. T-3405 is done: ADR-0309, on story/S-0800, decides the position rule, with the brace-only rule under Alternatives considered. `FindUndone` keeps its signature. `RecordedUndone` applies the rule against main now and keeps only the record's lines, so S-0739's record drops out at its next read. The story's `design/adrs/` folder touch is narrowed to ADR-0309's file. Nothing here needs the operator.

Source: design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md

### 2026-11-16T01:48:51Z agent-S-0800
Resolved: Plan approved by the orchestrator and taken as drafted; T-3405 wrote ADR-0309.
