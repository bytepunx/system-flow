---
id: T-2834
type: task
nature: remediation
title: The request answered tests poll the log for the event, and a held-back log line reproduces the race
status: backlog
parent: S-0649
owner: alex
created: 2026-11-10T06:58:34Z
updated: 2026-11-10T06:58:34Z
transitions: []
stream: S-0649
tags: [flai, test]
touches: [flai/internal/channel/channel_test.go]
---
# T-2834 The request answered tests poll the log for the event, and a held-back log line reproduces the race

## Work

First layer: it waits for no task.

- Reproduce first. Give `TestEachRequestIsTimedBeneathTheTransport` a log destination that holds back each line holding `"request answered"`, for 300 ms or so, before it writes it to the `lockedBuffer`. With the test's one-shot read, the test then fails with "no request answered event" on every run and without load. Record that run in the task's Notes.
- Fix. Add one helper beside `lockedBuffer`, such as `answeredEvent(t, logs, within)`, that polls the log for the "request answered" event until a deadline of a few seconds and fails the test with the log when none comes. Use it in `TestEachRequestIsTimedBeneathTheTransport` and in `TestTheSendTimeCountsTheWaitForTheWriteLock`, in place of their own loops.
- Keep the held-back line in the first test, shorter than the helper's deadline, so the test covers an event logged after the reply on every run.
- Change no product code: the event is meant to come after the write (S-0611). If the fix proves to need it, write down why in the narrative's Decisions first.

## Done when

- With the held-back line and before the fix, the test fails with "no request answered event"; the run is in the task's Notes.
- After the fix, `go test ./internal/channel -run 'TestEachRequestIsTimedBeneathTheTransport|TestTheSendTimeCountsTheWaitForTheWriteLock' -count 20` passes in `flai/`, with the held-back line in place.
- `flai test flai/internal/channel/channel_test.go` passes.

## Notes
