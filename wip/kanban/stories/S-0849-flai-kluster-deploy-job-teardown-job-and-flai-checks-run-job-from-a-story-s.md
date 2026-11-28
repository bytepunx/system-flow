---
id: S-0849
type: story
nature: remediation
title: flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount
status: backlog
owner: alex
created: 2026-11-21T18:57:50Z
updated: 2026-11-28T00:26:08Z
transitions: []
tags: []
topics: [kluster, harness]
touches: [design/adrs, design/adrs/README.md, flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/internal/serve/kluster_requests.go, flai/internal/serve/kluster_requests_test.go, flai/cmd/kluster.go, flai/cmd/checks.go, flai/cmd/kluster_container_test.go, design/system/kluster.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0405-flai-kluster-deploy-job-teardown-job-and-flai-checks-run-job-from-a-story-s.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 128
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 27
          output: 8982
          cache_read: 2224464
          cache_write: 75508
          cost: 1.2287
    - kind: orchestrator
      seconds: 106
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 33
          output: 244
          cache_read: 1433289
          cache_write: 180411
          cost: 0.5142
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-21T18:57:50Z
  value: 25
  by: planner-S-0849
  at: 2026-11-28T00:04:06Z
forecast:
  duration: 1h30m
  delivery: 2026-11-28T09:07:00Z
  basis: "Its own forecast of 1h30m; 12th in the pull order with an in-progress limit of 3, behind S-0835, S-0837, S-0844, S-0814, S-0852, S-0813, S-0822, S-0848, S-0853, S-0854, S-0824, S-0845, S-0241 and S-0815."
  by: flai
  at: 2026-11-28T00:26:08Z
---
# S-0849 flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount

## Goal

This story remediates [I-0405](../../../design/issues/I-0405-flai-kluster-deploy-job-teardown-job-and-flai-checks-run-job-from-a-story-s.md), "flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0405 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0405 is closed with `flai issue close I-0405 --reason` saying what fixed it

## Tasks
- T-3655 Write the ADR: a story's container asks for its kluster jobs through a requests folder in its agent cache folder, which flai serve also takes from
- T-3656 Give the kluster request files a story's container requests folder, keeping the main checkout's for the host
- T-3657 flai serve's kluster request pass takes requests from each story's agent cache folder and answers there
- T-3658 flai kluster deploy --job, teardown --job, and flai checks run --job ask through the story's agent cache folder in its container, with a test that reproduces I-0405
- T-3659 Say in kluster.md and the users guide that a container's --job asks through its cache folder, and close I-0405

## Notes

Cost of delay inputs set by flai from I-0405. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-18T08:00:16Z, 3.5 days before this story; under one cycle counts as one).

### Planning

Proposed fix, the planner's, for the operator to confirm on the plan's thread: in a story's container, where the main checkout's `.flai-cache/kluster` is absent, `flai kluster deploy --job`, `teardown --job`, and `flai checks run --job` write their request in the story's agent cache folder, which the container mounts read-write already (ADR-0319), in a folder of its own beside the copy, `.flai-cache/agents/S-nnnn/kluster-requests/`. flai serve's request pass also takes requests from each story's folder, refuses one whose story is not the folder's, and writes its answer and the story's job state there too, so `--wait` works from the container. The folder sits outside `kluster/`, because a teardown that passes removes that copy whole. A new ADR records it, refining ADR-0292 §1 and ADR-0319. The alternative, a container's `--job` refusing at once and naming the MCP tool, is on the thread.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/kluster/request.go`, `flai/cmd/kluster.go` | design: I-0405's instance names them |
| `flai/cmd/checks.go` | design: I-0405 names it; its help names `.flai-cache/kluster/requests` |
| `flai/internal/serve/kluster_requests.go` | layout: the pass that takes and answers the requests |
| `flai/internal/kluster/request_test.go`, `flai/internal/serve/kluster_requests_test.go` | co-change, 16% each |
| `flai/cmd/kluster_container_test.go` | layout: the container tests of ADR-0319's readers, where the reproducing test fits |
| `docs/users/flai-reference.md` | co-change, 53%: the help text of the three commands changes |
| `docs/users/flai.md`, `design/system/kluster.md` | co-change and design: both say a container's `--job` does not reach flai serve |
| `design/adrs`, `design/adrs/README.md` | design: the new ADR and its generated index |
| I-0405 and `design/issues/summary.md` | criteria: closing the issue |

The folder touch `design/adrs` is kept because the new ADR's number and slug are not known until it is written, and other open stories also add ADRs. No task names it.

Left out: `flai/internal/mcpserver/kluster.go`, since the MCP tools run on the host and keep the main checkout's folder, provided the root-based functions in `request.go` stay as they are. `flai/internal/kluster/agentcopy.go` is left out too, since the requests folder sits outside the copy. `docs/operators/settings.md` is left out because no setting changes.

Forecast: 1h30m, raised from flai's 33m. flai worked 120 s per unit of size over 86 done remediation stories, times size 16. S-0834, the agent's copy, is the closest story: it took 85m by its usage. This one also adds an ADR and changes the kluster package, flai serve's pass, the command, and the docs. Delivery 2026-11-28T08:45:00Z, flai's 07:47 plus the added time.

Cost of delay: 25 USD a week, as flai worked it out from the operator's input, 10m lost per 168h cycle at 150 USD an hour. It stands: the MCP tools work around the defect, so little time is lost per occurrence.

No `after`: the story reads nothing from another open story. S-0835 also touches `flai/internal/serve/kluster_requests_test.go`, `design/system/kluster.md`, and `docs/users/flai.md`. A trial merge judges that overlap.
