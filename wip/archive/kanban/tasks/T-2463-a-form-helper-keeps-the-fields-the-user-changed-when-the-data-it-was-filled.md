---
id: T-2463
type: task
nature: feature
title: A form helper keeps the fields the user changed when the data it was filled from is read again
status: done
parent: S-0574
owner: alex
created: 2026-11-06T18:02:17Z
updated: 2026-11-06T18:05:00Z
transitions:
  - to: ready
    at: 2026-11-06T18:03:07Z
    by: system-flow
  - to: in-progress
    at: 2026-11-06T18:03:10Z
    by: system-flow
  - to: done
    at: 2026-11-06T18:05:00Z
    by: agent-S-0574
stream: S-0574
tags: []
touches: [flaiover/src/lib/form.ts, flaiover/src/lib/form.test.ts, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 110
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 5815
      cache_read: 2542990
      cache_write: 39980
      cost: 0.9271
---
# T-2463 A form helper keeps the fields the user changed when the data it was filled from is read again

## Work

Add `flaiover/src/lib/form.ts` with `rebase(seed, current, fresh)`: for each leaf of a plain object, keep the current value where it differs from the seed it was filled from (the user changed it), else take the fresh read. Keys the fresh read lacks go; keys it adds come in. Behaviour tests in `form.test.ts`. Say in `design/system/flaiover-dashboard.md` § Internal structure that a form filled from data a page reads again keeps the fields the user changed (S-0574).

First layer: the two components build on it.

## Done when

- [ ] `rebase` keeps edited leaves, takes untouched ones from the fresh read, drops vanished keys, adds new ones, and its tests pass.
- [ ] The design says the rule.

## Notes
