---
id: T-1643
type: task
nature: feature
title: The acceptance preview and flai accept take a passing batch's verification of the story's rebased head in place of the story's own verify
status: backlog
parent: S-0400
owner: alex
created: 2026-10-16T05:00:59Z
updated: 2026-10-16T05:00:59Z
transitions: []
stream: S-0400
tags: [cli]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go]
---
# T-1643 The acceptance preview and flai accept take a passing batch's verification of the story's rebased head in place of the story's own verify

## Work

- Let `flai accept` be told, by the queue, the batch it lands from. Its `--verified` check passes when the commit is the story's rebased head in that batch and the batch's verification passed. Every other blocker of the preview stays as it is.
- The acceptance merges the story's rebased head from the batch, so that what lands is what was verified.
- Test an acceptance from a passing batch, from a failing one, and from a batch whose commit is not the story's.
- First in its layer: it waits for no task.

## Done when

An acceptance from a passing batch goes through with the same commit, transitions, and notes as one by hand, the others are refused, and the tests pass.

## Notes
