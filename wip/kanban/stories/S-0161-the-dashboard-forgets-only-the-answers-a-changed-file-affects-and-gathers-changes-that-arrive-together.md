---
id: S-0161
type: story
nature: improvement
title: The dashboard forgets only the answers a changed file affects, and gathers changes that arrive together
status: in-progress
parent: E-0012
owner: alex
created: 2026-09-29T07:00:30Z
updated: 2026-09-29T19:40:45Z
transitions:
  - to: ready
    at: 2026-09-29T19:18:54Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:38:13Z
    by: agent-S-0161
tags: [dashboard]
topics: [server-side, back-end]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 444
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 100
      output: 759
      cache_read: 5351914
      cache_write: 135988
      cost: 2.1845
---
# S-0161 The dashboard forgets only the answers a changed file affects, and gathers changes that arrive together

## Goal

The dashboard forgets every answer it holds when any file of the project changes, and the board asks for its board, `/api/publish`, and the agents again at each `change` event, one per file, whatever the file. With agents working, a narrative or thread changes every few seconds, so pages keep asking flai for everything. Cause 6 of `design/system/server-performance.md`.

## Acceptance criteria
- [x] A change to a narrative, a thread, or a document does not make the board ask for `board.get` or `publish.preview` again; a change to a work item or `board.md` does.
- [x] Changes that arrive within half a second of each other make a page ask again once.
- [x] The answers kept for other paths survive a change that does not affect them: behaviour tests in `flaiover/src/lib/server` cover which paths forget which answers.

## Tasks
- T-0571 The dashboard's Repo forgets only the answers a changed path affects
- T-0572 Pages gather changes that arrive together, and the board asks again only when a work item changed

## Notes

Measured by S-0152: `/api/publish` 56 and `/api/unpushed` 58 requests in 7 minutes, about every 7.5 s, each 160 to 235 ms. `Repo.changed` in `flaiover/src/lib/server/repo.ts` clears every answer; `routes/board/+page.svelte` reloads at each `change`. Related: S-0154 (story pages receive live updates) listens to the same events.

Verified by `flaiover/src/lib/server/repo-forget.test.ts` (which paths forget which answers: a narrative, thread, or document change keeps `board.get`; a work item, `board.md`, or archive change forgets it), `flaiover/src/routes/board/board.svelte.test.ts` (the board asks for `/api/board`, `/api/publish`, and `/api/unpushed` only for work item or manifest changes, once for changes within 500 ms), and `flaiover/src/lib/events.test.ts` (gathering). The request rate in the running dashboard was not measured again: that needs an image built from this story in place of the operator's container.
