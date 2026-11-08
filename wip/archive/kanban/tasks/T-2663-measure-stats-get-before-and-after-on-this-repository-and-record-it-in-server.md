---
id: T-2663
type: task
nature: improvement
title: Measure stats.get before and after on this repository and record it in server-performance.md
status: done
parent: S-0607
owner: alex
created: 2026-11-07T23:49:47Z
updated: 2026-11-08T02:43:53Z
transitions:
  - to: ready
    at: 2026-11-08T02:35:13Z
    by: agent-S-0607
  - to: in-progress
    at: 2026-11-08T02:35:14Z
    by: agent-S-0607
  - to: done
    at: 2026-11-08T02:43:53Z
    by: agent-S-0607
stream: S-0607
tags: [cli, docs]
touches: [design/system/server-performance.md]
after: [T-2646, T-2650, T-2656]
usage:
  source: log
  seconds: 519
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 11182
      cache_read: 5280396
      cache_write: 325272
      cost: 3.0644
---
# T-2663 Measure stats.get before and after on this repository and record it in server-performance.md

## Work

Criterion 3: a warm 30-day `stats.get` on this repository takes at least 250 ms less than before, and `design/system/server-performance.md` records the before and after per window.

- Build the story's base commit and its head, each with `go build` into a folder under `/tmp`. Measure both on the main checkout, writing nothing, as S-0596 did: an in-process program that calls `statsread.Read` and `metrics.Compute`, or `flai hostapi --timing` where the guard allows it.
- Take the windows 1, 7, 30, and 365 days, warm, after one call that is left out, with at least 5 calls each. Give the median, and give the reads and the compute apart. Note the host's load average, since S-0596 found wall times inflated 2 to 3 times under load.
- Show criterion 2: two `stats.get` calls started together compute once.
- Add a section to `server-performance.md` after "A second pass at 596 stories (S-0596)", in the form of the sections on S-0471, S-0500, and S-0531: cause, fix, and what was measured, with a table of before and after per window. Bump the document's `updated`.
- If the 30-day saving is under 250 ms, say so in the narrative and on the story, with the figures, rather than tick criterion 3.

It waits for T-2646, T-2650, and T-2656, since it measures what they changed.

## Done when

- `server-performance.md` holds the before and after per window, and the 30-day saving is at least 250 ms.
- The markdown lint passes on the document.

## Notes
