---
id: T-1865
type: task
nature: feature
title: The dashboard's /api/components routes call components.list, status, start, and stop on the host, answering 501 from an older flai
status: done
parent: S-0439
owner: alex
created: 2026-10-16T21:34:43Z
updated: 2026-10-17T01:18:49Z
transitions:
  - to: ready
    at: 2026-10-17T01:11:14Z
    by: agent-S-0439
  - to: in-progress
    at: 2026-10-17T01:11:15Z
    by: agent-S-0439
  - to: done
    at: 2026-10-17T01:18:49Z
    by: agent-S-0439
stream: S-0439
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/components/+server.ts, flaiover/src/routes/api/components/components.test.ts, flaiover/src/routes/api/components/start/+server.ts, flaiover/src/routes/api/components/stop/+server.ts]
usage:
  source: log
  seconds: 454
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 19156
      cache_read: 2701347
      cache_write: 72265
      cost: 1.345
---
# T-1865 The dashboard's /api/components routes call components.list, status, start, and stop on the host, answering 501 from an older flai

## Work

- `repo.ts`: `components()`, `componentsStatus()`, `componentsStart(project?)`, and `componentsStop(project?)` over `Repo.ask`; the catalogue is not cached, since a file may be added at any time.
- `agent.ts`: the four methods in `REQUIRED_METHODS`.
- `GET /api/components` answers the catalogue and the status together. `POST /api/components/start` and `/stop` answer the status. A host flai without the methods answers 501 with the upgrade hint, as `/api/design/tokens` does, and the action being off answers 403 with the command to enable it.

Waits for nothing in this story; it needs S-0438's methods on the host.

## Done when

Route tests with a fake hub cover each route, the 501, and the 403. `flai test flaiover/src/routes/api/components flaiover/src/lib/server` passes, and `contract_test.go` passes against S-0438's methods.

## Notes

None.
