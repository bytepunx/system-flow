---
id: D-0393
title: SvelteKit refuses an export from a +server.ts file that is neither a request handler nor underscore-prefixed, so a helper kept there is named with a leading underscore
kind: decision
sources: [S-0051]
topics: [sveltekit, flaiover, dashboard]
paths: ["flaiover/src/routes/api/items/[id]/move/+server.ts"]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0393 SvelteKit refuses an export from a +server.ts file that is neither a request handler nor underscore-prefixed, so a helper kept there is named with a leading underscore

S-0051 kept its argument helper `_moveArgs` in the move endpoint because "SvelteKit rejects exports from `+server.ts` that are not handlers or underscore-prefixed". The build enforces it, so it holds for every endpoint in flaiover: prefix the helper with `_`, or move it to a module under `src/lib`.

Source: S-0051.
