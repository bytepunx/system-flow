---
id: I-0025
title: "A pattern kill meant for a scratch server stops the operator's dashboard: container processes match pkill on the host"
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-09-20T06:30:30Z
last_reported: 2026-09-20T06:30:30Z
updated: 2026-10-01T09:08:42Z
---

# I-0025 A pattern kill meant for a scratch server stops the operator's dashboard: container processes match pkill on the host

## Description
A pattern kill meant for a scratch server stops the operator's dashboard: container processes match pkill on the host

## Instances

### 2026-09-20T06:30:30Z
S-0070 browser check: pkill -f 'node build' for a scratch flaiover on port 4399 also signalled the node build in flaiover-system-flow (Docker Engine in WSL2, same user, same kernel). The container stayed Up but stopped listening: adapter-node's graceful shutdown, held open by SSE clients, so docker ps and restart policies saw nothing wrong. Restarted with flai dashboard stop and flai dashboard. Agent rule since: stop scratch processes by PID or run them in named containers. Worth considering: a health check on the dashboard container so a server that stopped listening is restarted, and flai dashboard status reporting 'not answering'.

## Remediation
Closed 2026-10-01T09:08:42Z: S-0184: the image has a HEALTHCHECK, flai dashboard status says running, not answering, or gone, and flai host restarts a dashboard that is gone or not answering from the image it recorded, with a back-off (ADR-0062). S-0170 had already made the server exit on SIGTERM and agents stop by PID (ADR-0058).
