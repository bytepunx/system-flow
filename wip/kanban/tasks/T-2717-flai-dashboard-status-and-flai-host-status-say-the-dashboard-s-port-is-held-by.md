---
id: T-2717
type: task
nature: remediation
title: flai dashboard status and flai host status say the dashboard's port is held, by what, and how to free it
status: in-progress
parent: S-0615
owner: alex
created: 2026-11-09T09:52:14Z
updated: 2026-11-09T10:16:16Z
transitions:
  - to: ready
    at: 2026-11-09T10:16:15Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T10:16:16Z
    by: system-flow
stream: S-0615
tags: [flai, dashboard, host]
touches: [flai/cmd/dashboard.go, flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/dashboard_watch_test.go]
after: [T-2716]
usage:
  source: log
  seconds: 3
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 1
      output: 12
      cache_read: 150044
      cache_write: 528
      cost: 0.0548
---
# T-2717 flai dashboard status and flai host status say the dashboard's port is held, by what, and how to free it

## Work

Make the two status commands say the new state of T-2716 in words an operator can act on (I-0243), in place of "not answering" and a restart count that grows every five minutes.

- `flai dashboard status` (`flai/cmd/dashboard.go`): for `port-held`, a line such as `flaiover's port 4242 is held by docker-proxy (pid 188220) for 172.17.0.2:3000, not by its container: docker no longer lists that container. flai host does not restart it. Stop that process, or restart the Docker daemon, and the watch starts it again`; without a holder found, say a process other than the container holds the port and name `sudo ss -ltnp 'sport = :4242'` to find it. `--json` gives `state: port-held` and the holder (its PID, host address, and container address) under its own key.
- `flai host status` (`flai/cmd/host.go`): the `dashboard:` line says `port held by another process` with the watch's detail, and `--json` carries the detail under `dashboard`, as T-2716 records it.

It waits for T-2716, which adds the state, its detail, and the changed `dashboardState`.

## Done when

- A test in `flai/cmd/dashboard_watch_test.go` gives `flai dashboard status` a running container with no published address, a port that answers, and a fake `/proc` holder, and finds the line with the holder's PID and what to do, and `--json` with `port-held` and the holder; and with no holder found, the line naming `ss`.
- A test in `flai/cmd/host_test.go` finds the `dashboard:` line of `flai host status` saying the port is held, with the detail, and the detail in `--json`.
- `flai test` on the four files passes.

## Notes
