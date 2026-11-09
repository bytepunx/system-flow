---
id: I-0243
title: flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes
class: defect
status: closed
count: 1
cost: 144h
first_reported: 2026-10-31T18:32:12Z
last_reported: 2026-10-31T18:32:12Z
updated: 2026-11-09T10:16:16Z
---

# I-0243 flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes

## Description
flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes

## Instances

### 2026-10-31T18:32:12Z
Since 2026-10-25T20:50Z flai host has restarted the dashboard every five minutes, 1,697 times by 2026-10-31, each with reason not-answering. The cause: a flaiover 0.55.2 container started 2026-10-25T20:23Z (id 1c706ffd) is no longer listed by docker ps, even with --all, but its processes still run: containerd-shim 188119, tini 188143, node server.js 188237, and docker-proxy 188220 holding 0.0.0.0:4242 to 172.17.0.2:3000. Every container flai host starts since then fails to publish 4242 (docker inspect shows an empty port map), gets the same bridge address 172.17.0.2 since Docker's IPAM does not know the orphan, and is shadowed by it, so the host's /_health probe of its own container fails and it restarts it, while browsers and flai serve reach the orphan. Found when a dev flai was installed with dashboard.allow_unsigned on: the orphan, started without FLAIOVER_ALLOW_UNSIGNED, kept refusing it as unsigned although every container the host started carried the allowance. flai host could tell the published port is held by a process that is not its container (docker inspect's port map empty while the port answers) and say so instead of restarting, and name the process to stop.

## Remediation

Story S-0615 remediates this issue, created from it at 2026-11-07T23:39:11Z.
Closed 2026-11-09T10:16:16Z: S-0615: flai host's watch now finds a dashboard whose container runs with no published port while its recorded host port answers, calls it port held, names the docker-proxy that holds the port from /proc, and does not restart it, logging it once (ADR-0267). TestLookDashboardFindsThePortHeldByAnOrphanedContainersProxy reproduces the cause, and TestHostDoesNotRestartADashboardWhosePortIsHeld shows the watch makes no restart until the holder is gone.
