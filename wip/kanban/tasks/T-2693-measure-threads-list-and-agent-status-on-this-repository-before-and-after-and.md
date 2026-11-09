---
id: T-2693
type: task
nature: improvement
title: Measure threads.list and agent.status on this repository before and after, and record cause 21's figures in server-performance.md
status: done
parent: S-0611
owner: alex
created: 2026-11-07T23:57:19Z
updated: 2026-11-09T17:26:21Z
transitions:
  - to: ready
    at: 2026-11-09T17:00:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:24:51Z
    by: system-flow
  - to: done
    at: 2026-11-09T17:26:21Z
    by: agent-S-0611
stream: S-0611
tags: [cli, dashboard]
touches: [design/system/server-performance.md]
after: [T-2689, T-2690, T-2691, T-2692]
usage:
  source: log
  seconds: 86
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 87
      cache_read: 3153245
      cache_write: 10899
      cost: 1.1549
---
# T-2693 Measure threads.list and agent.status on this repository before and after, and record cause 21's figures in server-performance.md

## Work

Criterion 4. Measure on this repository, with the main checkout's `bin/flai` for before and the story's build for after:

- `threads.list` as the board asks it, `{}` after the change and `{"all":true}` before, with `flai hostapi --timing`: its `bytes` and `duration_ms`.
- `agent.status` as `/api/host-agent` asks it, with no parameter, before and after; and with `finished: true` after, for the pages that list finished runs. `agent.status` needs flai serve's host state, so time it through flai serve's `request answered` events if `flai hostapi` cannot give it.
- Through flai serve, the `send_ms` of each of the two answers before the change cannot be had; give the after only.

Record the figures in `design/system/server-performance.md`: a section for cause 21 after the story, as the sections for causes 1 to 7 are written, and the sentence near its top that says timing ends when the answer is ready to write, which now has `send_ms` beside it.

Waits for T-2689, T-2690, T-2691, and T-2692: it measures what they change.

## Done when

- [ ] The board's `threads.list` is under 100 KB and `agent.status` under 50 KB on this repository, or the shortfall and its cause are recorded.
- [ ] `design/system/server-performance.md` records cause 21's before and after, with the host's load at the time.
- [ ] `flai test` passes on the changed paths.

## Notes

- Planned by planner-S-0611.
