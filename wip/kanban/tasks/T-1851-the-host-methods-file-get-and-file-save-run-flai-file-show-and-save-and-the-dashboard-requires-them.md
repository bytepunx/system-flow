---
id: T-1851
type: task
nature: feature
title: The host methods file.get and file.save run flai file show and save, and the dashboard requires them
status: done
parent: S-0436
owner: alex
created: 2026-10-16T21:31:49Z
updated: 2026-10-16T22:33:27Z
transitions:
  - to: ready
    at: 2026-10-16T22:25:44Z
    by: agent-S-0436
  - to: in-progress
    at: 2026-10-16T22:25:44Z
    by: agent-S-0436
  - to: done
    at: 2026-10-16T22:33:27Z
    by: agent-S-0436
stream: S-0436
tags: [cli, dashboard]
touches: [flai/internal/hostapi/files_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts]
after: [T-1850]
usage:
  source: log
  seconds: 463
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 86
      output: 391
      cache_read: 4489785
      cache_write: 143264
      cost: 2.152
---
# T-1851 The host methods file.get and file.save run flai file show and save, and the dashboard requires them

## Work

- `flai/internal/hostapi/files.go`: `file.get` `{path, story?}`, merged into the method table in `hostapi.go` as `designMethods` is, answering `{path, content, hash}` from T-1849's `FileShow`.
- `writes.go`: a `file.save` spec `{path, content, hash, story?, message?}` that runs `flai file save` with the content on standard input, as `doc.save`'s spec does, mapping exit 3 to Conflict and exit 4 to Refused, ungated as `doc.save` is.
- `flaiover/src/lib/server/agent.ts`: add both to `REQUIRED_METHODS`; `contract_test.go` then checks them.

Waits for T-1850: the save spec runs its command.

## Done when

Host API tests cover a read, a save that commits the token file and the stylesheet, a conflict, a refusal with its findings, a story's worktree, and a path out of reach. `flai test flai/internal/hostapi` and `flai test flaiover/src/lib/server` pass.

## Notes

The dashboard route `flaiover/src/routes/api/design/tokens/+server.ts` already calls both (S-0383); it needs no change.
