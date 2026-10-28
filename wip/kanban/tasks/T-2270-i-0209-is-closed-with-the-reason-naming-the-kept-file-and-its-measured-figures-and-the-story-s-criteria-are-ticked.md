---
id: T-2270
type: task
nature: improvement
title: I-0209 is closed with the reason naming the kept file and its measured figures, and the story's criteria are ticked
status: done
parent: S-0531
owner: alex
created: 2026-10-28T18:24:56Z
updated: 2026-10-28T18:36:55Z
transitions:
  - to: ready
    at: 2026-10-28T18:26:36Z
    by: agent-S-0531
  - to: in-progress
    at: 2026-10-28T18:36:40Z
    by: agent-S-0531
  - to: done
    at: 2026-10-28T18:36:55Z
    by: agent-S-0531
stream: S-0531
tags: [flai]
touches: [design/issues/I-0209-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md, design/issues/summary.md]
after: [T-2268, T-2269]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 59
      cache_read: 397885
      cache_write: 4171
      cost: 0.1361
---
# T-2270 I-0209 is closed with the reason naming the kept file and its measured figures, and the story's criteria are ticked

## Work

In the story's worktree, close I-0209 with `flai issue close I-0209 --reason "<reason>"`. The reason names the kept file under `.flai-cache`, its key, and the cold figures T-2268 recorded in `design/system/server-performance.md`, and says that the first call after a branch tip moved still reads the whole log. `flai issue close` rewrites `design/issues/summary.md`.

Tick criterion 1 once T-2267's tests and T-2268's figures show the cause gone, and criterion 2 once the issue is closed, with `flai criteria tick S-0531 <n>`.

This task waits for T-2268, whose figures the reason cites, and for T-2269, so that the documents the reason points to are written.

## Done when

- I-0209's status is closed, with the reason above.
- `design/issues/summary.md` no longer lists I-0209.
- Both criteria of S-0531 are ticked.

## Notes
