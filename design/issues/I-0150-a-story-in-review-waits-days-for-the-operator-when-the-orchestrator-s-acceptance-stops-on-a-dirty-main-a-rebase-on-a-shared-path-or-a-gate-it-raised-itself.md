---
id: I-0150
title: A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself
class: efficiency
status: closed
count: 1
cost: 96h
first_reported: 2026-10-15T22:48:32Z
last_reported: 2026-10-15T22:48:32Z
updated: 2026-10-23T05:41:17Z
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
          cache_read: 204810
          cache_write: 23259
          cost: 0.0083
---

# I-0150 A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself

## Description
A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself

## Instances

### 2026-10-15T22:48:32Z
Report: design/analysis/2026-10-15-bottlenecks.md.
Three stories sat in review from 2026-10-08/09 to 2026-10-12T14:33Z (259 h together) although the orchestrator had accept_reviews and was running: each acceptance stopped on something only the operator could clear, and the board delivered one story in three days.

## Impact
flai stats W42 review.total_seconds 973,317 over 58 stories, 933,886 of it on S-0341 (361,859 s), S-0342 (321,146 s), and S-0343 (250,881 s), each awaited alex. accept_reviews on since 2026-10-07; the orchestrator declined S-0341 on TH-0387 for a gate never built, and S-0342/S-0343 on S-0338's uncommitted acceptance (TH-0389, TH-0391, TH-0394), then flai accept S-0342 stopped mid-rebase on design/system/flai-cli.md, a claims.shared path (TH-0397; coordination W42 met_accept 1). cfd done 318 on 2026-10-08/09, 319 on 10-10/11, 350 on 10-12.

- time_lost_per_cycle: 96h

## Remediation

Found by the analysis in [2026-10-15-bottlenecks.md](../analysis/2026-10-15-bottlenecks.md).

Story S-0467 remediates this issue, created from it at 2026-10-21T11:57:59Z.
Closed 2026-10-23T05:41:17Z: S-0467, by ADR-0189, fixes each of the three causes.

- A dirty main: `flai accept --by orchestrator` leaves the main checkout's uncommitted files outside wip/ out of the acceptance, as the merge queue's landing does. Only the `dirty_main` blocker refuses it, when the story's change would overwrite one of those files.
- A stopped acceptance: the orchestrator finishes one whose move to done it made itself (S-0338's case).
- A rebase on a shared path it cannot settle: flai aborts the rebase and sends the story back to in-progress, with the paths and one thread. flai serve starts the story's agent on that thread (TH-0397's case).
- A gate it raised itself: the preview names a thread that only the orchestrator wrote on as `own_thread_open`, and the orchestrator resolves it with a reason before it accepts (TH-0387's case).
- Tests: `accept_orchestrator_dirty_test.go`, `accept_conflict_test.go`, `accept_own_thread_test.go`, `serve/agents_test.go`.
