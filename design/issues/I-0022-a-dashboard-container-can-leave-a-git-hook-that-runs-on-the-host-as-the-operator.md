---
id: I-0022
title: A dashboard container can leave a git hook that runs on the host as the operator
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-09-19T07:56:04Z
last_reported: 2026-09-19T07:56:04Z
updated: 2026-09-19T10:27:32Z
---

# I-0022 A dashboard container can leave a git hook that runs on the host as the operator

## Description
A dashboard container can leave a git hook that runs on the host as the operator

## Instances

### 2026-09-19T07:56:04Z
S-0052: found while weighing credentials for the container. The clone is mounted read-write with .git included; in a scratch project a pre-push hook written from inside a container of the published image ran on the host as the host user at the next git push. Needs a compromised container, not just the dashboard token. Queued as S-0064.

## Remediation
Closed 2026-09-19T10:27:32Z: S-0064, ADR-0027: flai dashboard mounts .git/hooks, .git/config and .git/info read-only over the clone; a hook and settings written from inside the container fail, tried end to end. Takes effect for a dashboard once it is restarted by a flai that has the fix.
