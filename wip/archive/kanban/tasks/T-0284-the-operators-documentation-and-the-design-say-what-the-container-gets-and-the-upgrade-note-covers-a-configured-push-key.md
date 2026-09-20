---
id: T-0284
type: task
nature: feature
title: The operators' documentation and the design say what the container gets, and the upgrade note covers a configured push key
status: done
parent: S-0077
owner: alex
created: 2026-09-20T12:39:12Z
updated: 2026-09-20T12:48:14Z
transitions:
  - to: ready
    at: 2026-09-20T12:45:52Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:45:52Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:48:14Z
    by: system-flow
stream: S-0077
tags: []
---
# T-0284 The operators' documentation and the design say what the container gets, and the upgrade note covers a configured push key

## Work
docs/operators/index.md rewritten where it explains mounts, the push key, read-only git paths, and what the container can write. design/system/flaiover-dashboard.md, flai-cli.md, pushing-from-the-board.md, dashboard-host-channel.md, design/tech/docker.md. Upgrade note for an operator with dashboard.push_key set.

## Done when
- make lint-md and make check pass

## Notes
