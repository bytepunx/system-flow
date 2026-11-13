---
id: T-3272
type: task
nature: improvement
title: Close I-0354 saying what fixed it
status: done
parent: S-0761
owner: alex
created: 2026-11-13T07:29:58Z
updated: 2026-11-13T09:06:04Z
transitions:
  - to: ready
    at: 2026-11-13T09:05:52Z
    by: agent-S-0761
  - to: in-progress
    at: 2026-11-13T09:05:52Z
    by: agent-S-0761
  - to: done
    at: 2026-11-13T09:06:04Z
    by: agent-S-0761
stream: S-0761
tags: [flai, issues]
touches: [design/issues/I-0354-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-3268, T-3270, T-3271]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 993
      cache_read: 563182
      cache_write: 14325
      cost: 0.2243
---
# T-3272 Close I-0354 saying what fixed it

## Work

Close I-0354 with `flai issue close I-0354 --reason "<reason>"`. The reason names the two fixes: S-0678's interval closed at its close by T-3268, and T-3269's ADR with what T-3270 built, naming the tests that reproduce it. Say it is the same finding as I-0347, which S-0723 closed, and why it came back: S-0723 could not repair S-0678 from its branch. If the ADR chose repair alone, the reason says so instead, and T-3270 and T-3271 will have been cancelled.

Close it last, after the close-out's checks of the other tasks, so that the head the close-out verifies is the one that closes it. It waits for T-3268 so the main checkout's check is clean, and for T-3270 and T-3271 so the reason names what they built.

## Done when

- I-0354 is closed with a reason naming the repair, the ADR, and the tests.
- `design/issues/summary.md` lists it closed.

## Notes
