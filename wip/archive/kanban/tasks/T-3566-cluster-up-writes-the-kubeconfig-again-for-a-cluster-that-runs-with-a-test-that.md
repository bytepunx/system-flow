---
id: T-3566
type: task
nature: remediation
title: Cluster.Up writes the kubeconfig again for a cluster that runs, with a test that reproduces I-0367
status: done
parent: S-0840
owner: alex
created: 2026-11-17T12:29:20Z
updated: 2026-11-18T07:56:39Z
transitions:
  - to: ready
    at: 2026-11-18T07:46:59Z
    by: agent-S-0840
  - to: in-progress
    at: 2026-11-18T07:46:59Z
    by: agent-S-0840
  - to: done
    at: 2026-11-18T07:56:39Z
    by: agent-S-0840
stream: S-0840
tags: [flai, kluster]
touches: [flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go]
after: [T-3564]
usage:
  source: log
  seconds: 580
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 3674
      cache_read: 2028391
      cache_write: 68889
      cost: 1.0303
---
# T-3566 Cluster.Up writes the kubeconfig again for a cluster that runs, with a test that reproduces I-0367

## Work

This task waits for T-3564: it builds the behaviour the operator agreed there, and changes it to fit the answer.

First write the test that reproduces I-0367 in `flai/internal/kluster/kluster_test.go`. Give the fake kluster script a mode in which `up` fails with `create cluster: cluster "<name>" already exists` while `status` lists the cluster as running. Today `Up` fails there and writes no kubeconfig.

Then change `Cluster.Up` in `flai/internal/kluster/kluster.go`. It reads `c.Status(ctx)` before `kluster up`:

- A cluster that exists and runs: skip `kluster up`, make the kubeconfig's folder, and run `kluster kubeconfig` as today. Say on `out` that the cluster was up already and that its kubeconfig is written again. Return a result or a flag that the command can report.
- A cluster that exists but is stopped: refuse with the error T-3564 agreed, writing no kubeconfig.
- A cluster that does not exist: run `kluster up` and then `kluster kubeconfig`, as today.

A failed `kluster status` fails `Up` with its own error. Keep `TestUp` and `TestUpThatFailsWritesNoKubeconfig` passing: they give a status with no row for the cluster. Add a case for the stopped cluster, and check the calls each case runs.

The flai serve job behind the dashboard's **Up** calls `Cluster.Up` through `hostapi.KlusterCluster`. Run `flai test` on `flai/internal/kluster` and `flai/internal/hostapi` to see that the job still passes.

## Done when

- A test in `flai/internal/kluster/kluster_test.go` runs `Up` on a cluster that exists and runs, while `kluster up` would fail with "already exists". It passes, and the kubeconfig is written under `.flai-cache/kluster`.
- A test covers a stopped cluster, which `Up` refuses with no kubeconfig written.
- `flai test flai/internal/kluster/kluster.go flai/internal/kluster/kluster_test.go` passes.

## Notes
