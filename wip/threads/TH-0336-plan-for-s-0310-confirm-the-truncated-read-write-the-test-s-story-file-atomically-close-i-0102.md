---
id: TH-0336
title: "Plan for S-0310: confirm the truncated read, write the test's story file atomically, close I-0102"
anchor:
  path: wip/kanban/stories/S-0310-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md
  item: S-0310
status: resolved
participants: [planner-S-0310, orchestrator]
created: 2026-10-07T23:31:09Z
updated: 2026-10-07T23:31:22Z
---

# TH-0336 Plan for S-0310: confirm the truncated read, write the test's story file atomically, close I-0102

On wip/kanban/stories/S-0310-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md.

## Entries

### 2026-10-07T23:31:09Z planner-S-0310
S-0310 is planned: three tasks, three layers, one after the other.

## Tasks

| Layer | Task | Nature | Touches | After |
|-------|------|--------|---------|-------|
| 1 | T-1279 Surface the tool's error in TestWaitForWorkAcrossAFolder and confirm the cause under load | research | `flai/internal/mcpserver/folder_test.go` | none |
| 2 | T-1280 Write the story file atomically in readyStoryIn, with a test that reproduces the truncated read | remediation | `flai/internal/mcpserver/folder_test.go` | T-1279 |
| 3 | T-1281 Close I-0102 with what fixed it | remediation | I-0102's file, `design/issues/summary.md` | T-1280 |

T-1280 waits for T-1279 because the fix depends on the cause it confirms, and both change `folder_test.go`. T-1281 waits for T-1280 because the reason names the fix.

## Proposed cause

- `handed over: map[]` means the second wait_for_work call answered a tool error, not a timeout. The test drops the error text.
- `readyStoryIn` rewrites the story with `os.WriteFile`, which truncates first. The held wait polls every 20ms.
- A poll that reads the truncated file fails the parse in `workitem`'s `store.list`, and `folderWaits.hold` returns that error.
- flai's own writes use `atomicfile.WriteFile`, so only the test helper has this window. A loaded host widens it.

## Figures

- Forecast 1h, delivery 2026-10-08T06:15:00Z. flai gave 12m; raised because reproducing a failure under load takes repeated package runs (178s each under load in I-0102).
- Cost of delay value 25 USD a week, from the 10m input, unchanged.

## Assumptions

- The fix is in the test helper. `folder.go` and `server.go` stay out of the touches.
- If T-1279 refutes the cause, the story's agent rewrites T-1280 for the cause it finds.

## Proposed: one story to add, not drafted

A held wait_for_work or wait_for_events fails at the first read error. A story file half-written by hand, for example from an editor that does not write atomically, would end an agent's wait the same way. A story could make the holds look again on the next poll and fail only when the error persists. That changes server behaviour, so I have not drafted it. Say if you want it drafted.

### 2026-10-07T23:31:22Z orchestrator
Approved. The three tasks cover both criteria:

- T-1279 confirms the cause.
- T-1280 removes it with a test that reproduces the truncated read (criterion 1).
- T-1281 closes I-0102 (criterion 2).

The serial layers are explained: T-1280 needs T-1279's finding, and both change `folder_test.go`. The fix stays in the test helper, so no recorded decision changes. The 1h forecast and the 25 USD a week value, from flai's input from the issue, stand.

On the story you propose, where held waits tolerate a half-written file: it is new scope and a change to server behaviour, so it is the operator's to ask for. I am not drafting it, and this plan does not depend on it. Operator: if you want it, say so on S-0310 or ask the planner for it.

Source: wip/kanban/stories/S-0310-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md

Resolved: Plan approved by the orchestrator under plan_backlog_stories; the proposed extra story is left to the operator
