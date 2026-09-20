---
id: S-0073
type: story
nature: feature
title: "The dashboard reads work items through the channel: project, board, items, threads, and file changes"
status: ready
parent: E-0003
owner: alex
created: 2026-09-20T07:26:50Z
updated: 2026-09-20T07:27:44Z
transitions:
  - to: ready
    at: 2026-09-20T07:27:44Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system]
---
# S-0073 The dashboard reads work items through the channel: project, board, items, threads, and file changes

## Goal
The first reads leave the mount. flai answers the project's identity and layout, the board, the item list with the archive, one item with its children, and threads as structured data over the channel, and tells the dashboard when files change. The dashboard stops reading and parsing those files itself.

## Acceptance criteria
- [ ] Methods for the project (name, key, owner, layout, the dashboard settings the server needs), the board as `flai board --json` gives it, items by type, state, and archive with bodies on request, one item with its children, threads with their entries and by anchor; each is the same data the CLI prints, from the same Go code, and each is tested against the fixtures
- [ ] A watcher in flai on the manifest's design, docs, and wip folders and the manifest sends a `change` notification per file, debounced as chokidar is today; `/api/events`, the stats and inbox caches, and the notifier are fed from it and behave as before
- [ ] `/api/manifest`, `/api/board`, `/api/items`, `/api/items/[id]`, `/api/threads`, and `/_ready` answer from the channel; `board.ts`'s own composition and `repo.ts`'s item and thread parsing are deleted, not kept beside the new path
- [ ] With no flai connected those routes answer 503 with a message the pages show: what is missing and the command that starts it
- [ ] A page load's requests run concurrently over the one connection; the board of this repository loads no slower than today, measured and recorded
- [ ] The design documents and the users' documentation say where the dashboard's data comes from

## Tasks

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on the channel story. The operator chose structured methods over restricted file reads (ADR-0029), so nothing here exposes a file.
