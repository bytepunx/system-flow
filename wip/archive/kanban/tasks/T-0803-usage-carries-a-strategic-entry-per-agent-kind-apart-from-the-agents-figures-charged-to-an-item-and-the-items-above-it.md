---
id: T-0803
type: task
nature: improvement
title: Usage carries a strategic entry per agent kind, apart from the agents' figures, charged to an item and the items above it
status: done
parent: S-0225
owner: alex
created: 2026-10-04T04:07:16Z
updated: 2026-10-04T04:12:43Z
transitions:
  - to: ready
    at: 2026-10-04T04:07:53Z
    by: agent-S-0225
  - to: in-progress
    at: 2026-10-04T04:07:54Z
    by: agent-S-0225
  - to: done
    at: 2026-10-04T04:12:43Z
    by: agent-S-0225
stream: S-0225
tags: []
touches: [flai/internal/usage, flai/internal/workitem, design/system/work-hierarchy.md]
usage:
  source: log
  seconds: 289
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 80
      output: 561
      cache_read: 4273989
      cache_write: 103122
      cost: 1.7854
---
# T-0803 Usage carries a strategic entry per agent kind, apart from the agents' figures, charged to an item and the items above it

## Work

Give an item's `usage` a `strategic` list, one entry per strategic agent kind (`planner`, `orchestrator`, `analyzer`, in that order), each with `kind`, `seconds`, `estimated: true`, and `models` in the shape agent models have (ADR of this story). It is kept apart from the agents' figures: `Tokens`, `Cost`, `Seconds`, `Models`, `Add`, and `Sum` stay the agents' alone, and `Empty` keeps meaning that the agents spent nothing; `Same` compares `strategic` too, and an item whose usage carries only `strategic` is written, read, and validated.

`workitem` gains a way to charge a kind's usage to an item and to each item above it, up to its epic, giving an item with no usage one with `source: sum`, no seconds, and no models; `RollUp`, which sums the agents' figures from the children, keeps a parent's `strategic` as it is. `front-matter-fields.txt` lists the keys of `usage` and of a strategic entry, and its test keeps them equal to the code. `design/system/work-hierarchy.md` shows the block and says how it is charged.

Waits for nothing: the serve, stats, and dashboard tasks all build on these types.

## Done when

- [ ] `usage.Usage` carries `Strategic`, written by `usageBlock` after `models`, parsed, and refused by `usageErrors` for an unknown or repeated kind, a negative count or cost, or `estimated` not true
- [ ] Charging adds to the item's entry for the kind and to the same entry on every item above it, and tests pin it on a story under an epic and on an epic
- [ ] `RollUp` and `Sum` leave `strategic` out of the agents' sums and keep a parent's own, with tests
- [ ] `front-matter-fields.txt` and its test cover the usage keys; `go test ./internal/usage/... ./internal/workitem/...` passes
- [ ] `design/system/work-hierarchy.md` describes the `strategic` entry

## Notes
