---
id: T-2562
type: task
nature: improvement
title: kluster.status says what setup still lacks, and kluster.workloads lists the cluster's nodes and pods
status: done
parent: S-0586
owner: alex
created: 2026-11-07T08:31:17Z
updated: 2026-11-07T09:18:09Z
transitions:
  - to: ready
    at: 2026-11-07T09:06:45Z
    by: agent-S-0586
  - to: in-progress
    at: 2026-11-07T09:06:46Z
    by: agent-S-0586
  - to: done
    at: 2026-11-07T09:18:09Z
    by: agent-S-0586
stream: S-0586
tags: [kluster, cli, flaiover]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/workloads/+server.ts, flaiover/src/routes/api/kluster/kluster.test.ts, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go]
after: [T-2560]
usage:
  source: log
  seconds: 683
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 102
      output: 33466
      cache_read: 12169386
      cache_write: 205005
      cost: 4.287
---
# T-2562 kluster.status says what setup still lacks, and kluster.workloads lists the cluster's nodes and pods

## Work

Waits for T-2560, whose block and setting keys these fields report.

Give the Klusters page what it needs to lead an operator through setup, in `flai/internal/hostapi/kluster.go`. Each stays a read, whatever the host action says, and runs no flai command.

- **`kluster.status` gains:**
  - `settings`: the block's keys as the manifest holds them, with each key's default.
  - `settings_on`: whether the `settings` host action is on, since the block's and the config's writes need it.
  - `config`: the config's path.
  - `findings`: each `manifest.kluster` sentence `KlusterErrors` gives.
  - `deployable`: whether `deploy` is set.
  - `tools`: for each of `kluster`, `docker`, and `kubectl`, whether it is on the host's `PATH` and its path, by `Exec.LookPath`.
- **`kluster.workloads`, a new read:** the cluster's nodes, and each namespace's pods with their phase, ready count, and restarts. It runs `kubectl get nodes` and `kubectl get pods --all-namespaces`, each `--output json`, with flai's kubeconfig and the cluster's context, within 10 seconds, only when the kubeconfig exists, as `access` reads the services. `workloads_error` says why it could not read them.

Add the route `flaiover/src/routes/api/kluster/workloads/+server.ts`, the method in `flaiover/src/lib/server/agent.ts`'s list, and the types in `flaiover/src/lib/server/repo.ts`. Every test fakes `kluster` and `kubectl`, as those of S-0577 and S-0581 do: none needs a cluster or Docker.

## Done when

- Tests cover `kluster.status` with no block, with the scaffold's block, with a full block, with a block `flai check` faults, and with each tool missing from `PATH`.
- Tests cover `kluster.workloads` with a fake `kubectl`, with no kubeconfig, and with `kubectl` missing.
- `flaiover/src/routes/api/kluster/kluster.test.ts` covers the new route.
- `flai test` passes on the paths this task changes.

## Notes
