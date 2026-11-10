---
id: T-2837
type: task
nature: remediation
title: Close I-0325 with what fixed it
status: done
parent: S-0650
owner: alex
created: 2026-11-10T06:58:54Z
updated: 2026-11-10T07:25:09Z
transitions:
  - to: ready
    at: 2026-11-10T07:25:02Z
    by: agent-S-0650
  - to: in-progress
    at: 2026-11-10T07:25:02Z
    by: agent-S-0650
  - to: done
    at: 2026-11-10T07:25:09Z
    by: agent-S-0650
stream: S-0650
tags: [issues]
touches: [design/issues/I-0325-testcomponentsserveendswhenviteexits-fails-under-load-the-fake-vite-exits.md, design/issues/summary.md]
after: [T-2836]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 777
      cache_read: 500072
      cache_write: 17905
      cost: 0.2588
---
# T-2837 Close I-0325 with what fixed it

## Work

Second layer: it waits for T-2836, because the reason names the fix T-2836 makes and the runs that show it holds.

Run `flai issue close I-0325 --reason "<reason>"` in the story's worktree. The reason says the crash-mode fake vite exited on a 300 ms timer that raced flai's readiness poll, that it now exits only after it has answered, and names the reproducing run and the passing runs from T-2836's Notes.

## Done when

- I-0325's status is closed, with the reason.
- `design/issues/summary.md` no longer lists I-0325.
- `flai check --strict` is clean.

## Notes
