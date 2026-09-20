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

Reviewed 2026-09-20 (S-0077, [ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)): the routes ADR-0027 left open, a tracked file the operator runs, a moved ref, and an ignored file the host executes (a built binary, `node_modules/.bin`, a tool cache), are closed too, by removal: the container is given no file of the project, so there is nothing for it to write. Tried from inside a container started the new way; see the story's narrative. What a compromised container can still do is ask flai on the host for the channel's named methods and use the two secrets it must hold; that is restated in ADR-0031's consequences and is not this issue.
