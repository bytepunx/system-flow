---
id: T-3248
type: task
nature: feature
title: metrics.md and the user guide say how a timer is defined and where flai stats reports it
status: done
parent: S-0752
owner: alex
created: 2026-11-12T16:09:27Z
updated: 2026-11-17T08:47:19Z
transitions:
  - to: ready
    at: 2026-11-17T08:36:53Z
    by: agent-S-0752
  - to: in-progress
    at: 2026-11-17T08:36:53Z
    by: agent-S-0752
  - to: done
    at: 2026-11-17T08:47:19Z
    by: agent-S-0752
stream: S-0752
tags: [cli]
touches: [design/system/metrics.md, docs/users/flai.md, docs/operators/runbooks/change-the-board.md]
after: [T-3243]
usage:
  source: log
  seconds: 626
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 12145
      cache_read: 3102298
      cache_write: 107560
      cost: 1.5703
---
# T-3248 metrics.md and the user guide say how a timer is defined and where flai stats reports it

## Work

Waits for T-3243: it documents the section that task names. Add to `design/system/metrics.md` how a workflow's own timer is computed and where `flai stats --json` gives it, and to `docs/users/flai.md` how an operator adds one and reads it. Regenerate `docs/users/flai-reference.md` with `make flai-reference` if the command's help changed.

## Done when

Both documents say how a timer is defined and read, and `TestReferenceIsCurrent` passes.

## Notes
