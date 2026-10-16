---
id: T-1777
type: task
nature: feature
title: flai accept accepts a formula instance with no branch by moving it to done and archiving it, and the orchestrator's acceptance takes its steps' exits as evidence
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:47:54Z
updated: 2026-10-16T05:47:54Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go]
after: [T-1774]
---
# T-1777 flai accept accepts a formula instance with no branch by moving it to done and archiving it, and the orchestrator's acceptance takes its steps' exits as evidence

## Work

`flai accept` merges a story's branch. For a story with `formula:` and no branch:

- skip the merge, the sync checks, and the tiers;
- keep the move to done, the archive, the commit, and the release step's nature rules.

The acceptance preview names no branch blocker for it. Under `--by orchestrator`, the evidence is each task's step and exit, and the four orchestrator blockers hold as they do for any story. Test both, and that a story with no branch and no `formula` is still refused.

Waits for the step runner, whose instances it accepts.

## Done when

`go test ./cmd/... ./internal/preview/...` passes. The existing accept tests are unchanged.

## Notes
