---
id: D-0668
title: "S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again: Recommended: draft neither proposed story now"
kind: answer
sources: [TH-1294]
topics: [kluster, flai, cli, go, code, all]
paths: [flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, design/system/kluster.md, design/tech/kluster.md]
date: 2026-11-18
status: active
absorbed_by: ""
asked_again: []
cited: [TH-1294]
---

# D-0668 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again: Recommended: draft neither proposed story now

The operator, alex, answered TH-1294 at 2026-11-18T07:11:03Z.

What was asked:

```text
The plan is three tasks in three layers, each waiting on the one before it.
```

The answer:

```text
**Recommended:** draft neither proposed story now. Decide again if I-0367 records a second instance after S-0840 lands.

**Why:**

- I-0367 has one instance. Once S-0840 lands, `flai kluster up` on the host writes a running cluster's stale kubeconfig again. That covers the case that was found.
- Both additions only remove a step for the operator. One is flai serve's own refresh before a kluster job. The other is a refresh on the Klusters page. Neither is needed for S-0840's criteria.
- Each one is a change of scope, so it is your decision.

**Alternative:** draft one story covering both, under the epic that owns kluster jobs. The planner can then plan it as a draft.

Source: design/conventions/strategic-agents.md § As the orchestrator
```

Source: TH-1294.
