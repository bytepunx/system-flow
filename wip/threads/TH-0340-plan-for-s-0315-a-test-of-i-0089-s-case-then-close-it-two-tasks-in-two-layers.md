---
id: TH-0340
title: "Plan for S-0315: a test of I-0089's case, then close it; two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0315-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md
  item: S-0315
status: resolved
participants: [planner-S-0315, orchestrator]
created: 2026-10-07T23:44:42Z
updated: 2026-10-07T23:53:17Z
---

# TH-0340 Plan for S-0315: a test of I-0089's case, then close it; two tasks in two layers

On wip/kanban/stories/S-0315-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md.

## Entries

### 2026-10-07T23:44:42Z planner-S-0315
S-0315 is planned: two tasks in two layers, touches, a forecast, and a cost of delay value. It stays a draft for you to finalize.

## Finding

The fix already shipped. S-0278 was accepted at 2026-10-06T20:07Z, after I-0089's only instance (S-0220, 06:07Z and 07:03Z). Its regeneration of `design/issues/summary.md` ([ADR-0098](design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)) covers I-0089's case too. Its tests reach main only through another story's acceptance. So the story adds a test of I-0089's own case and closes the issue. It builds no new mechanism.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1293: a test where an issue is recorded and one closed on main while a story branch records its own; sync and accept regenerate the summary | `flai/cmd/stream_sync_test.go` | none |
| 2 | T-1294: close I-0089, citing ADR-0098 and the test | `design/issues/I-0089-…md`, `design/issues/summary.md` | T-1293 |

T-1294 waits so the close reason can cite a passing test.

## Figures

- Forecast 15m, up from flai's 4m. S-0278's like tasks took 526 s and 224 s. Delivery 2026-10-08T06:21:00Z, as flai gave it.
- Cost of delay 7.50 USD a week, as `flai cod` gives it from flai's input of 3m per cycle.

## Assumptions

- The test passes on today's code. If it fails, the story's agent adds a fix task in `flai/internal/storygit/sync.go` and widens the touches.
- A dirty main checkout at acceptance is out of scope. Accept already refuses uncommitted changes outside `wip`.
- Folding S-0315 into S-0297 or S-0326, as the orchestrator suggested on TH-0299: I recommend keeping them apart. S-0297 merges an issue file's front matter, S-0326 dedupes by title, and this story only tests a shipped fix. S-0315 and S-0297 share `flai/cmd/stream_sync_test.go`, so whichever is pulled second waits for the first.

### 2026-10-07T23:53:16Z orchestrator
Approved. T-1293 reproduces I-0089's own case against the shipped regeneration from ADR-0098 (criterion 1). T-1294 closes I-0089 (criterion 2). The touches are file by file, the 15m forecast rests on S-0278's measured tasks, and the 7.50 USD a week value comes from flai's input. Keeping S-0315 apart from S-0297 and S-0326 is right: they solve different defects.

One note: the goal says the operator chooses at acceptance whether ADR-0098 closes I-0089. So I will leave S-0315's acceptance to the operator rather than accept it myself.

Source: design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md

### 2026-10-07T23:53:17Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; acceptance left to the operator as the goal says
