---
id: S-0604
type: story
nature: improvement
title: board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing
status: done
parent: E-0033
owner: alex
created: 2026-11-07T22:40:34Z
updated: 2026-11-09T08:07:31Z
transitions:
  - to: ready
    at: 2026-11-08T00:23:34Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-08T02:06:43Z
    by: agent-S-0604
  - to: review
    at: 2026-11-09T08:01:21Z
    by: agent-S-0604
  - to: done
    at: 2026-11-09T08:07:31Z
    by: orchestrator
blocked:
  - from: 2026-11-08T02:44:08Z
    until: 2026-11-08T02:44:08Z
    reason: awaiting the operator on TH-0912
  - from: 2026-11-08T02:44:24Z
    until: 2026-11-08T02:58:16Z
    reason: awaiting the operator on TH-0912
  - from: 2026-11-08T02:58:40Z
    until: 2026-11-08T03:30:55Z
    reason: awaiting the operator on TH-0912
  - from: 2026-11-08T03:31:12Z
    until: 2026-11-09T07:25:05Z
    reason: awaiting the operator on TH-0912
  - from: 2026-11-09T07:25:20Z
    until: 2026-11-09T07:35:38Z
    reason: awaiting the operator on TH-0912
tags: [cli]
touches: [flai/internal/release/history.go, flai/internal/mergequeue/eligible.go, flai/internal/hostapi/hostapi.go, flai/internal/release/release.go, flai/internal/release/history_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/mergequeue/eligible_test.go, flai/internal/hostapi/reads_test.go, flai/internal/storygit/undone.go, design/system/server-performance.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3842
  priced_by: harness
  turns:
    - day: 2026-11-08
      ceremony: 1
      hand_edits: 4
      work: 64
    - day: 2026-11-09
      ceremony: 1
      test_runs: 2
      hand_edits: 2
      work: 99
  models:
    - model: claude-opus-5-5
      input: 572
      output: 101652
      cache_read: 94924574
      cache_write: 2117905
      cost: 34.0834
  strategic:
    - kind: orchestrator
      seconds: 816
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 70
          output: 389
          cache_read: 3935176
          cache_write: 553524
          cost: 1.1478
cost_of_delay:
  inputs:
    penalty_per_week: 40
    by: orchestrator
    at: 2026-11-08T00:03:30Z
  value: 40
  by: planner-S-0604
  at: 2026-11-08T00:10:53Z
forecast:
  duration: 1h
  delivery: 2026-11-08T03:07:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 5, behind S-0599, S-0601, S-0603, S-0606 and S-0607."
  by: flai
  at: 2026-11-08T01:46:55Z
finalized:
  by: orchestrator
  at: 2026-11-08T00:20:42Z
---
# S-0604 board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing

## Goal

Cause 14 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `board.get` has 16,641 slow calls in 4 days at p50 800 ms, mostly `release.pending` (438 ms) and `merge.queue` (279 ms). `release`'s `kept` mutex is one for every repository (`flai/internal/release/history.go`) and is held across a `git show-ref`, so the 13 projects' boards queue on it: 23 to 88 ms alone, 250 ms with 8 at once. `merge.queue` reads every thread even with no story in review (`flai/internal/mergequeue/eligible.go`), 185 to 630 ms. Each story's undone paths are read one after another (`flai/internal/hostapi/hostapi.go`).

## Acceptance criteria

- [x] `release`'s kept history takes one lock per repository root, and holds none across a git process another repository's call would wait for.
- [x] While flai serve's main-refs watcher has seen no change, `release.pending` starts no git process.
- [x] `merge.queue` reads no thread while no story is in review, and the undone paths of the stories are read at once.
- [x] 13 `board.get` calls at once on the operator's host, one per served project: the other 12 answer under 300 ms, this repository's remaining time is shown to be causes 13 and 16, and server-performance.md records the before and after.

## Tasks
- T-2659 The release package's kept history takes one lock per repository root, so one root's git processes never hold another root's call
- T-2666 merge.queue reads no thread and runs no git status while no story is in review
- T-2670 board.get reads the stories' undone paths at once on a bounded pool, and reads a story in review's once
- T-2677 flai serve's main-refs watcher tells the release package a root's refs are unchanged, so release.pending starts no git process until a ref moves
- T-2681 Measure 13 board.get calls at once before and after on the operator's host, and record them and the per-root history in the design

## Notes

- `publish.preview` holds the same lock and gains too.
- Criterion 4 was reworded on TH-0912 (operator, 2026-11-09, confirming the orchestrator's recommendation, as D-0465 did for S-0424). Its original words, every one of the 13 boards under 300 ms, moved to the draft S-0636, which waits for S-0602 (cause 13) and S-0606 (cause 16): this repository's board took 342 to 761 ms after S-0604, from its thread read with stories in review and its item list and holds, not from cause 14.

### Planning

Touches, file by file, and where each came from:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/release/history.go` | declared | The `kept` history and its lock (`kept`, `readHistory`, the `show-ref` key) |
| `flai/internal/mergequeue/eligible.go` | declared | `ComputeAt` reads `threads.List` before it looks for a story in review |
| `flai/internal/hostapi/hostapi.go` | declared | `board.get`'s serial `undoesMain` loop over the ready, in-progress, and review cards |
| `flai/internal/release/release.go` | layout | `PendingBatch` takes the lock and calls `readHistory`, `computePlan`, and `list` under it |
| `flai/internal/release/history_test.go` | layout | The kept history's tests, including the one that counts git processes |
| `flai/internal/serve/mergequeue.go` | layout | The main-refs watcher (`mainRefs`, `queuePass.start`, `changed`); it lacks `.git/refs/tags`, which the `show-ref --tags` key reads |
| `flai/internal/serve/mergequeue_test.go` | layout | The watcher's signal to `release` is tested beside it |
| `flai/internal/mergequeue/eligible_test.go` | co-change (10%) | No thread read with no story in review |
| `flai/internal/hostapi/reads_test.go` | co-change (10%) | `board.get`'s tests (`TestBoardGetCarriesTheMergeQueue`) |
| `flai/internal/storygit/undone.go` | layout | `RecordedUndone` runs one `git diff` per recorded path and rewrites the record; read at once, it must stay safe per story |
| `design/system/server-performance.md` | design | Criterion 4: the before and after under cause 14 |
| `design/system/flai-cli.md` | co-change (14%), design | Its `release` section says `show-ref` runs on every call while unchanged; that changes |

Co-change files left out: `flai/internal/hostapi/writes.go`, `flai/cmd/serve_actions.go`, `flaiover/src/lib/server/agent.ts`, `flai/internal/mcpserver/server.go`, and the rest of `flai touches suggest`'s list. They change with these files for other reasons, and nothing in this story's criteria reaches them. No folder touch.

No `after`: the story reads nothing another open story writes. S-0602 keeps threads parsed and S-0606 keeps the item list; each makes `board.get` cheaper on its own, and the trial merge judges the shared `flai/internal/hostapi/hostapi.go`.

Forecast: 1h, raised from `flai forecast`'s 26m for the reasons its basis gives; delivery 02:32, flai's 01:58 moved by the 34 minutes added.

Cost of delay: 40 USD a week, as `flai cod S-0604` works it out from the story's one input, `penalty_per_week: 40`, which the orchestrator set on TH-0896 from TH-0888's table. No revenue rides on the story, and the value stands as computed.

### Accepted by the orchestrator

- Verified: cf467cbc3912661f277a2d228084955c77ee09ab
- At: 2026-11-09T08:07:31Z

Verdict: pass at cf467cbc3912661f277a2d228084955c77ee09ab; the close-out passed every step, and each criterion is met by the diff.
- 1: flai/internal/release/history.go, flai/internal/release/release.go, flai/internal/release/history_test.go
- 2: flai/internal/release/history.go, flai/internal/release/release.go, flai/internal/serve/mergequeue.go, flai/internal/release/history_test.go, flai/internal/serve/mergequeue_test.go
- 3: flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go
- 4: design/system/server-performance.md, design/system/flai-cli.md
