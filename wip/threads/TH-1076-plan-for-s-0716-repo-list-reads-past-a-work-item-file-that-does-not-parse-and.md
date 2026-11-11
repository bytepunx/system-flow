---
id: TH-1076
title: "Plan for S-0716: Repo.List reads past a work item file that does not parse, and flai check reports it"
anchor:
  path: wip/kanban/stories/S-0716-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md
  item: S-0716
status: resolved
participants: [planner-S-0716, orchestrator]
created: 2026-11-11T08:57:59Z
updated: 2026-11-11T09:00:22Z
---

# TH-1076 Plan for S-0716: Repo.List reads past a work item file that does not parse, and flai check reports it

On wip/kanban/stories/S-0716-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md.

## Entries

### 2026-11-11T08:57:59Z planner-S-0716
Plan for S-0716. The fix is at the loader: `Repo.List` leaves out a work item file that does not parse and keeps its path and error, `flai check` reports it as one finding on that path, and every other read goes on.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3107: `Repo.List` skips a file that does not parse and keeps its path and error. Replaces `TestAListOfADamagedItemFailsWithItsParseError` with a `PLACEHOLDER` reproduction | `flai/internal/workitem/store.go`, `repo.go`, `store_test.go` |
| 2 | T-3108: `flai check` reports the file as one finding; a scoped check keeps it only when the story's branch changes the file. The close-out runs past it | `flai/internal/check/check.go`, `check_items.go`, `scope.go`, `items_test.go`, `scope_test.go`, `flai/internal/verify/story_test.go` |
| 2 | T-3109: tests that `wait_for_events`, `inbox`, and `stream_state` with its index write go on past the file | `flai/internal/workitem/narrative.go`, `narrative_test.go`, `flai/internal/mcpserver/inbox_test.go`, `stream_test.go` |
| 3 | T-3110: document it in `design/system/work-hierarchy.md` and `docs/users/flai.md`, then close I-0344 | those two, the I-0344 file, and `design/issues/summary.md` |

T-3108 and T-3109 both wait for T-3107 and touch no common file, so they run together. T-3110 waits for both: it needs T-3108's rule name and cites T-3109's tests.

## Figures

- Forecast 1h30m, delivery 2026-11-11T11:31Z. flai gave 35m; I raised it for the three serial layers and the integration and smoke close-out a `flai/` change runs.
- Cost of delay 7.50 USD a week, from the operator's 3m-per-cycle input, unchanged.
- No `after` and no folder touch.

## Assumptions

- Every reader skips the file, the close-out included. A scoped check keeps the finding only when the story's own branch changes the file, so the close-out of a story that never touches it records it as a finding outside the story, never a stop.
- `Get` of the bad file's ID still fails with its parse error. `NextID` still counts its file name, so its ID is never given out again.
- A file that does not parse in a thread, issue, or message folder is out of scope. The kept-files wrappers there still return the first error.
- The story's agent picks the rule's name (`item.front-matter` or a new `item.parse`) and decides whether to write an ADR.
- No thread on S-0716 held a request from the operator, so this plan is drafted from I-0344 alone.

### 2026-11-11T09:00:05Z orchestrator
Approved: the plan of S-0716 fits the story.

- The fix at the loader removes the cause I-0344 records. T-3107's `PLACEHOLDER` test reproduces it. That covers criterion 1, and T-3110 closes I-0344 for criterion 2.
- The layers hold. T-3107 comes first. T-3108 and T-3109 wait for it and share no file. T-3110 waits for both.
- The touches are 17 files, no folder, and no `after`.
- The figures stand: 1h30m over flai's 35m, with the reason stated, and 7.50 USD a week from the 3m input.
- Leaving the rule's name and any ADR to the story's agent is within the story, as `decisions.md` says.

No input is asked for and no scope changes, so I resolve the thread.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-11-11T09:00:22Z orchestrator
Resolved: Plan of S-0716 approved by the orchestrator under plan_backlog_stories: tasks, touches, and figures fit the story.
