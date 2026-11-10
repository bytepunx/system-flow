---
id: TH-0971
title: "Plan for S-0649: the request answered tests poll the log for the event, then I-0324 closes"
anchor:
  path: wip/kanban/stories/S-0649-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md
  item: S-0649
status: resolved
participants: [planner-S-0649, orchestrator]
created: 2026-11-10T06:58:52Z
updated: 2026-11-10T07:01:30Z
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

### 2026-11-10T07:00:43Z orchestrator
The plan is approved.

- **Tasks.** T-2834 changes only `flai/internal/channel/channel_test.go`. T-2835 waits for it, because it closes I-0324 with what fixed it. The layers fit the two criteria.
- **Touches.** The test file, I-0324, and `design/issues/summary.md`. Leaving out `flai/internal/channel/channel.go` is right: S-0611 put the event after the write on purpose, and the fix changes only the test. S-0638's branch only bumps I-0324's count. It does not fix it, so the two stories do not duplicate each other.
- **Figures.** The 35m forecast, raised from flai's 5m, is explained by the repeated runs and the close-out under load. The 90 USD a week value is what `flai cod` gives.
- **The 36m input.** I left it as it is. flai set it, and I set no figure of my own. Under the `throughput` policy the ready column is ordered by forecast duration, so this input does not change where S-0649 goes. To count I-0324's 5 occurrences, the operator can run `flai edit S-0649 --time-lost-per-cycle 60m`, and the planner works out the value on its next run.

I am resolving this thread and finalizing S-0649.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T07:01:30Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0649.
