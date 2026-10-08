---
id: T-1428
type: task
nature: improvement
title: "Tests reproduce I-0092's instances not yet covered: an acceptance whose rebase meets main's bump of the same issue, and a sync whose bump meets main's close of it"
status: in-progress
parent: S-0297
owner: alex
created: 2026-10-08T09:44:05Z
updated: 2026-10-08T09:44:41Z
transitions:
  - to: ready
    at: 2026-10-08T09:44:41Z
    by: agent-S-0297
  - to: in-progress
    at: 2026-10-08T09:44:41Z
    by: agent-S-0297
stream: S-0297
tags: []
touches: [flai/cmd/stream_sync_test.go]
---

# T-1428 Tests reproduce I-0092's instances not yet covered: an acceptance whose rebase meets main's bump of the same issue, and a sync whose bump meets main's close of it

## Work

S-0326 (ADR-0126) made the rebase in `flai stream sync`, `flai task done`, and `flai accept` merge an issue file both sides changed. `flai/cmd/stream_sync_test.go` already reproduces two sync bumps (`TestSyncMergesAnIssueBothSidesBumped`) and the trial merge (`TestSyncTrialMergeLeavesOutIssueFiles`). Add, beside them and with their helpers (`bumpStories`, `syncIssuesJSON`, `issuesIn`, `hasInstances`):

1. An acceptance of a branch that bumped an issue after main's bump of it reached main: `flai accept` merges the issue, count and both instances, and the summary, and lands on main (I-0092's second instance, S-0254 and S-0228).
2. A sync of a branch that bumped an issue after main closed it: the sync merges it into a closed issue holding the bump, with the summary written again (I-0092's third instance, S-0275).

Waits for nothing.

## Done when

- Both tests pass, and each fails against `main`'s rebase without the issue merge (shown by reasoning on the stop it would meet, or by running against a scratch copy).

## Notes
