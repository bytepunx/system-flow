---
id: T-0312
type: task
nature: feature
title: Tried end to end with two scratch projects sharing one dashboard
status: done
parent: S-0080
owner: alex
created: 2026-09-20T21:09:25Z
updated: 2026-09-20T22:42:07Z
transitions:
  - to: ready
    at: 2026-09-20T22:31:46Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T22:31:46Z
    by: system-flow
  - to: done
    at: 2026-09-20T22:42:07Z
    by: system-flow
stream: S-0080
tags: []
---
# T-0312 Tried end to end with two scratch projects sharing one dashboard

## Work
Two scratch projects, own configuration, one flai serve, one container: both connect, the switcher lists both with their state, moving between them works, a write in one does not touch the other, stopping one project's registration leaves the container serving the other, stopping the last one stops the container.

## Done when
- What was tried and what was not is in the narrative; every unexpected answer explained before review

## Notes
