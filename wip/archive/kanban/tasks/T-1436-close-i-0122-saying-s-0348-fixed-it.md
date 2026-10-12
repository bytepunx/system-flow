---
id: T-1436
type: task
nature: remediation
title: Close I-0122 saying S-0348 fixed it
status: done
parent: S-0361
owner: alex
created: 2026-10-12T14:14:17Z
updated: 2026-10-12T14:21:15Z
transitions:
  - to: ready
    at: 2026-10-12T14:21:09Z
    by: agent-S-0361
  - to: in-progress
    at: 2026-10-12T14:21:09Z
    by: agent-S-0361
  - to: done
    at: 2026-10-12T14:21:15Z
    by: agent-S-0361
stream: S-0361
tags: [issues]
touches: [design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1435]
usage:
  source: log
  seconds: 6
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 2
      output: 322
      cache_read: 119117
      cache_write: 7158
      cost: 0.0875
---
# T-1436 Close I-0122 saying S-0348 fixed it

## Work

Close I-0122 from the story's worktree, as TH-0388 recommends:

```bash
flai issue close I-0122 --reason "Duplicate of I-0123; fixed by S-0348 (ADR-0133): a check scoped to a story leaves out board.wip-limit; reproduced by TestScopedCheckLeavesOutABoardWIPLimitOutsideTheStory"
```

The command rewrites the issue's status and `design/issues/summary.md`. Waits for T-1435, because the reason names the test that task confirms reproduces I-0122's instance.

## Done when

- I-0122's front matter reads `status: closed`, and its reason names S-0348 and ADR-0133.
- `design/issues/summary.md` no longer lists I-0122.
- `flai check --strict` scoped to S-0361 is clean.

## Notes
