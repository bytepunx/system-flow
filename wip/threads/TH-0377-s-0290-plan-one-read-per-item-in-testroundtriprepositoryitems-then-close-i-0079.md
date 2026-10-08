---
id: TH-0377
title: "S-0290 plan: one read per item in TestRoundTripRepositoryItems, then close I-0079"
anchor:
  path: wip/kanban/stories/S-0290-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md
  item: S-0290
status: resolved
participants: [planner-S-0290, alex]
created: 2026-10-08T08:42:45Z
updated: 2026-10-08T08:47:18Z
---

# TH-0377 S-0290 plan: one read per item in TestRoundTripRepositoryItems, then close I-0079

On wip/kanban/stories/S-0290-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md.

## Entries

### 2026-10-08T08:42:45Z planner-S-0290
**Recommendation:** raise S-0290's `time_lost_per_cycle` input from 6m to 1h47m, with `flai edit S-0290 --time-lost-per-cycle 1h47m` or on the dashboard. Then `flai cod` agrees with the value I set, 268 USD a week. The 6m counted one occurrence. I-0079 now has 9 occurrences of 5m each in 70.5 hours, which comes to 107m a week. I left the input alone because it is yours.

**Cause and fix.** In a story worktree, `Open` puts `wip/` in the main checkout. The test parses every item through `List`, then reads each file again and compares the two. All nine instances are an item rewritten between those two reads. `Repo.Save` writes atomically, so one read always gets a whole file. The test should parse and compare from that one read, and skip a file archived since the list.

**Tasks, in two layers:**

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1367 round-trips each item from one read of its file and skips a file gone since the list, with a reproduction test | `flai/internal/workitem/workitem_test.go` | none |
| 2 | T-1368 closes I-0079 with that fix as the reason | I-0079's file, `design/issues/summary.md` | T-1367 |

**Story figures:**

- Touches: those three files, no folders.
- Forecast: 30m, raised from flai's 14m for the integration tier the close-out runs. Delivery is 09:44Z.
- Cost of delay value: 268 USD a week.

**Assumptions:**

- I left out reading the worktree's own committed `wip/` instead. It would also be stable, but it checks a stale copy and needs a hook into `Open`.
- The fix changes only a test, so no design or user doc changes. No design document names the test.
- S-0321, in progress, also touches the I-0079 file. `design/issues` is in `claims.shared`, so that overlap holds nothing.
- The story was moved to ready while I planned it. I moved nothing.

### 2026-10-08T08:47:18Z alex
Resolved.
