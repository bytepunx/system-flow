---
id: T-0307
type: task
nature: feature
title: flaiover accepts more than one project connection at once, keyed by project, sharing the one credential
status: done
parent: S-0080
owner: alex
created: 2026-09-20T21:09:24Z
updated: 2026-09-20T21:30:42Z
transitions:
  - to: ready
    at: 2026-09-20T21:17:33Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T21:17:33Z
    by: system-flow
  - to: done
    at: 2026-09-20T21:30:42Z
    by: system-flow
stream: S-0080
tags: []
---
# T-0307 flaiover accepts more than one project connection at once, keyed by project, sharing the one credential

## Work
AgentHub becomes a registry: one shared credential, many simultaneous /agent connections, one per project key, a newer connection replacing only the one with the same key. repo() and agent() stay as they are called today (no change to existing routes) but resolve to the current request's project through AsyncLocalStorage, set in hooks.server.ts from the URL or a project query parameter. A read of the connected projects (key, name, and whether attended) is available server side, and a request naming a project key nothing is connected for is refused by the server, not only by client-side routing.

## Done when
- Server tests: two flai connections at once, each answering for its own project; one lost connection does not affect the other; an unknown key is refused
- The flaiover tests pass

## Notes
