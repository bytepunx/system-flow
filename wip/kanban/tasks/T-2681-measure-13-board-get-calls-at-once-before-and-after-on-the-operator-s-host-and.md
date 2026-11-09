---
id: T-2681
type: task
nature: improvement
title: Measure 13 board.get calls at once before and after on the operator's host, and record them and the per-root history in the design
status: done
parent: S-0604
owner: alex
created: 2026-11-07T23:50:38Z
updated: 2026-11-08T02:37:50Z
transitions:
  - to: ready
    at: 2026-11-08T02:32:30Z
    by: agent-S-0604
  - to: in-progress
    at: 2026-11-08T02:32:30Z
    by: agent-S-0604
  - to: done
    at: 2026-11-08T02:37:50Z
    by: agent-S-0604
stream: S-0604
tags: [cli]
touches: [design/system/server-performance.md, design/system/flai-cli.md]
after: [T-2677, T-2666, T-2670]
usage:
  source: log
  seconds: 320
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 61
      output: 10871
      cache_read: 10151390
      cache_write: 226492
      cost: 3.6449
---
# T-2681 Measure 13 board.get calls at once before and after on the operator's host, and record them and the per-root history in the design

## Work

Build two flai binaries in the same hour, one from `main` and one from the story branch. Run each as earlier rounds in `design/system/server-performance.md` did. Time 13 `board.get` calls at once, one per served project, through the host API's method table in one process with the main-refs watcher running. Give the phases `release.pending`, `merge.queue`, `undone.read`, and `board.view` too. Record the method, the host's load, the project count, and the figures before and after under cause 14 of S-0596's section.

Then rewrite the paragraph in `design/system/flai-cli.md` on `release.Pending`'s kept history. It must say the history is kept per root under that root's own lock, and that `show-ref` is skipped while flai serve's watcher saw no change. It must also say which callers still read git.

Waits for T-2677, T-2666, and T-2670: it measures what they deliver together.

## Done when

- [ ] `design/system/server-performance.md` records 13 `board.get` calls at once, before and after, with how they were measured. Each answers under 300 ms, or the record states the shortfall and its cause.
- [ ] `design/system/flai-cli.md` describes the per-root lock and the watcher's skip.

## Notes
