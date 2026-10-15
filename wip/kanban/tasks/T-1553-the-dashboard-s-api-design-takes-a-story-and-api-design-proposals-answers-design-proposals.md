---
id: T-1553
type: task
nature: feature
title: The dashboard's /api/design takes a story, and /api/design/proposals answers design.proposals
status: ready
parent: S-0384
owner: alex
created: 2026-10-13T17:48:21Z
updated: 2026-10-15T03:17:15Z
transitions:
  - to: ready
    at: 2026-10-15T03:17:15Z
    by: system-flow
stream: S-0384
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts, flaiover/src/routes/api/design/proposals/+server.ts, flaiover/src/routes/api/design/proposals/proposals.test.ts]
after: [T-1552]
---
# T-1553 The dashboard's /api/design takes a story, and /api/design/proposals answers design.proposals

## Work

Add `design.proposals` to `REQUIRED_METHODS`; in `repo.ts` pass a `story` to the `design.get` read, keyed apart per story, and add the proposals read, forgotten on any story or document change. `GET /api/design?story=S-nnnn` and `GET /api/design/proposals` answer them. Tests cover each, and a refused story. It waits for T-1552, whose methods it calls.

## Done when

- Both endpoints answer as the host does, and their tests pass.

## Notes

- Layer 3.
