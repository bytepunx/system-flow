---
id: T-1644
type: task
nature: feature
title: A red batch is bisected to the first story with which it fails, with each step kept in the queue's state
status: backlog
parent: S-0400
owner: alex
created: 2026-10-16T05:01:03Z
updated: 2026-10-16T05:01:03Z
transitions: []
stream: S-0400
tags: [cli]
touches: [flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go]
---
# T-1644 A red batch is bisected to the first story with which it fails, with each step kept in the queue's state

## Work

- `bisect.go`: given a failing batch of stories in order, and a function that verifies main with a prefix of them, find the first story with which the prefix fails, by halving. Return the stories before it, the story, and the stories after it.
- Record each step, the prefix and its result, in the state's batch, so the view can show it.
- A batch of one is its own answer. A prefix that passes when the whole batch failed points to the later half.
- Test with a fake verify: a failure at each position, and a failure that needs two stories together.
- First in its layer: it waits for no task.

## Done when

Bisection finds the failing story in every tested case, records its steps, and the tests pass.

## Notes
