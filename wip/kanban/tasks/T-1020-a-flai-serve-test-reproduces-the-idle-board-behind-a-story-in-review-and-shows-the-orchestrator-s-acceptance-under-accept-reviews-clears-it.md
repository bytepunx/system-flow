---
id: T-1020
type: task
nature: remediation
title: A flai serve test reproduces the idle board behind a story in review and shows the orchestrator's acceptance under accept_reviews clears it
status: backlog
parent: S-0296
owner: alex
created: 2026-10-06T11:50:06Z
updated: 2026-10-06T11:50:06Z
transitions: []
stream: S-0296
tags: [flai]
touches: [flai/internal/serve/review_wait_test.go, flai/internal/serve/orchestrate.go]
---
# T-1020 A flai serve test reproduces the idle board behind a story in review and shows the orchestrator's acceptance under accept_reviews clears it

## Work

Write `flai/internal/serve/review_wait_test.go`, beside `hold_test.go` and `orchestrate_test.go`, reproducing I-0088 in a scratch project: one story in review with every criterion ticked, and a ready story whose touches overlap it, so it is held.

- With `orchestration.permissions.accept_reviews` off, the launcher starts no story agent and the ready story stays held: this is the cause the issue describes.
- With `orchestrate` on and `accept_reviews` on, a stand-in orchestrator accepts the story in review the way S-0221 built it (`flai accept --by orchestrator --verified`), and the launcher then starts the held story's agent without a person.

Use the stand-in harness the existing serve tests use. If the launcher only notices the review on its one-minute look, so that the test has to wait for it, make it look at the orchestrator when a story moves to review, in `orchestrate.go`. The task waits for nothing in this story; S-0221, which the story waits for, provides `flai accept --by orchestrator`.

## Done when

- The test fails with `accept_reviews` off, the board idle behind the story in review, and passes with it on, the held story started.
- `scripts/flai-test.sh` passes, the race detector included.

## Notes
