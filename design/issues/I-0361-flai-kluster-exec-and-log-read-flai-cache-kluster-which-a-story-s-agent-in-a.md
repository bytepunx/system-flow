---
id: I-0361
title: flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-13T12:16:36Z
last_reported: 2026-11-13T12:16:36Z
updated: 2026-11-13T12:16:36Z
---

# I-0361 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount

## Description
flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount

## Instances

### 2026-11-13T12:16:36Z
Story: S-0763.
S-0763 built flai kluster exec and log for a story's agent against its own namespace. Both read the main checkout's .flai-cache/kluster (the record, the logs, and the namespace's kubeconfig), and the agent container mounts the story's worktree, its .flai-cache folder, .git, wip/, system-flow.yaml, and .flai-cache/verify only (agent-adapters.md), so an agent in a container finds no record, no log, and no kubeconfig there. A sub-agent's draft said the opposite, that a container agent reaches the logs only through flai kluster log; the docs now say the commands serve an agent on the host. A container agent needs the namespace's kubeconfig and logs mounted, or a host method behind exec and log.

## Remediation
