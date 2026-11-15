---
id: T-3378
type: task
nature: remediation
title: Make verify_behind_main send back only the tiers the overlapping main paths select
status: done
parent: S-0792
owner: alex
created: 2026-11-15T07:43:31Z
updated: 2026-11-15T17:05:07Z
transitions:
  - to: ready
    at: 2026-11-15T16:56:14Z
    by: agent-S-0792
  - to: in-progress
    at: 2026-11-15T16:56:14Z
    by: agent-S-0792
  - to: done
    at: 2026-11-15T17:05:07Z
    by: agent-S-0792
stream: S-0792
tags: [flai, acceptance]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept_stale_verify_test.go, design/adrs/0303-an-acceptance-s-verify-behind-main-sends-back-only-the-tiers-selected-by-the.md]
after: [T-3377]
usage:
  source: log
  seconds: 533
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 22
      output: 6766
      cache_read: 1634260
      cache_write: 73811
      cost: 0.9905
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 21939
          cache_write: 603
          cost: 0.0068
---
# T-3378 Make verify_behind_main send back only the tiers the overlapping main paths select

## Work

In `verifyBehindMain` in `flai/internal/preview/accept.go`, filter the paths `verify.MainPathsSince` returns through T-3377's overlap with the report's `Paths` before `verify.StaleTiers` selects tiers from them. Main paths that overlap nothing send back nothing. An error from the overlap is the blocker, naming the error, as an unreadable report is. Make `VerifyBehindMainBlocker` count and name the overlapping paths, not every path main brought in.

In `flai/internal/preview/accept_test.go`, keep `TestAnAcceptanceWaitsForAVerifyOverMainsNewPaths` blocking the I-0229 case, and add the case the operator reported: main changes a flai package the story's packages do not import, so no blocker though `integration` and `smoke` select it. In `flai/cmd/accept_stale_verify_test.go`, check that `TestAnAcceptanceIsRefusedWhileMainOutranTheStorysVerify` still holds under the rule (its main change is a sibling file in the story's package), and add an acceptance that goes through after a main change outside the story's change set.

It waits for T-3377, whose function it calls.

## Done when

- An acceptance whose story's verify main outran only on paths outside its change set has no `verify_behind_main`, from the operator, the orchestrator, and the dry run.
- The I-0229 case and a shared file still block, naming the overlapping paths.
- `flai test` on the changed files passes.

## Notes
