---
id: I-0361
title: flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-11-13T12:16:36Z
last_reported: 2026-11-13T12:16:36Z
updated: 2026-11-18T08:07:19Z
---

# I-0361 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount

## Description
flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount

## Instances

### 2026-11-13T12:16:36Z
Story: S-0763.
S-0763 built flai kluster exec and log for a story's agent against its own namespace. Both read the main checkout's .flai-cache/kluster (the record, the logs, and the namespace's kubeconfig), and the agent container mounts the story's worktree, its .flai-cache folder, .git, wip/, system-flow.yaml, and .flai-cache/verify only (agent-adapters.md), so an agent in a container finds no record, no log, and no kubeconfig there. A sub-agent's draft said the opposite, that a container agent reaches the logs only through flai kluster log; the docs now say the commands serve an agent on the host. A container agent needs the namespace's kubeconfig and logs mounted, or a host method behind exec and log.

## Remediation

Story S-0834 remediates this issue, created from it at 2026-11-17T08:17:29Z.
Closed 2026-11-18T08:07:19Z: S-0834 fixed it under ADR-0319. T-3547: each write of a story namespace's record copies its record, step logs, state, and kubeconfig, with the CA inlined, into the story's agent cache folder, .flai-cache/agents/S-nnnn/kluster, which its container mounts. T-3548: a story's container gets FLAI_KLUSTER_API_HOST=host.docker.internal off the host's network. T-3549: flai kluster log, env, and exec read that copy, and the kluster config from the worktree, where .flai-cache/kluster is absent. env and exec also write a kubeconfig naming the API server at that host, with tls-server-name kept. flai/cmd/kluster_container_test.go reproduces the cause over a container's layout. The --job requests from a container's shell are left to I-0405.
