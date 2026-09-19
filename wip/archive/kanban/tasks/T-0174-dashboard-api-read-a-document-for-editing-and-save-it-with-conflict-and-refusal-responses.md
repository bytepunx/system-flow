---
id: T-0174
type: task
nature: feature
title: "Dashboard API: read a document for editing and save it, with conflict and refusal responses"
status: done
parent: S-0040
owner: alex
created: 2026-09-19T02:11:29Z
updated: 2026-09-19T02:19:33Z
transitions:
  - to: ready
    at: 2026-09-19T02:17:59Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:17:59Z
    by: alex
  - to: done
    at: 2026-09-19T02:19:33Z
    by: alex
stream: S-0040
tags: []
touches: [flaiover/src/lib/server, flaiover/src/routes/api]
---

# T-0174 Dashboard API: read a document for editing and save it, with conflict and refusal responses

## Work
`GET /api/docs/edit?path=` returns `flai doc show`. `PUT /api/docs/file` with `{ path, content, hash, message? }` runs `flai doc save`, sending the content on stdin (extend `src/lib/server/flai.ts` with an input option, no shell). Map a conflict to 409 with `{ error, current, hash, diff }`, a check refusal or an ownership refusal to 422 with `{ error, findings }`, and leave other failures as they are. The token rules of ADR-0018 already cover the routes. Tests against a temp copy of the fixture with the built flai: save, conflict, refusal, body-mode refusal.

## Done when
- The endpoint tests pass
- flaiover lint and svelte-check are clean

## Notes
