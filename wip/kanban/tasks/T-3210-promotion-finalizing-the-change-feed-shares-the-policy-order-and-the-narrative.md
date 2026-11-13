---
id: T-3210
type: task
nature: feature
title: Promotion, finalizing, the change feed, shares, the policy order, and the narrative's checks read a state's kind
status: backlog
parent: S-0739
owner: alex
created: 2026-11-12T16:06:46Z
updated: 2026-11-12T16:06:46Z
transitions: []
stream: S-0739
tags: [cli]
touches: [flai/internal/workitem/promote.go, flai/internal/workitem/promotable.go, flai/internal/workitem/finalizable.go, flai/internal/workitem/changes.go, flai/internal/workitem/share.go, flai/internal/workitem/policy.go, flai/internal/workitem/narrative.go, flai/internal/workitem/promote_test.go]
after: [T-3203]
---
# T-3210 Promotion, finalizing, the change feed, shares, the policy order, and the narrative's checks read a state's kind

## Work

Waits for T-3203: its tests use the shared test workflow that task adds. Replace the references to a state's name in `promote.go`, `promotable.go`, `finalizable.go`, `changes.go`, `share.go`, `policy.go`, and `narrative.go` with reads of kinds: a story is promoted from the first queue state into the one the pull takes from, a draft is finalized in a queue state, a share lasts while both stories are open, and the policy orders the ordered states. Leave `item.go`'s constants as the baseline's names.

## Done when

The seven files name no state and `promote_test.go` and the existing tests pass on both workflows.

## Notes
