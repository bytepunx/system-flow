---
id: S-0615
type: story
nature: remediation
title: flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes
status: review
owner: alex
created: 2026-11-07T23:39:11Z
updated: 2026-11-09T10:45:56Z
transitions:
  - to: ready
    at: 2026-11-09T09:55:53Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T09:56:04Z
    by: agent-S-0615
  - to: review
    at: 2026-11-09T10:45:56Z
    by: agent-S-0615
tags: [flai, dashboard, host]
topics: [dashboard, host]
touches: [flai/internal/host/watch.go, flai/internal/host/watch_test.go, flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_porthold.go, flai/cmd/dashboard_porthold_test.go, design/system/flaiover-dashboard.md, design/adrs, docs/operators/index.md, docs/users/flai.md, design/issues/I-0243-flai-host-cannot-see-a-dashboard-container-docker-lost-track-of-which-keeps-the.md, design/issues/summary.md, flai/cmd/root.go, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3011
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 1
      hand_edits: 1
      work: 120
  models:
    - model: claude-opus-5-5
      input: 354
      output: 60816
      cache_read: 52495951
      cache_write: 723492
      cost: 16.3606
  strategic:
    - kind: planner
      seconds: 511
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 34
          output: 15989
          cache_read: 1285719
          cache_write: 116054
          cost: 0.0875
        - model: claude-opus-5-5
          input: 68
          output: 33407
          cache_read: 5053985
          cache_write: 181995
          cost: 3.1352
    - kind: orchestrator
      seconds: 225
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 58
          output: 404
          cache_read: 2162623
          cache_write: 282966
          cost: 0.6389
        - model: claude-sonnet-5-5
          input: 6
          output: 48
          cache_read: 42239
          cache_write: 34728
          cost: 0.0569
cost_of_delay:
  inputs:
    time_lost_per_cycle: 144h
    by: flai
    at: 2026-11-07T23:39:11Z
  value: 21600
  by: planner-S-0615
  at: 2026-11-09T09:53:42Z
forecast:
  duration: 2h
  delivery: 2026-11-09T12:23:00Z
  basis: "Its own forecast of 2h; 1st in the pull order with an in-progress limit of 5, behind S-0608, S-0610 and S-0614."
  by: flai
  at: 2026-11-09T09:55:56Z
finalized:
  by: orchestrator
  at: 2026-11-09T09:55:28Z
---
# S-0615 flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes

## Goal

This story remediates [I-0243](../../../design/issues/I-0243-flai-host-cannot-see-a-dashboard-container-docker-lost-track-of-which-keeps-the.md), "flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0243 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0243 is closed with `flai issue close I-0243 --reason` saying what fixed it

## Tasks
- T-2714 An ADR and the dashboard's design say flai host does not restart a dashboard whose published port another process holds
- T-2715 flai finds the docker-proxy that holds a host port from the process table
- T-2716 flai host's watch finds the dashboard's port held by another process and stops restarting it, with a test that reproduces I-0243
- T-2717 flai dashboard status and flai host status say the dashboard's port is held, by what, and how to free it
- T-2718 The operators guide and the flai guide say what the watch does when another process holds the dashboard's port
- T-2719 Close I-0243 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0243. time_lost_per_cycle 144h: 144h per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-31T18:32:12Z, 7.2 days before this story).

### Planning

Proposed solution, from I-0243's one instance and the code on main:

- The cause in the code: `dashboardState` in `flai/cmd/dashboard_watch.go` answers not answering whenever `probeAddress` (`docker port flaiover 3000/tcp`, in `flai/cmd/dashboard_upgrade.go`) gives no address, without asking whether the recorded published port answers. A container with an empty port map is always not answering, so the watch in `flai/internal/host/watch.go` restarts it after two looks, at the back-off's five-minute cap, for ever.
- A new watch state, `port-held`: docker runs the recorded container, `docker port` publishes nothing for it, and the host port it was recorded with (`dashboardRecord.Publish`) is held all the same, by something that is not the container. The watch does not restart on it, since a restart cannot win the port back, and logs it once each time the state is entered, at error, naming the holder.
- The holder: on Linux, the `docker-proxy` process whose command line (`/proc/<pid>/cmdline`, readable by any user) carries `-host-port <port>`, with its PID and its `-container-ip`. Where none is found, or there is no `/proc`, it says only that a process other than the container holds the port.
- `flai dashboard status` and `flai host status` say the port is held, by what, and what to do: stop that process (`sudo kill <pid>` of the docker-proxy and the orphan's shim, or restart the Docker daemon), then the watch starts the dashboard on its own.
- A new ADR refines ADR-0062: a dashboard whose published port another process holds is not restarted.

Tasks, in four layers:

| Layer | Tasks | Why it waits |
|-------|-------|--------------|
| 1 | T-2714 (ADR and design), T-2715 (holder lookup and port probe) | nothing; no file in common |
| 2 | T-2716 (the watch state, with the test that reproduces I-0243) | T-2715, whose lookup and probe it calls |
| 3 | T-2717 (the status lines), T-2719 (close I-0243) | T-2716's state and detail; T-2719 also T-2714's ADR, which its reason names |
| 4 | T-2718 (operators guide and flai guide) | T-2714's ADR, which it links, and T-2717's messages, which it quotes |

Touches, each a file but one:

| Touch | From |
|-------|------|
| `flai/cmd/dashboard_watch.go` | design: I-0243's cause, `dashboardState` and `lookDashboard` |
| `flai/internal/host/watch.go` | design: the watch's states and `watcher.look`, which restarts on every state but running and off |
| `flai/cmd/dashboard_watch_test.go` | co-change: 6 of 8 commits with the two above (75%) |
| `flai/cmd/dashboard_test.go` | co-change: 5 of 8 (62%); its `fakeRunner` answers `docker port` with an address always, so it must learn to answer none to reproduce the cause |
| `flai/cmd/dashboard.go` | co-change: 7 of 8 (88%); `flai dashboard status` calls `dashboardState` and says the new state |
| `flai/internal/host/watch_test.go` | co-change: 2 of 8 (25%); the watcher's tests |
| `flai/cmd/host.go`, `flai/cmd/host_test.go` | layout: `flai host status`'s `dashboard:` line and its tests |
| `flai/cmd/dashboard_porthold.go`, `flai/cmd/dashboard_porthold_test.go` | layout: new, the lookup of the port's holder, testable against a fake `/proc` tree |
| `design/system/flaiover-dashboard.md` | design: § Runtime says what the watch restarts |
| `design/adrs` | the new ADR refining ADR-0062 |
| `docs/operators/index.md` | design: § The dashboard's watch |
| `docs/users/flai.md` | co-change: 2 of 8 (25%); the `flai dashboard status` line |
| `design/issues/I-0243-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

The folder touch `design/adrs` stays: the ADR's number and name are given by `flai adr new` when T-2714 writes it, so no task can name its file now. The folder is in `claims.shared`, so it holds no ready story. Left out: `flai/cmd/dashboard_upgrade.go` (co-change 75%), since `probeAddress` keeps its contract and the new branch is in `dashboardState`; `docs/operators/settings.md` (38%), since no setting is added. The story has no `after`: it reads nothing another open story writes.

Figures:

- Forecast: `flai forecast` gave 30m (104 s per unit of size over 58 done large remediation stories, size 17). Raised to 2h: six tasks in four layers, two new Go files tested against a fake `/proc`, a state carried through two packages, an ADR, and a close-out that took 15 to 22 minutes on this host today (S-0612, S-0613). Delivery 2026-11-09T18:00:00Z, past flai's 10:32, since the story is a draft still to be finalized and promoted, behind six stories in the pull order.
- Cost of delay: `flai cod` gave 21600 USD a week from the input flai set from I-0243, 144h lost per 168h cycle at 150 USD an hour. The value stands as the input gives it. The 144h is the time the defect lasted, six days of restarts, rather than work lost, so the value is likely high; the plan's thread asks the operator whether to lower the input. The input is the operator's, and it is left as it is.
