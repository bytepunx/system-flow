---
id: T-1481
type: task
nature: remediation
title: Close I-0124 as fixed by the rebase check whose findings block acceptance until the operator restores the path or keeps the removal
status: done
parent: S-0371
owner: alex
created: 2026-10-12T17:16:13Z
updated: 2026-10-14T22:42:09Z
transitions:
  - to: ready
    at: 2026-10-14T22:42:04Z
    by: agent-S-0371
  - to: in-progress
    at: 2026-10-14T22:42:04Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T22:42:09Z
    by: agent-S-0371
stream: S-0371
tags: [cli, docs]
touches: [design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/summary.md]
after: [T-1480]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1573
      cache_read: 248188
      cache_write: 11633
      cost: 0.145
---
# T-1481 Close I-0124 as fixed by the rebase check whose findings block acceptance until the operator restores the path or keeps the removal

## Work

It waits for T-1480, the last of the work, so that the reason names what was built and documented.

Run `flai issue close I-0124 --reason "..."` in the story's worktree, the reason naming the check `storygit.Rebase` runs after a clean rebase and records, the warning of `flai stream sync` and `flai task done`, the acceptance blocker, `flai stream restore`, `--keep-removal` with its reason, the host writes, the ADR from T-1476, and the tests that reproduce I-0124 from T-1477, T-1574, T-1478, and T-1479. `flai issue close` regenerates `design/issues/summary.md`.

## Done when

- I-0124 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- Criterion 2 of S-0371 is met.

## Notes

Run it in the worktree, not the main checkout: `flai issue` writes under the checkout it runs in.
