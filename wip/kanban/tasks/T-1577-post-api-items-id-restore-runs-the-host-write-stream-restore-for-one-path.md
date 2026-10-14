---
id: T-1577
type: task
nature: feature
title: "POST /api/items/[id]/restore runs the host write stream.restore for one path"
status: backlog
parent: S-0388
owner: alex
created: 2026-10-14T21:08:46Z
updated: 2026-10-14T21:08:46Z
transitions: []
stream: S-0388
tags: [dashboard]
touches: ["flaiover/src/routes/api/items/[id]/restore/+server.ts", "flaiover/src/routes/api/items/[id]/restore/restore.test.ts"]
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
