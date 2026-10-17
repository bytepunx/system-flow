---
id: T-1883
type: task
nature: feature
title: The operator and user docs and the component library's design describe a story's dev server, the install step, the cap, and the review link
status: done
parent: S-0442
owner: alex
created: 2026-10-16T21:37:41Z
updated: 2026-10-17T03:42:25Z
transitions:
  - to: ready
    at: 2026-10-17T03:26:46Z
    by: agent-S-0442
  - to: in-progress
    at: 2026-10-17T03:26:47Z
    by: agent-S-0442
  - to: done
    at: 2026-10-17T03:42:25Z
    by: agent-S-0442
stream: S-0442
tags: [dashboard, cli]
touches: [docs/operators/index.md, docs/operators/settings.md, docs/users/flaiover.md, design/system/component-library.md, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, docs/users/flai.md]
after: [T-1880, T-1881, T-1882]
usage:
  source: log
  seconds: 938
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 88
      output: 35579
      cache_read: 6320308
      cache_write: 166255
      cost: 2.9503
---
# T-1883 The operator and user docs and the component library's design describe a story's dev server, the install step, the cap, and the review link

## Work

- `docs/operators/index.md` § The components host action: a story's dev server, the dependency install it may run, and when it stops.
- `docs/operators/settings.md`: a row for `components.max_servers`.
- `docs/users/flaiover.md`: the source picker and the review page's link.
- `design/system/flaiover-dashboard.md`: the picker, `story` on the components routes, and the review page's link.
- `design/system/component-library.md`: a story's worktree as built.

Waits for T-1880, T-1881, and T-1882, so that the docs describe all three as built.

## Done when

The five files describe the behaviour, the settings test passes, and `flai test` on them passes.

## Notes

None.
