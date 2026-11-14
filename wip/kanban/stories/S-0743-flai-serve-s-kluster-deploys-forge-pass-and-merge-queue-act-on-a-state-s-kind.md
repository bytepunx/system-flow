---
id: S-0743
type: story
nature: feature
title: flai serve's kluster deploys, forge pass, and merge queue act on a state's kind and a transition's effect instead of review and done
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:52Z
updated: 2026-11-14T17:31:41Z
transitions: []
tags: [cli]
topics: [cli, kluster]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/forge.go, flai/internal/serve/forge_test.go, flai/internal/serve/forgemerge.go, flai/internal/forge/state.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 108314
          cache_write: 18020
          cost: 0.0361
draft: true
forecast:
  duration: 45m
  delivery: 2026-11-14T20:35:00Z
  basis: "Its own forecast of 45m; 8th in the pull order with an in-progress limit of 3, behind S-0781, S-0241, S-0736, S-0737, S-0738, S-0739, S-0741 and S-0742."
  by: flai
  at: 2026-11-14T17:31:41Z
---
# S-0743 flai serve's kluster deploys, forge pass, and merge queue act on a state's kind and a transition's effect instead of review and done

## Goal

What flai serve does around review follows the workflow. A story is deployed into its kluster namespace, pushed to the forge with its pull request, and taken by the merge queue when it is in a state of the wait kind the acceptance leaves from, and a send back takes it to the state the workflow's send-back transition names, so a workflow with a state of its own before `review` is not deployed, pushed, or queued there.

## Acceptance criteria

- [ ] `flai/internal/serve/kluster.go`, `forgestatus.go`, `forge.go`, `forgemerge.go`, and `mergequeue.go`, `flai/internal/mergequeue`, `flai/internal/kluster/deploy.go`, and `flai/internal/forge/state.go` name no state: each reads a state's kind or a transition's effect from the workflow
- [ ] A story is deployed, pushed with its pull request, and eligible for the merge queue in the state the acceptance transition leaves from, and torn down on the transitions with the acceptance or the cancellation effect, as in `review`, `done`, and `cancelled` today
- [ ] A failing required check, a requested change, a bisected batch, and a failed deploy send the story back along the transition with the send-back effect, with today's reason and thread
- [ ] Tests run the deploy, the forge pass, eligibility, and the send back against a workflow with a state renamed and one added between `in-progress` and `review`, with the outcome the kinds and effects give
- [ ] With the baseline, every existing `serve`, `mergequeue`, `kluster`, and `forge` test passes unchanged

## Tasks

- T-3218 A story's kluster namespace is deployed in the state the acceptance leaves from and torn down on acceptance or cancellation, by effect
- T-3219 The forge pass pushes, reads, and sends back a story by the state the acceptance leaves from and the send-back effect
- T-3220 The merge queue takes stories in the state the acceptance leaves from and sends a bisected one back by the send-back effect

## Notes

- `serve/kluster.go` holds 16 of serve's references, `forgestatus.go` 6, `mergequeue.go` 5; `mergequeue` holds 7.

### Planning

- Touches, all files:
  - Layout: the files of `serve`, `mergequeue`, `kluster`, and `forge` that name a state, with their tests.
  - Co-change: `flai/internal/mergequeue/batch_test.go` (25% in `flai touches suggest`), `state_test.go`, and `kluster/deploy_test.go`.
  - Left out from `flai touches suggest`: `serve/serve.go`, `cmd/queue.go`, and `mergequeue/land.go` and `log.go` (5% to 13%). None names a state.
  - S-0660, under E-0035, also changes `serve/kluster.go`; the trial merge judges that.
- `after`: S-0737. The deploy, push, eligibility, and send back read the acceptance and send-back transitions and the kinds from the workflow S-0737 gives, and the tests use its test hook.
- Tasks: T-3218 (kluster), T-3219 (forge), and T-3220 (merge queue), one layer: no two share a path.
- Forecast 45m, `flai forecast`'s rate for size 25, up from its 40m for the first draft's 22 once the co-change tests were added.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
