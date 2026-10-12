---
id: T-1454
type: task
nature: improvement
title: Close I-0135 saying what fixed it
status: done
parent: S-0366
owner: alex
created: 2026-10-12T14:59:53Z
updated: 2026-10-12T15:37:54Z
transitions:
  - to: ready
    at: 2026-10-12T15:37:51Z
    by: agent-S-0366
  - to: in-progress
    at: 2026-10-12T15:37:52Z
    by: agent-S-0366
  - to: done
    at: 2026-10-12T15:37:54Z
    by: agent-S-0366
stream: S-0366
tags: [flai, issues]
touches: [design/issues/I-0135-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1452, T-1453]
usage:
  source: log
  seconds: 2
  estimated: true
  models: []
---
# T-1454 Close I-0135 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0135 --reason "<what fixed it>"`. The reason names the ADR from T-1451, the change to `ScopeToStory` from T-1452, and the tests that reproduce the three instances on TH-0389. It also says that `flai check` in the main checkout still warns on an open thread on an archived item. Check that `design/issues/summary.md` no longer lists I-0135.

It waits for T-1452 and T-1453, so the reason describes what the story shipped.

## Done when

- I-0135 is closed, with a reason naming the ADR and the fix.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes in the worktree.

## Notes
