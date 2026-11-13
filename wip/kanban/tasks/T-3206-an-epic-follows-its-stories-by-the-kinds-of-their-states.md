---
id: T-3206
type: task
nature: feature
title: An epic follows its stories by the kinds of their states
status: backlog
parent: S-0739
owner: alex
created: 2026-11-12T16:06:38Z
updated: 2026-11-12T16:06:38Z
transitions: []
stream: S-0739
tags: [cli]
touches: [flai/internal/workitem/follow.go, flai/internal/workitem/follow_test.go]
after: [T-3203]
---
# T-3206 An epic follows its stories by the kinds of their states

## Work

Waits for T-3203: its tests use the shared test workflow that task adds. Rewrite `follow.go`'s 18 references to a state's name as reads of kinds: an epic goes to the first queue state with its first story, to active work with the first started, to the wait state with the last open story, and to accepted when that one is accepted, as ADR-0076 says for the baseline. Test it on the baseline and on the shared fixture.

## Done when

`follow.go` names no state and `follow_test.go` passes on both workflows.

## Notes
