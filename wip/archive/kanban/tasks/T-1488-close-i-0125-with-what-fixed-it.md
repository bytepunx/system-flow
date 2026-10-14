---
id: T-1488
type: task
nature: improvement
title: Close I-0125 with what fixed it
status: done
parent: S-0372
owner: alex
created: 2026-10-12T17:23:44Z
updated: 2026-10-14T21:40:47Z
transitions:
  - to: ready
    at: 2026-10-14T21:40:30Z
    by: agent-S-0372
  - to: in-progress
    at: 2026-10-14T21:40:30Z
    by: agent-S-0372
  - to: done
    at: 2026-10-14T21:40:47Z
    by: agent-S-0372
stream: S-0372
tags: [issues]
touches: [design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/summary.md]
after: [T-1487]
usage:
  source: log
  seconds: 17
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 4450
      cache_read: 616231
      cache_write: 28546
      cost: 0.3728
---
# T-1488 Close I-0125 with what fixed it

## Work

Close I-0125 with `flai issue close I-0125 --reason`, run in the story's worktree. The reason names S-0372 and the ADR of T-1482, and says what each change remedies:

- the settled reply (T-1484, T-1485) for the first two instances;
- the notice without generated paths (T-1483) for the third;
- the blocked hold on wakes (T-1486) for the fourth.

It also names the tests that reproduce them. It waits for T-1487, the last of the work, so that the reason describes what was built and documented.

## Done when

- I-0125's status is closed, and its Remediation section carries the reason.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes

flai writes both files. Neither is edited by hand.
