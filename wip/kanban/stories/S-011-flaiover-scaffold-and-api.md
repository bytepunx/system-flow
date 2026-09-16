---
id: S-011
type: story
nature: feature
title: SvelteKit scaffold and repo reader API
status: backlog
parent: E-003
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T16:22:05Z
transitions: []
tags: []
---

# S-011 SvelteKit scaffold and repo reader API

## Goal
SvelteKit 2, Svelte 5, Tailwind 4, adapter-node, ssr off, with /api endpoints that read the manifest, items, and documents from PROJECT_DIR.

## Acceptance criteria
- [ ] pnpm dev serves the app against this repo
- [ ] /api/items returns parsed front matter for every item
- [ ] /api/docs/tree returns the design, docs, and wip trees
- [ ] File watcher invalidates caches and emits SSE

## Tasks

## Notes
