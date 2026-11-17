---
id: TH-1305
title: "I-0367 fix: flai kluster up writes a running cluster's kubeconfig again instead of failing on kluster up"
anchor:
  path: wip/kanban/stories/S-0840-flai-kluster-up-cannot-refresh-a-stale-kubeconfig-kluster-up-refuses-an.md
  item: S-0840
status: open
participants: [agent-S-0840]
created: 2026-11-17T13:23:43Z
updated: 2026-11-17T13:23:43Z
---

# TH-1305 I-0367 fix: flai kluster up writes a running cluster's kubeconfig again instead of failing on kluster up

On wip/kanban/stories/S-0840-flai-kluster-up-cannot-refresh-a-stale-kubeconfig-kluster-up-refuses-an.md.

## Entries

### 2026-11-17T13:23:43Z agent-S-0840
**Recommended: option 1.** `Cluster.Up` (`flai/internal/kluster/kluster.go`) reads `kluster status` before it runs `kluster up`. What it does next depends on what the status lists:

| Status lists the cluster | What `Up` does |
|--------------------------|----------------|
| Running | Skips `kluster up`. Writes one line to its output saying the cluster was up already. Then writes the kubeconfig again with `kluster kubeconfig --name <name> --output .flai-cache/kluster/kubeconfig` |
| Stopped | Runs nothing and refuses. The refusal says the cluster is stopped and names two fixes: start its k3d node containers with `docker start`, then run `flai kluster up` again; or run `flai kluster down` and then `flai kluster up` to make it anew |
| Absent | Runs `kluster up`, then writes the kubeconfig, as it does today |

`Up` keeps its signature. So the dashboard's **Up** job (`flai/internal/hostapi/kluster.go`) and flai serve's bring-up before a deploy (`flai/internal/serve/kluster.go`) get the fix with no change. The line about the cluster being up goes into their logs too.

**Why:** I-0367's instance failed on `create cluster: cluster "authstar-integration" already exists`. The refusals in `deploy.go` and `agentconfig.go`, and authstar's I-0023, already tell the operator to run `flai kluster up` for a stale kubeconfig. With this fix that advice works, and no message changes.

**Option 2:** a new command, `flai kluster kubeconfig`. Not recommended. It adds a command, a guard rule, and a reference entry, and it leaves every existing message pointing at an `up` that still fails.

**Out of scope:** flai serve writing the kubeconfig again before a job when the kept server does not answer. That is the issue's second suggestion. Your pending decision on TH-1294 covers it.

Answer 1 or 2, or name a change. I record the answer in `design/system/kluster.md` and `design/tech/kluster.md`, then build it.
