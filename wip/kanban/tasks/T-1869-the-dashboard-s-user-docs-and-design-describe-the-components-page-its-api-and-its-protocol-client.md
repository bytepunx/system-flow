---
id: T-1869
type: task
nature: feature
title: The dashboard's user docs and design describe the Components page, its API, and its protocol client
status: backlog
parent: S-0439
owner: alex
created: 2026-10-16T21:35:03Z
updated: 2026-10-16T21:35:03Z
transitions: []
stream: S-0439
tags: [dashboard]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-1868]
---
# T-1869 The dashboard's user docs and design describe the Components page, its API, and its protocol client

## Work

- `docs/users/flaiover.md`: the Components row in the site menu table, and a section on the page: what it shows, how to enable and start it, and what its errors mean.
- `design/system/flaiover-dashboard.md`: the route, the three API routes in its API table, and the protocol client.

Waits for T-1868, so that the docs describe the page as built.

## Done when

Both files describe the page as it behaves, and `flai test` on them passes.

## Notes

None.
