---
id: S-0161
type: story
nature: improvement
title: The dashboard forgets only the answers a changed file affects, and gathers changes that arrive together
status: backlog
parent: E-0012
owner: alex
created: 2026-09-29T07:00:30Z
updated: 2026-09-29T07:00:30Z
transitions: []
tags: [dashboard]
topics: [server-side, back-end]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0161 The dashboard forgets only the answers a changed file affects, and gathers changes that arrive together

## Goal

The dashboard forgets every answer it holds when any file of the project changes, and the board asks for its board, `/api/publish`, and the agents again at each `change` event, one per file, whatever the file. With agents working, a narrative or thread changes every few seconds, so pages keep asking flai for everything. Cause 6 of `design/system/server-performance.md`.

## Acceptance criteria
- [ ] A change to a narrative, a thread, or a document does not make the board ask for `board.get` or `publish.preview` again; a change to a work item or `board.md` does.
- [ ] Changes that arrive within half a second of each other make a page ask again once.
- [ ] The answers kept for other paths survive a change that does not affect them: behaviour tests in `flaiover/src/lib/server` cover which paths forget which answers.

## Tasks

## Notes

Measured by S-0152: `/api/publish` 56 and `/api/unpushed` 58 requests in 7 minutes, about every 7.5 s, each 160 to 235 ms. `Repo.changed` in `flaiover/src/lib/server/repo.ts` clears every answer; `routes/board/+page.svelte` reloads at each `change`. Related: S-0154 (story pages receive live updates) listens to the same events.
