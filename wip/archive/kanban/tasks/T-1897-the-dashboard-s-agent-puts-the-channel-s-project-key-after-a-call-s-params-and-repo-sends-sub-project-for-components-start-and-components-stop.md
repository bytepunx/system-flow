---
id: T-1897
type: task
nature: feature
title: The dashboard's agent puts the channel's project key after a call's params, and Repo sends sub_project for components.start and components.stop
status: done
parent: S-0448
owner: alex
created: 2026-10-18T02:28:33Z
updated: 2026-10-18T02:35:57Z
transitions:
  - to: ready
    at: 2026-10-18T02:28:58Z
    by: agent-S-0448
  - to: in-progress
    at: 2026-10-18T02:28:59Z
    by: agent-S-0448
  - to: done
    at: 2026-10-18T02:35:57Z
    by: agent-S-0448
stream: S-0448
tags: []
touches: [flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/repo.test.ts, flaiover/src/routes/api/components]
usage:
  source: log
  seconds: 418
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 9369
      cache_read: 1502330
      cache_write: 52479
      cost: 0.8383
---

# T-1897 The dashboard's agent puts the channel's project key after a call's params, and Repo sends sub_project for components.start and components.stop

## Work

- In `flaiover/src/lib/server/agent.ts`'s `ask`, put the channel's `project` key after the call's params, so that no method's params can overwrite it.
- Add a test in `agent.test.ts` that a call whose params carry a `project` of their own still sends the connection's key.
- In `repo.ts`, `componentsStart` and `componentsStop` send the sub-project as `sub_project`; their signatures, and the routes under `flaiover/src/routes/api/components` with their `project` body field, stay as they are. Add or update a test that the components writes send `sub_project` (in `repo.test.ts`, or the routes' `components.test.ts` where it already checks what reaches flai).
- Waits for nothing: it shares no path with T-1896.

## Done when

- The agent test fails on `main`'s `agent.ts` and passes on the change.
- The `sub_project` test fails on `main`'s `repo.ts` and passes.
- `flai test` on the changed paths passes.

## Notes
