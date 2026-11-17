---
id: T-3442
type: task
nature: feature
title: The dashboard's Publish and release_publish answer a publish a failed suite refused with the suite, the tier, and the log
status: backlog
parent: S-0811
owner: alex
created: 2026-11-16T23:48:09Z
updated: 2026-11-16T23:48:09Z
transitions: []
stream: S-0811
tags: [flai]
touches: [flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_test.go]
after: [T-3441]
---
# T-3442 The dashboard's Publish and release_publish answer a publish a failed suite refused with the suite, the tier, and the log

## Work

`publish.run` and the orchestrator's `release_publish` run `flai release --pending` through `Publisher` in `flai/internal/hostapi/writes_outcome.go`, and `publishExits` maps only exit 3 to a payload. Add T-3441's new exit to `publishExits`, so that a publish a failed suite refused answers with flai's message, which names the suite, the tier, and the log, rather than as an unexplained failure. Journal it with its own outcome, apart from a conflict.

It waits for T-3441, which defines the exit and its message.

## Done when

- A test in `flai/internal/hostapi/writes_test.go` shows that a publish whose run exits with T-3441's code answers an error that carries flai's message, and that the journal records it apart from a conflict and from a push the host action held.
- `flai test` passes on the changed paths.

## Notes
