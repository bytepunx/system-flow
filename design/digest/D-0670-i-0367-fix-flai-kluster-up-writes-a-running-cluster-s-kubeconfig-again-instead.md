---
id: D-0670
title: "I-0367 fix: flai kluster up writes a running cluster's kubeconfig again instead of failing on kluster up: Recommended: option 1, as agent-S-0840 wrote it"
kind: answer
sources: [TH-1305, I-0367, S-0840]
topics: [kluster, flai, cli, go, code, all]
paths: [flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, design/system/kluster.md, design/tech/kluster.md]
date: 2026-11-18
status: absorbed
absorbed_by: "design/system/kluster.md#the-projects-cluster"
asked_again: []
cited: []
---

# D-0670 I-0367 fix: flai kluster up writes a running cluster's kubeconfig again instead of failing on kluster up: Recommended: option 1, as agent-S-0840 wrote it

The operator, alex, answered TH-1305 at 2026-11-18T07:45:54Z: option 1, as agent-S-0840 proposed it. S-0840 built it and closed I-0367.

What was asked: how to fix I-0367, where `flai kluster up` could not refresh a stale kubeconfig, since `kluster up` refuses a cluster that exists before flai writes the kubeconfig. Option 1 makes `Cluster.Up` read `kluster status` first; option 2 adds a new `flai kluster kubeconfig` command.

The answer:

```text
**Recommended: option 1, as agent-S-0840 wrote it.** `Cluster.Up` reads `kluster status` first. For a running cluster it skips `kluster up` and writes the kubeconfig again. For a stopped one it refuses and names the two fixes. For an absent one it works as it does today.
```

The reasons given: option 1 keeps both steps `design/system/kluster.md` § The project's cluster names and fixes only their order; the refusals in `deploy.go` and `agentconfig.go` already tell the operator to run `flai kluster up` for a stale kubeconfig, so that advice now works with no message changed; option 2 adds a command, a guard rule, and a reference entry and leaves `up` failing. flai serve writing the kubeconfig again before a job stays out of scope (TH-1294, D-0668).

How S-0840 built it (`wip/archive/agents/S-0840.md`, Decisions):

> Cluster.Up keeps its signature, Up(ctx, out) error. It writes its own line to out when the cluster was running already. So the hostapi Up job and flai serve's bring-up get the fix without a change.

I-0367's fix: `flai kluster up`, the dashboard's **Up** job, and flai serve's bring-up all run this one `Cluster.Up`. `TestUpOfARunningClusterWritesItsKubeconfigAgain` in `flai/internal/kluster` reproduces the instance.

Source: TH-1305, I-0367, S-0840.
