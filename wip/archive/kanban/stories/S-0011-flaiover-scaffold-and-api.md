---
id: S-0011
type: story
nature: feature
title: SvelteKit scaffold and repo reader API
status: done
parent: E-0003
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T04:24:09Z
transitions:
  - to: ready
    at: 2026-09-17T04:02:09Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:02:09Z
    by: agent
  - to: review
    at: 2026-09-17T04:09:27Z
    by: agent
  - to: done
    at: 2026-09-17T04:24:09Z
    by: alex
tags: []
---

# S-0011 SvelteKit scaffold and repo reader API

## Goal
SvelteKit 2, Svelte 5, Tailwind 4, adapter-node, ssr off, with /api endpoints that read the manifest, items, and documents from PROJECT_DIR.

## Acceptance criteria
- [x] pnpm dev serves the app against this repo
- [x] /api/items returns parsed front matter for every item
- [x] /api/docs/tree returns the design, docs, and wip trees
- [x] File watcher invalidates caches and emits SSE

## Tasks
- T-0082 SvelteKit 2, Svelte 5, Tailwind 4, adapter-node scaffold with ssr off and the project conventions applied
- T-0083 Repo reader: manifest, items, docs tree and file, with a mtime cache and a watcher emitting SSE
- T-0084 API routes /api/manifest, /api/items, /api/items/:id, /api/docs/tree, /api/docs/file, /api/events with tests
- T-0085 Scripts, Makefile, CI workflow for flaiover; design and tech docs

## Notes
- Real pins: Vite 8, TypeScript 6, Vitest 4, Playwright 1.60, ESLint 10; design/tech updated from the lockfile.
