---
id: T-3206
type: task
nature: feature
title: An epic follows its stories by the kinds of their states
status: done
parent: S-0739
owner: alex
created: 2026-11-12T16:06:38Z
updated: 2026-11-15T00:54:49Z
transitions:
  - to: ready
    at: 2026-11-14T23:56:42Z
    by: agent-S-0739
  - to: in-progress
    at: 2026-11-15T00:26:04Z
    by: agent-S-0739
  - to: done
    at: 2026-11-15T00:54:49Z
    by: agent-S-0739
stream: S-0739
tags: [cli]
touches: [flai/internal/workitem/follow.go, flai/internal/workitem/follow_test.go, flai/internal/workitem/boardview.go]
after: [T-3203]
usage:
  source: log
  seconds: 1725
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 12541
      cache_read: 1591923
      cache_write: 84753
      cost: 1.0676
---
# T-3206 An epic follows its stories by the kinds of their states

## Work

Waits for T-3203: its tests use the shared test workflow that task adds. Rewrite `follow.go`'s 18 references to a state's name as reads of kinds: an epic goes to the first queue state with its first story, to active work with the first started, to the wait state with the last open story, and to accepted when that one is accepted, as ADR-0076 says for the baseline. Test it on the baseline and on the shared fixture.

## Done when

`follow.go` names no state and `follow_test.go` passes on both workflows.

## Notes
