---
id: T-2807
type: task
nature: remediation
title: Close I-0319 with what fixed it
status: done
parent: S-0640
owner: alex
created: 2026-11-09T22:08:00Z
updated: 2026-11-10T05:19:56Z
transitions:
  - to: ready
    at: 2026-11-10T05:19:45Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:19:46Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:19:56Z
    by: agent-S-0640
stream: S-0640
tags: [flaiover, testing]
touches: [design/issues/I-0319-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md, design/issues/summary.md]
after: [T-2805]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 441
      cache_read: 442810
      cache_write: 6478
      cost: 0.1492
---
# T-2807 Close I-0319 with what fixed it

## Work

Run `flai issue close I-0319 --reason "<text>"`. The reason names the change T-2805 made to the props-panel test, and the test's duration before and after it from T-2805's `## Notes`. It writes the issue file and regenerates `design/issues/summary.md`.

It waits for T-2805, because the reason states what T-2805 changed and measured.

## Done when

- I-0319's status is closed, with a reason that names the fix and both durations.
- `design/issues/summary.md` no longer lists I-0319.
- `flai check --strict` passes.

## Notes
