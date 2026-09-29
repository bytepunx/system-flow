---
id: T-0521
type: task
nature: feature
title: "Close out: criteria, design pointer, narrative, check"
status: done
parent: S-0145
owner: alex
created: 2026-09-29T02:12:50Z
updated: 2026-09-29T02:22:12Z
transitions:
  - to: ready
    at: 2026-09-29T02:21:21Z
    by: agent-S-0145
  - to: in-progress
    at: 2026-09-29T02:21:21Z
    by: agent-S-0145
  - to: done
    at: 2026-09-29T02:22:12Z
    by: agent-S-0145
stream: S-0145
tags: []
---

# T-0521 Close out: criteria, design pointer, narrative, check

## Work

Tick the criteria that were verified, add a pointer from `design/system/agent-context.md` to the draft ADR, close the narrative, run `flai check --strict` and the markdown lint on both checkouts, commit, and move the story to review.

## Done when

Every criterion is checked, the narrative's `## Current state` and `## Next steps` are true, `flai check --strict` and `make lint-md` pass, the worktree is clean, and S-0145 is in review.

## Notes
