---
id: T-3205
type: task
nature: feature
title: flai check's item rules and its scope read a state's kind instead of its name
status: done
parent: S-0741
owner: alex
created: 2026-11-12T16:06:30Z
updated: 2026-11-15T20:45:50Z
transitions:
  - to: ready
    at: 2026-11-15T20:04:28Z
    by: agent-S-0741
  - to: in-progress
    at: 2026-11-15T20:04:28Z
    by: agent-S-0741
  - to: done
    at: 2026-11-15T20:45:50Z
    by: agent-S-0741
stream: S-0741
tags: [cli]
touches: [flai/internal/check/check_items.go, flai/internal/check/items_test.go, flai/internal/check/scope.go, flai/internal/check/concurrent_test.go]
after: [T-3204]
usage:
  source: log
  seconds: 2482
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 28715
      cache_read: 5193175
      cache_write: 322521
      cost: 3.468
---
# T-3205 flai check's item rules and its scope read a state's kind instead of its name

## Work

Replace the 21 references to a state's name in `check_items.go` and the one in `scope.go` with reads of the state's kind from the repository's workflow, so a story in a wait state is checked as one in `review` is today and one in an active-work state as one in `in-progress`. Add tests that run the item rules on a workflow with a state renamed and one added between `in-progress` and `review`.

## Done when

`check_items.go` and `scope.go` name no state, and the item rules give the same findings on both workflows where the kinds agree.

## Notes
