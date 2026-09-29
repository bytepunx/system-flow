---
id: S-0162
type: story
nature: improvement
title: The items and documents pages ask for what they show, not the whole archive and every file's front matter
status: backlog
parent: E-0012
owner: alex
created: 2026-09-29T07:00:31Z
updated: 2026-09-29T07:00:31Z
transitions: []
tags: [dashboard]
topics: [server-side, back-end]
touches: [flaiover/src, flai/internal/hostapi]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0162 The items and documents pages ask for what they show, not the whole archive and every file's front matter

## Goal

Two answers are large and rebuilt on every ask: `items.list` with the archive and bodies is 1.28 MB, and `docs.tree` with the front matter of every file under `design`, `docs`, and `wip` is 772 KB after a 127 ms walk. Cause 7 of `design/system/server-performance.md`.

## Acceptance criteria
- [ ] The items page's first answer is under 200 KB on this repository: bodies are asked for when an item is opened, and the archive when it is shown.
- [ ] `docs.tree` is under 200 KB, or `flai serve` keeps it until a file under the three folders changes, and a warm answer takes under 10 ms.
- [ ] The readiness probe, `/_ready`, asks for no more than the manifest to know flai answers.

## Tasks

## Notes

Measured by S-0152: `items.list` 122 ms and 1,278,709 bytes; `docs.tree` 135 ms and 772,454 bytes. `/_ready` asks `repo().items()`, the whole archive with bodies.
