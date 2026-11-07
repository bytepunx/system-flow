---
id: T-2569
type: task
nature: improvement
title: The Deployments page starts and stops each service's port-forward, which flai serve runs on the host bound to 0.0.0.0
status: in-progress
parent: S-0586
owner: alex
created: 2026-11-07T09:54:08Z
updated: 2026-11-07T09:54:48Z
transitions:
  - to: ready
    at: 2026-11-07T09:54:46Z
    by: agent-S-0586
  - to: in-progress
    at: 2026-11-07T09:54:48Z
    by: agent-S-0586
stream: S-0586
tags: [kluster]
touches: [flai/internal/hostapi/kluster_forward.go, flai/internal/hostapi/kluster_forward_test.go, flai/internal/hostapi/hostapi.go, flaiover/src/routes/api/kluster/forward/+server.ts, flaiover/src/routes/api/kluster/forward/forward.test.ts, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts]
---
# T-2569 The Deployments page starts and stops each service's port-forward, which flai serve runs on the host bound to 0.0.0.0

## Work

The operator's answer 3 on TH-0861 (D-0511): Start and Stop for each port-forward on the Deployments page, run by flai serve, bound to `0.0.0.0` so that a hosted flai can reach them. It waits for no task: T-2562's `kluster.status` and the Deployments page's access rows, which S-0581 built, are what it builds on, and it shares no path with T-2559.

- **Host methods**, in a new `flai/internal/hostapi/kluster_forward.go`, behind the `kluster` host action as `kluster.deploy` is: `kluster.forward.start` `{namespace, service, port}` runs `kubectl port-forward --address 0.0.0.0 --namespace <ns> service/<name> <local>:<port>` with flai's kubeconfig and the cluster's context, on a free local port flai picks, and answers the forward with its local port; `kluster.forward.stop` `{namespace, service, port}` stops it; each namespace's access row in `kluster.deployments` says which of its service ports have a forward running and on which local port. The forwards are kept by flai serve in memory, one per namespace, service, and port; a second start answers the running one; a forward whose process exits is dropped and its last output kept for the row; all stop when flai serve stops, and a namespace's when it is torn down. Values are checked: the namespace a DNS label, the service a name the cluster lists, the port one it lists.
- **flaiover**: a `forward` route that posts start and stop, the methods in `agent.ts`'s list, the types in `repo.ts`, and on the Deployments page a Start or Stop button per service port in place of the `kubectl port-forward` command, with a link to the forward on the host's address and the local port once it runs. The page says once that a running forward is reachable from the host's network, with no sign-in of its own.

Every test fakes `kubectl`: none needs a cluster.

## Done when

- Go tests cover start, a second start, stop, a forward whose process exits, a refused value, the action off, and a teardown stopping its namespace's forwards, with the `--address 0.0.0.0` argument checked.
- `deployments.svelte.test.ts` covers Start, the running link, Stop, and a refusal shown beside its row; no text names a command for the operator to run.
- `flai test` passes on the paths this task changes.

## Notes
