---
id: T-0237
type: task
nature: feature
title: Try it in a browser against a scratch project, user documentation, criteria
status: backlog
parent: S-0060
owner: alex
created: 2026-09-19T09:39:14Z
updated: 2026-09-19T09:39:14Z
transitions: []
stream: S-0060
tags: []
---

# T-0237 Try it in a browser against a scratch project, user documentation, criteria

## Work
On the branch's dev server against a scratch git project: create a proposed ADR, see the file, the index row, and the commit; accept it; create one that supersedes and one that refines; a body the check refuses; the editor refusing an accepted ADR's body. `docs/users`, `flaiover-dashboard.md`. Tick the criteria for what was seen. Hold the story out of review until S-0059 is accepted and this branch is synced onto main.

## Done when
- The browser check is recorded in the narrative
- The documents describe it
- The criteria are ticked and `flai check --strict` is clean

## Notes
