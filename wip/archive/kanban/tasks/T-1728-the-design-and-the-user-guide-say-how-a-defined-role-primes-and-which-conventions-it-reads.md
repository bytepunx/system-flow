---
id: T-1728
type: task
nature: feature
title: The design and the user guide say how a defined role primes and which conventions it reads
status: done
parent: S-0415
owner: alex
created: 2026-10-16T05:43:56Z
updated: 2026-10-28T20:29:04Z
transitions:
  - to: ready
    at: 2026-10-28T20:28:01Z
    by: agent-S-0415
  - to: in-progress
    at: 2026-10-28T20:28:02Z
    by: agent-S-0415
  - to: done
    at: 2026-10-28T20:29:04Z
    by: agent-S-0415
stream: S-0415
tags: [cli]
touches: [design/system/agent-context.md, design/system/conventions.md, docs/users/flai.md, design/system/flai-cli.md]
after: [T-1726, T-1727]
usage:
  source: log
  seconds: 62
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5836
      cache_read: 1246086
      cache_write: 38065
      cost: 0.6705
---
# T-1728 The design and the user guide say how a defined role primes and which conventions it reads

## Work

- In `design/system/agent-context.md` § Strategic packs, say that a role's topics and item kind come from its role file.
- In `design/system/conventions.md` § Roles, say that a convention's `roles` may list any defined role.
- In `docs/users/flai.md`, describe `flai prime --role` for a defined role.

Last layer: it describes what T-1726 and T-1727 built.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
