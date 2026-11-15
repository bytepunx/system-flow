---
id: T-3369
type: task
nature: feature
title: The workflow names the states the acceptance leaves from, the send-back out of a state, and the closed kinds, and flai serve opens its repositories with the workflow its options give
status: done
parent: S-0743
owner: alex
created: 2026-11-15T00:54:10Z
updated: 2026-11-15T01:07:08Z
transitions:
  - to: ready
    at: 2026-11-15T00:54:18Z
    by: agent-S-0743
  - to: in-progress
    at: 2026-11-15T00:54:18Z
    by: agent-S-0743
  - to: done
    at: 2026-11-15T01:07:08Z
    by: agent-S-0743
stream: S-0743
tags: []
touches: [flai/internal/workflow/review.go, flai/internal/workflow/review_test.go, flai/internal/serve/serve.go, flai/internal/serve/flow.go, flai/internal/serve/flow_test.go]
usage:
  source: log
  seconds: 770
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 7772
      cache_read: 3038653
      cache_write: 81724
      cost: 1.2279
---
# T-3369 The workflow names the states the acceptance leaves from, the send-back out of a state, and the closed kinds, and flai serve opens its repositories with the workflow its options give

## Work

- Add `flai/internal/workflow/review.go`, methods on `Workflow` that read a role from kinds and effects, never a name: `AcceptFrom(typ)` (the states the transitions of typ with the accept effect leave from, canonical, in board order), `AwaitsAcceptance(typ, state)`, `SendBack(typ, from)` (the transition of typ out of from with the send-back effect), `Gives(typ, from, to, effect)` (whether the move from one state to the other gives the effect), `Closed(state)` (kind accepted or cancelled), and `Underway(state)` (kind active or waiting). Test them in `review_test.go` against the baseline and against a workflow with `review` renamed and a state added between `in-progress` and it.
- Add `Options.Workflow *workflow.Workflow` to `flai/internal/serve/serve.go`, nil for the baseline, and `flai/internal/serve/flow.go` with `openRepo(o, root)`, which opens the repository as `workitem.Open` does and gives it `o.Workflow`, and `o.flow()`. `flow_test.go` holds the serve tests' fixture: the other workflow's text and a helper that parses it into options.
- Waits for nothing. T-3218, T-3219, and T-3220 wait for it: each reads these methods and opens its repositories through `openRepo`.

## Done when

- The methods answer the baseline's `review`, `in-progress`, `done`, and `cancelled` roles as today, and the other workflow's renamed and added states by their kinds and effects, in `review_test.go`.
- `openRepo` gives the repository the options' workflow, and the baseline when it is nil, in `flow_test.go`.

## Notes
