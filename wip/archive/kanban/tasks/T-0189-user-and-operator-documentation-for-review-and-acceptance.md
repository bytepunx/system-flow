---
id: T-0189
type: task
nature: feature
title: User and operator documentation for review and acceptance
status: done
parent: S-0041
owner: alex
created: 2026-09-19T04:08:29Z
updated: 2026-09-19T04:16:22Z
transitions:
  - to: ready
    at: 2026-09-19T04:16:11Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:16:11Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:16:22Z
    by: system-flow
stream: S-0041
tags: []
touches: [docs]
---

# T-0189 User and operator documentation for review and acceptance

## Work
Update `docs/users/flaiover.md` (reviewing a story: what the page shows, accepting, sending back, what a failure looks like, that nothing is pushed from the container) and `docs/operators/index.md` (acceptance from the dashboard runs as the manifest's owner; what it needs: identity, the host-path mount, the excludes).

## Done when
- Both documents state the behaviour as built
- Markdown lint and `flai check --strict` pass

## Notes
