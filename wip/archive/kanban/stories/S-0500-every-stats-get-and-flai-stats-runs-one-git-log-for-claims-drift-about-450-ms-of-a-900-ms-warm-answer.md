---
id: S-0500
type: story
nature: improvement
title: Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer
status: done
owner: alex
created: 2026-10-26T04:01:29Z
updated: 2026-10-26T04:57:42Z
transitions:
  - to: ready
    at: 2026-10-26T04:28:07Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-26T04:28:17Z
    by: agent-S-0500
  - to: review
    at: 2026-10-26T04:43:51Z
    by: agent-S-0500
  - to: done
    at: 2026-10-26T04:57:42Z
    by: orchestrator
tags: [flai]
topics: [metrics, performance, cli]
touches: [flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, flai/internal/statsread/statsread_test.go, flai/internal/hostapi/timing_test.go, design/system/server-performance.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md, design/issues/summary.md, design/issues/I-0209-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 953
  priced_by: harness
  turns:
    - day: 2026-10-26
      ceremony: 5
      hand_edits: 2
      work: 31
  models:
    - model: claude-opus-5-5
      input: 196
      output: 55533
      cache_read: 9318995
      cache_write: 346661
      cost: 5.2965
  strategic:
    - kind: orchestrator
      seconds: 1426
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 254
          output: 4656
          cache_read: 58573711
          cache_write: 66023
          cost: 12.8565
        - model: claude-sonnet-5-5
          input: 8
          output: 49
          cache_read: 90646
          cache_write: 39263
          cost: 0.0888
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-10-26T04:01:29Z
  value: 0.04
  by: planner-S-0500
  at: 2026-10-26T04:26:07Z
forecast:
  duration: 1h
  delivery: 2026-10-26T19:11:00Z
  basis: "flai forecast's 16m (87 s per unit of size over 18 done improvement stories on claude-opus-5-5 in the medium band, size 11) raised to 1h for a kept-history memo with counting-runner tests, a Read-level test, a cold and warm measurement of stats.get on this repository with two binaries, and the design and docs updates, as S-0471's 18m became 52m of progress; delivery is flai's 2026-10-26T18:27Z, 22nd in the pull order, plus the 44m added"
  by: planner-S-0500
  at: 2026-10-26T04:26:07Z
finalized:
  by: orchestrator
  at: 2026-10-26T04:27:59Z
---
# S-0500 Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer

## Goal

This story remediates [I-0184](../../../design/issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md), "Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0184 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0184 is closed with `flai issue close I-0184 --reason` saying what fixed it

## Tasks
- T-2136 storygit.Committed keeps the stories' commits per repository root for the life of the process, keyed by the branch tips, and reads the log again only when a tip moves
- T-2137 stats.get and flai stats are measured again on this repository with the kept commits, and server-performance.md, flai-cli.md, and the flai stats guide record the fix
- T-2138 I-0184 is closed with the reason naming the kept commits, and the story's criteria are ticked

## Notes

Cost of delay inputs set by flai from I-0184. time_lost_per_cycle 1s: 450ms per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T12:50:20Z, 4.6 days before this story; under one cycle counts as one).

### Planning

Cause, from the code: `storygit.Committed` in `flai/internal/storygit/committed.go` runs one `git log --name-only` over the main branch and every `story/` branch at every call, and `statsread.Read` calls it at every `stats.get` and `flai stats`. Nothing keeps the result between calls, although the stories' commits change only when a branch tip moves. On this repository the log prints 40,938 lines and took 0.67 s on 2026-10-26. `release.Pending` solved the same shape with a history kept per repository root for the life of the process, keyed by `git show-ref` (S-0157, `flai/internal/release/history.go`).

Proposed solution: keep the stories' committed files per repository root for the life of the process, keyed by `git show-ref --heads`, the tips of the main branch and every story branch, which costs a few milliseconds. While the key is unchanged, `Committed` answers from what it kept, deep-copied, with one git process in place of the log; when a tip moves, it reads the whole log again. A whole read on change keeps the code simple and correct across the rebases `flai stream sync` makes, which move story tips to commits the kept history does not reach; reading only the commits added since the kept tips is a possible follow-up, not this story. `flai accept`'s own call after the acceptance commit sees a changed key and reads again, as it must.

Touches, all files, no folder touch:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/storygit/committed.go` | design | I-0184 names `storygit.Committed`, defined here; the kept commits and their key go here |
| `flai/internal/storygit/committed_test.go` | layout | its tests; the counting-runner tests that reproduce the cause go here |
| `flai/internal/statsread/statsread_test.go` | co-change | changed with `statsread.go` in 75 percent of its commits; a test that `Read` starts no log while nothing changed |
| `flai/internal/hostapi/timing_test.go` | layout | names the phases each read marks; `stats.get` may join it with `exec.git.show-ref` |
| `design/system/server-performance.md` | design | holds I-0184's figures; records the fix and the new measurement |
| `design/system/flai-cli.md` | design | § Internal structure describes the kept release history; the kept commits are described beside it |
| `docs/users/flai.md` | co-change | the `flai stats` section says the touches drift reads git; it says when git is read again |
| I-0184's file under `design/issues/` | declared | the goal links it; criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

`flai touches suggest` from the two code files also listed `flai/internal/metrics/metrics.go` (88 percent), `flai/cmd/stats.go` (50), `design/system/metrics.md` and `flai/cmd/check_stats_test.go` (38), and `docs/users/flai-reference.md` (25). None was taken: the fix changes no metric, no command flag, and no output, so the metrics, the stats command, its contract tests, and the generated reference stay as they are. `flai/internal/statsread/statsread.go` is not touched either: the memo lives in `storygit`, so every caller of `Committed` gains it.

Forecast: flai gives 16m, 87 s per unit of size over 18 done improvement stories on claude-opus-5-5 in the medium band, size 11. Raised to 1h: the story needs a kept history with counting-runner tests, a `Read`-level test, a cold and warm measurement on this repository with two binaries, and three documents, and S-0471, the same shape of work, took 52m of progress against a forecast of 18m. Delivery moves from flai's 2026-10-26T18:27Z by the same 44m.

Cost of delay: `flai cod` gives 0.04 USD a week from flai's input of 1s lost per 168h cycle. The value stands as computed. The input counts one occurrence, while every charts page load, every item change while the charts page is open, and every agent's `flai stats` pays the 450 ms; the plan's thread says so.

Tasks, in three layers:

1. T-2136: the kept commits in `storygit.Committed`, with the tests that reproduce the cause
2. T-2137, after T-2136: measure on this repository and record it in `server-performance.md`, `flai-cli.md`, and the `flai stats` guide
3. T-2138, after T-2137: close I-0184 and tick the criteria

### Accepted by the orchestrator

- Verified: 23c2a1a7955fde86c105ec264d37fbb6756c62b2
- Batch: 20261026-044440
- Enqueued verified: 23c2a1a7955fde86c105ec264d37fbb6756c62b2
- At: 2026-10-26T04:57:42Z

Verdict: Accept S-0500: both criteria are met, the fix and its counting-runner tests are sound, and nothing under .claude/ changed.
- 1: flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, flai/internal/statsread/statsread_test.go, design/system/server-performance.md, design/system/flai-cli.md, docs/users/flai.md
- 2: design/issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md, design/issues/summary.md
