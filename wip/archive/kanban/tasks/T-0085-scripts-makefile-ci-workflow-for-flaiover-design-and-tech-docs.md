---
id: T-0085
type: task
nature: feature
title: "Scripts, Makefile, CI workflow for flaiover; design and tech docs"
status: done
parent: S-0011
owner: alex
created: 2026-09-17T04:02:08Z
updated: 2026-09-17T04:09:27Z
transitions:
  - to: ready
    at: 2026-09-17T04:09:26Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:09:27Z
    by: agent
  - to: done
    at: 2026-09-17T04:09:27Z
    by: agent
stream: S-0011
tags: [dashboard, ci, docs]
---

# T-0085 Scripts, Makefile, CI workflow for flaiover; design and tech docs

## Work
scripts/flaiover-*.sh and Makefile targets (install, dev, build, test), .github/workflows/flaiover.yml (lint, unit tests, build; e2e later), design/tech entries with pinned versions (sveltekit, svelte, vite, tailwind, adapter-node, yaml, chokidar, vitest, playwright, pnpm, node), flaiover-dashboard.md updated to what exists.

## Done when
CI green; tech docs list the real versions.

## Notes
