---
id: T-1577
type: task
nature: feature
title: "POST /api/items/[id]/restore runs the host write stream.restore for one path"
status: done
parent: S-0388
owner: alex
created: 2026-10-14T21:08:46Z
updated: 2026-10-14T23:15:30Z
transitions:
  - to: ready
    at: 2026-10-14T23:10:19Z
    by: agent-S-0388
  - to: in-progress
    at: 2026-10-14T23:10:20Z
    by: agent-S-0388
  - to: done
    at: 2026-10-14T23:15:30Z
    by: agent-S-0388
stream: S-0388
tags: [dashboard]
touches: ["flaiover/src/routes/api/items/[id]/restore/+server.ts", "flaiover/src/routes/api/items/[id]/restore/restore.test.ts"]
usage:
  source: log
  seconds: 310
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 8609
      cache_read: 1249907
      cache_write: 46228
      cost: 0.6912
---
# T-1577 POST /api/items/[id]/restore runs the host write stream.restore for one path

## Work

- Add `flaiover/src/routes/api/items/[id]/restore/+server.ts`, after the pattern of `api/items/[id]/finalize/+server.ts`. A POST with `{ path }` calls `repo().write('stream.restore', …)` for the story and that path through `respond(...)`.
- Answer 400 when `path` is missing or empty. Pass flai's refusal text through verbatim (ADR-0016).
- Take the argument names from the `stream.restore` entry S-0371's T-1575 adds to `flai/internal/hostapi/writes.go`. Read them there before writing.
- Leave `flaiover/src/lib/server/agent.ts` alone: T-1575 adds `stream.restore` to `REQUIRED_METHODS`.
- It waits for no task of this story.

## Done when

- [ ] `restore.test.ts` covers a write with the story and the path, a missing path answered 400, and a refusal passed through
- [ ] `flai test` on the two files passes

## Notes
