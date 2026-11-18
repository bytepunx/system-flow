---
id: D-0670
title: "I-0367 fix: flai kluster up writes a running cluster's kubeconfig again instead of failing on kluster up: Recommended: option 1, as agent-S-0840 wrote it"
kind: answer
sources: [TH-1305]
topics: [kluster, flai, cli, go, code, all]
paths: [flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, design/system/kluster.md, design/tech/kluster.md]
date: 2026-11-18
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0670 I-0367 fix: flai kluster up writes a running cluster's kubeconfig again instead of failing on kluster up: Recommended: option 1, as agent-S-0840 wrote it

The operator, alex, answered TH-1305 at 2026-11-18T07:45:54Z.

What was asked:

```text
**Recommended: option 1.** `Cluster.Up` (`flai/internal/kluster/kluster.go`) reads `kluster status` before it runs `kluster up`. What it does next depends on what the status lists:
```

The answer:

```text
**Recommended: option 1, as agent-S-0840 wrote it.** `Cluster.Up` reads `kluster status` first. For a running cluster it skips `kluster up` and writes the kubeconfig again. For a stopped one it refuses and names the two fixes. For an absent one it works as it does today.

**Why:**

- `design/system/kluster.md` § The project's cluster says flai brings the cluster up with `kluster up` and writes its kubeconfig with `kluster kubeconfig --name <name> --output`. Option 1 keeps both steps and fixes only their order for a cluster that exists.
- The refusals in `deploy.go` and `agentconfig.go` already tell the operator to run `flai kluster up` for a stale kubeconfig. Option 1 makes that advice work, with no message changed.
- `Up` keeps its signature, so the dashboard's **Up** job and flai serve's bring-up get the fix with no change. That stays within S-0840's touches.
- Option 2 adds a command, a guard rule, and a reference entry, and leaves `up` failing.

**Why I am asking you:** the story's goal asks for the fix to be proposed before it is built. Refusing a stopped cluster is new behaviour no document records yet.

**Alternative:** option 2, a new `flai kluster kubeconfig` command.

The second suggestion in I-0367, flai serve writing the kubeconfig again before a job, stays out of scope, as TH-1294 asks.

Source: design/system/kluster.md § The project's cluster
```

Source: TH-1305.
