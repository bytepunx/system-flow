---
title: SvelteKit, Svelte, Vite, TypeScript, tests
updated: 2026-09-15
status: active
---

# SvelteKit stack

| Package | Version | Purpose |
|---------|---------|---------|
| `@sveltejs/kit` | 2 | App framework, routing, server endpoints |
| `svelte` | 5 | UI, runes for state |
| `@sveltejs/adapter-node` | 5 | Node server output for the Docker image |
| `vite` | 7 | Build and dev server |
| `typescript` | 5 | Everywhere, `strict` on |
| `vitest` | latest | Unit tests, including the metrics port fixture test |
| `@playwright/test` | latest | End-to-end against the template sample repo |

Decision: [ADR 0007](../adrs/0007-sveltekit-spa-with-node-adapter-in-docker.md).

## Why

- The brief asks for SvelteKit and Tailwind. Svelte 5 runes give simple client state for the board and filters.
- `adapter-node` because the app must read a mounted filesystem; a static adapter cannot. The UI still runs as an SPA with `ssr = false`.
- Vitest and Playwright are the SvelteKit defaults and need no extra integration.
