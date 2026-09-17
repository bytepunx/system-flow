---
id: T-0080
type: task
nature: feature
title: "Tests with a recording runner; manual smoke with a stand-in image"
status: done
parent: S-0009
owner: alex
created: 2026-09-17T03:50:04Z
updated: 2026-09-17T03:54:27Z
transitions:
  - to: ready
    at: 2026-09-17T03:54:27Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:54:27Z
    by: agent
  - to: done
    at: 2026-09-17T03:54:27Z
    by: agent
stream: S-0009
tags: [cli]
---

# T-0080 Tests with a recording runner; manual smoke with a stand-in image

## Work
cmd/dashboard_test.go with a recording fake Runner: missing docker error, pull when image absent, run arguments (bind address, mount, user, env, name), stop, status, and precedence of flags over manifest over config; then a manual smoke with nginx:alpine as a stand-in image against this repository.

## Done when
Tests green; the stand-in container serves on the configured port and stops cleanly.

## Notes
