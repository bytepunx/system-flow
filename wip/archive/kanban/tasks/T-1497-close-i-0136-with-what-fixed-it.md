---
id: T-1497
type: task
nature: improvement
title: Close I-0136 with what fixed it
status: done
parent: S-0374
owner: alex
created: 2026-10-12T17:34:31Z
updated: 2026-10-14T22:05:41Z
transitions:
  - to: ready
    at: 2026-10-14T22:05:35Z
    by: agent-S-0374
  - to: in-progress
    at: 2026-10-14T22:05:35Z
    by: agent-S-0374
  - to: done
    at: 2026-10-14T22:05:41Z
    by: agent-S-0374
stream: S-0374
tags: [cli]
touches: [design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md, design/issues/summary.md]
after: [T-1496]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 1498
      cache_read: 239001
      cache_write: 5450
      cost: 0.1135
---
# T-1497 Close I-0136 with what fixed it

## Work

Waits for T-1496, so the fix is built, tested, and documented before the issue says it is.

Run `flai issue close I-0136 --reason "..."` in the story's worktree, the reason naming the ADR, that `flai verify` now reuses across a clean rebase the tiers main's new paths do not select (T-1494), and that the close-out syncs and verifies again, at most three rounds (T-1495). `flai` regenerates `design/issues/summary.md`.

## Done when

- I-0136 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes

Layer 4. Ticks the story's second criterion.
