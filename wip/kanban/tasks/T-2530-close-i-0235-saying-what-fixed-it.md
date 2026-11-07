---
id: T-2530
type: task
nature: improvement
title: Close I-0235 saying what fixed it
status: done
parent: S-0555
owner: alex
created: 2026-11-07T03:47:56Z
updated: 2026-11-07T04:15:05Z
transitions:
  - to: ready
    at: 2026-11-07T04:14:45Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T04:14:46Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:15:05Z
    by: agent-S-0555
stream: S-0555
tags: [flai, issues]
touches: [design/issues/I-0235-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-2521, T-2527]
usage:
  source: log
  seconds: 14
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 48
      cache_read: 848392
      cache_write: 2317
      cost: 0.3041
---
# T-2530 Close I-0235 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0235 --reason "<what fixed it>"`. The reason names T-2518's ADR, T-2521's change to `ScopeToStories`, and the tests that reproduce the instances on S-0549. It also says that `flai check` in the main checkout still reports a story in progress with no narrative. Then check that `design/issues/summary.md` no longer lists I-0235.

It waits for T-2521 and T-2527, so that the reason describes what the story shipped.

## Done when

- I-0235 is closed, with a reason that names the ADR and the fix.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes in the worktree.

## Notes
