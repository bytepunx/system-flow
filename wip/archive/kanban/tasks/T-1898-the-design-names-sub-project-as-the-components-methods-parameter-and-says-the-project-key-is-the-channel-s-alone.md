---
id: T-1898
type: task
nature: feature
title: The design names sub_project as the components methods' parameter and says the project key is the channel's alone
status: done
parent: S-0448
owner: alex
created: 2026-10-18T02:28:39Z
updated: 2026-10-18T02:37:15Z
transitions:
  - to: ready
    at: 2026-10-18T02:28:59Z
    by: agent-S-0448
  - to: in-progress
    at: 2026-10-18T02:36:47Z
    by: agent-S-0448
  - to: done
    at: 2026-10-18T02:37:15Z
    by: agent-S-0448
stream: S-0448
tags: []
touches: [design/system/component-library.md, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, docs/users/flai-reference.md]
after: [T-1896, T-1897]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 3285
      cache_read: 526699
      cache_write: 18399
      cost: 0.2939
---

# T-1898 The design names sub_project as the components methods' parameter and says the project key is the channel's alone

## Work

- `design/system/component-library.md` and `design/system/dashboard-host-channel.md` name `sub_project` as the components methods' parameter and say why `project` is the channel's alone (I-0163).
- `design/system/flaiover-dashboard.md` says the routes' `project` body field reaches flai as `sub_project`.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference` only if it names the methods' parameter.
- Waits for T-1896 and T-1897, so that it describes what they landed.

## Done when

- The three design documents name `sub_project` and the channel's `project` key, with `updated` bumped.
- `flai test` on the changed paths passes.

## Notes
