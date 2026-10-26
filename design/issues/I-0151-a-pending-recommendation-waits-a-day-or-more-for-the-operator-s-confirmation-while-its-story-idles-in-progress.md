---
id: I-0151
title: A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress
class: efficiency
status: closed
count: 1
cost: 29h
first_reported: 2026-10-15T22:48:37Z
last_reported: 2026-10-15T22:48:37Z
updated: 2026-10-26T05:14:22Z
usage:
  strategic:
    - kind: analyzer
      seconds: 114
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 98
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23259
          cost: 0.0083
---

# I-0151 A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress

## Description
A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress

## Instances

### 2026-10-15T22:48:37Z
Report: design/analysis/2026-10-15-bottlenecks.md.
The orchestrator recommends within seconds, but a recommendation ends no wait (ADR-0091); the operator's confirmation took 25 to 113 hours on the four W42 waits, and S-0391 is waiting now on TH-0461 since 2026-10-15T05:09Z.

## Impact
flai stats W42 threads.total_seconds 1,003,576, threads.confirmed count 4 total_seconds 421,270 (mean 105,317 s each); waiting.longest: S-0232 TH-0324 405,614 s, S-0337 TH-0390 325,914 s, S-0380 TH-0445 94,433 s, S-0378 TH-0436 89,637 s, all awaited alex, each with the orchestrator's recommendation posted within a minute. S-0232 cycle_time_seconds 407,617 of which wait_threads_seconds 405,614. aging: S-0391 in progress 63,454 s over p85 on TH-0461's pending recommendation.

- time_lost_per_cycle: 117h

## Remediation

Found by the analysis in [2026-10-15-bottlenecks.md](../analysis/2026-10-15-bottlenecks.md).

Story S-0498 remediates this issue, created from it at 2026-10-26T04:01:18Z.
Closed 2026-10-26T05:14:22Z: Fixed by S-0498 (ADR-0205). flai serve now reminds the operator of each thread that awaits them, a pending recommendation or a question, once its last entry is older than threads.remind_after, and again every threads.remind_every. The reminder goes through a channel kept on the host (flai serve notify set --url/--command), and nothing is written on the thread, so no wait ends. This project sets 4h and 24h. flai serve also blocks a story whose agent ended asking on a thread, with the reason 'awaiting the operator on TH-nnnn', until the answer starts the agent again, so the idle shows on the board and counts as blocked time. TestARecommendationPending25HoursIsRemindedAtTheFirstLook in flai/internal/serve/remind_test.go and TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain in agents_test.go reproduce the two halves. S-0502 takes the design-choice waits with orchestration.permissions.confirm_after.
