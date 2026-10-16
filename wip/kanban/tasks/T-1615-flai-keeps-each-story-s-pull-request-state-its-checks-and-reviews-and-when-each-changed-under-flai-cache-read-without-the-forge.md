---
id: T-1615
type: task
nature: feature
title: flai keeps each story's pull request state, its checks and reviews and when each changed, under .flai-cache, read without the forge
status: backlog
parent: S-0394
owner: alex
created: 2026-10-16T04:36:38Z
updated: 2026-10-16T04:36:38Z
transitions: []
stream: S-0394
tags: [cli]
touches: [flai/internal/forge/state.go, flai/internal/forge/state_test.go]
---
# T-1615 flai keeps each story's pull request state, its checks and reviews and when each changed, under .flai-cache, read without the forge

## Work

- In `flai/internal/forge/state.go`, keep one file per project under `.flai-cache` with each story's pull request:
  - its number, URL, and head commit
  - each required check's name, state, and log URL
  - each review's author, state, and comment URLs
- Record each change of a check's or a review's state as an event with its time, appended, so that a reader can tell when the story waited on checks and when on reviewers. S-0396 reads these events.
- Write atomically with `flai/internal/atomicfile`, and read without any token or network.
- Test writing, reading, a change appended as an event, and a missing file read as no state.
- First in its layer: it waits for no task.

## Done when

The state and its events round-trip through the file, and the tests pass.

## Notes
