---
id: T-3205
type: task
nature: feature
title: flai check's item rules and its scope read a state's kind instead of its name
status: backlog
parent: S-0741
owner: alex
created: 2026-11-12T16:06:30Z
updated: 2026-11-12T16:06:30Z
transitions: []
stream: S-0741
tags: [cli]
touches: [flai/internal/check/check_items.go, flai/internal/check/items_test.go, flai/internal/check/scope.go]
---
# T-3205 flai check's item rules and its scope read a state's kind instead of its name

## Work

Replace the 21 references to a state's name in `check_items.go` and the one in `scope.go` with reads of the state's kind from the repository's workflow, so a story in a wait state is checked as one in `review` is today and one in an active-work state as one in `in-progress`. Add tests that run the item rules on a workflow with a state renamed and one added between `in-progress` and `review`.

## Done when

`check_items.go` and `scope.go` name no state, and the item rules give the same findings on both workflows where the kinds agree.

## Notes
