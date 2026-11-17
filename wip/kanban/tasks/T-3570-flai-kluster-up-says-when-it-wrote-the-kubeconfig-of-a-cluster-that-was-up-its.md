---
id: T-3570
type: task
nature: remediation
title: flai kluster up says when it wrote the kubeconfig of a cluster that was up, its reference is regenerated, and I-0367 is closed
status: backlog
parent: S-0840
owner: alex
created: 2026-11-17T12:29:29Z
updated: 2026-11-17T12:29:29Z
transitions: []
stream: S-0840
tags: [flai, kluster, docs]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md]
after: [T-3566]
---
# T-3570 flai kluster up says when it wrote the kubeconfig of a cluster that was up, its reference is regenerated, and I-0367 is closed

## Work

This task waits for T-3566: the command reports what `Cluster.Up` now returns, and its test runs against that behaviour.

In `flai/cmd/kluster.go`, change `newKlusterUpCmd`:

- Its `Long` help says that a cluster that runs already is not brought up again, and that only its kubeconfig is written again. Say this is how to refresh a kubeconfig after the cluster's API port moves. It also says that a stopped cluster is refused.
- Its last line says, for a cluster that was up already, that the cluster was up and its kubeconfig was written again, with the path and the context.

Add a case to `flai/cmd/kluster_test.go` beside `TestKlusterUp`, using the fake kluster from `runKluster`. The case is a cluster that `status` lists as running, whose `up` would fail with "already exists". It exits 0 and writes the kubeconfig.

Regenerate `docs/users/flai-reference.md` with `make flai-reference`, so that `TestReferenceIsCurrent` passes.

Then close the issue: `flai issue close I-0367 --reason "<what fixed it>"`. The reason names S-0840 and the change to `Cluster.Up`, and says that `flai kluster up` now writes the kubeconfig again for a cluster that runs. Then tick both of the story's criteria with `flai criteria tick`.

## Done when

- `flai kluster up` on a running cluster whose `kluster up` would fail with "already exists" exits 0, writes `.flai-cache/kluster/kubeconfig`, and says so. A test in `flai/cmd/kluster_test.go` shows it.
- `docs/users/flai-reference.md` matches the command's help, and `flai test docs/users/flai-reference.md flai/cmd/kluster.go` passes.
- I-0367 is closed with a reason naming what fixed it.

## Notes
