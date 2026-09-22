---
id: I-0035
title: A host action that stops its own answering container is killed mid-run by that connection's own cancellation
class: defect
status: closed
count: 1
cost: 40m
first_reported: 2026-09-22T22:24:30Z
last_reported: 2026-09-22T22:24:30Z
updated: 2026-09-22T22:25:22Z
---

# I-0035 A host action that stops its own answering container is killed mid-run by that connection's own cancellation

## Description
A host action that stops its own answering container is killed mid-run by that connection's own cancellation

## Instances

### 2026-09-22T22:24:30Z
S-0081, T-0325: restart, a swapping upgrade, and a last-project stop each stop the container answering the very request that asked for them. exec.CommandContext ties the subprocess to the request context, a child of the WebSocket connection's; that connection dying (as a direct effect of the command's own docker stop) cancels the context and SIGINTs the flai subprocess mid-operation. Only found against real Docker (Playwright against the actual /host page): a restart left no container running at all, stopped between docker stop and docker run. Fixed with internal/hostapi's new spec.detachTimeout, which runs these three writes with a background context instead of the request's.

## Remediation
Closed 2026-09-22T22:25:22Z: S-0081, commit 735702d: dashboard.restart, dashboard.upgrade, and dashboard.stop run with their own background context (spec.detachTimeout, internal/hostapi), verified fixed against real Docker
