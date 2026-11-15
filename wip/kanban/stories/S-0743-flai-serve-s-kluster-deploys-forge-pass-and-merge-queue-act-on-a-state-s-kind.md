---
id: S-0743
type: story
nature: feature
title: flai serve's kluster deploys, forge pass, and merge queue act on a state's kind and a transition's effect instead of review and done
status: in-progress
parent: E-0034
owner: alex
created: 2026-11-12T15:59:52Z
updated: 2026-11-15T01:51:11Z
transitions:
  - to: ready
    at: 2026-11-14T17:47:55Z
    by: alex
  - to: in-progress
    at: 2026-11-15T00:51:05Z
    by: agent-S-0743
tags: [cli]
topics: [cli, kluster]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/forge.go, flai/internal/serve/forge_test.go, flai/internal/serve/forgemerge.go, flai/internal/forge/state.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/state.go, flai/internal/serve/flow.go, flai/internal/serve/flow_test.go, flai/internal/serve/serve.go, flai/internal/workflow/review.go, flai/internal/workflow/review_test.go, flai/internal/kluster/kluster.go, flai/internal/mergequeue/preflight.go, flai/internal/serve/kluster_requests.go, flai/internal/serve/kluster_requests_test.go, design/system/workflow.md]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3565
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-15
      ceremony: 1
      work: 67
  models:
    - model: claude-opus-5-5
      input: 454
      output: 2622
      cache_read: 34292156
      cache_write: 1661731
      cost: 13.5535
  strategic:
    - kind: orchestrator
      seconds: 29
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 45
          cache_read: 330385
          cache_write: 51750
          cost: 0.1122
cost_of_delay:
  value: 80.17
  by: system-flow
  at: 2026-11-14T17:44:05Z
forecast:
  duration: 45m
  delivery: 2026-11-15T01:22:00Z
  basis: "Its own forecast of 45m; 5th in the pull order with an in-progress limit of 3, behind S-0738, S-0741, S-0746 and S-0739."
  by: flai
  at: 2026-11-14T23:48:54Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:53Z
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
- T-3369 The workflow names the states the acceptance leaves from, the send-back out of a state, and the closed kinds, and flai serve opens its repositories with the workflow its options give
- T-3371 An agent's own kluster deploy request and the merge queue's preflight read the workflow, not a state's name

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
- Cost of delay: 80.17 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 45m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
