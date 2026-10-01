---
id: S-0184
type: story
nature: improvement
title: The dashboard recovers from a stray stop and is not taken over by a second flai serve
status: in-progress
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:58:30Z
transitions:
  - to: ready
    at: 2026-10-01T08:32:11Z
    by: alex
  - to: in-progress
    at: 2026-10-01T08:54:33Z
    by: agent-S-0184
tags: [flai, dashboard]
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_watch.go, flai/cmd/host.go, flai/internal/host, flai/internal/channel, flai/internal/config, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/Dockerfile, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, docs/operators, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0025, design/issues/I-0029, design/issues/summary.md, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch_test.go, flai/cmd/dashboard_token.go, flai/cmd/host_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 255
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 130
      output: 856
      cache_read: 5669980
      cache_write: 216985
      cost: 2.3288
---
# S-0184 The dashboard recovers from a stray stop and is not taken over by a second flai serve

## Goal

Two issues leave the operator with a dashboard that is down or misdirected until they notice:

- **I-0025** (mitigated): a pattern kill meant for a scratch server stopped the dashboard container's process, because container processes match `pkill` on the host. The server now exits cleanly on SIGTERM and agents are stopped by PID (S-0170, ADR-0058), but the container runs `--detach --rm` with no restart policy or HEALTHCHECK (`flai/cmd/dashboard_upgrade.go` ~22), so a stray signal still takes the dashboard down until a manual `flai dashboard`, and `flai dashboard status` does not probe whether it answers.
- **I-0029** (mitigated): the single flai host (ADR-0040) refuses a second host, but a `flai serve` run by hand, or a host on another address, can still connect, and the dashboard evicts the older connection unconditionally (`flaiover/src/lib/server/agent.ts` ~224-228, "replaced by a newer connection"). T-0339's refuse-and-back-off work was abandoned with S-0084.

## Acceptance criteria
- [ ] The dashboard container has a HEALTHCHECK, and `flai dashboard status` reports running, not answering, or gone, from an HTTP probe
- [ ] The flai host restarts a dashboard that is gone or not answering, with a back-off and a log line saying why; the operator can turn this off
- [ ] The dashboard refuses a new flai connection for a project while the connection it has answers, with a close code that says so, and flai backs off on that code instead of retrying at once; a holder that stops answering is replaced
- [ ] Tests cover the probe, the restart, and the refused and the replaced connection
- [ ] The design (`design/system/flaiover-dashboard.md`, `design/system/dashboard-host-channel.md`) and the operator guide describe them
- [ ] I-0025 and I-0029 are closed with what fixed them

## Tasks
- T-0648 The dashboard image has a HEALTHCHECK and flai dashboard status says running, not answering, or gone
- T-0649 flai host restarts a dashboard that is gone or not answering, with a back-off, and dashboard.no_restart turns it off
- T-0650 The dashboard refuses a second flai connection while its holder answers, and flai backs off on that close code
- T-0651 Design and operator guide describe the watchdog and the refused connection; close I-0025 and I-0029

## Notes

The S-0084 worktree holds the uncommitted T-0339 attempt; read it before starting, then let it go.
