---
id: T-3546
type: task
nature: research
title: An ADR decides where flai kluster exec, env, and log read a container agent's namespace, and the operator accepts it
status: done
parent: S-0834
owner: alex
created: 2026-11-17T08:32:08Z
updated: 2026-11-18T07:10:55Z
transitions:
  - to: ready
    at: 2026-11-17T12:06:03Z
    by: agent-S-0834
  - to: in-progress
    at: 2026-11-17T12:06:03Z
    by: agent-S-0834
  - to: done
    at: 2026-11-18T07:10:55Z
    by: agent-S-0834
stream: S-0834
tags: [adr, kluster, container]
touches: [design/adrs/README.md, design/adrs/0319-flai-copies-a-story-namespace-s-record-step-logs-state-and-kubeconfig-into-the.md]
usage:
  source: log
  seconds: 118
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 9294
      cache_read: 2838615
      cache_write: 123790
      cost: 1.744
---
# T-3546 An ADR decides where flai kluster exec, env, and log read a container agent's namespace, and the operator accepts it

## Work

I-0361 asks for a fix to be proposed before it is built. Write a proposed ADR that refines ADR-0292 §2 and ADR-0253. It decides where `flai kluster exec`, `env`, and `log` find a story's namespace when the story's agent runs in a container.

The planner's recommendation is on the plan's thread:

- flai serve copies the namespace's own kubeconfig, its deployment record, and its step logs into the story's agent cache folder, `.flai-cache/agents/<S-nnnn>/kluster/`. The container mounts that folder already, so ADR-0253's mount list is unchanged.
- The copy is written whole, at the end of each deploy, checks run, and teardown. A teardown that passes removes it.
- The three commands read the copy when the main checkout's `.flai-cache/kluster` is absent. They read the kluster config from the worktree, and run `exec` in the worktree the container mounts.
- The kubeconfig reaches the cluster's API server on the container's network. On `bridge`, its server is `host.docker.internal`, with `tls-server-name` keeping the certificate's name. On `host`, it is unchanged.

Weigh the alternatives the explorer found:

- Mount `.flai-cache/kluster` read-only. This is refused, because it holds flai's cluster-wide kubeconfig and every namespace's token.
- Mount one file per story. A bind-mounted file goes stale when `atomicfile` renames over it, and `RemoveAll` recreates the log folder.
- A host MCP tool for `log`. It leaves `exec` and `env` unsolved.
- Run `exec` on the host. That runs the agent's commands outside its container.

Put the decision into `design/system/kluster.md` § What the agent runs against its namespace, "Reads". Add the ADR to `design/adrs/README.md`. Open a thread on S-0834 asking the operator to accept it.

This task waits for nothing. Every other task builds what it decides.

## Done when

- The ADR is written and listed in `design/adrs/README.md`. It names the copy's folder and files, when they are written and removed, how the commands find them, and how the kubeconfig reaches the API server on `bridge` and `host`.
- `design/system/kluster.md` says the same.
- The operator has accepted the ADR on a thread on S-0834, or changed it, and the ADR records what they decided.

## Notes
