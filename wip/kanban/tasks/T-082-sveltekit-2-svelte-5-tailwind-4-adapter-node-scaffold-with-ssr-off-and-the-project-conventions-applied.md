---
id: T-082
type: task
nature: feature
title: SvelteKit 2, Svelte 5, Tailwind 4, adapter-node scaffold with ssr off and the project conventions applied
status: done
parent: S-011
owner: alex
created: 2026-09-17T04:02:08Z
updated: 2026-09-17T04:09:25Z
transitions:
  - to: ready
    at: 2026-09-17T04:09:25Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:09:25Z
    by: agent
  - to: done
    at: 2026-09-17T04:09:25Z
    by: agent
stream: S-011
tags: [dashboard]
---

# T-082 SvelteKit 2, Svelte 5, Tailwind 4, adapter-node scaffold with ssr off and the project conventions applied

## Work
Create flaiover/ with SvelteKit 2, Svelte 5, TypeScript strict, Tailwind 4 via the Vite plugin, adapter-node, ssr = false at the root layout, pnpm with packageManager pinned, .nvmrc 24, ESLint and Prettier per the scaffold, vitest and playwright configured; a minimal app shell page that reads /api/manifest.

## Done when
pnpm dev serves the shell against this repository; pnpm build succeeds.

## Notes
