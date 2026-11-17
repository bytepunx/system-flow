---
id: TH-1286
title: "S-0834 plan: an ADR first, then a copy of the namespace in the story's cache folder, in five tasks over four layers, 3h, 25 USD a week"
anchor:
  path: wip/kanban/stories/S-0834-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md
  item: S-0834
status: open
participants: [planner-S-0834]
created: 2026-11-17T08:34:14Z
updated: 2026-11-17T08:34:14Z
---

# TH-1286 S-0834 plan: an ADR first, then a copy of the namespace in the story's cache folder, in five tasks over four layers, 3h, 25 USD a week

On wip/kanban/stories/S-0834-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md.

## Entries

### 2026-11-17T08:34:14Z planner-S-0834
**Plan: five tasks in four layers.** T-3546's ADR proposes the fix, and the operator accepts it before any code is written, as I-0361 asks.

| Layer | Task | Waits for, and why |
|-------|------|--------------------|
| 1 | T-3546 An ADR decides where `exec`, `env`, and `log` read a container agent's namespace, and the operator accepts it | nothing |
| 2 | T-3547 flai serve copies the namespace's kubeconfig, record, and step logs into `.flai-cache/agents/<S-nnnn>/kluster/`, and a teardown removes them | T-3546: builds what it decides |
| 2 | T-3548 The container is given the host name of the cluster's API server on its network (`host.docker.internal` on bridge) | T-3546: the variable's name and values |
| 3 | T-3549 `exec`, `env`, and `log` read the copy in a container, with a test over a container's layout that reproduces I-0361 | T-3547 and T-3548: reads the copy and the variable |
| 4 | T-3550 kluster.md, agent-adapters.md, and flai.md say what was built, and I-0361 is closed | T-3549 |

T-3547 and T-3548 share no file, so they can run together.

**Recommended fix, which T-3546 puts to you:**

- Copy the story's own kubeconfig, record, and step logs into the story's agent cache folder. The container mounts that folder already, so ADR-0253's mount list does not change.
- The commands read the copy when the main checkout's `.flai-cache/kluster` is absent.

Alternatives, and why I do not recommend them:

- Mount `.flai-cache/kluster` read-only. It holds flai's cluster-wide kubeconfig and every namespace's token.
- Mount one file per story. A bind-mounted file goes stale when `atomicfile` renames over it.
- A host MCP tool for `log`. It leaves `exec` and `env` unsolved.
- Run `exec` on the host. That takes the agent's commands out of its container.

**Assumptions:**

1. The explorer found that in a container, `exec` and `env` also fail at `kluster.New`. It reads the kluster config from the main checkout's root, which is not mounted. `exec` also runs in `<main>/.flai-cache/worktrees/<S-nnnn>`, which is not mounted either. T-3549 fixes both, since I-0361's cause covers them.
2. The k3d API server's certificate names the host its kubeconfig gives, so a `tls-server-name` keeps verification on bridge. If it does not, T-3546 says so, and the ADR settles it.
3. `kubectl` in the container is the project's image layer's to provide, not this story's. This repository runs no kluster, so its layer is left alone.
4. The forecast is 3h against flai's 6m, which came from two criteria and no touches. The cost of delay value is 25 USD a week, from flai's one input, unchanged.

**Proposed addition, not drafted:**

- From the code, the explorer inferred that `flai kluster deploy --job`, `teardown --job`, and `flai checks run --job` from a container's shell write their request under the main checkout's `.flai-cache/kluster/requests/`, which is not mounted. flai serve then never sees them. The MCP tools `kluster_deploy`, `kluster_checks`, and `kluster_teardown` do work, because they run on the host.
- This is not run or verified. I would record it as an issue for its own story rather than widen S-0834. Say if you want it folded in here instead.
