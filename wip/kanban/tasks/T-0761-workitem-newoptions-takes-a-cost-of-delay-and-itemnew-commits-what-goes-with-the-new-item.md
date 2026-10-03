---
id: T-0761
type: task
nature: improvement
title: workitem.NewOptions takes a cost of delay and itemnew commits what goes with the new item
status: done
parent: S-0203
owner: alex
created: 2026-10-03T18:07:41Z
updated: 2026-10-03T18:11:26Z
transitions:
  - to: ready
    at: 2026-10-03T18:08:34Z
    by: agent-S-0203
  - to: in-progress
    at: 2026-10-03T18:08:34Z
    by: agent-S-0203
  - to: done
    at: 2026-10-03T18:11:26Z
    by: agent-S-0203
stream: S-0203
tags: []
touches: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/itemnew]
usage:
  source: log
  seconds: 172
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 60
      output: 22724
      cache_read: 2462105
      cache_write: 88652
      cost: 1.4911
---

# T-0761 workitem.NewOptions takes a cost of delay and itemnew commits what goes with the new item

## Work

`workitem.NewOptions` gains `CostOfDelay *CostOfDelay`, set on the new item (epics and stories only, as `Validate` allows), so that it is in place when `flai check` runs over the new item.

`itemnew.Options` gains a hook run once the check keeps the item and before the autocommit: it writes what goes with the item and returns the paths to commit with it. An error from it removes the item, as a refusal does, and is returned. `flai issue story` uses it to link the story from the issue in the same commit. Waits for nothing: the first layer, beside T-0760.

## Done when

- A story created with a cost of delay carries it; a task given one is refused.
- The hook's paths are committed with the item; its error leaves nothing.
- Tests for both pass with `go test -race -short ./internal/workitem/ ./internal/itemnew/`.

## Notes
