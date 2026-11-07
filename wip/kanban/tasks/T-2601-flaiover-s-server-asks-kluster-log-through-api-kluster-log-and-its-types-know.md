---
id: T-2601
type: task
nature: feature
title: flaiover's server asks kluster.log through /api/kluster/log, and its types know pending steps and the components setting
status: in-progress
parent: S-0595
owner: alex
created: 2026-11-07T22:07:13Z
updated: 2026-11-07T22:24:49Z
transitions:
  - to: ready
    at: 2026-11-07T22:10:02Z
    by: agent-S-0595
  - to: in-progress
    at: 2026-11-07T22:24:49Z
    by: agent-S-0595
stream: S-0595
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/log/+server.ts, flaiover/src/routes/api/kluster/log/log.test.ts]
after: [T-2599]
---
# T-2601 flaiover's server asks kluster.log through /api/kluster/log, and its types know pending steps and the components setting

## Work

ADR-0263 §4. In `flaiover/src/lib/server/repo.ts` add `klusterLog({namespace, step, offset})`, a read through `ask('kluster.log', …)` wrapped in `upgraded(…)` as the other kluster methods are, with its type `KlusterLog` (the fields T-2599 answers). Widen `KlusterStep`'s outcome to include `pending`, and the setting kind type to include `components`. Add `GET /api/kluster/log?namespace=&step=&offset=` under `flaiover/src/routes/api/kluster/log/+server.ts`, validating the query (namespace a DNS label, step a step name, offset a whole number of zero or more) and answering flai's object; test it in `log.test.ts` with the fake host the other kluster route tests use. Waits for T-2599, whose answer it types.

## Done when

- `log.test.ts` covers a read, a bad query refused with 400, and an old host answered 501; it fails without the route.
- `flai test` on the changed paths passes (lint, vitest, svelte-check).

## Notes
