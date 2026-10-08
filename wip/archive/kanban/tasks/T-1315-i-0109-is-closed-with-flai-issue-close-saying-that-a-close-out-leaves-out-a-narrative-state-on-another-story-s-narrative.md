---
id: T-1315
type: task
nature: improvement
title: I-0109 is closed with flai issue close, saying that a close-out leaves out a narrative.state on another story's narrative
status: done
parent: S-0323
owner: alex
created: 2026-10-08T00:26:47Z
updated: 2026-10-08T06:17:06Z
transitions:
  - to: ready
    at: 2026-10-08T06:12:28Z
    by: agent-S-0323
  - to: in-progress
    at: 2026-10-08T06:16:22Z
    by: agent-S-0323
  - to: done
    at: 2026-10-08T06:17:06Z
    by: agent-S-0323
stream: S-0323
tags: [flai]
touches: [design/issues/I-0109-flai-check-finds-narrative-state-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0111-flai-check-finds-narrative-state-outside-the-story-at-close-out.md]
after: [T-1314]
usage:
  source: log
  seconds: 44
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2700
      cache_read: 802663
      cache_write: 24946
      cost: 0.4142
---
# T-1315 I-0109 is closed with flai issue close, saying that a close-out leaves out a narrative.state on another story's narrative

## Work

Close I-0109 once the cause is gone, for criterion 2. It waits for T-1314, because the issue closes on the fix, not on the plan.

- In the story's worktree, run `flai issue close I-0109 --reason "<reason>"`, the reason naming S-0323, T-1314, the ADR T-1313 wrote, and what changed: a check scoped to a story leaves out a `narrative.state` on another story's narrative, so a close-out records none.
- Check that `design/issues/summary.md` no longer lists I-0109.

## Done when

- I-0109's status is closed with the reason above, and `summary.md` is regenerated.
- `flai check --strict` passes.

## Notes
