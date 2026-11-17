---
id: S-0834
type: story
nature: remediation
title: flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount
status: backlog
owner: alex
created: 2026-11-17T08:17:29Z
updated: 2026-11-17T08:32:48Z
transitions: []
tags: [flai, kluster, container]
topics: [kluster, harness]
touches: [design/adrs/README.md, design/system/kluster.md, flai/internal/kluster/agentcopy.go, flai/internal/kluster/agentcopy_test.go, flai/internal/kluster/agentconfig.go, flai/internal/kluster/agentconfig_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, scripts/agent-container.sh, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/cmd/kluster_container_test.go, docs/users/flai-reference.md, design/system/agent-adapters.md, docs/users/flai.md, design/issues/I-0361-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-17T08:17:29Z
  value: 25
  by: planner-S-0834
  at: 2026-11-17T08:31:46Z
forecast:
  duration: 3h
  delivery: 2026-11-17T20:19:00Z
  basis: "Its own forecast of 3h; 22nd in the pull order with an in-progress limit of 3, behind S-0752, S-0806, S-0827, S-0754, S-0811, S-0809, S-0756, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813, S-0814, S-0815, S-0822, S-0823, S-0824, S-0828, S-0829, S-0830, S-0831, S-0832 and S-0833."
  by: flai
  at: 2026-11-17T08:32:48Z
---
# S-0834 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount

## Goal

This story remediates [I-0361](../../../design/issues/I-0361-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md), "flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0361 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0361 is closed with `flai issue close I-0361 --reason` saying what fixed it

## Tasks
- T-3546 An ADR decides where flai kluster exec, env, and log read a container agent's namespace, and the operator accepts it
- T-3547 flai serve copies a story namespace's kubeconfig, record, and step logs into the story's agent cache folder, and a teardown removes them
- T-3548 A story's container is told the host name by which it reaches the cluster's API server on its network

## Notes

Cost of delay inputs set by flai from I-0361. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T12:16:36Z, 3.8 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0834 on 2026-11-17. The plan and its assumptions are on the plan's thread on this story.

The proposed fix, which the first task's ADR settles with the operator: flai serve copies the story's own namespace kubeconfig, deployment record, and step logs into the story's agent cache folder, `.flai-cache/agents/<S-nnnn>/kluster/`, which the container mounts already. `flai kluster exec`, `env`, and `log` read that copy when the main checkout's `.flai-cache/kluster` is absent. They read the kluster config from the worktree, and run `exec` in the worktree the container mounts. ADR-0253's mount list stays as it is: no part of `.flai-cache/kluster` is mounted, since it holds flai's cluster-wide kubeconfig and every namespace's token.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `design/adrs/README.md`, `design/system/kluster.md` | design: the ADR that proposes the fix, its index, and the "Reads" bullet of kluster.md § What the agent runs against its namespace, which says the commands serve an agent on the host |
| `flai/internal/kluster/agentcopy.go`, `agentcopy_test.go` | layout: a new file for the copy into the agent cache folder |
| `flai/internal/kluster/agentconfig.go`, `agentconfig_test.go` | co-change and layout: `AgentKubeconfigPath` and `AgentEnv`, which `exec` and `env` read |
| `flai/internal/kluster/deploy.go`, `deploy_test.go` | co-change and layout: the deploy, checks, and teardown steps that write the record and logs, where the copy is made and removed |
| `flai/internal/kluster/kluster.go`, `kluster_test.go` | layout: `kluster.New` reads the kluster config from the main checkout's root, which is not mounted |
| `flai/internal/harness/container.go`, `container_test.go`, `scripts/agent-container.sh` | co-change: the container's variables, which gain how the container reaches the host's API server on its network |
| `flai/cmd/kluster.go`, `kluster_test.go`, `kluster_container_test.go` | declared by the issue and co-change: `exec`, `env`, and `log`, and a test over a container's layout that reproduces I-0361 |
| `docs/users/flai-reference.md` | co-change: the commands' help text changes, and the generated reference with it |
| `design/system/agent-adapters.md`, `docs/users/flai.md` | design: the container's mounts row and the user guide's exec and log sections |
| `design/issues/I-0361-…md`, `design/issues/summary.md` | the second criterion: `flai issue close` writes them |

No folder touch is kept. The ADR's own file is named when it is written, and `task_done` widens the task's touches to it.

Figures:

- **Forecast 3h, delivery 2026-11-18T12:00:00Z.** `flai forecast` said 6m, from a size of 2: two criteria and no touches. That size leaves out the work: an ADR and four code and docs tasks across `kluster`, `harness`, and `cmd`, 21 files. 3h follows the 30m to 45m a task of this kind took in E-0035's stories. Delivery is flai's 15:16 today, plus this duration, plus a night's wait for the ADR's acceptance.
- **Cost of delay 25 USD a week.** `flai cod` works it from the one input, 10m lost per 168h cycle at 150 USD an hour. It stands. This repository runs no kluster, so the loss falls on projects whose agents run in containers against a namespace, and one occurrence is all the issue records.

No `after`: the story reads nothing that another open story writes.
