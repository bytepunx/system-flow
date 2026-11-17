---
id: T-3474
type: task
nature: feature
title: flai check reports a suite the workflow file names that no tier has, and a non-manual tier's suite nothing names
status: in-progress
parent: S-0806
owner: alex
created: 2026-11-16T23:50:47Z
updated: 2026-11-17T08:48:40Z
transitions:
  - to: ready
    at: 2026-11-17T08:48:39Z
    by: agent-S-0806
  - to: in-progress
    at: 2026-11-17T08:48:40Z
    by: agent-S-0806
stream: S-0806
tags: [flai]
touches: [flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go]
after: [T-3469, T-3471]
---
# T-3474 flai check reports a suite the workflow file names that no tier has, and a non-manual tier's suite nothing names

## Work

In `flai/internal/check/check_workflow.go`, compare the suites the workflow file names, from T-3469's `suites.go`, with the suites of the manifest's tiers, each tier's `suite` as S-0805 gives it, its own name when it gives none (ADR-0312 §9):

- a suite the file names that no tier has is an error, on the line that names it;
- a suite of a tier that is not `manual`, which no transition, state, or `verify_at` point names, is a warning naming the tier.

`all` names every suite, so a file that names `all` anywhere raises no warning and `all` is never an unknown suite. Give the two findings codes beside the other `workflow.*` codes, and say in each what to do.

It waits for T-3469, whose suite list it reads, and for T-3471: the check's baseline case, the baseline with no findings, needs the baseline to name `all` first, or every tier of a project warns.

## Done when

- `workflow_test.go` in `flai/internal/check` covers an unknown suite, an unnamed suite of a tier, a `manual` tier that is not named, a file that names `all`, and the baseline with no findings.
- `flai check --strict` on this repository reports neither finding, and `flai test` on the changed paths passes.

## Notes
