---
id: T-2716
type: task
nature: remediation
title: flai host's watch finds the dashboard's port held by another process and stops restarting it, with a test that reproduces I-0243
status: done
parent: S-0615
owner: alex
created: 2026-11-09T09:51:46Z
updated: 2026-11-09T10:16:00Z
transitions:
  - to: ready
    at: 2026-11-09T10:06:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T10:06:10Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:16:00Z
    by: agent-S-0615
stream: S-0615
tags: [flai, dashboard, host]
touches: [flai/internal/host/watch.go, flai/internal/host/watch_test.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard.go, flai/cmd/root.go]
after: [T-2715]
usage:
  source: log
  seconds: 590
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 6391
      cache_read: 5516755
      cache_write: 76031
      cost: 1.7193
---
# T-2716 flai host's watch finds the dashboard's port held by another process and stops restarting it, with a test that reproduces I-0243

## Work

The cause of I-0243 in the code: `dashboardState` in `flai/cmd/dashboard_watch.go` answers `not-answering` whenever `probeAddress` gives no address, so a container whose port map is empty because something else holds the port is restarted at every back-off, for ever.

- `flai/internal/host/watch.go`: add the state `DashboardPortHeld = "port-held"` (its container runs, docker publishes nothing for it, and its recorded host port is held by another process). `watcher.look` treats it as no miss and no restart: misses go to zero and the back-off stands. Carry the holder with the state, for instance by `Look` returning a state and a detail and `DashboardWatch` gaining `detail` (JSON `detail,omitempty`), and have `watchDashboard` log `dashboard port held by another process` at error once each time the state is entered, with the detail, not at every look.
- `flai/cmd/dashboard_watch.go`: when docker runs the container and `probeAddress` fails or gives nothing, ask the port probe of T-2715 whether the recorded `publish` answers; when it does, the state is `port-held` and the detail names the holder from the lookup of T-2715 on `/proc`, or says a process other than the container holds the port when none is found. Otherwise it stays `not-answering`. `dashboardState` needs the recorded publish: pass it from `lookDashboard`, and from `flai dashboard status` in `flai/cmd/dashboard.go` the settings' bind and port, changing only that call there; what the status prints is T-2717's.
- `flai/cmd/dashboard_test.go`: let `fakeRunner` answer `docker port` with nothing, as docker does for a container with an empty port map, so a test can reproduce the cause.

It waits for T-2715, whose lookup and port probe it calls.

## Done when

- A test in `flai/cmd/dashboard_watch_test.go` reproduces I-0243: docker runs `flaiover`, `docker port` answers nothing, the recorded port answers, and a fake `/proc` names a `docker-proxy` for it; `lookDashboard` answers `port-held` with that holder, and the host's watch over a run of looks makes no restart. With the port not answering it is `not-answering` and restarted as before.
- `flai/internal/host/watch_test.go`: `port-held` never restarts, resets the misses, keeps the back-off, and is logged once per entry; the existing restart and back-off tests pass unchanged.
- `flai test` on the six files passes.

## Notes
