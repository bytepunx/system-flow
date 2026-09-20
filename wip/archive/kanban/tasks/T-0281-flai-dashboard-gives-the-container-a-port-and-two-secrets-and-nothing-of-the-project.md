---
id: T-0281
type: task
nature: feature
title: flai dashboard gives the container a port and two secrets, and nothing of the project
status: done
parent: S-0077
owner: alex
created: 2026-09-20T12:39:12Z
updated: 2026-09-20T12:44:47Z
transitions:
  - to: ready
    at: 2026-09-20T12:39:13Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:39:13Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:44:47Z
    by: system-flow
stream: S-0077
tags: []
---
# T-0281 flai dashboard gives the container a port and two secrets, and nothing of the project

## Work
runDashboard stops passing the project volume, PROJECT_DIR, the guard mounts, the git identity, the excludes file, the push key with its known hosts and passwd, and drops the clone audit. dashboard_guard.go and dashboard_pushkey.go go. --push-key, --push-known-hosts, and config dashboard.push_key and push_known_hosts are retired: set, they print what replaced them and the dashboard starts anyway. status and the JSON output drop mount, read_only, and push_key. Whether the container still runs as the host user is decided and recorded.

## Done when
- A test reads the docker arguments: no --volume, only the two secret mounts, no GIT_ or PROJECT_DIR variables
- A test for the retired settings' message
- make flai-test passes

## Notes
