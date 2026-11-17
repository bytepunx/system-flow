---
id: T-3527
type: task
nature: improvement
title: Measure board.get's undone.read with the batched diffs, record it in server-performance.md, and close I-0349
status: done
parent: S-0828
owner: alex
created: 2026-11-17T08:29:03Z
updated: 2026-11-17T12:13:41Z
transitions:
  - to: ready
    at: 2026-11-17T12:05:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T12:05:53Z
    by: system-flow
  - to: done
    at: 2026-11-17T12:13:41Z
    by: agent-S-0828
stream: S-0828
tags: [flai, docs, performance]
touches: [flai/internal/hostapi/hostapi.go, design/system/server-performance.md, design/issues/I-0349-board-get-s-undone-read-runs-one-git-diff-per-recorded-path-of-each-story-s.md, design/issues/summary.md]
after: [T-3525]
usage:
  source: log
  seconds: 466
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 32
      output: 293
      cache_read: 4077767
      cache_write: 24411
      cost: 1.5673
---
# T-3527 Measure board.get's undone.read with the batched diffs, record it in server-performance.md, and close I-0349

## Work

Waits for T-3525, whose batched `RecordedUndone` is what this task measures and documents.

- In `flai/internal/hostapi/hostapi.go`, make the comments of `undoneWorkers` and `boardUndone` say what a story's read runs now: one `git diff` for its branch and one per distinct old base, not a git process per recorded path.
- Measure `board.get`'s `undone.read` before and after, as § The merge queue's git read once (S-0698) of `design/system/server-performance.md` did: the git processes a call runs, from the timed runner's `exec.git.*` phases, and the phase's time, on a story with an undone record of many paths. Where no story on the host has such a record, build one in a scratch repository and say so. Never write a figure you did not observe.
- Add a section to `design/system/server-performance.md` after § The merge queue's git read once (S-0698), with its cause, its fix, and what was measured, and say there whether keeping the answer while the branch's and main's heads are unchanged, as S-0678 keeps the trial merges, is still worth a story.
- Close the issue with `flai issue close I-0349 --reason "<what fixed it>"`, naming S-0828 and the batched diffs.

## Done when

- The comments in `hostapi.go` match what `RecordedUndone` runs.
- `design/system/server-performance.md` records the before and after git processes and time of `undone.read`, and passes the markdown lint.
- I-0349 is closed with its reason, and `design/issues/summary.md` no longer lists it.

## Notes
