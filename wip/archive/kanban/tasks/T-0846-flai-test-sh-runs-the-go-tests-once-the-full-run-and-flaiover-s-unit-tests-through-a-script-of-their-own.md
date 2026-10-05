---
id: T-0846
type: task
nature: feature
title: flai-test.sh runs the Go tests once, the full run, and flaiover's unit tests through a script of their own
status: done
parent: S-0267
owner: alex
created: 2026-10-05T00:18:45Z
updated: 2026-10-05T00:19:53Z
transitions:
  - to: ready
    at: 2026-10-05T00:19:02Z
    by: agent-S-0267
  - to: in-progress
    at: 2026-10-05T00:19:03Z
    by: agent-S-0267
  - to: done
    at: 2026-10-05T00:19:53Z
    by: agent-S-0267
stream: S-0267
tags: []
touches: [scripts/flaiover-unit.sh, scripts/test.sh, scripts/flai-test.sh, Makefile]
usage:
  source: log
  seconds: 50
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2756
      cache_read: 490793
      cache_write: 13127
      cost: 0.2583
---
# T-0846 flai-test.sh runs the Go tests once, the full run, and flaiover's unit tests through a script of their own

## Work

Move the flaiover vitest block out of `scripts/test.sh` into `scripts/flaiover-unit.sh`, which builds `bin/flai` when it is missing or stale and runs `pnpm test:unit` only when `flaiover/node_modules` is present. `scripts/test.sh` runs the short Go tests and then that script, so `make test` behaves as before. `scripts/flai-test.sh` runs the lint, `scripts/flaiover-unit.sh`, `scripts/integration.sh`, and `scripts/smoke.sh`, no longer `scripts/test.sh`, so the Go tests run once, with the race detector, in their full form. The Makefile's `flai-test` help says so. This task waits for none.

## Done when

- `scripts/flai-test.sh` runs `go test` once, without `-short`, and still runs vitest when `flaiover/node_modules` is present
- `make test` and `make integration` run what they ran before

## Notes
