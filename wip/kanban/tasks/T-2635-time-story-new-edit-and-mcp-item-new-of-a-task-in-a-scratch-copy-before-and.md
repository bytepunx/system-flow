---
id: T-2635
type: task
nature: improvement
title: Time story new, edit, and MCP item_new of a task in a scratch copy, before and after, and record them in server-performance.md
status: done
parent: S-0598
owner: alex
created: 2026-11-07T23:48:19Z
updated: 2026-11-08T00:56:54Z
transitions:
  - to: ready
    at: 2026-11-08T00:35:35Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T00:35:40Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:56:54Z
    by: agent-S-0598
stream: S-0598
tags: [cli]
touches: [design/system/server-performance.md]
after: [T-2622, T-2625]
usage:
  source: log
  seconds: 1274
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 56
      output: 19315
      cache_read: 9371723
      cache_write: 191706
      cost: 3.4266
---
# T-2635 Time story new, edit, and MCP item_new of a task in a scratch copy, before and after, and record them in server-performance.md

## Work

Make a scratch copy of this repository outside the worktree, as S-0596 did in [server-performance.md § How it was measured](../../../design/system/server-performance.md#how-it-was-measured). In it, time `flai story new --body-stdin`, `flai edit`, and MCP `item_new` of a task with a body, five runs each, wall and CPU. Use a `flai` built from the main branch for the before and one built from the story's branch for the after.

Record both in `design/system/server-performance.md`: cause 8's row of "Where the time goes" and a section for S-0598 with how it was measured and the figures, as the S-0471, S-0500, and S-0531 sections do. If any of the three is 3 s or more, say what the time goes to, and fix it in T-2622 or T-2625 before you close this task.

This task waits for T-2622 and T-2625, the two writes it times.

## Done when

- Each of the three takes under 3 s wall in the scratch copy, measured, not estimated.
- `server-performance.md` records the before and after for each, the host, and the item count.

## Notes
