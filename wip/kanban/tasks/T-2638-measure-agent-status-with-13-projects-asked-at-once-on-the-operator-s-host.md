---
id: T-2638
type: task
nature: improvement
title: Measure agent.status with 13 projects asked at once on the operator's host, before and after, and record it in server-performance.md
status: done
parent: S-0603
owner: alex
created: 2026-11-07T23:48:41Z
updated: 2026-11-08T02:12:34Z
transitions:
  - to: ready
    at: 2026-11-08T02:06:31Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T02:06:35Z
    by: system-flow
  - to: done
    at: 2026-11-08T02:12:34Z
    by: agent-S-0603
stream: S-0603
tags: [cli]
touches: [design/system/server-performance.md]
after: [T-2621, T-2623, T-2629]
usage:
  source: log
  seconds: 359
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 11043
      cache_read: 6083344
      cache_write: 210654
      cost: 2.6148
---
# T-2638 Measure agent.status with 13 projects asked at once on the operator's host, before and after, and record it in server-performance.md

## Work

Criterion 3, its figures. Measure `agent.status` on this repository with the 13 served projects asked at once, on the operator's host, at the story's base and at its head: each call's wall time, its phases from T-2629, and its CPU where it can be taken. Use the method of `design/system/server-performance.md` § A second pass at 596 stories (S-0596) › How it was measured: an in-process program under `/tmp`, writing nothing. Note the host's load average beside the figures.

Record the before and after in `design/system/server-performance.md`, in a section under its causes for cause 12 (S-0603), as the sections for S-0156, S-0500, and S-0531 do: the cause, the fix, and what was measured. Bump its `updated`.

Waits for T-2621, T-2623, and T-2629: it measures what they change, and its phases come from T-2629.

## Done when

- Each of the 13 calls asked at once answers under 200 ms at the story's head, or the figures that missed are recorded with the reason.
- `design/system/server-performance.md` records the before and after, with the phases and the host's load.
- `flai test` on the changed paths passes.

## Notes
