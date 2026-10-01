---
id: S-0184
type: story
nature: improvement
title: The dashboard recovers from a stray stop and is not taken over by a second flai serve
status: ready
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:32:11Z
transitions:
  - to: ready
    at: 2026-10-01T08:32:11Z
    by: alex
tags: [flai, dashboard]
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/internal/host, flai/internal/channel, flaiover/src/lib/server/agent.ts, flaiover/Dockerfile]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

## Notes

The S-0084 worktree holds the uncommitted T-0339 attempt; read it before starting, then let it go.
