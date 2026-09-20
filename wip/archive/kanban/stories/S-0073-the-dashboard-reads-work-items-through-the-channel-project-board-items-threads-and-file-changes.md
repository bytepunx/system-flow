---
id: S-0073
type: story
nature: feature
title: "The dashboard reads work items through the channel: project, board, items, threads, and file changes"
status: done
parent: E-0003
owner: alex
created: 2026-09-20T07:26:50Z
updated: 2026-09-20T12:05:11Z
transitions:
  - to: ready
    at: 2026-09-20T07:27:44Z
    by: alex
  - to: in-progress
    at: 2026-09-20T07:50:50Z
    by: system-flow
  - to: review
    at: 2026-09-20T08:23:57Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:05:11Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system]
---
# S-0073 The dashboard reads work items through the channel: project, board, items, threads, and file changes

## Goal
The first reads leave the mount. flai answers the project's identity and layout, the board, the item list with the archive, one item with its children, and threads as structured data over the channel, and tells the dashboard when files change. The dashboard stops reading and parsing those files itself.

## Acceptance criteria
- [x] Methods for the project (name, key, owner, layout, the dashboard settings the server needs), the board as `flai board --json` gives it, items by type, state, and archive with bodies on request, one item with its children, threads with their entries and by anchor; each is the same data the CLI prints, from the same Go code, and each is tested against the fixtures
- [x] A watcher in flai on the manifest's design, docs, and wip folders and the manifest sends a `change` notification per file, debounced as chokidar is today; `/api/events`, the stats and inbox caches, and the notifier are fed from it and behave as before
- [x] `/api/manifest`, `/api/board`, `/api/items`, `/api/items/[id]`, `/api/threads`, and `/_ready` answer from the channel; `board.ts`'s own composition and `repo.ts`'s item and thread parsing are deleted, not kept beside the new path
- [x] With no flai connected those routes answer 503 with a message the pages show: what is missing and the command that starts it
- [x] A page load's requests run concurrently over the one connection; the board of this repository loads no slower than today, measured and recorded
- [x] The design documents and the users' documentation say where the dashboard's data comes from

## Tasks
- T-0263 flai: methods for the project, the board, items, one item, and threads, from the code the CLI prints with
- T-0264 flai: a watcher that tells the dashboard which files changed
- T-0265 flaiover: the manifest, board, items, item, threads, readiness, and events come from the channel; the file parsing for them is deleted
- T-0266 Measure the board's load against today's, try it end to end, and document where the data comes from

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on the channel story. The operator chose structured methods over restricted file reads (ADR-0029), so nothing here exposes a file.
