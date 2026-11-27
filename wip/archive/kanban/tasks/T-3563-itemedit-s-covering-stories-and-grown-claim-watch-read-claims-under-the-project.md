---
id: T-3563
type: task
nature: remediation
title: itemedit's covering stories and grown-claim watch read claims under the project's workflow
status: done
parent: S-0843
owner: alex
created: 2026-11-17T12:29:08Z
updated: 2026-11-18T10:21:28Z
transitions:
  - to: ready
    at: 2026-11-18T10:04:11Z
    by: agent-S-0843
  - to: in-progress
    at: 2026-11-18T10:04:11Z
    by: agent-S-0843
  - to: done
    at: 2026-11-18T10:21:28Z
    by: agent-S-0843
stream: S-0843
tags: [flai, workflow, claims]
touches: [flai/internal/itemedit/covers.go, flai/internal/itemedit/covers_test.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go]
usage:
  source: log
  seconds: 1037
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 4076
      cache_read: 1189762
      cache_write: 53872
      cost: 0.6416
---
# T-3563 itemedit's covering stories and grown-claim watch read claims under the project's workflow

## Work

Two callers in `flai/internal/itemedit` build their Holds with the baseline wrapper `workitem.NewHolds` although a workflow is at hand, so a story's claim keeps the touches of a task closed under a renamed state (I-0379):

- `CoveringIn(flow, ...)` in `covers.go` takes `flow` and judges which stories are under way by it, but reads each claim from `workitem.NewHolds(items, projects)`. Use `workitem.NewHoldsIn(flow, items, projects)`.
- The claim watch in `claim.go` reads the story's claim before a write from `workitem.NewHolds(items, repo.Manifest.Projects)`. Use `workitem.NewHoldsIn(repo.Flow(), items, repo.Manifest.Projects)`, so the claim before and after a write are read under one workflow.

Waits for nothing: `NewHoldsIn` exists on main already, and no other task of this story changes these files.

## Done when

- Neither file calls `workitem.NewHolds`.
- A test in `covers_test.go` shows `CoveringIn` under a workflow whose closed task state is not `done` leaves out a story whose only covering touch is on a task in that state.
- A test in `claim_test.go` shows the grown-claim message is not sent for a path only a task closed under such a workflow claimed.
- `flai test` on the four files passes.

## Notes
