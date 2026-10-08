---
id: T-1212
type: task
nature: feature
title: flaiover's repo reads messages from the host, and /api/messages serves them
status: done
parent: S-0336
owner: alex
created: 2026-10-07T20:16:56Z
updated: 2026-10-08T06:21:01Z
transitions:
  - to: ready
    at: 2026-10-08T06:17:39Z
    by: agent-S-0336
  - to: in-progress
    at: 2026-10-08T06:17:39Z
    by: agent-S-0336
  - to: done
    at: 2026-10-08T06:21:01Z
    by: agent-S-0336
stream: S-0336
tags: [flaiover]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/repo.test.ts, flaiover/src/routes/api/messages/+server.ts, flaiover/src/routes/api/messages/messages.test.ts]
after: [T-1211]
usage:
  source: log
  seconds: 202
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 234
      cache_read: 1574687
      cache_write: 65881
      cost: 0.7319
---
# T-1212 flaiover's repo reads messages from the host, and /api/messages serves them

## Work

Carry T-1211's reads into the dashboard's server. It waits for T-1211, whose reads it calls.

- `repo().messages()` and `repo().messagesFor(story)` in `repo.ts`, typed, beside `threads()` and `threadsFor()`.
- `GET /api/messages?story=<id>&all=1` in `src/routes/api/messages/+server.ts`, as `/api/threads` does; no POST.

## Done when

- Tests cover the repo calls and the route, with and without `story` and `all`.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
