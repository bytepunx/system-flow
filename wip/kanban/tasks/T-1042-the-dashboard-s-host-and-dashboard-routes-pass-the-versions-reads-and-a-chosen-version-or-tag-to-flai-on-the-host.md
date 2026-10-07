---
id: T-1042
type: task
nature: improvement
title: The dashboard's host and dashboard routes pass the versions reads and a chosen version or tag to flai on the host
status: in-progress
parent: S-0298
owner: alex
created: 2026-10-06T21:45:05Z
updated: 2026-10-07T14:33:12Z
transitions:
  - to: ready
    at: 2026-10-07T14:33:12Z
    by: agent-S-0298
  - to: in-progress
    at: 2026-10-07T14:33:12Z
    by: agent-S-0298
stream: S-0298
tags: [dashboard]
touches: [flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/routes/api/host/+server.ts, flaiover/src/routes/api/host/host.test.ts, flaiover/src/routes/api/dashboard/+server.ts, flaiover/src/routes/api/dashboard/dashboard.test.ts]
after: [T-1038]
---
# T-1042 The dashboard's host and dashboard routes pass the versions reads and a chosen version or tag to flai on the host

## Work

Add the host API methods the ADR names, `host.versions` and `dashboard.versions`, to the channel's method list in `flaiover/src/lib/server/agent.ts`. In `routes/api/host/+server.ts` add the action `versions` (a read, with the timeout `check` has) and let `upgrade` carry an optional `version`; in `routes/api/dashboard/+server.ts` add `versions` and let `upgrade` carry an optional `tag`. Pass each through as the method's parameter, unvalidated beyond its shape: flai on the host checks it is published. Keep the gating by the `host` and `dashboard` host actions and the no-retry handling of an upgrade whose connection drops. The routes' tests stand the channel in, so this waits only for T-1038 for the method and parameter names, not for the Go side.

## Done when

- `versions` answers both lists through the stood-in channel, and `upgrade` with a version or tag sends it, refused while its host action is off, tested in `host.test.ts` and `dashboard.test.ts`.
- `npm run test` and `npm run check` pass in `flaiover`.

## Notes
