---
id: T-3590
type: task
nature: improvement
title: "A test reproduces I-0386: a review story whose close-out stopped after main moved is sent back by the resolve job the orchestrator asked for"
status: backlog
parent: S-0845
owner: alex
created: 2026-11-17T12:34:18Z
updated: 2026-11-17T12:34:18Z
transitions: []
stream: S-0845
tags: [flai, test]
touches: [flai/internal/serve/resolvejob_test.go]
after: [T-3588]
---
# T-3590 A test reproduces I-0386: a review story whose close-out stopped after main moved is sent back by the resolve job the orchestrator asked for

## Work

Waits for T-3588, whose orchestrator request this test makes. In `flai/internal/serve/resolvejob_test.go`, beside `TestARebaseAndResolveThatRebasesCleanlyAndPassesLeavesTheStoryInReview` and `...CloseOutStopsSendsTheStoryBackNamingTheStep`, reproduce I-0386's S-0742 case: a story in review whose last close-out stopped at a step, main moved past its base, and a resolve request made by the orchestrator. Assert the job rebases, its close-out stops at that step, the story is moved to in-progress as `flai`, and the `Acceptance sent back:` thread names the step, with no entry by the operator. Add the clean case too: the same request, a passing close-out, and the story left in review.

## Done when

- The two tests pass, and the first fails on the code before T-3588.
- `flai test flai/internal/serve/resolvejob_test.go` passes.

## Notes
