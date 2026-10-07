---
id: T-1212
type: task
nature: feature
title: flaiover's repo reads messages from the host, and /api/messages serves them
status: backlog
parent: S-0336
owner: alex
created: 2026-10-07T20:16:56Z
updated: 2026-10-07T20:16:56Z
transitions: []
stream: S-0336
tags: [flaiover]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/repo.test.ts, flaiover/src/routes/api/messages/+server.ts, flaiover/src/routes/api/messages/messages.test.ts]
after: [T-1211]
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
