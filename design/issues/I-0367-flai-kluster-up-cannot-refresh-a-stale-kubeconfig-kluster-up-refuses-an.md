---
id: I-0367
title: "flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it"
class: blocker
status: closed
count: 1
cost: 30m
first_reported: 2026-11-13T20:28:29Z
last_reported: 2026-11-13T20:28:29Z
updated: 2026-11-18T08:01:38Z
---

# I-0367 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it

## Description
flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it

## Instances

### 2026-11-13T20:28:29Z
Found in bytepunx/authstar on 2026-11-13 with flai 1.95.0 (9b63bacf3). The k3d containers of the cluster authstar-integration restarted, and its API port moved from 127.0.0.1:41600 to 127.0.0.1:40820. flai's kept .flai-cache/kluster/kubeconfig still named 41600, so every kluster job failed with connection refused: S-0025's teardown and deploy, main's promote, and S-0006's agent-requested deploy (job S-0006-deploy-20261113T200933Z, failed at the deploy step). The deploy script's message and authstar's I-0023 both say to write the kubeconfig again with flai kluster up. The operator ran it, and it failed each time: kluster up answered 'create cluster: cluster "authstar-integration" already exists' and exited 1, so flai never reached kluster kubeconfig --name <name> --output, which design/system/kluster.md says writes the file. The operator's workaround, by hand: kluster kubeconfig --config integration/kluster.yaml --name authstar-integration --output .flai-cache/kluster/kubeconfig. Suggested fix: flai kluster up treats an existing cluster as up and writes the kubeconfig again, or a command of its own does (flai kluster kubeconfig); and flai serve could write it again itself before a job when the kept server address does not answer and kluster status says the cluster runs, since no agent may run up.

## Remediation

Story S-0840 remediates this issue, created from it at 2026-11-17T08:32:51Z.
Closed 2026-11-18T08:01:38Z: S-0840 fixed it, as the operator decided on TH-1305. Cluster.Up reads kluster status first. For a cluster that runs, it runs no kluster up (which refused a cluster that exists) and writes the kubeconfig again with kluster kubeconfig. For a stopped cluster it refuses and names the fixes. For an absent cluster it brings it up as before. flai kluster up, the dashboard's Up job, and flai serve's bring-up all get the fix. TestUpOfARunningClusterWritesItsKubeconfigAgain in flai/internal/kluster reproduces the instance.
