---
id: S-0665
type: story
nature: improvement
title: flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image
status: backlog
owner: alex
created: 2026-11-10T08:34:37Z
updated: 2026-11-10T08:54:23Z
transitions: []
tags: [serve, containers]
topics: [docker, code]
touches: [flai/internal/harness/harness.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_start.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/plan.go, flai/internal/serve/orchestrate.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, design/system/flai-cli.md, docs/users/flai.md, design/system/strategic-agents.md, design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md, design/issues/summary.md]
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
      seconds: 196
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 170
          output: 730
          cache_read: 9317567
          cache_write: 313234
          cost: 0.6455
        - model: claude-opus-5-5
          input: 553
          output: 15503
          cache_read: 42559750
          cache_write: 641675
          cost: 11.7787
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-10T08:34:37Z
  value: 37.5
  by: planner-S-0665
  at: 2026-11-10T08:51:29Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T14:39:00Z
  basis: "Its own forecast of 1h30m; 11th in the pull order with an in-progress limit of 5, behind S-0680, S-0682, S-0684, S-0685, S-0687, S-0674, S-0673, S-0672, S-0671, S-0241, S-0659, S-0660, S-0661, S-0662 and S-0663."
  by: flai
  at: 2026-11-10T08:54:23Z
---
# S-0665 flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image

## Goal

This story remediates [I-0274](../../../design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md), "flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0274 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0274 is closed with `flai issue close I-0274 --reason` saying what fixed it

## Tasks
- T-2942 The container adapter's Start names the container's image, name, and mounts
- T-2943 flai serve records a container run's image, container name, and mounts in agents.json, with a test reproducing I-0274
- T-2945 The serve journal's start line of a container run names its image and container
- T-2946 The design and the user guide say what a container run records and logs
- T-2948 I-0274 is closed with what fixed it

## Notes

Cost of delay inputs set by flai from I-0274. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T05:41:15Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0274's one instance (S-0567's planner run): the container adapter already computes the image, the container's name, and the mounts when it builds the `docker run` argv (`flai/internal/harness/container.go`, `Start` and `roleStart`). It then drops them, and `spawnWith` keeps only `docker` as the run's `command`. So the adapter's `harness.Start` should carry them as a field, and `startOf` (`flai/internal/serve/agents_start.go`), which every story and strategic run goes through, should copy them onto `AgentRun` as `image`, `container`, and `mounts` in `agents.json`. The journal's `started docker (container)` line should also name the image and container. Then a run removed by `--rm` can still be shown to have used the layer image. The host API's `agent.status` already returns the whole run, so the fields reach the dashboard's data. Showing them on a dashboard page is left out.

Touches, all files, no folder:

| Touch | From |
|-------|------|
| `flai/internal/harness/harness.go`, `flai/internal/harness/container.go`, `flai/internal/harness/container_test.go` | layout: `harness.Start` and the argv builder; co-change with `flai/internal/serve/container.go` (`flai touches suggest`: 50% and 14%) |
| `flai/internal/serve/agents.go`, `flai/internal/serve/agents_start.go`, `flai/internal/serve/container.go`, `flai/internal/serve/container_test.go` | layout: `AgentRun`, `startOf`, `containerReady` and `stopContainer`; co-change 43% for `agents.go` and `container_test.go` |
| `flai/internal/serve/plan.go`, `flai/internal/serve/orchestrate.go`, `flai/internal/serve/analyze.go`, `flai/internal/serve/curate.go` | layout: each writes its own `started %s (%s)` journal line |
| `design/system/flai-cli.md`, `docs/users/flai.md`, `design/system/strategic-agents.md` | design: where `harness_version` and the S-0567 live run are documented |
| I-0274's file, `design/issues/summary.md` | design: the second criterion's `flai issue close` |

`flai/internal/serve/role.go` is left out: `containerSubject` refuses a role a project defines on the container harness, so its start line never reads a container run. The story declared no touches before this plan. It has no `after`: it reads no other open story's output.

Forecast: `flai forecast` gave 31m, delivery 2026-11-10T12:03:00Z. I raised it to 1h30m, delivery 2026-11-10T13:05:00Z. The story has five tasks in four layers across two Go packages and five start paths, a fake-docker reproduction, and three docs. Its close-out also runs integration and smoke, as it does for any story that touches `flai/`.

Cost of delay: `flai cod` gives 37.50 USD a week, from flai's input of 15m lost per 168h cycle at 150 USD an hour. It stands. The input comes from I-0274's one occurrence, and nothing here gives a reason to change it.
