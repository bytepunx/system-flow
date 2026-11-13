---
id: T-3276
type: task
nature: feature
title: Each passing deploy makes the namespace's own kubeconfig from a flai-agent ServiceAccount bound to edit, and the teardown deletes it
status: in-progress
parent: S-0763
owner: alex
created: 2026-11-13T09:02:45Z
updated: 2026-11-13T09:08:16Z
transitions:
  - to: ready
    at: 2026-11-13T09:08:16Z
    by: agent-S-0763
  - to: in-progress
    at: 2026-11-13T09:08:16Z
    by: agent-S-0763
stream: S-0763
tags: [cli]
touches: [flai/internal/kluster/agentconfig.go, flai/internal/kluster/agentconfig_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 15
          cache_read: 81709
          cache_write: 2328
          cost: 0.0242
---
# T-3276 Each passing deploy makes the namespace's own kubeconfig from a flai-agent ServiceAccount bound to edit, and the teardown deletes it

## Work

Criteria 1 and 2, in `flai/internal/kluster` alone. `DeployBy` and `TeardownBy` are what `flai kluster`, flai serve's jobs, and the dashboard all call, so `flai/internal/serve` needs no change.

In a new `agentconfig.go`:

- `AgentKubeconfigPath(root, story)`: `.flai-cache/kluster/agents/s-nnnn.kubeconfig`.
- The step that makes it, run with the namespace's `Env` and flai's cluster kubeconfig:
  - `kubectl create serviceaccount flai-agent` in the namespace, and a RoleBinding `flai-agent` to the ClusterRole `edit`, both idempotent, as `kubectl apply` of a manifest or `create --dry-run=client -o yaml | apply` gives;
  - `kubectl create token flai-agent --duration 24h` in the namespace;
  - a kubeconfig with the cluster entry (server and certificate authority) of flai's kubeconfig for the context `k3d-<name>`, that token as its user, and the namespace as its context's namespace, written whole (`atomicfile`) with mode 0600.
- `Cluster.AgentEnv(story)`: `Env` for the story's namespace without `FLAI_KLUSTER_IMAGES`, and with `KUBECONFIG` replaced by `AgentKubeconfigPath`, not appended, since `Env` sets it to flai's kubeconfig already. T-3279's `exec` and S-0661's `env` read it.

In `deploy.go`, call the step in `runDeploy` once `c.steps` has passed, before the outcome is set to passed: for a story's namespace only, never for the default namespace a promotion deploys. Its output goes to the deploy step's log, and a failure fails the deploy at the `deploy` step with its error, so the record keeps the steps it has today and `kluster.log`, the Deployments page, and `plan()` need no new step name. A deploy again makes a new token and file. In the teardown's removal of the record, the state folder, and the logs, remove the namespace's kubeconfig too; a teardown that fails keeps it, as it keeps the record.

Tests, with `fakeKubectl` taught to answer `create token` with a token and to record each call: the calls after a passing deploy, the file's cluster, user, namespace, and mode; none after a failed deploy step; none for a promotion; a new file at a deploy again; the file gone after a teardown and kept after a failed one; `AgentEnv` without `FLAI_KLUSTER_IMAGES` and with one `KUBECONFIG`. The teardown test's expected call list (`deploy_test.go`, the teardown's) is unchanged, since the teardown runs no new kubectl call.

Waits for nothing: it is the first layer.

## Done when

- A passing deploy of a story writes `.flai-cache/kluster/agents/s-nnnn.kubeconfig` from a `flai-agent` ServiceAccount bound to `edit` with a 24-hour token, a deploy again writes a new one, and the teardown removes it.
- `Cluster.AgentEnv` gives the namespace's variables without `FLAI_KLUSTER_IMAGES` and with `KUBECONFIG` the namespace's kubeconfig.
- `go test ./internal/kluster/` passes.

## Notes
