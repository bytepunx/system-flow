---
id: T-3548
type: task
nature: remediation
title: A story's container is told the host name by which it reaches the cluster's API server on its network
status: done
parent: S-0834
owner: alex
created: 2026-11-17T08:32:37Z
updated: 2026-11-18T07:40:38Z
transitions:
  - to: ready
    at: 2026-11-18T07:11:52Z
    by: agent-S-0834
  - to: in-progress
    at: 2026-11-18T07:11:52Z
    by: agent-S-0834
  - to: done
    at: 2026-11-18T07:40:38Z
    by: agent-S-0834
stream: S-0834
tags: [harness, container, go]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, scripts/agent-container.sh, docs/operators/settings.md]
after: [T-3546]
usage:
  source: log
  seconds: 1726
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 53
      output: 314
      cache_read: 3279258
      cache_write: 107053
      cost: 1.9876
---
# T-3548 A story's container is told the host name by which it reaches the cluster's API server on its network

## Work

T-3546's ADR decides how the namespace's kubeconfig reaches the API server from a container, so this task waits for it. Use the variable's name and values the ADR gives. The planner proposed `FLAI_KLUSTER_API_HOST`.

- `containerEnv` in `flai/internal/harness/container.go` gives a story's container the variable:
  - `host.docker.internal` on `bridge`, the network the host names or the default when it names none, where `--add-host host.docker.internal:host-gateway` is already given;
  - empty on `host`.
- Pass it by name, with its value in the docker client's environment, as the other variables go.
- `scripts/agent-container.sh` gives the same variable on its own `bridge` and `host` networks.
- Add the variable's row to `docs/operators/settings.md`, whose test fails on an environment variable with no row.

T-3547 changes the kluster package at the same time; this task shares no file with it.

## Done when

- A test in `container_test.go` shows the variable in the `docker run` argv by name, with `host.docker.internal` on `bridge` and empty on `host`.
- `docs/operators/settings.md` has the variable's row.
- ShellCheck passes on `scripts/agent-container.sh`.
- `flai test` passes on the changed paths.

## Notes
