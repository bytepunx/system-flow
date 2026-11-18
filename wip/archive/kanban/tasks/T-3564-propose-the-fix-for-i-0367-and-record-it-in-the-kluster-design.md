---
id: T-3564
type: task
nature: remediation
title: Propose the fix for I-0367 and record it in the kluster design
status: done
parent: S-0840
owner: alex
created: 2026-11-17T12:29:10Z
updated: 2026-11-18T07:46:37Z
transitions:
  - to: ready
    at: 2026-11-17T13:23:26Z
    by: agent-S-0840
  - to: in-progress
    at: 2026-11-17T13:23:27Z
    by: agent-S-0840
  - to: done
    at: 2026-11-18T07:46:37Z
    by: agent-S-0840
stream: S-0840
tags: [kluster, design]
touches: [design/system/kluster.md, design/tech/kluster.md]
usage:
  source: log
  seconds: 104
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 6739
      cache_read: 3720267
      cache_write: 126349
      cost: 1.8897
---
# T-3564 Propose the fix for I-0367 and record it in the kluster design

## Work

The story's goal asks for a fix proposed from I-0367's instance before any code is written, so this task comes first and waits for no other.

Open a thread on S-0840 with the proposal, the recommended answer first. `Cluster.Up` reads `kluster status` (`Cluster.Status`) before it runs `kluster up`:

- When the status lists the cluster and it runs, `Up` skips `kluster up` and writes the kubeconfig again with `kluster kubeconfig --name <name> --output .flai-cache/kluster/kubeconfig`. `flai kluster up` says that the cluster was up already and that its kubeconfig was written again.
- When the status lists the cluster but it is stopped, `Up` refuses. The refusal names the cluster as stopped and says what to do. kluster has no start command, so the refusal names `flai kluster down` and `flai kluster up`, or the provider's own start command.
- When the status does not list the cluster, `Up` runs `kluster up` and then writes the kubeconfig, as it does today.

Name the alternative the issue suggests, a command of its own, `flai kluster kubeconfig`, and why it is not recommended. The operator's documented step for a stale kubeconfig is `flai kluster up` already (`deploy.go`, `agentconfig.go`, authstar's I-0023). Name as out of scope the issue's second suggestion, flai serve writing the kubeconfig again before a job.

Once the thread is answered, record the decision in `design/system/kluster.md`, in § The project's cluster and in the `flai kluster up` row of § Commands. Record it in the Calls row of `design/tech/kluster.md` too. Bump each file's `updated`.

## Done when

- The operator answered the proposal on a thread on S-0840.
- `design/system/kluster.md` § Commands says what `flai kluster up` does for a cluster that is running, one that is stopped, and one that is absent.
- `design/tech/kluster.md`'s Calls row says that `up` reads `status` first.

## Notes
