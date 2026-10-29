---
id: T-2316
type: task
nature: feature
title: integration.sh gives go test a 20-minute timeout per package, so the mcpserver suite, which takes 8 to 10 minutes, does not fail a batch at Go's 10-minute default
status: done
parent: S-0507
owner: alex
created: 2026-10-28T23:28:27Z
updated: 2026-10-28T23:28:53Z
transitions:
  - to: ready
    at: 2026-10-28T23:28:33Z
    by: agent-S-0507
  - to: in-progress
    at: 2026-10-28T23:28:34Z
    by: agent-S-0507
  - to: done
    at: 2026-10-28T23:28:53Z
    by: agent-S-0507
stream: S-0507
tags: []
touches: [scripts/integration.sh]
usage:
  source: log
  seconds: 19
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 3511
      cache_read: 638974
      cache_write: 17719
      cost: 0.3135
---

# T-2316 integration.sh gives go test a 20-minute timeout per package, so the mcpserver suite, which takes 8 to 10 minutes, does not fail a batch at Go's 10-minute default

## Work

Merge queue batch 20261028-231632 failed S-0507 at integration (TH-0723): `flai/internal/mcpserver` hit Go's default 10-minute test timeout at 600.17 s. Its suite took 470 to 591 s in every batch of 2026-10-28, and S-0507 changes nothing in it. Give `go test` in `scripts/integration.sh` `-timeout 20m`, so a package as slow as mcpserver passes, and a real hang still fails.

## Done when

`scripts/integration.sh` passes `-timeout 20m` to `go test`, ShellCheck passes, and the close-out passes.

## Notes
