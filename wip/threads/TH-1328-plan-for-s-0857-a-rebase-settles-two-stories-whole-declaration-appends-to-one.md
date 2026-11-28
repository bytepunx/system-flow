---
id: TH-1328
title: "Plan for S-0857: a rebase settles two stories' whole-declaration appends to one Go file"
anchor:
  path: wip/kanban/stories/S-0857-a-rebase-stops-on-a-go-test-file-where-two-stories-each-appended-a-test-at-the.md
  item: S-0857
status: resolved
participants: [planner-S-0857, claude-concierge]
created: 2026-11-27T23:59:12Z
updated: 2026-11-28T00:01:52Z
---

# TH-1328 Plan for S-0857: a rebase settles two stories' whole-declaration appends to one Go file

On wip/kanban/stories/S-0857-a-rebase-stops-on-a-go-test-file-where-two-stories-each-appended-a-test-at-the.md.

## Entries

### 2026-11-27T23:59:12Z planner-S-0857
S-0857 is planned: 15 touches, 6 tasks in 5 layers, forecast 1h30m, cost of delay 37.50 USD a week.

## Proposed fix

Every rebase flai runs settles a stop on a `.go` file when each side only adds whole top-level declarations at the same place. It keeps ours, then theirs. The merge reads the three stages, not git's markers. Where a side changes anything else, the stop waits as it does today. The cases that wait:

- a change to the import block
- both sides declaring the same name
- a result that does not parse or is not gofmt-clean

The change goes in `issues.FilesIn`, which stream sync, task done, accept, the merge queue's batch, the rebase and resolve job, and the trial merge all share.

## Tasks and layers

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-3623 ADR refining ADR-0175 | none |
| 2 | T-3625 `MergeGo` in `flai/internal/storygit/gomerge.go` | T-3623, the rule it builds |
| 3 | T-3628 Go kind in `storygit.Files` and `issues.FilesIn`, with real-git sync and trial tests | T-3625, which it calls |
| 4 | T-3630 rebase and resolve job test reproducing I-0409, and `conflictText` | T-3628 |
| 4 | T-3633 `design/system/workflow.md` and `docs/users/flai.md` | T-3628 |
| 5 | T-3634 close I-0409 | T-3630, T-3633 |

## Assumptions

- It covers every `.go` file, not only `_test.go`. Recommended, since the rule is the same. Narrowing it to `_test.go` changes only the matcher in T-3628.
- Import-block changes stay unsettled. Both sides adding imports would need its own merge, and it is not in the one instance.
- A new ADR is needed, because ADR-0175 says a shared path that is not markdown is never merged. Accepted ADRs are immutable.
- The story keeps the folder touch `design/adrs`, since the ADR's number is not known until it is written.
- The story has no `after`.

No split, merge, or drop proposed.

### 2026-11-28T00:01:51Z claude-concierge
For the operator, who delegated completing drafts and queueing ready work in this session: the plan stands as drafted, assumptions included.

### 2026-11-28T00:01:52Z claude-concierge
Resolved: plan reviewed and the story finalized on the operator's delegation
