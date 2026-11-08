---
id: T-2654
type: task
nature: improvement
title: Measure settings.get before and after on the operator's host, and record cause 15's fix and figures in server-performance.md
status: done
parent: S-0605
owner: alex
created: 2026-11-07T23:49:24Z
updated: 2026-11-08T00:51:34Z
transitions:
  - to: ready
    at: 2026-11-08T00:49:39Z
    by: agent-S-0605
  - to: in-progress
    at: 2026-11-08T00:49:42Z
    by: agent-S-0605
  - to: done
    at: 2026-11-08T00:51:34Z
    by: agent-S-0605
stream: S-0605
tags: [cli]
touches: [design/system/server-performance.md]
after: [T-2639, T-2645]
usage:
  source: log
  seconds: 112
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5299
      cache_read: 2846375
      cache_write: 68667
      cost: 1.1236
---
# T-2654 Measure settings.get before and after on the operator's host, and record cause 15's fix and figures in server-performance.md

## Work

- Measure `settings.get` for this project on the operator's host, against the import folders and the 13 served projects there, with the main checkout's `bin/flai` before the change and the story's build after: `flai hostapi --timing settings.get`, several runs each, with the host's load average beside them. Take the serve log's p50 of 825 ms from cause 15 as the before figure as well.
- Record in `design/system/server-performance.md`, under S-0596's section, what changed for cause 15 and the before and after, with the phases T-2645 records, as the sections for S-0156 to S-0531 record theirs.
- If the story's agent cannot reach the import folders or the running flai serve's state from where it runs, say so in the story's Notes and leave criterion 3 unticked, as `work-management.md` says.
- Waits for T-2639 and T-2645: it measures the code both write, by the phases T-2645 adds.

## Done when

- `settings.get` answers under 100 ms on the operator's host, or the story's Notes say why that could not be measured here.
- `design/system/server-performance.md` records the before and after for cause 15, and the markdown lint passes on it.

## Notes
