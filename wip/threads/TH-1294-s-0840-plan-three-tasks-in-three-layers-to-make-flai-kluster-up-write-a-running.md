---
id: TH-1294
title: "S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again"
anchor:
  path: wip/kanban/stories/S-0840-flai-kluster-up-cannot-refresh-a-stale-kubeconfig-kluster-up-refuses-an.md
  item: S-0840
status: resolved
participants: [planner-S-0840, orchestrator, alex, agent-S-0840]
created: 2026-11-17T12:29:47Z
updated: 2026-11-18T07:12:09Z
---

# TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again

On wip/kanban/stories/S-0840-flai-kluster-up-cannot-refresh-a-stale-kubeconfig-kluster-up-refuses-an.md.

## Entries

### 2026-11-17T12:29:47Z planner-S-0840
The plan is three tasks in three layers, each waiting on the one before it.

| Layer | Task | Touches | Waits for |
|-------|------|---------|-----------|
| 1 | T-3564 Propose the fix for I-0367 and record it in the kluster design | `design/system/kluster.md`, `design/tech/kluster.md` | none |
| 2 | T-3566 Cluster.Up writes the kubeconfig again for a cluster that runs, with a test that reproduces I-0367 | `flai/internal/kluster/kluster.go`, `flai/internal/kluster/kluster_test.go` | T-3564: it builds what the operator agreed there |
| 3 | T-3570 flai kluster up says so, its reference is regenerated, and I-0367 is closed | `flai/cmd/kluster.go`, `flai/cmd/kluster_test.go`, `docs/users/flai-reference.md` | T-3566: the command reports what Up returns |

The cause: `Cluster.Up` (`flai/internal/kluster/kluster.go`) returns when `kluster up` fails, so `kluster kubeconfig` never runs for a cluster that exists. The fix proposed: `Up` reads `kluster status` first. A cluster that runs gets its kubeconfig written again, with no `kluster up`. A stopped cluster is refused, and the refusal names the cause. An absent cluster is brought up as today. T-3564 puts this to the operator on a thread before any code is written, as the goal asks.

Figures: touches are 7 files, with no folders. The forecast is 40m, raised from flai's 23m for three layers and the integration and smoke tiers of the close-out. The cost of delay is 75 USD a week, from flai's inputs, unchanged. There is no `after`.

Assumptions:

- The dashboard's **Up** job calls the same `Cluster.Up`, so it gets the fix with no change in `flai/internal/hostapi`.
- The refusals in `deploy.go` and `agentconfig.go` say to write the kubeconfig again with `flai kluster up`. They become true with no change, so they are not touched.
- kluster has no start command, so a stopped cluster is refused rather than started.

Proposed additions, not drafted, for you to decide:

- A story for the issue's second suggestion. flai serve would write the kubeconfig again before a kluster job when the kept server does not answer and `kluster status` says the cluster runs. This matters because no agent may run `kluster up`.
- A story so that the Klusters page can refresh the kubeconfig. Today it offers **Up** only while the cluster is not running, so a running cluster's stale kubeconfig cannot be refreshed from the dashboard.

### 2026-11-17T12:32:58Z orchestrator
Approved. The plan fits S-0840.

- T-3566 makes `Cluster.Up` write the kubeconfig again for a cluster that is running, with a test that reproduces I-0367. That meets criterion 1.
- T-3570 closes I-0367. That meets criterion 2.
- T-3564 records the proposed fix in `design/system/kluster.md` and `design/tech/kluster.md` before any code is written, as the goal asks.
- The touches are seven files and no folder. Each task's touches are inside them.
- Each layer waits on the output of the layer before it.
- The forecast of 40m and the value of 75 USD a week are flai's figures, with the planner's stated reason for the forecast.

The two stories you propose would add scope beyond I-0367: flai serve writing the kubeconfig again before a kluster job, and a refresh on the Klusters page. I do not draft them. I post a recommendation on them below for the operator.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T12:33:21Z orchestrator (recommendation)
**Recommended:** draft neither proposed story now. Decide again if I-0367 records a second instance after S-0840 lands.

**Why:**

- I-0367 has one instance. Once S-0840 lands, `flai kluster up` on the host writes a running cluster's stale kubeconfig again. That covers the case that was found.
- Both additions only remove a step for the operator. One is flai serve's own refresh before a kluster job. The other is a refresh on the Klusters page. Neither is needed for S-0840's criteria.
- Each one is a change of scope, so it is your decision.

**Alternative:** draft one story covering both, under the epic that owns kluster jobs. The planner can then plan it as a draft.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-18T07:11:03Z alex
Confirmed the recommendation of 2026-11-17T12:33:21Z orchestrator.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-18T07:12:00Z agent-S-0840
Resolved: The operator confirmed the plan and the recommendation: no extra story is drafted now. Decide again if I-0367 records a second instance after S-0840 lands.

### 2026-11-18T07:12:09Z orchestrator
Resolved: Plan approved on 2026-11-17. The operator confirmed the recommendation (D-0668): draft neither proposed story now, and decide again if I-0367 records a second instance after S-0840 lands. Nothing is left open.
