---
id: T-2830
type: task
nature: remediation
title: Close I-0323 with the fix as its reason
status: done
parent: S-0648
owner: alex
created: 2026-11-10T06:37:42Z
updated: 2026-11-10T08:39:10Z
transitions:
  - to: ready
    at: 2026-11-10T08:39:00Z
    by: agent-S-0648
  - to: in-progress
    at: 2026-11-10T08:39:01Z
    by: agent-S-0648
  - to: done
    at: 2026-11-10T08:39:10Z
    by: agent-S-0648
stream: S-0648
tags: [issues]
touches: [design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-2827, T-2828, T-2829]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2578
      cache_read: 950523
      cache_write: 24309
      cost: 0.4036
---
# T-2830 Close I-0323 with the fix as its reason

## Work

Third layer: it waits for T-2827 and T-2828, whose fix the reason names, and for T-2829, whose removal clears the finding.

- Run `flai check --strict` in the story's worktree and confirm it reports no `item.duplicate-id`.
- Close the issue with `flai issue close I-0323 --reason "<reason>"`. The reason names the race, a usage `Save` that recreated a task file `Archive` had just moved, and the fix: `Save` refuses a moved item, `Archive` moves under the items lock, `Measure` writes to the item where it now is, and the stray T-2760 copy is removed. It names the tests of T-2827 and T-2828 by name.
- Check that `design/issues/summary.md` no longer lists I-0323.

## Done when

- I-0323's status is closed and its reason names the fix and the tests.
- `design/issues/summary.md` does not list I-0323.
- `flai check --strict` is clean in the story's worktree.

## Notes
