---
id: T-1820
type: task
nature: feature
title: The orchestrator's pack adds the answer and recommendation entries ranked against the threads awaiting the operator
status: done
parent: S-0430
owner: alex
created: 2026-10-16T06:04:49Z
updated: 2026-10-23T19:30:00Z
transitions:
  - to: ready
    at: 2026-10-23T19:24:51Z
    by: agent-S-0430
  - to: in-progress
    at: 2026-10-23T19:24:52Z
    by: agent-S-0430
  - to: done
    at: 2026-10-23T19:29:59Z
    by: agent-S-0430
stream: S-0430
tags: [cli]
touches: [flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/inbox/inbox.go, flai/internal/threads/threads.go]
after: [T-1819]
usage:
  source: log
  seconds: 307
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 18524
      cache_read: 2898893
      cache_write: 98360
      cost: 1.5382
---
# T-1820 The orchestrator's pack adds the answer and recommendation entries ranked against the threads awaiting the operator

## Work

In `ForStrategic` (`flai/internal/context/strategic.go`), for role `orchestrate`, read the threads awaiting the operator (those `inbox` would list for the operator: open, not opened by the operator, without a pending recommendation from the orchestrator), build a query from their titles and first entries, and add the `answer` and `recommendation` entries `digest.go`'s ranking puts first, within the budget, with the reason naming the thread they ranked for.

Waits for T-1819, whose selection and ranking it reuses.

## Done when

`strategic_test.go` shows an entry ranked first for a thread that asks its question again, a decision entry left out of the orchestrator's pack, and an orchestrator's pack with no threads awaiting that adds no entry; `flai test flai/internal/context` passes.

## Notes

The planner's and analyzer's packs are unchanged by this story.
