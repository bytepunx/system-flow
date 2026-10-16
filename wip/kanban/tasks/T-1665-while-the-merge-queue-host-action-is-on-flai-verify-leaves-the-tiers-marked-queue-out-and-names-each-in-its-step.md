---
id: T-1665
type: task
nature: improvement
title: While the merge-queue host action is on, flai verify leaves the tiers marked queue out and names each in its step
status: backlog
parent: S-0405
owner: alex
created: 2026-10-16T05:02:32Z
updated: 2026-10-16T05:02:32Z
transitions: []
stream: S-0405
tags: [cli]
touches: [flai/internal/verify/select.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go]
after: [T-1654]
---
# T-1665 While the merge-queue host action is on, flai verify leaves the tiers marked queue out and names each in its step

## Work

- In `select.go` and `story.go`, when the host action `merge-queue` is on for the project, leave each tier with `queue` out of a story's verify. Name each in the step as left to the queue, and pass without it. With the action off, run every tier as today.
- `--all`, and `flai test --all`, still run every tier.
- Keep ADR-0156's reuse working: a tier left out is neither reused nor sent back.
- Test with the action on and off, and with `--all`.
- Waits for T-1654, which adds the key.

## Done when

`flai verify` leaves the marked tiers out only while the queue runs, says so, and the tests pass.

## Notes
