---
id: T-0218
type: task
nature: feature
title: Dashboard preview and move for a research story, user documentation, and unblocking S-0052
status: done
parent: S-0053
owner: alex
created: 2026-09-19T07:15:38Z
updated: 2026-09-19T07:24:00Z
transitions:
  - to: ready
    at: 2026-09-19T07:20:23Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T07:20:23Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:24:00Z
    by: system-flow
stream: S-0053
tags: []
---

# T-0218 Dashboard preview and move for a research story, user documentation, and unblocking S-0052

## Work
The dashboard's acceptance preview (`/api/items/:id/acceptance`, `AcceptConfirm.svelte`) for a research story: it is acceptable, says research releases nothing and why, and lists components landing without a release; the move then accepts it, with the existing accepted-locally notice and push command when the container cannot push (S-0052 decides anything more). Tests for the preview. `docs/users/flai.md` and, if the confirmation's wording changes, `docs/users/flaiover.md`. Then `flai unblock S-0052`, which the operator blocked for this reason. Tick the criteria for what was verified.

## Done when
- The preview test passes and the dashboard build passes
- The user docs describe accepting a research story
- S-0052 is unblocked and the criteria are ticked

## Notes
