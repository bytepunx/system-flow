---
id: T-3549
type: task
nature: remediation
title: flai kluster exec, env, and log read the story's copy in a container, with a test over a container's layout that reproduces I-0361
status: done
parent: S-0834
owner: alex
created: 2026-11-17T08:32:57Z
updated: 2026-11-18T07:59:53Z
transitions:
  - to: ready
    at: 2026-11-18T07:40:59Z
    by: agent-S-0834
  - to: in-progress
    at: 2026-11-18T07:40:59Z
    by: agent-S-0834
  - to: done
    at: 2026-11-18T07:59:52Z
    by: agent-S-0834
stream: S-0834
tags: [kluster, cli, go]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_container_test.go, flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/internal/kluster/agentconfig.go, flai/internal/kluster/agentconfig_test.go, docs/users/flai-reference.md, flai/internal/kluster/log.go]
after: [T-3547, T-3548]
usage:
  source: log
  seconds: 1133
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 19405
      cache_read: 5766236
      cache_write: 214544
      cost: 2.862
---
# T-3549 flai kluster exec, env, and log read the story's copy in a container, with a test over a container's layout that reproduces I-0361

## Work

This task reads the copy T-3547 writes and the variable T-3548 gives the container, so it waits for both.

First write the test that reproduces I-0361, in `flai/cmd/kluster_container_test.go`, and watch it fail. Its fixture has a container's layout:

- a main checkout holding only `.git`, `wip/`, `system-flow.yaml`, and the story's `.flai-cache/agents/<S-nnnn>/kluster/` copy;
- no `.flai-cache/kluster` and no kluster config at its root;
- a story worktree holding the kluster config.

Today `exec` and `env` fail at `kluster.New`, and `log` says no job ran the step.

Then make it pass:

- When the main checkout's `.flai-cache/kluster` is absent, `klusterAgentEnv` and `readWholeLog` in `flai/cmd/kluster.go` read the story's copy (`AgentCopyDir`): its record, its kubeconfig, and its step logs.
- `kluster.New` in `flai/internal/kluster/kluster.go` reads the kluster config from the worktree when the main checkout's root does not hold it.
- `exec` runs in the worktree the container mounts, `repo.Root`, rather than `<main>/.flai-cache/worktrees/<S-nnnn>`, which is not mounted.
- When the variable T-3548 gives is set, `exec` and `env` set `KUBECONFIG` to a copy of the kubeconfig whose server host is that host, with `tls-server-name` keeping the host the certificate names. Write the copy under the story's cache folder, mode 0600.
- On the host, every command reads as it does today. The existing tests in `kluster_test.go` hold that.
- Change the commands' help text, which says `log` reads the logs of the checkout it runs in. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- `kluster_container_test.go` runs `exec`, `env`, and `log` over the container layout. `exec` runs in the worktree with the copy's kubeconfig, its server rewritten when the variable is set. `env` prints the copy's variables. `log` prints the copied step log.
- The tests fail without the change, and pass with it.
- `kluster_test.go` passes unchanged for the host.
- `TestReferenceIsCurrent` passes.
- `flai test` passes on the changed paths.

## Notes
