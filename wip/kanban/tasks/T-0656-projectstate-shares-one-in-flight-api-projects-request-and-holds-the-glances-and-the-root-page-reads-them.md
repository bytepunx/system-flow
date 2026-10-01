---
id: T-0656
type: task
nature: feature
title: projectState shares one in-flight /api/projects request and holds the glances, and the root page reads them
status: done
parent: S-0186
owner: alex
created: 2026-10-01T09:18:09Z
updated: 2026-10-01T09:20:07Z
transitions:
  - to: ready
    at: 2026-10-01T09:18:19Z
    by: agent-S-0186
  - to: in-progress
    at: 2026-10-01T09:18:19Z
    by: agent-S-0186
  - to: done
    at: 2026-10-01T09:20:07Z
    by: agent-S-0186
stream: S-0186
tags: []
touches: [flaiover/src/lib/project.svelte.ts, flaiover/src/routes/+page.svelte, flaiover/src/routes/overview.svelte.test.ts, flaiover/src/lib/project.svelte.test.ts]
usage:
  source: log
  seconds: 108
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 27
      output: 6069
      cache_read: 1369891
      cache_write: 26567
      cost: 0.5747
---
# T-0656 projectState shares one in-flight /api/projects request and holds the glances, and the root page reads them

## Work

- `projectState.refresh()` returns the request already in flight when there is one, so the layout's and the page's calls on one load make one `/api/projects` request.
- `Project` carries the glance fields (`review`, `threadsAwaiting`, `agentAttending`); the root page's list reads `projectState.list` and `loadGlances()` goes.
- Tests: concurrent `refresh()` calls make one request and a later call after it settles makes another; mounting the root page together with a layout-style `refresh()` makes one `/api/projects` request.

## Done when

- One root page load makes one `/api/projects` request, a test counts it, and the flaiover unit tests, `svelte-check`, and lint pass.

## Notes
