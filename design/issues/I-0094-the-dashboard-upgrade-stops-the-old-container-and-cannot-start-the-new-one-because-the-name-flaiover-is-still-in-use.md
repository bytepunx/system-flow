---
id: I-0094
title: The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use
class: defect
status: open
count: 2
cost: 6m
first_reported: 2026-10-06T21:01:06Z
last_reported: 2026-10-07T15:09:19Z
updated: 2026-10-07T18:59:46Z
---

# I-0094 The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use

## Description

`dashboard.upgrade` replaces the running flaiover container with one on the new image. It stops the old container and at once runs the new one under the same name. Docker refused the new one: the name was still held by the old container. The likely reason, not confirmed, is that the old container runs with `--rm`, as the new one does, and Docker had not finished removing it when the new `docker run` asked for the name. The upgrade then reports failure with the dashboard down.

flai host's own watch recovered it: it found the dashboard gone and started it about a minute later. To the operator, in the dashboard, the page they had just used to upgrade stopped answering.

## Instances

### 2026-10-06T21:01:06Z
Story: S-0223.
On 2026-10-06 the operator ran the dashboard upgrade from the dashboard at 20:36:56Z, after upgrading the host to flai 1.33.0. dashboard.upgrade took 45 seconds and failed: the previous image stopped, and docker run for the new one exited 125, 'Conflict. The container name /flaiover is already in use'. The dashboard was down until flai host found it gone and started it at 20:38:06Z, on 0.35.2. The operator then restarted every process at 20:48:06Z and reported the project as hanging or crashing. The cost is the time from the failure to that restart.

### 2026-10-07T15:09:19Z
On 2026-10-07 at 15:07:21Z flai dashboard upgrade, run in a shell from the operator's home directory on flai 1.35.0, stopped the container on 0.36.1 and failed to start the new one: docker run exited 125 with 'Conflict. The container name /flaiover is already in use'. flai host's watch found the dashboard gone and started it at 15:07:45Z on the image just pulled, flaiover 0.41.0, so the pull took effect after 24 seconds without a dashboard. The journal shows the same failure from the dashboard's own upgrade action on 2026-10-06 at 20:37Z and on 2026-10-07 at 01:09Z, each recovered by the watch; this is the fourth occurrence.

## Remediation

Directions to weigh: after stopping the old container, wait until Docker no longer lists the name before running the new one, with a short limit; or run the new container under a temporary name and rename it; and when the start fails, start the previous image again so that the operator is not left without a dashboard. A test with a stand-in for docker that keeps the name for a moment after `stop` would reproduce it.

Story S-0316 remediates this issue, created from it at 2026-10-07T18:59:46Z.
