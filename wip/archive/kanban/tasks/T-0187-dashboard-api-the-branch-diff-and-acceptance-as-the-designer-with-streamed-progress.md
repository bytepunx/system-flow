---
id: T-0187
type: task
nature: feature
title: "Dashboard API: the branch diff, and acceptance as the designer with streamed progress"
status: done
parent: S-0041
owner: alex
created: 2026-09-19T04:08:29Z
updated: 2026-09-19T04:12:35Z
transitions:
  - to: ready
    at: 2026-09-19T04:11:04Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:11:04Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:12:35Z
    by: system-flow
stream: S-0041
tags: []
touches: [flaiover/src/lib/server, flaiover/src/routes/api]
---

# T-0187 Dashboard API: the branch diff, and acceptance as the designer with streamed progress

## Work
`GET /api/items/:id/diff` returns `flai stream diff`. `POST /api/items/:id/accept` spawns flai (not `execFile`, which buffers), passes `--by` with `designer()`, and with `{ include_uncommitted: true }` passes `--yes` (S-0051); it writes one JSON line per flai info event as it arrives and a last line with the result, or with the error verbatim and the blockers when flai refuses. Add a streaming variant beside `flai()` in `src/lib/server/flai.ts` for it. Tests with a fake flai, a script standing in for the binary through `FLAI_BIN`, as the criterion asks: progress lines arrive before the result, the designer is passed as `--by`, a failure's message reaches the client verbatim with a non-200 final line, and the diff endpoint passes flai's JSON through.

## Done when
- The endpoint tests pass with the fake flai
- flaiover lint and svelte-check are clean

## Notes
