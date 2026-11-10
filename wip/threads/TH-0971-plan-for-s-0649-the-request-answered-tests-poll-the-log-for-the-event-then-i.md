---
id: TH-0971
title: "Plan for S-0649: the request answered tests poll the log for the event, then I-0324 closes"
anchor:
  path: wip/kanban/stories/S-0649-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md
  item: S-0649
status: open
participants: [planner-S-0649]
created: 2026-11-10T06:58:52Z
updated: 2026-11-10T06:58:52Z
---

# TH-0971 Plan for S-0649: the request answered tests poll the log for the event, then I-0324 closes

On wip/kanban/stories/S-0649-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md.

## Entries

### 2026-11-10T06:58:52Z planner-S-0649
Plan for S-0649, two tasks in two layers.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2834 The request answered tests poll the log for the event, and a held-back log line reproduces the race | `flai/internal/channel/channel_test.go` | none |
| 2 | T-2835 Close I-0324 with what fixed it | `design/issues/I-0324-…md`, `design/issues/summary.md` | T-2834, whose fix and runs the reason names |

Cause: S-0611 (2465e6c) logs the channel client's "request answered" event after the answer's write. `TestEachRequestIsTimedBeneathTheTransport` reads the log once, as soon as the fake dashboard has the reply, so under load it reads before the event is logged. `TestTheSendTimeCountsTheWaitForTheWriteLock` polls with a deadline already and does not fail.

Fix: one polling helper in the test file, used by both tests. A log destination that holds back the "request answered" line reproduces the failure on every run without load.

Figures: forecast 35m, delivery 2026-11-10T08:30:00Z, raised from `flai forecast`'s 5m for the repeated test runs and an integration-tier close-out under load. Cost of delay 90.00 USD a week, from `flai cod` as it stands.

Assumptions:

- The fix changes no product code. The event belongs after the write, as S-0611 decided.
- The file's other one-shot log reads (`refusals`, `allowedUnsigned`, `gateRefusals`) and MCP's timing test are out of scope: they do not read an event logged after the reply.
- No `after` on another story: S-0649 reads no other story's output.
- The run was started as asked again, but no thread on S-0649 existed, so this is the story's first plan.

For the operator: the time_lost_per_cycle input, 36m, was worked out from 3 occurrences. I-0324 counts 5 now, which would be 60m and 150.00 USD a week. The input is yours, so I left it as it is.
