---
id: T-1615
type: task
nature: feature
title: flai keeps each story's pull request state, its checks and reviews and when each changed, under .flai-cache, read without the forge
status: done
parent: S-0394
owner: alex
created: 2026-10-16T04:36:38Z
updated: 2026-10-17T16:14:35Z
transitions:
  - to: ready
    at: 2026-10-17T16:12:01Z
    by: agent-S-0394
  - to: in-progress
    at: 2026-10-17T16:12:01Z
    by: agent-S-0394
  - to: done
    at: 2026-10-17T16:14:35Z
    by: agent-S-0394
stream: S-0394
tags: [cli]
touches: [flai/internal/forge/state.go, flai/internal/forge/state_test.go]
usage:
  source: log
  seconds: 154
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 3841
      cache_read: 549922
      cache_write: 16612
      cost: 0.2809
    - model: claude-sonnet-5-5
      input: 38
      output: 14047
      cache_read: 1303893
      cache_write: 99792
      cost: 0.5204
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
