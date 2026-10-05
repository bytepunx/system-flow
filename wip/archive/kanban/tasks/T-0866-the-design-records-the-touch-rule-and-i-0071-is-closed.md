---
id: T-0866
type: task
nature: remediation
title: The design records the touch rule and I-0071 is closed
status: done
parent: S-0260
owner: alex
created: 2026-10-05T03:14:17Z
updated: 2026-10-05T03:22:29Z
transitions:
  - to: ready
    at: 2026-10-05T03:14:33Z
    by: agent-S-0260
  - to: in-progress
    at: 2026-10-05T03:21:04Z
    by: agent-S-0260
  - to: done
    at: 2026-10-05T03:22:29Z
    by: agent-S-0260
stream: S-0260
tags: []
touches: [design/system/work-hierarchy.md, design/issues, docs/users/flai.md]
after: [T-0865]
usage:
  source: log
  seconds: 85
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 57
      output: 14243
      cache_read: 2512284
      cache_write: 76652
      cost: 1.2967
---
# T-0866 The design records the touch rule and I-0071 is closed

## Work

Say in `design/system/work-hierarchy.md`'s touches paragraph what a touch may be, as T-0865 built it, and close I-0071 with `flai issue close I-0071 --reason` naming the fix, from the worktree. Waits for T-0865, whose rule it records.

## Done when

- [x] work-hierarchy.md states the touch rule and that every writer applies it
- [x] I-0071 is closed with a reason naming S-0260's fix

## Notes
