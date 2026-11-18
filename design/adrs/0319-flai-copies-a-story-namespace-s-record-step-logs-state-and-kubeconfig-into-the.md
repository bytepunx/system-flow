---
id: ADR-0319
title: "flai copies a story namespace's record, step logs, state, and kubeconfig into the story's agent cache folder, which its container mounts, and flai kluster log, env, and exec read that copy where the main checkout's .flai-cache/kluster is absent, naming the API server as host.docker.internal off the host's network"
status: accepted
date: 2026-11-18
supersedes: []
superseded_by: []
refines: [ADR-0253, ADR-0292]
topics: [kluster, harness]
---

# ADR-0319 flai copies a story namespace's record, step logs, state, and kubeconfig into the story's agent cache folder, which its container mounts, and flai kluster log, env, and exec read that copy where the main checkout's .flai-cache/kluster is absent, naming the API server as host.docker.internal off the host's network

## Context

[ADR-0292](0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) §2 opens `flai kluster env`, `exec`, and `log` to a story's agent for its own namespace. S-0763 built them. All three read the main checkout's `.flai-cache/kluster`, and `env` and `exec` also read the kluster config file from the main checkout's root ([I-0361](../issues/I-0361-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md)).

A story's container mounts, at their host paths, the story's worktree, its agent cache folder `.flai-cache/agents/S-nnnn` read-write, the main checkout's `.git` read-only with the story's registration read-write ([ADR-0246](0246-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-git-read-only.md)), and the main checkout's `wip/`, `system-flow.yaml`, `.flai-cache/verify`, and the checked-write lock ([ADR-0253](0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md), [ADR-0283](0283-a-close-out-s-record-step-confirms-each-finding-outside-the-story-on-a-work.md)). So in a container:

| What the commands read | Where | In the container |
|------------------------|-------|------------------|
| The namespace's record | `.flai-cache/kluster/deployments/s-nnnn.json` | absent |
| A step's log | `.flai-cache/kluster/logs/s-nnnn/<step>.log` | absent |
| The namespace's kubeconfig, and the certificate authority file it may name | `.flai-cache/kluster/agents/s-nnnn.kubeconfig` | absent |
| The state folder, `FLAI_KLUSTER_STATE` | `.flai-cache/kluster/state/s-nnnn/` | absent |
| The kluster config file | the main checkout's root, as `kluster.config` names it | absent; the worktree has it, tracked |
| The worktree `exec` runs in | `.flai-cache/worktrees/S-nnnn` | mounted |

`env` and `exec` stop at `kluster.New`, and `log` finds no log. `.flai-cache/kluster` also holds flai's kubeconfig for the whole cluster and every namespace's token, which no agent may hold (ADR-0292 §2).

The kubeconfig names the API server as the host sees it, such as `https://127.0.0.1:<port>`, where k3d publishes it. From a container on any network but the host's, `127.0.0.1` is the container itself. Every story's container already reaches the host as `host.docker.internal`: `docker run` is given `--add-host host.docker.internal:host-gateway` (`flai/internal/harness/container.go`, `dockerUser`).

## Decision

**flai copies a story namespace's record, step logs, state folder, and kubeconfig into the story's agent cache folder, which the story's container mounts, and `flai kluster log`, `env`, and `exec` read that copy where the main checkout's `.flai-cache/kluster` is absent. Off the host's network, the container is told the API server's host name, `host.docker.internal`, and `env` and `exec` give a kubeconfig that names it.**

### 1. The copy

- **Where.** `.flai-cache/agents/S-nnnn/kluster/`, with the layout of `.flai-cache/kluster` for the story's namespace alone: `deployments/s-nnnn.json`, `logs/s-nnnn/<step>.log`, `state/s-nnnn/`, and `agents/s-nnnn.kubeconfig`. Nothing of another namespace, the default one, or the cluster's kubeconfig is ever copied.
- **When.** Each time flai writes the namespace's record, during a deploy, a checks run, or a teardown, whoever runs the job, `flai serve` or the operator's `flai kluster` command. So the copy follows each step's end. A step's log is not followed while the step runs; the agent waits for the job's `kluster` event (ADR-0292 §5) and reads it then.
- **The kubeconfig.** The namespace's own, with its certificate authority file, when it names one, written into it as `certificate-authority-data`, so the copy names no file outside it. Files are written whole (`atomicfile`), 0600, in folders of mode 0700.
- **Removal.** A teardown that passes removes the copy whole, as it removes the record, the state, the logs, and the kubeconfig. A failed teardown keeps it.
- **For every story's namespace.** The copy is written whatever the story's harness, so an agent that moves to a container finds it at the next job.

### 2. The readers

- **Which root.** `flai kluster log`, `env`, and `exec` read the main checkout's `.flai-cache/kluster` when it is there, as on the host, and otherwise the story's copy, `<main>/.flai-cache/agents/S-nnnn/kluster`.
- **The kluster config.** Read from the main checkout's root, and otherwise from the story's worktree.
- **The variables.** As ADR-0292 §2 has them, with `KUBECONFIG` and `FLAI_KLUSTER_STATE` naming the copy's kubeconfig and state folder when the copy is read. A command that writes under the copy's state folder writes only the copy: the next job writes it again.
- **`exec`** runs in the story's worktree, as before; the container mounts it at its host path.
- **The guard** is unchanged. The copy holds no cluster kubeconfig, so the refusal of a line naming `.flai-cache/kluster/kubeconfig` still holds every agent from the cluster's.

### 3. The API server's host

- **The variable.** A story's container is given `FLAI_KLUSTER_API_HOST=host.docker.internal` on every network but `host`, and no such variable on `host`.
- **The kubeconfig it names.** Where the variable is set and the copy is read, `env` and `exec` write `agents/s-nnnn.container.kubeconfig` beside the copy's: the same file, with the server's host replaced by the variable's value and `tls-server-name` set to the host it replaced, so the API server's certificate is still verified against the name it was issued for. `KUBECONFIG` names that file. It is written again on each call, so it follows a new token.

### 4. What stays

- ADR-0253's mount list. No part of `.flai-cache/kluster` is mounted.
- ADR-0292's boundary: the namespace's kubeconfig, which the API server holds to the namespace, not the guard.

This refines ADR-0292 §2, whose line that an agent in a container reaches the logs only through `flai kluster log` did not hold, and its Consequence that S-0660 says how a container reaches the API server. It refines ADR-0253 in giving the agent cache folder a second use and no new mount.

## Consequences

- A story's agent in a container runs `flai kluster env`, `exec`, and `log` for its own namespace as on the host. I-0361 closes.
- Each story's namespace is kept twice on the host while it lives, in `.flai-cache/kluster` and in the story's cache folder. Logs are the bulk of it.
- The token sits in a second file, in the story's own folder, which only that story's container mounts. It is the token ADR-0292 gives that agent already.
- The container reaches the API server only where k3d publishes it on an address the network's gateway reaches, as on `0.0.0.0`, k3d's default. A config that binds it to `127.0.0.1` alone is reached from the `host` network only.
- That `tls-server-name` keeps verification on `bridge` rests on k3d's certificate naming the host the kubeconfig gives. This repository runs no kluster, so S-0834 tests the rewrite and the layout, not a live API server.
- `flai kluster deploy --job`, `teardown --job`, and `flai checks run --job` from a container's shell write their request under the main checkout's `.flai-cache/kluster/requests/`, which is not mounted. The MCP tools `kluster_deploy`, `kluster_checks`, and `kluster_teardown` run on the host and are what a container's agent uses. That is outside this decision.

## Alternatives considered

- **Mount `.flai-cache/kluster` read-only.** It holds flai's cluster kubeconfig and every namespace's token.
- **Mount the story's own files and folders from `.flai-cache/kluster`.** The record and the kubeconfig are written by rename, which leaves a bind-mounted file stale, and each deploy removes and makes again the log folder, which leaves a bind-mounted folder empty.
- **A host MCP tool for `log`.** It leaves `env` and `exec` unsolved.
- **Run `exec` on the host through flai serve.** It takes the agent's commands out of its container, and every `kubectl get` would wait in the job queue, which ADR-0292 rejected.
- **Rewrite the server's host when the copy is written.** The job does not know which network the story's container joins; the container does.
- **Copy only when the story's agent runs in a container.** An agent restarted in a container after a deploy would find no copy until the next job.
