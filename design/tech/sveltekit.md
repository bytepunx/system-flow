---
title: SvelteKit, Svelte, Vite, TypeScript, tests
updated: 2026-09-15
status: active
---

# SvelteKit stack

| Package | Version | Purpose |
|---------|---------|---------|
| `@sveltejs/kit` | 2.63 | App framework, routing, server endpoints; the adapter is configured in `vite.config.ts` |
| `svelte` | 5.56 | UI, runes for state (runes forced on for project files) |
| `@sveltejs/adapter-node` | 5.5 | Node server output for the Docker image |
| `vite` | 8 | Build and dev server |
| `typescript` | 6 | Everywhere, `strict` on |
| `vitest` | 4.1 | Unit tests, including the metrics port fixture test; a `client` project runs `*.svelte.test.ts` component tests in jsdom |
| `jsdom` | 30.1 | DOM for Svelte component tests (S-0045) |
| `@playwright/test` | 1.60 | End-to-end against the template sample repo |
| `yaml` | 2.9 | Front matter and manifest parsing on the server |
| `chokidar` | 5.0 | File watching behind `/api/events` |
| `eslint`, `prettier` | 10, 3.8 | Lint and format as scaffolded by `sv create` |

Exact versions are in `flaiover/pnpm-lock.yaml`. Scaffolded with `sv create` on 2026-09-17 (S-0011).

Decision: [ADR 0007](../adrs/0007-sveltekit-spa-with-node-adapter-in-docker.md).

## Why

- The brief asks for SvelteKit and Tailwind. Svelte 5 runes give simple client state for the board and filters.
- `adapter-node` because the app must read a mounted filesystem; a static adapter cannot. The UI still runs as an SPA with `ssr = false`.
- Vitest and Playwright are the SvelteKit defaults and need no extra integration.
