---
id: T-3256
type: task
nature: feature
title: The touches, activity, task plan, inbox, and palette helpers read a state's kind
status: done
parent: S-0758
owner: alex
created: 2026-11-12T16:10:00Z
updated: 2026-11-17T08:12:25Z
transitions:
  - to: ready
    at: 2026-11-17T08:01:28Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:01:28Z
    by: system-flow
  - to: done
    at: 2026-11-17T08:12:25Z
    by: agent-S-0758
stream: S-0758
tags: [dashboard]
touches: [flaiover/src/lib/touches.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/taskplan.ts, flaiover/src/lib/viz/palette.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts]
after: [T-3252]
usage:
  source: log
  seconds: 657
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 13844
      cache_read: 3065727
      cache_write: 66917
      cost: 1.3078
---
# T-3256 The touches, activity, task plan, inbox, and palette helpers read a state's kind

## Work

Waits for T-3252: these helpers read the module that task adds. Replace each comparison with a state's name in these files with the module's kind or colour: open from closed, active work, the wait for review, and accepted.

## Done when

None of these files compares a status with a state's name, and their tests pass.

## Notes
