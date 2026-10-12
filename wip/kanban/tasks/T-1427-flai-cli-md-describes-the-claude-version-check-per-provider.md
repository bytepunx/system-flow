---
id: T-1427
type: task
nature: improvement
title: flai-cli.md describes the claude version check per provider
status: done
parent: S-0360
owner: alex
created: 2026-10-08T08:59:51Z
updated: 2026-10-12T20:49:41Z
transitions:
  - to: ready
    at: 2026-10-12T20:48:56Z
    by: agent-S-0360
  - to: in-progress
    at: 2026-10-12T20:48:57Z
    by: agent-S-0360
  - to: done
    at: 2026-10-12T20:49:41Z
    by: agent-S-0360
stream: S-0360
tags: [cli]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md]
after: [T-1426]
usage:
  source: log
  seconds: 39
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 10
      output: 73
      cache_read: 970976
      cache_write: 13071
      cost: 0.4477
---
# T-1427 flai-cli.md describes the claude version check per provider

## Work

- `design/system/flai-cli.md`, where the check of each new `claude` version is described: it runs through the project's provider, the record per version and provider, and the skip on an unset key.

It waits for T-1426, so that it describes the check as built.

## Done when

- `flai test design/system/flai-cli.md` passes.

## Notes

Layer 2 of S-0360.
