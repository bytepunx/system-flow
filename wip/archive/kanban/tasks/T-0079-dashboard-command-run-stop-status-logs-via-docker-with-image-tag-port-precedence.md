---
id: T-0079
type: task
nature: feature
title: "dashboard command: run, stop, status, logs via docker with image, tag, port precedence"
status: done
parent: S-0009
owner: alex
created: 2026-09-17T03:50:03Z
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
tags: [cli, dashboard]
---

# T-0079 dashboard command: run, stop, status, logs via docker with image, tag, port precedence

## Work
cmd/dashboard.go: flai dashboard [run] [--image] [--tag] [--port] [--pull] [--attach] [--open]; precedence flags, then system-flow.yaml dashboard, then config; requires docker on PATH with an install hint; pulls when the image is missing or --pull; docker run --rm -d --name flaiover-<project> -p 127.0.0.1:<port>:3000 -v <root>:/project -e PROJECT_DIR=/project --user uid:gid; prints the URL; stop, status, logs subcommands.

## Done when
Command tree builds; every docker invocation goes through execx.Runner.

## Notes
