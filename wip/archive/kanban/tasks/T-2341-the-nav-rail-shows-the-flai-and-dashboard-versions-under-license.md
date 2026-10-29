---
id: T-2341
type: task
nature: feature
title: The nav rail shows the flai and dashboard versions under License
status: done
parent: S-0544
owner: alex
created: 2026-10-29T17:40:32Z
updated: 2026-10-29T17:43:30Z
transitions:
  - to: ready
    at: 2026-10-29T17:40:53Z
    by: agent-S-0544
  - to: in-progress
    at: 2026-10-29T17:40:55Z
    by: agent-S-0544
  - to: done
    at: 2026-10-29T17:43:30Z
    by: agent-S-0544
stream: S-0544
tags: []
touches: [flaiover/src/routes/api/agent/+server.ts, flaiover/src/routes/api/agent/agent.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 155
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 6680
      cache_read: 2330054
      cache_write: 53693
      cost: 1.0293
---
# T-2341 The nav rail shows the flai and dashboard versions under License

## Work

`GET /api/agent` answers `dashboard`, this dashboard's own version (`$lib/server/metrics`' `version`), in every answer, beside `flai`, the connected flai's. `HostFlaiStatus` carries it. `NavRail.svelte` shows a label under License with `flai <version>` and `dashboard <version>` on two lines, each line only when its version is known. Docs: the nav rail paragraph in `docs/users/flaiover.md` and the rail's foot in `design/system/flaiover-dashboard.md`.

Waits for nothing: it is the first layer.

## Done when

- The rail's foot shows both versions, newline separated, under License, and a test pins it.
- `/api/agent` answers `dashboard` whether or not a flai is connected, and a test pins it.
- The docs say so.

## Notes
