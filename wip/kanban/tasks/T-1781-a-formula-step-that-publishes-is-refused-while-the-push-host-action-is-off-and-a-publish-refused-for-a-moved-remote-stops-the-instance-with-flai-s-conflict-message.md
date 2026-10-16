---
id: T-1781
type: task
nature: improvement
title: A formula step that publishes is refused while the push host action is off, and a publish refused for a moved remote stops the instance with flai's conflict message
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:33Z
updated: 2026-10-16T05:48:33Z
transitions: []
stream: S-0423
tags: [cli]
touches: [flai/internal/serve/formula.go, flai/internal/serve/formula_publish_test.go]
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
