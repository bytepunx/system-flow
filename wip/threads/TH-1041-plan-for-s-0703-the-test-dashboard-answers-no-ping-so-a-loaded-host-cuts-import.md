---
id: TH-1041
title: "Plan for S-0703: the test dashboard answers no ping, so a loaded host cuts import.run; three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0703-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md
  item: S-0703
status: open
participants: [planner-S-0703]
created: 2026-11-10T10:35:18Z
updated: 2026-11-10T10:35:18Z
---

# TH-1041 Plan for S-0703: the test dashboard answers no ping, so a loaded host cuts import.run; three tasks in three layers

On wip/kanban/stories/S-0703-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md.

## Entries

### 2026-11-10T10:35:18Z planner-S-0703
S-0703 is planned: three tasks in three layers, each waiting on the one before. Its touches, forecast, and cost of delay are written, with the reasons in Notes › Planning.

No thread named "asked" was on S-0703, and it had no plan before, so this is its first plan.

Proposed cause:

- I-0267's only instance ran at 13346d1 (2026-11-07T02:07Z). S-0560's drain fix (6920a061c, T-2556) landed after it, at 10:05Z, so that halt path is closed.
- A second path remains. `flai/internal/channel/channeltest/dashboard.go` reads only inside `Ask` and `Next`, so it answers no ping. The serve tests keep the client's default ping: every 4 s, with a 4 s timeout. On a host at load 56, a gap with no read cancels the connection. `conn.CloseNow` then ends it without a close frame, which reads as "failed to read frame header: EOF".

Tasks:

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-3042: a serve test with a short `PingEvery` and a pause reproduces the EOF | `flai/internal/serve/candidates_test.go` | none |
| 2 | T-3044: the test dashboard reads its connection for its whole life, so it answers pings | `flai/internal/channel/channeltest/dashboard.go` | T-3042 |
| 3 | T-3047: the candidate tests pass repeatedly under load; I-0267 is closed | `design/issues/I-0267-…`, `design/issues/summary.md` | T-3044 |

Figures:

- Forecast: 1h30m, up from flai's 6m. flai's figure rests on size 2 alone. The story also runs the serve package under load and a close-out with the integration tier.
- Cost of delay: 22.50 USD a week, from your 9m-per-cycle input. It stands as flai worked it out.

Assumptions:

- The fix stays in test code. Production `flai/internal/channel` ping behaviour is unchanged: an unanswered ping still means the dashboard is gone.
- If T-3042 cannot reproduce the EOF through the ping, the story's agent asks here before T-3044.
- The `PingEvery: time.Hour` workaround in `flai/internal/hostapi/components_test.go` stays. Dropping it is optional cleanup, not in scope.

Nothing here needs your answer unless you want the fix in production code instead.
