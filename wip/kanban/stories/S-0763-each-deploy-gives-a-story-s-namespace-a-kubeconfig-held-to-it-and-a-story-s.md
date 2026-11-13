---
id: S-0763
type: story
nature: feature
title: Each deploy gives a story's namespace a kubeconfig held to it, and a story's agent runs flai kluster exec and log against its own namespace
status: in-progress
parent: E-0035
owner: alex
created: 2026-11-13T08:03:30Z
updated: 2026-11-13T09:06:55Z
transitions:
  - to: ready
    at: 2026-11-13T09:06:47Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-13T09:06:55Z
    by: agent-S-0763
tags: [cli]
topics: [kluster, deployment]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/kluster/agentconfig.go, flai/internal/kluster/agentconfig_test.go, flai/internal/kluster/log.go, flai/internal/kluster/log_test.go, flai/internal/hostapi/kluster.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/internal/guard/guard_rules.go, flai/internal/guard/call.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [S-0660]
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
      seconds: 528
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 72
          output: 23419
          cache_read: 2900832
          cache_write: 120148
          cost: 0.1556
        - model: claude-opus-5-5
          input: 90
          output: 27856
          cache_read: 9459441
          cache_write: 273897
          cost: 4.6405
    - kind: orchestrator
      seconds: 104
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 27
          output: 217
          cache_read: 1105393
          cache_write: 81319
          cost: 0.3419
cost_of_delay:
  value: 555.56
  by: planner-S-0763
  at: 2026-11-13T09:04:24Z
forecast:
  duration: 1h
  delivery: 2026-11-13T11:53:00Z
  basis: "Its own forecast of 1h; 7th in the pull order with an in-progress limit of 3, behind S-0735, S-0761, S-0760, S-0745, S-0759, S-0751, S-0757 and S-0740."
  by: flai
  at: 2026-11-13T09:06:51Z
---
# S-0763 Each deploy gives a story's namespace a kubeconfig held to it, and a story's agent runs flai kluster exec and log against its own namespace

## Goal

A story's agent reaches its own story's kluster namespace, and nothing else in the cluster, with its own commands. Each deploy makes a kubeconfig held to the namespace by Kubernetes RBAC, and the agent runs `flai kluster exec` and `flai kluster log` for its own story, as [ADR-0292](../../../design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) §2 decides.

## Acceptance criteria

- [ ] After the deploy step of a story's deploy passes, flai makes a ServiceAccount `flai-agent` in the namespace, a RoleBinding to the cluster's `edit` ClusterRole, and `.flai-cache/kluster/agents/s-nnnn.kubeconfig` with a token from `kubectl create token` valid for 24 hours; a deploy again makes a new one
- [ ] The teardown deletes the namespace's kubeconfig
- [ ] `flai kluster exec S-nnnn -- <command>` runs a command in the story's worktree with ADR-0254 §4's variables for its namespace, but `FLAI_KLUSTER_IMAGES`, and `KUBECONFIG` the namespace's kubeconfig, while the namespace has a record, a failed one included
- [ ] `flai kluster log S-nnnn --step <step>` prints a step's log as the host method `kluster.log` reads it
- [ ] `flai guard` lets a story's agent and its sub-agents run `exec` and `log` for their own story, and refuses them for another story's namespace or the default one, the `kluster` program but `status`, and a shell line naming `.flai-cache/kluster/kubeconfig`
- [ ] `design/system/kluster.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` and its reference say all of it

## Tasks
- T-3276 Each passing deploy makes the namespace's own kubeconfig from a flai-agent ServiceAccount bound to edit, and the teardown deletes it
- T-3277 flai kluster log prints a step's log of the story's namespace through the reader the host method kluster.log uses
- T-3278 flai guard lets a story's agent and its sub-agents run kluster exec and log for their own story, and refuses other namespaces, the kluster program but status, and flai's cluster kubeconfig
- T-3279 flai kluster exec runs a command in the story's worktree with its namespace's variables and its own kubeconfig while the namespace has a record
- T-3280 kluster.md, flai-cli.md, the user guide, and the settings index say how a story's namespace gets its own kubeconfig and how its agent runs exec and log

## Notes

- Split out of S-0660 on TH-1138: the operator confirmed the orchestrator's recommendation that S-0660 keep to its nine criteria, and that this story build the namespace's kubeconfig, `exec`, and `log`. `flai kluster env` is S-0661's (its T-2865), and S-0661 waits for this story, since what `env` gives includes the namespace's `KUBECONFIG`.
- `after`: S-0660, whose deploy and teardown requests and job queue this story's kubeconfig step joins.

### Planning

Touches, by where each came from:

| Touch | From |
|-------|------|
| `flai/internal/kluster/deploy.go`, `deploy_test.go` | Declared; the layout: `runDeploy` is where the deploy step passes, and the teardown removes the record, state, and logs there |
| `flai/internal/kluster/agentconfig.go`, `agentconfig_test.go` | Layout: a new file for the kubeconfig step, its path, and `Cluster.AgentEnv` (T-3276) |
| `flai/internal/kluster/log.go`, `log_test.go` | Layout: the `kluster.log` reader moves out of the host API so `flai kluster log` reads as it does (T-3277) |
| `flai/internal/hostapi/kluster.go` | Layout: `klusterLogIn` and `klusterLogDo` read the log there today, unexported |
| `flai/cmd/kluster.go`, `kluster_test.go` | Declared; the layout: the `flai kluster` commands |
| `flai/internal/guard/guard_rules.go`, `kluster_test.go` | Declared; the layout: `klustering`, `klusterWrite`, and `ownStory` |
| `flai/internal/guard/call.go` | Layout: where the Bash call reaches `klusterWrite`, for the two line rules |
| `design/system/kluster.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` | Declared; criterion 6 |
| `docs/operators/settings.md` | Co-change (21% of the commits that change the declared paths) and the layout: `scripts/flai-reference.sh` writes its flag index with the reference, and its row of the kluster variables names the commands that give them |

No folder touch. `flai/internal/serve` is left out: its jobs call `DeployBy` and `TeardownBy`, which T-3276 changes. `flai/internal/mcpserver` is left out: the criteria name no MCP tool for `exec` or `log`.

`after`: S-0660, done. T-3276 builds on its `DeployBy` and `TeardownBy` in `flai/internal/kluster/deploy.go`, which name who asked, and on the deploys its requests queue, which the kubeconfig step joins. T-3277 and T-3279 build on its `sessionStory` in `flai/cmd/kluster.go`, and T-3278 on its `ownStory` and request rules in `flai/internal/guard/guard_rules.go`.

Forecast: `flai forecast` gives 42m from size 23 (6 criteria, 17 touches). Raised to 1h: S-0660, the kluster story this was split from, took 1h23m of agent work for six tasks, and this one adds a kubectl step to the deploy, moves the `kluster.log` reader, adds two commands, and adds guard rules that read whole shell lines. The delivery is flai's play-out of the pull order with that duration: 31st, behind S-0661, S-0662, and S-0663, which wait for it.

Cost of delay: 555.56 USD a week, as `flai cod S-0763` works it out: S-0763's share of the operator's `revenue_per_week: 1250` on E-0035, split by forecast duration between the epic's two open stories without inputs of their own, S-0763 (1h) and S-0661 (1h15m). Not adjusted: the split follows the operator's input and flai's rule.
