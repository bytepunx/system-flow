---
id: T-3520
type: task
nature: remediation
title: settle reads the story again and records left review for a story no longer in review
status: backlog
parent: S-0822
owner: alex
created: 2026-11-17T08:19:37Z
updated: 2026-11-17T08:19:37Z
transitions: []
stream: S-0822
tags: [flai, serve]
touches: [flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 9
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 18
          cache_read: 121485
          cache_write: 26992
          cost: 0.0466
---
# T-3520 settle reads the story again and records left review for a story no longer in review

## Work

Fix the cause I-0396 describes. In `flai/internal/closeout/closeout.go`, add an outcome `OutcomeLeft = "left review"` beside `OutcomeStayed` and `OutcomeSentBack`.

In `flai/internal/serve/resolvejob.go`, `settle` reads the story again before it records `stayed in review`, as `sendBack` does. When the story is no longer in review, it records `OutcomeLeft` with a reason that names:

- the state the story is in;
- who moved it there, and when, from its last transition;
- for a passed close-out, the command that moves it back to review, `flai move S-nnnn review`, since its close-out passed at its head.

`sendBack`'s path for a story that left review sets `OutcomeLeft` too. Today it leaves the outcome empty, so the review page words it as a plain close-out.

In `flai/internal/serve/resolvejob_test.go`, add a test modelled on `TestARebaseAndResolveLeavesAStoryNoLongerInReviewAsItIs`. It moves the story to in-progress while the close-out runs, lets the close-out pass, and reads the job record: outcome `left review`, and a reason with the state, the mover, and the command. Keep the passing test that stays in review. Add a `closeout_test.go` case that reads the new outcome through `Status`.

This task waits for no other: it is the fix itself.

## Done when

- A resolve job whose story left review while it ran records outcome `left review`, never `stayed in review`, whether its close-out passed or failed.
- The reason names the story's state, who moved it and when, and, after a pass, the command that moves it back to review.
- The new tests pass with `flai test` on the changed paths.

## Notes
