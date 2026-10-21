---
id: T-1781
type: task
nature: improvement
title: A formula step that publishes is refused while the push host action is off, and a publish refused for a moved remote stops the instance with flai's conflict message
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:33Z
updated: 2026-10-21T13:10:50Z
transitions: []
stream: S-0423
tags: [cli]
touches: [flai/cmd/release.go, flai/cmd/release_push_action_test.go]
after: [T-1780]
---
# T-1781 A formula step that publishes is refused while the push host action is off, and a publish refused for a moved remote stops the instance with flai's conflict message

## Work

In the step runner:

- Refuse a step whose command is `flai release --pending`, `flai push`, or `git push` while the project's `push` host action is off. Use the refusal `release_publish` gives, naming `flai serve enable push`, and send the step down the failure path.
- When `flai release --pending` exits 3 with a `conflict:` message, put the message unchanged on the task and in the thread, as `release_publish` returns it.
- Journal the publish under `push` with the method `formula.publish` and the formula's name, so a release a formula cut is told apart.

Waits for the formula file, whose steps it runs in its test.

## Done when

`go test ./internal/serve/...` passes, with a test for each case.

## Notes

Replanned 2026-10-21 by planner-S-0423 on TH-0471's answer 5: there is no serve-side step runner, so `flai/internal/serve/formula.go` is gone from the touches. The refusal moves into `flai release --pending` itself, for an agent `flai serve` started while the `push` host action is off, read from the host's config as `cmd/verify.go`'s `hostActionOn` reads it; the `conflict:` exit 3 is already flai's, and S-0422's failure path carries it to the task and a thread. The rewrite of the Work to say so is asked on the plan thread. It waits for T-1780 because both change `flai/cmd/release.go`.
