---
id: T-1548
type: task
nature: feature
title: flai/internal/designsystem reads the design system from an open story's worktree and lists the open stories whose branches change it
status: done
parent: S-0384
owner: alex
created: 2026-10-13T17:47:55Z
updated: 2026-10-15T03:22:44Z
transitions:
  - to: ready
    at: 2026-10-15T03:17:11Z
    by: system-flow
  - to: in-progress
    at: 2026-10-15T03:17:11Z
    by: system-flow
  - to: done
    at: 2026-10-15T03:22:44Z
    by: system-flow
stream: S-0384
tags: [cli]
touches: [flai/internal/designsystem/proposals.go, flai/internal/designsystem/proposals_test.go]
usage:
  source: log
  seconds: 333
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 11182
      cache_read: 1604685
      cache_write: 48698
      cost: 0.8252
---
# T-1548 flai/internal/designsystem reads the design system from an open story's worktree and lists the open stories whose branches change it

## Work

In `flai/internal/designsystem/proposals.go`, `ShowStory` gives `Show` for the document in a story's worktree, refusing a story that is not in progress or in review, or has no worktree, and reusing how S-0378 finds a story's worktree for `flai doc show --story`. `Proposals` lists the open stories whose branch changes the document against its merge base with the main branch, each with its ID, title, and state. Tests use a temporary repository with a story branch that changes it and one that does not. It waits for nothing.

## Done when

- `go test ./internal/designsystem/` passes with both functions' cases.

## Notes

- Layer 1, beside T-1549 and T-1550, which are the dashboard's and share no file with it.
