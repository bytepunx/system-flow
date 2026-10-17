---
id: T-1859
type: task
nature: feature
title: flai embeds a Vite plugin and a harness page that render one component by path after the setup module, and writes the merged Vite config
status: done
parent: S-0438
owner: alex
created: 2026-10-16T21:33:39Z
updated: 2026-10-16T23:28:23Z
transitions:
  - to: ready
    at: 2026-10-16T23:12:21Z
    by: system-flow
  - to: in-progress
    at: 2026-10-16T23:12:21Z
    by: system-flow
  - to: done
    at: 2026-10-16T23:28:23Z
    by: system-flow
stream: S-0438
tags: [cli]
touches: [flai/internal/components/harness/plugin.mjs, flai/internal/components/harness/harness.html, flai/internal/components/harness/harness.mjs, flai/internal/components/harness.go, flai/internal/components/harness_test.go]
after: [T-1858]
usage:
  source: log
  seconds: 962
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 94
      output: 456
      cache_read: 5002073
      cache_write: 164090
      cost: 2.4017
---
# T-1859 flai embeds a Vite plugin and a harness page that render one component by path after the setup module, and writes the merged Vite config

## Work

- `harness/plugin.mjs`: a Vite plugin that serves `/__flai/components/` from `harness.html` through `server.transformIndexHtml`, and resolves a virtual module that imports the `setup` module and maps each component path to `import()`.
- `harness/harness.mjs`: it mounts the component named by `?c=` or a `select` message with Svelte's `mount`, inside a boundary that posts `error` with the file and line, and posts `ready`, `rendered`, and `size` (from a `ResizeObserver`). It accepts messages only from the origin it is given in the URL.
- `harness.go`: embeds the folder and writes it, with a `vite.config.mjs` that imports the sub-project's own config and merges flai's plugin with Vite's `mergeConfig`, under `.flai-cache/components/<project>/<sub-project>/`.

Waits for T-1858, whose ADR fixes the protocol.

## Done when

Tests check the written files and the merged config for a SvelteKit and a plain Svelte fixture. `flai test flai/internal/components` passes. The live run is T-1862's.

## Notes

If SvelteKit's plugin does not let the merged config serve a page outside its routes, try `configureServer` middleware ahead of SvelteKit's, and record which held in the story's `## Decisions`.
