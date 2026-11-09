---
id: S-0614
type: story
nature: remediation
title: The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged
status: done
owner: alex
created: 2026-11-07T23:39:04Z
updated: 2026-11-09T10:17:12Z
transitions:
  - to: ready
    at: 2026-11-09T09:37:02Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T09:37:50Z
    by: agent-S-0614
  - to: review
    at: 2026-11-09T10:15:50Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:17:12Z
    by: orchestrator
tags: [flai, merge-queue]
topics: [merge-queue, markdown-lint, forge]
touches: [flai/internal/mdlint/quote.go, flai/internal/mdlint/quote_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, design/system/flai-cli.md, design/issues/I-0242-the-merge-queue-s-send-back-thread-quotes-go-test-s-summary-line-with-its-tabs.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2292
  priced_by: harness
  turns:
    - day: 2026-11-09
      test_runs: 1
      hand_edits: 1
      work: 38
  models:
    - model: claude-opus-5-5
      input: 158
      output: 43747
      cache_read: 18615555
      cache_write: 343562
      cost: 7.1776
  strategic:
    - kind: planner
      seconds: 298
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 54
          output: 16362
          cache_read: 3565621
          cache_write: 168908
          cost: 2.3918
    - kind: orchestrator
      seconds: 196
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 44
          output: 329
          cache_read: 1756253
          cache_write: 209696
          cost: 0.512
        - model: claude-sonnet-5-5
          input: 5
          output: 372
          cache_read: 58218
          cache_write: 20932
          cost: 0.0586
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5h
    by: flai
    at: 2026-11-07T23:39:04Z
  value: 750
  by: planner-S-0614
  at: 2026-11-09T09:35:26Z
forecast:
  duration: 1h30m
  delivery: 2026-11-09T16:00:00Z
  basis: "flai forecast's 22m (119 s per unit over 42 medium remediation stories, size 11) raised to 1h30m for four Go tasks and a close-out, which took 17 to 22 minutes alone on this host today; delivery allows for finalizing, promotion, and six stories ahead in the pull order"
  by: planner-S-0614
  at: 2026-11-09T09:35:26Z
finalized:
  by: orchestrator
  at: 2026-11-09T09:36:32Z
---
# S-0614 The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged

## Goal

This story remediates [I-0242](../../../design/issues/I-0242-the-merge-queue-s-send-back-thread-quotes-go-test-s-summary-line-with-its-tabs.md), "The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0242 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0242 is closed with `flai issue close I-0242 --reason` saying what fixed it

## Tasks
- T-2710 mdlint.Fenced quotes tool output in a code block the markdown lint passes
- T-2711 The merge queue's send-back quotes with mdlint.Fenced and, when its thread is still refused, sends the story back with a thread that leaves the findings out
- T-2712 The pull request send-back's fenced quote of a check's output uses mdlint.Fenced
- T-2713 flai-cli.md says the send-back threads quote lint-safe and how the queue's send-back goes on when its thread is refused, and I-0242 is closed

## Notes

Cost of delay inputs set by flai from I-0242. time_lost_per_cycle 5h: 5h per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-31T16:43:40Z, 7.3 days before this story).

### Planning

Proposed solution, from I-0242's one instance and the code as it stands on main:

- The tab itself is fixed already: `fenced` in `flai/internal/mergequeue/bisect.go` writes tabs as four spaces, with `TestTheFailingTiersOutputIsQuotedWithoutTabs`. What remains of the cause is that any quote of tool output the thread writer's lint refuses still wedges the queue: `SendBack` returns the refusal, `flai/internal/mergequeue/land.go` returns it at every look, and the story is never sent back.
- One lint-safe fenced quote, `mdlint.Fenced`, beside `QuoteBareURLs` in `flai/internal/mdlint/quote.go`: line ends as `\n`, tabs as spaces, ANSI escapes and other control characters dropped, trailing spaces trimmed, runs of blank lines folded, and a fence longer than any backtick run. `flai/internal/digest/seed.go` already does most of this for its own quotes.
- The merge queue's send-back quotes with it, and when the thread is still refused, opens it again with the findings left out and the log folder named, so the story is sent back and the batch goes on.
- The pull request send-back's fenced quote, `quote` in `flai/internal/serve/forgestatus.go`, has the same flaw (no tab handling at all) and quotes with it too.

Touches, each a file:

| Touch | From |
|-------|------|
| `flai/internal/mdlint/quote.go`, `flai/internal/mdlint/quote_test.go` | layout: where `QuoteBareURLs`, the one exported quote helper, lives |
| `flai/internal/mergequeue/bisect.go` | design: I-0242 names it; `SendBack`, `openFailing`, `failingText`, `fenced` |
| `flai/internal/mergequeue/bisect_test.go` | co-change: 8 of 12 commits with `bisect.go` (67%) |
| `flai/internal/serve/forgestatus.go` | layout: the pull request send-back's `quote`, the same fenced quote without tab handling |
| `flai/internal/serve/forgestatus_test.go` | co-change: 3 of 12 commits (25%) |
| `design/system/flai-cli.md` | design: § The merge queue pass and § Reading the pull request back describe the two send-back threads |
| `design/issues/I-0242-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

`flai/internal/mergequeue/land.go` is left out: the retry loop there stays as it is once `SendBack` no longer fails on a quote. No folder touch is kept. The story has no `after`: it reads nothing another open story writes.

Figures:

- Forecast: `flai forecast` gave 22m (119 s per unit of size over 42 done medium remediation stories, size 11). Raised to 1h30m: four tasks of Go with tests, and a close-out that took 17 to 22 minutes on its own on this host today (S-0603, S-0612, S-0613). Delivery 2026-11-09T16:00:00Z, past flai's 10:02, since the story is a draft still to be finalized and promoted, behind six stories in the pull order.
- Cost of delay: `flai cod` gave 750 USD a week from the operator's input, 5h lost per 168h cycle at 150 USD an hour. It stands: the tab is fixed, but a refused quote still wedges the whole queue and blocks every hand acceptance behind the ADR-0221 gate, as the instance shows.

### Accepted by the orchestrator

- Verified: 642e5bb2d3682c967ae3708bed213930ecbd56df
- At: 2026-11-09T10:17:12Z

Verdict: pass. Both criteria are met at 642e5bb; flai verify passed every step at that commit.
- 1: flai/internal/mdlint/quote.go, flai/internal/mdlint/quote_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, design/system/flai-cli.md
- 2: design/issues/I-0242-the-merge-queue-s-send-back-thread-quotes-go-test-s-summary-line-with-its-tabs.md, design/issues/summary.md
