---
id: S-0763
type: story
nature: feature
title: Each deploy gives a story's namespace a kubeconfig held to it, and a story's agent runs flai kluster exec and log against its own namespace
status: backlog
parent: E-0035
owner: alex
created: 2026-11-13T08:03:30Z
updated: 2026-11-13T08:03:30Z
transitions: []
tags: [cli]
topics: [kluster, deployment]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0660]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

## Notes

- Split out of S-0660 on TH-1138: the operator confirmed the orchestrator's recommendation that S-0660 keep to its nine criteria, and that this story build the namespace's kubeconfig, `exec`, and `log`. `flai kluster env` is S-0661's (its T-2865), and S-0661 waits for this story, since what `env` gives includes the namespace's `KUBECONFIG`.
- `after`: S-0660, whose deploy and teardown requests and job queue this story's kubeconfig step joins.
